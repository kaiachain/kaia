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

package account

import (
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
)

// Compile time interface checks
var (
	_ Account = (*LegacyAccount)(nil)
	_ Account = (*ExternallyOwnedAccount)(nil)
	_ Account = (*SmartContractAccount)(nil)

	_ ProgramAccount = (*SmartContractAccount)(nil)

	_ AccountWithKey = (*ExternallyOwnedAccount)(nil)
	_ AccountWithKey = (*SmartContractAccount)(nil)
)

// TestAccountSerialization tests serialization of various account types.
func TestAccountSerialization(t *testing.T) {
	accs := []struct {
		Name string
		acc  Account
	}{
		{"EOA", genEOA()},
		{"EOAWithPublic", genEOAWithPublicKey()},
		{"SCA", genSCA()},
		{"SCAWithPublic", genSCAWithPublicKey()},
	}
	testcases := []struct {
		Name string
		fn   func(t *testing.T, acc Account)
	}{
		{"RLP", testAccountRLP},
		{"JSON", testAccountJSON},
	}
	for _, test := range testcases {
		for _, acc := range accs {
			Name := test.Name + "/" + acc.Name
			t.Run(Name, func(t *testing.T) {
				test.fn(t, acc.acc)
			})
		}
	}
}

func testAccountRLP(t *testing.T, acc Account) {
	enc := NewAccountSerializerWithAccount(acc)

	b, err := rlp.EncodeToBytes(enc)
	if err != nil {
		panic(err)
	}

	dec := NewAccountSerializer()

	if err := rlp.DecodeBytes(b, &dec); err != nil {
		panic(err)
	}

	if !acc.Equal(dec.account) {
		fmt.Println("acc")
		fmt.Println(acc)
		fmt.Println("dec.account")
		fmt.Println(dec.account)
		t.Errorf("acc != dec.account")
	}
}

func testAccountJSON(t *testing.T, acc Account) {
	enc := NewAccountSerializerWithAccount(acc)

	b, err := json.Marshal(enc)
	if err != nil {
		panic(err)
	}

	dec := NewAccountSerializer()

	if err := json.Unmarshal(b, &dec); err != nil {
		panic(err)
	}

	if !acc.Equal(dec.account) {
		fmt.Println("acc")
		fmt.Println(acc)
		fmt.Println("dec.account")
		fmt.Println(dec.account)
		t.Errorf("acc != dec.account")
	}
}

func genRandomHash() (h common.Hash) {
	hasher := sha3.NewKeccak256()

	r := rand.Uint64()
	rlp.Encode(hasher, r)
	hasher.Sum(h[:0])

	return h
}

func genEOA() *ExternallyOwnedAccount {
	humanReadable := false

	return newExternallyOwnedAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
	})
}

func genEOAWithPublicKey() *ExternallyOwnedAccount {
	humanReadable := false

	k, _ := crypto.GenerateKey()

	return newExternallyOwnedAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&k.PublicKey),
	})
}

func genSCA() *SmartContractAccount {
	humanReadable := false

	return newSmartContractAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyLegacy(),
		AccountValueKeyStorageRoot:   genRandomHash(),
		AccountValueKeyCodeHash:      genRandomHash().Bytes(),
		AccountValueKeyCodeInfo:      params.CodeInfo(0),
	})
}

func genSCAWithPublicKey() *SmartContractAccount {
	humanReadable := false

	k, _ := crypto.GenerateKey()

	return newSmartContractAccountWithMap(map[AccountValueKeyType]interface{}{
		AccountValueKeyNonce:         rand.Uint64(),
		AccountValueKeyBalance:       big.NewInt(rand.Int63n(10000)),
		AccountValueKeyHumanReadable: humanReadable,
		AccountValueKeyAccountKey:    accountkey.NewAccountKeyPublicWithValue(&k.PublicKey),
		AccountValueKeyStorageRoot:   genRandomHash(),
		AccountValueKeyCodeHash:      genRandomHash().Bytes(),
		AccountValueKeyCodeInfo:      params.CodeInfo(0),
	})
}

