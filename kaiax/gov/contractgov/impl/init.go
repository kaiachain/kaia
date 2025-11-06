package impl

import (
	"sync"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/rcrowley/go-metrics"
)

var (
	_ contractgov.ContractGovModule = (*contractGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)

	// Cache metrics
	cacheHits   = metrics.NewRegisteredCounter("gov/contractgov/cache/hits", nil)
	cacheMisses = metrics.NewRegisteredCounter("gov/contractgov/cache/misses", nil)
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

	// Cache for parameter sets by contract address + storage root
	paramSetCache *lru.Cache[[64]byte, gov.PartialParamSet] // contract addr + storage root -> param set
	cacheMutex    sync.RWMutex
}

func NewContractGovModule() *contractGovModule {
	paramCache, _ := lru.New[[64]byte, gov.PartialParamSet](100) // Cache up to 100 different storage roots
	return &contractGovModule{
		paramSetCache: paramCache,
	}
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
