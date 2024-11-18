// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package reward

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/consensus"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
)

const (
	chainHeadChanSize = 100
)

// addressType defined in AddressBook
const (
	addressTypeNodeID = iota
	addressTypeStakingAddr
	addressTypeRewardAddr
	addressTypePoCAddr // TODO-Kaia: PoC should be changed to KIF after changing AddressBook contract
	addressTypeKIRAddr // TODO-Kaia: KIR should be changed to KEF after changing AddressBook contract
)

var (
	addressBookContractAddress = system.AddressBookAddr

	logger = log.NewModuleLogger(log.Reward)
)

// blockChain is an interface for blockchain.Blockchain used in reward package.
type blockChain interface {
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	GetBlockByNumber(number uint64) *types.Block
	GetReceiptsByBlockHash(hash common.Hash) types.Receipts
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig
	CurrentHeader() *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetHeaderByNumber(number uint64) *types.Header
	State() (*state.StateDB, error)
	CurrentBlock() *types.Block

	StateCache() state.Database
	Processor() blockchain.Processor

	blockchain.ChainContext
}

// Cannot use governance.Engine because of cyclic dependency.
// Instead declare only the methods used by this package.
type governanceHelper interface {
	CurrentParams() *params.GovParamSet
	EffectiveParams(num uint64) (*params.GovParamSet, error)
}

type StakingManager struct {
	stakingInfoCache *lru.ARCCache
	stakingInfoDB    stakingInfoDB
	governanceHelper governanceHelper
	blockchain       blockChain
	chainHeadChan    chan blockchain.ChainHeadEvent
	chainHeadSub     event.Subscription

	// Preloaded stakingInfos are fetched before the GetStakingInfo request.
	// This is used when the state is available when preloaded, but not available
	// when GetStakingInfo is called. e.g. reexec loop in stateAtBlock.
	// Therefore preloaded staking infos must not be evicted,
	// and it should be only used temporarily, hence a separate mapping.
	preloadedInfo *common.RefCountingMap
}

var (
	// variables for sole StakingManager
	once           sync.Once
	stakingManager *StakingManager

	// errors for staking manager
	ErrStakingManagerNotSet = errors.New("staking manager is not set")
	ErrChainHeadChanNotSet  = errors.New("chain head channel is not set")
)

// NewStakingManager creates and returns StakingManager.
//
// On the first call, a StakingManager is created with given parameters.
// From next calls, the existing StakingManager is returned. (Parameters
// from the next calls will not affect.)
func NewStakingManager(bc blockChain, gh governanceHelper, db stakingInfoDB) *StakingManager {
	if bc != nil && gh != nil {
		// this is only called once
		once.Do(func() {
			cache, _ := lru.NewARC(128)
			stakingManager = &StakingManager{
				stakingInfoCache: cache,
				stakingInfoDB:    db,
				governanceHelper: gh,
				blockchain:       bc,
				chainHeadChan:    make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
				preloadedInfo:    common.NewRefCountingMap(),
			}

			// Before migration, staking information of current and before should be stored in DB.
			//
			// Staking information from block of StakingUpdateInterval ahead is needed to create a block.
			// If there is no staking info in either cache, db or state trie, the node cannot make a block.
			// The information in state trie is deleted after state trie migration.
			blockchain.RegisterMigrationPrerequisites(func(blockNum uint64) error {
				// Don't need to check if staking info is stored after kaia fork.
				if isKaiaForkEnabled(blockNum) {
					return nil
				}
				if err := checkStakingInfoStored(blockNum); err != nil {
					return err
				}
				return checkStakingInfoStored(blockNum + params.StakingUpdateInterval())
			})
		})
	} else {
		logger.Error("unable to set StakingManager", "blockchain", bc, "governanceHelper", gh)
	}

	return stakingManager
}

func GetStakingManager() *StakingManager {
	return stakingManager
}

