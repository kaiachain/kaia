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

interface ICnStakingV3 {
    /* ========== STRUCTS ========== */

    struct LockupConditions {
        uint256[] unlockTime;
        uint256[] unlockAmount;
        bool allReviewed;
        uint256 reviewedCount;
        mapping(address => bool) reviewedAdmin;
    }

    struct WithdrawalRequest {
        address to;
        uint256 value;
        uint256 withdrawableFrom;
        WithdrawalStakingState state;
    }

    /* ========== ENUM ========== */

    enum WithdrawalStakingState {
        Unknown,
        Transferred,
        Canceled
    }

    /* ========== EVENTS ========== */

    // Initialization
    event DeployCnStakingV3(
        string contractType,
        address nodeId,
        address rewardAddress,
        uint256[] unlockTime,
        uint256[] unlockAmount
    );

    event UpdateRewardAddress(address indexed rewardAddress);
    event AcceptRewardAddress(address indexed rewardAddress);
    event UpdateStakingTracker(address indexed stakingTracker);
    event UpdateVoterAddress(address indexed voterAddress);
    event UpdateGCId(uint256 indexed gcId);
    event ToggleRedelegation(bool isRedelegationEnabled);

    event DepositLockupStakingAndInit(address indexed from, uint256 value);
    event SetPublicDelegation(address indexed from, address publicDelegation, address rewardAddress);
    event ReviewInitialConditions(address indexed from);
    event CompleteReviewInitialConditions();

    event WithdrawLockupStaking(address indexed to, uint256 value);

    event DelegateKaia(address indexed from, uint256 value);
    event Redelegation(address indexed user, address indexed targetCnStakingV3, uint256 value);
    event HandleRedelegation(
        address indexed user,
        address indexed prevCnStakingV3,
        address indexed targetCnStakingV3,
        uint256 value
    );

    event ApproveStakingWithdrawal(
        uint256 indexed approvedWithdrawalId,
        address indexed to,
        uint256 value,
        uint256 withdrawableFrom
    );
    event CancelApprovedStakingWithdrawal(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value);
    event WithdrawApprovedStaking(uint256 indexed approvedWithdrawalId, address indexed to, uint256 value);

    /* ========== CONSTANT/IMMUTABLE GETTERS ========== */

    // Constants
    function CONTRACT_TYPE() external pure returns (string memory);

    function VERSION() external pure returns (uint256);

    function ADDRESS_BOOK_ADDRESS() external pure returns (address);

    function STAKE_LOCKUP() external pure returns (uint256);

    function OPERATOR_ROLE() external pure returns (bytes32);

    function ADMIN_ROLE() external pure returns (bytes32);

    function STAKER_ROLE() external pure returns (bytes32);

    function UNSTAKING_APPROVER_ROLE() external pure returns (bytes32);

    function UNSTAKING_CLAIMER_ROLE() external pure returns (bytes32);

    // Immutables
    function nodeId() external view returns (address);

    function isPublicDelegationEnabled() external view returns (bool);

    /* ========== OPERATION FUNCTIONS ========== */

    function setStakingTracker(address _tracker) external;

    function setGCId(uint256 _gcId) external;

    function setPublicDelegation(address _pdFactory, bytes memory _pdArgs) external;

    function reviewInitialConditions() external;

    function depositLockupStakingAndInit() external payable;

    // Contract configuration
    function withdrawLockupStaking(address payable _to, uint256 _value) external;

    function updateRewardAddress(address _rewardAddress) external;

    function updateStakingTracker(address _tracker) external;

    function updateVoterAddress(address _voterAddress) external;

    function toggleRedelegation() external;

    // Accept reward address, not by multisig
    function acceptRewardAddress(address _rewardAddress) external;

    /* ========== STAKING FUNCTIONS ========== */

    // Staking
    function delegate() external payable;

    receive() external payable;

    // Redelegation
    function redelegate(address _user, address _targetCnV3, uint256 _value) external;

    function handleRedelegation(address _user) external payable;

    // Withdrawal
    function approveStakingWithdrawal(address _to, uint256 _value) external returns (uint256);

    function cancelApprovedStakingWithdrawal(uint256 _approvedWithdrawalId) external;

    function withdrawApprovedStaking(uint256 _approvedWithdrawalId) external;

    /* ========== PUBLIC GETTERS ========== */

    function gcId() external view returns (uint256);

    function publicDelegation() external view returns (address);

    function lastRedelegation(address _account) external view returns (uint256);

    function rewardAddress() external view returns (address);

    function pendingRewardAddress() external view returns (address);

    function initialLockupStaking() external view returns (uint256);

    function remainingLockupStaking() external view returns (uint256);

    function isInitialized() external view returns (bool);

    function staking() external view returns (uint256);

    function unstaking() external view returns (uint256);

    function withdrawalRequestCount() external view returns (uint256);

    function isRedelegationEnabled() external view returns (bool);

    function stakingTracker() external view returns (address);

    function voterAddress() external view returns (address);

    function getLockupStakingInfo()
        external
        view
        returns (
            uint256[] memory unlockTime,
            uint256[] memory unlockAmount,
            uint256 initial,
            uint256 remaining,
            uint256 withdrawable
        );

    function getApprovedStakingWithdrawalIds(
        uint256 _from,
        uint256 _to,
        WithdrawalStakingState _state
    ) external view returns (uint256[] memory ids);

    function getApprovedStakingWithdrawalInfo(
        uint256 _index
    ) external view returns (address to, uint256 value, uint256 withdrawableFrom, WithdrawalStakingState state);
}
