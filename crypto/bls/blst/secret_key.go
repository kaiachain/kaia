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
	"crypto/rand"

	"github.com/kaiachain/kaia/v2/crypto/bls/types"
	blst "github.com/supranational/blst/bindings/go"
)

type secretKey struct {
	// Pointer to underlying blst struct, hence the name 'p'
	p *blstSecretKey
}

func GenerateKey(ikm []byte) (types.SecretKey, error) {
	// draft-irtf-cfrg-bls-signature-05 section 2.3. KeyGen
	// requires that IKM MUST be at least 32 bytes long, but it MAY be longer.
	if len(ikm) < 32 {
		return nil, types.ErrSecretKeyGen
	}

	p := blst.KeyGen(ikm)
	if p == nil || !p.Valid() {
		return nil, types.ErrSecretKeyGen
	}
	return &secretKey{p: p}, nil
}

func RandKey() (types.SecretKey, error) {
	ikm := make([]byte, 32)
	if _, err := rand.Read(ikm); err != nil {
		return nil, err
	}
	return GenerateKey(ikm)
}

func SecretKeyFromBytes(b []byte) (types.SecretKey, error) {
	if len(b) != types.SecretKeyLength {
		return nil, types.ErrSecretKeyLength(len(b))
	}

	p := new(blstSecretKey).Deserialize(b)
	if p == nil || !p.Valid() {
		return nil, types.ErrSecretKeyUnmarshal
	}
	return &secretKey{p: p}, nil
}

func (sk *secretKey) PublicKey() types.PublicKey {
	// must succeed because SecretKey always hold a valid scalar,
	p := new(blstPublicKey).From(sk.p) // blst_sk_to_pk2_in_g2
	return &publicKey{p: p}
}

func (sk *secretKey) Marshal() []byte {
	return sk.p.Serialize() // blst_p1_affine_serialize
}
