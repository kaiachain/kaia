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
	"bytes"
	"fmt"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/params"
)

const (
	// import { erc20Abi } from 'viem';
	erc20AbiJson = `[{"type":"event","name":"Approval","inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}]},{"type":"event","name":"Transfer","inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}]},{"type":"function","name":"allowance","stateMutability":"view","inputs":[{"name":"owner","type":"address"},{"name":"spender","type":"address"}],"outputs":[{"type":"uint256"}]},{"type":"function","name":"approve","stateMutability":"nonpayable","inputs":[{"name":"spender","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]},{"type":"function","name":"balanceOf","stateMutability":"view","inputs":[{"name":"account","type":"address"}],"outputs":[{"type":"uint256"}]},{"type":"function","name":"decimals","stateMutability":"view","inputs":[],"outputs":[{"type":"uint8"}]},{"type":"function","name":"name","stateMutability":"view","inputs":[],"outputs":[{"type":"string"}]},{"type":"function","name":"symbol","stateMutability":"view","inputs":[],"outputs":[{"type":"string"}]},{"type":"function","name":"totalSupply","stateMutability":"view","inputs":[],"outputs":[{"type":"uint256"}]},{"type":"function","name":"transfer","stateMutability":"nonpayable","inputs":[{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]},{"type":"function","name":"transferFrom","stateMutability":"nonpayable","inputs":[{"name":"sender","type":"address"},{"name":"recipient","type":"address"},{"name":"amount","type":"uint256"}],"outputs":[{"type":"bool"}]}]`
	// function swapForGas(address token, uint256 amountIn, uint256 minAmountOut, uint256 amountRepay) external
	routerAbiJson = `[{"inputs":[{"internalType":"address","name":"token","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint256","name":"minAmountOut","type":"uint256"},{"internalType":"uint256","name":"amountRepay","type":"uint256"}],"name":"swapForGas","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

var (
	erc20ApproveFunc = mustParseAbi(erc20AbiJson, "approve")
	routerSwapFunc   = mustParseAbi(routerAbiJson, "swapForGas")
)

type ApproveArgs struct {
	Sender  common.Address // tx.from
	Token   common.Address // tx.to
	Spender common.Address
	Amount  *big.Int
}

type SwapArgs struct {
	Sender       common.Address // tx.from
	Router       common.Address // tx.to
	Token        common.Address
	AmountIn     *big.Int
	MinAmountOut *big.Int
	AmountRepay  *big.Int
}

// IsApproveTx checks following conditions:
// A1. tx.to is a whitelisted ERC20 token.
// A2. tx.data is `approve(spender, amount)`.
// A3. spender is a whitelisted SwapRouter contract.
// A4. amount is nonzero.
func (g *GaslessModule) IsApproveTx(tx *types.Transaction) bool {
	args, ok := decodeApproveTx(tx)
	return ok && g.isApproveTx(args)
}

func (g *GaslessModule) isApproveTx(args *ApproveArgs) bool {
	return g.allowedTokens[args.Token] && // A1
		g.swapRouters[args.Spender] && // A3
		args.Amount.Sign() > 0 // A4
}

// IsSwapTx checks following conditions:
// S1. tx.to is a whitelisted SwapRouter contract.
// S2. tx.data is `swapForGas(token, amountIn, minAmountOut, amountRepay)`.
// S3. token is a whitelisted ERC20 token.
func (g *GaslessModule) IsSwapTx(tx *types.Transaction) bool {
	args, ok := decodeSwapTx(tx)
	return ok && g.isSwapTx(args)
}

func (g *GaslessModule) isSwapTx(args *SwapArgs) bool {
	return g.swapRouters[args.Router] && // S1
		g.allowedTokens[args.Token] // S3
}

func mustParseAbi(abiJson string, funcName string) abi.Method {
	abi, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		panic(fmt.Errorf("failed to parse abi: %w", err))
	}
	method, ok := abi.Methods[funcName]
	if !ok {
		panic(fmt.Errorf("method %s not found", funcName))
	}
	return method
}

func decodeApproveTx(tx *types.Transaction) (*ApproveArgs, bool) {
	to, inputs, ok := decodeFunctionCall(tx, erc20ApproveFunc)
	if !ok {
		return nil, false
	}
	spender, ok := inputs["spender"].(common.Address)
	if !ok {
		return nil, false
	}
	amount, ok := inputs["amount"].(*big.Int)
	if !ok {
		return nil, false
	}
	return &ApproveArgs{
		Sender:  tx.ValidatedSender(),
		Token:   to,
		Spender: spender,
		Amount:  amount,
	}, true
}

func decodeSwapTx(tx *types.Transaction) (args *SwapArgs, ok bool) {
	to, inputs, ok := decodeFunctionCall(tx, routerSwapFunc)
	if !ok {
		return nil, false
	}
	token, ok := inputs["token"].(common.Address)
	if !ok {
		return nil, false
	}
	amountIn, ok := inputs["amountIn"].(*big.Int)
	if !ok {
		return nil, false
	}
	minAmountOut, ok := inputs["minAmountOut"].(*big.Int)
	if !ok {
		return nil, false
	}
	amountRepay, ok := inputs["amountRepay"].(*big.Int)
	if !ok {
		return nil, false
	}
	return &SwapArgs{
		Sender:       tx.ValidatedSender(),
		Router:       to,
		Token:        token,
		AmountIn:     amountIn,
		MinAmountOut: minAmountOut,
		AmountRepay:  amountRepay,
	}, true
}

func decodeFunctionCall(tx *types.Transaction, method abi.Method) (common.Address, map[string]interface{}, bool) {
	if tx.Type() != types.TxTypeLegacyTransaction || // not legacy tx: unable to statically determine the max gas fee.
		tx.To() == nil || // not a contract call.
		len(tx.Data()) < 4 || // too short to be a contract call.
		!bytes.Equal(tx.Data()[:4], method.ID) { // not the target function.
		return common.Address{}, nil, false
	}

	inputs := make(map[string]interface{})
	err := method.Inputs.UnpackIntoMap(inputs, tx.Data()[4:])
	return *tx.To(), inputs, err == nil
}

// IsGaslessPattern checks following conditions:
// Ax. IsApproveTx conditions (if ApproveTx != nil)
// Sx. IsSwapTx conditions
// AP1. ApproveTx.from == SwapTx.from
// SP1. ApproveTx.to == SwapTx.token
// SP2. ApproveTx.amount >= SwapTx.amountIn
// SP3. SwapTx.amountRepay = RepayAmount(ApproveTx, SwapTx)
func (g *GaslessModule) IsExecutable(approveTxOrNil, swapTx *types.Transaction) bool {
	// Sx.
	swapArgs, ok := decodeSwapTx(swapTx)
	if !ok || !g.isSwapTx(swapArgs) {
		return false
	}

	// Conditions involving ApproveTx
	if approveTxOrNil != nil {
		// Ax.
		approveArgs, ok := decodeApproveTx(approveTxOrNil)
		if !ok || !g.isApproveTx(approveArgs) {
			return false
		}
		// AP1.
		if approveArgs.Sender != swapTx.ValidatedSender() {
			return false
		}
		// SP1.
		if approveArgs.Token != swapArgs.Token {
			return false
		}
		// SP2.
		if approveArgs.Amount.Cmp(swapArgs.AmountIn) < 0 {
			return false
		}
	}

	// SP3.
	if swapArgs.AmountRepay.Cmp(repayAmount(approveTxOrNil, swapTx)) != 0 {
		return false
	}

	return true
}

// MakeLendTx creates a transaction with following properties:
// L1. LendTx.type = 0x7802 (TxTypeEthereumDynamicFee)
// L2. LendTx.from = proposer
// L3. LendTx.to = SwapTx.from
// L4. LendTx.value = LendAmount(approveTxOrNil, swapTx)
func (g *GaslessModule) GetLendTxGenerator(approveTxOrNil, swapTx *types.Transaction) builder.TxGenerator {
	return func(nonce uint64) (*types.Transaction, error) {
		var (
			to      = swapTx.ValidatedSender()
			chainId = g.InitOpts.ChainConfig.ChainID
			signer  = types.LatestSignerForChainID(chainId)
			key     = g.InitOpts.NodeKey
		)

		tx, err := types.NewTransactionWithMap(types.TxTypeEthereumDynamicFee, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:     nonce,
			types.TxValueKeyTo:        &to,
			types.TxValueKeyAmount:    lendAmount(approveTxOrNil, swapTx),
			types.TxValueKeyGasLimit:  params.TxGas,
			types.TxValueKeyGasFeeCap: swapTx.GasFeeCap(),
			types.TxValueKeyGasTipCap: swapTx.GasTipCap(),
			types.TxValueKeyChainID:   chainId,
		})
		if err != nil {
			return nil, err
		}

		err = tx.Sign(signer, key)
		return tx, err
	}
}

func (g *GaslessModule) ExtractTxBundles(txs []*types.Transaction, prevBundles []*builder.Bundle) []*builder.Bundle {
	// TODO: implement me
	return nil
}

func lendAmount(approveTxOrNil, swapTx *types.Transaction) *big.Int {
	r := new(big.Int)

	// R2 = ApproveTx.Fee() if exists
	if approveTxOrNil != nil {
		r.Add(r, approveTxOrNil.Fee())
	}

	// R3 = SwapTx.Fee()
	r.Add(r, swapTx.Fee())

	// LendAmount = R2 + R3
	return r
}

func repayAmount(approveTxOrNil, swapTx *types.Transaction) *big.Int {
	// R1 = LendTx.Fee() = SwapTx.GasPrice() * TxGas
	r1 := new(big.Int).Mul(swapTx.GasPrice(), new(big.Int).SetUint64(params.TxGas))

	// RepayAmount = R1 + R2 + R3
	return new(big.Int).Add(r1, lendAmount(approveTxOrNil, swapTx))
}
