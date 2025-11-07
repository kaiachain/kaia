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
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	randao_impl "github.com/kaiachain/kaia/kaiax/randao/impl"
	reward_impl "github.com/kaiachain/kaia/kaiax/reward/impl"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_impl "github.com/kaiachain/kaia/kaiax/staking/impl"
	"github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	valset_impl "github.com/kaiachain/kaia/kaiax/valset/impl"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These variables are the global variables of the test blockchain.
var (
	nodeKeys []*ecdsa.PrivateKey
	addrs    []common.Address
)

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
	kaiaCompatibleBlock      *big.Int
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

// makeCommittedSeals returns a list of committed seals for the global variable nodeKeys.
func makeCommittedSeals(hash common.Hash) [][]byte {
	committedSeals := make([][]byte, len(nodeKeys))
	hashData := crypto.Keccak256(core.PrepareCommittedSeal(hash))
	for i, key := range nodeKeys {
		sig, _ := crypto.Sign(hashData, key)
		committedSeals[i] = make([]byte, types.IstanbulExtraSeal)
		copy(committedSeals[i][:], sig)
	}
	return committedSeals
}

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
	for i := 0; i < n; i++ {
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
func newBlockChain(n int, items ...interface{}) (*blockchain.BlockChain, *backend) {
	// generate a genesis block
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Config = params.TestChainConfig.Copy()
	genesis.Timestamp = uint64(time.Now().Unix())

	var (
		period   = istanbul.DefaultConfig.BlockPeriod
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
		case kaiaCompatibleBlock:
			genesis.Config.KaiaCompatibleBlock = v
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

	if len(nodeKeys) != n {
		setNodeKeys(n, nil)
	}

	// if governance mode is single, this address is the governing node address
	b := newTestBackendWithConfig(genesis.Config, period, nodeKeys[0])

	appendValidators(genesis, addrs)

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

	genesis.MustCommit(b.db)

	bc, err := blockchain.NewBlockChain(b.db, nil, genesis.Config, b, vm.Config{})
	if err != nil {
		panic(err)
	}

	// kaiax module setup
	mGov := gov_impl.NewGovModule()
	mReward := reward_impl.NewRewardModule()
	mValset := valset_impl.NewValsetModule()
	mRandao := randao_impl.NewRandaoModule()
	if mStaking == nil {
		mStaking = staking_impl.NewStakingModule()
	}

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
		}),
		mValset.Init(&valset_impl.InitOpts{
			Chain:         bc,
			ChainKv:       bc.StateCache().TrieDB().DiskDB().GetMiscDB(),
			GovModule:     mGov,
			StakingModule: mStaking,
		}),
		mRandao.Init(&randao_impl.InitOpts{
			ChainConfig: bc.Config(),
			Chain:       bc,
			Downloader:  fakeDownloader,
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

	b.RegisterKaiaxModules(mGov, mStaking, mValset, mRandao)
	b.RegisterConsensusModule(mReward, mGov)

	if b.Start(bc, bc.CurrentBlock, bc.HasBadBlock) != nil {
		panic(err)
	}

	return bc, b
}

func appendValidators(genesis *blockchain.Genesis, addrs []common.Address) {
	if len(genesis.ExtraData) < types.IstanbulExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func makeHeader(parent *types.Block, config *istanbul.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), new(big.Int).SetUint64(config.BlockPeriod)),
		BlockScore: defaultBlockScore,
	}
	if parent.Header().BaseFee != nil {
		// We don't have chainConfig so the BaseFee of the current block is set by parent's for test
		header.BaseFee = parent.Header().BaseFee
	}
	return header
}

func makeBlock(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	block := makeBlockWithoutSeal(chain, engine, parent)
	stopCh := make(chan struct{})
	result, err := engine.Seal(chain, block, stopCh)
	if err != nil {
		panic(err)
	}
	return result
}

// makeBlockWithSeal creates a block with the proposer seal as well as all committed seals of validators.
func makeBlockWithSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	blockWithoutSeal := makeBlockWithoutSeal(chain, engine, parent)

	// add proposer seal for the block
	block, err := engine.updateBlock(blockWithoutSeal)
	if err != nil {
		panic(err)
	}

	// write validators committed seals to the block
	header := block.Header()
	committedSeals := makeCommittedSeals(block.Hash())
	err = writeCommittedSeals(header, committedSeals)
	if err != nil {
		panic(err)
	}
	block = block.WithSeal(header)

	return block
}

func makeBlockWithoutSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	header := makeHeader(parent, engine.config)
	if err := engine.Prepare(chain, header); err != nil {
		panic(err)
	}
	state, _ := chain.StateAt(parent.Root())
	block, _ := engine.Finalize(chain, header, state, nil, nil)
	return block
}

func TestPrepare(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	header := makeHeader(chain.Genesis(), engine.config)
	err := engine.Prepare(chain, header)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	header.ParentHash = common.HexToHash("0x1234567890")
	err = engine.Prepare(chain, header)
	if err != consensus.ErrUnknownAncestor {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrUnknownAncestor)
	}
}

func TestSealStopChannel(t *testing.T) {
	chain, engine := newBlockChain(4)
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(istanbul.RequestEvent{})
	eventLoop := func() {
		select {
		case ev := <-eventSub.Chan():
			_, ok := ev.Data.(istanbul.RequestEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
			}
			stop <- struct{}{}
		}
		eventSub.Unsubscribe()
	}
	go eventLoop()

	finalBlock, err := engine.Seal(chain, block, stop)
	assert.NoError(t, err)
	assert.Nil(t, finalBlock)
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock, _ := engine.updateBlock(block)

	actualBlock, err := engine.Seal(chain, block, make(chan struct{}))
	assert.NoError(t, err)
	assert.Equal(t, expectedBlock.Hash(), actualBlock.Hash())
}

