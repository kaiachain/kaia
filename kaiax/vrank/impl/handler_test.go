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
	"slices"
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

func mustPop(t *testing.T, sub chan *vrank.BroadcastRequest) *vrank.BroadcastRequest {
	select {
	case req := <-sub:
		return req
	case <-time.After(2 * time.Second):
		t.Fatal("should broadcast")
	}
	return nil
}

func mustNotPop(t *testing.T, sub chan *vrank.BroadcastRequest) *vrank.BroadcastRequest {
	select {
	case <-sub:
		t.Fatal("should not broadcast")
	default:
	}
	return nil
}

// VRankScenario defines a single-block scenario for the full VRank cycle (optimistic: all messages delivered).
// If UnresponsiveCands is set, HandleVRankCandidate is not called for those candidates (message dropped),
// so they appear in cfReport for validators.
type VRankScenario struct {
	Name              string
	Nodes             []string   // all node names, e.g. ["N1", "N2", "N3"]
	Council           []string   // Council(1)
	Candidates        []string   // Candidates(1)
	Proposer          string     // Proposer(1, 0)
	UnresponsiveCands []string   // optional: candidates whose message is not delivered to validators
	ExpectedCfReports [][]string // ExpectedCfReports[i] = expected GetCfReport(1, 0) for Nodes[i]
}

// runVRankScenario runs the full cycle for block 1 and asserts GetCfReport(1, 0) for each node matches ExpectedCfReports[i].
func runVRankScenario(t *testing.T, s VRankScenario) {
	const blockNum = uint64(1)
	require.Len(t, s.ExpectedCfReports, len(s.Nodes), "ExpectedCfReports must have one entry per node")

	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)

	nameToCN := make(map[string]*CN)
	for _, name := range s.Nodes {
		nameToCN[name] = createCN(t, valset)
	}

	councilAddrs := make([]common.Address, 0, len(s.Council))
	for _, name := range s.Council {
		councilAddrs = append(councilAddrs, nameToCN[name].Addr)
	}
	candAddrs := make([]common.Address, 0, len(s.Candidates))
	for _, name := range s.Candidates {
		candAddrs = append(candAddrs, nameToCN[name].Addr)
	}
	proposerAddr := nameToCN[s.Proposer].Addr
	valset.EXPECT().GetCouncil(blockNum).Return(councilAddrs, nil).AnyTimes()
	valset.EXPECT().GetCandidates(blockNum).Return(candAddrs, nil).AnyTimes()
	valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(proposerAddr, nil).AnyTimes()

	block1 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	view1_0 := &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	pppMsg := &vrank.VRankPreprepare{Block: block1, View: view1_0}

	// 1. HandleIstanbulPreprepare: every council member
	for _, name := range s.Council {
		nameToCN[name].VRankModule.HandleIstanbulPreprepare(block1, view1_0)
	}

	// 2. HandleVRankPreprepare: each candidate receives preprepare (optimistic: we deliver directly)
	// 3. HandleVRankCandidate: each validator receives each candidate's message (unless in UnresponsiveCands)
	notDelivered := make(map[string]bool)
	for _, name := range s.UnresponsiveCands {
		notDelivered[name] = true
	}
	for _, candName := range s.Candidates {
		cand := nameToCN[candName]
		_ = cand.VRankModule.HandleVRankPreprepare(pppMsg)
		if notDelivered[candName] {
			continue
		}
		sig, _ := crypto.Sign(crypto.Keccak256(block1.Hash().Bytes()), cand.Key)
		candMsg := &vrank.VRankCandidate{
			BlockNumber: blockNum,
			Round:       uint8(view1_0.Round.Uint64()),
			BlockHash:   block1.Hash(),
			Sig:         sig,
		}
		for _, valName := range s.Council {
			err := nameToCN[valName].VRankModule.HandleVRankCandidate(candMsg)
			require.NoError(t, err)
		}
	}

	// 4. GetCfReport(1, 0) from each node and assert ExpectedCfReports[i]
	for i, nodeName := range s.Nodes {
		report, err := nameToCN[nodeName].VRankModule.GetCfReport(blockNum, 0)
		require.NoError(t, err)
		expectedNames := s.ExpectedCfReports[i]
		expectedAddrs := make([]common.Address, 0, len(expectedNames))
		for _, name := range expectedNames {
			expectedAddrs = append(expectedAddrs, nameToCN[name].Addr)
		}
		if len(expectedAddrs) == 0 {
			assert.Empty(t, report, "node %s: expected empty cfReport", nodeName)
		} else {
			require.Len(t, report, len(expectedAddrs), "node %s: cfReport length", nodeName)
			for _, addr := range expectedAddrs {
				assert.True(t, slices.Contains(report, addr), "node %s: cfReport should contain %s", nodeName, addr.Hex())
			}
		}
	}
}

