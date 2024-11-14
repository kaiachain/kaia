package impl

import (
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	_ contractgov.ContractGovModule = (*contractGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
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
	ChainConfig *params.ChainConfig
	Chain       chain
	hgm         headergov.HeaderGovModule
}

func NewContractGovModule() *contractGovModule {
	return &contractGovModule{}
}

func (c *contractGovModule) Init(opts *InitOpts) error {
	if opts == nil {
		return ErrInitNil
	}

	c.ChainConfig = opts.ChainConfig
	c.Chain = opts.Chain
	c.hgm = opts.Hgm
	if c.ChainConfig == nil || c.ChainConfig.Istanbul == nil || c.hgm == nil {
		return ErrInitNil
	}

	return nil
}

func (c *contractGovModule) Start() error {
	logger.Info("ContractGovModule started")
	return nil
}

func (c *contractGovModule) Stop() {
	logger.Info("ContractGovModule stopped")
}
