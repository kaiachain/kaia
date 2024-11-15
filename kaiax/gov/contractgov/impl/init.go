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
	if opts == nil || opts.ChainConfig == nil || opts.ChainConfig.Istanbul == nil || opts.Hgm == nil {
		return ErrInitNil
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
