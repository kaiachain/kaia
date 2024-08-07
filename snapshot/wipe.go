// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2021 The klaytn Authors
// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state/snapshot/wipe.go (2021/10/21).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package snapshot

import (
	"bytes"
	"time"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/rcrowley/go-metrics"
)

// wipeKeyRange deletes a range of keys from the database starting with prefix
// and having a specific total key length. The start and limit is optional for
// specifying a particular key range for deletion.
//
// Origin is included for wiping and limit is excluded if they are specified.
func wipeKeyRange(db database.DBManager, kind string, prefix []byte, origin []byte, limit []byte, keylen int, meter metrics.Meter, report bool) error {
	// Batch deletions together to avoid holding an iterator for too long
	var (
		batch = db.NewSnapshotDBBatch()
		items int
	)
	defer batch.Release()
	// Iterate over the key-range and delete all of them
	start, logged := time.Now(), time.Now()

	it := db.NewSnapshotDBIterator(prefix, origin)
	var stop []byte
	if limit != nil {
		stop = append(prefix, limit...)
	}
	for it.Next() {
		// Skip any keys with the correct prefix but wrong length (trie nodes)
		key := it.Key()
		if !bytes.HasPrefix(key, prefix) {
			break
		}
		if len(key) != keylen {
			continue
		}
		if stop != nil && bytes.Compare(key, stop) >= 0 {
			break
		}
		// Delete the key and periodically recreate the batch and iterator
		batch.Delete(key)
		items++

		if items%10000 == 0 {
			// Batch too large (or iterator too long lived, flush and recreate)
			it.Release()
			if err := batch.Write(); err != nil {
				return err
			}
			batch.Reset()
			seekPos := key[len(prefix):]
			it = db.NewSnapshotDBIterator(prefix, seekPos)

			if time.Since(logged) > 8*time.Second && report {
				logger.Info("Deleting state snapshot leftovers", "kind", kind, "wiped", items, "elapsed", common.PrettyDuration(time.Since(start)))
				logged = time.Now()
			}
		}
	}
	it.Release()
	if err := batch.Write(); err != nil {
		return err
	}
	if meter != nil {
		meter.Mark(int64(items))
	}
	if report {
		logger.Info("Deleted state snapshot leftovers", "kind", kind, "wiped", items, "elapsed", common.PrettyDuration(time.Since(start)))
	}
	return nil
}