func checkEncode(t *testing.T, account Account, expected string) {
	enc := NewAccountSerializerWithAccount(account)
	b, err := rlp.EncodeToBytes(enc)
	assert.Nil(t, err)
	assert.Equal(t, expected, hexutil.Encode(b))
}

func checkEncodeExt(t *testing.T, account Account, expected string) {
	enc := NewAccountSerializerExtWithAccount(account)
	b, err := rlp.EncodeToBytes(enc)
	assert.Nil(t, err)
	assert.Equal(t, expected, hexutil.Encode(b))
}

func checkDecode(t *testing.T, encoded string, expected Account) {
	b := common.FromHex(encoded)
	dec := NewAccountSerializer()
	err := rlp.DecodeBytes(b, &dec)
	assert.Nil(t, err)
	assert.True(t, dec.GetAccount().Equal(expected))
}

func checkDecodeExt(t *testing.T, encoded string, expected Account) {
	b := common.FromHex(encoded)
	dec := NewAccountSerializerExt()
	err := rlp.DecodeBytes(b, &dec)
	assert.Nil(t, err)
	assert.True(t, dec.GetAccount().Equal(expected))
}

func checkEncodeJSON(t *testing.T, account Account, expectedMap map[string]interface{}) {
	enc := NewAccountSerializerWithAccount(account)
	actual, err := json.Marshal(enc)
	assert.Nil(t, err)

	expected, err := json.Marshal(expectedMap)
	assert.Nil(t, err)
	assert.JSONEq(t, string(expected), string(actual))
}

func checkDecodeJSON(t *testing.T, j map[string]interface{}, expected Account) {
	b, err := json.Marshal(j)
	assert.Nil(t, err)

	dec := NewAccountSerializer()
	err = dec.UnmarshalJSON(b)
	assert.Nil(t, err)
	assert.True(t, dec.GetAccount().Equal(expected))
}

