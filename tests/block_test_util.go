// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from tests/block_test_util.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package tests

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/gxhash"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

// A BlockTest checks handling of entire blocks.
type BlockTest struct {
	json btJSON
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (t *BlockTest) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

type btJSON struct {
	Blocks     []btBlock               `json:"blocks"`
	Genesis    btHeader                `json:"genesisBlockHeader"`
	Pre        blockchain.GenesisAlloc `json:"pre"`
	Post       blockchain.GenesisAlloc `json:"postState"`
	BestBlock  common.UnprefixedHash   `json:"lastblockhash"`
	Network    string                  `json:"network"`
	SealEngine string                  `json:"sealEngine"`
}

type btBlock struct {
	BlockHeader     *btHeader
	ExpectException string
	Rlp             string
	UncleHeaders    []*btHeader
}

//go:generate gencodec -type btHeader -field-override btHeaderMarshaling -out gen_btheader.go

type btHeader struct {
	Bloom                 types.Bloom
	Coinbase              common.Address
	MixHash               common.Hash
	Nonce                 uint64
	Number                *big.Int
	Hash                  common.Hash
	ParentHash            common.Hash
	ReceiptTrie           common.Hash
	StateRoot             common.Hash
	TransactionsTrie      common.Hash
	UncleHash             common.Hash
	ExtraData             []byte
	Difficulty            *big.Int
	GasLimit              uint64
	GasUsed               uint64
	Timestamp             uint64
	BaseFeePerGas         *big.Int
	WithdrawalsRoot       *common.Hash
	BlobGasUsed           *uint64
	ExcessBlobGas         *uint64
	ParentBeaconBlockRoot *common.Hash
}

type btHeaderMarshaling struct {
	Nonce         math.HexOrDecimal64
	ExtraData     hexutil.Bytes
	Number        *math.HexOrDecimal256
	Difficulty    *math.HexOrDecimal256
	GasLimit      math.HexOrDecimal64
	GasUsed       math.HexOrDecimal64
	Timestamp     math.HexOrDecimal64
	BaseFeePerGas *math.HexOrDecimal256
	BlobGasUsed   *math.HexOrDecimal64
	ExcessBlobGas *math.HexOrDecimal64
}

type eestEngine struct {
	*gxhash.Gxhash
	baseFee  *big.Int
	gasLimit uint64
}

var _ consensus.Engine = &eestEngine{}

func (e *eestEngine) BeforeApplyMessage(evm *vm.EVM, msg *types.Transaction) {
	// Change GasLimit to the one in the eth header
	evm.Context.GasLimit = e.gasLimit

	if evm.ChainConfig().Rules(evm.Context.BlockNumber).IsCancun {
		// EIP-1052 must be activated for backward compatibility on Kaia. But EIP-2929 is activated instead of it on Ethereum
		vm.ChangeGasCostForTest(&evm.Config.JumpTable, vm.EXTCODEHASH, params.WarmStorageReadCostEIP2929)
	}

	// When istanbul is enabled, instrinsic gas is different from eth, so enable IsPrague to make them equal
	r := evm.ChainConfig().Rules(evm.Context.BlockNumber)
	if evm.ChainConfig().Rules(evm.Context.BlockNumber).IsIstanbul {
		r.IsPrague = true
	}
	updatedIntrinsicGas, dataTokens, _ := types.IntrinsicGas(msg.Data(), msg.AccessList(), msg.AuthorizationList(), msg.To() == nil, r)
	from, _ := msg.From()
	msg = types.NewMessage(from, msg.To(), msg.Nonce(), msg.GetTxInternalData().GetAmount(), msg.Gas(), msg.GasPrice(), msg.GasFeeCap(), msg.GasTipCap(), msg.Data(), true, updatedIntrinsicGas, dataTokens, msg.AccessList(), msg.AuthorizationList())

	// Gas prices are calculated in eth
	if e.baseFee != nil {
		evm.GasPrice = math.BigMin(new(big.Int).Add(msg.GasTipCap(), e.baseFee), msg.GasFeeCap())
	}
}

func (e *eestEngine) Initialize(chain consensus.ChainReader, header *types.Header, state *state.StateDB) {
	if chain.Config().IsPragueForkEnabled(header.Number) {
		context := blockchain.NewEVMBlockContext(header, chain, nil)
		vmenv := vm.NewEVM(context, vm.TxContext{}, state, chain.Config(), &vm.Config{})
		blockchain.ProcessParentBlockHash(header, vmenv, state, chain.Config().Rules(header.Number))
	}
}

func (e *eestEngine) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	ethReward := common.Big0
	for _, receipt := range receipts {
		for _, tx := range txs {
			if tx.Hash() != receipt.TxHash {
				continue
			}

			ethGasPrice := tx.GasPrice()
			if e.baseFee != nil {
				ethGasPrice = math.BigMin(new(big.Int).Add(tx.GasTipCap(), e.baseFee), tx.GasFeeCap())
			}

			ethReward = new(big.Int).Add(ethReward, calculateEthMiningReward(ethGasPrice, tx.GasFeeCap(), tx.GasTipCap(), e.baseFee, receipt.GasUsed, chain.Config().Rules(header.Number)))
		}
	}

	state.AddBalance(header.Rewardbase, ethReward)
	header.Root = state.IntermediateRoot(true)

	return types.NewBlock(header, txs, receipts), nil
}

