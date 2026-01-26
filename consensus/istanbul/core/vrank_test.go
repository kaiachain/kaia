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
		N            = 6
		quorum       = 4
		committee, _ = genValidators(N)
		testRound    = uint64(2)
		view         = istanbul.View{Sequence: big.NewInt(1), Round: big.NewInt(int64(testRound))}
		vrank        = NewVrank()
	)

	vrank.StartTimer()
	vrank.SetLatestView(view, committee, quorum)
	time.Sleep(1 * time.Millisecond)

	// Simulate messages for each round up to testRound
	for round := uint64(0); round <= testRound; round++ {
		for i := 0; i < quorum; i++ {
			r := (1 + time.Duration(rand.Int63n(10))) * time.Millisecond
			time.Sleep(r)
			vrank.AddPreprepare(committee[i], round, time.Now())
			time.Sleep(r)
			vrank.AddCommit(committee[i], round, time.Now())
		}

		// late messages
		for i := quorum; i < N; i++ {
			r := (1 + time.Duration(rand.Int63n(10))) * time.Millisecond
			time.Sleep(r)
			vrank.AddPreprepare(committee[i], round, time.Now())
			time.Sleep(r)
			vrank.AddCommit(committee[i], round, time.Now())
		}

		// round change messages (not for round 0)
		if round > 0 {
			for i := 0; i < N; i++ {
				r := time.Duration(round*10+uint64(rand.Int63n(10))) * time.Millisecond
				vrank.AddRoundChange(committee[i], round, time.Now().Add(r))
			}
			r := time.Duration(round*10+uint64(rand.Int63n(10))) * time.Millisecond
			vrank.AddMyRoundChange(round, time.Now().Add(r))
		}
	}

	vrank.Log()

	assert.NotEqual(t, vrank.timestamps[testRound].preprepareArrivalTime, int64(0))

	firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum := vrank.calcMetrics()
	assert.NotEqual(t, firstCommit, int64(0))
	assert.NotEqual(t, quorumCommit, int64(0))
	assert.NotEqual(t, avgCommitWithinQuorum, int64(0))
	assert.NotEqual(t, lastCommit, int64(0))
	// Count entries in sync.Map
	commitCount := 0
	vrank.timestamps[testRound].commitArrivalTimeMap.Range(func(_, _ any) bool {
		commitCount++
		return true
	})
	assert.Equal(t, N, commitCount)

	seq, round, preprepareArrivalTimes, commitArrivalTimes, myRoundChangeTimes, roundChangeArrivalTimes := vrank.buildLogData()
	assert.Equal(t, view.Sequence.Int64(), seq)
	assert.Equal(t, view.Round.Int64(), round)
	t.Logf("preprepareArrivalTimes: %v", preprepareArrivalTimes)
	t.Logf("commitArrivalTimes: %v", commitArrivalTimes)
	t.Logf("myRoundChangeTimes: %v", myRoundChangeTimes)
	t.Logf("roundChangeArrivalTimes: %v", roundChangeArrivalTimes)
}
