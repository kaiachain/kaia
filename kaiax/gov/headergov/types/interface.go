package types

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax"
	govtypes "github.com/kaiachain/kaia/kaiax/gov/types"
)

//go:generate mockgen -destination=kaiax/gov/headergov/mocks/headergov_mock.go github.com/kaiachain/kaia/kaiax/gov/headergov/types HeaderGovModule
type HeaderGovModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ConsensusModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	EffectiveParamSet(blockNum uint64) (ParamSet, error)
	EffectiveParamsPartial(blockNum uint64) (map[string]interface{}, error)
}

type GovData interface {
	Items() map[ParamEnum]interface{}
	Serialize() ([]byte, error)
}

type VoteData interface {
	Voter() common.Address
	Name() string
	Type() govtypes.ParamEnum
	Value() interface{}

	Serialize() ([]byte, error)
}

var (
	_ GovData  = (*govData)(nil)
	_ VoteData = (*voteData)(nil)
)
