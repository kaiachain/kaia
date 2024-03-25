// Copyright 2022 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/**
 * @dev Interface of the GovParam Contract
 */
interface IGovParam {
    struct Param {
        uint256 activation;
        bool exists;
        bytes val;
    }

    event SetParam(string name, bool exists, bytes value, uint256 activation);

    function setParam(
        string calldata name, bool exists, bytes calldata value,
        uint256 activation) external;

    function setParamIn(
        string calldata name, bool exists, bytes calldata value,
        uint256 relativeActivation) external;

    /// All (including soft-deleted) param names ever existed
    function paramNames(uint256 idx) external view returns (string memory);
    function getAllParamNames() external view returns (string[] memory);

    /// Raw checkpoints
    function checkpoints(string calldata name) external view
        returns(Param[] memory);
    function getAllCheckpoints() external view
        returns(string[] memory, Param[][] memory);

    /// Any given stored (including soft-deleted) params
    function getParam(string calldata name) external view
        returns(bool, bytes memory);
    function getParamAt(string calldata name, uint256 blockNumber) external view
        returns(bool, bytes memory);

    /// All existing params
    function getAllParams() external view
        returns (string[] memory, bytes[] memory);
    function getAllParamsAt(uint256 blockNumber) external view
        returns(string[] memory, bytes[] memory);
}
