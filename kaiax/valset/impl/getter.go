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
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

// GetCouncil returns the whole validator list for validating the block `num`.
func (v *ValsetModule) GetCouncil(num uint64) ([]common.Address, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	} else {
		return council.List(), nil
	}
}

// GetDemotedValidators are subtract of qualified from council(N)
func (v *ValsetModule) GetDemotedValidators(num uint64) ([]common.Address, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	demoted, err := v.getDemotedValidators(council, num)
	if err != nil {
		return nil, err
	}
	return demoted.List(), nil
}

func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	council, err := v.getCouncil(num)
	if err != nil {
		return nil, err
	}
	demoted, err := v.getDemotedValidators(council, num)
	if err != nil {
		return nil, err
	}
	return council.Subtract(demoted), nil
}

// GetCommittee returns the current block's committee.
func (v *ValsetModule) GetCommittee(num uint64, round uint64) ([]common.Address, error) {
	if num == 0 {
		return v.GetCouncil(0)
	}

	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return nil, err
	}
	return v.getCommittee(c, round)
}

func (v *ValsetModule) GetProposer(num, round uint64) (common.Address, error) {
	if num == 0 {
		return common.Address{}, nil
	}
	if header := v.Chain.GetHeaderByNumber(num); header != nil {
		if uint64(header.Round()) == round {
			return v.Chain.Engine().Author(header)
		}
	}
	// TODO-kaiax: Sync blockContext
	c, err := v.getBlockContext(num)
	if err != nil {
		return common.Address{}, err
	}
	return v.getProposer(c, round)
}

// GetValidatorSet returns the validator set for the given block number.
func (v *ValsetModule) GetValidatorSet(num uint64) (*istanbul.BlockValSet, error) {
	council, err := v.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	demoted, err := v.GetDemotedValidators(num)
	if err != nil {
		return nil, err
	}

	// Use the constructor function
	return istanbul.NewBlockValSet(council, demoted), nil
}

// GetConsensusInfo returns consensus information regarding the given block number.
func (v *ValsetModule) GetConsensusInfo(block *types.Block) (consensus.ConsensusInfo, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return consensus.ConsensusInfo{}, nil
	}

	if v.Chain == nil {
		return consensus.ConsensusInfo{}, errNoChainReader
	}

	// Get the committers of this block from committed seals
	extra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return consensus.ConsensusInfo{}, err
	}
	committers, err := v.recoverCommittedSeals(extra, block.Hash())
	if err != nil {
		return consensus.ConsensusInfo{}, err
	}

	round := block.Header().Round()
	// Get the committee list of this block (blockNumber, round)
	currentProposer, err := v.GetProposer(blockNumber, uint64(round))
	if err != nil {
		log.NewModuleLogger(log.ConsensusIstanbulBackend).Error("Failed to get proposer.", "blockNum", blockNumber, "round", uint64(round), "err", err)
		return consensus.ConsensusInfo{}, errInternalError
	}

	currentCommittee, err := v.GetCommittee(blockNumber, uint64(round))
	if err != nil {
		currentCommittee = []common.Address{}
	}

	// Get origin proposer at 0 round.
	var roundZeroProposer *common.Address
	roundZeroCommittee, err := v.GetCommittee(blockNumber, 0)
	if err == nil && len(roundZeroCommittee) > 0 {
		roundZeroProposer = &roundZeroCommittee[0]
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

// recoverCommittedSeals recovers the addresses of validators that committed seals for the given block.
func (v *ValsetModule) recoverCommittedSeals(extra *types.IstanbulExtra, headerHash common.Hash) ([]common.Address, error) {
	committers := make([]common.Address, len(extra.CommittedSeal))
	for idx, cs := range extra.CommittedSeal {
		committer, err := istanbul.GetSignatureAddress(istanbulCore.PrepareCommittedSeal(headerHash), cs)
		if err != nil {
			return nil, err
		}
		committers[idx] = committer
	}
	return committers, nil
}

// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
//
// This function is derived from consensus/istanbul/backend/engine.go:sigHash
// to avoid import cycle between valset and istanbul backend packages.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// CurrentHeader returns the current header of the chain.
func (v *ValsetModule) CurrentHeader() *types.Header {
	return v.Chain.CurrentBlock().Header()
}

// NewBlockchainContractBackend creates a new blockchain contract backend.
func (v *ValsetModule) NewBlockchainContractBackend() *backends.BlockchainContractBackend {
	if chain, ok := v.Chain.(backends.BlockChainForCaller); ok {
		return backends.NewBlockchainContractBackend(chain, nil, nil)
	}
	return nil
}

// ResolveRpcNumber resolves the RPC block number to a uint64.
func (v *ValsetModule) ResolveRpcNumber(number *rpc.BlockNumber, allowPending bool) (uint64, error) {
	headNum := v.CurrentHeader().Number.Uint64()
	var num uint64
	if number == nil || *number == rpc.LatestBlockNumber {
		num = headNum
	} else if *number == rpc.PendingBlockNumber {
		num = headNum + 1
	} else {
		num = uint64(number.Int64())
	}

	if num > headNum+1 { // May allow up to head + 1 to query the pending block.
		return 0, errUnknownBlock
	} else if num == headNum+1 && !allowPending {
		return 0, errPendingNotAllowed
	} else {
		return num, nil
	}
}
