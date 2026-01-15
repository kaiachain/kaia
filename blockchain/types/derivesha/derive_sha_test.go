// Modifications Copyright 2024 The Kaia Authors
// Copyright 2022 The klaytn Authors
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

package derivesha

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

var dummyList = types.Transactions([]*types.Transaction{
	types.NewTransaction(1, common.Address{}, big.NewInt(123), 21000, big.NewInt(25e9), nil),
})

type testGov struct{}

// Mimic the governance vote situation
var testGovSchedule = map[uint64]int{
	0: types.ImplDeriveShaOriginal,
	1: types.ImplDeriveShaOriginal,
	2: types.ImplDeriveShaOriginal,
	3: types.ImplDeriveShaSimple,
	4: types.ImplDeriveShaSimple,
	5: types.ImplDeriveShaSimple,
	6: types.ImplDeriveShaConcat,
	7: types.ImplDeriveShaConcat,
	8: types.ImplDeriveShaConcat,
}

func (e *testGov) GetParamSet(num uint64) gov.ParamSet {
	return gov.ParamSet{
		DeriveShaImpl: uint64(testGovSchedule[num]),
	}
}

func TestEmptyRoot(t *testing.T) {
	assert.Equal(t,
		DeriveShaOrig{}.DeriveTransactionsRoot(types.Transactions{}).Hex(),
		types.EmptyTxRootOriginal.Hex())
	assert.Equal(t,
		DeriveShaSimple{}.DeriveTransactionsRoot(types.Transactions{}).Hex(),
		types.EmptyTxRootSimple.Hex())
	assert.Equal(t,
		DeriveShaConcat{}.DeriveTransactionsRoot(types.Transactions{}).Hex(),
		types.EmptyTxRootConcat.Hex())
}

func TestMuxChainConfig(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

	for implType, impl := range impls {
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: implType}, nil)
		assert.Equal(t,
			DeriveTransactionsRootMux(dummyList, big.NewInt(0)),
			impl.DeriveTransactionsRoot(dummyList),
		)
		assert.Equal(t,
			EmptyRootHashMux(big.NewInt(0)),
			impl.DeriveTransactionsRoot(types.Transactions{}),
		)
	}
}

func TestMuxGovernance(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

	InitDeriveSha(
		&params.ChainConfig{DeriveShaImpl: testGovSchedule[0]},
		&testGov{})

	for num := uint64(0); num < 9; num++ {
		implType := testGovSchedule[num]
		impl := impls[implType]

		assert.Equal(t,
			DeriveTransactionsRootMux(dummyList, new(big.Int).SetUint64(num)),
			impl.DeriveTransactionsRoot(dummyList),
		)
		assert.Equal(t,
			EmptyRootHashMux(new(big.Int).SetUint64(num)),
			impl.DeriveTransactionsRoot(types.Transactions{}),
		)
	}
}

// createBlobTransactionWithSidecar creates a blob transaction with sidecar for testing.
func createBlobTransactionWithSidecar() *types.Transaction {
	emptyBlob := new(kzg4844.Blob)
	emptyBlobCommit, _ := kzg4844.BlobToCommitment(emptyBlob)
	emptyBlobProof, _ := kzg4844.ComputeBlobProof(emptyBlob, emptyBlobCommit)
	sidecar := types.NewBlobTxSidecar(types.BlobSidecarVersion0, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, []kzg4844.Proof{emptyBlobProof})

	to := common.HexToAddress("7b65B75d204aBed71587c9E519a89277766EE1d0")
	blobData, err := types.NewTxInternalDataWithMap(types.TxTypeEthereumBlob, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      uint64(1234),
		types.TxValueKeyTo:         to,
		types.TxValueKeyAmount:     big.NewInt(10),
		types.TxValueKeyGasLimit:   uint64(1000000),
		types.TxValueKeyGasFeeCap:  big.NewInt(25),
		types.TxValueKeyGasTipCap:  big.NewInt(25),
		types.TxValueKeyData:       []byte("1234"),
		types.TxValueKeyAccessList: types.AccessList{{Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), StorageKeys: []common.Hash{{0}}}},
		types.TxValueKeyBlobFeeCap: big.NewInt(25),
		types.TxValueKeyBlobHashes: []common.Hash{{0}},
		types.TxValueKeySidecar:    sidecar,
		types.TxValueKeyChainID:    big.NewInt(2),
	})
	if err != nil {
		panic(err)
	}

	return types.NewTx(blobData)
}

