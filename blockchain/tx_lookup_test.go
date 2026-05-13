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
