package compress

import (
	"fmt"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	compress_interface "github.com/kaiachain/kaia/kaiax/compress/interface"
	"github.com/kaiachain/kaia/storage/database"
)

func (rc *ReceiptCompressModule) Compress() {
	compress_interface.Compress(rc, database.ReceiptCompressType, rc.Dbm.CompressReceipts)
}

func (rc *ReceiptCompressModule) RewindTo(newBlock *types.Block) {}

func (rc *ReceiptCompressModule) RewindDelete(hash common.Hash, num uint64) {
	if err := rc.Dbm.DeleteReceiptsFromChunk(num, hash); err != nil {
		compress.Logger.Warn("[Receipt Compression] Failed to delete receipt", "blockNum", num, "blockHash", hash.String())
	}
}

func (rc *ReceiptCompressModule) TestCopyOriginData(copyTestDB database.Batch, from, to uint64) {
	// Copy origin receipts
	for i := from; i < to; i++ {
		hash := rc.Dbm.ReadCanonicalHash(i)
		originReceipts := rc.Dbm.ReadReceipts(hash, i)
		rc.Dbm.PutReceiptsToBatch(copyTestDB, hash, i, originReceipts)
	}
}

func (rc *ReceiptCompressModule) TestVerifyCompressionIntegrity(from, to uint64) error {
	for i := from; i < to; i++ {
		for _, originReceipt := range rc.Dbm.ReadReceipts(rc.Dbm.ReadCanonicalHash(i), i) {
			compressedReceipt, err := rc.Dbm.FindReceiptFromChunkWithTxHash(i, originReceipt.TxHash)
			if err != nil {
				return err
			}
			if !reflect.DeepEqual(originReceipt, compressedReceipt) {
				return fmt.Errorf("[Receipt Compression Test] Compressed receipt is not the same data with origin receipt data (number=%d, txHash=%s)", i, originReceipt.TxHash.String())
			}
		}
	}
	return nil
}
