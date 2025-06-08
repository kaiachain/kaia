package impl

import (
	"math/big"
	"testing"

	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/storage/database"
)

func BenchmarkGetBlsPubkey(b *testing.B) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	var (
		db     = database.NewMemoryDBManager()
		alloc  = testAllocStorage()
		config = testRandaoForkChainConfig(big.NewInt(0))
		addr   = common.HexToAddress("0x0000000000000000000000000000000000000001")
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	// Generate multiple blocks with the same storage state
	for i := 1; i <= 5; i++ {
		backend.Commit()
	}

	// Benchmark cases
	benchCases := []struct {
		name      string
		setupFunc func() *RandaoModule
	}{
		{
			name: "WithCache",
			setupFunc: func() *RandaoModule {
				mRandao := NewRandaoModule()
				fakeDownloader := &downloader.FakeDownloader{}
				mRandao.Init(&InitOpts{
					ChainConfig: config,
					Chain:       backend.BlockChain(),
					Downloader:  fakeDownloader,
				})
				return mRandao
			},
		},
		{
			name: "WithoutCache",
			setupFunc: func() *RandaoModule {
				mRandao := NewRandaoModule()
				// Disable cache (set size to 1)
				mRandao.blsPubkeyCache, _ = lru.NewARC(1)
				mRandao.storageRootCache, _ = lru.NewARC(1)

				fakeDownloader := &downloader.FakeDownloader{}
				mRandao.Init(&InitOpts{
					ChainConfig: config,
					Chain:       backend.BlockChain(),
					Downloader:  fakeDownloader,
				})
				return mRandao
			},
		},
	}

	for _, bc := range benchCases {
		b.Run(bc.name, func(b *testing.B) {
			mRandao := bc.setupFunc()

			// Warm up with first call (exclude initialization cost)
			_, _ = mRandao.GetBlsPubkey(addr, big.NewInt(1))

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				// Use different block numbers to test various cache patterns
				blockNum := big.NewInt(int64((i % 5) + 1))
				_, _ = mRandao.GetBlsPubkey(addr, blockNum)
			}
		})
	}
}
