// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/types/transaction_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
)

// The values in those tests are from the Transaction Tests
// at github.com/ethereum/tests.
var (
	testAddr = common.HexToAddress("b94f5374fce5edbc8e2a8697c15331677e6ebf0b")

	emptyTx = NewTransaction(
		0,
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
		big.NewInt(0), 0, big.NewInt(0),
		nil,
	)

	rightvrsTx, _ = NewTransaction(
		3,
		testAddr,
		big.NewInt(10),
		2000,
		big.NewInt(1),
		common.FromHex("5544"),
	).WithSignature(
		LatestSignerForChainID(common.Big1),
		common.Hex2Bytes("98ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4a8887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a301"),
	)

	accessListTx = TxInternalDataEthereumAccessList{
		ChainID:      big.NewInt(1),
		AccountNonce: 3,
		Recipient:    &testAddr,
		Amount:       big.NewInt(10),
		GasLimit:     25000,
		Price:        big.NewInt(1),
		Payload:      common.FromHex("5544"),
	}

	accessAddr   = common.HexToAddress("0x0000000000000000000000000000000000000001")
	dynamicFeeTx = TxInternalDataEthereumDynamicFee{
		ChainID:      big.NewInt(1),
		AccountNonce: 3,
		Recipient:    &testAddr,
		Amount:       big.NewInt(10),
		GasLimit:     25000,
		GasFeeCap:    big.NewInt(1),
		GasTipCap:    big.NewInt(1),
		Payload:      common.FromHex("5544"),
		AccessList:   AccessList{{Address: accessAddr, StorageKeys: []common.Hash{{0}}}},
	}

	emptyEip2718Tx = &Transaction{
		data: &accessListTx,
	}

	emptyEip1559Tx = &Transaction{
		data: &dynamicFeeTx,
	}

	signedEip2718Tx, _ = emptyEip2718Tx.WithSignature(
		NewEIP2930Signer(big.NewInt(1)),
		common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"),
	)

	signedEip1559Tx, _ = emptyEip1559Tx.WithSignature(
		NewLondonSigner(big.NewInt(1)),
		common.Hex2Bytes("c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b266032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d3752101"))
)

var testKaiaChainConfig = &params.ChainConfig{
	ChainID:                  new(big.Int).SetUint64(111111),
	IstanbulCompatibleBlock:  common.Big0,
	LondonCompatibleBlock:    common.Big0,
	EthTxTypeCompatibleBlock: common.Big0,
	MagmaCompatibleBlock:     common.Big0,
	KoreCompatibleBlock:      common.Big0,
	ShanghaiCompatibleBlock:  common.Big0,
	CancunCompatibleBlock:    common.Big0,
	KaiaCompatibleBlock:      common.Big0,
	UnitPrice:                25000000000, // 25 ston
}

func TestTransactionSigHash(t *testing.T) {
	signer := LatestSignerForChainID(common.Big1)
	if signer.Hash(emptyTx) != common.HexToHash("a715f8447b97e3105d2cc0a8aca1466fa3a02f7cc6d2f9a3fe89f2581c9111c5") {
		t.Errorf("empty transaction hash mismatch, ɡot %x", signer.Hash(emptyTx))
	}
	if signer.Hash(rightvrsTx) != common.HexToHash("bd63ce94e66c7ffbce3b61023bbf9ee6df36047525b123201dcb5c4332f105ae") {
		t.Errorf("RightVRS transaction hash mismatch, ɡot %x", signer.Hash(rightvrsTx))
	}
}

func TestEIP2718TransactionSigHash(t *testing.T) {
	s := NewEIP2930Signer(big.NewInt(1))
	if s.Hash(emptyEip2718Tx) != common.HexToHash("49b486f0ec0a60dfbbca2d30cb07c9e8ffb2a2ff41f29a1ab6737475f6ff69f3") {
		t.Errorf("empty EIP-2718 transaction hash mismatch, got %x", s.Hash(emptyEip2718Tx))
	}
	if s.Hash(signedEip2718Tx) != common.HexToHash("49b486f0ec0a60dfbbca2d30cb07c9e8ffb2a2ff41f29a1ab6737475f6ff69f3") {
		t.Errorf("signed EIP-2718 transaction hash mismatch, got %x", s.Hash(signedEip2718Tx))
	}
}

func TestEIP1559TransactionSigHash(t *testing.T) {
	s := NewLondonSigner(big.NewInt(1))
	if s.Hash(emptyEip1559Tx) != common.HexToHash("a52ce25a7d108740bce8fbb2dfa1f26793b2e8eea94a7700bedbae13cbdd8a0f") {
		t.Errorf("empty EIP-1559 transaction hash mismatch, got %x", s.Hash(emptyEip2718Tx))
	}
	if s.Hash(signedEip1559Tx) != common.HexToHash("a52ce25a7d108740bce8fbb2dfa1f26793b2e8eea94a7700bedbae13cbdd8a0f") {
		t.Errorf("signed EIP-1559 transaction hash mismatch, got %x", s.Hash(signedEip2718Tx))
	}
}

