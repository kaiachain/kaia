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
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	mock_api "github.com/kaiachain/kaia/api/mocks"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	mock_builder "github.com/kaiachain/kaia/kaiax/builder/mock"
	"github.com/kaiachain/kaia/params"
	mock_work "github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func init() {
	blockchain.InitDeriveSha(params.TestChainConfig)
}

func TestBuilderModule_PreAddTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBackend := mock_api.NewMockBackend(ctrl)
	mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
	mockTxBundlingModule.EXPECT().ExtractTxBundles(gomock.Any(), gomock.Any()).Return([]*builder.Bundle{}).AnyTimes()

	builderModule := &BuilderModule{
		InitOpts: InitOpts{
			Backend: mockBackend,
			Modules: []builder.TxBundlingModule{
				mockTxBundlingModule,
			},
		},
	}

	tests := []struct {
		name           string
		tx             *types.Transaction
		local          bool
		expectedError  error
		pendingBundles map[common.Hash]txAndTime
	}{
		{
			name:           "Valid transaction",
			tx:             types.NewTransaction(0, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
			expectedError:  nil,
			pendingBundles: make(map[common.Hash]txAndTime),
		},
		{
			name:          "Bundle transaction during lock period",
			tx:            types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
			expectedError: errors.New("Unable to add known bundle tx into tx pool during lock period"),
			pendingBundles: map[common.Hash]txAndTime{
				types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil).Hash(): {
					tx:   types.NewTransaction(1, common.Address{}, big.NewInt(0), 0, big.NewInt(0), nil),
					time: time.Now(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builderModule.pendingBundles = tt.pendingBundles
			err := builderModule.PreAddTx(nil, tt.tx, true)
			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestBuilderModule_IsModuleTx(t *testing.T) {
	b := &BuilderModule{
		pendingBundles: make(map[common.Hash]txAndTime),
	}

	tx1 := createTestTransaction(0)
	tx2 := createTestTransaction(1)

	// Add tx1 to pending bundles
	b.pendingBundles[tx1.Hash()] = txAndTime{
		tx:   tx1,
		time: time.Now(),
	}

	tests := []struct {
		name     string
		tx       *types.Transaction
		expected bool
	}{
		{
			name:     "transaction in pending bundles",
			tx:       tx1,
			expected: true,
		},
		{
			name:     "transaction not in pending bundles",
			tx:       tx2,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := b.IsModuleTx(tt.tx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBuilderModule_GetCheckBalance(t *testing.T) {
	b := &BuilderModule{}
	assert.Nil(t, b.GetCheckBalance())
}

func TestBuilderModule_IsReady(t *testing.T) {
	b := &BuilderModule{}
	assert.True(t, b.IsReady(nil, 0, nil))
}

func TestBuilderModule_PreReset_NewBundlesAndLocktime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	now := time.Now()

	tests := []struct {
		name            string
		existingBundles map[common.Hash]txAndTime
		pendingTxs      map[common.Address]types.Transactions
		expectedBundles map[common.Hash]txAndTime
	}{
		{
			name:            "No existing bundles, new bundles from unlocked txs",
			existingBundles: make(map[common.Hash]txAndTime),
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(0),
					createTestTransaction(1),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
		},
		{
			name: "Existing bundles within lock period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now,
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(1),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx: createTestTransaction(0),
				},
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
		},
		{
			name: "Existing bundles with expired lock period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-BundleLockPeriod - time.Second),
				},
			},
			pendingTxs: map[common.Address]types.Transactions{
				common.HexToAddress("0x1"): {
					createTestTransaction(1),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(1).Hash(): {
					tx: createTestTransaction(1),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock backend
			mockBackend := mock_api.NewMockBackend(ctrl)
			mockBackend.EXPECT().CurrentBlock().Return(types.NewBlock(&types.Header{
				BaseFee: common.Big0,
			}, nil, nil)).AnyTimes()
			mockBackend.EXPECT().ChainConfig().Return(&params.ChainConfig{}).AnyTimes()

			// Setup mock tx pool
			mockTxPool := mock_work.NewMockTxPool(ctrl)
			mockTxPool.EXPECT().UnlockedPending().Return(tt.pendingTxs, nil).AnyTimes()
			mockTxPool.EXPECT().GetCurrentState().Return(nil).AnyTimes()

			// Setup mock tx bundling module
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().ExtractTxBundles(gomock.Eq(flattenTxs(tt.pendingTxs)), gomock.Any()).Return(bundleTxs(tt.pendingTxs)).AnyTimes()

			// Setup BuilderModule
			builderModule := &BuilderModule{
				pendingBundles: copyTxAndTimeMap(tt.existingBundles),
				InitOpts: InitOpts{
					Backend: mockBackend,
					Modules: []builder.TxBundlingModule{
						mockTxBundlingModule,
					},
				},
			}

			builderModule.PreReset(mockTxPool, nil, nil)

			assert.Equal(t, len(tt.expectedBundles), len(builderModule.pendingBundles))

			for hash, expected := range tt.expectedBundles {
				_, ok := tt.existingBundles[hash]
				fmt.Println(ok)
				fmt.Println(hash.Hex())
				fmt.Println(len(tt.existingBundles))
				if existing, exists := tt.existingBundles[hash]; exists {
					// existing bundle tx time should not be changed
					assert.Equal(t, existing.tx.Hash(), builderModule.pendingBundles[hash].tx.Hash())
					assert.Equal(t, existing.time, builderModule.pendingBundles[hash].time)
				} else {
					// new bundle tx time should be set to current time
					assert.Equal(t, expected.tx.Hash(), builderModule.pendingBundles[hash].tx.Hash())
					assert.True(t, builderModule.pendingBundles[hash].time.After(now))
				}
			}
		})
	}
}

func TestBuilderModule_PreReset_BundleTimeout(t *testing.T) {
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
			name: "Bundle within timeout period",
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
			name: "Bundle exceeds timeout period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-BundleTimeout - time.Second),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-BundleTimeout - time.Second),
				},
			},
			expectedUnexecutable: true,
		},
		{
			name: "Bundle exceeds lock period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-BundleLockPeriod - time.Second),
				},
			},
			expectedBundles:      map[common.Hash]txAndTime{},
			expectedUnexecutable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock backend
			mockBackend := mock_api.NewMockBackend(ctrl)
			mockBackend.EXPECT().CurrentBlock().Return(types.NewBlock(&types.Header{
				BaseFee: common.Big0,
			}, nil, nil)).AnyTimes()
			mockBackend.EXPECT().ChainConfig().Return(&params.ChainConfig{}).AnyTimes()

			// Setup mock tx pool
			mockTxPool := mock_work.NewMockTxPool(ctrl)
			mockTxPool.EXPECT().UnlockedPending().Return(map[common.Address]types.Transactions{}, nil).AnyTimes()
			mockTxPool.EXPECT().GetCurrentState().Return(nil).AnyTimes()

			// Setup mock tx bundling module
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().ExtractTxBundles(gomock.Any(), gomock.Any()).Return([]*builder.Bundle{}).AnyTimes()

			// Setup BuilderModule
			builderModule := &BuilderModule{
				pendingBundles: copyTxAndTimeMap(tt.existingBundles),
				InitOpts: InitOpts{
					Backend: mockBackend,
					Modules: []builder.TxBundlingModule{
						mockTxBundlingModule,
					},
				},
			}

			builderModule.PreReset(mockTxPool, nil, nil)

			// Verify the number of bundles
			assert.Equal(t, len(tt.expectedBundles), len(builderModule.pendingBundles))

			// Verify each bundle's state
			for hash, expected := range tt.expectedBundles {
				actual, exists := builderModule.pendingBundles[hash]
				assert.True(t, exists)
				assert.Equal(t, expected.tx.Hash(), actual.tx.Hash())
				assert.Equal(t, expected.time, actual.time)
				assert.Equal(t, tt.expectedUnexecutable, actual.tx.IsMarkedUnexecutable())
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

func flattenTxs(txs map[common.Address]types.Transactions) []*types.Transaction {
	flattenedTxs := []*types.Transaction{}
	for _, txs := range txs {
		flattenedTxs = append(flattenedTxs, txs...)
	}
	return flattenedTxs
}

func bundleTxs(txs map[common.Address]types.Transactions) []*builder.Bundle {
	bundle := &builder.Bundle{
		BundleTxs: []*builder.TxOrGen{
			builder.NewTxOrGenFromGen(func(uint64) (*types.Transaction, error) { return nil, nil }, common.Hash{}),
		},
		TargetTxHash: common.Hash{},
	}
	for _, tx := range flattenTxs(txs) {
		bundle.BundleTxs = append(bundle.BundleTxs, builder.NewTxOrGenFromTx(tx))
	}
	return []*builder.Bundle{bundle}
}

func copyTxAndTimeMap(m map[common.Hash]txAndTime) map[common.Hash]txAndTime {
	newMap := make(map[common.Hash]txAndTime)
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
