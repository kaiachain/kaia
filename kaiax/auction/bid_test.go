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

package auction

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/stretchr/testify/require"
)

var data = bidData{
	TargetTxHash:  common.HexToHash("0xf3c03c891206b24f5d2ff65b460df9b58c652279a3e0faed865dde4c46fe9dab"),
	BlockNumber:   11,
	Sender:        common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
	To:            common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"),
	Nonce:         0,
	Bid:           new(big.Int).SetBytes(common.Hex2Bytes("8ac7230489e80000")),
	CallGasLimit:  10000000,
	Data:          common.Hex2Bytes("d09de08a"),
	SearcherSig:   common.Hex2Bytes("2162312ceb6a69efdb73c98ee96e56d0aea1ea019184c372022ab378151112c0747066e9a9d224a822dbf31d59de492502d69d7cfc789464fa84aaac0d53f6a11b"),
	AuctioneerSig: common.Hex2Bytes("63ca36c4f6a3522b59070539453ff92011463940f98930b34a80b06a5b6b45fa136f8e79957e56e41de19cb340f2f1f7db31f964e5d5f26b1d8df13aeb2b390c1b"),
}

var testBid = &Bid{
	bidData: data,
}

func TestBidEIP712Encode(t *testing.T) {
	digest := testBid.GetHashTypedData(big.NewInt(31337), common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"))
	require.Equal(t, common.Hex2Bytes("da9b3f7a46d0b5e6875970b19ef7c60e2f969e5b44f6a4701b9889694df6fe0d"), digest)
}

func TestBidGetEthSignedMessageHash(t *testing.T) {
	digest := testBid.GetEthSignedMessageHash()
	require.Equal(t, common.Hex2Bytes("a328ed8cc9e6941076a892efd7687278bfd6f85b4b89ad196d0eef5215eb0059"), digest)
}

func TestBidValidateSearcherSig(t *testing.T) {
	err := testBid.ValidateSearcherSig(big.NewInt(31337), common.HexToAddress("0xDc64a140Aa3E981100a9becA4E685f962f0cF6C9"))
	require.NoError(t, err)
	// Do not modify the original bid.
	require.Equal(t, uint8(27), testBid.SearcherSig[crypto.RecoveryIDOffset])
}

func TestBidValidateAuctioneerSig(t *testing.T) {
	err := testBid.ValidateAuctioneerSig(common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"))
	require.NoError(t, err)
	// Do not modify the original bid.
	require.Equal(t, uint8(27), testBid.AuctioneerSig[crypto.RecoveryIDOffset])
}

func TestBidEncodeRLP(t *testing.T) {
	bid := testBid
	var buf bytes.Buffer
	err := bid.EncodeRLP(&buf)
	require.NoError(t, err)
	require.Equal(t, bid.Hash(), crypto.Keccak256Hash(buf.Bytes()))
}

func TestBidDecodeRLP(t *testing.T) {
	bid := testBid
	decoded := &Bid{}
	var buf bytes.Buffer
	err := bid.EncodeRLP(&buf)
	require.NoError(t, err)
	err = decoded.DecodeRLP(&buf)
	require.NoError(t, err)
	require.Equal(t, bid.bidData, decoded.bidData)
}

func TestBidHash(t *testing.T) {
	bid := testBid
	hash := bid.Hash()
	require.Equal(t, bid.hash.Load(), hash)
}
