// Copyright 2024 The Kaia Authors
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

package gasless

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/gasless GaslessModule
type GaslessModule interface {
	kaiax.BaseModule

	// IsApproveTx checks if the transaction is a GaslessApproveTx, i.e. a transaction that approves an whitelisted ERC20 token to the SwapRouter contract.
	// An ApproveTx can be inserted to txpool.queue even if sender's balance is insufficient.
	IsApproveTx(tx *types.Transaction) bool

	// IsSwapTx checks if the transaction is a GaslessSwapTx, i.e. a transaction that invokes the SwapRouter contract.
	// An SwapTx can be inserted to txpool.queue even if sender's balance is insufficient.
	IsSwapTx(tx *types.Transaction) bool

	// IsExecutable checks if the given approve and swap transactions are ready to be executed.
	// A (ApproveTx, SwapTx) pair or a (SwapTx) is executable even if sender's balance is insufficient.
	IsExecutable(approveTxOrNil, swapTx *types.Transaction) bool

	// GetMakeLendTxFunc returns a function that creates a signed lend transaction that can fund the given approve and swap transactions.
	// TODO: define makeTxFunc type elsewhere (e.g. bundling module)
	GetMakeLendTxFunc(approveTxOrNil, swapTx *types.Transaction) func(nonce uint64) (*types.Transaction, error)
}
