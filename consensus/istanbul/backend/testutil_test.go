// Modifications Copyright 2020 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
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

package backend

import (
	"crypto/ecdsa"
	"errors"
	"flag"
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/system"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/blockchain/vm"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/consensus/istanbul"
	"github.com/kaiachain/kaia/v2/consensus/istanbul/core"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/crypto/bls"
	"github.com/kaiachain/kaia/v2/datasync/downloader"
	gov_impl "github.com/kaiachain/kaia/v2/kaiax/gov/impl"
	randao_impl "github.com/kaiachain/kaia/v2/kaiax/randao/impl"
	staking_impl "github.com/kaiachain/kaia/v2/kaiax/staking/impl"
	valset_impl "github.com/kaiachain/kaia/v2/kaiax/valset/impl"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/rlp"
	"github.com/kaiachain/kaia/v2/storage/database"
)

var (
	testBaseConfig   *params.ChainConfig
	testKoreConfig   *params.ChainConfig
	testRandaoConfig *params.ChainConfig
)

func init() {
	testBaseConfig = &params.ChainConfig{
		Istanbul:   params.GetDefaultIstanbulConfig(),
		Governance: params.GetDefaultGovernanceConfig(),
	}

	testKoreConfig = testBaseConfig.Copy()
	testKoreConfig.IstanbulCompatibleBlock = common.Big0
	testKoreConfig.LondonCompatibleBlock = common.Big0
	testKoreConfig.EthTxTypeCompatibleBlock = common.Big0
	testKoreConfig.MagmaCompatibleBlock = common.Big0
	testKoreConfig.KoreCompatibleBlock = common.Big0

	testRandaoConfig = testKoreConfig.Copy()
	testRandaoConfig.ShanghaiCompatibleBlock = common.Big0
	testRandaoConfig.CancunCompatibleBlock = common.Big0
	testRandaoConfig.RandaoCompatibleBlock = common.Big0
}

type testOverrides struct {
	node0Key       *ecdsa.PrivateKey // Override node[0] key
	node0BlsKey    bls.SecretKey     // Override node[0] bls key
	blockPeriod    *uint64           // Override block period. If not set, 1 second is used.
	stakingAmounts []uint64          // Override staking amounts. If not set, 0 for all nodes.
}

type testContext struct {
	config      *params.ChainConfig
	nodeKeys    []*ecdsa.PrivateKey // Generated node keys
	nodeAddrs   []common.Address    // Generated node addrs
	nodeBlsKeys []bls.SecretKey     // Generated node bls keys

	chain  *blockchain.BlockChain
	engine *backend
}

