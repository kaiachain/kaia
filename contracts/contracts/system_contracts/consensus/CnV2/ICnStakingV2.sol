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

interface ICnStakingV2 {
    // Initialization
    event DeployContract(
        string contractType,
        address contractValidator,
        address nodeId,
        address rewardAddress,
        address[] cnAdminList,
        uint256 requirement,
        uint256[] unlockTime,
        uint256[] unlockAmount
    );
    event ReviewInitialConditions(address indexed from);
    event CompleteReviewInitialConditions();
    event DepositLockupStakingAndInit(address from, uint256 value);

    // Multisig operation in general
    event SubmitRequest(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ConfirmRequest(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address[] confirmers
    );
    event RevokeConfirmation(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg,
        address[] confirmers
    );
    event CancelRequest(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ExecuteRequestSuccess(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ExecuteRequestFailure(
        uint256 indexed id,
        address indexed from,
        Functions functionId,
        bytes32 firstArg,
        bytes32 secondArg,
        bytes32 thirdArg
    );
    event ClearRequest();

    // Specific multisig operations
    event AddAdmin(address indexed admin);
    event DeleteAdmin(address indexed admin);
    event UpdateRequirement(uint256 requirement);
    event WithdrawLockupStaking(address indexed to, uint256 value);
    event ApproveStakingWithdrawal(uint256 approvedWithdrawalId, address to, uint256 value, uint256 withdrawableFrom);
    event CancelApprovedStakingWithdrawal(uint256 approvedWithdrawalId, address to, uint256 value);
    event UpdateRewardAddress(address rewardAddress);
    event UpdateStakingTracker(address stakingTracker);
    event UpdateVoterAddress(address voterAddress);
    event UpdateGCId(uint256 gcId);

    // Public functions
    event StakeKlay(address from, uint256 value);
    event WithdrawApprovedStaking(uint256 approvedWithdrawalId, address to, uint256 value);
    event AcceptRewardAddress(address rewardAddress);

    // Emitted from AddressBook
    event ReviseRewardAddress(address cnNodeId, address prevRewardAddress, address curRewardAddress);

    enum RequestState {
        Unknown,
        NotConfirmed,
        Executed,
        ExecutionFailed,
        Canceled
    }
    enum Functions {
        Unknown,
        AddAdmin,
        DeleteAdmin,
        UpdateRequirement,
        ClearRequest,
        WithdrawLockupStaking,
        ApproveStakingWithdrawal,
        CancelApprovedStakingWithdrawal,
        UpdateRewardAddress,
        UpdateStakingTracker,
        UpdateVoterAddress
    }
    enum WithdrawalStakingState {
        Unknown,
        Transferred,
        Canceled
    }

    // Constants
    function MAX_ADMIN() external view returns (uint256);

    function CONTRACT_TYPE() external view returns (string memory);

    function VERSION() external view returns (uint256);

    function ADDRESS_BOOK_ADDRESS() external view returns (address);

    function STAKE_LOCKUP() external view returns (uint256);

    // Initialization
    function setStakingTracker(address _tracker) external;

    function setGCId(uint256 _gcId) external;

    function reviewInitialConditions() external;

    function depositLockupStakingAndInit() external payable;

    // Submit multisig request
    function submitAddAdmin(address _admin) external;

    function submitDeleteAdmin(address _admin) external;

    function submitUpdateRequirement(uint256 _requirement) external;

    function submitClearRequest() external;

    function submitWithdrawLockupStaking(address payable _to, uint256 _value) external;

    function submitApproveStakingWithdrawal(address _to, uint256 _value) external;

    function submitCancelApprovedStakingWithdrawal(uint256 _approvedWithdrawalId) external;

    function submitUpdateRewardAddress(address _rewardAddress) external;

    function submitUpdateStakingTracker(address _tracker) external;

    function submitUpdateVoterAddress(address _voterAddress) external;

    // Specific multisig operations
    function addAdmin(address _admin) external;

    function deleteAdmin(address _admin) external;

    function updateRequirement(uint256 _requirement) external;

    function clearRequest() external;

    function withdrawLockupStaking(address payable _to, uint256 _value) external;

    function approveStakingWithdrawal(address _to, uint256 _value) external;

    function cancelApprovedStakingWithdrawal(uint256 _approvedWithdrawalId) external;

    function updateRewardAddress(address _rewardAddress) external;

    function updateStakingTracker(address _tracker) external;

    function updateVoterAddress(address _voterAddress) external;

    // Confirm multisig request
    function confirmRequest(
        uint256 _id,
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external;

    function revokeConfirmation(
        uint256 _id,
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external;

    // Public functions
    function stakeKlay() external payable;

    receive() external payable;

    function withdrawApprovedStaking(uint256 _approvedWithdrawalId) external;

    function acceptRewardAddress(address _rewardAddress) external;

    // Getters
    function contractValidator() external returns (address);

    function adminList(uint256 idx) external returns (address);

    function requirement() external returns (uint256);

    function isAdmin(address) external returns (bool);

    function lastClearedId() external returns (uint256);

    function requestCount() external returns (uint256);

    function initialLockupStaking() external returns (uint256);

    function remainingLockupStaking() external returns (uint256);

    function isInitialized() external returns (bool);

    function staking() external returns (uint256);

    function unstaking() external returns (uint256);

    function withdrawalRequestCount() external returns (uint256);

    function gcId() external view returns (uint256);

    function nodeId() external view returns (address);

    function rewardAddress() external view returns (address);

    function pendingRewardAddress() external view returns (address);

    function stakingTracker() external view returns (address);

    function voterAddress() external view returns (address);

    function getReviewers() external view returns (address[] memory reviewers);

    function getState()
        external
        view
        returns (
            address contractValidator,
            address nodeId,
            address rewardAddress,
            address[] memory adminList,
            uint256 requirement,
            uint256[] memory unlockTime,
            uint256[] memory unlockAmount,
            bool allReviewed,
            bool isInitialized
        );

    function getRequestIds(
        uint256 _from,
        uint256 _to,
        RequestState _state
    ) external view returns (uint256[] memory ids);

    function getRequestInfo(
        uint256 _id
    )
        external
        view
        returns (
            Functions functionId,
            bytes32 firstArg,
            bytes32 secondArg,
            bytes32 thirdArg,
            address proposer,
            address[] memory confirmers,
            RequestState state
        );

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
