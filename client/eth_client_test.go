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

package client

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Interface compliance checks
var (
	_ = kaia.ChainStateReader(&EthClient{})
	_ = kaia.ContractCaller(&EthClient{})
	_ = kaia.LogFilterer(&EthClient{})
	_ = kaia.TransactionSender(&EthClient{})
	_ = kaia.GasPricer(&EthClient{})
	_ = kaia.PendingStateReader(&EthClient{})
	_ = kaia.PendingContractCaller(&EthClient{})
	_ = kaia.GasEstimator(&EthClient{})
	_ = bind.ContractBackend(&EthClient{})
	_ = bind.DeployBackend(&EthClient{})
)

// AnvilTestSuite tests EthClient against an Anvil server
type AnvilTestSuite struct {
	suite.Suite

	// Server
	serverURL string
	cmd       *exec.Cmd
	cancel    context.CancelFunc

	// Client
	ethClient *EthClient

	chainConfig *params.ChainConfig
	testerKey   *ecdsa.PrivateKey
	testerAddr  common.Address
	testerNonce uint64
}

func TestAnvilTestSuite(t *testing.T) {
	suite.Run(t, new(AnvilTestSuite))
}

// SetupSuite starts anvil server and prepares test fixtures
func (s *AnvilTestSuite) SetupSuite() {
	s.chainConfig = genesisChainConfig.Copy()
	s.initTestAccount()

	if err := s.startServer(); err != nil {
		s.T().Skip("Anvil server setup failed:", err)
		return
	}

	if err := s.connectClient(); err != nil {
		s.TearDownSuite()
		s.T().Skip("Client connection failed:", err)
		return
	}

	s.T().Logf("Anvil server started on %s", s.serverURL)
}

func (s *AnvilTestSuite) initTestAccount() {
	s.testerKey, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	s.testerAddr = crypto.PubkeyToAddress(s.testerKey.PublicKey)
}

// TearDownSuite cleans up resources
func (s *AnvilTestSuite) TearDownSuite() {
	if s.ethClient != nil {
		s.ethClient.Close()
	}
	if s.cancel != nil {
		s.cancel()
	}
}

// startServer launches anvil on a random port
func (s *AnvilTestSuite) startServer() error {
	if _, err := exec.LookPath("anvil"); err != nil {
		return errors.New("anvil not found in PATH")
	}

	port := 20000 + rand.Intn(30000)
	s.serverURL = fmt.Sprintf("http://127.0.0.1:%d", port)

	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.cmd = exec.CommandContext(ctx, "anvil",
		"--chain-id", "1337",
		"--host", "127.0.0.1",
		"--port", strconv.Itoa(port),
	)

	if err := s.cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start anvil: %w", err)
	}

	if err := s.waitForServer(fmt.Sprintf("127.0.0.1:%d", port), 5*time.Second); err != nil {
		cancel()
		return fmt.Errorf("anvil server not ready: %w", err)
	}

	return nil
}

// waitForServer polls until the server is accepting connections
func (s *AnvilTestSuite) waitForServer(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return errors.New("timeout waiting for server")
}

// connectClient establishes connection to the anvil server
func (s *AnvilTestSuite) connectClient() error {
	client, err := tryConnectEth(s.serverURL)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	s.ethClient = client
	return nil
}

// =============================================================================
// Tests
// =============================================================================

func (s *AnvilTestSuite) TestBlockchainAccess() {
	ctx := context.Background()
	ec := s.ethClient
	t := s.T()

	tx := s.sendKaiaToTesterLegacyTx()

	bn, err := ec.BlockNumber(ctx)
	require.NoError(t, err)

	block, err := ec.BlockByNumber(ctx, bn)
	require.NoError(t, err)
	assert.Equal(t, bn.Uint64(), block.Header().Number.Uint64())

	block, err = ec.BlockByHash(ctx, block.Header().Hash())
	require.NoError(t, err)
	assert.Equal(t, bn.Uint64(), block.Header().Number.Uint64())

	header, err := ec.HeaderByNumber(ctx, bn)
	require.NoError(t, err)
	assert.Equal(t, bn.Uint64(), header.Number.Uint64())

	apiTx, _, err := ec.TransactionByHash(ctx, tx.Hash())
	require.NoError(t, err)
	assert.Equal(t, tx.Hash(), apiTx.Hash)

	cnt, err := ec.TransactionCount(ctx, header.Hash())
	require.NoError(t, err)
	assert.Equal(t, uint(1), cnt)

	apiTx, err = ec.TransactionInBlock(ctx, header.Hash(), 0)
	require.NoError(t, err)
	assert.Equal(t, tx.Hash(), apiTx.Hash)

	receipt, err := ec.TransactionReceipt(ctx, tx.Hash())
	require.NoError(t, err)
	assert.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	receiptMap, err := ec.TransactionReceiptRpcOutput(ctx, tx.Hash())
	require.NoError(t, err)
	status, err := strconv.ParseUint(receiptMap["status"].(string), 0, 64)
	require.NoError(t, err)
	assert.Equal(t, types.ReceiptStatusSuccessful, uint(status))
}

