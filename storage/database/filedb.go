// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package database

type item struct {
	key []byte
	val []byte
}

// fileDB inserts an item, which has key and value in byte slice.
// It inserts the item to somewhere and returns the location of the item.
// An item can be retrieved with the returned location, URI.
type fileDB interface {
	write(items item) (string, error)
	read(key []byte) ([]byte, error)
	delete(key []byte) error
	deleteBucket()
}
