// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/vm/gas.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package vm

import "github.com/holiman/uint256"

// Gas costs
const (
	GasZero        uint64 = 0  // G_zero
	GasQuickStep   uint64 = 2  // G_base
	GasFastestStep uint64 = 3  // G_verylow
	GasFastStep    uint64 = 5  // G_low
	GasMidStep     uint64 = 8  // G_mid
	GasSlowStep    uint64 = 10 // G_high or G_exp
	GasExtStep     uint64 = 20 // G_blockhash
)

// callGas returns the actual gas cost of the call.
//
// The returned gas is gas - base * 63 / 64.
func callGas(availableGas, base uint64, callCost *uint256.Int) (uint64, error) {
	availableGas = availableGas - base
	gas := availableGas - availableGas/64
	// If the bit length exceeds 64 bit we know that the newly calculated "gas" for EIP150
	// is smaller than the requested amount. Therefore we return the new gas instead
	// of returning an error.
	if !callCost.IsUint64() || gas < callCost.Uint64() {
		return gas, nil
	}

	return callCost.Uint64(), nil
}
