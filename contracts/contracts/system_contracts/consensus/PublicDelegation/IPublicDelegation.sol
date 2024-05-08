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

import "../CnV3/ICnStakingV3.sol";
import "./IKIP163.sol";

interface IPublicDelegation is IKIP163 {
    /* ========== STRUCT ========== */

    struct PDConstructorArgs {
        address owner;
        address commissionTo;
        uint256 commissionRate;
        string gcName;
    }

    /* ========== ENUM ========== */

    /// @dev Current state of the withdrawal request
    enum WithdrawalRequestState {
        Undefined,
        Requested,
        Withdrawable,
        Withdrawn,
        PendingCancel,
        Canceled
    }

    /* ========== EVENTS ========== */

    // Initialization
    event DeployContract(string _contractType, address _baseCnStakingV3, PDConstructorArgs _pdArgs);

    // Operations
    event UpdateCommissionTo(address indexed _prevCommissionTo, address indexed _commissionTo);
    event UpdateCommissionRate(uint256 indexed _prevCommissionRate, uint256 indexed _commissionRate);
    event SendCommission(address indexed _commissionTo, uint256 _commission);

    // Staking
    event Staked(address indexed _user, uint256 _assets, uint256 _shares);
    event Redeemed(address indexed _user, address indexed _recipient, uint256 _assets, uint256 _shares);
    event Redelegated(address indexed _user, address indexed _targetCnV3, uint256 _assets);
    event RequestWithdrawal(
        address indexed _user,
        address indexed _recipient,
        uint256 indexed _requestId,
        uint256 _assets
    );
    event RequestCancelWithdrawal(address indexed _user, uint256 indexed _requestId);
    event Claimed(address indexed _user, uint256 indexed _requestId);

    /* ========== CONSTANT/IMMUTABLE GETTERS ========== */

    function MAX_COMMISSION_RATE() external pure returns (uint256);

    function COMMISSION_DENOMINATOR() external pure returns (uint256);

    function CONTRACT_TYPE() external pure returns (string memory);

    function VERSION() external pure returns (uint256);

    function baseCnStakingV3() external view returns (ICnStakingV3);

    /* ========== OPERATION FUNCTIONS ========== */

    function updateCommissionTo(address _commissionTo) external;

    function updateCommissionRate(uint256 _commissionRate) external;

    /* ========== PUBLIC FUNCTIONS ========== */

    // Staking
    function stake() external payable;

    function stakeFor(address _recipient) external payable; // Defined in IKIP163

    receive() external payable;

    // Withdrawal
    function withdraw(address _recipient, uint256 _assets) external;

    function redeem(address _recipient, uint256 _shares) external;

    function cancelApprovedStakingWithdrawal(uint256 _requestId) external;

    function claim(uint256 _requestId) external;

    // Redelegation
    function redelegateByAssets(address _targetCnV3, uint256 _assets) external;

    function redelegateByShares(address _targetCnV3, uint256 _shares) external;

    // Sweep
    function sweep() external;

    /* ========== PUBLIC VIEWS ========== */

    function commissionTo() external view returns (address);

    function commissionRate() external view returns (uint256);

    function userRequestIds(address _owner, uint256 _index) external view returns (uint256);

    function requestIdToOwner(uint256 _requestId) external view returns (address);

    function getCurrentWithdrawalRequestState(uint256 _requestId) external view returns (WithdrawalRequestState);

    function getUserRequestCount(address _owner) external view returns (uint256);

    function getUserRequestIdsWithState(
        address _owner,
        WithdrawalRequestState _state
    ) external view returns (uint256[] memory);

    function getUserRequestIds(address _owner) external view returns (uint256[] memory);

    function maxRedeem(address _owner) external view returns (uint256);

    function maxWithdraw(address _owner) external view returns (uint256);

    function reward() external view returns (uint256); // Defined in IKIP163

    function totalAssets() external view returns (uint256);

    function convertToShares(uint256 _assets) external view returns (uint256);

    function convertToAssets(uint256 _shares) external view returns (uint256);

    function previewDeposit(uint256 _assets) external view returns (uint256);

    function previewWithdraw(uint256 _assets) external view returns (uint256);

    function previewRedeem(uint256 _shares) external view returns (uint256);
}