func (e *eestEngine) applyHeader(h TestHeader) {
	e.baseFee = h.BaseFee
	e.gasLimit = h.GasLimit
}

func (t *BlockTest) Run() error {
	config, ok := Forks[t.json.Network]
	if !ok {
		return UnsupportedForkError{t.json.Network}
	}
	config.SetDefaults()
	// Since we calculate the baseFee differently than eth, we will set it to 0 to turn off the gas fee.
	config.Governance.KIP71 = &params.KIP71Config{
		LowerBoundBaseFee:         0,
		UpperBoundBaseFee:         0,
		GasTarget:                 0,
		MaxBlockGasUsedForBaseFee: 0,
		BaseFeeDenominator:        0,
	}
	config.Governance.Reward = &params.RewardConfig{
		DeferredTxFee: true,
	}
	blockchain.InitDeriveSha(config)

	// import pre accounts & construct test genesis block & state root
	db := database.NewMemoryDBManager()
	gblock, err := t.genesis(config).Commit(common.Hash{}, db)
	if err != nil {
		return err
	}

	st := MakePreState(db, t.json.Pre, true, config.Rules(gblock.Number()))
	simulatedRoot, err := useEthStateRoot(st)
	if err != nil {
		return err
	}
	if simulatedRoot != t.json.Genesis.StateRoot {
		return fmt.Errorf("genesis block state root does not match test: computed=%x, test=%x", gblock.Root().Bytes()[:6], t.json.Genesis.StateRoot[:6])
	}

	tracer := vm.NewStructLogger(nil)
	chain, err := blockchain.NewBlockChain(db, nil, config, &eestEngine{Gxhash: gxhash.NewShared()}, vm.Config{Debug: true, Tracer: tracer, ComputationCostLimit: params.OpcodeComputationCostLimitInfinite})
	if err != nil {
		return err
	}
	defer chain.Stop()

	_, err = t.insertBlocksFromTx(chain, *gblock, db, tracer)
	if err != nil {
		return err
	}

	newDB, err := chain.State()
	if err != nil {
		return err
	}
	if err = t.validatePostState(newDB); err != nil {
		return fmt.Errorf("post state validation failed: %v", err)
	}

	return nil
}

func (t *BlockTest) genesis(config *params.ChainConfig) *blockchain.Genesis {
	return &blockchain.Genesis{
		Config:     config,
		Timestamp:  t.json.Genesis.Timestamp,
		ParentHash: t.json.Genesis.ParentHash,
		ExtraData:  t.json.Genesis.ExtraData,
		GasUsed:    t.json.Genesis.GasUsed,
		BlockScore: t.json.Genesis.Number,
		Alloc:      t.json.Pre,
	}
}

/*
See https://github.com/ethereum/tests/wiki/Blockchain-Tests-II

	Whether a block is valid or not is a bit subtle, it's defined by presence of
	blockHeader and transactions fields. If they are missing, the block is
	invalid and we must verify that we do not accept it.

	Since some tests mix valid and invalid blocks we need to check this for every block.

	If a block is invalid it does not necessarily fail the test, if it's invalidness is
	expected we are expected to ignore it and continue processing and then validate the
	post state.
*/
func (t *BlockTest) insertBlocks(bc *blockchain.BlockChain, preBlock *types.Block) ([]btBlock, error) {
	validBlocks := make([]btBlock, 0)
	latestParentHash := preBlock.Hash()
	latestRoot := preBlock.Root()
	// insert the test blocks, which will execute all transactions
	for _, b := range t.json.Blocks {
		cb, err := b.decode(latestParentHash, latestRoot)
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, fmt.Errorf("Block RLP decoding failed when expected to succeed: %v", err)
			}
		}
		// RLP decoding worked, try to insert into chain:
		latestParentHash = cb.Hash()
		latestRoot = cb.Root()
		blocks := types.Blocks{cb}
		i, err := bc.InsertChain(blocks)
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, fmt.Errorf("Block #%v insertion into chain failed: %v", blocks[i].Number(), err)
			}
		}
		if b.BlockHeader == nil {
			return nil, fmt.Errorf("Block insertion should have failed")
		}

		// validate RLP decoding by checking all values against test file JSON
		if err = validateHeader(b.BlockHeader, cb.Header()); err != nil {
			return nil, fmt.Errorf("Deserialised block header validation failed: %v", err)
		}
		validBlocks = append(validBlocks, b)
	}
	return validBlocks, nil
}

