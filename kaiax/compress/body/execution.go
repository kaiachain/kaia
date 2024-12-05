// Copyright 2024 The Kaia Authors
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

package compress

import (
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/compress"
	compress_interface "github.com/kaiachain/kaia/kaiax/compress/interface"
	"github.com/kaiachain/kaia/storage/database"
)

func (bc *BodyCompressModule) Compress() {
	compress_interface.Compress(bc, database.BodyCompressType, bc.Dbm.CompressBody)
}

func (bc *BodyCompressModule) RewindTo(newBlock *types.Block) {}

func (bc *BodyCompressModule) RewindDelete(hash common.Hash, num uint64) {
	if err := bc.Dbm.DeleteBodyFromChunk(num, hash); err != nil {
		compress.Logger.Warn("[Body Compression] Failed to delete body", "blockNum", num, "blockHash", hash.String())
	}
}

func (bc *BodyCompressModule) TestCopyOriginData(copyTestDB database.Batch, from, to uint64) {
	// Copy origin body
	for i := from; i < to; i++ {
		hash := bc.Dbm.ReadCanonicalHash(i)
		originBody := bc.Dbm.ReadBody(hash, i)
		bc.Dbm.PutBodyToBatch(copyTestDB, hash, i, originBody)
	}
}

func (bc *BodyCompressModule) TestVerifyCompressionIntegrity(from, to uint64) error {
	for i := from; i < to; i++ {
		originBody := bc.Dbm.ReadBody(bc.Dbm.ReadCanonicalHash(i), i)
		for _, originTx := range originBody.Transactions {
			compressedTx, err := bc.Dbm.FindTxFromChunkWithTxHash(i, originTx.Hash())
			if err != nil {
				return err
			}
			if !originTx.Equal(compressedTx) {
				return fmt.Errorf("[Body Compression Test] Compressed body is not the same data with origin body data (number=%d, txHash=%s)", i, originTx.Hash().String())
			}
		}
	}
	return nil
}