// GetStakingInfo returns a stakingInfo on the staking block of the given block number.
// Note that staking block is the block on which the associated staking information is stored and used during an interval.
// - Before kaia fork: staking block is calculated by params.CalcStakingBlockNumber(blockNum)
// - After kaia fork: staking block is the previous block of the given block number.
func GetStakingInfo(blockNum uint64) *StakingInfo {
	stakingBlockNumber := blockNum
	var stakingInfo *StakingInfo
	if isKaiaForkEnabled(blockNum) {
		if blockNum > 0 {
			stakingBlockNumber--
		}
		stakingInfo = GetStakingInfoForKaiaBlock(stakingBlockNumber)
	} else {
		stakingBlockNumber = params.CalcStakingBlockNumber(blockNum)
		stakingInfo = GetStakingInfoOnStakingBlock(stakingBlockNumber)
	}

	logger.Debug("Staking information is requested", "blockNum", blockNum, "staking block number", stakingBlockNumber)
	return stakingInfo
}

// GetStakingInfoForKaiaBlock returns a corresponding kaia StakingInfo for a given block number.
// Note that the given block number is a kaia staking info for the next block.
func GetStakingInfoForKaiaBlock(blockNum uint64) *StakingInfo {
	if stakingManager == nil {
		logger.Error("unable to GetStakingInfo", "err", ErrStakingManagerNotSet)
		return nil
	}

	// Check if the next block is a kaia block.
	if !isKaiaForkEnabled(blockNum + 1) {
		logger.Error("invalid block number for kaia staking info", "block number", blockNum)
		return nil
	}

	// Get staking info from cache
	if cachedStakingInfo := getStakingInfoFromCache(blockNum); cachedStakingInfo != nil {
		return cachedStakingInfo
	}

	stakingInfo, err := updateKaiaStakingInfo(blockNum)
	if err != nil {
		logger.Error("failed to update kaia stakingInfo", "block number", blockNum, "err", err)
		return nil
	}

	logger.Debug("Get stakingInfo from header.", "block number", blockNum, "stakingInfo", stakingInfo)
	return stakingInfo
}

// GetStakingInfoOnStakingBlock returns a corresponding StakingInfo for a staking block number.
// If the given number is not on the staking block, it returns nil.
//
// Fixup for Gini coefficients:
// Kaia core stores Gini: -1 in its database.
// We ensure GetStakingInfoOnStakingBlock() to always return meaningful Gini.
// - If cache hit                               -> fillMissingGini -> modifies cached in-memory object
// - If db hit                                  -> fillMissingGini -> write to cache
// - If read contract -> write to db (gini: -1) -> fillMissingGini -> write to cache
func GetStakingInfoOnStakingBlock(stakingBlockNumber uint64) *StakingInfo {
	if stakingManager == nil {
		logger.ErrorWithStack("unable to GetStakingInfo", "err", ErrStakingManagerNotSet)
		return nil
	}

	// Return staking info of a previous block if kaia fork is enabled.
	if isKaiaForkEnabled(stakingBlockNumber) {
		logger.Error("unable to use GetStakingInfoOnStakingBlock to get staking info for kaia fork", "staking block number", stakingBlockNumber)
		return nil
	}

	// shortcut if given block is not on staking update interval
	if !params.IsStakingUpdateInterval(stakingBlockNumber) {
		return nil
	}

	// Get staking info from cache
	if cachedStakingInfo := getStakingInfoFromCache(stakingBlockNumber); cachedStakingInfo != nil {
		return cachedStakingInfo
	}

	// Get staking info from DB
	if storedStakingInfo, err := getStakingInfoFromDB(stakingBlockNumber); storedStakingInfo != nil && err == nil {
		logger.Debug("StakingInfoDB hit.", "staking block number", stakingBlockNumber, "stakingInfo", storedStakingInfo)
		addStakingInfoToCache(storedStakingInfo)
		return storedStakingInfo
	} else {
		logger.Debug("failed to get stakingInfo from DB", "err", err, "staking block number", stakingBlockNumber)
	}

	// Calculate staking info from block header and updates it to cache and db
	calcStakingInfo, err := updateStakingInfo(stakingBlockNumber)
	if calcStakingInfo == nil {
		logger.Error("failed to update stakingInfo", "staking block number", stakingBlockNumber, "err", err)
		return nil
	}

	logger.Debug("Get stakingInfo from header.", "staking block number", stakingBlockNumber, "stakingInfo", calcStakingInfo)
	return calcStakingInfo
}

