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
// This file is derived from quorum/consensus/istanbul/backend/backend_test.go (2020/04/16).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

var (
	testCommitteeSize = uint64(21)
	testSigningData   = []byte("dummy data")
	maxBlockNum       = int64(100)
)

type keys []*ecdsa.PrivateKey

func (slice keys) Len() int {
	return len(slice)
}

func (slice keys) Less(i, j int) bool {
	return strings.Compare(crypto.PubkeyToAddress(slice[i].PublicKey).String(), crypto.PubkeyToAddress(slice[j].PublicKey).String()) < 0
}

func (slice keys) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func getTestConfig() *params.ChainConfig {
	config := params.TestChainConfig.Copy()
	config.Governance = params.GetDefaultGovernanceConfig()
	config.Istanbul = params.GetDefaultIstanbulConfig()
	return config
}

// TestBackend_GetTargetReceivers checks if the gossiping targets are same as council members
func TestBackend_GetTargetReceivers(t *testing.T) {
	stakes := []uint64{5000000, 5000000, 5000000, 5000000}
	ctrl, mStaking := makeMockStakingManager(t, stakes, 0)
	defer ctrl.Finish()

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(5)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	configItems = append(configItems, mStaking)

	chain, istBackend := newBlockChain(len(stakes), configItems...)
	chain.RegisterExecutionModule(istBackend.govModule)
	defer istBackend.Stop()

	// Test for blocks from 0 to maxBlockNum
	// from 0 to 4: before istanbul hard fork
	// from 5 to 100: after istanbul hard fork
	var previousBlock, currentBlock *types.Block = nil, chain.CurrentBlock()
	for i := int64(0); i < maxBlockNum; i++ {
		previousBlock = currentBlock
		currentBlock = makeBlockWithSeal(chain, istBackend, previousBlock)
		_, err := chain.InsertChain(types.Blocks{currentBlock})
		assert.NoError(t, err)
	}

	for i := int64(0); i < maxBlockNum; i++ {
		// Test for round 0 to round 14
		for round := int64(0); round < 14; round++ {
			currentCouncilState, err := istBackend.GetCommitteeStateByRound(uint64(i), uint64(round))
			assert.NoError(t, err)

			// skip if the testing node is not in a committee
			isInSubList := currentCouncilState.Committee().Contains(istBackend.Address())
			if isInSubList == false {
				continue
			}

			nextCouncilState, err := istBackend.GetCommitteeStateByRound(uint64(i), uint64(round)+1)
			assert.NoError(t, err)

			// Receiving the receiver list of a message
			targets := istBackend.getTargetReceivers()
			targetAddrs := new(valset.AddressSet)
			for addr := range targets {
				targetAddrs.Add(addr)
			}

			// number of message receivers have to be smaller than or equal to the number of the current committee and the next committee
			// committees[0]: current round's committee
			// committees[1]: next view's committee
			committees := make([]*valset.AddressSet, 2)
			committees[0], committees[1] = currentCouncilState.Committee(), nextCouncilState.Committee()
			assert.True(t, len(targets) <= committees[0].Len()+committees[1].Len())

			// Check all nodes in the current and the next round are included in the target list
			assert.Equal(t, 0, committees[0].Subtract(targetAddrs).Subtract(valset.NewAddressSet([]common.Address{istBackend.Address()})).Len())
			assert.Equal(t, 0, committees[1].Subtract(targetAddrs).Subtract(valset.NewAddressSet([]common.Address{istBackend.Address()})).Len())

			// Check if a validator not in the current/next committee is included in target list
			assert.Equal(t, 0, targetAddrs.Subtract(committees[0]).Subtract(committees[1]).Len())
		}
	}
}

func newTestBackend() (b *backend) {
	config := getTestConfig()
	config.Istanbul.ProposerPolicy = params.WeightedRandom
	return newTestBackendWithConfig(config, istanbul.DefaultConfig.BlockPeriod, nil)
}

func newTestBackendWithConfig(chainConfig *params.ChainConfig, blockPeriod uint64, key *ecdsa.PrivateKey) (b *backend) {
	dbm := database.NewMemoryDBManager()
	if key == nil {
		// if key is nil, generate new key for a test account
		key, _ = crypto.GenerateKey()
	}
	if chainConfig.Governance.GovernanceMode == "single" {
		// if governance mode is single, set the node key to the governing node.
		chainConfig.Governance.GoverningNode = crypto.PubkeyToAddress(key.PublicKey)
	}
	g := governance.NewMixedEngine(chainConfig, dbm)
	istanbulConfig := &istanbul.Config{
		Epoch:          chainConfig.Istanbul.Epoch,
		ProposerPolicy: istanbul.ProposerPolicy(chainConfig.Istanbul.ProposerPolicy),
		SubGroupSize:   chainConfig.Istanbul.SubGroupSize,
		BlockPeriod:    blockPeriod,
		Timeout:        10000,
	}

	backend := New(&BackendOpts{
		IstanbulConfig: istanbulConfig,
		Rewardbase:     common.HexToAddress("0x2A35FE72F847aa0B509e4055883aE90c87558AaD"),
		PrivateKey:     key,
		DB:             dbm,
		Governance:     g,
		NodeType:       common.CONSENSUSNODE,
	}).(*backend)
	return backend
}

