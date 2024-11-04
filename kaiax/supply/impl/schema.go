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
	lastSupplyCheckpointNumberKey = []byte("lastSupplyCheckpointNumber")
	supplyCheckpointPrefix        = []byte("supplyCheckpoint")
)

// supplyCheckpoint is the disk format for checkpoints (i.e. periodically committed AccReward).
type supplyCheckpoint struct {
	Minted   []byte
	BurntFee []byte
}

func supplyCheckpointKey(blockNumber uint64) []byte {
	return append(supplyCheckpointPrefix, common.Int64ToByteBigEndian(blockNumber)...)
}

func ReadLastSupplyCheckpointNumber(db database.Database) uint64 {
	b, err := db.Get(lastSupplyCheckpointNumberKey)
	if err != nil || len(b) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func WriteLastSupplyCheckpointNumber(db database.Database, num uint64) {
	data := common.Int64ToByteBigEndian(num)
	if err := db.Put(lastSupplyCheckpointNumberKey, data); err != nil {
		logger.Crit("Failed to write last supply checkpoint number", "err", err)
	}
}

func ReadSupplyCheckpoint(db database.Database, num uint64) *supply.AccReward {
	b, err := db.Get(supplyCheckpointKey(num))
	if err != nil || len(b) == 0 {
		return nil
	}
	stored := &supplyCheckpoint{}
	if err := rlp.DecodeBytes(b, stored); err != nil {
		logger.Crit("Failed to deserialize supply checkpoint", "err", err)
	}
	return &supply.AccReward{
		Minted:   new(big.Int).SetBytes(stored.Minted),
		BurntFee: new(big.Int).SetBytes(stored.BurntFee),
	}
}

func WriteSupplyCheckpoint(db database.Database, num uint64, checkpoint *supply.AccReward) {
	stored := &supplyCheckpoint{
		Minted:   checkpoint.Minted.Bytes(),
		BurntFee: checkpoint.BurntFee.Bytes(),
	}
	b, err := rlp.EncodeToBytes(stored)
	if err != nil {
		logger.Crit("Failed to serialize supply checkpoint", "err", err)
	}
	if err := db.Put(supplyCheckpointKey(num), b); err != nil {
		logger.Crit("Failed to write supply checkpoint", "err", err)
	}
}

func DeleteSupplyCheckpoint(db database.Database, num uint64) {
	if err := db.Delete(supplyCheckpointKey(num)); err != nil {
		logger.Crit("Failed to delete supply checkpoint", "err", err)
	}
}
