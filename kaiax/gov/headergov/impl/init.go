package impl

import (
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ headergov.HeaderGovModule = (*headerGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxGov)
)

type chain interface {
	GetHeaderByNumber(number uint64) *types.Header
	CurrentBlock() *types.Block
	State() (*state.StateDB, error)
}

type ValsetModule interface {
	GetCouncil(num uint64) ([]common.Address, error)
	GetProposer(num uint64, round uint64) (common.Address, error)
}

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	ValSet      ValsetModule
	NodeAddress common.Address
}

type headerGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	ValSet      ValsetModule

	groupedVotes headergov.GroupedVotesMap
	governances  headergov.GovDataMap
	history      headergov.History
	mu           *sync.RWMutex

	epoch uint64

	// for APIs
	nodeAddress common.Address
	myVotes     []headergov.VoteData // queue

	quit atomic.Int32 // stop the migration thread
	wg   sync.WaitGroup
}

func NewHeaderGovModule() *headerGovModule {
	return &headerGovModule{}
}

func (h *headerGovModule) Init(opts *InitOpts) error {
	if err := validateOpts(opts); err != nil {
		return err
	}

	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	h.ValSet = opts.ValSet
	h.nodeAddress = opts.NodeAddress

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return ErrZeroEpoch
	}

	// init memory
	h.groupedVotes = make(map[uint64]headergov.VotesInEpoch)
	h.governances = make(map[uint64]headergov.GovData)
	h.history = make(headergov.History)
	h.myVotes = make([]headergov.VoteData, 0)
	h.mu = &sync.RWMutex{}

	h.initSchema()

	govs := readGovDataFromDB(h.Chain, h.ChainKv)
	for blockNum, gov := range govs {
		h.AddGov(blockNum, gov)
	}

	votes := readVoteDataFromDB(h.Chain, h.ChainKv)
	for blockNum, vote := range votes {
		h.AddVote(blockNum, vote)
	}

	return nil
}

func (h *headerGovModule) initSchema() {
	// 1. Init gov from legacy governance DB, if empty.
	if ReadGovDataBlockNums(h.ChainKv) == nil {
		legacyGovBlockNums := ReadLegacyIdxHistory(h.ChainKv)
		if legacyGovBlockNums == nil {
			legacyGovBlockNums = StoredUint64Array{0}
		}
		WriteGovDataBlockNums(h.ChainKv, legacyGovBlockNums)
	}

	// 2. Init vote and lowest vote scanned epoch index, if empty.
	if ReadLowestVoteScannedEpochIdx(h.ChainKv) == nil {
		latestEpochIdx := calcEpochIdx(h.Chain.CurrentBlock().NumberU64(), h.epoch)
		h.accumulateVotesInEpoch(latestEpochIdx)
	}
}

func (h *headerGovModule) Start() error {
	logger.Info("HeaderGovModule started")

	// Reset the quit state.
	h.quit.Store(0)
	h.wg.Add(1)
	go h.migrate()

	return nil
}

func (h *headerGovModule) Stop() {
	logger.Info("HeaderGovModule stopped")
	h.quit.Store(1)
	h.wg.Wait()
}

func (h *headerGovModule) migrate() {
	defer h.wg.Done()

	// Scan all epochs in the background including 0th epoch
	pBorder := ReadLowestVoteScannedEpochIdx(h.ChainKv)
	if pBorder == nil {
		logger.Crit("Unexpected nil: lowest vote scanned epoch index")
		return
	}

	border := *pBorder
	for int64(border) > 0 {
		if h.quit.Load() == 1 {
			return
		}

		border -= 1
		h.accumulateVotesInEpoch(border)
	}

	if border == 0 {
		logger.Info("HeaderGovModule migrate complete")
	}
}

func (h *headerGovModule) isKoreHF(num uint64) bool {
	return h.ChainConfig.IsKoreForkEnabled(new(big.Int).SetUint64(num))
}

func (h *headerGovModule) PushMyVotes(vote headergov.VoteData) {
	h.myVotes = append(h.myVotes, vote)
}

func (h *headerGovModule) PopMyVotes(idx int) {
	h.myVotes = append(h.myVotes[:idx], h.myVotes[idx+1:]...)
}

// scanAllVotesInEpoch scans all votes from headers in the given epoch.
func (h *headerGovModule) scanAllVotesInEpoch(epochIdx uint64) map[uint64]headergov.VoteData {
	rangeStart := calcEpochStartBlock(epochIdx, h.epoch)
	rangeEnd := calcEpochStartBlock(epochIdx+1, h.epoch)
	if currBlock := h.Chain.CurrentBlock().NumberU64(); rangeEnd > currBlock {
		rangeEnd = currBlock + 1 // so currBlock is scanned
	}

	votes := make(map[uint64]headergov.VoteData)
	for blockNum := rangeStart; blockNum < rangeEnd; blockNum++ {
		header := h.Chain.GetHeaderByNumber(blockNum)
		if header == nil {
			logger.Warn("Missing header found", "num", blockNum)
			continue
		}

		if len(header.Vote) == 0 {
			continue
		}

		vote, err := headergov.VoteBytes(header.Vote).ToVoteData()
		if err != nil {
			logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			continue
		}

		if vote == nil {
			continue
		}

		// Only governance params are collected. Validator params are ignored.
		if _, ok := gov.Params[vote.Name()]; ok {
			votes[blockNum] = vote
		}
	}

	return votes
}

