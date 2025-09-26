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

// Package faker implements a fake consensus engine for testing purposes.
// It accepts all blocks as valid, with optional failure conditions.
package faker

import (
	"errors"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

// ByzantiumBlockReward is the block reward for testing purposes.
var ByzantiumBlockReward = big.NewInt(3e+18)

// Faker is a consensus engine that accepts all blocks as valid.
type Faker struct {
	failBlock uint64        // Block number to fail at (0 = no failure)
	failDelay time.Duration // Delay before returning in Seal
	fullFake  bool          // Accept all blocks without validation
}

// NewFaker creates a fake consensus engine that accepts all blocks.
func NewFaker() *Faker {
	return &Faker{}
}

// NewFakeFailer creates a fake consensus engine that fails at a specific block.
func NewFakeFailer(fail uint64) *Faker {
	return &Faker{failBlock: fail}
}

// NewFakeDelayer creates a fake consensus engine that delays block sealing.
func NewFakeDelayer(delay time.Duration) *Faker {
	return &Faker{failDelay: delay}
}

// NewFullFaker creates a fake consensus engine that accepts all blocks without any validation.
func NewFullFaker() *Faker {
	return &Faker{fullFake: true}
}

// NewShared creates a shared fake consensus engine.
func NewShared() *Faker {
	return &Faker{}
}

// Author returns a fixed address for testing.
func (f *Faker) Author(header *types.Header) (common.Address, error) {
	return params.AuthorAddressForTesting, nil
}

// CanVerifyHeadersConcurrently returns true for concurrent verification.
func (f *Faker) CanVerifyHeadersConcurrently() bool {
	return true
}

// PreprocessHeaderVerification is not used for faker.
func (f *Faker) PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error) {
	panic("not implemented for faker")
}

// GetConsensusInfo returns empty consensus info.
func (f *Faker) GetConsensusInfo(block *types.Block) (consensus.ConsensusInfo, error) {
	return consensus.ConsensusInfo{}, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules.
func (f *Faker) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if f.fullFake {
		return nil
	}

	number := header.Number.Uint64()

	// Check if we should fail this block
	if f.failBlock != 0 && number == f.failBlock {
		return consensus.ErrUnknownAncestor
	}

	// Short circuit if the header is known
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	// For genesis block, skip parent check
	if number == 0 {
		return nil
	}

	// Check parent existence
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// All other headers are valid in fake mode
	return nil
}

// VerifyHeaders verifies a batch of headers concurrently.
func (f *Faker) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort, results := make(chan struct{}), make(chan error, len(headers))
	// If we're running a full engine faking, accept all headers as valid
	if f.fullFake || len(headers) == 0 {
		go func() {
			for i := 0; i < len(headers); i++ {
				results <- nil
			}
		}()
		return abort, results
	}
	go func() {
		for i := range headers {
			select {
			case <-abort:
				return
			default:
				err := f.verifyHeaderWorker(chain, headers, seals, i)
				results <- err
			}
		}
	}()
	return abort, results
}

// verifyHeaderWorker is a helper for batch header verification.
// Similar to gxhash, it checks if previous headers in the batch can serve as parents.
func (f *Faker) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	header := headers[index]
	number := header.Number.Uint64()

	// Check if we should fail this block
	if f.failBlock != 0 && number == f.failBlock {
		return consensus.ErrUnknownAncestor
	}

	// Short circuit if the header is known
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}

	// For genesis block, skip parent check
	if number == 0 {
		return nil
	}

	// Find parent - either from previous header in batch or from chain
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(header.ParentHash, number-1)
	} else if headers[index-1].Hash() == header.ParentHash {
		parent = headers[index-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}

	if parent == nil {
		return consensus.ErrUnknownAncestor
	}

	// All other headers are valid in fake mode
	return nil
}

// VerifySeal implements consensus.Engine, checking the seal validity.
func (f *Faker) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	if f.failBlock != 0 && header.Number.Uint64() == f.failBlock {
		return errors.New("seal verification failed")
	}
	return nil
}

// Prepare prepares the header for mining.
func (f *Faker) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// Handle nil Number field
	if header.Number == nil {
		return errors.New("header number is nil")
	}

	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	header.BlockScore = f.CalcBlockScore(chain, header.Time.Uint64(), parent)
	return nil
}

// Initialize runs any pre-transaction state modifications.
func (f *Faker) Initialize(chain consensus.ChainReader, header *types.Header, state *state.StateDB) {
	// No initialization needed for faker
}

// Finalize processes the block without modifying state.
func (f *Faker) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate block rewards
	f.accumulateRewards(state, header)

	// Set header fields
	header.Root = state.IntermediateRoot(true)
	header.ReceiptHash = types.DeriveSha(types.Receipts(receipts), header.Number)
	header.Bloom = types.CreateBloom(receipts)

	return types.NewBlock(header, txs, receipts), nil
}

// accumulateRewards credits the coinbase of the given block with the mining reward.
func (f *Faker) accumulateRewards(state *state.StateDB, header *types.Header) {
	reward := ByzantiumBlockReward
	state.AddBalance(params.AuthorAddressForTesting, reward)
}

// Seal generates a new sealing request for the given block.
func (f *Faker) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	// Add delay if configured
	if f.failDelay > 0 {
		select {
		case <-time.After(f.failDelay):
		case <-stop:
			return nil, nil
		}
	}

	// Check if we should fail
	if f.failBlock != 0 && block.NumberU64() == f.failBlock {
		return nil, errors.New("seal failed")
	}

	return block, nil
}

// CalcBlockScore calculates the block score.
func (f *Faker) CalcBlockScore(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(1)
}

// APIs returns RPC APIs (none for faker).
func (f *Faker) APIs(chain consensus.ChainReader) []rpc.API {
	return nil
}

// Protocol returns the consensus protocol.
func (f *Faker) Protocol() consensus.Protocol {
	return consensus.KaiaProtocol
}

// PurgeCache is a no-op for faker.
func (f *Faker) PurgeCache() {
	// No cache to purge for faker
}
