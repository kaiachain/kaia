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

	"github.com/kaiachain/kaia/common"
	abv2data "github.com/kaiachain/kaia/contracts/bindings/abv2data"
	addressbookv2contract "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
)

// defaultGenesisStake is staked per validator by MakeABv2AllocConfig.
var defaultGenesisStake = new(big.Int).Mul(big.NewInt(5_000_000), big.NewInt(params.KAIA))

// ABv2NodeSpec is what a caller supplies per genesis validator. Everything else in
// NodeInfo is derived by makeABv2NodeInfos.
type ABv2NodeSpec struct {
	NodeID  common.Address
	BlsInfo addressbookv2contract.BlsPublicKeyInfo
}

// MakeABv2AllocConfig fills an AllocPermissionlessConfig with the genesis defaults for the
// given validators, staking defaultGenesisStake each. Callers that expose any of these as
// flags or governance (homi) overwrite those fields on the returned config.
func MakeABv2AllocConfig(owner common.Address, specs []ABv2NodeSpec) *AllocPermissionlessConfig {
	nodeIds, nodeInfos := makeABv2NodeInfos(specs)

	stakeAmts := make([]*big.Int, len(specs))
	for i := range stakeAmts {
		stakeAmts[i] = new(big.Int).Set(defaultGenesisStake)
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
			PfsThreshold:            big.NewInt(valset.DefaultPfsThreshold),
			CfsThreshold:            big.NewInt(valset.DefaultCfsThreshold),
			PauseTimeout:            big.NewInt(int64(valset.DefaultValPausedTimeout.Seconds())),
			IdleTimeout:             big.NewInt(int64(valset.DefaultValIdleTimeout.Seconds())),
			MaxNodeCount:            big.NewInt(valset.DefaultMaxNodeCount),
			MaxValActivePausedCount: big.NewInt(valset.DefaultMaxValActivePausedCount),
			MaxCandReadyCount:       big.NewInt(valset.DefaultMaxCandReadyCount),
			KefAddress:              owner,
			KifAddress:              owner,
			KpfAddress:              owner,
		},
	}
}

// makeABv2NodeInfos builds the genesis validator list for AllocPermissionlessConfig.
// RewardAddress is derived rather than reused: the ABv2DataContract constructor rejects
// a reward address equal to the node address.
func makeABv2NodeInfos(specs []ABv2NodeSpec) ([]common.Address, []addressbookv2contract.NodeInfo) {
	nodeIds := make([]common.Address, len(specs))
	infos := make([]addressbookv2contract.NodeInfo, len(specs))
	for i, spec := range specs {
		nodeIds[i] = spec.NodeID
		infos[i] = addressbookv2contract.NodeInfo{
			Manager:       spec.NodeID,
			RewardAddress: common.BytesToAddress(crypto.Keccak256(spec.NodeID.Bytes())),
			VoterAddress:  spec.NodeID,
			TimeoutAt:     new(big.Int),
			GcId:          big.NewInt(int64(i + 1)),
			BlsInfo:       spec.BlsInfo,
			State:         valset.ValActive.ToUint8(),
		}
	}
	return nodeIds, infos
}
