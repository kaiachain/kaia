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
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	mock_kaiax "github.com/kaiachain/kaia/kaiax/mock"
	"github.com/kaiachain/kaia/params"
	mock_builder "github.com/kaiachain/kaia/work/builder/mock"
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
		knownTxs      *knownTxs
		expectedError error
	}{
		{
			name:          "Transaction not in knownTxs",
			tx:            createTestTransaction(0),
			knownTxs:      &knownTxs{},
			expectedError: nil,
		},
		{
			name: "Transaction in pending during KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    time.Now(),
					promotedTime: time.Now(),
					status:       TxStatusPending,
				},
			},
			expectedError: ErrUnableToAddKnownBundleTx,
		},
		{
			name: "Transaction in queue during KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    time.Now(),
					promotedTime: time.Now(),
					status:       TxStatusQueue,
				},
			},
			expectedError: ErrUnableToAddKnownBundleTx,
		},
		{
			name: "Transaction in pending after KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now().Add(-KnownTxTimeout),
					status:       TxStatusPending,
				},
			},
			expectedError: nil,
		},
		{
			name: "Transaction in queue after KnownTxTimeout period",
			tx:   createTestTransaction(0),
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now().Add(-KnownTxTimeout),
					status:       TxStatusQueue,
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(true).AnyTimes()
			mockTxBundlingModule.EXPECT().GetMaxBundleTxsInQueue().Return(uint(200)).AnyTimes()

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
			mockTxBundlingModule.EXPECT().IsBundleTx(testTx).Return(true).AnyTimes()
			mockTxBundlingModule.EXPECT().GetMaxBundleTxsInQueue().Return(uint(200)).AnyTimes()

			builderModule := NewBuilderWrappingModule(mockTxBundlingModule)

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

func TestPreAddTx_BundleTxQueueSizeLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name                    string
		tx                      *types.Transaction
		isBundleTx              bool
		maxBundleTxsInQueue     uint
		initialQueueCount       int
		expectedError           error
		expectedQueueCountAfter int
	}{
		{
			name:                    "Non-bundle transaction should not be added to queue",
			tx:                      createTestTransaction(0),
			isBundleTx:              false,
			maxBundleTxsInQueue:     10,
			initialQueueCount:       10,
			expectedError:           nil,
			expectedQueueCountAfter: 10,
		},
		{
			name:                    "Bundle transaction with empty queue should be added",
			tx:                      createTestTransaction(0),
			isBundleTx:              true,
			maxBundleTxsInQueue:     10,
			initialQueueCount:       0,
			expectedError:           nil,
			expectedQueueCountAfter: 1,
		},
		{
			name:                    "Bundle transaction with space in queue should be added",
			tx:                      createTestTransaction(0),
			isBundleTx:              true,
			maxBundleTxsInQueue:     10,
			initialQueueCount:       5,
			expectedError:           nil,
			expectedQueueCountAfter: 6,
		},
		{
			name:                    "Bundle transaction at queue limit should be added",
			tx:                      createTestTransaction(0),
			isBundleTx:              true,
			maxBundleTxsInQueue:     10,
			initialQueueCount:       10,
			expectedError:           ErrBundleTxQueueFull,
			expectedQueueCountAfter: 10,
		},
		{
			name:                    "Bundle transaction exceeding queue limit should return error",
			tx:                      createTestTransaction(0),
			isBundleTx:              true,
			maxBundleTxsInQueue:     10,
			initialQueueCount:       11,
			expectedError:           ErrBundleTxQueueFull,
			expectedQueueCountAfter: 11, // Should not change
		},
		{
			name:                    "Bundle transaction with zero queue limit should be added",
			tx:                      createTestTransaction(0),
			isBundleTx:              true,
			maxBundleTxsInQueue:     0,
			initialQueueCount:       0,
			expectedError:           ErrBundleTxQueueFull,
			expectedQueueCountAfter: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock bundling module
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(tt.isBundleTx).Times(1)

			// Only expect GetMaxBundleTxsInQueue if it's a bundle transaction
			if tt.isBundleTx {
				mockTxBundlingModule.EXPECT().GetMaxBundleTxsInQueue().Return(tt.maxBundleTxsInQueue).Times(1)
			}

			// Create initial knownTxs with queue items if needed
			knownTxs := &knownTxs{}
			for i := 0; i < tt.initialQueueCount; i++ {
				tx := createTestTransaction(uint64(i + 100)) // Use different nonce to avoid hash conflicts
				knownTxs.add(tx, TxStatusQueue)
			}

			// Create builder module
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         knownTxs,
			}

			// Call PreAddTx
			err := builderModule.PreAddTx(tt.tx, true)

			// Verify error
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			// Verify queue count
			actualQueueCount := builderModule.knownTxs.numQueue()
			assert.Equal(t, tt.expectedQueueCountAfter, actualQueueCount,
				"Queue count mismatch. Expected: %d, Actual: %d", tt.expectedQueueCountAfter, actualQueueCount)

			// If transaction was successfully added, verify it's in the queue
			if tt.expectedError == nil && tt.isBundleTx {
				knownTx, exists := builderModule.knownTxs.get(tt.tx.Hash())
				assert.True(t, exists, "Transaction should be in knownTxs")
				assert.Equal(t, TxStatusQueue, knownTx.status, "Transaction should have queue status")
				assert.Equal(t, tt.tx.Hash(), knownTx.tx.Hash(), "Transaction hash should match")
			}
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
				knownTxs:         new(knownTxs),
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
				knownTxs:         new(knownTxs),
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
		knownTxs         *knownTxs
		expectedTxs      *knownTxs
	}{
		{
			name:             "No txs transactions",
			txs:              make(map[uint64]*types.Transaction),
			isBundleTxResult: false,
			knownTxs:         new(knownTxs),
			expectedTxs:      new(knownTxs),
		},
		{
			name:             "Non-bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: false,
			knownTxs:         new(knownTxs),
			expectedTxs:      new(knownTxs),
		},
		{
			name:             "New bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs:         new(knownTxs),
			expectedTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
			},
		},
		{
			name:             "Existing bundle transaction",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(0)},
			isBundleTxResult: true,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    oldTime,
					promotedTime: oldTime,
				},
			},
			expectedTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    oldTime,
					promotedTime: oldTime,
				},
			},
		},
		{
			name:             "KnownTx has txs",
			txs:              map[uint64]*types.Transaction{0: createTestTransaction(4)},
			isBundleTxResult: true,
			knownTxs: &knownTxs{
				createTestTransaction(3).Hash(): {
					tx:           createTestTransaction(3),
					promotedTime: oldTime,
				},
			},
			expectedTxs: &knownTxs{
				createTestTransaction(3).Hash(): {
					tx:           createTestTransaction(3),
					promotedTime: oldTime,
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
				knownTxs:         tt.knownTxs.Copy(),
			}

			builderModule.IsReady(tt.txs, 0, nil)

			// Verify knownTxs state
			assert.Equal(t, len(*tt.expectedTxs), len(*builderModule.knownTxs))
			for hash, expected := range *tt.expectedTxs {
				actual, exists := (*builderModule.knownTxs)[hash]
				assert.True(t, exists)
				assert.Equal(t, expected.tx.Hash(), actual.tx.Hash())

				// For new transactions, verify the time is recent
				if _, exists := (*tt.knownTxs)[hash]; !exists {
					assert.True(t, time.Since(actual.promotedTime) < time.Second, "New transaction time should be recent")
				} else {
					// For existing transactions, verify the time is preserved
					assert.Equal(t, expected.addedTime, actual.addedTime, "Existing transaction's addedTime should be preserved")
					assert.Equal(t, expected.promotedTime, actual.promotedTime, "Existing transaction's promotedTime should be preserved")
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
		knownTxs       *knownTxs
		additionalTxs  map[uint64]*types.Transaction
		readyTxs       []*types.Transaction
		expectedResult bool
	}{
		{
			name:           "Max bundle txs is zero",
			maxBundleTxs:   0,
			knownTxs:       new(knownTxs),
			expectedResult: false,
		},
		{
			name:         "Below max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusPending,
				},
			},
			expectedResult: true,
		},
		{
			name:         "At max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			expectedResult: false,
		},
		{
			name:         "At max bundle txs limit with ready bundle txs",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(1)},
			expectedResult: true,
		},
		{
			name:         "At max bundle txs limit with ready non-bundle txs",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(1000)},
			expectedResult: false,
		},
		{
			name:         "Above max bundle txs limit",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(2).Hash(): {
					tx:           createTestTransaction(2),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			expectedResult: false,
		},
		{
			name:         "Above max bundle txs limit with ready txs",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(2).Hash(): {
					tx:           createTestTransaction(2),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			readyTxs:       []*types.Transaction{createTestTransaction(2)},
			expectedResult: true,
		},
		{
			name:         "KnownTxs with unexecutable transactions",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: unexecutable,
					status:       TxStatusPending,
				},
			},
			expectedResult: true,
		},
		{
			name:         "KnownTxs with queue and demoted transactions",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: now,
					status:       TxStatusQueue,
				},
				createTestTransaction(2).Hash(): {
					tx:           createTestTransaction(2),
					promotedTime: now,
					status:       TxStatusDemoted,
				},
			},
			expectedResult: true,
		},
		{
			name:         "Additional bundle tx within limit",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{902: createTestTransaction(902)},
			expectedResult: true,
		},
		{
			name:         "Additional non-bundle tx within limit",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{901: createTestTransaction(1000)},
			expectedResult: true,
		},
		{
			name:         "Multiple additional txs with non-bundle tx",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			additionalTxs:  map[uint64]*types.Transaction{899: createTestTransaction(899), 901: createTestTransaction(1000)},
			expectedResult: true,
		},
		{
			name:         "Known transaction in knownTxs",
			maxBundleTxs: 2,
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now,
					status:       TxStatusPending,
				},
				testTx.Hash(): {
					tx:           testTx,
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			expectedResult: false,
		},
		{
			name:           "No max bundle txs limit",
			maxBundleTxs:   math.MaxUint64,
			knownTxs:       new(knownTxs),
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
				knownTxs:         tt.knownTxs.Copy(),
			}

			// Mark some transactions as unexecutable to test that they are not counted
			for _, knownTx := range *tt.knownTxs {
				if knownTx.promotedTime.After(now) {
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
				knownTxs:         new(knownTxs),
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
		name            string
		existingBundles *knownTxs
		expectedBundles *knownTxs
		expectedDrop    []common.Hash
	}{
		{
			name: "Bundle tx within PendingTimeout period",
			existingBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			expectedBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now,
					status:       TxStatusPending,
				},
			},
			expectedDrop: []common.Hash{},
		},
		{
			name: "Bundle tx within QueueTimeout period",
			existingBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now,
					status:       TxStatusQueue,
				},
			},
			expectedBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now,
					status:       TxStatusQueue,
				},
			},
			expectedDrop: []common.Hash{},
		},
		{
			name: "Bundle tx exceeds PendingTimeout period",
			existingBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now.Add(-PendingTimeout),
					status:       TxStatusPending,
				},
			},
			expectedBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now,
					promotedTime: now.Add(-PendingTimeout),
					status:       TxStatusPending,
				},
			},
			expectedDrop: []common.Hash{createTestTransaction(0).Hash()},
		},
		{
			name: "Bundle tx exceeds QueueTimeout period",
			existingBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now.Add(-QueueTimeout),
					promotedTime: now,
					status:       TxStatusQueue,
				},
			},
			expectedBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					addedTime:    now.Add(-QueueTimeout),
					promotedTime: now,
					status:       TxStatusQueue,
				},
			},
			expectedDrop: []common.Hash{createTestTransaction(0).Hash()},
		},
		{
			name: "Bundle tx exceeds KnownTxTimeout period",
			existingBundles: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: now.Add(-KnownTxTimeout),
				},
			},
			expectedBundles: new(knownTxs),
			expectedDrop:    []common.Hash{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			// Setup BuilderModule
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         tt.existingBundles.Copy(),
			}

			drops := builderModule.PreReset(nil, nil)

			// Verify the number of bundles
			assert.Equal(t, len(*tt.expectedBundles), len(*builderModule.knownTxs))

			// Verify each bundle's state
			for hash, expected := range *tt.expectedBundles {
				actual, exists := (*builderModule.knownTxs)[hash]
				assert.True(t, exists)
				assert.Equal(t, expected.tx.Hash(), actual.tx.Hash())
				assert.Equal(t, expected.addedTime, actual.addedTime)
				assert.Equal(t, expected.promotedTime, actual.promotedTime)
				assert.Equal(t, tt.expectedDrop, drops)
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
			builderModule := NewBuilderWrappingModule(mockTxBundlingModule)

			mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
			if tt.hasTxPoolModule {
				mockTxPoolModule.EXPECT().PreReset(nil, nil).Return(nil).Times(1)
				builderModule.txPoolModule = mockTxPoolModule
			} else {
				mockTxPoolModule.EXPECT().PreReset(nil, nil).Return(nil).Times(0)
			}

			builderModule.PreReset(nil, nil)
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

			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         new(knownTxs),
			}

			mockTxPoolModule := mock_kaiax.NewMockTxPoolModule(ctrl)
			if tt.hasTxPoolModule {
				mockTxPoolModule.EXPECT().PostReset(nil, nil, gomock.Any(), gomock.Any()).Return().Times(1)
				builderModule.txPoolModule = mockTxPoolModule
			} else {
				mockTxPoolModule.EXPECT().PostReset(nil, nil, gomock.Any(), gomock.Any()).Return().Times(0)
			}

			builderModule.PostReset(nil, nil, make(map[common.Address]types.Transactions), make(map[common.Address]types.Transactions))
		})
	}
}

