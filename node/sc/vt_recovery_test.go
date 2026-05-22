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
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	sctoken "github.com/kaiachain/kaia/contracts/bindings/testing/sc_erc20"
	scnft "github.com/kaiachain/kaia/contracts/bindings/testing/sc_erc721"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

type testInfo struct {
	t                 *testing.T
	sim               *backends.SimulatedBackend
	sc                *SubBridge
	bm                *BridgeManager
	localInfo         *BridgeInfo
	remoteInfo        *BridgeInfo
	tokenLocalAddr    common.Address
	tokenRemoteAddr   common.Address
	tokenLocalBridge  *sctoken.ServiceChainToken
	tokenRemoteBridge *sctoken.ServiceChainToken
	nftLocalAddr      common.Address
	nftRemoteAddr     common.Address
	nftLocalBridge    *scnft.ServiceChainNFT
	nftRemoteBridge   *scnft.ServiceChainNFT
	nodeAuth          *bind.TransactOpts
	chainAuth         *bind.TransactOpts
	aliceAuth         *bind.TransactOpts
	recoveryCh        chan bool
	mu                sync.Mutex
	nftIndex          int64
}

const (
	testGasLimit     = 1000000
	testAmount       = 321
	testToken        = 123
	testNFT          = int64(7321)
	testChargeToken  = 100000
	testTimeout      = 10 * time.Second
	testTxCount      = 7
	testBlockOffset  = 3 // +2 for genesis and bridge contract, +1 by a hardcoded hint
	testPendingCount = 3
)

type operations struct {
	request     func(*testInfo, *BridgeInfo)
	handle      func(*testInfo, *BridgeInfo, IRequestValueTransferEvent)
	dummyHandle func(*testInfo, *BridgeInfo)
}

type recoveryFixtureOptions struct {
	withToken bool
	withNFT   bool
	mintNFT   bool
}

type recoveryDirection uint8

const (
	childToParent recoveryDirection = iota
	parentToChild
)

type simpleRecoveryCase struct {
	name          string
	tokenType     uint8
	direction     recoveryDirection
	requestBridge func(*testInfo) *BridgeInfo
	handleBridge  func(*testInfo) *BridgeInfo
}

var ops = map[uint8]*operations{
	KAIA: {
		request:     requestKLAYTransfer,
		handle:      handleKLAYTransfer,
		dummyHandle: dummyHandleRequestKLAYTransfer,
	},
	ERC20: {
		request:     requestTokenTransfer,
		handle:      handleTokenTransfer,
		dummyHandle: dummyHandleRequestTokenTransfer,
	},
	ERC721: {
		request:     requestNFTTransfer,
		handle:      handleNFTTransfer,
		dummyHandle: dummyHandleRequestNFTTransfer,
	},
}

func prepareRecoveryWithOptions(t *testing.T, vtcallback func(*testInfo), fixtureOpts recoveryFixtureOptions) *testInfo {
	t.Helper()
	info := prepare(t, vtcallback, fixtureOpts)
	t.Cleanup(func() {
		info.sim.Close()
	})
	return info
}

// TestBasicKLAYTransferRecovery tests each methods of the value transfer recovery.
func TestBasicKLAYTransferRecovery(t *testing.T) {
	// 1. Init dummy chain and do some value transfers.
	info := prepareRecoveryWithOptions(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KAIA].request(info, info.localInfo)
		}
	}, recoveryFixtureOptions{})
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	// 2. Update recovery hint.
	err := vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	// 3. Request events by using the hint.
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}

	// 4. Check pending events.
	t.Log("check pending tx", "len", len(vtr.childEvents))
	count := 0
	for _, ev := range vtr.childEvents {
		assert.Equal(t, info.nodeAuth.From, ev.GetFrom())
		assert.Equal(t, info.aliceAuth.From, ev.GetTo())
		assert.Equal(t, big.NewInt(testAmount), ev.GetValueOrTokenId())
		assert.Condition(t, func() bool {
			return uint64(testBlockOffset) <= ev.GetRaw().BlockNumber
		})
		count++
	}
	assert.Equal(t, testPendingCount, count)

	// 5. Recover pending events
	info.recoveryCh <- true
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	ops[KAIA].dummyHandle(info, info.remoteInfo)

	// 6. Check empty pending events.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, 0, len(vtr.childEvents))

	assert.Equal(t, nil, vtr.Recover()) // nothing to recover
}

