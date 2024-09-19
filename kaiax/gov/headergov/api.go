package headergov

import (
	"encoding/json"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/rpc"
)

func (h *headerGovModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   NewHeaderGovAPI(h),
			Public:    true,
		},
	}
}

type headerGovAPI struct {
	h *headerGovModule
}

func NewHeaderGovAPI(s *headerGovModule) *headerGovAPI {
	return &headerGovAPI{s}
}

func (api *headerGovAPI) Vote(name string, value interface{}) (string, error) {
	blockNumber := api.h.Chain.CurrentBlock().NumberU64()
	gp, err := api.h.EffectiveParamSet(blockNumber + 1)
	if err != nil {
		return "", err
	}

	gMode := gp.GovernanceMode
	if gMode == "single" && api.h.NodeAddress != gp.GoverningNode {
		return "", ErrVotePermissionDenied
	}

	vote := NewVoteData(api.h.NodeAddress, name, value)
	if vote == nil {
		return "", ErrInvalidKeyValue
	}

	err = api.h.VerifyVote(blockNumber+1, vote)
	if err != nil {
		return "", err
	}

	// TODO-kaiax: add removevalidator vote check

	api.h.PushMyVotes(vote)
	return "Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
}

func (api *headerGovAPI) IdxCache() []uint64 {
	return api.h.cache.GovBlockNums()
}

type MyVotesAPI struct {
	BlockNum uint64
	Casted   bool
	Key      string
	Value    interface{}
}

func (api *headerGovAPI) MyVotes() []MyVotesAPI {
	epochIdx := calcEpochIdx(api.h.Chain.CurrentBlock().NumberU64(), api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]MyVotesAPI, 0)
	for blockNum, vote := range votesInEpoch {
		if vote.Voter() == api.h.NodeAddress {
			ret = append(ret, MyVotesAPI{
				BlockNum: blockNum,
				Casted:   true,
				Key:      vote.Name(),
				Value:    vote.Value(),
			})
		}
	}

	for _, vote := range api.h.myVotes {
		ret = append(ret, MyVotesAPI{
			BlockNum: 0, // TODO: remove
			Casted:   false,
			Key:      vote.Name(),
			Value:    vote.Value(),
		})
	}

	return ret
}

// PendingVotes returns all pending votes in the current epoch.
func (api *headerGovAPI) PendingVotes() []VoteData {
	epochIdx := calcEpochIdx(api.h.Chain.CurrentBlock().NumberU64(), api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]VoteData, 0)
	for _, vote := range votesInEpoch {
		if vote.Voter() == api.h.NodeAddress {
			ret = append(ret, vote)
		}
	}

	return ret
}

func (api *headerGovAPI) NodeAddress() common.Address {
	return api.h.NodeAddress
}

func (api *headerGovAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return api.getParams(num)
}

func (api *headerGovAPI) getParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.h.Chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	gp, err := api.h.EffectiveParamSet(blockNumber)
	if err != nil {
		return nil, err
	}
	return EnumMapToStrMap(gp.ToEnumMap()), nil
}

func (api *headerGovAPI) VotesInEpoch(blockNum uint64) string {
	epochIdx := calcEpochIdx(blockNum, api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)
	j, _ := json.Marshal(votesInEpoch)
	return string(j)
}

func (api *headerGovAPI) Status() (string, error) {
	type StatusApi struct {
		GroupedVotes map[uint64]VotesInEpoch `json:"groupedVotes"`
		Governances  map[uint64]GovData      `json:"governances"`
		GovHistory   History                 `json:"govHistory"`
		NodeAddress  common.Address          `json:"nodeAddress"`
		MyVotes      []VoteData              `json:"myVotes"`
	}
	publicCache := StatusApi{
		GroupedVotes: api.h.cache.GroupedVotes(),
		Governances:  api.h.cache.Govs(),
		GovHistory:   api.h.cache.History(),
		NodeAddress:  api.h.NodeAddress,
		MyVotes:      api.h.myVotes,
	}

	cacheJson, err := json.Marshal(publicCache)
	if err != nil {
		logger.Error("kaiax: Failed to marshal cache", "err", err)
		return "", err
	}

	return string(cacheJson), nil
}
