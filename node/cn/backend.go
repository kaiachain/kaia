// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from eth/backend.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package cn

import (
	"errors"
	"fmt"
	"math/big"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/bloombits"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	istanbulBackend "github.com/kaiachain/kaia/consensus/istanbul/backend"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	reward_impl "github.com/kaiachain/kaia/kaiax/reward/impl"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_impl "github.com/kaiachain/kaia/kaiax/staking/impl"
	supply_impl "github.com/kaiachain/kaia/kaiax/supply/impl"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/node/cn/filters"
	"github.com/kaiachain/kaia/node/cn/gasprice"
	"github.com/kaiachain/kaia/node/cn/tracers"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/reward"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/work"
)

var errCNLightSync = errors.New("can't run cn.CN in light sync mode")

//go:generate mockgen -destination=mocks/lesserver_mock.go -package=mocks github.com/kaiachain/kaia/node/cn LesServer
type LesServer interface {
	Start(srvr p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *blockchain.ChainIndexer)
}

// Miner is an interface of work.Miner used by ServiceChain.
//
//go:generate mockgen -destination=mocks/miner_mock.go -package=mocks github.com/kaiachain/kaia/node/cn Miner
type Miner interface {
	Start()
	Stop()
	Register(agent work.Agent)
	Mining() bool
	HashRate() (tot int64)
	SetExtra(extra []byte) error
	Pending() (*types.Block, *state.StateDB)
	PendingBlock() *types.Block
	kaiax.ExecutionModuleHost // Because miner executes blocks, inject ExecutionModule.
}

// BackendProtocolManager is an interface of cn.ProtocolManager used from cn.CN and cn.ServiceChain.
//
//go:generate mockgen -destination=protocolmanager_mock_test.go -package=cn github.com/kaiachain/kaia/node/cn BackendProtocolManager
type BackendProtocolManager interface {
	Downloader() ProtocolManagerDownloader
	SetWsEndPoint(wsep string)
	GetSubProtocols() []p2p.Protocol
	ProtocolVersion() int
	ReBroadcastTxs(transactions types.Transactions)
	SetAcceptTxs()
	NodeType() common.ConnType
	Start(maxPeers int)
	Stop()
	SetSyncStop(flag bool)
	staking.StakingModuleHost
}

// CN implements the Kaia consensus node service.
type CN struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Handlers
	txPool          work.TxPool
	blockchain      work.BlockChain
	protocolManager BackendProtocolManager
	lesServer       LesServer

	// DB interfaces
	chainDB database.DBManager // Block chain database

	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager accounts.AccountManager

	bloomRequests     chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer      *blockchain.ChainIndexer       // Bloom indexer operating during block imports
	closeBloomHandler chan struct{}

	APIBackend *CNAPIBackend

	miner    Miner
	gasPrice *big.Int

	rewardbase  common.Address
	nodeAddress common.Address

	networkId     uint64
	netRPCService *api.PublicNetAPI

	lock sync.RWMutex // Protects the variadic fields (e.g. gas price)

	components []interface{}

	governance governance.Engine

	// kaiax modules
	baseModules    []kaiax.BaseModule
	jsonRpcModules []kaiax.JsonRpcModule
	stakingModule  staking.StakingModule // TODO-kaiax: temporary for governance/api.go. Remove it after having kaiax/reward.
}

func (s *CN) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

// senderTxHashIndexer subscribes chainEvent and stores senderTxHash to txHash mapping information.
func senderTxHashIndexer(db database.DBManager, chainEvent <-chan blockchain.ChainEvent, subscription event.Subscription) {
	defer subscription.Unsubscribe()

	for {
		select {
		case event := <-chainEvent:
			var err error
			batch := db.NewSenderTxHashToTxHashBatch()
			for _, tx := range event.Block.Transactions() {
				senderTxHash, ok := tx.SenderTxHash()

				// senderTxHash and txHash are the same if tx is not a fee-delegated tx.
				// Do not store mapping between senderTxHash and txHash in this case.
				if !ok {
					continue
				}

				txHash := tx.Hash()
				db.PutSenderTxHashToTxHashToBatch(batch, senderTxHash, txHash)
			}

			if err == nil {
				batch.Write()
				batch.Release()
			}

		case <-subscription.Err():
			return
		}
	}
}

