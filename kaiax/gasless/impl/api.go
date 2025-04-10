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
	"context"
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

func (b *GaslessModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewGaslessAPI(b),
			Public:    true,
		},
	}
}

type GaslessAPI struct {
	b *GaslessModule
}

func NewGaslessAPI(b *GaslessModule) *GaslessAPI {
	return &GaslessAPI{b}
}

// GaslessTxResult represents the result of checking if a transaction is a gasless transaction
type GaslessTxResponse struct {
	IsGasless bool    `json:"isGasless"` // Whether the transaction is a gasless transaction
	Reason    *string `json:"reason"`    // Reason why the transaction is not a gasless transaction, empty if it is
}

func ToResponse(err error) *GaslessTxResponse {
	var pErrStr *string
	isExecutable := true
	if err != nil {
		errStr := err.Error()
		pErrStr = &errStr
		isExecutable = false
	}
	return &GaslessTxResponse{
		IsGasless: isExecutable,
		Reason:    pErrStr,
	}
}

// IsGaslessTx checks if the given raw transactions form a valid gasless transaction
// It returns a detailed result explaining why a transaction is not a valid gasless transaction if it's not
func (s *GaslessAPI) IsGaslessTx(ctx context.Context, rawTxs []hexutil.Bytes) *GaslessTxResponse {
	if len(rawTxs) == 0 {
		return ToResponse(errors.New("no transactions provided"))
	}

	// Decode the raw transactions
	txs := make([]*types.Transaction, 0, len(rawTxs))
	for i, rawTx := range rawTxs {
		if len(rawTx) == 0 {
			return ToResponse(fmt.Errorf("empty transaction at index %d", i))
		}

		// Handle Ethereum transaction envelope
		if 0 < rawTx[0] && rawTx[0] < 0x7f {
			rawTx = append([]byte{byte(types.EthereumTxTypeEnvelope)}, rawTx...)
		}

		tx := new(types.Transaction)
		if err := rlp.DecodeBytes(rawTx, tx); err != nil {
			return ToResponse(fmt.Errorf("failed to decode transaction at index %d: %v", i, err))
		}

		txs = append(txs, tx)
	}

	// Check if the transactions form a valid gasless transaction
	// Case 1: A single swap transaction
	if len(txs) == 1 {
		swapTx := txs[0]
		if !s.b.IsSwapTx(swapTx) {
			return ToResponse(errors.New("transaction is not a swap transaction"))
		}

		return ToResponse(s.b.VerifyExecutable(nil, swapTx))
	}

	// Case 2: An approve transaction followed by a swap transaction
	if len(txs) == 2 {
		approveTx := txs[0]
		swapTx := txs[1]

		if !s.b.IsApproveTx(approveTx) {
			return ToResponse(errors.New("first transaction is not an approve transaction"))
		}

		if !s.b.IsSwapTx(swapTx) {
			return ToResponse(errors.New("second transaction is not a swap transaction"))
		}

		err := s.b.VerifyExecutable(approveTx, swapTx)
		return ToResponse(err)
	}

	return ToResponse(fmt.Errorf("expected 1 or 2 transactions, got %d", len(txs)))
}

type GaslessInfoResult struct {
	IsDisabled    bool             `json:"isDisabled"`
	SwapRouter    common.Address   `json:"swapRouter"`
	AllowedTokens []common.Address `json:"allowedTokens"`
}

func (s *GaslessAPI) GaslessInfo() *GaslessInfoResult {
	at := []common.Address{}
	for addr := range s.b.allowedTokens {
		at = append(at, addr)
	}
	return &GaslessInfoResult{
		IsDisabled:    s.b.IsDisabled(),
		SwapRouter:    s.b.swapRouter,
		AllowedTokens: at,
	}
}
