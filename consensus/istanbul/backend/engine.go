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
	"bytes"
	"encoding/hex"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	consensuscommon "github.com/kaiachain/kaia/consensus/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
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

// Author retrieves/caches the Kaia address of the account that minted the given block.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := cacheSignatureAddresses(sigHash(header).Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}
	return addr, nil
}

// computeSignatureAddrs extracts/caches signer and committer addresses from header seals.
func (sb *backend) Committers(header *types.Header) ([]common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	// The length of Committed seals should be larger than 0
	if len(istanbulExtra.CommittedSeal) == 0 {
		return nil, istanbul.ErrEmptyCommittedSeals
	}

	committers := make([]common.Address, 0, len(istanbulExtra.CommittedSeal))
	proposalSeal := istanbulCore.PrepareCommittedSeal(header.Hash())
	for _, seal := range istanbulExtra.CommittedSeal {
		addr, err := cacheSignatureAddresses(proposalSeal, seal)
		if err != nil {
			return nil, istanbul.ErrInvalidSignature
		}
		committers = append(committers, addr)
	}
	return committers, nil
}

func (sb *backend) VerifySeals(header *types.Header, sigCacheMode bool) error {
	if header.Number == nil {
		return consensus.ErrUnknownBlock
	}
	if header.Number.Uint64() == 0 {
		return nil
	}
	if _, err := types.ExtractIstanbulExtra(header); err != nil {
		return istanbul.ErrInvalidExtraDataFormat
	}
	number := header.Number.Uint64()
	signer, err := sb.Author(header)
	if err != nil {
		return err
	}
	committers, err := sb.Committers(header)
	if err != nil {
		return err
	}

	if sigCacheMode {
		return nil
	}

	if sb.valsetModule == nil || sb.govModule == nil {
		return istanbul.ErrNoEssentialModule
	}

	// check whether the signer is in the validator set.
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return err
	}
	if !valset.NewAddressSet(qualified).Contains(signer) {
		return istanbul.ErrUnauthorized
	}

	// Retrieve the snapshot needed to verify this header.
	council, err := sb.valsetModule.GetCouncil(number)
	if err != nil {
		return err
	}

	// Every validator can have only one seal. If more than one seals are signed by a
	// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
	councilSet := valset.NewAddressSet(council).Copy()
	validSeal := 0
	for _, addr := range committers {
		if councilSet.Remove(addr) {
			validSeal++
		} else {
			return istanbul.ErrInvalidCommittedSeals
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	committeeSize := sb.govModule.GetParamSet(number).CommitteeSize
	f := consensuscommon.CalcFaultTolerance(len(qualified), committeeSize)
	if validSeal <= 2*f {
		return istanbul.ErrInvalidCommittedSeals
	}
	return nil
}

// PrepareExtra builds Istanbul extra-data validators section for the given header.
func (sb *backend) PrepareExtra(header *types.Header, _ *types.Header) ([]byte, error) {
	if sb.valsetModule == nil {
		return nil, istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return prepareExtra(header, qualified)
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	// Bail out if we're unauthorized to sign a block
	if sb.valsetModule == nil {
		return nil, istanbul.ErrNoEssentialModule
	}
	qualified, err := sb.valsetModule.GetQualifiedValidators(number)
	if err != nil {
		return nil, err
	}
	if !valset.NewAddressSet(qualified).Contains(sb.address) {
		return nil, istanbul.ErrUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	block, err = sb.updateBlock(block)
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
			block = types.SetRoundToBlock(block, result.Round)
			if block.Hash() == result.Block.Hash() {
				return result.Block, nil
			}
		}
	}
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(block *types.Block) (*types.Block, error) {
	header := block.Header()
	// sign the hash
	seal, err := sb.Sign(sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
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

// RegisterKaiaxModules sets kaiax modules of the Istanbul backend
func (sb *backend) RegisterKaiaxModules(mGov gov.GovModule, mValset valset.ValsetModule) {
	sb.govModule = mGov
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

// GetConsensusInfo returns consensus information regarding the given block number.
func (sb *backend) GetConsensusInfo(block *types.Block) (consensus.ConsensusInfo, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return consensus.ConsensusInfo{}, nil
	}

	if sb.chain == nil {
		return consensus.ConsensusInfo{}, errNoChainReader
	}

	// get the committers of this block from committed seals
	committers, err := sb.Committers(block.Header())
	if err != nil {
		if err != istanbul.ErrEmptyCommittedSeals {
			return consensus.ConsensusInfo{}, err
		}
		committers = []common.Address{}
	}

	round := block.Header().Round()
	// get the committee list of this block (blockNumber, round)
	currentProposer, err := sb.valsetModule.GetProposer(blockNumber, uint64(round))
	if err != nil {
		logger.Error("Failed to get proposer.", "blockNum", blockNumber, "round", uint64(round), "err", err)
		return consensus.ConsensusInfo{}, istanbul.ErrInternalError
	}

	var currentCommittee []common.Address
	if sb.valsetModule != nil {
		currentCommittee, _ = sb.valsetModule.GetCommittee(block.NumberU64(), uint64(round))
	}

	// Uncomment to validate if committers are in the committee
	// for _, recovered := range committers {
	// 	found := false
	// 	for _, calculated := range currentCommittee {
	// 		if recovered == calculated {
	// 			found = true
	// 		}
	// 	}
	// 	if !found {
	// 		return consensus.ConsensusInfo{}, errInvalidCommittedSeals
	// 	}
	// }

	// get origin proposer at 0 round.
	var roundZeroProposer *common.Address
	if sb.valsetModule != nil {
		addr, err := sb.valsetModule.GetProposer(blockNumber, 0)
		if err == nil {
			roundZeroProposer = &addr
		}
	}

	cInfo := consensus.ConsensusInfo{
		SigHash:        sigHash(block.Header()),
		Proposer:       currentProposer,
		OriginProposer: roundZeroProposer,
		Committee:      currentCommittee,
		Committers:     committers,
		Round:          round,
	}

	return cInfo, nil
}

func (sb *backend) PurgeCache() {
	// TODO-kaiax: Implement this
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.IstanbulExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.IstanbulExtraVanity])

	ist := &types.IstanbulExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.IstanbulExtraSeal != 0 {
		return istanbul.ErrInvalidSignature
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return istanbul.ErrInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.IstanbulExtraSeal {
			return istanbul.ErrInvalidCommittedSeals
		}
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(istanbulExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}
