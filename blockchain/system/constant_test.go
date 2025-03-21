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
package system

import (
	"testing"

	"github.com/kaiachain/kaia/crypto"
	"github.com/stretchr/testify/assert"
)

func TestRuntimeCodeRegression(t *testing.T) {
	tcs := []struct {
		code []byte
		hash string
	}{
		{MainnetCreditCode, "0x24dccf9f86d49ffe0385d6fd43ceed51a71d53d9e72df9d7943a24128b4916ec"},
		{MainnetCreditV2Code, "0xb45837dfb0d4edd411a8962780361c0b984e2a25a5a03be465ae9731a5d5c0ab"},
		{RegistryCode, "0xfd3c2152828579b98068570231554ed4bacf528f50ff1bf9fce6300ec023f720"},
		{Kip113Code, "0x236841ea654b0f18e83e934ba0f69b4ab215f0b6ffbeee288797ce67c89aea25"},
		{ERC1967ProxyCode, "0x7bd49b148f3b1ffd97fb2ef2fdc773271822fa8306d3bcba626fbd412ed21c12"},
	}

	for _, tc := range tcs {
		codeHash := crypto.Keccak256Hash(tc.code)
		assert.Equal(t, tc.hash, codeHash.Hex())
	}
}
