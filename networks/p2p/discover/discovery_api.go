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

package discover

import (
	"fmt"
	"strings"

	"github.com/kaiachain/kaia/networks/rpc"
)

func (d *discovery2) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "bootnode",
			Version:   "1.0",
			Service:   NewDiscoveryAPI(d.tab),
		},
	}
}

type DiscoveryAPI struct {
	tab *Table2
}

func NewDiscoveryAPI(tab *Table2) *DiscoveryAPI {
	return &DiscoveryAPI{tab: tab}
}

// Refresh triggers a refresh of the discovery table.
func (api *DiscoveryAPI) Refresh() {
	api.tab.Refresh()
}

// Nodes returns all nodes in the table. If ntype is given (e.g. "cn", "pn", "en", "bn"),
// returns only nodes of that type. Otherwise returns all nodes.
func (api *DiscoveryAPI) Nodes(ntype *string) []*Node {
	if ntype != nil {
		nt := ParseNodeType(strings.ToLower(*ntype))
		if s := api.tab.storages[nt]; s != nil {
			return s.all()
		}
		return nil
	}
	var nodes []*Node
	for _, s := range api.tab.storages {
		nodes = append(nodes, s.all()...)
	}
	return nodes
}

// GetNode returns a node from the database by its ID.
func (api *DiscoveryAPI) GetNode(id NodeID) *Node {
	return api.tab.db.node(id)
}

// PutNode adds a node to both the table and the database.
func (api *DiscoveryAPI) PutNode(kni string) error {
	n, err := ParseNode(kni)
	if err != nil {
		return err
	}
	api.tab.addNode(n)
	return api.tab.db.updateNode(n)
}

// DeleteNode removes a node from both the table and the database.
func (api *DiscoveryAPI) DeleteNode(kni string) error {
	n, err := ParseNode(kni)
	if err != nil {
		return err
	}
	api.tab.deleteNode(n)
	return api.tab.db.deleteNode(n.ID)
}

// GetAuthorizedNodes returns the list of authorized nodes.
func (api *DiscoveryAPI) GetAuthorizedNodes() []*Node {
	return api.tab.GetAuthorizedNodes()
}

// PutAuthorizedNodes adds nodes to the authorized list. knis is a comma-separated list of KNI URIs.
func (api *DiscoveryAPI) PutAuthorizedNodes(knis string) error {
	nodes, err := parseNodeList(knis)
	if err != nil {
		return err
	}
	api.tab.PutAuthorizedNodes(nodes)
	return nil
}

// DeleteAuthorizedNodes removes nodes from the authorized list. knis is a comma-separated list of KNI URIs.
func (api *DiscoveryAPI) DeleteAuthorizedNodes(knis string) error {
	nodes, err := parseNodeList(knis)
	if err != nil {
		return err
	}
	api.tab.DeleteAuthorizedNodes(nodes)
	return nil
}

func parseNodeList(knis string) ([]*Node, error) {
	kniList := strings.Split(knis, ",")
	nodes := make([]*Node, 0, len(kniList))
	for _, kni := range kniList {
		n, err := ParseNode(strings.TrimSpace(kni))
		if err != nil {
			return nil, fmt.Errorf("invalid KNI: %s: %v", kni, err)
		}
		nodes = append(nodes, n)
	}
	return nodes, nil
}
