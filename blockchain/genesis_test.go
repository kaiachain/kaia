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

	"github.com/davecgh/go-spew/spew"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDefaultGenesisBlock tests the genesis block generation functions: DefaultGenesisBlock, DefaultKairosGenesisBlock
func TestDefaultGenesisBlock(t *testing.T) {
	block := genMainnetGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.MainnetGenesisHash {
		t.Errorf("wrong Mainnet genesis hash, got %v, want %v", block.Hash(), params.MainnetGenesisHash)
	}
	block = genKairosGenesisBlock().ToBlock(common.Hash{}, nil)
	if block.Hash() != params.KairosGenesisHash {
		t.Errorf("wrong Kairos genesis hash, got %v, want %v", block.Hash(), params.KairosGenesisHash)
	}
}

// TestHardCodedChainConfigUpdate tests the public network's chainConfig update.
func TestHardCodedChainConfigUpdate(t *testing.T) {
	mainnetGenesisBlock, kairosGenesisBlock := genMainnetGenesisBlock(), genKairosGenesisBlock()
	tests := []struct {
		name             string
		newHFBlock       *big.Int
		originHFBlock    *big.Int
		fn               func(database.DBManager, *big.Int) (*params.ChainConfig, common.Hash, error)
		wantConfig       *params.ChainConfig // expect value of the SetupGenesisBlock's first return value
		wantHash         common.Hash
		wantErr          error
		wantStoredConfig *params.ChainConfig // expect value of the stored config in DB
		resetFn          func(*big.Int)
	}{
		{
			name:       "Mainnet chainConfig update",
			newHFBlock: big.NewInt(3),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				mainnetGenesisBlock.MustCommit(db)
				mainnetGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, mainnetGenesisBlock, params.MainnetNetworkId, false, false)
			},
			wantHash:         params.MainnetGenesisHash,
			wantConfig:       mainnetGenesisBlock.Config,
			wantStoredConfig: mainnetGenesisBlock.Config,
		},
		// TODO-Kaia: add more Mainnet test cases after Mainnet hard fork block numbers are added
		{
			// Because of the fork-ordering check logic, the istanbulCompatibleBlock should be less than the londonCompatibleBlock
			name:       "Kairos chainConfig update - correct hard-fork block number order",
			newHFBlock: big.NewInt(79999999),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				kairosGenesisBlock.MustCommit(db)
				kairosGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, kairosGenesisBlock, params.KairosNetworkId, false, false)
			},
			wantHash:         params.KairosGenesisHash,
			wantConfig:       kairosGenesisBlock.Config,
			wantStoredConfig: kairosGenesisBlock.Config,
		},
		{
			// This test fails because the new istanbulCompatibleBlock(90909999) is larger than londonCompatibleBlock(80295291)
			name:       "Kairos chainConfig update - wrong hard-fork block number order",
			newHFBlock: big.NewInt(90909999),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				kairosGenesisBlock.MustCommit(db)
				kairosGenesisBlock.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, kairosGenesisBlock, params.KairosNetworkId, false, false)
			},
			wantHash:         common.Hash{},
			wantConfig:       kairosGenesisBlock.Config,
			wantStoredConfig: nil,
			wantErr: fmt.Errorf("unsupported fork ordering: %v enabled at %v, but %v enabled at %v",
				"istanbulBlock", big.NewInt(90909999), "londonBlock", big.NewInt(80295291)),
		},
		{
			name:       "incompatible config in DB",
			newHFBlock: big.NewInt(3),
			fn: func(db database.DBManager, newHFBlock *big.Int) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := genMainnetGenesisBlock()
				genesisBlock := genesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, genesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				blocks, _ := GenerateChain(genesis.Config, genesisBlock, faker.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newConfig := *genesis
				newConfig.Config.IstanbulCompatibleBlock = newHFBlock
				return SetupGenesisBlock(db, &newConfig, params.MainnetNetworkId, true, false)
			},
			wantHash:         params.MainnetGenesisHash,
			wantConfig:       mainnetGenesisBlock.Config,
			wantStoredConfig: params.MainnetChainConfig,
			wantErr: &params.ConfigCompatError{
				What:         "Istanbul Block",
				StoredConfig: params.MainnetChainConfig.IstanbulCompatibleBlock,
				NewConfig:    big.NewInt(3),
				RewindTo:     2,
			},
		},
	}

	for _, test := range tests {
		db := database.NewMemoryDBManager()
		config, hash, err := test.fn(db, test.newHFBlock)

		// Check the return values
		assert.Equal(t, test.wantErr, err, test.name+": err is mismatching")
		assert.Equal(t, test.wantConfig, config, test.name+": config is mismatching")
		assert.Equal(t, test.wantHash, hash, test.name+": hash is mismatching")

		// Check stored genesis block
		if test.wantHash != (common.Hash{}) {
			stored := db.ReadBlock(test.wantHash, 0)
			assert.Equal(t, test.wantHash, stored.Hash(), test.name+": stored genesis block is not compatible")
		}

		// Check stored chainConfig
		storedChainConfig, err := db.ReadChainConfig(test.wantHash)
		assert.NoError(t, err)
		assert.Equal(t, test.wantStoredConfig, storedChainConfig, test.name+": stored chainConfig is not compatible")
	}
}

