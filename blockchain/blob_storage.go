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
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/rlp"
)

const (
	BLOCKS_PER_BUCKET       = 1000
	BLOB_SIDECARS_RETENTION = 21 * 24 * time.Hour
)

type BlobStorageConfig struct {
	baseDir   string
	retention time.Duration
}

func DefaultBlobStorageConfig(baseDir string) BlobStorageConfig {
	return BlobStorageConfig{
		baseDir:   filepath.Join(baseDir, "blob"),
		retention: BLOB_SIDECARS_RETENTION,
	}
}

var (
	ErrBlobBlockNumberNil          = errors.New("block number is nil")
	ErrBlobSidecarNil              = errors.New("sidecar is nil")
	ErrBlobFailedToCreateDirectory = errors.New("failed to create blob directory")
	ErrBlobWriteFailed             = errors.New("failed to write blob file")
	ErrBlobNotFound                = errors.New("blob file not found")
	ErrBlobReadFailed              = errors.New("failed to read blob file")
)

type BlobStorage struct {
	config BlobStorageConfig
}

func NewBlobStorage(config BlobStorageConfig) *BlobStorage {
	return &BlobStorage{
		config: config,
	}
}

func (b *BlobStorage) Save(blockNumber *big.Int, txIndex int, sidecar *types.BlobTxSidecar) error {
	if blockNumber == nil {
		return ErrBlobBlockNumberNil
	}
	if sidecar == nil {
		return ErrBlobSidecarNil
	}

	// Get filename
	dir, filename := b.GetFilename(blockNumber, txIndex)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logger.Error("failed to create directory", "dir", dir, "err", err)
		return ErrBlobFailedToCreateDirectory
	}

	// RLP encode the sidecar
	encoded, err := rlp.EncodeToBytes(sidecar)
	if err != nil {
		return fmt.Errorf("failed to RLP encode sidecar: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, encoded, 0o644); err != nil {
		logger.Error("failed to write file", "file", filename, "err", err)
		return ErrBlobWriteFailed
	}

	return nil
}

func (b *BlobStorage) Get(blockNumber *big.Int, txIndex int) (*types.BlobTxSidecar, error) {
	if blockNumber == nil {
		return nil, ErrBlobBlockNumberNil
	}

	// Get filename
	_, filename := b.GetFilename(blockNumber, txIndex)

	// Read file
	encoded, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Warn("blob file not found", "file", filename)
			return nil, ErrBlobNotFound
		}
		logger.Error("failed to read file", "file", filename, "err", err)
		return nil, ErrBlobReadFailed
	}

	// RLP decode the sidecar
	var sidecar types.BlobTxSidecar
	if err := rlp.DecodeBytes(encoded, &sidecar); err != nil {
		return nil, fmt.Errorf("failed to RLP decode sidecar: %w", err)
	}

	return &sidecar, nil
}

// Prune removes all buckets that are older than `retentionBucketThreshold`.
// Example:
// at 1814400 +  500: do nothing
// at 1814400 + 1000: prune [0, 999]
// at 1814400 + 1500: do nothing
// at 1814400 + 2000: prune [0, 999], [1000, 1999]
// at 1814400 + 2500: do nothing
func (b *BlobStorage) Prune(current *big.Int) error {
	if current == nil {
		return ErrBlobBlockNumberNil
	}

	retentionBlockNumber := b.GetRetentionBlockNumber(current)
	if retentionBlockNumber == nil || retentionBlockNumber.Sign() < 0 {
		// no target blocks to prune
		return nil
	}

	if new(big.Int).Mod(retentionBlockNumber, big.NewInt(BLOCKS_PER_BUCKET)) != big.NewInt(0) {
		return nil
	}

	// Calculate retention bucket number
	retentionBucketThreshold := b.getBucketIdx(retentionBlockNumber)
	if retentionBucketThreshold == nil || retentionBucketThreshold.Sign() <= 0 {
		// no target blocks to prune
		return nil
	}

	// Get all bucket directories in the base directory
	entries, err := os.ReadDir(b.config.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	capacity := calculateCapacity(len(entries), b.config.retention)
	dirsToDelete := make([]string, 0, capacity)

	// Process each bucket directory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Parse bucket directory name as bucket number
		bucketNum, err := strconv.ParseUint(entry.Name(), 10, 64)
		if err != nil {
			// Skip non-numeric directory names
			continue
		}

		bucketNumBig := big.NewInt(int64(bucketNum))

		// Compare bucketNum with retentionBucketThreshold
		if bucketNumBig.Cmp(retentionBucketThreshold) < 0 {
			subDirPath := filepath.Join(b.config.baseDir, entry.Name())
			dirsToDelete = append(dirsToDelete, subDirPath)
		}
	}

	// Remove directories
	for _, dir := range dirsToDelete {
		if err := os.RemoveAll(dir); err != nil {
			// Ignore errors when removing directories
			continue
		}
	}

	return nil
}

