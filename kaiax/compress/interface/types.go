package compress

import (
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
)

type CompressModuleInterface interface {
	GetChain() compress.BlockChain
	GetDbm() database.DBManager

	// Unit test functions
	TestCopyOriginData(copyTestDB database.Batch, from, to uint64)
	TestVerifyCompressionIntegrity(from, to uint64) error
}
