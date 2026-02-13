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
	councilPermissionlessPrefix      = []byte("councilPermissionless")
	istanbulSnapshotKeyPrefix        = []byte("snapshot")

	validatorStateChangeBlockNums = []byte("validatorStateChangeBlockNums")

	voteNumMu        = &sync.RWMutex{}
	stateChangeNumMu = &sync.RWMutex{}
)

func councilKey(num uint64) []byte {
	return append(councilPermissionedPrefix, common.Int64ToByteLittleEndian(num)...)
}

func councilKeyPermissionless(num uint64) []byte {
	return append(councilPermissionlessPrefix, common.Int64ToByteLittleEndian(num)...)
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

func ReadValidatorStateChangeBlockNums(db database.Database) []uint64 {
	return readVoteOrStateChangeBlockNums(db, validatorStateChangeBlockNums)
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

func writeValidatorStateChangeBlockNums(db database.Database, nums []uint64) {
	writeVoteOrStageChangeBlockNums(db, validatorStateChangeBlockNums, nums)
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

// insertValidatorStateChangeBlockNum inserts a new block num where validator set has been updated
func insertValidatorStateChangeBlockNum(db database.Database, num uint64) {
	stateChangeNumMu.Lock()
	defer stateChangeNumMu.Unlock()

	nums := ReadValidatorStateChangeBlockNums(db)

	// Skip if num already exists in the array
	if slices.Contains(nums, num) {
		return
	}

	nums = append(nums, num)
	writeValidatorStateChangeBlockNums(db, nums)
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

// trimValidatorStateChangeBlockNums deletes all block nums greater than or equal to `since`.
func trimValidatorStateChangeBlockNums(db database.Database, since uint64) {
	stateChangeNumMu.Lock()
	defer stateChangeNumMu.Unlock()

	nums := ReadValidatorStateChangeBlockNums(db)
	if nums == nil {
		return
	}

	nums = slices.DeleteFunc(nums, func(n uint64) bool { return n >= since })
	writeValidatorStateChangeBlockNums(db, nums)
}

func ReadPermissionlessCouncilKeyExist(db database.Database, num uint64) bool {
	b, err := db.Get(councilKeyPermissionless(num))
	if err != nil || len(b) == 0 {
		return false
	}
	return true
}

func ReadCouncil(db database.Database, num uint64) *ValidatorList {
	if council := readCouncilPermissionless(db, num); council != nil {
		return NewValidatorList(council)
	}
	// if no validator of permissionless format is not returned, fallback to legacy format without hardfork branch
	return ConvertLegacyToValidatorList(readCouncilPermissioned(db, num))
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

func readCouncilPermissionless(db database.Database, num uint64) valset.ValidatorChartMap {
	b, err := db.Get(councilKeyPermissionless(num))
	if err != nil || len(b) == 0 {
		return nil
	}

	var result valset.ValidatorChartMap
	if err = json.Unmarshal(b, &result); err != nil {
		logger.Error("Malformed council", "num", num, "err", err)
		return nil
	}
	return result
}

func writeCouncil(db database.Database, num uint64, validators *ValidatorList) {
	var (
		b   []byte
		key []byte
		err error
	)
	if validators.isLegacy {
		b, err = json.Marshal(validators.List())
		key = councilKey(num)
	} else {
		b, err = json.Marshal(validators.permlessVals)
		key = councilKeyPermissionless(num)
	}
	if err != nil {
		logger.Crit("Failed to marshal council", "num", num, "err", err)
	}
	if err = db.Put(key, b); err != nil {
		logger.Crit("Failed to write council", "num", num, "err", err)
	}
	return
}

func deleteCouncil(db database.Database, num uint64) {
	if err := db.Delete(councilKey(num)); err != nil {
		logger.Crit("Failed to delete council", "num", num, "err", err)
	}
	if err := db.Delete(councilKeyPermissionless(num)); err != nil {
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
