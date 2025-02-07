// Copyright 2024 The Kaia Authors
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
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gasless"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	_ gasless.GaslessModule = (*GaslessModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGasless)
)

type InitOpts struct {
	ChainConfig *params.ChainConfig
	NodeKey     *ecdsa.PrivateKey
}

type GaslessModule struct {
	InitOpts

	swapRouters   map[common.Address]bool
	allowedTokens map[common.Address]bool
}

func NewGaslessModule() *GaslessModule {
	return &GaslessModule{}
}

func (g *GaslessModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.NodeKey == nil {
		return ErrInitUnexpectedNil
	}
	g.InitOpts = *opts

	g.swapRouters = map[common.Address]bool{
		common.HexToAddress("0x1234"): true,
	}
	g.allowedTokens = map[common.Address]bool{
		common.HexToAddress("0xabcd"): true,
	}
	return nil
}

func (g *GaslessModule) Start() error {
	return nil
}

func (g *GaslessModule) Stop() {
}
