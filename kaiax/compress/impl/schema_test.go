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
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestItemSchema(t *testing.T) {
	var (
		ukv = database.NewMemDB()
		ckv = database.NewMemDB()
		hs  = NewHeaderSchema(ukv, ckv)

		from = uint64(0x1111)
		num  = uint64(0x1234)
		to   = uint64(0x2222)
		hash = common.HexToHash("0xabcd")
	)

	assert.Equal(t, ukv, hs.uncompressedDb())
	assert.Equal(t, ckv, hs.compressedDb())

	// "h" || num || hash
	assert.Equal(t, hexutil.MustDecode("0x680000000000001234000000000000000000000000000000000000000000000000000000000000abcd"), hs.uncompressedKey(num, hash))
	// "Compress-h" || to || from
	assert.Equal(t, hexutil.MustDecode("0x436f6d707265737365642d6800000000000022220000000000001111"), chunkKey(hs.compressedKeyPrefix(), from, to))
}
