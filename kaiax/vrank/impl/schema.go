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
	"maps"
	"slices"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	vrankCheckpointPrefix = []byte("vrankCheckpoint")
	lastCheckpointKey     = []byte("vrankLastCheckpoint") // latest checkpoint block number
)

type pfsStorage struct {
	Addrs  []common.Address
	Scores []uint64
}

type cpMatrixStorage struct {
	Candidates []common.Address
	Reporters  []common.Address
	Matrix     [][]uint64
}

type vrankCheckpointStorage struct {
	PFS      pfsStorage
	CPMatrix cpMatrixStorage
}

func scoreCheckpointKey(blockNum uint64) []byte {
	return append(vrankCheckpointPrefix, common.Int64ToByteBigEndian(blockNum)...)
}

func readCheckpointStorage(db database.Database, blockNum uint64) *vrankCheckpointStorage {
	b, err := db.Get(scoreCheckpointKey(blockNum))
	if err != nil || len(b) == 0 {
		return nil
	}
	var stored vrankCheckpointStorage
	if err := rlp.DecodeBytes(b, &stored); err != nil {
		logger.Crit("Failed to deserialize checkpoint", "blockNum", blockNum, "err", err)
	}
	return &stored
}

// ReadCheckpointPFS returns the PFS map stored at blockNum, or nil if not found.
func ReadCheckpointPFS(db database.Database, blockNum uint64) map[common.Address]uint64 {
	stored := readCheckpointStorage(db, blockNum)
	if stored == nil {
		return nil
	}
	return deserializePFS(stored.PFS)
}

// ReadCheckpointCPMatrix returns the CPMatrix stored at blockNum, or nil if not found.
func ReadCheckpointCPMatrix(db database.Database, blockNum uint64) vrank.CPMatrix {
	stored := readCheckpointStorage(db, blockNum)
	if stored == nil {
		return nil
	}
	return deserializeCFS(stored.CPMatrix)
}

// ReadCheckpoint returns the PFS map and CPMatrix stored at blockNum.
// Returns (nil, nil) if no checkpoint exists for that block.
func ReadCheckpoint(db database.Database, blockNum uint64) (map[common.Address]uint64, vrank.CPMatrix) {
	stored := readCheckpointStorage(db, blockNum)
	if stored == nil {
		return nil, nil
	}
	return deserializePFS(stored.PFS), deserializeCFS(stored.CPMatrix)
}

// WriteCheckpoint persists the PFS map and cpMatrix together at blockNum.
func WriteCheckpoint(db database.Database, blockNum uint64, pfs map[common.Address]uint64, cpMatrix vrank.CPMatrix) {
	b, err := rlp.EncodeToBytes(vrankCheckpointStorage{
		PFS:      serializePFS(pfs),
		CPMatrix: serializeCFS(cpMatrix),
	})
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

func serializePFS(pfs map[common.Address]uint64) pfsStorage {
	addrs := slices.Collect(maps.Keys(pfs))
	slices.SortFunc(addrs, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })
	scores := make([]uint64, len(addrs))
	for i, addr := range addrs {
		scores[i] = pfs[addr]
	}
	return pfsStorage{
		Addrs:  addrs,
		Scores: scores,
	}
}

func serializeCFS(cpMatrix vrank.CPMatrix) cpMatrixStorage {
	candidates := slices.Collect(maps.Keys(cpMatrix))
	slices.SortFunc(candidates, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })

	reporterSet := make(map[common.Address]struct{})
	for _, candidate := range candidates {
		for _, reporter := range slices.Collect(maps.Keys(cpMatrix[candidate])) {
			reporterSet[reporter] = struct{}{}
		}
	}
	reporters := slices.Collect(maps.Keys(reporterSet))
	slices.SortFunc(reporters, func(a, b common.Address) int { return bytes.Compare(a.Bytes(), b.Bytes()) })

	matrix := make([][]uint64, len(candidates))
	for i, candidate := range candidates {
		matrix[i] = make([]uint64, len(reporters))
		for j, reporter := range reporters {
			matrix[i][j] = cpMatrix[candidate][reporter]
		}
	}

	return cpMatrixStorage{
		Candidates: candidates,
		Reporters:  reporters,
		Matrix:     matrix,
	}
}

func deserializePFS(stored pfsStorage) map[common.Address]uint64 {
	pfs := make(map[common.Address]uint64, len(stored.Addrs))
	for i, addr := range stored.Addrs {
		if i >= len(stored.Scores) {
			break
		}
		pfs[addr] = stored.Scores[i]
	}
	return pfs
}

func deserializeCFS(stored cpMatrixStorage) vrank.CPMatrix {
	cpMatrix := make(vrank.CPMatrix, len(stored.Candidates))
	for rowIdx, candidate := range stored.Candidates {
		cpMatrix[candidate] = make(map[common.Address]uint64)
		if rowIdx >= len(stored.Matrix) {
			break
		}
		row := stored.Matrix[rowIdx]
		for colIdx, reporter := range stored.Reporters {
			if colIdx >= len(row) {
				break
			}
			cpMatrix[candidate][reporter] = row[colIdx]
		}
	}
	return cpMatrix
}
