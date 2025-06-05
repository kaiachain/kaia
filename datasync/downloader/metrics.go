// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/downloader/metrics.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package downloader

import (
	kaiametrics "github.com/kaiachain/kaia/v2/metrics"
	"github.com/rcrowley/go-metrics"
)

var (
	headerInMeter      = metrics.NewRegisteredMeter("klay/downloader/headers/in", nil)
	headerReqTimer     = kaiametrics.NewRegisteredHybridTimer("klay/downloader/headers/req", nil)
	headerDropMeter    = metrics.NewRegisteredMeter("klay/downloader/headers/drop", nil)
	headerTimeoutMeter = metrics.NewRegisteredMeter("klay/downloader/headers/timeout", nil)

	bodyInMeter      = metrics.NewRegisteredMeter("klay/downloader/bodies/in", nil)
	bodyReqTimer     = kaiametrics.NewRegisteredHybridTimer("klay/downloader/bodies/req", nil)
	bodyDropMeter    = metrics.NewRegisteredMeter("klay/downloader/bodies/drop", nil)
	bodyTimeoutMeter = metrics.NewRegisteredMeter("klay/downloader/bodies/timeout", nil)

	receiptInMeter      = metrics.NewRegisteredMeter("klay/downloader/receipts/in", nil)
	receiptReqTimer     = kaiametrics.NewRegisteredHybridTimer("klay/downloader/receipts/req", nil)
	receiptDropMeter    = metrics.NewRegisteredMeter("klay/downloader/receipts/drop", nil)
	receiptTimeoutMeter = metrics.NewRegisteredMeter("klay/downloader/receipts/timeout", nil)

	stakingInfoInMeter      = metrics.NewRegisteredMeter("klay/downloader/stakinginfo/in", nil)
	stakingInfoReqTimer     = kaiametrics.NewRegisteredHybridTimer("klay/downloader/stakinginfo/req", nil)
	stakingInfoDropMeter    = metrics.NewRegisteredMeter("klay/downloader/stakinginfo/drop", nil)
	stakingInfoTimeoutMeter = metrics.NewRegisteredMeter("klay/downloader/stakinginfo/timeout", nil)

	stateInMeter   = metrics.NewRegisteredMeter("klay/downloader/states/in", nil)
	stateDropMeter = metrics.NewRegisteredMeter("klay/downloader/states/drop", nil)

	throttleCounter = metrics.NewRegisteredCounter("klay/downloader/throttle", nil)
)
