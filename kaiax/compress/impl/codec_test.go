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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	codecs := []Codec{
		NewZstdCodec(),
	}

	for _, codec := range codecs {
		randBytes := make([]byte, 1024*1024)
		rand.Read(randBytes)

		compressed, err := codec.compress(randBytes)
		assert.NoError(t, err)
		assert.NotNil(t, compressed)

		decompressed, err := codec.decompress(compressed)
		assert.NoError(t, err)
		assert.Equal(t, randBytes, decompressed)
	}
}
