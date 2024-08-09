package staking

import (
	"encoding/json"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	stakingInfoPrefix = []byte("stakingInfo")
)

type storedStakingInfo struct {
	StakingInfo

	// Legacy fields were stored in the database before binary upgrades.
	KIRAddr common.Address `json:"KIRAddr"` // KIRAddr -> KCFAddr from v1.10.2
	PoCAddr common.Address `json:"PoCAddr"` // PoCAddr -> KFFAddr from v1.10.2
	KCFAddr common.Address `json:"kcfAddr"` // KCFAddr -> KEFAddr from Kaia v1.0.0
	KFFAddr common.Address `json:"kffAddr"` // KFFAddr -> KIFAddr from Kaia v1.0.0
}

func makeKey(prefix []byte, num uint64) []byte {
	byteKey := common.Int64ToByteLittleEndian(num)
	return append(prefix, byteKey...)
}

func ReadStakingInfo(db database.Database, num uint64) *StakingInfo {
	b, err := db.Get(makeKey(stakingInfoPrefix, num))
	if err != nil || len(b) == 0 {
		return nil
	}

	si := new(StakingInfo)
	if err := json.Unmarshal(b, si); err != nil {
		logger.Error("Invalid StakingInfo JSON", "num", num, "err", err)
		return nil
	}
	return si
}

func WriteStakingInfo(db database.Database, num uint64, si *StakingInfo) {
	b, err := json.Marshal(si)
	if err != nil {
		logger.Error("Failed to marshal StakingInfo", "num", num, "err", err)
		return
	}

	if err := db.Put(makeKey(stakingInfoPrefix, num), b); err != nil {
		logger.Crit("Failed to write StakingInfo", "num", num, "err", err)
	}
}

func DeleteStakingInfo(db database.Database, num uint64) {
	if err := db.Delete(makeKey(stakingInfoPrefix, num)); err != nil {
		logger.Crit("Failed to delete StakingInfo", "num", num, "err", err)
	}
}
