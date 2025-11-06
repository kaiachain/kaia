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
	"encoding/json"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	emptyBlob          = new(kzg4844.Blob)
	emptyBlobCommit, _ = kzg4844.BlobToCommitment(emptyBlob)
	emptyBlobProof, _  = kzg4844.ComputeBlobProof(emptyBlob, emptyBlobCommit)
)

var (
	to             = common.HexToAddress("7b65B75d204aBed71587c9E519a89277766EE1d0")
	key, from      = defaultTestKey()
	feePayer       = common.HexToAddress("5A0043070275d9f6054307Ee7348bD660849D90f")
	nonce          = uint64(1234)
	amount         = big.NewInt(10)
	gasLimit       = uint64(1000000)
	gasPrice       = big.NewInt(25)
	gasTipCap      = big.NewInt(25)
	gasFeeCap      = big.NewInt(25)
	accesses       = AccessList{{Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), StorageKeys: []common.Hash{{0}}}}
	sidecar        = NewBlobTxSidecar(BlobSidecarVersion0, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, []kzg4844.Proof{emptyBlobProof})
	authorizations = []SetCodeAuthorization{{ChainID: *uint256.NewInt(2), Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), Nonce: nonce, V: uint8(0), R: *uint256.NewInt(0), S: *uint256.NewInt(0)}}
)

// TestTransactionSerialization tests RLP/JSON serialization for TxInternalData
func TestTransactionSerialization(t *testing.T) {
	txs := []struct {
		Name string
		tx   TxInternalData
	}{
		{"OriginalTx", genLegacyTransaction()},
		{"SmartContractDeploy", genSmartContractDeployTransaction()},
		{"FeeDelegatedSmartContractDeploy", genFeeDelegatedSmartContractDeployTransaction()},
		{"FeeDelegatedSmartContractDeployWithRatio", genFeeDelegatedSmartContractDeployWithRatioTransaction()},
		{"ValueTransfer", genValueTransferTransaction()},
		{"ValueTransferMemo", genValueTransferMemoTransaction()},
		{"FeeDelegatedValueTransferMemo", genFeeDelegatedValueTransferMemoTransaction()},
		{"FeeDelegatedValueTransferMemoWithRatio", genFeeDelegatedValueTransferMemoWithRatioTransaction()},
		{"ChainDataAnchoring", genChainDataTransaction()},
		{"FeeDelegatedChainDataAnchoring", genFeeDelegatedChainDataTransaction()},
		{"FeeDelegatedChainDataAnchoringWithRatio", genFeeDelegatedChainDataWithRatioTransaction()},
		{"AccountUpdate", genAccountUpdateTransaction()},
		{"FeeDelegatedAccountUpdate", genFeeDelegatedAccountUpdateTransaction()},
		{"FeeDelegatedAccountUpdateWithRatio", genFeeDelegatedAccountUpdateWithRatioTransaction()},
		{"FeeDelegatedValueTransfer", genFeeDelegatedValueTransferTransaction()},
		{"SmartContractExecution", genSmartContractExecutionTransaction()},
		{"FeeDelegatedSmartContractExecution", genFeeDelegatedSmartContractExecutionTransaction()},
		{"FeeDelegatedSmartContractExecutionWithRatio", genFeeDelegatedSmartContractExecutionWithRatioTransaction()},
		{"FeeDelegatedValueTransferWithRatio", genFeeDelegatedValueTransferWithRatioTransaction()},
		{"Cancel", genCancelTransaction()},
		{"FeeDelegatedCancel", genFeeDelegatedCancelTransaction()},
		{"FeeDelegatedCancelWithRatio", genFeeDelegatedCancelWithRatioTransaction()},
		{"AccessList", genAccessListTransaction()},
		{"DynamicFee", genDynamicFeeTransaction()},
		{"Blob", genBlobTransaction()},
		{"SetCode", genSetCodeTransaction()},
	}

	testcases := []struct {
		Name string
		fn   func(t *testing.T, tx TxInternalData)
	}{
		{"RLP", testTransactionRLP},
		{"JSON", testTransactionJSON},
		{"RPC", testTransactionRPC},
	}

	txMap := make(map[TxType]TxInternalData)
	for _, test := range testcases {
		for _, tx := range txs {
			txMap[tx.tx.Type()] = tx.tx
			Name := test.Name + "/" + tx.Name
			t.Run(Name, func(t *testing.T) {
				test.fn(t, tx.tx)
			})
		}
	}

	// Below code checks whether serialization for all tx implementations is done or not.
	// If no serialization, make test failed.
	for i := TxTypeLegacyTransaction; i < TxTypeEthereumLast; i++ {
		if i == TxTypeKaiaLast {
			i = TxTypeEthereumAccessList
		}

		tx, err := NewTxInternalData(i)
		// TxTypeAccountCreation is not supported now
		if i == TxTypeAccountCreation {
			continue
		}
		if err == nil {
			if _, ok := txMap[tx.Type()]; !ok {
				t.Errorf("No serialization test for tx %s", tx.Type().String())
			}
		}
	}
}

