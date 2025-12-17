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
	testcases := []struct {
		name        string
		blockNumber *big.Int
		txIndex     int
		sidecar     *types.BlobTxSidecar
		wantErr     bool
		errContains string
		verify      func(*testing.T, *BlobStorage, *big.Int, int)
	}{
		{
			name:        "success",
			blockNumber: big.NewInt(1000),
			txIndex:     0,
			sidecar:     createTestSidecar(t, 1),
			wantErr:     false,
			verify: func(t *testing.T, storage *BlobStorage, blockNumber *big.Int, txIndex int) {
				_, filename := storage.GetFilename(blockNumber, txIndex)
				_, err := os.Stat(filename)
				require.NoError(t, err)
				assert.True(t, !os.IsNotExist(err))
			},
		},
		{
			name:        "nil block number",
			blockNumber: nil,
			txIndex:     0,
			sidecar:     createTestSidecar(t, 1),
			wantErr:     true,
			errContains: "block number is nil",
		},
		{
			name:        "nil sidecar",
			blockNumber: big.NewInt(1000),
			txIndex:     0,
			sidecar:     nil,
			wantErr:     true,
			errContains: "sidecar is nil",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			config := BlobStorageConfig{
				BaseDir:   tmpDir,
				Retention: 21 * 24 * time.Hour,
			}
			storage := NewBlobStorage(config)

			err := storage.Save(tc.blockNumber, tc.txIndex, tc.sidecar)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, storage, tc.blockNumber, tc.txIndex)
				}
			}
		})
	}
}

func TestBlobStorage_Get(t *testing.T) {
	testcases := []struct {
		name        string
		blockNumber *big.Int
		txIndex     int
		setup       func(*BlobStorage, *big.Int, int) *types.BlobTxSidecar
		wantErr     bool
		errContains string
		verify      func(*testing.T, *types.BlobTxSidecar, *types.BlobTxSidecar)
	}{
		{
			name:        "success",
			blockNumber: big.NewInt(2000),
			txIndex:     1,
			setup: func(storage *BlobStorage, blockNumber *big.Int, txIndex int) *types.BlobTxSidecar {
				originalSidecar := createTestSidecar(t, 2)
				err := storage.Save(blockNumber, txIndex, originalSidecar)
				require.NoError(t, err)
				return originalSidecar
			},
			wantErr: false,
			verify: func(t *testing.T, original, retrieved *types.BlobTxSidecar) {
				require.NotNil(t, retrieved)
				assert.Equal(t, original.Version, retrieved.Version)
				assert.Equal(t, len(original.Blobs), len(retrieved.Blobs))
				assert.Equal(t, len(original.Commitments), len(retrieved.Commitments))
				assert.Equal(t, len(original.Proofs), len(retrieved.Proofs))

				for i := range original.Blobs {
					assert.Equal(t, original.Blobs[i], retrieved.Blobs[i])
					assert.Equal(t, original.Commitments[i], retrieved.Commitments[i])
					assert.Equal(t, original.Proofs[i], retrieved.Proofs[i])
				}
			},
		},
		{
			name:        "nil block number",
			blockNumber: nil,
			txIndex:     0,
			wantErr:     true,
			errContains: "block number is nil",
		},
		{
			name:        "file not found",
			blockNumber: big.NewInt(9999),
			txIndex:     99,
			wantErr:     true,
			errContains: "blob file not found",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			config := BlobStorageConfig{
				BaseDir:   tmpDir,
				Retention: 21 * 24 * time.Hour,
			}
			storage := NewBlobStorage(config)

			var originalSidecar *types.BlobTxSidecar
			if tc.setup != nil {
				originalSidecar = tc.setup(storage, tc.blockNumber, tc.txIndex)
			}

			retrievedSidecar, err := storage.Get(tc.blockNumber, tc.txIndex)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, originalSidecar, retrievedSidecar)
				}
			}
		})
	}
}

