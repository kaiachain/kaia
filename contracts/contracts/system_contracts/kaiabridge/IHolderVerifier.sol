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
pragma solidity 0.8.24;

interface IHolderVerifier {
    /* ========== EVENTS ========== */

    event RecordAdded(string indexed fnsaAddr, uint256 conyBalance);
    event RecordsAdded(uint256 count);
    event ProvisionRequested(
        string indexed fnsaAddr,
        address indexed kaiaAddr,
        uint256 conyBalance,
        uint256 kaiaAmount,
        uint64 seq
    );

    /* ========== VIEWS ========== */

    function getRecord(string calldata fnsaAddr) external view returns (uint256 conyBalance, uint64 seq);

    function getRecords(
        uint256 startIdx,
        uint256 maxCount
    ) external view returns (string[] memory fnsaAddrs, uint256[] memory conyBalances, uint64[] memory seqs);

    function getRecordCount() external view returns (uint256);

    function isProvisioned(string memory fnsaAddr) external view returns (bool);

    function allConyBalances() external view returns (uint256);

    function provisionedConyBalances() external view returns (uint256);

    function provisionedAccounts() external view returns (uint256);

    /* ========== MUTATIVE FUNCTIONS ========== */

    function addRecord(string calldata fnsaAddr, uint256 conyBalance) external;

    function addRecords(string[] calldata fnsaAddrs, uint256[] calldata conyBalances) external;

    function requestProvision(
        bytes calldata publicKey,
        string calldata fnsaAddress,
        bytes32 messageHash,
        bytes calldata signature
    ) external;
}
