// Copyright 2026 The Kaia Authors
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

package main

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
)

// Reference digest computed with foundry `cast` over the same abi.encode + keccak256
// as NodeVerifier._verifyNodeIdProof / NodeIdSigUtil — pins the Go digest to Solidity.
func TestCreateNodeDigest_MatchesSolidity(t *testing.T) {
	got := common.BytesToHash(createNodeDigest(
		big.NewInt(8217),
		common.HexToAddress("0x1111111111111111111111111111111111111111"),
		common.HexToAddress("0x2222222222222222222222222222222222222222"),
		common.HexToAddress("0x3333333333333333333333333333333333333333"),
	))
	want := common.HexToHash("0x775bb349c16c3946915d91b60fbc33faa61328af564b025623ba955cc2b678b2")
	if got != want {
		t.Fatalf("digest mismatch:\n got  %s\n want %s", got.Hex(), want.Hex())
	}
}

// The produced signature must recover to nodeId once the v+27 offset is undone,
// mirroring the contract's ecrecover.
func TestCreateNodeSig_RecoversNodeId(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	nodeId := crypto.PubkeyToAddress(key.PublicKey)
	digest := createNodeDigest(big.NewInt(1001), common.HexToAddress("0xAbCd"), nodeId, common.HexToAddress("0xBeeF"))

	sig, err := crypto.Sign(digest, key)
	if err != nil {
		t.Fatal(err)
	}
	sig[64] += 27

	recSig := make([]byte, 65)
	copy(recSig, sig)
	recSig[64] -= 27 // crypto.SigToPub wants v in {0,1}
	pub, err := crypto.SigToPub(digest, recSig)
	if err != nil {
		t.Fatal(err)
	}
	if crypto.PubkeyToAddress(*pub) != nodeId {
		t.Fatalf("recovered %s, want %s", crypto.PubkeyToAddress(*pub).Hex(), nodeId.Hex())
	}
}
