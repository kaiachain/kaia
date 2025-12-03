// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from core/genesis_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultGenesisBlock tests the genesis block generation functions: DefaultGenesisBlock, DefaultKairosGenesisBlock
func TestDefaultGenesisBlock(t *testing.T) {
	block := DefaultGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.MainnetGenesisHash {
		t.Errorf("wrong Mainnet genesis hash, got %v, want %v", block.Hash(), params.MainnetGenesisHash)
	}
	block = DefaultKairosGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.KairosGenesisHash {
		t.Errorf("wrong Kairos genesis hash, got %v, want %v", block.Hash(), params.KairosGenesisHash)
	}
}

// TestHardCodedChainConfigUpdate tests the public network's chainConfig update.
func TestHardCodedChainConfigUpdate(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	type ExpectedReturn struct {
		config *params.ChainConfig
		hash   common.Hash
		err    error
	}
	type ExpectedDB struct {
		storedCfg *params.ChainConfig
		ghash     common.Hash
	}

	updateConfig := func(cfg *params.ChainConfig, blockNumber uint64) *params.ChainConfig {
		cfg = cfg.Copy()
		cfg.IstanbulCompatibleBlock = big.NewInt(int64(blockNumber))
		return cfg
	}

	tests := []struct {
		name         string
		fn           func(database.DBManager) (*params.ChainConfig, common.Hash, error)
		updateConfig func(cfg *params.ChainConfig, blockNumber uint64) *params.ChainConfig
		wantReturn   ExpectedReturn
		wantDB       ExpectedDB
	}{
		{
			// genesis.Config = returned config = stored config
			name: "Mainnet chainConfig update",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genesis := DefaultGenesisBlock()
				genesis.MustCommit(db)
				genesis.Config = updateConfig(params.MainnetChainConfig, 3)
				return SetupGenesisBlock(db, genesis)
			},
			wantReturn: ExpectedReturn{
				config: updateConfig(params.MainnetChainConfig, 3),
				hash:   params.MainnetGenesisHash,
				err:    nil,
			},
			wantDB: ExpectedDB{
				storedCfg: updateConfig(params.MainnetChainConfig, 3),
				ghash:     params.MainnetGenesisHash,
			},
		},
		// TODO-Kaia: add more Mainnet test cases after Mainnet hard fork block numbers are added
		{
			// Because of the fork-ordering check logic, the istanbulCompatibleBlock should be less than the londonCompatibleBlock
			// genesis.Config = returned config = stored config
			name: "Kairos chainConfig update - correct hard-fork block number order",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genesis := DefaultKairosGenesisBlock()
				genesis.MustCommit(db)
				genesis.Config = updateConfig(params.KairosChainConfig, 79999999)
				return SetupGenesisBlock(db, genesis)
			},
			wantReturn: ExpectedReturn{
				config: updateConfig(params.KairosChainConfig, 79999999),
				hash:   params.KairosGenesisHash,
				err:    nil,
			},
			wantDB: ExpectedDB{
				storedCfg: updateConfig(params.KairosChainConfig, 79999999),
				ghash:     params.KairosGenesisHash,
			},
		},
		{
			// This test fails because the new istanbulCompatibleBlock(90909999) is larger than londonCompatibleBlock(80295291)
			name: "Kairos chainConfig update - wrong hard-fork block number order",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genesis := DefaultKairosGenesisBlock()
				genesis.MustCommit(db)
				genesis.Config = updateConfig(params.KairosChainConfig, 90909999)
				return SetupGenesisBlock(db, genesis)
			},
			wantReturn: ExpectedReturn{
				config: nil,
				hash:   common.Hash{},
				err: fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
					"istanbulBlock", big.NewInt(90909999), "londonBlock", big.NewInt(80295291)),
			},
			wantDB: ExpectedDB{ // not overwritten
				storedCfg: params.KairosChainConfig,
				ghash:     params.KairosGenesisHash,
			},
		},
		{
			name: "incompatible config in DB",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := DefaultGenesisBlock()
				genesisBlock := genesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, genesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				blocks, _ := GenerateChain(genesis.Config, genesisBlock, faker.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newGenesis := genesis.copy()
				newGenesis.Config = updateConfig(params.MainnetChainConfig, 2)
				return SetupGenesisBlock(db, newGenesis)
			},
			wantReturn: ExpectedReturn{
				config: updateConfig(params.MainnetChainConfig, 2),
				hash:   params.MainnetGenesisHash,
				err: &params.ConfigCompatError{
					What:         "Istanbul Block",
					StoredConfig: params.MainnetChainConfig.IstanbulCompatibleBlock,
					NewConfig:    big.NewInt(2),
					RewindTo:     1,
				},
			},
			wantDB: ExpectedDB{
				storedCfg: params.MainnetChainConfig,
				ghash:     params.MainnetGenesisHash,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := database.NewMemoryDBManager()
			config, hash, err := test.fn(db)

			// Check the return values
			assert.Equal(t, test.wantReturn.config, config, test.name+": config is mismatching")
			assert.Equal(t, test.wantReturn.hash, hash, test.name+": hash is mismatching")
			assert.Equal(t, test.wantReturn.err, err, test.name+": err is mismatching")

			// Check DB
			ghash := db.ReadCanonicalHash(0)
			storedCfg, err := db.ReadChainConfig(ghash)
			assert.NoError(t, err)
			assert.Equal(t, test.wantDB.storedCfg, storedCfg, test.name+": stored chainConfig is mismatching")
			assert.Equal(t, test.wantDB.ghash, ghash, test.name+": stored genesis block is not compatible")
		})
	}
}

