// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/gasprice/gasprice.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package gasprice

import (
	"context"
	"math/big"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/misc"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"golang.org/x/exp/slices"
)

const sampleNumber = 3 // Number of transactions sampled in a block

type Config struct {
	Blocks           int
	Percentile       int
	MaxHeaderHistory int
	MaxBlockHistory  int
	MaxPrice         *big.Int `toml:",omitempty"`
}

// OracleBackend includes all necessary background APIs for oracle.
type OracleBackend interface {
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	GetBlockReceipts(ctx context.Context, hash common.Hash) types.Receipts
	ChainConfig() *params.ChainConfig
	CurrentBlock() *types.Block
}

type TxPool interface {
	GasPrice() *big.Int
}

// Oracle recommends gas prices based on the content of recent
// blocks. Suitable for both light and full clients.
type Oracle struct {
	backend   OracleBackend
	lastHead  common.Hash
	lastPrice *big.Int
	maxPrice  *big.Int
	cacheLock sync.RWMutex
	fetchLock sync.Mutex
	txPool    TxPool
	govModule gov.GovModule

	checkBlocks, maxEmpty, maxBlocks  int
	percentile                        int
	maxHeaderHistory, maxBlockHistory int

	historyCache *lru.Cache
}

// NewOracle returns a new oracle.
func NewOracle(backend OracleBackend, config Config, txPool TxPool, govModule gov.GovModule) *Oracle {
	blocks := config.Blocks
	if blocks < 1 {
		blocks = 1
	}
	percent := config.Percentile
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	maxPrice := config.MaxPrice
	if maxPrice == nil || maxPrice.Int64() <= 0 {
		maxPrice = big.NewInt(params.DefaultGPOMaxPrice)
		logger.Warn("Sanitizing invalid gasprice oracle price cap", "provided", config.MaxPrice, "updated", maxPrice)
	}
	maxHeaderHistory := config.MaxHeaderHistory
	if maxHeaderHistory < 1 {
		maxHeaderHistory = 1
		logger.Warn("Sanitizing invalid gasprice oracle max header history", "provided", config.MaxHeaderHistory, "updated", maxHeaderHistory)
	}
	maxBlockHistory := config.MaxBlockHistory
	if maxBlockHistory < 1 {
		maxBlockHistory = 1
		logger.Warn("Sanitizing invalid gasprice oracle max block history", "provided", config.MaxBlockHistory, "updated", maxBlockHistory)
	}
	cache, _ := lru.New(2048)

	return &Oracle{
		backend:          backend,
		lastPrice:        common.Big0,
		maxPrice:         maxPrice,
		checkBlocks:      blocks,
		maxEmpty:         blocks / 2,
		maxBlocks:        blocks * 5,
		percentile:       percent,
		maxHeaderHistory: maxHeaderHistory,
		maxBlockHistory:  maxBlockHistory,
		txPool:           txPool,
		govModule:        govModule,
		historyCache:     cache,
	}
}

// Tx gas price requirements has changed over the hardforks
//
// | Fork              | gasPrice                     | maxFeePerGas                 | maxPriorityFeePerGas         |
// |------------------ |----------------------------- |----------------------------- |----------------------------- |
// | Before EthTxType  | must be fixed UnitPrice (1)  | N/A (2)                      | N/A (2)                      |
// | After EthTxType   | must be fixed UnitPrice (1)  | must be fixed UnitPrice (3)  | must be fixed UnitPrice (3)  |
// | After Magma       | BaseFee or higher (4)        | BaseFee or higher (4)        | Ignored (4)                  |
//
// (1) If tx.type != 2 && !rules.IsMagma: https://github.com/kaiachain/kaia/blob/v1.11.1/blockchain/tx_pool.go#L729
// (2) If tx.type == 2 && !rules.IsEthTxType: https://github.com/kaiachain/kaia/blob/v1.11.1/blockchain/tx_pool.go#L670
// (3) If tx.type == 2 && !rules.IsMagma: https://github.com/kaiachain/kaia/blob/v1.11.1/blockchain/tx_pool.go#L710
// (4) If tx.type == 2 && rules.IsMagma: https://github.com/kaiachain/kaia/blob/v1.11.1/blockchain/tx_pool.go#L703
//
// The suggested prices needs to match the requirements.
//
// | Fork              | SuggestPrice (for gasPrice and maxFeePerGas)                | SuggestTipCap (for maxPriorityFeePerGas)                                              |
// |------------------ |------------------------------------------------------------ |-------------------------------------------------------------------------------------- |
// | Before Magma      | Fixed UnitPrice                                             | Fixed UnitPrice                                                                       |
// | After Magma       | BaseFee * 2                                                 | Zero                                                                                  |
// | After Kaia        | BaseFee * 1.10 or 1.15 + SuggestTipCap                      | Zero if nextBaseFee is lower bound, 60% percentile of last 20 blocks otherwise.       |

