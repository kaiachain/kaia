// Copyright 2025 The Kaia Authors
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
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlobStorage_Save(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	// Create test sidecar
	blockNumber := big.NewInt(1000)
	txIndex := 0
	sidecar := createTestSidecar(t, 1)

	// Save the sidecar
	err := storage.Save(blockNumber, txIndex, sidecar)
	require.NoError(t, err)

	// Verify file exists
	_, filename := storage.GetFilename(blockNumber, txIndex)
	_, err = os.Stat(filename)
	require.NoError(t, err)
	assert.True(t, !os.IsNotExist(err))
}

func TestBlobStorage_Save_NilBlockNumber(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	txIndex := 0
	sidecar := createTestSidecar(t, 1)

	err := storage.Save(nil, txIndex, sidecar)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "block number is nil")
}

func TestBlobStorage_Save_NilSidecar(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	blockNumber := big.NewInt(1000)
	txIndex := 0

	err := storage.Save(blockNumber, txIndex, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "sidecar is nil")
}

func TestBlobStorage_Get(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	// Create and save test sidecar
	blockNumber := big.NewInt(2000)
	txIndex := 1
	originalSidecar := createTestSidecar(t, 2)

	err := storage.Save(blockNumber, txIndex, originalSidecar)
	require.NoError(t, err)

	// Get the sidecar
	retrievedSidecar, err := storage.Get(blockNumber, txIndex)
	require.NoError(t, err)
	require.NotNil(t, retrievedSidecar)

	// Verify sidecar data
	assert.Equal(t, originalSidecar.Version, retrievedSidecar.Version)
	assert.Equal(t, len(originalSidecar.Blobs), len(retrievedSidecar.Blobs))
	assert.Equal(t, len(originalSidecar.Commitments), len(retrievedSidecar.Commitments))
	assert.Equal(t, len(originalSidecar.Proofs), len(retrievedSidecar.Proofs))

	// Verify blob data matches
	for i := range originalSidecar.Blobs {
		assert.Equal(t, originalSidecar.Blobs[i], retrievedSidecar.Blobs[i])
		assert.Equal(t, originalSidecar.Commitments[i], retrievedSidecar.Commitments[i])
		assert.Equal(t, originalSidecar.Proofs[i], retrievedSidecar.Proofs[i])
	}
}

func TestBlobStorage_Get_NilBlockNumber(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	txIndex := 0

	_, err := storage.Get(nil, txIndex)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "block number is nil")
}

func TestBlobStorage_Get_FileNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	blockNumber := big.NewInt(9999)
	txIndex := 99

	_, err := storage.Get(blockNumber, txIndex)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blob file not found")
}

func TestBlobStorage_Prune(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 10 * time.Second, // Short retention for testing
	}
	storage := NewBlobStorage(config)

	// Save old block (should be pruned)
	// Block 100 is in epoch 0: 100/1000 = 0
	oldBlockNumber := big.NewInt(100)
	oldSidecar := createTestSidecar(t, 1)
	err := storage.Save(oldBlockNumber, 0, oldSidecar)
	require.NoError(t, err)

	// Save recent block (should be kept)
	// Block 2000 is in epoch 2: 2000/1000 = 2
	recentBlockNumber := big.NewInt(2000)
	recentSidecar := createTestSidecar(t, 1)
	err = storage.Save(recentBlockNumber, 0, recentSidecar)
	require.NoError(t, err)

	// Prune with current block number that makes old block exceed retention
	// retention is 10 seconds, so block 2010 - 10 = 2000
	// getEpochIdx(2000) = 2, retentionEpochThreshold = 2
	// epoch 0: 0 < 2, so it should be pruned
	// epoch 2: 2 >= 2, so it should be kept
	currentBlockNumber := big.NewInt(2010)
	err = storage.Prune(currentBlockNumber)
	require.NoError(t, err)

	// Verify old block is deleted
	_, err = storage.Get(oldBlockNumber, 0)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "blob file not found")

	// Verify recent block still exists
	_, err = storage.Get(recentBlockNumber, 0)
	require.NoError(t, err)
}

func TestBlobStorage_Prune_NilBlockNumber(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	err := storage.Prune(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "block number is nil")
}

func TestBlobStorage_Prune_NegativeRetention(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	// Use a very small block number that would result in negative retention
	blockNumber := big.NewInt(100)
	err := storage.Prune(blockNumber)
	require.NoError(t, err)
	// Should return without error and do nothing
}

func TestBlobStorage_Prune_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 10 * time.Second,
	}
	storage := NewBlobStorage(config)

	// Prune on empty directory should not error
	blockNumber := big.NewInt(1000)
	err := storage.Prune(blockNumber)
	require.NoError(t, err)
}

func TestBlobStorage_GetFilename(t *testing.T) {
	tmpDir := t.TempDir()
	config := BlobStorageConfig{
		baseDir:   tmpDir,
		retention: 21 * 24 * time.Hour,
	}
	storage := NewBlobStorage(config)

	blockNumber := big.NewInt(12345)
	txIndex := 5

	dir, filename := storage.GetFilename(blockNumber, txIndex)

	// Verify directory structure
	expectedEpoch := "12" // 12345 / 1000 = 12
	expectedDir := filepath.Join(tmpDir, expectedEpoch)
	assert.Equal(t, expectedDir, dir)

	// Verify filename
	expectedFilename := filepath.Join(expectedDir, "12345_5.bin")
	assert.Equal(t, expectedFilename, filename)
}

// createTestSidecar creates a test BlobTxSidecar with the specified number of blobs
func createTestSidecar(t *testing.T, numBlobs int) *types.BlobTxSidecar {
	blobs := make([]kzg4844.Blob, numBlobs)
	commitments := make([]kzg4844.Commitment, numBlobs)
	proofs := make([]kzg4844.Proof, numBlobs)

	// Fill with test data (simple pattern for testing)
	for i := 0; i < numBlobs; i++ {
		// Create a simple blob with pattern
		var blob kzg4844.Blob
		for j := 0; j < len(blob) && j < 100; j++ {
			blob[j] = byte(i + j)
		}
		blobs[i] = blob

		// Create simple commitment and proof
		var commitment kzg4844.Commitment
		var proof kzg4844.Proof
		for j := 0; j < len(commitment); j++ {
			commitment[j] = byte(i + j)
			proof[j] = byte(i + j + 10)
		}
		commitments[i] = commitment
		proofs[i] = proof
	}

	return types.NewBlobTxSidecar(types.BlobSidecarVersion1, blobs, commitments, proofs)
}
