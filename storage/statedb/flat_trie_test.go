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

package statedb

import (
	"testing"

	"github.com/erigontech/erigon-lib/commitment"
	"github.com/erigontech/erigon-lib/kaiatrie"
	"github.com/erigontech/erigon-lib/state"
)

func Test_FlatTrie_Import(t *testing.T) {
	t.Log(commitment.ModeDirect)
	t.Log(len(state.Schema))
	t.Log(kaiatrie.ModeRawBytes)
}
