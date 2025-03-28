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
	"math/big"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

const (
	RPC_AUCTION_HASH_PROP         = "bidHash"
	RPC_AUCTION_TARGET_DECODE_ERR = "errTargetTxDecode"
	RPC_AUCTION_TARGET_SEND_ERR   = "errTargetTxSend"
	RPC_AUCTION_BID_VALIDATE_ERR  = "errValidateBid"
)

func (a *AuctionModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "auction",
			Version:   "1.0",
			Service:   newAuctionAPI(a),
			Public:    false,
		},
	}
}

type AuctionAPI struct {
	a *AuctionModule
}

func newAuctionAPI(a *AuctionModule) *AuctionAPI {
	return &AuctionAPI{a: a}
}

type RPCOutput map[string]any

// BidInput is the same format with `BidData`, replacing `[]byte` type to `hexutil.Bytes`
type BidInput struct {
	TargetTxRaw   hexutil.Bytes  `json:"targetTxRaw"`
	TargetTxHash  common.Hash    `json:"targetTxHash"`
	BlockNumber   uint64         `json:"blockNumber"`
	Sender        common.Address `json:"sender"`
	To            common.Address `json:"to"`
	Nonce         uint64         `json:"nonce"`
	Bid           *big.Int       `json:"bid"`
	CallGasLimit  uint64         `json:"callGasLimit"`
	Data          hexutil.Bytes  `json:"data"`
	SearcherSig   hexutil.Bytes  `json:"searcherSig"`
	AuctioneerSig hexutil.Bytes  `json:"auctioneerSig"`
}

func toBid(bidInput BidInput) *auction.Bid {
	bidData := auction.BidData{
		TargetTxHash:  bidInput.TargetTxHash,
		BlockNumber:   bidInput.BlockNumber,
		Sender:        bidInput.Sender,
		To:            bidInput.To,
		Nonce:         bidInput.Nonce,
		Bid:           bidInput.Bid,
		CallGasLimit:  bidInput.CallGasLimit,
		Data:          bidInput.Data,
		SearcherSig:   bidInput.SearcherSig,
		AuctioneerSig: bidInput.AuctioneerSig,
	}
	return &auction.Bid{BidData: bidData}
}

func toTx(targetTxRaw []byte) (*types.Transaction, error) {
	if 0 < targetTxRaw[0] && targetTxRaw[0] < 0x7f {
		targetTxRaw = append([]byte{byte(types.EthereumTxTypeEnvelope)}, targetTxRaw...)
	}
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(targetTxRaw, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (api *AuctionAPI) SubmitBid(ctx context.Context, bidInput BidInput) RPCOutput {
	//  1. directly send traget transaction (target tx can be empty)
	if len(bidInput.TargetTxRaw) > 0 {
		targetTx, errTxDecode := toTx(bidInput.TargetTxRaw)
		if errTxDecode != nil {
			return makeRPCOutput(common.Hash{}, errTxDecode, nil, nil)
		}
		if targetTx.Hash() != bidInput.TargetTxHash {
			return makeRPCOutput(common.Hash{}, auction.ErrInvalidTargetTxHash, nil, nil)
		}
		errTargetTxSend := api.a.Backend.SendTx(ctx, targetTx)
		// ignore `nonce too low` error against target tx validation
		if errTargetTxSend != nil && errTargetTxSend != blockchain.ErrNonceTooLow {
			return makeRPCOutput(common.Hash{}, nil, errTargetTxSend, nil)
		}
	}

	// 2. add bid
	bid := toBid(bidInput)
	bidHash, errValidateBid := api.a.bidPool.AddBid(bid)
	return makeRPCOutput(bidHash, nil, nil, errValidateBid)
}

func makeRPCOutput(bidHash common.Hash, errTargetTxDecode, errTargetTxSend, errValidateBid error) RPCOutput {
	var (
		m                                                           = make(map[string]any)
		errTargetTxDecodeStr, errTargetTxSendStr, errValidateBidStr string
	)
	if errTargetTxDecode != nil {
		errTargetTxDecodeStr = errTargetTxDecode.Error()
	}
	if errTargetTxSend != nil {
		errTargetTxSendStr = errTargetTxSend.Error()
	}
	if errValidateBid != nil {
		errValidateBidStr = errValidateBid.Error()
	}
	m[RPC_AUCTION_HASH_PROP] = bidHash
	m[RPC_AUCTION_TARGET_DECODE_ERR] = errTargetTxDecodeStr
	m[RPC_AUCTION_TARGET_SEND_ERR] = errTargetTxSendStr
	m[RPC_AUCTION_BID_VALIDATE_ERR] = errValidateBidStr
	return m
}