// PreloadStakingInfo preloads staking info for the given headers.
// It first finds the first block that does not have state, and then
// it regenerates the state from the nearest block that has state to the target block to preload staking info.
// Note that the state is saved every 128 blocks to disk in full node.
func PreloadStakingInfo(headers []*types.Header, stakingModule staking.StakingModule) (uint64, error) {
	// If no headers to preload, do nothing
	if len(headers) == 0 {
		return 0, nil
	}

	if stakingManager == nil {
		return 0, ErrStakingManagerNotSet
	}

	var (
		current  *types.Block
		database state.Database
		target   = headers[len(headers)-1].Number.Uint64()
		bc       = stakingManager.blockchain
	)

	database = state.NewDatabaseWithExistingCache(bc.StateCache().TrieDB().DiskDB(), bc.StateCache().TrieDB().TrieNodeCache())

	// Find the first block that does not have state
	i := 0
	for i < len(headers) {
		if _, err := state.New(headers[i].Root, database, nil, nil); err != nil {
			break
		}
		i++
	}
	// Early return if all blocks have state
	if i == len(headers) {
		return 0, nil
	}

	// Find the nearest block that has state
	origin := headers[i].Number.Uint64() - headers[i].Number.Uint64()%128
	current = bc.GetBlockByNumber(origin)
	if current == nil {
		return 0, fmt.Errorf("block %d not found", origin)
	}
	statedb, err := state.New(current.Header().Root, database, nil, nil)
	if err != nil {
		return 0, err
	}

	var (
		parent     common.Hash
		preloadRef = stakingModule.AllocPreloadRef()
	)

	// Include target since we want staking info at `target`, not for `target`.
	for current.NumberU64() <= target {
		stakingModule.PreloadFromState(preloadRef, current.Header(), statedb)
		if current.NumberU64() == target {
			break
		}
		// Retrieve the next block to regenerate and process it
		next := current.NumberU64() + 1
		if current = bc.GetBlockByNumber(next); current == nil {
			return preloadRef, fmt.Errorf("block #%d not found", next)
		}
		_, _, _, _, _, err := bc.Processor().Process(current, statedb, vm.Config{})
		if err != nil {
			return preloadRef, fmt.Errorf("processing block %d failed: %v", current.NumberU64(), err)
		}
		// Finalize the state so any modifications are written to the trie
		root, err := statedb.Commit(true)
		if err != nil {
			return preloadRef, err
		}
		if err := statedb.Reset(root); err != nil {
			return preloadRef, fmt.Errorf("state reset after block %d failed: %v", current.NumberU64(), err)
		}
		database.TrieDB().ReferenceRoot(root)
		if !common.EmptyHash(parent) {
			database.TrieDB().Dereference(parent)
		}
		if current.Root() != root {
			err = fmt.Errorf("mistmatching state root block expected %x reexecuted %x", current.Root(), root)
			// Logging here because something went wrong when the state roots disagree even if the execution was successful.
			logger.Error("incorrectly regenerated historical state", "block", current.NumberU64(), "err", err)
			return preloadRef, fmt.Errorf("incorrectly regenerated historical state for block %d: %v", current.NumberU64(), err)
		}
		parent = root
	}

	return preloadRef, nil
}

// PreloadStakingInfoWithState fetches the stakingInfo based on the given state trie
// and then stores it in the preloaded map. Because preloaded map does not have eviction policy,
// it should be removed manually after use. Note that preloaded info is for the next block of the given header.
func PreloadStakingInfoWithState(header *types.Header, statedb *state.StateDB) error {
	if stakingManager == nil {
		return ErrStakingManagerNotSet
	}

	if !isKaiaForkEnabled(header.Number.Uint64() + 1) {
		return nil // no need to preload staking info before kaia fork because we have it in the database.
	}

	if header.Root != statedb.IntermediateRoot(false) { // Sanity check
		return errors.New("must supply the statedb for the given header") // this should not happen
	}

	num := header.Number.Uint64()
	info, err := getStakingInfoFromMultiCallAtState(num, statedb, header)
	if err != nil {
		return fmt.Errorf("staking info preload failed. root err: %v", err)
	}
	if info != nil {
		stakingManager.preloadedInfo.Add(num, info)
	}
	logger.Trace("preloaded staking info", "staking block number", num)
	return nil
}

// UnloadStakingInfo removes a stakingInfo from the preloaded map.
// Must be executed after PreloadStakingInfoWithState(Header{num}, state).
func UnloadStakingInfo(num uint64) {
	if stakingManager == nil {
		logger.Error("unable to GetStakingInfo", "err", ErrStakingManagerNotSet)
		return
	}

	stakingManager.preloadedInfo.Remove(num)
}

