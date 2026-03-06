// Copyright 2025 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title SlotMath
/// @notice Slot limit calculations for validator slot limit calculations
library SlotMath {
    /// @notice Returns the maximum number of total changeable slots per epoch
    function f(uint256 n) internal pure returns (uint256) {
        if (n == 0) return 0;
        return (n - 1) / 3;
    }

    /// @notice Returns the maximum number of validators allowed in ValPaused or ValExiting state each
    /// @dev Returns max(1, F/2) where F is the maximum number of total changeable slots per epoch
    ///      N = epoch validator count (VA + VP at epoch start). Returns 0 if N < 2.
    /// @param n The number of epoch validators
    /// @return The maximum count for each of ValPaused and ValExiting states
    function maxSlotAvailable(uint256 n) internal pure returns (uint256) {
        if (n < 2) return 0;
        uint256 cnt = f(n) / 2;
        return cnt < 1 ? 1 : cnt;
    }

    /// @notice Returns the minimum number of active validators required for BFT liveness
    /// @dev Takes the stricter of two constraints:
    ///      - BFT liveness: 2F+1 active validators (from N = 3F+1)
    ///      - Budget-based: n - 2 * maxSlotAvailable(n)
    ///      For most n these are equal; for n=4 the BFT threshold (3) is stricter
    ///      than the budget-based minimum (2), preventing a pause-exit-pause cycle
    ///      from draining active validators below the BFT liveness threshold.
    /// @param n The epoch validator count
    /// @return The minimum active validator count
    function minActiveCount(uint256 n) internal pure returns (uint256) {
        if (n == 0) return 0;
        uint256 bftMin = 2 * f(n) + 1;
        uint256 budgetMin = n - 2 * maxSlotAvailable(n);
        uint256 cnt = bftMin > budgetMin ? bftMin : budgetMin;
        return cnt < 1 ? 1 : cnt;
    }
}
