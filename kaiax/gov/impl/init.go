package impl

import (
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
)

var (
	_ gov.GovModule = (*GovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
)

type GovModule struct {
	hgm headergov.HeaderGovModule
	cgm contractgov.ContractGovModule
}

type InitOpts struct {
	Hgm headergov.HeaderGovModule
	Cgm contractgov.ContractGovModule
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
