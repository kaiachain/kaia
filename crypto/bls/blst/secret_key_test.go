// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package blst

import (
	"testing"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/crypto/bls/types"
	"github.com/stretchr/testify/assert"
)

var (
	// https://github.com/ethereum/bls12-381-tests
	// sign/sign_case_142f678a8d05fcd1.json
	testSecretKeyBytes = common.FromHex("0x47b8192d77bf871b62e87859d653922725724a5c031afeabc60bcef5ff665138")

	// Standard test mnemonic test..junk
	testEcPrivateKeyBytes = common.FromHex("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
)

func TestGenerateKey(t *testing.T) {
	sk1, err := GenerateKey(testEcPrivateKeyBytes)
	assert.Nil(t, err)

	sk2, err := GenerateKey(testEcPrivateKeyBytes)
	assert.Nil(t, err)

	// GenerateKey is deterministic
	assert.Equal(t, sk1.Marshal(), sk2.Marshal())
}

func TestRandKey(t *testing.T) {
	sk, err := RandKey()
	assert.Nil(t, err)
	assert.Equal(t, types.SecretKeyLength, len(sk.Marshal()))
}

func TestSecretKeyFromBytes(t *testing.T) {
	b := testSecretKeyBytes
	sk, err := SecretKeyFromBytes(b)
	assert.Nil(t, err)
	assert.Equal(t, b, sk.Marshal())

	_, err = SecretKeyFromBytes([]byte{1, 2, 3, 4})
	assert.NotNil(t, err)

	// Valid secret key must be 1 <= SK < r
	// as per draft-irtf-cfrg-bls-signature-05#2.3. KeyGen
	zero := make([]byte, types.SecretKeyLength)
	_, err = SecretKeyFromBytes(zero)
	assert.Equal(t, types.ErrSecretKeyUnmarshal, err)

	order := common.FromHex(types.CurveOrderHex)
	_, err = SecretKeyFromBytes(order)
	assert.Equal(t, types.ErrSecretKeyUnmarshal, err)
}

func TestSecretKeyPublicKey(t *testing.T) {
	var (
		// https://github.com/ethereum/bls12-381-tests
		// sign/sign_case_84d45c9c7cca6b92.json
		// verify/verify_valid_case_195246ee3bd3b6ec.json
		skb = common.FromHex("0x328388aff0d4a5b7dc9205abd374e7e98f3cd9f3418edb4eafda5fb16473d216")
		pkb = common.FromHex("0xb53d21a4cfd562c469cc81514d4ce5a6b577d8403d32a394dc265dd190b47fa9f829fdd7963afdf972e5e77854051f6f")
	)

	sk, err := SecretKeyFromBytes(skb)
	assert.Nil(t, err)

	pk := sk.PublicKey()
	assert.Equal(t, pkb, pk.Marshal())
}
