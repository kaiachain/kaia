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

package impl

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{Number: big.NewInt(0)}, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) State() (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

type dummyGovModule struct {
	chainConfig *params.ChainConfig
}

func (d *dummyGovModule) GetParamSet(blockNum uint64) gov.ParamSet {
	return gov.ParamSet{UnitPrice: d.chainConfig.UnitPrice}
}

type AccountKeyPickerForTest struct {
	AddrKeyMap map[common.Address]accountkey.AccountKey
}

func (a *AccountKeyPickerForTest) GetKey(addr common.Address) accountkey.AccountKey {
	return a.AddrKeyMap[addr]
}

func (a *AccountKeyPickerForTest) SetKey(addr common.Address, key accountkey.AccountKey) {
	a.AddrKeyMap[addr] = key
}

func (a *AccountKeyPickerForTest) Exist(addr common.Address) bool {
	return a.AddrKeyMap[addr] != nil
}

func makeTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	if privKey == nil {
		var err error
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	p := &AccountKeyPickerForTest{
		AddrKeyMap: make(map[common.Address]accountkey.AccountKey),
	}
	p.SetKey(addr, accountkey.NewAccountKeyLegacy())

	signer := types.LatestSignerForChainID(big.NewInt(1))
	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, data)
	tx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)

	return tx
}

func makeApproveTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, approveArgs ApproveArgs) *types.Transaction {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}

	data := append([]byte{}, common.Hex2Bytes("095ea7b3")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Spender.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(approveArgs.Amount.Bytes()), 32)...)
	approveTx := makeTx(t, privKey, nonce, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(1), data)

	return approveTx
}

func makeSwapTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, swapArgs SwapArgs) *types.Transaction {
	var err error
	if privKey == nil {
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}

	data := append([]byte{}, common.Hex2Bytes("43bab9f7")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.Token.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountIn.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.MinAmountOut.Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(swapArgs.AmountRepay.Bytes()), 32)...)
	swapTx := makeTx(t, privKey, nonce, common.HexToAddress("0x1234"), big.NewInt(0), 1000000, big.NewInt(1), data)

	return swapTx
}

func flattenPoolTxs(structured map[common.Address]types.Transactions) map[common.Hash]bool {
	flattened := map[common.Hash]bool{}
	for _, txs := range structured {
		for _, tx := range txs {
			flattened[tx.Hash()] = true
		}
	}
	return flattened
}

func flattenBundleTxs(txOrGens []*builder.TxOrGen) ([]common.Hash, error) {
	nodeNonce := uint64(0)
	hashes := []common.Hash{}
	for _, txOrGen := range txOrGens {
		tx, err := txOrGen.GetTx(nodeNonce)
		if err != nil {
			return nil, err
		}
		if txOrGen.IsTxGenerator() {
			nodeNonce += 1
		}
		hashes = append(hashes, tx.Hash())
	}
	return hashes, nil
}

type testTxPool struct {
	statedb *state.StateDB
}

func (pool *testTxPool) GetCurrentState() *state.StateDB {
	return pool.statedb
}
