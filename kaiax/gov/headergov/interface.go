package headergov

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
)

//go:generate mockgen -destination=kaiax/gov/headergov/mock/headergov_mock.go github.com/kaiachain/kaia/kaiax/gov/headergov HeaderGovModule
type HeaderGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	EffectiveParamSet(blockNum uint64) gov.ParamSet
	EffectiveParamsPartial(blockNum uint64) gov.PartialParamSet
	NodeAddress() common.Address
	GetLatestValidatorVote(num uint64) (uint64, VoteData)
	GetMyVotes() []VoteData
	PopMyVotes(idx int)
}

type GovData interface {
	Items() gov.PartialParamSet
	ToGovBytes() (GovBytes, error)
}

type VoteData interface {
	Voter() common.Address
	Name() gov.ParamName
	Value() any

	ToVoteBytes() (VoteBytes, error)
}

var (
	_ GovData  = (*govData)(nil)
	_ VoteData = (*voteData)(nil)
)
