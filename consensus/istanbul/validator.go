// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
//
// This file is derived from quorum/consensus/istanbul/validator.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package istanbul

import (
	"math"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/valset"
)

type BlockValSet struct {
	council   *valset.AddressSet // council = demoted + qualified
	qualified *valset.AddressSet
	demoted   *valset.AddressSet
}

func NewBlockValSet(council, demoted []common.Address) *BlockValSet {
	councilSet := valset.NewAddressSet(council)
	demotedSet := valset.NewAddressSet(demoted)
	qualifiedSet := councilSet.Subtract(demotedSet)

	return &BlockValSet{councilSet, qualifiedSet, demotedSet}
}
func (cs *BlockValSet) Council() *valset.AddressSet   { return cs.council }
func (cs *BlockValSet) Qualified() *valset.AddressSet { return cs.qualified }
func (cs *BlockValSet) Demoted() *valset.AddressSet   { return cs.demoted }
func (cs *BlockValSet) CheckValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	// 1. Get signature address
	signer, err := GetSignatureAddress(data, sig)
	if err != nil {
		logger.Error("Failed to get signer address", "err", err)
		return common.Address{}, err
	}

	// 2. Check validator
	if cs.Qualified().Contains(signer) {
		return signer, nil
	}

	return common.Address{}, ErrUnauthorizedAddress
}

type RoundCommitteeState struct {
	*BlockValSet
	committee *valset.AddressSet
	proposer  common.Address

	// pre-calculated values
	committeeSize        uint64
	requiredMessageCount int
	f                    int
}

func NewRoundCommitteeState(set *BlockValSet, committeeSize uint64, committee []common.Address, proposer common.Address) *RoundCommitteeState {
	committeeSet := valset.NewAddressSet(committee)
	reqMsgCnt := requiredMessageCount(set.Qualified().Len(), committeeSize)
	fNum := f(set.Qualified().Len(), committeeSize)

	return &RoundCommitteeState{set, committeeSet, proposer, committeeSize, reqMsgCnt, fNum}
}
func (cs *RoundCommitteeState) ValSet() *BlockValSet          { return cs.BlockValSet }
func (cs *RoundCommitteeState) Committee() *valset.AddressSet { return cs.committee }
func (cs *RoundCommitteeState) NonCommittee() *valset.AddressSet {
	return cs.qualified.Subtract(cs.committee)
}
func (cs *RoundCommitteeState) Proposer() common.Address            { return cs.proposer }
func (cs *RoundCommitteeState) IsProposer(addr common.Address) bool { return cs.proposer == addr }
func (cs *RoundCommitteeState) CommitteeSize() uint64               { return cs.committeeSize }
func (cs *RoundCommitteeState) RequiredMessageCount() int           { return cs.requiredMessageCount }
func (cs *RoundCommitteeState) F() int                              { return cs.f }

// requiredMessageCount returns a minimum required number of consensus messages to proceed
func requiredMessageCount(qualifiedSize int, committeeSize uint64) int {
	var size int
	if qualifiedSize > int(committeeSize) {
		size = int(committeeSize)
	} else {
		size = qualifiedSize
	}
	// For less than 4 validators, quorum size equals validator count.
	if size < 4 {
		return size
	}
	// Adopted QBFT quorum implementation
	// https://github.com/Consensys/quorum/blob/master/consensus/istanbul/qbft/core/core.go#L312
	return int(math.Ceil(float64(2*size) / 3))
}

// f returns a maximum endurable number of byzantine fault nodes
func f(qualifiedSize int, committeeSize uint64) int {
	if qualifiedSize > int(committeeSize) {
		return int(math.Ceil(float64(committeeSize)/3)) - 1
	} else {
		return int(math.Ceil(float64(qualifiedSize)/3)) - 1
	}
}
