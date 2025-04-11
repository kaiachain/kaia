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
	"github.com/kaiachain/kaia/common"
	"github.com/urfave/cli/v2"
)

var (
	AllowedTokensFlag = &cli.StringSliceFlag{
		Name:     "gasless.allowed-tokens",
		Usage:    "allow token addresses for gasless module, allow all tokens if all",
		Value:    cli.NewStringSlice("all"),
		Aliases:  []string{"genesis.module.gasless.allowed-tokens"},
		Category: "KAIAX",
	}
	DisableFlag = &cli.BoolFlag{
		Name:     "gasless.disable",
		Usage:    "disable gasless module",
		Value:    false,
		Aliases:  []string{"genesis.module.gasless.disable"},
		Category: "KAIAX",
	}
	GaslessTxSlotsFlag = &cli.IntFlag{
		Name:     "gasless.gasless-tx-slots",
		Usage:    "number of gasless transaction pair slots in tx pool",
		Value:    100,
		Aliases:  []string{"genesis.module.gasless.gasless-tx-slots"},
		Category: "KAIAX",
	}
)

type GaslessConfig struct {
	// all tokens are allowed if AllowedTokens is nil while all are disallowed if empty slice
	AllowedTokens []common.Address `toml:",omitempty"`

	// disable gasless module
	Disable bool

	// max number of gasless tx pairs in tx pool
	GaslessTxSlots int
}

func GetGaslessConfig(ctx *cli.Context) *GaslessConfig {
	config := makeDefaultGaslessConfig()

	if tokens := ctx.StringSlice(AllowedTokensFlag.Name); tokens != nil {
		config.AllowedTokens = []common.Address{}
		for _, addr := range tokens {
			if addr == "all" {
				config.AllowedTokens = nil
				break
			}
			config.AllowedTokens = append(config.AllowedTokens, common.HexToAddress(addr))
		}
	}

	config.Disable = ctx.Bool(DisableFlag.Name)

	if slot := ctx.Int(GaslessTxSlotsFlag.Name); slot > 0 {
		config.GaslessTxSlots = slot
	}

	return config
}

func makeDefaultGaslessConfig() *GaslessConfig {
	return &GaslessConfig{
		AllowedTokens:  nil,
		Disable:        false,
		GaslessTxSlots: 100,
	}
}
