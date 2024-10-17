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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
)

func (v *ValsetModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   newValidatorAPI(v),
			Public:    true,
		},
		// TODO-kaiax-valset: add more namespaces. Istanbul? But it's not open to public. Anyway, move getDemoted to istanbul namespace.
	}
}

type ValidatorAPI struct {
	v *ValsetModule
}

// TODO-kaiax-valset: change the number parameter type to rpc.BlockNumberOrHash.
func newValidatorAPI(v *ValsetModule) *ValidatorAPI {
	return &ValidatorAPI{v}
}

// TODO-kaiax-valset: currently, this API returns the council list which is the result of block N-1. Is it correct?
func (api *ValidatorAPI) GetCouncil(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.v.chain, number)
	if err != nil {
		return nil, err
	}
	return api.v.GetCouncilAddressList(header.Number.Uint64() - 1)
}

func (api *ValidatorAPI) GetCouncilSize(number *rpc.BlockNumber) (int, error) {
	council, err := api.GetCouncil(number)
	if err != nil {
		return -1, err
	}
	return len(council), nil
}

func (api *ValidatorAPI) GetCommittee(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.v.chain, number)
	if err != nil {
		return nil, err
	}
	return api.v.GetCommitteeAddressList(header.Number.Uint64(), uint64(header.Round()))
}

func (api *ValidatorAPI) GetCommitteeSize(number *rpc.BlockNumber) (int, error) {
	committee, err := api.GetCommittee(number)
	if err != nil {
		return -1, err
	}
	return len(committee), nil
}

func (api *ValidatorAPI) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.v.chain, number)
	if err != nil {
		return nil, err
	}
	prevBlockResult, err := api.v.getBlockResultsByNumber(header.Number.Uint64() - 1)
	if err != nil {
		return nil, err
	}
	qualified, _ := splitByMinimumStakingAmount(prevBlockResult)
	return qualified.sortedAddressList(true), nil
}

// GetValidatorsAtHash retrieves the list of qualified validators with the given block hash.
func (api *ValidatorAPI) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.v.chain.GetHeaderByHash(hash)
	if header != nil {
		return nil, errNilHeader
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetValidators(&rpcBlockNumber)
}

func (api *ValidatorAPI) GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.v.chain, number)
	if err != nil {
		return nil, err
	}
	prevBlockResult, err := api.v.getBlockResultsByNumber(header.Number.Uint64() - 1)
	if err != nil {
		return nil, err
	}
	_, demoted := splitByMinimumStakingAmount(prevBlockResult)
	return demoted.sortedAddressList(true), nil
}

// GetDemotedValidatorsAtHash retrieves the list of demoted validators with the given block hash.
func (api *ValidatorAPI) GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.v.chain.GetHeaderByHash(hash)
	if header != nil {
		return nil, errNilHeader
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetDemotedValidators(&rpcBlockNumber)
}

// Retrieve the header at requested block number
func headerByRpcNumber(chain chain, number *rpc.BlockNumber) (*types.Header, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = chain.CurrentBlock().Header()
	} else if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get council list of the pending block.", "number", number)
		return nil, errPendingNotAllowed
	} else {
		header = chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return header, nil
}
