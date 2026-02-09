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

func ensurePop(t *testing.T, sub chan *vrank.BroadcastRequest) *vrank.BroadcastRequest {
	select {
	case req := <-sub:
		return req
	case <-time.After(2 * time.Second):
		t.Fatal("should broadcast")
	}
	return nil
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
	req := ensurePop(t, val.sub)
	assert.Equal(t, []common.Address{cand.Addr}, req.Targets)
	assert.Equal(t, VRankPreprepareMsg, req.Code)
	assert.Equal(t, &pppMsg, req.Msg)

	// candidate correctly broadcast VRankCandidate upon receiving VRankPreprepare
	cand.VRankModule.HandleVRankPreprepare(&pppMsg)
	req = ensurePop(t, cand.sub)
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

func TestHandleIstanbulPreprepare(t *testing.T) {
	var (
		block1  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view1_0 = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	)

	t.Run("permissionless fork is disabled", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		val.VRankModule.ChainConfig.PermissionlessCompatibleBlock = nil
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		assert.Equal(t, time.Time{}, val.VRankModule.prepreparedTime)
		select {
		case <-val.sub:
			t.Fatal("under disabled permissionless fork, it should not broadcast")
		default:
		}
	})

	t.Run("non-proposers should not broadcast", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		var vals []*CN
		for i := 0; i < 3; i++ {
			vals = append(vals, createCN(t, valset))
		}
		proposer := vals[0]
		nonProposer := vals[1]
		candidate := vals[2]

		valset.EXPECT().GetCouncil(uint64(2)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(len(vals))
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(len(vals))
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(len(vals))

		proposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		nonProposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		candidate.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		assert.NotEqual(t, time.Time{}, proposer.VRankModule.prepreparedTime)
		assert.NotEqual(t, time.Time{}, nonProposer.VRankModule.prepreparedTime)
		assert.Equal(t, time.Time{}, candidate.VRankModule.prepreparedTime) // candidate should not preprepare

		req := ensurePop(t, proposer.sub)
		assert.Equal(t, []common.Address{candidate.Addr}, req.Targets)

		select {
		case <-nonProposer.sub:
			t.Fatal("non-proposer should not broadcast")
		case <-candidate.sub:
			t.Fatal("candidate should not broadcast")
		default:
		}
	})
}
