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
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

var CalcDeferredRewardTimer time.Duration

var logger = log.NewModuleLogger(log.Reward)

var (
	errInvalidFormat = errors.New("invalid ratio format")
	errParsingRatio  = errors.New("parsing ratio fail")
)

type BalanceAdder interface {
	AddBalance(addr common.Address, v *big.Int)
}

// Cannot use governance.Engine because of cyclic dependency.
// Instead declare only the methods used by this package.
type governanceHelper interface {
	EffectiveParams(num uint64) (*params.GovParamSet, error)
}

type rewardConfig struct {
	// hardfork rules
	rules params.Rules

	// values calculated from block header and transactions
	// sum of actual fees paid. sum( tx.effectiveGasPrice * tx.gasUsed )
	// can be optimized before Kaia (baseFee * blockGasUsed), before Magma (unitPrice * blockGasUsed)
	totalFee *big.Int

	// values from GovParamSet
	mintingAmount *big.Int
	minimumStake  *big.Int
	deferredTxFee bool

	// parsed ratio
	cnRatio    *big.Int
	kifRatio   *big.Int
	kefRatio   *big.Int
	totalRatio *big.Int

	// parsed KIP82 ratio
	cnProposerRatio *big.Int
	cnStakingRatio  *big.Int
	cnTotalRatio    *big.Int
}

// Used in block reward distribution and klay_getReward RPC API
type RewardSpec struct {
	Minted   *big.Int                    `json:"minted"`   // the amount newly minted
	TotalFee *big.Int                    `json:"totalFee"` // total tx fee spent
	BurntFee *big.Int                    `json:"burntFee"` // the amount burnt
	Proposer *big.Int                    `json:"proposer"` // the amount allocated to the block proposer
	Stakers  *big.Int                    `json:"stakers"`  // total amount allocated to stakers
	KIF      *big.Int                    `json:"kif"`      // the amount allocated to KIF
	KEF      *big.Int                    `json:"kef"`      // the amount allocated to KEF
	Rewards  map[common.Address]*big.Int `json:"rewards"`  // mapping from reward recipient to amounts
}

func NewRewardSpec() *RewardSpec {
	return &RewardSpec{
		Minted:   big.NewInt(0),
		TotalFee: big.NewInt(0),
		BurntFee: big.NewInt(0),
		Proposer: big.NewInt(0),
		Stakers:  big.NewInt(0),
		KIF:      big.NewInt(0),
		KEF:      big.NewInt(0),
		Rewards:  make(map[common.Address]*big.Int),
	}
}

func (spec *RewardSpec) Add(delta *RewardSpec) {
	spec.Minted.Add(spec.Minted, delta.Minted)
	spec.TotalFee.Add(spec.TotalFee, delta.TotalFee)
	spec.BurntFee.Add(spec.BurntFee, delta.BurntFee)
	spec.Proposer.Add(spec.Proposer, delta.Proposer)
	spec.Stakers.Add(spec.Stakers, delta.Stakers)
	spec.KIF.Add(spec.KIF, delta.KIF)
	spec.KEF.Add(spec.KEF, delta.KEF)

	for addr, amount := range delta.Rewards {
		incrementRewardsMap(spec.Rewards, addr, amount)
	}
}

// Used in klay_totalSupply RPC API
// A minified version of RewardSpec that contains the total amount only.
type TotalReward struct {
	Minted   *big.Int
	BurntFee *big.Int
}

// TODO: this is for legacy, will be removed
type RewardDistributor struct{}

func NewRewardDistributor(gh governanceHelper) *RewardDistributor {
	return &RewardDistributor{}
}

// DistributeBlockReward distributes a given block's reward at the end of block processing
func DistributeBlockReward(b BalanceAdder, rewards map[common.Address]*big.Int) {
	for addr, amount := range rewards {
		b.AddBalance(addr, amount)
	}
}

