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

package forkid

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// TestCreation tests that different genesis and fork rule combinations result in
// the correct fork ID.
func TestCreation(t *testing.T) {
	type testcase struct {
		head uint64
		want ID
	}
	tests := []struct {
		config  *params.ChainConfig
		genesis common.Hash
		cases   []testcase
	}{
		// Mainnet test cases
		{
			params.MainnetChainConfig,
			params.MainnetGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0xdd58eb45), Next: 86816005}},          // Unsynced
				{86816004, ID{Hash: checksumToBytes(0xdd58eb45), Next: 86816005}},   // Last Genesis block
				{86816005, ID{Hash: checksumToBytes(0x7b1131bb), Next: 99841497}},   // First Istanbul+London+EthTxType block
				{99841496, ID{Hash: checksumToBytes(0x7b1131bb), Next: 99841497}},   // Last Istanbul+London+EthTxType block
				{99841497, ID{Hash: checksumToBytes(0x8b3961e6), Next: 119750400}},  // First Magma block
				{119750399, ID{Hash: checksumToBytes(0x8b3961e6), Next: 119750400}}, // Last Magma block
				{119750400, ID{Hash: checksumToBytes(0x171d8904), Next: 135456000}}, // First Kore+Kip103 block
				{135455999, ID{Hash: checksumToBytes(0x171d8904), Next: 135456000}}, // Last Kore+Kip103 block
				{135456000, ID{Hash: checksumToBytes(0x68717c3d), Next: 147534000}}, // First Shanghai block
				{147533999, ID{Hash: checksumToBytes(0x68717c3d), Next: 147534000}}, // Last Shanghai block
				{147534000, ID{Hash: checksumToBytes(0x75771543), Next: 162900480}}, // First Cancun+Randao block
				{162900479, ID{Hash: checksumToBytes(0x75771543), Next: 162900480}}, // Last Cancun+Randao block
				{162900480, ID{Hash: checksumToBytes(0x3ab6dcda), Next: 190670000}}, // First Kaia+Kip160 block
				{190669999, ID{Hash: checksumToBytes(0x3ab6dcda), Next: 190670000}}, // Last Kaia+Kip160 block
				{190670000, ID{Hash: checksumToBytes(0xc00bab0e), Next: 0}},         // First Prague block
				{198489578, ID{Hash: checksumToBytes(0xc00bab0e), Next: 0}},         // Today Prague block
			},
		},
		// Kairos test cases
		{
			params.KairosChainConfig,
			params.KairosGenesisHash,
			[]testcase{
				{0, ID{Hash: checksumToBytes(0x5e10f192), Next: 75373312}},          // Unsynced
				{9, ID{Hash: checksumToBytes(0x5e10f192), Next: 75373312}},          // Last Genesis block
				{75373312, ID{Hash: checksumToBytes(0x3640372e), Next: 80295291}},   // First Istanbul block
				{80295290, ID{Hash: checksumToBytes(0x3640372e), Next: 80295291}},   // Last Istanbul block
				{80295291, ID{Hash: checksumToBytes(0xb1d92301), Next: 86513895}},   // First London block
				{86513894, ID{Hash: checksumToBytes(0xb1d92301), Next: 86513895}},   // Last London block
				{86513895, ID{Hash: checksumToBytes(0xf2be985d), Next: 98347376}},   // First EthTxType block
				{98347375, ID{Hash: checksumToBytes(0xf2be985d), Next: 98347376}},   // Last EthTxType block
				{98347376, ID{Hash: checksumToBytes(0xf99f3ae7), Next: 111736800}},  // First Magma block
				{111736799, ID{Hash: checksumToBytes(0xf99f3ae7), Next: 111736800}}, // Last Magma block
				{111736800, ID{Hash: checksumToBytes(0x40c759b1), Next: 119145600}}, // First Kore block
				{119145599, ID{Hash: checksumToBytes(0x40c759b1), Next: 119145600}}, // Last Kore block
				{119145600, ID{Hash: checksumToBytes(0x989f8b4a), Next: 131608000}}, // First Kip103 block
				{131607999, ID{Hash: checksumToBytes(0x989f8b4a), Next: 131608000}}, // Last Kip103 block
				{131608000, ID{Hash: checksumToBytes(0xff2bb36d), Next: 141367000}}, // First Shanghai block
				{141366999, ID{Hash: checksumToBytes(0xff2bb36d), Next: 141367000}}, // Last Shanghai block
				{141367000, ID{Hash: checksumToBytes(0x897592ea), Next: 156660000}}, // First Cancun+Randao block
				{156659999, ID{Hash: checksumToBytes(0x897592ea), Next: 156660000}}, // Last Cancun+Randao block
				{156660000, ID{Hash: checksumToBytes(0xf00fbab3), Next: 187930000}}, // First Kaia+Kip160 block
				{187929999, ID{Hash: checksumToBytes(0xf00fbab3), Next: 187930000}}, // Last Kaia+Kip160 block
				{187930000, ID{Hash: checksumToBytes(0x0d29cd98), Next: 0}},         // First Prague block
				{198978293, ID{Hash: checksumToBytes(0x0d29cd98), Next: 0}},         // Last Prague block
			},
		},
	}
	for i, tt := range tests {
		for j, ttt := range tt.cases {
			if have := NewID(tt.config, tt.genesis, ttt.head); have != ttt.want {
				t.Errorf("test %d, case %d: fork ID mismatch: have %x, want %x", i, j, have, ttt.want)
			}
		}
	}
}

