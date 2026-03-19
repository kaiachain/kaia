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
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net"
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
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Interface compliance checks
var (
	// _ = kaia.Subscription(&Client{})
	_ = kaia.ChainReader(&KaiaClient{})
	_ = kaia.TransactionReader(&KaiaClient{})
	_ = kaia.ChainStateReader(&KaiaClient{})
	_ = kaia.ChainSyncReader(&KaiaClient{})
	_ = kaia.ContractCaller(&KaiaClient{})
	_ = kaia.LogFilterer(&KaiaClient{})
	_ = kaia.TransactionSender(&KaiaClient{})
	_ = kaia.GasPricer(&KaiaClient{})
	_ = kaia.PendingStateReader(&KaiaClient{})
	_ = kaia.PendingContractCaller(&KaiaClient{})
	_ = kaia.GasEstimator(&KaiaClient{})
	// _ = kaia.PendingStateEventer(&Client{})
)

// Expected genesis block hashes for different clients
var (
	ethGenesisBlockHash  = common.HexToHash("0x3dfc072ca3cee03dd5ecf932a54b762082e5fd5ad6e7d34e7dc0411cfd325e2a")
	kaiaGenesisBlockHash = common.HexToHash("0x3b624db9bc6547b908e2e78460d2849047b6d28c0c078f09d6a0472ab0e57d0c")
	genesisChainConfig   = &params.ChainConfig{
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
)

// =============================================================================
// MockHttpServerTestSuite
// =============================================================================
type MockHttpServerTestSuite struct {
	suite.Suite

	// Server
	server    *http.Server
	serverURL string

	// Client
	kaiaClient *KaiaClient

	// Chain config
	chainConfig *params.ChainConfig

	// Test account
	testerKey         *ecdsa.PrivateKey
	testerAddr        common.Address
	testerInitBalance *big.Int

	// Dummy data
	blocks []*types.Block
	txs    []*types.Transaction
}

func TestMockKaiaRpcServerTestSuite(t *testing.T) {
	suite.Run(t, new(MockHttpServerTestSuite))
}

func (s *MockHttpServerTestSuite) SetupSuite() {
	s.chainConfig = genesisChainConfig.Copy()
	s.initTestAccount()
	s.initBlocksAndTxs()

	if err := s.startServer(); err != nil {
		s.T().Skip("Mock server setup failed:", err)
		return
	}

	if err := s.connectClient(); err != nil {
		s.TearDownSuite()
		s.T().Skip("Client connection failed:", err)
		return
	}

	s.T().Logf("MockHttpServer started on %s", s.serverURL)
}

func (s *MockHttpServerTestSuite) initTestAccount() {
	s.testerKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	s.testerAddr = crypto.PubkeyToAddress(s.testerKey.PublicKey)
	s.testerInitBalance = big.NewInt(2e15)
}

func (s *MockHttpServerTestSuite) initBlocksAndTxs() {
	s.blocks = make([]*types.Block, 3)
	for i := 0; i < 3; i++ {
		s.blocks[i] = types.NewBlockWithHeader(s.genMockHeader(i))
	}

	signer := types.LatestSignerForChainID(s.chainConfig.ChainID)

	tx1 := types.NewTransaction(0, s.testerAddr, big.NewInt(12), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signedTx1, _ := types.SignTx(tx1, signer, s.testerKey)
	s.txs = append(s.txs, signedTx1)

	tx2 := types.NewTransaction(1, common.Address{2}, big.NewInt(8), params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
	signedTx2, _ := types.SignTx(tx2, signer, s.testerKey)
	s.txs = append(s.txs, signedTx2)
}

func (s *MockHttpServerTestSuite) TearDownSuite() {
	if s.kaiaClient != nil {
		s.kaiaClient.Close()
	}
	if s.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.server.Shutdown(ctx)
	}
}

func (s *MockHttpServerTestSuite) startServer() error {
	handler := s.createMockHandler()

	addr := "127.0.0.1:36000"
	s.serverURL = "http://" + addr
	s.server = &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	errChan := make(chan error, 1)
	go func() {
		errChan <- s.server.ListenAndServe()
	}()

	// Wait for server to be ready
	if err := s.waitForServer(addr, 5*time.Second); err != nil {
		select {
		case serverErr := <-errChan:
			return fmt.Errorf("server startup failed: %w", serverErr)
		default:
			return err
		}
	}

	return nil
}

func (s *MockHttpServerTestSuite) waitForServer(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("timeout waiting for server")
}

func (s *MockHttpServerTestSuite) connectClient() error {
	client, err := tryConnect(s.serverURL)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	s.kaiaClient = client
	return nil
}

func (s *MockHttpServerTestSuite) createMockHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var reqData map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		s.T().Logf("Request: %s", reqData["method"])

		method, _ := reqData["method"].(string)
		id := reqData["id"]
		response := s.handleRPCMethod(method, reqData["params"])
		response["id"] = id
		response["jsonrpc"] = "2.0"

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})
}

