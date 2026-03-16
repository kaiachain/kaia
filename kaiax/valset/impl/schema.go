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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"slices"
	"sync"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	validatorVoteBlockNums           = []byte("validatorVoteBlockNums")
	lowestScannedValidatorVoteNumKey = []byte("lowestScannedValidatorVoteNum")
	councilPermissionedPrefix        = []byte("council")
	istanbulSnapshotKeyPrefix        = []byte("snapshot")

	voteNumMu = &sync.RWMutex{}
)

func councilKey(num uint64) []byte {
	return append(councilPermissionedPrefix, common.Int64ToByteLittleEndian(num)...)
}

func istanbulSnapshotKey(hash common.Hash) []byte {
	return append(istanbulSnapshotKeyPrefix, hash[:]...)
}

func readVoteOrStateChangeBlockNums(db database.Database, key []byte) []uint64 {
	b, err := db.Get(key)
	if err != nil || len(b) == 0 {
		return nil
	}

	var nums []uint64
	if err := json.Unmarshal(b, &nums); err != nil {
		logger.Error(fmt.Sprintf("Malformed %s block nums", string(key)), "err", err)
		return nil
	}
	return nums
}

func ReadValidatorVoteBlockNums(db database.Database) []uint64 {
	return readVoteOrStateChangeBlockNums(db, validatorVoteBlockNums)
}

func writeVoteOrStageChangeBlockNums(db database.Database, key []byte, nums []uint64) {
	slices.Sort(nums)
	b, err := json.Marshal(nums)
	if err != nil {
		logger.Crit(fmt.Sprintf("Failed to marshal %s block nums", string(key)), "err", err)
	}
	if err = db.Put(key, b); err != nil {
		logger.Crit(fmt.Sprintf("Failed to write %s block nums", string(key)), "err", err)
	}
}

func writeValidatorVoteBlockNums(db database.Database, nums []uint64) {
	writeVoteOrStageChangeBlockNums(db, validatorVoteBlockNums, nums)
}

// insertValidatorVoteBlockNums inserts a new block num into the validator vote block nums.
func insertValidatorVoteBlockNums(db database.Database, num uint64) {
	voteNumMu.Lock()
	defer voteNumMu.Unlock()

	nums := ReadValidatorVoteBlockNums(db)

	// Skip if num already exists in the array
	if slices.Contains(nums, num) {
		return
	}

	nums = append(nums, num)
	writeValidatorVoteBlockNums(db, nums)
}

// trimValidatorVoteBlockNums deletes all block nums greater than or equal to `since`.
func trimValidatorVoteBlockNums(db database.Database, since uint64) {
	voteNumMu.Lock()
	defer voteNumMu.Unlock()

	nums := ReadValidatorVoteBlockNums(db)
	if nums == nil {
		return
	}

	nums = slices.DeleteFunc(nums, func(n uint64) bool { return n >= since })
	writeValidatorVoteBlockNums(db, nums)
}

func ReadCouncilPermissiond(db database.Database, permissionedNum uint64) valset.CommonAddressSet {
	permissionedCouncil := readCouncilPermissioned(db, permissionedNum)
	if permissionedCouncil == nil {
		return nil
	}
	return valset.NewCommonAddressSet(permissionedCouncil)
}

func readCouncilPermissioned(db database.Database, num uint64) []common.Address {
	b, err := db.Get(councilKey(num))
	if err != nil || len(b) == 0 {
		return nil
	}

	var addrs []common.Address
	if err = json.Unmarshal(b, &addrs); err != nil {
		logger.Error("Malformed council", "num", num, "err", err)
		return nil
	}
	return addrs
}

func writeCouncilPermissioned(db database.Database, num uint64, validators valset.CommonAddressSet) {
	key := councilKey(num)
	marshalAndWrite(db, num, key, validators)
}

func marshalAndWrite(db database.Database, num uint64, key []byte, validators valset.CommonAddressSet) {
	b, err := validators.Marshal()
	if err != nil {
		logger.Crit("Failed to marshal council", "num", num, "err", err)
	}
	if err = db.Put(key, b); err != nil {
		logger.Crit("Failed to write council", "num", num, "err", err)
	}
}

func deleteCouncil(db database.Database, num uint64) {
	if err := db.Delete(councilKey(num)); err != nil {
		logger.Crit("Failed to delete council", "num", num, "err", err)
	}
}

func ReadLowestScannedVoteNum(db database.Database) *uint64 {
	b, err := db.Get(lowestScannedValidatorVoteNumKey)
	if err != nil || len(b) == 0 {
		return nil
	}
	if len(b) != 8 {
		logger.Error("Malformed lowest scanned snapshot num", "length", len(b))
		return nil
	}
	ret := binary.BigEndian.Uint64(b)
	return &ret
}

func writeLowestScannedVoteNum(db database.Database, num uint64) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, num)
	if err := db.Put(lowestScannedValidatorVoteNumKey, b); err != nil {
		logger.Crit("Failed to write lowest scanned snapshot num", "num", num, "err", err)
	}
}

// Only relevant fields of the JSON-encoded istanbul snapshot.
type istanbulSnapshotStorage struct {
	Validators        []common.Address `json:"validators"` // qualified validators
	DemotedValidators []common.Address `json:"demotedValidators"`
}

func ReadIstanbulSnapshot(db database.Database, hash common.Hash) []common.Address {
	b, err := db.Get(istanbulSnapshotKey(hash))
	if err != nil || len(b) == 0 {
		return nil
	}

	snap := new(istanbulSnapshotStorage)
	if err := json.Unmarshal(b, snap); err != nil {
		logger.Error("Malformed istanbul snapshot", "hash", hash.String(), "err", err)
		return nil
	}
	return append(snap.Validators, snap.DemotedValidators...)
}