func TestVerifyHeader(t *testing.T) {
	var configItems []interface{}
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, LondonCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, EthTxTypeCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, magmaCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, koreCompatibleBlock(new(big.Int).SetUint64(0)))
	chain, engine := newBlockChain(1, configItems...)
	defer engine.Stop()

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	block, _ = engine.updateBlock(block)
	err := engine.VerifyHeader(chain, block.Header(), false)
	if err != errEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errEmptyCommittedSeals)
	}

	// short extra data
	header := block.Header()
	header.Extra = []byte{}
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}
	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.BlockScore = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidBlockScore {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidBlockScore)
	}

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = new(big.Int).Add(chain.Genesis().Time(), new(big.Int).SetUint64(engine.config.BlockPeriod-1))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(now().Unix()), new(big.Int).SetUint64(10))
	err = engine.VerifyHeader(chain, header, false)
	if err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}

	// TODO-Kaia: add more tests for header.Governance, header.Rewardbase, header.Vote
}

func TestVerifySeal(t *testing.T) {
	ctrl, mStaking := makeMockStakingManager(t, nil, 0)
	defer ctrl.Finish()
	chain, engine := newBlockChain(1, mStaking)
	defer engine.Stop()

	genesis := chain.Genesis()

	// cannot verify genesis
	err := engine.VerifySeal(chain, genesis.Header())
	if err != errUnknownBlock {
		t.Errorf("error mismatch: have %v, want %v", err, errUnknownBlock)
	}
	block := makeBlock(chain, engine, genesis)

	// clean cache before testing
	signatureAddresses.Purge()

	// change block content
	header := block.Header()
	header.Number = big.NewInt(4)
	block1 := block.WithSeal(header)
	err = engine.VerifySeal(chain, block1.Header())
	if err != errUnauthorized {
		t.Errorf("error mismatch: have %v, want %v", err, errUnauthorized)
	}

	// clean cache before testing
	signatureAddresses.Purge()

	// unauthorized users but still can get correct signer address
	engine.privateKey, _ = crypto.GenerateKey()
	err = engine.VerifySeal(chain, block.Header())
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

func TestVerifyHeaders(t *testing.T) {
	var configItems []interface{}
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block

	ctrl, mStaking := makeMockStakingManager(t, nil, 0)
	ctrl.Finish()

	chain, engine := newBlockChain(1, append(configItems, mStaking)...)
	chain.RegisterExecutionModule(engine.govModule)
	defer engine.Stop()

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	var previousBlock, currentBlock *types.Block = nil, chain.Genesis()
	for i := 0; i < size; i++ {
		// 100 headers with 50 of them empty committed seals, 50 of them invalid committed seals.
		previousBlock = currentBlock
		currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
		_, err := chain.InsertChain(types.Blocks{currentBlock})
		assert.NoError(t, err)

		blocks = append(blocks, currentBlock)
		headers = append(headers, blocks[i].Header())
	}

	// Set time to avoid future block errors
	now = func() time.Time {
		return time.Unix(headers[size-1].Time.Int64(), 0)
	}
	defer func() {
		now = time.Now
	}()

	// Helper function to verify headers and collect results
	verifyHeadersAndCollectResults := func(t *testing.T, testHeaders []*types.Header, expectErrors bool) int {
		abort, results := engine.VerifyHeaders(chain, testHeaders, nil)
		defer close(abort)

		timeout := time.NewTimer(2 * time.Second)
		defer timeout.Stop()

		index := 0
		errorCount := 0
		for {
			select {
			case err := <-results:
				if err != nil {
					errorCount++
					// These errors are expected in the test setup
					if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals && err != consensus.ErrUnknownAncestor {
						if !expectErrors {
							t.Errorf("unexpected error: %v", err)
						}
					}
				}
				index++
				if index == len(testHeaders) {
					return errorCount
				}
			case <-timeout.C:
				t.Error("timeout waiting for header verification results")
				return errorCount
			}
		}
	}

	// Test 1: Verify valid headers using VerifyHeaders (skip first header to avoid empty parents)
	t.Run("ValidHeaders", func(t *testing.T) {
		// Use headers[1:] to skip the first header that causes empty parents slice
		errorCount := verifyHeadersAndCollectResults(t, headers[1:], false)
		if errorCount > 0 {
			t.Logf("Valid headers test completed with %d expected errors", errorCount)
		}
	})

	// Test 2: Verify headers with invalid block number using VerifyHeaders
	t.Run("InvalidBlockNumber", func(t *testing.T) {
		// Create a copy of headers and modify one
		testHeaders := make([]*types.Header, len(headers)-1)
		copy(testHeaders, headers[1:])
		testHeaders[1].Number = big.NewInt(999999) // Invalid block number

		errorCount := verifyHeadersAndCollectResults(t, testHeaders, true)
		if errorCount == 0 {
			t.Error("expected errors for invalid block number, but got none")
		}
	})
}

func TestPrepareExtra(t *testing.T) {
	validators := make([]common.Address, 4)
	validators[0] = common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a"))
	validators[1] = common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212"))
	validators[2] = common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6"))
	validators[3] = common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440"))

	vanity := make([]byte, types.IstanbulExtraVanity)
	expectedResult := append(vanity, hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")...)

	h := &types.Header{
		Extra: vanity,
	}

	payload, err := prepareExtra(h, validators)
	if err != nil {
		t.Errorf("error mismatch: have %v, want: nil", err)
	}
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}

	// append useless information to extra-data
	h.Extra = append(vanity, make([]byte, 15)...)

	payload, err = prepareExtra(h, validators)
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}
}

func TestWriteSeal(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-3)...)
	expectedIstExtra := &types.IstanbulExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          expectedSeal,
		CommittedSeal: [][]byte{},
	}
	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := writeSeal(h, expectedSeal)
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify istanbul extra-data
	istExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedSeal := append(expectedSeal, make([]byte, 1)...)
	err = writeSeal(h, unexpectedSeal)
	if err != errInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSignature)
	}
}

func TestWriteCommittedSeals(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-3)...)
	expectedIstExtra := &types.IstanbulExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          []byte{},
		CommittedSeal: [][]byte{expectedCommittedSeal},
	}
	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := writeCommittedSeals(h, [][]byte{expectedCommittedSeal})
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify istanbul extra-data
	istExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	err = writeCommittedSeals(h, [][]byte{unexpectedCommittedSeal})
	if err != errInvalidCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidCommittedSeals)
	}
}

