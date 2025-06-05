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
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/consensus/istanbul"
	"github.com/kaiachain/kaia/v2/fork"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/stretchr/testify/assert"
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
			mockBackend, mockCtrl := newMockBackend(t, validatorAddrs)
			if tc.valid {
				mockBackend.EXPECT().Sign(gomock.Any()).Return(nil, nil).AnyTimes()
				mockBackend.EXPECT().Broadcast(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
			}

			istConfig := istanbul.DefaultConfig.Copy()
			istConfig.ProposerPolicy = istanbul.WeightedRandom

			istCore := New(mockBackend, istConfig).(*core)
			assert.NoError(t, istCore.Start())

			lastProposal, _ := mockBackend.LastProposal()
			proposal, err := genBlock(lastProposal.(*types.Block), validatorKeyMap[validatorAddrs[0]])
			assert.NoError(t, err)

			istCore.current.round.Set(big.NewInt(tc.round))
			istCore.current.Preprepare = &istanbul.Preprepare{
				View:     istCore.currentView(),
				Proposal: proposal,
			}
			istCore.sendCommit()
			istCore.Stop()
			mockCtrl.Finish()
		}
	}
}
