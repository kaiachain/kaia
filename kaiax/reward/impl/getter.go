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

package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/kaiax/reward"
	"github.com/kaiachain/kaia/kaiax/staking"
)

// Below outlines the relationship between the getters and their helper functions.
//
// GetRewardSummary
// - loadBlockData
// - - getTotalFee = F
// - getRewardSummary = MR + DF + NDF
//
// GetBlockReward
// - loadBlockData
// - - getTotalFee = F
// - getDeferredReward = MR + DF
// - specWithNonDeferredFee = MR + DF + NDF
//
// GetDeferredReward
// - getTotalFee = F
// - getDeferredReward = MR + DF
// - - getDeferredRewardSimple   for istanbul.policy != 2
// - - getDeferredRewardFull     for istanbul.policy == 2
// - - - getDeferredRewardFullKore    for IsKore
// - - - getDeferredRewardFullLegacy  for !IsKore

func (r *RewardModule) GetRewardSummary(num uint64) (*reward.RewardSummary, error) {
	config, _, totalFee, err := r.loadBlockData(num)
	if err != nil {
		return nil, err
	}
	return getRewardSummary(config, totalFee), nil
}

func (r *RewardModule) loadBlockData(num uint64) (*reward.RewardConfig, *types.Header, *big.Int, error) {
	block := r.Chain.GetBlockByNumber(num)
	if block == nil {
		return nil, nil, nil, reward.ErrNoBlock
	}
	receipts := r.Chain.GetReceiptsByBlockHash(block.Hash())
	if receipts == nil {
		return nil, nil, nil, reward.ErrNoReceipts
	}
	header := block.Header()
	txs := block.Transactions()

	config, err := reward.NewRewardConfig(r.ChainConfig, r.GovModule, header)
	if err != nil {
		return nil, nil, nil, err
	}
	totalFee, err := getTotalFee(config, header, txs, receipts)
	if err != nil {
		return nil, nil, nil, err
	}

	return config, header, totalFee, nil
}

func getRewardSummary(config *reward.RewardConfig, totalFee *big.Int) *reward.RewardSummary {
	minted := new(big.Int).Set(config.MintingAmount)

	burntFee := big.NewInt(0)
	if config.IsSimple { // simplified getDeferredRewardSimple
		if config.Rules.IsMagma {
			burntFee = getBurnAmountMagma(totalFee)
		}
	} else { // simplified getDeferredRewardFull
		if config.Rules.IsKore {
			burntFee = getBurnAmountKore(config, totalFee)
		} else if config.Rules.IsMagma {
			burntFee = getBurnAmountMagma(totalFee)
		}
	}

	summary := reward.NewRewardSummary()
	summary.Minted = minted
	summary.TotalFee = totalFee
	summary.BurntFee = burntFee
	return summary
}

func (r *RewardModule) GetBlockReward(num uint64) (*reward.RewardSpec, error) {
	config, header, totalFee, err := r.loadBlockData(num)
	if err != nil {
		return nil, err
	}

	spec, err := r.getDeferredReward(config, header, totalFee)
	if err != nil {
		return nil, err
	}

	return r.specWithNonDeferredFee(spec, config, header, totalFee)
}

// specWithNonDeferredFee adds non-deferred fees to the reward spec.
func (r *RewardModule) specWithNonDeferredFee(spec *reward.RewardSpec, config *reward.RewardConfig, header *types.Header, totalFee *big.Int) (*reward.RewardSpec, error) {
	if config.DeferredTxFee {
		return spec, nil // nothing to do under deferred mode
	}

	newSpec := spec.Copy()

	if config.Rules.IsMagma {
		burntFee := getBurnAmountMagma(totalFee)
		distributedFee := new(big.Int).Sub(totalFee, burntFee)

		newSpec.TotalFee.Add(newSpec.TotalFee, totalFee)
		newSpec.BurntFee.Add(newSpec.BurntFee, burntFee)
		newSpec.Proposer.Add(newSpec.Proposer, distributedFee)

		// Since Magma, non-deferred fees are assigned to header.Rewardbase.
		newSpec.IncReceipient(config.Rewardbase, distributedFee)
	} else {
		distributedFee := new(big.Int).Set(totalFee)

		newSpec.TotalFee.Add(newSpec.TotalFee, totalFee)
		newSpec.Proposer.Add(newSpec.Proposer, distributedFee)

		// Before Magma, non-deferred fees are assigned to evm.Coinbase which originates from Engine().Author(header).
		coinbase, err := r.Chain.Engine().Author(header)
		if err != nil {
			return nil, err
		}
		newSpec.IncReceipient(coinbase, distributedFee)
	}

	return newSpec, nil
}

