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
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	blstypes "github.com/kaiachain/kaia/crypto/bls/types"
	mock_randao "github.com/kaiachain/kaia/kaiax/randao/mock"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/kaiax/vrank"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func signVRankCandidate(t *testing.T, m *VRankModule, key *ecdsa.PrivateKey, blsKey bls.SecretKey, blockNum uint64, round uint8, blockHash common.Hash) ([crypto.SignatureLength]byte, [blstypes.SignatureLength]byte) {
	t.Helper()
	sigHash := m.vrankCandidateSigHash(blockNum, round, blockHash)
	sig, err := crypto.Sign(sigHash.Bytes(), key)
	require.NoError(t, err)
	blsSig := bls.Sign(blsKey, sigHash.Bytes()).Marshal()
	return [crypto.SignatureLength]byte(sig), [blstypes.SignatureLength]byte(blsSig)
}

func signVRankPreprepare(t *testing.T, m *VRankModule, key *ecdsa.PrivateKey, blockNum uint64, round uint8, blockHash common.Hash) [crypto.SignatureLength]byte {
	sigHash := m.vrankPreprepareSigHash(blockNum, round, blockHash)
	sig, err := crypto.Sign(sigHash.Bytes(), key)
	require.NoError(t, err)
	return [crypto.SignatureLength]byte(sig)
}

// newCNMulti builds n CNs sharing one fresh valset/randao mock. Each subtest gets its own
// pair (no cross-subtest leakage). Always installs withGenesis so HandleVRankCandidate's
// CurrentHeader() call has a non-nil head; pass extra options to override.
func newCNMulti(t *testing.T, n int, options ...vrankOpt) ([]*CN, *mock_valset.MockValsetModule, *mock_randao.MockRandaoModule) {
	ctrl := gomock.NewController(t)
	valset := mock_valset.NewMockValsetModule(ctrl)
	randao := mock_randao.NewMockRandaoModule(ctrl)
	ret := make([]*CN, n)
	for i := range n {
		opts := append([]vrankOpt{withValset(valset), withRandao(randao), withGenesis()}, options...)
		ret[i] = newCN(t, opts...)
	}
	return ret, valset, randao
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
	ExpectedCfReports [][]string // ExpectedCfReports[i] = expected EvaluateCandidates(1, 0) for Nodes[i]
}