func testTransactionRLP(t *testing.T, tx TxInternalData) {
	enc := newTxInternalDataSerializerWithValues(tx)

	signer := MakeSigner(params.TestChainConfig, big.NewInt(2))
	rawTx := &Transaction{data: tx}
	rawTx.Sign(signer, key)

	if _, ok := tx.(TxInternalDataFeePayer); ok {
		rawTx.SignFeePayer(signer, key)
	}

	b, err := rlp.EncodeToBytes(enc)
	if err != nil {
		panic(err)
	}

	if tx.Type().IsEthTypedTransaction() {
		assert.Equal(t, byte(EthereumTxTypeEnvelope), b[0])
	}

	dec := newTxInternalDataSerializer()

	if err := rlp.DecodeBytes(b, &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

func testTransactionJSON(t *testing.T, tx TxInternalData) {
	enc := newTxInternalDataSerializerWithValues(tx)

	signer := MakeSigner(params.TestChainConfig, big.NewInt(2))
	rawTx := &Transaction{data: tx}
	rawTx.Sign(signer, key)

	if _, ok := tx.(TxInternalDataFeePayer); ok {
		rawTx.SignFeePayer(signer, key)
	}

	b, err := json.Marshal(enc)
	if err != nil {
		panic(err)
	}

	dec := newTxInternalDataSerializer()

	if err := json.Unmarshal(b, &dec); err != nil {
		panic(err)
	}

	if !tx.Equal(dec.tx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, dec.tx)
	}
}

// Copied from api/api_public_blockchain.go
func newRPCTransaction(tx *Transaction, blockHash common.Hash, blockNumber uint64, index uint64) map[string]interface{} {
	var from common.Address
	if tx.IsEthereumTransaction() {
		signer := LatestSignerForChainID(tx.ChainId())
		from, _ = Sender(signer, tx)
	} else {
		from, _ = tx.From()
	}

	output := tx.MakeRPCOutput()

	output["blockHash"] = blockHash
	output["blockNumber"] = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
	output["from"] = from
	output["hash"] = tx.Hash()
	output["transactionIndex"] = hexutil.Uint(index)

	return output
}

func testTransactionRPC(t *testing.T, tx TxInternalData) {
	signer := LatestSignerForChainID(big.NewInt(2))
	rawTx := &Transaction{data: tx}
	rawTx.Sign(signer, key)

	if _, ok := tx.(TxInternalDataFeePayer); ok {
		rawTx.SignFeePayer(signer, key)
	}

	h := rawTx.Hash()
	tx.SetHash(&h)

	// Copied from newRPCTransaction
	rpcout := newRPCTransaction(rawTx, common.Hash{}, 0, 0)
	if tx.Type().IsEthTypedTransaction() {
		if _, ok := rpcout["chainId"]; !ok {
			t.Fatalf("The chainId field must be presented.")
		}
	}

	b, err := json.Marshal(rpcout)
	if err != nil {
		panic(err)
	}

	decTx := &Transaction{}

	if err := json.Unmarshal(b, decTx); err != nil {
		panic(err)
	}

	if !rawTx.Equal(decTx) {
		t.Fatalf("tx != dec.tx\ntx=%v\ndec.tx=%v", tx, decTx)
	}
}

func genLegacyTransaction() TxInternalData {
	txdata, err := NewTxInternalDataWithMap(TxTypeLegacyTransaction, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyData:     []byte("1234"),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return txdata
}

func genAccessListTransaction() TxInternalData {
	tx, err := NewTxInternalDataWithMap(TxTypeEthereumAccessList, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:      nonce,
		TxValueKeyTo:         &to,
		TxValueKeyAmount:     amount,
		TxValueKeyGasLimit:   gasLimit,
		TxValueKeyGasPrice:   gasPrice,
		TxValueKeyData:       []byte("1234"),
		TxValueKeyAccessList: accesses,
		TxValueKeyChainID:    big.NewInt(2),
	})
	if err != nil {
		panic(err)
	}

	return tx
}

func genDynamicFeeTransaction() TxInternalData {
	tx, err := NewTxInternalDataWithMap(TxTypeEthereumDynamicFee, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:      nonce,
		TxValueKeyTo:         &to,
		TxValueKeyAmount:     amount,
		TxValueKeyGasLimit:   gasLimit,
		TxValueKeyGasFeeCap:  gasFeeCap,
		TxValueKeyGasTipCap:  gasTipCap,
		TxValueKeyData:       []byte("1234"),
		TxValueKeyAccessList: accesses,
		TxValueKeyChainID:    big.NewInt(2),
	})
	if err != nil {
		panic(err)
	}

	return tx
}

