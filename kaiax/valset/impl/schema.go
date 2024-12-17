package impl

import (
	"encoding/binary"
	"encoding/json"
	"slices"
	"sync"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	validatorVoteBlockNums           = []byte("validatorVoteBlockNums")
	lowestScannedValidatorVoteNumKey = []byte("lowestScannedValidatorVoteNum")
	councilPrefix                    = []byte("council")
	istanbulSnapshotKeyPrefix        = []byte("snapshot")

	mu = &sync.RWMutex{}
)

func councilKey(num uint64) []byte {
	return append(councilPrefix, common.Int64ToByteLittleEndian(num)...)
}

func istanbulSnapshotKey(hash common.Hash) []byte {
	return append(istanbulSnapshotKeyPrefix, hash[:]...)
}

func ReadValidatorVoteBlockNums(db database.Database) []uint64 {
	b, err := db.Get(validatorVoteBlockNums)
	if err != nil || len(b) == 0 {
		return nil
	}

	var nums []uint64
	if err := json.Unmarshal(b, &nums); err != nil {
		logger.Error("Malformed valset vote block nums", "err", err)
		return nil
	}
	return nums
}

func writeValidatorVoteBlockNums(db database.Database, nums []uint64) {
	slices.Sort(nums)
	b, err := json.Marshal(nums)
	if err != nil {
		logger.Crit("Failed to marshal valset vote block nums", "err", err)
	}
	if err = db.Put(validatorVoteBlockNums, b); err != nil {
		logger.Crit("Failed to write valset vote block nums", "err", err)
	}
}

// insertValidatorVoteBlockNums inserts a new block num into the validator vote block nums.
func insertValidatorVoteBlockNums(db database.Database, num uint64) {
	mu.Lock()
	defer mu.Unlock()

	nums := ReadValidatorVoteBlockNums(db)

	// Skip if num already exists in the array
	for _, n := range nums {
		if n == num {
			return
		}
	}

	nums = append(nums, num)
	writeValidatorVoteBlockNums(db, nums)
}

// trimValidatorVoteBlockNums deletes all block nums greater than or equal to `since`.
func trimValidatorVoteBlockNums(db database.Database, since uint64) {
	mu.Lock()
	defer mu.Unlock()

	nums := ReadValidatorVoteBlockNums(db)
	if nums == nil {
		return
	}

	nums = slices.DeleteFunc(nums, func(n uint64) bool { return n >= since })
	writeValidatorVoteBlockNums(db, nums)
}

func ReadCouncil(db database.Database, num uint64) []common.Address {
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

func writeCouncil(db database.Database, num uint64, addrs []common.Address) {
	b, err := json.Marshal(addrs)
	if err != nil {
		logger.Crit("Failed to marshal council", "num", num, "err", err)
	}
	if err = db.Put(councilKey(num), b); err != nil {
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
