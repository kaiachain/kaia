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

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckpointRoundTrip_PreservesZeroFailureCandidates(t *testing.T) {
	db := database.NewMemDB()
	cp := testCheckpointInterval
	P1, C1, C2 := addrN(1), addrN(10), addrN(11)

	WriteCheckpoint(db, cp,
		map[common.Address]uint64{P1: 1},
		vrank.CPMatrix{
			C1: {P1: 2},
			C2: {},
		},
	)

	pfs := ReadCheckpointPFS(db, cp)
	cpMatrix := ReadCheckpointCPMatrix(db, cp)
	require.NotNil(t, pfs)
	require.NotNil(t, cpMatrix)
	assert.Equal(t, uint64(1), pfs[P1])
	assert.Contains(t, cpMatrix, C1)
	assert.Contains(t, cpMatrix, C2)
	assert.Equal(t, uint64(2), cpMatrix[C1][P1])
	assert.Equal(t, uint64(0), cpMatrix[C2][P1])
}
