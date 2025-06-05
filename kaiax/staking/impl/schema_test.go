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
	"testing"

	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/kaiax/staking"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	si := &staking.StakingInfo{
		SourceBlockNum:   2880,
		NodeIds:          []common.Address{common.HexToAddress("0x159ae5ccda31b77475c64d88d4499c86f77b7ecc")},
		StakingContracts: []common.Address{common.HexToAddress("0x70e051c46ea76b9af9977407bb32192319907f9e")},
		RewardAddrs:      []common.Address{common.HexToAddress("0xd155d4277c99fa837c54a37a40a383f71a3d082a")},
		KEFAddr:          common.HexToAddress("0x673003e5f9a852d3dc85b83d16ef62d45497fb96"),
		KIFAddr:          common.HexToAddress("0x576dc0c2afeb1661da3cf53a60e76dd4e32c7ab1"),
		StakingAmounts:   []uint64{5000000},
	}

	db := database.NewMemDB()
	num := uint64(100)

	assert.Nil(t, ReadStakingInfo(db, num))
	WriteStakingInfo(db, num, si)
	assert.Equal(t, si, ReadStakingInfo(db, num))
	DeleteStakingInfo(db, num)
	assert.Nil(t, ReadStakingInfo(db, num))
}
