package gov

import (
	gov_types "github.com/kaiachain/kaia/kaiax/gov/types"
	"github.com/kaiachain/kaia/log"
)

var (
	_ gov_types.GovModule = (*GovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
)

type GovModule struct {
	hgm HeaderGovModule
	cgm ContractGovModule
}

type InitOpts struct {
	Hgm HeaderGovModule
	Cgm ContractGovModule
}

func NewGovModule() *GovModule {
	return &GovModule{}
}

func (m *GovModule) Init(opts *InitOpts) error {
	m.hgm = opts.Hgm
	m.cgm = opts.Cgm
	return nil
}

func (m *GovModule) Start() error {
	logger.Info("GovModule started")
	return nil
}

func (m *GovModule) Stop() {
	logger.Info("GovModule stopped")
}
