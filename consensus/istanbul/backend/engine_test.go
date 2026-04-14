// Modifications Copyright 2024 The Kaia Authors
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
//
// This file is derived from quorum/consensus/istanbul/backend/engine_test.go (2020/04/16).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	randao_impl "github.com/kaiachain/kaia/kaiax/randao/impl"
	reward_impl "github.com/kaiachain/kaia/kaiax/reward/impl"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_impl "github.com/kaiachain/kaia/kaiax/staking/impl"
	"github.com/kaiachain/kaia/kaiax/staking/mock"
	system_impl "github.com/kaiachain/kaia/kaiax/system/impl"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_impl "github.com/kaiachain/kaia/kaiax/valset/impl"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

// These variables are the global variables of the test blockchain.
var (
	nodeKeys       []*ecdsa.PrivateKey
	addrs          []common.Address
	testGovModules sync.Map // map[*backend]gov.GovModule
)

func govModuleOf(b *backend) gov.GovModule {
	m, ok := testGovModules.Load(b)
	if !ok {
		panic("missing gov module for test backend")
	}
	return m.(gov.GovModule)
}

// These are the types in order to add a custom configuration of the test chain.
// You may need to create a configuration type if necessary.
type (
	istanbulCompatibleBlock  *big.Int
	LondonCompatibleBlock    *big.Int
	EthTxTypeCompatibleBlock *big.Int
	magmaCompatibleBlock     *big.Int
	koreCompatibleBlock      *big.Int
	shanghaiCompatibleBlock  *big.Int
	cancunCompatibleBlock    *big.Int
	randaoCompatibleBlock    *big.Int
	kaiaCompatibleBlock      *big.Int
	pragueCompatibleBlock    *big.Int
	osakaCompatibleBlock     *big.Int
)

type (
	minimumStake           *big.Int
	mintingAmount          *big.Int
	lowerBoundBaseFee      uint64
	upperBoundBaseFee      uint64
	stakingUpdateInterval  uint64
	proposerUpdateInterval uint64
	proposerPolicy         uint64
	governanceMode         string
	epoch                  uint64
	subGroupSize           uint64
	blockPeriod            uint64
)

// Include a node from the global nodeKeys and addrs
func includeNode(addr common.Address, key *ecdsa.PrivateKey) {
	for _, a := range addrs {
		if a.String() == addr.String() {
			// already exists
			return
		}
	}
	nodeKeys = append(nodeKeys, key)
	addrs = append(addrs, addr)
}

// Exclude a node from the global nodeKeys and addrs
func excludeNodeByAddr(target common.Address) {
	for i, a := range addrs {
		if a.String() == target.String() {
			nodeKeys = append(nodeKeys[:i], nodeKeys[i+1:]...)
			addrs = append(addrs[:i], addrs[i+1:]...)
			break
		}
	}
}

func enableVotes(paramNames []gov.ParamName) {
	for _, paramName := range paramNames {
		gov.Params[paramName].VoteForbidden = false
	}
}

func disableVotes(paramNames []gov.ParamName) {
	for _, paramName := range paramNames {
		gov.Params[paramName].VoteForbidden = true
	}
}

