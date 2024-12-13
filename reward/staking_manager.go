// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package reward

import (
	"fmt"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	logger = log.NewModuleLogger(log.Reward)
)

// blockChain is an interface for blockchain.Blockchain used in reward package.
type blockChain interface {
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	GetBlockByNumber(number uint64) *types.Block
	GetReceiptsByBlockHash(hash common.Hash) types.Receipts
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig
	CurrentHeader() *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetHeaderByNumber(number uint64) *types.Header
	State() (*state.StateDB, error)
	CurrentBlock() *types.Block

	StateCache() state.Database
	Processor() blockchain.Processor

	blockchain.ChainContext
}

// PreloadStakingInfo preloads staking info for the given headers.
// It first finds the first block that does not have state, and then
// it regenerates the state from the nearest block that has state to the target block to preload staking info.
// Note that the state is saved every 128 blocks to disk in full node.
func PreloadStakingInfo(bc blockChain, headers []*types.Header, stakingModule staking.StakingModule) (uint64, error) {
	// If no headers to preload, do nothing
	if len(headers) == 0 {
		return 0, nil
	}

	var (
		current  *types.Block
		database state.Database
		target   = headers[len(headers)-1].Number.Uint64()
	)

	database = state.NewDatabaseWithExistingCache(bc.StateCache().TrieDB().DiskDB(), bc.StateCache().TrieDB().TrieNodeCache())

	// Find the first block that does not have state
	i := 0
	for i < len(headers) {
		if _, err := state.New(headers[i].Root, database, nil, nil); err != nil {
			break
		}
		i++
	}
	// Early return if all blocks have state
	if i == len(headers) {
		return 0, nil
	}

	// Find the nearest block that has state
	origin := headers[i].Number.Uint64() - headers[i].Number.Uint64()%128
	current = bc.GetBlockByNumber(origin)
	if current == nil {
		return 0, fmt.Errorf("block %d not found", origin)
	}
	statedb, err := state.New(current.Header().Root, database, nil, nil)
	if err != nil {
		return 0, err
	}

	var (
		parent     common.Hash
		preloadRef = stakingModule.AllocPreloadRef()
	)

	// Include target since we want staking info at `target`, not for `target`.
	for current.NumberU64() <= target {
		stakingModule.PreloadFromState(preloadRef, current.Header(), statedb)
		if current.NumberU64() == target {
			break
		}
		// Retrieve the next block to regenerate and process it
		next := current.NumberU64() + 1
		if current = bc.GetBlockByNumber(next); current == nil {
			return preloadRef, fmt.Errorf("block #%d not found", next)
		}
		_, _, _, _, _, err := bc.Processor().Process(current, statedb, vm.Config{})
		if err != nil {
			return preloadRef, fmt.Errorf("processing block %d failed: %v", current.NumberU64(), err)
		}
		// Finalize the state so any modifications are written to the trie
		root, err := statedb.Commit(true)
		if err != nil {
			return preloadRef, err
		}
		if err := statedb.Reset(root); err != nil {
			return preloadRef, fmt.Errorf("state reset after block %d failed: %v", current.NumberU64(), err)
		}
		database.TrieDB().ReferenceRoot(root)
		if !common.EmptyHash(parent) {
			database.TrieDB().Dereference(parent)
		}
		if current.Root() != root {
			err = fmt.Errorf("mistmatching state root block expected %x reexecuted %x", current.Root(), root)
			// Logging here because something went wrong when the state roots disagree even if the execution was successful.
			logger.Error("incorrectly regenerated historical state", "block", current.NumberU64(), "err", err)
			return preloadRef, fmt.Errorf("incorrectly regenerated historical state for block %d: %v", current.NumberU64(), err)
		}
		parent = root
	}

	return preloadRef, nil
}