func checkSyncMode(config *Config) error {
	// TODO-Kaia: allow snap sync after resolving the staking info sync issue
	if config.SyncMode == downloader.SnapSync {
		return errors.New("snap sync is temporarily disabled")
	}
	if !config.SyncMode.IsValid() {
		return fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	if config.SyncMode == downloader.LightSync {
		return errCNLightSync
	}
	return nil
}

func setEngineType(chainConfig *params.ChainConfig) {
	if chainConfig.Clique != nil {
		types.EngineType = types.Engine_Clique
	}
	if chainConfig.Istanbul != nil {
		types.EngineType = types.Engine_IBFT
	}
}

// New creates a new CN object (including the
// initialisation of the common CN object)
func New(ctx *node.ServiceContext, config *Config) (*CN, error) {
	if err := checkSyncMode(config); err != nil {
		return nil, err
	}

	chainDB := CreateDB(ctx, config, "chaindata")

	chainConfig, genesisHash, genesisErr := blockchain.SetupGenesisBlock(chainDB, config.Genesis, config.NetworkId, config.IsPrivate, false)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	setEngineType(chainConfig)

	// load governance state
	chainConfig.SetDefaults()
	// latest values will be applied to chainConfig after NewMixedEngine call
	governance := governance.NewMixedEngine(chainConfig, chainDB)
	logger.Info("Initialised chain configuration", "config", chainConfig)

	config.GasPrice = new(big.Int).SetUint64(chainConfig.UnitPrice)

	cn := &CN{
		config:            config,
		chainDB:           chainDB,
		chainConfig:       chainConfig,
		eventMux:          ctx.EventMux,
		accountManager:    ctx.AccountManager,
		engine:            CreateConsensusEngine(ctx, config, chainConfig, chainDB, governance, ctx.NodeType()),
		networkId:         config.NetworkId,
		gasPrice:          config.GasPrice,
		rewardbase:        config.Rewardbase,
		bloomRequests:     make(chan chan *bloombits.Retrieval),
		bloomIndexer:      NewBloomIndexer(chainDB, params.BloomBitsBlocks),
		closeBloomHandler: make(chan struct{}),
		governance:        governance,
	}

	// istanbul BFT. Derive and set node's address using nodekey
	if cn.chainConfig.Istanbul != nil {
		cn.nodeAddress = crypto.PubkeyToAddress(ctx.NodeKey().PublicKey)
		governance.SetNodeAddress(cn.nodeAddress)
	}

	logger.Info("Initialising Klaytn protocol", "versions", cn.engine.Protocol().Versions, "network", config.NetworkId)

	if !config.SkipBcVersionCheck {
		if err := blockchain.CheckBlockChainVersion(chainDB); err != nil {
			return nil, err
		}
	}
	var (
		vmConfig    = config.getVMConfig()
		cacheConfig = &blockchain.CacheConfig{
			ArchiveMode:          config.NoPruning,
			CacheSize:            config.TrieCacheSize,
			BlockInterval:        config.TrieBlockInterval,
			TriesInMemory:        config.TriesInMemory,
			LivePruningRetention: config.LivePruningRetention,
			TrieNodeCacheConfig:  &config.TrieNodeCacheConfig,
			SenderTxHashIndexing: config.SenderTxHashIndexing,
			SnapshotCacheSize:    config.SnapshotCacheSize,
			SnapshotAsyncGen:     config.SnapshotAsyncGen,
		}
	)

	bc, err := blockchain.NewBlockChain(chainDB, cacheConfig, cn.chainConfig, cn.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	bc.SetCanonicalBlock(config.StartBlockNumber)

	// Write the live pruning flag to database if the node is started for the first time
	if config.LivePruning && !chainDB.ReadPruningEnabled() {
		if bc.CurrentBlock().NumberU64() > 0 {
			return nil, errors.New("cannot enable live pruning after chain has advanced")
		}
		logger.Info("Writing live pruning flag to database")
		chainDB.WritePruningEnabled()
	}

	// Live pruning is enabled according to the flag in database
	// regardless of the command line flag --state.live-pruning
	// But live pruning is disabled when --state.live-pruning-retention=0
	if chainDB.ReadPruningEnabled() && config.LivePruningRetention != 0 {
		logger.Info("Live pruning is enabled", "retention", config.LivePruningRetention)
	} else if !chainDB.ReadPruningEnabled() {
		logger.Info("Live pruning is disabled because flag not stored in database")
	} else if config.LivePruningRetention == 0 {
		logger.Info("Live pruning is disabled because retention is set to zero")
	}

	cn.blockchain = bc
	governance.SetBlockchain(cn.blockchain)
	if err := governance.UpdateParams(cn.blockchain.CurrentBlock().NumberU64()); err != nil {
		return nil, err
	}
	blockchain.InitDeriveShaWithGov(cn.chainConfig, governance)

	// Synchronize proposerpolicy & useGiniCoeff
	pset, err := governance.EffectiveParams(bc.CurrentBlock().NumberU64() + 1)
	if err != nil {
		return nil, err
	}
	if cn.blockchain.Config().Istanbul != nil {
		cn.blockchain.Config().Istanbul.ProposerPolicy = pset.Policy()
	}
	if cn.blockchain.Config().Governance.Reward != nil {
		cn.blockchain.Config().Governance.Reward.UseGiniCoeff = pset.UseGiniCoeff()
	}

	if config.SenderTxHashIndexing {
		ch := make(chan blockchain.ChainEvent, 255)
		chainEventSubscription := cn.blockchain.SubscribeChainEvent(ch)
		go senderTxHashIndexer(chainDB, ch, chainEventSubscription)
	}

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		logger.Error("Rewinding chain to upgrade configuration", "err", compat)
		cn.blockchain.SetHead(compat.RewindTo)
		chainDB.WriteChainConfig(genesisHash, cn.chainConfig)
	}
	cn.bloomIndexer.Start(cn.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}
	// TODO-Kaia-ServiceChain: add account creation prevention in the txPool if TxTypeAccountCreation is supported.
	config.TxPool.NoAccountCreation = config.NoAccountCreation
	cn.txPool = blockchain.NewTxPool(config.TxPool, cn.chainConfig, bc)
	governance.SetTxPool(cn.txPool)

	// Permit the downloader to use the trie cache allowance during fast sync
	cacheLimit := cacheConfig.TrieNodeCacheConfig.LocalCacheSizeMiB
	if cn.protocolManager, err = NewProtocolManager(cn.chainConfig, config.SyncMode, config.NetworkId, cn.eventMux, cn.txPool, cn.engine, cn.blockchain, chainDB, cacheLimit, ctx.NodeType(), config); err != nil {
		return nil, err
	}

	if err := cn.setAcceptTxs(); err != nil {
		logger.Error("Failed to decode IstanbulExtra", "err", err)
	}

	cn.protocolManager.SetWsEndPoint(config.WsEndpoint)

	if ctx.NodeType() == common.CONSENSUSNODE {
		logger.Info("Loaded node keys",
			"nodeAddress", crypto.PubkeyToAddress(ctx.NodeKey().PublicKey),
			"nodePublicKey", hexutil.Encode(crypto.FromECDSAPub(&ctx.NodeKey().PublicKey)),
			"blsPublicKey", hexutil.Encode(ctx.BlsNodeKey().PublicKey().Marshal()))

		if _, err := cn.Rewardbase(); err != nil {
			logger.Error("Cannot determine the rewardbase address", "err", err)
		}
	}

	// Governance states which are not yet applied to the db remains at in-memory storage
	// It disappears during the node restart, so restoration is needed before the sync starts
	// By calling CreateSnapshot, it restores the gov state snapshots and apply the votes in it
	// Particularly, the gov.changeSet is also restored here.
	// Temporarily set chain since snapshot needs state since kaia hardfork
	logger.Info("Start creating istanbul snapshot")
	cn.createSnapshot()
	logger.Info("Finished creating istanbul snapshot")

	// set worker
	if config.WorkerDisable {
		cn.miner = work.NewFakeWorker()
		// Istanbul backend can be accessed by APIs to call its methods even though the core of the
		// consensus engine doesn't run.
		istBackend, ok := cn.engine.(consensus.Istanbul)
		if ok {
			istBackend.SetChain(cn.blockchain)
		}
	} else {
		// TODO-Kaia improve to handle drop transaction on network traffic in PN and EN
		cn.miner = work.New(cn, cn.chainConfig, cn.EventMux(), cn.engine, ctx.NodeType(), crypto.PubkeyToAddress(ctx.NodeKey().PublicKey), cn.config.TxResendUseLegacy)
	}

	// istanbul BFT
	cn.miner.SetExtra(makeExtraData(config.ExtraData))

	cn.APIBackend = &CNAPIBackend{cn, nil}

	gpoParams := config.GPO

	// NOTE-Kaia Now we use latest unitPrice
	//         So let's override gpoParams.Default with config.GasPrice
	gpoParams.Default = config.GasPrice

	cn.APIBackend.gpo = gasprice.NewOracle(cn.APIBackend, gpoParams, cn.txPool, cn.governance)
	//@TODO Kaia add core component
	cn.addComponent(cn.blockchain)
	cn.addComponent(cn.txPool)
	cn.addComponent(cn.APIs())
	cn.addComponent(cn.ChainDB())
	cn.addComponent(cn.engine)

	if err := cn.SetupKaiaxModules(); err != nil {
		logger.Error("Failed to setup kaiax modules", "err", err)
	}

	if config.AutoRestartFlag {
		daemonPath := config.DaemonPathFlag
		restartInterval := config.RestartTimeOutFlag
		if restartInterval <= time.Second {
			logger.Crit("Invalid auto-restart timeout", "timeout", restartInterval)
		}

		// Restarts the node with the same configuration if blockNumber is not changed for a specific time.
		restartTimer := time.AfterFunc(restartInterval, func() {
			logger.Warn("Restart node", "command", daemonPath+" restart")
			cmd := exec.Command(daemonPath, "restart")
			cmd.Run()
		})
		logger.Info("Initialize auto-restart feature", "timeout", restartInterval, "daemonPath", daemonPath)

		go func() {
			blockChecker := time.NewTicker(time.Second)
			prevBlockNum := cn.blockchain.CurrentBlock().NumberU64()

			for range blockChecker.C {
				currentBlockNum := cn.blockchain.CurrentBlock().NumberU64()

				if prevBlockNum != currentBlockNum {
					prevBlockNum = currentBlockNum
					restartTimer.Reset(restartInterval)
				}
			}
		}()
	}

	// Only for KES nodes
	if config.TrieNodeCacheConfig.RedisSubscribeBlockEnable {
		go cn.blockchain.BlockSubscriptionLoop(cn.txPool.(*blockchain.TxPool))
	}

	if config.DBType == database.RocksDB && config.RocksDBConfig.Secondary {
		go cn.blockchain.CurrentBlockUpdateLoop(cn.txPool.(*blockchain.TxPool))
	}

	return cn, nil
}

func (s *CN) createSnapshot() {
	var (
		currBlock = s.blockchain.CurrentBlock()
		headers   []*types.Header
		err       error
	)
	// Temporarily supply blockchain for `Finalize`, staking module for `snapshot` and reward module for `processor`.
	s.blockchain.Engine().(consensus.Istanbul).SetChain(s.blockchain)
	mStaking := staking_impl.NewStakingModule()
	mReward := reward_impl.NewRewardModule()
	mStaking.Init(&staking_impl.InitOpts{
		ChainKv:     s.chainDB.GetMiscDB(),
		ChainConfig: s.chainConfig,
		Chain:       s.blockchain,
	})
	mReward.Init(&reward_impl.InitOpts{
		ChainConfig:   s.chainConfig,
		Chain:         s.blockchain,
		GovModule:     reward_impl.FromLegacy(s.governance),
		StakingModule: mStaking,
	})
	s.blockchain.Engine().(consensus.Istanbul).RegisterStakingModule(mStaking)
	s.blockchain.Engine().(consensus.Istanbul).RegisterConsensusModule(mReward)
	if headers, err = s.Engine().GetKaiaHeadersForSnapshotApply(s.blockchain, currBlock.NumberU64(), currBlock.Hash(), nil); err != nil {
		logger.Error("Failed to get headers to apply", "err", err)
	} else {
		preloadRef, err := reward.PreloadStakingInfo(s.blockchain, headers, mStaking)
		if err != nil {
			logger.Error("Preload staking info failed", "err", err)
		}
		defer func() {
			mStaking.FreePreloadRef(preloadRef)
		}()
	}
	if err := s.Engine().CreateSnapshot(s.blockchain, currBlock.NumberU64(), currBlock.Hash(), headers); err != nil {
		logger.Error("CreateSnapshot failed", "err", err)
	}
	s.blockchain.Engine().(consensus.Istanbul).SetChain(nil)
	s.blockchain.Engine().(consensus.Istanbul).RegisterStakingModule(nil)
	s.blockchain.Engine().(consensus.Istanbul).UnregisterConsensusModule(mReward)
}

// setAcceptTxs sets AcceptTxs flag in 1CN case to receive tx propagation.
func (s *CN) setAcceptTxs() error {
	if s.chainConfig.Istanbul != nil {
		istanbulExtra, err := types.ExtractIstanbulExtra(s.blockchain.Genesis().Header())
		if err != nil {
			return err
		} else {
			if len(istanbulExtra.Validators) == 1 {
				s.protocolManager.SetAcceptTxs()
			}
		}
	}
	return nil
}

// add component which may be used in another service component
func (s *CN) addComponent(component interface{}) {
	s.components = append(s.components, component)
}

func (s *CN) Components() []interface{} {
	return s.components
}

func (s *CN) SetComponents(component []interface{}) {
	// do nothing
}

func (s *CN) SetupKaiaxModules() error {
	// Declare modules

	var (
		mStaking = staking_impl.NewStakingModule()
		mReward  = reward_impl.NewRewardModule()
		mSupply  = supply_impl.NewSupplyModule()
		mGov     = gov_impl.NewGovModule()
	)

	// Initialize modules
	err := errors.Join(
		mStaking.Init(&staking_impl.InitOpts{
			ChainKv:     s.chainDB.GetMiscDB(),
			ChainConfig: s.chainConfig,
			Chain:       s.blockchain,
		}),
		mReward.Init(&reward_impl.InitOpts{
			ChainConfig:   s.chainConfig,
			Chain:         s.blockchain,
			GovModule:     reward_impl.FromLegacy(s.governance),
			StakingModule: mStaking,
		}),
		mSupply.Init(&supply_impl.InitOpts{
			ChainKv:      s.chainDB.GetMiscDB(),
			ChainConfig:  s.chainConfig,
			Chain:        s.blockchain,
			RewardModule: mReward,
		}),
		mGov.Init(&gov_impl.InitOpts{
			ChainConfig: s.chainConfig,
			ChainKv:     s.chainDB.GetMiscDB(),
			Chain:       s.blockchain,
			NodeAddress: s.nodeAddress,
		}),
	)
	if err != nil {
		return err
	}

	// Register modules to respective components
	// TODO-kaiax: Organize below lines.
	s.RegisterBaseModules(mStaking, mReward, mSupply, mGov)
	s.RegisterJsonRpcModules(mStaking, mReward, mSupply, mGov)
	s.miner.RegisterExecutionModule(mStaking, mSupply, mGov)
	s.blockchain.RegisterExecutionModule(mSupply, mGov)
	s.blockchain.RegisterRewindableModule(mStaking, mSupply, mGov)
	if engine, ok := s.engine.(consensus.Istanbul); ok {
		engine.RegisterStakingModule(mStaking)
		engine.RegisterConsensusModule(mReward, mGov)
	}
	s.protocolManager.RegisterStakingModule(mStaking)

	s.stakingModule = mStaking
	return nil
}

func (s *CN) RegisterBaseModules(modules ...kaiax.BaseModule) {
	// Add to s.modules so s.Start() and s.Stop() can call them.
	s.baseModules = append(s.baseModules, modules...)
}

func (s *CN) RegisterJsonRpcModules(modules ...kaiax.JsonRpcModule) {
	// Add to s.modules so s.Start() and s.Stop() can call them.
	s.jsonRpcModules = append(s.jsonRpcModules, modules...)
}

// istanbul BFT
func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"klay",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.GetMaximumExtraDataSize() {
		logger.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.GetMaximumExtraDataSize())
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) database.DBManager {
	dbc := &database.DBConfig{
		Dir: name, DBType: config.DBType, ParallelDBWrite: config.ParallelDBWrite, SingleDB: config.SingleDB, NumStateTrieShards: config.NumStateTrieShards,
		LevelDBCacheSize: config.LevelDBCacheSize, LevelDBCompression: config.LevelDBCompression,
		PebbleDBCacheSize: config.PebbleDBCacheSize, OpenFilesLimit: database.GetOpenFilesLimit(),
		LevelDBBufferPool: config.LevelDBBufferPool, EnableDBPerfMetrics: config.EnableDBPerfMetrics, RocksDBConfig: &config.RocksDBConfig, DynamoDBConfig: &config.DynamoDBConfig,
	}
	return ctx.OpenDatabase(dbc)
}