func genBlobTransaction() TxInternalData {
	tx, err := NewTxInternalDataWithMap(TxTypeEthereumBlob, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:      nonce,
		TxValueKeyTo:         to,
		TxValueKeyAmount:     amount,
		TxValueKeyGasLimit:   gasLimit,
		TxValueKeyGasFeeCap:  gasFeeCap,
		TxValueKeyGasTipCap:  gasTipCap,
		TxValueKeyData:       []byte("1234"),
		TxValueKeyAccessList: accesses,
		TxValueKeyBlobFeeCap: gasFeeCap,
		TxValueKeyBlobHashes: []common.Hash{{0}},
		TxValueKeySidecar:    sidecar,
		TxValueKeyChainID:    big.NewInt(2),
	})
	if err != nil {
		panic(err)
	}

	return tx
}

func genSetCodeTransaction() TxInternalData {
	tx, err := NewTxInternalDataWithMap(TxTypeEthereumSetCode, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:             nonce,
		TxValueKeyTo:                to,
		TxValueKeyAmount:            amount,
		TxValueKeyGasLimit:          gasLimit,
		TxValueKeyGasFeeCap:         gasFeeCap,
		TxValueKeyGasTipCap:         gasTipCap,
		TxValueKeyData:              []byte("1234"),
		TxValueKeyAccessList:        accesses,
		TxValueKeyAuthorizationList: authorizations,
		TxValueKeyChainID:           big.NewInt(2),
	})
	if err != nil {
		panic(err)
	}

	return tx
}

func genValueTransferTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeValueTransfer, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genValueTransferMemoTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeValueTransferMemo, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		TxValueKeyData:     []byte(string("hello")),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedValueTransferMemoTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedValueTransferMemo, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		TxValueKeyData:     []byte(string("hello")),
		TxValueKeyFeePayer: feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedValueTransferMemoWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedValueTransferMemoWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:              nonce,
		TxValueKeyTo:                 to,
		TxValueKeyAmount:             amount,
		TxValueKeyGasLimit:           gasLimit,
		TxValueKeyGasPrice:           gasPrice,
		TxValueKeyFrom:               from,
		TxValueKeyData:               []byte(string("hello")),
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genSmartContractDeployTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeSmartContractDeploy, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:         nonce,
		TxValueKeyAmount:        amount,
		TxValueKeyGasLimit:      gasLimit,
		TxValueKeyGasPrice:      gasPrice,
		TxValueKeyHumanReadable: true,
		TxValueKeyTo:            &to,
		TxValueKeyFrom:          from,
		// The binary below is a compiled binary of KlaytnReward.sol.
		TxValueKeyData:       common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"),
		TxValueKeyCodeFormat: params.CodeFormatEVM,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedSmartContractDeployTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedSmartContractDeploy, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:         nonce,
		TxValueKeyAmount:        amount,
		TxValueKeyGasLimit:      gasLimit,
		TxValueKeyGasPrice:      gasPrice,
		TxValueKeyHumanReadable: true,
		TxValueKeyTo:            &to,
		TxValueKeyFrom:          from,
		TxValueKeyFeePayer:      feePayer,
		// The binary below is a compiled binary of KlaytnReward.sol.
		TxValueKeyData:       common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"),
		TxValueKeyCodeFormat: params.CodeFormatEVM,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedSmartContractDeployWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedSmartContractDeployWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:         nonce,
		TxValueKeyAmount:        amount,
		TxValueKeyGasLimit:      gasLimit,
		TxValueKeyGasPrice:      gasPrice,
		TxValueKeyHumanReadable: true,
		TxValueKeyTo:            &to,
		TxValueKeyFrom:          from,
		TxValueKeyFeePayer:      feePayer,
		// The binary below is a compiled binary of KlaytnReward.sol.
		TxValueKeyData:               common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029"),
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
		TxValueKeyCodeFormat:         params.CodeFormatEVM,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genChainDataTransaction() TxInternalData {
	data := &AnchoringDataInternalType0{
		common.HexToHash("0"), common.HexToHash("1"),
		common.HexToHash("2"), common.HexToHash("3"),
		common.HexToHash("4"), big.NewInt(5), big.NewInt(6), big.NewInt(7),
	}
	encodedCCTxData, err := rlp.EncodeToBytes(data)
	if err != nil {
		panic(err)
	}
	blockTxData := &AnchoringData{0, encodedCCTxData}

	anchoredData, err := rlp.EncodeToBytes(blockTxData)
	if err != nil {
		panic(err)
	}

	txdata, err := NewTxInternalDataWithMap(TxTypeChainDataAnchoring, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:        nonce,
		TxValueKeyFrom:         from,
		TxValueKeyGasLimit:     gasLimit,
		TxValueKeyGasPrice:     gasPrice,
		TxValueKeyAnchoredData: anchoredData,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return txdata
}

func genFeeDelegatedChainDataTransaction() TxInternalData {
	data := &AnchoringDataInternalType0{
		common.HexToHash("0"), common.HexToHash("1"),
		common.HexToHash("2"), common.HexToHash("3"),
		common.HexToHash("4"), big.NewInt(5), big.NewInt(6), big.NewInt(7),
	}
	encodedCCTxData, err := rlp.EncodeToBytes(data)
	if err != nil {
		panic(err)
	}
	blockTxData := &AnchoringData{0, encodedCCTxData}

	anchoredData, err := rlp.EncodeToBytes(blockTxData)
	if err != nil {
		panic(err)
	}

	txdata, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedChainDataAnchoring, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:        nonce,
		TxValueKeyFrom:         from,
		TxValueKeyGasLimit:     gasLimit,
		TxValueKeyGasPrice:     gasPrice,
		TxValueKeyAnchoredData: anchoredData,
		TxValueKeyFeePayer:     feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return txdata
}

func genFeeDelegatedChainDataWithRatioTransaction() TxInternalData {
	data := &AnchoringDataInternalType0{
		common.HexToHash("0"), common.HexToHash("1"),
		common.HexToHash("2"), common.HexToHash("3"),
		common.HexToHash("4"), big.NewInt(5), big.NewInt(6), big.NewInt(7),
	}
	encodedCCTxData, err := rlp.EncodeToBytes(data)
	if err != nil {
		panic(err)
	}
	blockTxData := &AnchoringData{0, encodedCCTxData}

	anchoredData, err := rlp.EncodeToBytes(blockTxData)
	if err != nil {
		panic(err)
	}

	txdata, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedChainDataAnchoringWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:              nonce,
		TxValueKeyFrom:               from,
		TxValueKeyGasLimit:           gasLimit,
		TxValueKeyGasPrice:           gasPrice,
		TxValueKeyAnchoredData:       anchoredData,
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return txdata
}

//func genAccountCreationTransaction() TxInternalData {
//	d, err := NewTxInternalDataWithMap(TxTypeAccountCreation, map[TxValueKeyType]interface{}{
//		TxValueKeyNonce:         nonce,
//		TxValueKeyTo:            to,
//		TxValueKeyAmount:        amount,
//		TxValueKeyGasLimit:      gasLimit,
//		TxValueKeyGasPrice:      gasPrice,
//		TxValueKeyFrom:          from,
//		TxValueKeyHumanReadable: false,
//		TxValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&key.PublicKey),
//	})
//
//	if err != nil {
//		// Since we do not have testing.T here, call panic() instead of t.Fatal().
//		panic(err)
//	}
//
//	return d
//}

func genFeeDelegatedValueTransferTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedValueTransfer, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		TxValueKeyFeePayer: feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedValueTransferWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedValueTransferWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:              nonce,
		TxValueKeyTo:                 to,
		TxValueKeyAmount:             amount,
		TxValueKeyGasLimit:           gasLimit,
		TxValueKeyGasPrice:           gasPrice,
		TxValueKeyFrom:               from,
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genSmartContractExecutionTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeSmartContractExecution, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		// A abi-packed bytes calling "reward" of KlaytnReward.sol with an address "bc5951f055a85f41a3b62fd6f68ab7de76d299b2".
		TxValueKeyData: common.Hex2Bytes("6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2"),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedSmartContractExecutionTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedSmartContractExecution, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		// A abi-packed bytes calling "reward" of KlaytnReward.sol with an address "bc5951f055a85f41a3b62fd6f68ab7de76d299b2".
		TxValueKeyData:     common.Hex2Bytes("6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2"),
		TxValueKeyFeePayer: feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedSmartContractExecutionWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedSmartContractExecutionWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyTo:       to,
		TxValueKeyAmount:   amount,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		// A abi-packed bytes calling "reward" of KlaytnReward.sol with an address "bc5951f055a85f41a3b62fd6f68ab7de76d299b2".
		TxValueKeyData:               common.Hex2Bytes("6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2"),
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genAccountUpdateTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeAccountUpdate, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:      nonce,
		TxValueKeyGasLimit:   gasLimit,
		TxValueKeyGasPrice:   gasPrice,
		TxValueKeyFrom:       from,
		TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&key.PublicKey),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedAccountUpdateTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedAccountUpdate, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:      nonce,
		TxValueKeyGasLimit:   gasLimit,
		TxValueKeyGasPrice:   gasPrice,
		TxValueKeyFrom:       from,
		TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(&key.PublicKey),
		TxValueKeyFeePayer:   feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedAccountUpdateWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedAccountUpdateWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:              nonce,
		TxValueKeyGasLimit:           gasLimit,
		TxValueKeyGasPrice:           gasPrice,
		TxValueKeyFrom:               from,
		TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(&key.PublicKey),
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genCancelTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeCancel, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedCancelTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedCancel, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:    nonce,
		TxValueKeyGasLimit: gasLimit,
		TxValueKeyGasPrice: gasPrice,
		TxValueKeyFrom:     from,
		TxValueKeyFeePayer: feePayer,
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

