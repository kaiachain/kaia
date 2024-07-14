// Modifications Copyright 2022 The klaytn Authors
// Copyright 2021 The go-ethereum Authors
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

// This file is derived from eth/state_accessor.go (2022/08/08).
// Modified and improved for the klaytn development.

package cn

import (
	"errors"
	"fmt"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/reward"
	statedb2 "github.com/kaiachain/kaia/storage/statedb"
)

// stateAtBlock retrieves the state database associated with a certain block.
// If no state is locally available for the given block, a number of blocks
// are attempted to be reexecuted to generate the desired state. The optional
// base layer statedb can be passed then it's regarded as the statedb of the
// parent block.
// Parameters:
//   - block: The block for which we want the state (== state at the stateRoot of the parent)
//   - reexec: The maximum number of blocks to reprocess trying to obtain the desired state
//   - base: If the caller is tracing multiple blocks, the caller can provide the parent state
//     continuously from the callsite.
//   - checklive: if true, then the live 'blockchain' state database is used. If the caller want to
//     perform Commit or other 'save-to-disk' changes, this should be set to false to avoid
//     storing trash persistently
//   - preferDisk: this arg can be used by the caller to signal that even though the 'base' is provided,
//     it would be preferrable to start from a fresh state, if we have it on disk.
func (cn *CN) stateAtBlock(block *types.Block, reexec uint64, base *state.StateDB, checkLive bool, preferDisk bool) (statedb *state.StateDB, err error) {
	var (
		current  *types.Block
		database state.Database
		report   = true
		origin   = block.NumberU64()
	)
	// Check the live database first if we have the state fully available, use that.
	if checkLive {
		statedb, err = cn.blockchain.StateAt(block.Root())
		if err == nil {
			return statedb, nil
		}
	}
	if base != nil {
		if preferDisk {
			// Create an ephemeral trie.Database for isolating the live one. Otherwise
			// the internal junks created by tracing will be persisted into the disk.
			database = state.NewDatabaseWithExistingCache(cn.ChainDB(), cn.blockchain.StateCache().TrieDB().TrieNodeCache())
			if statedb, err = state.New(block.Root(), database, nil, nil); err == nil {
				logger.Info("Found disk backend for state trie", "root", block.Root(), "number", block.Number())
				return statedb, nil
			}
		}
		// The optional base statedb is given, mark the start point as parent block
		statedb, database, report = base, base.Database(), false
		current = cn.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	} else {
		// Otherwise try to reexec blocks until we find a state or reach our limit
		current = block

		// Create an ephemeral trie.Database for isolating the live one. Otherwise
		// the internal junks created by tracing will be persisted into the disk.
		database = state.NewDatabaseWithExistingCache(cn.ChainDB(), cn.blockchain.StateCache().TrieDB().TrieNodeCache())

		for i := uint64(0); i < reexec; i++ {
			if current.NumberU64() == 0 {
				return nil, errors.New("genesis state is missing")
			}
			parent := cn.blockchain.GetBlock(current.ParentHash(), current.NumberU64()-1)
			if parent == nil {
				return nil, fmt.Errorf("missing block %v %d", current.ParentHash(), current.NumberU64()-1)
			}
			current = parent

			statedb, err = state.New(current.Root(), database, nil, nil)
			if err == nil {
				break
			}
		}
		if err != nil {
			switch err.(type) {
			case *statedb2.MissingNodeError:
				return nil, fmt.Errorf("historical state unavailable. tried regeneration but not possible, possibly due to state migration/pruning or global state saving interval is bigger than reexec value (reexec=%d)", reexec)
			default:
				return nil, err
			}
		}
	}
	// State was available at historical point, regenerate
	var (
		start  = time.Now()
		logged time.Time
		parent common.Hash

		preloadedStakingBlockNums = make([]uint64, 0, origin-current.NumberU64())
	)
	defer func() {
		for _, num := range preloadedStakingBlockNums {
			reward.UnloadStakingInfo(num)
		}
	}()
	for current.NumberU64() < origin {
		// Print progress logs if long enough time elapsed
		if report && time.Since(logged) > 8*time.Second {
			logger.Info("Regenerating historical state", "block", current.NumberU64()+1, "target", origin, "remaining", origin-current.NumberU64()-1, "elapsed", time.Since(start))
			logged = time.Now()
		}
		// Quit the state regeneration if time limit exceeds
		if cn.config.DisableUnsafeDebug && time.Since(start) > cn.config.StateRegenerationTimeLimit {
			return nil, fmt.Errorf("this request has queried old states too long since it exceeds the state regeneration time limit(%s)", cn.config.StateRegenerationTimeLimit.String())
		}
		// Preload StakingInfo from the current block and state. Needed for next block's engine.Finalize() post-Kaia.
		preloadedStakingBlockNums = append(preloadedStakingBlockNums, current.NumberU64())
		if err := reward.PreloadStakingInfoWithState(current.Header(), statedb); err != nil {
			return nil, fmt.Errorf("preloading staking info from block %d failed: %v", current.NumberU64(), err)
		}
		// Retrieve the next block to regenerate and process it
		next := current.NumberU64() + 1
		if current = cn.blockchain.GetBlockByNumber(next); current == nil {
			return nil, fmt.Errorf("block #%d not found", next)
		}
		_, _, _, _, _, err := cn.blockchain.Processor().Process(current, statedb, vm.Config{})
		if err != nil {
			return nil, fmt.Errorf("processing block %d failed: %v", current.NumberU64(), err)
		}
		// Finalize the state so any modifications are written to the trie
		root, err := statedb.Commit(true)
		if err != nil {
			return nil, err
		}
		if err := statedb.Reset(root); err != nil {
			return nil, fmt.Errorf("state reset after block %d failed: %v", current.NumberU64(), err)
		}
		database.TrieDB().ReferenceRoot(root)
		if !common.EmptyHash(parent) {
			database.TrieDB().Dereference(parent)
		}
		if current.Header().Root != root {
			err = fmt.Errorf("mistmatching state root block expected %x reexecuted %x", current.Header().Root, root)
			// Logging here because something went wrong when the state roots disagree even if the execution was successful.
			logger.Error("incorrectly regenerated historical state", "block", current.NumberU64(), "err", err)
			return nil, fmt.Errorf("incorrectly regenerated historical state for block %d: %v", current.NumberU64(), err)
		}
		parent = root
	}
	if report {
		nodes, _, imgs := database.TrieDB().Size()
		logger.Info("Historical state regenerated", "block", current.NumberU64(), "elapsed", time.Since(start), "nodes", nodes, "preimages", imgs)
	}

	return statedb, nil
}

