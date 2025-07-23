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

interface IAuctionDepositVault {
    /* ========== STRUCT ========== */

    struct WithdrawReserveInfo {
        /// @dev The time when a withdraw reservation is invoked
        uint256 at;
        /// @dev Amount of withdrawl
        uint256 amount;
    }

    /* ========== EVENT ========== */

    /// @dev Emitted when a fee vault is changed
    event ChangeAuctionFeeVault(address oldFeeVault, address newFeeVault);

    /// @dev Emitted when a minimum deposit amount is changed
    event ChangeMinDepositAmount(uint256 oldAmount, uint256 newAmount);

    /// @dev Emitted when a minimum withdrawal locktime is changed
    event ChangeMinWithdrawLocktime(uint256 oldLocktime, uint256 newLocktime);

    /// @dev Emitted when a deposit is taken
    event VaultDeposit(
        address searcher,
        uint256 amount,
        uint256 totalAmount,
        uint256 nonce
    );

    /// @dev Emitted when a withdrawl is taken
    event VaultWithdraw(address searcher, uint256 amount, uint256 nonce);

    /// @dev Emitted when a withdrawl resveration is taken
    event VaultReserveWithdraw(address searcher, uint256 amount, uint256 nonce);

    /// @dev Emitted when a deposit balance is insufficient
    event InsufficientBalance(
        address searcher,
        uint256 balance,
        uint256 amount
    );

    /// @dev Emitted when a bid is taken
    event TakenBid(address searcher, uint256 amount);

    /// @dev Emitted when a bid fails
    event TakenBidFailed(address searcher, uint256 amount);

    /// @dev Emitted when a gas is reimbused
    event TakenGas(address searcher, uint256 gasAmount);

    /// @dev Emitted when a gas reimbursement fails
    event TakenGasFailed(address searcher, uint256 gasAmount);

    /* ========== FUNCTION INTERFACE  ========== */

    function changeMinDepositAmount(uint256 newMinAmount) external;

    function changeMinWithdrawLocktime(uint256 newLocktime) external;

    function changeAuctionFeeVault(address newAuctionFeeVault) external;

    function deposit() external payable;

    function depositFor(address searcher) external payable;

    function reserveWithdraw() external;

    function withdraw() external;

    function takeBid(address searcher, uint256 amount) external returns (bool);

    function takeGas(address searcher, uint256 gasUsed) external returns (bool);

    function depositBalances(address searcher) external view returns (uint256);

    function getDepositAddrsLength() external view returns (uint256);

    function getDepositAddrs(
        uint256 start,
        uint256 limit
    ) external view returns (address[] memory searchers);

    function isMinDepositOver(address searcher) external view returns (bool);

    function getAllAddrsOverMinDeposit(
        uint256 start,
        uint256 limit
    )
        external
        view
        returns (
            address[] memory searchers,
            uint256[] memory depositAmounts,
            uint256[] memory nonces
        );
}
