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

package forkid

import (
	"encoding/binary"
	"hash/crc32"
	"math/big"
	"reflect"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// ID is a fork identifier as defined by EIP-2124.
type ID struct {
	Hash [4]byte // CRC32 checksum of the genesis block and passed fork block numbers
	Next uint64  // Block number of the next upcoming fork, or 0 if no forks are known
}

// NewID calculates the Kaia fork ID from the chain config, genesis hash, head.
func NewID(config *params.ChainConfig, genesis common.Hash, head uint64) ID {
	// Calculate the starting checksum from the genesis hash
	hash := crc32.ChecksumIEEE(genesis[:])

	// Calculate the current fork checksum and the next fork block
	var next uint64
	for _, fork := range gatherForks(config) {
		if fork <= head {
			// Fork already passed, checksum the previous hash and the fork number
			hash = checksumUpdate(hash, fork)
			continue
		}
		next = fork
		break
	}
	return ID{Hash: checksumToBytes(hash), Next: next}
}

// LatestForkCompatibleBlock returns the latest fork compatible block or genesis(0) if no forks are known.
func LatestForkCompatibleBlock(config *params.ChainConfig, head *big.Int) *big.Int {
	latestForkCompatibleBlock := common.Big0
	for _, fork := range gatherForks(config) {
		if new(big.Int).SetUint64(fork).Cmp(head) <= 0 {
			latestForkCompatibleBlock = new(big.Int).SetUint64(fork)
			continue
		} else {
			break
		}
	}
	return latestForkCompatibleBlock
}

// NextForkCompatibleBlock returns the next fork compatible block or nil if no forks are known.
func NextForkCompatibleBlock(config *params.ChainConfig, head *big.Int) *big.Int {
	var nextForkCompatibleBlock *big.Int
	for _, fork := range gatherForks(config) {
		if new(big.Int).SetUint64(fork).Cmp(head) <= 0 {
			continue
		}
		nextForkCompatibleBlock = new(big.Int).SetUint64(fork)
		break
	}
	return nextForkCompatibleBlock
}

// LastForkCompatibleBlock returns the last fork compatible block or genesis(0) if no forks are known.
func LastForkCompatibleBlock(config *params.ChainConfig) *big.Int {
	forks := gatherForks(config)
	if len(forks) == 0 {
		return common.Big0
	}
	return new(big.Int).SetUint64(forks[len(forks)-1])
}

func BlobConfig(config *params.ChainConfig, head uint64) *params.BlobConfig {
	return nil
}

// checksumUpdate calculates the next IEEE CRC32 checksum based on the previous
// one and a fork block number (equivalent to CRC32(original-blob || fork)).
func checksumUpdate(hash uint32, fork uint64) uint32 {
	var blob [8]byte
	binary.BigEndian.PutUint64(blob[:], fork)
	return crc32.Update(hash, crc32.IEEETable, blob[:])
}

// checksumToBytes converts a uint32 checksum into a [4]byte array.
func checksumToBytes(hash uint32) [4]byte {
	var blob [4]byte
	binary.BigEndian.PutUint32(blob[:], hash)
	return blob
}

// gatherForks gathers all the known forks and creates a sorted list out of them.
func gatherForks(config *params.ChainConfig) []uint64 {
	// Gather all the fork block numbers via reflection
	kind := reflect.TypeFor[params.ChainConfig]()
	conf := reflect.ValueOf(config).Elem()

	var forks []uint64
	for i := 0; i < kind.NumField(); i++ {
		// Fetch the next field and skip non-fork rules
		field := kind.Field(i)
		if !strings.HasSuffix(field.Name, "CompatibleBlock") {
			continue
		}
		if field.Type != reflect.TypeFor[*big.Int]() {
			continue
		}
		// Extract the fork rule block number and aggregate it
		rule := conf.Field(i).Interface().(*big.Int)
		if rule != nil {
			forks = append(forks, rule.Uint64())
		}
	}
	// Sort the fork block numbers to permit chronologival XOR
	for i := 0; i < len(forks); i++ {
		for j := i + 1; j < len(forks); j++ {
			if forks[i] > forks[j] {
				forks[i], forks[j] = forks[j], forks[i]
			}
		}
	}
	// Deduplicate block numbers applying multiple forks
	for i := 1; i < len(forks); i++ {
		if forks[i] == forks[i-1] {
			forks = append(forks[:i], forks[i+1:]...)
			i--
		}
	}
	// Skip any forks in block 0, that's the genesis ruleset
	if len(forks) > 0 && forks[0] == 0 {
		forks = forks[1:]
	}
	return forks
}