// stateAtTransaction returns the execution environment of a certain transaction.
func (cn *CN) stateAtTransaction(block *types.Block, txIndex int, reexec uint64) (blockchain.Message, vm.BlockContext, vm.TxContext, *state.StateDB, error) {
	// Short circuit if it's genesis block.
	if block.NumberU64() == 0 {
		return nil, vm.BlockContext{}, vm.TxContext{}, nil, errors.New("no transaction in genesis")
	}
	// Create the parent state database
	parent := cn.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.BlockContext{}, vm.TxContext{}, nil, fmt.Errorf("parent %#x not found", block.ParentHash())
	}
	// Lookup the statedb of parent block from the live database,
	// otherwise regenerate it on the flight.
	statedb, err := cn.stateAtBlock(parent, reexec, nil, true, false)
	if err != nil {
		return nil, vm.BlockContext{}, vm.TxContext{}, nil, err
	}
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.BlockContext{}, vm.TxContext{}, statedb, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner(cn.blockchain.Config(), block.Number())
	for idx, tx := range block.Transactions() {
		// Assemble the transaction call message and return if the requested offset
		msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, block.NumberU64())
		if err != nil {
			logger.Warn("stateAtTransition failed", "hash", tx.Hash(), "block", block.NumberU64(), "err", err)
			return nil, vm.BlockContext{}, vm.TxContext{}, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}

		txContext := blockchain.NewEVMTxContext(msg, block.Header(), cn.chainConfig)
		blockContext := blockchain.NewEVMBlockContext(block.Header(), cn.blockchain, nil)
		if idx == txIndex {
			return msg, blockContext, txContext, statedb, nil
		}
		// Not yet the searched for transaction, execute on top of the current state
		vmenv := vm.NewEVM(blockContext, txContext, statedb, cn.blockchain.Config(), &vm.Config{})
		if _, err := blockchain.ApplyMessage(vmenv, msg); err != nil {
			return nil, vm.BlockContext{}, vm.TxContext{}, nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		// Ensure any modifications are committed to the state
		// Since Kaia is forked after EIP158/161 (a.k.a Spurious Dragon), deleting empty object is always effective
		statedb.Finalise(true, true)
	}
	return nil, vm.BlockContext{}, vm.TxContext{}, nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}
