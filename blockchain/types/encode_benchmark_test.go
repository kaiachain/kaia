// Modifications Copyright 2024 The Kaia Authors
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
	"bytes"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/rlp"
)

func BenchmarkTxEncode(b *testing.B) {
	txs := []struct {
		Name string
		tx   TxInternalData
	}{
		{"OriginalTx", genLegacyTransaction()},
		{"ValueTransfer", genValueTransferTransaction()},
		{"FeeDelegatedValueTransfer", genFeeDelegatedValueTransferTransaction()},
		{"FeeDelegatedValueTransferWithRatio", genFeeDelegatedValueTransferWithRatioTransaction()},
		{"AccountUpdate", genAccountUpdateTransaction()},
		{"SmartContractDeploy", genSmartContractDeployTransaction()},
		{"SmartContractExecution", genSmartContractExecutionTransaction()},
		{"ChainDataTx", genChainDataTransaction()},
	}
	testcases := []struct {
		Name string
		fn   func(b *testing.B, tx TxInternalData)
	}{
		{"Encode", benchmarkEncode},
		{"EncodeToBytes", benchmarkEncodeToBytes},
		{"EncodeInterfaceSlice", benchmarkInterfaceSlice},
		{"EncodeInterfaceSliceLoop", benchmarkInterfaceSliceLoop},
	}

	for _, tx := range txs {
		for _, test := range testcases {
			Name := test.Name + "/" + tx.Name
			b.Run(Name, func(b *testing.B) {
				test.fn(b, tx.tx)
			})
		}
	}
}

func benchmarkEncode(b *testing.B, txInternal TxInternalData) {
	var i int
	tx := &Transaction{data: txInternal}
	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		buffer := new(bytes.Buffer)
		if err := rlp.Encode(buffer, tx); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func benchmarkEncodeToBytes(b *testing.B, txInternal TxInternalData) {
	var i int
	tx := &Transaction{data: txInternal}
	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		_, err := rlp.EncodeToBytes(tx)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func benchmarkInterfaceSlice(b *testing.B, txInternal TxInternalData) {
	var i int
	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		ifs := getInterfaceSlice(txInternal)
		_, err := rlp.EncodeToBytes(ifs)
		if err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
}

func benchmarkInterfaceSliceLoop(b *testing.B, txInternal TxInternalData) {
	var i int
	b.ResetTimer()
	for i = 0; i < b.N; i++ {
		buffer := new(bytes.Buffer)
		ifs := getInterfaceSlice(txInternal)
		for _, it := range ifs {
			if err := rlp.Encode(buffer, it); err != nil {
				b.Fatal(err)
			}
		}
	}
	b.StopTimer()
}

func getInterfaceSlice(tx TxInternalData) []interface{} {
	return tx.(SliceMaker).MakeInterfaceSlice()
}

type SliceMaker interface {
	MakeInterfaceSlice() []interface{}
}

func (v *TxInternalDataLegacy) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.Payload,
		v.V,
		v.R,
		v.S,
	}
}

func (v *TxInternalDataValueTransfer) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.TxSignatures,
	}
}

func (v *TxInternalDataFeeDelegatedValueTransfer) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.TxSignatures,
		v.FeePayer,
		v.FeePayerSignatures,
	}
}

func (v *TxInternalDataFeeDelegatedValueTransferWithRatio) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.FeeRatio,
		v.TxSignatures,
		v.FeePayer,
		v.FeePayerSignatures,
	}
}

func (v *TxInternalDataAccountCreation) MakeInterfaceSlice() []interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(v.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.HumanReadable,
		keyEnc,
		v.TxSignatures,
	}
}

func (v *TxInternalDataAccountUpdate) MakeInterfaceSlice() []interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(v.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.From,
		keyEnc,
		v.TxSignatures,
	}
}

func (v *TxInternalDataSmartContractDeploy) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.Payload,
		v.HumanReadable,
		v.TxSignatures,
	}
}

func (v *TxInternalDataSmartContractExecution) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.Recipient,
		v.Amount,
		v.From,
		v.Payload,
		v.TxSignatures,
	}
}

func (v *TxInternalDataChainDataAnchoring) MakeInterfaceSlice() []interface{} {
	return []interface{}{
		v.Type(),
		v.AccountNonce,
		v.Price,
		v.GasLimit,
		v.From,
		v.Payload,
		v.TxSignatures,
	}
}
