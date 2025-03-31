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
	"crypto/ecdsa"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	_ auction.AuctionModule = &AuctionModule{}

	logger = log.NewModuleLogger(log.KaiaxAuction)
)

type ProtocolManagerDownloader interface {
	Synchronising() bool
}

type apiBackend interface {
	SendTx(ctx context.Context, signedTx *types.Transaction) error
}

type InitOpts struct {
	ChainConfig *params.ChainConfig
	Chain       backends.BlockChainForCaller
	Backend     apiBackend
	Downloader  ProtocolManagerDownloader

	NodeKey *ecdsa.PrivateKey
}

type AuctionModule struct {
	InitOpts

	bidPool *BidPool
}

const (
	AuctionEarlyDeadline = 150 * time.Millisecond
)

func NewAuctionModule() *AuctionModule {
	return &AuctionModule{}
}

func (a *AuctionModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.Chain == nil || opts.Backend == nil || opts.Downloader == nil || opts.NodeKey == nil {
		return auction.ErrInitUnexpectedNil
	}

	a.bidPool = NewBidPool(opts.ChainConfig, opts.Chain)
	if a.bidPool == nil {
		return auction.ErrInitUnexpectedNil
	}

	a.InitOpts = *opts

	return nil
}

func (a *AuctionModule) Start() error {
	a.bidPool.start()
	return nil
}

func (a *AuctionModule) Stop() {
	// Clear the existing auction pool.
	a.bidPool.stop()
}
