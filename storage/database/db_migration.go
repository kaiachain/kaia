// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package database

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"
)

const (
	reportCycle = IdealBatchSize * 20
)

// ErrDBMigrationInterrupted is returned by copyDB when migration is aborted
// because the quit channel was closed (e.g. by SIGINT or SIGTERM).
// Callers can use errors.Is to distinguish an interrupted migration from a
// successful completion.
var ErrDBMigrationInterrupted = errors.New("db migration interrupted")

// copyDB migrates a DB to another DB.
// This feature uses Iterator. A src DB should have implementation of Iteratee to use this function.
func copyDB(name string, srcDB, dstDB Database, quit chan struct{}) error {
	// create src iterator and dst batch
	srcIter := srcDB.NewIterator(nil, nil)
	dstBatch := dstDB.NewBatch()
	defer dstBatch.Release()

	// vars for log
	start := time.Now()
	fetched := 0

	for fetched = 0; srcIter.Next(); fetched++ {
		// fetch keys and values
		// Contents of srcIter.Key() and srcIter.Value() should not be modified, and
		// only valid until the next call to Next.
		key := make([]byte, len(srcIter.Key()))
		val := make([]byte, len(srcIter.Value()))
		copy(key, srcIter.Key())
		copy(val, srcIter.Value())

		// write fetched keys and values to DB
		// If dstDB is dynamoDB, Put will Write when the number items reach dynamoBatchSize.
		if err := dstBatch.Put(key, val); err != nil {
			return fmt.Errorf("failed to put batch: %w", err)
		}

		if dstBatch.ValueSize() > IdealBatchSize {
			if err := dstBatch.Write(); err != nil {
				return err
			}
			dstBatch.Reset()
		}

		// make a report
		if fetched%reportCycle == 0 {
			logger.Info("DB migrated",
				"db", name, "fetchedTotal", fetched, "elapsedTotal", time.Since(start))
		}

		// check for quit signal from OS
		select {
		case <-quit:
			logger.Warn("exit called", "db", name, "fetchedTotal", fetched, "elapsedTotal", time.Since(start))
			return ErrDBMigrationInterrupted
		default:
		}
	}

	if err := dstBatch.Write(); err != nil {
		return fmt.Errorf("failed to write items: %w", err)
	}
	dstBatch.Reset()

	logger.Info("Finish DB migration", "db", name, "fetchedTotal", fetched, "elapsedTotal", time.Since(start))

	srcIter.Release()
	if err := srcIter.Error(); err != nil { // any accumulated error from iterator
		return fmt.Errorf("failed to iterate: %w", err)
	}

	return nil
}

// StartDBMigration migrates a DB to another DB.
// (e.g. LevelDB -> LevelDB, LevelDB -> BadgerDB, LevelDB -> DynamoDB)
// Do not migrate db while a node is executing.
func (dbm *databaseManager) StartDBMigration(dstdbm DBManager) error {
	// settings for quit signal from os
	quit := make(chan struct{})
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logger.Info("Got interrupt, shutting down...")
		close(quit)
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				logger.Info("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()

	// from non single DB
	if !dbm.config.SingleDB {
		errChan := make(chan error, databaseEntryTypeSize)
		for et := range databaseEntryTypeSize {
			srcDB := dbm.getDatabase(et)

			dstDB := dstdbm.getDatabase(MiscDB)
			if !dstdbm.GetDBConfig().SingleDB {
				dstDB = dstdbm.getDatabase(et)
			}

			if srcDB == nil {
				logger.Warn("skip nil src db", "db", dbBaseDirs[et])
				errChan <- nil
				continue
			}

			if dstDB == nil {
				logger.Warn("skip nil dst db", "db", dbBaseDirs[et])
				errChan <- nil
				continue
			}

			dbIdx := et
			go func() {
				errChan <- copyDB(dbBaseDirs[dbIdx], srcDB, dstDB, quit)
			}()
		}

		var errs []error
		for range databaseEntryTypeSize {
			if err := <-errChan; err != nil {
				logger.Error("copyDB got an error", "err", err)
				errs = append(errs, err)
			}
		}

		// Reset state trie DB path if migrated state trie path ("statetrie_migrated_XXXXXX") is set
		dstdbm.setDBDir(DBEntryType(StateTrieDB), "")

		return errors.Join(errs...)
	}

	// single DB -> single DB
	srcDB := dbm.getDatabase(0)
	dstDB := dstdbm.getDatabase(0)

	if err := copyDB("single", srcDB, dstDB, quit); err != nil {
		return err
	}

	// If the current src DB is misc DB, clear all db dir on dst
	// TODO: If DB Migration supports non-single db, change the checking logic
	if path.Base(dbm.config.Dir) == dbBaseDirs[MiscDB] {
		for i := range uint8(databaseEntryTypeSize) {
			dstdbm.setDBDir(DBEntryType(i), "")
		}
	}

	return nil
}
