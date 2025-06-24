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
	"math"
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

//go:generate mockgen -source=../tx_pool.go -destination=../mock/tx_pool.go -package=mock_kaiax TxPool

func init() {
	blockchain.InitDeriveSha(params.TestChainConfig)
}

func TestPreAddTx_KnownTxTimeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		tx            *types.Transaction
		knownTxs      map[common.Hash]knownTx
		expectedError error
	}{
		{
			name:          "Transaction not in knownTxs",
			tx:            createTestTransaction(0),
			knownTxs:      make(map[common.Hash]knownTx),
			expectedError: nil,
		},
		{
			name: "Transaction during KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now(),
				},
			},
			expectedError: ErrUnableToAddKnownBundleTx,
		},
		{
			name: "Transaction after KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]knownTx{
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
				knownTxs:         make(map[common.Hash]knownTx),
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
				knownTxs:         make(map[common.Hash]knownTx),
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
				knownTxs:         make(map[common.Hash]knownTx),
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
		knownTxs         map[common.Hash]knownTx
		expectedTxs      map[common.Hash]knownTx
	}{
		{
			name:             "No txs transactions",
			txs:              make(map[uint64]*types.Transaction),
			isBundleTxResult: false,
			knownTxs:         make(map[common.Hash]knownTx),
			expectedTxs:      make(map[common.Hash]knownTx),
		},
		{
			name:             "Non-bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: false,
			knownTxs:         map[common.Hash]knownTx{},
			expectedTxs:      map[common.Hash]knownTx{},
		},
		{
			name:             "New bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs:         map[common.Hash]knownTx{},
			expectedTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
			},
		},
		{
			name:             "Existing bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: oldTime,
				},
			},
			expectedTxs: map[common.Hash]knownTx{
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
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(3).Hash(): {
					tx:   createTestTransaction(3),
					time: oldTime,
				},
			},
			expectedTxs: map[common.Hash]knownTx{
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
			mockTxBundlingModule.EXPECT().GetMaxBundleTxsInPending().Return(uint(math.MaxUint64)).AnyTimes()

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyknownTxMap(tt.knownTxs),
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

func TestIsReady_MaxBundleTxs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()
	unexecutable := now.Add(time.Hour)
	testTx := createTestTransaction(900)

	// we assume that tx.Nonce() < 1000 is a bundle tx
	tests := []struct {
		name           string
		maxBundleTxs   uint
		knownTxs       map[common.Hash]knownTx
		additionalTxs  map[uint64]*types.Transaction
		readyTxs       []*types.Transaction
		expectedResult bool
	}{
		{
			name:           "Max bundle txs is zero",
			maxBundleTxs:   0,
			knownTxs:       make(map[common.Hash]knownTx),
			expectedResult: false,
		},
		{
			name:         "Below max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now(),
				},
			},
			expectedResult: true,
		},
		{
			name:         "At max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: now,
				},
			},
			expectedResult: false,
		},
		{
			name:         "At max bundle txs limit with ready bundle txs",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: now,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(1)},
			expectedResult: true,
		},
		{
			name:         "At max bundle txs limit with ready non-bundle txs",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: now,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(1000)},
			expectedResult: false,
		},
		{
			name:         "Above max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: now,
				},
				createTestTransaction(2).Hash(): {
					tx:   createTestTransaction(2),
					time: now,
				},
			},
			expectedResult: false,
		},
		{
			name:         "Above max bundle txs limit with ready txs",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: now,
				},
				createTestTransaction(2).Hash(): {
					tx:   createTestTransaction(2),
					time: now,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(2)},
			expectedResult: true,
		},
		{
			name:         "KnownTxs with unexecutable transactions",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:   createTestTransaction(1),
					time: unexecutable,
				},
			},
			expectedResult: true,
		},
		{
			name:         "KnownTxs with demoted transactions",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				createTestTransaction(1).Hash(): {
					tx:        createTestTransaction(1),
					time:      now,
					isDemoted: true,
				},
			},
			expectedResult: true,
		},
		{
			name:         "Additional bundle tx within limit",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{902: createTestTransaction(902)},
			expectedResult: true,
		},
		{
			name:         "Additional non-bundle tx within limit",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{901: createTestTransaction(1000)},
			expectedResult: true,
		},
		{
			name:         "Multiple additional txs with non-bundle tx",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{899: createTestTransaction(899), 901: createTestTransaction(1000)},
			expectedResult: true,
		},
		{
			name:         "Known transaction in knownTxs",
			maxBundleTxs: 2,
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
				testTx.Hash(): {
					tx:   testTx,
					time: now,
				},
			},
			expectedResult: true,
		},
		{
			name:           "No max bundle txs limit",
			maxBundleTxs:   math.MaxUint64,
			knownTxs:       make(map[common.Hash]knownTx),
			additionalTxs:  map[uint64]*types.Transaction{901: createTestTransaction(901)},
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockIsBundleTx := mockTxBundlingModule.EXPECT().IsBundleTx(gomock.Any()).DoAndReturn(func(tx *types.Transaction) bool {
				return !(tx == nil || tx.Nonce() >= 1000)
			}).Times(1)
			if tt.maxBundleTxs != math.MaxUint64 {
				mockIsBundleTx.AnyTimes()
			}
			mockTxBundlingModule.EXPECT().GetMaxBundleTxsInPending().Return(tt.maxBundleTxs).AnyTimes()

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyknownTxMap(tt.knownTxs),
			}

			// Mark some transactions as unexecutable to test that they are not counted
			for _, knownTx := range tt.knownTxs {
				if knownTx.time.After(now) {
					knownTx.tx.MarkUnexecutable(true)
				}
			}

			txs := map[uint64]*types.Transaction{
				900: testTx,
			}
			for nonce, tx := range tt.additionalTxs {
				txs[nonce] = tx
			}
			result := builderModule.IsReady(txs, 900, tt.readyTxs)
			assert.Equal(t, tt.expectedResult, result)
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
				knownTxs:         make(map[common.Hash]knownTx),
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
		existingBundles      map[common.Hash]knownTx
		expectedBundles      map[common.Hash]knownTx
		expectedUnexecutable bool
	}{
		{
			name: "Bundle tx within PendingTimeout period",
			existingBundles: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			expectedBundles: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			expectedUnexecutable: false,
		},
		{
			name: "Bundle tx exceeds PendingTimeout period",
			existingBundles: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-PendingTimeout),
				},
			},
			expectedBundles: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-PendingTimeout),
				},
			},
			expectedUnexecutable: true,
		},
		{
			name: "Bundle tx exceeds KnownTxTimeout period",
			existingBundles: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-KnownTxTimeout),
				},
			},
			expectedBundles:      map[common.Hash]knownTx{},
			expectedUnexecutable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			// Setup BuilderModule
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyknownTxMap(tt.existingBundles),
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
				knownTxs:         make(map[common.Hash]knownTx),
			}

			mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
			if tt.hasTxPoolModule {
				mockTxPoolModule.EXPECT().PreReset(nil, nil).Return().Times(1)
				builderModule.txPoolModule = mockTxPoolModule
			} else {
				mockTxPoolModule.EXPECT().PreReset(nil, nil).Return().Times(0)
			}

			builderModule.PreReset(nil, nil)
		})
	}
}

