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

package blockchain

import (
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func TestBlockChainTxLookup(t *testing.T) {
	var bc BlockChain
	if bc.TxLookup() != nil {
		t.Fatal("unexpected tx lookup on zero BlockChain")
	}

	called := false
	bc.SetTxLookup(func(common.Hash) *types.Transaction {
		called = true
		return nil
	})

	lookup := bc.TxLookup()
	if lookup == nil {
		t.Fatal("expected tx lookup to be registered")
	}
	lookup(common.Hash{})
	if !called {
		t.Fatal("registered tx lookup was not invoked")
	}

	bc.SetTxLookup(nil)
	if bc.TxLookup() == nil {
		t.Fatal("nil tx lookup registration should be ignored")
	}
}
