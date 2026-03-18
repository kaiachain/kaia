// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title IRegistry
/// @notice Minimal interface for the Kaia system registry at 0x401.
/// @dev Full implementation at contracts-klaytn-v1.12/contracts/Registry.sol
interface IRegistry {
    /// @notice Returns the active address for a registered system contract name.
    /// @param name The registered name (e.g., "ABv2DataContract")
    /// @return The active contract address, or address(0) if none
    function getActiveAddr(string memory name) external view returns (address);
}
