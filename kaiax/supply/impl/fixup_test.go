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

package supply

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Helpers

// newFixupModule creates a bare SupplyModule with only the fields needed by
// applyBlobFeeFixupIfNeeded: ChainKv and ChainConfig.ChainID.
func newFixupModule(db database.Database, chainID *big.Int) *SupplyModule {
	s := NewSupplyModule()
	s.ChainKv = db
	s.ChainConfig = &params.ChainConfig{ChainID: chainID}
	return s
}

// writeCheckpoints writes identical fake AccReward entries at every 128-block
// boundary from start to end (inclusive) and sets lastAccRewardNumber to end.
func writeCheckpoints(db database.Database, start, end uint64) {
	fake := &supply.AccReward{
		TotalMinted: big.NewInt(1000),
		BurntFee:    big.NewInt(50), // intentionally low — simulates missing blob fee
	}
	for num := start; num <= end; num += checkpointInterval {
		WriteAccReward(db, num, fake)
	}
	WriteLastAccRewardNumber(db, end)
}

// ----------------------------------------------------------------------------
// Tests

func TestApplyBlobFeeFixupIfNeeded(t *testing.T) {
	// Kairos  (chainID=1001): firstBadCheckpoint=209134080, safeReset=209133952
	// Mainnet (chainID=8217): firstBadCheckpoint=213333120, safeReset=213332992

	t.Run("kairos: fixup runs and resets pointer", func(t *testing.T) {
		db := database.NewMemDB()
		// Write a pre-Osaka checkpoint (safe) and two bad ones.
		writeCheckpoints(db, 209133952, 209134208) // 209133952, 209134080, 209134208

		s := newFixupModule(db, big.NewInt(int64(params.KairosNetworkId)))
		s.applyBlobFeeFixupIfNeeded()

		assert.Equal(t, uint64(209133952), ReadLastAccRewardNumber(db), "pointer must be reset to safeReset")
		assert.NotNil(t, ReadAccReward(db, 209133952), "pre-Osaka checkpoint must be preserved")
		assert.Nil(t, ReadAccReward(db, 209134080), "first bad checkpoint must be deleted")
		assert.Nil(t, ReadAccReward(db, 209134208), "subsequent bad checkpoint must be deleted")
		assert.True(t, ReadSupplyBlobFeeFixupDone(db), "done flag must be set")
	})

	t.Run("mainnet: fixup runs and resets pointer", func(t *testing.T) {
		db := database.NewMemDB()
		writeCheckpoints(db, 213332992, 213333248) // 213332992, 213333120, 213333248

		s := newFixupModule(db, big.NewInt(int64(params.MainnetNetworkId)))
		s.applyBlobFeeFixupIfNeeded()

		assert.Equal(t, uint64(213332992), ReadLastAccRewardNumber(db))
		assert.NotNil(t, ReadAccReward(db, 213332992))
		assert.Nil(t, ReadAccReward(db, 213333120))
		assert.Nil(t, ReadAccReward(db, 213333248))
		assert.True(t, ReadSupplyBlobFeeFixupDone(db))
	})

	t.Run("before bad range: no-op, done flag not set", func(t *testing.T) {
		// Simulates Mainnet before Osaka activation: lastAccNum < firstBadCheckpoint.
		db := database.NewMemDB()
		writeCheckpoints(db, 213332992, 213332992) // only one checkpoint, before the bad range
		require.Equal(t, uint64(213332992), ReadLastAccRewardNumber(db))

		s := newFixupModule(db, big.NewInt(int64(params.MainnetNetworkId)))
		s.applyBlobFeeFixupIfNeeded()

		// Pointer and checkpoint must be untouched.
		assert.Equal(t, uint64(213332992), ReadLastAccRewardNumber(db))
		assert.NotNil(t, ReadAccReward(db, 213332992))
		// Done flag must NOT be set so the check re-runs after Osaka activates.
		assert.False(t, ReadSupplyBlobFeeFixupDone(db))
	})

	t.Run("already done: immediate no-op", func(t *testing.T) {
		db := database.NewMemDB()
		writeCheckpoints(db, 209133952, 209134208)
		WriteSupplyBlobFeeFixupDone(db) // pre-set done flag

		s := newFixupModule(db, big.NewInt(int64(params.KairosNetworkId)))
		s.applyBlobFeeFixupIfNeeded()

		// Nothing must change — bad checkpoints remain but the pointer is also untouched.
		assert.Equal(t, uint64(209134208), ReadLastAccRewardNumber(db), "pointer must not be changed")
		assert.NotNil(t, ReadAccReward(db, 209134080), "bad checkpoint must not be deleted when already done")
	})

	t.Run("idempotent: second call is a no-op", func(t *testing.T) {
		db := database.NewMemDB()
		writeCheckpoints(db, 209133952, 209134208)

		s := newFixupModule(db, big.NewInt(int64(params.KairosNetworkId)))
		s.applyBlobFeeFixupIfNeeded()
		// First call: resets to 209133952 and marks done.
		assert.Equal(t, uint64(209133952), ReadLastAccRewardNumber(db))

		// Write a new correct checkpoint above the old reset point (as catchup would).
		WriteAccReward(db, 209134080, &supply.AccReward{TotalMinted: big.NewInt(2000), BurntFee: big.NewInt(200)})
		WriteLastAccRewardNumber(db, 209134080)

		s.applyBlobFeeFixupIfNeeded()
		// Second call must not delete the newly correct checkpoint.
		assert.Equal(t, uint64(209134080), ReadLastAccRewardNumber(db))
		assert.NotNil(t, ReadAccReward(db, 209134080))
	})

	t.Run("unknown network: marks done immediately, no-op", func(t *testing.T) {
		db := database.NewMemDB()
		writeCheckpoints(db, 100, 100)

		s := newFixupModule(db, big.NewInt(99999)) // non-production chain ID
		s.applyBlobFeeFixupIfNeeded()

		assert.Equal(t, uint64(100), ReadLastAccRewardNumber(db), "pointer must be unchanged")
		assert.True(t, ReadSupplyBlobFeeFixupDone(db), "done flag must be set to skip future checks")
	})
}
