// Modifications Copyright 2025 The Kaia Authors
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

package tests

import (
	"context"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	uniswapFactoryContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/factory"
	uniswapRouterContracts "github.com/kaiachain/kaia/contracts/contracts/libs/uniswap/router"
	kip149contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip149"
	gaslessContract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip247"
	testingContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	testingGaslessContracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts/gasless"
	gaslessImpl "github.com/kaiachain/kaia/kaiax/gasless/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	bigKaia = big.NewInt(params.KAIA)
	bigGkei = big.NewInt(params.Gkei)
)

func TestGasless(t *testing.T) {
	log.EnableLogForTest(log.LvlError, log.LvlError)

	// prepare chain configuration
	config := params.MainnetChainConfig.Copy()
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = big.NewInt(0)
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.PragueCompatibleBlock = big.NewInt(0)

	config.Istanbul.SubGroupSize = 1
	config.Istanbul.ProposerPolicy = uint64(istanbul.RoundRobin)

	fullNode, node, validator, _, workspace := newBlockchain(t, config, nil)
	defer func() {
		os.RemoveAll(workspace)
		if err := fullNode.Stop(); err != nil {
			t.Fatal(err)
		}
	}()

	numAccounts := 2
	_, accounts, _ := createAccount(t, numAccounts, validator)

	var (
		owner      = validator
		transactor = backends.NewBlockchainContractBackend(node.BlockChain(), node.TxPool().(*blockchain.TxPool), nil)
		chain      = node.BlockChain().(*blockchain.BlockChain)
	)

	/* ------------- Deploy contracts ------------- */
	testTokenAddr, testTokenContract := deployTestToken(t, chain, transactor, owner, owner.Addr)
	wkaiaAddr, wkaiaContract := deployWKAIA(t, chain, transactor, owner)
	factoryAddr, factoryContract := deployUniswapV2Factory(t, chain, transactor, owner, owner.Addr)
	routerAddr, routerContract := deployUniswapV2Router02(t, chain, transactor, owner, factoryAddr, wkaiaAddr)
	gsrAddr, gsrContract := deployGaslessSwapRouter(t, chain, transactor, owner, wkaiaAddr)

	/* ------------- Register GaslessSwapRouter address in Registry ------------- */
	// send register tx
	targetBlockNum := new(big.Int).Add(node.BlockChain().CurrentHeader().Number, big.NewInt(10))
	registry, err := kip149contract.NewRegistry(system.RegistryAddr, transactor)
	if err != nil {
		t.Fatal(err)
	}
	registerTx, err := registry.Register(bind.NewKeyedTransactor(owner.Keys[0]), gaslessImpl.GaslessSwapRouterName, gsrAddr, targetBlockNum)
	if err != nil {
		t.Fatal(err)
	}
	registerTxReceipt := waitReceipt(chain, registerTx.Hash())
	if registerTxReceipt == nil || registerTxReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to registor GaslessSwapRouter address")
	}

	// wait target block
	targetHeader := waitBlock(chain, targetBlockNum.Uint64())
	require.NotNil(t, targetHeader)

	/* ------------- Set up initial liquidity ------------- */
	contracts := contractsForGasless{
		testTokenAddr:     testTokenAddr,
		testTokenContract: testTokenContract,
		wkaiaAddr:         wkaiaAddr,
		wkaiaContract:     wkaiaContract,
		factoryAddr:       factoryAddr,
		factoryContract:   factoryContract,
		routerAddr:        routerAddr,
		routerContract:    routerContract,
		gsrAddr:           gsrAddr,
		gsrContract:       gsrContract,
	}
	setupLiquidity(t, owner, contracts, chain)

	/* ------------------------------------ Main test process ------------------------------------- */
	// In the test below, we swap 1 Token -> `amountsOut` WKAIA.
	swapAmmount := new(big.Int).Mul(big.NewInt(1), bigKaia)
	amountsOut, err := routerContract.GetAmountsOut(&bind.CallOpts{}, swapAmmount, []common.Address{testTokenAddr, wkaiaAddr})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("amountsOut: %s", amountsOut)

	var (
		gasPriceBN         = new(big.Int).Mul(big.NewInt(50), bigGkei)
		R1                 = new(big.Int).Mul(big.NewInt(21000), gasPriceBN)
		R2                 = new(big.Int).Mul(big.NewInt(100000), gasPriceBN)
		R3                 = new(big.Int).Mul(big.NewInt(500000), gasPriceBN)
		ammontRepay        = new(big.Int).Add(R1, new(big.Int).Add(R2, R3))
		amountRepaySwap    = new(big.Int).Add(R1, R3)
		transferToken      = new(big.Int).Mul(big.NewInt(100), bigKaia)
		swapExpectedOutput = amountsOut[1]
		margin             = new(big.Int).Div(swapExpectedOutput, big.NewInt(100))
		minAmountOut       = new(big.Int).Add(ammontRepay, margin)
		deadline           = new(big.Int).Add(chain.CurrentBlock().Time(), big.NewInt(300))
	)

	//// Reject when token balance is zero, approval is zero.

	// reject approveTx when token balance is zero
	_, err = sendApproveTx(t, testTokenContract, accounts[1], gsrAddr, abi.MaxUint256)
	assert.ErrorContains(t, err, "insufficient sender token balance")
	assert.ErrorContains(t, err, "have=0, want=nonzero")

	// reject swapTx when token isn't yet approved
	_, err = sendSwapTx(t, gsrContract, accounts[1], testTokenAddr, swapAmmount, minAmountOut, ammontRepay, deadline)
	assert.ErrorContains(t, err, "insufficient approval: approval=0")

	//// After having token balance, ApproveTx + SwapTx pair succeeds.

	// transfer test token
	testTokenTransferTx, err := testTokenContract.Transfer(bind.NewKeyedTransactor(owner.Keys[0]), accounts[0].Addr, transferToken)
	if err != nil {
		t.Fatal(err)
	}
	testTokenTransferReceipt := waitReceipt(chain, testTokenTransferTx.Hash())
	if testTokenTransferReceipt == nil {
		t.Fatal("failed to transfer test token")
	}
	owner.Nonce += 1

	preSwapBalanceOfTestAcc, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, accounts[0].Addr)
	preSwapBalanceOfOwner, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	t.Log("test acc balance: ", preSwapBalanceOfTestAcc)
	t.Log("owner balance: ", preSwapBalanceOfOwner)

	// success send normal approveTx
	approveTx, err := sendApproveTx(t, testTokenContract, accounts[0], gsrAddr, abi.MaxUint256)
	if err != nil {
		t.Fatal(err)
	}
	accounts[0].Nonce += 1

	// success send normal swapTx
	swapTx, err := sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, swapAmmount, minAmountOut, ammontRepay, deadline)
	if err != nil {
		t.Fatal(err)
	}
	accounts[0].Nonce += 1

	// check if account[0] without kaia can send tx
	approveTxReceipt := waitReceipt(chain, approveTx.Hash())
	require.NotNil(t, approveTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, approveTxReceipt.Status, "approveTx failed")

	swapTxReceipt := waitReceipt(chain, swapTx.Hash())
	require.NotNil(t, swapTxReceipt)
	require.Equal(t, types.ReceiptStatusSuccessful, swapTxReceipt.Status, "swapTx failed")

	//// Verify the effect of ApproveTx + SwapTx pair.

	// verify test acc's kaia balances
	// since gasPrice may be less than the theoretical value, we check that the current balance is greater than or equal to FinalUserAmount.
	// expected: (pre balance) = 0
	// expected: (current balance) >= FinalUserAmount
	_, _, gaslessBlockNum, _ := chain.GetTxAndLookupInfo(swapTx.Hash())
	preState, _, err := node.APIBackend.StateAndHeaderByNumber(context.Background(), rpc.BlockNumber(gaslessBlockNum-1))
	if err != nil {
		t.Fatal(err)
	}
	currentState, _, err := node.APIBackend.StateAndHeaderByNumber(context.Background(), rpc.BlockNumber(gaslessBlockNum))
	if err != nil {
		t.Fatal(err)
	}
	swappedForGasEvent, err := gsrContract.ParseSwappedForGas(*swapTxReceipt.Logs[len(swapTxReceipt.Logs)-1]) // SwappedForGas is issued at the end of swapForGas
	if err != nil {
		t.Fatal(err)
	}
	require.True(t, preState.GetBalance(accounts[0].Addr).Cmp(common.Big0) == 0)
	require.True(t, currentState.GetBalance(accounts[0].Addr).Cmp(swappedForGasEvent.FinalUserAmount) != -1)

	// verify test token balances
	// expected: (current balance) = (pre balance) - swapAmmount
	currentBalanceOfTestAcc, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, accounts[0].Addr)
	require.True(t, currentBalanceOfTestAcc.Cmp(new(big.Int).Sub(preSwapBalanceOfTestAcc, swapAmmount)) == 0)
	t.Logf("final token balance is correct: %v", currentBalanceOfTestAcc.Cmp(new(big.Int).Sub(preSwapBalanceOfTestAcc, swapAmmount)) == 0)

	//// Reject obviously reverting SwapTx.

	// reject swapTx when minAmountOut < amountRepay
	_, err = sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, swapAmmount, common.Big0, amountRepaySwap, deadline)
	assert.ErrorContains(t, err, "insufficient minAmountOut")

	// reject swapTx when amountIn < router.GetAmountIn(minAmountOut)
	_, err = sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, common.Big0, minAmountOut, amountRepaySwap, deadline)
	assert.ErrorContains(t, err, "insufficient amountIn")

	// reject swapTx when balance < amountIn
	// the test acc first received `transferToken` but used up some. So it has less than `transferToken`.
	_, err = sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, transferToken, minAmountOut, amountRepaySwap, deadline)
	assert.ErrorContains(t, err, "insufficient balance")

	// reject swapTx when deadline is in the past
	_, err = sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, swapAmmount, minAmountOut, amountRepaySwap, common.Big1)
	assert.ErrorContains(t, err, "insufficient deadline: deadline=1")

	// reject swapTx originating from an EOA with code
	sendSetCodeTx(t, chain, transactor, accounts[0])
	_, err = sendSwapTx(t, gsrContract, accounts[0], testTokenAddr, swapAmmount, minAmountOut, amountRepaySwap, deadline)
	assert.ErrorContains(t, err, "sender with code is not allowed")
}

