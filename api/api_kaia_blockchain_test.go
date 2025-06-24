// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
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

package api

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mock_api "github.com/kaiachain/kaia/api/mocks"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/params"
)

func testInitForKaiaApi(t *testing.T) (*gomock.Controller, *mock_api.MockBackend, *KaiaBlockChainAPI) {
	mockCtrl := gomock.NewController(t)
	mockBackend := mock_api.NewMockBackend(mockCtrl)

	blockchain.InitDeriveSha(params.TestChainConfig)

	api := NewKaiaBlockChainAPI(mockBackend)
	return mockCtrl, mockBackend, api
}

func TestKaiaAPI_EstimateGas(t *testing.T) {
	mockCtrl, mockBackend, api := testInitForKaiaApi(t)
	defer mockCtrl.Finish()

	testEstimateGas(t, mockBackend, func(ethArgs EthTransactionArgs, overrides *EthStateOverride) (hexutil.Uint64, error) {
		// Testcases are written in EthTransactionArgs. Convert to Kaia CallArgs
		args := CallArgs{
			From:                 ethArgs.from(),
			To:                   ethArgs.To,
			GasPrice:             ethArgs.GasPrice,
			MaxFeePerGas:         ethArgs.MaxFeePerGas,
			MaxPriorityFeePerGas: ethArgs.MaxPriorityFeePerGas,
			Data:                 ethArgs.data(),
		}
		if ethArgs.Gas != nil {
			args.Gas = *ethArgs.Gas
		}
		if ethArgs.Value != nil {
			args.Value = *ethArgs.Value
		}
		return api.EstimateGas(context.Background(), args, nil, overrides)
	})
}
