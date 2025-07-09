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
	"slices"

	"github.com/kaiachain/kaia/kaiax"
)

// WrapAndConcatenateBundlingModules wraps bundling modules and concatenates them.
// given: mTxPool = [ A, B, C ], mTxBundling = [ B, D ]
// want : mTxPool = [ A, WB, C, WD ] (W: wrapped)
func WrapAndConcatenateBundlingModules(mTxBundling []kaiax.TxBundlingModule, mTxPool []kaiax.TxPoolModule) []kaiax.TxPoolModule {
	ret := make([]kaiax.TxPoolModule, 0, len(mTxBundling)+len(mTxPool))

	for _, module := range mTxPool {
		if txb, ok := module.(kaiax.TxBundlingModule); ok {
			ret = append(ret, NewBuilderWrappingModule(txb))
		} else {
			ret = append(ret, module)
		}
	}

	for _, module := range mTxBundling {
		if txp, ok := module.(kaiax.TxPoolModule); !(ok && slices.Contains(mTxPool, txp)) {
			ret = append(ret, NewBuilderWrappingModule(module))
		}
	}

	return ret
}
