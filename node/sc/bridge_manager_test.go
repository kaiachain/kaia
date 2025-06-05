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

package sc

import (
	"context"
	crand "crypto/rand"
	"encoding/hex"
	"log"
	"math/big"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/v2/accounts"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind"
	"github.com/kaiachain/kaia/v2/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/v2/accounts/keystore"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/blockchain/vm"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/contracts/contracts/service_chain/bridge"
	sctoken "github.com/kaiachain/kaia/v2/contracts/contracts/testing/sc_erc20"
	scnft "github.com/kaiachain/kaia/v2/contracts/contracts/testing/sc_erc721"
	scnft_no_uri "github.com/kaiachain/kaia/v2/contracts/contracts/testing/sc_erc721_no_uri"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/node/sc/bridgepool"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/rlp"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/stretchr/testify/assert"
)

// WaitGroupWithTimeOut waits the given wait group until the timout duration.
func WaitGroupWithTimeOut(wg *sync.WaitGroup, duration time.Duration, t *testing.T) {
	c := make(chan struct{})
	go func() {
		wg.Wait()
		c <- struct{}{}
	}()
	t.Log("start to wait group")
	select {
	case <-c:
		t.Log("waiting group is done")
	case <-time.After(duration):
		t.Fatal("timed out waiting group")
	}
}

// CheckReceipt can check if the tx receipt has expected status.
func CheckReceipt(b bind.DeployBackend, tx *types.Transaction, duration time.Duration, expectedStatus uint, t *testing.T) {
	timeoutContext, cancelTimeout := context.WithTimeout(context.Background(), duration)
	defer cancelTimeout()

	receipt, err := bind.WaitMined(timeoutContext, b, tx)
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedStatus, receipt.Status)
}

func handleValueTransfer(t *testing.T, ev IRequestValueTransferEvent, bridgeInfo *BridgeInfo, wg *sync.WaitGroup, backend *backends.SimulatedBackend) {
	var (
		tokenType      = ev.GetTokenType()
		valueOrTokenId = ev.GetValueOrTokenId()
		from           = ev.GetFrom()
		to             = ev.GetTo()
		contractAddr   = ev.GetRaw().Address
		tokenAddr      = ev.GetTokenAddress()
		requestNonce   = ev.GetRequestNonce()
		txHash         = ev.GetRaw().TxHash
	)
	t.Log("Request Event",
		"type", tokenType,
		"amount", valueOrTokenId,
		"from", from.String(),
		"to", to.String(),
		"contract", contractAddr.String(),
		"token", tokenAddr.String(),
		"requestNonce", requestNonce)

	bridge := bridgeInfo.bridge
	done, err := bridge.HandledRequestTx(nil, txHash)
	assert.NoError(t, err)
	assert.Equal(t, false, done)

	// insert the value transfer request event to the bridge info's event list.
	bridgeInfo.AddRequestValueTransferEvents([]IRequestValueTransferEvent{ev})

	// handle the value transfer request event in the event list.
	bridgeInfo.processingPendingRequestEvents()

	backend.Commit() // block
	wg.Done()
	done, err = bridge.HandledRequestTx(nil, txHash)
	assert.NoError(t, err)
	assert.Equal(t, true, done)
}

// TestBridgeManager tests the event/method of Token/NFT/Bridge contracts.
// TODO-Kaia-Servicechain needs to refine this test.
// - consider parent/child chain simulated backend.
// - separate each test
func TestBridgeManager(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(10)

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	pAuth := bacc.cAccount.GenerateTransactOpts()
	cAuth := bacc.pAccount.GenerateTransactOpts()

	// Generate a new random account and a funded simulator
	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		alice.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}
	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bridgeManager, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	testToken := big.NewInt(123)
	testKAIA := big.NewInt(321)

	// 1. Deploy Bridge Contract
	addr, err := bridgeManager.DeployBridgeTest(sim, 10000, false)
	if err != nil {
		log.Fatalf("Failed to deploy new bridge contract: %v", err)
	}
	bridgeInfo, _ := bridgeManager.GetBridgeInfo(addr)
	bridge := bridgeInfo.bridge
	t.Log("===== BridgeContract Addr ", addr.Hex())
	sim.Commit() // block

	// 2. Deploy Token Contract
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(alice, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployGXToken: %v", err)
	}
	sim.Commit() // block

	// 3. Deploy NFT Contract
	nftAddr, tx, nft, err := scnft.DeployServiceChainNFT(alice, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}
	sim.Commit() // block

	// Register the owner as a signer
	_, err = bridge.RegisterOperator(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, cAuth.From)
	assert.NoError(t, err)
	sim.Commit() // block

	// Register tokens on the bridgeInfo
	bridgeInfo.RegisterToken(tokenAddr, tokenAddr)
	bridgeInfo.RegisterToken(nftAddr, nftAddr)

	// Register tokens on the bridge
	bridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr)
	bridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, nftAddr, nftAddr)
	sim.Commit() // block

	cTokenAddr, err := bridge.RegisteredTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)
	cNftAddr, err := bridge.RegisteredTokens(nil, nftAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cNftAddr, nftAddr)

	balance, _ := sim.BalanceAt(context.Background(), pAuth.From, nil)
	t.Logf("auth(%v) KAIA balance : %v\n", pAuth.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), cAuth.From, nil)
	t.Logf("auth2(%v) KAIA balance : %v\n", cAuth.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), alice.From, nil)
	t.Logf("auth3(%v) KAIA balance : %v\n", alice.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), bob.From, nil)
	t.Logf("auth4(%v) KAIA balance : %v\n", bob.From.String(), balance)

	// 4. Subscribe Bridge Contract
	bridgeManager.SubscribeEvent(addr)

	reqVTevCh := make(chan RequestValueTransferEvent)
	reqVTencodedEvCh := make(chan RequestValueTransferEncodedEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	bridgeManager.SubscribeReqVTev(reqVTevCh)
	bridgeManager.SubscribeReqVTencodedEv(reqVTencodedEvCh)
	bridgeManager.SubscribeHandleVTev(handleValueTransferEventCh)

	go func() {
		for {
			select {
			case ev := <-reqVTevCh:
				handleValueTransfer(t, ev, bridgeInfo, &wg, sim)
			case ev := <-reqVTencodedEvCh:
				handleValueTransfer(t, ev, bridgeInfo, &wg, sim)
			case ev := <-handleValueTransferEventCh:
				t.Log("Handle value transfer event",
					"bridgeAddr", ev.Raw.Address.Hex(),
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"owner", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"handleNonce", ev.HandleNonce)
				wg.Done()
			}
		}
	}()

	nftTokenIDs := []uint64{4437, 4438, 4439}
	testURIs := []string{"", "testURI", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	// 6. Register (Mint) an NFT to Alice
	{
		for i := 0; i < len(nftTokenIDs); i++ {
			tx, err = nft.MintWithTokenURI(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, alice.From, big.NewInt(int64(nftTokenIDs[i])), testURIs[i])
			assert.NoError(t, err)
			t.Log("Register NFT Transaction", tx.Hash().Hex())
			sim.Commit() // block
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenIDs[i])))
			assert.Equal(t, nil, err)
			assert.Equal(t, alice.From, owner)
		}
	}

	// 7. Request ERC20 Transfer from Alice to Bob
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, testToken, bob.From, big.NewInt(0), nil)
		assert.NoError(t, err)
		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8. RequestKLAYTransfer from Alice to Bob
	{
		tx, err = bridge.RequestKLAYTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, Value: testKAIA, GasLimit: testGasLimit}, bob.From, testKAIA, nil)
		assert.NoError(t, err)
		t.Log("DepositKLAY Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9. Request NFT transfer from Alice to Bob
	{
		for i := 0; i < len(nftTokenIDs); i++ {
			tx, err = nft.RequestValueTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, big.NewInt(int64(nftTokenIDs[i])), bob.From, nil)
			assert.NoError(t, err)
			t.Log("nft.RequestValueTransfer Transaction", tx.Hash().Hex())
			sim.Commit() // block
			CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

			uri, err := nft.TokenURI(nil, big.NewInt(int64(nftTokenIDs[i])))
			assert.NoError(t, err)
			assert.Equal(t, testURIs[i], uri)
			t.Log("URI length: ", len(testURIs[i]), len(uri))
		}
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(nil, bob.From)
		assert.Equal(t, nil, err)
		assert.Equal(t, testToken.String(), balance.String())
	}

	// 11. Check KAIA balance
	{
		balance, err = sim.BalanceAt(context.Background(), bob.From, nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, testKAIA.String(), balance.String())
	}

	// 12. Check NFT owner sent by RequestValueTransfer()
	{
		for i := 0; i < len(nftTokenIDs); i++ {
			owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenIDs[i])))
			assert.Equal(t, nil, err)
			assert.Equal(t, bob.From, owner)
		}
	}
	bridgeManager.Stop()
}

