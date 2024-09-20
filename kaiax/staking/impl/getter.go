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

package staking

import (
	"github.com/kaiachain/kaia/kaiax/staking"
)

func (s *StakingModule) GetStakingInfo(num uint64) (*staking.StakingInfo, error) {
	return nil, nil
}

func sourceBlockNum(num uint64, isKaia bool, interval uint64) uint64 {
	if isKaia {
		if num == 0 {
			return 0
		} else {
			return num - 1
		}
	} else {
		if num <= 2*interval {
			return 0
		} else {
			// Simplified from the previous implementation:
			// if (num % interval) == 0, return num - 2*interval
			// else return num - interval - (num % interval)
			return roundDown(num-1, interval) - interval
		}
	}
}

func roundDown(n, p uint64) uint64 {
	return n - (n % p)
}
