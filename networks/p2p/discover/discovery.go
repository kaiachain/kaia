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

type Discovery2 interface {
	// Lifecycle. Used by p2p.Server.
	Close()

	// Trigger a refresh of the discovery table and wait for completion.
	Refresh()

	// Serve discovered nodes. Used by p2p.DialSched.
	RandomNodes(buf []*Node, nType NodeType) int

	// Update authorized nodes. To be used by permless.
	PutAuthorizedNodes(nodes []*Node)
}

func NewDiscovery2(cfg *Config) (Discovery2, error) {
	return nil, nil
}