// updateKaiaStakingInfo updates kaia staking info in cache created from given block number.
// From Kaia fork, not only the staking block number but also the calculation of staking amounts is changed,
// so we need separate update function for kaia staking info.
func updateKaiaStakingInfo(blockNum uint64) (*StakingInfo, error) {
	if stakingManager == nil {
		return nil, ErrStakingManagerNotSet
	}

	stakingInfo, err := getStakingInfoFromMultiCall(blockNum)
	if err != nil {
		return nil, err
	}

	addStakingInfoToCache(stakingInfo)
	logger.Debug("Add a new stakingInfo to stakingInfoCache", "blockNum", blockNum)

	logger.Debug("Added stakingInfo", "stakingInfo", stakingInfo)
	return stakingInfo, nil
}

// updateStakingInfo updates staking info in cache and db created from given block number.
func updateStakingInfo(blockNum uint64) (*StakingInfo, error) {
	if stakingManager == nil {
		return nil, ErrStakingManagerNotSet
	}

	stakingInfo, err := getStakingInfoFromAddressBook(blockNum)
	if err != nil {
		return nil, err
	}

	// Add to DB before setting Gini; DB will contain {Gini: -1}
	if err := AddStakingInfoToDB(stakingInfo); err != nil {
		logger.Debug("failed to write staking info to db", "err", err, "stakingInfo", stakingInfo)
		return stakingInfo, err
	}

	addStakingInfoToCache(stakingInfo)
	logger.Info("Add a new stakingInfo to stakingInfoCache and stakingInfoDB", "blockNum", blockNum)

	logger.Debug("Added stakingInfo", "stakingInfo", stakingInfo)
	return stakingInfo, nil
}

// NOTE: Even if the AddressBook contract code is erroneous and it returns unexpected result, this function should not return error in order not to stop block proposal.
// getStakingInfoFromMultiCall returns stakingInfo fetched from MultiCall contract.
// The MultiCall contract gets types and staking addresses from AddressBook contract, and balances of staking addresses.
func getStakingInfoFromMultiCall(blockNum uint64) (*StakingInfo, error) {
	header := stakingManager.blockchain.GetHeaderByNumber(blockNum)
	if header == nil {
		return nil, fmt.Errorf("failed to get header by number %d", blockNum)
	}

	statedb, err := stakingManager.blockchain.StateAt(header.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to get state at number %d. root err: %s", blockNum, err)
	}

	return getStakingInfoFromMultiCallAtState(blockNum, statedb, header)
}

func getStakingInfoFromMultiCallAtState(blockNum uint64, statedb *state.StateDB, header *types.Header) (*StakingInfo, error) {
	caller := backends.NewBlockchainContractBackend(stakingManager.blockchain, nil, nil)
	code, err := caller.CodeAt(context.Background(), addressBookContractAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve code of AddressBook contract. root err: %s", err)
	}
	if code == nil {
		// This is an expected behavior when the addressBook contract is not installed.
		logger.Info("The addressBook is not installed. Use empty stakingInfo")
		return newEmptyStakingInfo(blockNum), nil
	}

	// Get staking info from multicall contract
	contract, err := system.NewMultiCallContractCaller(statedb, stakingManager.blockchain, header)
	if err != nil {
		return nil, fmt.Errorf("failed to create multicall contract caller. root err: %s", err)
	}

	res, err := contract.MultiCallStakingInfo(&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(blockNum)})
	if err != nil {
		return nil, fmt.Errorf("failed to call MultiCall contract. root err: %s", err)
	}

	types := res.TypeList
	addrs := res.AddressList
	stakingAmounts := res.StakingAmounts

	if len(types) == 0 && len(addrs) == 0 {
		// This is an expected behavior when the addressBook contract is not activated yet.
		// Note that multicall contract calls address book internally.
		logger.Info("The addressBook is not yet activated. Use empty stakingInfo")
		return newEmptyStakingInfo(blockNum), nil
	}

	if len(types) != len(addrs) {
		return nil, fmt.Errorf("length of type list and address list differ. len(type)=%d, len(addrs)=%d", len(types), len(addrs))
	}

	return newStakingInfo(stakingManager.blockchain, stakingManager.governanceHelper, blockNum, types, addrs, stakingAmounts...)
}

