package impl

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/reward"
)

func (g *GovModule) APIs() []rpc.API {
	return append(g.hgm.APIs(), []rpc.API{
		{
			Namespace: "governance",
			Version:   "1.0",
			Service:   NewGovAPI(g),
			Public:    true,
		},
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewKaiaAPI(g),
			Public:    true,
		},
	}...)
}

type GovAPI struct {
	g *GovModule
}

type KaiaAPI struct {
	g *GovModule
}

type VotesAPI struct {
	BlockNum uint64
	Key      string
	Value    interface{}
}

type MyVotesAPI struct {
	BlockNum uint64
	Key      string
	Value    interface{}
	Casted   bool
}

type StatusAPI struct {
	GroupedVotes map[uint64]headergov.VotesInEpoch `json:"groupedVotes"`
	Governances  map[uint64]headergov.GovData      `json:"governances"`
	GovHistory   headergov.History                 `json:"govHistory"`
	NodeAddress  common.Address                    `json:"nodeAddress"`
	MyVotes      []headergov.VoteData              `json:"myVotes"`
}

type AccumulatedRewards struct {
	FirstBlockTime string   `json:"firstBlockTime"`
	LastBlockTime  string   `json:"lastBlockTime"`
	FirstBlock     *big.Int `json:"firstBlock"`
	LastBlock      *big.Int `json:"lastBlock"`

	// TotalMinted + TotalTxFee - TotalBurntTxFee = TotalProposerRewards + TotalStakingRewards + TotalKIFRewards + TotalKEFRewards
	TotalMinted          *big.Int                    `json:"totalMinted"`
	TotalTxFee           *big.Int                    `json:"totalTxFee"`
	TotalBurntTxFee      *big.Int                    `json:"totalBurntTxFee"`
	TotalProposerRewards *big.Int                    `json:"totalProposerRewards"`
	TotalStakingRewards  *big.Int                    `json:"totalStakingRewards"`
	TotalKIFRewards      *big.Int                    `json:"totalKIFRewards"`
	TotalKEFRewards      *big.Int                    `json:"totalKEFRewards"`
	Rewards              map[common.Address]*big.Int `json:"rewards"`
}

func NewGovAPI(g *GovModule) *GovAPI {
	return &GovAPI{g}
}

func (api *GovAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.g, num)
}

func (api *GovAPI) GetRewardsAccumulated(first rpc.BlockNumber, last rpc.BlockNumber) (*AccumulatedRewards, error) {
	currentBlock := api.g.chain.CurrentBlock().NumberU64()
	firstBlock := currentBlock
	if first >= rpc.EarliestBlockNumber {
		firstBlock = uint64(first.Int64())
	}

	lastBlock := currentBlock
	if last >= rpc.EarliestBlockNumber {
		lastBlock = uint64(last.Int64())
	}

	if firstBlock > lastBlock {
		return nil, errors.New("the last block number should be equal or larger the first block number")
	}

	if lastBlock > currentBlock {
		return nil, errors.New("the last block number should be equal or less than the current block number")
	}

	blockCount := lastBlock - firstBlock + 1
	if blockCount > 604800 { // 7 days. naive resource protection
		return nil, errors.New("block range should be equal or less than 604800")
	}
	// initialize structures before request a job
	accumRewards := &AccumulatedRewards{}
	blockRewards := reward.NewRewardSpec()
	mu := sync.Mutex{} // protect blockRewards

	numWorkers := runtime.NumCPU()
	reqCh := make(chan uint64, numWorkers)
	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}

	// introduce the worker pattern to prevent resource exhaustion
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// the minimum digit of request is period to avoid current access to an accArray item
			for num := range reqCh {
				bn := rpc.BlockNumber(num)
				kaiaApi := NewKaiaAPI(api.g)
				blockReward, err := kaiaApi.GetRewards(&bn)
				if err != nil {
					errCh <- err
					return
				}

				mu.Lock()
				blockRewards.Add(blockReward)
				mu.Unlock()
			}
		}()
	}

	// write the information of the first block
	header := api.g.chain.GetHeaderByNumber(firstBlock)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", firstBlock)
	}
	accumRewards.FirstBlock = header.Number
	accumRewards.FirstBlockTime = time.Unix(header.Time.Int64(), 0).String()

	// write the information of the last block
	header = api.g.chain.GetHeaderByNumber(lastBlock)
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", lastBlock)
	}
	accumRewards.LastBlock = header.Number
	accumRewards.LastBlockTime = time.Unix(header.Time.Int64(), 0).String()

	go func() {
		defer close(reqCh)
		for num := firstBlock; num <= lastBlock; num++ {
			reqCh <- num
		}
	}()

	// generate a goroutine to return error early
	go func() {
		wg.Wait()
		close(errCh)
	}()

	if err := <-errCh; err != nil {
		return nil, err
	}

	// collect the accumulated rewards information
	accumRewards.Rewards = blockRewards.Rewards
	accumRewards.TotalMinted = blockRewards.Minted
	accumRewards.TotalTxFee = blockRewards.TotalFee
	accumRewards.TotalBurntTxFee = blockRewards.BurntFee
	accumRewards.TotalProposerRewards = blockRewards.Proposer
	accumRewards.TotalStakingRewards = blockRewards.Stakers
	accumRewards.TotalKIFRewards = blockRewards.KIF
	accumRewards.TotalKEFRewards = blockRewards.KEF

	return accumRewards, nil
}

