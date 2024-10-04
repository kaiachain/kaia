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

	"github.com/kaiachain/kaia/params"
)

var (
	big100 = big.NewInt(100)
)

type RewardConfig struct {
	rules params.Rules

	// Governance parameters
	mintingAmount *big.Int          // reward.mintingamount
	minimumStake  *big.Int          // reward.minimumstake
	deferredTxFee bool              // reward.deferredtxfee
	rewardRatio   *RewardRatio      // reward.ratio
	kip82Ratio    *RewardKip82Ratio // reward.kip82ratio
}

// TODO: use kaiax/gov
func NewRewardConfig(chainConfig *params.ChainConfig, num uint64) (*RewardConfig, error) {
	rc := &RewardConfig{}

	rc.rules = chainConfig.Rules(new(big.Int).SetUint64(num))

	rc.mintingAmount = new(big.Int).Set(chainConfig.Governance.Reward.MintingAmount)
	rc.minimumStake = new(big.Int).Set(chainConfig.Governance.Reward.MinimumStake)
	rc.deferredTxFee = chainConfig.Governance.Reward.DeferredTxFee

	if ratio, err := NewRewardRatio(chainConfig.Governance.Reward.Ratio); err != nil {
		return nil, err
	} else {
		rc.rewardRatio = ratio
	}

	if kip82Ratio, err := NewRewardKip82Ratio(chainConfig.Governance.Reward.Kip82Ratio); err != nil {
		return nil, err
	} else {
		rc.kip82Ratio = kip82Ratio
	}

	return rc, nil
}

// Parsed and validated reward.ratio parameter.
type RewardRatio struct {
	g int64 // Validators (GC)
	x int64 // Fund1 (KIF, KFF, KGF, PoC)
	y int64 // Fund2 (KEF, KCF, KIR)
}

func NewRewardRatio(ratio string) (*RewardRatio, error) {
	parts := strings.Split(ratio, "/")
	if len(parts) != 3 {
		return nil, errMalformedRewardRatio(ratio)
	}

	g, err1 := strconv.ParseInt(parts[0], 10, 64)
	x, err2 := strconv.ParseInt(parts[1], 10, 64)
	y, err3 := strconv.ParseInt(parts[2], 10, 64)
	if err1 != nil || err2 != nil || err3 != nil || g+x+y != 100 || g < 0 || x < 0 || y < 0 {
		return nil, errMalformedRewardRatio(ratio)
	}
	return &RewardRatio{g: g, x: x, y: y}, nil
}

// Split splits the amount into three parts according to the ratio.
func (r *RewardRatio) Split(amount *big.Int) (*big.Int, *big.Int, *big.Int) {
	gAmount := new(big.Int).Mul(amount, big.NewInt(r.g))
	gAmount = gAmount.Div(gAmount, big100)

	xAmount := new(big.Int).Mul(amount, big.NewInt(r.x))
	xAmount = xAmount.Div(xAmount, big100)

	yAmount := new(big.Int).Mul(amount, big.NewInt(r.y))
	yAmount = yAmount.Div(yAmount, big100)

	return gAmount, xAmount, yAmount
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
