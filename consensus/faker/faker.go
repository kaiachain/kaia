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
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/networks/p2p"
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
	address   common.Address
	eventMux  *event.TypeMux
}

// NewFaker creates a fake consensus engine that accepts all blocks.
func NewFaker() *Faker {
	return &Faker{address: params.AuthorAddressForTesting, eventMux: new(event.TypeMux)}
}

// NewFakerWithFixedSealer creates a fake consensus engine with a fixed proposer/sealer address.
func NewFakerWithFixedSealer(address common.Address) *Faker {
	return &Faker{address: address, eventMux: new(event.TypeMux)}
}

// NewFakeFailer creates a fake consensus engine that fails at a specific block.
func NewFakeFailer(fail uint64) *Faker {
	return &Faker{failBlock: fail, address: params.AuthorAddressForTesting, eventMux: new(event.TypeMux)}
}

// NewFakeDelayer creates a fake consensus engine that delays block sealing.
func NewFakeDelayer(delay time.Duration) *Faker {
	return &Faker{failDelay: delay, address: params.AuthorAddressForTesting, eventMux: new(event.TypeMux)}
}

// NewFullFaker creates a fake consensus engine that accepts all blocks without any validation.
func NewFullFaker() *Faker {
	return &Faker{fullFake: true, address: params.AuthorAddressForTesting, eventMux: new(event.TypeMux)}
}

// NewShared creates a shared fake consensus engine.
func NewShared() *Faker {
	return &Faker{address: params.AuthorAddressForTesting, eventMux: new(event.TypeMux)}
}

// Author returns a fixed address for testing.
func (f *Faker) Author(header *types.Header) (common.Address, error) {
	return f.address, nil
}

func (f *Faker) MakeCommittedSeal(header *types.Header) ([]byte, error) {
	return nil, f.verifyFailure(header)
}

func (f *Faker) MakeAuthorSeal(header *types.Header) ([]byte, error) {
	return nil, f.verifyFailure(header)
}

func (f *Faker) Validators(header *types.Header) ([]common.Address, error) {
	if err := f.verifyFailure(header); err != nil {
		return nil, err
	}
	return []common.Address{f.address}, nil
}

func (f *Faker) Round(header *types.Header) (byte, error) {
	return 0, f.verifyFailure(header)
}

func (f *Faker) RawSeals(header *types.Header) ([]byte, [][]byte, error) {
	return nil, nil, f.verifyFailure(header)
}

func (f *Faker) WriteRound(header *types.Header, round int64) {
	if header == nil {
		return
	}
	if len(header.Extra) >= 32 {
		header.Extra[31] = byte(round)
	}
}

// Start is a no-op for faker.
func (f *Faker) Start(chain consensus.ChainReader, executor consensus.Executor) error {
	return nil
}

// Stop is a no-op for faker.
func (f *Faker) Stop() error {
	return nil
}

// Sealer returns faker itself as a sealer implementation.
func (f *Faker) Sealer() consensus.Sealer { return f }

// RegisterKaiaxModules is a no-op for faker.
func (f *Faker) RegisterKaiaxModules(mGov gov.GovModule, mValset valset.ValsetModule) {}

// RegisterVRankModule is a no-op for faker.
func (f *Faker) RegisterVRankModule(mVRank vrank.VRankModule) {}

func (f *Faker) SigHash(_ *types.Header) common.Hash {
	return common.Hash{}
}

func (f *Faker) HeaderHash(_ *types.Header) common.Hash {
	return common.Hash{}
}

// WriteValidators is a no-op for faker.
func (f *Faker) WriteValidators(header *types.Header, _ []common.Address) error {
	if header == nil {
		return nil
	}
	return f.verifyFailure(header)
}

// WriteAuthorSeal is a no-op for faker.
func (f *Faker) WriteAuthorSeal(header *types.Header, _ []byte) error {
	return f.verifyFailure(header)
}

// WriteCommittedSeals is a no-op for faker.
func (f *Faker) WriteCommittedSeals(header *types.Header, _ [][]byte) error {
	return f.verifyFailure(header)
}