func TestSetupGenesis(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	type ExpectedReturn struct {
		config *params.ChainConfig
		hash   common.Hash
		err    error
	}
	type ExpectedDB struct {
		storedCfg *params.ChainConfig
		ghash     common.Hash
	}

	// NOTE: Do NOT move {mainnet,kairosGenesis,custom,new}Genesis.Config pointers.
	// Because they are referenced by expected values.
	// It is only safe to change config.IstanbulCompatibleBlock.
	var (
		mainnetGenesis    = DefaultGenesisBlock()
		kairosGenesis     = DefaultKairosGenesisBlock()
		customGenesisHash = common.HexToHash("0x4eb4035b7a09619a9950c9a4751cc331843f2373ef38263d676b4a132ba4059c")
		customChainId     = uint64(4343)
		customGenesis     = genCustomGenesisBlock(customChainId)
		newGenesis        = customGenesis.copy()
	)
	tests := []struct {
		name           string
		fn             func(database.DBManager) (*params.ChainConfig, common.Hash, error)
		expectedReturn ExpectedReturn
		expectedDB     ExpectedDB
	}{
		/*
			{
				name: "genesis without ChainConfig",
				fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
					return SetupGenesisBlock(db, new(Genesis))
				},
				expectedReturn: ExpectedReturn{
					err: errGenesisNoConfig,
				},
				expectedDB: ExpectedDB{
					ghash:     common.Hash{},
					storedCfg: nil,
				},
			},
		*/
		{
			name: "no block in DB, genesis == nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "no block in DB, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "no block in DB, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.KairosGenesisHash,
				config: params.KairosChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "no block in DB, genesis is customGenesis",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: customGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Mainnet block in DB, storedCfg is nil, genesis is nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.MainnetGenesisHash)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, storedCfg is nil, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.MainnetGenesisHash)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, storedCfg is nil, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.MainnetGenesisHash)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: params.KairosGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, storedCfg is nil, genesis is Custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.MainnetGenesisHash)
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: customGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Kairos block in DB, storedCfg is nil, genesis is nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.KairosGenesisHash)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   params.KairosGenesisHash,
				config: params.KairosChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, storedCfg is nil, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.KairosGenesisHash)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.KairosGenesisHash,
				config: params.KairosChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, storedCfg is nil, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.KairosGenesisHash)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: params.MainnetGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, storedCfg is nil, genesis is Custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				writeNilChainConfig(db, params.KairosGenesisHash)
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: customGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Custom block in DB, storedCfg is nil, genesis is nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				writeNilChainConfig(db, customGenesisHash)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: customGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Custom block in DB, storedCfg is nil, genesis is Custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				writeNilChainConfig(db, customGenesisHash)
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: customGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Custom block in DB, storedCfg is nil, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				writeNilChainConfig(db, customGenesisHash)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: customGenesisHash, New: params.MainnetGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Custom block in DB, storedCfg is nil, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				writeNilChainConfig(db, customGenesisHash)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: customGenesisHash, New: params.KairosGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Mainnet block in DB, genesis == nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.MainnetGenesisHash,
				config: params.MainnetChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: params.KairosGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Mainnet block in DB, genesis is custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genCustomGenesisBlock(customChainId))
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: customGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.MainnetGenesisHash,
				storedCfg: params.MainnetChainConfig,
			},
		},
		{
			name: "Kairos block in DB, genesis == nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   params.KairosGenesisHash,
				config: params.KairosChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   params.KairosGenesisHash,
				config: params.KairosChainConfig,
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: params.MainnetGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "Kairos block in DB, genesis is custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				DefaultKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: customGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     params.KairosGenesisHash,
				storedCfg: params.KairosChainConfig,
			},
		},
		{
			name: "custom block in DB, genesis == nil",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, nil)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: customGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "custom block in DB, genesis is custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, customGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: customGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "custom block in DB, genesis is Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, DefaultGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: customGenesisHash, New: params.MainnetGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "custom block in DB, genesis is Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, DefaultKairosGenesisBlock())
			},
			expectedReturn: ExpectedReturn{
				hash:   common.Hash{},
				config: nil,
				err:    &GenesisMismatchError{Stored: customGenesisHash, New: params.KairosGenesisHash},
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "custom block in DB, compatible newGenesis, head=1 < Istanbul=2 < NewIstanbul=5",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := customGenesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, customGenesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				head := 1
				blocks, _ := GenerateChain(customGenesis.Config, genesis, faker.NewFaker(), db, head, nil)
				bc.InsertChain(blocks)
				newGenesis.Config.IstanbulCompatibleBlock = big.NewInt(5)
				return SetupGenesisBlock(db, newGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: newGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: newGenesis.Config,
			},
		},
		{
			name: "custom block in DB, compatible newGenesis, head=1 < NewIstanbul=2 < Istanbul=5",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.Config.IstanbulCompatibleBlock = big.NewInt(5)
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := customGenesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, customGenesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				head := 1
				blocks, _ := GenerateChain(customGenesis.Config, genesis, faker.NewFaker(), db, head, nil)
				bc.InsertChain(blocks)
				newGenesis.Config.IstanbulCompatibleBlock = big.NewInt(2)
				return SetupGenesisBlock(db, newGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: newGenesis.Config,
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: newGenesis.Config,
			},
		},
		{
			name: "custom block in DB, incompatible newGenesis, Istanbul=2 < head=4 < NewIstanbul=5",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := customGenesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, customGenesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				head := 4
				blocks, _ := GenerateChain(customGenesis.Config, genesis, faker.NewFaker(), db, head, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newGenesis.Config.IstanbulCompatibleBlock = big.NewInt(5)
				return SetupGenesisBlock(db, newGenesis)
			},
			expectedReturn: ExpectedReturn{
				hash:   customGenesisHash,
				config: newGenesis.Config,
				err: &params.ConfigCompatError{
					What:         "Istanbul Block",
					StoredConfig: big.NewInt(2),
					NewConfig:    big.NewInt(5),
					RewindTo:     1,
				},
			},
			expectedDB: ExpectedDB{
				ghash:     customGenesisHash,
				storedCfg: customGenesis.Config,
			},
		},
		{
			name: "Mainnet chainConfig update",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				mainnetGenesis.MustCommit(db)
				mainnetGenesis.Config.IstanbulCompatibleBlock = big.NewInt(3)
				return SetupGenesisBlock(db, mainnetGenesis)
			},
			expectedReturn: ExpectedReturn{
				config: mainnetGenesis.Config,
				hash:   params.MainnetGenesisHash,
				err:    nil,
			},
			expectedDB: ExpectedDB{
				storedCfg: mainnetGenesis.Config,
				ghash:     params.MainnetGenesisHash,
			},
		},
		{
			// Because of the fork-ordering check logic, the istanbulCompatibleBlock should be less than the londonCompatibleBlock
			// genesis.Config = returned config = stored config
			name: "Kairos chainConfig update - correct hard-fork block number order",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				kairosGenesis.MustCommit(db)
				kairosGenesis.Config.IstanbulCompatibleBlock = big.NewInt(79999999)
				return SetupGenesisBlock(db, kairosGenesis)
			},
			expectedReturn: ExpectedReturn{
				config: kairosGenesis.Config,
				hash:   params.KairosGenesisHash,
				err:    nil,
			},
			expectedDB: ExpectedDB{
				storedCfg: kairosGenesis.Config,
				ghash:     params.KairosGenesisHash,
			},
		},
		{
			// This test fails because the new istanbulCompatibleBlock(90909999) is larger than londonCompatibleBlock(80295291)
			name: "Kairos chainConfig update - wrong hard-fork block number order",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				kairosGenesis.MustCommit(db)
				kairosGenesis.Config.IstanbulCompatibleBlock = big.NewInt(90909999)
				return SetupGenesisBlock(db, kairosGenesis)
			},
			expectedReturn: ExpectedReturn{
				config: nil,
				hash:   common.Hash{},
				err: fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
					"istanbulBlock", big.NewInt(90909999), "londonBlock", big.NewInt(80295291)),
			},
			expectedDB: ExpectedDB{ // not overwritten
				storedCfg: params.KairosChainConfig,
				ghash:     params.KairosGenesisHash,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db := database.NewMemoryDBManager()
			config, hash, err := test.fn(db)

			// Check the return values
			assert.Equal(t, test.expectedReturn.config, config)
			assert.Equal(t, test.expectedReturn.hash, hash)
			assert.Equal(t, test.expectedReturn.err, err)

			// Check DB
			ghash := db.ReadCanonicalHash(0)
			storedCfg, err := db.ReadChainConfig(ghash)
			assert.NoError(t, err)
			assert.Equal(t, test.expectedDB.storedCfg, storedCfg, test.name+": stored chainConfig is mismatching")
			assert.Equal(t, test.expectedDB.ghash, ghash, test.name+": stored genesis block is not compatible")

			// reset hardfork blocks
			mainnetGenesis.Config.IstanbulCompatibleBlock = new(big.Int).Set(params.MainnetChainConfig.IstanbulCompatibleBlock)
			kairosGenesis.Config.IstanbulCompatibleBlock = new(big.Int).Set(params.KairosChainConfig.IstanbulCompatibleBlock)
			customGenesis.Config.IstanbulCompatibleBlock = big.NewInt(2)
			newGenesis.Config.IstanbulCompatibleBlock = big.NewInt(2)
		})
	}
}

