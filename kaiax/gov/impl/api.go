package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

func (g *GovModule) APIs() []rpc.API {
	ret := append(g.Hgm.APIs(), g.Cgm.APIs()...)
	return append(ret, []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   NewGovAPI(g),
			Public:    true,
		},
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewKaiaAPI(g),
			Public:    true,
		},
	}...)
}

type GovAPI struct {
	g *GovModule
}

type KaiaAPI struct {
	g *GovModule
}

func NewGovAPI(g *GovModule) *GovAPI {
	return &GovAPI{g}
}

func (api *GovAPI) GetParams(num *rpc.BlockNumber) (gov.PartialParamSet, error) {
	return getParams(api.g, num)
}

func (api *GovAPI) NodeAddress() (common.Address, error) {
	return api.g.Hgm.NodeAddress(), nil
}

func NewKaiaAPI(g *GovModule) *KaiaAPI {
	return &KaiaAPI{g}
}

func (api *KaiaAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.g, num)
}

func (api *KaiaAPI) GetParams(num *rpc.BlockNumber) (gov.PartialParamSet, error) {
	return getParams(api.g, num)
}

func patchDeprecatedParams(gp gov.ParamSet, rule params.Rules) gov.ParamSet {
	// To avoid confusion, override some parameters that are deprecated after hardforks.
	// e.g., stakingupdateinterval is shown as 86400 but actually irrelevant (i.e. updated every block)
	if rule.IsKore {
		// Gini option deprecated since Kore, as All committee members have an equal chance
		// of being elected block proposers.
		gp.UseGiniCoeff = false
	}
	if rule.IsRandao {
		// Block proposer is randomly elected at every block with Randao,
		// no more precalculated proposer list.
		gp.ProposerUpdateInterval = 1
	}
	if rule.IsKaia {
		// Staking information updated every block since Kaia.
		gp.StakingUpdateInterval = 1
	}
	return gp
}

func getChainConfig(g *GovModule, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = g.Chain.CurrentBlock().NumberU64()
	} else {
		blocknum = num.Uint64()
	}

	ret := g.Chain.Config().Copy()
	pset := g.GetParamSet(blocknum)
	rule := ret.Rules(new(big.Int).SetUint64(blocknum))
	pset = patchDeprecatedParams(pset, rule)

	// patch governance parameters
	ret.Istanbul.Epoch = pset.Epoch
	ret.Istanbul.ProposerPolicy = pset.ProposerPolicy
	ret.Istanbul.SubGroupSize = pset.CommitteeSize
	ret.UnitPrice = pset.UnitPrice
	ret.DeriveShaImpl = int(pset.DeriveShaImpl)
	ret.Governance = &params.GovernanceConfig{
		GoverningNode:    pset.GoverningNode,
		GovernanceMode:   pset.GovernanceMode,
		GovParamContract: pset.GovParamContract,
		Reward: &params.RewardConfig{
			MintingAmount:          pset.MintingAmount,
			Ratio:                  pset.Ratio,
			Kip82Ratio:             pset.Kip82Ratio,
			UseGiniCoeff:           pset.UseGiniCoeff,
			DeferredTxFee:          pset.DeferredTxFee,
			StakingUpdateInterval:  pset.StakingUpdateInterval,
			ProposerUpdateInterval: pset.ProposerUpdateInterval,
			MinimumStake:           pset.MinimumStake,
		},
		KIP71: &params.KIP71Config{
			LowerBoundBaseFee:         pset.LowerBoundBaseFee,
			UpperBoundBaseFee:         pset.UpperBoundBaseFee,
			GasTarget:                 pset.GasTarget,
			MaxBlockGasUsedForBaseFee: pset.MaxBlockGasUsedForBaseFee,
			BaseFeeDenominator:        pset.BaseFeeDenominator,
		},
	}
	return ret
}

func getParams(g *GovModule, num *rpc.BlockNumber) (gov.PartialParamSet, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = g.Chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	rule := g.Chain.Config().Rules(new(big.Int).SetUint64(blockNumber))
	gp := g.GetParamSet(blockNumber)
	gp = patchDeprecatedParams(gp, rule)
	ret := gp.ToMap()
	return ret, nil
}

func (api *KaiaAPI) NodeAddress() common.Address {
	return api.g.Hgm.NodeAddress()
}