func deployTestToken(t *testing.T, chain *blockchain.BlockChain, transactor *backends.BlockchainContractBackend, owner *TestAccountType, initialHolder common.Address,
) (common.Address, *testingGaslessContracts.TestToken) {
	addr, tx, contract, err := testingGaslessContracts.DeployTestToken(bind.NewKeyedTransactor(owner.Keys[0]), transactor, initialHolder)
	if err != nil {
		t.Fatal(err)
	}

	receipt := waitReceipt(chain, tx.Hash())
	if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deploy TestToken")
	}

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("TestToken deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return addr, contract
}

func deployWKAIA(t *testing.T, chain *blockchain.BlockChain, transactor *backends.BlockchainContractBackend, owner *TestAccountType,
) (common.Address, *testingContracts.WKAIA) {
	addr, tx, contract, err := testingContracts.DeployWKAIA(bind.NewKeyedTransactor(owner.Keys[0]), transactor)
	if err != nil {
		t.Fatal(err)
	}

	receipt := waitReceipt(chain, tx.Hash())
	if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deploy WKAIA")
	}

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("WKAIA deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return addr, contract
}

func deployUniswapV2Factory(t *testing.T, chain *blockchain.BlockChain, transactor *backends.BlockchainContractBackend, owner *TestAccountType, feeToSetter common.Address,
) (common.Address, *uniswapFactoryContracts.UniswapV2Factory) {
	addr, tx, contract, err := uniswapFactoryContracts.DeployUniswapV2Factory(bind.NewKeyedTransactor(owner.Keys[0]), transactor, feeToSetter)
	if err != nil {
		t.Fatal(err)
	}

	receipt := waitReceipt(chain, tx.Hash())
	if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deploy UniswapV2Factory")
	}

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("UniswapV2Factory deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return addr, contract
}