func genFeeDelegatedCancelWithRatioTransaction() TxInternalData {
	d, err := NewTxInternalDataWithMap(TxTypeFeeDelegatedCancelWithRatio, map[TxValueKeyType]interface{}{
		TxValueKeyNonce:              nonce,
		TxValueKeyGasLimit:           gasLimit,
		TxValueKeyGasPrice:           gasPrice,
		TxValueKeyFrom:               from,
		TxValueKeyFeePayer:           feePayer,
		TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
	})
	if err != nil {
		// Since we do not have testing.T here, call panic() instead of t.Fatal().
		panic(err)
	}

	return d
}

type serializeTC struct {
	Name         string
	Type         TxType
	Map          map[TxValueKeyType]interface{}
	ChainID      uint64
	SenderSigs   txSigsHex
	FeePayerSigs txSigsHex

	// Expected serializations
	// See https://docs.kaia.io/build/transactions/ for the definitions of the RLP encodings.
	SigRLP          string
	SigFeePayerRLP  string
	SenderTxHashRLP string
	TxHashRLP       string
	TxJson          string
	RpcJson         string
}

// txSigHex is for human-readable test cases.
type txSigHex struct {
	v uint64
	r string
	s string
}

type txSigsHex []txSigHex

func (s txSigsHex) TxSignatures() TxSignatures {
	sigs := make([]*TxSignature, len(s))
	for i, sig := range s {
		sigs[i] = &TxSignature{
			V: new(big.Int).SetUint64(sig.v),
			R: new(big.Int).SetBytes(common.HexToHash(sig.r).Bytes()),
			S: new(big.Int).SetBytes(common.HexToHash(sig.s).Bytes()),
		}
	}
	return sigs
}

func addrPtr(addr common.Address) *common.Address {
	return &addr
}

