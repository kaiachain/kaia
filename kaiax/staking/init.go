package staking

import (
	"errors"
	"fmt"
	"math/big"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	staking_types "github.com/kaiachain/kaia/kaiax/staking/types"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

type StakingInfo = staking_types.StakingInfo

var (
	_ (staking_types.StakingModule) = (*StakingModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXStaking)

	errZeroStakingInterval = errors.New("staking interval cannot be zero")
	errInvalidABookResult  = errors.New("invalid result from an AddressBook call")
)

func errCannotCallABook(inner error) error {
	return fmt.Errorf("failed to make an AddressBook call: %w", inner)
}

type chain interface {
	backends.BlockChainForCaller
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
}

type StakingModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain

	stakingInterval   uint64
	cachedStakingInfo *lru.ARCCache
}

func NewStakingModule() *StakingModule {
	cache, _ := lru.NewARC(128)
	return &StakingModule{
		cachedStakingInfo: cache,
	}
}

func (s *StakingModule) Init(opts *InitOpts) error {
	s.ChainKv = opts.ChainKv
	s.ChainConfig = opts.ChainConfig
	s.Chain = opts.Chain

	// StakingInterval is first determined by the Genesis config, then never changes.
	s.stakingInterval = opts.ChainConfig.Governance.Reward.StakingUpdateInterval
	if s.stakingInterval == 0 {
		return errZeroStakingInterval
	}
	return nil
}

func (s *StakingModule) Start() error {
	logger.Info("StakingModule started")
	return nil
}

func (s *StakingModule) Stop() {
	logger.Info("StakingModule stopped")
}

func (s *StakingModule) isKaia(num uint64) bool {
	return s.ChainConfig.IsKaiaForkEnabled(new(big.Int).SetUint64(num))
}
