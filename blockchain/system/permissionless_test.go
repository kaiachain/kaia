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
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	abv2data "github.com/kaiachain/kaia/contracts/bindings/abv2data"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func makeTestPermissionlessConfig(t *testing.T, n int) (*AllocPermissionlessConfig, []*ecdsa.PrivateKey) {
	t.Helper()

	ownerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	owner := crypto.PubkeyToAddress(ownerKey.PublicKey)

	nodeKeys := make([]*ecdsa.PrivateKey, n)
	nodeIds := make([]common.Address, n)
	nodeInfos := make([]addressbookv2contract.NodeInfo, n)
	stakeAmts := make([]*big.Int, n)
	stakeAmt := new(big.Int).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA))

	for i := range n {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		addr := crypto.PubkeyToAddress(key.PublicKey)
		blsSk, err := bls.DeriveFromECDSA(key)
		require.NoError(t, err)

		nodeKeys[i] = key
		nodeIds[i] = addr
		stakeAmts[i] = new(big.Int).Set(stakeAmt)
		nodeInfos[i] = addressbookv2contract.NodeInfo{
			Manager:       addr,
			RewardAddress: common.BytesToAddress(crypto.Keccak256(addr.Bytes())),
			VoterAddress:  addr,
			TimeoutAt:     new(big.Int),
			GcId:          big.NewInt(int64(i + 1)),
			BlsInfo: addressbookv2contract.BlsPublicKeyInfo{
				PublicKey: blsSk.PublicKey().Marshal(),
				Pop:       bls.PopProve(blsSk).Marshal(),
			},
			State: valset.ValActive.ToUint8(),
		}
	}

	return &AllocPermissionlessConfig{
		Owner:     owner,
		NodeIds:   nodeIds,
		NodeInfos: nodeInfos,
		StakeAmts: stakeAmts,
		DataConfig: abv2data.IABv2DataContractInitData{
			InitialOwner:            owner,
			InitialSuspender:        owner,
			InitialConfigurator:     owner,
			PfsThreshold:            big.NewInt(2),
			CfsThreshold:            big.NewInt(300),
			PauseTimeout:            big.NewInt(8 * 3600),
			IdleTimeout:             big.NewInt(30 * 86400),
			MaxNodeCount:            big.NewInt(100),
			MaxValActivePausedCount: big.NewInt(50),
			MaxCandReadyCount:       big.NewInt(3),
			KefAddress:              common.HexToAddress("0x1111"),
			KifAddress:              common.HexToAddress("0x2222"),
			KpfAddress:              common.HexToAddress("0x3333"),
		},
	}, nodeKeys
}

func TestAllocPermissionlessInstallsInitializedABv2(t *testing.T) {
	config, _ := makeTestPermissionlessConfig(t, 4)

	alloc, err := AllocPermissionless(config)
	require.NoError(t, err)
	require.Contains(t, alloc, RegistryAddr)
	require.Contains(t, alloc, AddressBookAddr)
	require.NotEmpty(t, alloc[AddressBookAddr].Code)

	for _, info := range config.NodeInfos {
		assert.NotContains(t, alloc, info.Manager)
		require.Contains(t, alloc, info.StakingContract)
		assert.Equal(t, config.StakeAmts[0], alloc[info.StakingContract].Balance)
	}

	backend := backends.NewSimulatedBackend(blockchain.GenesisAlloc(alloc))
	caller, err := addressbookv2contract.NewAddressBookV2Caller(AddressBookAddr, backend)
	require.NoError(t, err)

	suspender, err := caller.GetSuspender(&bind.CallOpts{})
	require.NoError(t, err)
	assert.Equal(t, config.DataConfig.InitialSuspender, suspender)

	profiles, err := caller.GetAllProfiles(&bind.CallOpts{})
	require.NoError(t, err)
	require.Len(t, profiles, len(config.NodeIds))
	for i, profile := range profiles {
		assert.Equal(t, config.NodeIds[i], profile.NodeId)
		assert.Equal(t, config.NodeInfos[i].StakingContract, profile.StakingContract)
		assert.Equal(t, valset.ValActive.ToUint8(), profile.State)
	}
}

func TestAllocPermissionlessMismatchedLengths(t *testing.T) {
	config, _ := makeTestPermissionlessConfig(t, 2)
	config.StakeAmts = config.StakeAmts[:1]

	_, err := AllocPermissionless(config)
	require.ErrorContains(t, err, "mismatched lengths")
}
