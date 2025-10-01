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
	"fmt"
	"math/big"
	"math/rand"
	"os/exec"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func launchAnvilServer(t *testing.T) (string, func()) {
	randPort := rand.Intn(30000) + 20000
	serverURL := fmt.Sprintf("http://127.0.0.1:%d", randPort)

	// Check if anvil is installed
	_, err := exec.LookPath("anvil")
	if err != nil {
		t.Skip("anvil not found in PATH, skipping test")
	}

	// Start anvil in background
	cmd := exec.Command("anvil", "--chain-id", "1337", "--host", "127.0.0.1", "--port", strconv.Itoa(randPort))
	err = cmd.Start()
	if err != nil {
		t.Skipf("Failed to start anvil: %v, skipping test", err)
	}

	t.Logf("Started anvil server with PID: %d", cmd.Process.Pid)
	var once sync.Once

	// Create cleanup function
	cleanup := func() {
		once.Do(func() {
			if cmd.Process != nil {
				t.Logf("Killing anvil server (PID: %d)", cmd.Process.Pid)
				cmd.Process.Kill()
				cmd.Wait() // Wait for process to actually terminate
			}
		})
	}

	// Kill anvil after 30 seconds as safety timeout
	go func() {
		time.Sleep(30 * time.Second)
		cleanup()
	}()

	// Give anvil time to start up
	time.Sleep(2 * time.Second)

	return serverURL, cleanup
}

func TestEthClient_MockServer(t *testing.T) {
	quitChan := make(chan struct{})
	defer close(quitChan)

	serverURL := launchMockServer(t, quitChan)

	client, err := DialContext(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Client connected to mock server")
	defer client.Close()

	kaiaHeader, err := client.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}

	ethclient, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Eth client connected to mock server")
	ethHeader, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0x3b624db9bc6547b908e2e78460d2849047b6d28c0c078f09d6a0472ab0e57d0c", kaiaHeader.Hash().Hex())
	assert.Equal(t, "0x1c6ef781e4f30626053500c374498f78e3138128603e6f9c92bff0292613c5bb", ethHeader.Hash().Hex())
}

func TestEthClient_AnvilServer(t *testing.T) {
	serverURL, cleanup := launchAnvilServer(t)
	defer cleanup()
	ethclient, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		t.Fatal(err)
	}

	// Test if server actually responds with a simple call
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if chainId, err := ethclient.ChainID(ctx); err != nil {
		t.Log("anvil server is not responding:", err)
		t.Skip("skip this test")
		return
	} else if chainId.Cmp(big.NewInt(1337)) != 0 {
		t.Fatal("the server must have chain id 1337, but got", chainId)
		return
	}

	t.Log("Eth client connected to anvil server", serverURL)

	_, err = ethclient.BlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	header, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	cnt, err := ethclient.TransactionCount(context.Background(), header.Hash())
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, uint(0), cnt)

	testAddrInitialBalance := big.NewInt(1e18)
	tx := func() *types.Transaction {
		richKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
		richAddr := crypto.PubkeyToAddress(richKey.PublicKey)
		if err != nil {
			t.Fatal(err)
		}
		nonce, err := ethclient.NonceAt(context.Background(), richAddr, nil)
		if err != nil {
			t.Fatal(err)
		}
		tx := types.NewTransaction(nonce, testAddr, testAddrInitialBalance, params.TxGas, new(big.Int).SetUint64(params.DefaultLowerBoundBaseFee), nil)
		signer := types.LatestSignerForChainID(genesisConfig.ChainID)
		signedTx, _ := types.SignTx(tx, signer, richKey)
		return signedTx
	}()
	assert.Equal(t, header.ParentHash, common.Hash{})
	assert.NotEqual(t, header.Hash(), common.Hash{})

	hash, err := ethclient.SendRawTransaction(context.Background(), tx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	receipt, err := ethclient.TransactionReceiptRpcOutput(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	status, err := strconv.ParseUint(receipt["status"].(string), 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.ReceiptStatusSuccessful, uint(status), "tx %s failed", hash.Hex())

	balance, err := ethclient.BalanceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, balance.Cmp(testAddrInitialBalance) >= 0)

	nonce, err := ethclient.NonceAt(context.Background(), testAddr, nil)
	if err != nil {
		t.Fatal(err)
	}
	dynamicTx := func() *types.Transaction {
		tx := types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			ChainID:      genesisConfig.ChainID,
			AccountNonce: nonce,
			Recipient:    &testAddr,
			Amount:       big.NewInt(10),
			GasLimit:     25000,
			GasFeeCap:    big.NewInt(50e9),
			GasTipCap:    big.NewInt(25e9),
		})
		signer := types.LatestSignerForChainID(genesisConfig.ChainID)
		signedTx, _ := types.SignTx(tx, signer, testKey)
		return signedTx
	}()
	assert.Equal(t, header.ParentHash, common.Hash{})
	assert.NotEqual(t, header.Hash(), common.Hash{})
	_, err = ethclient.SendRawTransaction(context.Background(), dynamicTx)
	if err != nil {
		t.Fatal(err)
	}

	nonce++
	deployTx := func() *types.Transaction {
		// contract Storage { uint256 number = 1337; * @dev Return value @return value of 'number' */ function retrieve() public view returns (uint256){ return number; } }
		bytecode := common.Hex2Bytes("60806040526105395f553480156013575f5ffd5b5060af80601f5f395ff3fe6080604052348015600e575f5ffd5b50600436106026575f3560e01c80632e64cec114602a575b5f5ffd5b60306044565b604051603b91906062565b60405180910390f35b5f5f54905090565b5f819050919050565b605c81604c565b82525050565b5f60208201905060735f8301846055565b9291505056fea2646970667358221220bbed5c2a1719068dca0cf4da53d280029c147463a0b8f8319bc3494906ad27a964736f6c634300081e0033")
		tx := types.NewContractCreation(nonce, big.NewInt(0), 1e6, big.NewInt(25e9), bytecode)
		signer := types.LatestSignerForChainID(genesisConfig.ChainID)
		signedTx, _ := types.SignTx(tx, signer, testKey)
		return signedTx
	}()
	hash, err = ethclient.SendRawTransaction(context.Background(), deployTx)
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	receipt, err = ethclient.TransactionReceiptRpcOutput(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	status, err = strconv.ParseUint(receipt["status"].(string), 0, 64)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, types.ReceiptStatusSuccessful, uint(status), "tx %s failed", hash.Hex())

	contractAddr := crypto.CreateAddress(testAddr, nonce)
	calldata := common.Hex2Bytes("2e64cec1") // retrieve()(uint256)
	ret, err := ethclient.CallContract(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata}, nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(ret))

	accesslist, _, _, err := ethclient.CreateAccessList(context.Background(), kaia.CallMsg{To: &contractAddr, Data: calldata})
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, accesslist.StorageKeys())

	storage, err := ethclient.StorageAt(context.Background(), contractAddr, (*accesslist)[0].StorageKeys[0], nil)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000539", common.Bytes2Hex(storage))
}