// TestBridgeManagerERC721_notSupportURI tests if bridge can handle an ERC721 which does not support URI.
func TestBridgeManagerERC721_notSupportURI(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	// pAuth := bacc.cAccount.GenerateTransactOpts()
	cAuth := bacc.pAccount.GenerateTransactOpts()

	// Generate a new random account and a funded simulator
	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		alice.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bridgeManager, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	// Deploy Bridge Contract
	addr, err := bridgeManager.DeployBridgeTest(sim, 10000, false)
	if err != nil {
		log.Fatalf("Failed to deploy new bridge contract: %v", err)
	}
	bridgeInfo, _ := bridgeManager.GetBridgeInfo(addr)
	bridge := bridgeInfo.bridge
	t.Log("===== BridgeContract Addr ", addr.Hex())
	sim.Commit() // block

	// Deploy NFT Contract
	nftTokenID := uint64(4438)
	nftAddr, tx, nft, err := scnft_no_uri.DeployServiceChainNFTNoURI(alice, sim, addr)
	if err != nil {
		log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
	}

	nft_uri, err := scnft.NewServiceChainNFT(nftAddr, sim)
	if err != nil {
		log.Fatalf("Failed to get NFT object: %v", err)
	}

	sim.Commit() // block

	// Register the owner as a signer
	_, err = bridge.RegisterOperator(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, cAuth.From)
	assert.NoError(t, err)
	sim.Commit() // block

	// Register tokens on the bridgeInfo
	bridgeInfo.RegisterToken(nftAddr, nftAddr)

	// Register tokens on the bridge
	bridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, nftAddr, nftAddr)
	sim.Commit() // block

	cNftAddr, err := bridge.RegisteredTokens(nil, nftAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cNftAddr, nftAddr)

	// Subscribe Bridge Contract
	bridgeManager.SubscribeEvent(addr)

	reqVTevCh := make(chan RequestValueTransferEvent)
	reqVTencodedEvCh := make(chan RequestValueTransferEncodedEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	bridgeManager.SubscribeReqVTev(reqVTevCh)
	bridgeManager.SubscribeReqVTencodedEv(reqVTencodedEvCh)
	bridgeManager.SubscribeHandleVTev(handleValueTransferEventCh)

	go func() {
		for {
			select {
			case ev := <-reqVTevCh:
				handleValueTransfer(t, ev, bridgeInfo, &wg, sim)
			case ev := <-reqVTencodedEvCh:
				handleValueTransfer(t, ev, bridgeInfo, &wg, sim)
			case ev := <-handleValueTransferEventCh:
				t.Log("Handle value transfer event",
					"bridgeAddr", ev.Raw.Address.Hex(),
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"owner", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"handleNonce", ev.HandleNonce)
				wg.Done()
			}
		}
	}()

	// Register (Mint) an NFT to Alice
	{
		tx, err = nft.Mint(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, alice.From, big.NewInt(int64(nftTokenID)))
		assert.NoError(t, err)
		t.Log("Register NFT Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		assert.Equal(t, alice.From, owner)
	}

	// Request NFT transfer from Alice to Bob
	{
		tx, err = nft.RequestValueTransfer(&bind.TransactOpts{From: alice.From, Signer: alice.Signer, GasLimit: testGasLimit}, big.NewInt(int64(nftTokenID)), bob.From, nil)
		assert.NoError(t, err)
		t.Log("nft.RequestValueTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
		uri, err := nft_uri.TokenURI(nil, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, vm.ErrExecutionReverted, err)
		assert.Equal(t, "", uri)
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// Check NFT owner
	{
		owner, err := nft.OwnerOf(nil, big.NewInt(int64(nftTokenID)))
		assert.Equal(t, nil, err)
		assert.Equal(t, bob.From, owner)
	}

	bridgeManager.Stop()
}

// TestBridgeManagerWithFee tests the KAIA/ERC20 transfer with fee.
func TestBridgeManagerWithFee(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(7 * 2)

	// Generate a new random account and a funded simulator
	AliceKey, _ := crypto.GenerateKey()
	Alice := bind.NewKeyedTransactor(AliceKey)

	BobKey, _ := crypto.GenerateKey()
	Bob := bind.NewKeyedTransactor(BobKey)

	receiverKey, _ := crypto.GenerateKey()
	receiver := bind.NewKeyedTransactor(receiverKey)

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	pAuth := bacc.cAccount.GenerateTransactOpts()
	cAuth := bacc.pAccount.GenerateTransactOpts()

	// Create Simulated backend
	initialValue := int64(10000000000)
	alloc := blockchain.GenesisAlloc{
		Alice.From:            {Balance: big.NewInt(initialValue)},
		bacc.cAccount.address: {Balance: big.NewInt(initialValue)},
		bacc.pAccount.address: {Balance: big.NewInt(initialValue)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
	}
	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bridgeManager, err := NewBridgeManager(sc)

	testToken := int64(100000)
	testKAIA := int64(100000)
	KAIAFee := int64(500)
	ERC20Fee := int64(500)

	// 1. Deploy Bridge Contract
	pBridgeAddr, err := bridgeManager.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)
	pBridgeInfo, _ := bridgeManager.GetBridgeInfo(pBridgeAddr)
	pBridge := pBridgeInfo.bridge
	t.Log("===== BridgeContract Addr ", pBridgeAddr.Hex())
	sim.Commit() // block

	// 2. Deploy Token Contract
	tokenAddr, tx, token, err := sctoken.DeployServiceChainToken(pAuth, sim, pBridgeAddr)
	assert.NoError(t, err)
	sim.Commit() // block

	// Set value transfer fee
	{
		nilReceiver, err := pBridge.FeeReceiver(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, common.Address{}, nilReceiver)
	}

	pBridge.SetFeeReceiver(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, receiver.From)
	sim.Commit() // block

	{
		recv, err := pBridge.FeeReceiver(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, receiver.From, recv)
	}

	{
		fee, err := pBridge.FeeOfKLAY(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, big.NewInt(0).String(), fee.String())
	}

	{
		fee, err := pBridge.FeeOfERC20(nil, tokenAddr)
		assert.Equal(t, nil, err)
		assert.Equal(t, big.NewInt(0).String(), fee.String())
	}

	cn, err := pBridge.ConfigurationNonce(nil)
	assert.NoError(t, err)
	_, err = pBridge.RegisterOperator(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, cAuth.From)
	assert.NoError(t, err)
	pBridge.SetKLAYFee(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, big.NewInt(KAIAFee), cn)
	pBridge.SetERC20Fee(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, big.NewInt(ERC20Fee), cn+1)
	sim.Commit() // block

	{
		fee, err := pBridge.FeeOfKLAY(nil)
		assert.Equal(t, nil, err)
		assert.Equal(t, KAIAFee, fee.Int64())
	}

	{
		fee, err := pBridge.FeeOfERC20(nil, tokenAddr)
		assert.Equal(t, nil, err)
		assert.Equal(t, ERC20Fee, fee.Int64())
	}

	// Register tokens on the bridgeInfo
	pBridgeInfo.RegisterToken(tokenAddr, tokenAddr)

	// Register tokens on the bridge
	pBridge.RegisterToken(&bind.TransactOpts{From: cAuth.From, Signer: cAuth.Signer, GasLimit: testGasLimit}, tokenAddr, tokenAddr)
	sim.Commit() // block

	cTokenAddr, err := pBridge.RegisteredTokens(nil, tokenAddr)
	assert.Equal(t, err, nil)
	assert.Equal(t, cTokenAddr, tokenAddr)

	balance, _ := sim.BalanceAt(context.Background(), Alice.From, nil)
	t.Logf("Alice(%v) KAIA balance : %v\n", Alice.From.String(), balance)

	balance, _ = sim.BalanceAt(context.Background(), Bob.From, nil)
	t.Logf("Bob(%v) KAIA balance : %v\n", Bob.From.String(), balance)

	// 4. Subscribe Bridge Contract
	bridgeManager.SubscribeEvent(pBridgeAddr)

	reqVTevCh := make(chan RequestValueTransferEvent)
	handleValueTransferEventCh := make(chan *HandleValueTransferEvent)
	bridgeManager.SubscribeReqVTev(reqVTevCh)
	bridgeManager.SubscribeHandleVTev(handleValueTransferEventCh)

	go func() {
		for {
			select {
			case ev := <-reqVTevCh:
				t.Log("Request value transfer event",
					"type", ev.GetTokenType(),
					"amount", ev.GetValueOrTokenId(),
					"from", ev.GetFrom().String(),
					"to", ev.GetTo().String(),
					"contract", ev.GetRaw().Address.String(),
					"token", ev.GetTokenAddress().String(),
					"requestNonce", ev.GetRequestNonce(),
					"fee", ev.GetFee().String())

				// insert the value transfer request event to the bridge info's event list.
				pBridgeInfo.AddRequestValueTransferEvents([]IRequestValueTransferEvent{ev})

				// handle the value transfer request event in the event list.
				pBridgeInfo.processingPendingRequestEvents()

				sim.Commit() // block
				wg.Done()

			case ev := <-handleValueTransferEventCh:
				t.Log("Handle value transfer event",
					"bridgeAddr", ev.Raw.Address.Hex(),
					"type", ev.TokenType,
					"amount", ev.ValueOrTokenId,
					"owner", ev.To.String(),
					"contract", ev.Raw.Address.String(),
					"token", ev.TokenAddress.String(),
					"handleNonce", ev.HandleNonce)
				wg.Done()
			}
		}
	}()

	// 5. transfer from parentAcc to Alice for charging and check balances
	{
		tx, err = token.Transfer(&bind.TransactOpts{From: pAuth.From, Signer: pAuth.Signer, GasLimit: testGasLimit}, Alice.From, big.NewInt(initialValue))
		if err != nil {
			log.Fatalf("Failed to Transfer for charging: %v", err)
		}
		t.Log("Transfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)

		balance, err = token.BalanceOf(nil, pAuth.From)
		assert.Equal(t, nil, err)
		t.Log("parentAcc token balance", balance.String())

		balance, err = token.BalanceOf(nil, Alice.From)
		assert.Equal(t, nil, err)
		t.Log("Alice token balance", balance.String())

		balance, err = token.BalanceOf(nil, Bob.From)
		assert.Equal(t, nil, err)
		t.Log("Bob token balance", balance.String())
	}

	// 7-1. Request ERC20 Transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee), nil)
		assert.NoError(t, err)
		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 7-2. Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(0), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-3. Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee-1), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 7-4. Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.RequestValueTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, big.NewInt(testToken), Bob.From, big.NewInt(ERC20Fee+1), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-1. Approve/Request ERC20 Transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee), nil)
		assert.Equal(t, nil, err)

		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 8-2. Approve/Request ERC20 Transfer from Alice to Bob with insufficient zero feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(0), nil)
		assert.Equal(t, nil, err)

		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-3. Approve/Request ERC20 Transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee-1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee-1), nil)
		assert.Equal(t, nil, err)

		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 8-4. Approve/Request ERC20 Transfer from Alice to Bob with enough feeLimit
	{
		tx, err = token.Approve(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, pBridgeAddr, big.NewInt(testToken+ERC20Fee+1))
		assert.Equal(t, nil, err)

		tx, err = pBridge.RequestERC20Transfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, GasLimit: testGasLimit}, tokenAddr, Bob.From, big.NewInt(testToken), big.NewInt(ERC20Fee+1), nil)
		assert.Equal(t, nil, err)

		t.Log("RequestValueTransfer Transaction", tx.Hash().Hex())
		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-1. Request KAIA transfer from Alice to Bob with same feeLimit with fee
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKAIA + KAIAFee), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKAIA), nil)
		if err != nil {
			log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
		}
		t.Log("RequestKLAYTransfer Transaction", tx.Hash().Hex())

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-2. Request KAIA transfer from Alice to Bob with zero feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKAIA), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKAIA), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-3. Request KAIA transfer from Alice to Bob with insufficient feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKAIA + (KAIAFee - 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKAIA), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// 9-4. Request KAIA transfer from Alice to Bob with enough feeLimit
	{
		tx, err = pBridge.RequestKLAYTransfer(&bind.TransactOpts{From: Alice.From, Signer: Alice.Signer, Value: big.NewInt(testKAIA + (KAIAFee + 1)), GasLimit: testGasLimit}, Bob.From, big.NewInt(testKAIA), nil)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// 9-4. Request KAIA transfer from Alice to Alice through () payable method
	{
		nonce, _ := sim.PendingNonceAt(context.Background(), Alice.From)
		gasPrice, _ := sim.SuggestGasPrice(context.Background())
		unsignedTx := types.NewTransaction(nonce, pBridgeAddr, big.NewInt(testKAIA+KAIAFee), testGasLimit, gasPrice, []byte{})

		chainID, _ := sim.ChainID(context.Background())
		tx, err = types.SignTx(unsignedTx, types.LatestSignerForChainID(chainID), AliceKey)
		sim.SendTransaction(context.Background(), tx)
		assert.Equal(t, nil, err)

		sim.Commit() // block

		CheckReceipt(sim, tx, 1*time.Second, types.ReceiptStatusSuccessful, t)
	}

	// Wait a few second for wait group
	WaitGroupWithTimeOut(&wg, 3*time.Second, t)

	// 10. Check Token balance
	{
		balance, err = token.BalanceOf(nil, Alice.From)
		assert.Equal(t, nil, err)
		t.Log("Alice token balance", balance.String())
		assert.Equal(t, initialValue-(testToken+ERC20Fee)*4, balance.Int64())

		balance, err = token.BalanceOf(nil, Bob.From)
		assert.Equal(t, nil, err)
		t.Log("Bob token balance", balance.String())
		assert.Equal(t, testToken*4, balance.Int64())

		balance, err = token.BalanceOf(nil, receiver.From)
		assert.Equal(t, nil, err)
		t.Log("Fee receiver token balance", balance.String())
		assert.Equal(t, ERC20Fee*4, balance.Int64())
	}

	// 11. Check KAIA balance
	{
		balance, _ = sim.BalanceAt(context.Background(), Alice.From, nil)
		t.Log("Alice KAIA balance :", balance)
		assert.Equal(t, initialValue-(testKAIA+KAIAFee)*2-KAIAFee, balance.Int64())

		balance, _ = sim.BalanceAt(context.Background(), Bob.From, nil)
		t.Log("Bob KAIA balance :", balance)
		assert.Equal(t, big.NewInt(testKAIA*2).String(), balance.String())

		balance, _ = sim.BalanceAt(context.Background(), receiver.From, nil)
		t.Log("receiver KAIA balance :", balance)
		assert.Equal(t, KAIAFee*3, balance.Int64())
	}

	bridgeManager.Stop()
}

