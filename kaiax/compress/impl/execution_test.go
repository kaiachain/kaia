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

	"github.com/stretchr/testify/assert"
)

func TestRewind(t *testing.T) {
	var (
		c, answers = makeTestModule(t)
		currNum    = uint64(1000)
		endNum     = uint64(872) // last block that can be compressed (curr - retention)
		rewindNum  = uint64(777)
	)

	// Start compression threads and wait for completion.
	c.Start()
	waitCompletion(t, c, endNum)

	// Rewind the chain. Note that setHeadBeyondRoot shall Stop() the module before rewinding.
	c.Stop()
	for i := currNum - 1; i > rewindNum; i-- {
		c.RewindDelete(c.DBM.ReadCanonicalHash(i), i)
	}

	for _, schema := range c.schemas {
		// 1. The nextNum must be updated.
		// The database state returns to the state just after inserting the rewindNum.
		// [-chunk-][-chunk-][-chunk-][-chunk-]------
		// |                                   |    |
		// genesis                        nextNum rewindNum = (new) currNum
		nextNum := *readNextNum(schema)
		assert.LessOrEqual(t, nextNum, rewindNum, schema.name())

		// 2. Uncompressed and compressed data must be correct.
		validateSchema(t, c, schema, nextNum, currNum, answers[schema.name()])
	}
}
