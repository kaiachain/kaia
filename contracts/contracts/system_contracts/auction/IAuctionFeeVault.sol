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

pragma solidity 0.8.25;

interface IAuctionFeeVault {
    /* ========== EVENTS ========== */

    event FeePaybackFailed(address indexed receiver, uint256 amount);

    event FeeDeposit(
        address indexed sender,
        uint256 amount,
        uint256 paybackAmount,
        uint256 validatorPaybackAmount
    );

    event FeeWithdrawal(uint256 amount);

    event SearcherPaybackRateUpdated(uint256 searcherPaybackRate);

    event ValidatorPaybackRateUpdated(uint256 validatorPaybackRate);

    event RewardAddressRegistered(
        address indexed nodeId,
        address indexed reward
    );

    /* ========== FUNCTION INTERFACE ========== */

    function takeBid(address searcher) external payable;

    function withdraw(address to) external;

    function registerRewardAddress(address nodeId, address rewardAddr) external;

    function setSearcherPaybackRate(uint256 _searcherPaybackRate) external;

    function setValidatorPaybackRate(uint256 _validatorPaybackRate) external;

    function getRewardAddr(address nodeId) external view returns (address);
}