func TestVRankModule(t *testing.T) {
	scenarios := []VRankScenario{
		{
			Name:       "happy path and non-council get empty report",
			Nodes:      []string{"N1", "N2", "N3"},
			Council:    []string{"N1"},
			Candidates: []string{"N2"},
			Proposer:   "N1",
			ExpectedCfReports: [][]string{
				{}, // N1: Council(1), valid CfReport
				{}, // N2: not in Council(1)
				{}, // N3: not in Council(1)
			},
		},
		{
			Name:              "Unresponsive N3 should be in cfReport",
			Nodes:             []string{"N1", "N2", "N3"},
			Council:           []string{"N1", "N2"},
			Candidates:        []string{"N3"},
			Proposer:          "N1",
			UnresponsiveCands: []string{"N3"},
			ExpectedCfReports: [][]string{
				{"N3"}, // N1: Council(1), cfReport contains N3
				{"N3"}, // N2: Council(1), cfReport contains N3
				{},     // N3: not in Council(1)
			},
		},
		{
			Name:              "Unresponsive N4 should be in cfReport",
			Nodes:             []string{"N1", "N2", "N3", "N4"},
			Council:           []string{"N1", "N2"},
			Candidates:        []string{"N3", "N4"},
			Proposer:          "N1",
			UnresponsiveCands: []string{"N4"},
			ExpectedCfReports: [][]string{
				{"N4"}, // N1: Council(1)
				{"N4"}, // N2: Council(1)
				{},     // N3: not in Council(1)
				{},     // N4: not in Council(1)
			},
		},
		{
			Name:              "Unresponsive N3 and N4 should be in cfReport",
			Nodes:             []string{"N1", "N2", "N3", "N4"},
			Council:           []string{"N1", "N2"},
			Candidates:        []string{"N3", "N4"},
			Proposer:          "N1",
			UnresponsiveCands: []string{"N3", "N4"},
			ExpectedCfReports: [][]string{
				{"N3", "N4"}, // N1: Council(1)
				{"N3", "N4"}, // N2: Council(1)
				{},           // N3: not in Council(1)
				{},           // N4: not in Council(1)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.Name, func(t *testing.T) {
			runVRankScenario(t, s)
		})
	}
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
		prepreparedTime, _, _ := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		mustNotPop(t, val.sub)
	})

	t.Run("the proposer should not start collection when not in the next council", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		proposer, validator, candidate := createCN(t, valset), createCN(t, valset), createCN(t, valset)

		// proposer is not in the next council, so it should only broadcast and does not start collection.
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{validator.Addr}, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(2)
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(2)

		proposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		prepreparedTime, _, _ := proposer.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		mustPop(t, proposer.sub) // proposer should broadcast

		validator.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		prepreparedTime, _, _ = validator.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.False(t, prepreparedTime.IsZero())
		mustNotPop(t, validator.sub)
	})

	t.Run("non-proposers including candidate should not broadcast", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		proposer, nonProposer, candidate := createCN(t, valset), createCN(t, valset), createCN(t, valset)

		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(3)
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(3)

		proposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		nonProposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		candidate.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		prepreparedTime, _, _ := proposer.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.False(t, prepreparedTime.IsZero())
		prepreparedTime, _, _ = nonProposer.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.False(t, prepreparedTime.IsZero())
		prepreparedTime, _, _ = candidate.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())

		req := mustPop(t, proposer.sub)
		assert.Equal(t, []common.Address{candidate.Addr}, req.Targets)

		mustNotPop(t, nonProposer.sub)
		mustNotPop(t, candidate.sub)
	})
}

