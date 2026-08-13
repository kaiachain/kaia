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

package core

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/stretchr/testify/assert"
)

func TestProposalNumberMatchesView(t *testing.T) {
	mk := func(num, seq *big.Int) *bft.Preprepare {
		return &bft.Preprepare{
			View:     &bft.View{Round: big.NewInt(0), Sequence: seq},
			Proposal: types.NewBlockWithHeader(&types.Header{Number: num}),
		}
	}
	overflow := new(big.Int).Add(new(big.Int).Lsh(big.NewInt(1), 64), big.NewInt(5)) // 2^64+5

	assert.True(t, proposalNumberMatchesView(mk(big.NewInt(5), big.NewInt(5))), "matching number should pass")
	assert.False(t, proposalNumberMatchesView(mk(big.NewInt(9), big.NewInt(5))), "in-range mismatch should fail")
	assert.False(t, proposalNumberMatchesView(mk(overflow, big.NewInt(5))), "out-of-uint64 mismatch should fail")
}
