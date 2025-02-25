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
	"crypto/rand"
	"math/big"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/require"
)

type mockBlockHashReader struct{}

func (m *mockBlockHashReader) ReadCanonicalHash(num uint64) common.Hash {
	return common.BigToHash(big.NewInt(int64(num)))
}

func TestCompressContext(t *testing.T) {
	var (
		ukv     = database.NewMemDB()
		ckv     = database.NewMemDB()
		mockDbm = &mockBlockHashReader{}
		schema  = NewHeaderSchema(ukv, ckv)
		codec   = NewZstdCodec()
		quit    = atomic.Int32{}
		ctx     = newChunkContext(mockDbm, schema, codec, 10, 1024, 0)
	)

	ctx.until(10, &quit)
}

// Demonstrates that the directory size won't reduce after deletion, but reduces after compacting.
func TestCompact(t *testing.T) {
	dbc := &database.DBConfig{Dir: t.TempDir()}
	kv, err := database.NewLevelDB(dbc, database.BodyDB)
	require.Nil(t, err)

	// Insert and measure
	for i := uint64(0); i < 10000; i++ {
		randval := make([]byte, 1000)
		rand.Read(randval)
		kv.Put(common.Int64ToByteBigEndian(i), randval)
	}
	dirSize(t, dbc.Dir)

	// Delete 80% of the data and measure
	for i := uint64(0); i < 8000; i++ {
		kv.Delete(common.Int64ToByteBigEndian(i))
	}
	dirSize(t, dbc.Dir)

	// Compact and measure
	kv.Compact(nil, nil)
	dirSize(t, dbc.Dir)

	kv.Close()
}

func dirSize(t *testing.T, dir string) int {
	out, err := exec.Command("du", dir).Output()
	require.Nil(t, err)
	size := strings.Split(string(out), "\t")[0]
	sizeInt, err := strconv.Atoi(size)
	require.Nil(t, err)
	t.Logf("dir size: %d", sizeInt)
	return sizeInt
}