// SuggestPrice returns the recommended gas price.
// This value is intended to be used as gasPrice or maxFeePerGas.
func (gpo *Oracle) SuggestPrice(ctx context.Context) (*big.Int, error) {
	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}

	nextNum := new(big.Int).Add(gpo.backend.CurrentBlock().Number(), common.Big1)
	if gpo.backend.ChainConfig().IsKaiaForkEnabled(nextNum) {
		// After Kaia, include suggested tip
		baseFee := gpo.txPool.GasPrice()
		suggestedTip, err := gpo.SuggestTipCap(ctx)
		if err != nil {
			return nil, err
		}
		// If network is relaxed, give a buffer of 10% to the suggested tip. Otherwise, 15%.
		if suggestedTip.Cmp(common.Big0) == 0 {
			baseFee.Mul(baseFee, big.NewInt(110))
		} else {
			baseFee.Mul(baseFee, big.NewInt(115))
		}
		baseFee.Div(baseFee, big.NewInt(100))
		return new(big.Int).Add(baseFee, suggestedTip), nil
	} else if gpo.backend.ChainConfig().IsMagmaForkEnabled(nextNum) {
		// After Magma, return the twice of BaseFee as a buffer.
		baseFee := gpo.txPool.GasPrice()
		return new(big.Int).Mul(baseFee, common.Big2), nil
	} else {
		// Before Magma, return the fixed UnitPrice.
		unitPrice := gpo.txPool.GasPrice()
		return unitPrice, nil
	}
}

// SuggestTipCap returns the recommended gas tip cap.
// This value is intended to be used as maxPriorityFeePerGas.
func (gpo *Oracle) SuggestTipCap(ctx context.Context) (*big.Int, error) {
	if gpo.txPool == nil {
		// If txpool is not set, just return 0. This is used for testing.
		return common.Big0, nil
	}

	nextNum := new(big.Int).Add(gpo.backend.CurrentBlock().Number(), common.Big1)
	if gpo.backend.ChainConfig().IsKaiaForkEnabled(nextNum) {
		// After Kaia, return using fee history.
		// If the next baseFee is lower bound, return 0.
		// Otherwise, by default config, this will return 60% percentile of last 20 blocks.
		// See node/cn/config.go for the default config.
		header := gpo.backend.CurrentBlock().Header()
		headHash := header.Hash()
		// If the latest gasprice is still available, return it.
		if lastPrice, ok := gpo.readCacheChecked(headHash); ok {
			return new(big.Int).Set(lastPrice), nil
		}
		if gpo.isRelaxedNetwork(header) {
			gpo.writeCache(headHash, common.Big0)
			return common.Big0, nil
		}
		price, err := gpo.suggestTipCapUsingFeeHistory(ctx)
		if err == nil {
			gpo.writeCache(headHash, price)
		}
		return price, err
	} else if gpo.backend.ChainConfig().IsMagmaForkEnabled(nextNum) {
		// After Magma, return zero
		return common.Big0, nil
	} else {
		// Before Magma, return the fixed UnitPrice.
		unitPrice := gpo.txPool.GasPrice()
		return unitPrice, nil
	}
}

