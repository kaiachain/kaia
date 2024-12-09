// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2018 The go-ethereum Authors
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
// This file is derived from core/rawdb/schema.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package database

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/rcrowley/go-metrics"
)

// The fields below define the low level database schema prefixing.
var (
	// databaseVerisionKey tracks the current database version.
	databaseVerisionKey = []byte("DatabaseVersion")

	// headHeaderKey tracks the latest know header's hash.
	headHeaderKey = []byte("LastHeader")

	// headBlockKey tracks the latest know full block's hash.
	headBlockKey       = []byte("LastBlock")
	headBlockBackupKey = []byte("LastBlockBackup")

	// headFastBlockKey tracks the latest known incomplete block's hash duirng fast sync.
	headFastBlockKey       = []byte("LastFast")
	headFastBlockBackupKey = []byte("LastFastBackup")

	// fastTrieProgressKey tracks the number of trie entries imported during fast sync.
	fastTrieProgressKey = []byte("TrieSync")

	validSectionKey = []byte("count")

	sectionHeadKeyPrefix = []byte("shead")

	// snapshotKeyPrefix is a governance snapshot prefix
	snapshotKeyPrefix = []byte("snapshot")

	// snapshotJournalKey tracks the in-memory diff layers across restarts.
	snapshotJournalKey = []byte("SnapshotJournal")

	// SnapshotGeneratorKey tracks the snapshot generation marker across restarts.
	SnapshotGeneratorKey = []byte("SnapshotGenerator")

	// snapshotDisabledKey flags that the snapshot should not be maintained due to initial sync.
	snapshotDisabledKey = []byte("SnapshotDisabled")

	// snapshotRecoveryKey tracks the snapshot recovery marker across restarts.
	snapshotRecoveryKey = []byte("SnapshotRecovery")

	// snapshotSyncStatusKey tracks the snapshot sync status across restarts.
	snapshotSyncStatusKey = []byte("SnapshotSyncStatus")

	// snapshotRootKey tracks the hash of the last snapshot.
	snapshotRootKey = []byte("SnapshotRoot")

	// badBlockKey tracks the list of bad blocks seen by local
	badBlockKey = []byte("InvalidBlock")

	// Data item prefixes (use single byte to avoid mixing data types, avoid `i`, used for indexes).
	headerPrefix       = []byte("h") // headerPrefix + num (uint64 big endian) + hash -> header
	headerTDSuffix     = []byte("t") // headerPrefix + num (uint64 big endian) + hash + headerTDSuffix -> td
	headerHashSuffix   = []byte("n") // headerPrefix + num (uint64 big endian) + headerHashSuffix -> hash
	headerNumberPrefix = []byte("H") // headerNumberPrefix + hash -> num (uint64 big endian)

	blockBodyPrefix     = []byte("b") // blockBodyPrefix + num (uint64 big endian) + hash -> block body
	blockReceiptsPrefix = []byte("r") // blockReceiptsPrefix + num (uint64 big endian) + hash -> block receipts

	txLookupPrefix        = []byte("l") // txLookupPrefix + hash -> transaction/receipt lookup metadata
	SnapshotAccountPrefix = []byte("a") // SnapshotAccountPrefix + account hash -> account trie value
	SnapshotStoragePrefix = []byte("o") // SnapshotStoragePrefix + account hash + storage hash -> storage trie value
	codePrefix            = []byte("c") // codePrefix + code hash -> contract code

	preimagePrefix = []byte("secure-key-")  // preimagePrefix + hash -> preimage
	configPrefix   = []byte("klay-config-") // config prefix for the db

	pruningEnabledKey        = []byte("PruningEnabled")
	pruningMarkPrefix        = []byte("Pruning-")                                // KIP-111 pruning markings
	pruningMarkValue         = []byte{0x01}                                      // A nonempty value to store a pruning mark
	pruningMarkKeyLen        = len(pruningMarkPrefix) + 8 + common.ExtHashLength // prefix + num (uint64) + node hash
	lastPrunedBlockNumberKey = []byte("lastPrunedBlockNumber")

	// Chain index prefixes (use `i` + single byte to avoid mixing data types).
	BloomBitsIndexPrefix = []byte("iB") // BloomBitsIndexPrefix is the data table of a chain indexer to track its progress

	preimageCounter    = metrics.NewRegisteredCounter("db/preimage/total", nil)
	preimageHitCounter = metrics.NewRegisteredCounter("db/preimage/hits", nil)

	childChainTxHashPrefix          = []byte("ccTxHash")
	lastServiceChainTxReceiptKey    = []byte("LastServiceChainTxReceipt")
	lastIndexedBlockKey             = []byte("LastIndexedBlockKey")
	receiptFromParentChainKeyPrefix = []byte("receiptFromParentChain")

	parentOperatorFeePayerPrefix = []byte("parentOperatorFeePayer")
	childOperatorFeePayerPrefix  = []byte("childOperatorFeePayer")

	valueTransferTxHashPrefix = []byte("vt-tx-hash-key-") // Prefix + hash -> hash

	// bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash -> bloom bits
	bloomBitsPrefix = []byte("B")

	senderTxHashToTxHashPrefix = []byte("SenderTxHash")

	governancePrefix     = []byte("governance")
	governanceHistoryKey = []byte("governanceIdxHistory")
	governanceStateKey   = []byte("governanceState")

	databaseDirPrefix  = []byte("databaseDirectory")
	migrationStatusKey = []byte("migrationStatus")

	stakingInfoPrefix = []byte("stakingInfo")

	supplyCheckpointPrefix        = []byte("supplyCheckpoint")
	lastSupplyCheckpointNumberKey = []byte("lastSupplyCheckpointNumber")

	chaindatafetcherCheckpointKey = []byte("chaindatafetcherCheckpoint")
)

