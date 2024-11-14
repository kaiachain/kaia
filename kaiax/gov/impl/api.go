package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

func (g *GovModule) APIs() []rpc.API {
	ret := append(g.hgm.APIs(), g.cgm.APIs()...)
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
	return api.g.hgm.NodeAddress(), nil
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

func getChainConfig(g *GovModule, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = g.chain.CurrentBlock().NumberU64()
	} else {
		blocknum = num.Uint64()
	}

	pset := g.EffectiveParamSet(blocknum)
	latestConfig := g.chain.Config()
	config := pset.ToGovParamSet().ToChainConfig()
	config.ChainID = latestConfig.ChainID
	config.IstanbulCompatibleBlock = latestConfig.IstanbulCompatibleBlock
	config.LondonCompatibleBlock = latestConfig.LondonCompatibleBlock
	config.EthTxTypeCompatibleBlock = latestConfig.EthTxTypeCompatibleBlock
	config.MagmaCompatibleBlock = latestConfig.MagmaCompatibleBlock
	config.KoreCompatibleBlock = latestConfig.KoreCompatibleBlock
	config.ShanghaiCompatibleBlock = latestConfig.ShanghaiCompatibleBlock
	config.CancunCompatibleBlock = latestConfig.CancunCompatibleBlock
	config.KaiaCompatibleBlock = latestConfig.KaiaCompatibleBlock
	config.Kip103CompatibleBlock = latestConfig.Kip103CompatibleBlock
	config.Kip103ContractAddress = latestConfig.Kip103ContractAddress
	config.Kip160CompatibleBlock = latestConfig.Kip160CompatibleBlock
	config.Kip160ContractAddress = latestConfig.Kip160ContractAddress
	config.RandaoCompatibleBlock = latestConfig.RandaoCompatibleBlock

	return config
}

// checkStateForStakingInfo checks the state of block for the given block number for staking info
func checkStateForStakingInfo(g *GovModule, blockNumber uint64) error {
	if blockNumber == 0 {
		return nil
	}

	// The staking info at blockNumber is calculated by the state of previous block
	blockNumber--
	if !g.chain.Config().IsKaiaForkEnabled(big.NewInt(int64(blockNumber + 1))) {
		return nil
	}
	header := g.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return gov.ErrUnknownBlock
	}
	_, err := g.chain.StateAt(header.Root)
	return err
}

func getParams(g *GovModule, num *rpc.BlockNumber) (gov.PartialParamSet, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = g.chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	gp := g.EffectiveParamSet(blockNumber)
	ret := gp.ToMap()
	// To avoid confusion, override some parameters that are deprecated after hardforks.
	// e.g., stakingupdateinterval is shown as 86400 but actually irrelevant (i.e. updated every block)
	rule := g.chain.Config().Rules(new(big.Int).SetUint64(blockNumber))
	if rule.IsKore {
		// Gini option deprecated since Kore, as All committee members have an equal chance
		// of being elected block proposers.
		if _, ok := ret[gov.RewardUseGiniCoeff]; ok {
			ret[gov.RewardUseGiniCoeff] = false
		}
	}
	if rule.IsRandao {
		// Block proposer is randomly elected at every block with Randao,
		// no more precalculated proposer list.
		if _, ok := ret[gov.RewardProposerUpdateInterval]; ok {
			ret[gov.RewardProposerUpdateInterval] = 1
		}
	}
	if rule.IsKaia {
		// Staking information updated every block since Kaia.
		if _, ok := ret[gov.RewardStakingUpdateInterval]; ok {
			ret[gov.RewardStakingUpdateInterval] = 1
		}
	}

	return ret, nil
}

func (api *KaiaAPI) NodeAddress() common.Address {
	return api.g.hgm.NodeAddress()
}