// This test checks signature operations on access list transactions.
func TestEIP2930Signer(t *testing.T) {
	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		keyAddr = crypto.PubkeyToAddress(key.PublicKey)
		signer1 = NewEIP2930Signer(big.NewInt(1))
		signer2 = NewEIP2930Signer(big.NewInt(2))
		tx0     = NewTx(&TxInternalDataEthereumAccessList{AccountNonce: 1, ChainID: new(big.Int)})
		tx1     = NewTx(&TxInternalDataEthereumAccessList{ChainID: big.NewInt(1), AccountNonce: 1, V: new(big.Int), R: new(big.Int), S: new(big.Int)})
		tx2, _  = SignTx(NewTx(&TxInternalDataEthereumAccessList{ChainID: big.NewInt(2), AccountNonce: 1}), signer2, key)
	)

	tests := []struct {
		tx             *Transaction
		signer         Signer
		wantSignerHash common.Hash
		wantSenderErr  error
		wantSignErr    error
		wantHash       common.Hash // after signing
	}{
		{
			tx:             tx0,
			signer:         signer1,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantSenderErr:  ErrInvalidChainId,
			wantHash:       common.HexToHash("1ccd12d8bbdb96ea391af49a35ab641e219b2dd638dea375f2bc94dd290f2549"),
		},
		{
			tx:             tx1,
			signer:         signer1,
			wantSenderErr:  ErrInvalidSig,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantHash:       common.HexToHash("1ccd12d8bbdb96ea391af49a35ab641e219b2dd638dea375f2bc94dd290f2549"),
		},
		{
			// This checks what happens when trying to sign an unsigned tx for the wrong chain.
			tx:             tx1,
			signer:         signer2,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("846ad7672f2a3a40c1f959cd4a8ad21786d620077084d84c8d7c077714caa139"),
			wantSignErr:    ErrInvalidChainId,
		},
		{
			// This checks what happens when trying to re-sign a signed tx for the wrong chain.
			tx:             tx2,
			signer:         signer1,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("367967247499343401261d718ed5aa4c9486583e4d89251afce47f4a33c33362"),
			wantSignErr:    ErrInvalidChainId,
		},
	}

	for i, test := range tests {
		sigHash := test.signer.Hash(test.tx)
		if sigHash != test.wantSignerHash {
			t.Errorf("test %d: wrong sig hash: got %x, want %x", i, sigHash, test.wantSignerHash)
		}
		sender, err := Sender(test.signer, test.tx)
		if err != test.wantSenderErr {
			t.Errorf("test %d: wrong Sender error %q", i, err)
		}
		if err == nil && sender != keyAddr {
			t.Errorf("test %d: wrong sender address %x", i, sender)
		}
		signedTx, err := SignTx(test.tx, test.signer, key)
		if err != test.wantSignErr {
			t.Fatalf("test %d: wrong SignTx error %q", i, err)
		}
		if signedTx != nil {
			if signedTx.Hash() != test.wantHash {
				t.Errorf("test %d: wrong tx hash after signing: got %x, want %x", i, signedTx.Hash(), test.wantHash)
			}
		}
	}
}

func TestHomesteadSigner(t *testing.T) {
	rlpTx := common.Hex2Bytes("f87e8085174876e800830186a08080ad601f80600e600039806000f350fe60003681823780368234f58015156014578182fd5b80825250506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222")

	tx, err := decodeTx(rlpTx)
	assert.NoError(t, err)

	addr, err := EIP155Signer{}.Sender(tx)
	assert.NoError(t, err)
	assert.Equal(t, "0x4c8D290a1B368ac4728d83a9e8321fC3af2b39b1", addr.String())
}

// This test checks signature operations on dynamic fee transactions.
func TestLondonSigner(t *testing.T) {
	var (
		key, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		keyAddr = crypto.PubkeyToAddress(key.PublicKey)
		signer1 = NewLondonSigner(big.NewInt(1))
		signer2 = NewLondonSigner(big.NewInt(2))
		tx0     = NewTx(&TxInternalDataEthereumDynamicFee{AccountNonce: 1, ChainID: new(big.Int)})
		tx1     = NewTx(&TxInternalDataEthereumDynamicFee{ChainID: big.NewInt(1), AccountNonce: 1, V: new(big.Int), R: new(big.Int), S: new(big.Int)})
		tx2, _  = SignTx(NewTx(&TxInternalDataEthereumDynamicFee{ChainID: big.NewInt(2), AccountNonce: 1}), signer2, key)
	)

	tests := []struct {
		tx             *Transaction
		signer         Signer
		wantSignerHash common.Hash
		wantSenderErr  error
		wantSignErr    error
		wantHash       common.Hash // after signing
	}{
		{
			tx:             tx0,
			signer:         signer1,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantSenderErr:  ErrInvalidChainId,
			wantHash:       common.HexToHash("a2c6373b7eed946fd4165a0d8503aa26afc8e99f09e2be58b332fbbedc279f7a"),
		},
		{
			tx:             tx1,
			signer:         signer1,
			wantSenderErr:  ErrInvalidSig,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantHash:       common.HexToHash("a2c6373b7eed946fd4165a0d8503aa26afc8e99f09e2be58b332fbbedc279f7a"),
		},
		{
			// This checks what happens when trying to sign an unsigned tx for the wrong chain.
			tx:             tx1,
			signer:         signer2,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("b6afee4d44e0392fb5d3204b350596d6677440bced7ebd998db73c9671527c57"),
			wantSignErr:    ErrInvalidChainId,
		},
		{
			// This checks what happens when trying to re-sign a signed tx for the wrong chain.
			tx:             tx2,
			signer:         signer1,
			wantSenderErr:  ErrInvalidChainId,
			wantSignerHash: common.HexToHash("b0759fc55582f3e60ded82843dcc17733d8c65f543d2cf2613a47a5c6ac9fc48"),
			wantSignErr:    ErrInvalidChainId,
		},
	}

	for i, test := range tests {
		sigHash := test.signer.Hash(test.tx)
		if sigHash != test.wantSignerHash {
			t.Errorf("test %d: wrong sig hash: got %x, want %x", i, sigHash, test.wantSignerHash)
		}
		sender, err := Sender(test.signer, test.tx)
		if err != test.wantSenderErr {
			t.Errorf("test %d: wrong Sender error %q", i, err)
		}
		if err == nil && sender != keyAddr {
			t.Errorf("test %d: wrong sender address %x", i, sender)
		}
		signedTx, err := SignTx(test.tx, test.signer, key)
		if err != test.wantSignErr {
			t.Fatalf("test %d: wrong SignTx error %q", i, err)
		}
		if signedTx != nil {
			if signedTx.Hash() != test.wantHash {
				t.Errorf("test %d: wrong tx hash after signing: got %x, want %x", i, signedTx.Hash(), test.wantHash)
			}
		}
	}
}