func setNodeKeys(n int, governingNode *ecdsa.PrivateKey) ([]*ecdsa.PrivateKey, []common.Address) {
	nodeKeys = make([]*ecdsa.PrivateKey, n)
	addrs = make([]common.Address, n)
	for i := range n {
		if i == 0 && governingNode != nil {
			nodeKeys[i] = governingNode
		} else {
			nodeKeys[i], _ = crypto.GenerateKey()
		}
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}
	return nodeKeys, addrs
}

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain(t *testing.T, n int, items ...interface{}) (*blockchain.BlockChain, *backend) {
	// Keep PrepareHeader timestamp logic aligned with block generation interval in tests.
	oldInterval := params.BlockGenerationInterval
	t.Cleanup(func() {
		params.BlockGenerationInterval = oldInterval
	})

	// generate a genesis block
	genesis := blockchain.DefaultTestGenesisBlock()
	genesis.Config = params.TestChainConfig.Copy()

	var (
		period   = uint64(params.DefaultBlockGenerationInterval)
		mStaking staking.StakingModule
		err      error
	)
	// force enable Istanbul engine and governance
	for _, item := range items {
		switch v := item.(type) {
		case istanbulCompatibleBlock:
			genesis.Config.IstanbulCompatibleBlock = v
		case LondonCompatibleBlock:
			genesis.Config.LondonCompatibleBlock = v
		case EthTxTypeCompatibleBlock:
			genesis.Config.EthTxTypeCompatibleBlock = v
		case magmaCompatibleBlock:
			genesis.Config.MagmaCompatibleBlock = v
		case koreCompatibleBlock:
			genesis.Config.KoreCompatibleBlock = v
		case shanghaiCompatibleBlock:
			genesis.Config.ShanghaiCompatibleBlock = v
		case cancunCompatibleBlock:
			genesis.Config.CancunCompatibleBlock = v
		case randaoCompatibleBlock:
			genesis.Config.RandaoCompatibleBlock = v
		case kaiaCompatibleBlock:
			genesis.Config.KaiaCompatibleBlock = v
		case pragueCompatibleBlock:
			genesis.Config.PragueCompatibleBlock = v
		case osakaCompatibleBlock:
			genesis.Config.OsakaCompatibleBlock = v
		case proposerPolicy:
			genesis.Config.Istanbul.ProposerPolicy = uint64(v)
		case epoch:
			genesis.Config.Istanbul.Epoch = uint64(v)
		case subGroupSize:
			genesis.Config.Istanbul.SubGroupSize = uint64(v)
		case minimumStake:
			genesis.Config.Governance.Reward.MinimumStake = v
		case stakingUpdateInterval:
			genesis.Config.Governance.Reward.StakingUpdateInterval = uint64(v)
		case proposerUpdateInterval:
			genesis.Config.Governance.Reward.ProposerUpdateInterval = uint64(v)
		case mintingAmount:
			genesis.Config.Governance.Reward.MintingAmount = v
		case governanceMode:
			genesis.Config.Governance.GovernanceMode = string(v)
		case lowerBoundBaseFee:
			genesis.Config.Governance.KIP71.LowerBoundBaseFee = uint64(v)
		case upperBoundBaseFee:
			genesis.Config.Governance.KIP71.UpperBoundBaseFee = uint64(v)
		case blockPeriod:
			period = uint64(v)
		case *params.ChainConfig:
			genesis.Config = v
		case *mock.MockStakingModule:
			mStaking = v
		}
	}
	genesis.Config.SetDefaults()
	params.BlockGenerationInterval = int64(period)

	now := uint64(time.Now().Unix())
	if now > period {
		genesis.Timestamp = now - period
	} else {
		genesis.Timestamp = 0
	}

	if len(nodeKeys) != n {
		setNodeKeys(n, nil)
	}

	// if governance mode is single, this address is the governing node address
	b := newTestBackendWithConfig(genesis.Config, nodeKeys[0])

	genesisHeader := &types.Header{
		Number: new(big.Int).SetUint64(genesis.Number),
		Extra:  append([]byte(nil), genesis.ExtraData...),
	}
	if err := b.sealer.WriteValidators(genesisHeader, addrs); err != nil {
		panic(err)
	}
	genesis.ExtraData = append([]byte(nil), genesisHeader.Extra...)

	// Set up Registry and KIP113 contracts for Randao fork if RandaoCompatibleBlock is set
	if genesis.Config.RandaoCompatibleBlock != nil {
		// Generate BLS keys for all nodes
		nodeBlsKeys := make([]bls.SecretKey, n)
		for i := range n {
			nodeBlsKeys[i], _ = bls.DeriveFromECDSA(nodeKeys[i])
		}

		allocRegistryStorage := system.AllocRegistry(&params.RegistryConfig{
			Records: map[string]common.Address{
				"KIP113": system.Kip113LogicAddrMock,
			},
			Owner: common.HexToAddress("0xffff"),
		})
		infos := make(map[common.Address]system.BlsPublicKeyInfo)
		for i, addr := range addrs {
			infos[addr] = system.BlsPublicKeyInfo{
				PublicKey: nodeBlsKeys[i].PublicKey().Marshal(),
				Pop:       bls.PopProve(nodeBlsKeys[i]).Marshal(),
			}
		}
		allocKip113Storage := system.AllocKip113Proxy(system.AllocKip113Init{
			Infos: infos,
			Owner: common.HexToAddress("0xffff"),
		})
		if genesis.Alloc == nil {
			genesis.Alloc = make(blockchain.GenesisAlloc)
		}
		genesis.Alloc[system.RegistryAddr] = blockchain.GenesisAccount{
			Code:    system.RegistryMockCode,
			Balance: big.NewInt(0),
			Storage: allocRegistryStorage,
		}
		genesis.Alloc[system.Kip113LogicAddrMock] = blockchain.GenesisAccount{
			Code:    system.Kip113MockCode,
			Balance: big.NewInt(0),
			Storage: allocKip113Storage,
		}
	}

	genesisGov := make(gov.PartialParamSet)
	for name, param := range gov.Params {
		val, err := param.ChainConfigValue(genesis.Config)
		if err != nil {
			panic(err)
		}
		err = genesisGov.Add(string(name), val)
		if err != nil {
			panic(err)
		}
	}

	genesis.Governance, err = headergov.NewGovData(genesisGov).ToGovBytes()
	if err != nil {
		panic(err)
	}

	dbm := database.NewMemoryDBManager()
	genesis.MustCommit(dbm)

	chainSealer := istanbul.NewSealerImpl(b.privateKey)
	types.SetHeaderHashFn(chainSealer.HeaderHash)
	bc, err := blockchain.NewBlockChain(dbm, nil, genesis.Config, chainSealer, vm.Config{})
	if err != nil {
		panic(err)
	}

	// kaiax module setup
	mGov := gov_impl.NewGovModule()
	mReward := reward_impl.NewRewardModule()
	mValset := valset_impl.NewValsetModule()
	mRandao := randao_impl.NewRandaoModule()
	mSystem := system_impl.NewSystemModule()
	if mStaking == nil {
		mStaking = staking_impl.NewStakingModule()
	}
	blsSecretKey, _ := bls.DeriveFromECDSA(b.privateKey)

	fakeDownloader := downloader.NewFakeDownloader()
	if err = errors.Join(
		mGov.Init(&gov_impl.InitOpts{
			Chain:       bc,
			ChainKv:     bc.StateCache().TrieDB().DiskDB().GetMiscDB(),
			ChainConfig: genesis.Config,
			Valset:      mValset,
			NodeAddress: b.address,
		}),
		mReward.Init(&reward_impl.InitOpts{
			ChainConfig:   bc.Config(),
			Chain:         bc,
			GovModule:     mGov,
			StakingModule: mStaking, // Irrelevant in ProposerPolicy=0. Won't inject mock.
			ValsetModule:  mValset,
		}),
		mValset.Init(&valset_impl.InitOpts{
			Chain:         bc,
			ChainKv:       bc.StateCache().TrieDB().DiskDB().GetMiscDB(),
			GovModule:     mGov,
			StakingModule: mStaking,
		}),
		mRandao.Init(&randao_impl.InitOpts{
			ChainConfig:  bc.Config(),
			Chain:        bc,
			Downloader:   fakeDownloader,
			BlsSecretKey: blsSecretKey,
		}),
		mSystem.Init(&system_impl.InitOpts{
			Chain: bc,
		}),
		func() error {
			if stakingImpl, ok := mStaking.(*staking_impl.StakingModule); ok {
				return stakingImpl.Init(&staking_impl.InitOpts{
					ChainKv:     bc.StateCache().TrieDB().DiskDB().GetMiscDB(),
					ChainConfig: genesis.Config,
					Chain:       bc,
				})
			}
			return nil
		}(),
	); err != nil {
		panic(err)
	}

	bc.RegisterKaiaxModules(
		mGov,
		mValset,
		nil,
		nil,
		[]kaiax.HeaderModule{mReward, mGov, mRandao},
		[]kaiax.BlockStateModule{mReward, mSystem},
	)
	b.RegisterKaiaxModules(mGov, mValset)

	if b.Start(bc, nil) != nil {
		panic(err)
	}

	testGovModules.Store(b, mGov)
	t.Cleanup(func() {
		testGovModules.Delete(b)
	})
	return bc, b
}