func deployUniswapV2Router02(t *testing.T, chain *blockchain.BlockChain, transactor *backends.BlockchainContractBackend, owner *TestAccountType, factory, wkaia common.Address,
) (common.Address, *uniswapRouterContracts.UniswapV2Router02) {
	addr, tx, contract, err := uniswapRouterContracts.DeployUniswapV2Router02(bind.NewKeyedTransactor(owner.Keys[0]), transactor, factory, wkaia)
	if err != nil {
		t.Fatal(err)
	}

	receipt := waitReceipt(chain, tx.Hash())
	if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deploy UniswapV2Router02")
	}

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("UniswapV2Router02 deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return addr, contract
}

func deployGaslessSwapRouter(t *testing.T, chain *blockchain.BlockChain, transactor *backends.BlockchainContractBackend, owner *TestAccountType, wkaiaAddr common.Address,
) (common.Address, *gaslessContract.GaslessSwapRouter) {
	addr, tx, contract, err := gaslessContract.DeployGaslessSwapRouter(bind.NewKeyedTransactor(owner.Keys[0]), transactor, wkaiaAddr)
	if err != nil {
		t.Fatal(err)
	}

	receipt := waitReceipt(chain, tx.Hash())
	if receipt == nil || receipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deploy GaslessSwapRouter")
	}

	_, _, num, _ := chain.GetTxAndLookupInfo(tx.Hash())
	t.Logf("GaslessSwapRouter deployed at block=%2d, addr=%s", num, addr.Hex())

	owner.Nonce++
	return addr, contract
}

