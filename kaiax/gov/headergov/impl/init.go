package impl

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	_ headergov.HeaderGovModule = (*headerGovModule)(nil)

	logger = log.NewModuleLogger(log.KaiaXGov)
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

//go:generate mockgen -destination=kaiax/headergov/mocks/headergov_mock.go github.com/kaiachain/kaia/kaiax/headergov HeaderGovModule
type headerGovModule struct {
	ChainKv     database.Database
	ChainConfig *params.ChainConfig
	Chain       chain
	nodeAddress common.Address
	myVotes     []headergov.VoteData // queue

	epoch uint64
	cache *headergov.HeaderCache
}

func NewHeaderGovModule() *headerGovModule {
	return &headerGovModule{}
}

func (h *headerGovModule) Init(opts *InitOpts) error {
	h.ChainKv = opts.ChainKv
	h.ChainConfig = opts.ChainConfig
	h.Chain = opts.Chain
	h.nodeAddress = opts.NodeAddress
	h.myVotes = make([]headergov.VoteData, 0)
	if h.ChainKv == nil || h.ChainConfig == nil || h.ChainConfig.Istanbul == nil || h.Chain == nil {
		return ErrInitNil
	}

	h.epoch = h.ChainConfig.Istanbul.Epoch
	if h.epoch == 0 {
		return ErrZeroEpoch
	}

	votes := readVoteDataFromDB(h.Chain, h.ChainKv)
	govs := readGovDataFromDB(h.Chain, h.ChainKv)

	h.cache = headergov.NewHeaderGovCache()
	for blockNum, vote := range votes {
		h.cache.AddVote(calcEpochIdx(blockNum, h.epoch), blockNum, vote)
	}
	for blockNum, gov := range govs {
		h.cache.AddGov(blockNum, gov)
	}

	return nil
}

func (s *headerGovModule) Start() error {
	logger.Info("HeaderGovModule started")
	return nil
}

func (s *headerGovModule) Stop() {
	logger.Info("HeaderGovModule stopped")
}

func (s *headerGovModule) isKoreHF(num uint64) bool {
	return s.ChainConfig.IsKoreForkEnabled(new(big.Int).SetUint64(num))
}

func (s *headerGovModule) PushMyVotes(vote headergov.VoteData) {
	s.myVotes = append(s.myVotes, vote)
}

func (s *headerGovModule) PopMyVotes(idx int) {
	s.myVotes = append(s.myVotes[:idx], s.myVotes[idx+1:]...)
}

func readVoteDataFromDB(chain chain, db database.Database) map[uint64]headergov.VoteData {
	voteBlocks := ReadVoteDataBlockNums(db)
	votes := make(map[uint64]headergov.VoteData)
	if voteBlocks != nil {
		for _, blockNum := range *voteBlocks {
			header := chain.GetHeaderByNumber(blockNum)
			parsedVote, err := headergov.DeserializeHeaderVote(header.Vote)
			if err != nil {
				panic(err)
			}

			votes[blockNum] = parsedVote
		}
	}

	return votes
}

func readGovDataFromDB(chain chain, db database.Database) map[uint64]headergov.GovData {
	govBlocks := ReadGovDataBlockNums(db)
	govs := make(map[uint64]headergov.GovData)

	// TODO: in production, govBlocks must not be nil. Remove this after implementing kcn init and data migration.
	if govBlocks == nil {
		govBlocks = &StoredUint64Array{0}
	}

	for _, blockNum := range *govBlocks {
		header := chain.GetHeaderByNumber(blockNum)

		parsedGov, err := headergov.DeserializeHeaderGov(header.Governance)
		if err != nil {
			logger.Error("Failed to parse vote", "num", blockNum, "err", err)
			continue
		}

		govs[blockNum] = parsedGov
	}

	return govs
}

func calcEpochIdx(blockNum uint64, epoch uint64) uint64 {
	return blockNum / epoch
}