func TestSetupGenesis(t *testing.T) {
	var (
		customGenesisHash = common.HexToHash("0x4eb4035b7a09619a9950c9a4751cc331843f2373ef38263d676b4a132ba4059c")
		customChainId     = uint64(4343)
		customGenesis     = genCustomGenesisBlock(customChainId)
	)
	tests := []struct {
		name       string
		fn         func(database.DBManager) (*params.ChainConfig, common.Hash, error)
		wantConfig *params.ChainConfig
		wantHash   common.Hash
		wantErr    error
	}{
		{
			name: "genesis without ChainConfig",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, new(Genesis), params.UnusedNetworkId, false, false)
			},
			wantErr:    errGenesisNoConfig,
			wantConfig: params.TestChainConfig,
		},
		{
			name: "no block in DB, genesis == nil, Mainnet networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil, params.MainnetNetworkId, false, false)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "no block in DB, genesis == nil, Kairos networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, nil, params.KairosNetworkId, false, false)
			},
			wantHash:   params.KairosGenesisHash,
			wantConfig: params.KairosChainConfig,
		},
		{
			name: "no block in DB, genesis == customGenesis, private network",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "Mainnet block in DB, genesis == nil, Mainnet networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genMainnetGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil, params.MainnetNetworkId, false, false)
			},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "Kairos block in DB, genesis == nil, Kairos networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, nil, params.KairosNetworkId, false, false)
			},
			wantHash:   params.KairosGenesisHash,
			wantConfig: params.KairosChainConfig,
		},
		{
			name: "custom block in DB, genesis == nil, custom networkId",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, nil, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "Mainnet block in DB, genesis == Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genMainnetGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genKairosGenesisBlock(), params.KairosNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: params.KairosGenesisHash},
			wantHash:   params.KairosGenesisHash,
			wantConfig: params.KairosChainConfig,
		},
		{
			name: "Kairos block in DB, genesis == Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genMainnetGenesisBlock(), params.MainnetNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: params.MainnetGenesisHash},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "Mainnet block in DB, genesis == custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genMainnetGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, genCustomGenesisBlock(customChainId), customChainId, true, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.MainnetGenesisHash, New: customGenesisHash},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "Kairos block in DB, genesis == custom",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				genKairosGenesisBlock().MustCommit(db)
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantErr:    &GenesisMismatchError{Stored: params.KairosGenesisHash, New: customGenesisHash},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "custom block in DB, genesis == Mainnet",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, genMainnetGenesisBlock(), params.MainnetNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: customGenesisHash, New: params.MainnetGenesisHash},
			wantHash:   params.MainnetGenesisHash,
			wantConfig: params.MainnetChainConfig,
		},
		{
			name: "custom block in DB, genesis == Kairos",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, genKairosGenesisBlock(), params.KairosNetworkId, false, false)
			},
			wantErr:    &GenesisMismatchError{Stored: customGenesisHash, New: params.KairosGenesisHash},
			wantHash:   params.KairosGenesisHash,
			wantConfig: params.KairosChainConfig,
		},
		{
			name: "compatible config in DB",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				customGenesis.MustCommit(db)
				return SetupGenesisBlock(db, customGenesis, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
		},
		{
			name: "incompatible config in DB",
			fn: func(db database.DBManager) (*params.ChainConfig, common.Hash, error) {
				// Commit the 'old' genesis block with Istanbul transition at #2.
				// Advance to block #4, past the Istanbul transition block of customGenesis.
				genesis := customGenesis.MustCommit(db)

				bc, _ := NewBlockChain(db, nil, customGenesis.Config, faker.NewFullFaker(), vm.Config{})
				defer bc.Stop()

				blocks, _ := GenerateChain(customGenesis.Config, genesis, faker.NewFaker(), db, 4, nil)
				bc.InsertChain(blocks)
				// This should return a compatibility error.
				newConfig := *customGenesis
				newConfig.Config.IstanbulCompatibleBlock = big.NewInt(3)
				return SetupGenesisBlock(db, &newConfig, customChainId, true, false)
			},
			wantHash:   customGenesisHash,
			wantConfig: customGenesis.Config,
			wantErr: &params.ConfigCompatError{
				What:         "Istanbul Block",
				StoredConfig: big.NewInt(2),
				NewConfig:    big.NewInt(3),
				RewindTo:     1,
			},
		},
	}

	for _, test := range tests {
		db := database.NewMemoryDBManager()
		config, hash, err := test.fn(db)
		// Check the return values.
		if !reflect.DeepEqual(err, test.wantErr) {
			spew := spew.ConfigState{DisablePointerAddresses: true, DisableCapacities: true}
			t.Errorf("%s: returned error %#v, want %#v", test.name, spew.NewFormatter(err), spew.NewFormatter(test.wantErr))
		}
		if !reflect.DeepEqual(config, test.wantConfig) {
			t.Errorf("%s:\nreturned %v\nwant     %v", test.name, config, test.wantConfig)
		}
		if hash != test.wantHash {
			t.Errorf("%s: returned hash %s, want %s", test.name, hash.Hex(), test.wantHash.Hex())
		} else if err == nil {
			// Check database content.
			stored := db.ReadBlock(test.wantHash, 0)
			if stored.Hash() != test.wantHash {
				t.Errorf("%s: block in DB has hash %s, want %s", test.name, stored.Hash(), test.wantHash)
			}
		}
	}
}

// Test that SetupGenesisBlock writes the genesis state trie if not present.
func TestGenesisRestoreState(t *testing.T) {
	db := database.NewMemoryDBManager()

	// Setup first to Commit the Genesis block
	_, hash, err := SetupGenesisBlock(db, nil, params.MainnetNetworkId, false, false)
	assert.Nil(t, err)
	assert.Equal(t, params.MainnetGenesisHash, hash)

	// Simulate state migration by deleting the state trie root node
	header := db.ReadHeader(hash, 0)
	root := header.Root.ExtendZero()
	db.DeleteTrieNode(root)

	// Setup again to restore the state trie
	_, hash, err = SetupGenesisBlock(db, nil, params.MainnetNetworkId, false, false)
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

func genMainnetGenesisBlock() *Genesis {
	genesis := DefaultGenesisBlock()
	genesis.Config = params.MainnetChainConfig.Copy()
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config)
	return genesis
}

func genKairosGenesisBlock() *Genesis {
	genesis := DefaultKairosGenesisBlock()
	genesis.Config = params.KairosChainConfig.Copy()
	genesis.Governance = SetGenesisGovernance(genesis)
	InitDeriveSha(genesis.Config)
	return genesis
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