// Some cases taken from kaia-sdk: https://github.com/kaiachain/kaia-sdk/blob/dev/js-ext-core/test/tx.spec.ts
// Some cases taken from viem: https://github.com/wevm/viem/blob/main/vectors/src/transaction.json.gz
// - Viem vectors can be filtered with (e.g.):
//   o = JSON.parse(fs.readFileSync('transaction.json'))
//   o.filter((x) => x.name.substr(0,6) == "legacy").filter((x) => { l = x.transaction.data?.length; return (l > 0 && l < 100) }).filter((x) => { t = x.transaction; return (t.data && t.to && t.gas && t.gasPrice && t.nonce && t.value) })
// - Viem vectors don't contain SigRLP. Correct SigRLPs were manually generated using (e.g.):
//   rlp encode '["0x24ae","0x07a3d3","0x3cd980f3","0x1cd12f7edecac097265d53754b782004bf0b8fb7","0x7aa930ce8a1a","0xe02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e","0x01","",""]'
// - Viem vectors are manually prefixed with EthereumTxTypeEnvelope(0x78).
var serializeTCs = []serializeTC{
	{
		Name: "00_Legacy", // Viem "legacy: 6840"
		Type: TxTypeLegacyTransaction,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(9390),
			TxValueKeyGasPrice: big.NewInt(500691),
			TxValueKeyGasLimit: uint64(1020887283),
			TxValueKeyTo:       common.HexToAddress("0x1cd12f7edecac097265d53754b782004bf0b8fb7"),
			TxValueKeyAmount:   big.NewInt(134867086903834),
			TxValueKeyData:     hexutil.MustDecode("0xe02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e"),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 28, r: "0x8c3d94705e12d605d1144ac4c29bfcde87a06cf8424f2addddbc64668e9d78ba", s: "0x233281c375ddedb00ce1f5fb167f80e690701db29f77f1c372cbe72895147299"}},
		SigRLP:     "0xf8448224ae8307a3d3843cd980f3941cd12f7edecac097265d53754b782004bf0b8fb7867aa930ce8a1a98e02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e018080",
		TxHashRLP:  "0xf8848224ae8307a3d3843cd980f3941cd12f7edecac097265d53754b782004bf0b8fb7867aa930ce8a1a98e02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e1ca08c3d94705e12d605d1144ac4c29bfcde87a06cf8424f2addddbc64668e9d78baa0233281c375ddedb00ce1f5fb167f80e690701db29f77f1c372cbe72895147299",
	},
	{
		Name: "01_AccessList", // Viem "eip2930: 7007"
		Type: TxTypeEthereumAccessList,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:  big.NewInt(9),
			TxValueKeyNonce:    uint64(17013),
			TxValueKeyGasPrice: big.NewInt(0),
			TxValueKeyGasLimit: uint64(3212800333),
			TxValueKeyTo:       addrPtr(common.HexToAddress("0x610590408a97e8f84152d45cbe2ec0b08059598f")),
			TxValueKeyAmount:   big.NewInt(0),
			TxValueKeyData:     hexutil.MustDecode("0xd962f328314a9cd384c349a461fdba53ead623afb4f751df94a604"),
			TxValueKeyAccessList: AccessList{
				{
					Address: common.HexToAddress("0x67d560ba27b75d467fcf658e02ac3765c3634056"),
					StorageKeys: []common.Hash{
						common.HexToHash("0xa12c9f0600000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x5b00000000000000000000000000000000000000000000000000000000000000"),
						common.HexToHash("0x082b41eb5f7244029351c2533948000000000000000000000000000000000000"),
					},
				},
			},
		},
		ChainID:    9,
		SenderSigs: []txSigHex{{v: 1, r: "0xc71988016460ca44c352598cc116cf2ac1b9293387f2f8287aa889b6cf55b1e7", s: "0x634101cdd4f96d642fe8509166cc63729722cb51c80f1624b4fe693106b0664c"}},
		SigRLP:     "0x01f8ba098242758084bf7f714d94610590408a97e8f84152d45cbe2ec0b08059598f809bd962f328314a9cd384c349a461fdba53ead623afb4f751df94a604f87cf87a9467d560ba27b75d467fcf658e02ac3765c3634056f863a0a12c9f0600000000000000000000000000000000000000000000000000000000a05b00000000000000000000000000000000000000000000000000000000000000a0082b41eb5f7244029351c2533948000000000000000000000000000000000000",
		TxHashRLP:  "0x7801f8fd098242758084bf7f714d94610590408a97e8f84152d45cbe2ec0b08059598f809bd962f328314a9cd384c349a461fdba53ead623afb4f751df94a604f87cf87a9467d560ba27b75d467fcf658e02ac3765c3634056f863a0a12c9f0600000000000000000000000000000000000000000000000000000000a05b00000000000000000000000000000000000000000000000000000000000000a0082b41eb5f7244029351c253394800000000000000000000000000000000000001a0c71988016460ca44c352598cc116cf2ac1b9293387f2f8287aa889b6cf55b1e7a0634101cdd4f96d642fe8509166cc63729722cb51c80f1624b4fe693106b0664c",
	},
	{
		Name: "08_ValueTransfer", // kaia-sdk
		Type: TxTypeValueTransfer,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 0x25, r: "0xf3d0cd43661cabf53425535817c5058c27781f478cb5459874feaa462ed3a29a", s: "0x6748abe186269ff10b8100a4b7d7fea274b53ea2905acbf498dc8b5ab1bf4fbc"}},
		SigRLP:     "0xf839b5f4088204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b018080",
		TxHashRLP:  "0x08f87a8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a0f3d0cd43661cabf53425535817c5058c27781f478cb5459874feaa462ed3a29aa06748abe186269ff10b8100a4b7d7fea274b53ea2905acbf498dc8b5ab1bf4fbc",
		TxJson: `{
				"typeInt":8,
				"type": "TxTypeValueTransfer",
				"nonce": "0x4d2",
				"gasPrice": "0x19",
				"gas": "0xf4240",
				"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
				"value": "0xa",
				"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
				"signatures": [{"V": "0x25", "R": "0xf3d0cd43661cabf53425535817c5058c27781f478cb5459874feaa462ed3a29a", "S": "0x6748abe186269ff10b8100a4b7d7fea274b53ea2905acbf498dc8b5ab1bf4fbc"}],
				"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
			}`,
		RpcJson: `{
				"typeInt":8,
				"type": "TxTypeValueTransfer",
				"nonce": "0x4d2",
				"gasPrice": "0x19",
				"gas": "0xf4240",
				"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
				"value": "0xa",
				"signatures": [{"V": "0x25", "R": "0xf3d0cd43661cabf53425535817c5058c27781f478cb5459874feaa462ed3a29a", "S": "0x6748abe186269ff10b8100a4b7d7fea274b53ea2905acbf498dc8b5ab1bf4fbc"}]
			}`,
	},
	{
		Name: "09_FeeDelegatedValueTransfer", // kaia-sdk
		Type: TxTypeFeeDelegatedValueTransfer,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer: common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x25, r: "0x9f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956", s: "0x6bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0xf45cf8d7f88c08e6b6ec0b3b562f34ca94283e4689021987abb6b0772ddfd80a", s: "0x298fe2c5aeabb6a518f4cbb5ff39631a5d88be505d3923374f65fdcf63c2955b"}},
		SigRLP:          "0xf839b5f4098204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b018080",
		SigFeePayerRLP:  "0xf84eb5f4098204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x09f87a8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a09f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956a06bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7",
		TxHashRLP:       "0x09f8d68204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a09f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956a06bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0f45cf8d7f88c08e6b6ec0b3b562f34ca94283e4689021987abb6b0772ddfd80aa0298fe2c5aeabb6a518f4cbb5ff39631a5d88be505d3923374f65fdcf63c2955b",
	},
}