// TxLookupEntry is a positional metadata to help looking up the data content of
// a transaction or receipt given only its hash.
type TxLookupEntry struct {
	BlockHash  common.Hash
	BlockIndex uint64
	Index      uint64
}

// headerKey = headerPrefix + num (uint64 big endian) + hash
func headerKey(number uint64, hash common.Hash) []byte {
	return append(append(headerPrefix, common.Int64ToByteBigEndian(number)...), hash.Bytes()...)
}

// headerTDKey = headerPrefix + num (uint64 big endian) + hash + headerTDSuffix
func headerTDKey(number uint64, hash common.Hash) []byte {
	return append(headerKey(number, hash), headerTDSuffix...)
}

// headerHashKey = headerPrefix + num (uint64 big endian) + headerHashSuffix
func headerHashKey(number uint64) []byte {
	return append(append(headerPrefix, common.Int64ToByteBigEndian(number)...), headerHashSuffix...)
}

// headerNumberKey = headerNumberPrefix + hash
func headerNumberKey(hash common.Hash) []byte {
	return append(headerNumberPrefix, hash.Bytes()...)
}

// blockBodyKey = blockBodyPrefix + num (uint64 big endian) + hash
func blockBodyKey(number uint64, hash common.Hash) []byte {
	return append(append(blockBodyPrefix, common.Int64ToByteBigEndian(number)...), hash.Bytes()...)
}

// blockReceiptsKey = blockReceiptsPrefix + num (uint64 big endian) + hash
func blockReceiptsKey(number uint64, hash common.Hash) []byte {
	return append(append(blockReceiptsPrefix, common.Int64ToByteBigEndian(number)...), hash.Bytes()...)
}

// TxLookupKey = txLookupPrefix + hash
func TxLookupKey(hash common.Hash) []byte {
	return append(txLookupPrefix, hash.Bytes()...)
}

// AccountSnapshotKey = SnapshotAccountPrefix + hash
func AccountSnapshotKey(hash common.Hash) []byte {
	return append(SnapshotAccountPrefix, hash.Bytes()...)
}

// StorageSnapshotKey = SnapshotStoragePrefix + account hash + storage hash
func StorageSnapshotKey(accountHash, storageHash common.Hash) []byte {
	return append(append(SnapshotStoragePrefix, accountHash.Bytes()...), storageHash.Bytes()...)
}

// StorageSnapshotsKey = SnapshotStoragePrefix + account hash + storage hash
func StorageSnapshotsKey(accountHash common.Hash) []byte {
	return append(SnapshotStoragePrefix, accountHash.Bytes()...)
}

