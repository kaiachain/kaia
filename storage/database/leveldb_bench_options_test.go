// Modifications Copyright 2024 The Kaia Authors
// Copyright 2018 The klaytn Authors
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

package database

import (
	"math/rand"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func genTempDirForTestDB(b *testing.B) string {
	dir, err := os.MkdirTemp("", "kaia-test-db-bench-")
	if err != nil {
		b.Fatalf("cannot create temporary directory: %v", err)
	}
	return dir
}

func getKaiaLDBOptions() *opt.Options {
	return getLevelDBOptions(&DBConfig{LevelDBCacheSize: 128, OpenFilesLimit: 128})
}

func getKaiaLDBOptionsForGetX(x int) *opt.Options {
	opts := getKaiaLDBOptions()
	opts.WriteBuffer *= x
	opts.BlockCacheCapacity *= x
	opts.OpenFilesCacheCapacity *= x
	opts.DisableBlockCache = true

	return opts
}

func getKaiaLDBOptionsForPutX(x int) *opt.Options {
	opts := getKaiaLDBOptions()
	opts.BlockCacheCapacity *= x
	opts.BlockRestartInterval *= x

	opts.BlockSize *= x
	opts.CompactionExpandLimitFactor *= x
	opts.CompactionL0Trigger *= x
	opts.CompactionTableSize *= x

	opts.CompactionSourceLimitFactor *= x
	opts.Compression = opt.DefaultCompression

	return opts
}

func getKaiaLDBOptionsForBatchX(x int) *opt.Options {
	opts := getKaiaLDBOptions()
	opts.BlockCacheCapacity *= x
	opts.BlockRestartInterval *= x

	opts.BlockSize *= x
	opts.CompactionExpandLimitFactor *= x
	opts.CompactionL0Trigger *= x
	opts.CompactionTableSize *= x

	opts.CompactionSourceLimitFactor *= x
	opts.Compression = opt.DefaultCompression

	return opts
}

// readTypeFunc determines requested index
func benchmarkKaiaOptionsGet(b *testing.B, opts *opt.Options, valueLength, numInsertions, numGets int, readTypeFunc func(int, int) int) {
	b.StopTimer()
	b.ReportAllocs()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	for i := 0; i < numInsertions; i++ {
		bs := []byte(strconv.Itoa(i))
		db.Put(bs, randStrBytes(valueLength))
	}

	b.StartTimer()
	for k := 0; k < b.N; k++ {
		for i := 0; i < numGets; i++ {
			bs := []byte(strconv.Itoa(readTypeFunc(i, numInsertions)))
			_, err := db.Get(bs)
			if err != nil {
				b.Fatalf("get failed: %v", err)
			}
		}
	}
}

func randomRead(currIndex, numInsertions int) int {
	return rand.Intn(numInsertions)
}

func sequentialRead(currIndex, numInsertions int) int {
	return numInsertions - currIndex - 1
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func zipfRead(currIndex, numInsertions int) int {
	zipf := rand.NewZipf(r, 3.14, 2.72, uint64(numInsertions))
	zipfNum := zipf.Uint64()
	return numInsertions - int(zipfNum) - 1
}

const (
	getKaiaValueLegnth   = 250
	getKaiaNumInsertions = 1000 * 100
	getKaiaNumGets       = 1000
)

var getKaiaOptions = [...]struct {
	name          string
	valueLength   int
	numInsertions int
	numGets       int
	opts          *opt.Options
}{
	{"X1", getKaiaValueLegnth, getKaiaNumInsertions, getKaiaNumGets, getKaiaLDBOptionsForGetX(1)},
	{"X2", getKaiaValueLegnth, getKaiaNumInsertions, getKaiaNumGets, getKaiaLDBOptionsForGetX(2)},
	{"X4", getKaiaValueLegnth, getKaiaNumInsertions, getKaiaNumGets, getKaiaLDBOptionsForGetX(4)},
	{"X8", getKaiaValueLegnth, getKaiaNumInsertions, getKaiaNumGets, getKaiaLDBOptionsForGetX(8)},
	//{"X16", getKaiaValueLegnth, getKaiaNumInsertions, getKaianumGets, getKaiaLDBOptionsForGetX(16)},
	//{"X32", getKaiaValueLegnth, getKaiaNumInsertions, getKaianumGets, getKaiaLDBOptionsForGetX(32)},
	//{"X64", getKaiaValueLegnth, getKaiaNumInsertions, getKaianumGets, getKaiaLDBOptionsForGetX(64)},
}

func Benchmark_KaiaOptions_Get_Random(b *testing.B) {
	for _, bm := range getKaiaOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKaiaOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, randomRead)
		})
	}
}

func Benchmark_KaiaOptions_Get_Sequential(b *testing.B) {
	for _, bm := range getKaiaOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKaiaOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, sequentialRead)
		})
	}
}

