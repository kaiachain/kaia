package impl

import (
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/contractgov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

var (
	_ gov.GovModule = (*GovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

//go:generate mockgen -destination=kaiax/gov/impl/mock/blockchain_mock.go github.com/kaiachain/kaia/kaiax/gov/impl BlockChain
type BlockChain interface {
	CurrentBlock() *types.Block
	Config() *params.ChainConfig
	GetHeaderByNumber(val uint64) *types.Header
	StateAt(root common.Hash) (*state.StateDB, error)
	GetReceiptsByBlockHash(blockHash common.Hash) types.Receipts
	GetBlock(hash common.Hash, number uint64) *types.Block
}

type GovModule struct {
	hgm   headergov.HeaderGovModule
	cgm   contractgov.ContractGovModule
	chain BlockChain
}

type InitOpts struct {
	Hgm   headergov.HeaderGovModule
	Cgm   contractgov.ContractGovModule
	Chain BlockChain
}

func NewGovModule() *GovModule {
	return &GovModule{}
}

func (m *GovModule) Init(opts *InitOpts) error {
	if opts == nil {
		return gov.ErrInitNil
	}

	m.hgm = opts.Hgm
	m.cgm = opts.Cgm
	m.chain = opts.Chain

	if m.hgm == nil || m.cgm == nil || m.chain == nil || m.chain.Config() == nil {
		return gov.ErrInitNil
	}
	return nil
}

func (m *GovModule) Start() error {
	logger.Info("GovModule started")
	return errors.Join(
		m.hgm.Start(),
		m.cgm.Start(),
	)
}

func (m *GovModule) Stop() {
	logger.Info("GovModule stopped")
	m.hgm.Stop()
	m.cgm.Stop()
}

func (m *GovModule) isKoreHF(num uint64) bool {
	return m.chain.Config().IsKoreForkEnabled(new(big.Int).SetUint64(num))
}