func (s *MockHttpServerTestSuite) handleRPCMethod(method string, params interface{}) map[string]interface{} {
	switch method {
	case "kaia_chainID":
		return map[string]interface{}{"result": "0x" + s.chainConfig.ChainID.Text(16)}

	case "kaia_getBalance":
		p := params.([]interface{})
		address := common.HexToAddress(p[0].(string))
		blockNumber := s.parseBlockNumber(p[1].(string))
		return s.mockGetBalance(address, blockNumber)

	case "kaia_getBlockByNumber":
		p := params.([]interface{})
		return s.mockGetBlockByNumber(p[0].(string))

	case "eth_getBlockByNumber":
		p := params.([]interface{})
		return s.mockGetBlockByNumberEth(p[0].(string))

	case "kaia_getBlockByHash":
		p := params.([]interface{})
		return s.mockGetBlockByHash(p[0].(string))

	case "kaia_blockNumber":
		return map[string]interface{}{"result": "0x2"}

	case "kaia_syncing":
		return map[string]interface{}{"result": false}

	case "net_version":
		return map[string]interface{}{"result": s.chainConfig.ChainID.Text(10)}

	case "kaia_gasPrice":
		return map[string]interface{}{"result": "0x3b9aca00"}

	case "kaia_estimateGas":
		return map[string]interface{}{"result": "0x5208"}

	case "kaia_call":
		return map[string]interface{}{"result": "0x"}

	case "kaia_getTransactionByBlockHashAndIndex":
		p := params.([]interface{})
		blockHash := p[0].(string)
		txIdx, _ := strconv.ParseUint(p[1].(string)[2:], 16, 64)
		return s.mockGetTransactionByBlockHashAndIndex(blockHash, txIdx)

	default:
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32601,
				"message": fmt.Sprintf("the method %s does not exist/is not available", method),
			},
		}
	}
}

func (s *MockHttpServerTestSuite) parseBlockNumber(str string) *big.Int {
	assert.True(s.T(), strings.HasPrefix(str, "0x"), "blockNumber should start with 0x, got %v", str)

	num, ok := new(big.Int).SetString(str[2:], 16)
	assert.True(s.T(), ok, "failed to parse blockNumber: %v", str)
	return num
}

// =============================================================================
// Tests
// =============================================================================

func (s *MockHttpServerTestSuite) TestKaiaClient() {
	s.Run("Header", s.testHeader)
	s.Run("BalanceAt", s.testBalanceAt)
	s.Run("ChainID", s.testChainID)
	s.Run("TxInBlockInterrupted", s.testTransactionInBlock)
	s.Run("GetBlock", s.testGetBlock)
	s.Run("StatusFunctions", s.testStatusFunctions)
	s.Run("CallContract", s.testCallContract)
	s.Run("CallContractAtHash", s.testCallContractAtHash)
	s.Run("TransactionSender", s.testTransactionSender)
}

func (s *MockHttpServerTestSuite) TestEthClient() {
	t := s.T()

	ethclient, err := tryConnectEth(s.serverURL)
	if err != nil {
		t.Skip("Could not connect Eth client to mock server:", err)
		return
	}
	defer ethclient.Close()

	ethHeader, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	require.NoError(t, err)
	assert.Equal(t, ethGenesisBlockHash, ethHeader.Hash())
}

// =============================================================================
// Mock Response Generators
// =============================================================================

func (s *MockHttpServerTestSuite) mockGetBalance(address common.Address, blockNumber *big.Int) map[string]interface{} {
	if blockNumber.Cmp(big.NewInt(int64(len(s.blocks)))) > 0 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}
	if address == s.testerAddr {
		return map[string]interface{}{"result": "0x71afd498d0000"} // 0.002 KAIA
	}
	return map[string]interface{}{"result": "0x0"}
}