// Test that SetupGenesisBlock writes the genesis state trie if not present.
func TestGenesisRestoreState(t *testing.T) {
	db := database.NewMemoryDBManager()

	// Setup first to Commit the Genesis block
	_, hash, err := SetupGenesisBlock(db, nil)
	assert.Nil(t, err)
	assert.Equal(t, params.MainnetGenesisHash, hash)

	// Simulate state migration by deleting the state trie root node
	header := db.ReadHeader(hash, 0)
	root := header.Root.ExtendZero()
	db.DeleteTrieNode(root)

	// Setup again to restore the state trie
	_, hash, err = SetupGenesisBlock(db, nil)
	assert.Nil(t, err)
	assert.Equal(t, params.MainnetGenesisHash, hash)
	ok, _ := db.HasTrieNode(root)
	assert.True(t, ok)
}

// TestCodeInfo tests the genesis accounts
func TestGenesisAccount(t *testing.T) {
	// test account address
	addr := common.HexToAddress("0x0100000000000000000000000000000000000000")

	// accumulate forks from Kairos config
	v := reflect.ValueOf(*params.KairosChainConfig.Copy())
	forks := map[string]*big.Int{"EmptyForkCompatibleBlock": big.NewInt(0)}
	for i := 0; i < v.NumField(); i++ {
		if strings.HasSuffix(v.Type().Field(i).Name, "CompatibleBlock") {
			forks[v.Type().Field(i).Name] = v.Field(i).Interface().(*big.Int)
		}
	}

	tcs := map[string]struct {
		account GenesisAccount
		test    func(*testing.T, params.Rules, params.VmVersion, bool, account.Account)
	}{
		"simple EOA": {
			account: GenesisAccount{
				Balance: big.NewInt(0),
			},
			test: func(t *testing.T, r params.Rules, vmv params.VmVersion, ok bool, acc account.Account) {
				require.False(t, ok)
				require.Equal(t, params.VmVersion0, vmv)
				require.Equal(t, account.ExternallyOwnedAccountType, acc.Type())
			},
		},
		"simple SCA": {
			account: GenesisAccount{
				Balance: big.NewInt(0),
				Code:    hexutil.MustDecode("0x00"),
			},
			test: func(t *testing.T, r params.Rules, vmv params.VmVersion, ok bool, acc account.Account) {
				if !r.IsIstanbul {
					require.True(t, ok)
					require.Equal(t, params.VmVersion0, vmv)
					require.Equal(t, account.SmartContractAccountType, acc.Type())
				} else if r.IsIstanbul {
					require.True(t, ok)
					require.Equal(t, params.VmVersion1, vmv)
					require.Equal(t, account.SmartContractAccountType, acc.Type())
				}
			},
		},
		"account with delegation code": {
			account: GenesisAccount{
				Balance: big.NewInt(0),
				Code:    types.AddressToDelegation(common.HexToAddress("0x000000000000000000000000000000000000000")),
			},
			test: func(t *testing.T, r params.Rules, vmv params.VmVersion, ok bool, acc account.Account) {
				if !r.IsIstanbul {
					require.True(t, ok)
					require.Equal(t, params.VmVersion0, vmv)
					require.Equal(t, account.SmartContractAccountType, acc.Type())
				} else if r.IsIstanbul && !r.IsPrague {
					require.True(t, ok)
					require.Equal(t, params.VmVersion1, vmv)
					require.Equal(t, account.SmartContractAccountType, acc.Type())
				} else {
					require.True(t, ok)
					require.Equal(t, params.VmVersion1, vmv)
					require.Equal(t, account.ExternallyOwnedAccountType, acc.Type())
				}
			},
		},
		"account with empty code but non-empty storage": {
			account: GenesisAccount{
				Balance: big.NewInt(0),
				Storage: map[common.Hash]common.Hash{{1}: {1}},
			},
			test: func(t *testing.T, r params.Rules, vmv params.VmVersion, ok bool, acc account.Account) {
				if !r.IsPrague {
					require.True(t, ok)
					require.Equal(t, params.VmVersion0, vmv)
					require.Equal(t, account.SmartContractAccountType, acc.Type())
				} else {
					require.False(t, ok)
					require.Equal(t, params.VmVersion0, vmv)
					require.Equal(t, account.ExternallyOwnedAccountType, acc.Type())
				}
			},
		},
	}

	// test all cases for each fork
	for tcn, tc := range tcs {
		for fork, n := range forks {
			t.Run(fmt.Sprintf("%s fork %s", tcn, fork), func(t *testing.T) {
				genesis := &Genesis{
					Config: &params.ChainConfig{
						ChainID: new(big.Int).SetUint64(1),
					},
					Alloc: GenesisAlloc{
						addr: tc.account,
					},
				}

				// make config for fork
				for f2, n2 := range forks {
					var n3 *big.Int
					if n == nil {
						n3 = big.NewInt(0)
					} else if n2 != nil && n.Cmp(n2) >= 0 {
						n3 = big.NewInt(0)
					}
					r := reflect.ValueOf(genesis.Config)
					f := reflect.Indirect(r).FieldByName(f2)
					if f.Kind() != reflect.Invalid {
						f.Set(reflect.ValueOf(n3))
					}
				}

				genesis.Config.SetDefaultsForGenesis()
				genesis.Governance = SetGenesisGovernance(genesis)
				InitDeriveSha(genesis.Config)

				db := database.NewMemoryDBManager()
				gblock := genesis.ToBlock(common.Hash{}, db)
				stateDB, _ := state.New(gblock.Root(), state.NewDatabase(db), nil, nil)

				rules := genesis.Config.Rules(new(big.Int).SetUint64(genesis.Number))
				vmVersion, ok := stateDB.GetVmVersion(addr)
				account := stateDB.GetAccount(addr)

				tc.test(t, rules, vmVersion, ok, account)
			})
		}
	}
}

func genCustomGenesisBlock(customChainId uint64) *Genesis {
	genesis := &Genesis{
		Config: &params.ChainConfig{
			ChainID:                 new(big.Int).SetUint64(customChainId),
			IstanbulCompatibleBlock: big.NewInt(2),
			DeriveShaImpl:           types.ImplDeriveShaConcat,
		},
		Alloc: GenesisAlloc{
			common.HexToAddress("0x0100000000000000000000000000000000000000"): {
				Balance: big.NewInt(1), Storage: map[common.Hash]common.Hash{{1}: {1}},
			},
		},
	}
	genesis.Config.SetDefaultsForGenesis()
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config)
	return genesis
}

func writeNilChainConfig(db database.DBManager, hash common.Hash) {
	miscdb := db.GetMiscDB()
	configPrefix := []byte("klay-config-")
	key := append(configPrefix, hash.Bytes()...)
	miscdb.Delete(key)
}