func TestPostReset_TransactionStatusUpdates(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name           string
		knownTxs       *knownTxs
		queue          map[common.Address]types.Transactions
		pending        map[common.Address]types.Transactions
		expectedStatus map[common.Hash]int // hash -> expected status
	}{
		{
			name:           "No known transactions",
			knownTxs:       new(knownTxs),
			queue:          make(map[common.Address]types.Transactions),
			pending:        make(map[common.Address]types.Transactions),
			expectedStatus: make(map[common.Hash]int),
		},
		{
			name: "All known transactions in pending",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			queue: make(map[common.Address]types.Transactions),
			pending: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(1),
				},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusPending,
				createTestTransaction(1).Hash(): TxStatusPending,
			},
		},
		{
			name: "Some known transactions not in pending",
			knownTxs: &knownTxs{
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
			queue: make(map[common.Address]types.Transactions),
			pending: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(2),
				},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusPending,
				createTestTransaction(1).Hash(): TxStatusDemoted,
				createTestTransaction(2).Hash(): TxStatusPending,
			},
		},
		{
			name: "No known transactions in pending",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			queue: make(map[common.Address]types.Transactions),
			pending: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(2),
					createTestTransaction(3),
				},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusDemoted,
				createTestTransaction(1).Hash(): TxStatusDemoted,
			},
		},
		{
			name: "Already demoted transactions",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			queue: make(map[common.Address]types.Transactions),
			pending: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(1),
				},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusDemoted,
				createTestTransaction(1).Hash(): TxStatusPending,
			},
		},
		{
			name: "All known txs in queue",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
			queue: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(1),
				},
			},
			pending: make(map[common.Address]types.Transactions),
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusQueue,
				createTestTransaction(1).Hash(): TxStatusQueue,
			},
		},
		{
			name: "Mix of known txs in queue and pending",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
				createTestTransaction(2).Hash(): {
					tx: createTestTransaction(2),
				},
				createTestTransaction(3).Hash(): {
					tx: createTestTransaction(3),
				},
			},
			queue: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
				},
			},
			pending: map[common.Address]types.Transactions{
				common.HexToAddress("0x2"): {
					createTestTransaction(1),
					createTestTransaction(2),
				},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusQueue,
				createTestTransaction(1).Hash(): TxStatusPending,
				createTestTransaction(2).Hash(): TxStatusPending,
				createTestTransaction(3).Hash(): TxStatusDemoted,
			},
		},
		{
			name: "Transaction in queue should be marked as queue status",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusPending, // Start with pending status
				},
			},
			queue: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
			},
			pending: map[common.Address]types.Transactions{},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusQueue,
			},
		},
		{
			name: "Transaction in pending should be marked as pending status",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusQueue, // Start with queue status
				},
			},
			queue: map[common.Address]types.Transactions{},
			pending: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusPending,
			},
		},
		{
			name: "Transaction in both queue and pending should be marked as queue status (queue takes precedence)",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusDemoted, // Start with demoted status
				},
			},
			queue: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
			},
			pending: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusQueue,
			},
		},
		{
			name: "Transaction not in queue or pending should be marked as demoted status",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusPending, // Start with pending status
				},
			},
			queue:   map[common.Address]types.Transactions{},
			pending: map[common.Address]types.Transactions{},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusDemoted,
			},
		},
		{
			name: "Multiple transactions with different statuses",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusDemoted,
				},
				createTestTransaction(1).Hash(): {
					tx:           createTestTransaction(1),
					promotedTime: time.Now(),
					status:       TxStatusQueue,
				},
				createTestTransaction(2).Hash(): {
					tx:           createTestTransaction(2),
					promotedTime: time.Now(),
					status:       TxStatusPending,
				},
				createTestTransaction(3).Hash(): {
					tx:           createTestTransaction(3),
					promotedTime: time.Now(),
					status:       TxStatusQueue,
				},
			},
			queue: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
				{2}: {createTestTransaction(1)},
			},
			pending: map[common.Address]types.Transactions{
				{3}: {createTestTransaction(2)},
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusQueue,   // In queue
				createTestTransaction(1).Hash(): TxStatusQueue,   // In queue
				createTestTransaction(2).Hash(): TxStatusPending, // In pending
				createTestTransaction(3).Hash(): TxStatusDemoted, // Not in queue or pending
			},
		},
		{
			name:     "Empty knownTxs should not cause any issues",
			knownTxs: &knownTxs{},
			queue: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(0)},
			},
			pending: map[common.Address]types.Transactions{
				{1}: {createTestTransaction(1)},
			},
			expectedStatus: map[common.Hash]int{},
		},
		{
			name: "Transaction with different address in queue/pending should not match",
			knownTxs: &knownTxs{
				createTestTransaction(0).Hash(): {
					tx:           createTestTransaction(0),
					promotedTime: time.Now(),
					status:       TxStatusPending,
				},
			},
			queue: map[common.Address]types.Transactions{
				{2}: {createTestTransaction(1)}, // Different transaction
			},
			pending: map[common.Address]types.Transactions{
				{3}: {createTestTransaction(2)}, // Different transaction
			},
			expectedStatus: map[common.Hash]int{
				createTestTransaction(0).Hash(): TxStatusDemoted, // Should be demoted since not found
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock bundling module
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)

			// Create builder module
			builderModule := &BuilderWrappingModule{
				txBundlingModule: mockTxBundlingModule,
				knownTxs:         tt.knownTxs.Copy(),
			}

			// Create old and new headers for PostReset
			oldHead := &types.Header{Number: big.NewInt(100)}
			newHead := &types.Header{Number: big.NewInt(101)}

			// Call PostReset
			builderModule.PostReset(oldHead, newHead, tt.queue, tt.pending)

			// Verify that all known transactions have the expected status
			for hash, expectedStatus := range tt.expectedStatus {
				knownTx, exists := builderModule.knownTxs.get(hash)
				assert.True(t, exists, "Transaction should exist in knownTxs")
				assert.Equal(t, expectedStatus, knownTx.status,
					"Transaction %s should have status %d, got %d", hash.String(), expectedStatus, knownTx.status)
			}

			// Verify that the number of transactions in knownTxs hasn't changed
			assert.Equal(t, len(*tt.knownTxs), len(*builderModule.knownTxs),
				"Number of known transactions should not change")
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
