// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from node/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package node

import (
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
)

// KaiaNodeAPI offers helper utils
type KaiaNodeAPI struct {
	stack *Node
}

// NewKaiaNodeAPI creates a new Web3Service instance
func NewKaiaNodeAPI(stack *Node) *KaiaNodeAPI {
	return &KaiaNodeAPI{stack}
}

// ClientVersion returns the node name
func (s *KaiaNodeAPI) ClientVersion() string {
	return s.stack.Server().Name()
}

// Sha3 applies the Kaia sha3 implementation on the input.
// It assumes the input is hex encoded.
func (s *KaiaNodeAPI) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
