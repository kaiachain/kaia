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
	"math/big"
	"slices"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/rcrowley/go-metrics"
)

type vrank struct {
	roundStartTime           time.Time
	view                     istanbul.View
	committee                []common.Address
	preprepareArrivalTimeMap map[common.Address]time.Duration
	commitArrivalTimeMap     map[common.Address]time.Duration

	// metrics
	firstPreprepare           int64
	quorumPreprepare          int64
	avgPreprepareWithinQuorum int64
	lastPreprepare            int64

	firstCommit           int64
	quorumCommit          int64
	avgCommitWithinQuorum int64
	lastCommit            int64
}

var (
	// VRank metrics
	vrankFirstPreprepareArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_preprepare", nil)
	vrankQuorumPreprepareArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_preprepare", nil)
	vrankAvgPreprepareArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_preprepare_within_quorum", nil)
	vrankLastPreprepareArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_preprepare", nil)

	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	VRankLogFrequency = uint64(0) // Will be set to the value of VRankLogFrequencyFlag in SetKaiaConfig()

	Vrank *vrank // vrank instance is newly created every time a new round starts
)

const (
	vrankArrivedEarly = iota
	vrankArrivedLate
	vrankNotArrived
)

const (
	vrankNotArrivedPlaceholder = -1
)

func NewVrank(view istanbul.View, committee []common.Address) *vrank {
	return &vrank{
		roundStartTime:           time.Now(),
		view:                     view,
		committee:                committee,
		preprepareArrivalTimeMap: make(map[common.Address]time.Duration),
		commitArrivalTimeMap:     make(map[common.Address]time.Duration),
		firstCommit:              int64(0),
		quorumCommit:             int64(0),
		avgCommitWithinQuorum:    int64(0),
	}
}

func (v *vrank) TimeSinceRoundStart() time.Duration {
	return time.Now().Sub(v.roundStartTime)
}

func (v *vrank) AddPreprepare(msg *istanbul.Preprepare, src common.Address) {
	if v.isTargetPreprepare(msg, src) {
		t := v.TimeSinceRoundStart()
		v.preprepareArrivalTimeMap[src] = t
	}
}

func (v *vrank) AddCommit(msg *istanbul.Subject, src common.Address) {
	if v.isTargetCommit(msg, src) {
		t := v.TimeSinceRoundStart()
		v.commitArrivalTimeMap[src] = t
	}
}

// HandlePreprepared is called once when the state is changed to Preprepared
func (v *vrank) HandlePreprepared(blockNum *big.Int) {
	if v.view.Sequence.Cmp(blockNum) != 0 {
		return
	}

	if len(v.preprepareArrivalTimeMap) > 0 {
		_, arrivalTimes := sortByArrivalTimes(v.preprepareArrivalTimeMap)
		v.firstPreprepare = arrivalTimes[0]
		v.quorumPreprepare = arrivalTimes[len(arrivalTimes)-1]

		sum := int64(0)
		for _, arrivalTime := range arrivalTimes {
			sum += int64(arrivalTime)
		}
		v.avgPreprepareWithinQuorum = sum / int64(len(v.preprepareArrivalTimeMap))
	}
}

// HandleCommitted is called once when the state is changed to Committed
func (v *vrank) HandleCommitted(blockNum *big.Int) {
	if v.view.Sequence.Cmp(blockNum) != 0 {
		return
	}

	if len(v.commitArrivalTimeMap) > 0 {
		_, arrivalTimes := sortByArrivalTimes(v.commitArrivalTimeMap)
		v.firstCommit = arrivalTimes[0]
		v.quorumCommit = arrivalTimes[len(arrivalTimes)-1]

		sum := int64(0)
		for _, arrivalTime := range arrivalTimes {
			sum += int64(arrivalTime)
		}
		v.avgCommitWithinQuorum = sum / int64(len(v.commitArrivalTimeMap))
	}
}

// Log logs accumulated data in a compressed form
func (v *vrank) Log() {
	if len(v.preprepareArrivalTimeMap) > 0 {
		_, preprepareArrivalTimes := sortByArrivalTimes(v.preprepareArrivalTimeMap)
		v.lastPreprepare = preprepareArrivalTimes[len(preprepareArrivalTimes)-1]
	}

	if len(v.commitArrivalTimeMap) > 0 {
		_, commitArrivalTimes := sortByArrivalTimes(v.commitArrivalTimeMap)
		v.lastCommit = commitArrivalTimes[len(commitArrivalTimes)-1]
	}

	v.updateMetrics()

	// Skip logging if VRankLogFrequency is 0 or not in the logging frequency
	if VRankLogFrequency == 0 || v.view.Sequence.Uint64()%VRankLogFrequency != 0 {
		return
	}

	seq, round, msgArrivalTimes := v.buildLogData()
	logger.Warn("VRank", "seq", seq, "round", round, "msgArrivalTimes(preprepare,commit)", msgArrivalTimes)
}

func (v *vrank) buildLogData() (seq int64, round int64, msgArrivalTimes []string) {
	sortedCommittee := valset.NewAddressSet(v.committee).List()
	msgArrivalTimes = make([]string, len(sortedCommittee))
	for i, addr := range sortedCommittee {
		preprepareTime := "-" // not arrived
		if t, ok := v.preprepareArrivalTimeMap[addr]; ok {
			preprepareTime = encodeDuration(t)
		}
		commitTime := "-" // not arrived
		if t, ok := v.commitArrivalTimeMap[addr]; ok {
			commitTime = encodeDuration(t)
		}
		// Format: "(preprepare,commit)"
		msgArrivalTimes[i] = fmt.Sprintf("(%s,%s)", preprepareTime, commitTime)
	}

	return v.view.Sequence.Int64(), v.view.Round.Int64(), msgArrivalTimes
}

func (v *vrank) updateMetrics() {
	if v.firstPreprepare != int64(0) {
		vrankFirstPreprepareArrivalTimeGauge.Update(v.firstPreprepare)
	}
	if v.quorumPreprepare != int64(0) {
		vrankQuorumPreprepareArrivalTimeGauge.Update(v.quorumPreprepare)
	}
	if v.avgPreprepareWithinQuorum != int64(0) {
		vrankAvgPreprepareArrivalTimeWithinQuorumGauge.Update(v.avgPreprepareWithinQuorum)
	}
	if v.lastPreprepare != int64(0) {
		vrankLastPreprepareArrivalTimeGauge.Update(v.lastPreprepare)
	}

	if v.firstCommit != int64(0) {
		vrankFirstCommitArrivalTimeGauge.Update(v.firstCommit)
	}
	if v.quorumCommit != int64(0) {
		vrankQuorumCommitArrivalTimeGauge.Update(v.quorumCommit)
	}
	if v.avgCommitWithinQuorum != int64(0) {
		vrankAvgCommitArrivalTimeWithinQuorumGauge.Update(v.avgCommitWithinQuorum)
	}
	if v.lastCommit != int64(0) {
		vrankLastCommitArrivalTimeGauge.Update(v.lastCommit)
	}
}

func (v *vrank) isTargetPreprepare(msg *istanbul.Preprepare, src common.Address) bool {
	if msg.View == nil || msg.View.Sequence == nil || msg.View.Round == nil {
		return false
	}
	if msg.View.Cmp(&v.view) != 0 {
		return false
	}
	_, ok := v.preprepareArrivalTimeMap[src]
	if ok {
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