func TestTransactionEncode(t *testing.T) {
	txb, err := rlp.EncodeToBytes(rightvrsTx)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	should := common.FromHex("f86103018207d094b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a82554426a098ff921201554726367d2be8c804a7ff89ccf285ebc57dff8ae4c44b9c19ac4aa08887321be575c8095f789dd4c743dfe42c1820f9231f98a962b210e3ac2452a3")
	if !bytes.Equal(txb, should) {
		t.Errorf("encoded RLP mismatch, ɡot %x", txb)
	}
}

func TestEIP2718TransactionEncode(t *testing.T) {
	// RLP representation
	{
		have, err := rlp.EncodeToBytes(signedEip2718Tx)
		if err != nil {
			t.Fatalf("encode error: %v", err)
		}
		want := common.FromHex("7801f8630103018261a894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a825544c001a0c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b2660a032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d37521")
		if !bytes.Equal(have, want) {
			t.Errorf("encoded RLP mismatch, got %x", have)
		}
	}
}

func TestEIP1559TransactionEncode(t *testing.T) {
	// RLP representation
	{
		have, err := rlp.EncodeToBytes(signedEip1559Tx)
		if err != nil {
			t.Fatalf("encode error: %v", err)
		}
		want := common.FromHex("7802f89d010301018261a894b94f5374fce5edbc8e2a8697c15331677e6ebf0b0a825544f838f7940000000000000000000000000000000000000001e1a0000000000000000000000000000000000000000000000000000000000000000001a0c9519f4f2b30335884581971573fadf60c6204f59a911df35ee8a540456b2660a032f1e8e2c5dd761f9e4f88f41c8310aeaba26a8bfcdacfedfa12ec3862d37521")
		if !bytes.Equal(have, want) {
			t.Errorf("encoded RLP mismatch, got %x", have)
		}
	}
}

func TestEffectiveGasPrice(t *testing.T) {
	gasPrice := big.NewInt(1000)
	gasFeeCap, gasTipCap := big.NewInt(4000), big.NewInt(1000)

	legacyTx := NewTx(&TxInternalDataLegacy{Price: gasPrice})
	dynamicTx := NewTx(&TxInternalDataEthereumDynamicFee{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap})

	header := new(Header)
	have := legacyTx.EffectiveGasPrice(header, testKaiaChainConfig)
	want := gasPrice
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasPrice(header, testKaiaChainConfig)
	te := dynamicTx.GetTxInternalData().(TxInternalDataBaseFee)
	want = te.GetGasFeeCap()
	assert.Equal(t, want, have)

	header.BaseFee = big.NewInt(2000)
	have = legacyTx.EffectiveGasPrice(header, testKaiaChainConfig)
	want = header.BaseFee
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasPrice(header, testKaiaChainConfig)
	want = header.BaseFee
	assert.Equal(t, want, have)

	header.BaseFee = big.NewInt(0)
	have = legacyTx.EffectiveGasPrice(header, params.TestChainConfig)
	want = header.BaseFee
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasPrice(header, params.TestChainConfig)
	want = header.BaseFee
	assert.Equal(t, want, have)
}

func TestEffectiveGasTip(t *testing.T) {
	legacyTx := NewTx(&TxInternalDataLegacy{Price: big.NewInt(1000)})
	dynamicTx := NewTx(&TxInternalDataEthereumDynamicFee{GasFeeCap: big.NewInt(4000), GasTipCap: big.NewInt(1000)})

	// after magma hardfork
	baseFee := big.NewInt(2000)
	have := legacyTx.EffectiveGasTip(baseFee)
	want := big.NewInt(0) // from kaia codebase, legacyTxType also give a tip.
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	// before magma hardfork
	baseFee = nil
	have = legacyTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	have = dynamicTx.EffectiveGasTip(baseFee)
	want = big.NewInt(1000)
	assert.Equal(t, want, have)

	a := new(big.Int)
	assert.Equal(t, 0, a.BitLen())
}

func decodeTx(data []byte) (*Transaction, error) {
	var tx Transaction
	t, err := &tx, rlp.Decode(bytes.NewReader(data), &tx)

	return t, err
}

func defaultTestKey() (*ecdsa.PrivateKey, common.Address) {
	key, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return key, addr
}

