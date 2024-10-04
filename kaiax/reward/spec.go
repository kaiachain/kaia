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
	"math/big"

	"github.com/kaiachain/kaia/common"
)

type RewardSummary struct {
	Minted   *big.Int `json:"minted"`
	TotalFee *big.Int `json:"totalFee"`
	BurntFee *big.Int `json:"burntFee"`
}

type RewardSpec struct {
	RewardSummary
	Proposer *big.Int                    `json:"proposer"`
	Stakers  *big.Int                    `json:"stakers"`
	KIF      *big.Int                    `json:"kif"`
	KEF      *big.Int                    `json:"kef"`
	Rewards  map[common.Address]*big.Int `json:"rewards"`
}

func NewRewardSpec() *RewardSpec {
	return &RewardSpec{
		RewardSummary: RewardSummary{
			Minted:   big.NewInt(0),
			TotalFee: big.NewInt(0),
			BurntFee: big.NewInt(0),
		},
		Proposer: big.NewInt(0),
		Stakers:  big.NewInt(0),
		KIF:      big.NewInt(0),
		KEF:      big.NewInt(0),
		Rewards:  make(map[common.Address]*big.Int),
	}
}

func (spec *RewardSpec) Add(delta *RewardSpec) {
	spec.Minted.Add(spec.Minted, delta.Minted)
	spec.TotalFee.Add(spec.TotalFee, delta.TotalFee)
	spec.BurntFee.Add(spec.BurntFee, delta.BurntFee)
	spec.Proposer.Add(spec.Proposer, delta.Proposer)
	spec.Stakers.Add(spec.Stakers, delta.Stakers)
	spec.KIF.Add(spec.KIF, delta.KIF)
	spec.KEF.Add(spec.KEF, delta.KEF)

	for addr, amount := range delta.Rewards {
		incrementRewardsMap(spec.Rewards, addr, amount)
	}
}

func incrementRewardsMap(m map[common.Address]*big.Int, addr common.Address, amount *big.Int) {
	_, ok := m[addr]
	if !ok {
		m[addr] = big.NewInt(0)
	}
	m[addr] = m[addr].Add(m[addr], amount)
}
