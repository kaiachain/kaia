package impl

import (
	"maps"
	"math/big"
	"slices"

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
		voter     = api.h.nodeAddress
		nextBlock = api.h.Chain.CurrentBlock().NumberU64() + 1
		gp        = api.h.GetParamSet(nextBlock)
		gMode     = gp.GovernanceMode
	)

	if gMode == "single" && voter != gp.GoverningNode {
		return "", ErrVotePermissionDenied
	}

	// Fail-fast check at queue time: block voting if not currently eligible.
	// The canonical enforcement is at VerifyVote (proposal time), but blocking
	// here gives immediate feedback and avoids queuing votes that will never apply.
	if api.h.ValSet != nil {
		voters, err := api.h.ValSet.GetHeaderGovVoters(nextBlock)
		if err != nil {
			return "", err
		}
		if !slices.Contains(voters, voter) {
			return "", ErrVotePermissionDenied
		}
	}

	vote := headergov.NewVoteData(voter, name, value)
	if vote == nil {
		return "", ErrInvalidKeyValue
	}

	if gov.DeprecatedAt(vote.Name(), api.h.ChainConfig.Rules(new(big.Int).SetUint64(nextBlock))) {
		return "", ErrDeprecatedVote
	}

	err := api.h.checkConsistency(nextBlock, vote)
	if err != nil {
		return "", err
	}

	// TODO-kaiax: add removevalidator vote check

	api.h.PushMyVotes(vote)
	return "(kaiax) Your vote is prepared. It will be put into the block header or applied when your node generates a block as a proposer. Note that your vote may be duplicate.", nil
}

func (api *headerGovAPI) IdxCache() []uint64 {
	return api.h.GovBlockNums()
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
	pendingVotes := api.h.myVotesSnapshot()

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

	for _, vote := range pendingVotes {
		ret = append(ret, MyVotesResponse{
			BlockNum: 0,
			Casted:   false,
			Key:      string(vote.Name()),
			Value:    vote.Value(),
		})
	}

	return ret
}

func (api *headerGovAPI) Status() StatusResponse {
	api.h.mu.RLock()
	defer api.h.mu.RUnlock()

	groupedVotes := make(map[uint64]headergov.VotesInEpoch, len(api.h.groupedVotes))
	for epochIdx, votes := range api.h.groupedVotes {
		copiedVotes := make(headergov.VotesInEpoch, len(votes))
		maps.Copy(copiedVotes, votes)
		groupedVotes[epochIdx] = copiedVotes
	}

	governances := make(map[uint64]headergov.GovData, len(api.h.governances))
	maps.Copy(governances, api.h.governances)

	govHistory := make(headergov.History, len(api.h.history))
	maps.Copy(govHistory, api.h.history)

	myVotes := make([]headergov.VoteData, len(api.h.myVotes))
	copy(myVotes, api.h.myVotes)

	return StatusResponse{
		GroupedVotes: groupedVotes,
		Governances:  governances,
		GovHistory:   govHistory,
		NodeAddress:  api.h.nodeAddress,
		MyVotes:      myVotes,
	}
}
