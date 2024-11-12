package impl

import (
	"encoding/binary"
	"encoding/json"
	"sort"
	"sync"

	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey         = []byte("governanceVoteDataBlockNums")
	govDataBlockNumsKey          = []byte("governanceDataBlockNums")
	lowestVoteScannedBlockNumKey = []byte("governanceLowestVoteScannedBlockNum") // grows downwards
	mu                           = &sync.RWMutex{}
)

type StoredUint64Array []uint64

// readStoredUint64ArrayNoLock should be called only when the caller holds the lock.
func readStoredUint64ArrayNoLock(db database.Database, key []byte) *StoredUint64Array {
	b, err := db.Get(key)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new(StoredUint64Array)
	if err := json.Unmarshal(b, ret); err != nil {
		logger.Error("Invalid voteDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

// writeStoredUint64ArrayNoLock should be called only when the caller holds the lock.
func writeStoredUint64ArrayNoLock(db database.Database, key []byte, data *StoredUint64Array) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(key, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func readStoredUint64Array(db database.Database, key []byte) *StoredUint64Array {
	mu.RLock()
	defer mu.RUnlock()

	return readStoredUint64ArrayNoLock(db, key)
}

func writeStoredUint64Array(db database.Database, key []byte, data *StoredUint64Array) {
	mu.Lock()
	defer mu.Unlock()

	writeStoredUint64ArrayNoLock(db, key, data)
}

func ReadVoteDataBlockNums(db database.Database) *StoredUint64Array {
	return readStoredUint64Array(db, voteDataBlockNumsKey)
}

func WriteVoteDataBlockNums(db database.Database, voteDataBlockNums *StoredUint64Array) {
	writeStoredUint64Array(db, voteDataBlockNumsKey, voteDataBlockNums)
}

func InsertVoteDataBlockNum(db database.Database, blockNum uint64) {
	mu.Lock()
	defer mu.Unlock()

	blockNums := readStoredUint64ArrayNoLock(db, voteDataBlockNumsKey)
	if blockNums == nil {
		blockNums = new(StoredUint64Array)
	}

	// Check if blockNum already exists in the array
	for _, num := range *blockNums {
		if num == blockNum {
			return
		}
	}

	*blockNums = append(*blockNums, blockNum)
	// Sort the block numbers in ascending order
	sort.Slice(*blockNums, func(i, j int) bool {
		return (*blockNums)[i] < (*blockNums)[j]
	})

	writeStoredUint64ArrayNoLock(db, voteDataBlockNumsKey, blockNums)
}

func ReadGovDataBlockNums(db database.Database) *StoredUint64Array {
	return readStoredUint64Array(db, govDataBlockNumsKey)
}

func WriteGovDataBlockNums(db database.Database, govDataBlockNums *StoredUint64Array) {
	writeStoredUint64Array(db, govDataBlockNumsKey, govDataBlockNums)
}

func ReadLowestVoteScannedBlockNum(db database.Database) *uint64 {
	mu.RLock()
	defer mu.RUnlock()

	b, err := db.Get(lowestVoteScannedBlockNumKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	if len(b) != 8 {
		logger.Error("Invalid lowestVoteScannedBlockNum data length", "length", len(b))
		return nil
	}

	ret := binary.BigEndian.Uint64(b)
	return &ret
}

func WriteLowestVoteScannedBlockNum(db database.Database, lowestVoteScannedBlockNum uint64) {
	mu.Lock()
	defer mu.Unlock()

	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, lowestVoteScannedBlockNum)
	db.Put(lowestVoteScannedBlockNumKey, b)
}
