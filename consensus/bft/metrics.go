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

package bft

import (
	"github.com/rcrowley/go-metrics"
)

// CoreMetrics bundles the consensus-core instrumentation shared by the
// istanbul and kaiabft engines. Both emit under the same namespace
// ("consensus/istanbul/core") so dashboards cover them uniformly.
type CoreMetrics struct {
	Round         metrics.Meter
	CurrentRound  metrics.Gauge
	Sequence      metrics.Meter
	ConsensusTime metrics.Gauge
	CouncilSize   metrics.Gauge
	CommitteeSize metrics.Gauge
	HashLock      metrics.Gauge
}

// NewCoreMetrics registers core metrics under prefix and returns them bundled.
func NewCoreMetrics(prefix string) *CoreMetrics {
	return &CoreMetrics{
		Round:         metrics.NewRegisteredMeter(prefix+"/round", nil),
		CurrentRound:  metrics.NewRegisteredGauge(prefix+"/currentRound", nil),
		Sequence:      metrics.NewRegisteredMeter(prefix+"/sequence", nil),
		ConsensusTime: metrics.NewRegisteredGauge(prefix+"/timer", nil),
		CouncilSize:   metrics.NewRegisteredGauge(prefix+"/councilSize", nil),
		CommitteeSize: metrics.NewRegisteredGauge(prefix+"/committeeSize", nil),
		HashLock:      metrics.NewRegisteredGauge(prefix+"/hashLock", nil),
	}
}