func SenderTxHashToTxHashKey(senderTxHash common.Hash) []byte {
	return append(senderTxHashToTxHashPrefix, senderTxHash.Bytes()...)
}

// preimageKey = preimagePrefix + hash
func preimageKey(hash common.Hash) []byte {
	return append(preimagePrefix, hash.Bytes()...)
}

// CodeKey = codePrefix + hash
func CodeKey(hash common.Hash) []byte {
	return append(codePrefix, hash.Bytes()...)
}

// IsCodeKey reports whether the given byte slice is the key of contract code,
// if so return the raw code hash as well.
func IsCodeKey(key []byte) (bool, []byte) {
	if bytes.HasPrefix(key, codePrefix) && len(key) == common.HashLength+len(codePrefix) {
		return true, key[len(codePrefix):]
	}
	return false, nil
}

// configKey = configPrefix + hash
func configKey(hash common.Hash) []byte {
	return append(configPrefix, hash.Bytes()...)
}

func sectionHeadKey(encodedSection []byte) []byte {
	return append(sectionHeadKeyPrefix, encodedSection...)
}

func snapshotKey(hash common.Hash) []byte {
	return append(snapshotKeyPrefix, hash[:]...)
}

func childChainTxHashKey(ccBlockHash common.Hash) []byte {
	return append(childChainTxHashPrefix, ccBlockHash.Bytes()...)
}

func receiptFromParentChainKey(blockHash common.Hash) []byte {
	return append(receiptFromParentChainKeyPrefix, blockHash.Bytes()...)
}

func valueTransferTxHashKey(rTxHash common.Hash) []byte {
	return append(valueTransferTxHashPrefix, rTxHash.Bytes()...)
}

// bloomBitsKey = bloomBitsPrefix + bit (uint16 big endian) + section (uint64 big endian) + hash
func BloomBitsKey(bit uint, section uint64, hash common.Hash) []byte {
	key := append(append(bloomBitsPrefix, make([]byte, 10)...), hash.Bytes()...)

	binary.BigEndian.PutUint16(key[1:], uint16(bit))
	binary.BigEndian.PutUint64(key[3:], section)

	return key
}

func makeKey(prefix []byte, num uint64) []byte {
	byteKey := common.Int64ToByteLittleEndian(num)
	return append(prefix, byteKey...)
}

func databaseDirKey(dbEntryType uint64) []byte {
	return append(databaseDirPrefix, common.Int64ToByteBigEndian(dbEntryType)...)
}

// TrieNodeKey = if Legacy, hash32. Otherwise, exthash
func TrieNodeKey(hash common.ExtHash) []byte {
	if hash.IsZeroExtended() {
		return hash.Unextend().Bytes()
	} else {
		return hash.Bytes()
	}
}

type PruningMark struct {
	Number uint64
	Hash   common.ExtHash
}

// TriePruningMarkKey = prefix + number + hash
// Block number comes first to sort the entries by block numbers.
// Later we can iterate through the marks and extract any given block number range.
func pruningMarkKey(mark PruningMark) []byte {
	bNumber := make([]byte, 8)
	binary.BigEndian.PutUint64(bNumber, mark.Number)
	bHash := mark.Hash.Bytes()
	return append(append(pruningMarkPrefix, bNumber...), bHash...)
}

func parsePruningMarkKey(key []byte) PruningMark {
	if len(key) != pruningMarkKeyLen {
		logger.Crit("Invalid pruningMarkKey", "key", hexutil.Encode(key))
	}
	prefixLen := len(pruningMarkPrefix)
	bNumber := key[prefixLen : prefixLen+8]
	bHash := key[prefixLen+8:]
	return PruningMark{
		Number: binary.BigEndian.Uint64(bNumber),
		Hash:   common.BytesToExtHash(bHash),
	}
}

type SupplyCheckpoint struct {
	Minted   *big.Int
	BurntFee *big.Int
}

func (c *SupplyCheckpoint) Copy() *SupplyCheckpoint {
	return &SupplyCheckpoint{
		Minted:   new(big.Int).Set(c.Minted),
		BurntFee: new(big.Int).Set(c.BurntFee),
	}
}

func supplyCheckpointKey(blockNumber uint64) []byte {
	return append(supplyCheckpointPrefix, common.Int64ToByteBigEndian(blockNumber)...)
}
