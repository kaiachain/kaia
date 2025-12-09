// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/genesis.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var (
	errGenesisNoConfig = errors.New("genesis has no chain configuration")
	errNoGenesis       = errors.New("genesis block is not provided")
)

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	Governance []byte              `json:"governanceData"`
	BlockScore *big.Int            `json:"blockScore"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// copy copies the genesis.
func (g *Genesis) copy() *Genesis {
	if g != nil {
		cpy := *g
		if g.Config != nil {
			conf := g.Config.Copy()
			cpy.Config = conf
		}
		return &cpy
	}
	return nil
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	BlockScore *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// findBlockWithState returns the latest block with state.
func findBlockWithState(db database.DBManager) *types.Block {
	headBlock := db.ReadBlockByHash(db.ReadHeadBlockHash())
	if headBlock == nil {
		headBlock = db.ReadBlockByHash(db.ReadHeadBlockBackupHash())
		if headBlock == nil {
			logger.Crit("failed to read head block by head block hash")
		}
	}

	startBlock := headBlock
	for _, err := state.New(headBlock.Root(), state.NewDatabase(db), nil, nil); err != nil; {
		if headBlock.NumberU64() == 0 {
			logger.Crit("failed to find state from the head block to the genesis block",
				"headBlockNum", headBlock.NumberU64(),
				"headBlockHash", headBlock.Hash().String(), "headBlockRoot", headBlock.Root().String())
		}
		headBlock = db.ReadBlockByNumber(headBlock.NumberU64() - 1)
		if headBlock == nil {
			logger.Crit("failed to read previous block by head block number")
		}
		logger.Warn("found previous block", "blockNum", headBlock.NumberU64())
	}
	logger.Info("found the latest block with state",
		"blockNum", headBlock.NumberU64(), "startedNum", startBlock.NumberU64())
	return headBlock
}

// SetupGenesisBlock writes or updates the genesis block and ChainConfig in db.
// Note that a "genesis block" is composed of a regular Block and ChainConfig, hence the return types.
// The block and ChainConfig that will be returned/written to DB are:
// (1) stored ghash == nil: Commit the provided genesis and config, defaulting to Mainnet if absent.
// (2) stored ghash != nil && stored config == nil: Same as above.
// (3) stored ghash != nil && stored config != nil: Update ChainConfigDB only if it's compatible with the existing config.
// For the case (3), the stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// Notes:
// - `genesis != nil` is the normal case for Mainnet/Kairos `cn.New()`.
// - `genesis != nil` is the normal case for other networks `initGenesis()`.
// - `stored ghash != nil && genesis == nil` is the normal case for other networks `cn.New()`.
func SetupGenesisBlock(db database.DBManager, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {
	if genesis != nil && genesis.Config == nil {
		return nil, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no ghash genesis block.
	ghash := db.ReadCanonicalHash(0)
	if (ghash == common.Hash{}) {
		if genesis == nil {
			logger.Info("Writing default Mainnet genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			logger.Info("Writing custom genesis block")
		}
		InitDeriveSha(genesis.Config)
		block, err := genesis.Commit(common.Hash{}, db)
		if err != nil {
			return nil, common.Hash{}, err
		}
		return genesis.Config, block.Hash(), err
	}

	// Genesis block exists, and another genesis is supplied. Abort if they are different, because we don't want to overwrite the genesis block.
	if genesis != nil {
		// This is the usual path which does not overwrite genesis block with the new one.
		// Make sure the provided genesis is equal to the stored one.
		InitDeriveSha(genesis.Config)
		hash := genesis.ToBlock(common.Hash{}, nil).Hash()
		if hash != ghash {
			return nil, common.Hash{}, &GenesisMismatchError{ghash, hash}
		}
	}

	// The genesis block is present in the database but the corresponding state might not.
	// Because the trie can be partially corrupted, we always commit the trie.
	// It can happen in a state migrated database or live pruned database.
	if db.GetDomainsManager() == nil { // FlatTrie disallows re-commiting the block lower than the head block.
		commitGenesisState(genesis, db)
	}

	storedCfg, err := db.ReadChainConfig(ghash)
	if err != nil {
		logger.Crit("Failed to read chain config", "err", err)
	}
	// Genesis block exists, but no ChainConfig. Re-commit to store the ChainConfig.
	if storedCfg == nil {
		// Ensure the stored genesis block matches with the given genesis. Private
		// networks must explicitly specify the genesis in the config file, mainnet
		// genesis will be used as default and the initialization will always fail.
		if genesis == nil {
			logger.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			logger.Info("Writing custom genesis block")
		}

		InitDeriveSha(genesis.Config)
		block, err := genesis.Commit(common.Hash{}, db)
		if err != nil {
			return nil, common.Hash{}, err
		}

		return genesis.Config, block.Hash(), err
	} else {
		if storedCfg.Governance == nil {
			logger.Crit("Failed to read governance. storedcfg.Governance == nil")
		}
		if storedCfg.Governance.Reward == nil {
			logger.Crit("Failed to read governance. storedcfg.Governance.Reward == nil")
		}
	}

	// Get the existing chain configuration.
	newCfg := configOrDefault(genesis, ghash, storedCfg)
	if err := newCfg.CheckConfigForkOrder(); err != nil {
		return nil, common.Hash{}, err
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := db.ReadHeaderNumber(db.ReadHeadHeaderHash())
	if height == nil {
		return nil, common.Hash{}, errors.New("missing block number for head header hash")
	}
	compatErr := storedCfg.CheckCompatible(newCfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		return newCfg, ghash, compatErr
	}
	db.WriteChainConfig(ghash, newCfg)
	return newCfg, ghash, nil
}

func configOrDefault(g *Genesis, ghash common.Hash, storedCfg *params.ChainConfig) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetGenesisHash:
		return params.MainnetChainConfig
	case ghash == params.KairosGenesisHash:
		return params.KairosChainConfig
	default:
		return storedCfg
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(baseStateRoot common.Hash, db database.DBManager) *types.Block {
	if db == nil {
		// If db == nil, do not write to the real database. Here we supply a memory database as a placeholder.
		db = database.NewMemoryDBManager()
	}
	stateDB, _ := state.New(baseStateRoot, state.NewDatabase(db), nil, &statedb.TrieOpts{
		BaseBlockNumber: 0,
		CommitGenesis:   true,
	})
	rules := params.Rules{}
	if g.Config != nil {
		rules = g.Config.Rules(new(big.Int).SetUint64(g.Number))
	}
	for addr, account := range g.Alloc {
		originalCode := stateDB.GetCode(addr)
		_, ok := types.ParseDelegation(account.Code)
		isEOAWithCode := ok && rules.IsPrague
		isEOAWithoutCode := len(account.Code) == 0 && len(account.Storage) != 0 && rules.IsPrague
		isSCAWithCode := !isEOAWithCode && len(account.Code) != 0
		switch {
		case isEOAWithCode:
			// EOA with code. Usually SetCodeTx creates it. Unit tests may create it here in the genesis.
			stateDB.SetCodeToEOA(addr, account.Code, rules)
		case isEOAWithoutCode:
			// Represents an EOA that had code then nullified with another SetCodeTx. Unit tests may create it here in the genesis.
			stateDB.CreateEOA(addr, false, accountkey.NewAccountKeyLegacy())
		case isSCAWithCode:
			// Regular genesis smart contract account.
			stateDB.CreateSmartContractAccount(addr, params.CodeFormatEVM, rules)
			stateDB.SetCode(addr, account.Code)
		}
		// If account.Code is nil and originalCode is not nil,
		// just update the code and don't change the other states
		if len(account.Code) != 0 && originalCode != nil {
			logger.Warn("this address already has a not nil code, now the code of this address has been changed", "addr", addr.String())
			continue
		}
		for key, value := range account.Storage {
			stateDB.SetState(addr, key, value)
		}
		stateDB.AddBalance(addr, account.Balance)
		stateDB.SetNonce(addr, account.Nonce)
	}
	root := stateDB.IntermediateRoot(false)
	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		TimeFoS:    0,
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		Governance: g.Governance,
		GasUsed:    g.GasUsed,
		BlockScore: g.BlockScore,
		Root:       root,
	}
	if g.BlockScore == nil {
		head.BlockScore = params.GenesisBlockScore
	}
	if g.Config != nil && g.Config.IsMagmaForkEnabled(common.Big0) {
		if g.Config.Governance != nil && g.Config.Governance.KIP71 != nil {
			head.BaseFee = new(big.Int).SetUint64(g.Config.Governance.KIP71.LowerBoundBaseFee)
		} else {
			head.BaseFee = new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee)
		}
	}
	if g.Config != nil && g.Config.IsRandaoForkEnabled(common.Big0) {
		head.RandomReveal = params.ZeroRandomReveal
		head.MixHash = params.ZeroMixHash
	}

	stateDB.Commit(false)
	stateDB.Database().TrieDB().Commit(root, true, g.Number)

	return types.NewBlock(head, nil, nil)
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(baseStateRoot common.Hash, db database.DBManager) (*types.Block, error) {
	block := g.ToBlock(baseStateRoot, db)
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	db.WriteTd(block.Hash(), block.NumberU64(), g.BlockScore)
	db.WriteBlock(block)
	db.WriteReceipts(block.Hash(), block.NumberU64(), nil)
	db.WriteCanonicalHash(block.Hash(), block.NumberU64())
	db.WriteHeadBlockHash(block.Hash())
	db.WriteHeadHeaderHash(block.Hash())

	config := g.Config
	if config == nil {
		config = params.TestChainConfig
	}
	if err := config.CheckConfigForkOrder(); err != nil {
		return nil, err
	}
	db.WriteChainConfig(block.Hash(), config)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db database.DBManager) *types.Block {
	config := g.Config
	if config == nil {
		config = params.TestChainConfig
	}
	InitDeriveSha(config)

	block, err := g.Commit(common.Hash{}, db)
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given kei balance.
func GenesisBlockForTesting(db database.DBManager, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns the Mainnet genesis block.
// It is also used for default genesis block.
func DefaultGenesisBlock() *Genesis {
	ret := &Genesis{}
	if err := json.Unmarshal(mainnetGenesisJson, &ret); err != nil {
		logger.Crit("Error in Unmarshalling Mainnet Genesis Json", "err", err)
	}
	ret.Config = params.MainnetChainConfig.Copy()
	ret.Governance = SetGenesisGovernance(ret)
	InitDeriveSha(ret.Config)
	return ret
}

// DefaultKairosGenesisBlock returns the Kairos genesis block.
func DefaultKairosGenesisBlock() *Genesis {
	ret := &Genesis{}
	if err := json.Unmarshal(kairosGenesisJson, &ret); err != nil {
		logger.Crit("Error in Unmarshalling Kairos Genesis Json", "err", err)
	}
	ret.Config = params.KairosChainConfig.Copy()
	ret.Governance = SetGenesisGovernance(ret)
	InitDeriveSha(ret.Config)
	return ret
}

func DefaultTestGenesisBlock() *Genesis {
	ret := DefaultGenesisBlock()
	ret.Governance = nil
	return ret
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

func commitGenesisState(genesis *Genesis, db database.DBManager) {
	if genesis == nil {
		genesis = DefaultGenesisBlock()
	}
	// Run genesis.ToBlock() to calls StateDB.Commit() which writes the state trie.
	// But do not run genesis.Commit() which overwrites HeaderHash.
	InitDeriveSha(genesis.Config)
	genesis.ToBlock(common.Hash{}, db).Hash()
	logger.Info("Restored state trie for the genesis block")
}

type GovernanceSet map[string]interface{}

func SetGenesisGovernance(genesis *Genesis) []byte {
	g := make(GovernanceSet)
	governance := genesis.Config.Governance
	g["governance.governancemode"] = governance.GovernanceMode
	g["governance.governingnode"] = governance.GoverningNode
	g["governance.unitprice"] = genesis.Config.UnitPrice
	g["reward.mintingamount"] = governance.Reward.MintingAmount.String()
	g["reward.minimumstake"] = governance.Reward.MinimumStake.String()
	g["reward.ratio"] = governance.Reward.Ratio
	g["reward.useginicoeff"] = governance.Reward.UseGiniCoeff
	g["reward.deferredtxfee"] = governance.Reward.DeferredTxFee
	g["reward.stakingupdateinterval"] = governance.Reward.StakingUpdateInterval
	g["reward.proposerupdateinterval"] = governance.Reward.ProposerUpdateInterval
	g["istanbul.epoch"] = genesis.Config.Istanbul.Epoch
	g["istanbul.policy"] = genesis.Config.Istanbul.ProposerPolicy
	g["istanbul.committeesize"] = genesis.Config.Istanbul.SubGroupSize

	data, err := json.Marshal(g)
	if err != nil {
		logger.Crit("Error in marshaling governance data", "err", err)
	}
	ret, err := rlp.EncodeToBytes(data)
	if err != nil {
		logger.Crit("Error in RLP Encoding governance data", "err", err)
	}
	return ret
}