// Under non-deferred mode, transaction fees are ignored.
func (r *RewardModule) GetDeferredReward(header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) (*reward.RewardSpec, error) {
	config, err := reward.NewRewardConfig(r.ChainConfig, r.GovModule, header)
	if err != nil {
		return nil, err
	}
	totalFee, err := getTotalFee(config, header, txs, receipts)
	if err != nil {
		return nil, err
	}

	return r.getDeferredReward(config, header, totalFee)
}

func (r *RewardModule) getDeferredReward(config *reward.RewardConfig, header *types.Header, totalFee *big.Int) (*reward.RewardSpec, error) {
	if config.IsSimple {
		return getDeferredRewardSimple(config, totalFee)
	} else {
		si, err := r.StakingModule.GetStakingInfo(header.Number.Uint64())
		if err != nil {
			return nil, err
		}
		return getDeferredRewardFull(config, totalFee, si)
	}
}

// getTotalFee calculates the total transaction fees in the block.
func getTotalFee(config *reward.RewardConfig, header *types.Header, txs []*types.Transaction, receipts []*types.Receipt) (*big.Int, error) {
	if config.Rules.IsKaia {
		// sum { tx[i].gasUsed * tx[i].effectiveGasPrice }
		// = block.gasUsed * block.baseFeePerGas + sum { tx[i].gasUsed * tx[i].effectiveGasTip }
		if len(txs) != len(receipts) {
			return nil, reward.ErrTxReceiptsLenMismatch
		}
		totalFee := new(big.Int).Mul(big.NewInt(int64(header.GasUsed)), header.BaseFee)
		for i, tx := range txs {
			tip := new(big.Int).Mul(big.NewInt(int64(receipts[i].GasUsed)), tx.EffectiveGasTip(header.BaseFee))
			totalFee = totalFee.Add(totalFee, tip)
		}
		return totalFee, nil
	} else if config.Rules.IsMagma {
		// Optimized to block.gasUsed * block.baseFeePerGas
		return new(big.Int).Mul(big.NewInt(int64(header.GasUsed)), header.BaseFee), nil
	} else {
		// Optimized to block.gasUsed * governance.unitprice
		return new(big.Int).Mul(big.NewInt(int64(header.GasUsed)), config.UnitPrice), nil
	}
}

// getDeferredRewardSimple is for Simple policy.
func getDeferredRewardSimple(config *reward.RewardConfig, totalFee *big.Int) (*reward.RewardSpec, error) {
	spec := reward.NewRewardSpec()
	minted := new(big.Int).Set(config.MintingAmount)

	// Non-deferred mode
	if !config.DeferredTxFee {
		var proposer *big.Int
		if config.Rules.IsMagma {
			// In non-deferred mode, no fees to distribute here at the end of block processing.
			// Just distribute the minting reward to the proposer and stop.
			proposer = new(big.Int).Set(minted)
			totalFee = big.NewInt(0)
		} else {
			// But Simple policy had a bug where transaction fees were distributed to the proposer here at the end of block processing
			// despite configured to non-deferred mode. To keep the backward compatibility, the buggy behavior retains until Magma.
			proposer = new(big.Int).Add(minted, totalFee)
		}
		spec.Minted = new(big.Int).Set(minted)
		spec.TotalFee = totalFee
		spec.BurntFee = big.NewInt(0)
		spec.Proposer = proposer
		spec.IncReceipient(config.Rewardbase, proposer)
		return spec, nil
	}

	// Deferred mode
	burntFee := big.NewInt(0)
	if config.Rules.IsMagma {
		burntFee = getBurnAmountMagma(totalFee)
	}
	proposer := new(big.Int).Add(minted, totalFee)
	proposer.Sub(proposer, burntFee)

	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Proposer = proposer
	spec.IncReceipient(config.Rewardbase, proposer)
	return spec, nil
}

