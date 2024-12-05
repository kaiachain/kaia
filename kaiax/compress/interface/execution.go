package compress

import (
	"path/filepath"
	"time"

	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
)

func Compress[T CompressModuleInterface](m T, compressTyp database.CompressionType, compressFn database.CompressFn) {
	for {
		var (
			curBlkNum             = m.GetChain().CurrentBlock().NumberU64()
			residualBlkCnt        = curBlkNum % database.CompressMigrationThreshold
			nextCompressionBlkNum = m.GetDbm().ReadSubsequentCompressionBlkNumber(compressTyp)
			// Do not wait if next compression block number is far awway. Start migration right now
			noWait = curBlkNum > nextCompressionBlkNum && curBlkNum-nextCompressionBlkNum > database.CompressMigrationThreshold
		)

		if residualBlkCnt != 0 && !noWait {
			time.Sleep(time.Second * time.Duration(database.CompressMigrationThreshold-residualBlkCnt))
			continue
		}
		from, to := uint64(0), uint64(0)
		for {
			subsequentBlkNumber, err := compressFn(from, to, curBlkNum, true)
			if err != nil {
				compress.Logger.Warn("[Compression] failed to compress chunk", "err", err)
				break
			}
			if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
				compress.Logger.Info("[Compression] compression is completed", "from", from, "to", to, "subsequentBlkNumber", subsequentBlkNumber)
				break
			}
			from = subsequentBlkNumber
		}
	}
}

func TestCompress[T CompressModuleInterface](m T, compressTyp database.CompressionType, compressFn database.CompressFn, from, to uint64, startNum *uint64, tempDir string) error {
	dbConfig := m.GetDbm().GetDBConfig()
	copyTestDB, err := database.TestCreateNewDB(dbConfig, filepath.Join(dbConfig.Dir, tempDir))
	if err != nil {
		return err
	}
	defer copyTestDB.Release()
	if startNum != nil {
		m.GetDbm().WriteSubsequentCompressionBlkNumber(compressTyp, *startNum)
	}
	curBlkNum := m.GetChain().CurrentBlock().NumberU64()
	for {
		// Copy origin receipts
		m.TestCopyOriginData(copyTestDB, from, to)
		subsequentBlkNumber, err := compressFn(from, to, curBlkNum, true)
		if err != nil {
			return err
		}
		if err := m.TestVerifyCompressionIntegrity(from, to); err != nil {
			return err
		}
		if subsequentBlkNumber >= curBlkNum || subsequentBlkNumber >= to {
			compress.Logger.Info("[Compression] compression is completed", "from", from, "to", to, "subsequentBlkNumber", subsequentBlkNumber)
			break
		}
		from = subsequentBlkNumber
	}
	if _, err := database.WriteBatches(copyTestDB); err != nil {
		return err
	}
	return nil
}
