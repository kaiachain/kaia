// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
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

package main

import (
	"github.com/kaiachain/kaia/networks/p2p/discover"
)

type BootnodeAPI struct {
	bn *BN
}

func NewBootnodeAPI(b *BN) *BootnodeAPI {
	return &BootnodeAPI{bn: b}
}

func (api *BootnodeAPI) GetAuthorizedNodes() []*discover.Node {
	return api.bn.GetAuthorizedNodes()
}

type BootnodeRegistryAPI struct {
	bn *BN
}

func NewBootnodeRegistryAPI(b *BN) *BootnodeRegistryAPI {
	return &BootnodeRegistryAPI{bn: b}
}

func (api *BootnodeRegistryAPI) Name() string {
	return api.bn.Name()
}

func (api *BootnodeRegistryAPI) Resolve(target discover.NodeID, targetType discover.NodeType) *discover.Node {
	return api.bn.Resolve(target, targetType)
}

func (api *BootnodeRegistryAPI) Lookup(target discover.NodeID, targetType discover.NodeType) []*discover.Node {
	return api.bn.Lookup(target, targetType)
}

func (api *BootnodeRegistryAPI) ReadRandomNodes(nType discover.NodeType) []*discover.Node {
	var buf []*discover.Node
	api.bn.ReadRandomNodes(buf, nType)
	return buf
}

func (api *BootnodeRegistryAPI) CreateUpdateNodeOnDB(nodekni string) error {
	return api.bn.CreateUpdateNodeOnDB(nodekni)
}

func (api *BootnodeRegistryAPI) CreateUpdateNodeOnTable(nodekni string) error {
	return api.bn.CreateUpdateNodeOnTable(nodekni)
}

func (api *BootnodeRegistryAPI) GetNodeFromDB(id discover.NodeID) (*discover.Node, error) {
	return api.bn.GetNodeFromDB(id)
}

func (api *BootnodeRegistryAPI) GetTableEntries() []*discover.Node {
	return api.bn.GetTableEntries()
}

func (api *BootnodeRegistryAPI) GetTableReplacements() []*discover.Node {
	return api.bn.GetTableReplacements()
}

func (api *BootnodeRegistryAPI) DeleteNodeFromDB(nodekni string) error {
	return api.bn.DeleteNodeFromDB(nodekni)
}

func (api *BootnodeRegistryAPI) DeleteNodeFromTable(nodekni string) error {
	return api.bn.DeleteNodeFromTable(nodekni)
}

func (api *BootnodeRegistryAPI) PutAuthorizedNodes(rawurl string) error {
	return api.bn.PutAuthorizedNodes(rawurl)
}

func (api *BootnodeRegistryAPI) DeleteAuthorizedNodes(rawurl string) error {
	return api.bn.DeleteAuthorizedNodes(rawurl)
}
