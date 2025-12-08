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
	"crypto/ecdsa"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	sidecarV0      = NewBlobTxSidecar(BlobSidecarVersion0, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, []kzg4844.Proof{emptyBlobProof})
	sidecarV1      = NewBlobTxSidecar(BlobSidecarVersion1, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, []kzg4844.Proof{emptyBlobProof})
	authorizations = []SetCodeAuthorization{{ChainID: *uint256.NewInt(2), Address: common.HexToAddress("0x0000000000000000000000000000000000000001"), Nonce: nonce, V: uint8(0), R: *uint256.NewInt(0), S: *uint256.NewInt(0)}}
)

// TestTransactionSerialization tests the serialization of the transaction.
func TestTransactionSerialization(t *testing.T) {
	for _, tc := range serializeTCs {
		t.Run(tc.Name, func(t *testing.T) {
			// Create the transaction object
			txData, err := NewTxInternalDataWithMap(tc.Type, tc.Map)
			require.Nil(t, err)

			txData.SetSignature(tc.SenderSigs.TxSignatures())
			if txDataFeePayer, ok := txData.(TxInternalDataFeePayer); ok {
				txDataFeePayer.SetFeePayerSignatures(tc.FeePayerSigs.TxSignatures())
			}

			// Tx -> SigHash, TxHashRLP
			expectedSigHash := crypto.Keccak256Hash(hexutil.MustDecode(tc.SigRLP))
			assert.Equal(t, expectedSigHash, txData.SigHash(big.NewInt(int64(tc.ChainID))), "SigHash")

			ser := newTxInternalDataSerializerWithValues(txData)
			txHashRLP, err := rlp.EncodeToBytes(ser)
			require.Nil(t, err)
			assert.Equal(t, tc.TxHashRLP, "0x"+hex.EncodeToString(txHashRLP), "TxHash")

			// Tx -> FeePayerSigHash, SenderTxHash (only for fee-delegated txs)
			if txDataFeePayer, ok := txData.(TxInternalDataFeePayer); ok {
				expectedSigFeePayerHash := crypto.Keccak256Hash(hexutil.MustDecode(tc.FeePayerSigRLP))
				assert.Equal(t, expectedSigFeePayerHash, txDataFeePayer.FeePayerSigHash(big.NewInt(int64(tc.ChainID))), "FeePayerSigHash")

				expectedSenderTxHash := crypto.Keccak256Hash(hexutil.MustDecode(tc.SenderTxHashRLP))
				assert.Equal(t, expectedSenderTxHash, txDataFeePayer.SenderTxHash(), "SenderTxHash")
			}

			// TxHashRLP -> Tx -> TxHashRLP
			decRLP := newTxInternalDataSerializer()
			require.Nil(t, rlp.DecodeBytes(hexutil.MustDecode(tc.TxHashRLP), decRLP))
			txHashRLP2, err := rlp.EncodeToBytes(decRLP)
			require.Nil(t, err)
			assert.Equal(t, tc.TxHashRLP, "0x"+hex.EncodeToString(txHashRLP2), "TxHash round trip")

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

			// Tx -> RpcJson
			rpcFields := txData.MakeRPCOutput()
			rpcJson, err := json.Marshal(rpcFields)
			require.Nil(t, err)
			assert.JSONEq(t, tc.RpcJson, string(rpcJson), "RpcJson")
		})
	}
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
	FeePayerSigRLP  string
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

func hexToECDSAPublicKey(hexkey string) *ecdsa.PublicKey {
	// decode hex to bytes
	pubBytes := hexutil.MustDecode(hexkey)

	pubkey, _ := crypto.DecompressPubkey(pubBytes)
	return pubkey
}

// Some cases taken from kaia-sdk: https://github.com/kaiachain/kaia-sdk/blob/dev/js-ext-core/test/tx.spec.ts
// Some cases taken from viem: https://github.com/wevm/viem/blob/main/vectors/src/transaction.json.gz
//   - Viem vectors can be filtered with (e.g.):
//     o = JSON.parse(fs.readFileSync('transaction.json'))
//     o.filter((x) => x.name.substr(0,6) == "legacy").filter((x) => { l = x.transaction.data?.length; return (l > 0 && l < 100) }).filter((x) => { t = x.transaction; return (t.data && t.to && t.gas && t.gasPrice && t.nonce && t.value) })
//   - Viem vectors don't contain SigRLP. Correct SigRLPs were manually generated using (e.g.):
//     rlp encode '["0x24ae","0x07a3d3","0x3cd980f3","0x1cd12f7edecac097265d53754b782004bf0b8fb7","0x7aa930ce8a1a","0xe02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e","0x01","",""]'
//   - Viem vectors are manually prefixed with EthereumTxTypeEnvelope(0x78).
//
// Eip7702 case taken from eest blockchain_test
//   - eest test can be filtered with the way like vectors
//
// Related to ChainDataAnchoring test cases are manually generated
// TxTypeAccountCreation, TxTypeBatch are not supported.
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
		TxJson: `{
			"nonce": "0x24ae",
			"gasPrice": "0x7a3d3",
			"gas": "0x3cd980f3",
			"to": "0x1cd12f7edecac097265d53754b782004bf0b8fb7",
			"value": "0x7aa930ce8a1a",
			"input": "0xe02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e",
			"signatures": [{"V": "0x1c", "R": "0x8c3d94705e12d605d1144ac4c29bfcde87a06cf8424f2addddbc64668e9d78ba", "S": "0x233281c375ddedb00ce1f5fb167f80e690701db29f77f1c372cbe72895147299"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 0,
			"type": "TxTypeLegacyTransaction",
			"gas": "0x3cd980f3",
			"gasPrice": "0x7a3d3",
			"input": "0xe02ea655070ac8dce2e299bb782e344c55b17755c8c1e70e",
			"nonce": "0x24ae",
			"signatures": [{"V": "0x1c", "R": "0x8c3d94705e12d605d1144ac4c29bfcde87a06cf8424f2addddbc64668e9d78ba", "S": "0x233281c375ddedb00ce1f5fb167f80e690701db29f77f1c372cbe72895147299"}],
			"to": "0x1cd12f7edecac097265d53754b782004bf0b8fb7",
			"value": "0x7aa930ce8a1a"
		}`,
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
		TxJson: `{
			"typeInt": 30721,
			"type": "TxTypeEthereumAccessList",
			"chainId": "0x9",
			"nonce": "0x4275",
			"gasPrice": "0x0",
			"gas": "0xbf7f714d",
			"to": "0x610590408a97e8f84152d45cbe2ec0b08059598f",
			"value": "0x0",
			"input": "0xd962f328314a9cd384c349a461fdba53ead623afb4f751df94a604",
			"accessList": [{"address": "0x67d560ba27b75d467fcf658e02ac3765c3634056", "storageKeys": ["0xa12c9f0600000000000000000000000000000000000000000000000000000000", "0x5b00000000000000000000000000000000000000000000000000000000000000", "0x082b41eb5f7244029351c2533948000000000000000000000000000000000000"]}],
			"signatures": [{"V": "0x1", "R": "0xc71988016460ca44c352598cc116cf2ac1b9293387f2f8287aa889b6cf55b1e7", "S": "0x634101cdd4f96d642fe8509166cc63729722cb51c80f1624b4fe693106b0664c"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 30721,
			"type": "TxTypeEthereumAccessList",
			"chainId": "0x9",
			"gas": "0xbf7f714d",
			"gasPrice": "0x0",
			"input": "0xd962f328314a9cd384c349a461fdba53ead623afb4f751df94a604",
			"nonce": "0x4275",
			"signatures": [{"V": "0x1", "R": "0xc71988016460ca44c352598cc116cf2ac1b9293387f2f8287aa889b6cf55b1e7", "S": "0x634101cdd4f96d642fe8509166cc63729722cb51c80f1624b4fe693106b0664c"}],
			"to": "0x610590408a97e8f84152d45cbe2ec0b08059598f",
			"value": "0x0",
			"accessList": [{"address": "0x67d560ba27b75d467fcf658e02ac3765c3634056", "storageKeys": ["0xa12c9f0600000000000000000000000000000000000000000000000000000000", "0x5b00000000000000000000000000000000000000000000000000000000000000", "0x082b41eb5f7244029351c2533948000000000000000000000000000000000000"]}]
		}`,
	},
	{
		Name: "02_DynamicFee", // Viem "eip1559: 8386"
		Type: TxTypeEthereumDynamicFee,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:   big.NewInt(9),
			TxValueKeyNonce:     uint64(9778),
			TxValueKeyGasFeeCap: big.NewInt(1063),
			TxValueKeyGasTipCap: big.NewInt(0),
			TxValueKeyGasLimit:  uint64(0),
			TxValueKeyTo:        addrPtr(common.HexToAddress("0x529df59d69f3f13fc94af90af33e4f794d76a929")),
			TxValueKeyAmount:    new(big.Int).SetUint64(13922405099808277777),
			TxValueKeyData:      hexutil.MustDecode("0x3dde5092cfa391fc32d21999924a09e1f8b7f326dfdb312a423ca561"),
			TxValueKeyAccessList: AccessList{
				{
					Address: common.HexToAddress("0xee235c66c42e7e595890f50269770e80edfb198a"),
					StorageKeys: []common.Hash{
						common.HexToHash("0x70a80801a9955fad0a39bb6a7fd4000000000000000000000000000000000000"),
						common.HexToHash("0x657d828bf2a683f027ec06d606d8d40000000000000000000000000000000000"),
						common.HexToHash("0xf4e6d428979de9a52b6f73000000000000000000000000000000000000000000"),
					},
				},
			},
		},
		ChainID:    9,
		SenderSigs: []txSigHex{{v: 28, r: "0xd154b702f73f724567492930d9bf832cfa93d7b060e5da2ccfeac7f88ba502dd", s: "0x7330827cc70f8fb4024c1376b39f64cc8ae1430b7ccb178f258de88ba492dd11"}},
		SigRLP:     "0x02f8c209822632808204278094529df59d69f3f13fc94af90af33e4f794d76a92988c13651ac9928c5119c3dde5092cfa391fc32d21999924a09e1f8b7f326dfdb312a423ca561f87cf87a94ee235c66c42e7e595890f50269770e80edfb198af863a070a80801a9955fad0a39bb6a7fd4000000000000000000000000000000000000a0657d828bf2a683f027ec06d606d8d40000000000000000000000000000000000a0f4e6d428979de9a52b6f73000000000000000000000000000000000000000000",
		TxHashRLP:  "0x7802f9010509822632808204278094529df59d69f3f13fc94af90af33e4f794d76a92988c13651ac9928c5119c3dde5092cfa391fc32d21999924a09e1f8b7f326dfdb312a423ca561f87cf87a94ee235c66c42e7e595890f50269770e80edfb198af863a070a80801a9955fad0a39bb6a7fd4000000000000000000000000000000000000a0657d828bf2a683f027ec06d606d8d40000000000000000000000000000000000a0f4e6d428979de9a52b6f730000000000000000000000000000000000000000001ca0d154b702f73f724567492930d9bf832cfa93d7b060e5da2ccfeac7f88ba502dda07330827cc70f8fb4024c1376b39f64cc8ae1430b7ccb178f258de88ba492dd11",
		TxJson: `{
			"typeInt": 30722,
			"type": "TxTypeEthereumDynamicFee",
			"chainId": "0x9",
			"nonce": "0x2632",
			"maxPriorityFeePerGas": "0x0",
			"maxFeePerGas": "0x427",
			"gas": "0x0",
			"to": "0x529df59d69f3f13fc94af90af33e4f794d76a929",
			"input": "0x3dde5092cfa391fc32d21999924a09e1f8b7f326dfdb312a423ca561",
			"value": "0xc13651ac9928c511",
			"accessList": [{"address": "0xee235c66c42e7e595890f50269770e80edfb198a", "storageKeys": ["0x70a80801a9955fad0a39bb6a7fd4000000000000000000000000000000000000", "0x657d828bf2a683f027ec06d606d8d40000000000000000000000000000000000", "0xf4e6d428979de9a52b6f73000000000000000000000000000000000000000000"]}],
			"signatures": [{"V": "0x1c", "R": "0xd154b702f73f724567492930d9bf832cfa93d7b060e5da2ccfeac7f88ba502dd", "S": "0x7330827cc70f8fb4024c1376b39f64cc8ae1430b7ccb178f258de88ba492dd11"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 30722,
			"type": "TxTypeEthereumDynamicFee",
			"chainId": "0x9",
			"nonce": "0x2632",
			"maxPriorityFeePerGas": "0x0",
			"maxFeePerGas": "0x427",
			"gas": "0x0",
			"to": "0x529df59d69f3f13fc94af90af33e4f794d76a929",
			"input": "0x3dde5092cfa391fc32d21999924a09e1f8b7f326dfdb312a423ca561",
			"value": "0xc13651ac9928c511",
			"accessList": [{"address": "0xee235c66c42e7e595890f50269770e80edfb198a", "storageKeys": ["0x70a80801a9955fad0a39bb6a7fd4000000000000000000000000000000000000", "0x657d828bf2a683f027ec06d606d8d40000000000000000000000000000000000", "0xf4e6d428979de9a52b6f73000000000000000000000000000000000000000000"]}],
			"signatures": [{"V": "0x1c", "R": "0xd154b702f73f724567492930d9bf832cfa93d7b060e5da2ccfeac7f88ba502dd", "S": "0x7330827cc70f8fb4024c1376b39f64cc8ae1430b7ccb178f258de88ba492dd11"}]
		}`,
	},
	{
		Name: "03_Blob", // Viem "eip4844: 6457"
		Type: TxTypeEthereumBlob,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:   big.NewInt(6),
			TxValueKeyNonce:     uint64(6294),
			TxValueKeyTo:        common.HexToAddress("0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1"),
			TxValueKeyData:      hexutil.MustDecode("0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0"),
			TxValueKeyAmount:    big.NewInt(0),
			TxValueKeyGasFeeCap: big.NewInt(0),
			TxValueKeyGasTipCap: big.NewInt(0),
			TxValueKeyGasLimit:  uint64(0),
			TxValueKeyAccessList: AccessList{
				{
					Address: common.HexToAddress("0x6092415c41b602d192c02d8bb5b2ee62fbab3b70"),
					StorageKeys: []common.Hash{
						common.HexToHash("0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"),
					},
				},
			},
			TxValueKeyBlobFeeCap: new(big.Int).SetUint64(17602539720540508054),
			TxValueKeyBlobHashes: []common.Hash{
				common.HexToHash("0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000"),
				common.HexToHash("0x01f263630289db00000000000000000000000000000000000000000000000000"),
				common.HexToHash("0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"),
			},
			TxValueKeySidecar: (*BlobTxSidecar)(nil),
		},
		ChainID:    6,
		SenderSigs: []txSigHex{{v: 27, r: "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", s: "0x0def24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}},
		SigRLP:     "0x03f8e30682189680808094df3ca4eaf9017d01a26ef475e651faa9b1296da1809d09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0f838f7946092415c41b602d192c02d8bb5b2ee62fbab3b70e1a0a2c53cdc4de0f875229c19c1d05f5f000000000000000000000000000000000088f448c89d13854f96f863a0012730cf6ab975c7c39a00000000000000000000000000000000000000000000a001f263630289db00000000000000000000000000000000000000000000000000a00159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000",
		TxHashRLP:  "0x7803f901260682189680808094df3ca4eaf9017d01a26ef475e651faa9b1296da1809d09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0f838f7946092415c41b602d192c02d8bb5b2ee62fbab3b70e1a0a2c53cdc4de0f875229c19c1d05f5f000000000000000000000000000000000088f448c89d13854f96f863a0012730cf6ab975c7c39a00000000000000000000000000000000000000000000a001f263630289db00000000000000000000000000000000000000000000000000a00159b494f64c2c6adac876c9d3ea38f46c9aca7d8693000000000000000000001ba0f529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2ba00def24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746",
		TxJson: `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": null,
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": null,
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}]
		}`,
	},
	{
		Name: "03_Blob with sidecar v0", // Viem "eip4844: 6457 with sidecar v0"
		Type: TxTypeEthereumBlob,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:   big.NewInt(6),
			TxValueKeyNonce:     uint64(6294),
			TxValueKeyTo:        common.HexToAddress("0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1"),
			TxValueKeyData:      hexutil.MustDecode("0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0"),
			TxValueKeyAmount:    big.NewInt(0),
			TxValueKeyGasFeeCap: big.NewInt(0),
			TxValueKeyGasTipCap: big.NewInt(0),
			TxValueKeyGasLimit:  uint64(0),
			TxValueKeyAccessList: AccessList{
				{
					Address: common.HexToAddress("0x6092415c41b602d192c02d8bb5b2ee62fbab3b70"),
					StorageKeys: []common.Hash{
						common.HexToHash("0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"),
					},
				},
			},
			TxValueKeyBlobFeeCap: new(big.Int).SetUint64(17602539720540508054),
			TxValueKeyBlobHashes: []common.Hash{
				common.HexToHash("0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000"),
				common.HexToHash("0x01f263630289db00000000000000000000000000000000000000000000000000"),
				common.HexToHash("0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"),
			},
			TxValueKeySidecar: sidecarV0,
		},
		ChainID:    6,
		SenderSigs: []txSigHex{{v: 27, r: "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", s: "0x0def24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}},
		SigRLP:     "0x03f8e30682189680808094df3ca4eaf9017d01a26ef475e651faa9b1296da1809d09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0f838f7946092415c41b602d192c02d8bb5b2ee62fbab3b70e1a0a2c53cdc4de0f875229c19c1d05f5f000000000000000000000000000000000088f448c89d13854f96f863a0012730cf6ab975c7c39a00000000000000000000000000000000000000000000a001f263630289db00000000000000000000000000000000000000000000000000a00159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000",
		TxHashRLP:  blobtxWithSidecarV0HashRLP,
		TxJson:     fmt.Sprintf(blobtxWithSidecarV0TxJson, blob0),
		RpcJson:    fmt.Sprintf(blobtxWithSidecarV0RPCJson, blob0),
	},
	{
		Name: "03_Blob with sidecar v1", // Viem "eip4844: 6457 with sidecar v1"
		Type: TxTypeEthereumBlob,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:   big.NewInt(6),
			TxValueKeyNonce:     uint64(6294),
			TxValueKeyTo:        common.HexToAddress("0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1"),
			TxValueKeyData:      hexutil.MustDecode("0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0"),
			TxValueKeyAmount:    big.NewInt(0),
			TxValueKeyGasFeeCap: big.NewInt(0),
			TxValueKeyGasTipCap: big.NewInt(0),
			TxValueKeyGasLimit:  uint64(0),
			TxValueKeyAccessList: AccessList{
				{
					Address: common.HexToAddress("0x6092415c41b602d192c02d8bb5b2ee62fbab3b70"),
					StorageKeys: []common.Hash{
						common.HexToHash("0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"),
					},
				},
			},
			TxValueKeyBlobFeeCap: new(big.Int).SetUint64(17602539720540508054),
			TxValueKeyBlobHashes: []common.Hash{
				common.HexToHash("0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000"),
				common.HexToHash("0x01f263630289db00000000000000000000000000000000000000000000000000"),
				common.HexToHash("0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"),
			},
			TxValueKeySidecar: sidecarV1,
		},
		ChainID:    6,
		SenderSigs: []txSigHex{{v: 27, r: "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", s: "0x0def24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}},
		SigRLP:     "0x03f8e30682189680808094df3ca4eaf9017d01a26ef475e651faa9b1296da1809d09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0f838f7946092415c41b602d192c02d8bb5b2ee62fbab3b70e1a0a2c53cdc4de0f875229c19c1d05f5f000000000000000000000000000000000088f448c89d13854f96f863a0012730cf6ab975c7c39a00000000000000000000000000000000000000000000a001f263630289db00000000000000000000000000000000000000000000000000a00159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000",
		TxHashRLP:  blobtxWithSidecarV1HashRLP,
		TxJson:     fmt.Sprintf(blobtxWithSidecarV1TxJson, blob0),
		RpcJson:    fmt.Sprintf(blobtxWithSidecarV1RPCJson, blob0),
	},
	{
		Name: "04_SetCode", // EEST "test_set_code_to_system_contract[fork_Osaka-call_opcode_CALL-evm_code_type_LEGACY-system_contract_0x000f3df6d732807ef1319fb7b8bb8522d0beac02-blockchain_test]"
		Type: TxTypeEthereumSetCode,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyChainID:    big.NewInt(1),
			TxValueKeyNonce:      uint64(0),
			TxValueKeyGasFeeCap:  big.NewInt(7),
			TxValueKeyGasTipCap:  big.NewInt(0),
			TxValueKeyGasLimit:   uint64(500000),
			TxValueKeyTo:         common.HexToAddress("0x891dce8073514e62bb564c41f75df3062c381573"),
			TxValueKeyAmount:     big.NewInt(0),
			TxValueKeyData:       hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
			TxValueKeyAccessList: AccessList{},
			TxValueKeyAuthorizationList: []SetCodeAuthorization{
				{
					ChainID: *uint256.NewInt(1),
					Address: common.HexToAddress("0x000f3df6d732807ef1319fb7b8bb8522d0beac02"),
					Nonce:   0,
					V:       1,
					R:       *uint256.MustFromBig(hexutil.MustDecodeBig("0xd14d99564d380653121fa874e9e44f25a92d23c03a9e0dbb43e0e7e1e0995847")),
					S:       *uint256.MustFromBig(hexutil.MustDecodeBig("0x492693b6721976ffca5c44674720163a2e36bcd06d8336408496ac9b96f92380")),
				},
			},
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 1, r: "0xdcde3e1249cfff2b0593bf178f683f1b682fe7c8700310067a819b89b8363e52", s: "0x663694ed951ecdba3a9c3c652ce093d12d192a3acbd544267aa70e7004021bb2"}},
		SigRLP:     "0x04f89e018080078307a12094891dce8073514e62bb564c41f75df3062c38157380a00000000000000000000000000000000000000000000000000000000000000001c0f85cf85a0194000f3df6d732807ef1319fb7b8bb8522d0beac028001a0d14d99564d380653121fa874e9e44f25a92d23c03a9e0dbb43e0e7e1e0995847a0492693b6721976ffca5c44674720163a2e36bcd06d8336408496ac9b96f92380",
		TxHashRLP:  "0x7804f8e1018080078307a12094891dce8073514e62bb564c41f75df3062c38157380a00000000000000000000000000000000000000000000000000000000000000001c0f85cf85a0194000f3df6d732807ef1319fb7b8bb8522d0beac028001a0d14d99564d380653121fa874e9e44f25a92d23c03a9e0dbb43e0e7e1e0995847a0492693b6721976ffca5c44674720163a2e36bcd06d8336408496ac9b96f9238001a0dcde3e1249cfff2b0593bf178f683f1b682fe7c8700310067a819b89b8363e52a0663694ed951ecdba3a9c3c652ce093d12d192a3acbd544267aa70e7004021bb2",
		TxJson: `{
			"typeInt": 30724,
			"type": "TxTypeEthereumSetCode",
			"chainId": "0x1",
			"nonce": "0x0",
			"maxFeePerGas": "0x7",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x7a120",
			"to": "0x891dce8073514e62bb564c41f75df3062c381573",
			"value": "0x0",
			"input": "0x0000000000000000000000000000000000000000000000000000000000000001",
			"accessList": [],
			"authorizationList": [{"chainId": "0x1", "address": "0x000f3df6d732807ef1319fb7b8bb8522d0beac02", "nonce": "0x0", "yParity": "0x1", "r": "0xd14d99564d380653121fa874e9e44f25a92d23c03a9e0dbb43e0e7e1e0995847", "s": "0x492693b6721976ffca5c44674720163a2e36bcd06d8336408496ac9b96f92380"}],
			"signatures": [{"V": "0x1", "R": "0xdcde3e1249cfff2b0593bf178f683f1b682fe7c8700310067a819b89b8363e52", "S": "0x663694ed951ecdba3a9c3c652ce093d12d192a3acbd544267aa70e7004021bb2"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 30724,
			"type": "TxTypeEthereumSetCode",
			"chainId": "0x1",
			"nonce": "0x0",
			"maxFeePerGas": "0x7",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x7a120",
			"to": "0x891dce8073514e62bb564c41f75df3062c381573",
			"value": "0x0",
			"input": "0x0000000000000000000000000000000000000000000000000000000000000001",
			"accessList": [],
			"authorizationList": [{"chainId": "0x1", "address": "0x000f3df6d732807ef1319fb7b8bb8522d0beac02", "nonce": "0x0", "yParity": "0x1", "r": "0xd14d99564d380653121fa874e9e44f25a92d23c03a9e0dbb43e0e7e1e0995847", "s": "0x492693b6721976ffca5c44674720163a2e36bcd06d8336408496ac9b96f92380"}],
			"signatures": [{"V": "0x1", "R": "0xdcde3e1249cfff2b0593bf178f683f1b682fe7c8700310067a819b89b8363e52", "S": "0x663694ed951ecdba3a9c3c652ce093d12d192a3acbd544267aa70e7004021bb2"}]
		}`,
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
		FeePayerSigRLP:  "0xf84eb5f4098204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x09f87a8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a09f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956a06bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7",
		TxHashRLP:       "0x09f8d68204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a09f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956a06bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0f45cf8d7f88c08e6b6ec0b3b562f34ca94283e4689021987abb6b0772ddfd80aa0298fe2c5aeabb6a518f4cbb5ff39631a5d88be505d3923374f65fdcf63c2955b",
		TxJson: `{
			"typeInt": 9,
			"type": "TxTypeFeeDelegatedValueTransfer",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"signatures": [{"V": "0x25", "R": "0x9f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956", "S": "0x6bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xf45cf8d7f88c08e6b6ec0b3b562f34ca94283e4689021987abb6b0772ddfd80a", "S": "0x298fe2c5aeabb6a518f4cbb5ff39631a5d88be505d3923374f65fdcf63c2955b"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 9,
			"type": "TxTypeFeeDelegatedValueTransfer",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"signatures": [{"V": "0x25", "R": "0x9f8e49e2ad84b0732984398749956e807e4b526c786af3c5f7416b293e638956", "S": "0x6bf88342092f6ff9fabe31739b2ebfa1409707ce54a54693e91a6b9bb77df0e7"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xf45cf8d7f88c08e6b6ec0b3b562f34ca94283e4689021987abb6b0772ddfd80a", "S": "0x298fe2c5aeabb6a518f4cbb5ff39631a5d88be505d3923374f65fdcf63c2955b"}]
		}`,
	},
	{
		Name: "0a_FeeDelegatedValueTransferWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedValueTransferWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(0x19),
			TxValueKeyGasLimit:           uint64(0xf4240),
			TxValueKeyTo:                 common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:             big.NewInt(10),
			TxValueKeyFrom:               common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x25, r: "0xdde32b8241f039a82b124fe94d3e556eb08f0d6f26d07dcc0f3fca621f1090ca", s: "0x1c8c336b358ab6d3a2bbf25de2adab4d01b754e2fb3b9b710069177d54c1e956"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0x091ecf53f91bb97bb694f2f2443f3563ac2b646d651497774524394aae396360", s: "0x44228b88f275aa1ec1bab43681d21dc7e3a676786ed1906f6841d0a1a188f88a"}},
		SigRLP:          "0xf83ab6f50a8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b1e018080",
		FeePayerSigRLP:  "0xf84fb6f50a8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b1e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x0af87b8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b1ef845f84325a0dde32b8241f039a82b124fe94d3e556eb08f0d6f26d07dcc0f3fca621f1090caa01c8c336b358ab6d3a2bbf25de2adab4d01b754e2fb3b9b710069177d54c1e956",
		TxHashRLP:       "0x0af8d78204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b1ef845f84325a0dde32b8241f039a82b124fe94d3e556eb08f0d6f26d07dcc0f3fca621f1090caa01c8c336b358ab6d3a2bbf25de2adab4d01b754e2fb3b9b710069177d54c1e956945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0091ecf53f91bb97bb694f2f2443f3563ac2b646d651497774524394aae396360a044228b88f275aa1ec1bab43681d21dc7e3a676786ed1906f6841d0a1a188f88a",
		TxJson: `{
			"typeInt": 10,
			"type": "TxTypeFeeDelegatedValueTransferWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x25", "R": "0xdde32b8241f039a82b124fe94d3e556eb08f0d6f26d07dcc0f3fca621f1090ca", "S": "0x1c8c336b358ab6d3a2bbf25de2adab4d01b754e2fb3b9b710069177d54c1e956"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x91ecf53f91bb97bb694f2f2443f3563ac2b646d651497774524394aae396360", "S": "0x44228b88f275aa1ec1bab43681d21dc7e3a676786ed1906f6841d0a1a188f88a"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 10,
			"type": "TxTypeFeeDelegatedValueTransferWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x25", "R": "0xdde32b8241f039a82b124fe94d3e556eb08f0d6f26d07dcc0f3fca621f1090ca", "S": "0x1c8c336b358ab6d3a2bbf25de2adab4d01b754e2fb3b9b710069177d54c1e956"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x91ecf53f91bb97bb694f2f2443f3563ac2b646d651497774524394aae396360", "S": "0x44228b88f275aa1ec1bab43681d21dc7e3a676786ed1906f6841d0a1a188f88a"}]
		}`,
	},
	{
		Name: "10_ValueTransferMemo", // kaia-sdk
		Type: TxTypeValueTransferMemo,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyData:     hexutil.MustDecode("0x68656c6c6f"),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 0x25, r: "0x7d2b0c89ee8afa502b3186413983bfe9a31c5776f4f820210cffe44a7d568d1c", s: "0x2b1cbd587c73b0f54969f6b76ef2fd95cea0c1bb79256a75df9da696278509f3"}},
		SigRLP:     "0xf841b83cf83a108204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f018080",
		TxHashRLP:  "0x10f8808204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6ff845f84325a07d2b0c89ee8afa502b3186413983bfe9a31c5776f4f820210cffe44a7d568d1ca02b1cbd587c73b0f54969f6b76ef2fd95cea0c1bb79256a75df9da696278509f3",
		TxJson: `{
			"typeInt": 16,
			"type": "TxTypeValueTransferMemo",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x68656c6c6f",
			"signatures": [{"V": "0x25", "R": "0x7d2b0c89ee8afa502b3186413983bfe9a31c5776f4f820210cffe44a7d568d1c", "S": "0x2b1cbd587c73b0f54969f6b76ef2fd95cea0c1bb79256a75df9da696278509f3"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 16,
			"type": "TxTypeValueTransferMemo",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x68656c6c6f",
			"signatures": [{"V": "0x25", "R": "0x7d2b0c89ee8afa502b3186413983bfe9a31c5776f4f820210cffe44a7d568d1c", "S": "0x2b1cbd587c73b0f54969f6b76ef2fd95cea0c1bb79256a75df9da696278509f3"}]
		}`,
	},
	{
		Name: "11_FeeDelegatedValueTransferMemo", // kaia-sdk
		Type: TxTypeFeeDelegatedValueTransferMemo,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer: common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyData:     hexutil.MustDecode("0x68656c6c6f"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x64e213aef0167fbd853f8f9989ef5d8b912a77457395ccf13d7f37009edd5c5b", s: "0x5d0c2e55e4d8734fe2516ed56ac628b74c0eb02aa3b6eda51e1e25a1396093e1"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0x87390ac14d3c34440b6ddb7b190d3ebde1a07d9a556e5a82ce7e501f24a060f9", s: "0x37badbcb12cda1ed67b12b1831683a08a3adadee2ea760a07a46bdbb856fea44"}},
		SigRLP:          "0xf841b83cf83a118204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f018080",
		FeePayerSigRLP:  "0xf856b83cf83a118204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x11f8808204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6ff845f84326a064e213aef0167fbd853f8f9989ef5d8b912a77457395ccf13d7f37009edd5c5ba05d0c2e55e4d8734fe2516ed56ac628b74c0eb02aa3b6eda51e1e25a1396093e1",
		TxHashRLP:       "0x11f8dc8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6ff845f84326a064e213aef0167fbd853f8f9989ef5d8b912a77457395ccf13d7f37009edd5c5ba05d0c2e55e4d8734fe2516ed56ac628b74c0eb02aa3b6eda51e1e25a1396093e1945a0043070275d9f6054307ee7348bd660849d90ff845f84326a087390ac14d3c34440b6ddb7b190d3ebde1a07d9a556e5a82ce7e501f24a060f9a037badbcb12cda1ed67b12b1831683a08a3adadee2ea760a07a46bdbb856fea44",
		TxJson: `{
			"typeInt": 17,
			"type": "TxTypeFeeDelegatedValueTransferMemo",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x68656c6c6f",
			"signatures": [{"V": "0x26", "R": "0x64e213aef0167fbd853f8f9989ef5d8b912a77457395ccf13d7f37009edd5c5b", "S": "0x5d0c2e55e4d8734fe2516ed56ac628b74c0eb02aa3b6eda51e1e25a1396093e1"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x87390ac14d3c34440b6ddb7b190d3ebde1a07d9a556e5a82ce7e501f24a060f9", "S": "0x37badbcb12cda1ed67b12b1831683a08a3adadee2ea760a07a46bdbb856fea44"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 17,
			"type": "TxTypeFeeDelegatedValueTransferMemo",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x68656c6c6f",
			"signatures": [{"V": "0x26", "R": "0x64e213aef0167fbd853f8f9989ef5d8b912a77457395ccf13d7f37009edd5c5b", "S": "0x5d0c2e55e4d8734fe2516ed56ac628b74c0eb02aa3b6eda51e1e25a1396093e1"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x87390ac14d3c34440b6ddb7b190d3ebde1a07d9a556e5a82ce7e501f24a060f9", "S": "0x37badbcb12cda1ed67b12b1831683a08a3adadee2ea760a07a46bdbb856fea44"}]
		}`,
	},
	{
		Name: "12_FeeDelegatedValueTransferWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedValueTransferMemoWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(0x19),
			TxValueKeyGasLimit:           uint64(0xf4240),
			TxValueKeyTo:                 common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:             big.NewInt(10),
			TxValueKeyFrom:               common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
			TxValueKeyData:               hexutil.MustDecode("0x68656c6c6f"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x769f0afdc310289f9b24decb5bb765c8d7a87a6a4ae28edffb8b7085bbd9bc78", s: "0x6a7b970eea026e60ac29bb52aee10661a4222e6bdcdfb3839a80586e584586b4"}},
		FeePayerSigs:    []txSigHex{{v: 0x25, r: "0xc1c54bdc72ce7c08821329bf50542535fac74f4bba5de5b7881118a461d52834", s: "0x3a3a64878d784f9af91c2e3ab9c90f17144c47cfd9951e3588c75063c0649ecd"}},
		SigRLP:          "0xf842b83df83b128204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f1e018080",
		FeePayerSigRLP:  "0xf857b83df83b128204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f1e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x12f8818204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f1ef845f84326a0769f0afdc310289f9b24decb5bb765c8d7a87a6a4ae28edffb8b7085bbd9bc78a06a7b970eea026e60ac29bb52aee10661a4222e6bdcdfb3839a80586e584586b4",
		TxHashRLP:       "0x12f8dd8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0b8568656c6c6f1ef845f84326a0769f0afdc310289f9b24decb5bb765c8d7a87a6a4ae28edffb8b7085bbd9bc78a06a7b970eea026e60ac29bb52aee10661a4222e6bdcdfb3839a80586e584586b4945a0043070275d9f6054307ee7348bd660849d90ff845f84325a0c1c54bdc72ce7c08821329bf50542535fac74f4bba5de5b7881118a461d52834a03a3a64878d784f9af91c2e3ab9c90f17144c47cfd9951e3588c75063c0649ecd",
		TxJson: `{
			"typeInt": 18,
			"type": "TxTypeFeeDelegatedValueTransferMemoWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x68656c6c6f",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x769f0afdc310289f9b24decb5bb765c8d7a87a6a4ae28edffb8b7085bbd9bc78", "S": "0x6a7b970eea026e60ac29bb52aee10661a4222e6bdcdfb3839a80586e584586b4"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x25", "R": "0xc1c54bdc72ce7c08821329bf50542535fac74f4bba5de5b7881118a461d52834", "S": "0x3a3a64878d784f9af91c2e3ab9c90f17144c47cfd9951e3588c75063c0649ecd"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 18,
			"type": "TxTypeFeeDelegatedValueTransferMemoWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x68656c6c6f",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x769f0afdc310289f9b24decb5bb765c8d7a87a6a4ae28edffb8b7085bbd9bc78", "S": "0x6a7b970eea026e60ac29bb52aee10661a4222e6bdcdfb3839a80586e584586b4"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x25", "R": "0xc1c54bdc72ce7c08821329bf50542535fac74f4bba5de5b7881118a461d52834", "S": "0x3a3a64878d784f9af91c2e3ab9c90f17144c47cfd9951e3588c75063c0649ecd"}]
		}`,
	},
	{
		Name: "20_AccountUpdate", // kaia-sdk
		Type: TxTypeAccountUpdate,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:      uint64(1234),
			TxValueKeyGasPrice:   big.NewInt(0x19),
			TxValueKeyGasLimit:   uint64(0xf4240),
			TxValueKeyFrom:       common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(hexToECDSAPublicKey("0x033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d")),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 0x25, r: "0xf7d479628f05f51320f0842193e3f7ae55a5b49d3645bf55c35bee1e8fd2593a", s: "0x4de8eab5338fdc86e96f8c49ed516550f793fc2c4007614ce3d2a6b33cf9e451"}},
		SigRLP:     "0xf849b844f842208204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d018080",
		TxHashRLP:  "0x20f8888204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33df845f84325a0f7d479628f05f51320f0842193e3f7ae55a5b49d3645bf55c35bee1e8fd2593aa04de8eab5338fdc86e96f8c49ed516550f793fc2c4007614ce3d2a6b33cf9e451",
		TxJson: `{
			"typeInt": 32,
			"type": "TxTypeAccountUpdate",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"signatures": [{"V": "0x25", "R": "0xf7d479628f05f51320f0842193e3f7ae55a5b49d3645bf55c35bee1e8fd2593a", "S": "0x4de8eab5338fdc86e96f8c49ed516550f793fc2c4007614ce3d2a6b33cf9e451"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 32,
			"type": "TxTypeAccountUpdate",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"signatures": [{"V": "0x25", "R": "0xf7d479628f05f51320f0842193e3f7ae55a5b49d3645bf55c35bee1e8fd2593a", "S": "0x4de8eab5338fdc86e96f8c49ed516550f793fc2c4007614ce3d2a6b33cf9e451"}]
		}`,
	},
	{
		Name: "21_FeeDelegatedAccountUpdate", // kaia-sdk
		Type: TxTypeFeeDelegatedAccountUpdate,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:      uint64(1234),
			TxValueKeyGasPrice:   big.NewInt(0x19),
			TxValueKeyGasLimit:   uint64(0xf4240),
			TxValueKeyFrom:       common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:   common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyAccountKey: accountkey.NewAccountKeyPublicWithValue(hexToECDSAPublicKey("0x033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d")),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0xab69d9adca15d9763c4ce6f98b35256717c6e932007658f19c5a255de9e70dda", s: "0x26aa676a3a1a6e96aff4a3df2335788d614d54fb4db1c3c48551ce1fa7ac5e52"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0xf295cd69b4144d9dbc906ba144933d2cc535d9d559f7a92b4672cc5485bf3a60", s: "0x784b8060234ffd64739b5fc2f2503939340ab4248feaa6efcf62cb874345fe40"}},
		SigRLP:          "0xf849b844f842218204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d018080",
		FeePayerSigRLP:  "0xf85eb844f842218204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x21f8888204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33df845f84326a0ab69d9adca15d9763c4ce6f98b35256717c6e932007658f19c5a255de9e70ddaa026aa676a3a1a6e96aff4a3df2335788d614d54fb4db1c3c48551ce1fa7ac5e52",
		TxHashRLP:       "0x21f8e48204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33df845f84326a0ab69d9adca15d9763c4ce6f98b35256717c6e932007658f19c5a255de9e70ddaa026aa676a3a1a6e96aff4a3df2335788d614d54fb4db1c3c48551ce1fa7ac5e52945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0f295cd69b4144d9dbc906ba144933d2cc535d9d559f7a92b4672cc5485bf3a60a0784b8060234ffd64739b5fc2f2503939340ab4248feaa6efcf62cb874345fe40",
		TxJson: `{
			"typeInt": 33,
			"type": "TxTypeFeeDelegatedAccountUpdate",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"signatures": [{"V": "0x26", "R": "0xab69d9adca15d9763c4ce6f98b35256717c6e932007658f19c5a255de9e70dda", "S": "0x26aa676a3a1a6e96aff4a3df2335788d614d54fb4db1c3c48551ce1fa7ac5e52"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xf295cd69b4144d9dbc906ba144933d2cc535d9d559f7a92b4672cc5485bf3a60", "S": "0x784b8060234ffd64739b5fc2f2503939340ab4248feaa6efcf62cb874345fe40"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 33,
			"type": "TxTypeFeeDelegatedAccountUpdate",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"signatures": [{"V": "0x26", "R": "0xab69d9adca15d9763c4ce6f98b35256717c6e932007658f19c5a255de9e70dda", "S": "0x26aa676a3a1a6e96aff4a3df2335788d614d54fb4db1c3c48551ce1fa7ac5e52"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xf295cd69b4144d9dbc906ba144933d2cc535d9d559f7a92b4672cc5485bf3a60", "S": "0x784b8060234ffd64739b5fc2f2503939340ab4248feaa6efcf62cb874345fe40"}]
		}`,
	},
	{
		Name: "22_FeeDelegatedAccountUpdateWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedAccountUpdateWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(0x19),
			TxValueKeyGasLimit:           uint64(0xf4240),
			TxValueKeyFrom:               common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
			TxValueKeyAccountKey:         accountkey.NewAccountKeyPublicWithValue(hexToECDSAPublicKey("0x033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d")),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x0e5929f96dec2b41343a9e6f0150eef08741fe7dcece88cc5936c49ed19051dc", s: "0x5a07b07017190e0baba32bdf6352f5a358a2798ed3c56e704a63819b87cf8e3f"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0xcf8d102de7c6b0a41d3f02aefb7e419522341734c98af233408298d0c424c04b", s: "0x0286f89cab4668f728d7c269997116a49b80cec8776fc64e60588a9268571e35"}},
		SigRLP:          "0xf84ab845f843228204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d1e018080",
		FeePayerSigRLP:  "0xf85fb845f843228204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d1e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x22f8898204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d1ef845f84326a00e5929f96dec2b41343a9e6f0150eef08741fe7dcece88cc5936c49ed19051dca05a07b07017190e0baba32bdf6352f5a358a2798ed3c56e704a63819b87cf8e3f",
		TxHashRLP:       "0x22f8e58204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0ba302a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d1ef845f84326a00e5929f96dec2b41343a9e6f0150eef08741fe7dcece88cc5936c49ed19051dca05a07b07017190e0baba32bdf6352f5a358a2798ed3c56e704a63819b87cf8e3f945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0cf8d102de7c6b0a41d3f02aefb7e419522341734c98af233408298d0c424c04ba00286f89cab4668f728d7c269997116a49b80cec8776fc64e60588a9268571e35",
		TxJson: `{
			"typeInt": 34,
			"type": "TxTypeFeeDelegatedAccountUpdateWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0xe5929f96dec2b41343a9e6f0150eef08741fe7dcece88cc5936c49ed19051dc", "S": "0x5a07b07017190e0baba32bdf6352f5a358a2798ed3c56e704a63819b87cf8e3f"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xcf8d102de7c6b0a41d3f02aefb7e419522341734c98af233408298d0c424c04b", "S": "0x286f89cab4668f728d7c269997116a49b80cec8776fc64e60588a9268571e35"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 34,
			"type": "TxTypeFeeDelegatedAccountUpdateWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"key": "0x02a1033a514176466fa815ed481ffad09110a2d344f6c9b78c1d14afc351c3a51be33d",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0xe5929f96dec2b41343a9e6f0150eef08741fe7dcece88cc5936c49ed19051dc", "S": "0x5a07b07017190e0baba32bdf6352f5a358a2798ed3c56e704a63819b87cf8e3f"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xcf8d102de7c6b0a41d3f02aefb7e419522341734c98af233408298d0c424c04b", "S": "0x286f89cab4668f728d7c269997116a49b80cec8776fc64e60588a9268571e35"}]
		}`,
	},
	{
		Name: "28_SmartContractDeploy", // kaia-sdk
		Type: TxTypeSmartContractDeploy,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:         uint64(532),
			TxValueKeyGasPrice:      big.NewInt(50000000000),
			TxValueKeyGasLimit:      uint64(122000),
			TxValueKeyTo:            (*common.Address)(nil),
			TxValueKeyAmount:        common.Big0,
			TxValueKeyFrom:          common.HexToAddress("0xA2a8854b1802D8Cd5De631E690817c253d6a9153"),
			TxValueKeyData:          hexutil.MustDecode("0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033"),
			TxValueKeyHumanReadable: false,
			TxValueKeyCodeFormat:    params.CodeFormatEVM,
		},
		ChainID:    1001,
		SenderSigs: []txSigHex{{v: 0x7f6, r: "0x71f1da31b7a50b34af48479cca07341bcfe8a3d9cb0b930c942b2ca15e7c928a", s: "0x586178ed946103af9b25f343fb0f8e454c1b49561b692100dfc21ba668567f22"}},
		SigRLP:     "0xf9014bb90143f9014028820214850ba43b74008301dc90808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c6343000812003380808203e98080",
		TxHashRLP:  "0x28f90188820214850ba43b74008301dc90808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c634300081200338080f847f8458207f6a071f1da31b7a50b34af48479cca07341bcfe8a3d9cb0b930c942b2ca15e7c928aa0586178ed946103af9b25f343fb0f8e454c1b49561b692100dfc21ba668567f22",
		TxJson: `{
			"typeInt": 40,
			"type": "TxTypeSmartContractDeploy",
			"nonce": "0x214",
			"gasPrice": "0xba43b7400",
			"gas": "0x1dc90",
			"to": null,
			"value": "0x0",
			"from": "0xa2a8854b1802d8cd5de631e690817c253d6a9153",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f6", "R": "0x71f1da31b7a50b34af48479cca07341bcfe8a3d9cb0b930c942b2ca15e7c928a", "S": "0x586178ed946103af9b25f343fb0f8e454c1b49561b692100dfc21ba668567f22"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 40,
			"type": "TxTypeSmartContractDeploy",
			"nonce": "0x214",
			"gasPrice": "0xba43b7400",
			"gas": "0x1dc90",
			"to": null,
			"value": "0x0",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f6", "R": "0x71f1da31b7a50b34af48479cca07341bcfe8a3d9cb0b930c942b2ca15e7c928a", "S": "0x586178ed946103af9b25f343fb0f8e454c1b49561b692100dfc21ba668567f22"}]
		}`,
	},
	{
		Name: "29_FeeDelegatedSmartContractDeploy", // kaia-sdk
		Type: TxTypeFeeDelegatedSmartContractDeploy,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:         uint64(563),
			TxValueKeyGasPrice:      big.NewInt(50000000000),
			TxValueKeyGasLimit:      uint64(325793),
			TxValueKeyTo:            (*common.Address)(nil),
			TxValueKeyAmount:        common.Big0,
			TxValueKeyFrom:          common.HexToAddress("0xa2a8854b1802d8cd5de631e690817c253d6a9153"),
			TxValueKeyFeePayer:      common.HexToAddress("0xCb0eb737dfda52756495A5e08A9b37AAB3b271dA"),
			TxValueKeyData:          hexutil.MustDecode("0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033"),
			TxValueKeyHumanReadable: false,
			TxValueKeyCodeFormat:    params.CodeFormatEVM,
		},
		ChainID:         1001,
		SenderSigs:      []txSigHex{{v: 0x7f6, r: "0x735b4c96ba68f0853c2ca6836b8fd8246226a453ae82494a00e3e2d1aef3829a", s: "0x05919cbccf2a7a9533719d71502510018f313eb2cef504a4386efe7b615ce570"}},
		FeePayerSigs:    []txSigHex{{v: 0x7f5, r: "0x7799cedd67d7f9b603f2fae6e746aff154530a33d96cd35ee57fad66dd70015f", s: "0x107e893f829df641a00e8c713d2ec795b7153af205d7b6733ec240a5ae3935d8"}},
		SigRLP:          "0xf9014bb90143f9014029820233850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c6343000812003380808203e98080",
		FeePayerSigRLP:  "0xf90160b90143f9014029820233850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033808094cb0eb737dfda52756495a5e08a9b37aab3b271da8203e98080",
		SenderTxHashRLP: "0x29f90188820233850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c634300081200338080f847f8458207f6a0735b4c96ba68f0853c2ca6836b8fd8246226a453ae82494a00e3e2d1aef3829aa005919cbccf2a7a9533719d71502510018f313eb2cef504a4386efe7b615ce570",
		TxHashRLP:       "0x29f901e6820233850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c634300081200338080f847f8458207f6a0735b4c96ba68f0853c2ca6836b8fd8246226a453ae82494a00e3e2d1aef3829aa005919cbccf2a7a9533719d71502510018f313eb2cef504a4386efe7b615ce57094cb0eb737dfda52756495a5e08a9b37aab3b271daf847f8458207f5a07799cedd67d7f9b603f2fae6e746aff154530a33d96cd35ee57fad66dd70015fa0107e893f829df641a00e8c713d2ec795b7153af205d7b6733ec240a5ae3935d8",
		TxJson: `{
			"typeInt": 41,
			"type": "TxTypeFeeDelegatedSmartContractDeploy",
			"nonce": "0x233",
			"gasPrice": "0xba43b7400",
			"gas": "0x4f8a1",
			"to": null,
			"value": "0x0",
			"from": "0xa2a8854b1802d8cd5de631e690817c253d6a9153",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f6", "R": "0x735b4c96ba68f0853c2ca6836b8fd8246226a453ae82494a00e3e2d1aef3829a", "S": "0x5919cbccf2a7a9533719d71502510018f313eb2cef504a4386efe7b615ce570"}],
			"feePayer": "0xcb0eb737dfda52756495a5e08a9b37aab3b271da",
			"feePayerSignatures": [{"V": "0x7f5", "R": "0x7799cedd67d7f9b603f2fae6e746aff154530a33d96cd35ee57fad66dd70015f", "S": "0x107e893f829df641a00e8c713d2ec795b7153af205d7b6733ec240a5ae3935d8"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 41,
			"type": "TxTypeFeeDelegatedSmartContractDeploy",
			"nonce": "0x233",
			"gasPrice": "0xba43b7400",
			"gas": "0x4f8a1",
			"to": null,
			"value": "0x0",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f6", "R": "0x735b4c96ba68f0853c2ca6836b8fd8246226a453ae82494a00e3e2d1aef3829a", "S": "0x5919cbccf2a7a9533719d71502510018f313eb2cef504a4386efe7b615ce570"}],
			"feePayer": "0xcb0eb737dfda52756495a5e08a9b37aab3b271da",
			"feePayerSignatures": [{"V": "0x7f5", "R": "0x7799cedd67d7f9b603f2fae6e746aff154530a33d96cd35ee57fad66dd70015f", "S": "0x107e893f829df641a00e8c713d2ec795b7153af205d7b6733ec240a5ae3935d8"}]
		}`,
	},
	{
		Name: "2a_FeeDelegatedSmartContractDeployWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedSmartContractDeployWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(564),
			TxValueKeyGasPrice:           big.NewInt(50000000000),
			TxValueKeyGasLimit:           uint64(325793),
			TxValueKeyTo:                 (*common.Address)(nil),
			TxValueKeyAmount:             common.Big0,
			TxValueKeyFrom:               common.HexToAddress("0xA2a8854b1802D8Cd5De631E690817c253d6a9153"),
			TxValueKeyFeePayer:           common.HexToAddress("0xCb0eb737dfda52756495A5e08A9b37AAB3b271dA"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
			TxValueKeyData:               hexutil.MustDecode("0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033"),
			TxValueKeyHumanReadable:      false,
			TxValueKeyCodeFormat:         params.CodeFormatEVM,
		},
		ChainID:         1001,
		SenderSigs:      []txSigHex{{v: 0x7f5, r: "0x78763173066acc9396ea8a1b3f65bc6ade4c41d5180cc5f6a546a59ff434c87c", s: "0x5a34f172f5a2f741babd51e11be67de6ceff299baaf9cc87b666b5f0762b00a8"}},
		FeePayerSigs:    []txSigHex{{v: 0x7f6, r: "0xb5534d6fb6edc18f5e923194fbe6d5a0e5816eca15634f257dd1bfc200171ac1", s: "0x3db6b188a27fb5bb96b6bd18a556cccc8e1de7e2a6c6e8474dd91107a3a99b35"}},
		SigRLP:          "0xf9014cb90144f901412a820234850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033801e808203e98080",
		FeePayerSigRLP:  "0xf90161b90144f901412a820234850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033801e8094cb0eb737dfda52756495a5e08a9b37aab3b271da8203e98080",
		SenderTxHashRLP: "0x2af90189820234850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033801e80f847f8458207f5a078763173066acc9396ea8a1b3f65bc6ade4c41d5180cc5f6a546a59ff434c87ca05a34f172f5a2f741babd51e11be67de6ceff299baaf9cc87b666b5f0762b00a8",
		TxHashRLP:       "0x2af901e7820234850ba43b74008304f8a1808094a2a8854b1802d8cd5de631e690817c253d6a9153b90116608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033801e80f847f8458207f5a078763173066acc9396ea8a1b3f65bc6ade4c41d5180cc5f6a546a59ff434c87ca05a34f172f5a2f741babd51e11be67de6ceff299baaf9cc87b666b5f0762b00a894cb0eb737dfda52756495a5e08a9b37aab3b271daf847f8458207f6a0b5534d6fb6edc18f5e923194fbe6d5a0e5816eca15634f257dd1bfc200171ac1a03db6b188a27fb5bb96b6bd18a556cccc8e1de7e2a6c6e8474dd91107a3a99b35",
		TxJson: `{
			"typeInt": 42,
			"type": "TxTypeFeeDelegatedSmartContractDeployWithRatio",
			"nonce": "0x234",
			"gasPrice": "0xba43b7400",
			"gas": "0x4f8a1",
			"to": null,
			"value": "0x0",
			"from": "0xa2a8854b1802d8cd5de631e690817c253d6a9153",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"feeRatio": "0x1e",
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f5", "R": "0x78763173066acc9396ea8a1b3f65bc6ade4c41d5180cc5f6a546a59ff434c87c", "S": "0x5a34f172f5a2f741babd51e11be67de6ceff299baaf9cc87b666b5f0762b00a8"}],
			"feePayer": "0xcb0eb737dfda52756495a5e08a9b37aab3b271da",
			"feePayerSignatures": [{"V": "0x7f6", "R": "0xb5534d6fb6edc18f5e923194fbe6d5a0e5816eca15634f257dd1bfc200171ac1", "S": "0x3db6b188a27fb5bb96b6bd18a556cccc8e1de7e2a6c6e8474dd91107a3a99b35"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 42,
			"type": "TxTypeFeeDelegatedSmartContractDeployWithRatio",
			"nonce": "0x234",
			"gasPrice": "0xba43b7400",
			"gas": "0x4f8a1",
			"to": null,
			"value": "0x0",
			"input": "0x608060405234801561001057600080fd5b5060f78061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c80633fb5c1cb1460415780638381f58a146053578063d09de08a14606d575b600080fd5b6051604c3660046083565b600055565b005b605b60005481565b60405190815260200160405180910390f35b6051600080549080607c83609b565b9190505550565b600060208284031215609457600080fd5b5035919050565b60006001820160ba57634e487b7160e01b600052601160045260246000fd5b506001019056fea2646970667358221220e0f4e7861cb6d7acf0f61d34896310975b57b5bc109681dbbfb2e548ef7546b364736f6c63430008120033",
			"humanReadable": false,
			"feeRatio": "0x1e",
			"codeFormat": "0x0",
			"signatures": [{"V": "0x7f5", "R": "0x78763173066acc9396ea8a1b3f65bc6ade4c41d5180cc5f6a546a59ff434c87c", "S": "0x5a34f172f5a2f741babd51e11be67de6ceff299baaf9cc87b666b5f0762b00a8"}],
			"feePayer": "0xcb0eb737dfda52756495a5e08a9b37aab3b271da",
			"feePayerSignatures": [{"V": "0x7f6", "R": "0xb5534d6fb6edc18f5e923194fbe6d5a0e5816eca15634f257dd1bfc200171ac1", "S": "0x3db6b188a27fb5bb96b6bd18a556cccc8e1de7e2a6c6e8474dd91107a3a99b35"}]
		}`,
	},
	{
		Name: "30_SmartContractExecution", // kaia-sdk
		Type: TxTypeSmartContractExecution,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyData:     hexutil.MustDecode("0x6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2"),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 0x26, r: "0xe4276df1a779274fbb04bc18a0184809eec1ce9770527cebb3d64f926dc1810b", s: "0x4103b828a0671a48d64fe1a3879eae229699f05a684d9c5fd939015dcdd9709b"}},
		SigRLP:     "0xf860b85bf859308204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2018080",
		TxHashRLP:  "0x30f89f8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2f845f84326a0e4276df1a779274fbb04bc18a0184809eec1ce9770527cebb3d64f926dc1810ba04103b828a0671a48d64fe1a3879eae229699f05a684d9c5fd939015dcdd9709b",
		TxJson: `{
			"typeInt": 48,
			"type": "TxTypeSmartContractExecution",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2",
			"signatures": [{"V": "0x26", "R": "0xe4276df1a779274fbb04bc18a0184809eec1ce9770527cebb3d64f926dc1810b", "S": "0x4103b828a0671a48d64fe1a3879eae229699f05a684d9c5fd939015dcdd9709b"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 48,
			"type": "TxTypeSmartContractExecution",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x6353586b000000000000000000000000bc5951f055a85f41a3b62fd6f68ab7de76d299b2",
			"signatures": [{"V": "0x26", "R": "0xe4276df1a779274fbb04bc18a0184809eec1ce9770527cebb3d64f926dc1810b", "S": "0x4103b828a0671a48d64fe1a3879eae229699f05a684d9c5fd939015dcdd9709b"}]
		}`,
	},
	{
		Name: "31_FeeDelegatedSmartContractExecution", // kaia-sdk
		Type: TxTypeFeeDelegatedSmartContractExecution,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyTo:       common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:   big.NewInt(10),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer: common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyData:     hexutil.MustDecode("0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x25, r: "0x253aea7d2c37160da45e84afbb45f6b3341cf1e8fc2df4ecc78f14adb512dc4f", s: "0x22465b74015c2a8f8501186bb5e200e6ce44be52e9374615a7e7e21c41bc27b5"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0xe7c51db7b922c6fa2a941c9687884c593b1b13076bdf0c473538d826bf7b9d1a", s: "0x5b0de2aabb84b66db8bf52d62f3d3b71b592e3748455630f1504c20073624d80"}},
		SigRLP:          "0xf860b85bf859318204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d018080",
		FeePayerSigRLP:  "0xf875b85bf859318204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x31f89f8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723df845f84325a0253aea7d2c37160da45e84afbb45f6b3341cf1e8fc2df4ecc78f14adb512dc4fa022465b74015c2a8f8501186bb5e200e6ce44be52e9374615a7e7e21c41bc27b5",
		TxHashRLP:       "0x31f8fb8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723df845f84325a0253aea7d2c37160da45e84afbb45f6b3341cf1e8fc2df4ecc78f14adb512dc4fa022465b74015c2a8f8501186bb5e200e6ce44be52e9374615a7e7e21c41bc27b5945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0e7c51db7b922c6fa2a941c9687884c593b1b13076bdf0c473538d826bf7b9d1aa05b0de2aabb84b66db8bf52d62f3d3b71b592e3748455630f1504c20073624d80",
		TxJson: `{
			"typeInt": 49,
			"type": "TxTypeFeeDelegatedSmartContractExecution",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d",
			"signatures": [{"V": "0x25", "R": "0x253aea7d2c37160da45e84afbb45f6b3341cf1e8fc2df4ecc78f14adb512dc4f", "S": "0x22465b74015c2a8f8501186bb5e200e6ce44be52e9374615a7e7e21c41bc27b5"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xe7c51db7b922c6fa2a941c9687884c593b1b13076bdf0c473538d826bf7b9d1a", "S": "0x5b0de2aabb84b66db8bf52d62f3d3b71b592e3748455630f1504c20073624d80"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 49,
			"type": "TxTypeFeeDelegatedSmartContractExecution",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d",
			"signatures": [{"V": "0x25", "R": "0x253aea7d2c37160da45e84afbb45f6b3341cf1e8fc2df4ecc78f14adb512dc4f", "S": "0x22465b74015c2a8f8501186bb5e200e6ce44be52e9374615a7e7e21c41bc27b5"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0xe7c51db7b922c6fa2a941c9687884c593b1b13076bdf0c473538d826bf7b9d1a", "S": "0x5b0de2aabb84b66db8bf52d62f3d3b71b592e3748455630f1504c20073624d80"}]
		}`,
	},
	{
		Name: "32_FeeDelegatedSmartContractExecutionWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedSmartContractExecutionWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(0x19),
			TxValueKeyGasLimit:           uint64(0xf4240),
			TxValueKeyTo:                 common.HexToAddress("0x7b65B75d204aBed71587c9E519a89277766EE1d0"),
			TxValueKeyAmount:             big.NewInt(10),
			TxValueKeyFrom:               common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
			TxValueKeyData:               hexutil.MustDecode("0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x74ccfee18dc28932396b85617c53784ee366303bce39a2401d8eb602cf73766f", s: "0x4c937a5ab9401d2cacb3f39ba8c29dbcd44588cc5c7d0b6b4113cfa7b7d9427b"}},
		FeePayerSigs:    []txSigHex{{v: 0x25, r: "0x4a4997524694d535976d7343c1e3a260f99ba53fcb5477e2b96216ec96ebb565", s: "0x0f8cb31a35399d2b0fbbfa39f259c819a15370706c0449952c7cfc682d200d7c"}},
		SigRLP:          "0xf861b85cf85a328204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d1e018080",
		FeePayerSigRLP:  "0xf876b85cf85a328204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d1e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x32f8a08204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d1ef845f84326a074ccfee18dc28932396b85617c53784ee366303bce39a2401d8eb602cf73766fa04c937a5ab9401d2cacb3f39ba8c29dbcd44588cc5c7d0b6b4113cfa7b7d9427b",
		TxHashRLP:       "0x32f8fc8204d219830f4240947b65b75d204abed71587c9e519a89277766ee1d00a94a94f5374fce5edbc8e2a8697c15331677e6ebf0ba46353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d1ef845f84326a074ccfee18dc28932396b85617c53784ee366303bce39a2401d8eb602cf73766fa04c937a5ab9401d2cacb3f39ba8c29dbcd44588cc5c7d0b6b4113cfa7b7d9427b945a0043070275d9f6054307ee7348bd660849d90ff845f84325a04a4997524694d535976d7343c1e3a260f99ba53fcb5477e2b96216ec96ebb565a00f8cb31a35399d2b0fbbfa39f259c819a15370706c0449952c7cfc682d200d7c",
		TxJson: `{
			"typeInt": 50,
			"type": "TxTypeFeeDelegatedSmartContractExecutionWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"input": "0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x74ccfee18dc28932396b85617c53784ee366303bce39a2401d8eb602cf73766f", "S": "0x4c937a5ab9401d2cacb3f39ba8c29dbcd44588cc5c7d0b6b4113cfa7b7d9427b"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x25", "R": "0x4a4997524694d535976d7343c1e3a260f99ba53fcb5477e2b96216ec96ebb565", "S": "0xf8cb31a35399d2b0fbbfa39f259c819a15370706c0449952c7cfc682d200d7c"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 50,
			"type": "TxTypeFeeDelegatedSmartContractExecutionWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"to": "0x7b65b75d204abed71587c9e519a89277766ee1d0",
			"value": "0xa",
			"input": "0x6353586b0000000000000000000000000fcda0f2efbe1b4e61b487701ce4f2f8abc3723d",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x74ccfee18dc28932396b85617c53784ee366303bce39a2401d8eb602cf73766f", "S": "0x4c937a5ab9401d2cacb3f39ba8c29dbcd44588cc5c7d0b6b4113cfa7b7d9427b"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x25", "R": "0x4a4997524694d535976d7343c1e3a260f99ba53fcb5477e2b96216ec96ebb565", "S": "0xf8cb31a35399d2b0fbbfa39f259c819a15370706c0449952c7cfc682d200d7c"}]
		}`,
	},
	{
		Name: "38_Cancel", // kaia-sdk
		Type: TxTypeCancel,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 0x25, r: "0xfb2c3d53d2f6b7bb1deb5a09f80366a5a45429cc1e3956687b075a9dcad20434", s: "0x5c6187822ee23b1001e9613d29a5d6002f990498d2902904f7f259ab3358216e"}},
		SigRLP:     "0xe39fde388204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b018080",
		TxHashRLP:  "0x38f8648204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84325a0fb2c3d53d2f6b7bb1deb5a09f80366a5a45429cc1e3956687b075a9dcad20434a05c6187822ee23b1001e9613d29a5d6002f990498d2902904f7f259ab3358216e",
		TxJson: `{
			"typeInt": 56,
			"type": "TxTypeCancel",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"signatures": [{"V": "0x25", "R": "0xfb2c3d53d2f6b7bb1deb5a09f80366a5a45429cc1e3956687b075a9dcad20434", "S": "0x5c6187822ee23b1001e9613d29a5d6002f990498d2902904f7f259ab3358216e"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 56,
			"type": "TxTypeCancel",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"signatures": [{"V": "0x25", "R": "0xfb2c3d53d2f6b7bb1deb5a09f80366a5a45429cc1e3956687b075a9dcad20434", "S": "0x5c6187822ee23b1001e9613d29a5d6002f990498d2902904f7f259ab3358216e"}]
		}`,
	},
	{
		Name: "39_FeeDelegatedCancel", // kaia-sdk
		Type: TxTypeFeeDelegatedCancel,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:    uint64(1234),
			TxValueKeyGasPrice: big.NewInt(0x19),
			TxValueKeyGasLimit: uint64(0xf4240),
			TxValueKeyFrom:     common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer: common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x8409f5441d4725f90905ad87f03793857d124de7a43169bc67320cd2f020efa9", s: "0x60af63e87bdc565d7f7de906916b2334336ee7b24d9a71c9521a67df02e7ec92"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0x044d5b25e8c649a1fdaa409dc3817be390ad90a17c25bc17c89b6d5d248495e0", s: "0x73938e690d27b5267c73108352cf12d01de7fd0077b388e94721aa1fa32f85ec"}},
		SigRLP:          "0xe39fde398204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b018080",
		FeePayerSigRLP:  "0xf8389fde398204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x39f8648204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84326a08409f5441d4725f90905ad87f03793857d124de7a43169bc67320cd2f020efa9a060af63e87bdc565d7f7de906916b2334336ee7b24d9a71c9521a67df02e7ec92",
		TxHashRLP:       "0x39f8c08204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0bf845f84326a08409f5441d4725f90905ad87f03793857d124de7a43169bc67320cd2f020efa9a060af63e87bdc565d7f7de906916b2334336ee7b24d9a71c9521a67df02e7ec92945a0043070275d9f6054307ee7348bd660849d90ff845f84326a0044d5b25e8c649a1fdaa409dc3817be390ad90a17c25bc17c89b6d5d248495e0a073938e690d27b5267c73108352cf12d01de7fd0077b388e94721aa1fa32f85ec",
		TxJson: `{
			"typeInt": 57,
			"type": "TxTypeFeeDelegatedCancel",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"signatures": [{"V": "0x26", "R": "0x8409f5441d4725f90905ad87f03793857d124de7a43169bc67320cd2f020efa9", "S": "0x60af63e87bdc565d7f7de906916b2334336ee7b24d9a71c9521a67df02e7ec92"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x44d5b25e8c649a1fdaa409dc3817be390ad90a17c25bc17c89b6d5d248495e0", "S": "0x73938e690d27b5267c73108352cf12d01de7fd0077b388e94721aa1fa32f85ec"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 57,
			"type": "TxTypeFeeDelegatedCancel",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"signatures": [{"V": "0x26", "R": "0x8409f5441d4725f90905ad87f03793857d124de7a43169bc67320cd2f020efa9", "S": "0x60af63e87bdc565d7f7de906916b2334336ee7b24d9a71c9521a67df02e7ec92"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x44d5b25e8c649a1fdaa409dc3817be390ad90a17c25bc17c89b6d5d248495e0", "S": "0x73938e690d27b5267c73108352cf12d01de7fd0077b388e94721aa1fa32f85ec"}]
		}`,
	},
	{
		Name: "3a_FeeDelegatedCancelWithRatio", // kaia-sdk
		Type: TxTypeFeeDelegatedCancelWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(0x19),
			TxValueKeyGasLimit:           uint64(0xf4240),
			TxValueKeyFrom:               common.HexToAddress("0xa94f5374Fce5edBC8E2a8697C15331677e6EbF0B"),
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 0x26, r: "0x72efa47960bef40b536c72d7e03ceaf6ca5f6061eb8a3eda3545b1a78fe52ef5", s: "0x62006ddaf874da205f08b3789e2d014ae37794890fc2e575bf75201563a24ba9"}},
		FeePayerSigs:    []txSigHex{{v: 0x26, r: "0x6ba5ef20c3049323fc94defe14ca162e28b86aa64f7cf497ac8a5520e9615614", s: "0x4a0a0fc61c10b416759af0ce4ce5c09ca1060141d56d958af77050c9564df6bf"}},
		SigRLP:          "0xe4a0df3a8204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b1e018080",
		FeePayerSigRLP:  "0xf839a0df3a8204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b1e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x3af8658204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b1ef845f84326a072efa47960bef40b536c72d7e03ceaf6ca5f6061eb8a3eda3545b1a78fe52ef5a062006ddaf874da205f08b3789e2d014ae37794890fc2e575bf75201563a24ba9",
		TxHashRLP:       "0x3af8c18204d219830f424094a94f5374fce5edbc8e2a8697c15331677e6ebf0b1ef845f84326a072efa47960bef40b536c72d7e03ceaf6ca5f6061eb8a3eda3545b1a78fe52ef5a062006ddaf874da205f08b3789e2d014ae37794890fc2e575bf75201563a24ba9945a0043070275d9f6054307ee7348bd660849d90ff845f84326a06ba5ef20c3049323fc94defe14ca162e28b86aa64f7cf497ac8a5520e9615614a04a0a0fc61c10b416759af0ce4ce5c09ca1060141d56d958af77050c9564df6bf",
		TxJson: `{
			"typeInt": 58,
			"type": "TxTypeFeeDelegatedCancelWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"from": "0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x72efa47960bef40b536c72d7e03ceaf6ca5f6061eb8a3eda3545b1a78fe52ef5", "S": "0x62006ddaf874da205f08b3789e2d014ae37794890fc2e575bf75201563a24ba9"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x6ba5ef20c3049323fc94defe14ca162e28b86aa64f7cf497ac8a5520e9615614", "S": "0x4a0a0fc61c10b416759af0ce4ce5c09ca1060141d56d958af77050c9564df6bf"}],
			"hash": null
		}`,
		RpcJson: `{
			"typeInt": 58,
			"type": "TxTypeFeeDelegatedCancelWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x19",
			"gas": "0xf4240",
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x26", "R": "0x72efa47960bef40b536c72d7e03ceaf6ca5f6061eb8a3eda3545b1a78fe52ef5", "S": "0x62006ddaf874da205f08b3789e2d014ae37794890fc2e575bf75201563a24ba9"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x26", "R": "0x6ba5ef20c3049323fc94defe14ca162e28b86aa64f7cf497ac8a5520e9615614", "S": "0x4a0a0fc61c10b416759af0ce4ce5c09ca1060141d56d958af77050c9564df6bf"}]
		}`,
	},
	{
		Name: "48_ChainDataAnchoring", // Manual creation
		Type: TxTypeChainDataAnchoring,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:        uint64(1234),
			TxValueKeyGasPrice:     big.NewInt(25 * params.Gkei),
			TxValueKeyGasLimit:     uint64(50000000),
			TxValueKeyFrom:         common.HexToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			TxValueKeyAnchoredData: []byte{1, 2, 3, 4},
		},
		ChainID:    1,
		SenderSigs: []txSigHex{{v: 1, r: "0x2", s: "0x3"}},
		SigRLP:     "0xeeaae9488204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304018080",
		TxHashRLP:  "0x48ed8204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304c4c3010203",
		TxJson: `{
			"typeInt": 72,
			"type": "TxTypeChainDataAnchoring",
			"nonce": "0x4d2",
			"gasPrice": "0x5d21dba00",
			"gas": "0x2faf080",
			"from": "0x23a519a88e79fbc0bab796f3dce3ff79a2373e30",
			"input": "0x01020304",
			"inputJSON": null,
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 72,
			"type": "TxTypeChainDataAnchoring",
			"gas": "0x2faf080",
			"gasPrice": "0x5d21dba00",
			"nonce": "0x4d2",
			"input": "0x01020304",
			"inputJSON": null,
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}]
		}`,
	},
	{
		Name: "49_FeeDelegatedChainDataAnchoring", // Manual creation
		Type: TxTypeFeeDelegatedChainDataAnchoring,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:        uint64(1234),
			TxValueKeyGasPrice:     big.NewInt(25 * params.Gkei),
			TxValueKeyGasLimit:     uint64(50000000),
			TxValueKeyFrom:         common.HexToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			TxValueKeyAnchoredData: []byte{1, 2, 3, 4},
			TxValueKeyFeePayer:     common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 1, r: "0x2", s: "0x3"}},
		FeePayerSigs:    []txSigHex{{v: 1, r: "0x2", s: "0x3"}},
		SigRLP:          "0xeeaae9498204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304018080",
		FeePayerSigRLP:  "0xf843aae9498204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x49ed8204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304c4c3010203",
		TxHashRLP:       "0x49f8478204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e308401020304c4c3010203945a0043070275d9f6054307ee7348bd660849d90fc4c3010203",
		TxJson: `{
			"typeInt": 73,
			"type": "TxTypeFeeDelegatedChainDataAnchoring",
			"nonce": "0x4d2",
			"gasPrice": "0x5d21dba00",
			"gas": "0x2faf080",
			"from": "0x23a519a88e79fbc0bab796f3dce3ff79a2373e30",
			"input": "0x01020304",
			"inputJSON": null,
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 73,
			"type": "TxTypeFeeDelegatedChainDataAnchoring",
			"gas": "0x2faf080",
			"gasPrice": "0x5d21dba00",
			"nonce": "0x4d2",
			"input": "0x01020304",
			"inputJSON": null,
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}]
		}`,
	},
	{
		Name: "4a_FeeDelegatedChainDataAnchoringWithRatio", // Manual creation
		Type: TxTypeFeeDelegatedChainDataAnchoringWithRatio,
		Map: map[TxValueKeyType]interface{}{
			TxValueKeyNonce:              uint64(1234),
			TxValueKeyGasPrice:           big.NewInt(25 * params.Gkei),
			TxValueKeyGasLimit:           uint64(50000000),
			TxValueKeyFrom:               common.HexToAddress("0x23a519a88e79fbc0bab796f3dce3ff79a2373e30"),
			TxValueKeyAnchoredData:       []byte{1, 2, 3, 4},
			TxValueKeyFeePayer:           common.HexToAddress("0x5A0043070275d9f6054307Ee7348bD660849D90f"),
			TxValueKeyFeeRatioOfFeePayer: FeeRatio(30),
		},
		ChainID:         1,
		SenderSigs:      []txSigHex{{v: 1, r: "0x2", s: "0x3"}},
		FeePayerSigs:    []txSigHex{{v: 1, r: "0x2", s: "0x3"}},
		SigRLP:          "0xefabea4a8204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e3084010203041e018080",
		FeePayerSigRLP:  "0xf844abea4a8204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e3084010203041e945a0043070275d9f6054307ee7348bd660849d90f018080",
		SenderTxHashRLP: "0x4aee8204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e3084010203041ec4c3010203",
		TxHashRLP:       "0x4af8488204d28505d21dba008402faf0809423a519a88e79fbc0bab796f3dce3ff79a2373e3084010203041ec4c3010203945a0043070275d9f6054307ee7348bd660849d90fc4c3010203",
		TxJson: `{
			"typeInt": 74,
			"type": "TxTypeFeeDelegatedChainDataAnchoringWithRatio",
			"nonce": "0x4d2",
			"gasPrice": "0x5d21dba00",
			"gas": "0x2faf080",
			"from": "0x23a519a88e79fbc0bab796f3dce3ff79a2373e30",
			"input": "0x01020304",
			"inputJSON": null,
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"hash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		RpcJson: `{
			"typeInt": 74,
			"type": "TxTypeFeeDelegatedChainDataAnchoringWithRatio",
			"gas": "0x2faf080",
			"gasPrice": "0x5d21dba00",
			"nonce": "0x4d2",
			"input": "0x01020304",
			"inputJSON": null,
			"feeRatio": "0x1e",
			"signatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}],
			"feePayer": "0x5a0043070275d9f6054307ee7348bd660849d90f",
			"feePayerSignatures": [{"V": "0x1", "R": "0x2", "S": "0x3"}]
		}`,
	},
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
		TxValueKeySidecar:    sidecarV0,
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

//go:embed testdata/blob0.bin
var blob0 string

//go:embed testdata/blobtxWithSidecarV0HashRLP.rlp
var blobtxWithSidecarV0HashRLP string

//go:embed testdata/blobtxWithSidecarV1HashRLP.rlp
var blobtxWithSidecarV1HashRLP string

var (
	blobtxWithSidecarV0TxJson = `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": {"Version": 0, "Blobs": ["%s"], "Commitments": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"], "Proofs": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"]},
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}],
			"hash": null
		}`
	blobtxWithSidecarV0RPCJson = `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": {"Version": 0, "Blobs": ["%s"], "Commitments": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"], "Proofs": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"]},
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}]
		}`
	blobtxWithSidecarV1TxJson = `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": {"Version": 1, "Blobs": ["%s"], "Commitments": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"], "Proofs": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"]},
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}],
			"hash": null
		}`
	blobtxWithSidecarV1RPCJson = `{
			"typeInt": 30723,
			"type": "TxTypeEthereumBlob",
			"chainId": "0x6",
			"nonce": "0x1896",
			"maxFeePerGas": "0x0",
			"maxPriorityFeePerGas": "0x0",
			"gas": "0x0",
			"to": "0xdf3ca4eaf9017d01a26ef475e651faa9b1296da1",
			"value": "0x0",
			"input": "0x09d34150cb13b7867ed4a95638b03d3c4ff4d065b901f4351e89091cd0",
			"accessList": [{"address": "0x6092415c41b602d192c02d8bb5b2ee62fbab3b70", "storageKeys": ["0xa2c53cdc4de0f875229c19c1d05f5f0000000000000000000000000000000000"]}],
			"blobFeeCap": "0xf448c89d13854f96",
			"blobHashes": ["0x012730cf6ab975c7c39a00000000000000000000000000000000000000000000", "0x01f263630289db00000000000000000000000000000000000000000000000000", "0x0159b494f64c2c6adac876c9d3ea38f46c9aca7d869300000000000000000000"],
			"sidecar": {"Version": 1, "Blobs": ["%s"], "Commitments": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"], "Proofs": ["0xc00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"]},
			"signatures": [{"V": "0x1b", "R": "0xf529d0d7d2687fef8d097aafec3d8363ec5d69e29140c9603fba0179a2518b2b", "S": "0xdef24511874e80989362ae91d7509ec5eb09f81c8a0c039f4e7a66ea86e6746"}]
		}`
)