func TestHandleVRankPreprepare(t *testing.T) {
	var (
		block1  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view1_0 = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	)

	t.Run("permissionless fork is disabled", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		cand := createCN(t, valset)
		cand.VRankModule.ChainConfig.PermissionlessCompatibleBlock = nil
		cand.VRankModule.HandleVRankPreprepare(&vrank.VRankPreprepare{Block: block1, View: view1_0})
		mustNotPop(t, cand.sub)
	})

	t.Run("validators should not broadcast", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		proposer, nonProposer, candidate := createCN(t, valset), createCN(t, valset), createCN(t, valset)

		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(3)
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(3)

		proposer.VRankModule.HandleVRankPreprepare(&vrank.VRankPreprepare{Block: block1, View: view1_0})
		nonProposer.VRankModule.HandleVRankPreprepare(&vrank.VRankPreprepare{Block: block1, View: view1_0})
		candidate.VRankModule.HandleVRankPreprepare(&vrank.VRankPreprepare{Block: block1, View: view1_0})

		req := mustPop(t, candidate.sub)
		assert.Equal(t, []common.Address{proposer.Addr, nonProposer.Addr}, req.Targets)

		mustNotPop(t, nonProposer.sub)
		mustNotPop(t, proposer.sub)
	})
}

func TestHandleVRankCandidate(t *testing.T) {
	var (
		block1  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view1_0 = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	)

	t.Run("permissionless fork is disabled", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		val.VRankModule.ChainConfig.PermissionlessCompatibleBlock = nil
		val.VRankModule.HandleVRankCandidate(&vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: []byte{}})
		mustNotPop(t, val.sub)
	})

	t.Run("no nodes should broadcast", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		proposer, nonProposer, candidate := createCN(t, valset), createCN(t, valset), createCN(t, valset)
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: []byte{}}

		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(3)
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(3)

		proposer.VRankModule.HandleVRankCandidate(&msg)
		nonProposer.VRankModule.HandleVRankCandidate(&msg)
		candidate.VRankModule.HandleVRankCandidate(&msg)

		mustNotPop(t, proposer.sub)
		mustNotPop(t, nonProposer.sub)
		mustNotPop(t, candidate.sub)
	})

	t.Run("the proposer should not collect when not in the next council", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		proposer, validator, candidate := createCN(t, valset), createCN(t, valset), createCN(t, valset)
		sig, _ := crypto.Sign(crypto.Keccak256(block1.Hash().Bytes()), candidate.Key)
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: sig}

		// proposer is not in the next council, so it should only broadcast and does not start collection.
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{validator.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(2)
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(2)

		proposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0) // this won't happen in production
		proposer.VRankModule.HandleVRankCandidate(&msg)
		prepreparedTime, _, candMap := proposer.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		assert.Nil(t, candMap)

		validator.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		validator.VRankModule.HandleVRankCandidate(&msg)
		prepreparedTime, _, candMap = validator.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.False(t, prepreparedTime.IsZero())
		assert.Equal(t, 1, len(candMap))
	})

	t.Run("future messages", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val, cand := createCN(t, valset), createCN(t, valset)

		sig, _ := crypto.Sign(crypto.Keccak256(block1.Hash().Bytes()), cand.Key)
		block2 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		invalidSig, _ := crypto.Sign(crypto.Keccak256(block2.Hash().Bytes()), cand.Key)

		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{cand.Addr}, nil).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		tcs := []struct {
			name    string
			msg     *vrank.VRankCandidate
			wantErr error
		}{
			{
				name: "future block number",
				msg:  &vrank.VRankCandidate{BlockNumber: 2, Round: 0, BlockHash: block1.Hash(), Sig: sig}, wantErr: nil,
			},
			{
				name: "future round",
				msg:  &vrank.VRankCandidate{BlockNumber: 1, Round: 1, BlockHash: block1.Hash(), Sig: sig}, wantErr: nil,
			},
			{
				name: "future block hash",
				msg:  &vrank.VRankCandidate{BlockNumber: 1, Round: 0, BlockHash: block2.Hash(), Sig: sig}, wantErr: vrank.ErrMsgFromNonCandidate,
			},
			{
				name: "future signature and message mismatch",
				msg:  &vrank.VRankCandidate{BlockNumber: 1, Round: 0, BlockHash: block1.Hash(), Sig: invalidSig}, wantErr: vrank.ErrMsgFromNonCandidate,
			},
			{
				name: "invalid signature (signature of block 2)",
				msg:  &vrank.VRankCandidate{BlockNumber: 2, Round: 0, BlockHash: block1.Hash(), Sig: invalidSig}, wantErr: vrank.ErrMsgFromNonCandidate,
			},
		}

		for _, tc := range tcs {
			t.Run(tc.name, func(t *testing.T) {
				err := val.VRankModule.HandleVRankCandidate(tc.msg)
				assert.Equal(t, tc.wantErr, err)
			})
		}
	})

	t.Run("duplicate message", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val, cand := createCN(t, valset), createCN(t, valset)

		sig, _ := crypto.Sign(crypto.Keccak256(block1.Hash().Bytes()), cand.Key)
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: sig}

		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{cand.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		var receivedAt time.Time
		for range 3 {
			err := val.VRankModule.HandleVRankCandidate(&msg)
			assert.NoError(t, err)
			prepreparedTime, _, candMap := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
			assert.False(t, prepreparedTime.IsZero())
			assert.Equal(t, 1, len(candMap))
			cm := candMap[cand.Addr]
			assert.Greater(t, cm.ReceivedAt.Sub(prepreparedTime), time.Duration(0))
			if receivedAt.IsZero() {
				receivedAt = cm.ReceivedAt
			} else {
				assert.Equal(t, receivedAt, cm.ReceivedAt, "ReceivedAt should not change on duplicate")
			}
		}
	})
}

