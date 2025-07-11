// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"fmt"

	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/p2p"
)

// NetAPI offers network related RPC methods
type NetAPI struct {
	net            p2p.Server
	networkVersion uint64
}

// NewNetAPI creates a new net API instance.
func NewNetAPI(net p2p.Server, networkVersion uint64) *NetAPI {
	return &NetAPI{net, networkVersion}
}

// Listening returns an indication if the node is listening for network connections.
func (s *NetAPI) Listening() bool {
	return true // always listening
}

// PeerCount returns the number of connected peers.
func (s *NetAPI) PeerCount() hexutil.Uint {
	if s.net == nil {
		return 0
	}
	return hexutil.Uint(s.net.PeerCount())
}

// PeerCountByType returns the number of connected specific types of nodes.
func (s *NetAPI) PeerCountByType() map[string]uint {
	if s.net == nil {
		return make(map[string]uint)
	}
	return s.net.PeerCountByType()
}

// Version returns the current Kaia protocol version.
func (s *NetAPI) Version() string {
	return fmt.Sprintf("%d", s.networkVersion)
}

// NetworkID returns the network identifier set by the command-line option --networkid.
func (s *NetAPI) NetworkID() uint64 {
	return s.networkVersion
}