// runVRankScenario runs the full cycle for block 1 and asserts EvaluateCandidates(1, 0) for each node matches ExpectedCfReports[i].
func runVRankScenario(t *testing.T, s VRankScenario) {
	const blockNum = uint64(1)
	require.Len(t, s.ExpectedCfReports, len(s.Nodes), "ExpectedCfReports must have one entry per node")

	// Anchor shared mocks on the first CN's bundle, then have every other CN re-use them so
	// all nodes observe the same network view.
	nameToCN := make(map[string]*CN)
	var valset *mock_valset.MockValsetModule
	var randao *mock_randao.MockRandaoModule
	for _, name := range s.Nodes {
		if valset == nil {
			cn := newCN(t, withGenesis())
			valset, randao = cn.Valset, cn.Randao
			nameToCN[name] = cn
			continue
		}
		nameToCN[name] = newCN(t, withValset(valset), withRandao(randao), withGenesis())
	}

	addrsOf := func(names []string) []common.Address {
		addrs := make([]common.Address, 0, len(names))
		for _, name := range names {
			addrs = append(addrs, nameToCN[name].Addr)
		}
		return addrs
	}
	valset.EXPECT().GetCommittee(blockNum, uint64(0)).Return(addrsOf(s.Council), nil).AnyTimes()
	valset.EXPECT().GetCandTesting(blockNum).Return(addrsOf(s.Candidates), nil).AnyTimes()
	valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(nameToCN[s.Proposer].Addr, nil).AnyTimes()

	block1 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	view1_0 := &bft.View{Sequence: big.NewInt(1), Round: common.Big0}
	proposerCN := nameToCN[s.Proposer]
	pppSig := signVRankPreprepare(t, proposerCN.VRankModule, proposerCN.Key, blockNum, 0, block1.Hash())
	pppMsg := &vrank.VRankPreprepare{Block: block1, View: view1_0, Sig: pppSig}

	// 1. HandleIstanbulPreprepare: every council member
	for _, name := range s.Council {
		nameToCN[name].VRankModule.HandleIstanbulPreprepare(block1, view1_0)
	}

	// 2. HandleVRankPreprepare: each candidate receives preprepare (optimistic: we deliver directly)
	// 3. HandleVRankCandidate: each validator receives each candidate's message (unless in UnresponsiveCands)
	for _, candName := range s.Candidates {
		if slices.Contains(s.UnresponsiveCands, candName) {
			continue // skip VRank message deliver
		}

		cand := nameToCN[candName]
		_ = cand.VRankModule.HandleVRankPreprepare(pppMsg)

		sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, blockNum, uint8(view1_0.Round.Uint64()), block1.Hash())
		candMsg := &vrank.VRankCandidate{
			BlockNumber: blockNum,
			Round:       uint8(view1_0.Round.Uint64()),
			BlockHash:   block1.Hash(),
			Sig:         sig,
			BlsSig:      blsSig,
		}
		for _, valName := range s.Council {
			err := nameToCN[valName].VRankModule.HandleVRankCandidate(candMsg)
			require.NoError(t, err)
		}
	}

	// 4. EvaluateCandidates(1, 0) from each node and assert ExpectedCfReports[i]
	for i, nodeName := range s.Nodes {
		report, err := nameToCN[nodeName].VRankModule.EvaluateCandidates(blockNum, 0)
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

func TestVRankModuleHandlers(t *testing.T) {
	scenarios := []VRankScenario{
		{
			Name:       "happy path and non-council get empty report",
			Nodes:      []string{"N1", "N2", "N3"},
			Council:    []string{"N1"},
			Candidates: []string{"N2"},
			Proposer:   "N1",
			ExpectedCfReports: [][]string{
				{}, // N1: Council(1), valid Report
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
		view1_0 = &bft.View{Sequence: big.NewInt(1), Round: common.Big0}
	)

	t.Run("permissionless fork is disabled", func(t *testing.T) {
		val := newCN(t, withHardfork("osaka"), withGenesis())
		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		prepreparedTime, _, _ := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		mustNotPop(t, val.sub)
	})

	t.Run("the proposer should not start collection when not in the next council", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, validator, candidate := cns[0], cns[1], cns[2]

		// proposer is not in the next council, so it should only broadcast and does not start collection.
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{validator.Addr}, nil).Times(2)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(2)
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(2)

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
		cns, valset, _ := newCNMulti(t, 3)
		proposer, nonProposer, candidate := cns[0], cns[1], cns[2]

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(3)
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(3)

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
		view1_0 = &bft.View{Sequence: big.NewInt(1), Round: common.Big0}
	)

	t.Run("permissionless fork is disabled", func(t *testing.T) {
		cand := newCN(t, withHardfork("osaka"), withGenesis())
		cand.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		cand.VRankModule.HandleVRankPreprepare(&vrank.VRankPreprepare{Block: block1, View: view1_0})
		mustNotPop(t, cand.sub)
	})

	t.Run("validators should not broadcast", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, nonProposer, candidate := cns[0], cns[1], cns[2]

		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).AnyTimes()
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).AnyTimes()

		pppSig := signVRankPreprepare(t, proposer.VRankModule, proposer.Key, block1.NumberU64(), 0, block1.Hash())
		pppMsg := &vrank.VRankPreprepare{Block: block1, View: view1_0, Sig: pppSig}

		proposer.VRankModule.HandleVRankPreprepare(pppMsg)
		nonProposer.VRankModule.HandleVRankPreprepare(pppMsg)
		candidate.VRankModule.HandleVRankPreprepare(pppMsg)

		req := mustPop(t, candidate.sub)
		assert.Equal(t, []common.Address{proposer.Addr, nonProposer.Addr}, req.Targets)

		mustNotPop(t, nonProposer.sub)
		mustNotPop(t, proposer.sub)
	})

	t.Run("candidate should broadcast to the round-specific committee", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 4)
		proposer, candidate, round1Val1, round1Val2 := cns[0], cns[1], cns[2], cns[3]

		view1_1 := &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(1)}
		pppSig := signVRankPreprepare(t, proposer.VRankModule, proposer.Key, block1.NumberU64(), 1, block1.Hash())
		pppMsg := &vrank.VRankPreprepare{Block: block1, View: view1_1, Sig: pppSig}

		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(1)).Return(proposer.Addr, nil).AnyTimes()
		valset.EXPECT().GetCommittee(uint64(1), uint64(1)).Return([]common.Address{round1Val1.Addr, round1Val2.Addr}, nil).Times(1)

		err := candidate.VRankModule.HandleVRankPreprepare(pppMsg)
		require.NoError(t, err)

		req := mustPop(t, candidate.sub)
		assert.Equal(t, []common.Address{round1Val1.Addr, round1Val2.Addr}, req.Targets)
		assert.Len(t, req.Msg.(*vrank.VRankCandidate).BlsSig, 96)
	})

	t.Run("exact replay must not rebroadcast VRankCandidate", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, candidate, validator := cns[0], cns[1], cns[2]

		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).AnyTimes()
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{validator.Addr}, nil).Times(1)

		pppSig := signVRankPreprepare(t, proposer.VRankModule, proposer.Key, block1.NumberU64(), 0, block1.Hash())
		pppMsg := &vrank.VRankPreprepare{Block: block1, View: view1_0, Sig: pppSig}

		require.NoError(t, candidate.VRankModule.HandleVRankPreprepare(pppMsg))
		req := mustPop(t, candidate.sub)
		assert.Equal(t, []common.Address{validator.Addr}, req.Targets)

		require.NoError(t, candidate.VRankModule.HandleVRankPreprepare(pppMsg))
		mustNotPop(t, candidate.sub)
	})

	t.Run("same view with different block hash must not rebroadcast VRankCandidate", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, candidate, validator := cns[0], cns[1], cns[2]

		altBlock := types.NewBlockWithHeader(&types.Header{
			Number:     big.NewInt(1),
			ParentHash: common.HexToHash("0x01"),
		})

		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).AnyTimes()
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{validator.Addr}, nil).Times(1)

		pppSig1 := signVRankPreprepare(t, proposer.VRankModule, proposer.Key, block1.NumberU64(), 0, block1.Hash())
		pppMsg1 := &vrank.VRankPreprepare{Block: block1, View: view1_0, Sig: pppSig1}
		pppSig2 := signVRankPreprepare(t, proposer.VRankModule, proposer.Key, altBlock.NumberU64(), 0, altBlock.Hash())
		pppMsg2 := &vrank.VRankPreprepare{Block: altBlock, View: view1_0, Sig: pppSig2}

		require.NoError(t, candidate.VRankModule.HandleVRankPreprepare(pppMsg1))
		req := mustPop(t, candidate.sub)
		assert.Equal(t, []common.Address{validator.Addr}, req.Targets)

		require.NoError(t, candidate.VRankModule.HandleVRankPreprepare(pppMsg2))
		mustNotPop(t, candidate.sub)
	})

	t.Run("non-proposer signature should be rejected by candidate", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, nonProposer, candidate := cns[0], cns[1], cns[2]

		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).AnyTimes()

		// signed by nonProposer instead of proposer
		badSig := signVRankPreprepare(t, nonProposer.VRankModule, nonProposer.Key, block1.NumberU64(), 0, block1.Hash())
		badMsg := &vrank.VRankPreprepare{Block: block1, View: view1_0, Sig: badSig}

		err := candidate.VRankModule.HandleVRankPreprepare(badMsg)
		assert.ErrorIs(t, err, vrank.ErrMsgFromNonProposer)
		mustNotPop(t, candidate.sub)
	})
}

