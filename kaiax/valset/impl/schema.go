package impl

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	valSetVoteBlockNumsKey             = []byte("valSetVoteBlockNums")
	lowestScannedCheckpointIntervalKey = []byte("lowestScannedCheckpointIntervalKey")
	councilAddressPrefix               = []byte("councilAddress")
	istanbulSnapshotKeyPrefix          = []byte("snapshot") // snapshotKeyPrefix is a governance snapshot prefix

	mu = &sync.RWMutex{}
)

func councilAddressKey(num uint64) []byte {
	return append(councilAddressPrefix, common.Int64ToByteLittleEndian(num)...)
}

func istanbulSnapshotKey(hash common.Hash) []byte {
	return append(istanbulSnapshotKeyPrefix, hash[:]...)
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

func writeLowestScannedCheckpointIntervalNum(db database.Database, scanned uint64) error {
	b, err := json.Marshal(scanned)
	if err != nil {
		return err
	}
	if err = db.Put(lowestScannedCheckpointIntervalKey, b); err != nil {
		return err
	}
	return nil
}

func readLowestScannedCheckpointIntervalNum(db database.Database) (uint64, error) {
	b, err := db.Get(lowestScannedCheckpointIntervalKey)
	if err != nil || len(b) == 0 {
		return 0, fmt.Errorf("failed to read lowest scanned blocknumber from db. error=%v, b=%v", err, string(b))
	}
	var lowestScannedNum uint64
	if err = json.Unmarshal(b, &lowestScannedNum); err != nil {
		return 0, fmt.Errorf("failed to unmarshal encoded lowestScannedNum. err=%v", err)
	}
	return lowestScannedNum, nil
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

func readCouncilAddressListFromIstanbulSnapshot(db database.Database, blockHash common.Hash) ([]common.Address, error) {
	startBlockHash := blockHash

	b, err := db.Get(istanbulSnapshotKey(startBlockHash))
	if err != nil || len(b) == 0 {
		return nil, fmt.Errorf("failed to read istanbul snapshot from db. hash:%s, err:%v, b:%v", startBlockHash.String(), err, string(b))
	}

	snap := new(istanbulBackend.Snapshot)
	if err := json.Unmarshal(b, snap); err != nil {
		return nil, err
	}

	c := make(valset.AddressList, len(snap.ValSet.List())+len(snap.ValSet.DemotedList()))
	for _, val := range append(snap.ValSet.List(), snap.ValSet.DemotedList()...) {
		c = append(c, val.Address())
	}

	sort.Sort(c)
	return c, nil
}

func readCouncilAddressListFromValSetCouncilDB(db database.Database, num uint64) ([]common.Address, error) {
	var (
		voteBlocks = ReadValidatorVoteDataBlockNums(db)
		voteBlock  = uint64(0)
	)

	if voteBlocks == nil {
		return nil, errEmptyVoteBlock
	}
	for i := len(voteBlocks) - 1; i >= 0; i-- {
		if voteBlocks[i] <= num {
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

func WriteCouncilAddressListToDb(db database.Database, voteBlk uint64, councilAddressList valset.AddressList) error {
	copiedList := make(valset.AddressList, len(councilAddressList))
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
