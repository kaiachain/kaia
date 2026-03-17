// Copyright 2026 The Kaia Authors
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

package impl

import (
	"bytes"
	"encoding/binary"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	vrankCheckpointPrefix = []byte("vrankCheckpoint")
	lastCheckpointKey     = []byte("vrankLastCheckpoint") // latest checkpoint block number
)

type pfsEntry struct {
	Addr  common.Address
	Score uint64
}

type cpMatrixEntry struct {
	Candidate common.Address
	Reporter  common.Address
	Score     uint64
}

type vrankCheckpointStorage struct {
	PFS      []pfsEntry
	CPMatrix []cpMatrixEntry
}

func scoreCheckpointKey(blockNum uint64) []byte {
	return append(vrankCheckpointPrefix, common.Int64ToByteBigEndian(blockNum)...)
}

// ReadCheckpoint returns the PFS map and cpMatrix stored at blockNum.
// Returns (nil, nil) if no checkpoint exists for that block.
func ReadCheckpoint(db database.Database, blockNum uint64) (map[common.Address]uint64, map[common.Address]map[common.Address]uint64) {
	b, err := db.Get(scoreCheckpointKey(blockNum))
	if err != nil || len(b) == 0 {
		return nil, nil
	}
	var stored vrankCheckpointStorage
	if err := rlp.DecodeBytes(b, &stored); err != nil {
		logger.Crit("Failed to deserialize checkpoint", "blockNum", blockNum, "err", err)
	}

	pfs := make(map[common.Address]uint64, len(stored.PFS))
	for _, entry := range stored.PFS {
		pfs[entry.Addr] = entry.Score
	}

	cpMatrix := make(map[common.Address]map[common.Address]uint64)
	for _, entry := range stored.CPMatrix {
		if _, ok := cpMatrix[entry.Candidate]; !ok {
			cpMatrix[entry.Candidate] = make(map[common.Address]uint64)
		}
		cpMatrix[entry.Candidate][entry.Reporter] = entry.Score
	}
	return pfs, cpMatrix
}

// WriteCheckpoint persists the PFS map and cpMatrix together at blockNum.
func WriteCheckpoint(db database.Database, blockNum uint64, pfs map[common.Address]uint64, cpMatrix map[common.Address]map[common.Address]uint64) {
	pfsEntries := make([]pfsEntry, 0, len(pfs))
	for addr, score := range pfs {
		pfsEntries = append(pfsEntries, pfsEntry{Addr: addr, Score: score})
	}
	slices.SortFunc(pfsEntries, func(a, b pfsEntry) int { return bytes.Compare(a.Addr.Bytes(), b.Addr.Bytes()) })

	candidates := make([]common.Address, 0, len(cpMatrix))
	for candidate := range cpMatrix {
		candidates = append(candidates, candidate)
	}
	slices.SortFunc(candidates, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })

	cfsEntries := make([]cpMatrixEntry, 0)
	for _, candidate := range candidates {
		reporters := make([]common.Address, 0, len(cpMatrix[candidate]))
		for reporter := range cpMatrix[candidate] {
			reporters = append(reporters, reporter)
		}
		slices.SortFunc(reporters, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
		for _, reporter := range reporters {
			cfsEntries = append(cfsEntries, cpMatrixEntry{
				Candidate: candidate,
				Reporter:  reporter,
				Score:     cpMatrix[candidate][reporter],
			})
		}
	}

	b, err := rlp.EncodeToBytes(vrankCheckpointStorage{PFS: pfsEntries, CPMatrix: cfsEntries})
	if err != nil {
		logger.Crit("Failed to serialize checkpoint", "blockNum", blockNum, "err", err)
	}
	if err := db.Put(scoreCheckpointKey(blockNum), b); err != nil {
		logger.Crit("Failed to write checkpoint", "blockNum", blockNum, "err", err)
	}
}

// DeleteCheckpoint removes the checkpoint stored at blockNum.
func DeleteCheckpoint(db database.Database, blockNum uint64) {
	if err := db.Delete(scoreCheckpointKey(blockNum)); err != nil {
		logger.Crit("Failed to delete checkpoint", "blockNum", blockNum, "err", err)
	}
}

func calcCheckpointBlock(blockNum uint64) uint64 {
	return blockNum - (blockNum % scoreCheckpointInterval)
}

// ReadLastCheckpoint returns the block number of the most recently written checkpoint and true,
// or (0, false) if no checkpoint pointer has been written.
// Note: 0 is a valid checkpoint block number, so callers must check the bool.
func ReadLastCheckpoint(db database.Database) (uint64, bool) {
	b, err := db.Get(lastCheckpointKey)
	if err != nil || len(b) == 0 {
		return 0, false
	}
	return binary.BigEndian.Uint64(b), true
}

// WriteLastCheckpoint records blockNum as the most recently written checkpoint.
func WriteLastCheckpoint(db database.Database, blockNum uint64) {
	if err := db.Put(lastCheckpointKey, common.Int64ToByteBigEndian(blockNum)); err != nil {
		logger.Crit("Failed to write last checkpoint", "blockNum", blockNum, "err", err)
	}
}

// DeleteLastCheckpoint removes the last-checkpoint pointer from the DB.
func DeleteLastCheckpoint(db database.Database) {
	if err := db.Delete(lastCheckpointKey); err != nil {
		logger.Crit("Failed to delete last checkpoint", "err", err)
	}
}