// Return a filename like "11111/11111234_0.bin", in dirname and filename components.
func (b *BlobStorage) GetFilename(blockNumber *big.Int, txIndex int) (string, string) {
	if blockNumber == nil {
		// Return empty strings if blockNumber is nil to avoid panic
		// Caller should check for nil before calling this method
		return "", ""
	}

	// Create bucket directory based on block number to avoid too many files in one directory
	bucket := b.getBucketIdx(blockNumber)
	if bucket == nil {
		// Return empty strings if bucket is nil
		return "", ""
	}
	dir := filepath.Join(b.config.baseDir, bucket.String())
	return dir, filepath.Join(dir, fmt.Sprintf("%d_%d.bin", blockNumber.Uint64(), txIndex))
}

func (b *BlobStorage) GetRetentionBlockNumber(blockNumber *big.Int) *big.Int {
	if blockNumber == nil {
		// Return nil if blockNumber is nil
		return nil
	}

	// Convert retention period to seconds (assuming 1 block per second)
	retentionSeconds := int64(b.config.retention.Seconds())
	retentionBlocks := big.NewInt(retentionSeconds)

	// Calculate the block number to retain by subtracting retention period from current block number
	retentionBlockNumber := new(big.Int).Sub(blockNumber, retentionBlocks)

	return retentionBlockNumber
}

// getBucketIdx calculates the bucket number for a given block number
// Buckets are created by dividing block number by BLOCKS_PER_BUCKET
func (b *BlobStorage) getBucketIdx(blockNumber *big.Int) *big.Int {
	if blockNumber == nil {
		// Return nil if blockNumber is nil
		return nil
	}

	return new(big.Int).Div(blockNumber, big.NewInt(BLOCKS_PER_BUCKET))
}

// calculateCapacity calculates the appropriate capacity for the dirsToDelete slice
// based on the number of entries and the retention period.
//
// It estimates the maximum number of buckets that might be deleted by:
// - Calculating retention period in seconds (assuming 1 block per second)
// - Dividing by BLOCKS_PER_BUCKET to get the number of buckets (buckets are block number / BLOCKS_PER_BUCKET)
// - Adding a 2x buffer for safety to account for variations
// - Capping at maxCap (10000) to balance memory efficiency and filesystem performance
//
// Parameters:
//   - numEntries: The actual number of directory entries found
//   - retention: The retention period duration
//
// Returns:
//   - The calculated capacity, which is the minimum of numEntries and maxExpectedBuckets
func calculateCapacity(numEntries int, retention time.Duration) int {
	const maxCap = 10000
	// Calculate expected maximum buckets based on retention period
	// Buckets are created by dividing block number by BLOCKS_PER_BUCKET
	// Assuming 1 block per second: max buckets = retention_seconds / BLOCKS_PER_BUCKET
	// Add 2x buffer for safety
	retentionSeconds := int64(retention.Seconds())
	maxExpectedBuckets := min(maxCap, int((retentionSeconds/BLOCKS_PER_BUCKET)*2))
	capacity := min(maxExpectedBuckets, numEntries)
	return capacity
}
