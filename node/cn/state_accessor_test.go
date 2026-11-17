// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.
package cn

import (
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
)

// TestStateAtTransactionReleasesOnError checks that stateAtTransaction drops its
// reference to the parent state when it fails, so the state can still be
// garbage-collected.
func TestStateAtTransactionReleasesOnError(t *testing.T) {
	key, _ := crypto.GenerateKey()
	config := params.TestChainConfig
	signer := types.LatestSignerForChainID(config.ChainID)
	gspec := &blockchain.Genesis{Config: config, Alloc: blockchain.GenesisAlloc{
		crypto.PubkeyToAddress(key.PublicKey): {Balance: big.NewInt(params.KAIA)},
	}}
	newTx := func(nonce uint64, signer types.Signer) *types.Transaction {
		tx, err := types.SignTx(types.NewTransaction(nonce, common.HexToAddress("0x1234"), big.NewInt(1), params.TxGas, big.NewInt(25*params.Gkei), nil), signer, key)
		assert.NoError(t, err)
		return tx
	}

	setup := func() (*CN, *blockchain.BlockChain) {
		db, gendb := database.NewMemoryDBManager(), database.NewMemoryDBManager()
		genesis := gspec.MustCommit(db)
		gspec.MustCommit(gendb)
		blocks, _ := blockchain.GenerateChain(config, genesis, faker.NewFaker(), gendb, 2, func(i int, b *blockchain.BlockGen) {
			b.AddTx(newTx(uint64(i), signer))
		})
		chain, err := blockchain.NewBlockChain(db, &blockchain.CacheConfig{
			CacheSize:           512,
			BlockInterval:       blockchain.DefaultBlockInterval,
			TriesInMemory:       blockchain.DefaultTriesInMemory,
			TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
			SnapshotCacheSize:   512,
		}, config, faker.NewFaker(), vm.Config{})
		assert.NoError(t, err)
		_, err = chain.InsertChain(blocks)
		assert.NoError(t, err)
		return &CN{blockchain: chain, chainConfig: config, config: &Config{}}, chain
	}
	inMemory := func(chain *blockchain.BlockChain, root common.Hash) bool {
		for _, h := range chain.StateCache().TrieDB().Nodes() {
			if h == root.ExtendZero() {
				return true
			}
		}
		return false
	}
	// assertReleased checks that stateAtTransaction dropped exactly its own reference to root.
	assertReleased := func(t *testing.T, chain *blockchain.BlockChain, root common.Hash) {
		assert.True(t, inMemory(chain, root), "the chain's own reference was dropped too")
		// Drop the chain's own reference, as its garbage collection would.
		chain.StateCache().TrieDB().Dereference(root)
		assert.False(t, inMemory(chain, root), "the state is still referenced")
	}

	t.Run("tx index out of range", func(t *testing.T) {
		cn, chain := setup()
		head := chain.CurrentBlock()
		_, _, _, _, _, err := cn.stateAtTransaction(head, 99, 0, nil, true, false)
		assert.Error(t, err)
		assertReleased(t, chain, chain.GetBlockByHash(head.ParentHash()).Root())
	})
	t.Run("tx with a bad signature", func(t *testing.T) {
		cn, chain := setup()
		head := chain.CurrentBlock()
		header := types.CopyHeader(head.Header())
		header.ParentHash = head.Hash()
		header.Number = new(big.Int).Add(head.Number(), common.Big1)
		block := types.NewBlock(header, types.Transactions{newTx(2, types.LatestSignerForChainID(big.NewInt(9999)))}, nil)
		_, _, _, _, _, err := cn.stateAtTransaction(block, 0, 0, nil, true, false)
		assert.Error(t, err)
		assertReleased(t, chain, head.Root())
	})
	t.Run("InitializeState failure", func(t *testing.T) {
		cn, chain := setup()
		chain.Processor().(*blockchain.StateProcessor).RegisterBlockStateModule(failingInitModule{})
		head := chain.CurrentBlock()
		_, _, _, _, _, err := cn.stateAtTransaction(head, 0, 0, nil, true, false)
		assert.Error(t, err)
		assertReleased(t, chain, chain.GetBlockByHash(head.ParentHash()).Root())
	})
}

// failingInitModule is a block state module whose InitializeState always fails.
type failingInitModule struct{}

func (failingInitModule) InitializeState(*types.Header, *state.StateDB) error {
	return errors.New("initialize state failed")
}

func (failingInitModule) FinalizeState(*types.Header, *state.StateDB, []*types.Transaction, []*types.Receipt) error {
	return nil
}
