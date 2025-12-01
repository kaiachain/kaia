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
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/rlp"
)

type BlobStorageConfig struct {
	baseDir   string
	retention time.Duration
}

func DefaultBlobStorageConfig(c node.Config) BlobStorageConfig {
	return BlobStorageConfig{
		baseDir:   filepath.Join(c.ChainDataDir, "blob"),
		retention: 21 * 24 * time.Hour, // TODO Should use params.BLOB_SIDECARS_RETENTION
	}
}

var (
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

func (b *BlobStorage) Save(blockNumber big.Int, txIndex int, sidecar *types.BlobTxSidecar) error {
	if sidecar == nil {
		return errors.New("sidecar is nil")
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

func (b *BlobStorage) Get(blockNumber big.Int, txIndex int) (*types.BlobTxSidecar, error) {
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

func (b *BlobStorage) Prune(blockNumber big.Int) error {
	retentionBlockNumber := b.GetRetentionBlockNumber(blockNumber)
	if retentionBlockNumber.Sign() < 0 {
		// no target blocks to prune
		return nil
	}

	// Calculate retention subdirectory number once
	// Subtract 1 because we want to delete the previous subdirectory of the retention subdirectory
	retentionSubDir := new(big.Int).Sub(b.getSubDir(retentionBlockNumber), big.NewInt(1))

	// Get all subdirectories in the base directory
	entries, err := os.ReadDir(b.config.baseDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read base directory: %w", err)
	}

	var dirsToDelete []string

	// Process each subdirectory
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Parse subdirectory name as block number
		subDirNum, err := strconv.ParseUint(entry.Name(), 10, 64)
		if err != nil {
			// Skip non-numeric directory names
			continue
		}

		subDirNumBig := big.NewInt(int64(subDirNum))

		// Compare subDirNum with retentionSubDir
		if subDirNumBig.Cmp(retentionSubDir) >= 0 {
			continue
		}

		// Delete the entire subdirectory
		subDirPath := filepath.Join(b.config.baseDir, entry.Name())
		dirsToDelete = append(dirsToDelete, subDirPath)
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
func (b *BlobStorage) GetFilename(blockNumber big.Int, txIndex int) (string, string) {
	// Create subdirectory based on block number to avoid too many files in one directory
	subDir := b.getSubDir(blockNumber)
	dir := filepath.Join(b.config.baseDir, subDir.String())
	return dir, filepath.Join(dir, fmt.Sprintf("%d_%d.bin", blockNumber.Uint64(), txIndex))
}

func (b *BlobStorage) GetRetentionBlockNumber(blockNumber big.Int) big.Int {
	// Convert retention period to seconds (assuming 1 block per second)
	retentionSeconds := int64(b.config.retention.Seconds())
	retentionBlocks := big.NewInt(retentionSeconds)

	// Calculate the block number to retain by subtracting retention period from current block number
	retentionBlockNumber := new(big.Int).Sub(&blockNumber, retentionBlocks)

	return *retentionBlockNumber
}

// getSubDir calculates the subdirectory number for a given block number
// Subdirectories are created by dividing block number by 1000
func (b *BlobStorage) getSubDir(blockNumber big.Int) *big.Int {
	return new(big.Int).Div(&blockNumber, big.NewInt(1000))
}