func Benchmark_KaiaOptions_Get_Zipf(b *testing.B) {
	for _, bm := range getKaiaOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKaiaOptionsGet(b, bm.opts, bm.valueLength, bm.numInsertions, bm.numGets, zipfRead)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////// Put Insertion Tests Beginning //////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func benchmarkKaiaOptionsPut(b *testing.B, opts *opt.Options, valueLength, numInsertions int) {
	b.StopTimer()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		for k := 0; k < numInsertions; k++ {
			db.Put(randStrBytes(32), randStrBytes(valueLength))
		}
	}
}

func Benchmark_KaiaOptions_Put(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	const (
		putKaiaValueLegnth   = 250
		putKaiaNumInsertions = 1000 * 10
	)

	putKaiaOptions := [...]struct {
		name          string
		valueLength   int
		numInsertions int
		opts          *opt.Options
	}{
		{"X1", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(1)},
		{"X2", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(2)},
		{"X4", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(4)},
		{"X8", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(8)},
		//{"X16", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(16)},
		//{"X32", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(32)},
		//{"X64", putKaiaValueLegnth, putKaiaNumInsertions, getKaiaLDBOptionsForPutX(64)},
	}
	for _, bm := range putKaiaOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKaiaOptionsPut(b, bm.opts, bm.valueLength, bm.numInsertions)
		})
	}
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////// Put Insertion Tests End /////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////// SHARDED PUT INSERTION TESTS BEGINNING //////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func removeDirs(dirs []string) {
	for _, dir := range dirs {
		os.RemoveAll(dir)
	}
}

func genDatabases(b *testing.B, dirs []string, opts *opt.Options) []*levelDB {
	databases := make([]*levelDB, len(dirs), len(dirs))
	for i := 0; i < len(dirs); i++ {
		databases[i], _ = NewLevelDBWithOption(dirs[i], opts)
	}
	return databases
}

func closeDBs(databases []*levelDB) {
	for _, db := range databases {
		db.Close()
	}
}

func genKeysAndValues(valueLength, numInsertions int) ([][]byte, [][]byte) {
	keys := make([][]byte, numInsertions, numInsertions)
	values := make([][]byte, numInsertions, numInsertions)
	for i := 0; i < numInsertions; i++ {
		keys[i] = randStrBytes(32)
		values[i] = randStrBytes(valueLength)
	}
	return keys, values
}

///////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////// Batch Insertion Tests Beginning /////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

func benchmarkKaiaOptionsBatch(b *testing.B, opts *opt.Options, valueLength, numInsertions int) {
	b.StopTimer()
	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		keys, values := genKeysAndValues(valueLength, numInsertions)
		b.StartTimer()
		batch := db.NewBatch()
		for k := 0; k < numInsertions; k++ {
			batch.Put(keys[k], values[k])
		}
		batch.Write()
	}
}

func Benchmark_KaiaOptions_Batch(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()

	const (
		batchValueLength       = 250
		batchKaiaNumInsertions = 1000 * 10
	)

	putKaiaOptions := [...]struct {
		name          string
		valueLength   int
		numInsertions int
		opts          *opt.Options
	}{
		{"X1", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(1)},
		{"X2", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(2)},
		{"X4", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(4)},
		{"X8", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(8)},
		//{"X16", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(16)},
		//{"X32", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(32)},
		//{"X64", batchValueLength, batchKaiaNumInsertions, getKaiaLDBOptionsForBatchX(64)},
	}
	for _, bm := range putKaiaOptions {
		b.Run(bm.name, func(b *testing.B) {
			benchmarkKaiaOptionsBatch(b, bm.opts, bm.valueLength, bm.numInsertions)
		})
	}
}

// TODO-Kaia = Add a test for checking GoRoutine Overhead

///////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////// Batch Insertion Tests End /////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

///////////////////////////////////////////////////////////////////////////////////////////
////////////////////////// Ideal Batch Size Tests Begins //////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////

type idealBatchBM struct {
	name      string
	totalSize int
	batchSize int
	rowSize   int
}

func benchmarkIdealBatchSize(b *testing.B, bm idealBatchBM) {
	b.StopTimer()

	dir := genTempDirForTestDB(b)
	defer os.RemoveAll(dir)

	opts := getKaiaLDBOptions()
	db, err := NewLevelDBWithOption(dir, opts)
	require.NoError(b, err)
	defer db.Close()

	b.StartTimer()

	var wg sync.WaitGroup
	numBatches := bm.totalSize / bm.batchSize
	wg.Add(numBatches)
	for i := 0; i < numBatches; i++ {
		batch := db.NewBatch()
		for k := 0; k < bm.batchSize; k++ {
			batch.Put(randStrBytes(32), randStrBytes(bm.rowSize))
		}

		go func(currBatch Batch) {
			defer wg.Done()
			currBatch.Write()
		}(batch)
	}
	wg.Wait()
}

