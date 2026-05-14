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

package valset

import (
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

var (
	addr1 = common.HexToAddress("0x0001")
	addr2 = common.HexToAddress("0x0002")
	addr3 = common.HexToAddress("0x0003")
)

func TestMarkSuspended(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive},
		addr3: {State: ValPaused},
	}
	m.MarkSuspended([]common.Address{addr1, addr3})

	assert.True(t, m[addr1].Suspended, "addr1 should be suspended")
	assert.False(t, m[addr2].Suspended, "addr2 should not be suspended")
	assert.True(t, m[addr3].Suspended, "addr3 should be suspended")
}

func TestMarkSuspended_Empty(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive},
	}
	m.MarkSuspended(nil)

	assert.False(t, m[addr1].Suspended, "no suspended list → all false")
}

func TestExcludeSuspended(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive, Suspended: true},
		addr2: {State: ValActive, Suspended: false},
		addr3: {State: ValActive, Suspended: true},
	}
	filtered := m.ExcludeSuspended()

	assert.Len(t, filtered, 1)
	assert.NotNil(t, filtered[addr2])
	assert.Nil(t, filtered[addr1])
	assert.Nil(t, filtered[addr3])
}

func TestCountByState(t *testing.T) {
	m := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive},
		addr3: {State: ValPaused},
	}
	assert.Equal(t, uint64(2), m.CountByState(ValActive))
	assert.Equal(t, uint64(1), m.CountByState(ValPaused))
	assert.Equal(t, uint64(0), m.CountByState(ValExiting))
}

func TestCouncil(t *testing.T) {
	addr4 := common.HexToAddress("0x0004")
	addr5 := common.HexToAddress("0x0005")
	m := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValPaused},
		addr3: {State: ValReady},
		addr4: {State: ValInactive},
		addr5: {State: CandTesting},
	}
	council := m.Council()
	assert.Len(t, council, 2)
	assert.NotNil(t, council[addr1], "ValActive in council")
	assert.NotNil(t, council[addr2], "ValPaused in council")
	assert.Nil(t, council[addr3], "ValReady not in council")
	assert.Nil(t, council[addr4], "ValInactive not in council")
	assert.Nil(t, council[addr5], "CandTesting not in council")
}

func TestCommittee(t *testing.T) {
	addr4 := common.HexToAddress("0x0004")
	m := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValActive, Suspended: true},
		addr3: {State: ValPaused},
		addr4: {State: ValReady},
	}
	committee := m.Committee()
	assert.Len(t, committee, 1)
	assert.NotNil(t, committee[addr1], "unsuspended ValActive in committee")
	assert.Nil(t, committee[addr2], "suspended ValActive not in committee")
	assert.Nil(t, committee[addr3], "ValPaused not in committee")
	assert.Nil(t, committee[addr4], "ValReady not in committee")
}

func TestCNPeers(t *testing.T) {
	addr4 := common.HexToAddress("0x0004")
	addr5 := common.HexToAddress("0x0005")
	addr6 := common.HexToAddress("0x0006")
	addr7 := common.HexToAddress("0x0007")
	m := NodeMap{
		addr1: {State: ValActive},
		addr2: {State: ValReady},
		addr3: {State: ValPaused},
		addr4: {State: CandReady},
		addr5: {State: CandTesting},
		addr6: {State: ValInactive},
		addr7: {State: Registered},
	}
	peers := m.CNPeers()
	assert.Len(t, peers, 5)
	assert.NotNil(t, peers[addr1], "ValActive in CNPeers")
	assert.NotNil(t, peers[addr2], "ValReady in CNPeers")
	assert.NotNil(t, peers[addr3], "ValPaused in CNPeers")
	assert.NotNil(t, peers[addr4], "CandReady in CNPeers")
	assert.NotNil(t, peers[addr5], "CandTesting in CNPeers")
	assert.Nil(t, peers[addr6], "ValInactive not in CNPeers")
	assert.Nil(t, peers[addr7], "Registered not in CNPeers")
}

func TestIsRewardEligible(t *testing.T) {
	assert.True(t, ValActive.IsRewardEligible())
	assert.False(t, ValPaused.IsRewardEligible())
	assert.False(t, ValReady.IsRewardEligible())
	assert.False(t, ValInactive.IsRewardEligible())
	assert.False(t, ValExiting.IsRewardEligible())
	assert.False(t, CandReady.IsRewardEligible())
	assert.False(t, CandTesting.IsRewardEligible())
	assert.False(t, Registered.IsRewardEligible())
}

// TestCopy_DeepCopy verifies the deep-copy property that pure transition methods
// (ApplyEpoch, ApplyViolation, ApplyTimeout) rely on: mutations through the returned
// map must not affect the original. If Copy() were shallow, ApplyEpoch's internal
// mutations would leak back to the input and break the pure-function contract.
func TestCopy_DeepCopy(t *testing.T) {
	t1 := time.Unix(1000, 0)
	original := NodeMap{
		addr1: {State: ValActive, StakingAmount: 5_000_000, IdleTimeout: t1, Suspended: false},
		addr2: {State: ValPaused, StakingAmount: 6_000_000, PausedTimeout: t1, Suspended: true},
	}

	cp := original.Copy()

	// Sanity: copy is structurally equal.
	assert.Equal(t, len(original), len(cp))
	assert.Equal(t, original[addr1].State, cp[addr1].State)
	assert.Equal(t, original[addr2].StakingAmount, cp[addr2].StakingAmount)

	// Mutate every field through the copy.
	cp[addr1].State = ValExiting
	cp[addr1].StakingAmount = 999
	cp[addr1].IdleTimeout = time.Unix(2000, 0)
	cp[addr1].Suspended = true
	cp[addr2].State = ValInactive
	cp[addr2].PausedTimeout = time.Unix(3000, 0)

	// Original must be untouched.
	assert.Equal(t, ValActive, original[addr1].State, "Copy must be deep: State")
	assert.Equal(t, uint64(5_000_000), original[addr1].StakingAmount, "Copy must be deep: StakingAmount")
	assert.Equal(t, t1, original[addr1].IdleTimeout, "Copy must be deep: IdleTimeout")
	assert.False(t, original[addr1].Suspended, "Copy must be deep: Suspended")
	assert.Equal(t, ValPaused, original[addr2].State, "Copy must be deep: State (addr2)")
	assert.Equal(t, t1, original[addr2].PausedTimeout, "Copy must be deep: PausedTimeout")

	// Inserting into the copy must not appear in the original.
	cp[addr3] = &Node{State: ValReady}
	assert.Nil(t, original[addr3], "Copy must not share map backing")
}

func TestCopy_NilReceiver(t *testing.T) {
	var v NodeMap
	assert.Nil(t, v.Copy())
}

func TestCopy_NilValuePreserved(t *testing.T) {
	v := NodeMap{addr1: nil, addr2: {State: ValActive}}
	cp := v.Copy()
	assert.Nil(t, cp[addr1])
	assert.NotNil(t, cp[addr2])
}
