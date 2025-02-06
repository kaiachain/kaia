// Copyright 2024 The Kaia Authors
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

package supply

import (
	"github.com/kaiachain/kaia/kaiax"
)

//go:generate mockgen -destination=./mock/module.go -package=mock github.com/kaiachain/kaia/kaiax/supply SupplyModule
type SupplyModule interface {
	kaiax.BaseModule
	kaiax.JsonRpcModule
	kaiax.ExecutionModule
	kaiax.RewindableModule

	// GetTotalSupply returns the total supply at the given block number.
	// Returns (nil, err) if essential components are missing.
	// Returns (ts, err) if partial components (e.g. canonical burn amounts) are missing.
	// Otherwise, returns (ts, nil).
	GetTotalSupply(num uint64) (*TotalSupply, error)
}
