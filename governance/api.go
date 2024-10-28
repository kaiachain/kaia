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

package governance

import (
	"errors"
	"math/big"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
)

type GovernanceAPI struct {
	governance    Engine // Node interfaced by this API
	stakingModule staking.StakingModule
}

type returnTally struct {
	Key                string
	Value              interface{}
	ApprovalPercentage float64
}

func NewGovernanceAPI(gov Engine, stakingModule staking.StakingModule) *GovernanceAPI {
	return &GovernanceAPI{governance: gov, stakingModule: stakingModule}
}

type GovernanceKaiaAPI struct {
	governance    Engine
	chain         blockChain
	stakingModule staking.StakingModule
}

func NewGovernanceKaiaAPI(gov Engine, chain blockChain, stakingModule staking.StakingModule) *GovernanceKaiaAPI {
	return &GovernanceKaiaAPI{governance: gov, chain: chain, stakingModule: stakingModule}
}

var (
	errUnknownBlock           = errors.New("Unknown block")
	errNotAvailableInThisMode = errors.New("In current governance mode, voting power is not available")
	errSetDefaultFailure      = errors.New("Failed to set a default value")
	errPermissionDenied       = errors.New("You don't have the right to vote")
	errRemoveSelf             = errors.New("You can't vote on removing yourself")
	errInvalidKeyValue        = errors.New("Your vote couldn't be placed. Please check your vote's key and value")
	errInvalidLowerBound      = errors.New("lowerboundbasefee cannot be set exceeding upperboundbasefee")
	errInvalidUpperBound      = errors.New("upperboundbasefee cannot be set lower than lowerboundbasefee")
)

func (api *GovernanceKaiaAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.governance, num)
}

func (api *GovernanceKaiaAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.governance, num)
}

func (api *GovernanceKaiaAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

// Vote injects a new vote for governance targets such as unitprice and governingnode.
func (api *GovernanceAPI) Vote(key string, val interface{}) (string, error) {
	blockNumber := api.governance.BlockChain().CurrentBlock().NumberU64()
	pset, err := api.governance.EffectiveParams(blockNumber + 1)
	if err != nil {
		return "", err
	}
	gMode := pset.GovernanceModeInt()
	gNode := pset.GoverningNode()

	if gMode == params.GovernanceMode_Single && gNode != api.governance.NodeAddress() {
		return "", errPermissionDenied
	}
	vote, ok := api.governance.ValidateVote(&GovernanceVote{Key: strings.ToLower(key), Value: val})
	if !ok {
		return "", errInvalidKeyValue
	}
	if vote.Key == "governance.removevalidator" {
		if api.isRemovingSelf(val.(string)) {
			return "", errRemoveSelf
		}
	}
	if vote.Key == "kip71.lowerboundbasefee" {
		if vote.Value.(uint64) > pset.UpperBoundBaseFee() {
			return "", errInvalidLowerBound
		}
	}
	if vote.Key == "kip71.upperboundbasefee" {
		if vote.Value.(uint64) < pset.LowerBoundBaseFee() {
			return "", errInvalidUpperBound
		}
	}
	if api.governance.AddVote(key, val) {
		return "Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
	}
	return "", errInvalidKeyValue
}

func (api *GovernanceAPI) isRemovingSelf(val string) bool {
	for _, str := range strings.Split(val, ",") {
		str = strings.Trim(str, " ")
		if common.HexToAddress(str) == api.governance.NodeAddress() {
			return true
		}
	}
	return false
}

func (api *GovernanceAPI) ShowTally() []*returnTally {
	ret := []*returnTally{}

	for _, val := range api.governance.GetGovernanceTalliesCopy() {
		item := &returnTally{
			Key:                val.Key,
			Value:              val.Value,
			ApprovalPercentage: float64(val.Votes) / float64(api.governance.TotalVotingPower()) * 100,
		}
		ret = append(ret, item)
	}

	return ret
}

func (api *GovernanceAPI) TotalVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.TotalVotingPower()) / 1000.0, nil
}

func (api *GovernanceAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.governance, num)
}

func getParams(governance Engine, num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	pset, err := governance.EffectiveParams(blockNumber)
	if err != nil {
		return nil, err
	}
	sm := pset.StrMap()

	// To avoid confusion, override some parameters that are deprecated after hardforks.
	// e.g., stakingupdateinterval is shown as 86400 but actually irrelevant (i.e. updated every block)
	rule := governance.BlockChain().Config().Rules(new(big.Int).SetUint64(blockNumber))
	if rule.IsKore {
		// Gini option deprecated since Kore, as All committee members have an equal chance
		// of being elected block proposers.
		if _, ok := sm["reward.useginicoeff"]; ok {
			sm["reward.useginicoeff"] = false
		}
	}
	if rule.IsRandao {
		// Block proposer is randomly elected at every block with Randao,
		// no more precalculated proposer list.
		if _, ok := sm["reward.proposerupdateinterval"]; ok {
			sm["reward.proposerupdateinterval"] = 1
		}
	}
	if rule.IsKaia {
		// Staking information updated every block since Kaia.
		if _, ok := sm["reward.stakingupdateinterval"]; ok {
			sm["reward.stakingupdateinterval"] = 1
		}
	}
	return sm, nil
}