// TestKLAYTransferLongRangeRecovery tests a long block range recovery.
func TestKLAYTransferLongRangeRecovery(t *testing.T) {
	oldMaxPendingTxs := maxPendingTxs
	oldFilterLogsStride := filterLogsStride
	maxPendingTxs = 2
	filterLogsStride = 20
	defer func() {
		maxPendingTxs = oldMaxPendingTxs
		filterLogsStride = oldFilterLogsStride
	}()

	// 1. Init dummy chain and do some value transfers.
	info := prepareRecoveryWithOptions(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KAIA].request(info, info.localInfo)
			for range filterLogsStride {
				info.sim.Commit()
			}
		}
	}, recoveryFixtureOptions{})

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	steps := []struct {
		doRecover    bool
		expectedHint uint64
	}{
		{doRecover: false, expectedHint: uint64(testTxCount - testPendingCount)},
		{doRecover: true, expectedHint: uint64(testTxCount - testPendingCount + maxPendingTxs)},
		{doRecover: true, expectedHint: uint64(testTxCount)},
	}

	for i, step := range steps {
		if step.doRecover {
			if i == 1 {
				info.recoveryCh <- true
			}
			err := vtr.Recover()
			if err != nil {
				t.Fatal("fail to recover the value transfer")
			}
		}

		deadline := time.Now().Add(testTimeout)
		for {
			info.sim.Commit()
			if err := vtr.updateRecoveryHint(); err != nil {
				t.Fatal("fail to update value transfer hint")
			}
			if vtr.child2parentHint.requestNonce == uint64(testTxCount) &&
				vtr.child2parentHint.handleNonce == step.expectedHint {
				break
			}
			if time.Now().After(deadline) {
				t.Fatal("timeout waiting for child2parent hint update")
			}
			runtime.Gosched()
		}
		assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
		assert.Equal(t, step.expectedHint, vtr.child2parentHint.handleNonce)
	}
}

func TestRecoveryScenarios(t *testing.T) {
	info := prepareRecoveryWithOptions(t, nil, recoveryFixtureOptions{
		withToken: true,
		withNFT:   true,
		mintNFT:   true,
	})

	testCases := []simpleRecoveryCase{
		{
			name:          "basic_token_transfer_recovery",
			tokenType:     ERC20,
			direction:     childToParent,
			requestBridge: func(info *testInfo) *BridgeInfo { return info.localInfo },
			handleBridge:  func(info *testInfo) *BridgeInfo { return info.remoteInfo },
		},
		{
			name:          "basic_nft_transfer_recovery",
			tokenType:     ERC721,
			direction:     childToParent,
			requestBridge: func(info *testInfo) *BridgeInfo { return info.localInfo },
			handleBridge:  func(info *testInfo) *BridgeInfo { return info.remoteInfo },
		},
		{
			name:          "method_recover",
			tokenType:     KAIA,
			direction:     childToParent,
			requestBridge: func(info *testInfo) *BridgeInfo { return info.localInfo },
			handleBridge:  func(info *testInfo) *BridgeInfo { return info.remoteInfo },
		},
		{
			name:          "scenario_main_chain_recovery",
			tokenType:     KAIA,
			direction:     parentToChild,
			requestBridge: func(info *testInfo) *BridgeInfo { return info.remoteInfo },
			handleBridge:  func(info *testInfo) *BridgeInfo { return info.localInfo },
		},
		{
			name:          "scenario_automatic_recovery",
			tokenType:     KAIA,
			direction:     childToParent,
			requestBridge: func(info *testInfo) *BridgeInfo { return info.localInfo },
			handleBridge:  func(info *testInfo) *BridgeInfo { return info.remoteInfo },
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("%02d_%s", idx, tc.name), func(t *testing.T) {
			requestBridge := tc.requestBridge(info)
			for i := 0; i < testTxCount; i++ {
				ops[tc.tokenType].request(info, requestBridge)
			}

			vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)
			err := vtr.updateRecoveryHint()
			assert.NoError(t, err)

			if tc.direction == parentToChild {
				assert.NotEqual(t, vtr.parent2childHint.requestNonce, vtr.parent2childHint.handleNonce)
			} else {
				assert.NotEqual(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
			}

			info.recoveryCh <- true
			err = vtr.Recover()
			assert.NoError(t, err)
			ops[tc.tokenType].dummyHandle(info, tc.handleBridge(info))

			err = vtr.updateRecoveryHint()
			assert.NoError(t, err)
			if tc.direction == parentToChild {
				assert.Equal(t, vtr.parent2childHint.requestNonce, vtr.parent2childHint.handleNonce)
				return
			}
			assert.Equal(t, vtr.child2parentHint.requestNonce, vtr.child2parentHint.handleNonce)
		})
	}
}

