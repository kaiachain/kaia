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
	"encoding/json"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/storage/database"
)

var stakingInfoPrefix = []byte("stakingInfo")

func stakingInfoKey(num uint64) []byte {
	return append(stakingInfoPrefix, common.Int64ToByteLittleEndian(num)...)
}

func ReadStakingInfo(db database.Database, num uint64) *staking.StakingInfo {
	b, err := db.Get(stakingInfoKey(num))
	if err != nil || len(b) == 0 {
		return nil
	}

	var sl *staking.StakingInfoLegacy
	if err := json.Unmarshal(b, sl); err != nil {
		logger.Error("Malformed staking info", "num", num, "err", err)
		return nil
	}
	return sl.ToStakingInfo()
}

func WriteStakingInfo(db database.Database, num uint64, si *staking.StakingInfo) {
	b, err := json.Marshal(si)
	if err != nil {
		logger.Error("Failed to marshal StakingInfo", "num", num, "err", err)
		return
	}

	if err := db.Put(stakingInfoKey(num), b); err != nil {
		logger.Crit("Failed to write StakingInfo", "num", num, "err", err)
	}
}

func DeleteStakingInfo(db database.Database, num uint64) {
	if err := db.Delete(stakingInfoKey(num)); err != nil {
		logger.Crit("Failed to delete StakingInfo", "num", num, "err", err)
	}
}