func TestAccountSerializer(t *testing.T) {
	var (
		k, _         = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		commonFields = &AccountCommon{
			nonce:         42,
			balance:       big.NewInt(0x12345678),
			humanReadable: false,
			key:           accountkey.NewAccountKeyLegacy(),
		}
		commonFieldsEmpty = &AccountCommon{
			nonce:         0,
			balance:       big.NewInt(0),
			humanReadable: false,
			key:           accountkey.NewAccountKeyLegacy(),
		}
		commonFieldsUpdated = &AccountCommon{
			nonce:         42,
			balance:       big.NewInt(0x12345678),
			humanReadable: false,
			key:           accountkey.NewAccountKeyPublicWithValue(&k.PublicKey),
		}
		legacyKeyJson = map[string]interface{}{
			"keyType": 1,
			"key":     map[string]interface{}{},
		}
		pubKeyJson = map[string]interface{}{
			"keyType": 2,
			"key": map[string]interface{}{
				"x": "0x8318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75",
				"y": "0x3547f11ca8696646f2f3acb08e31016afac23e630c5d11f59f61fef57b0d2aa5",
			},
		}
		roothash    = common.HexToHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")
		roothashExt = common.HexToExtHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01")
		codehash    = common.HexToHash("aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd").Bytes()
		codehashB64 = "qqqqqru7u7vMzMzM3d3d3aqqqqq7u7u7zMzMzN3d3d0="
		codeinfo    = params.CodeInfo(0x10) // VmVersion=1 (Istanbul+), CodeFormat=0 (EVM)
	)

	testcases := []struct {
		desc string
		acc  Account

		// RLP: storageRoot must be unextended (32-bytes)
		rlp string
		// RLPExt: storageRoot (if exists) is unextended if originally zero-extended (32-bytes), kept extended if it was nonzero-extended (37-bytes)
		rlpExt string
		// JSON: storageRoot (if exists) must be unextended (32-bytes)
		json map[string]interface{}
	}{
		{
			"Empty EOA",
			&ExternallyOwnedAccount{
				AccountCommon: commonFieldsEmpty,
			},
			// 01 ["","","","01",[]]
			"0x01c580808001c0",
			"0x01c580808001c0",
			map[string]interface{}{
				"accType": 1,
				"account": map[string]interface{}{
					"nonce":         0,
					"balance":       "0x0",
					"humanReadable": false,
					"key":           legacyKeyJson,
				},
			},
		},
		{
			"Nonempty EOA",
			&ExternallyOwnedAccount{
				AccountCommon: commonFields,
			},
			// 01 ["0x2a","0x12345678","","0x01",[]]
			"0x01c92a84123456788001c0",
			"0x01c92a84123456788001c0",
			map[string]interface{}{
				"accType": 1,
				"account": map[string]interface{}{
					"nonce":         42,
					"balance":       "0x12345678",
					"humanReadable": false,
					"key":           legacyKeyJson,
				},
			},
		},
		{
			"AccountUpdated EOA",
			&ExternallyOwnedAccount{
				AccountCommon: commonFieldsUpdated,
			},
			// 01 ["0x2a","0x12345678","","0x02","0x038318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75"]
			"0x01ea2a84123456788002a1038318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75",
			"0x01ea2a84123456788002a1038318535b54105d4a7aae60c08fc45f9687181b4fdfc625bd1a753fa7397fed75",
			map[string]interface{}{
				"accType": 1,
				"account": map[string]interface{}{
					"nonce":         42,
					"balance":       "0x12345678",
					"humanReadable": false,
					"key":           pubKeyJson,
				},
			},
		},
		{
			"SCA with Zero-extended StorageRoot",
			&SmartContractAccount{
				AccountCommon: commonFields,
				storageRoot:   roothash.ExtendZero(),
				codeHash:      codehash,
				codeInfo:      codeinfo,
			},
			// 02 [["0x2a","0x12345678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]
			"0x02f84dc92a84123456788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			"0x02f84dc92a84123456788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			map[string]interface{}{
				"accType": 2,
				"account": map[string]interface{}{
					"nonce":         42,
					"balance":       "0x12345678",
					"humanReadable": false,
					"key":           legacyKeyJson,
					"storageRoot":   roothash.Hex(),
					"codeHash":      codehashB64,
					"codeFormat":    0,
					"vmVersion":     1,
				},
			},
		},
		{
			"SCA with Nonzero-extended StorageRoot",
			&SmartContractAccount{
				AccountCommon: commonFields,
				storageRoot:   roothashExt,
				codeHash:      codehash,
				codeInfo:      codeinfo,
			},
			// [["0x2a","0x12345678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]
			"0x02f84dc92a84123456788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			// [["0x2a","0x12345678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]
			"0x02f854c92a84123456788001c0a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01a0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			map[string]interface{}{
				"accType": 2,
				"account": map[string]interface{}{
					"nonce":         42,
					"balance":       "0x12345678",
					"humanReadable": false,
					"key":           legacyKeyJson,
					"storageRoot":   roothash.Hex(),
					"codeHash":      codehashB64,
					"codeFormat":    0,
					"vmVersion":     1,
				},
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.desc, func(t *testing.T) {
			// If original StorageRoot is nonzero-extended, skip the unextended RLP and JSON tests
			// as those encodings lose the extension information.
			hasNonzeroExtendedStorageRoot := false
			if pa := GetProgramAccount(tc.acc); pa != nil {
				hasNonzeroExtendedStorageRoot = !pa.GetStorageRoot().IsZeroExtended()
			}

			// obj -> rlp
			checkEncode(t, tc.acc, tc.rlp)
			checkEncodeExt(t, tc.acc, tc.rlpExt)

			// rlp -> obj
			if !hasNonzeroExtendedStorageRoot {
				checkDecode(t, tc.rlp, tc.acc)
			}
			checkDecodeExt(t, tc.rlpExt, tc.acc)

			// obj -> json
			checkEncodeJSON(t, tc.acc, tc.json)

			// json -> obj
			if !hasNonzeroExtendedStorageRoot {
				checkDecodeJSON(t, tc.json, tc.acc)
			}
		})
	}
}

