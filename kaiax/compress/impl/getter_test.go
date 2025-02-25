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
	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	var (
		n1, n2, n3, n4, n5 = uint64(1), uint64(2), uint64(3), uint64(4), uint64(5)
		h1, h2, h3, h4, h5 = common.HexToHash("0x1"), common.HexToHash("0x2"), common.HexToHash("0x3"), common.HexToHash("0x4"), common.HexToHash("0x5")
		d1, d2, d3, d4     = []byte{1}, []byte{2}, []byte{3}, []byte{4}
		c12                = []ChunkItem{{Num: n1, Hash: h1, Uncompressed: d1}, {Num: n2, Hash: h2, Uncompressed: d2}}
		c34                = []ChunkItem{{Num: n3, Hash: h3, Uncompressed: d3}, {Num: n4, Hash: h4, Uncompressed: d4}}
		cache              = newChunkCache(1, 1)
	)

	// Cache the first chunk's first item.
	cache.add(n1, n2, c12, n1, h1, d1)

	// 1. exact item match.
	data, ok := cache.get(n1, h1)
	assert.True(t, ok)
	assert.Equal(t, d1, data)

	// 2. containing chunk match.
	data, ok = cache.get(n2, h2)
	assert.True(t, ok)
	assert.Equal(t, d2, data)

	// 3. (num,hash) mismatch.
	_, ok = cache.get(n1, h2)
	assert.False(t, ok)

	// 4. uncached item.
	_, ok = cache.get(n5, h5)
	assert.False(t, ok)

	// Cache the second chunk's second item. Evicting the first items.
	cache.add(n3, n4, c34, n4, h4, d4)

	// 1. exact item match.
	data, ok = cache.get(n4, h4)
	assert.True(t, ok)
	assert.Equal(t, d4, data)

	// 2. containing chunk match.
	data, ok = cache.get(n3, h3)
	assert.True(t, ok)
	assert.Equal(t, d3, data)

	// 3. (num,hash) mismatch.
	_, ok = cache.get(n3, h4)
	assert.False(t, ok)

	// 4. uncached item.
	_, ok = cache.get(n1, h1)
	assert.False(t, ok)
}
