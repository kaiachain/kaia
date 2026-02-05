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

package common

import "math"

// CalcQuorumSize returns the minimum number of messages required to proceed (QBFT quorum).
// qualifiedLen is the number of qualified validators, committeeSize is the committee size.
// For less than 4 validators, quorum size equals the effective size;
// otherwise it is ceil(2 * min(qualifiedLen, committeeSize) / 3).
func CalcQuorumSize(qualifiedLen int, committeeSize uint64) int {
	size := qualifiedLen
	if size > int(committeeSize) {
		size = int(committeeSize)
	}
	if size < 4 {
		return size
	}
	return int(math.Ceil(float64(2*size) / 3))
}

// CalcFaultTolerance returns the maximum endurable number of byzantine fault nodes.
// qualifiedLen is the number of qualified validators, committeeSize is the committee size.
func CalcFaultTolerance(qualifiedLen int, committeeSize uint64) int {
	if qualifiedLen > int(committeeSize) {
		return int(math.Ceil(float64(committeeSize)/3)) - 1
	}
	return int(math.Ceil(float64(qualifiedLen)/3)) - 1
}
