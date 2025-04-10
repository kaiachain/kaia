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
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	mock_builder "github.com/kaiachain/kaia/kaiax/builder/mock"
	mock_kaiax "github.com/kaiachain/kaia/kaiax/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func init() {
	blockchain.InitDeriveSha(params.TestChainConfig)
}

func TestPreAddTx_KnownTxTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		tx            *types.Transaction
		knownTxs      map[common.Hash]txAndTime
		expectedError error
	}{
		{
			name:          "Transaction not in knownTxs",
			tx:            createTestTransaction(0),
			knownTxs:      make(map[common.Hash]txAndTime),
			expectedError: nil,
		},
		{
			name: "Transaction during KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now(),
				},
			},
			expectedError: errors.New("Unable to add known bundle tx into tx pool during lock period"),
		},
		{
			name: "Transaction after KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now().Add(-KnownTxTimeout),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(true).AnyTimes()

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         tt.knownTxs,
			}

			err := builderModule.PreAddTx(tt.tx, true)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestPreAddTx_TxPoolModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	preAddTxError := errors.New("tx pool error")

	tests := []struct {
		name            string
		hasTxPoolModule bool
		preAddTxResult  error
		expectedError   error
	}{
		{
			name:            "TxPoolModule is not set",
			hasTxPoolModule: false,
			preAddTxResult:  nil,
			expectedError:   nil,
		},
		{
			name:            "TxPoolModule is not set and PreAddTx returns error",
			hasTxPoolModule: false,
			preAddTxResult:  preAddTxError,
			expectedError:   nil,
		},
		{
			name:            "TxPoolModule is set",
			hasTxPoolModule: true,
			preAddTxResult:  nil,
			expectedError:   nil,
		},
		{
			name:            "TxPoolModule is set but PreAddTx returns error",
			hasTxPoolModule: true,
			preAddTxResult:  preAddTxError,
			expectedError:   preAddTxError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTx := createTestTransaction(0)
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]txAndTime),
			}

			mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)

			if tt.hasTxPoolModule {
				mockTxPoolModule.EXPECT().PreAddTx(testTx, true).Return(tt.preAddTxResult)
				builderModule.txPoolModule = mockTxPoolModule
			}

			err := builderModule.PreAddTx(testTx, true)
			assert.ErrorIs(t, err, tt.expectedError)
		})
	}
}

func TestIsModuleTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		hasTxPoolModule bool
		isBundleTx      bool
		isModuleTx      bool
		expectedResult  bool
	}{
		{
			name:            "TxPoolModule is not set and IsBundleTx returns true",
			hasTxPoolModule: false,
			isBundleTx:      true,
			isModuleTx:      false,
			expectedResult:  true,
		},
		{
			name:            "TxPoolModule is not set and IsBundleTx returns false",
			hasTxPoolModule: false,
			isBundleTx:      false,
			isModuleTx:      false,
			expectedResult:  false,
		},
		{
			name:            "TxPoolModule is set and IsModuleTx returns true",
			hasTxPoolModule: true,
			isBundleTx:      false,
			isModuleTx:      true,
			expectedResult:  true,
		},
		{
			name:            "TxPoolModule is set and IsModuleTx returns false",
			hasTxPoolModule: true,
			isBundleTx:      false,
			isModuleTx:      false,
			expectedResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testTx := createTestTransaction(0)
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(testTx).Return(tt.isBundleTx).AnyTimes()

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]txAndTime),
			}

			if tt.hasTxPoolModule {
				mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
				mockTxPoolModule.EXPECT().IsModuleTx(testTx).Return(tt.isModuleTx)
				builderModule.txPoolModule = mockTxPoolModule
			}

			result := builderModule.IsModuleTx(testTx)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGetCheckBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		hasTxPoolModule bool
		checkBalance    func(tx *types.Transaction) error
		expectedResult  func(tx *types.Transaction) error
	}{
		{
			name:            "TxPoolModule is not set",
			hasTxPoolModule: false,
			checkBalance:    nil,
			expectedResult:  nil,
		},
		{
			name:            "TxPoolModule is not set with check balance function",
			hasTxPoolModule: false,
			checkBalance: func(tx *types.Transaction) error {
				return errors.New("balance check error")
			},
			expectedResult: nil,
		},
		{
			name:            "TxPoolModule is set with check balance function",
			hasTxPoolModule: true,
			checkBalance: func(tx *types.Transaction) error {
				return nil
			},
			expectedResult: func(tx *types.Transaction) error {
				return nil
			},
		},
		{
			name:            "TxPoolModule is set with error check balance function",
			hasTxPoolModule: true,
			checkBalance: func(tx *types.Transaction) error {
				return errors.New("balance check error")
			},
			expectedResult: func(tx *types.Transaction) error {
				return errors.New("balance check error")
			},
		},
		{
			name:            "TxPoolModule is set without check balance function",
			hasTxPoolModule: true,
			checkBalance:    nil,
			expectedResult:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]txAndTime),
			}

			if tt.hasTxPoolModule {
				mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
				mockTxPoolModule.EXPECT().GetCheckBalance().Return(tt.checkBalance)
				builderModule.txPoolModule = mockTxPoolModule
			}

			result := builderModule.GetCheckBalance()
			if tt.expectedResult == nil {
				assert.Nil(t, result)
				return
			}

			// Test the returned function with a sample transaction
			testTx := createTestTransaction(0)
			expectedErr := tt.expectedResult(testTx)
			actualErr := result(testTx)
			if expectedErr == nil {
				assert.NoError(t, actualErr)
			} else {
				assert.EqualError(t, actualErr, expectedErr.Error())
			}
		})
	}
}

func TestIsReady_KnownTxs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	oldTime := now.Add(-time.Hour)

	tests := []struct {
		name             string
		txs              map[uint64]*types.Transaction
		isBundleTxResult bool
		knownTxs         map[common.Hash]txAndTime
		expectedTxs      map[common.Hash]txAndTime
	}{
		{
			name:             "No txs transactions",
			txs:              make(map[uint64]*types.Transaction),
			isBundleTxResult: false,
			knownTxs:         make(map[common.Hash]txAndTime),
			expectedTxs:      make(map[common.Hash]txAndTime),
		},
		{
			name:             "Non-bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: false,
			knownTxs:         map[common.Hash]txAndTime{},
			expectedTxs:      map[common.Hash]txAndTime{},
		},
		{
			name:             "New bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs:         map[common.Hash]txAndTime{},
			expectedTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
			},
		},
		{
			name:             "Existing bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: oldTime,
				},
			},
			expectedTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: oldTime,
				},
			},
		},
		{
			name:             "KnownTx has txs",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(4)},
			isBundleTxResult: true,
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(3).Hash(): {
					tx:   createTestTransaction(3),
					time: oldTime,
				},
			},
			expectedTxs: map[common.Hash]txAndTime{
				createTestTransaction(3).Hash(): {
					tx:   createTestTransaction(3),
					time: oldTime,
				},
				createTestTransaction(4).Hash(): {
					tx: createTestTransaction(4),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(gomock.Any()).Return(tt.isBundleTxResult).AnyTimes()

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyTxAndTimeMap(tt.knownTxs),
			}

			builderModule.IsReady(tt.txs, 0, nil)

			// Verify knownTxs state
			assert.Equal(t, len(tt.expectedTxs), len(builderModule.knownTxs))
			for hash, expected := range tt.expectedTxs {
				actual, exists := builderModule.knownTxs[hash]
				assert.True(t, exists)
				assert.Equal(t, expected.tx.Hash(), actual.tx.Hash())

				// For new transactions, verify the time is recent
				if _, exists := tt.knownTxs[hash]; !exists {
					assert.True(t, time.Since(actual.time) < time.Second, "New transaction time should be recent")
				} else {
					// For existing transactions, verify the time is preserved
					assert.Equal(t, expected.time, actual.time, "Existing transaction time should be preserved")
				}
			}
		})
	}
}

