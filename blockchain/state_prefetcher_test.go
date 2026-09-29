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

package blockchain

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

type prefetchTestKeyPicker map[common.Address]accountkey.AccountKey

func (p prefetchTestKeyPicker) GetKey(addr common.Address) accountkey.AccountKey {
	return p[addr]
}

func (p prefetchTestKeyPicker) Exist(addr common.Address) bool {
	return p[addr] != nil
}

// TestCopyTxForPrefetchDoesNotOverwriteOriginalValidatedGas verifies that
// prefetch cannot overwrite state-dependent validatedGas cached on the original
// block transaction.
func TestCopyTxForPrefetchDoesNotOverwriteOriginalValidatedGas(t *testing.T) {
	require.NoError(t, fork.SetHardForkBlockNumberConfig(&params.ChainConfig{
		IstanbulCompatibleBlock: big.NewInt(100),
	}))

	const blockNumber = uint64(0)
	signer := types.LatestSignerForChainID(big.NewInt(1))
	senderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	extraKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	from := crypto.PubkeyToAddress(senderKey.PublicKey)
	to := common.HexToAddress("0x0000000000000000000000000000000000000100")
	txData, err := types.NewTxInternalDataWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    uint64(0),
		types.TxValueKeyTo:       to,
		types.TxValueKeyAmount:   big.NewInt(1),
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: big.NewInt(1),
		types.TxValueKeyFrom:     from,
	})
	require.NoError(t, err)

	tx := types.NewTx(txData)
	require.NoError(t, tx.SignWithKeys(signer, []*ecdsa.PrivateKey{senderKey}))

	stalePicker := prefetchTestKeyPicker{
		from: accountkey.NewAccountKeyPublicWithValue(&senderKey.PublicKey),
	}
	currentPicker := prefetchTestKeyPicker{
		from: accountkey.NewAccountKeyWeightedMultiSigWithValues(1, accountkey.WeightedPublicKeys{
			accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&senderKey.PublicKey)),
			accountkey.NewWeightedPublicKey(1, (*accountkey.PublicKeySerializable)(&extraKey.PublicKey)),
		}),
	}

	_, err = tx.AsMessageWithAccountKeyPicker(signer, currentPicker, blockNumber)
	require.NoError(t, err)
	originalGas := tx.ValidatedGas().IntrinsicGas

	prefetchTx := copyTxForPrefetch(tx)
	_, err = prefetchTx.AsMessageWithAccountKeyPicker(signer, stalePicker, blockNumber)
	require.NoError(t, err)

	require.Equal(t, params.TxValidationGasPerKey, originalGas-prefetchTx.ValidatedGas().IntrinsicGas)
	require.Equal(t, originalGas, tx.ValidatedGas().IntrinsicGas)
}

type prefetchTestChain struct {
	cfg *params.ChainConfig
}

func (c prefetchTestChain) CurrentBlock() *types.Block                  { return nil }
func (c prefetchTestChain) GetHeader(common.Hash, uint64) *types.Header { return nil }
func (c prefetchTestChain) Config() *params.ChainConfig                 { return c.cfg }
func (c prefetchTestChain) Sealer() consensus.Sealer                    { return nil }

// TestPrefetchDoesNotCacheSenderOnOriginalTx verifies that the state-warming
// Prefetch validates a throwaway copy, so the sender cached during prefetch
// never lands on the transaction the block will execute.
func TestPrefetchDoesNotCacheSenderOnOriginalTx(t *testing.T) {
	cfg := params.TestChainConfig.Copy()
	require.NoError(t, fork.SetHardForkBlockNumberConfig(cfg))

	senderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	from := crypto.PubkeyToAddress(senderKey.PublicKey)
	txData, err := types.NewTxInternalDataWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    uint64(0),
		types.TxValueKeyTo:       common.HexToAddress("0x0000000000000000000000000000000000000100"),
		types.TxValueKeyAmount:   big.NewInt(1),
		types.TxValueKeyGasLimit: uint64(100000),
		types.TxValueKeyGasPrice: big.NewInt(1),
		types.TxValueKeyFrom:     from,
	})
	require.NoError(t, err)
	tx := types.NewTx(txData)
	require.NoError(t, tx.SignWithKeys(types.LatestSignerForChainID(cfg.ChainID), []*ecdsa.PrivateKey{senderKey}))
	require.Equal(t, common.Address{}, tx.ValidatedSender())

	header := &types.Header{Number: big.NewInt(1)}
	block := types.NewBlock(header, []*types.Transaction{tx}, nil)
	stateDB, err := state.New(types.EmptyRootHash, state.NewDatabase(database.NewMemoryDBManager()), nil, nil)
	require.NoError(t, err)

	// The copy must validate against the same state; otherwise the wrap
	// would be hiding a prefetch failure instead of isolating the cache.
	_, err = copyTxForPrefetch(tx).ValidateSender(types.MakeSigner(cfg, header.Number), stateDB, header.Number.Uint64())
	require.NoError(t, err)

	newStatePrefetcher(prefetchTestChain{cfg: cfg}).Prefetch(block, stateDB, vm.Config{}, nil)

	require.Equal(t, common.Address{}, tx.ValidatedSender())
}