func NewRewardConfig(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) (*rewardConfig, error) {
	cnRatio, kifRatio, kefRatio, totalRatio, err := parseRewardRatio(pset.Ratio())
	if err != nil {
		return nil, err
	}

	var cnProposerRatio, cnStakingRatio, cnTotalRatio int64
	if rules.IsKore {
		cnProposerRatio, cnStakingRatio, cnTotalRatio, err = parseRewardKip82Ratio(pset.Kip82Ratio())
		if err != nil {
			return nil, err
		}
	}

	return &rewardConfig{
		// hardfork rules
		rules: rules,

		// values calculated from block header
		totalFee: GetTotalTxFee(header, txs, receipts, rules, pset),

		// values from GovParamSet
		mintingAmount: new(big.Int).Set(pset.MintingAmountBig()),
		minimumStake:  new(big.Int).Set(pset.MinimumStakeBig()),
		deferredTxFee: pset.DeferredTxFee(),

		// parsed ratio
		cnRatio:    big.NewInt(cnRatio),
		kifRatio:   big.NewInt(kifRatio),
		kefRatio:   big.NewInt(kefRatio),
		totalRatio: big.NewInt(totalRatio),

		// parsed KIP82 ratio
		cnProposerRatio: big.NewInt(cnProposerRatio),
		cnStakingRatio:  big.NewInt(cnStakingRatio),
		cnTotalRatio:    big.NewInt(cnTotalRatio),
	}, nil
}

func GetTotalTxFee(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) *big.Int {
	totalFee := new(big.Int).SetUint64(header.GasUsed)
	if rules.IsMagma {
		totalFee = totalFee.Mul(totalFee, header.BaseFee)
	} else {
		totalFee = totalFee.Mul(totalFee, new(big.Int).SetUint64(pset.UnitPrice()))
	}

	if txs == nil || receipts == nil || !rules.IsKaia {
		return totalFee
	}

	if len(txs) != len(receipts) {
		logger.Error("GetTotalTxFee: txs and receipts length mismatch", "txs", len(txs), "receipts", len(receipts))
		return totalFee
	}
	// Since kaia, tip is added
	for i, tx := range txs {
		tip := new(big.Int).Mul(big.NewInt(int64(receipts[i].GasUsed)), tx.EffectiveGasTip(header.BaseFee))
		totalFee = totalFee.Add(totalFee, tip)
	}

	return totalFee
}

// config.Istanbul must have been set
func IsRewardSimple(pset *params.GovParamSet) bool {
	return pset.Policy() != uint64(istanbul.WeightedRandom)
}

// CalcRewardParamBlock returns the block number with which governance parameters must be fetched
// This mimics the legacy reward config cache before Kore
func CalcRewardParamBlock(num, epoch uint64, rules params.Rules) uint64 {
	if !rules.IsKore && num%epoch == 0 {
		return num - epoch
	}
	return num
}

// GetTotalReward returns the total rewards in this block, i.e. (minted - burntFee)
// Used in klay_totalSupply RPC API
func GetTotalReward(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) (*TotalReward, error) {
	total := new(TotalReward)
	rc, err := NewRewardConfig(header, txs, receipts, rules, pset)
	if err != nil {
		return nil, err
	}

	if IsRewardSimple(pset) {
		spec, err := CalcDeferredRewardSimple(header, txs, receipts, rules, pset)
		if err != nil {
			return nil, err
		}
		total.Minted = spec.Minted
		total.BurntFee = spec.BurntFee
	} else {
		// Lighter version of CalcDeferredReward where only minted and burntFee are calculated.
		// No need to lookup the staking info here.
		total.Minted = rc.mintingAmount
		_, _, burntFee := calcDeferredFee(rc)
		total.BurntFee = burntFee
	}

	// If not DeferredTxFee, fees are already added to the proposer during TX execution.
	// Therefore, there are no fees to distribute here at the end of block processing.
	// As such, the CalcDeferredRewardSimple() returns zero burntFee.
	// Here we calculate the fees burnt during the TX execution.
	if !rc.deferredTxFee && rules.IsMagma {
		txFee := GetTotalTxFee(header, txs, receipts, rules, pset)
		txFeeBurn := getBurnAmountMagma(txFee)
		total.BurntFee = total.BurntFee.Add(total.BurntFee, txFeeBurn)
	}

	return total, nil
}

