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
	"math/big"
	"sync"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gasless"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	logger = log.NewModuleLogger(log.KaiaxGasless)

	_ kaiax.BaseModule       = (*GaslessModule)(nil)
	_ kaiax.ExecutionModule  = (*GaslessModule)(nil)
	_ kaiax.TxPoolModule     = (*GaslessModule)(nil)
	_ kaiax.TxBundlingModule = (*GaslessModule)(nil)
)

type InitOpts struct {
	ChainConfig   *params.ChainConfig
	GaslessConfig *gasless.GaslessConfig
	NodeKey       *ecdsa.PrivateKey
	Chain         backends.BlockChainForCaller
	NodeType      common.ConnType // if CN, minimum balance required to enable gasless module
}

type GaslessModule struct {
	InitOpts
	swapRouter    common.Address
	allowedTokens map[common.Address]bool
	signer        types.Signer

	currentStateMu sync.Mutex     // even simple GetNonce affects statedb's internal state, hence can't use RWMutex.
	currentState   *state.StateDB // latest state for nonce lookup

	knownTxsMu sync.RWMutex
	knownTxs   *knownTxs
}

func NewGaslessModule() *GaslessModule {
	return &GaslessModule{
		allowedTokens: map[common.Address]bool{},
		knownTxs:      &knownTxs{},
	}
}

func (g *GaslessModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.GaslessConfig == nil || opts.NodeKey == nil || opts.Chain == nil {
		return ErrInitUnexpectedNil
	}

	g.InitOpts = *opts
	g.signer = types.LatestSignerForChainID(g.ChainConfig.ChainID)
	currentState, err := g.Chain.State()
	if err != nil {
		return err
	}
	g.setCurrentState(currentState)

	// Disable module if CN (lender) does not have sufficient balance
	if g.NodeType == common.CONSENSUSNODE {
		nodeAddr := crypto.PubkeyToAddress(opts.NodeKey.PublicKey)
		balance := g.getCurrentStateBalance(nodeAddr)
		if balance.Cmp(GaslessLenderMinBal) < 0 {
			g.GaslessConfig.Disable = true
			logger.Warn("disabling gasless module due to insufficient balance", "node", nodeAddr.Hex(), "balance", balance.String())
		}
	}

	return g.updateAddresses(g.Chain.CurrentBlock().Header())
}

func (g *GaslessModule) IsDisabled() bool {
	return g.GaslessConfig.Disable
}

func (g *GaslessModule) Start() error {
	return nil
}

func (g *GaslessModule) Stop() {
}

func (g *GaslessModule) setCurrentState(state *state.StateDB) {
	g.currentStateMu.Lock()
	defer g.currentStateMu.Unlock()

	g.currentState = state
}

func (g *GaslessModule) getCurrentStateNonce(addr common.Address) uint64 {
	g.currentStateMu.Lock()
	defer g.currentStateMu.Unlock()

	return g.currentState.GetNonce(addr)
}

func (g *GaslessModule) getCurrentStateBalance(addr common.Address) *big.Int {
	g.currentStateMu.Lock()
	defer g.currentStateMu.Unlock()

	return g.currentState.GetBalance(addr)
}

func (g *GaslessModule) getCurrentHasCode(addr common.Address) bool {
	g.currentStateMu.Lock()
	defer g.currentStateMu.Unlock()

	return g.currentState.GetCodeHash(addr) != types.EmptyCodeHash
}