// TestDeriveTransactionsRootExcludesSidecar tests that DeriveTransactionsRoot
// excludes blob sidecars when calculating the hash after InitDeriveSha().
func TestDeriveTransactionsRootExcludesSidecar(t *testing.T) {
	t.Run("InitDeriveSha execution excludes sidecar", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha with Original implementation
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		// Create blob transaction with sidecar
		txWithSidecar := createBlobTransactionWithSidecar()
		assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "transaction should have sidecar")

		// Create the same transaction without sidecar manually
		txWithoutSidecar := txWithSidecar.WithoutBlobTxSidecar()
		assert.Nil(t, txWithoutSidecar.BlobTxSidecar(), "transaction should not have sidecar")

		// Create transaction lists
		listWithSidecar := types.Transactions{txWithSidecar}
		listWithoutSidecar := types.Transactions{txWithoutSidecar}

		// Calculate roots using types.DeriveTransactionsRoot (which is set by InitDeriveSha)
		rootWithSidecar := types.DeriveTransactionsRoot(listWithSidecar, big.NewInt(0))
		rootWithoutSidecar := types.DeriveTransactionsRoot(listWithoutSidecar, big.NewInt(0))

		// Both should produce the same hash because DeriveTransactionsRoot excludes sidecars
		assert.Equal(t, rootWithSidecar, rootWithoutSidecar,
			"DeriveTransactionsRoot should exclude sidecars, so both lists should produce the same hash")
	})

	t.Run("Multiple implementation types", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		implementations := []struct {
			name     string
			implType int
		}{
			{"Original", types.ImplDeriveShaOriginal},
			{"Simple", types.ImplDeriveShaSimple},
			{"Concat", types.ImplDeriveShaConcat},
		}

		for _, impl := range implementations {
			t.Run(impl.name, func(t *testing.T) {
				// Initialize DeriveSha with specific implementation
				InitDeriveSha(&params.ChainConfig{DeriveShaImpl: impl.implType}, nil)

				// Create blob transaction with sidecar
				txWithSidecar := createBlobTransactionWithSidecar()
				assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "transaction should have sidecar")

				// Create the same transaction without sidecar manually
				txWithoutSidecar := txWithSidecar.WithoutBlobTxSidecar()
				assert.Nil(t, txWithoutSidecar.BlobTxSidecar(), "transaction should not have sidecar")

				// Create transaction lists
				listWithSidecar := types.Transactions{txWithSidecar}
				listWithoutSidecar := types.Transactions{txWithoutSidecar}

				// Calculate roots using types.DeriveTransactionsRoot
				rootWithSidecar := types.DeriveTransactionsRoot(listWithSidecar, big.NewInt(0))
				rootWithoutSidecar := types.DeriveTransactionsRoot(listWithoutSidecar, big.NewInt(0))

				// Both should produce the same hash because DeriveTransactionsRoot excludes sidecars
				assert.Equal(t, rootWithSidecar, rootWithoutSidecar,
					"%s implementation should exclude sidecars", impl.name)
			})
		}
	})

	t.Run("Mixed transaction list", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		// Create blob transaction with sidecar
		txWithSidecar := createBlobTransactionWithSidecar()
		assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "blob transaction should have sidecar")

		// Create legacy transaction
		legacyTx := types.NewTransaction(1, common.Address{}, big.NewInt(123), 21000, big.NewInt(25e9), nil)

		// Create mixed lists
		mixedListWithSidecar := types.Transactions{txWithSidecar, legacyTx}
		mixedListWithoutSidecar := types.Transactions{txWithSidecar.WithoutBlobTxSidecar(), legacyTx}

		// Calculate roots using types.DeriveTransactionsRoot
		rootWithSidecar := types.DeriveTransactionsRoot(mixedListWithSidecar, big.NewInt(0))
		rootWithoutSidecar := types.DeriveTransactionsRoot(mixedListWithoutSidecar, big.NewInt(0))

		// Both should produce the same hash
		assert.Equal(t, rootWithSidecar, rootWithoutSidecar,
			"DeriveTransactionsRoot should exclude sidecars even in mixed transaction lists")
	})

	t.Run("Edge cases", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		t.Run("Empty transaction list", func(t *testing.T) {
			emptyList := types.Transactions{}
			root := types.DeriveTransactionsRoot(emptyList, big.NewInt(0))
			expectedRoot := types.GetEmptyRootHash(big.NewInt(0))
			assert.Equal(t, expectedRoot, root, "empty list should produce empty root hash")
		})

		t.Run("Blob transaction without sidecar", func(t *testing.T) {
			// Create blob transaction with sidecar first, then remove it
			txWithSidecar := createBlobTransactionWithSidecar()
			txWithoutSidecar := txWithSidecar.WithoutBlobTxSidecar()
			assert.Nil(t, txWithoutSidecar.BlobTxSidecar(), "transaction should not have sidecar")

			// Create list with blob transaction without sidecar
			list := types.Transactions{txWithoutSidecar}
			root := types.DeriveTransactionsRoot(list, big.NewInt(0))

			// Should not panic and produce a valid hash
			assert.NotEqual(t, common.Hash{}, root, "should produce a valid hash")
		})

		t.Run("Multiple blob transactions with sidecars", func(t *testing.T) {
			// Create multiple blob transactions with sidecars
			tx1 := createBlobTransactionWithSidecar()
			tx2 := createBlobTransactionWithSidecar()

			assert.NotNil(t, tx1.BlobTxSidecar(), "tx1 should have sidecar")
			assert.NotNil(t, tx2.BlobTxSidecar(), "tx2 should have sidecar")

			// Create lists
			listWithSidecars := types.Transactions{tx1, tx2}
			listWithoutSidecars := types.Transactions{tx1.WithoutBlobTxSidecar(), tx2.WithoutBlobTxSidecar()}

			// Calculate roots using types.DeriveTransactionsRoot
			rootWithSidecars := types.DeriveTransactionsRoot(listWithSidecars, big.NewInt(0))
			rootWithoutSidecars := types.DeriveTransactionsRoot(listWithoutSidecars, big.NewInt(0))

			// Both should produce the same hash
			assert.Equal(t, rootWithSidecars, rootWithoutSidecars,
				"multiple blob transactions with sidecars should produce the same hash as without sidecars")
		})
	})
}