func TestRecipientEmpty(t *testing.T) {
	_, addr := defaultTestKey()
	tx, err := decodeTx(common.Hex2Bytes("f84980808080800126a0f18ba0124c1ed46fef6673ff5f614bafbb33a23ad92874bfa3cb3abad56d9a72a046690eb704a07384224d12e991da61faceefede59c6741d85e7d72e097855eaf"))
	if err != nil {
		t.Fatal(err)
	}

	signer := LatestSignerForChainID(common.Big1)

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if addr != from {
		t.Error("derived address doesn't match")
	}
}

func TestRecipientNormal(t *testing.T) {
	_, addr := defaultTestKey()

	tx, err := decodeTx(common.Hex2Bytes("f85d808080940000000000000000000000000000000000000000800126a0c1f2953a2277033c693f3d352b740479788672ba21e76d567557aa069b7e5061a06e798331dbd58c7438fe0e0a64b3b17c8378c726da3613abae8783b5dccc9944"))
	if err != nil {
		t.Fatal(err)
	}

	signer := LatestSignerForChainID(common.Big1)

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}

	if addr != from {
		t.Fatal("derived address doesn't match")
	}
}

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func TestTransactionPriceNonceSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	signer := LatestSignerForChainID(common.Big1)
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups, nil)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != 25*25 {
		t.Errorf("expected %d transactions, found %d", 25*25, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)

			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}

		// If the next tx has different from account, the price must be lower than the current one
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)
			if fromi != fromNext && txi.GasPrice().Cmp(next.GasPrice()) < 0 {
				t.Errorf("invalid ɡasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GasPrice(), i+1, fromNext[:4], next.GasPrice())
			}
		}
	}
}

func TestGasOverflow(t *testing.T) {
	// AccountCreation
	// calculate gas for account creation
	numKeys := new(big.Int).SetUint64(accountkey.MaxNumKeysForMultiSig)
	gasPerKey := new(big.Int).SetUint64(params.TxAccountCreationGasPerKey)
	defaultGas := new(big.Int).SetUint64(params.TxAccountCreationGasDefault)
	txGas := new(big.Int).SetUint64(params.TxGasAccountCreation)
	totalGas := new(big.Int).Add(txGas, new(big.Int).Add(defaultGas, new(big.Int).Mul(numKeys, gasPerKey)))
	assert.Equal(t, true, totalGas.BitLen() <= 64)

	// ValueTransfer
	// calculate gas for validation of multisig accounts.
	gasPerKey = new(big.Int).SetUint64(params.TxValidationGasPerKey)
	defaultGas = new(big.Int).SetUint64(params.TxValidationGasDefault)
	txGas = new(big.Int).SetUint64(params.TxGas)
	totalGas = new(big.Int).Add(txGas, new(big.Int).Add(defaultGas, new(big.Int).Mul(numKeys, gasPerKey)))
	assert.Equal(t, true, totalGas.BitLen() <= 64)

	// TODO-Kaia-Gas: Need to find a way of checking integer overflow for smart contract execution.
}

// TODO-Kaia-FailedTest This test is failed in Kaia
/*
// TestTransactionJSON tests serializing/de-serializing to/from JSON.
func TestTransactionJSON(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	signer := NewEIP155Signer(common.Big1)

	transactions := make([]*Transaction, 0, 50)
	for i := uint64(0); i < 25; i++ {
		var tx *Transaction
		switch i % 2 {
		case 0:
			tx = NewTransaction(i, common.Address{1}, common.Big0, 1, common.Big2, []byte("abcdef"))
		case 1:
			tx = NewContractCreation(i, common.Big0, 1, common.Big2, []byte("abcdef"))
		}
		transactions = append(transactions, tx)

		signedTx, err := SignTx(tx, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}

		transactions = append(transactions, signedTx)
	}

	for _, tx := range transactions {
		data, err := json.Marshal(tx)
		if err != nil {
			t.Fatalf("json.Marshal failed: %v", err)
		}

		var parsedTx *Transaction
		if err := json.Unmarshal(data, &parsedTx); err != nil {
			t.Fatalf("json.Unmarshal failed: %v", err)
		}

		// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
		if tx.Hash() != parsedTx.Hash() {
			t.Errorf("parsed tx differs from original tx, want %v, ɡot %v", tx, parsedTx)
		}
		if tx.ChainId().Cmp(parsedTx.ChainId()) != 0 {
			t.Errorf("invalid chain id, want %d, ɡot %d", tx.ChainId(), parsedTx.ChainId())
		}
	}
}
*/

