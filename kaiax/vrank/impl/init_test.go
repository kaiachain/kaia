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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInit_CatchUpFromCheckpoint verifies that Init replays only the tail blocks after the
// stored DB checkpoint and correctly populates both pfsCache and cpMatrixCache.
//
// The tail block (checkpoint+1) has round=1 and cfReport=[C1]:
//   - PFS:      pfReport=[P1] (proposer at round 0) → pfs[P1] incremented
//   - cpMatrix: reporter=P2   (proposer at round 1) → cpMatrix[C1][P2] incremented
//
// Both caches are seeded from the checkpoint (P1:1 / C1:{P1:1}) so the assertions
// confirm both the carry-over from the checkpoint AND the tail-block contribution.
func TestInit_CatchUpFromCheckpoint(t *testing.T) {
	const checkpoint = scoreCheckpointInterval // must be a scoreCheckpointInterval multiple
	P1, P2, C1 := addrN(1), addrN(2), addrN(10)

	// Only the two headers around the tail are needed: catchUp loads from the DB
	// checkpoint and only processes blockNum = checkpoint+1.
	headers := map[uint64]*types.Header{
		checkpoint:     makeHeaderWithRound(checkpoint, 0),
		checkpoint + 1: makeHeaderWithVRank(checkpoint+1, 1, []common.Address{C1}),
	}

	db := database.NewMemDB()
	WriteCheckpoint(db, checkpoint,
		map[common.Address]uint64{P1: 1},
		map[common.Address]map[common.Address]uint64{C1: {P1: 1}},
	)
	WriteLastCheckpoint(db, checkpoint)

	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	// round=0 → P1 (pfReport proposer); round=1 → P2 (cfReport reporter).
	// GetCandidates is NOT called: the checkpoint branch skips newCPMatrix.
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ uint64, round uint64) (common.Address, error) {
			if round == 1 {
				return P2, nil
			}
			return P1, nil
		},
	).AnyTimes()

	module := NewVRankModule()
	require.NoError(t, module.Init(&InitOpts{
		NodeKey:     key,
		Valset:      valset,
		ChainConfig: params.TestKaiaConfig("permissionless"),
		Chain:       &testChain{headers: headers},
		ChainKv:     db,
	}))

	// pfsCache: checkpoint seeded P1=1; tail block adds 1 more (round=1 → pfReport=[P1]).
	pfsCached, ok := module.pfsCache.Get(checkpoint + 1)
	require.True(t, ok, "pfsCache should be populated at head")
	pfs := pfsCached.(map[common.Address]uint64)
	assert.Equal(t, uint64(2), pfs[P1], "P1 score: 1 from checkpoint + 1 from tail block")

	// cpMatrixCache: checkpoint seeded C1:{P1:1}; tail block adds C1:{P2:1} (reporter=P2 at round=1).
	cpCached, ok := module.cpMatrixCache.Get(checkpoint + 1)
	require.True(t, ok, "cpMatrixCache should be populated at head")
	cpMatrix := cpCached.(map[common.Address]map[common.Address]uint64)
	assert.Equal(t, uint64(1), cpMatrix[C1][P1], "P1 contribution should carry over from checkpoint")
	assert.Equal(t, uint64(1), cpMatrix[C1][P2], "P2 should be credited as reporter in tail block")
}
