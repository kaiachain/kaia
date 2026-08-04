// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from params/bootnodes.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package params

// KIP-311 defines BN as the single bootstrap role for all node types, replacing
// the pre-fork role-specific bootstrap nodes. These lists are therefore shared by
// CN, PN and EN; do not reintroduce a per-node-type split.

// MainnetBootnodes are the URLs of bootnodes running on the Kaia main network.
var MainnetBootnodes = []string{
	"kni://b286e4140aea469992146c299f8915e34d59a014c29a045de41e578014761114176f42dece34a54ce1e37f92ab981fe1bb14e54fe23576e1182778c40e6272a2@ston146.node.kaia.io:32323?ntype=bn",
	"kni://a6b61ee786952e3f8b681867d2535485622e80d0b9b7b89f26b2c631e59a4246b5f879487fbde7c324c3308ece0cdb1d7738430bdffce4f7f8c4f5a09eef80a3@ston738.node.kaia.io:32323?ntype=bn",
	"kni://b25838727eb6b4631c8f8910b4f6376fe28041f251ee21129078a61d18d62d0dc2601be5a97eab8bdb5772309f861fddb7192935483813ef20e5716a34266f16@ston903.node.kaia.io:32323?ntype=bn",
}

// KairosBootnodes are the URLs of bootnodes running on the Kairos network.
var KairosBootnodes = []string{
	"kni://979159c738bb0c8c60b36267c56d2b4d4a995326be666460c3d612856caab522ebe6f81ea5cdbb605051f12cbf8f787ce0f172256545a5b3400c751afbdcd0c8@ston36-kairos.node.kaia.io:32323?discport=32323&ntype=bn",
}
