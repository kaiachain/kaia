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

package reward

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const (
	AddrNotFoundInCouncilNodes = -1
	maxStakingLimit            = uint64(100000000000)
	DefaultGiniCoefficient     = -1.0
)

var (
	maxStakingLimitBigInt = big.NewInt(0).SetUint64(maxStakingLimit)

	ErrAddrNotInStakingInfo = errors.New("Address is not in stakingInfo")
)

// StakingInfo contains staking information.
type StakingInfo struct {
	BlockNum uint64 `json:"blockNum"` // Block number where staking information of Council is fetched

	// Information retrieved from AddressBook smart contract
	CouncilNodeAddrs    []common.Address `json:"councilNodeAddrs"`    // NodeIds of Council
	CouncilStakingAddrs []common.Address `json:"councilStakingAddrs"` // Address of Staking account which holds staking balance
	CouncilRewardAddrs  []common.Address `json:"councilRewardAddrs"`  // Address of Council account which will get block reward

	KEFAddr common.Address `json:"kefAddr"` // Address of KEF contract
	KIFAddr common.Address `json:"kifAddr"` // Address of KIF contract

	UseGini bool    `json:"useGini"`
	Gini    float64 `json:"gini"` // gini coefficient

	// Derived from CouncilStakingAddrs
	CouncilStakingAmounts []uint64 `json:"councilStakingAmounts"` // Staking amounts of Council
}

func FromKaiax(si *staking.StakingInfo) *StakingInfo {
	return &StakingInfo{
		BlockNum:              si.SourceBlockNum,
		CouncilNodeAddrs:      si.NodeIds,
		CouncilStakingAddrs:   si.StakingContracts,
		CouncilRewardAddrs:    si.RewardAddrs,
		KEFAddr:               si.KEFAddr,
		KIFAddr:               si.KIFAddr,
		CouncilStakingAmounts: si.StakingAmounts,
	}
}

func FromKaiaxWithGini(si *staking.StakingInfo, useGini bool, minStake uint64) *StakingInfo {
	return &StakingInfo{
		BlockNum:              si.SourceBlockNum,
		CouncilNodeAddrs:      si.NodeIds,
		CouncilStakingAddrs:   si.StakingContracts,
		CouncilRewardAddrs:    si.RewardAddrs,
		KEFAddr:               si.KEFAddr,
		KIFAddr:               si.KIFAddr,
		CouncilStakingAmounts: si.StakingAmounts,
	}
}

