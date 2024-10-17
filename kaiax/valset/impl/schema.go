package impl

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	valSetVoteBlockNumsKey = []byte("valSetVoteBlockNums")
	councilAddressPrefix   = []byte("councilAddress")
	mu                     = &sync.RWMutex{}
)

func councilAddressKey(num uint64) []byte {
	return append(councilAddressPrefix, common.Int64ToByteLittleEndian(num)...)
}

func ReadValidatorVoteDataBlockNums(db database.Database) *[]uint64 {
	mu.Lock()
	defer mu.Unlock()

	b, err := db.Get(valSetVoteBlockNumsKey)
	if err != nil || len(b) == 0 {
		return nil
	}

	ret := new([]uint64)
	if err := json.Unmarshal(b, ret); err != nil {
		logger.Error("Invalid valSetVoteDataBlocks JSON", "err", err)
		return nil
	}
	return ret
}

func WriteValidatorVoteDataBlockNums(db database.Database, data *[]uint64) error {
	mu.Lock()
	defer mu.Unlock()

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = db.Put(valSetVoteBlockNumsKey, b); err != nil {
		return err
	}
	return nil
}

func ReadCouncilAddressListFromDb(db database.Database, voteBlk uint64) ([]common.Address, error) {
	mu.Lock()
	defer mu.Unlock()

	b, err := db.Get(councilAddressKey(voteBlk))
	if err != nil || len(b) == 0 {
		return nil, fmt.Errorf("failed to read council addresses from db at voteBlk %d. error=%v, b=%v", voteBlk, err, string(b))
	}

	var set []common.Address
	if err = json.Unmarshal(b, &set); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encoded council addresses at voteBlk %d. err=%v", voteBlk, err)
	}
	return set, nil
}

func WriteCouncilAddressListToDb(db database.Database, voteBlk uint64, councilAddressList []common.Address) error {
	mu.Lock()
	defer mu.Unlock()

	b, err := json.Marshal(councilAddressList)
	if err != nil {
		return err
	}
	if err = db.Put(councilAddressKey(voteBlk), b); err != nil {
		return err
	}
	return nil
}