type rewardList struct {
	kaiaReward *big.Int
	ethReward  *big.Int
}

func (t *BlockTest) insertBlocksFromTx(bc *blockchain.BlockChain, gBlock types.Block, db database.DBManager, tracer *vm.StructLogger) ([]btBlock, error) {
	validBlocks := make([]btBlock, 0)
	preBlock := &gBlock

	// insert the test blocks, which will execute all transactions
	for _, b := range t.json.Blocks {
		txs, header, err := b.decodeTx()
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, fmt.Errorf("Block RLP decoding failed when expected to succeed: %v", err)
			}
		}

		if e := bc.Engine().(interface{ applyHeader(TestHeader) }); e != nil {
			e.applyHeader(header)
		}

		// var maxFeePerGas *big.Int
		blocks, _ := blockchain.GenerateChain(bc.Config(), preBlock, bc.Engine(), db, 1, func(i int, b *blockchain.BlockGen) {
			b.SetRewardbase(common.Address(header.Coinbase))
			for _, tx := range txs {
				b.AddTxWithChainEvenHasError(bc, tx)
			}
		})
		preBlock = blocks[0]

		if header.GasUsed != blocks[0].GasUsed() {
			return nil, fmt.Errorf("Unexpected GasUsed error (Expected: %v, Actual: %v)", header.GasUsed, blocks[0].GasUsed())
		}

		i, err := bc.InsertChain(blocks)
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, fmt.Errorf("Block #%v insertion into chain failed: %v", blocks[i].Number(), err)
			}
		}
		if b.BlockHeader == nil {
			return nil, errors.New("Block insertion should have failed")
		}

		validBlocks = append(validBlocks, b)
	}
	return validBlocks, nil
}

func validateHeader(h *btHeader, h2 *types.Header) error {
	if h.Bloom != h2.Bloom {
		return fmt.Errorf("Bloom: want: %x have: %x", h.Bloom, h2.Bloom)
	}
	if h.Number.Cmp(h2.Number) != 0 {
		return fmt.Errorf("Number: want: %v have: %v", h.Number, h2.Number)
	}
	if h.ParentHash != h2.ParentHash {
		return fmt.Errorf("Parent hash: want: %x have: %x", h.ParentHash, h2.ParentHash)
	}
	if h.ReceiptTrie != h2.ReceiptHash {
		return fmt.Errorf("Receipt hash: want: %x have: %x", h.ReceiptTrie, h2.ReceiptHash)
	}
	if h.TransactionsTrie != h2.TxHash {
		return fmt.Errorf("Tx hash: want: %x have: %x", h.TransactionsTrie, h2.TxHash)
	}
	if h.StateRoot != h2.Root {
		return fmt.Errorf("State hash: want: %x have: %x", h.StateRoot, h2.Root)
	}
	if !bytes.Equal(h.ExtraData, h2.Extra) {
		return fmt.Errorf("Extra data: want: %x have: %x", h.ExtraData, h2.Extra)
	}
	if h.GasUsed != h2.GasUsed {
		return fmt.Errorf("GasUsed: want: %d have: %d", h.GasUsed, h2.GasUsed)
	}
	if h.Timestamp != h2.Time.Uint64() {
		return fmt.Errorf("TimestampGa: want: %v have: %v", h.Timestamp, h2.Time)
	}
	return nil
}