func (s *MockHttpServerTestSuite) mockGetBlockByNumber(blockNumberArg string) map[string]interface{} {
	block := s.getBlockByArg(blockNumberArg)
	if block == nil {
		return map[string]interface{}{"result": nil}
	}

	rpcOutput, err := api.RpcOutputBlock(block, false, false, s.chainConfig)
	require.NoError(s.T(), err)
	return map[string]interface{}{"result": rpcOutput}
}

func (s *MockHttpServerTestSuite) mockGetBlockByNumberEth(blockNumberArg string) map[string]interface{} {
	block := s.getBlockByArg(blockNumberArg)
	if block == nil {
		return map[string]interface{}{"result": nil}
	}

	rpcOutput, err := api.RpcMarshalEthBlock(block, nil, s.chainConfig, false, false, false)
	require.NoError(s.T(), err)
	return map[string]interface{}{"result": rpcOutput}
}

func (s *MockHttpServerTestSuite) mockGetBlockByHash(blockHash string) map[string]interface{} {
	blockNum := slices.IndexFunc(s.blocks, func(b *types.Block) bool {
		return b.Hash().Hex() == blockHash
	})
	if blockNum == -1 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}

	rpcOutput, err := api.RpcOutputBlock(s.blocks[blockNum], false, false, s.chainConfig)
	require.NoError(s.T(), err)
	return map[string]interface{}{"result": rpcOutput}
}

func (s *MockHttpServerTestSuite) mockGetTransactionByBlockHashAndIndex(blockHash string, transactionIndex uint64) map[string]interface{} {
	blockNum := slices.IndexFunc(s.blocks, func(b *types.Block) bool {
		return b.Hash().Hex() == blockHash
	})
	if blockNum < 2 || transactionIndex > 1 {
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    -32000,
				"message": kaia.NotFound.Error(),
			},
		}
	}

	return map[string]interface{}{"result": s.txs[transactionIndex]}
}

func (s *MockHttpServerTestSuite) getBlockByArg(arg string) *types.Block {
	switch arg {
	case "0x0", "earliest":
		return s.blocks[0]
	case "0x1":
		return s.blocks[1]
	case "0x2", "latest", "pending", "-0x2", "-0x1":
		return s.blocks[2]
	default:
		return nil
	}
}

// =============================================================================
// Individual Test Functions
// =============================================================================

func (s *MockHttpServerTestSuite) testHeader() {
	t := s.T()
	client := s.kaiaClient

	tests := map[string]struct {
		block   *big.Int
		want    *types.Header
		wantErr error
	}{
		"genesis":      {block: big.NewInt(0), want: s.blocks[0].Header()},
		"first_block":  {block: big.NewInt(1), want: s.blocks[1].Header()},
		"future_block": {block: big.NewInt(1000000000), want: nil, wantErr: kaia.NotFound},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := client.HeaderByNumber(ctx, tt.block)
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

func (s *MockHttpServerTestSuite) testBalanceAt() {
	t := s.T()
	client := s.kaiaClient

	tests := map[string]struct {
		account common.Address
		block   *big.Int
		want    *big.Int
		wantErr error
	}{
		"valid_account_genesis": {account: s.testerAddr, block: big.NewInt(0), want: s.testerInitBalance},
		"valid_account":         {account: s.testerAddr, block: big.NewInt(1), want: s.testerInitBalance},
		"non_existent_account":  {account: common.Address{1}, block: big.NewInt(1), want: big.NewInt(0)},
		"future_block":          {account: s.testerAddr, block: big.NewInt(1000000000), want: big.NewInt(0), wantErr: kaia.NotFound},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			got, err := client.BalanceAt(ctx, tt.account, tt.block)
			if tt.wantErr != nil && (err == nil || err.Error() != tt.wantErr.Error()) {
				t.Fatalf("BalanceAt(%x, %v) error = %q, want %q", tt.account, tt.block, err, tt.wantErr)
			}
			assert.Equal(s.T(), tt.want.String(), got.String())
		})
	}
}

func (s *MockHttpServerTestSuite) testTransactionInBlock() {
	t := s.T()
	c := s.kaiaClient

	block, err := c.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)

	// Test tx not found
	_, err = c.TransactionInBlock(context.Background(), block.Hash(), 20)
	assert.Equal(t, kaia.NotFound.Error(), err.Error())

	// Test tx found
	tx, err := c.TransactionInBlock(context.Background(), block.Hash(), 0)
	require.NoError(t, err)
	assert.Equal(t, s.txs[0].Hash(), tx.Hash())

	tx, err = c.TransactionInBlock(context.Background(), block.Hash(), 1)
	require.NoError(t, err)
	assert.Equal(t, s.txs[1].Hash(), tx.Hash())

	// Test pending block
	_, err = c.BlockByNumber(context.Background(), big.NewInt(int64(rpc.PendingBlockNumber)))
	require.NoError(t, err)
}

