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
		scoreCheckpointInterval:     makeHeaderWithRound(scoreCheckpointInterval, 0),
		scoreCheckpointInterval + 1: makeHeaderWithRound(scoreCheckpointInterval+1, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	v.pfsCache.Add(uint64(1), map[common.Address]uint64{addrN(0): 1})
	v.pfsCache.Add(scoreCheckpointInterval+1, map[common.Address]uint64{addrN(0): 2})
	v.cpMatrixCache.Add(uint64(1), map[common.Address]map[common.Address]uint64{addrN(1): {addrN(2): 1}})
	v.cpMatrixCache.Add(scoreCheckpointInterval+1, map[common.Address]map[common.Address]uint64{addrN(1): {addrN(2): 2}})

	v.RewindTo(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(scoreCheckpointInterval))}))

	_, ok := v.pfsCache.Get(uint64(1))
	assert.False(t, ok, "pfsCache entry below rewind point must be purged")
	_, ok = v.pfsCache.Get(scoreCheckpointInterval + 1)
	assert.False(t, ok, "pfsCache entry above rewind point must be purged")
	_, ok = v.cpMatrixCache.Get(uint64(1))
	assert.False(t, ok, "cpMatrixCache entry below rewind point must be purged")
	_, ok = v.cpMatrixCache.Get(scoreCheckpointInterval + 1)
	assert.False(t, ok, "cpMatrixCache entry above rewind point must be purged")
}

// TestRewindDelete_UpdatesLastCheckpoint verifies that after RewindDelete removes
// the highest checkpoint, the lastCheckpoint pointer rolls back to the previous one,
// and that removing all checkpoints deletes the pointer entirely.
func TestRewindDelete_UpdatesLastCheckpoint(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	_ = valset
	db := database.NewMemDB()

	cp1 := scoreCheckpointInterval
	cp2 := 2 * scoreCheckpointInterval

	headers := map[uint64]*types.Header{
		0:   makeHeaderWithRound(0, 0),
		cp1: makeHeaderWithRound(cp1, 0),
		cp2: makeHeaderWithRound(cp2, 0),
	}
	v := newTestModuleWithHeaders(t, valset, db, headers)

	WriteCheckpoint(db, cp1,
		map[common.Address]uint64{addrN(0): 1},
		map[common.Address]map[common.Address]uint64{addrN(1): {addrN(2): 1}},
	)
	WriteCheckpoint(db, cp2,
		map[common.Address]uint64{addrN(0): 2},
		map[common.Address]map[common.Address]uint64{addrN(1): {addrN(2): 2}},
	)
	WriteLastCheckpoint(db, cp2)

	// Rewind to cp1: removes cp2 checkpoint, rolls back last pointer to cp1.
	v.RewindDelete(common.Hash{}, cp1)
	pfs2, _ := ReadCheckpoint(db, cp2)
	assert.Nil(t, pfs2, "cp2 checkpoint should be deleted")
	pfs1, _ := ReadCheckpoint(db, cp1)
	assert.NotNil(t, pfs1, "cp1 checkpoint should survive")
	lastCP, hasCP := ReadLastCheckpoint(db)
	assert.True(t, hasCP)
	assert.Equal(t, cp1, lastCP)

	// Rewind to cp1-1: removes cp1 checkpoint, no surviving checkpoint → delete pointer.
	v.RewindDelete(common.Hash{}, cp1-1)
	pfs1, _ = ReadCheckpoint(db, cp1)
	assert.Nil(t, pfs1, "cp1 checkpoint should be deleted")
	_, hasCP = ReadLastCheckpoint(db)
	assert.False(t, hasCP, "lastCheckpoint should be deleted when no checkpoints survive")
}

// TestRewindDelete_MultiIntervalReorg is a regression test for the bug where
// RemoveScoresAfter used CurrentHeader() as the starting point for checkpoint
// deletion. In production, CurrentHeader() already points to the new (lower)
// head when RewindDelete is called, so calcCheckpointBlock(CurrentHeader()) ≤
// blockNum and the deletion loop would exit immediately without removing any
// stale checkpoints.
//
// This test simulates that scenario: the testChain only contains headers up to
// cp1 (the new head), while checkpoints at cp2 and cp3 remain in the DB from
// the old (longer) chain. RewindDelete must delete both cp2 and cp3.
func TestRewindDelete_MultiIntervalReorg(t *testing.T) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	_ = valset
	db := database.NewMemDB()

	cp1 := scoreCheckpointInterval
	cp2 := 2 * scoreCheckpointInterval
	cp3 := 3 * scoreCheckpointInterval

	// CurrentHeader() returns cp1 — the chain has already rewound to cp1.
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

	// Rewind to cp1: cp2 and cp3 must be deleted; cp1 must survive.
	v.RewindDelete(common.Hash{}, cp1)

	pfs3, _ := ReadCheckpoint(db, cp3)
	assert.Nil(t, pfs3, "cp3 checkpoint should be deleted after 2-interval reorg")
	pfs2, _ := ReadCheckpoint(db, cp2)
	assert.Nil(t, pfs2, "cp2 checkpoint should be deleted after 2-interval reorg")
	pfs1, _ := ReadCheckpoint(db, cp1)
	assert.NotNil(t, pfs1, "cp1 checkpoint should survive")
	lastCP, hasCP := ReadLastCheckpoint(db)
	assert.True(t, hasCP)
	assert.Equal(t, cp1, lastCP)
}
