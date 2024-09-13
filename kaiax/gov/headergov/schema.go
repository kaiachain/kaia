package headergov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/storage/database"
)

var (
	voteDataBlockNumsKey = []byte("governanceVoteDataBlockNums")
	govDataBlockNumsKey  = []byte("governanceDataBlockNums")
)

type StoredUint64Array []uint64

func readStoredUint64Array(db database.Database, key []byte) *StoredUint64Array {
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

func writeStoredUint64Array(db database.Database, key []byte, data *StoredUint64Array) {
	b, err := json.Marshal(data)
	if err != nil {
		logger.Error("Failed to marshal voteDataBlocks", "err", err)
		return
	}

	if err := db.Put(key, b); err != nil {
		logger.Crit("Failed to write voteDataBlocks", "err", err)
	}
}

func ReadVoteDataBlockNums(db database.Database) *StoredUint64Array {
	return readStoredUint64Array(db, voteDataBlockNumsKey)
}

func WriteVoteDataBlockNums(db database.Database, voteDataBlockNums *StoredUint64Array) {
	writeStoredUint64Array(db, voteDataBlockNumsKey, voteDataBlockNums)
}

func ReadGovDataBlockNums(db database.Database) *StoredUint64Array {
	return readStoredUint64Array(db, govDataBlockNumsKey)
}

func WriteGovDataBlockNums(db database.Database, govDataBlockNums *StoredUint64Array) {
	writeStoredUint64Array(db, govDataBlockNumsKey, govDataBlockNums)
}