// TestMethodStop tests the Stop method for stop the internal goroutine.
func TestMethodStop(t *testing.T) {
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, nil, nil)
	assert.NoError(t, vtr.Stop())

	vtr.isRunning = true
	assert.NoError(t, vtr.Stop())
	assert.Equal(t, false, vtr.isRunning)
}

// TestFlagVTRecovery tests the disabled vtrecovery option.
func TestFlagVTRecovery(t *testing.T) {
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: false}, nil, nil)
	err := vtr.Start()
	assert.Equal(t, ErrVtrDisabled, err)
}

// TestAlreadyStartedVTRecovery tests the already started VTR error cases.
func TestAlreadyStartedVTRecovery(t *testing.T) {
	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true, VTRecoveryInterval: 60}, nil, nil)
	vtr.isRunning = true
	err := vtr.Start()
	assert.Equal(t, ErrVtrAlreadyStarted, err)
}

// TestMultiOperatorRequestRecovery tests value transfer recovery for the multi-operator.
func TestMultiOperatorRequestRecovery(t *testing.T) {
	// 1. Init dummy chain and do some value transfers.
	info := prepareRecoveryWithOptions(t, func(info *testInfo) {
		for i := 0; i < testTxCount; i++ {
			ops[KAIA].request(info, info.localInfo)
		}
	}, recoveryFixtureOptions{})

	// 2. Set multi-operator.
	cAcc := info.nodeAuth
	pAcc := info.chainAuth
	opts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
	_, err := info.localInfo.bridge.RegisterOperator(opts, pAcc.From)
	assert.NoError(t, err)
	opts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	_, err = info.remoteInfo.bridge.RegisterOperator(opts, cAcc.From)
	assert.NoError(t, err)
	info.sim.Commit()

	// 3. Set operator threshold.
	opts = &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
	_, err = info.localInfo.bridge.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	opts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
	_, err = info.remoteInfo.bridge.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()

	vtr := NewValueTransferRecovery(&SCConfig{VTRecovery: true}, info.localInfo, info.remoteInfo)

	// 4. Update recovery hint.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	// 5. Request events by using the hint.
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}

	// 6. Check pending events.
	t.Log("check pending tx", "len", len(vtr.childEvents))
	count := 0
	for _, ev := range vtr.childEvents {
		assert.Equal(t, info.nodeAuth.From, ev.GetFrom())
		assert.Equal(t, info.aliceAuth.From, ev.GetTo())
		assert.Equal(t, big.NewInt(testAmount), ev.GetValueOrTokenId())
		assert.Condition(t, func() bool {
			return uint64(testBlockOffset) <= ev.GetRaw().BlockNumber
		})
		count++
	}
	assert.Equal(t, testPendingCount, count)

	// 7. Recover pending events
	info.recoveryCh <- true
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	ops[KAIA].dummyHandle(info, info.remoteInfo)

	// 8. Recover from the other operator (value transfer is not recovered yet).
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	t.Log("value transfer hint", vtr.child2parentHint)
	assert.Equal(t, uint64(testTxCount), vtr.child2parentHint.requestNonce)
	assert.Equal(t, uint64(testTxCount-testPendingCount), vtr.child2parentHint.handleNonce)

	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, testPendingCount, len(vtr.childEvents))
	assert.Equal(t, nil, vtr.recoverPendingEvents())
	info.remoteInfo.account = info.localInfo.account // other operator
	ops[KAIA].dummyHandle(info, info.remoteInfo)
	deadline := time.Now().Add(testTimeout)
	for {
		info.sim.Commit()
		if err := vtr.updateRecoveryHint(); err != nil {
			t.Fatal("fail to update value transfer hint")
		}
		if err := vtr.retrievePendingEvents(); err != nil {
			t.Fatal("fail to retrieve pending events from the bridge contract")
		}
		if len(vtr.childEvents) == 0 {
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timeout waiting for recovered pending requests to be drained")
		}
		runtime.Gosched()
	}

	// 9. Check results.
	err = vtr.updateRecoveryHint()
	if err != nil {
		t.Fatal("fail to update value transfer hint")
	}
	err = vtr.retrievePendingEvents()
	if err != nil {
		t.Fatal("fail to retrieve pending events from the bridge contract")
	}
	assert.Equal(t, 0, len(vtr.childEvents))
	assert.Equal(t, nil, vtr.Recover()) // nothing to recover
}

