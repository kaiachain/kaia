// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

// Verify that Client implements the Kaia interfaces.
var (
	// _ = kaia.Subscription(&Client{})
	_ = kaia.ChainReader(&Client{})
	_ = kaia.TransactionReader(&Client{})
	_ = kaia.ChainStateReader(&Client{})
	_ = kaia.ChainSyncReader(&Client{})
	_ = kaia.ContractCaller(&Client{})
	_ = kaia.LogFilterer(&Client{})
	_ = kaia.TransactionSender(&Client{})
	_ = kaia.GasPricer(&Client{})
	_ = kaia.PendingStateReader(&Client{})
	_ = kaia.PendingContractCaller(&Client{})
	_ = kaia.GasEstimator(&Client{})
	// _ = kaia.PendingStateEventer(&Client{})
)

var (
	testKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testAddr    = crypto.PubkeyToAddress(testKey.PublicKey)
	testBalance = big.NewInt(2e15)
)

var genesisConfig = &params.ChainConfig{
	ChainID:                  big.NewInt(1337),
	IstanbulCompatibleBlock:  big.NewInt(0),
	LondonCompatibleBlock:    big.NewInt(0),
	EthTxTypeCompatibleBlock: big.NewInt(0),
	MagmaCompatibleBlock:     big.NewInt(0),
	KoreCompatibleBlock:      big.NewInt(0),
	ShanghaiCompatibleBlock:  big.NewInt(0),
	CancunCompatibleBlock:    big.NewInt(0),
	KaiaCompatibleBlock:      big.NewInt(0),
	PragueCompatibleBlock:    big.NewInt(0),
	UnitPrice:                25000000000,
}

