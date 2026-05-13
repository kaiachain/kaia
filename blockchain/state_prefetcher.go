// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2020 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state_prefetcher.go (2019/04/02).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"sync/atomic"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
)

// statePrefetcher is a basic Prefetcher, which blindly executes a block on top
// of an arbitrary state with the goal of prefetching potentially useful state
// data from disk before the main block processor start executing.
type statePrefetcher struct {
	chain ChainContext // Canonical block chain
}

// newStatePrefetcher initialises a new statePrefetcher.
func newStatePrefetcher(chain ChainContext) *statePrefetcher {
	return &statePrefetcher{
		chain: chain,
	}
}

// Prefetch warms state caches via prefetchTxState (no EVM).
func (p *statePrefetcher) Prefetch(block *types.Block, stateDB *state.StateDB, _ vm.Config, interrupt *uint32) {
	signer := types.MakeSigner(p.chain.Config(), block.Header().Number)
	blockNumber := block.NumberU64()
	for _, tx := range block.Transactions() {
		if interrupt != nil && atomic.LoadUint32(interrupt) == 1 {
			return
		}
		prefetchTxState(stateDB, signer, tx, blockNumber)
	}
}