// TestBasicJournal tests basic journal functionality.
func TestBasicJournal(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)

	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KAIA)},
		auth2.From:            {Balance: big.NewInt(params.KAIA)},
		auth4.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}
	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	localAddr, err := bm.DeployBridgeTest(sim, 10000, true)
	assert.NoError(t, err)
	remoteAddr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)

	bm.SetJournal("", localAddr, remoteAddr)

	ps := sc.BridgePeerSet()
	ps.peers["test"] = nil

	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	localInfo, ok := bm.GetBridgeInfo(localAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, localInfo.subscribed)

	remoteInfo, ok := bm.GetBridgeInfo(remoteAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, remoteInfo.subscribed)
}

// TestMethodRestoreBridges tests restoring bridges from the journal.
func TestMethodRestoreBridges(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true
	config.VTRecoveryInterval = 60

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KAIA)},
		auth2.From:            {Balance: big.NewInt(params.KAIA)},
		auth4.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	var bridgeAddrs [4]common.Address
	for i := 0; i < 4; i++ {
		if i%2 == 0 {
			bridgeAddrs[i], err = bm.DeployBridgeTest(sim, 10000, true)
		} else {
			bridgeAddrs[i], err = bm.DeployBridgeTest(sim, 10000, false)
		}
		if err != nil {
			t.Fatal("deploy bridge test failed", bridgeAddrs[i])
		}
		bm.DeleteBridgeInfo(bridgeAddrs[i])
	}
	sim.Commit()

	// Set journal
	bm.SetJournal("", bridgeAddrs[0], bridgeAddrs[1])
	bm.journal.cache[bridgeAddrs[0]].Subscribed = true
	bm.SetJournal("", bridgeAddrs[2], bridgeAddrs[3])
	bm.journal.cache[bridgeAddrs[2]].Subscribed = true

	ps := sc.BridgePeerSet()
	ps.peers["test"] = nil

	// Call RestoreBridges
	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	// Duplicated RestoreBridges
	if err := bm.RestoreBridges(); err != nil {
		t.Fatal("bm restoring bridges failed")
	}

	// Case 1: check bridge contract creation.
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.NotEqual(t, nil, info.bridge)
	}

	// Case 2: check subscription
	for i := 0; i < 4; i++ {
		info, _ := bm.GetBridgeInfo(bridgeAddrs[i])
		assert.Equal(t, true, info.subscribed)
	}

	// Case 3: check recovery
	recovery1 := bm.recoveries[bridgeAddrs[0]]
	assert.NotEqual(t, nil, recovery1)
	recovery1.Start()
	assert.Equal(t, nil, recovery1.WaitRunningStatus(true, 5*time.Second))
	recovery2 := bm.recoveries[bridgeAddrs[2]]
	assert.NotEqual(t, nil, recovery2)
	recovery2.Start()
	assert.Equal(t, nil, recovery2.WaitRunningStatus(true, 5*time.Second))

	bm.stopAllRecoveries()
	bm.Stop()
}

// TestMethodGetAllBridge tests a method GetAllBridge.
func TestMethodGetAllBridge(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	testBridge1 := common.BytesToAddress([]byte("test1"))
	testBridge2 := common.BytesToAddress([]byte("test2"))

	bm.journal.insert("", testBridge1, testBridge2)
	bm.journal.insert("", testBridge2, testBridge1)

	bridges := bm.GetAllBridge()
	assert.Equal(t, 2, len(bridges))

	bm.Stop()
}

// TestErrorDuplication tests if duplication of journal insertion is ignored or not.
func TestErrorDuplication(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	err = bm.journal.insert("", localAddr, remoteAddr)
	assert.Equal(t, nil, err)
	err = bm.journal.insert("", remoteAddr, localAddr)
	assert.Equal(t, nil, err)

	// try duplicated insert.
	err = bm.journal.insert("", localAddr, remoteAddr)
	assert.NotEqual(t, nil, err)
	err = bm.journal.insert("", remoteAddr, localAddr)
	assert.NotEqual(t, nil, err)

	// check cache size for checking duplication
	bridges := bm.GetAllBridge()
	assert.Equal(t, 2, len(bridges))

	bm.Stop()
}

// TestMethodSetJournal tests if duplication of journal insertion is ignored or not.
func TestMethodSetJournal(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	// Simple insert case
	err = bm.SetJournal("", localAddr, remoteAddr)
	assert.Equal(t, nil, err)

	// Error case
	err = bm.SetJournal("", localAddr, remoteAddr)
	assert.NotEqual(t, nil, err)

	// Check the number of bridge elements for checking duplication
	bridges := bm.GetAllBridge()
	assert.Equal(t, 1, len(bridges))

	bm.Stop()
}

// TestErrorDuplicatedSetBridgeInfo tests if duplication of bridge info insertion.
func TestErrorDuplicatedSetBridgeInfo(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KAIA)},
		auth2.From:            {Balance: big.NewInt(params.KAIA)},
		auth4.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)
	addr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)
	bridgeInfo, _ := bm.GetBridgeInfo(addr)

	// Try to call duplicated SetBridgeInfo
	err = bm.SetBridgeInfo(addr, bridgeInfo.bridge, common.Address{}, nil, sc.bridgeAccounts.pAccount, false, false)
	assert.NotEqual(t, nil, err)
	bm.Stop()
}

