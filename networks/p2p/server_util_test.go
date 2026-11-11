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

package p2p

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/stretchr/testify/assert"
)

func TestEffectiveConnTypeTreatsPNAsEN(t *testing.T) {
	assert.Equal(t, common.ENDPOINTNODE, EffectiveConnType(common.PROXYNODE))
	assert.Equal(t, common.ENDPOINTNODE, EffectiveConnType(common.ENDPOINTNODE))
	assert.Equal(t, common.CONSENSUSNODE, EffectiveConnType(common.CONSENSUSNODE))
	assert.Equal(t, common.BOOTNODE, EffectiveConnType(common.BOOTNODE))
}

func TestEffectiveNodeTypeTreatsPNAsEN(t *testing.T) {
	assert.Equal(t, discover.NodeTypeEN, discover.EffectiveNodeType(discover.NodeTypePN))
	assert.Equal(t, discover.NodeTypeEN, discover.EffectiveNodeType(discover.NodeTypeEN))
	assert.Equal(t, discover.NodeTypeCN, discover.EffectiveNodeType(discover.NodeTypeCN))
	assert.Equal(t, discover.NodeTypeBN, discover.EffectiveNodeType(discover.NodeTypeBN))
}
