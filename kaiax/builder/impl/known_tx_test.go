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