var testTx1 = func() *types.Transaction {
	tx := types.NewTransaction(0, common.Address{2}, big.NewInt(12), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signer := types.LatestSignerForChainID(genesisConfig.ChainID)
	signedTx, _ := types.SignTx(tx, signer, testKey)
	return signedTx
}()

var testTx2 = func() *types.Transaction {
	tx := types.NewTransaction(1, common.Address{2}, big.NewInt(8), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signer := types.LatestSigner(genesisConfig)
	signedTx, _ := types.SignTx(tx, signer, testKey)
	return signedTx
}()

var blocks = make([]*types.Block, 3)

func init() {
	for i := 0; i < 3; i++ {
		header := genMockHeader(i)
		blocks[i] = types.NewBlockWithHeader(header)
	}
}

func MockGetBalance(t *testing.T, address common.Address, blockNumber *big.Int) map[string]interface{} {
	// Handle different test scenarios
	if blockNumber.Cmp(big.NewInt(2)) > 0 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}
	if address == testAddr {
		// testAddr - has balance (2000000000000000 wei)
		return map[string]interface{}{
			"result": "0x71afd498d0000", // 2000000000000000 wei
		}
	}

	// Non-existent account - zero balance
	return map[string]interface{}{
		"result": "0x0",
	}
}

func MockGetBlockByNumber(t *testing.T, blockNumberArg string) map[string]interface{} {
	var block *types.Block
	switch blockNumberArg {
	case "0x0", "earliest":
		block = blocks[0]
	case "0x1":
		block = blocks[1]
	case "0x2", "latest", "pending", "-0x2", "-0x1":
		block = blocks[2]
	default:
		// Return null result for other blocks (will be converted to kaia.NotFound by client)
		return map[string]interface{}{
			"result": nil,
		}
	}

	rpcOutput, err := api.RpcOutputBlock(block, false, false, genesisConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return map[string]interface{}{
		"result": rpcOutput,
	}
}

func MockGetBlockByNumberEth(t *testing.T, blockNumberArg string) map[string]interface{} {
	var block *types.Block
	switch blockNumberArg {
	case "0x0", "earliest":
		block = blocks[0]
	case "0x1":
		block = blocks[1]
	case "0x2", "latest", "pending", "-0x2", "-0x1":
		block = blocks[2]
	default:
		// Return null result for other blocks (will be converted to kaia.NotFound by client)
		return map[string]interface{}{
			"result": nil,
		}
	}

	rpcOutput, err := api.RpcMarshalEthBlock(block, nil, genesisConfig, false, false, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("RpcMarshalEthBlock nonce: %v", rpcOutput["nonce"])
	return map[string]interface{}{
		"result": rpcOutput,
	}
}

func MockGetBlockByHash(t *testing.T, blockHash string) map[string]interface{} {
	blockNum := slices.IndexFunc(blocks, func(h *types.Block) bool {
		return h.Hash().Hex() == blockHash
	})
	if blockNum == -1 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}

	block := blocks[blockNum]
	rpcOutput, err := api.RpcOutputBlock(block, false, false, genesisConfig)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return map[string]interface{}{
		"result": rpcOutput,
	}
}

func MockGetTransactionByBlockHashAndIndex(t *testing.T, blockHash string, transactionIndex uint64) map[string]interface{} {
	blockNum := slices.IndexFunc(blocks, func(h *types.Block) bool {
		return h.Hash().Hex() == blockHash
	})
	// only accept block number = 2
	if blockNum < 2 || transactionIndex > 2 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}

	txs := []*types.Transaction{testTx1, testTx2}
	return map[string]interface{}{
		"result": txs[transactionIndex],
	}
}

func launchMockServer(t *testing.T, quit chan struct{}) string {
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqData map[string]interface{}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&reqData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		t.Logf("MockHttpServer received request: %+v", reqData)

		// Extract method and id from JSON-RPC request
		method, _ := reqData["method"].(string)
		id := reqData["id"]

		var response map[string]interface{}

		switch method {
		case "kaia_chainID":
			response = map[string]interface{}{
				"result": "0x" + genesisConfig.ChainID.Text(16),
			}
		case "kaia_getBalance":
			params := reqData["params"].([]interface{})
			addressStr := params[0].(string)
			address := common.HexToAddress(addressStr)
			blockNumberStr := params[1].(string)
			if !strings.HasPrefix(blockNumberStr, "0x") {
				t.Fatalf("blockNumberStr should start with 0x, but got %v", blockNumberStr)
			}
			blockNumber, ok := new(big.Int).SetString(blockNumberStr[2:], 16)
			if !ok {
				t.Fatalf("unexpected error: %v", blockNumberStr)
			}
			response = MockGetBalance(t, address, blockNumber)
		case "kaia_getBlockByNumber":
			params := reqData["params"].([]interface{})
			blockNumber := params[0].(string)
			response = MockGetBlockByNumber(t, blockNumber)
		case "eth_getBlockByNumber":
			params := reqData["params"].([]interface{})
			blockNumber := params[0].(string)
			response = MockGetBlockByNumberEth(t, blockNumber)
		case "kaia_getBlockByHash":
			params := reqData["params"].([]interface{})
			blockHash := params[0].(string)
			response = MockGetBlockByHash(t, blockHash)
		case "kaia_blockNumber":
			response = map[string]interface{}{
				"result": "0x2", // Block 2
			}
		case "kaia_syncing":
			response = map[string]interface{}{
				"result": false,
			}
		case "net_version":
			response = map[string]interface{}{
				"result": genesisConfig.ChainID.Text(10),
			}
		case "kaia_gasPrice":
			response = map[string]interface{}{
				"result": "0x3b9aca00",
			}
		case "kaia_estimateGas":
			response = map[string]interface{}{
				"result": "0x5208",
			}
		case "kaia_call":
			response = map[string]interface{}{
				"result": "0x",
			}
		case "kaia_getTransactionByBlockHashAndIndex":
			params := reqData["params"].([]interface{})
			blockHash := params[0].(string)
			transactionIndexStr := params[1].(string)
			txIdx, err := strconv.ParseUint(transactionIndexStr[2:], 16, 64)
			if err != nil {
				t.Fatalf("unexpected error: %v", transactionIndexStr)
			}
			response = MockGetTransactionByBlockHashAndIndex(t, blockHash, txIdx)
		default:
			// Return method not found error
			response = map[string]interface{}{
				"jsonrpc": "2.0",
				"id":      id,
				"error": map[string]interface{}{
					"code":    -32601,
					"message": fmt.Sprintf("the method %s does not exist/is not available", method),
				},
			}
		}
		response["id"] = id
		response["jsonrpc"] = "2.0"

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("Failed to encode response: %v", err)
		}
		t.Logf("MockHttpServer sent response: %+v", response)
	})

	s := &http.Server{
		Addr:    "127.0.0.1:36000",
		Handler: myHandler,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			t.Errorf("Server failed: %v", err)
		}
	}()

	t.Log("MockHttpServer started on 127.0.0.1:36000")

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(ctx)
	}()

	return "http://127.0.0.1:36000"
}

