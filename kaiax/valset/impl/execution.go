// Copyright 2024 The Kaia Authors
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

package impl

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

func (v *ValsetModule) PostInsertBlock(block *types.Block) error {
	header := block.Header()
	num := header.Number.Uint64()

	// Ingest validator vote
	council, err := v.getCouncil(num)
	if err != nil {
		return err
	}
	governingNode := v.GovModule.EffectiveParamSet(num).GoverningNode
	if applyVote(header, council, governingNode) {
		insertValidatorVoteBlockNums(v.ChainKv, num)
		writeCouncil(v.ChainKv, num, council.List())
	}

	return nil
}

func (v *ValsetModule) RewindTo(block *types.Block) {
	trimValidatorVoteBlockNums(v.ChainKv, block.Header().Number.Uint64())
}

func (v *ValsetModule) RewindDelete(hash common.Hash, num uint64) {
	deleteCouncil(v.ChainKv, num)
}