func TestIntrinsicGas(t *testing.T) {
	// testData contains two kind of members
	// inputString - test input data
	// expectGas - expect gas according to the specific condition.
	//            it differs depending on whether the contract is created or not,
	//            or whether it has passed through the Istanbul compatible block.
	testData := []struct {
		inputString string
		expectGas1  uint64 // contractCreate - false, isIstanbul - false
		expectGas2  uint64 // contractCreate - false, isIstanbul - true
		expectGas3  uint64 // contractCreate - false, isPrague   - true
		expectGas4  uint64 // contractCreate - true,  isIstanbul - false
		expectGas5  uint64 // contractCreate - true,  isIstanbul - true
		expectGas6  uint64 // contractCreate - true, isPrague   - true
	}{
		{"0000", 21008, 21200, 21008, 53008, 53200, 53010},
		{"1000", 21072, 21200, 21020, 53072, 53200, 53022},
		{"0100", 21072, 21200, 21020, 53072, 53200, 53022},
		{"ff3d", 21136, 21200, 21032, 53136, 53200, 53034},
		{"0000a6bc", 21144, 21400, 21040, 53144, 53400, 53042},
		{"fd00fd00", 21144, 21400, 21040, 53144, 53400, 53042},
		{"", 21000, 21000, 21000, 53000, 53000, 53000},
	}
	for _, tc := range testData {
		var (
			data []byte // input data entered through the tx argument
			gas  uint64 // the gas varies depending on what comes in as a condition(contractCreate & IsIstanbulForkEnabled)
			err  error  // in this unittest, every testcase returns nil error.
		)

		data, err = hex.DecodeString(tc.inputString) // decode input string to hex data
		assert.Equal(t, nil, err)

		// TODO-Kaia: Add test for EIP-7623
		gas, err = IntrinsicGas(data, nil, nil, false, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas1, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, nil, false, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas2, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, nil, false, params.Rules{IsIstanbul: true, IsShanghai: true, IsPrague: true})
		assert.Equal(t, tc.expectGas3, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, nil, true, params.Rules{IsIstanbul: false})
		assert.Equal(t, tc.expectGas4, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, nil, true, params.Rules{IsIstanbul: true})
		assert.Equal(t, tc.expectGas5, gas)
		assert.Equal(t, nil, err)

		gas, err = IntrinsicGas(data, nil, nil, true, params.Rules{IsIstanbul: true, IsShanghai: true, IsPrague: true})
		assert.Equal(t, tc.expectGas6, gas)
		assert.Equal(t, nil, err)
	}
}

// Tests that if multiple transactions have the same price, the ones seen earlier
// are prioritized to avoid network spam attacks aiming for a specific ordering.
func TestTransactionTimeSort(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)

		tx, _ := SignTx(NewTransaction(0, common.Address{}, big.NewInt(100), 100, big.NewInt(1), nil), signer, key)
		tx.time = time.Unix(0, int64(len(keys)-start))

		groups[addr] = append(groups[addr], tx)
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups, nil)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != len(keys) {
		t.Errorf("expected %d transactions, found %d", len(keys), len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)

			if txi.GasPrice().Cmp(next.GasPrice()) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GasPrice(), i+1, fromNext[:4], next.GasPrice())
			}
			// Make sure time order is ascending if the txs have the same gas price
			if txi.GasPrice().Cmp(next.GasPrice()) == 0 && txi.time.After(next.time) {
				t.Errorf("invalid received time ordering: tx #%d (A=%x T=%v) > tx #%d (A=%x T=%v)", i, fromi[:4], txi.time, i+1, fromNext[:4], next.time)
			}
		}
	}
}

// TestTransactionTimeSortDifferentGasPrice tests that although multiple transactions have the different price, the ones seen earlier
// are prioritized to avoid network spam attacks aiming for a specific ordering.
func TestTransactionTimeSortDifferentGasPrice(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 5)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := LatestSignerForChainID(big.NewInt(1))

	// Generate a batch of transactions with overlapping prices, but different creation times
	groups := map[common.Address]Transactions{}
	gasPrice := big.NewInt(1)
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)

		tx, _ := SignTx(NewTransaction(0, common.Address{}, big.NewInt(100), 100, gasPrice, nil), signer, key)
		tx.time = time.Unix(0, int64(len(keys)-start))

		groups[addr] = append(groups[addr], tx)

		gasPrice = gasPrice.Add(gasPrice, big.NewInt(1))
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups, nil)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != len(keys) {
		t.Errorf("expected %d transactions, found %d", len(keys), len(txs))
	}
	for i, tx := range txs {
		from, _ := Sender(signer, tx)
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)

			// Make sure time order is ascending.
			if tx.time.After(next.time) {
				t.Errorf("invalid received time ordering: tx #%d (A=%x T=%v) > tx #%d (A=%x T=%v)", i, from[:4], tx.time, i+1, fromNext[:4], next.time)
			}
		}
	}
}

