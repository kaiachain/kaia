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

package staking

import (
	"math"
	"sort"

	"github.com/kaiachain/kaia/common"
)

// The gini coefficient of an empty set. As Gini is mathematically undefined with an empty set,
// so here we use -1 to notify the user.
var EmptyGini float64 = -1.0

// StakingInfo is the staking info to be used for block processing.
type StakingInfo struct {
	// The source block number where the staking info is captured.
	SourceBlockNum uint64 `json:"blockNum"`

	// The AddressBook triplets
	NodeIds          []common.Address `json:"councilNodeAddrs"`
	StakingContracts []common.Address `json:"councilStakingAddrs"`
	RewardAddrs      []common.Address `json:"councilRewardAddrs"`

	// Treasury fund addresses
	KEFAddr common.Address `json:"kefAddr"` // KEF contract address (or KCF, KIR)
	KIFAddr common.Address `json:"kifAddr"` // KIF contract address (or KFF, KGF, PoC)

	// Staking amounts of each staking contracts, in KAIA, rounded down.
	StakingAmounts []uint64 `json:"councilStakingAmounts"`

	// Computed fields
	consolidatedNodes  *[]consolidatedNode
	cachedGini         *float64
	cachedGiniMinStake uint64 // The minimum staking amount used to compute Gini coefficient.
}

// consolidatedNode is the refined staking information suitable for proposer selection.
// Sometimes a node would register multiple NodeIds in AddressBook,
// in which each entry has different StakingAddr and same RewardAddr.
// We treat those entries with common RewardAddr as one node.
//
// For example,
//
//	NodeAddrs      = [N1, N2, N3]
//	StakingAddrs   = [S1, S2, S3]
//	RewardAddrs    = [R1, R1, R3]
//	StakingAmounts = [A1, A2, A3]
//
// can be consolidated into
//
//	CN1 = {[N1,N2], [S1,S2], R1, A1+A2}
//	CN3 = {[N3],    [S3],    R3, A3}
type consolidatedNode struct {
	NodeIds          []common.Address
	StakingContracts []common.Address
	RewardAddr       common.Address // The common RewardAddr
	StakingAmount    uint64         // Sum of the staking amounts
}

// StakingInfoLegacy may have been persisted to database by the past versions.
// Past database may contain Gini fields, but they are ignored. Gini shall be computed on-demand.
// StakingInfoLegacy should only be used to read from database. A more compact StakingInfo
// shall be used when writing a new entry.
type StakingInfoLegacy struct {
	StakingInfo

	// Legacy treasury fund address fields for backward-compatibility
	KIRAddr common.Address `json:"KIRAddr"` // = KEFAddr
	PoCAddr common.Address `json:"PoCAddr"` // = KIFAddr
	KCFAddr common.Address `json:"kcfAddr"` // = KEFAddr
	KFFAddr common.Address `json:"kffAddr"` // = KIFAddr
}

// StakingInfoResponse is the response type for APIs
type StakingInfoResponse struct {
	StakingInfoLegacy

	UseGini bool    `json:"useGini"` // Whether the gini coefficient is used at the requested block number
	Gini    float64 `json:"gini"`    // The gini coefficient at the requested block number. Returned regardless of `UseGini` value.
}

func (si *StakingInfo) ConsolidatedNodes() []consolidatedNode {
	if si.consolidatedNodes == nil {
		si.consolidatedNodes = si.consolidateNodes()
	}
	return *si.consolidatedNodes
}

func (si *StakingInfo) consolidateNodes() *[]consolidatedNode {
	// because Go map is not ordered, rList keeps track of the occurrence order of RewardAddrs.
	// to later arrange the consolidatedNodes.
	cmap := make(map[common.Address]*consolidatedNode)
	rList := make([]common.Address, 0, len(si.RewardAddrs))

	for i := range si.NodeIds {
		r := si.RewardAddrs[i]
		if cn, ok := cmap[r]; ok {
			cn.NodeIds = append(cn.NodeIds, si.NodeIds[i])
			cn.StakingContracts = append(cn.StakingContracts, si.StakingContracts[i])
			cn.StakingAmount += si.StakingAmounts[i]
		} else {
			cmap[r] = &consolidatedNode{
				NodeIds:          []common.Address{si.NodeIds[i]},
				StakingContracts: []common.Address{si.StakingContracts[i]},
				RewardAddr:       r,
				StakingAmount:    si.StakingAmounts[i],
			}
			rList = append(rList, r)
		}
	}

	carr := make([]consolidatedNode, 0, len(cmap))
	for _, r := range rList {
		carr = append(carr, *cmap[r])
	}
	return &carr
}

// Returns the Gini coefficient among the staking amounts that are greater than or equal to minStake.
// The amounts are first consolidated by RewardAddr, filtered by minStake, and then summarized to Gini.
func (si *StakingInfo) Gini(minStake uint64) float64 {
	// Cache hits only if the same minStake is used.
	if si.cachedGini == nil || si.cachedGiniMinStake != minStake {
		g := si.computeGini(minStake)
		si.cachedGini = &g
		si.cachedGiniMinStake = minStake
	}
	return *si.cachedGini
}

func (si *StakingInfo) computeGini(minStake uint64) float64 {
	cnodes := si.ConsolidatedNodes()
	amounts := make(sort.Float64Slice, 0, len(cnodes))

	for _, cnode := range cnodes {
		if cnode.StakingAmount >= minStake {
			amounts = append(amounts, float64(cnode.StakingAmount))
		}
	}

	return computeGini(amounts)
}

func computeGini(amounts sort.Float64Slice) float64 {
	if len(amounts) == 0 {
		return EmptyGini
	}

	sort.Sort(amounts)

	sumOfAbsoluteDifferences := float64(0)
	subSum := float64(0)

	for i, x := range amounts {
		temp := x*float64(i) - subSum
		sumOfAbsoluteDifferences = sumOfAbsoluteDifferences + temp
		subSum = subSum + x
	}

	result := sumOfAbsoluteDifferences / subSum / float64(len(amounts))
	result = math.Round(result*100) / 100
	return result
}

// Parse from persisted data
func (sl *StakingInfoLegacy) ToStakingInfo() *StakingInfo {
	si := &sl.StakingInfo

	// Try to fill treasury fund addresses from legacy fields that may have been persisted in the past.
	if !common.EmptyAddress(sl.KCFAddr) {
		si.KEFAddr = sl.KCFAddr
	}
	if !common.EmptyAddress(sl.KFFAddr) {
		si.KIFAddr = sl.KFFAddr
	}
	if !common.EmptyAddress(sl.KIRAddr) {
		si.KEFAddr = sl.KIRAddr
	}
	if !common.EmptyAddress(sl.PoCAddr) {
		si.KIFAddr = sl.PoCAddr
	}

	return si
}

// Convert to API response
func (si *StakingInfo) ToResponse(useGini bool, minStake uint64) *StakingInfoResponse {
	return &StakingInfoResponse{
		StakingInfoLegacy: StakingInfoLegacy{
			StakingInfo: *si,
			KIRAddr:     si.KEFAddr,
			PoCAddr:     si.KIFAddr,
			KCFAddr:     si.KEFAddr,
			KFFAddr:     si.KIFAddr,
		},
		UseGini: useGini,
		Gini:    si.Gini(minStake),
	}
}