// CreateConsensusEngine creates the required type of consensus engine instance for a Kaia service
func CreateConsensusEngine(ctx *node.ServiceContext, config *Config, chainConfig *params.ChainConfig, db database.DBManager, gov governance.Engine, nodetype common.ConnType) consensus.Engine {
	// Only istanbul  BFT is allowed in the main net. PoA is supported by service chain
	if chainConfig.Governance == nil {
		chainConfig.Governance = params.GetDefaultGovernanceConfig()
	}
	return istanbulBackend.New(&istanbulBackend.BackendOpts{
		IstanbulConfig: &config.Istanbul,
		Rewardbase:     config.Rewardbase,
		PrivateKey:     ctx.NodeKey(),
		BlsSecretKey:   ctx.BlsNodeKey(),
		DB:             db,
		Governance:     gov,
		NodeType:       nodetype,
	})
}

// APIs returns the collection of RPC services the ethereum package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
func (s *CN) APIs() []rpc.API {
	var (
		nonceLock                = new(api.AddrLocker)
		publicBlockChainAPI      = api.NewPublicBlockChainAPI(s.APIBackend)
		publicKaiaAPI            = api.NewPublicKaiaAPI(s.APIBackend)
		publicTransactionPoolAPI = api.NewPublicTransactionPoolAPI(s.APIBackend, nonceLock)
		publicAccountAPI         = api.NewPublicAccountAPI(s.APIBackend.AccountManager())
	)

	apis := []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicKaiaAPI,
			Public:    true,
		}, {
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicBlockChainAPI,
			Public:    true,
		}, {
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicTransactionPoolAPI,
			Public:    true,
		}, {
			Namespace: "txpool",
			Version:   "1.0",
			Service:   api.NewPublicTxPoolAPI(s.APIBackend),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   api.NewPublicDebugAPI(s.APIBackend),
			Public:    false,
		}, {
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicAccountAPI,
			Public:    true,
		}, {
			Namespace: "personal",
			Version:   "1.0",
			Service:   api.NewPrivateAccountAPI(s.APIBackend, nonceLock),
			Public:    false,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   api.NewPrivateDebugAPI(s.APIBackend),
			Public:    false,
			IPCOnly:   s.config.DisableUnsafeDebug,
		},
	}

	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	publicFilterAPI := filters.NewPublicFilterAPI(s.APIBackend, false)
	publicDownloaderAPI := downloader.NewPublicDownloaderAPI(s.protocolManager.Downloader(), s.eventMux)
	privateDownloaderAPI := downloader.NewPrivateDownloaderAPI(s.protocolManager.Downloader())

	ethAPI := api.NewEthereumAPI(
		publicFilterAPI,
		publicKaiaAPI,
		publicBlockChainAPI,
		publicTransactionPoolAPI,
		publicAccountAPI,
		s.nodeAddress,
	)

	// Append all the local APIs and return
	apis = append(apis, []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewPublicKaiaAPI(s),
			Public:    true,
		}, {
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicDownloaderAPI,
			Public:    true,
		}, {
			Namespace: "kaia",
			Version:   "1.0",
			Service:   publicFilterAPI,
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   publicDownloaderAPI,
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   privateDownloaderAPI,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    false,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   tracers.NewAPI(s.APIBackend),
			Public:    false,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   tracers.NewUnsafeAPI(s.APIBackend),
			Public:    false,
			IPCOnly:   s.config.DisableUnsafeDebug,
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   ethAPI,
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
			Public:    false,
			IPCOnly:   s.config.DisableUnsafeDebug,
		},
	}...)

	// Append APIs exposed by JsonRpcModules
	for _, module := range s.jsonRpcModules {
		apis = append(apis, module.APIs()...)
	}

	return apis
}