type contractsForGasless struct {
	testTokenAddr     common.Address
	testTokenContract *testingGaslessContracts.TestToken
	wkaiaAddr         common.Address
	wkaiaContract     *testingContracts.WKAIA
	factoryAddr       common.Address
	factoryContract   *uniswapFactoryContracts.UniswapV2Factory
	routerAddr        common.Address
	routerContract    *uniswapRouterContracts.UniswapV2Router02
	gsrAddr           common.Address
	gsrContract       *gaslessContract.GaslessSwapRouter
}

func setupLiquidity(t *testing.T, owner *TestAccountType, contracts contractsForGasless, chain *blockchain.BlockChain) {
	var (
		testTokenAddr     = contracts.testTokenAddr
		testTokenContract = contracts.testTokenContract
		wkaiaAddr         = contracts.wkaiaAddr
		wkaiaContract     = contracts.wkaiaContract
		factoryAddr       = contracts.factoryAddr
		factoryContract   = contracts.factoryContract
		routerAddr        = contracts.routerAddr
		routerContract    = contracts.routerContract
		gsrContract       = contracts.gsrContract
		initialLiquidity  = new(big.Int).Mul(big.NewInt(1000), bigKaia)
	)

	/* ------------- create pair ------------- */
	createPairTx, err := factoryContract.CreatePair(bind.NewKeyedTransactor(owner.Keys[0]), testTokenAddr, wkaiaAddr)
	if err != nil {
		t.Fatal(err)
	}
	createPairReceipt := waitReceipt(chain, createPairTx.Hash())
	if createPairReceipt == nil || createPairReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to create pair")
	}
	owner.Nonce += 1

	/* ------------- deposit ------------- */
	optsForDeposit := bind.NewKeyedTransactor(owner.Keys[0])
	optsForDeposit.Value = initialLiquidity
	optsForDeposit.GasLimit = 300000
	depositTx, err := wkaiaContract.Deposit(optsForDeposit)
	if err != nil {
		t.Fatal(err)
	}
	depositReceipt := waitReceipt(chain, depositTx.Hash())
	if depositReceipt == nil || depositReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to deposit")
	}
	owner.Nonce += 1

	/* ------------- approve(TestToken) ------------- */
	testTokenApproveTx, err := testTokenContract.Approve(bind.NewKeyedTransactor(owner.Keys[0]), routerAddr, initialLiquidity)
	if err != nil {
		t.Fatal(err)
	}
	testTokenApproveReceipt := waitReceipt(chain, testTokenApproveTx.Hash())
	if testTokenApproveReceipt == nil || testTokenApproveReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to approve(TestToken)")
	}
	owner.Nonce += 1

	/* ------------- approve(WKAIA) ------------- */
	wkaiaApproveTx, err := wkaiaContract.Approve(bind.NewKeyedTransactor(owner.Keys[0]), routerAddr, initialLiquidity)
	if err != nil {
		t.Fatal(err)
	}
	wkaiaApproveReceipt := waitReceipt(chain, wkaiaApproveTx.Hash())
	if wkaiaApproveReceipt == nil || wkaiaApproveReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to approve(WKAIA)")
	}
	owner.Nonce += 1

	balanceOfWKAIA, _ := wkaiaContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	balanceOfTestToken, _ := testTokenContract.BalanceOf(&bind.CallOpts{}, owner.Addr)
	wallowance, _ := wkaiaContract.Allowance(&bind.CallOpts{}, owner.Addr, routerAddr)
	tallowance, _ := testTokenContract.Allowance(&bind.CallOpts{}, owner.Addr, routerAddr)
	t.Log("balance of tokens: ", balanceOfWKAIA, balanceOfTestToken)
	t.Log("allowances of tokens: ", wallowance, tallowance)

	/* ------------- add liquidity ------------- */
	optsForAddLiquidity := bind.NewKeyedTransactor(owner.Keys[0])
	optsForAddLiquidity.GasLimit = 3000000
	deadline := time.Now().Unix() + 60*20
	addLiquidityTx, err := routerContract.AddLiquidity(optsForAddLiquidity, testTokenAddr, wkaiaAddr,
		initialLiquidity, initialLiquidity, common.Big0, common.Big0, owner.Addr, big.NewInt(deadline))
	if err != nil {
		t.Fatal(err)
	}
	addLiquidityReceipt := waitReceipt(chain, addLiquidityTx.Hash())
	if addLiquidityReceipt == nil || addLiquidityReceipt.Status != types.ReceiptStatusSuccessful {
		t.Log(addLiquidityReceipt)
		t.Fatal("failed to add liquidity")
	}
	owner.Nonce += 1

	/* ------------- add token to gsr ------------- */
	optsForAddToken := bind.NewKeyedTransactor(owner.Keys[0])
	optsForAddToken.GasLimit = 300000
	addTokenTx, err := gsrContract.AddToken(optsForAddToken, testTokenAddr, factoryAddr, routerAddr)
	if err != nil {
		t.Fatal(err)
	}
	addTokenReceipt := waitReceipt(chain, addTokenTx.Hash())
	if addTokenReceipt == nil || addTokenReceipt.Status != types.ReceiptStatusSuccessful {
		t.Fatal("failed to add token to gsr")
	}
	owner.Nonce += 1
}