// TestTransactionCoding tests serializing/de-serializing to/from rlp and JSON.
func TestTransactionCoding(t *testing.T) {
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("could not generate key: %v", err)
	}
	var (
		signer     = LatestSignerForChainID(common.Big1)
		addr       = common.HexToAddress("0x0000000000000000000000000000000000000001")
		recipient  = common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
		accesses   = AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}}
		txDataList = []func(uint64) TxInternalData{
			func(i uint64) TxInternalData {
				// Legacy tx.
				return &TxInternalDataLegacy{
					AccountNonce: i,
					Recipient:    &recipient,
					GasLimit:     1,
					Price:        big.NewInt(2),
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Legacy tx contract creation.
				return &TxInternalDataLegacy{
					AccountNonce: i,
					GasLimit:     1,
					Price:        big.NewInt(2),
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Tx with non-zero access list.
				return &TxInternalDataEthereumAccessList{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					Recipient:    &recipient,
					GasLimit:     123457,
					Price:        big.NewInt(10),
					AccessList:   accesses,
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Tx with empty access list.
				return &TxInternalDataEthereumAccessList{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					Recipient:    &recipient,
					GasLimit:     123457,
					Price:        big.NewInt(10),
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Contract creation with access list.
				return &TxInternalDataEthereumAccessList{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					GasLimit:     123457,
					Price:        big.NewInt(10),
					AccessList:   accesses,
				}
			},
			func(i uint64) TxInternalData {
				// Tx with non-zero access list.
				return &TxInternalDataEthereumDynamicFee{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					Recipient:    &recipient,
					GasLimit:     123457,
					GasFeeCap:    big.NewInt(10),
					GasTipCap:    big.NewInt(10),
					AccessList:   accesses,
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Tx with dynamic fee.
				return &TxInternalDataEthereumDynamicFee{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					Recipient:    &recipient,
					GasLimit:     123457,
					GasFeeCap:    big.NewInt(10),
					GasTipCap:    big.NewInt(10),
					Payload:      []byte("abcdef"),
				}
			},
			func(i uint64) TxInternalData {
				// Contract creation with dynamic fee tx.
				return &TxInternalDataEthereumDynamicFee{
					ChainID:      big.NewInt(1),
					AccountNonce: i,
					GasLimit:     123457,
					GasFeeCap:    big.NewInt(10),
					GasTipCap:    big.NewInt(10),
					AccessList:   accesses,
				}
			},
			func(i uint64) TxInternalData {
				// Tx with non-zero access list.
				return &TxInternalDataEthereumSetCode{
					ChainID:           uint256.NewInt(1),
					AccountNonce:      i,
					Recipient:         recipient,
					GasLimit:          123457,
					GasFeeCap:         big.NewInt(10),
					GasTipCap:         big.NewInt(10),
					AccessList:        accesses,
					AuthorizationList: authorizations,
				}
			},
			func(i uint64) TxInternalData {
				// Tx with set code.
				return &TxInternalDataEthereumSetCode{
					ChainID:           uint256.NewInt(1),
					AccountNonce:      i,
					Recipient:         recipient,
					GasLimit:          123457,
					GasFeeCap:         big.NewInt(10),
					GasTipCap:         big.NewInt(10),
					AuthorizationList: authorizations,
				}
			},
		}
	)
	for i := 0; i < 500; i++ {
		txData := txDataList[i%len(txDataList)](uint64(i))
		transaction := Transaction{data: txData}
		tx, err := SignTx(&transaction, signer, key)
		if err != nil {
			t.Fatalf("could not sign transaction: %v", err)
		}
		// RLP
		parsedTx, err := encodeDecodeBinary(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)

		// JSON
		parsedTx, err = encodeDecodeJSON(tx)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(parsedTx, tx)
	}
}

func encodeDecodeJSON(tx *Transaction) (*Transaction, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return nil, fmt.Errorf("json encoding failed: %v", err)
	}
	parsedTx := &Transaction{}
	if err := json.Unmarshal(data, &parsedTx); err != nil {
		return nil, fmt.Errorf("json decoding failed: %v", err)
	}
	return parsedTx, nil
}

func encodeDecodeBinary(tx *Transaction) (*Transaction, error) {
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("rlp encoding failed: %v", err)
	}
	parsedTx := &Transaction{}
	if err := parsedTx.UnmarshalBinary(data); err != nil {
		return nil, fmt.Errorf("rlp decoding failed: %v", err)
	}
	return parsedTx, nil
}

func assertEqual(orig *Transaction, cpy *Transaction) error {
	// compare nonce, price, gaslimit, recipient, amount, payload, V, R, S
	if want, got := orig.Hash(), cpy.Hash(); want != got {
		return fmt.Errorf("parsed tx differs from original tx, want %v, got %v", want, got)
	}
	if want, got := orig.ChainId(), cpy.ChainId(); want.Cmp(got) != 0 {
		return fmt.Errorf("invalid chain id, want %d, got %d", want, got)
	}

	if orig.Type().IsEthTypedTransaction() && cpy.Type().IsEthTypedTransaction() {
		tOrig := orig.data.(TxInternalDataEthTyped)
		tCpy := cpy.data.(TxInternalDataEthTyped)

		if !reflect.DeepEqual(tOrig.GetAccessList(), tCpy.GetAccessList()) {
			return fmt.Errorf("access list wrong!")
		}
	}

	return nil
}

func TestIsSorted(t *testing.T) {
	signer := LatestSignerForChainID(big.NewInt(1))

	key, _ := crypto.GenerateKey()
	batches := make(txByPriceAndTime, 10)

	for i := 0; i < 10; i++ {
		tx, _ := SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(i)), nil), signer, key)
		txWithFee, _ := newTxWithMinerFee(tx, common.Address{}, big.NewInt(int64(i)))
		batches[i] = txWithFee
	}

	// Shuffle transactions.
	rand.Shuffle(len(batches), func(i, j int) {
		batches[i], batches[j] = batches[j], batches[i]
	})

	sort.Sort(txByPriceAndTime(batches))
	assert.True(t, sort.IsSorted(txByPriceAndTime(batches)))
}

