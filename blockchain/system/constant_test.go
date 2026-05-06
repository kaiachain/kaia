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

// TestRuntimeCodeRegression pins the keccak256 hash of every system-contract
// runtime bytecode this codebase deploys at a well-known address. The
// bytecode is part of the canonical chain state once installed; changing it
// flips the post-install state root and forks the chain. Consensus surface.
//
// DO NOT MODIFY THE EXPECTED HASHES BELOW.
//
// If a regenerated abigen output makes this test fail, investigate why
// (solc version, optimizer, source change) and either revert the
// regeneration or route the new bytecode through a new constant — the way
// ERC1967ProxyV5Code was added alongside the existing ERC1967ProxyCode.
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
		{ERC1967ProxyV5Code, "0x3426e8c58f22c64051f94b923a1ceff79723296d2d5e578be030b91d24d6eb2a"},
		{UniswapV2FactoryCode, "0xbab145d02e7005f0d84c6c1639d39b799b0ea16df99ebbdaf5a14d9da820b4e0"},
		{UniswapV2Router02Code, "0x8078c0090b05e0bee0587064947604e217146cc295dcb119a2c0217d6e88dac5"},
	}

	for _, tc := range tcs {
		codeHash := crypto.Keccak256Hash(tc.code)
		assert.Equal(t, tc.hash, codeHash.Hex())
	}
}
