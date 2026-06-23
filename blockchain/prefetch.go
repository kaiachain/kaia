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
	"context"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

type stateAtReader interface {
	StateAt(root common.Hash) (*state.StateDB, error)
}

// PrefetchBlockState warms the trieDB node cache via a disposable StateDB.
func PrefetchBlockState(ctx context.Context, stateReader stateAtReader, parentRoot common.Hash, blockNumber uint64, txs []*types.Transaction, signer types.Signer) {
	if stateReader == nil || len(txs) == 0 {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			logger.Warn("PrefetchBlockState recovered from panic", "err", r)
		}
	}()
	statedb, err := stateReader.StateAt(parentRoot)
	if err != nil {
		return
	}
	for _, tx := range txs {
		select {
		case <-ctx.Done():
			return
		default:
		}
		prefetchTxState(statedb, signer, tx, blockNumber)
	}
}

// prefetchTxState warms the trie nodes execution will need for one tx
// without running the EVM: sender + (fee-payer) account keys via
// ValidateSender, and recipient bytecode/storage-trie root.
func prefetchTxState(statedb *state.StateDB, signer types.Signer, tx *types.Transaction, blockNumber uint64) {
	if _, err := tx.ValidateSender(signer, statedb, blockNumber); err != nil {
		return
	}
	from := tx.ValidatedSender()
	statedb.Exist(from)
	if tx.IsFeeDelegatedTransaction() {
		if _, err := tx.ValidateFeePayer(signer, statedb, blockNumber); err == nil {
			statedb.Exist(tx.ValidatedFeePayer())
		}
	}
	// EIP-7702: sender may carry delegated code; warm its code + storage too.
	if statedb.GetCodeHash(from) != types.EmptyCodeHash {
		statedb.GetCode(from)
		statedb.GetState(from, common.Hash{})
	}
	to := tx.To()
	if to == nil {
		return
	}
	if statedb.GetCodeHash(*to) != types.EmptyCodeHash {
		statedb.GetCode(*to)
		statedb.GetState(*to, common.Hash{})
	}
}
