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
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a regular transaction for testing (not approve or swap)
func makeRegularTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int) *types.Transaction {
	// Add some regular method signature (e.g., "transfer")
	data := append([]byte{}, common.Hex2Bytes("a9059cbb")...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(common.HexToAddress("0x1234567890123456789012345678901234567890").Bytes()), 32)...)
	data = append(data, common.Hex2BytesFixed(hex.EncodeToString(amount.Bytes()), 32)...)
	return makeTx(t, privKey, nonce, to, big.NewInt(0), gasLimit, gasPrice, data)
}

// Helper function to create a contract deployment transaction for testing
func makeDeployTx(t *testing.T, privKey *ecdsa.PrivateKey, nonce uint64, gasLimit uint64, gasPrice *big.Int, data []byte) *types.Transaction {
	if privKey == nil {
		var err error
		privKey, err = crypto.GenerateKey()
		require.NoError(t, err)
	}

	signer := types.LatestSignerForChainID(big.NewInt(1))
	tx := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data)
	tx, err := types.SignTx(tx, signer, privKey)
	require.NoError(t, err)

	return tx
}

// TestGaslessAPIIsGaslessTx tests the GaslessAPI.IsGaslessTx method
func TestGaslessAPI_isGaslessTx(t *testing.T) {
	// Create a private key for testing
	privKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	sender := crypto.PubkeyToAddress(privKey.PublicKey)

	// Create a different private key for testing
	otherPrivKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	// Create addresses for testing
	tokenAddr := common.HexToAddress("0xabcd")
	routerAddr := common.HexToAddress("0x1234")
	differentTokenAddr := common.HexToAddress("0xdead1234dead1234dead1234dead1234dead1234")

	// Create a simulated backend for testing
	dbm := database.NewMemoryDBManager()
	alloc := testAllocStorage()
	backend := backends.NewSimulatedBackendWithDatabase(dbm, alloc, testChainConfig)

	// Setup test stateDB with a proper database
	stateDB, _ := backend.BlockChain().State()
	stateDB.SetNonce(sender, 0) // For approve tx

	// Create mock txpool
	txpool := &testTxPool{
		statedb: stateDB,
	}

	// Create and initialize GaslessModule
	gaslessModule := NewGaslessModule()
	nodeKey, _ := crypto.GenerateKey()
	err = gaslessModule.Init(&InitOpts{
		ChainConfig:   testChainConfig,
		GaslessConfig: testGaslessConfig,
		NodeKey:       nodeKey,
		Chain:         backend.BlockChain(),
		TxPool:        txpool,
	})
	require.NoError(t, err)

	// Override token and router maps for testing
	gaslessModule.allowedTokens = map[common.Address]bool{
		tokenAddr: true,
	}
	gaslessModule.swapRouter = routerAddr

	// Create GaslessAPI
	api := NewGaslessAPI(gaslessModule)

	// Create test transactions
	validApproveTx := makeApproveTx(t, privKey, 0, ApproveArgs{
		Token:   tokenAddr,
		Spender: routerAddr,
		Amount:  abi.MaxUint256,
	})
	// Create a standalone swap transaction for testing
	standaloneSwapTx := makeSwapTx(t, privKey, 0, SwapArgs{
		Router:       routerAddr,
		Token:        tokenAddr,
		AmountIn:     big.NewInt(500000),
		MinAmountOut: big.NewInt(100),
		AmountRepay:  big.NewInt(1021000),
		Deadline:     big.NewInt(300),
	})

	// Create a swap transaction for testing
	validSwapTx := makeSwapTx(t, privKey, 1, SwapArgs{
		Router:       routerAddr,
		Token:        tokenAddr,
		AmountIn:     big.NewInt(500000),
		MinAmountOut: big.NewInt(100),
		AmountRepay:  big.NewInt(2021000),
		Deadline:     big.NewInt(300),
	})
	invalidTokenSwapTx := makeSwapTx(t, privKey, 1, SwapArgs{
		Router:       routerAddr,
		Token:        differentTokenAddr,
		AmountIn:     big.NewInt(500000),
		MinAmountOut: big.NewInt(100),
		AmountRepay:  big.NewInt(1000000),
		Deadline:     big.NewInt(300),
	})
	differentSenderSwapTx := makeSwapTx(t, otherPrivKey, 0, SwapArgs{
		Router:       routerAddr,
		Token:        tokenAddr,
		AmountIn:     big.NewInt(500000),
		MinAmountOut: big.NewInt(100),
		AmountRepay:  big.NewInt(1000000),
		Deadline:     big.NewInt(300),
	})
	// Define gas parameters
	gasPrice := big.NewInt(1)
	approveGasLimit := uint64(1000000)
	swapGasLimit := uint64(1000000)

	regularTx := makeRegularTx(t, privKey, 0, tokenAddr, big.NewInt(1000), approveGasLimit, gasPrice)
	deployTx := makeDeployTx(t, privKey, 0, swapGasLimit, gasPrice, []byte{0x60, 0x80, 0x60, 0x40})

	// Test cases
	testCases := []struct {
		name           string
		txs            []*types.Transaction
		expectedResult bool
		reasonContains string
	}{
		{
			name:           "Valid standalone swap transaction",
			txs:            []*types.Transaction{standaloneSwapTx},
			expectedResult: true,
			reasonContains: "",
		},
		{
			name:           "Valid approve + swap pair",
			txs:            []*types.Transaction{validApproveTx, validSwapTx},
			expectedResult: true,
			reasonContains: "",
		},
		{
			name:           "Invalid - contract deployment",
			txs:            []*types.Transaction{deployTx},
			expectedResult: false,
			reasonContains: "transaction is not a swap transaction",
		},
		{
			name:           "Invalid - regular transaction",
			txs:            []*types.Transaction{regularTx},
			expectedResult: false,
			reasonContains: "transaction is not a swap transaction",
		},
		{
			name:           "Invalid - different token addresses in approve+swap pair",
			txs:            []*types.Transaction{validApproveTx, invalidTokenSwapTx},
			expectedResult: false,
			reasonContains: "second transaction is not a swap transaction",
		},
		{
			name:           "Invalid - different senders in approve+swap pair",
			txs:            []*types.Transaction{validApproveTx, differentSenderSwapTx},
			expectedResult: false,
			reasonContains: "approve and swap transactions have different senders",
		},
		{
			name:           "Invalid - too many transactions",
			txs:            []*types.Transaction{validApproveTx, validSwapTx, validSwapTx},
			expectedResult: false,
			reasonContains: "expected 1 or 2 transactions",
		},
		{
			name:           "Invalid - first transaction not approve in pair",
			txs:            []*types.Transaction{regularTx, validSwapTx},
			expectedResult: false,
			reasonContains: "first transaction is not an approve transaction",
		},
		{
			name:           "Invalid - second transaction not swap in pair",
			txs:            []*types.Transaction{validApproveTx, regularTx},
			expectedResult: false,
			reasonContains: "second transaction is not a swap transaction",
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert transactions to raw bytes
			rawTxs := make([]hexutil.Bytes, len(tc.txs))
			for i, tx := range tc.txs {
				data, err := rlp.EncodeToBytes(tx)
				require.NoError(t, err)
				rawTxs[i] = data
			}

			// Call IsGaslessTx
			result := api.IsGaslessTx(context.Background(), rawTxs)

			// Print debug information
			var reasonStr string
			if result.Reason != nil {
				reasonStr = *result.Reason
				t.Logf("IsGasless: %v, Reason: %s", result.IsGasless, reasonStr)
			} else {
				t.Logf("IsGasless: %v, Reason: nil", result.IsGasless)
			}

			// Check result
			assert.Equal(t, tc.expectedResult, result.IsGasless, "IsGasless flag should match expected result")
			if tc.reasonContains != "" {
				assert.NotNil(t, result.Reason, "Reason should not be nil for invalid transactions")
				assert.Contains(t, *result.Reason, tc.reasonContains, "Reason should contain expected message")
			} else if !tc.expectedResult {
				assert.NotNil(t, result.Reason, "Reason should not be nil for invalid transactions")
				assert.NotEmpty(t, *result.Reason, "Reason should not be empty for invalid transactions")
			} else {
				assert.Nil(t, result.Reason, "Reason should be nil for valid transactions")
			}
		})
	}
}
