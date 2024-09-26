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

	EffectiveParamSet(blockNum uint64) (gov.ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[gov.ParamEnum]any, error)
	NodeAddress() common.Address
}

type GovData interface {
	Items() map[gov.ParamEnum]any
	ToGovBytes() (GovBytes, error)
}

type VoteData interface {
	Voter() common.Address
	Name() string
	Enum() gov.ParamEnum
	Value() any

	ToVoteBytes() (VoteBytes, error)
}

var (
	_ GovData  = (*govData)(nil)
	_ VoteData = (*voteData)(nil)
)