func TestRewardDistribution(t *testing.T) {
	type vote struct {
		name  string
		value interface{}
	}
	type expected = map[int]uint64
	type testcase struct {
		length   int // total number of blocks to simulate
		votes    map[int]vote
		expected expected
	}

	mintAmount := uint64(1)
	koreBlock := uint64(9)
	testEpoch := 3

	testcases := []testcase{
		{
			12,
			map[int]vote{
				1: {"reward.mintingamount", "2"}, // activated at block 7 (activation is before-Kore)
				4: {"reward.mintingamount", "3"}, // activated at block 9 (activation is after-Kore)
			},
			map[int]uint64{
				1:  1,
				2:  2,
				3:  3,
				4:  4,
				5:  5,
				6:  6,
				7:  8, // 2 is minted from now
				8:  10,
				9:  13, // 3 is minted from now
				10: 16,
				11: 19,
				12: 22,
				13: 25,
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, epoch(testEpoch))
	configItems = append(configItems, mintingAmount(new(big.Int).SetUint64(mintAmount)))
	configItems = append(configItems, koreCompatibleBlock(new(big.Int).SetUint64(koreBlock)))
	configItems = append(configItems, shanghaiCompatibleBlock(new(big.Int).SetUint64(koreBlock)))
	configItems = append(configItems, cancunCompatibleBlock(new(big.Int).SetUint64(koreBlock)))
	configItems = append(configItems, kaiaCompatibleBlock(new(big.Int).SetUint64(koreBlock)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block

	chain, engine := newBlockChain(1, configItems...)
	chain.RegisterExecutionModule(engine.govModule)
	defer engine.Stop()

	assert.Equal(t, uint64(testEpoch), engine.govModule.GetParamSet(0).Epoch)
	assert.Equal(t, mintAmount, engine.govModule.GetParamSet(0).MintingAmount.Uint64())

	var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

	for _, tc := range testcases {
		// Place a vote if a vote is scheduled in upcoming block
		// Note that we're building (head+1)'th block here.
		for num := 0; num <= tc.length; num++ {
			if vote, ok := tc.votes[num+1]; ok {
				v := headergov.NewVoteData(engine.address, vote.name, vote.value)
				require.NotNil(t, v, fmt.Sprintf("vote is nil for %v %v", vote.name, vote.value))
				engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(v)
			}

			// Create a block
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)

			// check balance
			addr := currentBlock.Rewardbase()
			state, err := chain.State()
			assert.NoError(t, err)
			bal := state.GetBalance(addr)

			assert.Equal(t, tc.expected[num+1], bal.Uint64(), "wrong at block %d", num+1)
		}
	}
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

// Asserts that if all (key,value) pairs of `subset` exists in `set`
func assertMapSubset[M ~map[K]any, K comparable](t *testing.T, subset, set M) {
	for k, v := range subset {
		assert.Equal(t, set[k], v)
	}
}

//func setEngineKeyAsProposer(engine *backend, num, round uint64) error {
//	cState, err := engine.GetCommitteeStateByRound(num, round)
//	if err != nil {
//		return err
//	}
//	for idx, addr := range addrs {
//		if addr == cState.Proposer() {
//			engine.privateKey = nodeKeys[idx]
//			engine.address = addr
//			break
//		}
//	}
//	return nil
//}

func Test_AfterMinimumStakingVotes(t *testing.T) {
	// temporarily enable forbidden votes
	enableVotes([]gov.ParamName{gov.RewardMinimumStake, gov.GovernanceGovernanceMode})
	defer disableVotes([]gov.ParamName{gov.RewardMinimumStake, gov.GovernanceGovernanceMode})

	type vote struct {
		key   string
		value interface{}
	}
	type expected struct {
		blocks     []uint64
		validators []int
		demoted    []int
	}
	type testcase struct {
		stakingAmounts []uint64
		votes          []vote
		expected       []expected
	}

	testcases := []testcase{
		{
			// test the validators are updated properly when minimum staking is changed in none mode
			[]uint64{8000000, 7000000, 6000000, 5000000},
			[]vote{
				{"governance.governancemode", "none"}, // voted on epoch 1, applied from 6-8
				{"reward.minimumstake", "5500000"},    // voted on epoch 2, applied from 9-11
				{"reward.minimumstake", "6500000"},    // voted on epoch 3, applied from 12-14
				{"reward.minimumstake", "7500000"},    // voted on epoch 4, applied from 15-17
				{"reward.minimumstake", "8500000"},    // voted on epoch 5, applied from 18-20
				{"reward.minimumstake", "7500000"},    // voted on epoch 6, applied from 21-23
				{"reward.minimumstake", "6500000"},    // voted on epoch 7, applied from 24-26
				{"reward.minimumstake", "5500000"},    // voted on epoch 8, applied from 27-29
				{"reward.minimumstake", "4500000"},    // voted on epoch 9, applied from 30-32
			},
			[]expected{
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
				{[]uint64{12, 13, 14}, []int{0, 1}, []int{2, 3}},
				{[]uint64{15, 16, 17}, []int{0}, []int{1, 2, 3}},
				{[]uint64{18, 19, 20}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{21, 22, 23}, []int{0}, []int{1, 2, 3}},
				{[]uint64{24, 25, 26}, []int{0, 1}, []int{2, 3}},
				{[]uint64{27, 28, 29}, []int{0, 1, 2}, []int{3}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			// test the validators (including governing node) are updated properly when minimum staking is changed in single mode
			[]uint64{5000000, 6000000, 7000000, 8000000},
			[]vote{
				{"reward.minimumstake", "8500000"}, // voted on epoch 1, applied from 6-8
				{"reward.minimumstake", "7500000"}, // voted on epoch 2, applied from 9-11
				{"reward.minimumstake", "6500000"}, // voted on epoch 3, applied from 12-14
				{"reward.minimumstake", "5500000"}, // voted on epoch 4, applied from 15-17
				{"reward.minimumstake", "4500000"}, // voted on epoch 5, applied from 18-20
				{"reward.minimumstake", "5500000"}, // voted on epoch 6, applied from 21-23
				{"reward.minimumstake", "6500000"}, // voted on epoch 7, applied from 24-26
				{"reward.minimumstake", "7500000"}, // voted on epoch 8, applied from 27-29
				{"reward.minimumstake", "8500000"}, // voted on epoch 9, applied from 30-32
			},
			[]expected{
				// 0 is governing node, so it is included in the validators all the time
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 3}, []int{1, 2}},
				{[]uint64{12, 13, 14}, []int{0, 2, 3}, []int{1}},
				{[]uint64{15, 16, 17, 18, 19, 20, 21, 22, 23}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{24, 25, 26}, []int{0, 2, 3}, []int{1}},
				{[]uint64{27, 28, 29}, []int{0, 3}, []int{1, 2}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			// test the validators are updated properly if governing node is changed
			[]uint64{6000000, 6000000, 5000000, 5000000},
			[]vote{
				{"reward.minimumstake", "5500000"}, // voted on epoch 1, applied from 6-8
				{"governance.governingnode", 2},    // voted on epoch 2, applied from 9-11
			},
			[]expected{
				// 0 is governing node, so it is included in the validators all the time
				{[]uint64{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{6, 7, 8}, []int{0, 1}, []int{2, 3}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
			},
		},
	}

	testEpoch := 3
	var configItems []interface{}
	configItems = append(configItems, params.TestKaiaConfig("magma"))
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(testEpoch))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block

	for _, tc := range testcases {
		ctrl, mStaking := makeMockStakingManager(t, tc.stakingAmounts, 0)
		chain, engine := newBlockChain(len(tc.stakingAmounts), append(configItems, mStaking)...)
		chain.RegisterExecutionModule(engine.govModule)

		var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

		for _, v := range tc.votes {
			// vote a vote in each epoch
			if v.key == "governance.governingnode" {
				idx := v.value.(int)
				v.value = addrs[idx].String()
			}
			// assert.NoError(t, setEngineKeyAsProposer(engine, currentBlock.NumberU64()+1, 0))

			vote := headergov.NewVoteData(engine.address, v.key, v.value)
			require.NotNil(t, vote, fmt.Sprintf("vote is nil for %v %v", v.key, v.value))
			engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(vote)

			for i := 0; i < testEpoch; i++ {
				previousBlock = currentBlock
				currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
				_, err := chain.InsertChain(types.Blocks{currentBlock})
				assert.NoError(t, err)
			}
		}

		// insert blocks on extra epoch
		for i := 0; i < 2*testEpoch; i++ {
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		for _, e := range tc.expected {
			for _, num := range e.blocks {
				valSet, err := engine.GetValidatorSet(num + 1)
				assert.NoError(t, err)

				expectedValidators := makeExpectedResult(e.validators, addrs)
				expectedDemoted := makeExpectedResult(e.demoted, addrs)

				assert.Equal(t, expectedValidators, valSet.Qualified().List(), "blockNum:%d", num+1)
				assert.Equal(t, expectedDemoted, valSet.Demoted().List(), "blockNum:%d", num+1)
			}
		}

		ctrl.Finish()
		engine.Stop()
	}
}

func Test_AfterKaia_BasedOnStaking(t *testing.T) {
	type testcase struct {
		stakingAmounts     []uint64 // test staking amounts of each validator
		isKaiaCompatible   bool     // whether or not if the inserted block is kaia compatible
		expectedValidators []int    // the indices of expected validators
		expectedDemoted    []int    // the indices of expected demoted validators
	}

	genesisStakingAmounts := []uint64{5000000, 5000000, 5000000, 5000000}

	testcases := []testcase{
		// The following testcases are the ones before kaia incompatible change
		// Validators doesn't be changed due to staking interval
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		// The following testcases are the ones after kaia incompatible change
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			true,
			[]int{3},
			[]int{0, 1, 2},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			true,
			[]int{2, 3},
			[]int{0, 1},
		},
	}

	testNum := 4
	ms := uint64(5500000)
	configItems := makeSnapshotTestConfigItems(10, 10)
	for _, tc := range testcases {
		if !tc.isKaiaCompatible {
			configItems = append(configItems, params.TestKaiaConfig("istanbul"))
		}
		configItems = append(configItems, minimumStake(new(big.Int).SetUint64(ms)))
		setNodeKeys(testNum, nil)
		mockCtrl := gomock.NewController(t)
		mStaking := mock.NewMockStakingModule(mockCtrl)

		mStaking.EXPECT().GetStakingInfo(uint64(0)).Return(makeTestStakingInfo(genesisStakingAmounts, 0), nil).AnyTimes()
		if tc.isKaiaCompatible {
			mStaking.EXPECT().GetStakingInfo(uint64(1)).Return(makeTestStakingInfo(genesisStakingAmounts, 0), nil).AnyTimes()
			mStaking.EXPECT().GetStakingInfo(uint64(2)).Return(makeTestStakingInfo(tc.stakingAmounts, 1), nil).AnyTimes()
		} else {
			mStaking.EXPECT().GetStakingInfo(uint64(1)).Return(makeTestStakingInfo(genesisStakingAmounts, 0), nil).AnyTimes()
			mStaking.EXPECT().GetStakingInfo(uint64(2)).Return(makeTestStakingInfo(genesisStakingAmounts, 0), nil).AnyTimes()
		}

		chain, engine := newBlockChain(testNum, append(configItems, mStaking)...)

		block := makeBlockWithSeal(chain, engine, chain.Genesis())
		_, err := chain.InsertChain(types.Blocks{block})
		assert.NoError(t, err)

		valSet, err := engine.GetValidatorSet(block.NumberU64() + 1)
		assert.NoError(t, err)

		expectedValidators := makeExpectedResult(tc.expectedValidators, addrs)
		expectedDemoted := makeExpectedResult(tc.expectedDemoted, addrs)

		assert.Equal(t, expectedValidators, valSet.Qualified().List())
		assert.Equal(t, expectedDemoted, valSet.Demoted().List())

		mockCtrl.Finish()
		engine.Stop()
	}
}

func Test_BasedOnStaking(t *testing.T) {
	type testcase struct {
		stakingAmounts       []uint64 // test staking amounts of each validator
		isIstanbulCompatible bool     // whether or not if the inserted block is istanbul compatible
		isSingleMode         bool     // whether or not if the governance mode is single
		expectedValidators   []int    // the indices of expected validators
		expectedDemoted      []int    // the indices of expected demoted validators
	}

	testcases := []testcase{
		// The following testcases are the ones before istanbul incompatible change
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 6000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		// The following testcases are the ones after istanbul incompatible change
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			true,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{6000000, 5000000, 5000000, 5000000},
			true,
			false,
			[]int{0},
			[]int{1, 2, 3},
		},
		{
			[]uint64{6000000, 5000000, 5000000, 6000000},
			true,
			false,
			[]int{0, 3},
			[]int{1, 2},
		},
		{
			[]uint64{6000000, 5000000, 6000000, 6000000},
			true,
			false,
			[]int{0, 2, 3},
			[]int{1},
		},
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			true,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5500001, 5500000, 5499999, 0},
			true,
			false,
			[]int{0, 1},
			[]int{2, 3},
		},
		// The following testcases are the ones for testing governing node in single mode
		// The first staking amount is of the governing node
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 6000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 2, 3},
			[]int{1},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			true,
			true,
			[]int{0, 3},
			[]int{1, 2},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
	}

	ms := uint64(5500000)
	configItems := makeSnapshotTestConfigItems(1, 1)
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(ms)))
	configItems = append(configItems, LondonCompatibleBlock(nil))
	configItems = append(configItems, EthTxTypeCompatibleBlock(nil))
	configItems = append(configItems, magmaCompatibleBlock(nil))
	configItems = append(configItems, koreCompatibleBlock(nil))
	configItems = append(configItems, shanghaiCompatibleBlock(nil))
	configItems = append(configItems, cancunCompatibleBlock(nil))
	configItems = append(configItems, kaiaCompatibleBlock(nil))

	for _, tc := range testcases {
		if !tc.isIstanbulCompatible {
			configItems = append(configItems, istanbulCompatibleBlock(nil))
		} else {
			configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
		}
		if tc.isSingleMode {
			configItems = append(configItems, governanceMode("single"))
		}
		mockCtrl, mStaking := makeMockStakingManager(t, tc.stakingAmounts, 0)
		if tc.isIstanbulCompatible {
			mStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(makeTestStakingInfo(tc.stakingAmounts, 0), nil).AnyTimes()
		}
		chain, engine := newBlockChain(len(tc.stakingAmounts), append(configItems, mStaking)...)

		block := makeBlockWithSeal(chain, engine, chain.Genesis())
		_, err := chain.InsertChain(types.Blocks{block})
		assert.NoError(t, err)

		councilState, err := engine.GetValidatorSet(block.NumberU64() + 1)
		assert.NoError(t, err)

		expectedValidators := makeExpectedResult(tc.expectedValidators, addrs)
		expectedDemoted := makeExpectedResult(tc.expectedDemoted, addrs)

		assert.Equal(t, expectedValidators, councilState.Qualified().List())
		assert.Equal(t, expectedDemoted, councilState.Demoted().List())

		mockCtrl.Finish()
		engine.Stop()
	}
}

func Test_AddRemove(t *testing.T) {
	type vote struct {
		key   string
		value interface{}
	}
	type expected struct {
		validators []int // expected validator indexes at given block
	}
	type testcase struct {
		length   int // total number of blocks to simulate
		votes    map[int]vote
		expected map[int]expected
	}

	testcases := []testcase{
		{ // Singular change
			5,
			map[int]vote{
				1: {"governance.removevalidator", 3},
				3: {"governance.addvalidator", 3},
			},
			map[int]expected{
				0: {[]int{0, 1, 2, 3}},
				1: {[]int{0, 1, 2, 3}},
				2: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2}},
				4: {[]int{0, 1, 2, 3}},
			},
		},
		{ // Plural change
			5,
			map[int]vote{
				1: {"governance.removevalidator", []int{1, 2, 3}},
				3: {"governance.addvalidator", []int{1, 2}},
			},
			map[int]expected{
				0: {[]int{0, 1, 2, 3}},
				1: {[]int{0, 1, 2, 3}},
				2: {[]int{0}},
				3: {[]int{0}},
				4: {[]int{0, 1, 2}},
			},
		},
		{ // Around checkpoint interval (i.e. every 1024 block)
			params.CheckpointInterval + 10,
			map[int]vote{
				params.CheckpointInterval - 5: {"governance.removevalidator", 3},
				params.CheckpointInterval - 1: {"governance.removevalidator", 2},
				params.CheckpointInterval + 0: {"governance.removevalidator", 1},
				params.CheckpointInterval + 1: {"governance.addvalidator", 1},
				params.CheckpointInterval + 2: {"governance.addvalidator", 2},
				params.CheckpointInterval + 3: {"governance.addvalidator", 3},
			},
			map[int]expected{
				0:                             {[]int{0, 1, 2, 3}},
				1:                             {[]int{0, 1, 2, 3}},
				params.CheckpointInterval - 4: {[]int{0, 1, 2}},
				params.CheckpointInterval + 0: {[]int{0, 1}},
				params.CheckpointInterval + 1: {[]int{0}},
				params.CheckpointInterval + 2: {[]int{0, 1}},
				params.CheckpointInterval + 3: {[]int{0, 1, 2}},
				params.CheckpointInterval + 4: {[]int{0, 1, 2, 3}},
				params.CheckpointInterval + 9: {[]int{0, 1, 2, 3}},
			},
		},
		{ // multiple addvalidator & removevalidator
			10,
			map[int]vote{
				0: {"governance.removevalidator", 3},
				2: {"governance.addvalidator", 3},
				4: {"governance.addvalidator", 3},
				6: {"governance.removevalidator", 3},
				8: {"governance.removevalidator", 3},
			},
			map[int]expected{
				1: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2, 3}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2}},
				9: {[]int{0, 1, 2}},
			},
		},
		{ // multiple removevalidator & addvalidator
			10,
			map[int]vote{
				0: {"governance.removevalidator", 3},
				2: {"governance.removevalidator", 3},
				4: {"governance.addvalidator", 3},
				6: {"governance.addvalidator", 3},
			},
			map[int]expected{
				1: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2, 3}},
			},
		},
		{ // multiple addvalidators & removevalidators
			10,
			map[int]vote{
				0: {"governance.removevalidator", []int{2, 3}},
				2: {"governance.addvalidator", []int{2, 3}},
				4: {"governance.addvalidator", []int{2, 3}},
				6: {"governance.removevalidator", []int{2, 3}},
				8: {"governance.removevalidator", []int{2, 3}},
			},
			map[int]expected{
				1: {[]int{0, 1}},
				3: {[]int{0, 1, 2, 3}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1}},
				9: {[]int{0, 1}},
			},
		},
		{ // multiple removevalidators & addvalidators
			10,
			map[int]vote{
				0: {"governance.removevalidator", []int{2, 3}},
				2: {"governance.removevalidator", []int{2, 3}},
				4: {"governance.addvalidator", []int{2, 3}},
				6: {"governance.addvalidator", []int{2, 3}},
			},
			map[int]expected{
				1: {[]int{0, 1}},
				3: {[]int{0, 1}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2, 3}},
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, subGroupSize(4))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	stakes := []uint64{4000000, 4000000, 4000000, 4000000}

	for _, tc := range testcases {
		// Create test blockchain
		ctrl, mStaking := makeMockStakingManager(t, stakes, 0)
		chain, engine := newBlockChain(len(stakes), append(configItems, mStaking)...)
		chain.RegisterExecutionModule(engine.valsetModule, engine.govModule)

		// Backup the globals. The globals `nodeKeys` and `addrs` will be
		// modified according to validator change votes.
		allNodeKeys := make([]*ecdsa.PrivateKey, len(nodeKeys))
		allAddrs := make([]common.Address, len(addrs))
		copy(allNodeKeys, nodeKeys)
		copy(allAddrs, addrs)

		var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

		// Create blocks with votes
		for i := 0; i < tc.length; i++ {
			if v, ok := tc.votes[i]; ok { // If a vote is scheduled in this block,
				if idx, ok := v.value.(int); ok {
					addr := allAddrs[idx]
					vote := headergov.NewVoteData(engine.address, v.key, addr)
					require.NotNil(t, vote, fmt.Sprintf("vote is nil for %v %v", v.key, addr))
					engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(vote)
				} else {
					addrList := makeExpectedResult(v.value.([]int), allAddrs)
					vote := headergov.NewVoteData(engine.address, v.key, addrList)
					require.NotNil(t, vote, fmt.Sprintf("vote is nil for %v %v", v.key, addrList))
					engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(vote)
				}
				// t.Logf("Voting at block #%d for %s, %v", i, v.key, v.value)
			}

			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)

			// After a voting, reflect the validator change to the globals
			if v, ok := tc.votes[i]; ok {
				var indices []int
				if idx, ok := v.value.(int); ok {
					indices = []int{idx}
				} else {
					indices = v.value.([]int)
				}
				if v.key == "governance.addvalidator" {
					for _, i := range indices {
						includeNode(allAddrs[i], allNodeKeys[i])
					}
				}
				if v.key == "governance.removevalidator" {
					for _, i := range indices {
						excludeNodeByAddr(allAddrs[i])
					}
				}
			}
		}

		// Calculate historical validators
		for i := 0; i < tc.length; i++ {
			if _, ok := tc.expected[i]; !ok {
				continue
			}
			valSet, err := engine.GetValidatorSet(uint64(i) + 1)
			assert.NoError(t, err)

			expectedValidators := makeExpectedResult(tc.expected[i].validators, allAddrs)
			assert.Equal(t, expectedValidators, valSet.Qualified().List())
		}

		ctrl.Finish()
		engine.Stop()
	}
}

