// Copyright 2025 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package state

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This test mimics the behavior of a block miner with live pruning enabled.
func TestRollback(t *testing.T) {
	dbm := database.NewMemoryDBManager()
	dbm.WritePruningEnabled()

	var (
		sdb  = NewDatabase(dbm)
		acc1 = common.HexToAddress("0x0000000000000000000000000000000000000aaa")
		acc2 = common.HexToAddress("0x0000000000000000000000000000000000000bbb")
		acc3 = common.HexToAddress("0x0000000000000000000000000000000000000ccc")

		root1 common.Hash
		root2 common.Hash
	)

	{ // Block 1: Store baseline state.
		state, _ := New(common.Hash{}, sdb, nil, nil)
		state.AddBalance(acc1, big.NewInt(10))
		state.AddBalance(acc2, big.NewInt(20))
		state.AddBalance(acc3, big.NewInt(30))
		root1, _ = state.Commit(true)
		t.Logf("end block 1, root %x", root1)
	}

	{ // Block 2: Build a bundle-containing block.
		// worker.makeCurrent: PrunableStateAt(parentHash, parentNum)
		state, err := New(root1, sdb, nil, &statedb.TrieOpts{PruningBlockNumber: 1})
		assert.NoError(t, err)

		{ // Tx 0: Regular transaction that succeeds
			txSnap := state.Snapshot()              // EVM.Call: StateDB.Snapshot()
			state.AddBalance(acc1, big.NewInt(100)) // bc.ApplyTransaction: ApplyMessage()
			_ = txSnap                              // EVM.Call: no RevertToSnapshot()
			state.Finalise(true, false)             // bc.ApplyTransaction: Finalise(true, false)
		}

		{ // Tx 1: Regular transaction that fails
			txSnap := state.Snapshot()              // EVM.Call: StateDB.Snapshot()
			state.AddBalance(acc1, big.NewInt(999)) // bc.ApplyTransaction: ApplyMessage()
			state.RevertToSnapshot(txSnap)          // EVM.Call: RevertToSnapshot()
			state.Finalise(true, false)             // bc.ApplyTransaction: Finalise(true, false)
		}

		{ // Execute bundle transactions.
			// worker.commitBundleTransaction
			snapshot := state.Copy()
			state.StartPruningSnapshot()

			{ // Tx 2: Bundle transaction that succeeds
				txSnap := state.Snapshot()              // EVM.Call: StateDB.Snapshot()
				state.AddBalance(acc2, big.NewInt(100)) // bc.ApplyTransaction: ApplyMessage()
				_ = txSnap                              // EVM.Call: no RevertToSnapshot()
				state.Finalise(true, false)             // bc.ApplyTransaction: Finalise(true, false)
			}

			{ // Tx 3: Bundle transaction that fails
				txSnap := state.Snapshot()              // EVM.Call: StateDB.Snapshot()
				state.AddBalance(acc3, big.NewInt(100)) // bc.ApplyTransaction: ApplyMessage()
				state.RevertToSnapshot(txSnap)          // EVM.Call: RevertToSnapshot()
				state.Finalise(true, false)             // bc.ApplyTransaction: Finalise(true, false)
			}

			// worker.commitBundleTransaction: restoreEnv
			state.RevertPruningSnapshot()
			state.Set(snapshot)
		}

		// At the end of the day, only Tx 0 is applied.
		assert.Equal(t, uint64(110), state.GetBalance(acc1).Uint64())
		assert.Equal(t, uint64(20), state.GetBalance(acc2).Uint64())
		assert.Equal(t, uint64(30), state.GetBalance(acc3).Uint64())

		// Finalize the block.
		root2 = state.IntermediateRoot(true) // worker.commitNewWork: engine.Finalize
		assert.NoError(t, err)

		// After consensus, commit the block.
		// worker.wait: bc.WriteBlockWithState: bc.writeStateTrie
		root22, err := state.Commit(true)
		assert.NoError(t, err)
		assert.Equal(t, root2, root22)

		err = sdb.TrieDB().Commit(root2, true, 2)
		assert.NoError(t, err)

		t.Logf("end block 2, root %s", root2.Hex())
	}

	{ // Simulate the passage of time, in that
		// - db.pruningMarks are eventually written to the diskDB.
		sdb.TrieDB().Cap(0)

		// - in-memory trie cache is evicted (governed by --state.cache-size).
		sdb = NewDatabase(dbm)

		// - bc.pruneTrieNodeLoop() deleted (after retention) as dictated by the pruning marks.
		marks := dbm.ReadPruningMarks(0, 99)
		for _, mark := range marks {
			t.Logf("delete trie node (%s, %d)", mark.Hash.Hex(), mark.Number)
			// NOTE: if you keep the root node (e.g. if it survives in in-memory trie cache),
			// you can observe that GetBalance() below returns 0.
			// if mark.Hash.IsZeroExtended() {
			//   continue
			// }
			dbm.DeleteTrieNode(mark.Hash)
		}
	}

	{ // After that, some trie nodes that represent the latest state, must be intact.
		// i.e. the states must not be pruned.
		t.Log("query block 2")
		state, err := New(root2, sdb, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, uint64(110), state.GetBalance(acc1).Uint64(), "acc1")
		assert.Equal(t, uint64(20), state.GetBalance(acc2).Uint64(), "acc2")
		assert.Equal(t, uint64(30), state.GetBalance(acc3).Uint64(), "acc3")
		assert.NoError(t, state.Error()) // No error during GetBalance() calls.
	}
}
