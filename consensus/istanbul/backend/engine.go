// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/backend/engine.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"encoding/hex"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/rpc"
)

const (
	inmemoryPeers    = 200
	inmemoryMessages = 4096
)

var (
	inmemoryBlocks             = 2048 // Number of blocks to precompute validators' addresses
	inmemoryValidatorsPerBlock = 30   // Approximate number of validators' addresses from ecrecover
	signatureAddresses, _      = lru.NewARC(inmemoryBlocks * inmemoryValidatorsPerBlock)
)

// cacheSignatureAddresses extracts the address from the given data and signature and cache them for later usage.
func cacheSignatureAddresses(data []byte, sig []byte) (common.Address, error) {
	sigStr := hex.EncodeToString(sig)
	if addr, ok := signatureAddresses.Get(sigStr); ok {
		return addr.(common.Address), nil
	}
	addr, err := istanbul.GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	signatureAddresses.Add(sigStr, addr)
	return addr, err
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()
	block, err := sb.updateBlock(block)
	if err != nil {
		return nil, err
	}

	// Initialize seal state; returns nil if block was already committed before Seal started
	commitCh := sb.initSealState(number, block.Hash())
	if commitCh == nil {
		return nil, nil
	}
	defer sb.cleanupSealState()

	// post block into Istanbul engine
	go sb.EventMux().Post(istanbul.RequestEvent{
		Proposal: block,
	})

	logger.Debug("[Seal] Waiting for commitCh", "blockNum", number)

	for {
		select {
		case result := <-commitCh:
			logger.Debug("[Seal] Received from commitCh", "blockNum", number, "resultNil", result == nil)
			if result == nil {
				return nil, nil
			}
			// if the block hash and the hash from channel are the same,
			// return the result. Otherwise, keep waiting the next hash.
			header := block.Header()
			sb.sealer.WriteRound(header, result.Round)
			block = block.WithSeal(header)
			if block.Hash() == result.Block.Hash() {
				return result.Block, nil
			}
		}
	}
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	if header != nil && header.Number != nil && header.Number.Uint64() > 0 {
		parent := sb.chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
		if parent == nil {
			return nil, consensus.ErrUnknownAncestor
		}
		if parent.Hash() != header.ParentHash || parent.Number.Uint64()+1 != header.Number.Uint64() {
			return nil, consensus.ErrUnknownAncestor
		}
	}
	if sb.valsetModule == nil {
		return nil, istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	if !valset.NewAddressSet(qualified).Contains(sb.address) {
		return nil, consensus.ErrUnauthorized
	}
	authorSeal, err := sb.sealer.MakeAuthorSeal(header)
	if err != nil {
		return nil, err
	}
	if err := sb.sealer.WriteAuthorSeal(header, authorSeal); err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "istanbul",
			Version:   "1.0",
			Service:   &API{chain: chain, istanbul: sb},
			Public:    true,
		},
	}
}

// SetChain sets chain of the Istanbul backend
func (sb *backend) SetChain(chain consensus.ChainReader) {
	sb.chain = chain
}

// SignalPeerRegistrable unblocks ValidatePeerType so that peers can be registered.
// It is safe to call multiple times; only the first call has effect (sync.Once).
//
// Call sites:
//   - Mining nodes: called automatically via defer in Start(), after coreStarted=true.
//     Also covers Start() failure paths to prevent goroutine leaks.
//   - WorkerDisable nodes: caller must invoke this explicitly after SetChain(),
//     since Start() is never called.
//   - Stop(): called to unblock any peers waiting during shutdown.
func (sb *backend) SignalPeerRegistrable() {
	sb.chainInitOnce.Do(func() { close(sb.chainInitCh) })
}

// RegisterKaiaxModules sets kaiax modules of the Istanbul backend
func (sb *backend) RegisterKaiaxModules(mGov gov.GovModule, mValset valset.ValsetModule) {
	sb.valsetModule = mValset
	sb.core.RegisterKaiaxModules(mValset, mGov)
}

// Start starts the Istanbul backend core.
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool, executor consensus.Executor) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return istanbul.ErrStartedEngine
	}

	// clear previous data
	sb.sealMu.Lock()
	sb.proposedBlockHash = common.Hash{}
	sb.sealSkippedNum = 0
	if sb.commitCh != nil {
		close(sb.commitCh)
		sb.commitCh = nil
	}
	sb.sealMu.Unlock()

	// Ensure peers are unblocked even if Start() fails partway through.
	// This defer is intentionally placed before SetChain: if Start() fails
	// before SetChain, sb.chain remains nil, which is safely handled by
	// the nil guard in ValidatePeerType.
	defer sb.SignalPeerRegistrable()

	sb.SetChain(chain)
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock
	sb.executor = executor

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop stops the Istanbul backend core.
func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	// Unblock any peers waiting in ValidatePeerType during shutdown.
	sb.SignalPeerRegistrable()
	if !sb.coreStarted {
		return istanbul.ErrStoppedEngine
	}
	// Close commitCh to stop any pending Seal() calls
	sb.sealMu.Lock()
	if sb.commitCh != nil {
		close(sb.commitCh)
		sb.commitCh = nil
	}
	sb.sealMu.Unlock()
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

func (sb *backend) PurgeCache() {
	// TODO-kaiax: Implement this
}