// TestScenarioSubUnsub tests subscription and unsubscription scenario.
func TestScenarioSubUnsub(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KAIA)},
		auth2.From:            {Balance: big.NewInt(params.KAIA)},
		auth4.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	localAddr, err := bm.DeployBridgeTest(sim, 10000, true)
	if err != nil {
		t.Fatal("deploy bridge test failed", localAddr)
	}

	bm.SubscribeEvent(localAddr)
	bridgeInfo, ok := bm.GetBridgeInfo(localAddr)
	assert.Equal(t, true, ok)
	assert.Equal(t, true, bridgeInfo.subscribed)
	bm.UnsubscribeEvent(localAddr)
	assert.Equal(t, false, bridgeInfo.subscribed)

	// Journal is irrelevant to the bridge unsubscription.
	journal := bm.journal.cache[localAddr]
	assert.NotEqual(t, nil, journal)
}

// TestErrorEmptyAccount tests empty account error in case of journal insertion.
func TestErrorEmptyAccount(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sc := &SubBridge{
		config: &SCConfig{DataDir: tempDir, VTRecovery: true},
		peers:  newBridgePeerSet(),
	}
	bm, err := NewBridgeManager(sc)
	if err != nil {
		t.Fatalf("fail to create bridge manager %v", err)
	}

	localAddr := common.BytesToAddress([]byte("test1"))
	remoteAddr := common.BytesToAddress([]byte("test2"))

	err = bm.journal.insert("", localAddr, common.Address{})
	assert.NotEqual(t, nil, err)

	err = bm.journal.insert("", common.Address{}, remoteAddr)
	assert.NotEqual(t, nil, err)

	bm.Stop()
}

// TestErrorDupSubscription tests duplicated subscription error.
func TestErrorDupSubscription(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	key4, _ := crypto.GenerateKey()
	auth4 := bind.NewKeyedTransactor(key4)
	config := &SCConfig{}
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		auth.From:             {Balance: big.NewInt(params.KAIA)},
		auth2.From:            {Balance: big.NewInt(params.KAIA)},
		auth4.From:            {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// 1. Prepare manager and subscribe event
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	addr, err := bm.DeployBridgeTest(sim, 10000, false)
	bridgeInfo, _ := bm.GetBridgeInfo(addr)
	bridge := bridgeInfo.bridge
	t.Log("===== BridgeContract Addr ", addr.Hex())
	sim.Commit() // block

	bm.bridges[addr], err = NewBridgeInfo(sc, addr, bridge, common.Address{}, nil, bacc.cAccount, true, true, sim)

	bm.journal.cache[addr] = &BridgeJournal{"", addr, addr, true, false}

	bm.SubscribeEvent(addr)
	err = bm.SubscribeEvent(addr)
	assert.NotEqual(t, nil, err)

	bm.Stop()
}

// TestAnchoringBasic tests the following:
// 1. generate anchoring tx
// 2. decode anchoring tx
// 3. start anchoring from the current block
// 4. accumulated tx counts
func TestAnchoringBasic(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoring")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sim, sc, bAcc, _, _, _ := generateAnchoringEnv(t, tempDir)
	defer sim.Close()

	assert.Equal(t, uint64(0), sc.handler.txCountStartingBlockNumber)
	assert.Equal(t, uint64(0), sc.handler.latestTxCountAddedBlockNumber)
	assert.Equal(t, uint64(1), sc.handler.chainTxPeriod)

	auth := bAcc.pAccount.GenerateTransactOpts()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk := sim.BlockChain().CurrentBlock()

	// nil block
	{
		err := sc.handler.blockAnchoringManager(nil)
		assert.Error(t, ErrInvalidBlock, err)
	}

	{
		err := sc.handler.generateAndAddAnchoringTxIntoTxPool(nil)
		assert.Error(t, ErrInvalidBlock, err)
	}
	// Generate anchoring tx again for the curBlk.
	err = sc.handler.blockAnchoringManager(curBlk)
	assert.NoError(t, err)

	pending := sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 1, len(pending))
	var tx *types.Transaction
	for _, v := range pending {
		assert.Equal(t, 1, len(v))
		tx = v[0]
	}
	assert.Equal(t, uint64(0), sc.handler.txCount)

	assert.Equal(t, curBlk.NumberU64(), sc.handler.latestTxCountAddedBlockNumber)
	compareBlockAndAnchoringTx(t, curBlk, tx)
}

// TestAnchoringBasicWithFeePayer tests the following with feePayer:
// 1. generate anchoring tx
// 2. decode anchoring tx
// 3. start anchoring from the current block
// 4. accumulated tx counts
func TestAnchoringBasicWithFeePayer(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoring")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sim, sc, bAcc, parentOperator, feePayer, tester := generateAnchoringEnv(t, tempDir)
	defer sim.Close()

	invalidAccount := common.HexToAddress("0x1")
	bAcc.SetParentOperatorFeePayer(feePayer.Address)

	assert.Equal(t, uint64(0), sc.handler.txCountStartingBlockNumber)
	assert.Equal(t, uint64(0), sc.handler.latestTxCountAddedBlockNumber)
	assert.Equal(t, uint64(1), sc.handler.chainTxPeriod)

	// fail to generate anchoring tx with invalid parent operator
	{
		pAccBackup := bAcc.pAccount.address
		bAcc.pAccount.address = invalidAccount

		curBlk := sim.BlockChain().CurrentBlock()
		err = sc.handler.generateAndAddAnchoringTxIntoTxPool(curBlk)
		assert.Error(t, err, accounts.ErrUnknownAccount)

		bAcc.pAccount.address = pAccBackup
	}

	// fail to generate anchoring tx with invalid feePayer
	{
		bAcc.SetParentOperatorFeePayer(invalidAccount)

		curBlk := sim.BlockChain().CurrentBlock()
		err = sc.handler.generateAndAddAnchoringTxIntoTxPool(curBlk)
		assert.Error(t, err, accounts.ErrUnknownAccount)

		bAcc.SetParentOperatorFeePayer(feePayer.Address)
	}

	_, _, _, err = bridge.DeployBridge(tester, sim, true) // dummy tx
	sim.Commit()
	curBlk := sim.BlockChain().CurrentBlock()

	// Generate anchoring tx again for the curBlk.
	sc.handler.blockAnchoringManager(curBlk)
	pending := sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 1, len(pending))
	var tx *types.Transaction
	for _, v := range pending {
		assert.Equal(t, 1, len(v))
		tx = v[0]

		// Check Balance
		feePayerBalanceBefore, err := sim.BalanceAt(context.Background(), feePayer.Address, nil)
		assert.NoError(t, err)
		parentOperatorBalanceBefore, err := sim.BalanceAt(context.Background(), parentOperator.address, nil)
		assert.NoError(t, err)

		sim.SendTransaction(context.Background(), tx)
		sim.Commit()

		// Check Balance
		feePayerBalanceAfter, err := sim.BalanceAt(context.Background(), feePayer.Address, nil)
		assert.NoError(t, err)
		parentOperatorBalanceAfter, err := sim.BalanceAt(context.Background(), parentOperator.address, nil)
		assert.NoError(t, err)

		receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
		assert.NoError(t, err)

		fee := new(big.Int).SetUint64(receipt.GasUsed * params.DefaultUnitPrice)

		assert.Equal(t, new(big.Int).Sub(feePayerBalanceBefore, fee).String(), feePayerBalanceAfter.String())
		t.Log("feePayer paid ", fee)
		assert.Equal(t, parentOperatorBalanceBefore, parentOperatorBalanceAfter)
	}

	assert.Equal(t, uint64(0), sc.handler.txCount)
	assert.Equal(t, curBlk.NumberU64(), sc.handler.latestTxCountAddedBlockNumber)
	compareBlockAndAnchoringTx(t, curBlk, tx)
}

// TestAnchoringBasicWithBridgeTxPoolMock tests the following :
// - BridgeTxPool addLocal() fail case.
func TestAnchoringBasicWithBridgeTxPoolMock(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoring")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	sim, sc, bAcc, _, feePayer, _ := generateAnchoringEnv(t, tempDir)
	defer sim.Close()

	// mock BridgeTxPool
	{
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()
		mockBridgeTxPool := NewMockBridgeTxPool(mockCtrl)
		sc.bridgeTxPool = mockBridgeTxPool
		mockBridgeTxPool.EXPECT().AddLocal(gomock.Any()).Return(bridgepool.ErrKnownTx)
	}

	bAcc.SetParentOperatorFeePayer(feePayer.Address)

	assert.Equal(t, uint64(0), sc.handler.txCountStartingBlockNumber)
	assert.Equal(t, uint64(0), sc.handler.latestTxCountAddedBlockNumber)
	assert.Equal(t, uint64(1), sc.handler.chainTxPeriod)

	curBlk := sim.BlockChain().CurrentBlock()

	// Generate anchoring tx with mocked BridgeTxPool returns a error
	err = sc.handler.blockAnchoringManager(curBlk)
	assert.Equal(t, bridgepool.ErrKnownTx, err)
}

