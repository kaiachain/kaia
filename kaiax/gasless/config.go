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

package gasless

import (
	"math"

	"github.com/kaiachain/kaia/common"
	"github.com/urfave/cli/v2"
)

var (
	AllowedTokensFlag = &cli.StringSliceFlag{
		Name:     "gasless.allowed-tokens",
		Usage:    "allow token addresses for gasless module, allow all tokens if all",
		Value:    cli.NewStringSlice("all"),
		Aliases:  []string{"kaiax.module.gasless.allowed-tokens"},
		Category: "KAIAX",
	}
	DisableFlag = &cli.BoolFlag{
		Name:     "gasless.disable",
		Usage:    "disable gasless module",
		Value:    false,
		Aliases:  []string{"kaiax.module.gasless.disable"},
		Category: "KAIAX",
	}
	MaxBundleTxsInPendingFlag = &cli.IntFlag{
		Name:     "gasless.max-bundle-txs-in-pending",
		Usage:    "max number of gasless bundle txs in pending queue. Default value is 100. No limit if negative value",
		Value:    100,
		Aliases:  []string{"kaiax.module.gasless.max-bundle-txs-in-pending"},
		Category: "KAIAX",
	}
	MaxBundleTxsInQueueFlag = &cli.IntFlag{
		Name:     "gasless.max-bundle-txs-in-queue",
		Usage:    "max number of gasless bundle txs in queue. Default value is 200. No limit if negative value",
		Value:    200,
		Aliases:  []string{"kaiax.module.gasless.max-bundle-txs-in-queue"},
		Category: "KAIAX",
	}
)

type GaslessConfig struct {
	// all tokens are allowed if AllowedTokens is nil while all are disallowed if empty slice
	AllowedTokens         []common.Address `toml:",omitempty"`
	Disable               bool
	MaxBundleTxsInPending uint
	MaxBundleTxsInQueue   uint
}

func DefaultGaslessConfig() *GaslessConfig {
	return &GaslessConfig{
		AllowedTokens:         nil,
		Disable:               false,
		MaxBundleTxsInPending: 100,
		MaxBundleTxsInQueue:   200,
	}
}

func SetGaslessConfig(ctx *cli.Context, cfg *GaslessConfig) {
	if tokens := ctx.StringSlice(AllowedTokensFlag.Name); tokens != nil {
		cfg.AllowedTokens = []common.Address{}
		for _, addr := range tokens {
			if addr == "all" {
				cfg.AllowedTokens = nil
				break
			}
			cfg.AllowedTokens = append(cfg.AllowedTokens, common.HexToAddress(addr))
		}
	}

	cfg.Disable = ctx.Bool(DisableFlag.Name)

	// use default value if size is zero
	if size := ctx.Int(MaxBundleTxsInPendingFlag.Name); size > 0 {
		cfg.MaxBundleTxsInPending = uint(size)
	} else {
		cfg.MaxBundleTxsInPending = math.MaxUint64
	}

	if size := ctx.Int(MaxBundleTxsInQueueFlag.Name); size > 0 {
		cfg.MaxBundleTxsInQueue = uint(size)
	} else {
		cfg.MaxBundleTxsInQueue = math.MaxUint64
	}
}
