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
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestCore_sendPrepare(t *testing.T) {
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

		istCore.sendPrepare()
		istCore.Stop()
		mockCtrl.Finish()
	}
}

func BenchmarkMsgCmp(b *testing.B) {
	getEmptySubject := func() istanbul.Subject {
		return istanbul.Subject{
			View: &istanbul.View{
				Round:    big.NewInt(0),
				Sequence: big.NewInt(0),
			},
			Digest:   common.HexToHash("1"),
			PrevHash: common.HexToHash("2"),
		}
	}
	s1, s2 := getEmptySubject(), getEmptySubject()

	// Worst
	b.Run("reflect.DeepEqual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflect.DeepEqual(s1, s2)
		}
	})

	// Better
	b.Run("EqualImpl", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			s1.Equal(&s2)
		}
	})
}

func TestSubjectCmp(t *testing.T) {
	genRandomHash := func(n int) common.Hash {
		b := make([]byte, n)
		_, err := crand.Read(b)
		assert.Nil(t, err)
		return common.BytesToHash(b)
	}
	genRandomInt := func(min, max int) int64 {
		return int64(rand.Intn(max-min) + min)
	}
	genSubject := func(min, max int) *istanbul.Subject {
		round, seq := big.NewInt(genRandomInt(min, max)), big.NewInt(genRandomInt(min, max))
		digest, prevHash := genRandomHash(max), genRandomHash(max)
		return &istanbul.Subject{
			View: &istanbul.View{
				Round:    round,
				Sequence: seq,
			},
			Digest:   digest,
			PrevHash: prevHash,
		}
	}
	copySubject := func(s *istanbul.Subject) *istanbul.Subject {
		r := new(istanbul.Subject)
		v := new(istanbul.View)
		r.Digest = s.Digest
		r.PrevHash = s.PrevHash
		v.Round = new(big.Int).SetUint64(s.View.Round.Uint64())
		v.Sequence = new(big.Int).SetUint64(s.View.Sequence.Uint64())
		r.View = v
		return r
	}

	min, max, n := 1, 9999, 10000
	var identity bool
	var s1, s2 *istanbul.Subject
	for i := 0; i < n; i++ {
		s1 = genSubject(min, max)
		if rand.Intn(2) == 0 {
			identity = true
			s2 = copySubject(s1)
		} else {
			identity = false
			s2 = genSubject(max+1, max*2)
		}
		e := s1.Equal(s2)
		if identity {
			assert.Equal(t, e, true)
		} else {
			assert.Equal(t, e, false)
		}
		assert.Equal(t, e, reflect.DeepEqual(s1, s2))
	}
}

func TestNilSubjectCmp(t *testing.T) {
	sbj := istanbul.Subject{
		View: &istanbul.View{
			Round:    big.NewInt(0),
			Sequence: big.NewInt(0),
		},
		Digest:   common.HexToHash("1"),
		PrevHash: common.HexToHash("2"),
	}
	var nilSbj *istanbul.Subject = nil

	assert.Equal(t, sbj.Equal(nil), false)
	assert.Equal(t, sbj.Equal(nilSbj), false)
	assert.Equal(t, nilSbj.Equal(&sbj), false)
	assert.Equal(t, nilSbj.Equal(nilSbj), true)
	assert.Equal(t, nilSbj.Equal(nil), true)
}