func TestIsReady_TxPoolModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		hasTxPoolModule bool
		isReadyResult   bool
		expectedResult  bool
	}{
		{
			name:            "TxPoolModule is not set",
			hasTxPoolModule: false,
			isReadyResult:   false,
			expectedResult:  true,
		},
		{
			name:            "TxPoolModule is set and returns true",
			hasTxPoolModule: true,
			isReadyResult:   true,
			expectedResult:  true,
		},
		{
			name:            "TxPoolModule is set and returns false",
			hasTxPoolModule: true,
			isReadyResult:   false,
			expectedResult:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]txAndTime),
			}

			// Create a test transaction map
			testTxs := map[uint64]*types.Transaction{
				0: createTestTransaction(0),
			}

			// Set up mock expectations
			mockTxBundlingModule.EXPECT().IsBundleTx(gomock.Any()).Return(false).AnyTimes()

			if tt.hasTxPoolModule {
				mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
				mockTxPoolModule.EXPECT().IsReady(testTxs, uint64(0), nil).Return(tt.isReadyResult)
				builderModule.txPoolModule = mockTxPoolModule
			}

			result := builderModule.IsReady(testTxs, 0, nil)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestPreReset_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()

	tests := []struct {
		name                 string
		existingBundles      map[common.Hash]txAndTime
		expectedBundles      map[common.Hash]txAndTime
		expectedUnexecutable bool
	}{
		{
			name: "Bundle tx within PendingTimeout period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			expectedUnexecutable: false,
		},
		{
			name: "Bundle tx exceeds PendingTimeout period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-PendingTimeout),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-PendingTimeout),
				},
			},
			expectedUnexecutable: true,
		},
		{
			name: "Bundle tx exceeds KnownTxTimeout period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-KnownTxTimeout),
				},
			},
			expectedBundles:      map[common.Hash]txAndTime{},
			expectedUnexecutable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			// Setup BuilderModule
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyTxAndTimeMap(tt.existingBundles),
			}

			builderModule.PreReset(nil, nil)

			// Verify the number of bundles
			assert.Equal(t, len(tt.expectedBundles), len(builderModule.knownTxs))

			// Verify each bundle's state
			for hash, expected := range tt.expectedBundles {
				actual, exists := builderModule.knownTxs[hash]
				assert.True(t, exists)
				assert.Equal(t, expected.tx.Hash(), actual.tx.Hash())
				assert.Equal(t, expected.time, actual.time)
				assert.Equal(t, tt.expectedUnexecutable, actual.tx.IsMarkedUnexecutable())
			}
		})
	}
}

func TestPreReset_TxPoolModule(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		hasTxPoolModule bool
	}{
		{
			name:            "TxPoolModule is not set",
			hasTxPoolModule: false,
		},
		{
			name:            "TxPoolModule is set",
			hasTxPoolModule: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]txAndTime),
			}

			if tt.hasTxPoolModule {
				mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
				mockTxPoolModule.EXPECT().PreReset(nil, nil).Do(func(txs *types.TransactionsByPriceAndNonce, next uint64, ready types.Transactions) {
					panic("tx pool is called")
				})
				builderModule.txPoolModule = mockTxPoolModule
			}

			if tt.hasTxPoolModule {
				assert.Panics(t, func() {
					builderModule.PreReset(nil, nil)
				})
			} else {
				assert.NotPanics(t, func() {
					builderModule.PreReset(nil, nil)
				})
			}
		})
	}
}

func createTestTransaction(nonce uint64) *types.Transaction {
	return types.NewTransaction(
		nonce,
		common.HexToAddress("0x1"),
		common.Big0,
		21000,
		common.Big0,
		nil,
	)
}

func copyTxAndTimeMap(m map[common.Hash]txAndTime) map[common.Hash]txAndTime {
	newMap := make(map[common.Hash]txAndTime)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