// accumulateVotesInEpoch scans and saves votes to cache and DB.
// Because this function updates lowestVoteScannedEpochIdx, it verifies epochIdx.
func (h *headerGovModule) accumulateVotesInEpoch(epochIdx uint64) {
	pBorder := ReadLowestVoteScannedEpochIdx(h.ChainKv)
	if pBorder != nil && *pBorder != epochIdx+1 {
		logger.Error("Invalid epochIdx", "epochIdx", epochIdx, "lowestScanned", *pBorder)
		return
	}

	votes := h.scanAllVotesInEpoch(epochIdx)
	for blockNum, vote := range votes {
		h.AddVote(blockNum, vote)
		InsertVoteDataBlockNum(h.ChainKv, blockNum)
	}

	WriteLowestVoteScannedEpochIdx(h.ChainKv, epochIdx)
	logger.Info("Accumulated votes", "epochIdx", epochIdx)
}

func readVoteDataFromDB(chain chain, db database.Database) map[uint64]headergov.VoteData {
	voteBlocks := ReadVoteDataBlockNums(db)
	votes := make(map[uint64]headergov.VoteData)
	for _, blockNum := range voteBlocks {
		header := chain.GetHeaderByNumber(blockNum)
		parsedVote, err := headergov.VoteBytes(header.Vote).ToVoteData()
		if err != nil {
			panic(err)
		}

		votes[blockNum] = parsedVote
	}

	return votes
}

func readGovDataFromDB(chain chain, db database.Database) map[uint64]headergov.GovData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make(map[uint64]headergov.GovData)

	for _, blockNum := range govBlocks {
		header := chain.GetHeaderByNumber(blockNum)

		parsedGov, err := headergov.GovBytes(header.Governance).ToGovData()
		if err != nil {
			// For tests, genesis' governance can be nil.
			if blockNum == 0 {
				continue
			}

			logger.Error("Failed to parse gov", "num", blockNum, "err", err)
			panic("failed to parse gov")
		}

		govs[blockNum] = parsedGov
	}

	return govs
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}

func calcEpochStartBlock(epochIdx uint64, epoch uint64) uint64 {
	return epochIdx * epoch
}

func validateOpts(opts *InitOpts) error {
	switch {
	case opts == nil:
		return errInitNil("opts")
	case opts.ChainKv == nil:
		return errInitNil("opts.ChainKv")
	case opts.ChainConfig == nil:
		return errInitNil("opts.ChainConfig")
	case opts.ChainConfig.Istanbul == nil:
		return errInitNil("opts.ChainConfig.Istanbul")
	case opts.Chain == nil:
		return errInitNil("opts.Chain")
	}

	return nil
}

// GetGenesisGovBytes returns the genesis governance bytes for initGenesis().
// It panics if the given chain config is invalid.
func GetGenesisGovBytes(config *params.ChainConfig) headergov.GovBytes {
	partialParamSet := make(gov.PartialParamSet)
	genesisParamNames := getGenesisParamNames(config)
	for _, name := range genesisParamNames {
		param := gov.Params[name]
		value, err := param.ChainConfigValue(config)
		if err != nil {
			panic(err)
		}

		err = partialParamSet.Add(string(name), value)
		if err != nil {
			panic(fmt.Errorf("failed to add param %s: %w", name, err))
		}
	}

	govData := headergov.NewGovData(partialParamSet)
	ret, err := govData.ToGovBytes()
	if err != nil {
		panic(err)
	}

	return ret
}

// getGenesisParamNames returns a subset of parameters to be included
// in the genesis block header's Governance field.
// The subset depends on genesis block's hardfork level and genesis chain config.
// For backward compatibility, some fields may be intentionally omitted.
func getGenesisParamNames(config *params.ChainConfig) []gov.ParamName {
	genesisParamNames := make([]gov.ParamName, 0)

	// DeriveShaImpl is excluded unless the genesis is Kore HF,
	// even though it has been always included in the genesis chain config.
	if config.Governance != nil {
		genesisParamNames = append(genesisParamNames, []gov.ParamName{
			gov.GovernanceGovernanceMode, gov.GovernanceGoverningNode, gov.GovernanceUnitPrice,
			gov.RewardMintingAmount, gov.RewardRatio, gov.RewardUseGiniCoeff,
			gov.RewardDeferredTxFee, gov.RewardMinimumStake,
			gov.RewardStakingUpdateInterval, gov.RewardProposerUpdateInterval,
		}...)
	}

	if config.Istanbul != nil {
		genesisParamNames = append(genesisParamNames, []gov.ParamName{
			gov.IstanbulEpoch, gov.IstanbulPolicy, gov.IstanbulCommitteeSize,
		}...)
	}

	if config.IsMagmaForkEnabled(common.Big0) &&
		config.Governance.KIP71 != nil {
		genesisParamNames = append(genesisParamNames, []gov.ParamName{
			gov.Kip71LowerBoundBaseFee, gov.Kip71UpperBoundBaseFee, gov.Kip71GasTarget,
			gov.Kip71BaseFeeDenominator, gov.Kip71MaxBlockGasUsedForBaseFee,
		}...)
	}

	if config.IsKoreForkEnabled(common.Big0) &&
		config.Governance != nil {
		genesisParamNames = append(genesisParamNames, gov.GovernanceDeriveShaImpl)
		if !common.EmptyAddress(config.Governance.GovParamContract) {
			genesisParamNames = append(genesisParamNames, gov.GovernanceGovParamContract)
		}
		if config.Governance.Reward != nil &&
			config.Governance.Reward.Kip82Ratio != "" {
			genesisParamNames = append(genesisParamNames, gov.RewardKip82Ratio)
		}
	}

	return genesisParamNames
}