// Tests RLP encoding against manually generated strings.
func TestSmartContractAccountExt(t *testing.T) {
	// To create testcases,
	// - Install https://github.com/ethereumjs/rlp
	//     npm install -g rlp
	// - In bash, run
	//     maketc(){ echo $(rlp encode "$1")$(rlp encode "$2" | cut -b3-); }
	//     maketc 2 '["0x1234","0x5678"]'
	var (
		commonFields = &AccountCommon{
			nonce:         0x1234,
			balance:       big.NewInt(0x5678),
			humanReadable: false,
			key:           accountkey.NewAccountKeyLegacy(),
		}
		hash     = common.HexToHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff")
		exthash  = common.HexToExtHash("00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01")
		codehash = common.HexToHash("aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd").Bytes()
		codeinfo = params.CodeInfo(0x10)

		// StorageRoot is hash32:  maketc 2 '[["0x1234","0x5678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]'
		scaUnextRLP = "0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10"
		scaUnext    = &SmartContractAccount{
			AccountCommon: commonFields,
			storageRoot:   hash.ExtendZero(),
			codeHash:      codehash,
			codeInfo:      codeinfo,
		}
		// StorageRoot is exthash: maketc 2 '[["0x1234","0x5678","","0x01",[]],"0x00112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01","0xaaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd","0x10"]'
		scaExtRLP = "0x02f854c98212348256788001c0a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01a0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10"
		scaExt    = &SmartContractAccount{
			AccountCommon: commonFields,
			storageRoot:   exthash,
			codeHash:      codehash,
			codeInfo:      codeinfo,
		}
	)
	checkEncode(t, scaUnext, scaUnextRLP)
	checkEncodeExt(t, scaUnext, scaUnextRLP) // zero extensions are always unextended

	checkEncode(t, scaExt, scaUnextRLP)  // Regular encoding still results in hash32. Use it for merkle hash.
	checkEncodeExt(t, scaExt, scaExtRLP) // Must use SerializeExt to preserve exthash. Use it for disk storage.

	checkDecode(t, scaUnextRLP, scaUnext)
	checkDecodeExt(t, scaUnextRLP, scaUnext)

	checkDecode(t, scaExtRLP, scaExt)
	checkDecodeExt(t, scaExtRLP, scaExt)
}

func TestUnextendRLP(t *testing.T) {
	// storage slot
	testcases := []struct {
		extended   string
		unextended string
	}{
		{ // storage slot (33B) kept as-is
			"0xa06700000000000000000000000000000000000000000000000000000000000000",
			"0xa06700000000000000000000000000000000000000000000000000000000000000",
		},
		{ // Short EOA (<=ExtHashLength) kept as-is
			"0x01c98212348256788001c0",
			"0x01c98212348256788001c0",
		},
		{ // Long EOA (>ExtHashLength) kept as-is
			"0x01ea8212348256788002a1030bc77753515dd61c66df6445ffffbedfc16b6b46c73eb09f01a970cb3bf0a8de",
			"0x01ea8212348256788002a1030bc77753515dd61c66df6445ffffbedfc16b6b46c73eb09f01a970cb3bf0a8de",
		},
		{ // SCA with Hash keps as-is
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
		},
		{ // SCA with ExtHash unextended
			"0x02f854c98212348256788001c0a700112233445566778899aabbccddeeff00112233445566778899aabbccddeeffccccddddeeee01a0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
			"0x02f84dc98212348256788001c0a000112233445566778899aabbccddeeff00112233445566778899aabbccddeeffa0aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd10",
		},
		{ // Long malform data kept as-is
			"0xdead0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"0xdead0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{ // Short malformed data kept as-is
			"0x80",
			"0x80",
		},
		{ // Short malformed data kept as-is
			"0x00",
			"0x00",
		},
		{ // Legacy account may crash DecodeRLP, but must not crash UnextendSerializedAccount.
			"0xf8448080a00000000000000000000000000000000000000000000000000000000000000000a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
			"0xf8448080a00000000000000000000000000000000000000000000000000000000000000000a0c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
	}

	for _, tc := range testcases {
		unextended := UnextendSerializedAccount(common.FromHex(tc.extended))
		assert.Equal(t, tc.unextended, hexutil.Encode(unextended))
	}
}
