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

package auction

import (
	"fmt"
	"math"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/urfave/cli/v2"
)

const (
	DefaultMaxBidPoolSize = math.MaxInt64
	DefaultEDOffset       = 200 * time.Millisecond
)

var (
	DisableFlag = &cli.BoolFlag{
		Name:     "auction.disable",
		Usage:    "disable auction module",
		Value:    false,
		Aliases:  []string{"kaiax.module.auction.disable"},
		Category: "KAIAX",
	}
	MaxBidPoolSizeFlag = &cli.Int64Flag{
		Name:     "auction.max-bid-pool-size",
		Usage:    "max number of bids in bid pool. Default value is max uint64.",
		Value:    DefaultMaxBidPoolSize,
		Aliases:  []string{"kaiax.module.auction.max-bid-pool-size"},
		Category: "KAIAX",
	}
	EDOffsetFlag = &cli.DurationFlag{
		Name:     "auction.ed-offset",
		Usage:    "Early deadline offset used in auction. Default value is 200 milliseconds.",
		Value:    DefaultEDOffset,
		Aliases:  []string{"kaiax.module.auction.ed-offset"},
		Category: "KAIAX",
		Action: func(ctx *cli.Context, d time.Duration) error {
			if d < 0 || d > 1*time.Second {
				return fmt.Errorf("auction.ed-offset must be between 0 and 1 second, got %s", d)
			}
			return nil
		},
	}
)

type AuctionConfig struct {
	Disable        bool
	MaxBidPoolSize int64
	EDOffset       time.Duration
}

func DefaultAuctionConfig() *AuctionConfig {
	return &AuctionConfig{
		Disable:        false,
		MaxBidPoolSize: DefaultMaxBidPoolSize,
		EDOffset:       DefaultEDOffset,
	}
}

func SetAuctionConfig(ctx *cli.Context, cfg *AuctionConfig, nodeType common.ConnType) {
	disable := ctx.Bool(DisableFlag.Name)
	// Disable auction module for non-consensus nodes
	if nodeType != common.CONSENSUSNODE {
		disable = true
	}

	cfg.Disable = disable
	cfg.MaxBidPoolSize = DefaultMaxBidPoolSize
	if ctx.IsSet(MaxBidPoolSizeFlag.Name) {
		cfg.MaxBidPoolSize = ctx.Int64(MaxBidPoolSizeFlag.Name)
	}
	cfg.EDOffset = DefaultEDOffset
	if ctx.IsSet(EDOffsetFlag.Name) {
		cfg.EDOffset = ctx.Duration(EDOffsetFlag.Name)
	}
}
