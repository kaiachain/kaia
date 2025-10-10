package impl

import (
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	_ contractgov.ContractGovModule = (*contractGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

type chain interface {
	consensus.ChainReader
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
