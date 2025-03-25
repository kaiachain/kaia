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
)

var (
	ErrInitUnexpectedNil         = errors.New("unexpected nil during module init")
	ErrGSRNotInstalled           = errors.New("Gasless swap router contract not installed")
	ErrDecodeSwapTx              = errors.New("failed to decode swap transaction")
	ErrSwapTxInvalid             = errors.New("swap transaction is not valid")
	ErrDecodeApproveTx           = errors.New("failed to decode approve transaction")
	ErrApproveTxInvalid          = errors.New("approve transaction is not valid")
	ErrDifferentSenders          = errors.New("approve and swap transactions have different senders")
	ErrDifferentTokens           = errors.New("approve transaction is for different token than swap transaction")
	ErrInsufficientApproveAmount = errors.New("approve transaction approves insufficient amount")
	ErrNonSequentialNonce        = errors.New("approve and swap transactions have non-sequential nonces")
	ErrApproveNonceNotCurrent    = errors.New("approve transaction nonce is not current")
	ErrSwapNonceNotCurrent       = errors.New("swap transaction nonce is not current")
	ErrIncorrectRepayAmount      = errors.New("swap transaction has incorrect amountRepay")
)

func printApproveTx(args *ApproveArgs) string {
	return fmt.Sprintf("ApproveTx{Sender: %s, Token: %s, Spender: %s, Amount: %s}",
		args.Sender.Hex(), args.Token.Hex(), args.Spender.Hex(), args.Amount.String())
}

func printSwapTx(args *SwapArgs) string {
	return fmt.Sprintf("SwapTx{Sender: %s, Router: %s, Token: %s, AmountIn: %s, MinAmountOut: %s, AmountRepay: %s}",
		args.Sender.Hex(), args.Router.Hex(), args.Token.Hex(), args.AmountIn.String(), args.MinAmountOut.String(), args.AmountRepay.String())
}
