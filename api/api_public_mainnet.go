// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
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
	"errors"
	"strings"

	"github.com/kaiachain/kaia/v2/accounts/abi"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/contracts/contracts/system_contracts/misc"
	"github.com/kaiachain/kaia/v2/networks/rpc"
)

// MainnetCredit contract is stored in the address zero.
var (
	mainnetCreditContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
	latestBlockNrOrHash          = rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	errNoCypressCreditContract   = errors.New("no mainnet credit contract")
)

type CreditOutput struct {
	Photo       string `json:"photo"`
	Names       string `json:"names"`
	EndingPhoto string `json:"endingPhoto"`
	EndingNames string `json:"endingNames"`
}

// callCypressCreditGetFunc executes funcName in CypressCreditContract and returns the output.
func (s *PublicBlockChainAPI) callCypressCreditGetFunc(ctx context.Context, parsed *abi.ABI, funcName string) (*string, error) {
	abiGet, err := parsed.Pack(funcName)
	if err != nil {
		return nil, err
	}

	args := CallArgs{
		To:   &mainnetCreditContractAddress,
		Data: abiGet,
	}
	ret, err := s.Call(ctx, args, latestBlockNrOrHash)
	if err != nil {
		return nil, err
	}

	output := new(string)
	err = parsed.UnpackIntoInterface(output, funcName, ret)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetCypressCredit calls getPhoto and getNames in the CypressCredit contract
// and returns all the results as a struct.
func (s *PublicBlockChainAPI) GetCypressCredit(ctx context.Context) (*CreditOutput, error) {
	if ok, err := s.IsContractAccount(ctx, mainnetCreditContractAddress, latestBlockNrOrHash); err != nil {
		return nil, err
	} else if !ok {
		return nil, errNoCypressCreditContract
	}

	parsed, err := abi.JSON(strings.NewReader(misc.CypressCreditV2ABI))
	if err != nil {
		return nil, err
	}

	output := new(CreditOutput)

	// getPhoto and getNames must exist from the Cypress genesis.
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getPhoto"); err == nil {
		output.Photo = *str
	} else {
		return nil, err
	}
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getNames"); err == nil {
		output.Names = *str
	} else {
		return nil, err
	}

	// getEndingPhoto and getEndingNames are added at some nonzero block. They are returned if they exist.
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getEndingPhoto"); err == nil {
		output.EndingPhoto = *str
	}
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getEndingNames"); err == nil {
		output.EndingNames = *str
	}

	return output, nil
}
