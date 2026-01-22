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
// MERCHANTABILITY or FITNESS FOR A PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package consensus

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// CommitteeStateProviderInterface defines the interface for committee state providers
type CommitteeStateProviderInterface interface {
	GetValidatorSet(num uint64) (*istanbul.BlockValSet, error)
	GetCommitteeStateByRound(num uint64, round uint64) (*istanbul.RoundCommitteeState, error)
	GetProposerByRound(num uint64, round uint64) (common.Address, error)
}

// CommitteeStateProvider provides committee state related functionality
type CommitteeStateProvider struct {
	valsetModule valset.ValsetModule
	govModule    gov.GovModule
}

// NewCommitteeStateProvider creates a new CommitteeStateProvider
func NewCommitteeStateProvider(valsetModule valset.ValsetModule, govModule gov.GovModule) *CommitteeStateProvider {
	return &CommitteeStateProvider{
		valsetModule: valsetModule,
		govModule:    govModule,
	}
}

// GetValidatorSet returns the validator set for the given block number
func (csp *CommitteeStateProvider) GetValidatorSet(num uint64) (*istanbul.BlockValSet, error) {
	council, err := csp.valsetModule.GetCouncil(num)
	if err != nil {
		return nil, err
	}

	demoted, err := csp.valsetModule.GetDemotedValidators(num)
	if err != nil {
		return nil, err
	}

	return istanbul.NewBlockValSet(council, demoted), nil
}

// GetCommitteeStateByRound returns the committee state for the given block number and round
func (csp *CommitteeStateProvider) GetCommitteeStateByRound(num uint64, round uint64) (*istanbul.RoundCommitteeState, error) {
	blockValSet, err := csp.GetValidatorSet(num)
	if err != nil {
		return nil, err
	}

	committee, err := csp.valsetModule.GetCommittee(num, round)
	if err != nil {
		return nil, err
	}

	proposer, err := csp.valsetModule.GetProposer(num, round)
	if err != nil {
		return nil, err
	}

	committeeSize := csp.govModule.GetParamSet(num).CommitteeSize
	return istanbul.NewRoundCommitteeState(blockValSet, committeeSize, committee, proposer), nil
}

// GetProposerByRound returns the proposer address for the given block number and round
func (csp *CommitteeStateProvider) GetProposerByRound(num uint64, round uint64) (common.Address, error) {
	proposer, err := csp.valsetModule.GetProposer(num, round)
	if err != nil {
		return common.Address{}, err
	}
	return proposer, nil
}