// suggestTipCapUsingFeeHistory returns a tip cap based on fee history.
func (oracle *Oracle) suggestTipCapUsingFeeHistory(ctx context.Context) (*big.Int, error) {
	head, _ := oracle.backend.HeaderByNumber(ctx, rpc.LatestBlockNumber)
	headHash := head.Hash()

	oracle.fetchLock.Lock()
	defer oracle.fetchLock.Unlock()

	lastHead, lastPrice := oracle.readCache()
	if headHash == lastHead {
		return new(big.Int).Set(lastPrice), nil
	}
	var (
		sent, exp int
		number    = head.Number.Uint64()
		result    = make(chan results, oracle.checkBlocks)
		quit      = make(chan struct{})
		results   []*big.Int
	)
	for sent < oracle.checkBlocks && number > 0 {
		go oracle.getBlockValues(ctx, number, sampleNumber, result, quit)
		sent++
		exp++
		number--
	}
	for exp > 0 {
		res := <-result
		if res.err != nil {
			close(quit)
			return new(big.Int).Set(lastPrice), res.err
		}
		exp--
		// Nothing returned. There are two special cases here:
		// - The block is empty
		// - All the transactions included are sent by the miner itself.
		// In these cases, use the latest calculated price for sampling.
		if len(res.values) == 0 {
			res.values = []*big.Int{lastPrice}
		}
		// Besides, in order to collect enough data for sampling, if nothing
		// meaningful returned, try to query more blocks. But the maximum
		// is 2*checkBlocks.
		if len(res.values) == 1 && len(results)+1+exp < oracle.checkBlocks*2 && number > 0 {
			go oracle.getBlockValues(ctx, number, sampleNumber, result, quit)
			sent++
			exp++
			number--
		}
		results = append(results, res.values...)
	}
	price := lastPrice
	if len(results) > 0 {
		slices.SortFunc(results, func(a, b *big.Int) int { return a.Cmp(b) })
		price = results[(len(results)-1)*oracle.percentile/100]
	}
	// NOTE: This maximum suggested gas tip can lead to suggesting insufficient gas tip,
	//       however, the possibility of gas tip exceeding 500 gkei would be very low given the block capacity of Kaia.
	//       On the other hand, referencing the user-submitted transactions as-is can lead to suggesting
	//       very high gas tip when there are only a few transactions with unnecessarily high gas tip.
	if price.Cmp(oracle.maxPrice) > 0 {
		price = new(big.Int).Set(oracle.maxPrice)
	}

	return new(big.Int).Set(price), nil
}

type results struct {
	values []*big.Int
	err    error
}

// getBlockValues calculates the specified number of lowest transaction gas tips
// in a given block and sends them to the result channel. If a transaction was
// sent by the miner itself(it doesn't make any sense to include this kind of
// transaction prices for sampling), nil gasprice is returned.
func (oracle *Oracle) getBlockValues(ctx context.Context, blockNum uint64, limit int, result chan results, quit chan struct{}) {
	block, err := oracle.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNum))
	if block == nil {
		select {
		case result <- results{nil, err}:
		case <-quit:
		}
		return
	}
	signer := types.MakeSigner(oracle.backend.ChainConfig(), block.Number())

	// Sort the transaction by effective tip in ascending sort.
	txs := block.Transactions()
	sortedTxs := make([]*types.Transaction, len(txs))
	copy(sortedTxs, txs)
	baseFee := block.Header().BaseFee
	slices.SortFunc(sortedTxs, func(a, b *types.Transaction) int {
		tip1 := a.EffectiveGasTip(baseFee)
		tip2 := b.EffectiveGasTip(baseFee)
		return tip1.Cmp(tip2)
	})

	var prices []*big.Int
	for _, tx := range sortedTxs {
		tip := tx.EffectiveGasTip(baseFee)
		sender, err := types.Sender(signer, tx)
		if err == nil && sender != block.Rewardbase() {
			prices = append(prices, tip)
			if len(prices) >= limit {
				break
			}
		}
	}
	select {
	case result <- results{prices, nil}:
	case <-quit:
	}
}

// isRelaxedNetwork returns true if the current network congestion is low to the point
// paying any tip is unnecessary. It returns true when the head block is after Magma fork
// and the next base fee is at the lower bound.
func (oracle *Oracle) isRelaxedNetwork(header *types.Header) bool {
	pset := oracle.govModule.EffectiveParamSet(header.Number.Uint64() + 1)
	kip71 := pset.ToKip71Config()
	nextBaseFee := misc.NextMagmaBlockBaseFee(header, kip71)
	return nextBaseFee.Cmp(big.NewInt(int64(pset.LowerBoundBaseFee))) <= 0
}

func (oracle *Oracle) readCacheChecked(headHash common.Hash) (*big.Int, bool) {
	lastHead, lastPrice := oracle.readCache()
	return lastPrice, headHash == lastHead
}

func (oracle *Oracle) readCache() (common.Hash, *big.Int) {
	oracle.cacheLock.RLock()
	defer oracle.cacheLock.RUnlock()
	return oracle.lastHead, oracle.lastPrice
}

func (oracle *Oracle) writeCache(head common.Hash, price *big.Int) {
	oracle.cacheLock.Lock()
	oracle.lastHead = head
	oracle.lastPrice = price
	oracle.cacheLock.Unlock()
}

func (oracle *Oracle) PurgeCache() {
	oracle.cacheLock.Lock()
	oracle.historyCache.Purge()
	oracle.cacheLock.Unlock()
}