func newTestContext(numNodes int, config *params.ChainConfig, overrides *testOverrides) *testContext {
	if config == nil {
		config = testBaseConfig
	}
	if overrides == nil {
		overrides = &testOverrides{}
	}
	if overrides.node0Key == nil {
		overrides.node0Key, _ = crypto.GenerateKey()
	}
	if overrides.node0BlsKey == nil {
		overrides.node0BlsKey, _ = bls.DeriveFromECDSA(overrides.node0Key)
	}
	if overrides.blockPeriod == nil {
		one := uint64(1)
		overrides.blockPeriod = &one
	}
	if overrides.stakingAmounts == nil {
		overrides.stakingAmounts = make([]uint64, numNodes)
	}

	// Create node keys
	var (
		nodeKeys    = make([]*ecdsa.PrivateKey, numNodes)
		nodeAddrs   = make([]common.Address, numNodes)
		nodeBlsKeys = make([]bls.SecretKey, numNodes)

		dbm = database.NewMemoryDBManager()
	)
	nodeKeys[0] = overrides.node0Key
	nodeAddrs[0] = crypto.PubkeyToAddress(nodeKeys[0].PublicKey)
	nodeBlsKeys[0] = overrides.node0BlsKey
	for i := 1; i < numNodes; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		nodeAddrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
		nodeBlsKeys[i], _ = bls.DeriveFromECDSA(nodeKeys[i])
	}

	// Create genesis block
	if config.Governance.GovernanceMode == "single" {
		config.Governance.GoverningNode = nodeAddrs[0]
	}
	genesis := blockchain.DefaultKairosGenesisBlock()
	genesis.Config = config
	genesis.ExtraData = makeGenesisExtra(nodeAddrs)
	genesis.Timestamp = uint64(time.Now().Unix())
	if config.IsRandaoForkEnabled(big.NewInt(0)) {
		allocRegistryStorage := system.AllocRegistry(&params.RegistryConfig{
			Records: map[string]common.Address{
				"KIP113": system.Kip113LogicAddrMock,
			},
			Owner: common.HexToAddress("0xffff"),
		})
		infos := make(map[common.Address]system.BlsPublicKeyInfo)
		for i, addr := range nodeAddrs {
			infos[addr] = system.BlsPublicKeyInfo{
				PublicKey: nodeBlsKeys[i].PublicKey().Marshal(),
				Pop:       bls.PopProve(nodeBlsKeys[i]).Marshal(),
			}
		}
		allocKip113Storage := system.AllocKip113Proxy(system.AllocKip113Init{
			Infos: infos,
			Owner: common.HexToAddress("0xffff"),
		})
		alloc := blockchain.GenesisAlloc{
			system.RegistryAddr: {
				Code:    system.RegistryMockCode,
				Balance: big.NewInt(0),
				Storage: allocRegistryStorage,
			},
			system.Kip113LogicAddrMock: {
				Code:    system.Kip113MockCode,
				Balance: big.NewInt(0),
				Storage: allocKip113Storage,
			},
		}
		genesis.Alloc = alloc
	}
	genesis.MustCommit(dbm)

	// Create istanbul engine
	istanbulConfig := &istanbul.Config{
		Timeout:        10000,
		BlockPeriod:    *overrides.blockPeriod,
		ProposerPolicy: istanbul.ProposerPolicy(config.Istanbul.ProposerPolicy),
		Epoch:          config.Istanbul.Epoch,
		SubGroupSize:   config.Istanbul.SubGroupSize,
	}

	mGov := gov_impl.NewGovModule()
	mRandao := randao_impl.NewRandaoModule()
	engine := New(&BackendOpts{
		IstanbulConfig: istanbulConfig,
		Rewardbase:     common.HexToAddress("0x2A35FE72F847aa0B509e4055883aE90c87558AaD"),
		PrivateKey:     nodeKeys[0],
		BlsSecretKey:   nodeBlsKeys[0],
		DB:             dbm,
		GovModule:      mGov,
		NodeType:       common.CONSENSUSNODE,
	}).(*backend)

	// Create blockchain
	cacheConfig := &blockchain.CacheConfig{
		ArchiveMode:       false,
		CacheSize:         512,
		BlockInterval:     blockchain.DefaultBlockInterval,
		TriesInMemory:     blockchain.DefaultTriesInMemory,
		SnapshotCacheSize: 0, // Disable state snapshot
	}
	chain, err := blockchain.NewBlockChain(dbm, cacheConfig, config, engine, vm.Config{})
	if err != nil {
		panic(err)
	}

	mStaking := staking_impl.NewStakingModule()
	mValset := valset_impl.NewValsetModule()
	fakeDownloader := downloader.NewFakeDownloader()
	if err = errors.Join(
		mGov.Init(&gov_impl.InitOpts{
			Chain:       chain,
			ChainKv:     dbm.GetMiscDB(),
			ChainConfig: config,
			NodeAddress: engine.Address(),
		}),
		mStaking.Init(&staking_impl.InitOpts{
			ChainKv:     dbm.GetMiscDB(),
			ChainConfig: config,
			Chain:       chain,
		}),
		mValset.Init(&valset_impl.InitOpts{
			ChainKv:       dbm.GetMiscDB(),
			Chain:         chain,
			StakingModule: mStaking,
			GovModule:     mGov,
		}),
		mRandao.Init(&randao_impl.InitOpts{
			ChainConfig: config,
			Chain:       chain,
			Downloader:  fakeDownloader,
		})); err != nil {
		panic(err)
	}
	engine.RegisterKaiaxModules(mGov, mStaking, mValset, mRandao)
	// Start the engine
	if err = engine.Start(chain, chain.CurrentBlock, chain.HasBadBlock); err != nil {
		panic(err)
	}

	return &testContext{
		config:    config,
		nodeKeys:  nodeKeys,
		nodeAddrs: nodeAddrs,

		chain:  chain,
		engine: engine,
	}
}

