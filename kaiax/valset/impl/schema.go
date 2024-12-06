package impl

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/kaiachain/kaia/common"
	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	valSetVoteBlockNumsKey             = []byte("valSetVoteBlockNums")
	lowestScannedCheckpointIntervalKey = []byte("lowestScannedCheckpointIntervalKey")
	councilPrefix                      = []byte("council")
	istanbulSnapshotKeyPrefix          = []byte("snapshot") // snapshotKeyPrefix is a governance snapshot prefix

	mu = &sync.RWMutex{}
)

func councilKey(num uint64) []byte {
	return append(councilPrefix, common.Int64ToByteLittleEndian(num)...)
}

func istanbulSnapshotKey(hash common.Hash) []byte {
	return append(istanbulSnapshotKeyPrefix, hash[:]...)
}

func readValsetVoteBlockNums(db database.Database) []uint64 {
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

func writeValsetVoteBlockNums(db database.Database, voteBlk uint64) error {
	voteBlks := []uint64{0}
	if voteBlk != 0 {
		voteBlks = readValsetVoteBlockNums(db)
		if len(voteBlks) == 0 {
			return errEmptyVoteBlock
		}
		if voteBlks[len(voteBlks)-1] > voteBlk {
			return fmt.Errorf("invalid voteBlk %d, last voteBlk %d", voteBlk, voteBlks[len(voteBlks)-1])
		}
		voteBlks = append(voteBlks, voteBlk)
	}

	b, err := json.Marshal(voteBlks)
	if err != nil {
		return err
	}
	if err = db.Put(valSetVoteBlockNumsKey, b); err != nil {
		return err
	}
	return nil
}

func readCouncil(db database.Database, num uint64) ([]common.Address, error) {
	mu.RLock()
	defer mu.RUnlock()

	var (
		voteBlocks = readValsetVoteBlockNums(db)
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

	b, err := db.Get(councilKey(voteBlock))
	if err != nil || len(b) == 0 {
		return nil, fmt.Errorf("failed to read council from db at voteBlk %d. error=%v, b=%v", voteBlock, err, string(b))
	}

	var set []common.Address
	if err = json.Unmarshal(b, &set); err != nil {
		return nil, fmt.Errorf("failed to unmarshal encoded council at voteBlk %d. err=%v", voteBlock, err)
	}
	if num > 0 && voteBlock == 0 {
		sort.Sort(valset.AddressList(set))
	}
	return set, nil
}

func writeCouncil(db database.Database, voteBlk uint64, council []common.Address) error {
	mu.Lock()
	defer mu.Unlock()

	copiedList := make([]common.Address, len(council))
	copy(copiedList, council)

	if err := writeValsetVoteBlockNums(db, voteBlk); err != nil {
		return err
	}

	if voteBlk > 0 {
		sort.Sort(valset.AddressList(copiedList))
	}
	b, err := json.Marshal(copiedList)
	if err != nil {
		return err
	}
	if err = db.Put(councilKey(voteBlk), b); err != nil {
		return err
	}
	return nil
}

func readIstanbulSnapshot(db database.Database, blockHash common.Hash) (*istanbulBackend.Snapshot, error) {
	startBlockHash := blockHash

	b, err := db.Get(istanbulSnapshotKey(startBlockHash))
	if err != nil || len(b) == 0 {
		return nil, fmt.Errorf("failed to read istanbul snapshot from db. hash:%s, err:%v, b:%v", startBlockHash.String(), err, string(b))
	}

	snap := new(istanbulBackend.Snapshot)
	if err := json.Unmarshal(b, snap); err != nil {
		return nil, err
	}

	return snap, nil
}

func readLowestScannedCheckpointInterval(db database.Database) (uint64, error) {
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

func writeLowestScannedCheckpointInterval(db database.Database, scanned uint64) error {
	b, err := json.Marshal(scanned)
	if err != nil {
		return err
	}
	if err = db.Put(lowestScannedCheckpointIntervalKey, b); err != nil {
		return err
	}
	return nil
}
