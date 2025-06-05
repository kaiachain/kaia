// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package staking

import (
	"encoding/json"
	"io"
	"math"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/rlp"
)

// P2PStakingInfo contains staking information which is a wrapped version of StakingInfo.
type P2PStakingInfo struct {
	BlockNum uint64 `json:"blockNum"` // Block number where staking information of Council is fetched

	// Information retrieved from AddressBook smart contract
	CouncilNodeAddrs    []common.Address `json:"councilNodeAddrs"`    // NodeIds of Council
	CouncilStakingAddrs []common.Address `json:"councilStakingAddrs"` // Address of Staking account which holds staking balance
	CouncilRewardAddrs  []common.Address `json:"councilRewardAddrs"`  // Address of Council account which will get block reward

	KEFAddr common.Address `json:"kefAddr"` // Address of KEF contract
	KIFAddr common.Address `json:"kifAddr"` // Address of KIF contract

	UseGini bool    `json:"useGini"` // configure whether Gini is used or not
	Gini    float64 `json:"gini"`    // gini coefficient

	// Derived from CouncilStakingAddrs
	CouncilStakingAmounts []uint64 `json:"councilStakingAmounts"` // Staking amounts of Council
}

func FromStakingInfo(si *StakingInfo) *P2PStakingInfo {
	return &P2PStakingInfo{
		BlockNum:              si.SourceBlockNum,
		CouncilNodeAddrs:      si.NodeIds,
		CouncilStakingAddrs:   si.StakingContracts,
		CouncilRewardAddrs:    si.RewardAddrs,
		KEFAddr:               si.KEFAddr,
		KIFAddr:               si.KIFAddr,
		CouncilStakingAmounts: si.StakingAmounts,
	}
}

func FromStakingInfoWithGini(si *StakingInfo, useGini bool, minStake uint64) *P2PStakingInfo {
	return &P2PStakingInfo{
		BlockNum:              si.SourceBlockNum,
		CouncilNodeAddrs:      si.NodeIds,
		CouncilStakingAddrs:   si.StakingContracts,
		CouncilRewardAddrs:    si.RewardAddrs,
		KEFAddr:               si.KEFAddr,
		KIFAddr:               si.KIFAddr,
		CouncilStakingAmounts: si.StakingAmounts,
		UseGini:               useGini,
		Gini:                  si.Gini(minStake),
	}
}

func ToStakingInfo(si *P2PStakingInfo) *StakingInfo {
	return &StakingInfo{
		SourceBlockNum:   si.BlockNum,
		NodeIds:          si.CouncilNodeAddrs,
		StakingContracts: si.CouncilStakingAddrs,
		RewardAddrs:      si.CouncilRewardAddrs,
		KEFAddr:          si.KEFAddr,
		KIFAddr:          si.KIFAddr,
		StakingAmounts:   si.CouncilStakingAmounts,
	}
}

// MarshalJSON supports json marshalling for both P2PStakingInfo and StakingInfo
// TODO-Kaia-Mantle: remove this marshal function when backward-compatibility for KIR/PoC, KCF/KFF is not needed
func (st P2PStakingInfo) MarshalJSON() ([]byte, error) {
	type extendedSt struct {
		BlockNum              uint64           `json:"blockNum"`
		CouncilNodeAddrs      []common.Address `json:"councilNodeAddrs"`
		CouncilStakingAddrs   []common.Address `json:"councilStakingAddrs"`
		CouncilRewardAddrs    []common.Address `json:"councilRewardAddrs"`
		KEFAddr               common.Address   `json:"kefAddr"`
		KIFAddr               common.Address   `json:"kifAddr"`
		UseGini               bool             `json:"useGini"`
		Gini                  float64          `json:"gini"`
		CouncilStakingAmounts []uint64         `json:"councilStakingAmounts"`

		// legacy fields of StakingInfo
		KIRAddr common.Address `json:"KIRAddr"` // KIRAddr -> KCFAddr from v1.10.2
		PoCAddr common.Address `json:"PoCAddr"` // PoCAddr -> KFFAddr from v1.10.2
		KCFAddr common.Address `json:"kcfAddr"` // KCFAddr -> KEFAddr from Kaia v1.0.0
		KFFAddr common.Address `json:"kffAddr"` // KFFAddr -> KIFAddr from Kaia v1.0.0
	}

	var ext extendedSt
	ext.BlockNum = st.BlockNum
	ext.CouncilNodeAddrs = st.CouncilNodeAddrs
	ext.CouncilStakingAddrs = st.CouncilStakingAddrs
	ext.CouncilRewardAddrs = st.CouncilRewardAddrs
	ext.KEFAddr = st.KEFAddr
	ext.KIFAddr = st.KIFAddr
	ext.UseGini = st.UseGini
	ext.Gini = st.Gini
	ext.CouncilStakingAmounts = st.CouncilStakingAmounts

	// LHS are for backward-compatibility of database
	ext.KCFAddr = st.KEFAddr
	ext.KFFAddr = st.KIFAddr
	ext.KIRAddr = st.KEFAddr
	ext.PoCAddr = st.KIFAddr

	return json.Marshal(&ext)
}

