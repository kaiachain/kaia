// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @dev Type-safe handle for an ABv2Storage field offset.
type ConfigSlot is uint256;

using {ABv2ConfigLib.updateUint, ABv2ConfigLib.updateAddress} for ConfigSlot global;

/// @title ABv2ConfigLib
library ABv2ConfigLib {
    error InvalidInput();

    /// @dev ERC-7201 storage location for ABv2Storage
    /// keccak256(abi.encode(uint256(keccak256("addressbookv2.storage.ABv2Storage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 internal constant STORAGE_LOCATION = 0x1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00;

    /* ========== Config slot offsets ========== */

    ConfigSlot internal constant EXIT_THRESHOLD = ConfigSlot.wrap(9);
    ConfigSlot internal constant PAUSE_TIMEOUT = ConfigSlot.wrap(10);
    ConfigSlot internal constant IDLE_TIMEOUT = ConfigSlot.wrap(11);
    ConfigSlot internal constant MAX_VALIDATOR_COUNT = ConfigSlot.wrap(12);
    ConfigSlot internal constant MAX_READY_CANDIDATE_COUNT = ConfigSlot.wrap(13);
    // offset 14 = epochValCount (system-managed, not configurable)
    ConfigSlot internal constant KEF_ADDRESS = ConfigSlot.wrap(15);
    ConfigSlot internal constant KIF_ADDRESS = ConfigSlot.wrap(16);
    ConfigSlot internal constant KPF_ADDRESS = ConfigSlot.wrap(17);

    /* ========== Typed update functions ========== */

    /// @notice Validates nonzero, swaps storage value, returns old value.
    function updateUint(ConfigSlot slot, uint256 newValue) internal returns (uint256 oldValue) {
        if (newValue == 0) revert InvalidInput();
        bytes32 s = bytes32(uint256(STORAGE_LOCATION) + ConfigSlot.unwrap(slot));
        assembly {
            oldValue := sload(s)
            sstore(s, newValue)
        }
    }

    /// @notice Validates nonzero address, swaps storage value, returns old value.
    function updateAddress(ConfigSlot slot, address newValue) internal returns (address oldValue) {
        if (newValue == address(0)) revert InvalidInput();
        bytes32 s = bytes32(uint256(STORAGE_LOCATION) + ConfigSlot.unwrap(slot));
        assembly {
            oldValue := sload(s)
            sstore(s, newValue)
        }
    }
}
