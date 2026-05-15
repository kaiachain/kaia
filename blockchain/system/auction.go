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
	"bytes"
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/common"
	contracts "github.com/kaiachain/kaia/contracts/bindings/auction"
	contractsv3 "github.com/kaiachain/kaia/contracts/bindings/auctionv3"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/params"
	"github.com/mitchellh/mapstructure"
)

var (
	abi, _   = contracts.IAuctionEntryPointMetaData.GetAbi()
	abiV3, _ = contractsv3.IAuctionEntryPointMetaData.GetAbi()
)

// ReadAuctioneer and ReadGasBufferEstimate use the v2.1 binding for both v2.1 and v3.0
// deployments. The auctioneer() and gasBufferEstimate() view functions have identical
// signatures across versions, so their selectors match — only the AuctionTx struct (used
// by call/getAuctionTxHash) differs.
func ReadAuctioneer(backend bind.ContractCaller, contractAddr common.Address, num *big.Int) (common.Address, error) {
	caller, err := contracts.NewIAuctionEntryPointCaller(contractAddr, backend)
	if err != nil {
		return common.Address{}, err
	}

	opts := &bind.CallOpts{BlockNumber: num}
	return caller.Auctioneer(opts)
}

func ReadGasBufferEstimate(backend bind.ContractCaller, contractAddr common.Address, num *big.Int) (uint64, error) {
	caller, err := contracts.NewIAuctionEntryPointCaller(contractAddr, backend)
	if err != nil {
		return 0, err
	}
	opts := &bind.CallOpts{BlockNumber: num}
	buffer, err := caller.GasBufferEstimate(opts)
	if err != nil {
		return 0, err
	}
	return buffer.Uint64(), nil
}

func EncodeAuctionCallData(bid *auction.Bid, chainConfig *params.ChainConfig) ([]byte, error) {
	if chainConfig != nil && chainConfig.IsPermissionlessForkEnabled(new(big.Int).SetUint64(bid.BlockNumber)) {
		maxGasPrice := bid.MaxGasPrice
		if maxGasPrice == nil {
			maxGasPrice = new(big.Int)
		}
		input := contractsv3.IAuctionEntryPointAuctionTx{
			TargetTxHash:  bid.TargetTxHash,
			BlockNumber:   new(big.Int).SetUint64(bid.BlockNumber),
			Sender:        bid.Sender,
			To:            bid.To,
			Nonce:         new(big.Int).SetUint64(bid.Nonce),
			Bid:           bid.Bid,
			MaxGasPrice:   maxGasPrice,
			CallGasLimit:  new(big.Int).SetUint64(bid.CallGasLimit),
			Data:          bid.Data,
			SearcherSig:   bid.SearcherSig,
			AuctioneerSig: bid.AuctioneerSig,
		}
		return abiV3.Pack("call", input)
	}

	input := contracts.IAuctionEntryPointAuctionTx{
		TargetTxHash:  bid.TargetTxHash,
		BlockNumber:   new(big.Int).SetUint64(bid.BlockNumber),
		Sender:        bid.Sender,
		To:            bid.To,
		Nonce:         new(big.Int).SetUint64(bid.Nonce),
		Bid:           bid.Bid,
		CallGasLimit:  new(big.Int).SetUint64(bid.CallGasLimit),
		Data:          bid.Data,
		SearcherSig:   bid.SearcherSig,
		AuctioneerSig: bid.AuctioneerSig,
	}
	return abi.Pack("call", input)
}

// DecodeAuctionCallData is the reverse function of EncodeAuctionCallData.
// It can be used to identify the calldata of a bid. The v2.1/v3.0 ABI is
// auto-detected by matching the method selector.
func DecodeAuctionCallData(encoded []byte) (*auction.Bid, error) {
	if len(encoded) <= 4 {
		return nil, errors.New("invalid encoded data")
	}

	if bytes.Equal(encoded[:4], abiV3.Methods["call"].ID) {
		decoded, err := abiV3.Methods["call"].Inputs.Unpack(encoded[4:])
		if err != nil {
			return nil, err
		}
		var bid contractsv3.IAuctionEntryPointAuctionTx
		if err := mapstructure.Decode(decoded[0], &bid); err != nil {
			return nil, err
		}
		bidData := auction.BidData{
			TargetTxHash:  bid.TargetTxHash,
			BlockNumber:   bid.BlockNumber.Uint64(),
			Sender:        bid.Sender,
			To:            bid.To,
			Nonce:         bid.Nonce.Uint64(),
			Bid:           bid.Bid,
			MaxGasPrice:   bid.MaxGasPrice,
			CallGasLimit:  bid.CallGasLimit.Uint64(),
			Data:          bid.Data,
			SearcherSig:   bid.SearcherSig,
			AuctioneerSig: bid.AuctioneerSig,
		}
		return &auction.Bid{BidData: bidData}, nil
	}

	decoded, err := abi.Methods["call"].Inputs.Unpack(encoded[4:])
	if err != nil {
		return nil, err
	}
	var bid contracts.IAuctionEntryPointAuctionTx
	if err := mapstructure.Decode(decoded[0], &bid); err != nil {
		return nil, err
	}
	bidData := auction.BidData{
		TargetTxHash:  bid.TargetTxHash,
		BlockNumber:   bid.BlockNumber.Uint64(),
		Sender:        bid.Sender,
		To:            bid.To,
		Nonce:         bid.Nonce.Uint64(),
		Bid:           bid.Bid,
		CallGasLimit:  bid.CallGasLimit.Uint64(),
		Data:          bid.Data,
		SearcherSig:   bid.SearcherSig,
		AuctioneerSig: bid.AuctioneerSig,
	}
	return &auction.Bid{BidData: bidData}, nil
}
