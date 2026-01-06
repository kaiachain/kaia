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
	"math"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/rcrowley/go-metrics"
)

type Vrank struct {
	startTime            time.Time
	view                 istanbul.View
	committee            []common.Address
	threshold            time.Duration
	commitArrivalTimeMap map[common.Address]time.Duration

	// metrics
	firstCommit           int64
	quorumCommit          int64
	avgCommitWithinQuorum int64
	lastCommit            int64
}

var (
	// VRank metrics
	vrankFirstCommitArrivalTimeGauge           = metrics.NewRegisteredGauge("vrank/first_commit", nil)
	vrankQuorumCommitArrivalTimeGauge          = metrics.NewRegisteredGauge("vrank/quorum_commit", nil)
	vrankAvgCommitArrivalTimeWithinQuorumGauge = metrics.NewRegisteredGauge("vrank/avg_commit_within_quorum", nil)
	vrankLastCommitArrivalTimeGauge            = metrics.NewRegisteredGauge("vrank/last_commit", nil)

	vrankDefaultThreshold = "300ms" // the time to receive 2f+1 commits in an ideal network

	VRankLogFrequency = uint64(0) // Will be set to the value of VRankLogFrequencyFlag in SetKaiaConfig()

	vrank *Vrank
)

const (
	vrankArrivedEarly = iota
	vrankArrivedLate
	vrankNotArrived
)

const (
	vrankNotArrivedPlaceholder = -1
)

func NewVrank(view istanbul.View, committee []common.Address) *Vrank {
	threshold, _ := time.ParseDuration(vrankDefaultThreshold)
	return &Vrank{
		startTime:             time.Now(),
		view:                  view,
		committee:             committee,
		threshold:             threshold,
		firstCommit:           int64(0),
		quorumCommit:          int64(0),
		avgCommitWithinQuorum: int64(0),
		lastCommit:            int64(0),
		commitArrivalTimeMap:  make(map[common.Address]time.Duration),
	}
}

func (v *Vrank) TimeSinceStart() time.Duration {
	return time.Now().Sub(v.startTime)
}

func (v *Vrank) AddCommit(msg *istanbul.Subject, src common.Address) {
	if v.isTargetCommit(msg, src) {
		t := v.TimeSinceStart()
		v.commitArrivalTimeMap[src] = t
	}
}

func (v *Vrank) HandleCommitted(blockNum *big.Int) {
	if v.view.Sequence.Cmp(blockNum) != 0 {
		return
	}

	if len(v.commitArrivalTimeMap) != 0 {
		sum := int64(0)
		firstCommitTime := time.Duration(math.MaxInt64)
		quorumCommitTime := time.Duration(0)
		for _, arrivalTime := range v.commitArrivalTimeMap {
			sum += int64(arrivalTime)
			if firstCommitTime > arrivalTime {
				firstCommitTime = arrivalTime
			}
			if quorumCommitTime < arrivalTime {
				quorumCommitTime = arrivalTime
			}
		}
		avg := sum / int64(len(v.commitArrivalTimeMap))
		v.avgCommitWithinQuorum = avg
		v.firstCommit = int64(firstCommitTime)
		v.quorumCommit = int64(quorumCommitTime)

		if quorumCommitTime != time.Duration(0) && v.threshold > quorumCommitTime {
			v.threshold = quorumCommitTime
		}
	}
}

// Log logs accumulated data in a compressed form
func (v *Vrank) Log() {
	v.updateMetrics()

	// Skip logging if VRankLogFrequency is 0 or not in the logging frequency
	if VRankLogFrequency == 0 || v.view.Sequence.Uint64()%VRankLogFrequency != 0 {
		return
	}

	logger.Info("VRank", "seq", v.view.Sequence.Int64(),
		"round", v.view.Round.Int64(),
	)
}

func (v *Vrank) updateMetrics() {
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

func (v *Vrank) isTargetCommit(msg *istanbul.Subject, src common.Address) bool {
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
// The returned string is at most 4 bytes
func encodeDuration(d time.Duration) string {
	if d > 10*time.Second {
		return fmt.Sprintf("%.0fs", d.Seconds())
	} else if d > time.Second {
		return fmt.Sprintf("%.1fs", d.Seconds())
	} else {
		return fmt.Sprintf("%d", d.Milliseconds())
	}
}

func encodeDurationBatch(ds []time.Duration) []string {
	ret := make([]string, len(ds))
	for i, d := range ds {
		ret[i] = encodeDuration(d)
	}
	return ret
}
