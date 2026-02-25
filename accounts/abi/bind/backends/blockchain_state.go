// Modifications Copyright 2026 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package backends

import (
	"context"
	"math/big"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/node/cn/filters"
)

// `StateBlockchainContractBackend` has the same internal behavior with BlockchainContractBackend
// except for the `ContractCaller` interface functions i.e., `CodeAt()` and `CallContract()`
type StateBlockchainContractBackend struct {
	state *state.StateDB
	*BlockchainContractCommonBackend
}

var (
	_ bind.ContractCaller     = (*StateBlockchainContractBackend)(nil)
	_ bind.ContractTransactor = (*StateBlockchainContractBackend)(nil)
	_ bind.ContractFilterer   = (*StateBlockchainContractBackend)(nil)
	_ bind.DeployBackend      = (*StateBlockchainContractBackend)(nil)
	_ bind.ContractBackend    = (*StateBlockchainContractBackend)(nil)
)

func NewStateBlockchainContractBackend(bc BlockChainForCaller, tp TxPoolForCaller, es *filters.EventSystem, state *state.StateDB) *StateBlockchainContractBackend {
	bcCommonBackend := newBlockchainContractCommonBackend(bc, tp, es)
	return &StateBlockchainContractBackend{state, bcCommonBackend}
}

func (b *StateBlockchainContractBackend) getBlock(num *big.Int) (*types.Block, error) {
	var block *types.Block
	if num == nil {
		block = b.bc.CurrentBlock()
	} else {
		header := b.bc.GetHeaderByNumber(num.Uint64())
		if header == nil {
			return nil, errBlockDoesNotExist
		}
		block = b.bc.GetBlock(header.Hash(), header.Number.Uint64())
	}
	if block == nil {
		return nil, errBlockDoesNotExist
	}
	return block, nil
}

// bind.ContractCaller defined methods

func (b *StateBlockchainContractBackend) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.BlockchainContractCommonBackend.CodeAt(ctx, b.state, account, blockNumber)
}

func (b *StateBlockchainContractBackend) CallContract(ctx context.Context, call kaia.CallMsg, blockNumber *big.Int) ([]byte, error) {
	block, err := b.getBlock(blockNumber)
	if err != nil {
		return nil, err
	}
	return b.BlockchainContractCommonBackend.CallContract(ctx, b.state, block, call, blockNumber)
}