func generateAnchoringEnv(t *testing.T, tempDir string) (*backends.SimulatedBackend, *SubBridge, *BridgeAccounts, *accountInfo, accounts.Account, *bind.TransactOpts) {
	config := &SCConfig{AnchoringPeriod: 1}
	config.DataDir = tempDir
	config.VTRecovery = true

	ks := keystore.NewKeyStore(tempDir, keystore.StandardScryptN, keystore.StandardScryptP)
	back := []accounts.Backend{
		ks,
	}
	am := accounts.NewManager(back...)
	bAcc, _ := NewBridgeAccounts(am, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bAcc.pAccount.chainID = big.NewInt(0)
	bAcc.cAccount.chainID = big.NewInt(0)
	parentOperator := bAcc.pAccount

	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)

	initBal := new(big.Int).Exp(big.NewInt(10), big.NewInt(50), nil)

	feePayer, err := ks.NewAccount("pwd")
	assert.NoError(t, err)
	ks.TimedUnlock(feePayer, "pwd", 0)

	alloc := blockchain.GenesisAlloc{
		alice.From:             {Balance: initBal},
		feePayer.Address:       {Balance: initBal},
		parentOperator.address: {Balance: initBal},
	}
	sim := backends.NewSimulatedBackendWithGasPrice(alloc, params.DefaultUnitPrice)

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bAcc,
	}
	sc.blockchain = sim.BlockChain()

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
	}

	sc.handler.setRemoteGasPrice(params.DefaultUnitPrice)

	sc.bridgeTxPool = bridgepool.NewBridgeTxPool(bridgepool.BridgeTxPoolConfig{
		Journal:     path.Join(tempDir, "bridge_transactions.rlp"),
		GlobalQueue: 1024,
	})

	return sim, sc, bAcc, parentOperator, feePayer, alice
}

func compareBlockAndAnchoringTx(t *testing.T, block *types.Block, tx *types.Transaction) {
	// Decoding the anchoring tx.
	assert.Equal(t, true, tx.Type().IsChainDataAnchoring())
	anchoringData := new(types.AnchoringData)
	data, err := tx.AnchoredData()
	assert.NoError(t, err)

	err = rlp.DecodeBytes(data, anchoringData)
	assert.NoError(t, err)
	assert.Equal(t, types.AnchoringDataType0, anchoringData.Type)
	anchoringDataInternal := new(types.AnchoringDataInternalType0)
	if err := rlp.DecodeBytes(anchoringData.Data, anchoringDataInternal); err != nil {
		logger.Error("writeChildChainTxHashFromBlock : failed to decode anchoring data")
	}

	// Check the current block is anchored.
	assert.Equal(t, new(big.Int).SetUint64(block.NumberU64()).String(), anchoringDataInternal.BlockNumber.String())
	assert.Equal(t, block.Hash(), anchoringDataInternal.BlockHash)
	assert.Equal(t, big.NewInt(1).String(), anchoringDataInternal.BlockCount.String())
	assert.Equal(t, big.NewInt(1).String(), anchoringDataInternal.TxCount.String())
}

// TestAnchoringStart tests the following:
// 1. set anchoring period 4
// 2. check if tx counting started immediately and accumulated correctly
func TestAnchoringStart(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoringPeriod")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	config := &SCConfig{AnchoringPeriod: 4}
	config.DataDir = tempDir
	config.VTRecovery = true

	bAcc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bAcc.pAccount.chainID = big.NewInt(0)
	bAcc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bAcc,
	}
	sc.blockchain = sim.BlockChain()

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}
	sc.bridgeTxPool = bridgepool.NewBridgeTxPool(bridgepool.BridgeTxPoolConfig{
		Journal:     path.Join(tempDir, "bridge_transactions.rlp"),
		GlobalQueue: 1024,
	})

	assert.Equal(t, uint64(0), sc.handler.txCountStartingBlockNumber)
	assert.Equal(t, uint64(4), sc.handler.chainTxPeriod)

	sim.Commit() // start with arbitrary block number.

	// 1. Fresh start with dummy tx and check tx count
	auth := bAcc.pAccount.GenerateTransactOpts()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk := sim.BlockChain().CurrentBlock()
	sc.handler.blockAnchoringManager(curBlk)
	assert.Equal(t, uint64(1), sc.handler.txCount)
	pending := sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 0, len(pending)) // the anchoring period has not yet been reached.

	// 2. Generate dummy txs and check tx count
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk = sim.BlockChain().CurrentBlock()
	sc.handler.blockAnchoringManager(curBlk)
	assert.Equal(t, uint64(3), sc.handler.txCount)
	assert.Equal(t, 0, len(pending)) // the anchoring period has not yet been reached.

	// 3. Generate dummy blocks and check anchoring tx
	sim.Commit() // block number 4
	curBlk = sim.BlockChain().CurrentBlock()
	sc.handler.blockAnchoringManager(curBlk)
	assert.Equal(t, uint64(0), sc.handler.txCount)
	pending = sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 1, len(pending))
	for _, v := range pending {
		decodeAndCheckAnchoringTx(t, v[0], curBlk, 3, 3)
		break
	}
}

// TestAnchoringPeriod tests the following:
// 1. set anchoring period 4
// 2. accumulate tx counts
func TestAnchoringPeriod(t *testing.T) {
	const (
		startTxCount = 100
	)
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoringPeriod")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	config := &SCConfig{AnchoringPeriod: 4}
	config.DataDir = tempDir
	config.VTRecovery = true

	bAcc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bAcc.pAccount.chainID = big.NewInt(0)
	bAcc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bAcc,
	}
	sc.blockchain = sim.BlockChain()

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}
	sc.bridgeTxPool = bridgepool.NewBridgeTxPool(bridgepool.BridgeTxPoolConfig{
		Journal:     path.Join(tempDir, "bridge_transactions.rlp"),
		GlobalQueue: 1024,
	})

	assert.Equal(t, uint64(0), sc.handler.txCountStartingBlockNumber)
	assert.Equal(t, uint64(4), sc.handler.chainTxPeriod)

	// Period 1
	sim.Commit()
	auth := bAcc.pAccount.GenerateTransactOpts()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk := sim.BlockChain().CurrentBlock()

	sc.handler.txCountStartingBlockNumber = curBlk.NumberU64() - 1
	sc.handler.txCount = startTxCount
	sc.handler.blockAnchoringManager(curBlk)
	assert.Equal(t, uint64(startTxCount+1), sc.handler.txCount)
	pending := sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 0, len(pending)) // the anchoring period has not yet been reached.

	// Generate anchoring tx for the curBlk.
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk = sim.BlockChain().CurrentBlock()
	sc.handler.blockAnchoringManager(curBlk)
	pending = sc.GetBridgeTxPool().Pending()
	assert.Equal(t, 1, len(pending))

	for _, v := range pending {
		decodeAndCheckAnchoringTx(t, v[0], curBlk, 4, startTxCount+4)
		break
	}

	// Period 2:
	assert.Equal(t, uint64(0), sc.handler.txCount)
	assert.Equal(t, uint64(5), sc.handler.txCountStartingBlockNumber)

	// Generate anchoring tx.
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	sim.Commit()
	sim.Commit()

	curBlk = sim.BlockChain().CurrentBlock()
	sc.handler.blockAnchoringManager(curBlk)
	pending = sc.GetBridgeTxPool().Pending()
	for _, v := range pending {
		decodeAndCheckAnchoringTx(t, v[1], curBlk, 4, 3)
		break
	}
}

// decodeAndCheckAnchoringTx decodes anchoring tx and check with a block.
func decodeAndCheckAnchoringTx(t *testing.T, tx *types.Transaction, blk *types.Block, blockCount, txCounts int64) {
	assert.Equal(t, types.TxTypeChainDataAnchoring, tx.Type())
	anchoringData := new(types.AnchoringData)
	data, err := tx.AnchoredData()
	assert.NoError(t, err)

	err = rlp.DecodeBytes(data, anchoringData)
	assert.NoError(t, err)
	assert.Equal(t, types.AnchoringDataType0, anchoringData.Type)
	anchoringDataInternal := new(types.AnchoringDataInternalType0)
	if err := rlp.DecodeBytes(anchoringData.Data, anchoringDataInternal); err != nil {
		logger.Error("writeChildChainTxHashFromBlock : failed to decode anchoring data")
	}

	// Check the current block is anchored.
	assert.Equal(t, new(big.Int).SetUint64(blk.NumberU64()).String(), anchoringDataInternal.BlockNumber.String())
	assert.Equal(t, blk.Hash(), anchoringDataInternal.BlockHash)
	assert.Equal(t, big.NewInt(blockCount).String(), anchoringDataInternal.BlockCount.String())
	assert.Equal(t, big.NewInt(txCounts).String(), anchoringDataInternal.TxCount.String())
}

// TestDecodingLegacyAnchoringTx tests the following:
// 1. generate AnchoringDataLegacy anchoring tx
// 2. decode AnchoringDataLegacy with a decoding method of a sub-bridge handler.
func TestDecodingLegacyAnchoringTx(t *testing.T) {
	const (
		startBlkNum  = 10
		startTxCount = 100
	)
	tempDir, err := os.MkdirTemp(os.TempDir(), "anchoring")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	config := &SCConfig{AnchoringPeriod: 1}
	config.DataDir = tempDir
	config.VTRecovery = true

	bAcc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bAcc.pAccount.chainID = big.NewInt(0)
	bAcc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bAcc,
	}
	sc.blockchain = sim.BlockChain()

	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	// Encoding anchoring tx.
	auth := bAcc.pAccount.GenerateTransactOpts()
	_, _, _, err = bridge.DeployBridge(auth, sim, true) // dummy tx
	sim.Commit()
	curBlk := sim.BlockChain().CurrentBlock()

	anchoringData := &types.AnchoringDataLegacy{
		BlockHash:     curBlk.Hash(),
		TxHash:        curBlk.Header().TxHash,
		ParentHash:    curBlk.Header().ParentHash,
		ReceiptHash:   curBlk.Header().ReceiptHash,
		StateRootHash: curBlk.Header().Root,
		BlockNumber:   curBlk.Header().Number,
	}
	data, err := rlp.EncodeToBytes(anchoringData)
	assert.NoError(t, err)

	// Decoding the anchoring tx.
	decodedData, err := types.DecodeAnchoringData(data)
	assert.Equal(t, curBlk.Hash(), decodedData.GetBlockHash())
	assert.Equal(t, curBlk.Header().Number.String(), decodedData.GetBlockNumber().String())
}