func makeHeader(parent *types.Block, chainConfig *params.ChainConfig) *types.Header {
	interval := uint64(params.BlockGenerationInterval)
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), new(big.Int).SetUint64(interval)),
		BlockScore: params.DefaultBlockScore,
	}
	if parent.Header().BaseFee != nil {
		// We don't have chainConfig so the BaseFee of the current block is set by parent's for test
		header.BaseFee = parent.Header().BaseFee
	}
	if chainConfig.IsOsakaForkEnabled(header.Number) {
		var excessBlobGas uint64
		if chainConfig.IsOsakaForkEnabled(parent.Number()) {
			excessBlobGas = chainConfig.LatestBlobConfig(header.Number).CalcExcessBlobGas(parent.ExcessBlobGas(), parent.BlobGasUsed())
		}
		header.BlobGasUsed = new(uint64)
		header.ExcessBlobGas = &excessBlobGas
	}
	return header
}

// makeBlockWithSeal creates a block with the proposer seal as well as all committed seals of validators.
func makeBlockWithSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	return sealBlock(engine, makeBlockWithoutSeal(chain, engine, parent))
}

func makeBlockWithoutSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	return makeBlockWithoutSealAndModifiedHeader(chain, engine, parent, nil)
}

// makeBlockWithoutSealAndModifiedHeader creates a block without seal, optionally with a modified header.
// The modifyHeader function is called before finalization.
func makeBlockWithoutSealAndModifiedHeader(chain *blockchain.BlockChain, engine *backend, parent *types.Block, modifyHeader func(*types.Header)) *types.Block {
	header := makeHeader(parent, chain.Config())
	if err := chain.PrepareHeader(header); err != nil {
		panic(err)
	}
	if modifyHeader != nil {
		modifyHeader(header)
	}
	state, _ := chain.StateAt(parent.Root())
	block, _ := chain.Processor().FinalizeState(header, state, nil, nil)
	return block
}