func assertTxMined(t *testing.T, sim *backends.SimulatedBackend, tx *types.Transaction) {
	t.Helper()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	assert.NoError(t, err)
	if receipt == nil {
		t.Fatal("receipt not found")
	}

	result := blockchain.ExecutionResult{
		VmExecutionStatus: receipt.Status,
	}
	assert.NoError(t, result.Unwrap())
}

// prepare generates dummy blocks for testing value transfer recovery.
func prepare(t *testing.T, vtcallback func(*testInfo), fixtureOpts recoveryFixtureOptions) *testInfo {
	// Setup configuration.
	config := &SCConfig{}
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()
	config.DataDir = tempDir
	config.VTRecovery = true

	bacc, _, _ := newBridgeAccountsForTest(t, nil, config.DataDir)

	cAcc := bacc.cAccount.GenerateTransactOpts()
	pAcc := bacc.pAccount.GenerateTransactOpts()

	// Generate a new random account and a funded simulator.
	aliceKey, _ := crypto.GenerateKey()
	aliceAuth := bind.NewKeyedTransactor(aliceKey)

	// Alloc genesis and create a simulator.
	alloc := blockchain.GenesisAlloc{
		cAcc.From:      {Balance: big.NewInt(params.KAIA)},
		pAcc.From:      {Balance: big.NewInt(params.KAIA)},
		aliceAuth.From: {Balance: big.NewInt(params.KAIA)},
	}
	sim := backends.NewSimulatedBackend(alloc)

	sc := &SubBridge{
		chainDB:        database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}),
		config:         config,
		peers:          newBridgePeerSet(),
		bridgeAccounts: bacc,
		localBackend:   sim,
		remoteBackend:  sim,
	}
	handler, err := NewSubBridgeHandler(sc)
	if err != nil {
		log.Fatalf("Failed to initialize the bridgeHandler : %v", err)
		return nil
	}
	sc.handler = handler
	sc.blockchain = sim.BlockChain()

	// Prepare manager and deploy bridge contract.
	bm, err := NewBridgeManager(sc)
	localAddr, err := bm.DeployBridgeTest(sim, 10000, true)
	assert.NoError(t, err)
	remoteAddr, err := bm.DeployBridgeTest(sim, 10000, false)
	assert.NoError(t, err)

	localInfo, _ := bm.GetBridgeInfo(localAddr)
	remoteInfo, _ := bm.GetBridgeInfo(remoteAddr)
	sim.Commit()

	var (
		tokenLocalAddr  common.Address
		tokenRemoteAddr common.Address
		tokenLocal      *sctoken.ServiceChainToken
		tokenRemote     *sctoken.ServiceChainToken
		nftLocalAddr    common.Address
		nftRemoteAddr   common.Address
		nftLocal        *scnft.ServiceChainNFT
		nftRemote       *scnft.ServiceChainNFT
		tx              *types.Transaction
	)

	if fixtureOpts.withToken {
		var tokenLocalTx *types.Transaction
		var tokenRemoteTx *types.Transaction
		var tl *sctoken.ServiceChainToken
		var tr *sctoken.ServiceChainToken
		tokenLocalAddr, tokenLocalTx, tl, err = sctoken.DeployServiceChainToken(cAcc, sim, localAddr)
		if err != nil {
			log.Fatalf("Failed to DeployServiceChainToken: %v", err)
		}
		tokenRemoteAddr, tokenRemoteTx, tr, err = sctoken.DeployServiceChainToken(pAcc, sim, remoteAddr)
		if err != nil {
			log.Fatalf("Failed to DeployServiceChainToken: %v", err)
		}
		tokenLocal, tokenRemote = tl, tr
		sim.Commit()
		assertTxMined(t, sim, tokenLocalTx)
		assertTxMined(t, sim, tokenRemoteTx)

		charge := big.NewInt(testChargeToken)
		txOpts := &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
		tx, err = tokenRemote.Transfer(txOpts, remoteAddr, charge)
		if err != nil {
			log.Fatalf("Failed to Transfer for charging: %v", err)
		}
		sim.Commit()
		assertTxMined(t, sim, tx)
	}

	if fixtureOpts.withNFT {
		var nftLocalTx *types.Transaction
		var nftRemoteTx *types.Transaction
		var nl *scnft.ServiceChainNFT
		var nr *scnft.ServiceChainNFT
		nftLocalAddr, nftLocalTx, nl, err = scnft.DeployServiceChainNFT(cAcc, sim, localAddr)
		if err != nil {
			log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
		}
		nftRemoteAddr, nftRemoteTx, nr, err = scnft.DeployServiceChainNFT(pAcc, sim, remoteAddr)
		if err != nil {
			log.Fatalf("Failed to DeployServiceChainNFT: %v", err)
		}
		nftLocal, nftRemote = nl, nr
		sim.Commit()
		assertTxMined(t, sim, nftLocalTx)
		assertTxMined(t, sim, nftRemoteTx)
	}

	if fixtureOpts.withToken || fixtureOpts.withNFT {
		nodeOpts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
		chainOpts := &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
		registerTokenTxs := make([]*types.Transaction, 0, 4)
		if fixtureOpts.withToken {
			tx, err = localInfo.bridge.RegisterToken(nodeOpts, tokenLocalAddr, tokenRemoteAddr)
			assert.NoError(t, err)
			registerTokenTxs = append(registerTokenTxs, tx)
			tx, err = remoteInfo.bridge.RegisterToken(chainOpts, tokenRemoteAddr, tokenLocalAddr)
			assert.NoError(t, err)
			registerTokenTxs = append(registerTokenTxs, tx)
		}
		if fixtureOpts.withNFT {
			tx, err = localInfo.bridge.RegisterToken(nodeOpts, nftLocalAddr, nftRemoteAddr)
			assert.NoError(t, err)
			registerTokenTxs = append(registerTokenTxs, tx)
			tx, err = remoteInfo.bridge.RegisterToken(chainOpts, nftRemoteAddr, nftLocalAddr)
			assert.NoError(t, err)
			registerTokenTxs = append(registerTokenTxs, tx)
		}
		sim.Commit()
		for _, tx := range registerTokenTxs {
			assertTxMined(t, sim, tx)
		}
	}

	if fixtureOpts.withNFT && fixtureOpts.mintNFT {
		mintTxs := make([]*types.Transaction, 0, testTxCount*2)
		for i := 0; i < testTxCount; i++ {
			txOpts := &bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}
			tx, err = nftLocal.MintWithTokenURI(txOpts, cAcc.From, big.NewInt(testNFT+int64(i)), "testURI")
			assert.NoError(t, err)
			mintTxs = append(mintTxs, tx)

			txOpts = &bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}
			tx, err = nftRemote.MintWithTokenURI(txOpts, remoteAddr, big.NewInt(testNFT+int64(i)), "testURI")
			assert.NoError(t, err)
			mintTxs = append(mintTxs, tx)
		}
		sim.Commit()
		for _, tx := range mintTxs {
			assertTxMined(t, sim, tx)
		}
	}

	// Register the owner as a signer
	tx, err = localInfo.bridge.RegisterOperator(&bind.TransactOpts{From: cAcc.From, Signer: cAcc.Signer, GasLimit: testGasLimit}, cAcc.From)
	assert.NoError(t, err)
	tx, err = remoteInfo.bridge.RegisterOperator(&bind.TransactOpts{From: pAcc.From, Signer: pAcc.Signer, GasLimit: testGasLimit}, pAcc.From)
	assert.NoError(t, err)
	sim.Commit()

	// Subscribe events.
	err = bm.SubscribeEvent(localAddr)
	if err != nil {
		t.Fatalf("local bridge manager event subscription failed")
	}
	err = bm.SubscribeEvent(remoteAddr)
	if err != nil {
		t.Fatalf("remote bridge manager event subscription failed")
	}

	// Prepare channel for event handling.
	requestVTCh := make(chan RequestValueTransferEvent)
	requestVTencodedCh := make(chan RequestValueTransferEncodedEvent)
	handleVTCh := make(chan *HandleValueTransferEvent)
	recoveryCh := make(chan bool)
	bm.SubscribeReqVTev(requestVTCh)
	bm.SubscribeReqVTencodedEv(requestVTencodedCh)
	bm.SubscribeHandleVTev(handleVTCh)

	info := testInfo{
		t, sim, sc, bm, localInfo, remoteInfo,
		tokenLocalAddr, tokenRemoteAddr, tokenLocal, tokenRemote,
		nftLocalAddr, nftRemoteAddr, nftLocal, nftRemote,
		cAcc, pAcc, aliceAuth, recoveryCh,
		sync.Mutex{},
		testNFT,
	}

	// Start a event handling loop.
	wg := sync.WaitGroup{}
	wg.Add((2 * testTxCount) - testPendingCount)
	isRecovery := false
	reqHandler := func(ev IRequestValueTransferEvent) {
		t.Log("request value transfer", "nonce", ev.GetRequestNonce())
		if ev.GetRequestNonce() >= (testTxCount - testPendingCount) {
			t.Log("missing handle value transfer", "nonce", ev.GetRequestNonce())
		} else {
			switch ev.GetTokenType() {
			case KAIA, ERC20, ERC721:
				break
			default:
				t.Errorf("received ev.TokenType is unknown: %v", ev.GetTokenType())
				return
			}

			if ev.GetRaw().Address == info.localInfo.address {
				ops[ev.GetTokenType()].handle(&info, info.remoteInfo, ev)
			} else {
				ops[ev.GetTokenType()].handle(&info, info.localInfo, ev)
			}
		}
	}
	go func() {
		for {
			select {
			case ev := <-recoveryCh:
				isRecovery = ev
			case ev := <-requestVTCh:
				reqHandler(ev)
				if !isRecovery {
					wg.Done()
				}
			case ev := <-requestVTencodedCh:
				reqHandler(ev)
				if !isRecovery {
					wg.Done()
				}
			case ev := <-handleVTCh:
				t.Log("handle value transfer", "nonce", ev.HandleNonce)
				if !isRecovery {
					wg.Done()
				}
			}
		}
	}()

	// Request value transfer if requested.
	if vtcallback != nil {
		vtcallback(&info)
		WaitGroupWithTimeOut(&wg, testTimeout, t)
	}

	return &info
}

func requestKLAYTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	opts := bi.account.GenerateTransactOpts()
	opts.Value = big.NewInt(testAmount)
	tx, err := bi.bridge.RequestKLAYTransfer(opts, info.aliceAuth.From, big.NewInt(testAmount), nil)
	if err != nil {
		log.Fatalf("Failed to RequestKLAYTransfer: %v", err)
	}
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

func handleKLAYTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testAmount), ev.GetValueOrTokenId())
	opts := bi.account.GenerateTransactOpts()
	tx, err := bi.bridge.HandleKLAYTransfer(opts, ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), ev.GetValueOrTokenId(), ev.GetRequestNonce(), ev.GetRaw().BlockNumber, ev.GetExtraData())
	if err != nil {
		log.Fatalf("\tFailed to HandleKLAYTransfer: %v", err)
	}
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

// TODO-Kaia-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestKLAYTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleKLAYTransfer(info, bi, ev.(RequestValueTransferEvent))
	}
	info.sim.Commit()
}

func requestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var tx *types.Transaction

	testToken := big.NewInt(testToken)
	opts := bi.account.GenerateTransactOpts()

	var err error
	if bi.onChildChain {
		tx, err = info.tokenLocalBridge.RequestValueTransfer(opts, testToken, info.chainAuth.From, big.NewInt(0), nil)
	} else {
		tx, err = info.tokenRemoteBridge.RequestValueTransfer(opts, testToken, info.nodeAuth.From, big.NewInt(0), nil)
	}

	if err != nil {
		log.Fatalf("Failed to RequestValueTransfer for charging: %v", err)
	}
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

func handleTokenTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	assert.Equal(info.t, new(big.Int).SetUint64(testToken), ev.GetValueOrTokenId())
	tx, err := bi.bridge.HandleERC20Transfer(
		bi.account.GenerateTransactOpts(), ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), info.tokenRemoteAddr, ev.GetValueOrTokenId(), ev.GetRequestNonce(), ev.GetRaw().BlockNumber, ev.GetExtraData())
	if err != nil {
		log.Fatalf("Failed to HandleERC20Transfer: %v", err)
	}
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

// TODO-Kaia-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestTokenTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleTokenTransfer(info, bi, ev.(RequestValueTransferEvent))
	}
	info.sim.Commit()
}

func requestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var tx *types.Transaction

	opts := bi.account.GenerateTransactOpts()
	// TODO-Kaia need to separate child / parent chain nftIndex.
	nftIndex := new(big.Int).SetInt64(info.nftIndex)

	var err error
	if bi.onChildChain {
		tx, err = info.nftLocalBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From, nil)
	} else {
		tx, err = info.nftRemoteBridge.RequestValueTransfer(opts, nftIndex, info.aliceAuth.From, nil)
	}

	if err != nil {
		log.Fatalf("Failed to requestNFTTransfer for charging: %v", err)
	}
	info.nftIndex++
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

func handleNFTTransfer(info *testInfo, bi *BridgeInfo, ev IRequestValueTransferEvent) {
	bi.account.Lock()
	defer bi.account.UnLock()

	var nftAddr common.Address

	if bi.onChildChain {
		nftAddr = info.nftLocalAddr
	} else {
		nftAddr = info.nftRemoteAddr
	}
	uri := GetURI(ev)
	tx, err := bi.bridge.HandleERC721Transfer(
		bi.account.GenerateTransactOpts(),
		ev.GetRaw().TxHash, ev.GetFrom(), ev.GetTo(), nftAddr, ev.GetValueOrTokenId(),
		ev.GetRequestNonce(), ev.GetRaw().BlockNumber, uri, ev.GetExtraData())
	if err != nil {
		log.Fatalf("Failed to handleERC721Transfer: %v", err)
	}
	info.sim.Commit()
	assertTxMined(info.t, info.sim, tx)
}

// TODO-Kaia-ServiceChain: use ChildChainEventHandler
func dummyHandleRequestNFTTransfer(info *testInfo, bi *BridgeInfo) {
	for _, ev := range bi.GetPendingRequestEvents() {
		handleNFTTransfer(info, bi, ev)
	}
	info.sim.Commit()
}
