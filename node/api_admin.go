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
	"encoding/hex"

	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/networks/p2p"
)

// AdminNodeAPI is the collection of administrative API methods exposed over
// both secure and unsecure RPC channels.
type AdminNodeAPI struct {
	node *Node // Node interfaced by this API
}

// NewAdminNodeAPI creates a new API definition for the public admin methods
// of the node itself.
func NewAdminNodeAPI(node *Node) *AdminNodeAPI {
	return &AdminNodeAPI{node: node}
}

// Peers retrieves all the information we know about each individual peer at the
// protocol granularity.
func (api *AdminNodeAPI) Peers() ([]*p2p.PeerInfo, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}
	return server.PeersInfo(), nil
}

// BlsPublicKeyInfoOutput has string fields unlike system.BlsPublicKeyInfo.
type BlsPublicKeyInfoOutput struct {
	PublicKey string `json:"publicKey"`
	Pop       string `json:"pop"`
}

// NodeInfoOutput extends p2p.NodeInfo with additional fields.
type NodeInfoOutput struct {
	p2p.NodeInfo

	// Node address derived from node key
	NodeAddress string `json:"nodeAddress"`

	// BLS public key information for consensus
	BlsPublicKeyInfo BlsPublicKeyInfoOutput `json:"blsPublicKeyInfo"`
}

// NodeInfo retrieves all the information we know about the host node at the
// protocol granularity.
func (api *AdminNodeAPI) NodeInfo() (*NodeInfoOutput, error) {
	server := api.node.Server()
	if server == nil {
		return nil, ErrNodeStopped
	}

	var (
		nodeKey  = api.node.config.NodeKey()
		nodeAddr = crypto.PubkeyToAddress(nodeKey.PublicKey)

		blsPriv = api.node.config.BlsNodeKey()
		blsPub  = blsPriv.PublicKey().Marshal()
		blsPop  = bls.PopProve(blsPriv).Marshal()
	)

	info := &NodeInfoOutput{
		NodeInfo:    *server.NodeInfo(),
		NodeAddress: nodeAddr.Hex(),
		BlsPublicKeyInfo: BlsPublicKeyInfoOutput{
			PublicKey: hex.EncodeToString(blsPub),
			Pop:       hex.EncodeToString(blsPop),
		},
	}
	return info, nil
}

// Datadir retrieves the current data directory the node is using.
func (api *AdminNodeAPI) Datadir() string {
	return api.node.DataDir()
}
