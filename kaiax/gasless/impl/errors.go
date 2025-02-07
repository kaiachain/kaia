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

package impl

import (
	"errors"
	"fmt"
)

var ErrInitUnexpectedNil = errors.New("unexpected nil during module init")

func printApproveTx(args *ApproveArgs) string {
	return fmt.Sprintf("ApproveTx{Sender: %s, Token: %s, Spender: %s, Amount: %s}",
		args.Sender.Hex(), args.Token.Hex(), args.Spender.Hex(), args.Amount.String())
}

func printSwapTx(args *SwapArgs) string {
	return fmt.Sprintf("SwapTx{Sender: %s, Router: %s, Token: %s, AmountIn: %s, MinAmountOut: %s, AmountRepay: %s}",
		args.Sender.Hex(), args.Router.Hex(), args.Token.Hex(), args.AmountIn.String(), args.MinAmountOut.String(), args.AmountRepay.String())
}
