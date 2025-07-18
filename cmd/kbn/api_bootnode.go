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