func TestPostReset_TxDemotion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name            string
		knownTxs        map[common.Hash]knownTx
		pendingTxs      map[common.Address]types.Transactions
		expectedDemoted map[common.Hash]bool
	}{
		{
			name:            "No known transactions",
			knownTxs:        make(map[common.Hash]knownTx),
			pendingTxs:      make(map[common.Address]types.Transactions),
			expectedDemoted: make(map[common.Hash]bool),
		},
		{
			name: "All known transactions in pending",
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(1),
				},
			},
			expectedDemoted: make(map[common.Hash]bool),
		},
		{
			name: "Some known transactions not in pending",
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
				createTestTransaction(2).Hash(): {
					tx: createTestTransaction(2),
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(2),
				},
			},
			expectedDemoted: map[common.Hash]bool{
				createTestTransaction(1).Hash(): true,
			},
		},
		{
			name: "No known transactions in pending",
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(2),
					createTestTransaction(3),
				},
			},
			expectedDemoted: map[common.Hash]bool{
				createTestTransaction(0).Hash(): true,
				createTestTransaction(1).Hash(): true,
			},
		},
		{
			name: "Already demoted transactions",
			knownTxs: map[common.Hash]knownTx{
				createTestTransaction(0).Hash(): {
					tx:        createTestTransaction(0),
					isDemoted: true,
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(1),
				},
			},
			expectedDemoted: map[common.Hash]bool{
				createTestTransaction(0).Hash(): true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxPool := mock_kaiax.NewMockTxPoolForCaller(ctrl)

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         copyknownTxMap(tt.knownTxs),
				txPool:           mockTxPool,
			}
			builderModule.PostReset(nil, nil, tt.pendingTxs)

			// Verify that transactions are marked as demoted correctly
			for hash, knownTx := range builderModule.knownTxs {
				expectedDemoted := tt.expectedDemoted[hash]
				assert.Equal(t, expectedDemoted, knownTx.isDemoted, "Transaction %x demoted state mismatch", hash)
			}
		})
	}
}

func TestPostReset_TxPoolModule(t *testing.T) {
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
			mockTxPool := mock_kaiax.NewMockTxPoolForCaller(ctrl)

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         make(map[common.Hash]knownTx),
				txPool:           mockTxPool,
			}

			mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
			if tt.hasTxPoolModule {
				mockTxPoolModule.EXPECT().PostReset(nil, nil, gomock.Any()).Return().Times(1)
				builderModule.txPoolModule = mockTxPoolModule
			} else {
				mockTxPoolModule.EXPECT().PostReset(nil, nil, gomock.Any()).Return().Times(0)
			}

			builderModule.PostReset(nil, nil, make(map[common.Address]types.Transactions))
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

func copyknownTxMap(m map[common.Hash]knownTx) map[common.Hash]knownTx {
	newMap := make(map[common.Hash]knownTx)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
