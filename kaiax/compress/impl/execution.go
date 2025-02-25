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

package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (c *CompressModule) RewindTo(newBlock *types.Block) {
}

func (c *CompressModule) RewindDelete(hash common.Hash, num uint64) {
	for _, schema := range c.schemas {
		c.decompressChunk(schema, num)
	}
}

// decompressChunk permanently decompresses the chunk that contains the given number.
func (c *CompressModule) decompressChunk(schema ItemSchema, num uint64) {
	cData, from, to, ok := readCompressed(schema, num)
	if !ok {
		return
	}
	chunk, err := decompressChunk(c.codec, cData)
	if err != nil {
		return
	}
	writeUncompressedBatch(schema, chunk)
	writeNextNum(schema, from)
	deleteCompressed(schema, from, to)
}
