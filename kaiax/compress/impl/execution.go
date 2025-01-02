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
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (c *CompressModule) RewindTo(newBlock *types.Block) {}

func (c *CompressModule) RewindDelete(hash common.Hash, num uint64) {
	// Ovewrite subsequent block number to new starting number which contains compression range before
	if newHeaderFromAfterRewind, shouldUpdate := c.deleteHeader(hash, num); shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, HeaderCompressType, newHeaderFromAfterRewind)
	}
	if newBodyFromAfterRewind, shouldUpdate := c.deleteBody(hash, num); shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, BodyCompressType, newBodyFromAfterRewind)
	}
	if newReceiptsFromAfterRewind, shouldUpdate := c.deleteReceipts(hash, num); shouldUpdate {
		writeSubsequentCompressionBlkNumber(c.Dbm, ReceiptCompressType, newReceiptsFromAfterRewind)
	}
}
