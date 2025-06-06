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
)

// TestRollBack tests whether pruningMarks are recorded correctly during bundle execution
// 1. make statedb copy thorugh state.Copy()
// 2. let pruning module mark candidates
// 3. check copied if statedb is not reflected
func TestRollback(t *testing.T) {
	// worker.go:commitBundleTransactiond
	dbm := database.NewMemoryDBManager()
	opts := &statedb.TrieOpts{}
	dbm.WritePruningEnabled()
	opts.PruningBlockNumber = 1

	sdb := NewDatabase(dbm)
	sdb.TrieDB().SavePruningMarksInBundle()
	state, _ := New(common.Hash{}, sdb, nil, opts)

	var (
		acc  = common.HexToAddress("0x0000000000000000000000000000000000000aaa")
		acc2 = common.HexToAddress("0x0000000000000000000000000000000000000bbb")
		acc3 = common.HexToAddress("0x0000000000000000000000000000000000000ccc")
	)
	// set retention
	// insert block required

	// ACC == 10
	state.AddBalance(acc2, big.NewInt(30))
	state.AddBalance(acc3, big.NewInt(40))
	state.AddBalance(acc, big.NewInt(10))
	root10, _ := state.Commit(true)
	t.Log("root", root10.Hex())

	state, _ = New(root10, sdb, nil, opts)
	snapshot := state.Copy()

	// ACC == 200
	state.SetBalance(acc, big.NewInt(200))

	root200, _ := state.Commit(true)
	t.Log("root", root200.Hex())
	state.Database().TrieDB().Cap(0)

	// ACC == 10
	sdb.TrieDB().RevertPruningMarksInBundle()
	state.Set(snapshot)

	marks := dbm.ReadPruningMarks(0, 2)
	for _, mark := range marks {
		t.Log("delete", mark.Hash.Hex())
		dbm.DeleteTrieNode(mark.Hash)
	}

	// Already opened
	t.Log("balance", state.GetBalance(acc))

	// Newly opening
	state, err := New(root10, sdb, nil, opts)
	if err != nil {
		t.Log("err", err)
		t.Fail()
	}
	if state != nil {
		t.Log("balance", state.GetBalance(acc))
	}
}