// sealBlock adds the proposer seal and committed seals to a block.
func sealBlock(engine *backend, blockWithoutSeal *types.Block) *types.Block {
	// add proposer seal for the block
	block, err := engine.updateBlock(blockWithoutSeal)
	if err != nil {
		panic(err)
	}

	// write validators committed seals to the block
	committedSeals := make([][]byte, len(nodeKeys))
	for i, key := range nodeKeys {
		committedSeals[i], err = istanbul.NewSealerImpl(key).MakeCommittedSeal(block.Header())
		if err != nil {
			panic(err)
		}
	}
	header := block.Header()
	if err := engine.sealer.WriteCommittedSeals(header, committedSeals); err != nil {
		panic(err)
	}
	return block.WithSeal(header)
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(t, 1)
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock, _ := engine.updateBlock(block)

	actualBlock, err := engine.Seal(chain, block)
	assert.NoError(t, err)
	assert.Equal(t, expectedBlock.Hash(), actualBlock.Hash())
}

func makeSnapshotTestConfigItems(stakingInterval, proposerInterval uint64) []interface{} {
	return []interface{}{
		stakingUpdateInterval(stakingInterval),
		proposerUpdateInterval(proposerInterval),
		proposerPolicy(params.WeightedRandom),
	}
}

func makeMockStakingManager(t *testing.T, amounts []uint64, blockNum uint64) (*gomock.Controller, *mock.MockStakingModule) {
	if len(nodeKeys) != len(amounts) {
		setNodeKeys(len(amounts), nil) // explicitly set the nodeKey
	}

	si := makeTestStakingInfo(amounts, blockNum)

	mockCtrl := gomock.NewController(t)
	mStaking := mock.NewMockStakingModule(mockCtrl)
	mStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(si, nil).AnyTimes()
	return mockCtrl, mStaking
}

func makeTestStakingInfo(amounts []uint64, blockNum uint64) *staking.StakingInfo {
	if amounts == nil {
		amounts = make([]uint64, len(nodeKeys))
	}
	si := &staking.StakingInfo{
		SourceBlockNum: blockNum,
	}
	for idx, key := range nodeKeys {
		addr := crypto.PubkeyToAddress(key.PublicKey)

		pk, _ := crypto.GenerateKey()
		rewardAddr := crypto.PubkeyToAddress(pk.PublicKey)

		si.NodeIds = append(si.NodeIds, addr)
		si.StakingContracts = append(si.StakingContracts, addr)
		si.RewardAddrs = append(si.RewardAddrs, rewardAddr)
		si.StakingAmounts = append(si.StakingAmounts, amounts[idx])
	}
	return si
}

func makeExpectedResult(indices []int, candidate []common.Address) []common.Address {
	expected := make([]common.Address, len(indices))
	for eIdx, cIdx := range indices {
		expected[eIdx] = candidate[cIdx]
	}
	return valset.NewAddressSet(expected).List()
}
