// Copyright 2025 The Kaia Authors
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
	"context"

	"github.com/erigontech/erigon-lib/kv"
)

func (k *FlatKVModule) Get(key []byte) ([]byte, error) {
	var ret []byte
	err := k.chaindb.View(context.Background(), func(tx kv.Tx) error {
		var err error
		ret, err = tx.GetOne("Headers", key)
		return err
	})
	return ret, err
}

func (k *FlatKVModule) Put(key []byte, value []byte) error {
	return k.chaindb.Update(context.Background(), func(tx kv.RwTx) error {
		return tx.Put("Headers", key, value)
	})
}

func (k *FlatKVModule) Delete(key []byte) error {
	return k.chaindb.Update(context.Background(), func(tx kv.RwTx) error {
		return tx.Delete("Headers", key)
	})
}
