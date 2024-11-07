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
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
)

func (v *ValsetModule) VerifyHeader(header *types.Header) error {
	blockNum := header.Number.Uint64()
	if blockNum == 0 {
		return nil
	}

	// TODO-kaiax-valset: put the verifySigner verification in here?
	name, vote, err := newVoteDataFromBytes(header.Vote)

	// if vote.key is in gov.Params, do nothing
	if _, ok := gov.Params[gov.ParamName(name)]; ok {
		return nil
	}

	// otherwise, verify the votebyte
	if err != nil {
		return err
	}

	var (
		c      subsetCouncilSlice
		author common.Address
	)

	c, err = v.GetCouncilAddressList(blockNum)
	if err != nil {
		return err
	}
	author, err = v.chain.Engine().Author(header)
	if err != nil {
		return err
	}

	// check the proposer is the voter
	if author != vote.voter {
		return errInvalidVoter
	}
	// check the voter is in valSet
	idx := c.getIdxByAddress(vote.voter)
	if idx == -1 {
		return errInvalidVoter
	}

	return nil
}

func (v *ValsetModule) PrepareHeader(header *types.Header) error {
	if len(v.myVotes) > 0 {
		header.Vote, _ = v.myVotes[0].ToVoteBytes()
	}
	return nil
}

func (v *ValsetModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	logger.Info("NoopModule FinalizeHeader", "blockNum", header.Number.Uint64())
	return nil
}
