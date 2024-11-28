package impl

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	contractgov_impl "github.com/kaiachain/kaia/kaiax/gov/contractgov/impl"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	headergov_impl "github.com/kaiachain/kaia/kaiax/gov/headergov/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ gov.GovModule = (*GovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

//go:generate mockgen -destination=mock/blockchain_mock.go github.com/kaiachain/kaia/kaiax/gov/impl BlockChain
type BlockChain interface {
	blockchain.ChainContext

	CurrentBlock() *types.Block
	Config() *params.ChainConfig
	GetHeaderByNumber(val uint64) *types.Header
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts
	GetBlock(hash common.Hash, number uint64) *types.Block
}

type GovModule struct {
	Chain       BlockChain
	ChainConfig *params.ChainConfig
	Hgm         headergov.HeaderGovModule
	Cgm         contractgov.ContractGovModule
}

type InitOpts struct {
	ChainConfig *params.ChainConfig
	ChainKv     database.Database
	Chain       BlockChain
	NodeAddress common.Address
}

func NewGovModule() *GovModule {
	return &GovModule{}
}

func (m *GovModule) Init(opts *InitOpts) error {
	if opts == nil {
		return gov.ErrInitNil
	}

	var (
		hgm = headergov_impl.NewHeaderGovModule()
		cgm = contractgov_impl.NewContractGovModule()
	)

	err := errors.Join(
		hgm.Init(&headergov_impl.InitOpts{
			ChainKv:     opts.ChainKv,
			ChainConfig: opts.ChainConfig,
			Chain:       opts.Chain,
			NodeAddress: opts.NodeAddress,
		}),
		cgm.Init(&contractgov_impl.InitOpts{
			ChainConfig: opts.ChainConfig,
			Chain:       opts.Chain,
			Hgm:         hgm,
		}),
	)
	if err != nil {
		return err
	}

	m.Chain = opts.Chain
	m.ChainConfig = opts.ChainConfig
	m.Hgm = hgm
	m.Cgm = cgm
	return nil
}

func (m *GovModule) Start() error {
	logger.Info("GovModule started")
	return errors.Join(
		m.Hgm.Start(),
		m.Cgm.Start(),
	)
}

func (m *GovModule) Stop() {
	logger.Info("GovModule stopped")
	m.Hgm.Stop()
	m.Cgm.Stop()
}

func (m *GovModule) isKoreHF(num uint64) bool {
	return m.Chain.Config().IsKoreForkEnabled(new(big.Int).SetUint64(num))
}
