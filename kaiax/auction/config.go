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
	"github.com/kaiachain/kaia/common"
	"github.com/urfave/cli/v2"
)

var DisableFlag = &cli.BoolFlag{
	Name:     "auction.disable",
	Usage:    "disable auction module",
	Value:    false,
	Aliases:  []string{"genesis.module.auction.disable"},
	Category: "KAIAX",
}

type AuctionConfig struct {
	Disable bool
}

func GetAuctionConfig(ctx *cli.Context, nodeType common.ConnType) *AuctionConfig {
	disable := ctx.Bool(DisableFlag.Name)
	// Disable auction module for non-consensus nodes
	if nodeType != common.CONSENSUSNODE {
		disable = true
	}

	return &AuctionConfig{
		Disable: disable,
	}
}
