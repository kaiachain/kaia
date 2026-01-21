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
	"cmp"
	"maps"
	"slices"
	"strconv"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/rcrowley/go-metrics"
)

type vrank struct {
	miningStartTime time.Time
	view            istanbul.View
	committee       []common.Address
	quorum          int
	timestamps      [MaxRoundChangeCount]msgArrivalTimes
}

type msgArrivalTimes struct {
	preprepareArrivalTime     time.Duration // node receives only one preprepare from proposer
	commitArrivalTimeMap      map[common.Address]time.Duration
	myRoundChangeTime         time.Duration
	roundChangeArrivalTimeMap map[common.Address]time.Duration
}

const (
	DefaultVRankLogFrequency = uint64(60)
	MaxRoundChangeCount      = uint64(8)
)

var (
	// VRank metrics
	vrankFirstPreprepareArrivalTimeGauge       = metrics.NewRegisteredGauge("vrank/first_preprepare", nil)
	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	VRankLogFrequency = DefaultVRankLogFrequency // Will be set to the value of VRankLogFrequencyFlag in SetKaiaConfig()

	Vrank *vrank
)

func NewVrank() *vrank {
	ret := &vrank{
		view:       istanbul.View{},
		committee:  []common.Address{},
		timestamps: [MaxRoundChangeCount]msgArrivalTimes{},
	}
	for i := range ret.timestamps {
		ret.timestamps[i] = *NewMsgArrivalTimes()
	}
	return ret
}

func NewMsgArrivalTimes() *msgArrivalTimes {
	return &msgArrivalTimes{
		preprepareArrivalTime:     time.Duration(0),
		commitArrivalTimeMap:      make(map[common.Address]time.Duration),
		myRoundChangeTime:         time.Duration(0),
		roundChangeArrivalTimeMap: make(map[common.Address]time.Duration),
	}
}

func (v *vrank) StartTimer() {
	v.miningStartTime = time.Now()
	v.view = istanbul.View{}
	v.committee = []common.Address{}
	v.quorum = 0
	for i := range v.timestamps {
		v.timestamps[i] = *NewMsgArrivalTimes()
	}
}

func (v *vrank) SetLatestView(view istanbul.View, committee []common.Address, quorum int) {
	v.view = view
	v.committee = committee
	v.quorum = quorum
}

