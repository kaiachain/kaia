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

package cn

import (
	"errors"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type cnPeerReadResult struct {
	addrs []common.Address
	err   error
}

type fakeCNPeerReader struct {
	results []cnPeerReadResult
	calls   []uint64
}

func (r *fakeCNPeerReader) GetCNPeers(num uint64) ([]common.Address, error) {
	r.calls = append(r.calls, num)
	if len(r.results) == 0 {
		return nil, nil
	}
	result := r.results[0]
	if len(r.results) > 1 {
		r.results = r.results[1:]
	}
	return result.addrs, result.err
}

type recordCNPeerSink struct {
	calls [][]common.Address
}

func (s *recordCNPeerSink) SetCNPeers(addrs []common.Address) {
	if addrs == nil {
		s.calls = append(s.calls, nil)
		return
	}
	cp := make([]common.Address, len(addrs))
	copy(cp, addrs)
	s.calls = append(s.calls, cp)
}

func TestCNPeerUpdaterSyncHeadUsesNextBlockNumber(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{{addrs: cnTestAddrs(1)}}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.syncHead(types.NewBlockWithHeader(&types.Header{Number: big.NewInt(10)}))

	require.Equal(t, []uint64{11}, reader.calls)
	require.Len(t, sink.calls, 1)
	assert.Equal(t, cnTestAddrs(1), sink.calls[0])
}

func TestCNPeerUpdaterSkipsDuplicatePushes(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{{addrs: cnTestAddrs(1, 2)}}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.sync(11)
	updater.sync(12)

	require.Len(t, sink.calls, 1)
	assert.Equal(t, cnTestAddrs(1, 2), sink.calls[0])
}

func TestCNPeerUpdaterNilDisablesFiltering(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{{addrs: nil}}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.sync(11)
	updater.sync(12)

	require.Len(t, sink.calls, 1)
	assert.Nil(t, sink.calls[0])
}

func TestCNPeerUpdaterNilAndEmptyInputsAreDistinct(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{
		{addrs: nil},
		{addrs: []common.Address{}},
	}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.sync(11)
	updater.sync(12)

	require.Len(t, sink.calls, 2)
	assert.Nil(t, sink.calls[0])
	assert.NotNil(t, sink.calls[1])
	assert.Empty(t, sink.calls[1])
}

func TestCNPeerUpdaterFailureDisablesBeforeFirstPush(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{{err: errors.New("missing cn peers")}}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.sync(11)

	require.Len(t, sink.calls, 1)
	assert.Nil(t, sink.calls[0])
}

func TestCNPeerUpdaterFailureDisablesFilteringAfterPreviousPush(t *testing.T) {
	reader := &fakeCNPeerReader{results: []cnPeerReadResult{
		{addrs: cnTestAddrs(1)},
		{err: errors.New("temporary read failure")},
	}}
	sink := &recordCNPeerSink{}
	updater := newCNPeerUpdater(reader, sink)

	updater.sync(11)
	updater.sync(12)

	require.Len(t, sink.calls, 2)
	assert.Equal(t, cnTestAddrs(1), sink.calls[0])
	assert.Nil(t, sink.calls[1])
}

func cnTestAddrs(n ...int) []common.Address {
	addrs := make([]common.Address, len(n))
	for i, num := range n {
		addrs[i] = common.BigToAddress(big.NewInt(int64(num)))
	}
	return addrs
}
