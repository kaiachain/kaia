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
	"sort"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/kaiax/valset"
)

func (v *ValsetModule) getCouncil(num uint64) (*valset.AddressSet, error) {
	if num == 0 {
		return v.getCouncilGenesis()
	}

	// First try to get from the (migrated) DB.
	if council, ok, err := v.getCouncilDB(num); err != nil {
		return nil, err
	} else if ok {
		return council, nil
	} else {
		// Then fall back to the legacy istanbul snapshot.
		council, _, err := v.getCouncilFromIstanbulSnapshot(num, false)
		return council, err
	}
}

// getCouncilGenesis parses the genesis council from the header's extraData.
func (v *ValsetModule) getCouncilGenesis() (*valset.AddressSet, error) {
	header := v.Chain.GetHeaderByNumber(0)
	if header == nil {
		return nil, errNoHeader
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return nil, err
	}
	return valset.NewAddressSet(istanbulExtra.Validators), nil
}

func (v *ValsetModule) getCouncilDB(num uint64) (*valset.AddressSet, bool, error) {
	pMinVoteNum := ReadLowestScannedVoteNum(v.ChainKv)
	if pMinVoteNum == nil {
		return nil, false, errNoLowestScannedNum
	}
	nums := ReadValidatorVoteBlockNums(v.ChainKv)
	if nums == nil {
		return nil, false, errNoVoteBlockNums
	}

	voteNum := lastNumLessThan(nums, num)
	if voteNum < *pMinVoteNum {
		// found voteNum is not one of the scanned vote nums, i.e. the migration is not yet complete.
		// Return false to indicate that the data is not yet available.
		return nil, false, nil
	} else {
		council := valset.NewAddressSet(ReadCouncil(v.ChainKv, voteNum))
		return council, true, nil
	}
}

// lastNumLessThan returns the last (rightmost) number in the list that is less than the given number.
// If no such number exists, it returns 0.
// Suppose nums = [10, 20, 30]. If num = 25, the result is 20. If num = 7, the result is 0.
func lastNumLessThan(nums []uint64, num uint64) uint64 {
	// idx is the smallest index that is greater than or equal to `num`.
	// idx-1 is the largest index that is less than `num`.
	idx := sort.Search(len(nums), func(i int) bool {
		return nums[i] >= num
	})
	if idx > 0 && nums[idx-1] < num {
		return nums[idx-1]
	} else {
		return 0
	}
}

// getCouncilFromIstanbulSnapshot re-generates the council at the given targetNum.
// Returns the council at targetNum, the nearest snapshot number, and error if any.
//
// The council is calculated from the nearest istanbul snapshot (at snapshotNum)
// plus the validator votes in the range [snapshotNum+1, targetNum-1]. Note that
// snapshot at snapshotNum already reflects the validator vote at snapshotNum,
// so we apply the votes starting from snapshotNum+1.
//
// If write is true, ValidatorVoteBlockNums and Council in the extended range
// [snapshotNum+1, targetNum] are written to the database. Note that this time
// the targetNum is included in the range for completeness. This property is
// useful for snapshot interval-wise migration.
func (v *ValsetModule) getCouncilFromIstanbulSnapshot(targetNum uint64, write bool) (*valset.AddressSet, uint64, error) {
	if targetNum == 0 {
		council, err := v.getCouncilGenesis()
		return council, 0, err
	}
	// Load council at the nearest istanbul snapshot. This is the result
	// applying the votes up to the snapshotNum.
	snapshotNum := roundDown(targetNum-1, istanbulCheckpointInterval)
	header := v.Chain.GetHeaderByNumber(snapshotNum)
	if header == nil {
		return nil, 0, errNoHeader
	}
	council := valset.NewAddressSet(ReadIstanbulSnapshot(v.ChainKv, header.Hash()))
	if council.Len() == 0 {
		return nil, 0, ErrNoIstanbulSnapshot(snapshotNum)
	}

	// Apply the votes in the interval [snapshotNum+1, targetNum-1].
	for n := snapshotNum + 1; n < targetNum; n++ {
		if err := v.applyBlock(council, n, write); err != nil {
			return nil, 0, err
		}
	}
	// Apply the vote at targetNum to write to database, but do not return the modified council.
	if write {
		// Apply the vote at targetNum and write to database, but do not affect the returning council.
		if err := v.applyBlock(council.Copy(), targetNum, true); err != nil {
			return nil, 0, err
		}
	}
	return council, snapshotNum, nil
}

func (v *ValsetModule) applyBlock(council *valset.AddressSet, num uint64, write bool) error {
	header := v.Chain.GetHeaderByNumber(num)
	if header == nil {
		return errNoHeader
	}
	governingNode := v.GovModule.EffectiveParamSet(num).GoverningNode
	if applyVote(header, council, governingNode) && write {
		insertValidatorVoteBlockNums(v.ChainKv, num)
		writeCouncil(v.ChainKv, num, council.List())
	}
	return nil
}

// applyVote modifies the given council *in-place* by the validator vote in the given header.
// governingNode, if specified, is not affected by the vote.
// Returns true if the council is modified, false otherwise.
func applyVote(header *types.Header, council *valset.AddressSet, governingNode common.Address) bool {
	voteKey, addresses, ok := parseValidatorVote(header)
	if !ok {
		return false
	}

	originalSize := council.Len()
	for _, address := range addresses {
		if address == governingNode {
			continue
		}
		switch voteKey {
		case gov.AddValidator:
			if !council.Contains(address) {
				council.Add(address)
			}
		case gov.RemoveValidator:
			if council.Contains(address) {
				council.Remove(address)
			}
		}
	}
	return originalSize != council.Len()
}

func parseValidatorVote(header *types.Header) (gov.ValidatorParamName, []common.Address, bool) {
	// Check that a vote exists and is a validator vote.
	voteBytes := headergov.VoteBytes(header.Vote)
	if len(voteBytes) == 0 {
		return "", nil, false
	}
	vote, err := voteBytes.ToVoteData()
	if err != nil {
		return "", nil, false
	}
	voteKey := gov.ValidatorParamName(vote.Name())
	_, isValidatorVote := gov.ValidatorParams[voteKey]
	if !isValidatorVote {
		return "", nil, false
	}

	// Type cast the vote value. It can be a single address or a list of addresses.
	var addresses []common.Address
	switch voteValue := vote.Value().(type) {
	case common.Address:
		addresses = []common.Address{voteValue}
	case []common.Address:
		addresses = voteValue
	default:
		return "", nil, false
	}

	return voteKey, addresses, true
}

func roundDown(n, p uint64) uint64 {
	return n - (n % p)
}
