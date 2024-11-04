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

package supply

import (
	"encoding/binary"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	// Using "supplyCheckpoint" for backward compatibility
	lastAccRewardNumberKey = []byte("lastSupplyCheckpointNumber")
	accRewardPrefix        = []byte("supplyCheckpoint")
)

// accRewardStorage is the disk format for checkpoints (i.e. periodically committed AccReward).
type accRewardStorage struct {
	Minted   []byte
	BurntFee []byte
}

func accRewardKey(blockNumber uint64) []byte {
	return append(accRewardPrefix, common.Int64ToByteBigEndian(blockNumber)...)
}

func ReadLastAccRewardNumber(db database.Database) uint64 {
	b, err := db.Get(lastAccRewardNumberKey)
	if err != nil || len(b) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func WriteLastAccRewardNumber(db database.Database, num uint64) {
	data := common.Int64ToByteBigEndian(num)
	if err := db.Put(lastAccRewardNumberKey, data); err != nil {
		logger.Crit("Failed to write highest acc reward number", "err", err)
	}
}

func ReadAccReward(db database.Database, num uint64) *supply.AccReward {
	b, err := db.Get(accRewardKey(num))
	if err != nil || len(b) == 0 {
		return nil
	}
	stored := &accRewardStorage{}
	if err := rlp.DecodeBytes(b, stored); err != nil {
		logger.Crit("Failed to deserialize acc reward", "err", err)
	}
	return &supply.AccReward{
		TotalMinted: new(big.Int).SetBytes(stored.Minted),
		BurntFee:    new(big.Int).SetBytes(stored.BurntFee),
	}
}

func WriteAccReward(db database.Database, num uint64, accReward *supply.AccReward) {
	stored := &accRewardStorage{
		Minted:   accReward.TotalMinted.Bytes(),
		BurntFee: accReward.BurntFee.Bytes(),
	}
	b, err := rlp.EncodeToBytes(stored)
	if err != nil {
		logger.Crit("Failed to serialize acc reward", "err", err)
	}
	if err := db.Put(accRewardKey(num), b); err != nil {
		logger.Crit("Failed to write acc reward", "err", err)
	}
}

func DeleteAccReward(db database.Database, num uint64) {
	if err := db.Delete(accRewardKey(num)); err != nil {
		logger.Crit("Failed to delete acc reward", "err", err)
	}
}
