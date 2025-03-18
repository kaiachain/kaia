package builder

import (
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

type TxOrGen struct {
	ConcreteTx  *types.Transaction
	TxGenerator TxGenerator
	Id          common.Hash
}

func (t *TxOrGen) IsConcreteTx() bool {
	return t.ConcreteTx != nil
}

func (t *TxOrGen) IsTxGenerator() bool {
	return t.TxGenerator != nil
}

func NewTxOrGenFromTx(tx *types.Transaction) *TxOrGen {
	return &TxOrGen{
		ConcreteTx: tx,
		Id:         tx.Hash(),
	}
}

func NewTxOrGenListFromTxs(txs ...*types.Transaction) []*TxOrGen {
	txOrGens := make([]*TxOrGen, len(txs))
	for i, tx := range txs {
		txOrGens[i] = NewTxOrGenFromTx(tx)
	}
	return txOrGens
}

func NewTxOrGenFromGen(generator TxGenerator, id common.Hash) *TxOrGen {
	return &TxOrGen{
		TxGenerator: generator,
		Id:          id,
	}
}

// NewTxOrGenList creates a list of TxOrGen from a list of interfaces.
// Used for testing.
func NewTxOrGenList(interfaces ...interface{}) []*TxOrGen {
	txOrGens := make([]*TxOrGen, len(interfaces))
	for i, tx := range interfaces {
		switch v := tx.(type) {
		case *types.Transaction:
			txOrGens[i] = NewTxOrGenFromTx(v)
		case *TxOrGen:
			txOrGens[i] = v
		default:
			panic(fmt.Sprintf("unknown type: %T", v))
		}
	}
	return txOrGens
}

func (t *TxOrGen) Equals(other *TxOrGen) bool {
	return t.Id == other.Id
}

func (t *TxOrGen) GetTx(nonce uint64) (*types.Transaction, error) {
	if t.IsConcreteTx() {
		return t.ConcreteTx, nil
	}
	return t.TxGenerator(nonce)
}
