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

package types

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Single-signature transaction types used to terminate the process through
// logger.Crit when SetSignature received a signature list whose length was not
// exactly one. Callers fed by external input, such as the fee-payer signing
// RPCs, must get an error instead so the node stays alive.
func TestSetSignature_SingleSignatureTypes(t *testing.T) {
	sig := &TxSignature{V: big.NewInt(27), R: big.NewInt(2), S: big.NewInt(3)}
	gens := map[string]func() TxInternalData{
		"legacy":     genLegacyTransaction,
		"accessList": genAccessListTransaction,
		"dynamicFee": genDynamicFeeTransaction,
		"blob":       genBlobTransaction,
		"setCode":    genSetCodeTransaction,
	}
	for name, gen := range gens {
		t.Run(name, func(t *testing.T) {
			txdata := gen()
			assert.ErrorIs(t, txdata.SetSignature(nil), ErrInvalidSignatureCount)
			assert.ErrorIs(t, txdata.SetSignature(TxSignatures{}), ErrInvalidSignatureCount)
			assert.ErrorIs(t, txdata.SetSignature(TxSignatures{sig, sig}), ErrInvalidSignatureCount)

			assert.NoError(t, txdata.SetSignature(TxSignatures{sig}))
			got := txdata.RawSignatureValues()
			if assert.Len(t, got, 1) {
				assert.Equal(t, sig.V, got[0].V)
				assert.Equal(t, sig.R, got[0].R)
				assert.Equal(t, sig.S, got[0].S)
			}
		})
	}
}

// Kaia typed transactions accept any number of signatures (multisig accounts).
func TestSetSignature_MultiSignatureTypes(t *testing.T) {
	sig := &TxSignature{V: big.NewInt(27), R: big.NewInt(2), S: big.NewInt(3)}
	txdata := genValueTransferTransaction()
	assert.NoError(t, txdata.SetSignature(TxSignatures{sig, sig}))
	assert.Len(t, txdata.RawSignatureValues(), 2)
}

// Transaction.SetSignature is the entry point used by the RPC layer; it must
// propagate the error rather than hide it.
func TestTransactionSetSignature_PropagatesError(t *testing.T) {
	sig := &TxSignature{V: big.NewInt(27), R: big.NewInt(2), S: big.NewInt(3)}
	tx := &Transaction{data: genLegacyTransaction()}
	assert.ErrorIs(t, tx.SetSignature(TxSignatures{sig, sig}), ErrInvalidSignatureCount)
	assert.NoError(t, tx.SetSignature(TxSignatures{sig}))
}