// GetBlockReward returns the actual reward amounts paid in this block
// Used in kaia_getReward RPC API
func GetBlockReward(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	var spec *RewardSpec
	var err error

	if IsRewardSimple(pset) {
		spec, err = CalcDeferredRewardSimple(header, txs, receipts, rules, pset)
		if err != nil {
			return nil, err
		}
	} else {
		spec, err = CalcDeferredReward(header, txs, receipts, rules, pset)
		if err != nil {
			return nil, err
		}
	}

	// Compensate the difference between CalcDeferredReward() and actual payment.
	// If not DeferredTxFee, CalcDeferredReward() assumes 0 total_fee, but
	// some non-zero fee already has been paid to the proposer.
	if !pset.DeferredTxFee() {
		if rules.IsMagma {
			txFee := GetTotalTxFee(header, txs, receipts, rules, pset)
			txFeeBurn := getBurnAmountMagma(txFee)
			txFeeRemained := new(big.Int).Sub(txFee, txFeeBurn)
			spec.BurntFee = txFeeBurn

			spec.Proposer = spec.Proposer.Add(spec.Proposer, txFeeRemained)
			spec.TotalFee = spec.TotalFee.Add(spec.TotalFee, txFee)
			incrementRewardsMap(spec.Rewards, header.Rewardbase, txFeeRemained)
		} else {
			txFee := GetTotalTxFee(header, nil, nil, rules, pset)
			spec.Proposer = spec.Proposer.Add(spec.Proposer, txFee)
			spec.TotalFee = spec.TotalFee.Add(spec.TotalFee, txFee)
			// get the proposer of this block.
			proposer, err := ecrecover(header)
			if err != nil {
				return nil, err
			}
			incrementRewardsMap(spec.Rewards, proposer, txFee)
		}
	}

	return spec, nil
}

// CalcDeferredRewardSimple distributes rewards to proposer after optional fee burning
// this behaves similar to the previous MintKAIA
// MintKAIA has been superseded because we need to split reward distribution
// logic into (1) calculation, and (2) actual distribution.
// CalcDeferredRewardSimple does the former and DistributeBlockReward does the latter
func CalcDeferredRewardSimple(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	rc, err := NewRewardConfig(header, txs, receipts, rules, pset)
	if err != nil {
		return nil, err
	}

	minted := rc.mintingAmount

	// If not DeferredTxFee, fees are already added to the proposer during TX execution.
	// Therefore, there are no fees to distribute here at the end of block processing.
	// However, before Magma, there was a bug that distributed tx fee regardless
	// of `deferredTxFee` flag. See https://github.com/kaiachain/kaia/issues/1692.
	// To maintain backward compatibility, we only fix the buggy logic after Magma
	// and leave the buggy logic before Magma.
	// However, the fees must be compensated to calculate actual rewards paid.

	// bug-fixed logic after Magma
	if !rc.deferredTxFee && rc.rules.IsMagma {
		proposer := new(big.Int).Set(minted)
		logger.Debug("CalcDeferredRewardSimple after Magma when deferredTxFee=false returns",
			"proposer", proposer)
		spec := NewRewardSpec()
		spec.Minted = minted
		spec.TotalFee = big.NewInt(0)
		spec.BurntFee = big.NewInt(0)
		spec.Proposer = proposer
		incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)
		return spec, nil
	}

	totalFee := rc.totalFee
	rewardFee := new(big.Int).Set(totalFee)
	burntFee := big.NewInt(0)

	if rc.rules.IsMagma {
		burnt := getBurnAmountMagma(rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	proposer := big.NewInt(0).Add(minted, rewardFee)

	logger.Debug("CalcDeferredRewardSimple returns",
		"proposer", proposer.Uint64(),
		"totalFee", totalFee.Uint64(),
		"burntFee", burntFee.Uint64(),
	)

	spec := NewRewardSpec()
	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Proposer = proposer
	incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)
	return spec, nil
}