// NOTE: Even if the AddressBook contract code is erroneous and it returns unexpected result, this function should not return error in order not to stop block proposal.
// getStakingInfoFromAddressBook returns stakingInfo fetched from AddressBook contract
// 1. If calling AddressBook contract fails, it returns error
// 2. If AddressBook is not activated, emptyStakingInfo is returned without error
// 3. If AddressBook is activated, it returns fetched stakingInfo
func getStakingInfoFromAddressBook(blockNum uint64) (*StakingInfo, error) {
	caller := backends.NewBlockchainContractBackend(stakingManager.blockchain, nil, nil)
	code, err := caller.CodeAt(context.Background(), addressBookContractAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve code of AddressBook contract. root err: %s", err)
	}
	if code == nil {
		// This is an expected behavior when the addressBook contract is not installed.
		logger.Info("The addressBook is not installed. Use empty stakingInfo")
		return newEmptyStakingInfo(blockNum), nil
	}

	contract, err := contract.NewAddressBookCaller(addressBookContractAddress, caller)
	if err != nil {
		return nil, fmt.Errorf("failed to call AddressBook contract. root err: %s", err)
	}

	types, addrs, err := contract.GetAllAddress(&bind.CallOpts{BlockNumber: new(big.Int).SetUint64(blockNum)})
	if err != nil {
		return nil, fmt.Errorf("failed to call AddressBook contract. root err: %s", err)
	}

	if len(types) == 0 && len(addrs) == 0 {
		// This is an expected behavior when the addressBook contract is not activated yet.
		logger.Info("The addressBook is not yet activated. Use empty stakingInfo")
		return newEmptyStakingInfo(blockNum), nil
	}

	if len(types) != len(addrs) {
		return nil, fmt.Errorf("length of type list and address list differ. len(type)=%d, len(addrs)=%d", len(types), len(addrs))
	}

	return newStakingInfo(stakingManager.blockchain, stakingManager.governanceHelper, blockNum, types, addrs)
}

func addStakingInfoToCache(stakingInfo *StakingInfo) {
	// Fill in Gini coeff before adding to cache
	if err := fillMissingGiniCoefficient(stakingInfo, stakingInfo.BlockNum); err != nil {
		logger.Warn("Cannot fill in gini coefficient", "staking block number", stakingInfo.BlockNum, "err", err)
	}
	stakingManager.stakingInfoCache.Add(stakingInfo.BlockNum, stakingInfo)
}

func getStakingInfoFromCache(blockNum uint64) *StakingInfo {
	if cachedStakingInfo, ok := stakingManager.stakingInfoCache.Get(blockNum); ok {
		logger.Debug("StakingInfoCache hit.", "staking block number", blockNum, "stakingInfo", cachedStakingInfo)
		// Fill in Gini coeff if not set. Modifies the cached object.
		if err := fillMissingGiniCoefficient(cachedStakingInfo.(*StakingInfo), blockNum); err != nil {
			logger.Warn("Cannot fill in gini coefficient", "staking block number", blockNum, "err", err)
		}
		return cachedStakingInfo.(*StakingInfo)
	}

	if info, ok := stakingManager.preloadedInfo.Get(blockNum); ok {
		info := info.(*StakingInfo)
		logger.Debug("preloadedInfo hit.", "staking block number", blockNum, "stakingInfo", info)
		// Fill in Gini coeff if not set. Modifies the cached object.
		if err := fillMissingGiniCoefficient(info, blockNum); err != nil {
			logger.Warn("Cannot fill in gini coefficient", "staking block number", blockNum, "err", err)
		}
		return info
	}

	return nil
}

// checkStakingInfoStored makes sure the given staking info is stored in cache and DB
func checkStakingInfoStored(blockNum uint64) error {
	if stakingManager == nil {
		return ErrStakingManagerNotSet
	}

	stakingBlockNumber := params.CalcStakingBlockNumber(blockNum)

	// skip checking if staking info is stored in DB
	if _, err := getStakingInfoFromDB(stakingBlockNumber); err == nil {
		return nil
	}

	// update staking info in DB and cache from address book
	_, err := updateStakingInfo(stakingBlockNumber)
	return err
}