// Make empty header
func (ctx *testContext) MakeHeader(parent *types.Block) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), new(big.Int).SetUint64(ctx.engine.config.BlockPeriod)),
		BlockScore: defaultBlockScore,
	}
	if parent.Header().BaseFee != nil {
		// Assume BaseFee does not change
		header.BaseFee = parent.Header().BaseFee
	}
	return header
}

// Block with no signature.
func (ctx *testContext) MakeBlock(parent *types.Block) *types.Block {
	chain, engine := ctx.chain, ctx.engine
	header := ctx.MakeHeader(parent)
	if err := engine.Prepare(chain, header); err != nil {
		panic(err)
	}
	state, _ := chain.StateAt(parent.Root())
	block, _ := engine.Finalize(chain, header, state, nil, nil)
	return block
}

// Block with proposer seal (no committed seals).
func (ctx *testContext) MakeBlockWithSeal(parent *types.Block) *types.Block {
	chain, engine := ctx.chain, ctx.engine
	block := ctx.MakeBlock(parent)
	result, err := engine.Seal(chain, block, make(chan struct{}))
	if err != nil {
		panic(err)
	}
	return result
}

// Block with proposer seal and all committed seals.
func (ctx *testContext) MakeBlockWithCommittedSeals(parent *types.Block) *types.Block {
	blockWithoutSeal := ctx.MakeBlock(parent)

	// add proposer seal for the block
	block, err := ctx.engine.updateBlock(blockWithoutSeal)
	if err != nil {
		panic(err)
	}

	// write validators committed seals to the block
	header := block.Header()
	committedSeals := ctx.MakeCommittedSeals(block.Hash())
	err = writeCommittedSeals(header, committedSeals)
	if err != nil {
		panic(err)
	}
	block = block.WithSeal(header)

	return block
}

func (ctx *testContext) MakeCommittedSeals(hash common.Hash) [][]byte {
	committedSeals := make([][]byte, len(ctx.nodeKeys))
	hashData := crypto.Keccak256(core.PrepareCommittedSeal(hash))
	for i, key := range ctx.nodeKeys {
		sig, _ := crypto.Sign(hashData, key)
		committedSeals[i] = make([]byte, types.IstanbulExtraSeal)
		copy(committedSeals[i][:], sig)
	}
	return committedSeals
}

func (ctx *testContext) Cleanup() {
	ctx.chain.Stop()
	ctx.engine.Stop()
}

func makeGenesisExtra(addrs []common.Address) []byte {
	extra := &types.IstanbulExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}
	encoded, err := rlp.EncodeToBytes(&extra)
	if err != nil {
		panic(err)
	}

	vanity := make([]byte, types.IstanbulExtraVanity)
	return append(vanity, encoded...)
}

func TestTestContext(t *testing.T) {
	ctx := newTestContext(1, nil, nil)
	defer ctx.Cleanup()
}

func TestMain(m *testing.M) {
	// Because api/debug/flag.go sets the global logger Info level,
	// and BlockChain test generates a lot of Info logs, override to Warn level here.
	flag.Parse() // needed for testing.Verbose()
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	m.Run()
}