// UnmarshalJSON supports json unmarshalling for both P2PStakingInfo and StakingInfo
func (st *P2PStakingInfo) UnmarshalJSON(input []byte) error {
	type extendedSt struct {
		BlockNum              uint64           `json:"blockNum"`
		CouncilNodeAddrs      []common.Address `json:"councilNodeAddrs"`
		CouncilStakingAddrs   []common.Address `json:"councilStakingAddrs"`
		CouncilRewardAddrs    []common.Address `json:"councilRewardAddrs"`
		KEFAddr               common.Address   `json:"kefAddr"`
		KIFAddr               common.Address   `json:"kifAddr"`
		UseGini               bool             `json:"useGini"`
		Gini                  float64          `json:"gini"`
		CouncilStakingAmounts []uint64         `json:"councilStakingAmounts"`

		// legacy fields of StakingInfo
		KIRAddr common.Address `json:"KIRAddr"` // KIRAddr -> KCFAddr from v1.10.2
		PoCAddr common.Address `json:"PoCAddr"` // PoCAddr -> KFFAddr from v1.10.2
		KCFAddr common.Address `json:"kcfAddr"` // KCFAddr -> KEFAddr from Kaia v1.0.0
		KFFAddr common.Address `json:"kffAddr"` // KFFAddr -> KIFAddr from Kaia v1.0.0
	}

	var ext extendedSt
	emptyAddr := common.Address{}

	if err := json.Unmarshal(input, &ext); err != nil {
		return err
	}

	st.BlockNum = ext.BlockNum
	st.CouncilNodeAddrs = ext.CouncilNodeAddrs
	st.CouncilStakingAddrs = ext.CouncilStakingAddrs
	st.CouncilRewardAddrs = ext.CouncilRewardAddrs
	st.KEFAddr = ext.KEFAddr
	st.KIFAddr = ext.KIFAddr
	st.UseGini = ext.UseGini
	st.Gini = ext.Gini
	st.CouncilStakingAmounts = ext.CouncilStakingAmounts

	if st.KEFAddr == emptyAddr {
		st.KEFAddr = ext.KCFAddr
		if st.KEFAddr == emptyAddr {
			st.KEFAddr = ext.KIRAddr
		}
	}
	if st.KIFAddr == emptyAddr {
		st.KIFAddr = ext.KFFAddr
		if st.KIFAddr == emptyAddr {
			st.KIFAddr = ext.PoCAddr
		}
	}

	return nil
}

type p2pStakingInfoRLP struct {
	BlockNum              uint64
	CouncilNodeAddrs      []common.Address
	CouncilStakingAddrs   []common.Address
	CouncilRewardAddrs    []common.Address
	KEFAddr               common.Address
	KIFAddr               common.Address
	UseGini               bool
	Gini                  uint64
	CouncilStakingAmounts []uint64
}

func (s *P2PStakingInfo) String() string {
	j, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (s *P2PStakingInfo) EncodeRLP(w io.Writer) error {
	// float64 is not rlp serializable, so it converts to bytes
	return rlp.Encode(w, &p2pStakingInfoRLP{s.BlockNum, s.CouncilNodeAddrs, s.CouncilStakingAddrs, s.CouncilRewardAddrs, s.KEFAddr, s.KIFAddr, s.UseGini, math.Float64bits(s.Gini), s.CouncilStakingAmounts})
}

func (s *P2PStakingInfo) DecodeRLP(st *rlp.Stream) error {
	var dec p2pStakingInfoRLP
	if err := st.Decode(&dec); err != nil {
		return err
	}
	s.BlockNum = dec.BlockNum
	s.CouncilNodeAddrs, s.CouncilStakingAddrs, s.CouncilRewardAddrs = dec.CouncilNodeAddrs, dec.CouncilStakingAddrs, dec.CouncilRewardAddrs
	s.KEFAddr, s.KIFAddr, s.UseGini, s.Gini = dec.KEFAddr, dec.KIFAddr, dec.UseGini, math.Float64frombits(dec.Gini)
	s.CouncilStakingAmounts = dec.CouncilStakingAmounts
	return nil
}
