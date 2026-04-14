// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title SlotMath
/// @notice Slot limit calculations for validator state transitions.
///         All formulas derive from a single BFT threshold: minActiveCount = ceil(2n/3).
///         totalBudget = n - minActiveCount = floor(n/3), maxSlotAvailable = ceil(totalBudget/2).
library SlotMath {
    /// @notice Returns the minimum number of active validators required for BFT liveness
    /// @dev ceil(2n/3) = (2n + 2) / 3 in integer arithmetic.
    /// @param n The epoch validator count (ValActive count at epoch start)
    /// @return The minimum active validator count
    function minActiveCount(uint256 n) internal pure returns (uint256) {
        return (2 * n + 2) / 3;
    }

    /// @notice Returns the maximum number of validators allowed in ValPaused or ValExiting state each
    /// @dev ceil(floor(n/3) / 2) = (n/3 + 1) / 2 in integer arithmetic.
    /// @param n The epoch validator count (ValActive count at epoch start)
    /// @return The maximum count for each of ValPaused and ValExiting states
    function maxSlotAvailable(uint256 n) internal pure returns (uint256) {
        return (n / 3 + 1) / 2;
    }
}
