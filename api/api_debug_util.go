// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/storage/database"
)

var (
	errCompactRangeZeroStep = errors.New("step cannot be zero")
	errCompactRangeStartEnd = errors.New("start must be less than end")
)

// DebugUtilAPI is the collection of Kaia APIs exposed over the private
// debugging endpoint.
type DebugUtilAPI struct {
	b Backend
}

// NewDebugUtilAPI creates a new API definition for the private debug methods
// of the Kaia service.
func NewDebugUtilAPI(b Backend) *DebugUtilAPI {
	return &DebugUtilAPI{b: b}
}

// ChaindbProperty returns leveldb properties of the chain database.
func (api *DebugUtilAPI) ChaindbProperty(property string) (string, error) {
	return api.b.ChainDB().Stat(property)
}

type ChaindbCompactArgs struct {
	Preset *string           `json:"preset"`
	Ranges []CompactionRange `json:"ranges"`
}

// ChaindbCompact flattens the entire key-value database into a single level,
// removing all unused slots and merging all keys.
func (api *DebugUtilAPI) ChaindbCompact(args *ChaindbCompactArgs) error {
	if args == nil {
		sDefault := "default"
		args = &ChaindbCompactArgs{Preset: &sDefault}
	}
	preset := args.Preset
	ranges := args.Ranges

	if preset == nil {
		ranges = compactionPresetDefault
	} else if *preset == "default" {
		ranges = compactionPresetDefault
	} else if *preset == "custom" {
		if len(ranges) == 0 {
			return errors.New("no ranges provided for custom preset")
		}
	} else if r, ok := compactionPresets[*preset]; ok {
		ranges = r
	} else {
		return fmt.Errorf("unknown preset: %s", *preset)
	}

	for _, r := range ranges {
		if err := r.compactRange(context.Background(), api.b.ChainDB()); err != nil {
			return err
		}
	}
	return nil
}

// SetHead rewinds the head of the blockchain to a previous block.
func (api *DebugUtilAPI) SetHead(number rpc.BlockNumber) error {
	if number == rpc.PendingBlockNumber ||
		number == rpc.LatestBlockNumber ||
		number.Uint64() > api.b.CurrentBlock().NumberU64() {
		return errors.New("Cannot rewind to future")
	}
	return api.b.SetHead(uint64(number))
}