// CalcDeferredReward calculates the deferred rewards,
// which are determined at the end of block processing.
func CalcDeferredReward(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt, rules params.Rules, pset *params.GovParamSet) (*RewardSpec, error) {
	defer func(start time.Time) {
		CalcDeferredRewardTimer = time.Since(start)
	}(time.Now())

	rc, err := NewRewardConfig(header, txs, receipts, rules, pset)
	if err != nil {
		return nil, err
	}

	var (
		minted      = rc.mintingAmount
		stakingInfo = GetStakingInfo(header.Number.Uint64())
	)

	totalFee, rewardFee, burntFee := calcDeferredFee(rc)
	proposer, stakers, kif, kef, splitRem := calcSplit(rc, minted, rewardFee)
	shares, shareRem := calcShares(stakingInfo, stakers, rc.minimumStake.Uint64())

	// Remainder from (CN, KIF, KEF) split goes to KIF
	kif = kif.Add(kif, splitRem)
	// Remainder from staker shares goes to Proposer
	// Then, deduct it from stakers so that `minted + totalFee - burntFee = proposer + stakers + kif + kef`
	proposer = proposer.Add(proposer, shareRem)
	stakers = stakers.Sub(stakers, shareRem)

	// if KIF or KEF is not set, proposer gets the portion
	if stakingInfo == nil || common.EmptyAddress(stakingInfo.KIFAddr) {
		logger.Debug("KIF empty, proposer gets its portion", "kif", kif)
		proposer = proposer.Add(proposer, kif)
		kif = big.NewInt(0)
	}
	if stakingInfo == nil || common.EmptyAddress(stakingInfo.KEFAddr) {
		logger.Debug("KEF empty, proposer gets its portion", "kef", kef)
		proposer = proposer.Add(proposer, kef)
		kef = big.NewInt(0)
	}

	spec := NewRewardSpec()
	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Proposer = proposer
	spec.Stakers = stakers
	spec.KIF = kif
	spec.KEF = kef

	incrementRewardsMap(spec.Rewards, header.Rewardbase, proposer)

	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.KIFAddr) {
		incrementRewardsMap(spec.Rewards, stakingInfo.KIFAddr, kif)
	}
	if stakingInfo != nil && !common.EmptyAddress(stakingInfo.KEFAddr) {
		incrementRewardsMap(spec.Rewards, stakingInfo.KEFAddr, kef)
	}

	for rewardAddr, rewardAmount := range shares {
		incrementRewardsMap(spec.Rewards, rewardAddr, rewardAmount)
	}
	logger.Debug("CalcDeferredReward() returns", "spec", spec)

	return spec, nil
}