func Benchmark_IdealBatchSize(b *testing.B) {
	b.StopTimer()
	// please change below rowSize to change the size of an input row
	// key = 32 bytes, value = rowSize bytes
	const rowSize = 250

	benchmarks := []idealBatchBM{
		// to run test with total size smaller than 1,000 rows
		// go test -bench=Benchmark_IdealBatchSize/SmallBatches
		{"SmallBatches_100Rows_10Batch_250Bytes", 100, 10, rowSize},
		{"SmallBatches_100Rows_20Batch_250Bytes", 100, 20, rowSize},
		{"SmallBatches_100Rows_25Batch_250Bytes", 100, 25, rowSize},
		{"SmallBatches_100Rows_50Batch_250Bytes", 100, 50, rowSize},
		{"SmallBatches_100Rows_100Batch_250Bytes", 100, 100, rowSize},

		{"SmallBatches_200Rows_10Batch_250Bytes", 200, 10, rowSize},
		{"SmallBatches_200Rows_20Batch_250Bytes", 200, 20, rowSize},
		{"SmallBatches_200Rows_25Batch_250Bytes", 200, 25, rowSize},
		{"SmallBatches_200Rows_50Batch_250Bytes", 200, 50, rowSize},
		{"SmallBatches_200Rows_100Batch_250Bytes", 200, 100, rowSize},

		{"SmallBatches_400Rows_10Batch_250Bytes", 400, 10, rowSize},
		{"SmallBatches_400Rows_20Batch_250Bytes", 400, 20, rowSize},
		{"SmallBatches_400Rows_25Batch_250Bytes", 400, 25, rowSize},
		{"SmallBatches_400Rows_50Batch_250Bytes", 400, 50, rowSize},
		{"SmallBatches_400Rows_100Batch_250Bytes", 400, 100, rowSize},

		{"SmallBatches_800Rows_10Batch_250Bytes", 800, 10, rowSize},
		{"SmallBatches_800Rows_20Batch_250Bytes", 800, 20, rowSize},
		{"SmallBatches_800Rows_25Batch_250Bytes", 800, 25, rowSize},
		{"SmallBatches_800Rows_50Batch_250Bytes", 800, 50, rowSize},
		{"SmallBatches_800Rows_100Batch_250Bytes", 800, 100, rowSize},

		// to run test with total size between than 1k rows ~ 10k rows
		// go test -bench=Benchmark_IdealBatchSize/LargeBatches
		{"LargeBatches_1kRows_100Batch_250Bytes", 1000, 100, rowSize},
		{"LargeBatches_1kRows_200Batch_250Bytes", 1000, 200, rowSize},
		{"LargeBatches_1kRows_250Batch_250Bytes", 1000, 250, rowSize},
		{"LargeBatches_1kRows_500Batch_250Bytes", 1000, 500, rowSize},
		{"LargeBatches_1kRows_1000Batch_250Bytes", 1000, 1000, rowSize},

		{"LargeBatches_2kRows_100Batch_250Bytes", 2000, 100, rowSize},
		{"LargeBatches_2kRows_200Batch_250Bytes", 2000, 200, rowSize},
		{"LargeBatches_2kRows_250Batch_250Bytes", 2000, 250, rowSize},
		{"LargeBatches_2kRows_500Batch_250Bytes", 2000, 500, rowSize},
		{"LargeBatches_2kRows_1000Batch_250Bytes", 2000, 1000, rowSize},

		{"LargeBatches_4kRows_100Batch_250Bytes", 4000, 100, rowSize},
		{"LargeBatches_4kRows_200Batch_250Bytes", 4000, 200, rowSize},
		{"LargeBatches_4kRows_250Batch_250Bytes", 4000, 250, rowSize},
		{"LargeBatches_4kRows_500Batch_250Bytes", 4000, 500, rowSize},
		{"LargeBatches_4kRows_1000Batch_250Bytes", 4000, 1000, rowSize},

		{"LargeBatches_8kRows_100Batch_250Bytes", 8000, 100, rowSize},
		{"LargeBatches_8kRows_200Batch_250Bytes", 8000, 200, rowSize},
		{"LargeBatches_8kRows_250Batch_250Bytes", 8000, 250, rowSize},
		{"LargeBatches_8kRows_500Batch_250Bytes", 8000, 500, rowSize},
		{"LargeBatches_8kRows_1000Batch_250Bytes", 8000, 1000, rowSize},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for m := 0; m < b.N; m++ {
				benchmarkIdealBatchSize(b, bm)
			}
		})
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStrBytes(n int) []byte {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

func getShardForTest(keys [][]byte, index, numShards int) int64 {
	return int64(index % numShards)
	// TODO-Kaia: CHANGE BELOW LOGIC FROM ROUND-ROBIN TO USE getShardForTest
	//key := keys[index]
	//hashString := strings.TrimPrefix(common.Bytes2Hex(key),"0x")
	//if len(hashString) > 15 {
	//	hashString = hashString[:15]
	//}
	//seed, _ := strconv.ParseInt(hashString, 16, 64)
	//shard := seed % int64(numShards)
	//
	//return shard
}
