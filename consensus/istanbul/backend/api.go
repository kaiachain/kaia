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
// This file is derived from quorum/consensus/istanbul/backend/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/networks/rpc"
)

// API is a user facing RPC API to dump Istanbul state
type API struct {
	chain    consensus.ChainReader
	istanbul *backend
}

// GetValidators retrieves the list of qualified validators with the given block number.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := resolveRpcNumber(api.chain, number, true)
	if err != nil {
		return nil, err
	}

	valSet, err := api.istanbul.GetValidatorSet(num)
	if err != nil {
		return nil, err
	}
	return valSet.Qualified().List(), nil
}

// GetDemotedValidators retrieves the list of authorized, but demoted validators with the given block number.
func (api *API) GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := resolveRpcNumber(api.chain, number, true)
	if err != nil {
		return nil, err
	}

	valSet, err := api.istanbul.GetValidatorSet(num)
	if err != nil {
		return nil, err
	}
	return valSet.Demoted().List(), nil
}

// GetValidatorsAtHash retrieves the list of authorized validators with the given block hash.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetValidators(&rpcBlockNumber)
}

// GetDemotedValidatorsAtHash retrieves the list of demoted validators with the given block hash.
func (api *API) GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetDemotedValidators(&rpcBlockNumber)
}

// Candidates returns the current candidates the node tries to uphold and vote on.
func (api *API) Candidates() map[common.Address]bool {
	api.istanbul.candidatesLock.RLock()
	defer api.istanbul.candidatesLock.RUnlock()

	proposals := make(map[common.Address]bool)
	for address, auth := range api.istanbul.candidates {
		proposals[address] = auth
	}
	return proposals
}

// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	api.istanbul.candidates[address] = auth
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	delete(api.istanbul.candidates, address)
}

func RecoverCommittedSeals(extra *types.IstanbulExtra, headerHash common.Hash) ([]common.Address, error) {
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

func (api *API) GetTimeout() uint64 {
	return istanbul.DefaultConfig.Timeout
}

func resolveRpcNumber(chain consensus.ChainReader, number *rpc.BlockNumber, allowPending bool) (uint64, error) {
	headNum := chain.CurrentHeader().Number.Uint64()
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