// Fill in StakingInfo.Gini value if not set.
func fillMissingGiniCoefficient(stakingInfo *StakingInfo, number uint64) error {
	if !stakingInfo.UseGini {
		return nil
	}
	if stakingInfo.Gini >= 0 {
		return nil
	}

	// We reach here if UseGini == true && Gini == -1. There are two such cases.
	// - Gini was never been calculated, so it is DefaultGiniCoefficient.
	// - Gini was calculated but there was no eligible node, so Gini = -1.
	// For the second case, in theory we won't have to recalculate Gini,
	// but there is no way to distinguish both. So we just recalculate.
	pset, err := stakingManager.governanceHelper.EffectiveParams(number)
	if err != nil {
		return err
	}
	minStaking := pset.MinimumStakeBig().Uint64()

	c := stakingInfo.GetConsolidatedStakingInfo()
	if c == nil {
		return errors.New("Cannot create ConsolidatedStakingInfo")
	}

	stakingInfo.Gini = c.CalcGiniCoefficientMinStake(minStaking)
	logger.Debug("Calculated missing Gini for stored StakingInfo", "number", number, "gini", stakingInfo.Gini)
	return nil
}

// isKaiaForkEnabled returns true if the kaia fork is enabled at the given block number.
func isKaiaForkEnabled(blockNum uint64) bool {
	if stakingManager == nil {
		return false
	}
	return stakingManager.blockchain.Config() != nil && stakingManager.blockchain.Config().IsKaiaForkEnabled(new(big.Int).SetUint64(blockNum))
}

func PurgeStakingInfoCache() {
	stakingManager.stakingInfoCache.Purge()
}

// TODO-Kaia-Reward the following methods are used for testing purpose, it needs to be moved into test files.
// Unlike NewStakingManager(), SetTestStakingManager*() do not trigger once.Do().
// This way you can avoid irreversible side effects during tests.

// SetTestStakingManagerWithChain sets a full-featured staking manager with blockchain, database and cache.
// Note that this method is used only for testing purpose.
func SetTestStakingManagerWithChain(bc blockChain, gh governanceHelper, db stakingInfoDB) {
	cache, _ := lru.NewARC(128)
	SetTestStakingManager(&StakingManager{
		stakingInfoCache: cache,
		stakingInfoDB:    db,
		governanceHelper: gh,
		blockchain:       bc,
		chainHeadChan:    make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		preloadedInfo:    common.NewRefCountingMap(),
	})
}

// SetTestStakingManagerWithDB sets the staking manager with the given database.
// Note that this method is used only for testing purpose.
func SetTestStakingManagerWithDB(testDB stakingInfoDB) {
	SetTestStakingManager(&StakingManager{
		blockchain:    &blockchain.BlockChain{},
		stakingInfoDB: testDB,
		preloadedInfo: common.NewRefCountingMap(),
	})
}

// SetTestStakingManagerWithStakingInfoCache sets the staking manager with the given test staking information.
// Note that this method is used only for testing purpose.
func SetTestStakingManagerWithStakingInfoCache(testInfo *StakingInfo) {
	cache, _ := lru.NewARC(128)
	cache.Add(testInfo.BlockNum, testInfo)
	SetTestStakingManager(&StakingManager{
		blockchain:       &blockchain.BlockChain{},
		stakingInfoCache: cache,
		preloadedInfo:    common.NewRefCountingMap(),
	})
}

// AddTestStakingInfoToCache adds the given test staking information to the cache.
// Note that it won't overwrite the existing cache.
func AddTestStakingInfoToCache(testInfo *StakingInfo) {
	if stakingManager == nil {
		return
	}
	stakingManager.stakingInfoCache.Add(testInfo.BlockNum, testInfo)
}

// SetTestStakingManager sets the staking manager for testing purpose.
// Note that this method is used only for testing purpose.
func SetTestStakingManager(sm *StakingManager) {
	stakingManager = sm
}

// SetTestAddressBookAddress is only for testing purpose.
func SetTestAddressBookAddress(addr common.Address) {
	addressBookContractAddress = common.HexToAddress(addr.Hex())
}

func TestGetStakingCacheSize() int {
	return GetStakingManager().stakingInfoCache.Len()
}

func TestGetStakingPreloadSize() int {
	return GetStakingManager().preloadedInfo.Len()
}
