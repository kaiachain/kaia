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
		{RegistryCode, "0x81e4d72a5f324997e38f750704bd64dcce9b2c4901f843b3c35457e178e904b8"},
		{Kip113Code, "0x3454be181730d70863b44c6e6a4089808908dd497e50d6c425777b1b8566700c"},
		{ERC1967ProxyCode, "0x3426e8c58f22c64051f94b923a1ceff79723296d2d5e578be030b91d24d6eb2a"},
		{UniswapV2FactoryCode, "0xf81ae5cf2f10963f6d768077d89e938fcfbaa4c14872776babcf0dee4328810a"},
		{UniswapV2Router02Code, "0x8078c0090b05e0bee0587064947604e217146cc295dcb119a2c0217d6e88dac5"},
	}

	for _, tc := range tcs {
		codeHash := crypto.Keccak256Hash(tc.code)
		assert.Equal(t, tc.hash, codeHash.Hex())
	}
}