func (t *BlockTest) validatePostState(statedb *state.StateDB) error {
	// validate post state accounts in test file against what we have in state db
	for addr, acct := range t.json.Post {
		// address is indirectly verified by the other fields, as it's the db key
		code2 := statedb.GetCode(addr)
		balance2 := statedb.GetBalance(addr)
		nonce2 := statedb.GetNonce(addr)
		if !bytes.Equal(code2, acct.Code) {
			return fmt.Errorf("account code mismatch for addr: %s want: %v have: %s", addr.String(), acct.Code, hex.EncodeToString(code2))
		}
		if balance2.Cmp(acct.Balance) != 0 {
			return fmt.Errorf("account balance mismatch for addr: %s, want: %d, have: %d", addr.String(), acct.Balance, balance2)
		}
		if nonce2 != acct.Nonce {
			return fmt.Errorf("account nonce mismatch for addr: %s want: %d have: %d", addr.String(), acct.Nonce, nonce2)
		}
	}

	err := t.validateStorage(statedb)
	if err != nil {
		return err
	}

	return nil
}

// validateStorage validates storage while considering the difference between Kana and Ethereum.
func (t *BlockTest) validateStorage(statedb *state.StateDB) error {
	beaconRootsAddress := common.HexToAddress("0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02")     // EIP-4788
	depositContractAddress := common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa") // EIP-6110

	// check the number of account
	accountNum := 0
	statedb.ForEachAccount(func(addr common.Address, data account.Account) {
		accountNum++
	})
	if accountNum != len(t.json.Post) {
		return fmt.Errorf("the number of accounts mismatch want: %v have: %v", len(t.json.Post), accountNum)
	}

	for addr, acct := range t.json.Post {
		// EIP-4788 and EIP-6110 aren't supported
		if addr == beaconRootsAddress || addr == depositContractAddress {
			continue
		}

		storageSize := 0
		statedb.ForEachStorage(addr, func(key, value common.Hash) bool {
			storageSize++
			return true
		})
		if storageSize != len(acct.Storage) {
			return fmt.Errorf("account storage size mismatch for addr: %s, want: %v, have: %v", addr, len(acct.Storage), storageSize)
		}

		// the size of HistoryStorageAddress is the same but the storage data is different
		if addr == params.HistoryStorageAddress {
			continue
		}

		for k, v := range acct.Storage {
			v2 := statedb.GetState(addr, k)
			if v2 != v {
				return fmt.Errorf("account storage mismatch for addr: %s, slot: %x, want: %x, have: %x", addr, k, v, v2)
			}
		}
	}
	return nil
}

func (t *BlockTest) validateImportedHeaders(cm *blockchain.BlockChain, validBlocks []btBlock) error {
	// to get constant lookup when verifying block headers by hash (some tests have many blocks)
	bmap := make(map[common.Hash]btBlock, len(t.json.Blocks))
	for _, b := range validBlocks {
		bmap[b.BlockHeader.Hash] = b
	}
	// iterate over blocks backwards from HEAD and validate imported
	// headers vs test file. some tests have reorgs, and we import
	// block-by-block, so we can only validate imported headers after
	// all blocks have been processed by BlockChain, as they may not
	// be part of the longest chain until last block is imported.
	for b := cm.CurrentBlock(); b != nil && b.NumberU64() != 0; b = cm.GetBlockByHash(b.Header().ParentHash) {
		if err := validateHeader(bmap[b.Hash()].BlockHeader, b.Header()); err != nil {
			return fmt.Errorf("Imported block header validation failed: %v", err)
		}
	}
	return nil
}

type TestHeader struct {
	ParentHash       common.Hash
	UncleHash        common.Hash
	Coinbase         []byte
	Root             common.Hash
	TxHash           common.Hash
	ReceiptHash      common.Hash
	Bloom            types.Bloom
	Difficulty       *big.Int
	Number           *big.Int
	GasLimit         uint64
	GasUsed          uint64
	Time             *big.Int
	Extra            []byte
	MixHash          common.Hash
	Nonce            []byte
	BaseFee          *big.Int     `rlp:"optional"`
	WithdrawalsHash  *common.Hash `rlp:"optional"`
	BlobGasUsed      *uint64      `rlp:"optional"`
	ExcessBlobGas    *uint64      `rlp:"optional"`
	ParentBeaconRoot *common.Hash `rlp:"optional"`
	RequestsHash     *common.Hash `rlp:"optional"`
}

