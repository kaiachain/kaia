package impl

import (
	"errors"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/state"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/contractgov"
	contractgov_impl "github.com/kaiachain/kaia/v2/kaiax/gov/contractgov/impl"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	headergov_impl "github.com/kaiachain/kaia/v2/kaiax/gov/headergov/impl"
	"github.com/kaiachain/kaia/v2/kaiax/valset"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
	"golang.org/x/exp/maps" // TODO: use "maps"
)

var (
	_ gov.GovModule = (*GovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

//go:generate mockgen -destination=./mock/blockchain_mock.go -package=mock_impl github.com/kaiachain/kaia/v2/kaiax/gov/impl BlockChain
type BlockChain interface {
	blockchain.ChainContext

	CurrentBlock() *types.Block
	Config() *params.ChainConfig
	GetHeaderByNumber(val uint64) *types.Header
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	GetBlock(hash common.Hash, number uint64) *types.Block
}

type GovModule struct {
	Fallback    gov.PartialParamSet
	ChainConfig *params.ChainConfig
	Chain       BlockChain
	Hgm         headergov.HeaderGovModule
	Cgm         contractgov.ContractGovModule
}

type InitOpts struct {
	ChainConfig *params.ChainConfig
	ChainKv     database.Database
	Chain       BlockChain
	Valset      valset.ValsetModule
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
			ValSet:      opts.Valset,
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

	m.Fallback = ChainConfigFallback(opts.ChainConfig)
	m.Chain = opts.Chain
	m.ChainConfig = opts.ChainConfig
	m.Hgm = hgm
	m.Cgm = cgm
	return nil
}

// ChainConfigFallback returns the set of parameters that have different values between ChainConfig and default values.
func ChainConfigFallback(chainConfig *params.ChainConfig) gov.PartialParamSet {
	fallback := make(gov.PartialParamSet)

	if chainConfig == nil {
		logger.Info("Using empty fallback")
		return fallback
	}

	var candidates []gov.ParamName

	// on Mainnet/Kairos, fallback candidates are only the initial params specified in `params.{Mainnet,Kairos}ChainConfig`.
	if chainId := chainConfig.ChainID; chainId != nil &&
		(chainId.Cmp(params.MainnetChainConfig.ChainID) == 0 ||
			chainId.Cmp(params.KairosChainConfig.ChainID) == 0) {
		candidates = []gov.ParamName{
			gov.GovernanceDeriveShaImpl, gov.GovernanceGoverningNode, gov.GovernanceGovernanceMode, gov.RewardMintingAmount,
			gov.RewardRatio, gov.RewardUseGiniCoeff, gov.RewardDeferredTxFee, gov.RewardStakingUpdateInterval,
			gov.RewardProposerUpdateInterval, gov.RewardMinimumStake, gov.IstanbulEpoch, gov.IstanbulPolicy,
			gov.IstanbulCommitteeSize, gov.GovernanceUnitPrice,
		}
		logger.Info("Using Mainnet/Kairos fallback")
	} else {
		// on private net, fallback candidates are all params.
		candidates = maps.Keys(gov.Params)
		logger.Info("Using default fallback")
	}

	for _, name := range candidates {
		param := gov.Params[name]
		value, err := param.ChainConfigValue(chainConfig)
		if err != nil {
			logger.Error("Failed to fetch value from ChainConfig", "name", name, "error", err)
			continue
		}

		if !reflect.DeepEqual(value, param.DefaultValue) {
			err := fallback.Add(string(name), value)
			if err != nil {
				logger.Error("Failed to add param to fallback", "name", name, "value", value, "error", err)
			}
		}
	}

	return fallback
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