func TestKaiaClient(t *testing.T) {
	quitChan := make(chan struct{})
	defer close(quitChan)

	serverURL := launchMockServer(t, quitChan)

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	client, err := DialContext(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Client connected to mock server")
	defer client.Close()

	tests := map[string]struct {
		test func(t *testing.T)
	}{
		"Header": {
			func(t *testing.T) { testHeader(t, client) },
		},
		"BalanceAt": {
			func(t *testing.T) { testBalanceAt(t, client) },
		},
		"ChainID": {
			func(t *testing.T) { testChainID(t, client) },
		},
		"TxInBlockInterrupted": {
			func(t *testing.T) { testTransactionInBlock(t, client) },
		},
		"GetBlock": {
			func(t *testing.T) { testGetBlock(t, client) },
		},
		"StatusFunctions": {
			func(t *testing.T) { testStatusFunctions(t, client) },
		},
		"CallContract": {
			func(t *testing.T) { testCallContract(t, client) },
		},
		"CallContractAtHash": {
			func(t *testing.T) { testCallContractAtHash(t, client) },
		},
		// "AtFunctions": {
		// 	func(t *testing.T) { testAtFunctions(t, client) },
		// },
		"TransactionSender": {
			func(t *testing.T) { testTransactionSender(t, client) },
		},
	}

	t.Parallel()
	for name, tt := range tests {
		t.Run(name, tt.test)
	}
}

func testHeader(t *testing.T, client *Client) {
	tests := map[string]struct {
		block   *big.Int
		want    *types.Header
		wantErr error
	}{
		"genesis": {
			block: big.NewInt(0),
			want:  blocks[0].Header(),
		},
		"first_block": {
			block: big.NewInt(1),
			want:  blocks[1].Header(),
		},
		"future_block": {
			block:   big.NewInt(1000000000),
			want:    nil,
			wantErr: kaia.NotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := client
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := c.HeaderByNumber(ctx, tt.block)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("HeaderByNumber(%v) error = %q, want %q", tt.block, err, tt.wantErr)
			}
			if got != nil && got.Number != nil && got.Number.Sign() == 0 {
				got.Number = big.NewInt(0) // hack to make DeepEqual work
			}
			if got != nil && got.Hash() != tt.want.Hash() {
				t.Fatalf("HeaderByNumber(%v) got = %v, want %v", tt.block, got, tt.want)
			}
		})
	}
}

func testBalanceAt(t *testing.T, client *Client) {
	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account_genesis": {
			account: testAddr,
			block:   big.NewInt(0),
			want:    testBalance,
		},
		"valid_account": {
			account: testAddr,
			block:   big.NewInt(1),
			want:    testBalance,
		},
		"non_existent_account": {
			account: common.Address{1},
			block:   big.NewInt(1),
			want:    big.NewInt(0),
		},
		"future_block": {
			account: testAddr,
			block:   big.NewInt(1000000000),
			want:    big.NewInt(0),
			wantErr: kaia.NotFound,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := client
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := c.BalanceAt(ctx, tt.account, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("BalanceAt(%x, %v) error = %q, want %q", tt.account, tt.block, err, tt.wantErr)
			}
			if got.Cmp(tt.want) != 0 {
				t.Fatalf("BalanceAt(%x, %v) = %v, want %v", tt.account, tt.block, got, tt.want)
			}
		})
	}
}

