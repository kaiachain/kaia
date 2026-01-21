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
	"strings"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/rcrowley/go-metrics"
)

type vrank struct {
	miningStartTime           time.Time
	view                      istanbul.View
	committee                 []common.Address
	quorum                    int
	preprepareArrivalTime     time.Duration // node receives only one preprepare from proposer
	commitArrivalTimeMap      map[common.Address]time.Duration
	myRoundChangeTimes        []time.Duration
	roundChangeArrivalTimeMap map[common.Address][]time.Duration // per node address per round
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
	return &vrank{
		myRoundChangeTimes:        make([]time.Duration, MaxRoundChangeCount),
		commitArrivalTimeMap:      make(map[common.Address]time.Duration),
		roundChangeArrivalTimeMap: make(map[common.Address][]time.Duration),
		view:                      istanbul.View{},
		committee:                 []common.Address{},
	}
}

func (v *vrank) StartTimer() {
	v.miningStartTime = time.Now()
	v.preprepareArrivalTime = time.Duration(0)
	v.myRoundChangeTimes = make([]time.Duration, MaxRoundChangeCount)
	v.commitArrivalTimeMap = make(map[common.Address]time.Duration)
	v.roundChangeArrivalTimeMap = make(map[common.Address][]time.Duration)
	v.view = istanbul.View{}
	v.committee = []common.Address{}
	v.quorum = 0
}

func (v *vrank) SetLatestView(view istanbul.View, committee []common.Address, quorum int) {
	v.view = view
	v.committee = committee
	v.quorum = quorum
}

func (v *vrank) AddPreprepare(src common.Address, timestamp time.Time) {
	if v.preprepareArrivalTime == time.Duration(0) {
		v.preprepareArrivalTime = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddCommit(src common.Address, timestamp time.Time) {
	if _, exists := v.commitArrivalTimeMap[src]; !exists {
		v.commitArrivalTimeMap[src] = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddMyRoundChange(round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount || round == 0 {
		return
	}
	if v.myRoundChangeTimes[round-1] == time.Duration(0) {
		v.myRoundChangeTimes[round-1] = timestamp.Sub(v.miningStartTime)
	}
}

func (v *vrank) AddRoundChange(src common.Address, round uint64, timestamp time.Time) {
	if round > MaxRoundChangeCount || round == 0 {
		return
	}
	if _, exists := v.roundChangeArrivalTimeMap[src]; !exists {
		v.roundChangeArrivalTimeMap[src] = make([]time.Duration, MaxRoundChangeCount)
	}
	if v.roundChangeArrivalTimeMap[src][round-1] == time.Duration(0) {
		v.roundChangeArrivalTimeMap[src][round-1] = timestamp.Sub(v.miningStartTime)
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

func (v *vrank) buildLogData() (seq int64, round int64, preprepareArrivalTime string, commitArrivalTimes []string, myRoundChangeTimes []string, roundChangeArrivalTimes []string) {
	sortedCommittee := valset.NewAddressSet(v.committee).List()
	preprepareArrivalTime = "-"
	if v.preprepareArrivalTime != time.Duration(0) {
		preprepareArrivalTime = encodeDuration(v.preprepareArrivalTime)
	}

	// Build commitArrivalTimes: [18 18 - 18]
	commitArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		commitTime := "-" // not arrived
		if t, ok := v.commitArrivalTimeMap[addr]; ok {
			commitTime = encodeDuration(t)
		}
		commitArrivalTimes[i] = commitTime
	}

	// Build myRoundChangeTimes: [] or [10000 20000 600000]
	for _, t := range v.myRoundChangeTimes {
		if t != time.Duration(0) {
			myRoundChangeTimes = append(myRoundChangeTimes, encodeDuration(t))
		} else {
			myRoundChangeTimes = append(myRoundChangeTimes, "-")
		}
	}

	// Build roundChangeArrivalTimes: [- - - -] or [10123,20456,60789 10321,20654,60987 ...]
	roundChangeArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		times, exists := v.roundChangeArrivalTimeMap[addr]
		if !exists || len(times) == 0 {
			roundChangeArrivalTimes[i] = "-"
			continue
		}
		// Collect non-zero times for this validator (comma-separated for each round)
		var validatorTimes []string
		for _, t := range times {
			if t != time.Duration(0) {
				validatorTimes = append(validatorTimes, encodeDuration(t))
			} else {
				validatorTimes = append(validatorTimes, "-")
			}
		}
		if len(validatorTimes) == 0 {
			roundChangeArrivalTimes[i] = "-"
		} else {
			roundChangeArrivalTimes[i] = strings.Join(validatorTimes, ",")
		}
	}

	return v.view.Sequence.Int64(), v.view.Round.Int64(), preprepareArrivalTime, commitArrivalTimes, myRoundChangeTimes, roundChangeArrivalTimes
}

func (v *vrank) calcMetrics() (int64, int64, int64, int64) {
	var firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum int64
	if len(v.commitArrivalTimeMap) > 0 {
		_, arrivalTimes := sortByArrivalTimes(v.commitArrivalTimeMap)
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
	if v.preprepareArrivalTime != time.Duration(0) {
		vrankFirstPreprepareArrivalTimeGauge.Update(int64(v.preprepareArrivalTime))
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