func (s *CN) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

func (s *CN) Rewardbase() (eb common.Address, err error) {
	s.lock.RLock()
	rewardbase := s.rewardbase
	s.lock.RUnlock()

	if rewardbase != (common.Address{}) {
		return rewardbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			rewardbase := accounts[0].Address

			s.lock.Lock()
			s.rewardbase = rewardbase
			s.lock.Unlock()

			logger.Info("Rewardbase automatically configured", "address", rewardbase)
			return rewardbase, nil
		}
	}

	return common.Address{}, fmt.Errorf("rewardbase must be explicitly specified")
}

func (s *CN) StartMining(local bool) error {
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so none will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		s.protocolManager.SetAcceptTxs()
	}
	go s.miner.Start()
	return nil
}

func (s *CN) StopMining()    { s.miner.Stop() }
func (s *CN) IsMining() bool { return s.miner.Mining() }
func (s *CN) Miner() Miner   { return s.miner }

func (s *CN) AccountManager() accounts.AccountManager { return s.accountManager }
func (s *CN) BlockChain() work.BlockChain             { return s.blockchain }
func (s *CN) TxPool() work.TxPool                     { return s.txPool }
func (s *CN) EventMux() *event.TypeMux                { return s.eventMux }
func (s *CN) Engine() consensus.Engine                { return s.engine }
func (s *CN) ChainDB() database.DBManager             { return s.chainDB }
func (s *CN) IsListening() bool                       { return true } // Always listening
func (s *CN) ProtocolVersion() int                    { return s.protocolManager.ProtocolVersion() }
func (s *CN) NetVersion() uint64                      { return s.networkId }
func (s *CN) Progress() kaia.SyncProgress             { return s.protocolManager.Downloader().Progress() }
func (s *CN) Governance() governance.Engine           { return s.governance }

func (s *CN) ReBroadcastTxs(transactions types.Transactions) {
	s.protocolManager.ReBroadcastTxs(transactions)
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *CN) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.GetSubProtocols()
	}
	return append(s.protocolManager.GetSubProtocols(), s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Kaia protocol implementation.
func (s *CN) Start(srvr p2p.Server) error {
	// Start kaiax modules in the order they were registered
	for _, module := range s.baseModules {
		if err := module.Start(); err != nil {
			return err
		}
	}

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = api.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers()
	// Start the networking layer and the light server if requested
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}

	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// Kaia protocol.
func (s *CN) Stop() error {
	// Stop all the peer-related stuff first.
	s.protocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}

	// Then stop everything else.
	for _, module := range s.baseModules {
		module.Stop()
	}
	s.bloomIndexer.Close()
	close(s.closeBloomHandler)
	s.txPool.Stop()
	s.miner.Stop()
	s.blockchain.Stop()
	s.chainDB.Close()
	s.eventMux.Stop()

	return nil
}
