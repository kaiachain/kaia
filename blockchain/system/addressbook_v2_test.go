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

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInstallAddressBookV2(t *testing.T) {
	backend := backends.NewSimulatedBackend(nil)
	defer backend.Close()

	header := backend.BlockChain().CurrentHeader()
	statedb, err := backend.BlockChain().StateAt(header.Root)
	require.NoError(t, err)

	logicAddr := common.HexToAddress("0xcafebebe")

	// Before install: no code at AddressBookAddr
	assert.Empty(t, statedb.GetCode(AddressBookAddr))

	err = InstallAddressBookV2(statedb, logicAddr)
	require.NoError(t, err)

	// Verify proxy code is set at 0x400
	assert.Equal(t, ERC1967ProxyV5Code, statedb.GetCode(AddressBookAddr))

	// Verify implementation slot points to logicAddr
	implSlot := statedb.GetState(AddressBookAddr, common.BytesToHash(ImplementationSlot))
	assert.Equal(t, lpad32(logicAddr), implSlot)
}

func TestEncodeNodeStateUpdate(t *testing.T) {
	rules := params.Rules{IsPrague: true} // enable intrinsic gas calculation

	tests := []struct {
		name       string
		validators valset.NodeStateMap
		epochSF    uint64
		// expected decoded values (sorted by address)
		wantAddrs    []common.Address
		wantStates   []uint8
		wantTimeouts []*big.Int
		wantEpochSF  *big.Int // nil means skip check (empty map case)
	}{
		{
			name:       "empty map",
			validators: valset.NodeStateMap{},
			wantAddrs:  []common.Address{},
		},
		{
			name: "single ValActive, no timeout",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {State: valset.ValActive},
			},
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValActive.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(0)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "ValPaused uses PausedTimeout",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {
					State:         valset.ValPaused,
					PausedTimeout: time.Unix(1000, 0),
				},
			},
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValPaused.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(1000)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "ValReady uses IdleTimeout",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {
					State:       valset.ValReady,
					IdleTimeout: time.Unix(2000, 0),
				},
			},
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValReady.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(2000)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "ValInactive uses IdleTimeout",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {
					State:       valset.ValInactive,
					IdleTimeout: time.Unix(3000, 0),
				},
			},
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValInactive.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(3000)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "ValExiting has zero timeout",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {State: valset.ValExiting},
			},
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValExiting.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(0)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "multiple validators sorted by address",
			validators: valset.NodeStateMap{
				common.HexToAddress("0xcc"): {State: valset.ValActive},
				common.HexToAddress("0xaa"): {
					State:         valset.ValPaused,
					PausedTimeout: time.Unix(500, 0),
				},
				common.HexToAddress("0xbb"): {
					State:       valset.ValInactive,
					IdleTimeout: time.Unix(999, 0),
				},
			},
			wantAddrs: []common.Address{
				common.HexToAddress("0xaa"),
				common.HexToAddress("0xbb"),
				common.HexToAddress("0xcc"),
			},
			wantStates:   []uint8{valset.ValPaused.ToUint8(), valset.ValInactive.ToUint8(), valset.ValActive.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(500), big.NewInt(999), big.NewInt(0)},
			wantEpochSF:  big.NewInt(0),
		},
		{
			name: "non-zero epochSF is encoded as 4th param",
			validators: valset.NodeStateMap{
				common.HexToAddress("0x01"): {State: valset.ValActive},
			},
			epochSF:      1,
			wantAddrs:    []common.Address{common.HexToAddress("0x01")},
			wantStates:   []uint8{valset.ValActive.ToUint8()},
			wantTimeouts: []*big.Int{big.NewInt(0)},
			wantEpochSF:  big.NewInt(1),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			from, msg, err := EncodeNodeStateUpdate(rules, tc.validators, tc.epochSF)
			require.NoError(t, err)

			// from is always SystemAddress
			assert.Equal(t, params.SystemAddress, from)

			// to is always AddressBookAddr
			assert.Equal(t, &AddressBookAddr, msg.To())

			if len(tc.validators) == 0 {
				return
			}

			// Decode the ABI-packed data to verify contents
			method, err := AddressBookV2ABI.MethodById(msg.Data()[:4])
			require.NoError(t, err)
			assert.Equal(t, "processSystemTransition", method.Name)

			args, err := method.Inputs.Unpack(msg.Data()[4:])
			require.NoError(t, err)
			require.Len(t, args, 4)

			decodedAddrs := args[0].([]common.Address)
			decodedStates := args[1].([]uint8)
			decodedTimeouts := args[2].([]*big.Int)
			decodedEpochSF := args[3].(*big.Int)

			assert.Equal(t, tc.wantAddrs, decodedAddrs, "addresses should be sorted")
			assert.Equal(t, tc.wantStates, decodedStates)
			for i, want := range tc.wantTimeouts {
				assert.Equal(t, 0, want.Cmp(decodedTimeouts[i]), "timeout[%d]: want %s got %s", i, want, decodedTimeouts[i])
			}
			if tc.wantEpochSF != nil {
				assert.Equal(t, 0, tc.wantEpochSF.Cmp(decodedEpochSF), "epochSF: want %s got %s", tc.wantEpochSF, decodedEpochSF)
			}
		})
	}
}

// TestReadAddressBookV2BlsAll tests that BLS pub/pop keys can be read from ABv2.
func TestReadAddressBookV2BlsAll(t *testing.T) {
	config, _ := MakeTestPermissionlessConfig(3)

	alloc, err := AllocPermissionless(config)
	require.NoError(t, err)

	simBackend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	defer simBackend.Close()

	chain := simBackend.BlockChain()
	header := chain.CurrentHeader()
	statedb, err := chain.StateAt(header.Root)
	require.NoError(t, err)

	readBackend, err := backends.NewStateBlockchainContractBackend(chain, statedb)
	require.NoError(t, err)

	blsInfos, err := ReadAddressBookV2BlsAll(readBackend, nil)
	require.NoError(t, err)
	assert.Len(t, blsInfos, 3)

	// Verify each nodeId has matching BLS keys from config
	for i, nodeId := range config.NodeIds {
		info, ok := blsInfos[nodeId]
		require.True(t, ok, "nodeId %s should have BLS info", nodeId.Hex())
		assert.Equal(t, config.NodeInfos[i].BlsInfo.PublicKey, info.PublicKey)
		assert.Equal(t, config.NodeInfos[i].BlsInfo.Pop, info.Pop)
		assert.NoError(t, info.VerifyErr)
	}
}
