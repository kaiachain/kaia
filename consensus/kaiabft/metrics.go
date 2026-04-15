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

package kaiabft

import (
	"github.com/rcrowley/go-metrics"
)

// Consensus core metrics. Names mirror consensus/istanbul/core/* so dashboards
// can cover both engines with a regex panel.
var (
	roundMeter         = metrics.NewRegisteredMeter("consensus/kaiabft/core/round", nil)
	currentRoundGauge  = metrics.NewRegisteredGauge("consensus/kaiabft/core/currentRound", nil)
	sequenceMeter      = metrics.NewRegisteredMeter("consensus/kaiabft/core/sequence", nil)
	consensusTimeGauge = metrics.NewRegisteredGauge("consensus/kaiabft/core/timer", nil)
	councilSizeGauge   = metrics.NewRegisteredGauge("consensus/kaiabft/core/councilSize", nil)
	committeeSizeGauge = metrics.NewRegisteredGauge("consensus/kaiabft/core/committeeSize", nil)
	hashLockGauge      = metrics.NewRegisteredGauge("consensus/kaiabft/core/hashLock", nil)
)
