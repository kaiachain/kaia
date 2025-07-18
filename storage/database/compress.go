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

package database

import "github.com/kaiachain/kaia/common"

type compressModule interface {
	FindFromCompressedHeader(num uint64, hash common.Hash) ([]byte, bool)
	FindFromCompressedBody(num uint64, hash common.Hash) ([]byte, bool)
	FindFromCompressedReceipts(num uint64, hash common.Hash) ([]byte, bool)
}

func (dbm *databaseManager) RegisterCompressModule(module compressModule) {
	dbm.compressModule = module
}

func (dbm *databaseManager) hasFromCompressedHeader(hash common.Hash, number uint64) bool {
	if dbm.compressModule != nil {
		if _, ok := dbm.compressModule.FindFromCompressedHeader(number, hash); ok {
			return true
		}
	}
	return false
}

func (dbm *databaseManager) hasFromCompressedBody(hash common.Hash, number uint64) bool {
	if dbm.compressModule != nil {
		if _, ok := dbm.compressModule.FindFromCompressedBody(number, hash); ok {
			return true
		}
	}
	return false
}

func (dbm *databaseManager) getFromCompressedHeader(hash common.Hash, number uint64) ([]byte, bool) {
	if dbm.compressModule != nil {
		return dbm.compressModule.FindFromCompressedHeader(number, hash)
	}
	return nil, false
}

func (dbm *databaseManager) getFromCompressedBody(hash common.Hash, number uint64) ([]byte, bool) {
	if dbm.compressModule != nil {
		return dbm.compressModule.FindFromCompressedBody(number, hash)
	}
	return nil, false
}

func (dbm *databaseManager) getFromCompressedReceipts(hash common.Hash, number uint64) ([]byte, bool) {
	if dbm.compressModule != nil {
		return dbm.compressModule.FindFromCompressedReceipts(number, hash)
	}
	return nil, false
}
