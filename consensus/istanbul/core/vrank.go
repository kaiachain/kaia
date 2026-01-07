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
	"fmt"
	"maps"
	"slices"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/rcrowley/go-metrics"
)

type vrank struct {
	miningStartTime       time.Time
	view                  istanbul.View
	committee             []common.Address
	quorum                int
	preprepareArrivalTime time.Duration // node receives only one preprepare from proposer
	commitArrivalTimeMap  map[common.Address]time.Duration
}

var (
	// VRank metrics
	vrankFirstPreprepareArrivalTimeGauge       = metrics.NewRegisteredGauge("vrank/first_preprepare", nil)
	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	VRankLogFrequency = uint64(0) // Will be set to the value of VRankLogFrequencyFlag in SetKaiaConfig()

	Vrank *vrank
)

const (
	vrankArrivedEarly = iota
	vrankArrivedLate
	vrankNotArrived
)

const (
	vrankNotArrivedPlaceholder = -1
)

func NewVrank(view istanbul.View, committee []common.Address, quorum int) *vrank {
	if quorum == 0 {
		return nil
	}

	return &vrank{
		view:                  view,
		committee:             committee,
		quorum:                quorum,
		commitArrivalTimeMap:  make(map[common.Address]time.Duration),
		preprepareArrivalTime: time.Duration(0),
	}
}

func (v *vrank) StartTimer() {
	v.miningStartTime = time.Now()
}

func (v *vrank) StartNewRound(view istanbul.View) {
	v.view = view
	v.commitArrivalTimeMap = make(map[common.Address]time.Duration)
	v.preprepareArrivalTime = time.Duration(0)
}

func (v *vrank) AddPreprepare(msg *istanbul.Preprepare, src common.Address) {
	if v.isTargetPreprepare(msg, src) {
		v.preprepareArrivalTime = time.Since(v.miningStartTime)
	}
}

func (v *vrank) AddCommit(msg *istanbul.Subject, src common.Address) {
	if v.isTargetCommit(msg, src) {
		v.commitArrivalTimeMap[src] = time.Since(v.miningStartTime)
	}
}

// Log logs accumulated data in a compressed form
func (v *vrank) Log() {
	v.updateMetrics()

	// Skip logging if VRankLogFrequency is 0 or not in the logging frequency
	if VRankLogFrequency == 0 || v.view.Sequence.Uint64()%VRankLogFrequency != 0 {
		return
	}

	seq, round, preprepareArrivalTime, commitArrivalTimes := v.buildLogData()
	logger.Warn("VRank", "seq", seq, "round", round,
		"preprepareArrivalTime", preprepareArrivalTime,
		"commitArrivalTimes", commitArrivalTimes)
}

func (v *vrank) buildLogData() (seq int64, round int64, preprepareArrivalTime string, commitArrivalTimes []string) {
	sortedCommittee := valset.NewAddressSet(v.committee).List()
	preprepareArrivalTime = encodeDuration(v.preprepareArrivalTime)
	commitArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		commitTime := "-" // not arrived
		if t, ok := v.commitArrivalTimeMap[addr]; ok {
			commitTime = encodeDuration(t)
		}
		commitArrivalTimes[i] = fmt.Sprintf("%s", commitTime)
	}

	return v.view.Sequence.Int64(), v.view.Round.Int64(), preprepareArrivalTime, commitArrivalTimes
}

func (v *vrank) calcMetrics() (int64, int64, int64, int64) {
	var firstCommit, lastCommit, quorumCommit, avgCommitWithinQuorum int64
	if len(v.commitArrivalTimeMap) > 0 {
		_, arrivalTimes := sortByArrivalTimes(v.commitArrivalTimeMap)
		firstCommit = arrivalTimes[0]
		lastCommit = arrivalTimes[len(arrivalTimes)-1]
		if len(arrivalTimes) >= v.quorum {
			quorumCommit = arrivalTimes[v.quorum-1]
			sum := int64(0)
			for _, arrivalTime := range arrivalTimes[v.quorum-1:] {
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

func (v *vrank) isTargetPreprepare(msg *istanbul.Preprepare, src common.Address) bool {
	if msg.View == nil || msg.View.Sequence == nil || msg.View.Round == nil {
		return false
	}
	if msg.View.Cmp(&v.view) != 0 {
		return false
	}
	return true
}

func (v *vrank) isTargetCommit(msg *istanbul.Subject, src common.Address) bool {
	if msg.View == nil || msg.View.Sequence == nil || msg.View.Round == nil {
		return false
	}
	if msg.View.Cmp(&v.view) != 0 {
		return false
	}
	_, ok := v.commitArrivalTimeMap[src]
	if ok {
		return false
	}
	return true
}

// encodeDuration encodes given duration into string
func encodeDuration(d time.Duration) string {
	return fmt.Sprintf("%d", d.Milliseconds())
}

func sortByArrivalTimes(arrivalTimeMap map[common.Address]time.Duration) ([]common.Address, []int64) {
	// Sort addresses by their arrival times
	sortedAddrs := slices.SortedFunc(maps.Keys(arrivalTimeMap), func(a, b common.Address) int {
		return int(arrivalTimeMap[a] - arrivalTimeMap[b])
	})

	retTimes := make([]int64, len(sortedAddrs))
	for i, addr := range sortedAddrs {
		retTimes[i] = int64(arrivalTimeMap[addr])
	}

	return sortedAddrs, retTimes
}
