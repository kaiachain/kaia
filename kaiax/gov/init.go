package gov

import (
	"github.com/kaiachain/kaia/log"
)

var (
	_ GovModule = (*govModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
)

type govModule struct {
	hgm HeaderGovModule
	cgm ContractGovModule
}

type InitOpts struct {
	hgm HeaderGovModule
	cgm ContractGovModule
}

func (m *govModule) Init(opts *InitOpts) error {
	m.hgm = opts.hgm
	m.cgm = opts.cgm
	return nil
}

func (m *govModule) Start() error {
	logger.Info("GovModule started")
	return nil
}

func (m *govModule) Stop() {
	logger.Info("GovModule stopped")
}