func TestGetCfReport(t *testing.T) {
	var (
		valset                 = mock_valset.NewMockValsetModule(gomock.NewController(t))
		block1                 = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		block2                 = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		view1_0                = &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
		view2_0                = &istanbul.View{Sequence: big.NewInt(2), Round: common.Big0}
		validators, candidates []*CN
		valAddrs               = make([]common.Address, 3)
		candAddrs              = make([]common.Address, 8)
		candMsgsBlock2         = make([]vrank.VRankCandidate, 8)
		earlybirdCands         = candAddrs[0:2] // sent VRankCandidate for block2 before validator preprepared block2
		ontimeCands            = candAddrs[2:4]
		liarCands              = candAddrs[4:6]
		lateCands              = candAddrs[6:8]
	)

	for i := 0; i < 3; i++ {
		validators = append(validators, createCN(t, valset))
		valAddrs[i] = validators[i].Addr
	}
	for i := 0; i < 8; i++ {
		candidates = append(candidates, createCN(t, valset))
		candAddrs[i] = candidates[i].Addr
	}

	valset.EXPECT().GetCouncil(gomock.Any()).Return(valAddrs, nil).AnyTimes()
	valset.EXPECT().GetCandidates(gomock.Any()).Return(candAddrs, nil).AnyTimes()
	valset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(validators[0].Addr, nil).AnyTimes()

	for i := 0; i < 8; i++ {
		sig, _ := crypto.Sign(crypto.Keccak256(block2.Hash().Bytes()), candidates[i].Key)
		candMsgsBlock2[i] = vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: block2.Hash(), Sig: sig}
	}

	// Initialize prepreparedView at least once to set `v.prepreparedView`.
	for _, v := range validators {
		v.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
	}

	// Earlybirds: candidates send VRankCandidate for block2 before validator has preprepared block2.
	for i := 0; i < 2; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	// Now validator preprepares for block2.
	for _, v := range validators {
		v.VRankModule.HandleIstanbulPreprepare(block2, view2_0)
	}

	// On-time: candidates send VRankCandidate for block2 after validator has preprepared block2 and before deadline.
	for i := 2; i < 4; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	// Liars: candidates send VRankCandidate for block2 with wrong BlockHash.
	for i := 4; i < 6; i++ {
		liarHash := common.Hash{byte(i)}
		sig, _ := crypto.Sign(crypto.Keccak256(liarHash.Bytes()), candidates[i].Key)
		liarMsg := vrank.VRankCandidate{BlockNumber: block2.NumberU64(), Round: uint8(view2_0.Round.Uint64()), BlockHash: liarHash, Sig: sig}
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&liarMsg)
			assert.NoError(t, err)
		}
	}
	time.Sleep(candidatePrepareDeadlineMs * time.Millisecond)
	// Late: candidates send VRankCandidate for block2 after deadline.
	for i := 6; i < 8; i++ {
		for _, v := range validators {
			err := v.VRankModule.HandleVRankCandidate(&candMsgsBlock2[i])
			assert.NoError(t, err)
		}
	}

	for _, v := range validators {
		report, err := v.VRankModule.GetCfReport(2, 0)
		assert.NoError(t, err)
		assert.Len(t, report, 4, "cfReport: 2 liars + 2 late")
		for _, addr := range ontimeCands {
			assert.False(t, slices.Contains(report, addr))
		}
		for _, addr := range earlybirdCands {
			assert.False(t, slices.Contains(report, addr))
		}
		for _, addr := range liarCands {
			assert.True(t, slices.Contains(report, addr))
		}
		for _, addr := range lateCands {
			assert.True(t, slices.Contains(report, addr))
		}
		report2, err := v.VRankModule.GetCfReport(2, 0)
		assert.NoError(t, err)
		assert.Equal(t, report, report2, "GetCfReport must be deterministic")
	}
}