func TestTransactionSerialization2(t *testing.T) {
	for _, tc := range serializeTCs {
		t.Run(tc.Name, func(t *testing.T) {
			// Create the transaction object
			txData, err := NewTxInternalDataWithMap(tc.Type, tc.Map)
			require.Nil(t, err)

			txData.SetSignature(tc.SenderSigs.TxSignatures())
			if txDataFeePayer, ok := txData.(TxInternalDataFeePayer); ok {
				txDataFeePayer.SetFeePayerSignatures(tc.FeePayerSigs.TxSignatures())
			}

			// Tx -> SigRLP, TxRLP
			sigRLP, err := getSigRLP(tc.Type, txData, tc.ChainID)
			require.Nil(t, err)
			assert.Equal(t, tc.SigRLP, hexutil.Encode(sigRLP), "SigRLP")

			txRLP, err := getTxHashRLP(txData)
			require.Nil(t, err)
			assert.Equal(t, tc.TxHashRLP, hexutil.Encode(txRLP), "TxRLP")

			// Tx -> SigFeePayerRLP, SenderTxHashRLP (only for fee-delegated txs)
			if tc.Type.IsFeeDelegatedTransaction() {
				sigFeePayerRLP, err := getSigFeePayerRLP(tc.Type, txData, tc.ChainID)
				require.Nil(t, err)
				assert.Equal(t, tc.SigFeePayerRLP, hexutil.Encode(sigFeePayerRLP), "SigFeePayerRLP")

				// since TxInternalData only has SenderTxHash(), we compare the hash instead.
				senderTxHash := txData.SenderTxHash()
				expected := crypto.Keccak256Hash(hexutil.MustDecode(tc.SenderTxHashRLP))
				assert.Equal(t, expected.Hex(), senderTxHash.Hex(), "SenderTxHash")
			}

			// TxRLP -> Tx -> TxRLP
			decRLP := newTxInternalDataSerializer()
			require.Nil(t, rlp.DecodeBytes(txRLP, decRLP))
			txRLP2, err := getTxHashRLP(decRLP.tx)
			require.Nil(t, err)
			assert.Equal(t, txRLP, txRLP2, "TxRLP round trip")

			if len(tc.TxJson) > 0 {
				// Tx -> Json
				txJson, err := json.Marshal(txData)
				require.Nil(t, err)
				assert.JSONEq(t, tc.TxJson, string(txJson), "Json")

				// Json -> Tx -> Json
				decJson := newTxInternalDataSerializer()
				require.Nil(t, json.Unmarshal([]byte(tc.TxJson), decJson))
				txJson2, err := json.Marshal(decJson.tx)
				require.Nil(t, err)
				assert.JSONEq(t, tc.TxJson, string(txJson2), "Json round trip")
			}

			if len(tc.RpcJson) > 0 {
				// Tx -> RpcJson
				rpcFields := txData.MakeRPCOutput()
				rpcJson, err := json.Marshal(rpcFields)
				require.Nil(t, err)
				assert.JSONEq(t, tc.RpcJson, string(rpcJson), "RpcJson")
			}
		})
	}
}

