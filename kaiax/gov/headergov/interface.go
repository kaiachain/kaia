package headergov

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
)

//go:generate mockgen -destination=./mock/headergov_mock.go -package=mock_headergov github.com/kaiachain/kaia/kaiax/gov/headergov HeaderGovModule
type HeaderGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	GetParamSet(blockNum uint64) gov.ParamSet
	GetPartialParamSet(blockNum uint64) gov.PartialParamSet
	NodeAddress() common.Address
	PushMyVotes(vote VoteData)
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
