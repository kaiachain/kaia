// Copyright 2024 The kaia Authors
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
pragma solidity 0.8.25;

interface ICLRegistry {
    /* ========== STRUCT ========== */

    struct CLInfo {
        /// @dev The node ID of the validator
        address nodeId;
        /// @dev The governance committee ID of the validator
        uint256 gcId;
        /// @dev The address of the CLDEX pool
        address clPool;
        /// @dev The address of the CLStaking
        address clStaking;
    }

    /* ========== EVENT ========== */

    /// @dev Emitted when a pair is registered
    event RegisterPair(address nodeId, uint256 indexed gcId, address indexed clPool, address indexed clStaking);

    /// @dev Emitted when a pair is retired
    event RetirePair(address nodeId, uint256 indexed gcId, address indexed clPool, address indexed clStaking);

    /// @dev Emitted when a pair is updated
    event UpdatePair(address nodeId, uint256 indexed gcId, address indexed clPool, address indexed clStaking);

    /// @dev Register CL pair(s)
    /// @param list A struct of CLInfo
    function addCLPair(CLInfo[] calldata list) external;

    /// @dev Retire CL pair
    /// @param gcId A unique GC ID
    function removeCLPair(uint256 gcId) external;

    /// @dev Update CL pair(s)
    /// @param list A struct of CLInfo
    function updateCLPair(CLInfo[] calldata list) external;

    /// @dev Returns the CL information of all registered validators
    /// @return The CL information of all registered validators
    function getAllCLs() external view returns (address[] memory, uint256[] memory, address[] memory, address[] memory);
}