// Modify the decode function
func (bb *btBlock) decode(latestParentHash common.Hash, latestRoot common.Hash) (*types.Block, error) {
	data, err := hexutil.Decode(bb.Rlp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %v", err)
	}

	fmt.Printf("Debug: Full RLP hex: %x\n", data)

	// First decode just the raw RLP list
	s := rlp.NewStream(bytes.NewReader(data), 0)
	kind, size, err := s.Kind()
	if err != nil {
		return nil, fmt.Errorf("failed to get RLP kind: %v", err)
	}
	fmt.Printf("Debug: RLP kind: %v, size: %d\n", kind, size)

	if kind != rlp.List {
		return nil, fmt.Errorf("expected RLP list, got %v", kind)
	}

	// Manual decoding approach
	if _, err := s.List(); err != nil {
		return nil, fmt.Errorf("failed to enter outer list: %v", err)
	}

	// Decode header
	var header TestHeader
	if err := s.Decode(&header); err != nil {
		return nil, fmt.Errorf("failed to decode header: %v", err)
	}

	// Decode transactions
	var txs []*types.Transaction
	if err := s.Decode(&txs); err != nil {
		return nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	// Convert header
	var rewardbase common.Address
	if len(header.Coinbase) > 0 {
		copy(rewardbase[:], header.Coinbase[:20])
	}

	block := types.NewBlockWithHeader(&types.Header{
		ParentHash:   latestParentHash,
		Rewardbase:   rewardbase,
		Root:         latestRoot,
		TxHash:       header.TxHash,
		ReceiptHash:  header.ReceiptHash,
		Bloom:        header.Bloom,
		BlockScore:   params.GenesisBlockScore,
		Number:       header.Number,
		GasUsed:      header.GasUsed,
		Time:         header.Time,
		TimeFoS:      0,
		Extra:        header.Extra,
		Governance:   []byte{},
		Vote:         []byte{},
		BaseFee:      header.BaseFee,
		RandomReveal: []byte{},
		MixHash:      header.MixHash[:],
	})

	return block.WithBody(txs), nil
}

// Modify the decode function
func (bb *btBlock) decodeTx() (types.Transactions, TestHeader, error) {
	data, err := hexutil.Decode(bb.Rlp)
	if err != nil {
		return nil, TestHeader{}, fmt.Errorf("failed to decode hex: %v", err)
	}

	// First decode just the raw RLP list
	s := rlp.NewStream(bytes.NewReader(data), 0)
	kind, _, err := s.Kind()
	if err != nil {
		return nil, TestHeader{}, fmt.Errorf("failed to get RLP kind: %v", err)
	}

	if kind != rlp.List {
		return nil, TestHeader{}, fmt.Errorf("expected RLP list, got %v", kind)
	}

	// Manual decoding approach
	if _, err := s.List(); err != nil {
		return nil, TestHeader{}, fmt.Errorf("failed to enter outer list: %v", err)
	}

	// Decode header
	var header TestHeader
	if err := s.Decode(&header); err != nil {
		return nil, TestHeader{}, fmt.Errorf("failed to decode header: %v", err)
	}

	// Decode transactions
	var txs types.Transactions
	if _, err := s.List(); err != nil {
		return nil, TestHeader{}, fmt.Errorf("failed to enter outer list: %v", err)
	}

	// Self decode to convert to kaia's eth tx type
	for {
		var tx types.Transaction
		kind, _, err := s.Kind()
		if err == rlp.EOL {
			break
		} else if err != nil {
			return nil, TestHeader{}, fmt.Errorf("failed to get RLP kind: %v", err)
		}

		txdata, _ := s.Raw()
		ethTxDataInKaia := []byte{}
		switch kind {
		case rlp.List: // case of legacy
			ethTxDataInKaia = txdata
		case rlp.String: // case of envelope
			var ethTypeIndex int
			if txdata[0] < 0xb7 {
				ethTypeIndex = 1
			} else {
				ethTypeIndex = int(txdata[0] - 0xb7 + 1)
			}
			switch txdata[ethTypeIndex] {
			case 1, 2, 4: // eth transaction types whick kaia support
				ethTxDataInKaia = append([]byte{byte(types.EthereumTxTypeEnvelope)}, txdata[ethTypeIndex:]...)
			default:
				ethTxDataInKaia = txdata
			}
		default:
			return nil, TestHeader{}, fmt.Errorf("failed to get RLP kind: %v", err)
		}

		streamForTx := rlp.NewStream(bytes.NewReader(ethTxDataInKaia), 0)
		if err := streamForTx.Decode(&tx); err == rlp.EOL {
			break
		} else if err != nil {
			return nil, TestHeader{}, fmt.Errorf("failed to decode transaction: %v", err)
		}
		txs = append(txs, &tx)
	}

	return txs, header, nil
}