func TestCheckForkCompatibleBlock(t *testing.T) {
	type expectedNums struct {
		currentForkCompatibleBlock *big.Int
		nextForkCompatibleBlock    *big.Int
		lastForkCompatibleBlock    *big.Int
	}
	type testcase struct {
		head         uint64
		expectedNums expectedNums
	}
	tests := []struct {
		config *params.ChainConfig
		cases  []testcase
	}{
		// Mainnet test cases
		{
			params.MainnetChainConfig,
			[]testcase{
				{0, expectedNums{big.NewInt(0), big.NewInt(86816005), big.NewInt(190670000)}},                  // head is Genesis block
				{86816005, expectedNums{big.NewInt(86816005), big.NewInt(99841497), big.NewInt(190670000)}},    // head is Istanbul+London+EthTxType block
				{99841497, expectedNums{big.NewInt(99841497), big.NewInt(119750400), big.NewInt(190670000)}},   // head is Magma block
				{119750400, expectedNums{big.NewInt(119750400), big.NewInt(135456000), big.NewInt(190670000)}}, // head is Kore+Kip103 block
				{135456000, expectedNums{big.NewInt(135456000), big.NewInt(147534000), big.NewInt(190670000)}}, // head is Shanghai block
				{147534000, expectedNums{big.NewInt(147534000), big.NewInt(162900480), big.NewInt(190670000)}}, // head is Cancun+Randao block
				{162900480, expectedNums{big.NewInt(162900480), big.NewInt(190670000), big.NewInt(190670000)}}, // head is Kaia+Kip160 block
				{190670000, expectedNums{big.NewInt(190670000), nil, big.NewInt(190670000)}},                   // head is Prague block
			},
		},
		// Kairos test cases
		{
			params.KairosChainConfig,
			[]testcase{
				{0, expectedNums{big.NewInt(0), big.NewInt(75373312), big.NewInt(187930000)}},                  // head is Genesis block
				{75373312, expectedNums{big.NewInt(75373312), big.NewInt(80295291), big.NewInt(187930000)}},    // head is Istanbul block
				{80295291, expectedNums{big.NewInt(80295291), big.NewInt(86513895), big.NewInt(187930000)}},    // head is London block
				{86513895, expectedNums{big.NewInt(86513895), big.NewInt(98347376), big.NewInt(187930000)}},    // head is EthTxType block
				{98347376, expectedNums{big.NewInt(98347376), big.NewInt(111736800), big.NewInt(187930000)}},   // head is Magma block
				{111736800, expectedNums{big.NewInt(111736800), big.NewInt(119145600), big.NewInt(187930000)}}, // head is Kore block
				{119145600, expectedNums{big.NewInt(119145600), big.NewInt(131608000), big.NewInt(187930000)}}, // head is Kip103 block
				{131608000, expectedNums{big.NewInt(131608000), big.NewInt(141367000), big.NewInt(187930000)}}, // head is Shanghai block
				{141367000, expectedNums{big.NewInt(141367000), big.NewInt(156660000), big.NewInt(187930000)}}, // head is Cancun+Randao block
				{156660000, expectedNums{big.NewInt(156660000), big.NewInt(187930000), big.NewInt(187930000)}}, // head is Kaia+Kip160 block
				{187930000, expectedNums{big.NewInt(187930000), nil, big.NewInt(187930000)}},                   // head is Prague block
			},
		},
	}
	for i, tt := range tests {
		for j, ttt := range tt.cases {
			latestForkCompatibleBlock := LatestForkCompatibleBlock(tt.config, new(big.Int).SetUint64(ttt.head))
			nextForkCompatibleBlock := NextForkCompatibleBlock(tt.config, new(big.Int).SetUint64(ttt.head))
			lastForkCompatibleBlock := LastForkCompatibleBlock(tt.config)
			if latestForkCompatibleBlock.Cmp(ttt.expectedNums.currentForkCompatibleBlock) != 0 {
				t.Errorf("test %d, case %d: latest fork compatible block mismatch: have %x, want %x", i, j, latestForkCompatibleBlock, ttt.expectedNums.currentForkCompatibleBlock)
			}
			if nextForkCompatibleBlock == nil {
				if ttt.expectedNums.nextForkCompatibleBlock != nil {
					t.Errorf("test %d, case %d: next fork compatible block is nil, but expected %x", i, j, ttt.expectedNums.nextForkCompatibleBlock)
				}
			} else if nextForkCompatibleBlock.Cmp(ttt.expectedNums.nextForkCompatibleBlock) != 0 {
				t.Errorf("test %d, case %d: next fork compatible block mismatch: have %x, want %x", i, j, nextForkCompatibleBlock, ttt.expectedNums.nextForkCompatibleBlock)
			}
			if lastForkCompatibleBlock.Cmp(ttt.expectedNums.lastForkCompatibleBlock) != 0 {
				t.Errorf("test %d, case %d: last fork compatible block mismatch: have %x, want %x", i, j, lastForkCompatibleBlock, ttt.expectedNums.lastForkCompatibleBlock)
			}
		}
	}
}
