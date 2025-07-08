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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/kaiax"
	mock_kaiax "github.com/kaiachain/kaia/kaiax/mock"
	mock_builder "github.com/kaiachain/kaia/work/builder/mock"
	"github.com/stretchr/testify/assert"
)

func TestWrapAndConcatenateBundlingModules(t *testing.T) {
	// Create mock modules
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create test modules
	mockTxPool1 := mock_kaiax.NewMockTxPoolModule(ctrl)
	mockTxPool2 := mock_kaiax.NewMockTxPoolModule(ctrl)
	mockTxBundling1 := mock_builder.NewMockTxBundlingModule(ctrl)
	mockTxBundling2 := mock_builder.NewMockTxBundlingModule(ctrl)

	// Create a module that implements both interfaces
	mockBoth1 := struct {
		*mock_kaiax.MockTxPoolModule
		*mock_builder.MockTxBundlingModule
	}{
		MockTxPoolModule:     mock_kaiax.NewMockTxPoolModule(ctrl),
		MockTxBundlingModule: mock_builder.NewMockTxBundlingModule(ctrl),
	}
	mockBoth2 := struct {
		*mock_kaiax.MockTxPoolModule
		*mock_builder.MockTxBundlingModule
	}{
		MockTxPoolModule:     mock_kaiax.NewMockTxPoolModule(ctrl),
		MockTxBundlingModule: mock_builder.NewMockTxBundlingModule(ctrl),
	}

	testCases := []struct {
		name        string
		mTxBundling []kaiax.TxBundlingModule
		mTxPool     []kaiax.TxPoolModule
		expected    []interface{}
	}{
		{
			name:        "No modules",
			mTxBundling: []kaiax.TxBundlingModule{},
			mTxPool:     []kaiax.TxPoolModule{},
			expected:    []interface{}{},
		},
		{
			name:        "Only TxPool modules",
			mTxBundling: []kaiax.TxBundlingModule{},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockTxPool2},
			expected:    []interface{}{mockTxPool1, mockTxPool2},
		},
		{
			name:        "Only TxBundling modules",
			mTxBundling: []kaiax.TxBundlingModule{mockTxBundling1, mockTxBundling2},
			mTxPool:     []kaiax.TxPoolModule{},
			expected:    []interface{}{mockTxBundling1, mockTxBundling2},
		},
		{
			name:        "Mixed modules",
			mTxBundling: []kaiax.TxBundlingModule{mockTxBundling1, mockTxBundling2},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockTxPool2},
			expected:    []interface{}{mockTxPool1, mockTxPool2, mockTxBundling1, mockTxBundling2},
		},
		{
			name:        "Overlapping modules",
			mTxBundling: []kaiax.TxBundlingModule{mockBoth2, mockTxBundling1, mockBoth1},
			mTxPool:     []kaiax.TxPoolModule{mockTxPool1, mockBoth1, mockBoth2},
			expected:    []interface{}{mockTxPool1, mockBoth1, mockBoth2, mockTxBundling1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := WrapAndConcatenateBundlingModules(tc.mTxBundling, tc.mTxPool)

			// Check the length of the result
			assert.Equal(t, len(tc.expected), len(result))

			// Verify that modules are properly wrapped
			for i, module := range result {
				if wrapped, ok := module.(*BuilderWrappingModule); ok {
					assert.Equal(t, wrapped.txBundlingModule, tc.expected[i])
				} else {
					assert.Equal(t, module, tc.expected[i])
				}
			}
		})
	}
}
