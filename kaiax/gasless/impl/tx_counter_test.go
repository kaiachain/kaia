// Copyright 2025 The Kaia Authors
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
	"time"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestKnownTxs(t *testing.T) {
	knownTxs := knownTxs{}
	tx := types.NewTransaction(0, common.Address{}, common.Big0, 0, common.Big0, nil)
	knownTxs.add(tx, TxStatusQueue)
	ktx, ok := knownTxs.get(tx.Hash())
	assert.True(t, ok)
	assert.Equal(t, ktx.status, TxStatusQueue)
}

func TestKnownTxs_Timer(t *testing.T) {
	knownTx := knownTx{}
	assert.Equal(t, time.Time{}, knownTx.addedTime)
	assert.Equal(t, time.Time{}, knownTx.promotedTime)

	addedTime := knownTx.startAddedTimeIfZero()
	assert.Equal(t, addedTime, knownTx.addedTime)
	assert.Equal(t, time.Time{}, knownTx.promotedTime)

	pendingTime := knownTx.startPromotedTimeIfZero()
	assert.Equal(t, addedTime, knownTx.addedTime)
	assert.Equal(t, pendingTime, knownTx.promotedTime)

	// startAddedTimeIfZero() and startPromotedTimeIfZero() should not change the time if it is already set
	knownTx.startAddedTimeIfZero()
	assert.Equal(t, addedTime, knownTx.addedTime)
	assert.Equal(t, pendingTime, knownTx.promotedTime)

	knownTx.startPromotedTimeIfZero()
	assert.Equal(t, addedTime, knownTx.addedTime)
	assert.Equal(t, pendingTime, knownTx.promotedTime)
}
