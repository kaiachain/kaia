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
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #984: a PRE-PREPARE whose proposal number differs from its view sequence
// is invalid.
func TestDevParity_PreprepareNumberMustMatchView(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	m := newTestMachine(b, 10, 0, nil, []common.Address{self})
	m.state = stateAcceptRequest
	m.proposer = self

	pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(0)}, Proposal: newEmptyBlock(9, common.Hash{1})}
	err := m.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp)}, self)
	assert.ErrorIs(t, err, bft.ErrInvalidMessage)
}

// dev #924: ROUND CHANGE admission is limited to the committee the thresholds
// are computed over.
func TestDevParity_RoundChangeRequiresCommittee(t *testing.T) {
	b, key := newTestBackend(t, &fakeChain{cfg: legacyConfig()})
	self := crypto.PubkeyToAddress(key.PublicKey)
	outsider := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	m := newTestMachine(b, 10, 0, nil, []common.Address{self})
	m.qualified = valset.NewAddressSet([]common.Address{self, outsider})

	rc := &bft.Subject{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(1)}}
	err := m.handleRoundChange(&bft.Message{Code: bft.MsgRoundChange, Msg: mustEncode(t, rc), Address: outsider}, outsider)
	assert.ErrorIs(t, err, errNotFromCommittee)
	assert.Nil(t, m.maxRoundChangeRound(1), "a rejected ROUND CHANGE must not be retained")
}

// dev #905: the quorum is computed over the committee actually selected, not
// the governance CommitteeSize parameter.
func TestDevParity_QuorumUsesSelectedCommittee(t *testing.T) {
	ctrl := gomock.NewController(t)
	b, _ := newTestBackend(t, &fakeChain{cfg: permissionlessConfig()})
	vals := make([]common.Address, 7)
	for i := range vals {
		vals[i] = common.BigToAddress(big.NewInt(int64(i + 1)))
	}
	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCouncil(uint64(10)).Return(vals, nil)
	mValset.EXPECT().GetDemotedValidators(uint64(10)).Return(nil, nil)
	mValset.EXPECT().GetCommittee(uint64(10), uint64(0)).Return(vals, nil)
	mValset.EXPECT().GetProposer(uint64(10), uint64(0)).Return(vals[0], nil)
	mGov := gov_mock.NewMockGovModule(ctrl)
	mGov.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{CommitteeSize: 2}).AnyTimes()
	b.RegisterKaiaxModules(mGov, mValset)

	m := newMachine(b)
	_, committee, _, committeeSize, required, f, err := m.getRoundCommitteeState(10, 0)
	require.NoError(t, err)
	assert.Equal(t, 7, committee.Len())
	assert.Equal(t, uint64(7), committeeSize)
	assert.Equal(t, 5, required) // ceil(2*7/3)
	assert.Equal(t, 2, f)        // ceil(7/3)-1
}