func TestBridgeAliasAPIs(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Generate a new random account and a funded simulator
	aliceKey, _ := crypto.GenerateKey()
	alice := bind.NewKeyedTransactor(aliceKey)
	bobKey, _ := crypto.GenerateKey()
	bob := bind.NewKeyedTransactor(bobKey)

	config := &SCConfig{}
	config.DataDir = tempDir

	bacc, _ := NewBridgeAccounts(nil, tempDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	alloc := blockchain.GenesisAlloc{
		alice.From:            {Balance: big.NewInt(params.KAIA)},
		bob.From:              {Balance: big.NewInt(params.KAIA)},
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		config:         config,
		peers:          newBridgePeerSet(),
		localBackend:   sim,
		remoteBackend:  sim,
		bridgeAccounts: bacc,
	}

	sc.APIBackend = &SubBridgeAPI{sc}
	sc.handler, err = NewSubBridgeHandler(sc)
	assert.NoError(t, err)

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)
	sc.handler.subbridge.bridgeManager = bm

	// 1. Deploy bridge contracts and register them
	cBridgeAddr := deployBridge(t, bm, sim, true)
	pBridgeAddr := deployBridge(t, bm, sim, false)

	// 2. Deploy token Contracts
	cTokenAddr, _, _, err := sctoken.DeployServiceChainToken(alice, sim, cBridgeAddr)
	assert.NoError(t, err)
	pTokenAddr, _, _, err := sctoken.DeployServiceChainToken(alice, sim, pBridgeAddr)
	assert.NoError(t, err)
	sim.Commit() // block

	cBridgeAddrStr := cBridgeAddr.String()
	pBridgeAddrStr := pBridgeAddr.String()
	cTokenAddrStr := cTokenAddr.String()
	pTokenAddrStr := pTokenAddr.String()
	// -------------------------- Done prepration --------------------------

	// -------------------------- API test with the raw addresss format --------------------------
	{
		// TEST 1-1 - Success (Register bridge, tokens and subscribe registered bridges)
		bridgePairs := bm.subBridge.APIBackend.ListBridge()
		assert.Equal(t, len(bridgePairs), 0)

		testBridgeAPIBasic(t, bm, cBridgeAddr, pBridgeAddr, nil, cTokenAddr, pTokenAddr)

		bridgePairs = bm.subBridge.APIBackend.ListBridge()
		assert.Equal(t, len(bridgePairs), 1)
		assert.Equal(t, bridgePairs[0].Subscribed, true)
	}

	{
		// TEST 1-2 - Failure
		// Duplicated journal
		testDuplicatedJournal(t, bm, cBridgeAddr, pBridgeAddr, nil)

		// Duplicated token
		testDuplicatedToken(t, bm, &cBridgeAddrStr, &pBridgeAddrStr, &cTokenAddrStr, &pTokenAddrStr)

		// Already subscribed
		testAlreadySubscribed(t, bm, &cBridgeAddrStr, &pBridgeAddrStr)
	}

	{
		// TEST 1-3 - Success (Unsubscribe bridge, deregister bridges and tokens)
		testUnsubscribeAndDeRegister(t, bm, &cBridgeAddrStr, &pBridgeAddrStr, &cTokenAddrStr, &pTokenAddrStr)
	}

	// -------------------------- API test with the bridge format --------------------------
	alias := "MYBRIDGE"
	changedAlias := "MYBRIDGE_v2"
	invalidBridgeAlias := "0xMYBRIDGE"
	{
		// TEST 2-1 - Success (Register bridge, tokens and subscribe registered bridges)
		bridgePairs := bm.subBridge.APIBackend.ListBridge()
		assert.Equal(t, len(bridgePairs), 0)

		testBridgeAPIBasic(t, bm, cBridgeAddr, pBridgeAddr, &alias, cTokenAddr, pTokenAddr)

		bridgePairs = bm.subBridge.APIBackend.ListBridge()
		assert.Equal(t, len(bridgePairs), 1)
		assert.Equal(t, bridgePairs[0].Subscribed, true)
		bridgePair := bm.subBridge.APIBackend.GetBridgePairByAlias(alias)
		assert.Equal(t, bridgePair.BridgeAlias, alias)
	}

	{
		// TEST 2-2 - Failure
		// Duplicated bridge alias
		testDuplicatedJournal(t, bm, cBridgeAddr, pBridgeAddr, &alias)

		// Duplicated token
		testDuplicatedToken(t, bm, &alias, &cTokenAddrStr, &pTokenAddrStr, nil)

		// Already subscribed
		testAlreadySubscribed(t, bm, &alias, nil)
	}

	{
		// TEST 2-3 - Success (change bridge alias)
		err = bm.subBridge.APIBackend.ChangeBridgeAlias(alias, changedAlias)
		assert.NoError(t, err)

		// Try to deregister with empty bridge alias
		err = bm.subBridge.APIBackend.DeregisterBridge(&alias, nil)
		assert.Equal(t, err, ErrEmptyBridgeAlias)

		bridgePair := bm.subBridge.APIBackend.GetBridgePairByAlias(alias)
		assert.Nil(t, bridgePair)

		bridgePair = bm.subBridge.APIBackend.GetBridgePairByAlias(changedAlias)
		assert.Equal(t, bridgePair.BridgeAlias, changedAlias)
	}

	{
		// TEST 2-4 - Success (Unsubscribe bridge, deregister bridges and tokens)
		testUnsubscribeAndDeRegister(t, bm, &changedAlias, &cTokenAddrStr, &pTokenAddrStr, nil)

		// TEST 2-5 - Failure (Unsubscribe register bridge already unsubscribed)
		err := bm.subBridge.APIBackend.UnsubscribeBridge(&changedAlias, nil)
		assert.Equal(t, err, ErrEmptyBridgeAlias)
	}

	{
		// TEST 2-6 - Failure (Try to create a bridge alias with invalid bridge alias name)
		err = bm.subBridge.APIBackend.RegisterBridge(cBridgeAddr, pBridgeAddr, &invalidBridgeAlias)
		assert.Equal(t, err, ErrNotAllowedAliasFormat)
	}

	{
		// TEST 3 - Concurrent API Call
		contractPairLen := 10
		bridgeAddrs := make([]common.Address, contractPairLen)
		tokenAddrs := make([]common.Address, contractPairLen)
		t.Logf("Prepare %d contracts\n", contractPairLen)
		// Preparation: Deploy bridge and token contracts
		for i := 0; i < contractPairLen/2; i++ {
			cBridgeAddr, pBridgeAddr = deployBridge(t, bm, sim, true), deployBridge(t, bm, sim, false)
			cIdx, pIdx := i*2, i*2+1
			bridgeAddrs[cIdx], bridgeAddrs[pIdx] = cBridgeAddr, pBridgeAddr

			cTokenAddr, _, _, err := sctoken.DeployServiceChainToken(alice, sim, cBridgeAddr)
			assert.NoError(t, err)
			pTokenAddr, _, _, err := sctoken.DeployServiceChainToken(alice, sim, pBridgeAddr)
			assert.NoError(t, err)
			tokenAddrs[cIdx], tokenAddrs[pIdx] = cTokenAddr, pTokenAddr

			t.Logf("Deployed bridge contracts %d, %d\n", cIdx, pIdx)
		}

		// Declare another bridge and token contracts that did not initialize
		fixedChildBridgeAddr, fixedParentBridgeAddr := deployBridge(t, bm, sim, true), deployBridge(t, bm, sim, false)

		const (
			BRIDGE_SETUP = iota
			FAILURE
			CLEANUP_BRIDGE
			ALIAS_BRIDGE_SETUP
			ALIAS_FAILURE
			ALIAS_CLEANUP_BRIDGE
			REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE
		)
		// DO NOT CHANGE THE TEST ORDER
		testCases := map[uint8]string{
			BRIDGE_SETUP:         "BRIDGE_SETUP",
			FAILURE:              "FAILURE",
			CLEANUP_BRIDGE:       "CLEANUP_BRIDGE",
			ALIAS_BRIDGE_SETUP:   "ALIAS_BRIDGE_SETUP",
			ALIAS_FAILURE:        "ALIAS_FAILURE",
			ALIAS_CLEANUP_BRIDGE: "ALIAS_CLEANUP_BRIDGE",
			REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE: "REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE",
		}

		for testNum := 0; testNum < len(testCases); testNum++ {
			wg := sync.WaitGroup{}
			wg.Add(contractPairLen / 2)
			for i := 0; i < len(bridgeAddrs); i += 2 {
				cIdx, pIdx := i, i+1
				go func(cIdx, pIdx int) {
					cBridgeAddr, pBridgeAddr := bridgeAddrs[cIdx], bridgeAddrs[pIdx]
					cBridgeAddrStr, pBridgeAddrStr := cBridgeAddr.String(), pBridgeAddr.String()
					cTokenAddr, pTokenAddr := tokenAddrs[cIdx], tokenAddrs[pIdx]
					cTokenAddrStr, pTokenAddrStr := cTokenAddr.String(), pTokenAddr.String()
					alias := "MYBRIDGE_v3" + strconv.Itoa(cIdx)

					switch testNum {
					case BRIDGE_SETUP:
						// TEST 3-1. `testBridgeAPIBasic` again with concurrent calls using raw-address-format APIS
						testBridgeAPIBasic(t, bm, cBridgeAddr, pBridgeAddr, nil, cTokenAddr, pTokenAddr)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case FAILURE:
						// TEST 3-2 - Failure
						testDuplicatedJournal(t, bm, cBridgeAddr, pBridgeAddr, nil)
						testDuplicatedToken(t, bm, &cBridgeAddrStr, &pBridgeAddrStr, &cTokenAddrStr, &pTokenAddrStr)
						testAlreadySubscribed(t, bm, &cBridgeAddrStr, &pBridgeAddrStr)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case ALIAS_BRIDGE_SETUP:
						// TEST 3-3. `testBridgeAPIBasic` again with concurrent calls using alias APIS
						testBridgeAPIBasic(t, bm, cBridgeAddr, pBridgeAddr, &alias, cTokenAddr, pTokenAddr)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case ALIAS_FAILURE:
						// TEST 3-4 - Failure
						testDuplicatedJournal(t, bm, cBridgeAddr, pBridgeAddr, &alias)
						testDuplicatedToken(t, bm, &alias, &cTokenAddrStr, &pTokenAddrStr, nil)
						testAlreadySubscribed(t, bm, &alias, nil)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case CLEANUP_BRIDGE:
						// TEST 3-5 - Success (Unsubscribe bridge, deregister bridges and tokens)
						testUnsubscribeAndDeRegister(t, bm, &cBridgeAddrStr, &pBridgeAddrStr, &cTokenAddrStr, &pTokenAddrStr)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case ALIAS_CLEANUP_BRIDGE:
						// TEST 3-6 - Success (Unsubscribe bridge, deregister bridges and tokens)
						testUnsubscribeAndDeRegister(t, bm, &alias, &cTokenAddrStr, &pTokenAddrStr, nil)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					case REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE:
						// TEST 3-7 - Use the fresh bridge that did not register any token contracts
						cbAddr, pbAddr := fixedChildBridgeAddr.String(), fixedParentBridgeAddr.String()
						testRegisterToken(t, bm, &cbAddr, &pbAddr, &cTokenAddrStr, &pTokenAddrStr)
						t.Log("passed:", testCases[uint8(testNum)], cIdx, pIdx)
					}
					wg.Done()
				}(cIdx, pIdx)
			}
			wg.Wait()
			t.Log("Test Done: ", testCases[uint8(testNum)])
			// Check the status of conccurent calls with a signle thread
			switch testNum {
			case BRIDGE_SETUP:
				checkBridgeSetup(t, bm, true, 1, contractPairLen/2)
			case CLEANUP_BRIDGE:
				checkBridgeSetup(t, bm, false, 0, 0)
			case ALIAS_BRIDGE_SETUP:
				checkBridgeSetup(t, bm, true, 1, contractPairLen/2)
			case ALIAS_CLEANUP_BRIDGE:
				checkBridgeSetup(t, bm, false, 0, 0)
				// Initiailize another two bridge contracts for the test `REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE`
				err = bm.subBridge.APIBackend.RegisterBridge(fixedChildBridgeAddr, fixedParentBridgeAddr, nil)
				assert.NoError(t, err)
			case REGISTER_MULTIPLE_TOKEN_WITH_SINGLE_BRIDGE:
				checkRegisterMultipleToken(t, bm, fixedChildBridgeAddr, fixedParentBridgeAddr, contractPairLen/2)
			}
		}
		t.Log("All Done")
	}
}