func TestSign(t *testing.T) {
	b := newTestBackend()
	defer b.Stop()

	sig, err := b.Sign(testSigningData)
	assert.NoError(t, err)

	// Check signature recover
	hashData := crypto.Keccak256([]byte(testSigningData))
	pubkey, _ := crypto.Ecrecover(hashData, sig)
	actualSigner := common.BytesToAddress(crypto.Keccak256(pubkey[1:])[12:])
	assert.Equal(t, b.address, actualSigner)
}

func TestCheckSignature(t *testing.T) {
	b := newTestBackend()
	defer b.Stop()

	// testAddr is derived from testPrivateKey.
	testPrivateKey, _ := crypto.HexToECDSA("bb047e5940b6d83354d9432db7c449ac8fca2248008aaa7271369880f9f11cc1")
	testAddr := common.HexToAddress("0x70524d664ffe731100208a0154e556f9bb679ae6")
	testInvalidAddr := common.HexToAddress("0x9535b2e7faaba5288511d89341d94a38063a349b")

	hashData := crypto.Keccak256([]byte(testSigningData))
	sig, err := crypto.Sign(hashData, testPrivateKey)
	assert.NoError(t, err)

	assert.NoError(t, b.CheckSignature(testSigningData, testAddr, sig))
	assert.Equal(t, errInvalidSignature, b.CheckSignature(testSigningData, testInvalidAddr, sig))
}

func TestCheckValidatorSignature(t *testing.T) {
	// generate validators
	setNodeKeys(5, nil)
	valSet := istanbul.NewBlockValSet(addrs, []common.Address{})

	// 1. Positive test: sign with validator's key should succeed
	hashData := crypto.Keccak256(testSigningData)
	for i, k := range nodeKeys {
		// Sign
		sig, err := crypto.Sign(hashData, k)
		assert.NoError(t, err)
		// CheckValidatorSignature should succeed
		addr, err := valSet.CheckValidatorSignature(testSigningData, sig)
		assert.NoError(t, err)
		assert.Equal(t, addrs[i], addr)
	}

	// 2. Negative test: sign with any key other than validator's key should return error
	key, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Sign
	sig, err := crypto.Sign(hashData, key)
	assert.NoError(t, err)
	// CheckValidatorSignature should return ErrUnauthorizedAddress
	addr, err := valSet.CheckValidatorSignature(testSigningData, sig)
	assert.Equal(t, istanbul.ErrUnauthorizedAddress, err)
	assert.True(t, common.EmptyAddress(addr))
}

func TestCommit(t *testing.T) {
	commitCh := make(chan *types.Block)

	// Case: it's a proposer, so the backend.commit will receive channel result from backend.Commit function
	for _, test := range []struct {
		expectedErr       error
		expectedSignature [][]byte
	}{
		{
			// normal case
			nil,
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-1)...)},
		},
		{
			// invalid signature
			errInvalidCommittedSeals,
			nil,
		},
	} {
		chain, engine := newBlockChain(1)

		block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
		expBlock, _ := engine.updateBlock(block)

		go func() {
			select {
			case result := <-engine.commitCh:
				commitCh <- result.Block
				return
			}
		}()

		engine.proposedBlockHash = expBlock.Hash()
		assert.Equal(t, test.expectedErr, engine.Commit(expBlock, test.expectedSignature))

		if test.expectedErr == nil {
			// to avoid race condition is occurred by goroutine
			select {
			case result := <-commitCh:
				assert.Equal(t, expBlock.Hash(), result.Hash())
			case <-time.After(10 * time.Second):
				t.Fatal("timeout")
			}
		}
		engine.Stop()
	}
}

func TestGetProposer(t *testing.T) {
	ctrl, mStaking := makeMockStakingManager(t, nil, 0)
	defer ctrl.Finish()

	chain, engine := newBlockChain(1, mStaking)
	defer engine.Stop()

	block := makeBlockWithSeal(chain, engine, chain.Genesis())
	_, err := chain.InsertChain(types.Blocks{block})
	assert.NoError(t, err)
	assert.Equal(t, engine.GetProposer(1), engine.Address())
}
