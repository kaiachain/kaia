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

import {BlsPublicKeyInfo} from "../types/Node.sol";

/// @title NodeVerifier
/// @notice Validates node registration inputs and manages address uniqueness registry.
/// @dev Used by AddressBookV2 (NodeActions) for createNode / deleteNode.
library NodeVerifier {
    error InvalidInput();
    error AddressAlreadyRegistered();

    bytes32 private constant ZERO48HASH = 0xc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293; // keccak256(hex"00"*48)
    bytes32 private constant ZERO96HASH = 0x46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21; // keccak256(hex"00"*96)

    /// @notice Validates all node registration inputs and registers addresses
    /// @dev Checks: zero addresses, mutual uniqueness, BLS sizes + non-zero, registry uniqueness.
    /// @param registry The address uniqueness registry
    /// @param nodeId The node address
    /// @param stakingContract The staking contract address
    /// @param rewardAddress The reward address
    /// @param blsInfo The BLS public key and proof-of-possession
    function registerNode(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress,
        BlsPublicKeyInfo memory blsInfo
    ) internal {
        // Zero address checks
        if (nodeId == address(0) || stakingContract == address(0) || rewardAddress == address(0)) revert InvalidInput();
        // Mutual uniqueness among registerable addresses
        if (nodeId == stakingContract || stakingContract == rewardAddress || nodeId == rewardAddress) {
            revert InvalidInput();
        }
        // BLS public key (48 bytes) and proof-of-possession (96 bytes)
        if (blsInfo.publicKey.length != 48 || blsInfo.pop.length != 96) revert InvalidInput();
        if (keccak256(blsInfo.publicKey) == ZERO48HASH || keccak256(blsInfo.pop) == ZERO96HASH) revert InvalidInput();

        // Registry uniqueness check
        if (registry[nodeId] || registry[stakingContract] || registry[rewardAddress]) revert AddressAlreadyRegistered();

        // Register addresses
        registry[nodeId] = true;
        registry[stakingContract] = true;
        registry[rewardAddress] = true;
    }

    /// @notice Unregisters previously registered addresses
    function unregisterAddresses(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress
    ) internal {
        registry[nodeId] = false;
        registry[stakingContract] = false;
        registry[rewardAddress] = false;
    }
}
