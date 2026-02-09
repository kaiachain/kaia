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
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CN struct {
	Key         *ecdsa.PrivateKey
	Addr        common.Address
	VRankModule *VRankModule
	sub         chan *vrank.BroadcastRequest
}

func createCN(t *testing.T, valset valset.ValsetModule) *CN {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	sub := make(chan *vrank.BroadcastRequest)
	module := NewVRankModule()
	err := module.Init(&InitOpts{
		NodeKey:     key,
		Valset:      valset,
		ChainConfig: params.TestKaiaConfig("permissionless"),
	})
	require.NoError(t, err)

	module.broadcastFeed.Subscribe(sub)
	err = module.Start()
	require.NoError(t, err)
	return &CN{
		Key:         key,
		Addr:        addr,
		VRankModule: module,
		sub:         sub,
	}
}

func TestVRankModule(t *testing.T) {
	var (
		valset = mock_valset.NewMockValsetModule(gomock.NewController(t))
		val    = createCN(t, valset)
		cand   = createCN(t, valset)

		block  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view   = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
		sig, _ = crypto.Sign(crypto.Keccak256(block.Hash().Bytes()), cand.Key)

		pppMsg  = vrank.VRankPreprepare{Block: block, View: view}
		candMsg = vrank.VRankCandidate{BlockNumber: block.NumberU64(), Round: uint8(view.Round.Uint64()), BlockHash: block.Hash(), Sig: sig}
	)

	t.Logf("val.Addr: %s, cand.Addr: %s", val.Addr.Hex(), cand.Addr.Hex())

	valset.EXPECT().GetCouncil(gomock.Any()).Return([]common.Address{val.Addr}, nil).AnyTimes()
	valset.EXPECT().GetCandidates(gomock.Any()).Return([]common.Address{cand.Addr}, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(val.Addr, nil).AnyTimes()

	// validator correctly broadcast VRankPreprepare upon receiving IstanbulPreprepare
	assert.Equal(t, time.Time{}, val.VRankModule.prepreparedTime)
	val.VRankModule.HandleIstanbulPreprepare(block, view)
	assert.NotEqual(t, time.Time{}, val.VRankModule.prepreparedTime)
	req := <-val.sub
	assert.Equal(t, []common.Address{cand.Addr}, req.Targets)
	assert.Equal(t, VRankPreprepareMsg, req.Code)
	assert.Equal(t, &pppMsg, req.Msg)

	// candidate correctly broadcast VRankCandidate upon receiving VRankPreprepare
	cand.VRankModule.HandleVRankPreprepare(&pppMsg)
	req = <-cand.sub
	assert.Equal(t, []common.Address{val.Addr}, req.Targets)
	assert.Equal(t, VRankCandidateMsg, req.Code)
	assert.Equal(t, &candMsg, req.Msg)

	// validator correctly collects VRankCandidate upon receiving VRankCandidate
	err := val.VRankModule.HandleVRankCandidate(&candMsg)
	assert.NoError(t, err)
	d, ok := val.VRankModule.candResponses.Load(cand.Addr)
	assert.True(t, ok)
	assert.NotEqual(t, time.Duration(0), d.(time.Duration))
}