func sendApproveTx(t *testing.T, testTokenContract *testingGaslessContracts.TestToken, owner *TestAccountType, gsrAddr common.Address, amount *big.Int) (*types.Transaction, error) {
	optsForApprove := bind.NewKeyedTransactor(owner.Keys[0])
	optsForApprove.GasLimit = 300000
	optsForApprove.Nonce = big.NewInt(int64(owner.Nonce))
	approveTx, err := testTokenContract.Approve(optsForApprove, gsrAddr, amount)
	if err != nil {
		return nil, err
	}
	t.Log("approveTxHash", approveTx.Hash().Hex())
	return approveTx, nil
}

func sendSwapTx(t *testing.T, gsrContract *gaslessContract.GaslessSwapRouter, owner *TestAccountType, testTokenAddr common.Address, swapAmmount *big.Int, minAmountOut *big.Int, ammontRepay *big.Int, deadline *big.Int) (*types.Transaction, error) {
	optsForSwap := bind.NewKeyedTransactor(owner.Keys[0])
	optsForSwap.GasLimit = 300000
	optsForSwap.Nonce = big.NewInt(int64(owner.Nonce))
	swapTx, err := gsrContract.SwapForGas(optsForSwap, testTokenAddr, swapAmmount, minAmountOut, ammontRepay, deadline)
	if err != nil {
		return nil, err
	}
	t.Log("swapTxHash", swapTx.Hash().Hex())
	return swapTx, nil
}

func sendSetCodeTx(t *testing.T, chain *blockchain.BlockChain, transactor bind.ContractBackend, sender *TestAccountType) {
	chainID, err := transactor.ChainID(context.Background())
	require.NoError(t, err)
	auth, err := types.SignSetCode(sender.Keys[0], types.SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(chainID),
		Address: common.HexToAddress("0x000000000000000000000000000000000000aaaa"),
		Nonce:   sender.GetNonce() + 1,
	})
	require.NoError(t, err)

	valueMap := map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:             sender.GetNonce(),
		types.TxValueKeyTo:                common.HexToAddress("0x0000000000000000000000000000000000001111"),
		types.TxValueKeyAmount:            common.Big0,
		types.TxValueKeyData:              []byte{},
		types.TxValueKeyGasLimit:          uint64(100000),
		types.TxValueKeyGasFeeCap:         big.NewInt(50000000000),
		types.TxValueKeyGasTipCap:         big.NewInt(50000000000),
		types.TxValueKeyAccessList:        types.AccessList{},
		types.TxValueKeyAuthorizationList: []types.SetCodeAuthorization{auth},
		types.TxValueKeyChainID:           chainID,
	}
	tx, err := types.NewTransactionWithMap(types.TxTypeEthereumSetCode, valueMap)
	require.NoError(t, err)
	signer := types.LatestSignerForChainID(chainID)
	err = tx.SignWithKeys(signer, sender.Keys)
	require.NoError(t, err)

	err = transactor.SendTransaction(context.Background(), tx)
	require.NoError(t, err)

	receipt := waitReceipt(chain, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	code, err := transactor.CodeAt(context.Background(), sender.Addr, nil)
	require.NoError(t, err)
	sender.Nonce += 2
	t.Logf("setcode complete: code=0x%x", code)
}
