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

package builder

import (
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
)

type TxGenerator func(nonce uint64) (*types.Transaction, error)

type TxOrGen struct {
	concreteTx  *types.Transaction
	txGenerator TxGenerator
	Id          common.Hash // For concreteTx, txHash. For txGenerator, use a unique and deterministic identifier.
}

func NewTxOrGenFromTx(tx *types.Transaction) *TxOrGen {
	if tx == nil {
		return nil
	}

	return &TxOrGen{
		concreteTx: tx,
		Id:         tx.Hash(),
	}
}

func NewTxOrGenFromGen(generator TxGenerator, id common.Hash) *TxOrGen {
	if generator == nil {
		return nil
	}

	return &TxOrGen{
		txGenerator: generator,
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
		case TxGenerator:
			txOrGens[i] = NewTxOrGenFromGen(v, common.Hash{byte(i + 1)})
		case *TxOrGen:
			txOrGens[i] = v
		default:
			panic(fmt.Sprintf("unknown type: %T", v))
		}
	}
	return txOrGens
}

func (t *TxOrGen) GetTx(nonce uint64) (*types.Transaction, error) {
	if t.IsConcreteTx() {
		return t.concreteTx, nil
	}
	return t.txGenerator(nonce)
}

func (t *TxOrGen) IsConcreteTx() bool {
	return t.concreteTx != nil
}

func (t *TxOrGen) IsTxGenerator() bool {
	return t.txGenerator != nil
}

func (t *TxOrGen) Equals(other *TxOrGen) bool {
	return t.Id == other.Id
}
