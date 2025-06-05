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

package reward

import (
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/kaiax"
)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/v2/kaiax/reward RewardModule
type RewardModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.TxProcessModule

	// GetDeferredReward returns the deferred reward specification to be distributed at the given block that is being created.
	GetDeferredReward(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) (*RewardSpec, error)

	// GetBlockReward retrospectively calculates the block reward distributed at the given block number.
	GetBlockReward(num uint64) (*RewardSpec, error)

	// GetRewardSummary retrospectively calculates the reward summary at the given block number.
	GetRewardSummary(num uint64) (*RewardSummary, error)
}