func NewKaiaAPI(g *GovModule) *KaiaAPI {
	return &KaiaAPI{g}
}

func (api *KaiaAPI) GetChainConfig(num *rpc.BlockNumber) *params.ChainConfig {
	return getChainConfig(api.g, num)
}

func (api *KaiaAPI) GetStakingInfo(num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	return getStakingInfo(api.g, num)
}

func (api *KaiaAPI) GetParams(num *rpc.BlockNumber) (map[string]interface{}, error) {
	return getParams(api.g, num)
}

func (api *KaiaAPI) GetRewards(num *rpc.BlockNumber) (*reward.RewardSpec, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = api.g.chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	// Check if the node has state to calculate the snapshot.
	err := checkStateForStakingInfo(api.g, blockNumber)
	if err != nil {
		return nil, err
	}

	header := api.g.chain.GetHeaderByNumber(blockNumber)
	block := api.g.chain.GetBlock(header.Hash(), blockNumber)
	if block == nil {
		return nil, errors.New("not found block")
	}
	txs, receipts := block.Transactions(), api.g.chain.GetReceiptsByBlockHash(header.Hash())
	if header == nil {
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}

	rules := api.g.chain.Config().Rules(new(big.Int).SetUint64(blockNumber))
	pset, err := api.g.EffectiveParamSet(blockNumber)
	if err != nil {
		return nil, err
	}
	rewardParamNum := reward.CalcRewardParamBlock(header.Number.Uint64(), pset.Epoch, rules)
	rewardParamSet, err := api.g.EffectiveParamSet(rewardParamNum)
	if err != nil {
		return nil, err
	}

	return reward.GetBlockReward(header, txs, receipts, rules, rewardParamSet.ToGovParamSet())
}

func getChainConfig(g *GovModule, num *rpc.BlockNumber) *params.ChainConfig {
	var blocknum uint64
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blocknum = g.chain.CurrentBlock().NumberU64()
	} else {
		blocknum = num.Uint64()
	}

	pset, err := g.EffectiveParamSet(blocknum)
	if err != nil {
		return nil
	}

	latestConfig := g.chain.Config()
	config := pset.ToGovParamSet().ToChainConfig()
	config.ChainID = latestConfig.ChainID
	config.IstanbulCompatibleBlock = latestConfig.IstanbulCompatibleBlock
	config.LondonCompatibleBlock = latestConfig.LondonCompatibleBlock
	config.EthTxTypeCompatibleBlock = latestConfig.EthTxTypeCompatibleBlock
	config.MagmaCompatibleBlock = latestConfig.MagmaCompatibleBlock
	config.KoreCompatibleBlock = latestConfig.KoreCompatibleBlock
	config.ShanghaiCompatibleBlock = latestConfig.ShanghaiCompatibleBlock
	config.CancunCompatibleBlock = latestConfig.CancunCompatibleBlock
	config.KaiaCompatibleBlock = latestConfig.KaiaCompatibleBlock
	config.Kip103CompatibleBlock = latestConfig.Kip103CompatibleBlock
	config.Kip103ContractAddress = latestConfig.Kip103ContractAddress
	config.Kip160CompatibleBlock = latestConfig.Kip160CompatibleBlock
	config.Kip160ContractAddress = latestConfig.Kip160ContractAddress
	config.RandaoCompatibleBlock = latestConfig.RandaoCompatibleBlock

	return config
}

func getStakingInfo(g *GovModule, num *rpc.BlockNumber) (*reward.StakingInfo, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = g.chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}
	// Check if the node has state to calculate the snapshot.
	err := checkStateForStakingInfo(g, blockNumber)
	if err != nil {
		return nil, err
	}

	return reward.GetStakingInfo(blockNumber), nil
}

// checkStateForStakingInfo checks the state of block for the given block number for staking info
func checkStateForStakingInfo(g *GovModule, blockNumber uint64) error {
	if blockNumber == 0 {
		return nil
	}

	// The staking info at blockNumber is calculated by the state of previous block
	blockNumber--
	if !g.chain.Config().IsKaiaForkEnabled(big.NewInt(int64(blockNumber + 1))) {
		return nil
	}
	header := g.chain.GetHeaderByNumber(blockNumber)
	if header == nil {
		return gov.ErrUnknownBlock
	}
	_, err := g.chain.StateAt(header.Root)
	return err
}

func getParams(g *GovModule, num *rpc.BlockNumber) (map[string]interface{}, error) {
	blockNumber := uint64(0)
	if num == nil || *num == rpc.LatestBlockNumber || *num == rpc.PendingBlockNumber {
		blockNumber = g.chain.CurrentBlock().NumberU64()
	} else {
		blockNumber = uint64(num.Int64())
	}

	gp, err := g.EffectiveParamSet(blockNumber)
	if err != nil {
		return nil, err
	}
	return gov.EnumMapToStrMap(gp.ToEnumMap()), nil
}

func (api *KaiaAPI) NodeAddress() common.Address {
	return api.g.hgm.NodeAddress()
}
