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

import "openzeppelin-contracts-5.0/access/Ownable.sol";
import "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";
import "./AuctionError.sol";
import "./IAuctionFeeVault.sol";

contract AuctionFeeVault is IAuctionFeeVault, Ownable, AuctionError {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    address public constant ADDRESS_BOOK =
        0x0000000000000000000000000000000000000400;

    uint256 public constant MAX_PAYBACK_RATE = 10000;

    /* ========== STATE VARIABLES ========== */

    uint256 public validatorPaybackRate;
    uint256 public searcherPaybackRate;
    uint256 public accumulatedBids;

    mapping(address => address) private _nodeIdToRewardAddr;

    /* ========== CONSTRUCTOR ========== */

    constructor(
        address initialOwner,
        uint256 _searcherPaybackRate,
        uint256 _validatorPaybackRate
    ) Ownable(initialOwner) {
        if (!_checkRate(_searcherPaybackRate, _validatorPaybackRate))
            revert InvalidInput();
        searcherPaybackRate = _searcherPaybackRate;
        validatorPaybackRate = _validatorPaybackRate;
    }

    /* ========== FEE MANAGEMENT ========== */

    /// @dev Take a bid from a searcher
    /// The sum of paybackAmount and validatorPayback is always less than or equal to originalAmount
    /// @param searcher The address of the searcher
    function takeBid(address searcher) external payable override {
        if (msg.value == 0) return;

        uint256 originalAmount = msg.value;
        accumulatedBids += originalAmount;

        uint256 searcherPaybackAmount = (originalAmount * searcherPaybackRate) /
            MAX_PAYBACK_RATE;
        if (searcherPaybackAmount > 0) {
            /// @dev Do not revert if the searcher payback fails
            /// Need to restrict the gas limit for deterministic gas calculation
            if (!payable(searcher).send(searcherPaybackAmount)) {
                searcherPaybackAmount = 0;
                emit FeeDepositFailed(searcher, originalAmount);
            }
        }

        uint256 validatorPayback = (originalAmount * validatorPaybackRate) /
            MAX_PAYBACK_RATE;
        if (validatorPayback > 0) {
            address rewardAddr = _nodeIdToRewardAddr[block.coinbase];
            if (rewardAddr != address(0)) {
                /// @dev Do not revert if the validator payback fails
                /// Need to restrict the gas limit for deterministic gas calculation
                if (!payable(rewardAddr).send(validatorPayback)) {
                    validatorPayback = 0;
                    emit FeeDepositFailed(block.coinbase, originalAmount);
                }
            } else {
                validatorPayback = 0;
            }
        }

        emit FeeDeposit(
            block.coinbase,
            originalAmount,
            searcherPaybackAmount,
            validatorPayback
        );
    }

    /// @dev Withdraw KAIA from the vault
    /// @param to The address to withdraw to
    function withdraw(address to) external override onlyOwner {
        uint256 amount = address(this).balance;
        (bool success, ) = to.call{value: amount}("");
        if (!success) revert WithdrawalFailed();

        emit FeeWithdrawal(amount);
    }

    /// @dev Set the searcher payback rate
    /// @param _searcherPaybackRate The searcher payback rate (10000 = 100%)
    function setSearcherPaybackRate(
        uint256 _searcherPaybackRate
    ) external override onlyOwner {
        if (!_checkRate(_searcherPaybackRate, validatorPaybackRate))
            revert InvalidInput();
        searcherPaybackRate = _searcherPaybackRate;

        emit SearcherPaybackRateUpdated(searcherPaybackRate);
    }

    /// @dev Set the validator payback rate
    /// @param _validatorPaybackRate The validator payback rate (10000 = 100%)
    function setValidatorPaybackRate(
        uint256 _validatorPaybackRate
    ) external override onlyOwner {
        if (!_checkRate(_validatorPaybackRate, searcherPaybackRate))
            revert InvalidInput();
        validatorPaybackRate = _validatorPaybackRate;

        emit ValidatorPaybackRateUpdated(validatorPaybackRate);
    }

    /* ========== REGISTRATION ========== */

    /// @dev Register the reward address for a node
    /// @param nodeId The CN node ID registered as a validator
    /// @param rewardAddr The reward recipient address
    function registerRewardAddress(
        address nodeId,
        address rewardAddr
    ) external override {
        /// @dev If there's no corresponding staking contract, it will revert
        (, address staking, ) = IAddressBook(ADDRESS_BOOK).getCnInfo(nodeId);

        if (!IStaking(staking).isAdmin(msg.sender)) revert OnlyStakingAdmin();

        _nodeIdToRewardAddr[nodeId] = rewardAddr;

        emit RewardAddressRegistered(nodeId, rewardAddr);
    }

    /* ========== HELPERS ========== */

    function _checkRate(
        uint256 _rateA,
        uint256 _rateB
    ) internal pure returns (bool) {
        return _rateA + _rateB <= MAX_PAYBACK_RATE;
    }

    /* ========== GETTERS ========== */

    function getRewardAddr(
        address nodeId
    ) external view override returns (address) {
        return _nodeIdToRewardAddr[nodeId];
    }
}

interface IAddressBook {
    function getCnInfo(
        address _cnNodeId
    ) external view returns (address, address, address);
}

interface IStaking {
    function isAdmin(address _admin) external view returns (bool);
}
