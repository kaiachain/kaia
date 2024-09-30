package impl

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
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

type VotesResponse struct {
	BlockNum uint64
	Key      string
	Value    any
}

type MyVotesResponse struct {
	BlockNum uint64
	Key      string
	Value    any
	Casted   bool
}

type StatusResponse struct {
	GroupedVotes map[uint64]headergov.VotesInEpoch `json:"groupedVotes"`
	Governances  map[uint64]headergov.GovData      `json:"governances"`
	GovHistory   headergov.History                 `json:"govHistory"`
	NodeAddress  common.Address                    `json:"nodeAddress"`
	MyVotes      []headergov.VoteData              `json:"myVotes"`
}

func NewHeaderGovAPI(s *headerGovModule) *headerGovAPI {
	return &headerGovAPI{s}
}

func (api *headerGovAPI) Vote(name string, value any) (string, error) {
	var (
		voter       = api.h.nodeAddress
		blockNumber = api.h.Chain.CurrentBlock().NumberU64()
		gp          = api.h.EffectiveParamSet(blockNumber + 1)
		gMode       = gp.GovernanceMode
	)

	if gMode == "single" && voter != gp.GoverningNode {
		return "", ErrVotePermissionDenied
	}

	vote := headergov.NewVoteData(voter, name, value)
	if vote == nil {
		return "", ErrInvalidKeyValue
	}

	err := api.h.checkConsistency(blockNumber+1, vote)
	if err != nil {
		return "", err
	}

	// TODO-kaiax: add removevalidator vote check

	api.h.PushMyVotes(vote)
	return "(kaiax) Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
}

func (api *headerGovAPI) IdxCache() []uint64 {
	return api.h.cache.GovBlockNums()
}

func (api *headerGovAPI) Votes(num *rpc.BlockNumber) []VotesResponse {
	var blockNum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNum = api.h.Chain.CurrentBlock().NumberU64()
	} else {
		blockNum = num.Uint64()
	}

	epochIdx := calcEpochIdx(blockNum, api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]VotesResponse, 0)
	for blockNum, vote := range votesInEpoch {
		ret = append(ret, VotesResponse{
			BlockNum: blockNum,
			Key:      string(vote.Name()),
			Value:    vote.Value(),
		})
	}
	return ret
}

func (api *headerGovAPI) MyVotes() []MyVotesResponse {
	epochIdx := calcEpochIdx(api.h.Chain.CurrentBlock().NumberU64(), api.h.epoch)
	votesInEpoch := api.h.getVotesInEpoch(epochIdx)

	ret := make([]MyVotesResponse, 0)
	for blockNum, vote := range votesInEpoch {
		if vote.Voter() == api.h.nodeAddress {
			ret = append(ret, MyVotesResponse{
				BlockNum: blockNum,
				Casted:   true,
				Key:      string(vote.Name()),
				Value:    vote.Value(),
			})
		}
	}

	for _, vote := range api.h.myVotes {
		ret = append(ret, MyVotesResponse{
			BlockNum: 0,
			Casted:   false,
			Key:      string(vote.Name()),
			Value:    vote.Value(),
		})
	}

	return ret
}

func (api *headerGovAPI) NodeAddress() common.Address {
	return api.h.nodeAddress
}

func (api *headerGovAPI) GetParams(num *rpc.BlockNumber) (map[gov.ParamName]any, error) {
	return api.getParams(num)
}

func (api *headerGovAPI) getParams(num *rpc.BlockNumber) (map[gov.ParamName]any, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.h.Chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	gp := api.h.EffectiveParamSet(blockNumber)
	return gp.ToMap(), nil
}

func (api *headerGovAPI) Status() StatusResponse {
	return StatusResponse{
		GroupedVotes: api.h.cache.GroupedVotes(),
		Governances:  api.h.cache.Govs(),
		GovHistory:   api.h.cache.History(),
		NodeAddress:  api.h.nodeAddress,
		MyVotes:      api.h.myVotes,
	}
}