func testBridgeAPIBasic(t *testing.T, bm *BridgeManager,
	cBridgeAddr, pBridgeAddr common.Address,
	alias *string,
	cTokenAddr, pTokenAddr common.Address,
) {
	// TEST 1 - Success (Register bridge, tokens and subscribe registered bridges)
	cBridgeAddrStr := cBridgeAddr.String()
	pBridgeAddrStr := pBridgeAddr.String()
	cTokenAddrStr := cTokenAddr.String()
	pTokenAddrStr := pTokenAddr.String()

	// Register Bridge
	err := bm.subBridge.APIBackend.RegisterBridge(cBridgeAddr, pBridgeAddr, alias)
	assert.NoError(t, err)

	// Register tokens
	if alias != nil {
		err = bm.subBridge.APIBackend.RegisterToken(alias, &cTokenAddrStr, &pTokenAddrStr, nil)
	} else {
		err = bm.subBridge.APIBackend.RegisterToken(&cBridgeAddrStr, &pBridgeAddrStr, &cTokenAddrStr, &pTokenAddrStr)
	}
	assert.NoError(t, err)

	// Subscribe bridges
	if alias != nil {
		err = bm.subBridge.APIBackend.SubscribeBridge(alias, nil)
	} else {
		err = bm.subBridge.APIBackend.SubscribeBridge(&cBridgeAddrStr, &pBridgeAddrStr)
	}
	assert.NoError(t, err)
}

func testDuplicatedJournal(t *testing.T, bm *BridgeManager, cBridgeAddr, pBridgeAddr common.Address, alias *string) {
	err := bm.subBridge.APIBackend.RegisterBridge(cBridgeAddr, pBridgeAddr, alias)
	if err != ErrDuplicatedJournal && err != ErrDuplicatedAlias {
		t.Fatal("Unexpected error", err)
	}
}

func testDuplicatedToken(t *testing.T, bm *BridgeManager, cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr *string) {
	err := bm.subBridge.APIBackend.RegisterToken(cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr)
	assert.Equal(t, err, ErrDuplicatedToken)
}

func testAlreadySubscribed(t *testing.T, bm *BridgeManager, cBridgeAddrStr, pBridgeAddrStr *string) {
	err := bm.subBridge.APIBackend.SubscribeBridge(cBridgeAddrStr, pBridgeAddrStr)
	assert.Equal(t, err, ErrAlreadySubscribed)
}

func testRegisterToken(t *testing.T, bm *BridgeManager, cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr *string) {
	err := bm.subBridge.APIBackend.RegisterToken(cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr)
	assert.NoError(t, err)
}

func testUnsubscribeAndDeRegister(t *testing.T, bm *BridgeManager,
	cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr *string,
) {
	findBridgePair := func(addrStr string) *BridgeJournal {
		bridgePairs := bm.subBridge.APIBackend.ListBridge()
		for _, bridgePair := range bridgePairs {
			if bridgePair.ChildAddress.String() == addrStr || bridgePair.BridgeAlias == addrStr {
				return bridgePair
			}
		}
		return nil
	}
	bridgePairLen := len(bm.subBridge.APIBackend.ListBridge())
	defer func() {
		err := bm.subBridge.APIBackend.DeregisterBridge(cBridgeAddrStr, pBridgeAddrStr)
		assert.NoError(t, err)
		assert.Equal(t, len(bm.subBridge.APIBackend.ListBridge()) < bridgePairLen, true)
	}()

	bridgePair := findBridgePair(*cBridgeAddrStr)
	assert.NotNil(t, bridgePair)
	assert.Equal(t, bridgePair.Subscribed, true)
	err := bm.subBridge.APIBackend.UnsubscribeBridge(cBridgeAddrStr, pBridgeAddrStr)
	assert.NoError(t, err)
	assert.Equal(t, bridgePair.Subscribed, false)

	cBi, ok := bm.GetBridgeInfo(bridgePair.ChildAddress)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(cBi.counterpartToken), 1)
	pBi, ok := bm.GetBridgeInfo(bridgePair.ParentAddress)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(pBi.counterpartToken), 1)

	err = bm.subBridge.APIBackend.DeregisterToken(cBridgeAddrStr, pBridgeAddrStr, cTokenAddrStr, pTokenAddrStr)
	assert.NoError(t, err)
	assert.Equal(t, len(cBi.counterpartToken), 0)
	assert.Equal(t, len(pBi.counterpartToken), 0)
}

func TestLegacyBridgeJournalDecode(t *testing.T) {
	// `encodedJournalHexStr` is an encoded legacy bridge journals. The code below generate the following hex.
	encodedJournalHexStr := "eb9485564429cce278d4399436f1af2f91e1be6f0bd494c12701e0cb09d6600f774be1dbb585ddc749f9da80eb9485564429cce278d4399436f1af2f91e1be6f0bd594c12701e0cb09d6600f774be1dbb585ddc749f9db01eb9485564429cce278d4399436f1af2f91e1be6f0bd694c12701e0cb09d6600f774be1dbb585ddc749f9dc01"
	/*
		legacyJournals := []BridgeJournal{
			{
				ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd4"),
				ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9da"),
				Subscribed:    false,
			},
			{
				ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd5"),
				ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9db"),
				Subscribed:    true,
			},
			{
				ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd6"),
				ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9dc"),
				Subscribed:    true,
			},
		}
			encodedBuf := new(bytes.Buffer)
			for i := 0; i < len(journals); i++ {
				err := rlp.Encode(encodedBuf, &journals[i])
				assert.NoError(t, err)
			}
			encodedString := hex.EncodeToString(encodedBuf.Bytes())
			fmt.Println(encodedString)
	*/

	legacyJournals := []BridgeJournal{
		{
			ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd4"),
			ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9da"),
			Subscribed:    false,
		},
		{
			ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd5"),
			ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9db"),
			Subscribed:    true,
		},
		{
			ChildAddress:  common.HexToAddress("0x85564429cce278d4399436f1af2f91e1be6f0bd6"),
			ParentAddress: common.HexToAddress("0xc12701e0cb09d6600f774be1dbb585ddc749f9dc"),
			Subscribed:    true,
		},
	}

	tempJournalPath := "legacy-journal-decoding-test"
	tempFile, err := os.CreateTemp(".", tempJournalPath)
	assert.NoError(t, err)

	encodedHex, err := hex.DecodeString(encodedJournalHexStr)
	assert.NoError(t, err)

	_, err = tempFile.Write(encodedHex)
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	readJournalIdx := 0
	load := func(gwjournal BridgeJournal) error {
		assert.Equal(t, len(gwjournal.BridgeAlias), 0)
		assert.Equal(t, gwjournal.ChildAddress.String(), legacyJournals[readJournalIdx].ChildAddress.String())
		assert.Equal(t, gwjournal.ParentAddress.String(), legacyJournals[readJournalIdx].ParentAddress.String())
		assert.Equal(t, gwjournal.Subscribed, legacyJournals[readJournalIdx].Subscribed)
		readJournalIdx++
		return nil
	}
	journalAddr := newBridgeAddrJournal(tempFile.Name())
	err = journalAddr.load(load)
	assert.NoError(t, err)
}

