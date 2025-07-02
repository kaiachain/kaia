package impl

import (
	"testing"

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