func TestBlobStorage_Prune(t *testing.T) {
	testcases := []struct {
		name        string
		retention   time.Duration
		blockNumber *big.Int
		setup       func(*BlobStorage) (*big.Int, *big.Int)
		wantErr     bool
		errContains string
		verify      func(*testing.T, *BlobStorage, *big.Int, *big.Int)
	}{
		{
			name:        "success with pruning",
			retention:   10 * time.Second, // Short retention for testing
			blockNumber: big.NewInt(2010),
			setup: func(storage *BlobStorage) (*big.Int, *big.Int) {
				// Save old block (should be pruned)
				// Block 100 is in bucket 0: 100/1000 = 0
				oldBlockNumber := big.NewInt(100)
				oldSidecar := createTestSidecar(t, 1)
				err := storage.Save(oldBlockNumber, 0, oldSidecar)
				require.NoError(t, err)

				// Save recent block (should be kept)
				// Block 2000 is in bucket 2: 2000/1000 = 2
				recentBlockNumber := big.NewInt(2000)
				recentSidecar := createTestSidecar(t, 1)
				err = storage.Save(recentBlockNumber, 0, recentSidecar)
				require.NoError(t, err)

				return oldBlockNumber, recentBlockNumber
			},
			wantErr: false,
			verify: func(t *testing.T, storage *BlobStorage, oldBlockNumber, recentBlockNumber *big.Int) {
				// Verify old block is deleted
				// retention is 10 seconds, so block 2010 - 10 = 2000
				// getBucketIdx(2000) = 2, retentionBucketThreshold = 2
				// bucket 0: 0 < 2, so it should be pruned
				// bucket 2: 2 >= 2, so it should be kept
				_, err := storage.Get(oldBlockNumber, 0)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "blob file not found")

				// Verify recent block still exists
				_, err = storage.Get(recentBlockNumber, 0)
				require.NoError(t, err)
			},
		},
		{
			name:        "success without pruning",
			retention:   10 * time.Second, // Short retention for testing
			blockNumber: big.NewInt(2011),
			setup: func(storage *BlobStorage) (*big.Int, *big.Int) {
				// Save old block (should be pruned)
				// Block 100 is in bucket 0: 100/1000 = 0
				oldBlockNumber := big.NewInt(100)
				oldSidecar := createTestSidecar(t, 1)
				err := storage.Save(oldBlockNumber, 0, oldSidecar)
				require.NoError(t, err)

				// Save recent block (should be kept)
				// Block 2000 is in bucket 2: 2000/1000 = 2
				recentBlockNumber := big.NewInt(2000)
				recentSidecar := createTestSidecar(t, 1)
				err = storage.Save(recentBlockNumber, 0, recentSidecar)
				require.NoError(t, err)

				return oldBlockNumber, recentBlockNumber
			},
			wantErr: false,
			verify: func(t *testing.T, storage *BlobStorage, oldBlockNumber, recentBlockNumber *big.Int) {
				// Verify old block is deleted
				// retention is 10 seconds, so block 2011 - 10 = 2001
				// 2001 % BLOCKS_PER_BUCKET != 0
				// so it should not be pruned

				_, err := storage.Get(oldBlockNumber, 0)
				require.NoError(t, err)

				// Verify recent block still exists
				_, err = storage.Get(recentBlockNumber, 0)
				require.NoError(t, err)
			},
		},
		{
			name:        "nil block number",
			retention:   21 * 24 * time.Hour,
			blockNumber: nil,
			wantErr:     true,
			errContains: "block number is nil",
		},
		{
			name:        "negative retention",
			retention:   21 * 24 * time.Hour,
			blockNumber: big.NewInt(100),
			wantErr:     false,
			// Should return without error and do nothing
		},
		{
			name:        "empty directory",
			retention:   10 * time.Second,
			blockNumber: big.NewInt(1000),
			wantErr:     false,
			// Prune on empty directory should not error
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			config := BlobStorageConfig{
				BaseDir:   tmpDir,
				Retention: tc.retention,
			}
			storage := NewBlobStorage(config)

			var oldBlockNumber, recentBlockNumber *big.Int
			if tc.setup != nil {
				oldBlockNumber, recentBlockNumber = tc.setup(storage)
			}

			err := storage.Prune(tc.blockNumber)
			if tc.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.errContains)
			} else {
				require.NoError(t, err)
				if tc.verify != nil {
					tc.verify(t, storage, oldBlockNumber, recentBlockNumber)
				}
			}
		})
	}
}

func TestBlobStorage_GetFilename(t *testing.T) {
	testcases := []struct {
		name           string
		blockNumber    *big.Int
		txIndex        int
		expectedBucket string
		expectedDir    string
		expectedFile   string
	}{
		{
			name:           "success",
			blockNumber:    big.NewInt(12345),
			txIndex:        5,
			expectedBucket: "12", // 12345 / 1000 = 12
			expectedDir:    "",   // Will be set in test
			expectedFile:   "12345_5.bin",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			config := BlobStorageConfig{
				BaseDir:   tmpDir,
				Retention: 21 * 24 * time.Hour,
			}
			storage := NewBlobStorage(config)

			dir, filename := storage.GetFilename(tc.blockNumber, tc.txIndex)

			// Verify directory structure
			expectedDir := filepath.Join(tmpDir, tc.expectedBucket)
			assert.Equal(t, expectedDir, dir)

			// Verify filename
			expectedFilename := filepath.Join(expectedDir, tc.expectedFile)
			assert.Equal(t, expectedFilename, filename)
		})
	}
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
