// Copyright 2018 The klaytn Authors
// Copyright 2018 Péter Szilágyi. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found
// in the LICENSE file.
//
// This file is derived from crypto/bn256/bn256_fast.go (2018/06/04).
// Modified and improved for the klaytn development.

//go:build amd64 || arm64
// +build amd64 arm64

package bn256

import gnark "github.com/kaiachain/kaia/crypto/bn256/gnark"

// G1 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G1 = gnark.G1

// G2 is an abstract cyclic group. The zero value is suitable for use as the
// output of an operation, but cannot be used as an input.
type G2 = gnark.G2

// PairingCheck calculates the Optimal Ate pairing for a set of points.
func PairingCheck(a []*G1, b []*G2) bool {
	return gnark.PairingCheck(a, b)
}