func TestEthClient_AnvilServerWithCleanup(t *testing.T) {
	// Launch anvil server with 30s timeout
	serverURL, cleanup := launchAnvilServer(t)
	defer cleanup() // Ensure cleanup happens when test ends

	// Connect to anvil server
	ethclient, err := DialContextEth(context.Background(), serverURL)
	if err != nil {
		t.Skip("Anvil server not available:", err)
		return
	}
	defer ethclient.Close()

	// Test if server actually responds
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result interface{}
	if err := ethclient.c.CallContext(ctx, &result, "net_version"); err != nil {
		t.Skip("Anvil server not responding:", err)
		return
	}

	t.Log("Connected to anvil server, chain ID:", result)

	// Test basic functionality
	ethHeader, err := ethclient.HeaderByNumber(context.Background(), big.NewInt(0))
	if err != nil {
		t.Fatal("Failed to get genesis block:", err)
	}

	t.Log("Genesis block hash:", ethHeader.Hash().Hex())
	assert.Equal(t, uint64(0), ethHeader.Number.Uint64())

	// Test contract deployment
	bytecode := common.Hex2Bytes("608060405234801561001057600080fd5b506101de806100206000396000f3006080604052600436106100615763ffffffff7c01000000000000000000000000000000000000000000000000000000006000350416631a39d8ef81146100805780636353586b146100a757806370a08231146100ca578063fd6b7ef8146100f8575b3360009081526001602052604081208054349081019091558154019055005b34801561008c57600080fd5b5061009561010d565b60408051918252519081900360200190f35b6100c873ffffffffffffffffffffffffffffffffffffffff60043516610113565b005b3480156100d657600080fd5b5061009573ffffffffffffffffffffffffffffffffffffffff60043516610147565b34801561010457600080fd5b506100c8610159565b60005481565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604081208054349081019091558154019055565b60016020526000908152604090205481565b336000908152600160205260408120805490829055908111156101af57604051339082156108fc029083906000818181858888f193505050501561019c576101af565b3360009081526001602052604090208190555b505600a165627a7a72305820627ca46bb09478a015762806cc00c431230501118c7c26c30ac58c4e09e51c4f0029")

	// Use anvil's default funded account
	richKey, err := crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	if err != nil {
		t.Fatal("Failed to parse rich account key:", err)
	}

	deployTx := types.NewContractCreation(0, big.NewInt(0), 1000000, big.NewInt(1e9), bytecode)
	signer := types.LatestSignerForChainID(big.NewInt(1337))
	signedDeployTx, err := types.SignTx(deployTx, signer, richKey)
	if err != nil {
		t.Fatal("Failed to sign deploy tx:", err)
	}

	hash, err := ethclient.SendRawTransaction(context.Background(), signedDeployTx)
	if err != nil {
		t.Fatal("Failed to deploy contract:", err)
	}

	t.Log("Contract deployed with tx hash:", hash.Hex())
}
