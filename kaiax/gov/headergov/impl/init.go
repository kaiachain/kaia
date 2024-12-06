package impl

import (
	"math/big"
	"sync"

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

type InitOpts struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	NodeAddress common.Address
}

//go:generate mockgen -destination=mocks/headergov_mock.go github.com/kaiachain/kaia/kaiax/headergov HeaderGovModule
type headerGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain

	groupedVotes headergov.GroupedVotesMap
	governances  headergov.GovDataMap
	history      headergov.History
	mu           *sync.RWMutex

	epoch uint64

	// for APIs
	nodeAddress common.Address
	myVotes     []headergov.VoteData // queue
}

func NewHeaderGovModule() *headerGovModule {
	return &headerGovModule{}
}

func (h *headerGovModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainKv == nil || opts.ChainConfig == nil || opts.ChainConfig.Istanbul == nil || opts.Chain == nil {
		return ErrInitNil
	}

	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	h.nodeAddress = opts.NodeAddress
	h.myVotes = make([]headergov.VoteData, 0)
	h.mu = &sync.RWMutex{}

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return ErrZeroEpoch
	}

	// 1. Init gov. If Gov DB is empty, migrate from legacy governance DB.
	if ReadGovDataBlockNums(h.ChainKv) == nil {
		legacyGovBlockNums := ReadLegacyIdxHistory(h.ChainKv)
		WriteGovDataBlockNums(h.ChainKv, legacyGovBlockNums)
	}

	h.groupedVotes = make(map[uint64]headergov.VotesInEpoch)
	h.governances = make(map[uint64]headergov.GovData)
	govs := readGovDataFromDB(h.Chain, h.ChainKv)
	h.history = make(headergov.History)
	for blockNum, gov := range govs {
		h.AddGov(blockNum, gov)
	}

	// 2. Init votes. If votes exist in the latest epoch, read from DB. Otherwise, accumulate.
	lowestVoteScannedBlockNumPtr := ReadLowestVoteScannedBlockNum(h.ChainKv)
	if lowestVoteScannedBlockNumPtr != nil {
		votes := readVoteDataFromDB(h.Chain, h.ChainKv)

		for blockNum, vote := range votes {
			h.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)
		}
	} else {
		latestEpochIdx := calcEpochIdx(h.Chain.CurrentBlock().NumberU64(), h.epoch)
		h.accumulateVotesInEpoch(latestEpochIdx)
	}

	return nil
}

func (h *headerGovModule) Start() error {
	logger.Info("HeaderGovModule started")

	lowestVoteScannedBlockNumPtr := ReadLowestVoteScannedBlockNum(h.ChainKv)
	if lowestVoteScannedBlockNumPtr == nil {
		return ErrLowestVoteScannedBlockNotFound
	}

	// Scan all epochs in the background including 0th epoch
	epochIdxIter := calcEpochIdx(*lowestVoteScannedBlockNumPtr, h.epoch)
	go func() {
		for int64(epochIdxIter) > 0 {
			epochIdxIter -= 1
			h.accumulateVotesInEpoch(epochIdxIter)
		}
	}()

	return nil
}

func (h *headerGovModule) Stop() {
	logger.Info("HeaderGovModule stopped")
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

	votes := make(map[uint64]headergov.VoteData)
	for blockNum := rangeStart; blockNum < rangeEnd; blockNum++ {
		header := h.Chain.GetHeaderByNumber(blockNum)
		if header == nil || len(header.Vote) == 0 {
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
// Because this function updates lowestVoteScannedBlockNum, it verifies epochIdx.
func (h *headerGovModule) accumulateVotesInEpoch(epochIdx uint64) {
	lowestVoteScannedBlockNumPtr := ReadLowestVoteScannedBlockNum(h.ChainKv)

	// assert epochIdx == lowestVoteScannedBlockNum - 1
	if lowestVoteScannedBlockNumPtr != nil && *lowestVoteScannedBlockNumPtr != epochIdx+1 {
		logger.Error("Invalid epochIdx", "epochIdx", epochIdx, "lowestScanned", *lowestVoteScannedBlockNumPtr)
		return
	}

	votes := h.scanAllVotesInEpoch(epochIdx)
	for blockNum, vote := range votes {
		h.AddVote(epochIdx, blockNum, vote)
		InsertVoteDataBlockNum(h.ChainKv, blockNum)
	}

	WriteLowestVoteScannedBlockNum(h.ChainKv, epochIdx)
	logger.Info("Accumulated votes", "epochIdx", epochIdx, "lowestScanned", epochIdx)
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

	// TODO: in production, govBlocks must not be nil. Remove this after implementing kcn init and data migration.
	if govBlocks == nil {
		govBlocks = StoredUint64Array{0}
	}

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
