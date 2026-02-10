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

package vrank

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/kaiax"
)

type VRankModule interface {
	kaiax.BaseModule

	HandleIstanbulPreprepare(block *types.Block, view *istanbul.View)
	HandleVRankPreprepare(msg *VRankPreprepare) error
	HandleVRankCandidate(msg *VRankCandidate) error
	GetCfReport(blockNum, round uint64) (CfReport, error)
}

type VRankModuleHost interface {
	RegisterVRankModule(module VRankModule)
}
