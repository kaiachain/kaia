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
	"testing"
	"time"

	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/valset"
)

func TestVrank(t *testing.T) {
	var (
		N            = 6
		quorum       = 4
		committee, _ = genValidators(N)
		view         = istanbul.View{Sequence: big.NewInt(0), Round: big.NewInt(0)}
		msg          = &istanbul.Subject{View: &view}
		vrank        = NewVrank(view, committee)

		expectedAssessList  []uint8
		expectedLateCommits []time.Duration
	)

	committee = valset.NewAddressSet(committee).List() // sort it
	for i := 0; i < quorum; i++ {
		vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedEarly)
	}
	vrank.HandleCommitted(view.Sequence)
	for i := quorum; i < N; i++ {
		vrank.AddCommit(msg, committee[i])
		expectedAssessList = append(expectedAssessList, vrankArrivedLate)
		expectedLateCommits = append(expectedLateCommits, vrank.commitArrivalTimeMap[committee[i]])
	}
}
