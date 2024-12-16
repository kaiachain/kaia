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
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
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

func (t *BlockTest) Run() error {
	config, ok := Forks[t.json.Network]
	if !ok {
		return UnsupportedForkError{t.json.Network}
	}
	config.SetDefaults()
	config.Governance.KIP71 = &params.KIP71Config{
		LowerBoundBaseFee:         0,
		UpperBoundBaseFee:         0,
		GasTarget:                 0,
		MaxBlockGasUsedForBaseFee: 0,
		BaseFeeDenominator:        0,
	}
	blockchain.InitDeriveSha(config)

	// import pre accounts & construct test genesis block & state root
	db := database.NewMemoryDBManager()
	gblock, err := t.genesis(config).Commit(common.Hash{}, db)
	if err != nil {
		return err
	}

	// Our header fields are not compatible with Eth, so you can skip this.
	// if gblock.Hash() != t.json.Genesis.Hash {
	// 	return fmt.Errorf("genesis block hash doesn't match test: computed=%x, test=%x", gblock.Hash().Bytes()[:6], t.json.Genesis.Hash[:6])
	// }

	st := MakePreState(db, t.json.Pre, true, config.Rules(gblock.Number()))
	simulatedRoot, err := useEthStateRoot(st)
	if err != nil {
		return err
	}
	if simulatedRoot != t.json.Genesis.StateRoot {
		return fmt.Errorf("genesis block state root does not match test: computed=%x, test=%x", gblock.Root().Bytes()[:6], t.json.Genesis.StateRoot[:6])
	}

	// TODO-Kaia: Replace gxhash with istanbul
	chain, err := blockchain.NewBlockChain(db, nil, config, gxhash.NewShared(), vm.Config{})
	if err != nil {
		return err
	}
	defer chain.Stop()

	// validBlocks, err := t.insertBlocks(chain, gblock)
	_, rewardMap, err := t.insertBlocksFromTx(chain, gblock, db)
	if err != nil {
		return err
	}

	// no need
	// cmlast := chain.CurrentBlock().Hash()
	// if common.Hash(t.json.BestBlock) != cmlast {
	// 	return fmt.Errorf("last block hash validation mismatch: want: %x, have: %x", t.json.BestBlock, cmlast)
	// }

	newDB, err := chain.State()
	if err != nil {
		return err
	}
	if err = t.validatePostState(newDB, rewardMap); err != nil {
		return fmt.Errorf("post state validation failed: %v", err)
	}

	// return t.validateImportedHeaders(chain, validBlocks)
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

func (t *BlockTest) insertBlocksFromTx(bc *blockchain.BlockChain, preBlock *types.Block, db database.DBManager) ([]btBlock, map[common.Address]rewardList, error) {
	validBlocks := make([]btBlock, 0)
	rewardMap := map[common.Address]rewardList{}
	// insert the test blocks, which will execute all transactions
	for _, b := range t.json.Blocks {
		txs, rewardBase, baseFeePerGas, err := b.decodeTx()
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, nil, fmt.Errorf("Block RLP decoding failed when expected to succeed: %v", err)
			}
		}
		// RLP decoding worked, try to insert into chain:
		kaiaReward := common.Big0
		ethReward := common.Big0
		// var maxFeePerGas *big.Int
		blocks, receiptsList := blockchain.GenerateChain(bc.Config(), preBlock, bc.Engine(), db, 1, func(i int, b *blockchain.BlockGen) {
			b.SetRewardbase(rewardBase)
			for _, tx := range txs {
				b.AddTx(tx)
			}
		})
		// The reward calculation is different for kaia and eth, and this will be deducted from the state later.
		for _, receipt := range receiptsList[0] {
			for _, tx := range blocks[0].Body().Transactions {
				if tx.Hash() != receipt.TxHash {
					continue
				}
				// Record kaia's reward.
				kaiaReward = new(big.Int).Add(kaiaReward, new(big.Int).Mul(new(big.Int).SetUint64(receipt.GasUsed), tx.GasPrice()))
				fmt.Println(kaiaReward, tx.GasPrice(), tx.EffectiveGasPrice(blocks[0].Header(), bc.Config()), "naazenaaze fkldsjakfjdklsajklf")

				// Record eth's reward.
				ethGasPrice := tx.GasPrice()
				if baseFeePerGas != nil {
					ethGasPrice = math.BigMin(new(big.Int).Add(tx.GasTipCap(), baseFeePerGas), tx.GasFeeCap())
				}
				ethReward = new(big.Int).Add(ethReward, calculateEthMiningReward(ethGasPrice, tx.GasFeeCap(), tx.GasTipCap(), baseFeePerGas,
					receipt.GasUsed, bc.Config().Rules(blocks[0].Header().Number)))
			}
		}
		rewardMap[rewardBase] = rewardList{
			kaiaReward: kaiaReward,
			ethReward:  ethReward,
		}

		i, err := bc.InsertChain(blocks)
		if err != nil {
			if b.BlockHeader == nil {
				continue // OK - block is supposed to be invalid, continue with next block
			} else {
				return nil, nil, fmt.Errorf("Block #%v insertion into chain failed: %v", blocks[i].Number(), err)
			}
		}
		if b.BlockHeader == nil {
			return nil, nil, fmt.Errorf("Block insertion should have failed")
		}

		validBlocks = append(validBlocks, b)
	}
	return validBlocks, rewardMap, nil
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

