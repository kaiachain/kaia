// Copyright 2026 The Kaia Authors
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

package database

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func filledMemDB(entries int, valueSize int) *MemDB {
	db := NewMemDB()
	val := bytes.Repeat([]byte("v"), valueSize)
	for i := 0; i < entries; i++ {
		_ = db.Put([]byte(fmt.Sprintf("k%06d", i)), val)
	}
	return db
}

// TestCopyDB_InterruptedReturnsSentinel verifies that copyDB returns
// ErrDBMigrationInterrupted when the quit channel is closed mid-copy,
// instead of nil. Without this, callers cannot distinguish an interrupted
// migration from a successful one.
func TestCopyDB_InterruptedReturnsSentinel(t *testing.T) {
	src := filledMemDB(64, 32)
	dst := NewMemDB()

	quit := make(chan struct{})
	close(quit)

	err := copyDB("single", src, dst, quit)
	require.Error(t, err, "copyDB must return an error when migration is interrupted")
	assert.True(t, errors.Is(err, ErrDBMigrationInterrupted),
		"copyDB must return ErrDBMigrationInterrupted, got: %v", err)
}

// TestCopyDB_CompletesReturnsNil is a positive control: a non-interrupted
// run must still return nil.
func TestCopyDB_CompletesReturnsNil(t *testing.T) {
	src := filledMemDB(8, 16)
	dst := NewMemDB()

	err := copyDB("single", src, dst, make(chan struct{}))
	require.NoError(t, err)
}

// failingBatch is a Batch whose Put always returns an error.
// Used to force copyDB into the error path without depending on a real backend.
type failingBatch struct{ err error }

func (b *failingBatch) Put(_, _ []byte) error       { return b.err }
func (b *failingBatch) Delete(_ []byte) error       { return nil }
func (b *failingBatch) ValueSize() int              { return 0 }
func (b *failingBatch) Write() error                { return nil }
func (b *failingBatch) Reset()                      {}
func (b *failingBatch) Release()                    {}
func (b *failingBatch) Replay(KeyValueWriter) error { return nil }

// failingBatchDB wraps a Database but returns a failingBatch from NewBatch.
type failingBatchDB struct {
	Database
	err error
}

func (db *failingBatchDB) NewBatch() Batch { return &failingBatch{err: db.err} }

// TestStartDBMigration_NonSingleSurfacesCopyError verifies that StartDBMigration
// returns a non-nil error when a per-DB copy fails in the non-single source
// path, instead of swallowing the error and reporting success.
func TestStartDBMigration_NonSingleSurfacesCopyError(t *testing.T) {
	newMgr := func() *databaseManager {
		return &databaseManager{
			config: &DBConfig{DBType: MemoryDB, SingleDB: false},
			dbs:    make([]Database, databaseEntryTypeSize),
			cm:     newCacheManager(),
		}
	}

	src := newMgr()
	dst := newMgr()

	// Populate one src slot with data; the matching dst slot returns batches
	// whose Put fails. All other slots remain nil and are skipped.
	srcSlot := NewMemDB()
	require.NoError(t, srcSlot.Put([]byte("k"), []byte("v")))
	src.dbs[MiscDB] = srcSlot
	dst.dbs[MiscDB] = &failingBatchDB{Database: NewMemDB(), err: errors.New("forced put failure")}

	err := src.StartDBMigration(dst)
	require.Error(t, err, "StartDBMigration must surface per-DB copy errors in the non-single path")
}
