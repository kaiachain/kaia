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

package testing

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/vrank"
)

// StubVRankModule is a no-op VRankModule for tests that need a non-nil VRankModule
// but don't require actual VRank scoring behavior.
type StubVRankModule struct {
	PFS map[common.Address]uint64
}

func (*StubVRankModule) Start() error { return nil }
func (*StubVRankModule) Stop()        {}

func (m *StubVRankModule) GetPFS(uint64) (map[common.Address]uint64, error) { return m.PFS, nil }
func (*StubVRankModule) GetCFS(uint64) (map[common.Address]uint64, error)   { return nil, nil }

func (*StubVRankModule) HandleIstanbulPreprepare(*types.Block, *istanbul.View) {}
func (*StubVRankModule) HandleVRankPreprepare(*vrank.VRankPreprepare) error     { return nil }
func (*StubVRankModule) HandleVRankCandidate(*vrank.VRankCandidate) error       { return nil }
func (*StubVRankModule) TallyCfReport(uint64, uint64) ([]common.Address, error) {
	return nil, nil
}
func (*StubVRankModule) SubscribeVRank(chan<- *vrank.VRankBroadcastEvent) event.Subscription {
	return nil
}
func (*StubVRankModule) VerifyHeader(*types.Header) error  { return nil }
func (*StubVRankModule) PrepareHeader(*types.Header) error { return nil }
func (*StubVRankModule) FinalizeHeader(*types.Header, *state.StateDB, []*types.Transaction, []*types.Receipt) error {
	return nil
}