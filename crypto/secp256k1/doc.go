// Copyright 2018 The klaytn Authors
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
// This file is derived from crypto/secp256k1/secp256.go (2018/06/04).
// Modified and improved for the klaytn development.

/*
Package secp256k1 wraps the github.com/erigontech/secp256k1 library.

A CGO implementation used to be in this directory. But was deleted and replaced by a simple wrapper to erigontech/secp256k1 package.
When we imported the erigontech/erigon-lib, the kaiachain/kaia/crypto/secp256k1 and erigontech/secp256k1 caused CGO symbol collision.
The resolution was using erigontech/secp256k1 only.
*/
package secp256k1
