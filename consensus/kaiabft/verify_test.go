// Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package kaiabft

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/crypto"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// dev #967, #1028: verify rejects a proposal the import path would refuse.
func TestDevParity_VerifyAppliesImportBodyRules(t *testing.T) {
	signer := types.LatestSignerForChainID(big.NewInt(1))
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	tx, err := types.SignTx(types.NewTransaction(0, common.Address{1}, big.NewInt(0), 21000, big.NewInt(1), nil), signer, key)
	require.NoError(t, err)

	b, _ := newTestBackend(t, &fakeChain{cfg: legacyConfig()})

	underpriced := types.NewBlock(&types.Header{Number: big.NewInt(1), Time: big.NewInt(0), BaseFee: big.NewInt(10)}, []*types.Transaction{tx}, nil)
	_, err = b.verify(underpriced)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GasPrice")

	blobGas := uint64(params.BlobTxBlobGasPerBlob)
	phantomBlobs := types.NewBlock(&types.Header{Number: big.NewInt(1), Time: big.NewInt(0), BlobGasUsed: &blobGas}, nil, nil)
	_, err = b.verify(phantomBlobs)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blob gas used mismatch")

	osaka := &fakeChain{cfg: &params.ChainConfig{OsakaCompatibleBlock: big.NewInt(0)}}
	b2, _ := newTestBackend(t, osaka)
	bulky, err := types.SignTx(types.NewTransaction(0, common.Address{1}, big.NewInt(0), 21000, big.NewInt(1), make([]byte, params.MaxBlockSize)), signer, key)
	require.NoError(t, err)
	oversized := types.NewBlock(&types.Header{Number: big.NewInt(1), Time: big.NewInt(0)}, []*types.Transaction{bulky}, nil)
	_, err = b2.verify(oversized)
	assert.ErrorIs(t, err, blockchain.ErrBlockOversized)
}

// The import body rules run before speculative execution: a body-invalid
// proposal from the proposer is rejected without the executor being touched.
func TestBodyInvalidProposalSkipsSpeculativeExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	signer := types.LatestSignerForChainID(big.NewInt(1))
	txKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	underpriced, err := types.SignTx(types.NewTransaction(0, common.Address{1}, big.NewInt(0), 21000, big.NewInt(1), nil), signer, txKey)
	require.NoError(t, err)

	// The chain head is the unsealed genesis so the ROUND CHANGE sent on
	// rejection can read lastProposal without a sealer Author lookup.
	genesis := newEmptyBlock(0, common.Hash{})
	chain := &fakeChain{cfg: legacyConfig(), current: genesis, headers: map[common.Hash]*types.Header{genesis.Hash(): genesis.Header()}}
	b, key := newTestBackend(t, chain)
	self := crypto.PubkeyToAddress(key.PublicKey)
	proposerKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	proposer := crypto.PubkeyToAddress(proposerKey.PublicKey)

	// Any executor call (Clone is the first one startSpeculativeExecution makes)
	// fails the test: the mock has no expectations. The rejection sends a ROUND
	// CHANGE, whose catch-up asks the valset module for the next proposer.
	b.specCache = blockchain.NewSpeculativeResultCache()
	b.executor = mocks.NewMockExecutor(ctrl)
	mValset := valset_mock.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetProposer(gomock.Any(), gomock.Any()).Return(proposer, nil).AnyTimes()
	b.valsetModule = mValset

	block := types.NewBlock(&types.Header{Number: big.NewInt(10), ParentHash: genesis.Hash(), Time: big.NewInt(0), BaseFee: big.NewInt(10)}, []*types.Transaction{underpriced}, nil)
	m := newTestMachine(b, 10, 0, nil, []common.Address{self, proposer})
	defer m.stopTimer()
	m.state = stateAcceptRequest
	m.proposer = proposer
	pp := &bft.Preprepare{View: &bft.View{Sequence: big.NewInt(10), Round: big.NewInt(0)}, Proposal: block}
	err = m.handlePreprepare(&bft.Message{Code: bft.MsgPreprepare, Msg: mustEncode(t, pp)}, proposer)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid GasPrice")
	assert.False(t, b.specCache.HasUsable(block.Hash()), "no speculative entry may be reserved for a rejected body")
	assert.Equal(t, int64(1), m.round.Int64(), "the rejection must move to the next round")
}
