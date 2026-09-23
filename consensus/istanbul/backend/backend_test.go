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
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

var (
	testSigningData = []byte("dummy data")
	maxBlockNum     = int64(100)
)

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
	configItems = append(configItems, istanbulCompatibleBlock(big.NewInt(5)))
	configItems = append(configItems, LondonCompatibleBlock(nil))
	configItems = append(configItems, EthTxTypeCompatibleBlock(nil))
	configItems = append(configItems, magmaCompatibleBlock(nil))
	configItems = append(configItems, koreCompatibleBlock(nil))
	configItems = append(configItems, shanghaiCompatibleBlock(nil))
	configItems = append(configItems, cancunCompatibleBlock(nil))
	configItems = append(configItems, randaoCompatibleBlock(nil))
	configItems = append(configItems, kaiaCompatibleBlock(nil))
	configItems = append(configItems, pragueCompatibleBlock(nil))
	configItems = append(configItems, osakaCompatibleBlock(nil))
	configItems = append(configItems, permissionlessCompatibleBlock(nil))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	configItems = append(configItems, mStaking)

	chain, istBackend := newBlockChain(t, len(stakes), configItems...)
	mGov := govModuleOf(istBackend)
	chain.RegisterKaiaxModules(
		mGov,
		istBackend.valsetModule,
		[]kaiax.ExecutionModule{mGov},
		nil,
		nil,
		nil,
	)
	defer istBackend.Stop()

	// Test for blocks from 0 to maxBlockNum
	// from 0 to 4: before istanbul hard fork
	// from 5 to 100: after istanbul hard fork
	var previousBlock, currentBlock *types.Block = nil, chain.CurrentBlock()
	for range maxBlockNum {
		previousBlock = currentBlock
		currentBlock = makeBlockWithSeal(chain, istBackend, previousBlock)
		_, err := chain.InsertChain(types.Blocks{currentBlock})
		assert.NoError(t, err)
	}

	for i := range maxBlockNum {
		// Test for round 0 to round 14
		for round := int64(0); round < 14; round++ {
			currentCommittee, err := istBackend.valsetModule.GetCommittee(uint64(i), uint64(round))
			assert.NoError(t, err)

			// skip if the testing node is not in a committee
			isInSubList := valset.NewAddressSet(currentCommittee).Contains(istBackend.Address())
			if !isInSubList {
				continue
			}

			nextCommittee, err := istBackend.valsetModule.GetCommittee(uint64(i), uint64(round)+1)
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
			committees[0] = valset.NewAddressSet(currentCommittee)
			committees[1] = valset.NewAddressSet(nextCommittee)
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
	config := params.TestChainConfig.Copy()
	return newTestBackendWithConfig(config, nil)
}

func newTestBackendWithConfig(chainConfig *params.ChainConfig, key *ecdsa.PrivateKey) (b *backend) {
	if key == nil {
		// if key is nil, generate new key for a test account
		key, _ = crypto.GenerateKey()
	}
	if chainConfig.Governance.GovernanceMode == "single" {
		// if governance mode is single, set the node key to the governing node.
		chainConfig.Governance.GoverningNode = crypto.PubkeyToAddress(key.PublicKey)
	}
	istanbulConfig := &istanbul.Config{
		Epoch:          chainConfig.Istanbul.Epoch,
		ProposerPolicy: istanbul.ProposerPolicy(chainConfig.Istanbul.ProposerPolicy),
		SubGroupSize:   chainConfig.Istanbul.SubGroupSize,
		Timeout:        10000,
	}

	backend := New(&BackendOpts{
		IstanbulConfig: istanbulConfig,
		PrivateKey:     key,
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
			[][]byte{append([]byte{1}, bytes.Repeat([]byte{0x00}, istanbul.IstanbulExtraSeal-1)...)},
		},
		{
			// invalid signature
			consensus.ErrInvalidCommittedSeals,
			nil,
		},
	} {
		chain, engine := newBlockChain(t, 1)

		block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
		expBlock, _ := engine.updateBlock(block)

		// Initialize commitCh via initSealState (simulating Seal() setup)
		sealCommitCh := engine.initSealState(expBlock.NumberU64(), expBlock.Hash())

		if test.expectedErr == nil {
			go func() {
				result := <-sealCommitCh
				if result != nil {
					commitCh <- result.Block
				}
			}()
		}

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

// buildBlockWithTx returns a sealed block carrying the given tx, built on top of the
// chain's current head.
func buildBlockWithTx(t *testing.T, chain *blockchain.BlockChain, engine *backend, txdata *types.TxInternalDataLegacy) *types.Block {
	t.Helper()
	parent := chain.CurrentBlock()
	header := makeHeader(parent, chain.Config())
	if err := chain.PrepareHeader(header); err != nil {
		t.Fatal(err)
	}
	state, err := chain.StateAt(parent.Root())
	if err != nil {
		t.Fatal(err)
	}
	emptyBlock, err := chain.Processor().FinalizeState(header, state, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	signedTx, err := types.SignTx(types.NewTx(txdata), types.LatestSignerForChainID(chain.Config().ChainID), key)
	if err != nil {
		t.Fatal(err)
	}
	// NewBlock recomputes the header TxHash from the supplied tx, so Verify's tx-root
	// check still passes.
	return sealBlock(engine, types.NewBlock(emptyBlock.Header(), []*types.Transaction{signedTx}, nil))
}

// buildOversizedBlock returns a sealed block whose RLP-encoded size exceeds
// params.MaxBlockSize. Verify does not execute transactions, so the injected tx only
// needs to inflate the encoded size.
func buildOversizedBlock(t *testing.T, chain *blockchain.BlockChain, engine *backend) *types.Block {
	t.Helper()
	recipient := common.Address{}
	block := buildBlockWithTx(t, chain, engine, &types.TxInternalDataLegacy{
		Price:     big.NewInt(1),
		GasLimit:  1,
		Recipient: &recipient,
		Amount:    big.NewInt(0),
		Payload:   make([]byte, params.MaxBlockSize), // pushes the whole block past the cap
	})
	if block.Size() <= params.MaxBlockSize {
		t.Fatalf("test setup: block is not oversized (size=%d, cap=%d)", uint64(block.Size()), params.MaxBlockSize)
	}
	return block
}

// TestBackend_VerifyEnforcesBlockSizeCap checks that backend.Verify applies the
// EIP-7934 RLP block-size cap (params.MaxBlockSize) on Osaka-enabled blocks, the
// same rule the import path enforces in BlockValidator.ValidateBody. Without this,
// honest validators would PREPARE/COMMIT an oversized proposal and only reject it
// later on import.
func TestBackend_VerifyEnforcesBlockSizeCap(t *testing.T) {
	// Osaka enabled: the cap applies.
	t.Run("osaka", func(t *testing.T) {
		chain, engine := newBlockChain(t, 1, osakaCompatibleBlock(big.NewInt(0)))
		defer chain.Stop()
		defer engine.Stop()

		// A normal-sized block must not be rejected by the size cap.
		normal := makeBlockWithSeal(chain, engine, chain.CurrentBlock())
		_, err := engine.Verify(normal)
		assert.NotErrorIs(t, err, blockchain.ErrBlockOversized)

		// An oversized proposal must be rejected up front in Verify.
		oversized := buildOversizedBlock(t, chain, engine)
		_, err = engine.Verify(oversized)
		assert.ErrorIs(t, err, blockchain.ErrBlockOversized)
	})

	// Pre-Osaka: the cap is gated on IsOsakaForkEnabled, so it must not apply.
	t.Run("pre-osaka", func(t *testing.T) {
		chain, engine := newBlockChain(t, 1, osakaCompatibleBlock(nil), permissionlessCompatibleBlock(nil))
		defer chain.Stop()
		defer engine.Stop()

		oversized := buildOversizedBlock(t, chain, engine)
		_, err := engine.Verify(oversized)
		assert.NotErrorIs(t, err, blockchain.ErrBlockOversized)
	})
}

// TestBackend_VerifyEnforcesBodyRules checks that backend.Verify applies the body rules
// ValidateBody enforces on import, running both paths on the same block.
func TestBackend_VerifyEnforcesBodyRules(t *testing.T) {
	t.Run("gasPrice", func(t *testing.T) {
		const baseFeeGkei = 30
		chain, engine := newBlockChain(t, 1,
			LondonCompatibleBlock(big.NewInt(0)),
			EthTxTypeCompatibleBlock(big.NewInt(0)),
			magmaCompatibleBlock(big.NewInt(0)),
			lowerBoundBaseFee(baseFeeGkei*params.Gkei),
			upperBoundBaseFee(baseFeeGkei*params.Gkei),
		)
		defer chain.Stop()
		defer engine.Stop()

		recipient := common.Address{}
		build := func(gasPrice uint64) *types.Block {
			return buildBlockWithTx(t, chain, engine, &types.TxInternalDataLegacy{
				Price:     new(big.Int).SetUint64(gasPrice),
				GasLimit:  21000,
				Recipient: &recipient,
				Amount:    big.NewInt(0),
			})
		}

		priced := build(baseFeeGkei * params.Gkei)
		_, err := engine.Verify(priced)
		assert.NoError(t, err)
		assert.NoError(t, chain.Validator().ValidateBody(priced))

		underpriced := build(baseFeeGkei*params.Gkei - 1)
		_, err = engine.Verify(underpriced)
		assert.ErrorContains(t, err, "invalid GasPrice")
		assert.ErrorContains(t, chain.Validator().ValidateBody(underpriced), "invalid GasPrice")
	})

	// ValidateHeader cannot catch this: it never sees the body.
	t.Run("blobGasUsed", func(t *testing.T) {
		chain, engine := newBlockChain(t, 1, osakaCompatibleBlock(big.NewInt(0)))
		defer chain.Stop()
		defer engine.Stop()

		normal := makeBlockWithSeal(chain, engine, chain.CurrentBlock())
		_, err := engine.Verify(normal)
		assert.NoError(t, err)

		header := normal.Header()
		claimed := uint64(params.BlobTxBlobGasPerBlob)
		header.BlobGasUsed = &claimed
		lying := sealBlock(engine, types.NewBlock(header, nil, nil))

		_, err = engine.Verify(lying)
		assert.ErrorContains(t, err, "blob gas used mismatch")
		assert.ErrorContains(t, chain.Validator().ValidateBody(lying), "blob gas used mismatch")
	})
}
