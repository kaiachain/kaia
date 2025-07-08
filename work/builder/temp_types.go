// TODO-Kaia: This file is temporarily used during refactoring.
package builder

import (
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
)

type (
	Bundle      = builder.Bundle
	TxOrGen     = builder.TxOrGen
	TxGenerator = builder.TxGenerator
)

func NewTxOrGenFromTx(tx *types.Transaction) *TxOrGen {
	return builder.NewTxOrGenFromTx(tx)
}

func NewTxOrGenFromGen(generator TxGenerator, id common.Hash) *TxOrGen {
	return builder.NewTxOrGenFromGen(generator, id)
}

func NewTxOrGenList(interfaces ...interface{}) []*TxOrGen {
	return builder.NewTxOrGenList(interfaces...)
}

type (
	TxBundlingModule     = builder.TxBundlingModule
	TxBundlingModuleHost = builder.TxBundlingModuleHost
)
