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
	"github.com/kaiachain/kaia/common"
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
	require.False(t, SanityCheckSignatures(validTestSignatures(maxSignatures+1), TxTypeValueTransfer))
}

func TestRecoverPubkeyRejectsTooManySignatures(t *testing.T) {
	tooManySignatures := validTestSignatures(int(accountkey.MaxNumKeysForMultiSig) + 1)
	vfuncCalled := false

	_, err := tooManySignatures.RecoverPubkey(common.Hash{}, true, func(v *big.Int) *big.Int {
		vfuncCalled = true
		return v
	})

	require.ErrorIs(t, err, kerrors.ErrMaxKeysExceed)
	require.False(t, vfuncCalled)
}

func TestTransactionDecodeChecksFeePayerSignatureLength(t *testing.T) {
	tests := []struct {
		name    string
		numSigs int
		wantErr bool
	}{
		{
			name:    "at limit",
			numSigs: int(accountkey.MaxNumKeysForMultiSig),
		},
		{
			name:    "over limit",
			numSigs: int(accountkey.MaxNumKeysForMultiSig) + 1,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			txData := newTxInternalDataFeeDelegatedValueTransfer()
			txData.SetSignature(validTestSignatures(1))
			txData.SetFeePayerSignatures(validTestSignatures(test.numSigs))
			tx := NewTx(txData)

			rawRLP, err := rlp.EncodeToBytes(tx)
			require.NoError(t, err)
			var decodedRLP Transaction
			err = rlp.DecodeBytes(rawRLP, &decodedRLP)
			if test.wantErr {
				require.ErrorIs(t, err, ErrInvalidSig)
			} else {
				require.NoError(t, err)
			}

			rawJSON, err := json.Marshal(tx)
			require.NoError(t, err)
			var decodedJSON Transaction
			err = json.Unmarshal(rawJSON, &decodedJSON)
			if test.wantErr {
				require.ErrorIs(t, err, ErrInvalidSig)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