func TestHandleVRankCandidate(t *testing.T) {
	var (
		block1  = types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		view1_0 = &bft.View{Sequence: big.NewInt(1), Round: common.Big0}
	)
	t.Run("permissionless fork is disabled", func(t *testing.T) {
		val := newCN(t, withHardfork("osaka"))
		val.VRankModule.HandleVRankCandidate(&vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: [crypto.SignatureLength]byte{}})
		mustNotPop(t, val.sub)
	})

	t.Run("no nodes should broadcast", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, nonProposer, candidate := cns[0], cns[1], cns[2]
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: [crypto.SignatureLength]byte{}}

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{proposer.Addr, nonProposer.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(3)
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(3)

		proposer.VRankModule.HandleVRankCandidate(&msg)
		nonProposer.VRankModule.HandleVRankCandidate(&msg)
		candidate.VRankModule.HandleVRankCandidate(&msg)

		mustNotPop(t, proposer.sub)
		mustNotPop(t, nonProposer.sub)
		mustNotPop(t, candidate.sub)
	})

	t.Run("the proposer should not collect when not in the next council", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 3)
		proposer, validator, candidate := cns[0], cns[1], cns[2]
		sig, blsSig := signVRankCandidate(t, candidate.VRankModule, candidate.Key, candidate.BlsKey, block1.NumberU64(), uint8(view1_0.Round.Uint64()), block1.Hash())
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: sig, BlsSig: blsSig}

		// proposer is not in the next council, so it should only broadcast and does not start collection.
		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{validator.Addr}, nil).Times(3)
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(proposer.Addr, nil).Times(2)
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candidate.Addr}, nil).Times(1)

		proposer.VRankModule.HandleIstanbulPreprepare(block1, view1_0) // this won't happen in production
		err := proposer.VRankModule.HandleVRankCandidate(&msg)
		require.ErrorIs(t, err, vrank.ErrPrepreparedViewNotSet)
		prepreparedTime, _, candMap := proposer.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.True(t, prepreparedTime.IsZero())
		assert.Nil(t, candMap)

		validator.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		err = validator.VRankModule.HandleVRankCandidate(&msg)
		require.NoError(t, err)
		prepreparedTime, _, candMap = validator.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.False(t, prepreparedTime.IsZero())
		assert.Equal(t, 1, len(candMap))
	})

	t.Run("future messages", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 2)
		val, cand := cns[0], cns[1]

		block2 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		sigFutureBlock, blsSigFutureBlock := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block2.NumberU64(), 0, block2.Hash())
		sigFutureRound, blsSigFutureRound := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block1.NumberU64(), 1, block1.Hash())

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{cand.Addr}, nil).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)

		tcs := []struct {
			name    string
			msg     *vrank.VRankCandidate
			wantErr error
		}{
			{
				name: "future block number",
				msg:  &vrank.VRankCandidate{BlockNumber: 2, Round: 0, BlockHash: block2.Hash(), Sig: sigFutureBlock, BlsSig: blsSigFutureBlock}, wantErr: nil,
			},
			{
				name: "future round",
				msg:  &vrank.VRankCandidate{BlockNumber: 1, Round: 1, BlockHash: block1.Hash(), Sig: sigFutureRound, BlsSig: blsSigFutureRound}, wantErr: nil,
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
		cns, valset, _ := newCNMulti(t, 2)
		val, cand := cns[0], cns[1]

		sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block1.NumberU64(), uint8(view1_0.Round.Uint64()), block1.Hash())
		msg := vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: uint8(view1_0.Round.Uint64()), BlockHash: block1.Hash(), Sig: sig, BlsSig: blsSig}

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{cand.Addr}, nil).AnyTimes()
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

	t.Run("stale messages should be discarded", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 2)
		val, cand := cns[0], cns[1]
		block2 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		view2_0 := &bft.View{Sequence: big.NewInt(2), Round: common.Big0}
		sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block1.NumberU64(), 0, block1.Hash())
		msg := &vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: 0, BlockHash: block1.Hash(), Sig: sig, BlsSig: blsSig}

		valset.EXPECT().GetCommittee(uint64(2), uint64(0)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(2), uint64(0)).Return(common.Address{}, nil).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block2, view2_0)

		err := val.VRankModule.HandleVRankCandidate(msg)
		require.NoError(t, err)

		_, _, candMap := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: 1, R: 0})
		assert.Nil(t, candMap, "stale messages should not be stored")
	})

	t.Run("epoch boundary future messages should be stored until evaluation", func(t *testing.T) {
		cns, valset, _ := newCNMulti(t, 2)
		val, cand := cns[0], cns[1]
		blockNum := params.DefaultVRankEpoch - 1
		nextBlockNum := params.DefaultVRankEpoch
		block := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(blockNum)})
		nextBlock := types.NewBlockWithHeader(&types.Header{Number: new(big.Int).SetUint64(nextBlockNum)})
		view := &bft.View{Sequence: new(big.Int).SetUint64(blockNum), Round: common.Big0}
		sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, nextBlockNum, 0, nextBlock.Hash())
		msg := &vrank.VRankCandidate{BlockNumber: nextBlockNum, Round: 0, BlockHash: nextBlock.Hash(), Sig: sig, BlsSig: blsSig}

		valset.EXPECT().GetCommittee(blockNum, uint64(0)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(blockNum, uint64(0)).Return(common.Address{}, nil).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block, view)

		err := val.VRankModule.HandleVRankCandidate(msg)
		require.NoError(t, err)

		_, _, candMap := val.VRankModule.collector.GetViewData(vrank.ViewKey{N: nextBlockNum, R: 0})
		require.Len(t, candMap, 1)
		assert.Equal(t, msg.BlockHash, candMap[cand.Addr].Msg.BlockHash)
	})

	newBLSSetup := func() (val, cand *CN) {
		cns, vs, _ := newCNMulti(t, 2)
		v, c := cns[0], cns[1]
		vs.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{v.Addr}, nil).AnyTimes()
		vs.EXPECT().GetProposer(uint64(1), uint64(0)).Return(v.Addr, nil).AnyTimes()
		vs.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{c.Addr}, nil).AnyTimes()
		v.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		return v, c
	}

	t.Run("wrong BLS key is rejected", func(t *testing.T) {
		val, cand := newBLSSetup()
		wrongBlsKey, _ := bls.RandKey()
		sig, blsSig := signVRankCandidate(t, cand.VRankModule, cand.Key, wrongBlsKey, block1.NumberU64(), 0, block1.Hash())
		msg := &vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: 0, BlockHash: block1.Hash(), Sig: sig, BlsSig: blsSig}
		assert.ErrorIs(t, val.VRankModule.HandleVRankCandidate(msg), vrank.ErrInvalidCandidateBlsSig)
	})

	t.Run("corrupt BLS sig bytes are rejected", func(t *testing.T) {
		val, cand := newBLSSetup()
		randomBlsSig := [blstypes.SignatureLength]byte{}
		sig, _ := signVRankCandidate(t, cand.VRankModule, cand.Key, cand.BlsKey, block1.NumberU64(), 0, block1.Hash())
		msg := &vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: 0, BlockHash: block1.Hash(), Sig: sig, BlsSig: randomBlsSig}
		assert.ErrorIs(t, val.VRankModule.HandleVRankCandidate(msg), vrank.ErrInvalidCandidateBlsSig)
	})

	t.Run("missing BLS pubkey from Randao is rejected", func(t *testing.T) {
		cns, valset, randao := newCNMulti(t, 1)
		val := cns[0]
		var (
			candKey, _    = crypto.GenerateKey()
			candBlsKey, _ = bls.DeriveFromECDSA(candKey)
			candAddr      = crypto.PubkeyToAddress(candKey.PublicKey)
		)

		valset.EXPECT().GetCommittee(uint64(1), uint64(0)).Return([]common.Address{val.Addr}, nil).AnyTimes()
		valset.EXPECT().GetProposer(uint64(1), uint64(0)).Return(val.Addr, nil).AnyTimes()
		valset.EXPECT().GetCandTesting(uint64(1)).Return([]common.Address{candAddr}, nil).AnyTimes()
		// nil in Randao
		randao.EXPECT().GetBlsPubkey(candAddr, gomock.Any()).Return(nil, assert.AnError).AnyTimes()

		val.VRankModule.HandleIstanbulPreprepare(block1, view1_0)
		var (
			ecdsaSig, blsSig = signVRankCandidate(t, val.VRankModule, candKey, candBlsKey, block1.NumberU64(), 0, block1.Hash())
			msg              = &vrank.VRankCandidate{BlockNumber: block1.NumberU64(), Round: 0, BlockHash: block1.Hash(), Sig: ecdsaSig, BlsSig: blsSig}
		)
		assert.ErrorIs(t, val.VRankModule.HandleVRankCandidate(msg), vrank.ErrInvalidCandidateBlsSig)
	})
}