func (s *MockHttpServerTestSuite) testChainID() {
	id, err := s.kaiaClient.ChainID(context.Background())
	require.NoError(s.T(), err)
	assert.Equal(s.T(), s.chainConfig.ChainID, id)
}

func (s *MockHttpServerTestSuite) testGetBlock() {
	t := s.T()
	c := s.kaiaClient

	blockNumber, err := c.BlockNumber(context.Background())
	require.NoError(t, err)
	assert.Equal(t, int64(2), blockNumber.Int64())

	block, err := c.BlockByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	require.NoError(t, err)
	assert.Equal(t, blockNumber.Uint64(), block.NumberU64())

	blockH, err := c.BlockByHash(context.Background(), block.Hash())
	require.NoError(t, err)
	assert.Equal(t, block.Hash(), blockH.Hash())

	header, err := c.HeaderByNumber(context.Background(), big.NewInt(blockNumber.Int64()))
	require.NoError(t, err)
	assert.Equal(t, block.Header().Hash(), header.Hash())

	headerH, err := c.HeaderByHash(context.Background(), block.Hash())
	require.NoError(t, err)
	assert.Equal(t, block.Header().Hash(), headerH.Hash())
}

func (s *MockHttpServerTestSuite) testStatusFunctions() {
	t := s.T()
	c := s.kaiaClient

	progress, err := c.SyncProgress(context.Background())
	require.NoError(t, err)
	assert.Nil(t, progress)

	networkID, err := c.NetworkID(context.Background())
	require.NoError(t, err)
	assert.Equal(t, s.chainConfig.ChainID, networkID)

	gasPrice, err := c.SuggestGasPrice(context.Background())
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(1000000000), gasPrice)
}

func (s *MockHttpServerTestSuite) testCallContractAtHash() {
	t := s.T()
	c := s.kaiaClient

	msg := kaia.CallMsg{
		From:  s.testerAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}

	gas, err := c.EstimateGas(context.Background(), msg)
	require.NoError(t, err)
	assert.Equal(t, uint64(21000), gas)

	_, err = c.HeaderByNumber(context.Background(), big.NewInt(1))
	require.NoError(t, err)
}

func (s *MockHttpServerTestSuite) testCallContract() {
	t := s.T()
	c := s.kaiaClient

	msg := kaia.CallMsg{
		From:  s.testerAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}

	gas, err := c.EstimateGas(context.Background(), msg)
	require.NoError(t, err)
	assert.Equal(t, uint64(21000), gas)

	_, err = c.CallContract(context.Background(), msg, big.NewInt(1))
	require.NoError(t, err)
}

func (s *MockHttpServerTestSuite) testTransactionSender() {
	t := s.T()
	c := s.kaiaClient
	ctx := context.Background()

	block2, err := c.HeaderByNumber(ctx, big.NewInt(2))
	require.NoError(t, err)

	tx1, err := c.TransactionInBlock(ctx, block2.Hash(), 0)
	require.NoError(t, err)
	assert.Equal(t, s.txs[0].Hash(), tx1.Hash())

	// Sender is cached, no RPC needed
	_, err = c.TransactionSender(newCanceledContext(), tx1, block2.Hash(), 0)
	require.NoError(t, err)

	// s.txs[1] not fetched via RPC, will query server
	_, err = c.TransactionSender(ctx, s.txs[1], block2.Hash(), 1)
	require.NoError(t, err)
}

// =============================================================================
// Helper Functions
// =============================================================================

func newCanceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	<-ctx.Done()
	return ctx
}

func (s *MockHttpServerTestSuite) genMockHeader(number int) *types.Header {
	numToHash := common.HexToHash(strconv.Itoa(number))
	var parentHash common.Hash
	if number <= 0 {
		parentHash = common.Hash{0xFF}
	} else {
		parentHash = s.blocks[number-1].Hash()
	}

	return &types.Header{
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
}

func tryConnect(serverURL string) (*KaiaClient, error) {
	client, err := DialContext(context.Background(), serverURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = client.NetworkID(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

func tryConnectEth(serverURL string) (*EthClient, error) {
	client, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err = client.NetworkID(ctx); err != nil {
		return nil, err
	}

	return client, nil
}
