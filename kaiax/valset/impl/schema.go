package impl

import (
	"encoding/json"
	"fmt"
	"sort"
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

func ReadValidatorVoteDataBlockNums(db database.Database) []uint64 {
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
	return *ret
}

func UpdateValidatorVoteDataBlockNums(db database.Database, voteBlk uint64) error {
	voteBlks := []uint64{0}
	if voteBlk != 0 {
		voteBlks = ReadValidatorVoteDataBlockNums(db)
		if len(voteBlks) == 0 {
			return errEmptyVoteBlock
		}
		if voteBlks[len(voteBlks)-1] > voteBlk {
			return fmt.Errorf("invalid voteBlk %d, last voteBlk %d", voteBlk, voteBlks[len(voteBlks)-1])
		}
		voteBlks = append(voteBlks, voteBlk)
	}

	mu.Lock()
	defer mu.Unlock()

	b, err := json.Marshal(voteBlks)
	if err != nil {
		return err
	}
	if err = db.Put(valSetVoteBlockNumsKey, b); err != nil {
		return err
	}
	return nil
}

// ReadCouncilAddressListFromDb gets voteBlk from valset DB
// TODO-kaia-valset: try fetch council from cache or iterate to process the votes between snapshotBlock and num.
func ReadCouncilAddressListFromDb(db database.Database, bn uint64) ([]common.Address, error) {
	var (
		voteBlocks = ReadValidatorVoteDataBlockNums(db)
		voteBlock  = uint64(0)
	)

	if voteBlocks == nil {
		return nil, errEmptyVoteBlock
	}
	for i := len(voteBlocks) - 1; i >= 0; i-- {
		if voteBlocks[i] <= bn {
			voteBlock = voteBlocks[i]
			break
		}
	}

	mu.Lock()
	defer mu.Unlock()

	b, err := db.Get(councilAddressKey(voteBlock))
	if err != nil || len(b) == 0 {
		return nil, fmt.Errorf("failed to read council addresses from db at voteBlk %d. error=%v, b=%v", voteBlock, err, string(b))
	}

	var set []common.Address
	if err = json.Unmarshal(b, &set); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encoded council addresses at voteBlk %d. err=%v", voteBlock, err)
	}
	return set, nil
}

func WriteCouncilAddressListToDb(db database.Database, voteBlk uint64, councilAddressList subsetCouncilSlice) error {
	copiedList := make(subsetCouncilSlice, len(councilAddressList))
	copy(copiedList, councilAddressList)

	sort.Sort(copiedList)
	if err := UpdateValidatorVoteDataBlockNums(db, voteBlk); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	b, err := json.Marshal(copiedList)
	if err != nil {
		return err
	}
	if err = db.Put(councilAddressKey(voteBlk), b); err != nil {
		return err
	}
	return nil
}