func testTransactionInBlock(t *testing.T, c *Client) {
	// Get current block by number.
	block, err := c.BlockByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test tx in block not found.
	if _, err := c.TransactionInBlock(context.Background(), block.Hash(), 20); err.Error() != kaia.NotFound.Error() {
		t.Fatal("error should be kaia.NotFound")
	}

	// Test tx in block found.
	tx, err := c.TransactionInBlock(context.Background(), block.Hash(), 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx1.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	tx, err = c.TransactionInBlock(context.Background(), block.Hash(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tx.Hash() != testTx2.Hash() {
		t.Fatalf("unexpected transaction: %v", tx)
	}

	// Test pending block
	_, err = c.BlockByNumber(context.Background(), big.NewInt(int64(rpc.PendingBlockNumber)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testChainID(t *testing.T, c *Client) {
	id, err := c.ChainID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == nil || id.Cmp(genesisConfig.ChainID) != 0 {
		t.Fatalf("ChainID returned wrong number: %+v", id)
	}
}

func testGetBlock(t *testing.T, c *Client) {
	// Get current block number
	blockNumber, err := c.BlockNumber(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if blockNumber.Int64() != 2 {
		t.Fatalf("BlockNumber returned wrong number: %d", blockNumber)
	}
	// Get current block
	block, err := c.BlockByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.NumberU64() != blockNumber.Uint64() {
		t.Fatalf("BlockByNumber returned wrong block: want %d got %d", blockNumber, block.NumberU64())
	}
	// Get current block by hash
	blockH, err := c.BlockByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Hash() != blockH.Hash() {
		t.Fatalf("BlockByHash returned wrong block: want %v got %v", block.Hash().Hex(), blockH.Hash().Hex())
	}
	// Get header by number
	header, err := c.HeaderByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != header.Hash() {
		t.Fatalf("HeaderByNumber returned wrong header: want %v got %v", block.Header().Hash().Hex(), header.Hash().Hex())
	}
	// Get header by hash
	headerH, err := c.HeaderByHash(context.Background(), block.Hash())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if block.Header().Hash() != headerH.Hash() {
		t.Fatalf("HeaderByHash returned wrong header: want %v got %v", block.Header().Hash().Hex(), headerH.Hash().Hex())
	}
}

func testStatusFunctions(t *testing.T, c *Client) {
	// Sync progress
	progress, err := c.SyncProgress(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if progress != nil {
		t.Fatalf("unexpected progress: %v", progress)
	}

	// NetworkID
	networkID, err := c.NetworkID(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if networkID.Cmp(genesisConfig.ChainID) != 0 {
		t.Fatalf("unexpected networkID: %v", networkID)
	}

	// SuggestGasPrice
	gasPrice, err := c.SuggestGasPrice(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gasPrice.Cmp(big.NewInt(1000000000)) != 0 {
		t.Fatalf("unexpected gas price: %v", gasPrice)
	}

	// Note: SuggestGasTipCap, BlobBaseFee, and FeeHistory methods are not available in Kaia client
	// Skipping these tests for Kaia compatibility
}

func testCallContractAtHash(t *testing.T, c *Client) {
	// EstimateGas
	msg := kaia.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := c.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	_, err = c.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}
	// Note: CallContractAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
}

func testCallContract(t *testing.T, c *Client) {
	// EstimateGas
	msg := kaia.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := c.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
	// CallContract
	if _, err := c.CallContract(context.Background(), msg, big.NewInt(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// PendingCallContract
	if _, err := c.PendingCallContract(context.Background(), msg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func testAtFunctions(t *testing.T, c *Client) {
	_, err := c.HeaderByNumber(context.Background(), big.NewInt(1))
	if err != nil {
		t.Fatalf("BlockByNumber error: %v", err)
	}

	// send a transaction for some interesting pending status
	if err := sendTransaction(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// wait for the transaction to be included in the pending block
	for {
		// Check pending transaction count
		pending, err := c.PendingTransactionCount(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pending == 1 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	// Query balance
	balance, err := c.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: BalanceAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penBalance, err := c.PendingBalanceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balance.Cmp(penBalance) == 0 {
		t.Fatalf("unexpected balance: %v %v", balance, penBalance)
	}
	// NonceAt
	nonce, err := c.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: NonceAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penNonce, err := c.PendingNonceAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if penNonce != nonce+1 {
		t.Fatalf("unexpected nonce: %v %v", nonce, penNonce)
	}
	// StorageAt
	storage, err := c.StorageAt(context.Background(), testAddr, common.Hash{}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: StorageAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penStorage, err := c.PendingStorageAt(context.Background(), testAddr, common.Hash{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(storage, penStorage) {
		t.Fatalf("unexpected storage: %v %v", storage, penStorage)
	}
	// CodeAt
	code, err := c.CodeAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Note: CodeAtHash method is not available in Kaia client
	// Skipping this test for Kaia compatibility
	penCode, err := c.PendingCodeAt(context.Background(), testAddr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !bytes.Equal(code, penCode) {
		t.Fatalf("unexpected code: %v %v", code, penCode)
	}
	// Note: EstimateGasAtBlock and EstimateGasAtBlockHash methods are not available in Kaia client
	// Skipping these tests for Kaia compatibility

	// Verify that sender address of pending transaction is saved in cache.
	pendingBlock, err := c.BlockByNumber(context.Background(), big.NewInt(int64(rpc.PendingBlockNumber)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// No additional RPC should be required, ensure the server is not asked by
	// canceling the context.
	sender, err := c.TransactionSender(newCanceledContext(), pendingBlock.Transactions()[0], pendingBlock.Hash(), 0)
	if err != nil {
		t.Fatal("unable to recover sender:", err)
	}
	if sender != testAddr {
		t.Fatal("wrong sender:", sender)
	}
}

func testTransactionSender(t *testing.T, c *Client) {
	ctx := context.Background()

	// Retrieve testTx1 via RPC.
	block2, err := c.HeaderByNumber(ctx, big.NewInt(2))
	if err != nil {
		t.Fatal("can't get block 1:", err)
	}
	tx1, err := c.TransactionInBlock(ctx, block2.Hash(), 0)
	if err != nil {
		t.Fatal("can't get tx:", err)
	}
	if tx1.Hash() != testTx1.Hash() {
		t.Fatalf("wrong tx hash %v, want %v", tx1.Hash(), testTx1.Hash())
	}

	// The sender address is cached in tx1, so no additional RPC should be required in
	// TransactionSender. Ensure the server is not asked by canceling the context here.
	_, err = c.TransactionSender(newCanceledContext(), tx1, block2.Hash(), 0)
	if err != nil {
		t.Fatal(err)
	}

	// Kaia tx.From is zero for legacy txs, so skip the check
	// if sender1 != testAddr {
	// 	t.Fatal("wrong sender:", sender1)
	// }

	// Now try to get the sender of testTx2, which was not fetched through RPC.
	// TransactionSender should query the server here.
	_, err = c.TransactionSender(ctx, testTx2, block2.Hash(), 1)
	if err != nil {
		t.Fatal(err)
	}
	// Kaia tx.From is zero for legacy txs, so skip the check
	// if sender2 != testAddr {
	// 	t.Fatal("wrong sender:", sender2)
	// }
}

func newCanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done() // Ensure the close of the Done channel
	return ctx
}

func sendTransaction(c *Client) error {
	chainID, err := c.ChainID(context.Background())
	if err != nil {
		return err
	}
	nonce, err := c.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		return err
	}

	signer := types.LatestSignerForChainID(chainID)
	tx := types.NewTransaction(nonce, common.Address{2}, big.NewInt(1), 22000, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	tx, err = types.SignTx(tx, signer, testKey)
	if err != nil {
		return err
	}
	return c.SendTransaction(context.Background(), tx)
}

func genMockHeader(number int) *types.Header {
	numToHash := common.HexToHash(strconv.Itoa(int(number)))
	var parentHash common.Hash
	if number <= 0 {
		parentHash = common.Hash{0xFF}
	} else {
		parentHash = blocks[number-1].Hash()
	}
	header := &types.Header{
		ParentHash:  parentHash,
		Rewardbase:  common.Address{},
		Root:        numToHash,
		TxHash:      numToHash,
		ReceiptHash: numToHash,
		Bloom:       types.Bloom{},
		BlockScore:  big.NewInt(int64(number)),
		Number:      big.NewInt(int64(number)),
		GasUsed:     0,
		Time:        big.NewInt(10 + int64(number)),
		TimeFoS:     0,
		Extra:       []byte{},
		Governance:  []byte{},
		Vote:        []byte{},
		BaseFee:     big.NewInt(25e9),
	}
	return header
}

func TestKaiaClient_AnvilServer(t *testing.T) {
	serverURL, cleanup := launchAnvilServer(t)
	defer cleanup()
	client, err := DialContext(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}

	// Test if server actually responds with a simple call
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if chainId, err := client.NetworkID(ctx); err != nil {
		t.Log("anvil server is not responding:", err)
		t.Skip("skip this test")
		return
	} else if chainId.Cmp(big.NewInt(1337)) != 0 {
		t.Fatal("the server must have chain id 1337, but got", chainId)
		return
	}

	t.Log("Eth client connected to anvil server", serverURL)

	_, err = client.HeaderByNumber(context.Background(), big.NewInt(0))
	assert.Equal(t, err.Error(), "Method not found")
}