func TestFilterTransactionWithBaseFee(t *testing.T) {
	signer := LatestSignerForChainID(big.NewInt(1))

	pending := make(map[common.Address]Transactions)
	keys := make([]*ecdsa.PrivateKey, 3)

	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}

	from1 := crypto.PubkeyToAddress(keys[0].PublicKey)
	txs1 := make(Transactions, 3)
	txs1[0], _ = SignTx(NewTransaction(uint64(0), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[0])
	txs1[1], _ = SignTx(NewTransaction(uint64(1), common.Address{}, big.NewInt(100), 100, big.NewInt(40), nil), signer, keys[0])
	txs1[2], _ = SignTx(NewTransaction(uint64(2), common.Address{}, big.NewInt(100), 100, big.NewInt(50), nil), signer, keys[0])
	pending[from1] = txs1

	from2 := crypto.PubkeyToAddress(keys[1].PublicKey)
	txs2 := make(Transactions, 4)
	txs2[0], _ = SignTx(NewTransaction(uint64(0), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[1])
	txs2[1], _ = SignTx(NewTransaction(uint64(1), common.Address{}, big.NewInt(100), 100, big.NewInt(20), nil), signer, keys[1])
	txs2[2], _ = SignTx(NewTransaction(uint64(2), common.Address{}, big.NewInt(100), 100, big.NewInt(40), nil), signer, keys[1])
	txs2[3], _ = SignTx(NewTransaction(uint64(3), common.Address{}, big.NewInt(100), 100, big.NewInt(40), nil), signer, keys[1])
	pending[from2] = txs2

	from3 := crypto.PubkeyToAddress(keys[2].PublicKey)
	txs3 := make(Transactions, 5)
	txs3[0], _ = SignTx(NewTransaction(uint64(0), common.Address{}, big.NewInt(100), 100, big.NewInt(10), nil), signer, keys[2])
	txs3[1], _ = SignTx(NewTransaction(uint64(1), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[2])
	txs3[2], _ = SignTx(NewTransaction(uint64(2), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[2])
	txs3[3], _ = SignTx(NewTransaction(uint64(3), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[2])
	txs3[4], _ = SignTx(NewTransaction(uint64(4), common.Address{}, big.NewInt(100), 100, big.NewInt(30), nil), signer, keys[2])
	pending[from3] = txs3

	baseFee := big.NewInt(30)
	pending = FilterTransactionWithBaseFee(pending, baseFee)

	assert.Equal(t, len(pending[from1]), 3)
	for i := 0; i < len(pending[from1]); i++ {
		assert.Equal(t, txs1[i], pending[from1][i])
	}

	assert.Equal(t, len(pending[from2]), 1)
	for i := 0; i < len(pending[from2]); i++ {
		assert.Equal(t, txs2[i], pending[from2][i])
	}

	assert.Equal(t, len(pending[from3]), 0)
}

// go test -bench=BenchmarkSortTxsByPriceAndTime -benchtime=10x
func BenchmarkSortTxsByPriceAndTime20000(b *testing.B) { benchmarkSortTxsByPriceAndTime(b, 20000) }
func BenchmarkSortTxsByPriceAndTime10000(b *testing.B) { benchmarkSortTxsByPriceAndTime(b, 10000) }
func BenchmarkSortTxsByPriceAndTime100(b *testing.B)   { benchmarkSortTxsByPriceAndTime(b, 100) }
func benchmarkSortTxsByPriceAndTime(b *testing.B, size int) {
	signer := LatestSignerForChainID(big.NewInt(1))

	key, _ := crypto.GenerateKey()

	// make size to be even
	if size%2 == 1 {
		size += 1
	}
	txs := make(Transactions, size)

	for i := 0; i < size; i += 2 {
		gasFeeCap := rand.Int63n(50)
		txs[i], _ = SignTx(NewTransaction(uint64(i), common.Address{}, big.NewInt(100), 100, big.NewInt(25*params.Gkei), nil), signer, key)
		txs[i+1], _ = SignTx(NewTx(&TxInternalDataEthereumDynamicFee{
			AccountNonce: uint64(i),
			Recipient:    &common.Address{},
			Amount:       big.NewInt(100),
			GasLimit:     100,
			GasFeeCap:    big.NewInt(int64(25*params.Gkei) + gasFeeCap),
			GasTipCap:    big.NewInt(gasFeeCap),
			Payload:      nil,
		}), signer, key)
	}

	// Benchmark importing the transactions into the queue
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Shuffle(size, func(i, j int) {
			txs[i], txs[j] = txs[j], txs[i]
		})
		txs = SortTxsByPriceAndTime(txs, big.NewInt(25*params.Gkei))
	}
}

// go test -bench=BenchmarkTxSortByPriceAndTime -benchtime=10x
func BenchmarkTxSortByPriceAndTime20000(b *testing.B) { benchmarkTxSortByPriceAndTime(b, 20000) }
func BenchmarkTxSortByPriceAndTime10000(b *testing.B) { benchmarkTxSortByPriceAndTime(b, 10000) }
func BenchmarkTxSortByPriceAndTime100(b *testing.B)   { benchmarkTxSortByPriceAndTime(b, 100) }
func benchmarkTxSortByPriceAndTime(b *testing.B, size int) {
	signer := LatestSignerForChainID(big.NewInt(1))

	key, _ := crypto.GenerateKey()
	batches := make(txByPriceAndTime, size)

	for i := 0; i < size; i++ {
		gasFeeCap := rand.Int63n(50)
		tx, _ := SignTx(NewTx(&TxInternalDataEthereumDynamicFee{
			AccountNonce: uint64(i),
			Recipient:    &common.Address{},
			Amount:       big.NewInt(100),
			GasLimit:     100,
			GasFeeCap:    big.NewInt(int64(25*params.Gkei) + gasFeeCap),
			GasTipCap:    big.NewInt(gasFeeCap),
			Payload:      nil,
		}), signer, key)
		txWithFee, _ := newTxWithMinerFee(tx, common.Address{}, big.NewInt(25*params.Gkei))
		batches[i] = txWithFee
	}

	// Benchmark importing the transactions into the queue
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rand.Shuffle(len(batches), func(i, j int) {
			batches[i], batches[j] = batches[j], batches[i]
		})
		sort.Sort(batches)
	}
}

func TestTransactionPriceNonceSortLegacy(t *testing.T) {
	t.Parallel()
	testTransactionPriceNonceSort(t, nil)
}

func TestTransactionPriceNonceSort1559(t *testing.T) {
	t.Parallel()
	testTransactionPriceNonceSort(t, big.NewInt(0))
	testTransactionPriceNonceSort(t, big.NewInt(5))
	testTransactionPriceNonceSort(t, big.NewInt(50))
}

// Tests that transactions can be correctly sorted according to their price in
// decreasing order, but at the same time with increasing nonces when issued by
// the same account.
func testTransactionPriceNonceSort(t *testing.T, baseFee *big.Int) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	signer := LatestSignerForChainID(common.Big1)

	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	expectedCount := 0
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		count := 25
		for i := 0; i < 25; i++ {
			var tx *Transaction
			gasFeeCap := rand.Intn(50)
			if baseFee == nil {
				tx = NewTx(&TxInternalDataLegacy{
					AccountNonce: uint64(start + i),
					Recipient:    &common.Address{},
					Amount:       big.NewInt(100),
					GasLimit:     100,
					Price:        big.NewInt(int64(gasFeeCap)),
					Payload:      nil,
				})
			} else {
				tx = NewTx(&TxInternalDataEthereumDynamicFee{
					AccountNonce: uint64(start + i),
					Recipient:    &common.Address{},
					Amount:       big.NewInt(100),
					GasLimit:     100,
					GasFeeCap:    big.NewInt(int64(gasFeeCap)),
					GasTipCap:    big.NewInt(int64(rand.Intn(gasFeeCap + 1))),
					Payload:      nil,
				})
				if count == 25 && int64(gasFeeCap) < baseFee.Int64() {
					count = i
				}
			}
			tx, err := SignTx(tx, signer, key)
			if err != nil {
				t.Fatalf("failed to sign tx: %s", err)
			}
			groups[addr] = append(groups[addr], tx)
		}
		expectedCount += count
	}
	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups, baseFee)

	txs := Transactions{}
	for tx := txset.Peek(); tx != nil; tx = txset.Peek() {
		txs = append(txs, tx)
		txset.Shift()
	}
	if len(txs) != expectedCount {
		t.Errorf("expected %d transactions, found %d", expectedCount, len(txs))
	}
	for i, txi := range txs {
		fromi, _ := Sender(signer, txi)

		// Make sure the nonce order is valid
		for j, txj := range txs[i+1:] {
			fromj, _ := Sender(signer, txj)
			if fromi == fromj && txi.Nonce() > txj.Nonce() {
				t.Errorf("invalid nonce ordering: tx #%d (A=%x N=%v) < tx #%d (A=%x N=%v)", i, fromi[:4], txi.Nonce(), i+j, fromj[:4], txj.Nonce())
			}
		}
		// If the next tx has different from account, the price must be lower than the current one
		if i+1 < len(txs) {
			next := txs[i+1]
			fromNext, _ := Sender(signer, next)
			tip := txi.EffectiveGasTip(baseFee)
			nextTip := next.EffectiveGasTip(baseFee)
			if fromi != fromNext && tip.Cmp(nextTip) < 0 {
				t.Errorf("invalid gasprice ordering: tx #%d (A=%x P=%v) < tx #%d (A=%x P=%v)", i, fromi[:4], txi.GasPrice(), i+1, fromNext[:4], next.GasPrice())
			}
		}
	}
}

