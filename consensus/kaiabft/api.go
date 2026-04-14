// Copyright 2026 The Kaia Authors
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

package kaiabft

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/networks/rpc"
)

// api exposes validator-query methods under the "istanbul" RPC namespace
// for backwards compatibility with tooling that expects istanbul_*.
type api struct {
	chain   consensus.ChainReader
	backend *backend
}

func (a *api) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := resolveRpcNumber(a.chain, number, true)
	if err != nil {
		return nil, err
	}
	if a.backend.valsetModule == nil {
		return nil, errNoModule
	}
	return a.backend.valsetModule.GetQualifiedValidators(num)
}

func (a *api) GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := resolveRpcNumber(a.chain, number, true)
	if err != nil {
		return nil, err
	}
	if a.backend.valsetModule == nil {
		return nil, errNoModule
	}
	return a.backend.valsetModule.GetDemotedValidators(num)
}

func (a *api) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := a.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, consensus.ErrUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return a.GetValidators(&rpcBlockNumber)
}

func (a *api) GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := a.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, consensus.ErrUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return a.GetDemotedValidators(&rpcBlockNumber)
}

func (a *api) GetTimeout() uint64 {
	return a.backend.timeout
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

	if num > headNum+1 {
		return 0, consensus.ErrUnknownBlock
	} else if num == headNum+1 && !allowPending {
		return 0, consensus.ErrUnknownBlock
	}
	return num, nil
}