func (s *AnvilTestSuite) TestStateAccess() {
	ctx := context.Background()
	ec := s.ethClient
	t := s.T()

	networkID, err := ec.NetworkID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1337), networkID.Uint64())

	chainID, err := ec.ChainID(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1337), chainID.Uint64())

	balance, err := ec.BalanceAt(ctx, s.testerAddr, nil)
	require.NoError(t, err)
	assert.True(t, balance.Cmp(common.Big0) >= 0)

	nonce, err := ec.NonceAt(ctx, s.testerAddr, nil)
	require.NoError(t, err)
	assert.Equal(t, s.testerNonce, nonce)

	header, err := ec.HeaderByNumber(ctx, nil)
	require.NoError(t, err)
	assert.NotEqual(t, header.Hash(), common.Hash{})

	s.testDynamicFeeTx()

	contractAddr, err := s.deployStorageContract()
	require.NoError(t, err)

	calldata := common.Hex2Bytes("2e64cec1") // retrieve()(uint256)
	ret, err := ec.CallContract(ctx, kaia.CallMsg{To: &contractAddr, Data: calldata}, nil)
	require.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(ret))

	gas, err := ec.EstimateGas(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata})
	require.NoError(s.T(), err)
	assert.Equal(s.T(), uint64(23473), gas)

	accesslist, _, _, err := ec.CreateAccessList(ctx, kaia.CallMsg{To: &contractAddr, Data: calldata})
	require.NoError(t, err)
	assert.Equal(t, 1, accesslist.StorageKeys())

	code, err := ec.CodeAt(ctx, contractAddr, nil)
	require.NoError(t, err)
	assert.Equal(t, 175, len(code))

	storage, err := ec.StorageAt(ctx, contractAddr, (*accesslist)[0].StorageKeys[0], nil)
	require.NoError(t, err)
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(storage))

	bn, _ := ec.BlockNumber(ctx)
	logs, err := ec.FilterLogs(ctx, kaia.FilterQuery{
		FromBlock: big.NewInt(0),
		ToBlock:   bn,
		Addresses: []common.Address{contractAddr},
	})
	require.NoError(t, err)
	assert.Equal(t, 0, len(logs))
}

func (s *AnvilTestSuite) TestKaiaClient() {
	t := s.T()
	kaiaClient, err := tryConnect(s.serverURL)
	if err != nil {
		t.Skip("Could not connect Kaia client to anvil server:", err)
		return
	}
	defer kaiaClient.Close()

	_, err = kaiaClient.HeaderByNumber(context.Background(), big.NewInt(0))
	assert.Equal(t, err.Error(), "Method not found")
}

// =============================================================================
// Test Helpers
// =============================================================================

func (s *AnvilTestSuite) sendKaiaToTesterLegacyTx() *types.Transaction {
	unsignedTx := types.NewTransaction(
		s.testerNonce, s.testerAddr, big.NewInt(1e18),
		params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil,
	)
	signer := types.LatestSignerForChainID(s.chainConfig.ChainID)
	signedTx, _ := types.SignTx(unsignedTx, signer, s.testerKey)
	_, err := s.ethClient.SendRawTransaction(context.Background(), signedTx)
	require.NoError(s.T(), err)
	time.Sleep(1 * time.Second)
	s.testerNonce++
	return signedTx
}

func (s *AnvilTestSuite) testDynamicFeeTx() *types.Transaction {
	tx := types.NewTx(&types.TxInternalDataEthereumDynamicFee{
		ChainID:      s.chainConfig.ChainID,
		AccountNonce: s.testerNonce,
		Recipient:    &s.testerAddr,
		Amount:       big.NewInt(10),
		GasLimit:     params.TxGas,
		GasFeeCap:    big.NewInt(50e9),
		GasTipCap:    big.NewInt(25e9),
	})

	signer := types.LatestSignerForChainID(s.chainConfig.ChainID)
	signedTx, _ := types.SignTx(tx, signer, s.testerKey)
	_, err := s.ethClient.SendRawTransaction(context.Background(), signedTx)
	require.NoError(s.T(), err)
	time.Sleep(1 * time.Second)

	s.testerNonce++
	return signedTx
}

func (s *AnvilTestSuite) deployStorageContract() (common.Address, error) {
	ctx := context.Background()

	// Storage contract: uint256 number = 1337; function retrieve() returns (uint256)
	bytecode := common.Hex2Bytes("60806040526105395f553480156013575f5ffd5b5060af80601f5f395ff3fe6080604052348015600e575f5ffd5b50600436106026575f3560e01c80632e64cec114602a575b5f5ffd5b60306044565b604051603b91906062565b60405180910390f35b5f5f54905090565b5f819050919050565b605c81604c565b82525050565b5f60208201905060735f8301846055565b9291505056fea2646970667358221220bbed5c2a1719068dca0cf4da53d280029c147463a0b8f8319bc3494906ad27a964736f6c634300081e0033")

	tx := types.NewContractCreation(s.testerNonce, big.NewInt(0), 1e6, big.NewInt(25e9), bytecode)

	signer := types.LatestSignerForChainID(s.chainConfig.ChainID)
	signedTx, _ := types.SignTx(tx, signer, s.testerKey)
	txhash, err := s.ethClient.SendRawTransaction(ctx, signedTx)
	require.NoError(s.T(), err)
	time.Sleep(1 * time.Second)

	receipt, err := s.ethClient.TransactionReceipt(ctx, txhash)
	if err != nil {
		return common.Address{}, err
	}
	assert.Equal(s.T(), types.ReceiptStatusSuccessful, receipt.Status)
	require.NotEqual(s.T(), common.Address{}, receipt.ContractAddress)

	s.testerNonce++
	return receipt.ContractAddress, nil
}