// PrintBlock retrieves a block and returns its pretty printed form.
func (api *DebugUtilAPI) PrintBlock(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (string, error) {
	block, _ := api.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if block == nil {
		blockNumberOrHashString, _ := blockNrOrHash.NumberOrHashString()
		return "", fmt.Errorf("block %v not found", blockNumberOrHashString)
	}
	return spew.Sdump(block), nil
}

type CompactionRange struct {
	Start    hexutil.Bytes `json:"start"`
	Step     hexutil.Bytes `json:"step"`
	End      hexutil.Bytes `json:"end"`
	DB       []string      `json:"db"` // if empty, all database.
	Interval string        `json:"interval"`
}

// CompactRange compacts over the given range, for all specified databases.
func (r *CompactionRange) compactRange(ctx context.Context, chainDB database.DBManager) error {
	var interval time.Duration
	if r.Interval != "" {
		if duration, err := time.ParseDuration(r.Interval); err != nil {
			return err
		} else {
			interval = duration
		}
	}

	fn := func(segStart, segEnd []byte) error {
		sRange := fmt.Sprintf("0x%x-0x%x", segStart, segEnd)
		logger.Info("Compacting database started", "db", r.DB, "range", sRange)
		tStart := time.Now()
		if err := r.compactSegment(ctx, chainDB, segStart, segEnd); err != nil {
			logger.Error("Database compaction failed", "db", r.DB, "range", sRange, "err", err)
			return err
		}
		tElapsed := common.PrettyDuration(time.Since(tStart))
		logger.Info("Compacting database completed", "db", r.DB, "range", sRange, "elapsed", tElapsed)
		time.Sleep(interval)
		return nil
	}
	if err := iterateRange(r.Start, r.Step, r.End, fn); err != nil {
		return err
	}
	return nil
}

// compactSegment performs the compaction for a single segment, for all specified databases.
func (r *CompactionRange) compactSegment(ctx context.Context, chainDB database.DBManager, segStart, segEnd []byte) error {
	if len(r.DB) == 0 {
		return chainDB.Compact(segStart, segEnd)
	}

	for _, db := range r.DB {
		if db, err := getDatabase(chainDB, db); err != nil {
			return err
		} else if compactErr := db.Compact(segStart, segEnd); compactErr != nil {
			return compactErr
		}
	}
	return nil
}

func getDatabase(dbm database.DBManager, dbName string) (database.Database, error) {
	var db database.Database
	switch dbName {
	case "misc":
		db = dbm.GetMiscDB()
	case "header":
		db = dbm.GetHeaderDB()
	case "body":
		db = dbm.GetBodyDB()
	case "receipts":
		db = dbm.GetReceiptsDB()
	case "txlookup":
		db = dbm.GetTxLookupEntryDB()
	case "statetrie":
		db = dbm.GetStateTrieDB()
	default:
		return nil, fmt.Errorf("unknown database name: %s", dbName)
	}
	if db == nil {
		return nil, fmt.Errorf("cannot get database: %s", dbName)
	}
	return db, nil
}

func padBytes(b []byte, padLen int) []byte {
	if len(b) >= padLen {
		return b
	} else {
		return append(b, bytes.Repeat([]byte{0}, padLen-len(b))...)
	}
}

func padBig(b *big.Int, padLen int) []byte {
	return padBytes(b.Bytes(), padLen)
}

func padRange(start, step, end []byte) (startBig, stepBig, endBig *big.Int, padLen int, err error) {
	padLen = max(len(start), len(step), len(end))

	// if start is empty, it means the range starts from zero.
	startBig = new(big.Int).SetBytes(padBytes(start, padLen))

	stepBig = new(big.Int).SetBytes(padBytes(step, padLen))
	if stepBig.Sign() == 0 {
		err = errCompactRangeZeroStep
		return
	}

	endInfinity := len(end) == 0
	if endInfinity {
		// infinity is represented as 0x100..00 (1 followed by padLen zero bytes), just outside the range of prefixes.
		endBig = new(big.Int).SetBytes(padBytes([]byte{1}, padLen+1))
	} else {
		endBig = new(big.Int).SetBytes(padBytes(end, padLen))
	}
	if !endInfinity && startBig.Cmp(endBig) > 0 {
		err = errCompactRangeStartEnd
		return
	}

	return
}

func iterateRange(start, step, end []byte, fn func(segStart, segEnd []byte) error) error {
	endWithNil := len(end) == 0
	rangeStart, rangeStep, rangeEnd, padLen, err := padRange(start, step, end)
	if err != nil {
		return err
	}

	segStart := rangeStart
	segEnd := new(big.Int).Add(rangeStart, rangeStep)

	for segEnd.Cmp(rangeEnd) < 0 {
		if err := fn(padBig(segStart, padLen), padBig(segEnd, padLen)); err != nil {
			return err
		}

		segStart.Add(segStart, rangeStep)
		segEnd.Add(segEnd, rangeStep)
	}

	// last segment. run to the rangeEnd.
	if endWithNil {
		return fn(padBig(segStart, padLen), nil)
	} else {
		return fn(padBig(segStart, padLen), padBig(rangeEnd, padLen))
	}
}

var (
	compactionPresets = map[string][]CompactionRange{
		"default":     compactionPresetDefault,
		"allbutstate": slices.Concat(compactionPresetHeader, compactionPresetBody, compactionPresetReceipts),
		"header":      compactionPresetHeader,
		"body":        compactionPresetBody,
		"receipts":    compactionPresetReceipts,
	}
	compactionPresetDefault = []CompactionRange{
		{
			Start: hexutil.MustDecode("0x"),
			Step:  hexutil.MustDecode("0x01"),
			End:   hexutil.MustDecode("0x"),
		},
	}
	compactionPresetHeader = []CompactionRange{
		{
			// headerNumberPrefix("H") + hash -> num (uint64 big endian)
			Start: hexutil.MustDecode("0x48"),
			Step:  hexutil.MustDecode("0x01"),
			End:   hexutil.MustDecode("0x49"),
			DB:    []string{"header"},
		},
		{
			// headerPrefix("h") + num (uint64 big endian) + hash -> header
			Start: hexutil.MustDecode("0x68"),
			Step:  hexutil.MustDecode("0x01"),
			End:   hexutil.MustDecode("0x69"),
			DB:    []string{"header"},
		},
	}
	compactionPresetBody = []CompactionRange{
		{
			// blockBodyPrefix("b") + num (uint64 big endian) + hash -> block body
			Start: hexutil.MustDecode("0x62"),
			Step:  hexutil.MustDecode("0x01"),
			End:   hexutil.MustDecode("0x63"),
			DB:    []string{"body"},
		},
	}
	compactionPresetReceipts = []CompactionRange{
		{
			// blockReceiptsPrefix("r") + num (uint64 big endian) + hash -> block receipts
			Start: hexutil.MustDecode("0x72"),
			Step:  hexutil.MustDecode("0x01"),
			End:   hexutil.MustDecode("0x73"),
			DB:    []string{"receipts"},
		},
	}
	// compactionPresetReceiptsStepwise = []CompactionRange{
	// 	{
	// 		// e.g. 0x000000000b5d64b0 = num 190,670,000
	// 		// step 0x0000000001000000 = num  16,777,216
	// 		Start: hexutil.MustDecode("0x72"), // blockReceiptsPrefix("r")
	// 		Step:  hexutil.MustDecode("0x00"+"0x0000000001"),
	// 		End:   hexutil.MustDecode("0x72"+"0x000000000f"),
	// 		DB:    []string{"receipts"},
	// 	},
	// }
)
