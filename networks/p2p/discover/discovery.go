// Modifications Copyright 2026 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
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

package discover

type Discovery2 interface {
	RandomNodes(buf []*Node, targetType NodeType) int
	ClosestNodes(targetID NodeID, targetType NodeType, max int) []*Node
	Close()
}

type discovery struct {
	udp *udp
	tab *table
}

func newDiscovery(cfg *Config) (*discovery, error) {
	_, udp, err := newUDP(cfg)
	if err != nil {
		return nil, err
	}
	tab, err := newTable2(cfg, udp)
	if err != nil {
		return nil, err
	}

	tab.Start()
	udp.Start(tab, cfg.Unhandled)

	return &discovery{
		udp: udp,
		tab: tab,
	}, nil
}

func (d *discovery) Close() {
	d.udp.close()
	d.tab.Close()
}
