// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package core

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/stretchr/testify/assert"
)

func TestVrank(t *testing.T) {
	var (
		N             = 6
		quorum        = 4
		committee, _  = genValidators(N)
		view          = istanbul.View{Sequence: big.NewInt(1), Round: big.NewInt(2)}
		preprepareMsg = &istanbul.Preprepare{View: &view}
		commitMsg     = &istanbul.Subject{View: &view}
		vrank         = NewVrank(view, committee, quorum)
	)

	vrank.StartTimer()
	time.Sleep(1 * time.Millisecond)

	for i := 0; i < quorum; i++ {
		r := (1 + time.Duration(rand.Int63n(10))) * time.Millisecond
		time.Sleep(r)
		vrank.AddPreprepare(preprepareMsg, committee[i])
		time.Sleep(r)
		vrank.AddCommit(commitMsg, committee[i])
	}

	// late messages
	for i := quorum; i < N; i++ {
		r := (1 + time.Duration(rand.Int63n(10))) * time.Millisecond
		time.Sleep(r)
		vrank.AddPreprepare(preprepareMsg, committee[i])
		time.Sleep(r)
		vrank.AddCommit(commitMsg, committee[i])
	}

	vrank.Log()

	assert.NotEqual(t, vrank.preprepareArrivalTime, int64(0))

	firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum := vrank.calcMetrics()
	assert.NotEqual(t, firstCommit, int64(0))
	assert.NotEqual(t, quorumCommit, int64(0))
	assert.NotEqual(t, avgCommitWithinQuorum, int64(0))
	assert.NotEqual(t, lastCommit, int64(0))
	assert.Equal(t, N, len(vrank.commitArrivalTimeMap))

	seq, round, _, commitArrivalTimes := vrank.buildLogData()
	assert.Equal(t, view.Sequence.Int64(), seq)
	assert.Equal(t, view.Round.Int64(), round)
	t.Logf("commitArrivalTimes: %v", commitArrivalTimes)
	assert.Equal(t, N, len(commitArrivalTimes))
}