func (t *BlockTest) validatePostState(statedb *state.StateDB, rewardMap map[common.Address]rewardList) error {
	// validate post state accounts in test file against what we have in state db
	for addr, acct := range t.json.Post {
		if rewardList, exist := rewardMap[addr]; exist {
			// In the case of rewardBaseAddress, the Kaia reward will be deducted once.
			statedb.SubBalance(addr, rewardList.kaiaReward)
			statedb.AddBalance(addr, rewardList.ethReward)
		}

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

// func (bb *btBlock) decode() (*types.Block, error) {
// 	data, err := hexutil.Decode(bb.Rlp)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var b types.Block
// 	err = rlp.DecodeBytes(data, &b)
// 	return &b, err
// }

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
func (bb *btBlock) decodeTx() (types.Transactions, common.Address, *big.Int, error) {
	data, err := hexutil.Decode(bb.Rlp)
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to decode hex: %v", err)
	}

	// First decode just the raw RLP list
	s := rlp.NewStream(bytes.NewReader(data), 0)
	kind, _, err := s.Kind()
	if err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to get RLP kind: %v", err)
	}

	if kind != rlp.List {
		return nil, common.Address{}, nil, fmt.Errorf("expected RLP list, got %v", kind)
	}

	// Manual decoding approach
	if _, err := s.List(); err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to enter outer list: %v", err)
	}

	// Decode header
	var header TestHeader
	if err := s.Decode(&header); err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to decode header: %v", err)
	}

	// Decode transactions
	var txs types.Transactions
	if err := s.Decode(&txs); err != nil {
		return nil, common.Address{}, nil, fmt.Errorf("failed to decode transactions: %v", err)
	}

	// Convert header
	var rewardbase common.Address
	if len(header.Coinbase) > 0 {
		copy(rewardbase[:], header.Coinbase[:20])
	}

	return txs, rewardbase, header.BaseFee, nil
}

// func useEthBlockHash(r params.Rules, json *btJSON) common.Hash {
// 	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L241-L249
// 	type ethHeader struct {
// 		ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
// 		UncleHash   common.Hash    `json:"sha3Uncles"       gencodec:"required"`
// 		Coinbase    common.Address `json:"miner"`
// 		Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
// 		TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
// 		ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
// 		Bloom       types.Bloom    `json:"logsBloom"        gencodec:"required"`
// 		Difficulty  *big.Int       `json:"difficulty"       gencodec:"required"`
// 		Number      *big.Int       `json:"number"           gencodec:"required"`
// 		GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
// 		GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
// 		Time        uint64         `json:"timestamp"        gencodec:"required"`
// 		Extra       []byte         `json:"extraData"        gencodec:"required"`
// 		MixDigest   common.Hash    `json:"mixHash"`
// 		Nonce       uint64         `json:"nonce"`

// 		// BaseFee was added by EIP-1559 and is ignored in legacy headers.
// 		BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`

// 		// WithdrawalsHash was added by EIP-4895 and is ignored in legacy headers.
// 		WithdrawalsHash *common.Hash `json:"withdrawalsRoot" rlp:"optional"`

// 		// BlobGasUsed was added by EIP-4844 and is ignored in legacy headers.
// 		BlobGasUsed *uint64 `json:"blobGasUsed" rlp:"optional"`

// 		// ExcessBlobGas was added by EIP-4844 and is ignored in legacy headers.
// 		ExcessBlobGas *uint64 `json:"excessBlobGas" rlp:"optional"`

// 		// ParentBeaconRoot was added by EIP-4788 and is ignored in legacy headers.
// 		ParentBeaconRoot *common.Hash `json:"parentBeaconBlockRoot" rlp:"optional"`

// 		// RequestsHash was added by EIP-7685 and is ignored in legacy headers.
// 		RequestsHash *common.Hash `json:"requestsRoot" rlp:"optional"`
// 	}

// 	header := ethHeader{
// 		ParentHash:       json.Genesis.ParentHash,
// 		UncleHash:        json.Genesis.UncleHash,
// 		Coinbase:         json.Genesis.Coinbase,
// 		Root:             json.Genesis.StateRoot,
// 		TxHash:           json.Genesis.Hash,
// 		ReceiptHash:      json.Genesis.ReceiptTrie,
// 		Bloom:            json.Genesis.Bloom,
// 		Difficulty:       json.Genesis.Difficulty,
// 		Number:           json.Genesis.Number,
// 		GasLimit:         json.Genesis.GasLimit,
// 		GasUsed:          json.Genesis.GasUsed,
// 		Time:             json.Genesis.Timestamp,
// 		Extra:            json.Genesis.ExtraData,
// 		MixDigest:        json.Genesis.MixHash,
// 		Nonce:            json.Genesis.Nonce,
// 		BaseFee:          json.Genesis.BaseFeePerGas,
// 		WithdrawalsHash:  json.Genesis.WithdrawalsRoot,
// 		BlobGasUsed:      json.Genesis.BlobGasUsed,
// 		ExcessBlobGas:    json.Genesis.ExcessBlobGas,
// 		ParentBeaconRoot: json.Genesis.ParentBeaconBlockRoot,
// 	}
// 	return header.Hash()
// }
