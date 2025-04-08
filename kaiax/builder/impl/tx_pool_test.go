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
	"github.com/kaiachain/kaia/kaiax/builder"
	mock_builder "github.com/kaiachain/kaia/kaiax/builder/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func init() {
	blockchain.InitDeriveSha(params.TestChainConfig)
}

func TestPreAddTx_LockPeriod(t *testing.T) {
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
			name: "Transaction during lock period",
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
			name: "Transaction after lock period",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now().Add(-KnownTxTimeout - time.Second),
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(true).AnyTimes()

			builderModule := &BuilderModule{
				InitOpts: InitOpts{
					Modules: []builder.TxBundlingModule{
						mockTxBundlingModule,
					},
				},
				knownTxs: tt.knownTxs,
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

func TestPreAddTx_BundleTransaction(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		tx         *types.Transaction
		knownTxs   map[common.Hash]txAndTime
		isBundleTx bool
	}{
		{
			name:       "Non-bundle transaction",
			tx:         createTestTransaction(0),
			knownTxs:   make(map[common.Hash]txAndTime),
			isBundleTx: false,
		},
		{
			name:       "New bundle transaction",
			tx:         createTestTransaction(0),
			knownTxs:   make(map[common.Hash]txAndTime),
			isBundleTx: true,
		},
		{
			name: "Existing bundle transaction",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now(),
				},
			},
			isBundleTx: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(tt.isBundleTx).AnyTimes()

			builderModule := &BuilderModule{
				InitOpts: InitOpts{
					Modules: []builder.TxBundlingModule{
						mockTxBundlingModule,
					},
				},
				knownTxs: tt.knownTxs,
			}

			builderModule.PreAddTx(tt.tx, true)

			if tt.isBundleTx {
				assert.Contains(t, builderModule.knownTxs, tt.tx.Hash())
				if existing, exists := tt.knownTxs[tt.tx.Hash()]; exists {
					// If transaction was already in knownTxs, time should be preserved
					assert.Equal(t, existing.time, builderModule.knownTxs[tt.tx.Hash()].time)
				} else {
					// New bundle transaction should have current time
					assert.True(t, time.Since(builderModule.knownTxs[tt.tx.Hash()].time) < time.Second)
				}
			} else {
				assert.NotContains(t, builderModule.knownTxs, tt.tx.Hash())
			}
		})
	}
}

func TestIsModuleTx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name       string
		tx         *types.Transaction
		knownTxs   map[common.Hash]txAndTime
		isBundleTx bool
		expected   bool
	}{
		{
			name: "Transaction in knownTxs",
			tx:   createTestTransaction(0),
			knownTxs: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: time.Now(),
				},
			},
			isBundleTx: false,
			expected:   true,
		},
		{
			name:       "Transaction not in knownTxs but is bundle tx",
			tx:         createTestTransaction(0),
			knownTxs:   make(map[common.Hash]txAndTime),
			isBundleTx: true,
			expected:   true,
		},
		{
			name:       "Transaction not in knownTxs and not bundle tx",
			tx:         createTestTransaction(0),
			knownTxs:   make(map[common.Hash]txAndTime),
			isBundleTx: false,
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockTxBundlingModule := mock_builder.NewMockTxBundlingModule(ctrl)
			mockTxBundlingModule.EXPECT().IsBundleTx(tt.tx).Return(tt.isBundleTx).AnyTimes()

			builderModule := &BuilderModule{
				InitOpts: InitOpts{
					Modules: []builder.TxBundlingModule{
						mockTxBundlingModule,
					},
				},
				knownTxs: tt.knownTxs,
			}

			result := builderModule.IsModuleTx(tt.tx)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetCheckBalance(t *testing.T) {
	b := &BuilderModule{}
	assert.Nil(t, b.GetCheckBalance())
}

func TestIsReady(t *testing.T) {
	b := &BuilderModule{}
	assert.True(t, b.IsReady(nil, 0, nil))
}

func TestPreReset(t *testing.T) {
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
					time: now.Add(-PendingTimeout - time.Second),
				},
			},
			expectedBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-PendingTimeout - time.Second),
				},
			},
			expectedUnexecutable: true,
		},
		{
			name: "Bundle exceeds lock period",
			existingBundles: map[common.Hash]txAndTime{
				createTestTransaction(0).Hash(): {
					tx:   createTestTransaction(0),
					time: now.Add(-KnownTxTimeout - time.Second),
				},
			},
			expectedBundles:      map[common.Hash]txAndTime{},
			expectedUnexecutable: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup BuilderModule
			builderModule := &BuilderModule{
				knownTxs: copyTxAndTimeMap(tt.existingBundles),
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
