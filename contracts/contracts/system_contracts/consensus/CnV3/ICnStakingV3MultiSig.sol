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

import {EnumerableSet} from "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";

interface ICnStakingV3MultiSig {
    /* ========== STRUCT ========== */

    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address requestProposer;
        EnumerableSet.AddressSet confirmers;
        /// @dev Use `getRequestState` to check the state of the request
        RequestState state;
    }

    /* ========== ENUMS ========== */

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
        UpdateVoterAddress,
        ToggleRedelegation
    }

    /* ========== CONSTANTS ========== */

    function MAX_ADMIN() external view returns (uint256);

    /* ========== EVENTS ========== */

    event DeployCnStakingV3MultiSig(
        string contractType,
        address contractValidator,
        address[] cnAdminList,
        uint256 requirement
    );

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

    // Specific multisig operations
    event AddAdmin(address indexed admin);
    event DeleteAdmin(address indexed admin);
    event UpdateRequirement(uint256 requirement);
    event ClearRequest();

    /* ========== MULTISIG FUNCTIONS ========== */

    // Submit request
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

    function submitToggleRedelegation() external;

    // Specific multisig operations
    function addAdmin(address _admin) external;

    function deleteAdmin(address _admin) external;

    function updateRequirement(uint256 _requirement) external;

    function clearRequest() external;

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

    /* ========== PUBLIC GETTERS ========== */

    function contractValidator() external view returns (address);

    function isAdmin(address _admin) external view returns (bool);

    function adminList(uint256 _pos) external view returns (address);

    function requirement() external view returns (uint256);

    function lastClearedId() external view returns (uint256);

    function requestCount() external view returns (uint256);

    function getReviewers() external view returns (address[] memory reviewers);

    function getRequestState(uint256 _id) external view returns (RequestState);

    function getState()
        external
        view
        returns (
            address contractValidator,
            address nodeId,
            address rewardAddress,
            address[] memory adminListArr,
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
}
