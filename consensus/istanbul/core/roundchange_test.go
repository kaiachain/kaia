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

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundChangeSetBoundsMessagesPerRound(t *testing.T) {
	addrs := []common.Address{
		common.HexToAddress("0x1"),
		common.HexToAddress("0x2"),
		common.HexToAddress("0x3"),
	}
	rcs := newRoundChangeSet(valset.NewAddressSet(addrs), 2)
	round := big.NewInt(1)
	currentRound := big.NewInt(0)

	for _, addr := range addrs[:2] {
		_, err := rcs.Add(currentRound, round, &bft.Message{Address: addr})
		require.NoError(t, err)
	}

	_, err := rcs.Add(currentRound, round, &bft.Message{Address: addrs[2]})
	require.ErrorIs(t, err, errRoundChangeMessageLimit)
	assert.Equal(t, 2, rcs.roundChanges[round.Uint64()].Size())

	// A duplicate sender may still replace its prior message after the limit.
	count, err := rcs.Add(currentRound, round, &bft.Message{Address: addrs[0]})
	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestRoundChangeSetCheckIgnoresChangedUnsignedAttachment(t *testing.T) {
	addr := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{addr}), 1)
	round := big.NewInt(1)
	msg := &bft.Message{
		PrevHash:  common.HexToHash("0x11"),
		Code:      bft.MsgRoundChange,
		Address:   addr,
		Msg:       []byte{0x01, 0x02},
		Signature: []byte{0x03},
		Evidence:  []byte{0x04},
	}
	_, err := rcs.Add(big.NewInt(0), round, msg.WithoutEvidence())
	require.NoError(t, err)

	duplicate := *msg
	duplicate.Evidence = []byte{0x06, 0x07}

	require.ErrorIs(t, rcs.Check(big.NewInt(0), round, &duplicate), errIgnored)
}

// TestRoundChangeSetReplacesOnlyWithHigherClaim checks that a sender's ROUND
// CHANGE for a round is replaced only by one with a strictly higher prepared
// claim. Otherwise a sender could alternate signed messages, each new to the
// set, to have its Evidence verified and relayed again on every switch.
func TestRoundChangeSetReplacesOnlyWithHigherClaim(t *testing.T) {
	addr := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{addr}), 1)
	round := big.NewInt(3)
	roundChange := func(claimRound int64) *bft.Message {
		payload := &bft.RoundChange{View: &bft.View{Sequence: big.NewInt(1), Round: round}}
		if claimRound >= 0 {
			payload.Prepared = &bft.PreparedClaim{Round: big.NewInt(claimRound), Digest: common.HexToHash("0xd1")}
		}
		encoded, err := bft.Encode(payload)
		require.NoError(t, err)
		return &bft.Message{Code: bft.MsgRoundChange, Address: addr, Msg: encoded}
	}
	_, err := rcs.Add(big.NewInt(0), round, roundChange(1))
	require.NoError(t, err)

	for _, tc := range []struct {
		name string
		msg  *bft.Message
		want error
	}{
		{"no claim", roundChange(-1), errIgnored},
		{"lower claim", roundChange(0), errIgnored},
		{"same claim", roundChange(1), errIgnored},
		{"higher claim", roundChange(2), nil},
	} {
		require.ErrorIs(t, rcs.Check(big.NewInt(0), round, tc.msg), tc.want, tc.name)
	}

	// Add, the only writer, enforces the same rule.
	_, err = rcs.Add(big.NewInt(0), round, roundChange(0))
	require.ErrorIs(t, err, errIgnored)
	_, err = rcs.Add(big.NewInt(0), round, roundChange(2))
	require.NoError(t, err)
	assert.Equal(t, roundChange(2).Msg, rcs.Values(round)[0].Msg)
}

// A zero quorum means no ROUND CHANGE can contribute to progress, so not even a
// freshly created bucket may retain one.
func TestRoundChangeSetRejectsEveryMessageWhenLimitIsZero(t *testing.T) {
	addr := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{addr}), 0)

	_, err := rcs.Add(big.NewInt(0), big.NewInt(1), &bft.Message{Address: addr})

	require.ErrorIs(t, err, errRoundChangeMessageLimit)
	assert.Empty(t, rcs.roundChanges)
}

func TestRoundChangeSetClearDropsOlderBuckets(t *testing.T) {
	addr := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{addr}), 1)
	for _, round := range []int64{1, 2} {
		_, err := rcs.Add(big.NewInt(0), big.NewInt(round), &bft.Message{Address: addr})
		require.NoError(t, err)
	}

	rcs.Clear(big.NewInt(2))

	assert.NotContains(t, rcs.roundChanges, uint64(1))
	assert.Contains(t, rcs.roundChanges, uint64(2))
}

func TestRoundChangeSetEnforcesFutureRoundWindow(t *testing.T) {
	src := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{src}), 1)
	currentRound := big.NewInt(0)

	_, err := rcs.Add(currentRound, big.NewInt(maxRoundChangeRoundsAhead), &bft.Message{Address: src})
	require.NoError(t, err)

	_, err = rcs.Add(currentRound, big.NewInt(maxRoundChangeRoundsAhead+1), &bft.Message{Address: src})
	require.ErrorIs(t, err, errRoundChangeTooFar)
	assert.Len(t, rcs.roundChanges, 1)
}

func TestRoundChangeSetRejectsRoundOutsideUint64(t *testing.T) {
	src := common.HexToAddress("0x1")
	rcs := newRoundChangeSet(valset.NewAddressSet([]common.Address{src}), 1)
	tooLarge := new(big.Int).Lsh(big.NewInt(1), 64)

	_, err := rcs.Add(big.NewInt(0), tooLarge, &bft.Message{Address: src})

	require.ErrorIs(t, err, bft.ErrInvalidMessage)
	assert.Empty(t, rcs.roundChanges)
}

func TestHandleRoundChangeEnforcesFutureRoundWindow(t *testing.T) {
	src := common.HexToAddress("0x1")
	qualified := valset.NewAddressSet([]common.Address{src})
	current := newRoundState(&bft.View{Sequence: big.NewInt(1), Round: big.NewInt(0)}, qualified, common.Hash{}, nil, nil, nil, nil)
	current.committee = qualified
	current.requiredMessageCount = 2
	c := &core{
		state:          StateAcceptRequest,
		logger:         logger.NewWith(),
		current:        current,
		roundChangeSet: newRoundChangeSet(qualified, current.requiredMessageCount),
	}

	boundary := roundChangeMessage(t, src, maxRoundChangeRoundsAhead)
	err := c.handleRoundChange(boundary, src)
	require.ErrorIs(t, err, errIgnored)
	assert.Contains(t, c.roundChangeSet.roundChanges, uint64(maxRoundChangeRoundsAhead))

	tooFar := roundChangeMessage(t, src, maxRoundChangeRoundsAhead+1)
	err = c.handleRoundChange(tooFar, src)
	require.ErrorIs(t, err, errRoundChangeTooFar)
	assert.Len(t, c.roundChangeSet.roundChanges, 1)
}

func roundChangeMessage(t *testing.T, src common.Address, round uint64) *bft.Message {
	t.Helper()
	payload, err := bft.Encode(&bft.RoundChange{View: &bft.View{
		Sequence: big.NewInt(1),
		Round:    new(big.Int).SetUint64(round),
	}})
	require.NoError(t, err)
	return &bft.Message{Address: src, Code: bft.MsgRoundChange, Msg: payload}
}