func checkBridgeSetup(t *testing.T, bm *BridgeManager, expectedSubscribed bool, expectedNumberOfToken, expectedBridgeLen int) {
	bridgePairs := bm.subBridge.APIBackend.ListBridge()
	assert.Equal(t, len(bridgePairs), expectedBridgeLen)
	for _, bridgePair := range bridgePairs {
		assert.Equal(t, bridgePair.Subscribed, expectedSubscribed)
		cbi, ok := bm.GetBridgeInfo(bridgePair.ChildAddress)
		assert.Equal(t, ok, true)
		pbi, ok := bm.GetBridgeInfo(bridgePair.ChildAddress)
		assert.Equal(t, ok, true)
		assert.Equal(t, len(cbi.counterpartToken), expectedNumberOfToken)
		assert.Equal(t, len(pbi.counterpartToken), expectedNumberOfToken)
	}
}

func checkRegisterMultipleToken(t *testing.T, bm *BridgeManager, cBridgeAddr, pBridgeAddr common.Address, expectedLen int) {
	cbi, ok := bm.GetBridgeInfo(cBridgeAddr)
	assert.Equal(t, ok, true)
	pbi, ok := bm.GetBridgeInfo(pBridgeAddr)
	assert.Equal(t, ok, true)
	assert.Equal(t, len(cbi.counterpartToken), expectedLen)
	assert.Equal(t, len(pbi.counterpartToken), expectedLen)
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := crand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func TestBridgeAddressType(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}
	sc.handler, err = NewSubBridgeHandler(sc)
	assert.NoError(t, err)
	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)

	auth := bacc.cAccount.GenerateTransactOpts()

	// Deploy Bridge Contract
	bridgeAddr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)
	sim.Commit() // block

	anotherBridgeAddr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)
	sim.Commit() // block

	// Case 1 - Success (The bridge address type is contract address)
	{
		// 1. Deploy Token Contract
		_, tx, token, err := sctoken.DeployServiceChainToken(auth, sim, bridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)

		// 2. Deploy NFT Contract
		_, tx, nft, err := scnft.DeployServiceChainNFT(auth, sim, bridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)

		tx, err = token.SetBridge(auth, anotherBridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)
		tx, err = nft.SetBridge(auth, anotherBridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)
	}

	// Case 2 - Failure (The bridge address type is not a contract address)
	{
		_, tx, _, err := sctoken.DeployServiceChainToken(auth, sim, auth.From)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusErrExecutionReverted, t)

		_, tx, _, err = scnft.DeployServiceChainNFT(auth, sim, auth.From)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}

	// Case 3 - Failure (The bridge address type is not a contract address)
	{
		_, tx, token, err := sctoken.DeployServiceChainToken(auth, sim, bridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)

		_, tx, nft, err := scnft.DeployServiceChainNFT(auth, sim, bridgeAddr)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusSuccessful, t)

		tx, err = token.SetBridge(auth, auth.From)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusErrExecutionReverted, t)
		tx, err = nft.SetBridge(auth, auth.From)
		assert.NoError(t, err)
		sim.Commit() // block
		CheckReceipt(sim, tx, time.Second, types.ReceiptStatusErrExecutionReverted, t)
	}
}

// DeployBridgeTest is a test-only function which deploys a bridge contract with some amount of KAIA.
func (bm *BridgeManager) DeployBridgeTest(backend *backends.SimulatedBackend, amountOfDeposit int64, local bool) (common.Address, error) {
	var acc *accountInfo

	// When the pending block of backend is updated, commit it
	// bm.DeployBridge will be waiting until the block is committed
	pendingBlock := backend.PendingBlock()
	go func() {
		for pendingBlock == backend.PendingBlock() {
			time.Sleep(100 * time.Millisecond)
		}
		backend.Commit()
		return
	}()

	// Set transfer value of the bridge account
	if local {
		acc = bm.subBridge.bridgeAccounts.cAccount
	} else {
		acc = bm.subBridge.bridgeAccounts.pAccount
	}

	auth := acc.GenerateTransactOpts()
	auth.Value = big.NewInt(amountOfDeposit)

	// Deploy a bridge contract
	deployedBridge, addr, err := bm.DeployBridge(auth, backend, local)
	if err != nil {
		return common.Address{}, err
	}

	// Set the bridge contract information to the BridgeManager
	err = bm.SetBridgeInfo(addr, deployedBridge, common.Address{}, nil, acc, local, false)
	if err != nil {
		return common.Address{}, err
	}
	return addr, err
}

// deployBridge deploys bridge contract and returns its address
func deployBridge(t *testing.T, bm *BridgeManager, backend *backends.SimulatedBackend, local bool) common.Address {
	var acc *accountInfo

	// When the pending block of backend is updated, commit it
	// bm.DeployBridge will be waiting until the block is committed
	pendingBlock := backend.PendingBlock()
	go func() {
		for pendingBlock == backend.PendingBlock() {
			time.Sleep(100 * time.Millisecond)
		}
		backend.Commit()
		return
	}()

	// Set transfer value of the bridge account
	if local {
		acc = bm.subBridge.bridgeAccounts.cAccount
	} else {
		acc = bm.subBridge.bridgeAccounts.pAccount
	}

	auth := acc.GenerateTransactOpts()
	auth.Value = big.NewInt(10000)

	// Deploy a bridge contract
	_, addr, err := bm.DeployBridge(auth, backend, local)
	assert.NoError(t, err)
	return addr
}

func isExpectedBalance(t *testing.T, bridgeManager *BridgeManager,
	pBridgeAddr, cBridgeAddr common.Address,
	expectedParentBridgeBalance, expectedChildBridgeBalance int64,
) {
	pBridgeBalance, err := bridgeManager.subBridge.APIBackend.GetParentBridgeContractBalance(pBridgeAddr)
	assert.NoError(t, err)
	cBridgeBalance, err := bridgeManager.subBridge.APIBackend.GetChildBridgeContractBalance(cBridgeAddr)
	assert.NoError(t, err)
	assert.Equal(t, pBridgeBalance.Int64(), expectedParentBridgeBalance)
	assert.Equal(t, cBridgeBalance.Int64(), expectedChildBridgeBalance)
}

func TestGetBridgeContractBalance(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bacc, _ := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	bacc.pAccount.chainID = big.NewInt(0)
	bacc.cAccount.chainID = big.NewInt(0)

	// Create Simulated backend
	alloc := blockchain.GenesisAlloc{
		bacc.pAccount.address: {Balance: big.NewInt(params.KAIA)},
		bacc.cAccount.address: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}
	sc.APIBackend = &SubBridgeAPI{sc}
	sc.handler, err = NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize bridgeHandler : %v", err)
		return
	}

	bm, err := NewBridgeManager(sc)
	assert.NoError(t, err)
	sc.handler.subbridge.bridgeManager = bm

	// Case 1 - Success
	{
		initialChildbridgeBalance, initialParentbridgeBalance := int64(100), int64(100)
		cBridgeAddr, err := bm.DeployBridgeTest(sim, initialChildbridgeBalance, true)
		assert.NoError(t, err)
		pBridgeAddr, err := bm.DeployBridgeTest(sim, initialParentbridgeBalance, false)
		assert.NoError(t, err)
		bm.SetJournal("", cBridgeAddr, pBridgeAddr)
		assert.NoError(t, err)
		sim.Commit()
		isExpectedBalance(t, bm, pBridgeAddr, cBridgeAddr, initialParentbridgeBalance, initialChildbridgeBalance)
	}

	// Case 2 - ? (Random)
	{
		for i := 0; i < 10; i++ {
			initialChildbridgeBalance, initialParentbridgeBalance := rand.Int63n(10000), rand.Int63n(10000)
			cBridgeAddr, err := bm.DeployBridgeTest(sim, initialChildbridgeBalance, true)
			assert.NoError(t, err)
			pBridgeAddr, err := bm.DeployBridgeTest(sim, initialParentbridgeBalance, false)
			assert.NoError(t, err)
			bm.SetJournal("", cBridgeAddr, pBridgeAddr)
			assert.NoError(t, err)
			sim.Commit()
			isExpectedBalance(t, bm, pBridgeAddr, cBridgeAddr, initialParentbridgeBalance, initialChildbridgeBalance)
		}
	}
}
