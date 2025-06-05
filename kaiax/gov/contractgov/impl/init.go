package impl

import (
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/state"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
)

var (
	_ contractgov.ContractGovModule = (*contractGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

type chain interface {
	blockchain.ChainContext

	GetHeaderByNumber(number uint64) *types.Header
	CurrentBlock() *types.Block
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig
	GetBlock(hash common.Hash, number uint64) *types.Block
}

type InitOpts struct {
	ChainConfig *params.ChainConfig
	Chain       chain
	Hgm         headergov.HeaderGovModule
}

type contractGovModule struct {
	InitOpts
}

func NewContractGovModule() *contractGovModule {
	return &contractGovModule{}
}

func (c *contractGovModule) Init(opts *InitOpts) error {
	if err := validateOpts(opts); err != nil {
		return err
	}

	c.InitOpts = *opts
	return nil
}

func (c *contractGovModule) Start() error {
	logger.Info("ContractGovModule started")
	return nil
}

func (c *contractGovModule) Stop() {
	logger.Info("ContractGovModule stopped")
}

func validateOpts(opts *InitOpts) error {
	switch {
	case opts == nil:
		return errInitNil("opts")
	case opts.ChainConfig == nil:
		return errInitNil("opts.ChainConfig")
	case opts.Chain == nil:
		return errInitNil("opts.Chain")
	case opts.Hgm == nil:
		return errInitNil("opts.Hgm")
	}
	return nil
}
