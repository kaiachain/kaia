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

package system

import (
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/common"
	contracts "github.com/kaiachain/kaia/contracts/contracts/system_contracts/auction"
	"github.com/kaiachain/kaia/kaiax/auction"
)

func ReadAuctioneer(backend bind.ContractCaller, contractAddr common.Address, num *big.Int) (common.Address, error) {
	caller, err := contracts.NewIAuctionEntryPointCaller(contractAddr, backend)
	if err != nil {
		return common.Address{}, err
	}

	opts := &bind.CallOpts{BlockNumber: num}
	return caller.Auctioneer(opts)
}

func EncodeAuctionCallData(bid *auction.Bid) ([]byte, error) {
	abi, err := contracts.IAuctionEntryPointMetaData.GetAbi()
	if err != nil {
		return nil, err
	}

	input := contracts.IAuctionEntryPointAuctionTx{
		TargetTxHash:  bid.TargetTxHash,
		BlockNumber:   big.NewInt(int64(bid.BlockNumber)),
		Sender:        bid.Sender,
		To:            bid.To,
		Nonce:         big.NewInt(int64(bid.Nonce)),
		Bid:           bid.Bid,
		CallGasLimit:  big.NewInt(int64(bid.CallGasLimit)),
		Data:          bid.Data,
		SearcherSig:   bid.SearcherSig,
		AuctioneerSig: bid.AuctioneerSig,
	}

	data, err := abi.Pack("call", input)
	if err != nil {
		return nil, err
	}

	return data, nil
}