// TODO-Kaia: move RLP helpers to Transaction type. Signers shall use the helpers.
func getSigRLP(txType TxType, txData TxInternalData, chainID uint64) ([]byte, error) {
	if txType.IsLegacyTransaction() {
		return getSigRLPLegacy(txType, txData, chainID)
	} else if txType.IsEthTypedTransaction() {
		return getSigRLPEth(txType, txData, chainID)
	} else {
		return getSigRLPKaia(txData, chainID)
	}
}

// SigRLPEth = type + RLP(x.SerializeForSign()...).
// Applies to Eth typed txs (1,2,3,...). See eip2930Signer.Hash() and EIP-2718.
func getSigRLPEth(txType TxType, txData TxInternalData, chainID uint64) ([]byte, error) {
	elems := txData.SerializeForSign()
	// elems[0] is always chainID for Eth typed txs.
	// Fill in the chainID if it was missing.
	if txData.ChainId() == nil || txData.ChainId().Sign() == 0 {
		elems[0] = chainID
	}

	encoded, err := rlp.EncodeToBytes(elems)
	if err != nil {
		return nil, err
	}
	return append([]byte{byte(txType)}, encoded...), nil
}

// SigRLPLegacy = RLP([ tx.SerializeForSign()..., ChainID, 0, 0 ]).
// Applies to Legacy tx (type 0). See EIP155Signer.Hash() and EIP-2718.
func getSigRLPLegacy(txType TxType, txData TxInternalData, chainID uint64) ([]byte, error) {
	return rlp.EncodeToBytes(append(txData.SerializeForSign(), chainID, uint(0), uint(0)))
}

// SigRLPKaia = RLP([ RLP(tx.SerializeForSign()), ChainID, 0, 0 ])
// Applies to Kaia typed txs (8,9,10...). See EIP155Signer.Hash().
func getSigRLPKaia(txData TxInternalData, chainID uint64) ([]byte, error) {
	innerElems := txData.SerializeForSign()
	innerRLP, err := rlp.EncodeToBytes(innerElems)
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes([]interface{}{
		innerRLP,
		chainID,
		uint(0),
		uint(0),
	})
}

// SigFeePayerRLP = RLP([ RLP(tx.SerializeForSign()), feePayer, ChainID, 0, 0 ]).
// Applies to fee-delegated txs. See EIP155Signer.HashFeePayer().
func getSigFeePayerRLP(txType TxType, txData TxInternalData, chainID uint64) ([]byte, error) {
	t, ok := txData.(TxInternalDataFeePayer)
	if !ok {
		return nil, errNotFeeDelegationTransaction
	}
	innerElems := txData.SerializeForSign()
	innerRLP, err := rlp.EncodeToBytes(innerElems)
	if err != nil {
		return nil, err
	}
	return rlp.EncodeToBytes([]interface{}{
		innerRLP,
		t.GetFeePayer(),
		chainID,
		uint(0),
		uint(0),
	})
}

// TxHashRLP = type + RLP(x)
func getTxHashRLP(txData TxInternalData) ([]byte, error) {
	ser := newTxInternalDataSerializerWithValues(txData)
	return rlp.EncodeToBytes(ser)
}
