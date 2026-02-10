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

package vrank

import (
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestCollector_GetViewData_UnknownView(t *testing.T) {
	c := NewCollector()
	at, _, m := c.GetViewData(ViewKey{N: 1, R: 0})
	assert.True(t, at.IsZero())
	assert.Nil(t, m)
}

func TestCollector_GetViewData_NoCandMsgAdded(t *testing.T) {
	var (
		c   = NewCollector()
		vk  = ViewKey{N: 1, R: 0}
		now = time.Now()
	)
	c.AddPrepreparedTime(vk, now, common.Hash{})
	at, _, m := c.GetViewData(vk)
	assert.False(t, at.IsZero())
	assert.Equal(t, now.UnixNano(), at.UnixNano())
	assert.Nil(t, m)
}

func TestCollector_AddCandMsg_DuplicateRejected(t *testing.T) {
	var (
		c    = NewCollector()
		vk   = ViewKey{N: 1, R: 0}
		addr = common.HexToAddress("0x01")
		msg  = &VRankCandidate{BlockNumber: 1, Round: 0}
		when = time.Now()
	)

	ok := c.AddCandMsg(vk, addr, when, msg)
	assert.True(t, ok)
	_, _, m := c.GetViewData(vk)
	assert.Len(t, m, 1)
	assert.Equal(t, msg, m[addr].Msg)

	ok = c.AddCandMsg(vk, addr, when.Add(time.Second), msg)
	assert.False(t, ok)
	_, _, m = c.GetViewData(vk)
	assert.Len(t, m, 1)
}

func TestCollector_RemoveOldViews(t *testing.T) {
	var (
		c         = NewCollector()
		views     = []ViewKey{{N: 1, R: 0}, {N: 1, R: 8}, {N: 2, R: 0}, {N: 2, R: 1}, {N: 3, R: 0}}
		threshold = ViewKey{N: 2, R: 1}
		wantGone  = []bool{true, true, true, false, false}
	)

	for _, v := range views {
		c.AddPrepreparedTime(v, time.Now(), common.Hash{})
		c.AddCandMsg(v, common.HexToAddress("0x01"), time.Now(), &VRankCandidate{BlockNumber: v.N, Round: v.R})
	}

	c.RemoveOldViews(threshold)
	for i, v := range views {
		at, _, m := c.GetViewData(v)
		if wantGone[i] {
			assert.True(t, at.IsZero(), "view %v should be removed", v)
			assert.Nil(t, m)
		} else {
			assert.False(t, at.IsZero(), "view %v should remain", v)
			assert.NotNil(t, m)
		}
	}
}

func TestViewKey_Cmp(t *testing.T) {
	a, b, c := ViewKey{N: 1, R: 0}, ViewKey{N: 2, R: 0}, ViewKey{N: 1, R: 1}
	assert.Less(t, a.Cmp(b), 0)
	assert.Greater(t, b.Cmp(a), 0)
	assert.Equal(t, 0, a.Cmp(a))
	assert.Less(t, a.Cmp(c), 0)
	assert.Greater(t, c.Cmp(a), 0)
}