func (f *Faker) verifyFailure(header *types.Header) error {
	if f.fullFake {
		return nil
	}
	if header == nil || header.Number == nil {
		return consensus.ErrUnknownBlock
	}
	if f.failBlock != 0 && header.Number.Uint64() == f.failBlock {
		return consensus.ErrUnknownAncestor
	}
	return nil
}

// Committers checks preprocess-time seal extraction for faker.
func (f *Faker) Committers(header *types.Header) ([]common.Address, error) {
	if err := f.verifyFailure(header); err != nil {
		return nil, err
	}
	return []common.Address{f.address}, nil
}

func (f *Faker) Vanity(header *types.Header) ([]byte, error) {
	if err := f.verifyFailure(header); err != nil {
		return nil, err
	}
	vanity := make([]byte, 32)
	copy(vanity, header.Extra)
	return vanity, nil
}

// VerifySeals checks consensus-specific seals for faker.
func (f *Faker) VerifySeals(header *types.Header) error {
	return f.verifyFailure(header)
}

// F returns the number of tolerated Byzantine faults for faker.
func (f *Faker) F(blockNum uint64, qualifiedlen, committeeSize int) int {
	return 0
}

// Quorum returns the quorum threshold for faker.
func (f *Faker) Quorum(blockNum uint64, qualifiedlen, committeeSize int) int {
	return 0
}

// NewChainHead is a no-op for faker.
func (f *Faker) NewChainHead() error {
	return nil
}

// HandleMsg is a no-op for faker.
func (f *Faker) HandleMsg(_ common.Address, _ p2p.Msg) (bool, error) {
	return false, nil
}

// SetBroadcaster is a no-op for faker.
func (f *Faker) SetBroadcaster(consensus.Broadcaster) {}

// Prepare prepares the header for mining.
func (f *Faker) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// Handle nil Number field
	if header.Number == nil {
		return errors.New("header number is nil")
	}

	header.BlockScore = big.NewInt(1) // Simplified difficulty for testing
	return nil
}

// Finalize processes the block without modifying state.
func (f *Faker) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate block rewards
	f.accumulateRewards(state, header)

	// Set header fields
	header.Root = state.IntermediateRoot(true)
	header.ReceiptHash = types.DeriveReceiptsRoot(types.Receipts(receipts), header.Number)
	header.Bloom = types.CreateBloom(receipts)

	return types.NewBlock(header, txs, receipts), nil
}

// accumulateRewards credits the coinbase of the given block with the mining reward.
func (f *Faker) accumulateRewards(state *state.StateDB, header *types.Header) {
	reward := ByzantiumBlockReward
	state.AddBalance(f.address, reward)
}

// Seal generates a new sealing request for the given block.
func (f *Faker) Seal(chain consensus.ChainReader, block *types.Block) (*types.Block, error) {
	// Add delay if configured
	if f.failDelay > 0 {
		time.Sleep(f.failDelay)
	}

	// Check if we should fail
	if f.failBlock != 0 && block.NumberU64() == f.failBlock {
		return nil, errors.New("seal failed")
	}

	return block, nil
}

// APIs returns RPC APIs (none for faker).
func (f *Faker) APIs(chain consensus.ChainReader) []rpc.API {
	return nil
}

// PurgeCache is a no-op for faker.
func (f *Faker) PurgeCache() {
	// No cache to purge for faker
}

// SubmitTransactions is not implemented for faker - it uses push() flow instead.
func (f *Faker) SubmitTransactions(txs *types.TransactionsByPriceAndNonce, state *state.StateDB, header *types.Header, mux *event.TypeMux, onPrepared func(*consensus.ExecutionResult)) (finalizeCh <-chan *consensus.ExecutionResult) {
	// Return nil channel - faker uses the legacy push() flow
	return nil
}

// fakerSequenceEvent is a dummy event for Faker's SubscribeNewSequence.
type fakerSequenceEvent struct{}

// SubscribeNewSequence returns a subscription that never receives events.
func (f *Faker) SubscribeNewSequence() *event.TypeMuxSubscription {
	return f.eventMux.Subscribe(fakerSequenceEvent{})
}
