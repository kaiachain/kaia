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

package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/engine"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyHeaderPrePermissionlessNoOp(t *testing.T) {
	ctrl := gomock.NewController(t)
	config := params.TestKaiaConfig("osaka")
	config.PermissionlessCompatibleBlock = big.NewInt(10)

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(config)

	v := NewValsetModule()
	v.Chain = mockChain

	// Missing Istanbul extra validators would fail Sealer().Validators(header),
	// so this confirms pre-HF headers bypass the permissionless check entirely.
	header := &types.Header{Number: big.NewInt(9)}
	assert.NoError(t, v.VerifyHeader(header, nil))
}

func TestVerifyHeaderPermissionlessValidatorOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	config := params.TestKaiaConfig("permissionless")
	sealer := engine.NewSealer(nil, nil)

	mockChain := chain_mock.NewMockBlockChain(ctrl)
	mockChain.EXPECT().Config().Return(config).AnyTimes()
	mockChain.EXPECT().Sealer().Return(sealer).AnyTimes()

	v := NewValsetModule()
	v.Chain = mockChain
	v.transitionResultCache.Add(uint64(1), &TransitionResult{
		Nodes: valset.NodeMap{
			numToAddr(1): &valset.Node{State: valset.ValActive},
			numToAddr(2): &valset.Node{State: valset.ValActive},
		},
	})

	valid := &types.Header{Number: big.NewInt(1)}
	require.NoError(t, sealer.WriteValidators(valid, numsToAddrs(1, 2)))
	assert.NoError(t, v.VerifyHeader(valid, nil))

	reordered := &types.Header{Number: big.NewInt(1)}
	require.NoError(t, sealer.WriteValidators(reordered, numsToAddrs(2, 1)))
	assert.ErrorIs(t, v.VerifyHeader(reordered, nil), errMismatchedValidators)
}
