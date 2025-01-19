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
	"math"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (c *CompressModule) RewindTo(newBlock *types.Block) {}

func (c *CompressModule) RewindDelete(hash common.Hash, num uint64) {
	if !c.Enable {
		return
	}
	var (
		headerRewinding   bool
		bodyRewinding     bool
		receiptsRewinding bool
	)
	if c.lastHeaderRewindNum > num {
		go c.deleteHeader(hash, num)
		headerRewinding = true
	}
	if c.lastBodyRewindNum > num {
		go c.deleteBody(hash, num)
		bodyRewinding = true
	}
	if c.lastReceiptsRewindNum > num {
		go c.deleteReceipts(hash, num)
		receiptsRewinding = true
	}
	var (
		hd ChunkDelete
		bd ChunkDelete
		rd ChunkDelete
	)
	if bodyRewinding {
		bd = <-c.bodyChunkDeleted
	}
	if receiptsRewinding {
		rd = <-c.receiptsChunkDeleted
	}
	if headerRewinding {
		hd = <-c.headerChunkDeleted
	}
	// Ovewrite subsequent block number to new starting number which contains compression range before
	if hd.shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, HeaderCompressType, hd.subsequentCompressBlockNum)
		c.lastHeaderRewindNum = hd.subsequentCompressBlockNum
	}
	if bd.shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, BodyCompressType, bd.subsequentCompressBlockNum)
		c.lastBodyRewindNum = hd.subsequentCompressBlockNum
	}
	if rd.shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, ReceiptCompressType, rd.subsequentCompressBlockNum)
		c.lastReceiptsRewindNum = hd.subsequentCompressBlockNum
	}
}

func (c *CompressModule) clearRewindHistory() {
	c.lastHeaderRewindNum = math.MaxUint64
	c.lastBodyRewindNum = math.MaxUint64
	c.lastReceiptsRewindNum = math.MaxUint64
}
