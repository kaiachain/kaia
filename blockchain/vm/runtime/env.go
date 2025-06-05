// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/vm/runtime/env.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package runtime

import (
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/vm"
)

func NewEnv(cfg *Config) *vm.EVM {
	txContext := vm.TxContext{
		Origin:   cfg.Origin,
		GasPrice: cfg.GasPrice,
	}
	blockContext := vm.BlockContext{
		CanTransfer: blockchain.CanTransfer,
		Transfer:    blockchain.Transfer,
		GetHash:     cfg.GetHashFn,

		Coinbase:    cfg.Coinbase,
		Rewardbase:  cfg.Rewardbase,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		BlockScore:  cfg.BlockScore,
		GasLimit:    cfg.GasLimit,
		BaseFee:     cfg.BaseFee,
	}
	return vm.NewEVM(blockContext, txContext, cfg.State, cfg.ChainConfig, &cfg.EVMConfig)
}