// getDeferredRewardFull is for non-Simple policy.
func getDeferredRewardFull(config *reward.RewardConfig, totalFee *big.Int, si *staking.StakingInfo) (*reward.RewardSpec, error) {
	// Non-deferred and deferred modes share most of the logic
	// except that in non-deferred mode the block fees are considered zero.
	var burntFee *big.Int
	if !config.DeferredTxFee {
		totalFee = big.NewInt(0)
		burntFee = big.NewInt(0)
	} else {
		if config.Rules.IsKore {
			burntFee = getBurnAmountKore(config, totalFee)
		} else if config.Rules.IsMagma {
			burntFee = getBurnAmountMagma(totalFee)
		} else {
			burntFee = big.NewInt(0)
		}
	}

	// Both non-deferred and deferred modes
	if config.Rules.IsKore {
		return getDeferredRewardFullKore(config, totalFee, burntFee, si)
	} else {
		return getDeferredRewardFullLegacy(config, totalFee, burntFee, si)
	}
}

// getDeferredRewardFullKore is for non-Simple policy and after Kore.
func getDeferredRewardFullKore(config *reward.RewardConfig, totalFee, burntFee *big.Int, si *staking.StakingInfo) (*reward.RewardSpec, error) {
	var (
		spec             = reward.NewRewardSpec()
		minted           = new(big.Int).Set(config.MintingAmount)
		distributableFee = new(big.Int).Sub(totalFee, burntFee)
	)

	// Distribute using RewardRatio first. Unlike Legacy, fees are not distributed here
	// because fees are exclusively allocated to proposer. By the way, remainder goes to KIF.
	validators, kif, kef := config.RewardRatio.Split(minted)
	proposer, stakers := config.Kip82Ratio.Split(validators)
	ratioRemainder := calcRemainder(minted, proposer, stakers, kif, kef)
	kif.Add(kif, ratioRemainder)

	// Further distribute using Kip82Ratio. By the way, remainder goes to proposer.
	stakersAlloc, kip82Remainder := assignStakingRewards(config, stakers, si)
	proposer.Add(proposer, kip82Remainder)
	stakers.Sub(stakers, kip82Remainder)

	// Proposer gets the fees.
	proposer.Add(proposer, distributableFee)

	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Stakers = stakers
	for addr, amount := range stakersAlloc {
		spec.IncReceipient(addr, amount)
	}
	spec = specWithProposerAndFunds(spec, config, proposer, kif, kef, si)
	return spec, nil
}

// getDeferredRewardFullLegacy is for non-Simple policy and before Kore.
func getDeferredRewardFullLegacy(config *reward.RewardConfig, totalFee, burntFee *big.Int, si *staking.StakingInfo) (*reward.RewardSpec, error) {
	var (
		spec             = reward.NewRewardSpec()
		minted           = new(big.Int).Set(config.MintingAmount)
		distributableFee = new(big.Int).Sub(totalFee, burntFee)
		totalReward      = new(big.Int).Add(minted, distributableFee)
	)

	// Distribute using RewardRatio. Remainder goes to KIF.
	proposer, kif, kef := config.RewardRatio.Split(totalReward)
	ratioRemainder := calcRemainder(totalReward, proposer, kif, kef)
	kif.Add(kif, ratioRemainder)

	spec.Minted = minted
	spec.TotalFee = totalFee
	spec.BurntFee = burntFee
	spec.Stakers = common.Big0 // No stakers reward before Kore
	spec = specWithProposerAndFunds(spec, config, proposer, kif, kef, si)
	return spec, nil
}

