// Copyright 2025 The Kaia Authors
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

package impl

import (
	"errors"
	"fmt"
	"math"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip247"
	"github.com/kaiachain/kaia/contracts/contracts/testing/sc_erc20"
	"github.com/kaiachain/kaia/kaiax"
)

var _ kaiax.TxPoolModule = (*GaslessModule)(nil)

func (g *GaslessModule) PreAddTx(tx *types.Transaction, local bool) error {
	return nil
}

func (g *GaslessModule) IsModuleTx(tx *types.Transaction) bool {
	if tx == nil {
		return false
	}
	return g.IsApproveTx(tx) || g.IsSwapTx(tx)
}

func (g *GaslessModule) GetCheckBalance() func(tx *types.Transaction) error {
	return func(tx *types.Transaction) error {
		if approveArgs, ok := decodeApproveTx(tx, g.signer); ok {
			return g.checkBalanceForApprove(approveArgs)
		}
		if swapArgs, ok := decodeSwapTx(tx, g.signer); ok {
			return g.checkBalanceForSwap(swapArgs, tx.Nonce())
		}
		return errors.New("not a gasless transaction") // should not happen because IsModuleTx is called before GetCheckBalance
	}
}

func (g *GaslessModule) checkBalanceForApprove(approveArgs *ApproveArgs) error {
	token := approveArgs.Token
	bc := backends.NewBlockchainContractBackend(g.Chain, nil, nil)

	tokenContract, err := sc_erc20.NewERC20(token, bc)
	if err != nil {
		return err
	}

	// tx.token.balanceOf(sender) > 0
	tokenBalance, err := tokenContract.BalanceOf(nil, approveArgs.Sender)
	if err != nil {
		return err
	}
	if tokenBalance.Sign() <= 0 {
		return fmt.Errorf("insufficient sender token balance: token=%s, have=%s, want=nonzero", token.Hex(), tokenBalance.String())
	}

	return nil
}

// tx.minAmountOut >= tx.amountRepay
// tx.amountIn >= gsr.getAmountIn(minAmountOut)
// tx.token.approval(sender, router) >= tx.amountIn
// tx.token.balanceOf(sender) >= tx.amountIn
// tx.deadline >= currentTimestamp
func (g *GaslessModule) checkBalanceForSwap(swapArgs *SwapArgs, swapNonce uint64) error {
	// If SwapTx.nonce is the sender's next nonce, then there is no room for ApproveTx proceeding SwapTx.
	senderNonce := g.getCurrentStateNonce(swapArgs.Sender)
	noApproveTxPreceeds := swapNonce == senderNonce

	token := swapArgs.Token
	bc := backends.NewBlockchainContractBackend(g.Chain, nil, nil)

	tokenContract, err := sc_erc20.NewERC20(token, bc)
	if err != nil {
		return err
	}

	// tx.minAmountOut >= tx.amountRepay
	minAmountOut := swapArgs.MinAmountOut
	amountRepay := swapArgs.AmountRepay
	if minAmountOut.Cmp(amountRepay) < 0 {
		return fmt.Errorf("insufficient minAmountOut: minAmountOut=%s, amountRepay=%s", minAmountOut.String(), amountRepay.String())
	}

	// tx.amountIn >= gsr.getAmountIn(minAmountOut)
	routerContract, err := kip247.NewGaslessSwapRouterCaller(g.swapRouter, bc)
	if err != nil {
		return err
	}
	// Required token amountIn, given the current exchange rate and the declared minAmountOut.
	requiredAmountIn, err := routerContract.GetAmountIn(nil, token, minAmountOut)
	if err != nil {
		return err
	}
	if swapArgs.AmountIn.Cmp(requiredAmountIn) < 0 {
		return fmt.Errorf("insufficient amountIn: have=%s, want=%s", swapArgs.AmountIn.String(), requiredAmountIn.String())
	}

	if noApproveTxPreceeds {
		// tx.token.allowance(sender, router) >= tx.amountIn
		approval, err := tokenContract.Allowance(nil, swapArgs.Sender, g.swapRouter)
		if err != nil {
			return err
		}
		if approval.Cmp(swapArgs.AmountIn) < 0 {
			return fmt.Errorf("insufficient approval: approval=%s, want=%s", approval.String(), swapArgs.AmountIn.String())
		}
	}

	// tx.token.balanceOf(sender) >= tx.amountIn
	balance, err := tokenContract.BalanceOf(nil, swapArgs.Sender)
	if err != nil {
		return err
	}
	if balance.Cmp(swapArgs.AmountIn) < 0 {
		return fmt.Errorf("insufficient balance: balance=%s, want=%s", balance.String(), swapArgs.AmountIn.String())
	}

	// tx.deadline >= currentTimestamp
	deadline := swapArgs.Deadline
	if deadline.Cmp(g.Chain.CurrentBlock().Time()) < 0 {
		return fmt.Errorf("insufficient deadline: deadline=%s, want=%s", deadline.String(), g.Chain.CurrentBlock().Time().String())
	}

	return nil
}

func (g *GaslessModule) IsReady(txs map[uint64]*types.Transaction, i uint64, ready types.Transactions) bool {
	tx := txs[i]

	if g.IsApproveTx(tx) && i < uint64(math.MaxUint64) {
		return g.isApproveTxReady(tx, txs[i+1])
	}

	if g.IsSwapTx(tx) {
		var prevTx *types.Transaction
		if len(ready) > 0 {
			prevTx = ready[len(ready)-1]
		}
		return g.isSwapTxReady(tx, prevTx)
	}

	return false
}

func (g *GaslessModule) PreReset(oldHead, newHead *types.Header) []common.Hash {
	return nil
}

func (g *GaslessModule) PostReset(oldHead, newHead *types.Header, queue, pending map[common.Address]types.Transactions) {
}

// isApproveTxReady assumes that the caller checked `g.IsApproveTx(approveTx)`
func (g *GaslessModule) isApproveTxReady(approveTx, nextTx *types.Transaction) bool {
	addr, err := types.Sender(g.signer, approveTx)
	if err != nil {
		return false
	}
	nonce := g.getCurrentStateNonce(addr)

	if approveTx.Nonce() != nonce {
		return false
	}
	if nextTx == nil || !g.IsSwapTx(nextTx) {
		return false
	}

	return g.IsExecutable(approveTx, nextTx)
}

// isSwapTxReady assumes that the caller checked `g.IsSwapTx(swapTx)`
func (g *GaslessModule) isSwapTxReady(swapTx, prevTx *types.Transaction) bool {
	addr, err := types.Sender(g.signer, swapTx)
	if err != nil {
		return false
	}
	nonce := g.getCurrentStateNonce(addr)

	var approveTx *types.Transaction
	if swapTx.Nonce() == nonce {
		approveTx = nil
	} else if swapTx.Nonce() == nonce+1 {
		if prevTx == nil || !g.IsApproveTx(prevTx) {
			return false
		}
		approveTx = prevTx
	} else {
		return false
	}

	return g.IsExecutable(approveTx, swapTx)
}
