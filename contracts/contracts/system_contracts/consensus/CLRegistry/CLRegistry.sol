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

import "openzeppelin-contracts-5.0/access/Ownable.sol";
import "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";
import "./ICLRegistry.sol";

contract CLRegistry is ICLRegistry, Ownable {
    using EnumerableSet for EnumerableSet.UintSet;

    /* ========== CONSTANTS ========== */

    address private constant ZERO_ADDRESS = address(0);

    /* ========== STATE VARIABLES ========== */

    EnumerableSet.UintSet private _gcIds;
    mapping(uint256 => CLInfo) public clPoolList; // gcId -> CL pair

    /* ========== CONSTRUCTOR ========== */

    constructor(address initialOwner) Ownable(initialOwner) {}

    /* ========== PAIR MANAGEMENT ========== */

    /// @dev See {ICLRegistry-addCLPair}
    function addCLPair(CLInfo[] calldata list) external onlyOwner {
        for (uint i = 0; i < list.length; i++) {
            uint256 gcId = list[i].gcId;
            require(_validateCLPairInput(list[i]), "CLRegistry::addCLPair: Invalid pair input");
            require(!_isExistPair(gcId), "CLRegistry::addCLPair: GC ID does exist");
            clPoolList[gcId] = list[i];
            emit RegisterPair(list[i].nodeId, list[i].gcId, list[i].clPool, list[i].clStaking);
            _addGCId(gcId);
        }
    }

    /// @dev See {ICLRegistry-removeCLPair}
    function removeCLPair(uint256 gcId) external onlyOwner {
        require(gcId != 0, "CLRegistry::removeCLPair: Invalid GC ID");
        require(_isExistPair(gcId), "CLRegistry::removeCLPair: GC ID does not exist");

        emit RetirePair(
            clPoolList[gcId].nodeId,
            clPoolList[gcId].gcId,
            clPoolList[gcId].clPool,
            clPoolList[gcId].clStaking
        );
        delete clPoolList[gcId];
        _removeGCId(gcId);
    }

    /// @dev See {ICLRegistry-updateCLPair}
    function updateCLPair(CLInfo[] calldata list) external onlyOwner {
        for (uint i = 0; i < list.length; i++) {
            uint256 gcId = list[i].gcId;
            require(_validateCLPairInput(list[i]), "CLRegistry::updateCLPair: Invalid pair input");
            require(_isExistPair(gcId), "CLRegistry::updateCLPair: GC ID does not exist");
            clPoolList[gcId] = list[i];
            emit UpdatePair(list[i].nodeId, list[i].gcId, list[i].clPool, list[i].clStaking);
        }
    }

    /// @dev See {ICLRegistry-getAllCLs}
    function getAllCLs()
        external
        view
        returns (address[] memory, uint256[] memory, address[] memory, address[] memory)
    {
        uint256 len = _gcIds.length();
        address[] memory nodeIds = new address[](len);
        uint256[] memory gcIds = new uint256[](len);
        address[] memory clPools = new address[](len);
        address[] memory clStakings = new address[](len);

        for (uint i = 0; i < len; i++) {
            CLInfo storage clInfo = clPoolList[_gcIds.at(i)];
            nodeIds[i] = clInfo.nodeId;
            gcIds[i] = clInfo.gcId;
            clPools[i] = clInfo.clPool;
            clStakings[i] = clInfo.clStaking;
        }
        return (nodeIds, gcIds, clPools, clStakings);
    }

    // @dev Return all GC IDs
    function getAllGCIds() public view returns (uint256[] memory) {
        return _gcIds.values();
    }

    /// @dev Validate property values of `CLInfo`
    function _validateCLPairInput(CLInfo calldata pairInput) internal pure returns (bool) {
        return
            pairInput.gcId != 0 &&
            pairInput.nodeId != ZERO_ADDRESS &&
            pairInput.clPool != ZERO_ADDRESS &&
            pairInput.clStaking != ZERO_ADDRESS;
    }

    /// @dev Return true if a pair exists with given `gcId`
    function _isExistPair(uint256 gcId) internal view returns (bool) {
        return clPoolList[gcId].gcId != 0;
    }

    /// @dev Add GC ID to the global GC ID list. If it exists already, do nothing
    function _addGCId(uint256 gcId) internal {
        _gcIds.add(gcId);
    }

    /// @dev Remove GC ID to the global GC ID list. If it does not exist, do nothing
    function _removeGCId(uint256 gcId) internal {
        _gcIds.remove(gcId);
    }
}