func TestEmptyHeap(t *testing.T) {
	heap := NewTransactionsByPriceAndNonce(nil, nil, nil)
	assert.Nil(t, heap.Peek())
	assert.True(t, heap.Empty())
	assert.NotPanics(t, func() {
		heap.Shift()
		heap.Pop()
		heap.Clear()
	})
}

func TestHeapCopy(t *testing.T) {
	// Generate a batch of accounts to start with
	keys := make([]*ecdsa.PrivateKey, 25)
	for i := 0; i < len(keys); i++ {
		keys[i], _ = crypto.GenerateKey()
	}
	baseFee := common.Big0

	signer := LatestSignerForChainID(common.Big1)
	// Generate a batch of transactions with overlapping values, but shifted nonces
	groups := map[common.Address]Transactions{}
	for start, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)
		for i := 0; i < 25; i++ {
			tx, _ := SignTx(NewTransaction(uint64(start+i), common.Address{}, big.NewInt(100), 100, big.NewInt(int64(start+i)), nil), signer, key)
			groups[addr] = append(groups[addr], tx)
		}
	}

	// Sort the transactions and cross check the nonce ordering
	txset := NewTransactionsByPriceAndNonce(signer, groups, baseFee)
	txsetCopy := txset.Copy()

	assert.False(t, txset.Empty())
	assert.False(t, txsetCopy.Empty())
	for i := 0; i < 25; i++ {
		txset.Pop()
	}
	assert.True(t, txset.Empty())
	assert.False(t, txsetCopy.Empty())
}