func TestGetCfReport_Errors(t *testing.T) {
	block1 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	view1_0 := &istanbul.View{Sequence: big.NewInt(1), Round: common.Big0}
	candAddr := common.HexToAddress("0xc4nd1d473")

	t.Run("CfReport should contain candAddr", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.GetCfReport(1, 0)
		require.NoError(t, err)
		assert.True(t, slices.Contains(report, candAddr))
	})

	t.Run("epoch header returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(vrankEpoch-1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(vrankEpoch-1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(vrankEpoch-1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(int64(vrankEpoch - 1))})
		view := &istanbul.View{Sequence: big.NewInt(vrankEpoch - 1), Round: common.Big0}
		val.VRankModule.HandleIstanbulPreprepare(block, view)

		report, err := val.VRankModule.GetCfReport(vrankEpoch-1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("round out of range returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.GetCfReport(1, 11) // maxRound is 10
		require.ErrorIs(t, err, vrank.ErrRoundOutOfRange)
		assert.Nil(t, report)

		report, err = val.VRankModule.GetCfReport(1, 10)
		assert.NotErrorIs(t, err, vrank.ErrRoundOutOfRange)
	})

	t.Run("non-validator returns empty report", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val, otherVal := createCN(t, valset), createCN(t, valset)
		// This node is not in the council for block 1.
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{otherVal.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.GetCfReport(1, 0)
		require.NoError(t, err)
		assert.Empty(t, report)
	})

	t.Run("preprepared time not set returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		// skip HandleIstanbulPreprepare

		prepreparedTime, _, _ := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		report, err := val.VRankModule.GetCfReport(1, 0)
		require.ErrorIs(t, err, vrank.ErrPrepreparedTimeNotSet)
		assert.Nil(t, report)
	})

	t.Run("GetCandidates failed returns error", func(t *testing.T) {
		valset := mock_valset.NewMockValsetModule(gomock.NewController(t))
		val := createCN(t, valset)
		valset.EXPECT().GetCouncil(uint64(1)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandidates(uint64(1)).Return(nil, assert.AnError).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		report, err := val.VRankModule.GetCfReport(1, 0)
		require.ErrorIs(t, err, vrank.ErrGetCandidateFailed)
		assert.Nil(t, report)
	})
}
