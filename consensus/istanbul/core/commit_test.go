// Copyright 2024 The Kaia Authors
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
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCore_sendCommit(t *testing.T) {
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{})
	defer fork.ClearHardForkBlockNumberConfig()

	validatorAddrs, validatorKeyMap := genValidators(6)

	for _, tc := range []struct {
		tcName string
		round  int64
		valid  bool
	}{
		{"valid case", 0, true},
		{"invalid case - not committee", 2, false},
	} {
		{
			mockBackend, mockCtrl, mockValset, mockGov := newMockBackend(t, validatorAddrs)
			if tc.valid {
				mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).AnyTimes()
				mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			}

			istConfig := istanbul.DefaultConfig.Copy()
			istConfig.ProposerPolicy = istanbul.WeightedRandom

			istCore := New(mockBackend, istConfig).(*core)
			istCore.RegisterKaiaxModules(mockValset, mockGov)
			assert.NoError(t, istCore.Start())

			lastProposal, _ := mockBackend.LastProposal()
			proposal, err := genBlock(lastProposal.(*types.Block), validatorKeyMap[validatorAddrs[0]])
			assert.NoError(t, err)

			istCore.current.round.Set(big.NewInt(tc.round))
			istCore.current.Preprepare = &bft.Preprepare{
				View:     istCore.currentView(),
				Proposal: proposal,
			}
			istCore.sendCommit()
			istCore.Stop()
			mockCtrl.Finish()
		}
	}
}

// TestSendCommitForOldBlockSealMatchesDigest is a regression test for the COMMIT
// committed-seal desync: finalizeMessage must sign the CommittedSeal over the
// digest carried in the COMMIT message's own Subject (the block it votes for),
// not over the node's current proposal. sendCommitForOldBlock rebroadcasts a
// COMMIT whose Subject.Digest is an old block hash, so the produced seal must
// attest to that old block hash.
func TestSendCommitForOldBlockSealMatchesDigest(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	sealKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	sealAddr := crypto.PubkeyToAddress(sealKey.PublicKey)

	// Capture the finalized COMMIT payload broadcast to peers.
	var payload []byte
	mockBackend := mock_istanbul.NewMockBackend(mockCtrl)
	mockBackend.EXPECT().Address().Return(sealAddr).AnyTimes()
	mockBackend.EXPECT().Sealer().Return(istanbul.NewSealerImpl(sealKey)).AnyTimes()
	mockBackend.EXPECT().IsPermissionlessAt(gomock.Any()).Return(false).AnyTimes()
	mockBackend.EXPECT().Sign(gomock.Any()).Return([]byte{0x01}, nil).AnyTimes() // message-payload signature, unused here
	mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).DoAndReturn(func(_ common.Hash, p []byte) error {
		payload = append([]byte(nil), p...)
		return nil
	}).AnyTimes()

	istCore := New(mockBackend, istanbul.DefaultConfig.Copy()).(*core)

	oldBlockHash := common.HexToHash("0xa11d")        // digest the COMMIT votes for
	currentProposalHash := common.HexToHash("0xb22d") // unrelated current proposal
	view := &bft.View{Round: big.NewInt(0), Sequence: big.NewInt(1)}

	istCore.sendCommitForOldBlock(view, oldBlockHash, common.HexToHash("0xc33d"))

	require.NotNil(t, payload, "a COMMIT must have been broadcast")
	var msg bft.Message
	require.NoError(t, rlp.DecodeBytes(payload, &msg))

	// The committed seal must recover to the sealer over the message's own
	// (old block) digest...
	got, err := istanbul.GetSignatureAddress(istanbul.PrepareCommittedSeal(oldBlockHash), msg.CommittedSeal)
	require.NoError(t, err)
	assert.Equal(t, sealAddr, got, "CommittedSeal must attest to the message's own (old block) digest")

	// ...and must NOT recover to the sealer over an unrelated current proposal digest.
	wrong, _ := istanbul.GetSignatureAddress(istanbul.PrepareCommittedSeal(currentProposalHash), msg.CommittedSeal)
	assert.NotEqual(t, sealAddr, wrong, "CommittedSeal must not attest to an unrelated (current proposal) digest")
}
