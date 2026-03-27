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
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

// TestRewindTo_PurgesCache verifies that RewindTo clears all in-memory pfsCache and
// cpMatrixCache entries regardless of which blocks they cover.
func TestRewindTo_PurgesCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	_ = valset
	db := database.NewMemDB()

	headers := map[uint64]*types.Header{
		0:                           makeHeaderWithRound(0, 0),
		testCheckpointInterval:     makeHeaderWithRound(testCheckpointInterval, 0),
		testCheckpointInterval + 1: makeHeaderWithRound(testCheckpointInterval+1, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	v.pfsCache.Add(uint64(1), map[common.Address]uint64{addrN(0): 1})
	v.pfsCache.Add(testCheckpointInterval+1, map[common.Address]uint64{addrN(0): 2})
	v.cpMatrixCache.Add(uint64(1), vrank.CPMatrix{addrN(1): {addrN(2): 1}})
	v.cpMatrixCache.Add(testCheckpointInterval+1, vrank.CPMatrix{addrN(1): {addrN(2): 2}})

	v.RewindTo(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(testCheckpointInterval))}))

	_, ok := v.pfsCache.Get(uint64(1))
	assert.False(t, ok, "pfsCache entry below rewind point must be purged")
	_, ok = v.pfsCache.Get(testCheckpointInterval + 1)
	assert.False(t, ok, "pfsCache entry above rewind point must be purged")
	_, ok = v.cpMatrixCache.Get(uint64(1))
	assert.False(t, ok, "cpMatrixCache entry below rewind point must be purged")
	_, ok = v.cpMatrixCache.Get(testCheckpointInterval + 1)
	assert.False(t, ok, "cpMatrixCache entry above rewind point must be purged")
}

// TestRewindDelete_UpdatesLastCheckpoint verifies the real host contract:
// RewindDelete is called once per deleted block. When the deleted block is a
// checkpoint and is the latest checkpoint, the pointer must roll back to the
// previous surviving checkpoint.
func TestRewindDelete_UpdatesLastCheckpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	_ = valset
	db := database.NewMemDB()

	cp1 := testCheckpointInterval
	cp2 := 2 * testCheckpointInterval

	headers := map[uint64]*types.Header{
		0:   makeHeaderWithRound(0, 0),
		cp1: makeHeaderWithRound(cp1, 0),
		cp2: makeHeaderWithRound(cp2, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	WriteCheckpoint(db, cp1,
		map[common.Address]uint64{addrN(0): 1},
		vrank.CPMatrix{addrN(1): {addrN(2): 1}},
	)
	WriteCheckpoint(db, cp2,
		map[common.Address]uint64{addrN(0): 2},
		vrank.CPMatrix{addrN(1): {addrN(2): 2}},
	)
	WriteLastCheckpoint(db, cp2)

	// Host deletes block cp2.
	v.RewindDelete(common.Hash{}, cp2)
	assert.Nil(t, ReadCheckpointPFS(db, cp2), "cp2 checkpoint should be deleted")
	assert.Nil(t, ReadCheckpointCPMatrix(db, cp2), "cp2 cpMatrix checkpoint should be deleted")
	assert.NotNil(t, ReadCheckpointPFS(db, cp1), "cp1 checkpoint should survive")
	lastCP, hasCP := ReadLastCheckpoint(db)
	assert.True(t, hasCP)
	assert.Equal(t, cp1, lastCP)

	// Host then deletes block cp1.
	v.RewindDelete(common.Hash{}, cp1)
	assert.Nil(t, ReadCheckpointPFS(db, cp1), "cp1 checkpoint should be deleted")
	assert.Nil(t, ReadCheckpointCPMatrix(db, cp1), "cp1 cpMatrix checkpoint should be deleted")
	_, hasCP = ReadLastCheckpoint(db)
	assert.False(t, hasCP, "lastCheckpoint should be deleted when no checkpoints survive")
}

// TestRewindDelete_MultiIntervalReorg verifies the host's actual per-block delete
// sequence. When rewinding from cp3 to cp1, the host calls RewindDelete(cp3)
// and then RewindDelete(cp2). Each deleted checkpoint must be removed exactly
// when its block is deleted.
func TestRewindDelete_MultiIntervalReorg(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	_ = valset
	db := database.NewMemDB()

	cp1 := testCheckpointInterval
	cp2 := 2 * testCheckpointInterval
	cp3 := 3 * testCheckpointInterval

	headers := map[uint64]*types.Header{
		0:   makeHeaderWithRound(0, 0),
		cp1: makeHeaderWithRound(cp1, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	// Simulate old checkpoints written before the reorg.
	WriteCheckpoint(db, cp1, map[common.Address]uint64{addrN(0): 1}, nil)
	WriteCheckpoint(db, cp2, map[common.Address]uint64{addrN(0): 2}, nil)
	WriteCheckpoint(db, cp3, map[common.Address]uint64{addrN(0): 3}, nil)
	WriteLastCheckpoint(db, cp3)

	// Host deletes cp3 first.
	v.RewindDelete(common.Hash{}, cp3)
	assert.Nil(t, ReadCheckpointPFS(db, cp3), "cp3 checkpoint should be deleted after 2-interval reorg")
	assert.Nil(t, ReadCheckpointCPMatrix(db, cp3), "cp3 cpMatrix checkpoint should be deleted after 2-interval reorg")
	assert.NotNil(t, ReadCheckpointPFS(db, cp2), "cp2 checkpoint should still exist until block cp2 is deleted")
	assert.Empty(t, ReadCheckpointCPMatrix(db, cp2), "cp2 cpMatrix checkpoint should remain empty because the test stored no cpMatrix payload")
	lastCP, hasCP := ReadLastCheckpoint(db)
	assert.True(t, hasCP)
	assert.Equal(t, cp2, lastCP)

	// Host deletes cp2 next.
	v.RewindDelete(common.Hash{}, cp2)
	assert.Nil(t, ReadCheckpointPFS(db, cp2), "cp2 checkpoint should be deleted after 2-interval reorg")
	assert.Nil(t, ReadCheckpointCPMatrix(db, cp2), "cp2 cpMatrix checkpoint should be deleted after 2-interval reorg")
	assert.NotNil(t, ReadCheckpointPFS(db, cp1), "cp1 checkpoint should survive")
	lastCP, hasCP = ReadLastCheckpoint(db)
	assert.True(t, hasCP)
	assert.Equal(t, cp1, lastCP)
}
