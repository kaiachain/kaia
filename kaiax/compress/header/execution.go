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

func (hc *HeaderCompressModule) Compress() {
	compress_interface.Compress(hc, database.HeaderCompressType, hc.Dbm.CompressHeader)
}

func (hc *HeaderCompressModule) RewindTo(newBlock *types.Block) {}

func (hc *HeaderCompressModule) RewindDelete(hash common.Hash, num uint64) {
	if err := hc.Dbm.DeleteHeaderFromChunk(num, hash); err != nil {
		compress.Logger.Warn("[Header Compression] Failed to delete header", "blockNum", num, "blockHash", hash.String())
	}
}

func (hc *HeaderCompressModule) TestCopyOriginData(copyTestDB database.Batch, from, to uint64) {
	// Copy origin header
	for i := from; i < to; i++ {
		hash := hc.Dbm.ReadCanonicalHash(i)
		originHeader := hc.Dbm.ReadHeader(hash, i)
		hc.Dbm.PutHeaderToBatch(copyTestDB, hash, i, originHeader)
	}
}

func (hc *HeaderCompressModule) TestVerifyCompressionIntegrity(from, to uint64) error {
	for i := from; i < to; i++ {
		originHeader := hc.Dbm.ReadHeader(hc.Dbm.ReadCanonicalHash(i), i)
		compressedHeader, err := hc.Dbm.FindHeaderFromChunkWithBlkHash(i, originHeader.Hash())
		if err != nil {
			return err
		}
		if originHeader.Hash() != compressedHeader.Hash() {
			return fmt.Errorf("[Header Compression Test] Compressed header is not the same data with origin header data (number=%d, txHash=%s)", i, originHeader.Hash().String())
		}
	}
	return nil
}
