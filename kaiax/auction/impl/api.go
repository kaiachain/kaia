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
	"strings"

	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/auction"
	gasless_impl "github.com/kaiachain/kaia/kaiax/gasless/impl"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node/cn/filters"
	"github.com/kaiachain/kaia/rlp"
)

const (
	RPC_AUCTION_ERROR_PROP        = "err"
	RPC_AUCTION_HASH_PROP         = "bidHash"
	RPC_AUCTION_TARGET_DECODE_ERR = "errTargetTxDecode"
	RPC_AUCTION_TARGET_SEND_ERR   = "errTargetTxSend"
	RPC_AUCTION_BID_VALIDATE_ERR  = "errValidateBid"
)

var (
	EMPTY_HASH          = common.Hash{}
	ErrEmptyTargetTxRaw = errors.New("Empty target tx raw")
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
	a         *AuctionModule
	f         *filters.KaiaFilterAPI
	publicAPI *api.KaiaBlockChainAPI
}

func newAuctionAPI(a *AuctionModule) *AuctionAPI {
	kaiaBlockChainAPI := api.NewKaiaBlockChainAPI(a.Backend.(api.Backend))
	return &AuctionAPI{a: a, f: filters.NewAuctionFilterAPI(a.Backend), publicAPI: kaiaBlockChainAPI}
}

type RPCOutput map[string]any

// BidInput is the same format with `BidData`, execpt adding new field `TargetTxRaw` and replacing `[]byte` type to `hexutil.Bytes`
type BidInput struct {
	TargetTxRaw   hexutil.Bytes  `json:"targetTxRaw"`
	TargetTxHash  common.Hash    `json:"targetTxHash"`
	BlockNumber   uint64         `json:"blockNumber"`
	Sender        common.Address `json:"sender"`
	To            common.Address `json:"to"`
	Nonce         uint64         `json:"nonce"`
	Bid           hexutil.Big    `json:"bid"`
	CallGasLimit  uint64         `json:"callGasLimit"`
	Data          hexutil.Bytes  `json:"data"`
	SearcherSig   hexutil.Bytes  `json:"searcherSig"`
	AuctioneerSig hexutil.Bytes  `json:"auctioneerSig"`
}

func ToBid(bidInput BidInput) *auction.Bid {
	bidData := auction.BidData{
		TargetTxHash:  bidInput.TargetTxHash,
		BlockNumber:   bidInput.BlockNumber,
		Sender:        bidInput.Sender,
		To:            bidInput.To,
		Nonce:         bidInput.Nonce,
		Bid:           bidInput.Bid.ToInt(),
		CallGasLimit:  bidInput.CallGasLimit,
		Data:          bidInput.Data,
		SearcherSig:   bidInput.SearcherSig,
		AuctioneerSig: bidInput.AuctioneerSig,
	}
	return &auction.Bid{BidData: bidData}
}

func toTx(targetTxRaw []byte) (*types.Transaction, error) {
	if len(targetTxRaw) == 0 {
		return nil, ErrEmptyTargetTxRaw
	}
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(targetTxRaw, tx); err != nil {
		return nil, err
	}
	return tx, nil
}

func (api *AuctionAPI) SubmitBid(ctx context.Context, bidInput BidInput) RPCOutput {
	if api.a.IsDisabled() {
		return makeRPCOutput(EMPTY_HASH, auction.ErrAuctionDisabled)
	}

	//  1. directly send target transaction
	targetTx, errTxDecode := toTx(bidInput.TargetTxRaw)
	if errTxDecode != nil {
		return makeRPCOutput(EMPTY_HASH, errTxDecode)
	}
	if targetTx.Hash() != bidInput.TargetTxHash {
		return makeRPCOutput(EMPTY_HASH, auction.ErrInvalidTargetTxHash)
	}
	errTargetTxSend := api.a.Backend.SendTx(ctx, targetTx)
	// ignore known transaction related errors against target tx validation
	if errTargetTxSend != nil && !(strings.HasPrefix(errTargetTxSend.Error(), "known transaction:") || errors.Is(errTargetTxSend, gasless_impl.ErrUnableToAddKnownBundleTx)) {
		return makeRPCOutput(EMPTY_HASH, errTargetTxSend)
	}

	// 2. add bid
	bid := ToBid(bidInput)
	bidHash, errValidateBid := api.a.bidPool.AddBid(bid)
	return makeRPCOutput(bidHash, errValidateBid)
}

func (api *AuctionAPI) NewPendingTransactions(ctx context.Context, fullTx *bool) (*rpc.Subscription, error) {
	return api.f.NewPendingTransactions(ctx, fullTx)
}

func (api *AuctionAPI) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	return api.f.NewHeads(ctx)
}

func (api *AuctionAPI) Logs(ctx context.Context, crit filters.FilterCriteria) (*rpc.Subscription, error) {
	return api.f.Logs(ctx, crit)
}

func (api *AuctionAPI) Call(ctx context.Context, args api.CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return api.publicAPI.Call(ctx, args, blockNrOrHash)
}

func makeRPCOutput(bidHash common.Hash, err error) RPCOutput {
	m := make(map[string]any)
	if err != nil {
		m[RPC_AUCTION_ERROR_PROP] = err.Error()
	}
	if bidHash != EMPTY_HASH {
		m[RPC_AUCTION_HASH_PROP] = bidHash
	}
	return m
}