// MarshalJSON supports json marshalling for both oldStakingInfo and StakingInfo
// TODO-Kaia-Mantle: remove this marshal function when backward-compatibility for KIR/PoC, KCF/KFF is not needed
func (st StakingInfo) MarshalJSON() ([]byte, error) {
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

// UnmarshalJSON supports json unmarshalling for both oldStakingInfo and StakingInfo
func (st *StakingInfo) UnmarshalJSON(input []byte) error {
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

// Refined staking information suitable for proposer selection.
// Sometimes a node would register multiple NodeAddrs
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
	NodeAddrs     []common.Address
	StakingAddrs  []common.Address
	RewardAddr    common.Address // common reward address
	StakingAmount uint64         // sum of staking amounts
}

type ConsolidatedStakingInfo struct {
	nodes     []consolidatedNode
	nodeIndex map[common.Address]int // nodeAddr -> index in []nodes
}

type stakingInfoRLP struct {
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

func newEmptyStakingInfo(blockNum uint64) *StakingInfo {
	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeAddrs:      make([]common.Address, 0, 0),
		CouncilStakingAddrs:   make([]common.Address, 0, 0),
		CouncilRewardAddrs:    make([]common.Address, 0, 0),
		KEFAddr:               common.Address{},
		KIFAddr:               common.Address{},
		CouncilStakingAmounts: make([]uint64, 0, 0),
		Gini:                  DefaultGiniCoefficient,
		UseGini:               false,
	}
	return stakingInfo
}

// newStakingInfo fills stakingAmount from effectiveStakings if set. Otherwise, stakingAmount is filled with the balances of staking contracts at a given block.
// Currently, effectiveStakings are same as `statedb.GetBalance(stakingAddr)`, which is the balance of staking contracts.
func newStakingInfo(bc blockChain, helper governanceHelper, blockNum uint64, types []uint8, addrs []common.Address, effectiveStakings ...*big.Int) (*StakingInfo, error) {
	var (
		nodeIds      = []common.Address{}
		stakingAddrs = []common.Address{}
		rewardAddrs  = []common.Address{}
		pocAddr      = common.Address{}
		kirAddr      = common.Address{}
	)

	// Parse and construct node information
	for i, addrType := range types {
		switch addrType {
		case addressTypeNodeID:
			nodeIds = append(nodeIds, addrs[i])
		case addressTypeStakingAddr:
			stakingAddrs = append(stakingAddrs, addrs[i])
		case addressTypeRewardAddr:
			rewardAddrs = append(rewardAddrs, addrs[i])
		case addressTypePoCAddr:
			pocAddr = addrs[i]
		case addressTypeKIRAddr:
			kirAddr = addrs[i]
		default:
			return nil, fmt.Errorf("invalid type from AddressBook: %d", addrType)
		}
	}

	// validate parsed node information
	if len(nodeIds) != len(stakingAddrs) ||
		len(nodeIds) != len(rewardAddrs) ||
		len(effectiveStakings) > 0 && len(stakingAddrs) != len(effectiveStakings) ||
		common.EmptyAddress(pocAddr) ||
		common.EmptyAddress(kirAddr) {
		// This is an expected behavior when the addressBook contract is not activated yet.
		logger.Info("The addressBook is not yet activated. Use empty stakingInfo")
		return newEmptyStakingInfo(blockNum), nil
	}

	// Get balance of stakingAddrs
	stakingAmounts := make([]uint64, len(stakingAddrs))
	if len(effectiveStakings) == 0 {
		intervalBlock := bc.GetBlockByNumber(blockNum)
		if intervalBlock == nil {
			logger.Trace("Failed to get the block by the given number", "blockNum", blockNum)
			return nil, errors.New(fmt.Sprintf("Failed to get the block by the given number. blockNum: %d", blockNum))
		}
		statedb, err := bc.StateAt(intervalBlock.Root())
		if err != nil {
			logger.Trace("Failed to make a state for interval block", "interval blockNum", blockNum, "err", err)
			return nil, err
		}
		for i, stakingAddr := range stakingAddrs {
			tempStakingAmount := big.NewInt(0).Div(statedb.GetBalance(stakingAddr), big.NewInt(0).SetUint64(params.KAIA))
			if tempStakingAmount.Cmp(maxStakingLimitBigInt) > 0 {
				tempStakingAmount.SetUint64(maxStakingLimit)
			}
			stakingAmounts[i] = tempStakingAmount.Uint64()
		}
	} else {
		for i, effectiveStaking := range effectiveStakings {
			tempStakingAmount := big.NewInt(0).Div(effectiveStaking, big.NewInt(0).SetUint64(params.KAIA))
			if tempStakingAmount.Cmp(maxStakingLimitBigInt) > 0 {
				tempStakingAmount.SetUint64(maxStakingLimit)
			}
			stakingAmounts[i] = tempStakingAmount.Uint64()
		}
	}

	pset, err := helper.EffectiveParams(blockNum)
	if err != nil {
		return nil, err
	}
	useGini := pset.UseGiniCoeff()
	gini := DefaultGiniCoefficient

	stakingInfo := &StakingInfo{
		BlockNum:              blockNum,
		CouncilNodeAddrs:      nodeIds,
		CouncilStakingAddrs:   stakingAddrs,
		CouncilRewardAddrs:    rewardAddrs,
		KEFAddr:               kirAddr,
		KIFAddr:               pocAddr,
		CouncilStakingAmounts: stakingAmounts,
		Gini:                  gini,
		UseGini:               useGini,
	}
	return stakingInfo, nil
}

func (s *StakingInfo) GetIndexByNodeAddress(nodeAddress common.Address) (int, error) {
	for i, addr := range s.CouncilNodeAddrs {
		if addr == nodeAddress {
			return i, nil
		}
	}
	return AddrNotFoundInCouncilNodes, ErrAddrNotInStakingInfo
}

func (s *StakingInfo) GetStakingAmountByNodeId(nodeAddress common.Address) (uint64, error) {
	i, err := s.GetIndexByNodeAddress(nodeAddress)
	if err != nil {
		return 0, err
	}
	return s.CouncilStakingAmounts[i], nil
}

func (s *StakingInfo) String() string {
	j, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (s *StakingInfo) EncodeRLP(w io.Writer) error {
	// float64 is not rlp serializable, so it converts to bytes
	return rlp.Encode(w, &stakingInfoRLP{s.BlockNum, s.CouncilNodeAddrs, s.CouncilStakingAddrs, s.CouncilRewardAddrs, s.KEFAddr, s.KIFAddr, s.UseGini, math.Float64bits(s.Gini), s.CouncilStakingAmounts})
}

func (s *StakingInfo) DecodeRLP(st *rlp.Stream) error {
	var dec stakingInfoRLP
	if err := st.Decode(&dec); err != nil {
		return err
	}
	s.BlockNum = dec.BlockNum
	s.CouncilNodeAddrs, s.CouncilStakingAddrs, s.CouncilRewardAddrs = dec.CouncilNodeAddrs, dec.CouncilStakingAddrs, dec.CouncilRewardAddrs
	s.KEFAddr, s.KIFAddr, s.UseGini, s.Gini = dec.KEFAddr, dec.KIFAddr, dec.UseGini, math.Float64frombits(dec.Gini)
	s.CouncilStakingAmounts = dec.CouncilStakingAmounts
	return nil
}

func (s *StakingInfo) GetConsolidatedStakingInfo() *ConsolidatedStakingInfo {
	c := &ConsolidatedStakingInfo{
		nodes:     make([]consolidatedNode, 0),
		nodeIndex: make(map[common.Address]int),
	}

	rewardIndex := make(map[common.Address]int) // temporarily map rewardAddr -> index in []nodes

	for j := 0; j < len(s.CouncilNodeAddrs); j++ {
		var (
			nodeAddr      = s.CouncilNodeAddrs[j]
			stakingAddr   = s.CouncilStakingAddrs[j]
			rewardAddr    = s.CouncilRewardAddrs[j]
			stakingAmount = s.CouncilStakingAmounts[j]
		)
		if idx, ok := rewardIndex[rewardAddr]; !ok {
			c.nodes = append(c.nodes, consolidatedNode{
				NodeAddrs:     []common.Address{nodeAddr},
				StakingAddrs:  []common.Address{stakingAddr},
				RewardAddr:    rewardAddr,
				StakingAmount: stakingAmount,
			})
			c.nodeIndex[nodeAddr] = len(c.nodes) - 1 // point to new element
			rewardIndex[rewardAddr] = len(c.nodes) - 1
		} else {
			c.nodes[idx].NodeAddrs = append(c.nodes[idx].NodeAddrs, nodeAddr)
			c.nodes[idx].StakingAddrs = append(c.nodes[idx].StakingAddrs, stakingAddr)
			c.nodes[idx].StakingAmount += stakingAmount
			c.nodeIndex[nodeAddr] = idx // point to existing element
		}
	}
	return c
}

func (c *ConsolidatedStakingInfo) GetAllNodes() []consolidatedNode {
	return c.nodes
}

func (c *ConsolidatedStakingInfo) GetConsolidatedNode(nodeAddr common.Address) *consolidatedNode {
	if idx, ok := c.nodeIndex[nodeAddr]; ok {
		return &c.nodes[idx]
	}
	return nil
}

// Calculate Gini coefficient of the StakingAmounts.
// Only amounts greater or equal to `minStake` are included in the calculation.
// Set `minStake` to 0 to calculate Gini coefficient of all amounts.
func (c *ConsolidatedStakingInfo) CalcGiniCoefficientMinStake(minStake uint64) float64 {
	var amounts []float64
	for _, node := range c.nodes {
		if node.StakingAmount >= minStake {
			amounts = append(amounts, float64(node.StakingAmount))
		}
	}

	if len(amounts) == 0 {
		return DefaultGiniCoefficient
	}
	return CalcGiniCoefficient(amounts)
}

func (c *ConsolidatedStakingInfo) String() string {
	j, err := json.Marshal(c.nodes)
	if err != nil {
		return err.Error()
	}
	return string(j)
}

type float64Slice []float64

func (p float64Slice) Len() int           { return len(p) }
func (p float64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func CalcGiniCoefficient(stakingAmount float64Slice) float64 {
	sort.Sort(stakingAmount)

	// calculate gini coefficient
	sumOfAbsoluteDifferences := float64(0)
	subSum := float64(0)

	for i, x := range stakingAmount {
		temp := x*float64(i) - subSum
		sumOfAbsoluteDifferences = sumOfAbsoluteDifferences + temp
		subSum = subSum + x
	}

	result := sumOfAbsoluteDifferences / subSum / float64(len(stakingAmount))
	result = math.Round(result*100) / 100

	return result
}
