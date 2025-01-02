package compress

import (
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
)

const GB_1 = blockchain.MB_1 * 1000

var (
	compressHeaderCache   common.Cache
	compressBodyCache     common.Cache
	compressReceiptsCache common.Cache
)

func initCache(capSize uint64) {
	if capSize == 0 {
		capSize = blockchain.DefaultCompressChunkCap
	}
	// the maximum cache size is limited to 1GB
	capped := GB_1 / capSize
	if capSize > GB_1 {
		capped = 1
	}
	cacheConfig := common.FIFOCacheConfig{CacheSize: int(capped), IsScaled: true}
	compressHeaderCache = common.NewCache(cacheConfig)
	compressBodyCache = common.NewCache(cacheConfig)
	compressReceiptsCache = common.NewCache(cacheConfig)
	logger.Info("[Compression Module] Cache size is set", "capSize", capSize, "nSlots", capped)
}

func getFromCache(compressTyp CompressionType, from, to uint64) (any, bool) {
	cacheKey := getCacheKey(HeaderCompressType, from, to)
	switch compressTyp {
	case HeaderCompressType:
		if compressHeaderCache != nil {
			return compressHeaderCache.Get(cacheKey)
		}
	case BodyCompressType:
		if compressBodyCache != nil {
			return compressBodyCache.Get(cacheKey)
		}
	case ReceiptCompressType:
		if compressReceiptsCache != nil {
			return compressReceiptsCache.Get(cacheKey)
		}
	}
	return nil, false
}

func addCache(compressTyp CompressionType, from, to uint64, compressed any) {
	cacheKey := getCacheKey(HeaderCompressType, from, to)
	switch compressTyp {
	case HeaderCompressType:
		if compressHeaderCache != nil {
			compressHeaderCache.Add(cacheKey, compressed)
		}
	case BodyCompressType:
		if compressBodyCache != nil {
			compressBodyCache.Add(cacheKey, compressed)
		}
	case ReceiptCompressType:
		if compressReceiptsCache != nil {
			compressReceiptsCache.Add(cacheKey, compressed)
		}
	}
}

func clearCache() {
	if compressHeaderCache != nil {
		compressHeaderCache.Purge()
	}
	if compressBodyCache != nil {
		compressBodyCache.Purge()
	}
	if compressReceiptsCache != nil {
		compressReceiptsCache.Purge()
	}
}
