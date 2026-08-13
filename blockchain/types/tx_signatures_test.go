// Modifications Copyright 2026 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/require"
)

func validTestSignatures(n int) TxSignatures {
	sig := &TxSignature{
		V: big.NewInt(37),
		R: big.NewInt(1),
		S: big.NewInt(1),
	}
	sigs := make(TxSignatures, n)
	for i := range sigs {
		sigs[i] = sig
	}
	return sigs
}

func TestSanityCheckSignaturesLength(t *testing.T) {
	maxSignatures := int(accountkey.MaxNumKeysForMultiSig)

	require.False(t, SanityCheckSignatures(nil, TxTypeValueTransfer))
	require.True(t, SanityCheckSignatures(validTestSignatures(1), TxTypeValueTransfer))
	require.True(t, SanityCheckSignatures(validTestSignatures(maxSignatures), TxTypeValueTransfer))
	require.True(t, SanityCheckSignatures(validTestSignatures(maxSignatures+1), TxTypeValueTransfer))
}

func TestTransactionSignatureListLength(t *testing.T) {
	maxSignatures := int(accountkey.MaxNumKeysForMultiSig)

	senderTx := NewTx(newTxInternalDataValueTransfer())
	senderTx.SetSignature(validTestSignatures(maxSignatures))
	require.NoError(t, senderTx.ValidateSignatureListLength())
	senderTx.SetSignature(validTestSignatures(maxSignatures + 1))
	require.ErrorIs(t, senderTx.ValidateSignatureListLength(), kerrors.ErrMaxKeysExceed)

	feePayerTxData := newTxInternalDataFeeDelegatedValueTransfer()
	feePayerTxData.SetSignature(validTestSignatures(1))
	feePayerTxData.SetFeePayerSignatures(validTestSignatures(maxSignatures))
	feePayerTx := NewTx(feePayerTxData)
	require.NoError(t, feePayerTx.ValidateSignatureListLength())
	feePayerTxData.SetFeePayerSignatures(validTestSignatures(maxSignatures + 1))
	require.ErrorIs(t, feePayerTx.ValidateSignatureListLength(), kerrors.ErrMaxKeysExceed)
}

func TestTransactionDecodeAllowsOversizedSignatureLists(t *testing.T) {
	maxSignatures := int(accountkey.MaxNumKeysForMultiSig)
	tests := []struct {
		name string
		tx   *Transaction
	}{
		{
			name: "sender signatures",
			tx: func() *Transaction {
				txData := newTxInternalDataValueTransfer()
				txData.SetSignature(validTestSignatures(maxSignatures + 1))
				return NewTx(txData)
			}(),
		},
		{
			name: "fee payer signatures",
			tx: func() *Transaction {
				txData := newTxInternalDataFeeDelegatedValueTransfer()
				txData.SetSignature(validTestSignatures(1))
				txData.SetFeePayerSignatures(validTestSignatures(maxSignatures + 1))
				return NewTx(txData)
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rawRLP, err := rlp.EncodeToBytes(test.tx)
			require.NoError(t, err)
			var decodedRLP Transaction
			require.NoError(t, rlp.DecodeBytes(rawRLP, &decodedRLP))
			require.ErrorIs(t, decodedRLP.ValidateSignatureListLength(), kerrors.ErrMaxKeysExceed)

			rawJSON, err := json.Marshal(test.tx)
			require.NoError(t, err)
			var decodedJSON Transaction
			require.NoError(t, json.Unmarshal(rawJSON, &decodedJSON))
			require.ErrorIs(t, decodedJSON.ValidateSignatureListLength(), kerrors.ErrMaxKeysExceed)
		})
	}
}
