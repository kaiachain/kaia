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

package reward

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
)

var big100 = big.NewInt(100)

type RewardConfig struct {
	Rules      params.Rules
	Rewardbase common.Address // Proposer's reward recipient address

	IsSimple               bool              // istanbul.policy != WeightedRandom in which case simple rules are used
	UnitPrice              *big.Int          // governance.unitprice
	MintingAmount          *big.Int          // reward.mintingamount
	MinimumStake           *big.Int          // reward.minimumstake
	DeferredTxFee          bool              // reward.deferredtxfee
	RewardRatio            *RewardRatio      // reward.ratio
	Kip82Ratio             *RewardKip82Ratio // reward.kip82ratio
	StakingRewardThreshold *big.Int          // reward.stakingrewardthreshold
	UseFlexReward          bool              // reward.useflexreward
}

// TODO-kaiax: Restore to gov.GovModule after introducing kaiax/gov
type GovModule interface {
	GetParamSet(blockNum uint64) gov.ParamSet
}

func NewRewardConfig(chainConfig *params.ChainConfig, govModule GovModule, header *types.Header) (*RewardConfig, error) {
	rc := &RewardConfig{}

	rc.Rules = chainConfig.Rules(header.Number)
	rc.Rewardbase = header.Rewardbase

	paramset := govModule.GetParamSet(header.Number.Uint64())
	rc.IsSimple = paramset.ProposerPolicy != uint64(istanbul.WeightedRandom)
	rc.UnitPrice = new(big.Int).SetUint64(paramset.UnitPrice)
	rc.MintingAmount = new(big.Int).Set(paramset.MintingAmount)
	rc.MinimumStake = new(big.Int).Set(paramset.MinimumStake)
	rc.DeferredTxFee = paramset.DeferredTxFee
	rc.StakingRewardThreshold = new(big.Int).Set(paramset.StakingRewardThreshold)
	rc.UseFlexReward = paramset.UseFlexReward

	if ratio, err := NewRewardRatio(paramset.Ratio); err != nil {
		return nil, err
	} else {
		rc.RewardRatio = ratio
	}

	if kip82Ratio, err := NewRewardKip82Ratio(paramset.Kip82Ratio); err != nil {
		return nil, err
	} else {
		rc.Kip82Ratio = kip82Ratio
	}

	return rc, nil
}

// Parsed and validated reward.ratio parameter.
// Supports both 3-part (g/x/y) and 4-part (g/x/y/z) formats.
type RewardRatio struct {
	g int64 // Validators (GC)
	x int64 // Fund1 (KIF, KFF, KGF, PoC)
	y int64 // Fund2 (KEF, KCF, KIR)
	z int64 // Fund3 (KPF) - optional, used with flexreward.
}

func NewRewardRatio(ratio string) (*RewardRatio, error) {
	parts := strings.Split(ratio, "/")
	if len(parts) != 3 && len(parts) != 4 {
		return nil, errMalformedRewardRatio(ratio)
	}

	g, err1 := strconv.ParseInt(parts[0], 10, 64)
	x, err2 := strconv.ParseInt(parts[1], 10, 64)
	y, err3 := strconv.ParseInt(parts[2], 10, 64)
	z := int64(0)
	var err4 error
	if len(parts) == 4 {
		z, err4 = strconv.ParseInt(parts[3], 10, 64)
	}
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil || g+x+y+z != 100 || g < 0 || x < 0 || y < 0 || z < 0 {
		return nil, errMalformedRewardRatio(ratio)
	}
	return &RewardRatio{g: g, x: x, y: y, z: z}, nil
}

// Split splits the amount into three parts (g, x, y) according to the ratio.
// The z (KPF) portion is not included; use SplitFlex() for 4-part splits.
func (r *RewardRatio) Split(amount *big.Int) (*big.Int, *big.Int, *big.Int) {
	gAmount := new(big.Int).Mul(amount, big.NewInt(r.g))
	gAmount = gAmount.Div(gAmount, big100)

	xAmount := new(big.Int).Mul(amount, big.NewInt(r.x))
	xAmount = xAmount.Div(xAmount, big100)

	yAmount := new(big.Int).Mul(amount, big.NewInt(r.y))
	yAmount = yAmount.Div(yAmount, big100)

	return gAmount, xAmount, yAmount
}

// SplitFlex splits the amount into four parts (g, x, y, z) according to the ratio.
// For 3-part ratios, z is zero.
func (r *RewardRatio) SplitFlex(amount *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	gAmount := new(big.Int).Mul(amount, big.NewInt(r.g))
	gAmount = gAmount.Div(gAmount, big100)

	xAmount := new(big.Int).Mul(amount, big.NewInt(r.x))
	xAmount = xAmount.Div(xAmount, big100)

	yAmount := new(big.Int).Mul(amount, big.NewInt(r.y))
	yAmount = yAmount.Div(yAmount, big100)

	zAmount := new(big.Int).Mul(amount, big.NewInt(r.z))
	zAmount = zAmount.Div(zAmount, big100)

	return gAmount, xAmount, yAmount, zAmount
}

// Parsed and validated reward.kip82ratio parameter.
type RewardKip82Ratio struct {
	p int64 // Proposer
	s int64 // Stakers
}

func NewRewardKip82Ratio(ratio string) (*RewardKip82Ratio, error) {
	parts := strings.Split(ratio, "/")
	if len(parts) != 2 {
		return nil, errMalformedRewardKip82Ratio(ratio)
	}

	p, err1 := strconv.ParseInt(parts[0], 10, 64)
	s, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil || p+s != 100 || p < 0 || s < 0 {
		return nil, errMalformedRewardKip82Ratio(ratio)
	}
	return &RewardKip82Ratio{p: p, s: s}, nil
}

// Split splits the amount into two parts according to the ratio.
func (r *RewardKip82Ratio) Split(amount *big.Int) (*big.Int, *big.Int) {
	pAmount := new(big.Int).Mul(amount, big.NewInt(r.p))
	pAmount = pAmount.Div(pAmount, big100)

	sAmount := new(big.Int).Mul(amount, big.NewInt(r.s))
	sAmount = sAmount.Div(sAmount, big100)

	return pAmount, sAmount
}