// getBurnAmountMagma returns the amount of fees to be burnt by Magma.
func getBurnAmountMagma(totalFee *big.Int) *big.Int {
	return new(big.Int).Div(totalFee, big.NewInt(2))
}

// getBurnAmountKore returns the amount of fees to be burnt by Kore.
// This includes Magma burnt amount (half of the total fee).
func getBurnAmountKore(config *reward.RewardConfig, totalFee *big.Int) *big.Int {
	firstHalf := new(big.Int).Div(totalFee, big.NewInt(2))
	secondHalf := new(big.Int).Sub(totalFee, firstHalf)

	validatorMintingReward, _, _ := config.RewardRatio.Split(config.MintingAmount)
	proposerMintingReward, _ := config.Kip82Ratio.Split(validatorMintingReward)

	return new(big.Int).Add(
		firstHalf, // half the fee is always burnt
		math.BigMin(secondHalf, proposerMintingReward), // the rest is burnt up to the proposer's minting reward
	)
}

// calcRemainder returns total - sum(parts).
func calcRemainder(total *big.Int, parts ...*big.Int) *big.Int {
	remaining := new(big.Int).Set(total)
	for _, part := range parts {
		remaining.Sub(remaining, part)
	}
	return remaining
}

// assignStakingRewards assigns staking rewards to stakers according to their staking amounts.
// Returns the allocation and the remainder.
func assignStakingRewards(config *reward.RewardConfig, stakersReward *big.Int, si *staking.StakingInfo) (map[common.Address]*big.Int, *big.Int) {
	var (
		cns            = si.ConsolidatedNodes()
		minStake       = config.MinimumStake.Uint64()
		totalExcessInt = uint64(0) // sum of excess stakes (the amount over minStake) over all stakers
	)
	for _, cn := range cns {
		if cn.StakingAmount > minStake {
			totalExcessInt += cn.StakingAmount - minStake
		}
	}

	var (
		totalExcess = new(big.Int).SetUint64(totalExcessInt)
		remaining   = new(big.Int).Set(stakersReward)
		alloc       = make(map[common.Address]*big.Int)
	)
	for _, cn := range cns {
		if cn.StakingAmount > minStake {
			// The KAIA unit will cancel out:
			// reward (kei) = excess (KAIA) * stakersReward (kei) / totalExcess (KAIA)
			excess := new(big.Int).SetUint64(cn.StakingAmount - minStake)
			reward := new(big.Int).Mul(excess, stakersReward)
			reward.Div(reward, totalExcess)
			if reward.Sign() > 0 {
				alloc[cn.RewardAddr] = reward
			}
			remaining.Sub(remaining, reward)
		}
	}
	return alloc, remaining
}

// specWithProposerAndFunds assigns proposer, kif, kef to the reward spec.
// This must be the last step of building the RewardSpec as it finalizes the Proposer, KEF, KIF fields.
func specWithProposerAndFunds(spec *reward.RewardSpec, config *reward.RewardConfig, proposer, kif, kef *big.Int, si *staking.StakingInfo) *reward.RewardSpec {
	newSpec := spec.Copy()

	// If KIF or KEF address is not set, proposer takes it.
	if common.EmptyAddress(si.KIFAddr) {
		newSpec.KIF = common.Big0
		proposer.Add(proposer, kif)
	} else {
		newSpec.KIF = kif
		newSpec.IncReceipient(si.KIFAddr, kif)
	}

	if common.EmptyAddress(si.KEFAddr) {
		newSpec.KEF = common.Big0
		proposer.Add(proposer, kef)
	} else {
		newSpec.KEF = kef
		newSpec.IncReceipient(si.KEFAddr, kef)
	}

	newSpec.Proposer = proposer
	newSpec.IncReceipient(config.Rewardbase, proposer)
	return newSpec
}