// createTestHeader creates a simple header for testing.
func createTestHeader(number *big.Int) *types.Header {
	return &types.Header{
		ParentHash:  common.Hash{},
		Rewardbase:  common.Address{},
		Root:        common.Hash{},
		TxHash:      common.Hash{},
		ReceiptHash: common.Hash{},
		BlockScore:  big.NewInt(1),
		Number:      number,
		GasUsed:     0,
		Time:        big.NewInt(0),
		TimeFoS:     0,
		Extra:       []byte{},
		Governance:  []byte{},
		Vote:        []byte{},
	}
}

// TestNewBlockExcludesBlobSidecar tests that NewBlock excludes blob sidecars
// when calculating Header.TxHash.
func TestNewBlockExcludesBlobSidecar(t *testing.T) {
	t.Run("NewBlock excludes sidecar", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha with Original implementation
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		// Create blob transaction with sidecar
		txWithSidecar := createBlobTransactionWithSidecar()
		assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "transaction should have sidecar")

		// Create header
		header := createTestHeader(big.NewInt(0))

		// Create block with sidecar
		blockWithSidecar := types.NewBlock(header, []*types.Transaction{txWithSidecar}, []*types.Receipt{})
		assert.NotNil(t, blockWithSidecar, "block should be created")
		txHashWithSidecar := blockWithSidecar.Header().TxHash

		// Verify original transaction still has sidecar
		assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "original transaction should still have sidecar")

		// Create block without sidecar
		txWithoutSidecar := txWithSidecar.WithoutBlobTxSidecar()
		assert.Nil(t, txWithoutSidecar.BlobTxSidecar(), "transaction should not have sidecar")

		blockWithoutSidecar := types.NewBlock(header, []*types.Transaction{txWithoutSidecar}, []*types.Receipt{})
		assert.NotNil(t, blockWithoutSidecar, "block should be created")
		txHashWithoutSidecar := blockWithoutSidecar.Header().TxHash

		// Both blocks should have the same TxHash because sidecars are excluded
		assert.Equal(t, txHashWithSidecar, txHashWithoutSidecar,
			"NewBlock should exclude sidecars, so both blocks should have the same TxHash")
	})

	t.Run("Multiple implementation types", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		implementations := []struct {
			name     string
			implType int
		}{
			{"Original", types.ImplDeriveShaOriginal},
			{"Simple", types.ImplDeriveShaSimple},
			{"Concat", types.ImplDeriveShaConcat},
		}

		for _, impl := range implementations {
			t.Run(impl.name, func(t *testing.T) {
				// Initialize DeriveSha with specific implementation
				InitDeriveSha(&params.ChainConfig{DeriveShaImpl: impl.implType}, nil)

				// Create blob transaction with sidecar
				txWithSidecar := createBlobTransactionWithSidecar()
				assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "transaction should have sidecar")

				// Create header
				header := createTestHeader(big.NewInt(0))

				// Create blocks
				blockWithSidecar := types.NewBlock(header, []*types.Transaction{txWithSidecar}, []*types.Receipt{})
				txHashWithSidecar := blockWithSidecar.Header().TxHash

				blockWithoutSidecar := types.NewBlock(header, []*types.Transaction{txWithSidecar.WithoutBlobTxSidecar()}, []*types.Receipt{})
				txHashWithoutSidecar := blockWithoutSidecar.Header().TxHash

				// Both blocks should have the same TxHash
				assert.Equal(t, txHashWithSidecar, txHashWithoutSidecar,
					"%s implementation should exclude sidecars", impl.name)
			})
		}
	})

	t.Run("Mixed transaction list", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		// Create blob transaction with sidecar
		txWithSidecar := createBlobTransactionWithSidecar()
		assert.NotNil(t, txWithSidecar.BlobTxSidecar(), "blob transaction should have sidecar")

		// Create legacy transaction
		legacyTx := types.NewTransaction(1, common.Address{}, big.NewInt(123), 21000, big.NewInt(25e9), nil)

		// Create header
		header := createTestHeader(big.NewInt(0))

		// Create blocks with mixed transactions
		mixedBlockWithSidecar := types.NewBlock(header, []*types.Transaction{txWithSidecar, legacyTx}, []*types.Receipt{})
		txHashWithSidecar := mixedBlockWithSidecar.Header().TxHash

		mixedBlockWithoutSidecar := types.NewBlock(header, []*types.Transaction{txWithSidecar.WithoutBlobTxSidecar(), legacyTx}, []*types.Receipt{})
		txHashWithoutSidecar := mixedBlockWithoutSidecar.Header().TxHash

		// Both blocks should have the same TxHash
		assert.Equal(t, txHashWithSidecar, txHashWithoutSidecar,
			"NewBlock should exclude sidecars even in mixed transaction lists")
	})

	t.Run("Edge cases", func(t *testing.T) {
		log.EnableLogForTest(log.LvlCrit, log.LvlInfo)

		// Initialize DeriveSha
		InitDeriveSha(&params.ChainConfig{DeriveShaImpl: types.ImplDeriveShaOriginal}, nil)

		t.Run("Empty transaction list", func(t *testing.T) {
			header := createTestHeader(big.NewInt(0))

			emptyBlock := types.NewBlock(header, []*types.Transaction{}, []*types.Receipt{})
			expectedTxHash := types.GetEmptyRootHash(header.Number)
			assert.Equal(t, expectedTxHash, emptyBlock.Header().TxHash,
				"empty transaction list should produce empty root hash")
		})

		t.Run("Multiple blob transactions with sidecars", func(t *testing.T) {
			// Create multiple blob transactions with sidecars
			tx1 := createBlobTransactionWithSidecar()
			tx2 := createBlobTransactionWithSidecar()

			assert.NotNil(t, tx1.BlobTxSidecar(), "tx1 should have sidecar")
			assert.NotNil(t, tx2.BlobTxSidecar(), "tx2 should have sidecar")

			// Create header
			header := createTestHeader(big.NewInt(0))

			// Create blocks
			blockWithSidecars := types.NewBlock(header, []*types.Transaction{tx1, tx2}, []*types.Receipt{})
			txHashWithSidecars := blockWithSidecars.Header().TxHash

			blockWithoutSidecars := types.NewBlock(header, []*types.Transaction{tx1.WithoutBlobTxSidecar(), tx2.WithoutBlobTxSidecar()}, []*types.Receipt{})
			txHashWithoutSidecars := blockWithoutSidecars.Header().TxHash

			// Both blocks should have the same TxHash
			assert.Equal(t, txHashWithSidecars, txHashWithoutSidecars,
				"multiple blob transactions with sidecars should produce the same hash as without sidecars")
		})
	})
}