// calcDeferredFee splits fee into (total, reward, burnt)
func calcDeferredFee(rc *rewardConfig) (*big.Int, *big.Int, *big.Int) {
	// If not DeferredTxFee, fees are already added to the proposer during TX execution.
	// Therefore, there are no fees to distribute here at the end of block processing.
	// However, the fees must be compensated to calculate actual rewards paid.
	if !rc.deferredTxFee {
		return big.NewInt(0), big.NewInt(0), big.NewInt(0)
	}

	totalFee := rc.totalFee
	rewardFee := new(big.Int).Set(totalFee)
	burntFee := big.NewInt(0)

	// after magma, burn half of gas
	if rc.rules.IsMagma {
		burnt := getBurnAmountMagma(rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	// after kore, burn fees up to proposer's minted reward
	if rc.rules.IsKore {
		burnt := getBurnAmountKore(rc, rewardFee)
		rewardFee = rewardFee.Sub(rewardFee, burnt)
		burntFee = burntFee.Add(burntFee, burnt)
	}

	logger.Debug("calcDeferredFee()",
		"totalFee", totalFee.Uint64(),
		"rewardFee", rewardFee.Uint64(),
		"burntFee", burntFee.Uint64(),
	)
	return totalFee, rewardFee, burntFee
}

func getBurnAmountMagma(fee *big.Int) *big.Int {
	return new(big.Int).Div(fee, big.NewInt(2))
}

func getBurnAmountKore(rc *rewardConfig, fee *big.Int) *big.Int {
	cn, _, _ := splitByRatio(rc, rc.mintingAmount)
	proposer, _ := splitByKip82Ratio(rc, cn)

	logger.Debug("getBurnAmountKore()",
		"fee", fee.Uint64(),
		"proposer", proposer.Uint64(),
	)

	if fee.Cmp(proposer) >= 0 {
		return proposer
	} else {
		return new(big.Int).Set(fee) // return copy of the parameter
	}
}

// calcSplit splits fee into (proposer, stakers, kif, kef, remaining)
// the sum of the output must be equal to (minted + fee)
func calcSplit(rc *rewardConfig, minted, fee *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int) {
	totalResource := big.NewInt(0)
	totalResource = totalResource.Add(minted, fee)

	if rc.rules.IsKore {
		cn, kif, kef := splitByRatio(rc, minted)
		proposer, stakers := splitByKip82Ratio(rc, cn)

		proposer = proposer.Add(proposer, fee)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kif)
		remaining = remaining.Sub(remaining, kef)
		remaining = remaining.Sub(remaining, proposer)
		remaining = remaining.Sub(remaining, stakers)

		logger.Debug("calcSplit() after kore",
			"[in] minted", minted.Uint64(),
			"[in] fee", fee.Uint64(),
			"[out] proposer", proposer.Uint64(),
			"[out] stakers", stakers.Uint64(),
			"[out] kif", kif.Uint64(),
			"[out] kef", kef.Uint64(),
			"[out] remaining", remaining.Uint64(),
		)
		return proposer, stakers, kif, kef, remaining
	} else {
		cn, kif, kef := splitByRatio(rc, totalResource)

		remaining := new(big.Int).Set(totalResource)
		remaining = remaining.Sub(remaining, kif)
		remaining = remaining.Sub(remaining, kef)
		remaining = remaining.Sub(remaining, cn)

		logger.Debug("calcSplit() before kore",
			"[in] minted", minted.Uint64(),
			"[in] fee", fee.Uint64(),
			"[out] cn", cn.Uint64(),
			"[out] kif", kif.Uint64(),
			"[out] kef", kef.Uint64(),
			"[out] remaining", remaining.Uint64(),
		)
		return cn, big.NewInt(0), kif, kef, remaining
	}
}

// splitByRatio splits by `ratio`. It ignores any remaining amounts.
func splitByRatio(rc *rewardConfig, source *big.Int) (*big.Int, *big.Int, *big.Int) {
	cn := new(big.Int).Mul(source, rc.cnRatio)
	cn = cn.Div(cn, rc.totalRatio)

	kif := new(big.Int).Mul(source, rc.kifRatio)
	kif = kif.Div(kif, rc.totalRatio)

	kef := new(big.Int).Mul(source, rc.kefRatio)
	kef = kef.Div(kef, rc.totalRatio)

	return cn, kif, kef
}

// splitByKip82Ratio splits by `kip82ratio`. It ignores any remaining amounts.
func splitByKip82Ratio(rc *rewardConfig, source *big.Int) (*big.Int, *big.Int) {
	proposer := new(big.Int).Mul(source, rc.cnProposerRatio)
	proposer = proposer.Div(proposer, rc.cnTotalRatio)

	stakers := new(big.Int).Mul(source, rc.cnStakingRatio)
	stakers = stakers.Div(stakers, rc.cnTotalRatio)

	return proposer, stakers
}

// calcShares distributes stake reward among staked CNs
func calcShares(stakingInfo *StakingInfo, stakeReward *big.Int, minStake uint64) (map[common.Address]*big.Int, *big.Int) {
	// if stakingInfo is nil, stakeReward goes to proposer
	if stakingInfo == nil {
		return make(map[common.Address]*big.Int), stakeReward
	}

	cns := stakingInfo.GetConsolidatedStakingInfo()

	totalStakesInt := uint64(0)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStake { // comparison in KAIA
			totalStakesInt += (node.StakingAmount - minStake)
		}
	}

	totalStakes := new(big.Int).SetUint64(totalStakesInt)
	remaining := new(big.Int).Set(stakeReward)
	shares := make(map[common.Address]*big.Int)

	for _, node := range cns.GetAllNodes() {
		if node.StakingAmount > minStake {
			effectiveStake := new(big.Int).SetUint64(node.StakingAmount - minStake)
			// The KAIA unit will cancel out:
			// rewardAmount (kei) = stakeReward (kei) * effectiveStake (KAIA) / totalStakes (KAIA)
			rewardAmount := new(big.Int).Mul(stakeReward, effectiveStake)
			rewardAmount = rewardAmount.Div(rewardAmount, totalStakes)
			remaining = remaining.Sub(remaining, rewardAmount)
			if rewardAmount.Cmp(big.NewInt(0)) > 0 {
				shares[node.RewardAddr] = rewardAmount
			}
		}
	}
	logger.Debug("calcShares()",
		"[in] stakeReward", stakeReward.Uint64(),
		"[out] remaining", remaining.Uint64(),
		"[out] shares", shares,
	)

	return shares, remaining
}

