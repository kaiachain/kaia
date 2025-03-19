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
	"crypto/ecdsa"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	gasless_cfg "github.com/kaiachain/kaia/kaiax/gasless/config"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var logger = log.NewModuleLogger(log.KaiaxGasless)

type InitOpts struct {
	ChainConfig *params.ChainConfig
	CNConfig    *gasless_cfg.CNConfig
	NodeKey     *ecdsa.PrivateKey
	Chain       backends.BlockChainForCaller
	TxPool      kaiax.TxPoolForCaller
}

type GaslessModule struct {
	InitOpts
	swapRouter    *common.Address
	allowedTokens map[common.Address]bool
	signer        types.Signer
}

func NewGaslessModule() *GaslessModule {
	return &GaslessModule{}
}

func (g *GaslessModule) Init(opts *InitOpts) (disabled bool, err error) {
	if opts == nil || opts.ChainConfig == nil || opts.CNConfig == nil || opts.NodeKey == nil || opts.Chain == nil || opts.TxPool == nil {
		return true, ErrInitUnexpectedNil
	}

	if opts.CNConfig.Disable {
		return true, nil
	}

	g.InitOpts = *opts
	g.signer = types.LatestSignerForChainID(g.ChainConfig.ChainID)

	err = g.updateAddresses(g.Chain.CurrentBlock().Number())
	if err != nil {
		return true, err
	}

	return false, nil
}

func (g *GaslessModule) Start() error {
	return nil
}

func (g *GaslessModule) Stop() {
}
