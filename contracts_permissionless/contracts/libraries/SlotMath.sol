// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title SlotMath
/// @notice Slot limit calculations for validator state transitions.
///
///         For n >= 4 (standard BFT with f = floor(n/3) > 0):
///           minActiveCount = ceil(2n/3) = (2n + 2) / 3
///           totalBudget    = n - minActiveCount = floor(n/3) = n / 3
///           maxSlotAvailable = ceil(totalBudget / 2) = (n/3 + 1) / 2
///
///         For n < 4 (BFT fault tolerance f = 0; all nodes must remain active):
///           minActiveCount   = n
///           maxSlotAvailable = 0
library SlotMath {
    /// @notice Returns the minimum number of active validators required for BFT liveness.
    /// @dev n >= 4: ceil(2n/3) = (2n + 2) / 3 in integer arithmetic.
    ///      n <  4: f = 0 under BFT (2/3 majority requires all nodes), so minActiveCount = n.
    /// @param n The epoch slot factor (ValActive count at epoch start)
    /// @return The minimum active validator count
    function minActiveCount(uint256 n) internal pure returns (uint256) {
        if (n < 4) return n;
        return (2 * n + 2) / 3;
    }

    /// @notice Returns the maximum number of validators allowed in ValPaused or ValExiting state each.
    /// @dev n >= 4: totalBudget = floor(n/3) = n/3; maxSlotAvailable = ceil(totalBudget/2) = (n/3 + 1) / 2.
    ///      n <  4: totalBudget = 0 (f = 0, no transitions out of ValActive allowed), so maxSlotAvailable = 0.
    /// @param n The epoch slot factor (ValActive count at epoch start)
    /// @return The maximum count for each of ValPaused and ValExiting states
    function maxSlotAvailable(uint256 n) internal pure returns (uint256) {
        if (n < 4) return 0;
        return (n / 3 + 1) / 2;
    }
}
