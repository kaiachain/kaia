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

package system

import (
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/contracts/bindings/multicall"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewNodeMapFromABv2(t *testing.T) {
	addr1 := common.HexToAddress("0x0001")
	addr2 := common.HexToAddress("0x0002")
	addr3 := common.HexToAddress("0x0003")

	profiles := []multicall.Profile{
		{NodeId: addr1, State: valset.ValReady.ToUint8(), TimeoutAt: big.NewInt(100)},
		{NodeId: addr2, State: valset.ValPaused.ToUint8(), TimeoutAt: big.NewInt(200)},
		{NodeId: addr3, State: valset.ValActive.ToUint8(), TimeoutAt: big.NewInt(300)},
	}
	stakingAmounts := []*big.Int{
		new(big.Int).Mul(big.NewInt(3), big.NewInt(params.KAIA)),
		new(big.Int).Mul(big.NewInt(4), big.NewInt(params.KAIA)),
		new(big.Int).Mul(big.NewInt(5), big.NewInt(params.KAIA)),
	}

	nodes, err := NewNodeMapFromABv2(profiles, stakingAmounts)
	require.NoError(t, err)

	assert.Equal(t, valset.ValReady, nodes[addr1].State)
	assert.Equal(t, uint64(3), nodes[addr1].StakingAmount)
	assert.Equal(t, time.Unix(100, 0), nodes[addr1].IdleTimeout)
	assert.True(t, nodes[addr1].PausedTimeout.IsZero())

	assert.Equal(t, valset.ValPaused, nodes[addr2].State)
	assert.Equal(t, uint64(4), nodes[addr2].StakingAmount)
	assert.Equal(t, time.Unix(200, 0), nodes[addr2].PausedTimeout)
	assert.True(t, nodes[addr2].IdleTimeout.IsZero())

	assert.Equal(t, valset.ValActive, nodes[addr3].State)
	assert.Equal(t, uint64(5), nodes[addr3].StakingAmount)
	assert.True(t, nodes[addr3].IdleTimeout.IsZero())
	assert.True(t, nodes[addr3].PausedTimeout.IsZero())
}

func TestNewNodeMapFromABv2LengthMismatch(t *testing.T) {
	_, err := NewNodeMapFromABv2([]multicall.Profile{{}}, nil)
	require.ErrorContains(t, err, "profile/staking amount length mismatch")
}

func TestABv2SnapshotTransitionParam(t *testing.T) {
	snapshot := &ABv2Snapshot{
		ABv2TransitionParam: ABv2TransitionParam{
			EpochVACount:            10,
			PfsThreshold:            20,
			CfsThreshold:            30,
			IdleTimeout:             40 * time.Second,
			PauseTimeout:            50 * time.Second,
			MaxValActivePausedCount: 60,
		},
	}

	param := snapshot.TransitionParam()

	assert.Equal(t, uint64(10), param.EpochVACount)
	assert.Equal(t, uint64(20), param.PfsThreshold)
	assert.Equal(t, uint64(30), param.CfsThreshold)
	assert.Equal(t, 40*time.Second, param.IdleTimeout)
	assert.Equal(t, 50*time.Second, param.PauseTimeout)
	assert.Equal(t, uint64(60), param.MaxValActivePausedCount)
}
