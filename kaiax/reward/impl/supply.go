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
	"errors"
	"math/big"

	"github.com/kaiachain/kaia/kaiax/reward"
)

// totalSupplyFromState traverses the state trie to sum up the total supply.
// It is extremely inefficient so it should only be used for genesis block and testing.
func (r *RewardModule) totalSupplyFromState(num uint64) (*big.Int, error) {
	header := r.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, reward.ErrNoBlock
	}
	stateDB, err := r.Chain.StateAt(header.Root)
	if err != nil {
		return nil, err
	}
	dump := stateDB.RawDump()

	totalSupply := new(big.Int)
	for _, account := range dump.Accounts {
		balance, ok := new(big.Int).SetString(account.Balance, 10)
		if !ok {
			return nil, errors.New("malformed state dump")
		}
		totalSupply.Add(totalSupply, balance)
	}
	return totalSupply, nil
}