// parseRewardRatio parses string `ratio` into ints
func parseRewardRatio(ratio string) (int64, int64, int64, int64, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardSliceCount {
		logger.Error("Invalid ratio format", "ratio", ratio)
		return 0, 0, 0, 0, errInvalidFormat
	}
	cn, err1 := strconv.ParseInt(s[0], 10, 64)
	kif, err2 := strconv.ParseInt(s[1], 10, 64)
	kef, err3 := strconv.ParseInt(s[2], 10, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		logger.Error("Could not parse ratio", "ratio", ratio)
		return 0, 0, 0, 0, errParsingRatio
	}
	return cn, kif, kef, cn + kif + kef, nil
}

// parseRewardKip82Ratio parses string `kip82ratio` into ints
func parseRewardKip82Ratio(ratio string) (int64, int64, int64, error) {
	s := strings.Split(ratio, "/")
	if len(s) != params.RewardKip82SliceCount {
		logger.Error("Invalid kip82ratio format", "ratio", ratio)
		return 0, 0, 0, errInvalidFormat
	}
	proposer, err1 := strconv.ParseInt(s[0], 10, 64)
	stakers, err2 := strconv.ParseInt(s[1], 10, 64)

	if err1 != nil || err2 != nil {
		logger.Error("Could not parse kip82ratio", "ratio", ratio)
		return 0, 0, 0, errParsingRatio
	}
	return proposer, stakers, proposer + stakers, nil
}

func incrementRewardsMap(m map[common.Address]*big.Int, addr common.Address, amount *big.Int) {
	_, ok := m[addr]
	if !ok {
		m[addr] = big.NewInt(0)
	}

	m[addr] = m[addr].Add(m[addr], amount)
}

// ecrecover extracts the Kaia account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, nil
	}

	sigHash, err := sigHash(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}
	return addr, nil
}

func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	if err := rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false)); err != nil {
		logger.Error("fail to encode", "err", err)
		return common.Hash{}, err
	}
	hasher.Sum(hash[:0])
	return hash, nil
}
