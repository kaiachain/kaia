// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from quorum/consensus/istanbul/utils.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package istanbul

import (
	"math"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/rlp"
)

var logger = log.NewModuleLogger(log.ConsensusIstanbul)

func RLPHash(v interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, v)
	hw.Sum(h[:0])
	return h
}

// GetSignatureAddress gets the signer address from the signature
func GetSignatureAddress(data []byte, sig []byte) (common.Address, error) {
	// 1. Keccak data
	hashData := crypto.Keccak256([]byte(data))
	// 2. Recover public key
	pubkey, err := crypto.SigToPub(hashData, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

// requiredMessageCount returns a minimum required number of consensus messages to proceed
func requiredMessageCount(qualifiedSize int, committeeSize uint64) int {
	var size int
	if qualifiedSize > int(committeeSize) {
		size = int(committeeSize)
	} else {
		size = qualifiedSize
	}
	// For less than 4 validators, quorum size equals validator count.
	if size < 4 {
		return size
	}
	// Adopted QBFT quorum implementation
	// https://github.com/Consensys/quorum/blob/master/consensus/istanbul/qbft/core/core.go#L312
	return int(math.Ceil(float64(2*size) / 3))
}

// f returns a maximum endurable number of byzantine fault nodes
func f(qualifiedSize int, committeeSize uint64) int {
	if qualifiedSize > int(committeeSize) {
		return int(math.Ceil(float64(committeeSize)/3)) - 1
	} else {
		return int(math.Ceil(float64(qualifiedSize)/3)) - 1
	}
}
