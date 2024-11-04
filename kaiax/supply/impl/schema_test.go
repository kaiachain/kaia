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

package supply

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func TestSchema(t *testing.T) {
	db := database.NewMemDB()

	// When data is empty
	assert.Equal(t, uint64(0), ReadLastSupplyCheckpointNumber(db))
	assert.Nil(t, ReadSupplyCheckpoint(db, 100))

	// Test LastSupplyCheckpointNumber
	WriteLastSupplyCheckpointNumber(db, 100)
	assert.Equal(t, uint64(100), ReadLastSupplyCheckpointNumber(db))

	// Test SupplyCheckpoint
	testcases := []struct {
		Minted           *big.Int
		BurntFee         *big.Int
		expectedEncoding string
	}{
		{big.NewInt(0), big.NewInt(0), "0xc28080"},
		{big.NewInt(10000), big.NewInt(0), "0xc482271080"},
		{big.NewInt(10000), big.NewInt(20000), "0xc6822710824e20"},
	}
	for _, tc := range testcases {
		accReward := &supply.AccReward{
			Minted:   tc.Minted,
			BurntFee: tc.BurntFee,
		}
		// Check read-write round trip
		WriteSupplyCheckpoint(db, 200, accReward)
		assert.Equal(t, accReward, ReadSupplyCheckpoint(db, 200))

		// Check encoding
		b, err := db.Get(supplyCheckpointKey(200))
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedEncoding, hexutil.Encode(b))
	}
}
