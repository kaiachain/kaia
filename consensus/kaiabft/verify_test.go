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

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
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