func (v *vrank) AddPreprepare(src common.Address, round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount {
		return
	}
	if v.timestamps[round].preprepareArrivalTime == time.Duration(0) {
		v.timestamps[round].preprepareArrivalTime = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddCommit(src common.Address, round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount {
		return
	}
	if _, exists := v.timestamps[round].commitArrivalTimeMap[src]; !exists {
		v.timestamps[round].commitArrivalTimeMap[src] = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddMyRoundChange(round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount {
		return
	}
	if v.timestamps[round].myRoundChangeTime == time.Duration(0) {
		v.timestamps[round].myRoundChangeTime = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddRoundChange(src common.Address, round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount {
		return
	}
	if _, exists := v.timestamps[round].roundChangeArrivalTimeMap[src]; !exists {
		v.timestamps[round].roundChangeArrivalTimeMap[src] = timestamp.Sub(v.miningStartTime)
	}
}

// Log logs accumulated data in a compressed form
func (v *vrank) Log() {
	// Skip if no data collected (view not set)
	if v.view.Sequence == nil || v.view.Round == nil {
		return
	}

	v.updateMetrics()

	// Skip logging if VRankLogFrequency is 0 or not in the logging frequency
	if VRankLogFrequency == 0 || v.view.Sequence.Uint64()%VRankLogFrequency != 0 {
		return
	}

	seq, round, preprepareArrivalTime, commitArrivalTimes, myRoundChangeTimes, roundChangeArrivalTimes := v.buildLogData()
	logger.Warn("VRank", "seq", seq, "round", round,
		"preprepareArrivalTime", preprepareArrivalTime,
		"commitArrivalTimes", commitArrivalTimes,
		"myRoundChangeTimes", myRoundChangeTimes,
		"roundChangeArrivalTimes", roundChangeArrivalTimes)
}

func (v *vrank) buildLogData() (seq int64, round int64, preprepareArrivalTimes string, commitArrivalTimes []string, myRoundChangeTimes string, roundChangeArrivalTimes []string) {
	if v.view.Round == nil {
		return
	}
	sortedCommittee := valset.NewAddressSet(v.committee).List()
	maxRound := v.view.Round.Uint64()

	// Initialize per-validator arrays
	commitArrivalTimes = make([]string, len(sortedCommittee))
	roundChangeArrivalTimes = make([]string, len(sortedCommittee))

	// Build incrementally: each round appends with comma
	// round 0: [a1 b1 c1], round 1: [a1,a2 b1,b2 c1,c2], etc.
	for r := uint64(0); r <= maxRound; r++ {
		pp, commits, myRC, rcs := v.timestamps[r].buildLogData(sortedCommittee)
		preprepareArrivalTimes = appendTime(preprepareArrivalTimes, pp)
		myRoundChangeTimes = appendTime(myRoundChangeTimes, myRC)
		for i := range sortedCommittee {
			commitArrivalTimes[i] = appendTime(commitArrivalTimes[i], commits[i])
			roundChangeArrivalTimes[i] = appendTime(roundChangeArrivalTimes[i], rcs[i])
		}
	}

	return v.view.Sequence.Int64(), v.view.Round.Int64(), preprepareArrivalTimes, commitArrivalTimes, myRoundChangeTimes, roundChangeArrivalTimes
}

// appendTime appends a time value to an existing string with comma separator
func appendTime(existing, newVal string) string {
	if existing == "" {
		return newVal
	}
	return existing + "," + newVal
}

func (m *msgArrivalTimes) buildLogData(sortedCommittee []common.Address) (preprepareArrivalTime string, commitArrivalTimes []string, myRoundChangeTime string, roundChangeArrivalTimes []string) {
	// preprepareArrivalTime
	preprepareArrivalTime = "-"
	if m.preprepareArrivalTime != time.Duration(0) {
		preprepareArrivalTime = encodeDuration(m.preprepareArrivalTime)
	}

	// commitArrivalTimes: one per validator
	commitArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		if t, ok := m.commitArrivalTimeMap[addr]; ok && t != time.Duration(0) {
			commitArrivalTimes[i] = encodeDuration(t)
		} else {
			commitArrivalTimes[i] = "-"
		}
	}

	// myRoundChangeTime
	myRoundChangeTime = "-"
	if m.myRoundChangeTime != time.Duration(0) {
		myRoundChangeTime = encodeDuration(m.myRoundChangeTime)
	}

	// roundChangeArrivalTimes: one per validator
	roundChangeArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		if t, ok := m.roundChangeArrivalTimeMap[addr]; ok && t != time.Duration(0) {
			roundChangeArrivalTimes[i] = encodeDuration(t)
		} else {
			roundChangeArrivalTimes[i] = "-"
		}
	}

	return preprepareArrivalTime, commitArrivalTimes, myRoundChangeTime, roundChangeArrivalTimes
}

func (v *vrank) calcMetrics() (int64, int64, int64, int64) {
	if v.view.Round == nil {
		return 0, 0, 0, 0
	}
	round := v.view.Round.Uint64()
	if round >= MaxRoundChangeCount {
		round = MaxRoundChangeCount - 1
	}
	commitMap := v.timestamps[round].commitArrivalTimeMap

	var firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum int64
	if len(commitMap) > 0 {
		_, arrivalTimes := sortByArrivalTimes(commitMap)
		firstCommit = arrivalTimes[0]
		lastCommit = arrivalTimes[len(arrivalTimes)-1]
		if v.quorum > 0 && len(arrivalTimes) >= v.quorum {
			quorumCommit = arrivalTimes[v.quorum-1]
			sum := int64(0)
			for _, arrivalTime := range arrivalTimes[:v.quorum] {
				sum += int64(arrivalTime)
			}
			avgCommitWithinQuorum = sum / int64(v.quorum)
		}
	}

	return firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum
}

func (v *vrank) updateMetrics() {
	if v.view.Round == nil {
		return
	}
	round := v.view.Round.Uint64()
	if round >= MaxRoundChangeCount {
		round = MaxRoundChangeCount - 1
	}
	if v.timestamps[round].preprepareArrivalTime != time.Duration(0) {
		vrankFirstPreprepareArrivalTimeGauge.Update(int64(v.timestamps[round].preprepareArrivalTime))
	}
	firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum := v.calcMetrics()
	if firstCommit != 0 {
		vrankFirstCommitArrivalTimeGauge.Update(firstCommit)
	}
	if lastCommit != 0 {
		vrankLastCommitArrivalTimeGauge.Update(lastCommit)
	}
	if quorumCommit != 0 {
		vrankQuorumCommitArrivalTimeGauge.Update(quorumCommit)
	}
	if avgCommitWithinQuorum != 0 {
		vrankAvgCommitArrivalTimeWithinQuorumGauge.Update(avgCommitWithinQuorum)
	}
}

// encodeDuration encodes given duration into string
func encodeDuration(d time.Duration) string {
	return strconv.FormatInt(d.Milliseconds(), 10)
}

func sortByArrivalTimes(arrivalTimeMap map[common.Address]time.Duration) ([]common.Address, []int64) {
	// Sort addresses by their arrival times
	sortedAddrs := slices.SortedFunc(maps.Keys(arrivalTimeMap), func(a, b common.Address) int {
		return cmp.Compare(arrivalTimeMap[a], arrivalTimeMap[b])
	})

	retTimes := make([]int64, len(sortedAddrs))
	for i, addr := range sortedAddrs {
		retTimes[i] = int64(arrivalTimeMap[addr])
	}

	return sortedAddrs, retTimes
}
