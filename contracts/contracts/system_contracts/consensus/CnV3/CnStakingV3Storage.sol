// Copyright 2024 The klaytn Authors
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
pragma solidity 0.8.25;

import "./ICnStakingV3.sol";

abstract contract CnStakingV3Storage is ICnStakingV3 {
    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "CnStakingContract";
    uint256 public constant VERSION = 3;
    uint256 public constant ONE_WEEK = 1 weeks;
    uint256 public constant STAKE_LOCKUP = ONE_WEEK;
    address public constant ADDRESS_BOOK_ADDRESS = 0x0000000000000000000000000000000000000400;

    // Role IDs
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant STAKER_ROLE = keccak256("STAKER_ROLE");
    bytes32 public constant UNSTAKING_APPROVER_ROLE = keccak256("UNSTAKING_APPROVER_ROLE");
    bytes32 public constant UNSTAKING_CLAIMER_ROLE = keccak256("UNSTAKING_CLAIMER_ROLE");

    /* ========== IMMUTABLES ========== */

    /// @dev Informational node id
    address public immutable nodeId;

    /// @dev Flag to enable public delegation
    bool public immutable isPublicDelegationEnabled;

    /* ========== STATE VARIABLES ========== */

    // GC Id for the contract
    uint256 public gcId;

    // Public delegation
    address public publicDelegation; // Public delegation contract
    mapping(address => uint256) public lastRedelegation; // Last delegation time

    // Reward address
    address public rewardAddress; // Reward address
    address public pendingRewardAddress; // Pending reward address

    // Initial lockup
    LockupConditions public lockupConditions;
    uint256 public initialLockupStaking;
    uint256 public remainingLockupStaking;
    bool public isInitialized;

    // Delegated stakes
    uint256 public staking;
    uint256 public unstaking;
    uint256 public withdrawalRequestCount;
    mapping(uint256 => WithdrawalRequest) public withdrawalRequestMap;

    // Redelegation activation flag
    bool public isRedelegationEnabled;

    // External accounts
    address public stakingTracker; // used to call refreshStake(), refreshVoter()
    address public voterAddress; // read by StakingTracker
}