func TestGovernance_Votes(t *testing.T) {
	enableVotes([]gov.ParamName{gov.RewardMinimumStake, gov.GovernanceGovernanceMode, gov.RewardUseGiniCoeff})
	defer disableVotes([]gov.ParamName{gov.RewardMinimumStake, gov.GovernanceGovernanceMode, gov.RewardUseGiniCoeff})

	type vote struct {
		key   string
		value interface{}
	}
	type governanceItem struct {
		vote
		appliedBlockNumber uint64 // if applied block number is 0, then it checks the item on current block
	}
	type testcase struct {
		votes    []vote
		expected []governanceItem
	}

	testcases := []testcase{
		{
			votes: []vote{
				{"governance.governancemode", "none"},     // voted on block 1
				{"istanbul.committeesize", uint64(4)},     // voted on block 2
				{"governance.unitprice", uint64(2000000)}, // voted on block 3
				{"reward.mintingamount", "96000000000"},   // voted on block 4
				{"reward.ratio", "34/33/33"},              // voted on block 5
				{"reward.useginicoeff", true},             // voted on block 6
				{"reward.minimumstake", "5000000"},        // voted on block 7
				{"reward.kip82ratio", "50/50"},            // voted on block 8
				{"governance.deriveshaimpl", uint64(2)},   // voted on block 9
			},
			expected: []governanceItem{
				{vote{"governance.governancemode", "none"}, 6},
				{vote{"istanbul.committeesize", uint64(4)}, 6},
				{vote{"governance.unitprice", uint64(2000000)}, 9},
				{vote{"reward.mintingamount", "96000000000"}, 9},
				{vote{"reward.ratio", "34/33/33"}, 9},
				{vote{"reward.useginicoeff", true}, 12},
				{vote{"reward.minimumstake", "5000000"}, 12},
				{vote{"reward.kip82ratio", "50/50"}, 12},
				{vote{"governance.deriveshaimpl", uint64(2)}, 15},
				// check governance items on current block
				{vote{"governance.governancemode", "none"}, 0},
				{vote{"istanbul.committeesize", uint64(4)}, 0},
				{vote{"governance.unitprice", uint64(2000000)}, 0},
				{vote{"reward.mintingamount", "96000000000"}, 0},
				{vote{"reward.ratio", "34/33/33"}, 0},
				{vote{"reward.useginicoeff", true}, 0},
				{vote{"reward.minimumstake", "5000000"}, 0},
				{vote{"reward.kip82ratio", "50/50"}, 0},
				{vote{"governance.deriveshaimpl", uint64(2)}, 0},
			},
		},
		{
			votes: []vote{
				{"governance.governancemode", "none"},   // voted on block 1
				{"governance.governancemode", "single"}, // voted on block 2
				{"governance.governancemode", "none"},   // voted on block 3
				{"governance.governancemode", "single"}, // voted on block 4
				{"governance.governancemode", "none"},   // voted on block 5
				{"governance.governancemode", "single"}, // voted on block 6
				{"governance.governancemode", "none"},   // voted on block 7
				{"governance.governancemode", "single"}, // voted on block 8
				{"governance.governancemode", "none"},   // voted on block 9
			},
			expected: []governanceItem{
				{vote{"governance.governancemode", "single"}, 6},
				{vote{"governance.governancemode", "none"}, 9},
				{vote{"governance.governancemode", "single"}, 12},
				{vote{"governance.governancemode", "none"}, 15},
			},
		},
		{
			votes: []vote{
				{"governance.governancemode", "none"},     // voted on block 1
				{"istanbul.committeesize", uint64(4)},     // voted on block 2
				{"governance.unitprice", uint64(2000000)}, // voted on block 3
				{"governance.governancemode", "single"},   // voted on block 4
				{"istanbul.committeesize", uint64(22)},    // voted on block 5
				{"governance.unitprice", uint64(2)},       // voted on block 6
				{"governance.governancemode", "none"},     // voted on block 7
			},
			expected: []governanceItem{
				// governance mode for all blocks
				{vote{"governance.governancemode", "single"}, 1},
				{vote{"governance.governancemode", "single"}, 2},
				{vote{"governance.governancemode", "single"}, 3},
				{vote{"governance.governancemode", "single"}, 4},
				{vote{"governance.governancemode", "single"}, 5},
				{vote{"governance.governancemode", "none"}, 6},
				{vote{"governance.governancemode", "none"}, 7},
				{vote{"governance.governancemode", "none"}, 8},
				{vote{"governance.governancemode", "single"}, 9},
				{vote{"governance.governancemode", "single"}, 10},
				{vote{"governance.governancemode", "single"}, 11},
				{vote{"governance.governancemode", "none"}, 12},
				{vote{"governance.governancemode", "none"}, 13},
				{vote{"governance.governancemode", "none"}, 14},
				{vote{"governance.governancemode", "none"}, 0}, // check on current

				// committee size for all blocks
				{vote{"istanbul.committeesize", uint64(21)}, 1},
				{vote{"istanbul.committeesize", uint64(21)}, 2},
				{vote{"istanbul.committeesize", uint64(21)}, 3},
				{vote{"istanbul.committeesize", uint64(21)}, 4},
				{vote{"istanbul.committeesize", uint64(21)}, 5},
				{vote{"istanbul.committeesize", uint64(4)}, 6},
				{vote{"istanbul.committeesize", uint64(4)}, 7},
				{vote{"istanbul.committeesize", uint64(4)}, 8},
				{vote{"istanbul.committeesize", uint64(22)}, 9},
				{vote{"istanbul.committeesize", uint64(22)}, 10},
				{vote{"istanbul.committeesize", uint64(22)}, 11},
				{vote{"istanbul.committeesize", uint64(22)}, 12},
				{vote{"istanbul.committeesize", uint64(22)}, 13},
				{vote{"istanbul.committeesize", uint64(22)}, 14},
				{vote{"istanbul.committeesize", uint64(22)}, 0}, // check on current

				// unitprice for all blocks
				{vote{"governance.unitprice", uint64(1)}, 1},
				{vote{"governance.unitprice", uint64(1)}, 2},
				{vote{"governance.unitprice", uint64(1)}, 3},
				{vote{"governance.unitprice", uint64(1)}, 4},
				{vote{"governance.unitprice", uint64(1)}, 5},
				{vote{"governance.unitprice", uint64(1)}, 6},
				{vote{"governance.unitprice", uint64(1)}, 7},
				{vote{"governance.unitprice", uint64(1)}, 8},
				{vote{"governance.unitprice", uint64(2000000)}, 9},
				{vote{"governance.unitprice", uint64(2000000)}, 10},
				{vote{"governance.unitprice", uint64(2000000)}, 11},
				{vote{"governance.unitprice", uint64(2)}, 12},
				{vote{"governance.unitprice", uint64(2)}, 13},
				{vote{"governance.unitprice", uint64(2)}, 14},
				{vote{"governance.unitprice", uint64(2)}, 0}, // check on current
			},
		},
		{
			votes: []vote{
				{}, // voted on block 1
				{"kip71.lowerboundbasefee", uint64(750000000000)}, // voted on block 2
				{}, // voted on block 3
				{}, // voted on block 4
				{"kip71.lowerboundbasefee", uint64(25000000000)}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"kip71.lowerboundbasefee", uint64(750000000000)}, 6},
				{vote{"kip71.lowerboundbasefee", uint64(25000000000)}, 9},
			},
		},
		{
			votes: []vote{
				{}, // voted on block 1
				{"kip71.upperboundbasefee", uint64(750000000000)}, // voted on block 2
				{}, // voted on block 3
				{}, // voted on block 4
				{"kip71.upperboundbasefee", uint64(25000000000)}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"kip71.upperboundbasefee", uint64(750000000000)}, 6},
				{vote{"kip71.upperboundbasefee", uint64(25000000000)}, 9},
			},
		},
		{
			votes: []vote{
				{}, // voted on block 1
				{"kip71.maxblockgasusedforbasefee", uint64(840000000)}, // voted on block 2
				{}, // voted on block 3
				{}, // voted on block 4
				{"kip71.maxblockgasusedforbasefee", uint64(84000000)}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"kip71.maxblockgasusedforbasefee", uint64(840000000)}, 6},
				{vote{"kip71.maxblockgasusedforbasefee", uint64(84000000)}, 9},
			},
		},
		{
			votes: []vote{
				{},                                    // voted on block 1
				{"kip71.gastarget", uint64(50000000)}, // voted on block 2
				{},                                    // voted on block 3
				{},                                    // voted on block 4
				{"kip71.gastarget", uint64(30000000)}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"kip71.gastarget", uint64(50000000)}, 6},
				{vote{"kip71.gastarget", uint64(30000000)}, 9},
			},
		},
		{
			votes: []vote{
				{},                                       // voted on block 1
				{"kip71.basefeedenominator", uint64(32)}, // voted on block 2
				{},                                       // voted on block 3
				{},                                       // voted on block 4
				{"kip71.basefeedenominator", uint64(64)}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"kip71.basefeedenominator", uint64(32)}, 6},
				{vote{"kip71.basefeedenominator", uint64(64)}, 9},
			},
		},
		{
			votes: []vote{
				{}, // voted on block 1
				{"governance.govparamcontract", common.HexToAddress("0x0000000000000000000000000000000000000000")}, // voted on block 2
				{}, // voted on block 3
				{}, // voted on block 4
				{"governance.govparamcontract", common.HexToAddress("0x0000000000000000000000000000000000000400")}, // voted on block 5
			},
			expected: []governanceItem{
				{vote{"governance.govparamcontract", common.HexToAddress("0x0000000000000000000000000000000000000000")}, 6},
				{vote{"governance.govparamcontract", common.HexToAddress("0x0000000000000000000000000000000000000400")}, 9},
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, params.TestKaiaConfig("ethTxType"))
	configItems = append(configItems, lowerBoundBaseFee(25000000000))
	configItems = append(configItems, upperBoundBaseFee(750000000000))
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	for _, tc := range testcases {
		mockCtrl, mStaking := makeMockStakingManager(t, nil, 0)
		chain, engine := newBlockChain(1, append(configItems, mStaking)...)
		chain.RegisterExecutionModule(engine.valsetModule, engine.govModule)

		// test initial governance items
		pset := engine.govModule.GetParamSet(chain.CurrentHeader().Number.Uint64() + 1)
		assert.Equal(t, uint64(3), pset.Epoch)
		assert.Equal(t, "single", pset.GovernanceMode)
		assert.Equal(t, uint64(21), pset.CommitteeSize)
		assert.Equal(t, uint64(1), pset.UnitPrice)
		assert.Equal(t, "0", pset.MintingAmount.String())
		assert.Equal(t, "100/0/0", pset.Ratio)
		assert.Equal(t, false, pset.UseGiniCoeff)
		assert.Equal(t, "2000000", pset.MinimumStake.String())

		// add votes and insert voted blocks
		var (
			previousBlock, currentBlock *types.Block = nil, chain.Genesis()
			err                         error
		)

		for _, v := range tc.votes {
			if len(v.key) > 0 {
				vote := headergov.NewVoteData(engine.address, v.key, v.value)
				require.NotNil(t, vote, fmt.Sprintf("vote is nil for %v %v", v.key, v.value))
				engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(vote)
				t.Logf("Adding vote(%s,%v) at block %d", v.key, v.value, currentBlock.NumberU64()+1)
			}

			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err = chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		// insert blocks until the vote is applied
		for i := 0; i < 6; i++ {
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err = chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		for _, item := range tc.expected {
			blockNumber := item.appliedBlockNumber
			if blockNumber == 0 {
				blockNumber = chain.CurrentBlock().NumberU64()
			}
			partialParamSet := engine.govModule.(*gov_impl.GovModule).Hgm.GetPartialParamSet(blockNumber + 1)
			switch val := partialParamSet[gov.ParamName(item.key)]; v := val.(type) {
			case *big.Int:
				require.Equal(t, item.value, v.String())
			default:
				require.Equal(t, item.value, v)
			}
		}

		mockCtrl.Finish()
		engine.Stop()
	}
}

func TestGovernance_GovModule(t *testing.T) {
	// Test that ReaderEngine (CurrentParams(), GetParamSet()) works.
	type vote struct {
		name  string
		value interface{}
	}
	type expected = map[string]interface{} // expected (subset of) governance items
	type testcase struct {
		length   int // total number of blocks to simulate
		votes    map[int]vote
		expected map[int]expected
	}

	testcases := []testcase{
		{
			8,
			map[int]vote{
				1: {"governance.unitprice", uint64(17)},
			},
			map[int]expected{
				0: {"governance.unitprice": uint64(1)},
				1: {"governance.unitprice": uint64(1)},
				2: {"governance.unitprice": uint64(1)},
				3: {"governance.unitprice": uint64(1)},
				4: {"governance.unitprice": uint64(1)},
				5: {"governance.unitprice": uint64(1)},
				6: {"governance.unitprice": uint64(1)},
				7: {"governance.unitprice": uint64(17)},
				8: {"governance.unitprice": uint64(17)},
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, params.TestKaiaConfig("ethTxType"))
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	stakes := []uint64{4000000, 4000000, 4000000, 4000000}

	for _, tc := range testcases {
		// Create test blockchain
		mockCtrl, mStaking := makeMockStakingManager(t, stakes, 0)
		chain, engine := newBlockChain(len(stakes), append(configItems, mStaking)...)
		chain.RegisterExecutionModule(engine.valsetModule, engine.govModule)

		var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

		// Create blocks with votes
		for num := 0; num <= tc.length; num++ {
			// Validate current params with CurrentParams() and CurrentSetCopy().
			// Check that both returns the expected result.
			pset := engine.govModule.GetParamSet(uint64(num + 1))
			assertMapSubset(t, tc.expected[num+1], pset.ToGovParamSet().StrMap())

			// Place a vote if a vote is scheduled in upcoming block
			// Note that we're building (head+1)'th block here.
			if vote, ok := tc.votes[num+1]; ok {
				v := headergov.NewVoteData(engine.address, vote.name, vote.value)
				require.NotNil(t, v, fmt.Sprintf("vote is nil for %v %v", vote.name, vote.value))
				engine.govModule.(*gov_impl.GovModule).Hgm.PushMyVotes(v)
			}

			// Create a block
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		// Validate historic parameters with GetParamSet() and GetPartialParamSet().
		// Check that both returns the expected result.
		for num := 0; num <= tc.length; num++ {
			pset := engine.govModule.GetParamSet(uint64(num))
			assertMapSubset(t, tc.expected[num], pset.ToGovParamSet().StrMap())

			partialParamSet := make(map[string]any)
			for k, v := range engine.govModule.(*gov_impl.GovModule).Hgm.GetPartialParamSet(uint64(num + 1)) {
				partialParamSet[string(k)] = v
			}
			assertMapSubset(t, tc.expected[num+1], partialParamSet)
		}

		mockCtrl.Finish()
		engine.Stop()
	}
}
