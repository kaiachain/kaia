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

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/istanbul"
	mock_istanbul "github.com/kaiachain/kaia/consensus/istanbul/mocks"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestHandlePreprepareOldBlockCommitRequiresStoredRound covers the COMMIT rebroadcast for an
// already-committed block: post-permissionless the answer is bound to the round this node stored.
func TestHandlePreprepareOldBlockCommitRequiresStoredRound(t *testing.T) {
	const blockNum = uint64(10)

	testcases := []struct {
		name            string
		permissionless  bool
		storedRound     byte
		requestedRound  int64
		expectBroadcast bool
	}{
		{"permissionless answers the stored round", true, 1, 1, true},
		{"permissionless refuses another round", true, 1, 5, false},
		{"before fork answers any round", false, 1, 5, true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			t.Cleanup(mockCtrl.Finish)

			key, err := crypto.GenerateKey()
			require.NoError(t, err)
			src := crypto.PubkeyToAddress(key.PublicKey)

			proposal := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(blockNum)})
			view := &bft.View{Round: big.NewInt(tc.requestedRound), Sequence: new(big.Int).SetUint64(blockNum)}

			broadcast := false
			mockBackend := mock_istanbul.NewMockBackend(mockCtrl)
			mockBackend.EXPECT().Address().Return(src).AnyTimes()
			mockBackend.EXPECT().Sealer().Return(istanbul.NewSealerImpl(key)).AnyTimes()
			mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).AnyTimes()
			mockBackend.EXPECT().IsPermissionlessAt(blockNum).Return(tc.permissionless).AnyTimes()
			mockBackend.EXPECT().ProposalRound(proposal.Hash(), proposal.Number()).Return(tc.storedRound, true)
			mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).DoAndReturn(func(common.Hash, []byte) error {
				broadcast = true
				return nil
			}).AnyTimes()

			mockValset := valset_mock.NewMockValsetModule(mockCtrl)
			mockValset.EXPECT().GetProposer(blockNum, uint64(tc.requestedRound)).Return(src, nil)

			istCore := New(mockBackend, istanbul.DefaultConfig.Copy()).(*core)
			istCore.RegisterKaiaxModules(mockValset, nil)
			// Already past blockNum, so the PRE-PREPARE arrives as an old message.
			istCore.current = newRoundState(
				&bft.View{Round: big.NewInt(0), Sequence: new(big.Int).SetUint64(blockNum + 1)},
				valset.NewAddressSet([]common.Address{src}), common.Hash{}, nil, nil, nil)

			payload, err := bft.Encode(&bft.Preprepare{View: view, Proposal: proposal})
			require.NoError(t, err)
			err = istCore.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: payload}, src)

			assert.Equal(t, tc.expectBroadcast, broadcast)
			if tc.expectBroadcast {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, errOldMessage)
			}
		})
	}
}