// Checks the state of block for the given block number for staking info
func checkStateForStakingInfo(governance Engine, blockNumber uint64) error {
	if blockNumber == 0 {
		return nil
	}

	// The staking info at blockNumber is calculated by the state of previous block
	blockNumber--
	if !governance.BlockChain().Config().IsKaiaForkEnabled(big.NewInt(int64(blockNumber + 1))) {
		return nil
	}
	header := governance.BlockChain().GetHeaderByNumber(blockNumber)
	if header == nil {
		return errUnknownBlock
	}
	_, err := governance.BlockChain().StateAt(header.Root)
	return err
}

func (api *GovernanceAPI) PendingChanges() map[string]interface{} {
	return api.governance.PendingChanges()
}

func (api *GovernanceAPI) Votes() []GovernanceVote {
	return api.governance.Votes()
}

func (api *GovernanceAPI) IdxCache() []uint64 {
	return api.governance.IdxCache()
}

func (api *GovernanceAPI) IdxCacheFromDb() []uint64 {
	return api.governance.IdxCacheFromDb()
}

// TODO-Kaia: Return error if invalid input is given such as pending or a too big number
func (api *GovernanceAPI) ItemCacheFromDb(num *rpc.BlockNumber) map[string]interface{} {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	ret, _ := api.governance.DB().ReadGovernance(blockNumber)
	return ret
}

type VoteList struct {
	Key      string
	Value    interface{}
	Casted   bool
	BlockNum uint64
}

func (api *GovernanceAPI) MyVotes() []*VoteList {
	ret := []*VoteList{}

	for k, v := range api.governance.GetVoteMapCopy() {
		item := &VoteList{
			Key:      k,
			Value:    v.Value,
			Casted:   v.Casted,
			BlockNum: v.Num,
		}
		ret = append(ret, item)
	}

	return ret
}

func (api *GovernanceAPI) MyVotingPower() (float64, error) {
	if !api.isGovernanceModeBallot() {
		return 0, errNotAvailableInThisMode
	}
	return float64(api.governance.MyVotingPower()) / 1000.0, nil
}

func (api *GovernanceAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.governance, num)
}

func getChainConfig(governance Engine, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = governance.BlockChain().CurrentBlock().NumberU64()
	} else {
		blocknum = num.Uint64()
	}

	pset, err := governance.EffectiveParams(blocknum)
	if err != nil {
		return nil
	}

	// Fill in the non-governance-parameter fields of ChainConfig
	latestConfig := governance.BlockChain().Config()
	config := pset.ToChainConfig()
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
	config.PragueCompatibleBlock = latestConfig.PragueCompatibleBlock

	// To avoid confusion, override some parameters that are deprecated after hardforks.
	// e.g., stakingupdateinterval is shown as 86400 but actually irrelevant (i.e. updated every block)
	rule := governance.BlockChain().Config().Rules(new(big.Int).SetUint64(blocknum))
	if rule.IsKore {
		// Gini option deprecated since Kore, as All committee members have an equal chance
		// of being elected block proposers.
		config.Governance.Reward.UseGiniCoeff = false
	}
	if rule.IsRandao {
		// Block proposer is randomly elected at every block with Randao,
		// no more precalculated proposer list.
		config.Governance.Reward.ProposerUpdateInterval = 1
	}
	if rule.IsKaia {
		// Staking information updated every block since Kaia.
		config.Governance.Reward.StakingUpdateInterval = 1
	}

	return config
}

func (api *GovernanceAPI) NodeAddress() common.Address {
	return api.governance.NodeAddress()
}

func (api *GovernanceAPI) isGovernanceModeBallot() bool {
	blockNumber := api.governance.BlockChain().CurrentBlock().NumberU64()
	pset, err := api.governance.EffectiveParams(blockNumber + 1)
	if err != nil {
		return false
	}
	gMode := pset.GovernanceModeInt()
	return gMode == params.GovernanceMode_Ballot
}

// Disabled APIs
// func (api *GovernanceKaiaAPI) GetTxGasHumanReadable(num *rpc.BlockNumber) (uint64, error) {
// 	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
// 		// If the value hasn't been set in governance, set it with default value
// 		if ret := api.governance.GetGovernanceValue(params.ConstTxGasHumanReadable); ret == nil {
// 			return api.setDefaultTxGasHumanReadable()
// 		} else {
// 			return ret.(uint64), nil
// 		}
// 	} else {
// 		blockNum := num.Int64()
//
// 		if blockNum > api.chain.CurrentBlock().NumberU64() {
// 			return 0, errUnknownBlock
// 		}
//
// 		if ret, err := api.governance.GetGovernanceItemAtNumber(uint64(blockNum), GovernanceKeyMapReverse[params.ConstTxGasHumanReadable]); err == nil && ret != nil {
// 			return ret.(uint64), nil
// 		} else {
// 			logger.Error("Failed to retrieve TxGasHumanReadable, sending default value", "err", err)
// 			return api.setDefaultTxGasHumanReadable()
// 		}
// 	}
// }
//
// func (api *GovernanceKaiaAPI) setDefaultTxGasHumanReadable() (uint64, error) {
// 	err := api.governance.currentSet.SetValue(params.ConstTxGasHumanReadable, params.TxGasHumanReadable)
// 	if err != nil {
// 		return 0, errSetDefaultFailure
// 	} else {
// 		return params.TxGasHumanReadable, nil
// 	}
// }
