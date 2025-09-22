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
import "./IAuctionEntryPoint.sol";
import "./IAuctionFeeVault.sol";
import "./IAuctionDepositVault.sol";
import "./AuctionError.sol";

contract AuctionDepositVault is IAuctionDepositVault, AuctionError, Ownable {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    address public constant REGISTRY_ADDRESS =
        0x0000000000000000000000000000000000000401;

    /* ========== STATE VARIABLES ========== */

    uint256 public minDepositAmount = 10 ether;
    uint256 public minWithdrawLockTime = 60; // 60 seconds

    address public auctionFeeVault;

    mapping(address => uint256) public depositBalances;
    mapping(address => WithdrawReserveInfo) public withdrawReservations;

    EnumerableSet.AddressSet private _depositAddrs;

    /* ========== MODIFIER ========== */

    modifier onlyEntryPoint() {
        /// @dev If the entry point is not registered, it will return zero address, anyway revert
        if (msg.sender != _getAuctionEntryPointAddress())
            revert OnlyEntryPoint();
        _;
    }

    modifier noWithdrawReservation(address searcher) {
        if (withdrawReservations[searcher].at != 0)
            revert WithdrawReservationExists();
        _;
    }

    /* ========== CONSTRUCTOR ========== */

    constructor(
        address initialOwner,
        address _auctionFeeVault
    ) Ownable(initialOwner) notNull(_auctionFeeVault) {
        auctionFeeVault = _auctionFeeVault;
    }

    /* ========== VAULT IMPLEMENTATION ========== */

    /// @dev Deposit KAIA into the vault
    function deposit()
        external
        payable
        override
        noWithdrawReservation(msg.sender)
    {
        _deposit(msg.sender, msg.value);
    }

    function depositFor(
        address searcher
    ) external payable override noWithdrawReservation(searcher) {
        _deposit(searcher, msg.value);
    }

    /// @dev Reserve a withdrawal
    function reserveWithdraw()
        external
        override
        noWithdrawReservation(msg.sender)
    {
        address searcher = msg.sender;
        uint256 amount = depositBalances[searcher];
        if (amount == 0) revert ZeroDepositAmount();

        withdrawReservations[searcher].at =
            block.timestamp +
            minWithdrawLockTime;
        withdrawReservations[searcher].amount = amount;

        emit VaultReserveWithdraw(searcher, amount, _getNonce(searcher));
    }

    /// @dev Withdraw KAIA from the vault
    function withdraw() external override {
        address searcher = msg.sender;
        uint256 amount = withdrawReservations[searcher].amount;
        uint256 withdrawableAt = withdrawReservations[searcher].at;
        if (withdrawableAt == 0 || block.timestamp <= withdrawableAt)
            revert WithdrawalNotAllowedYet();

        if (depositBalances[searcher] < amount) {
            // if reserved amount is less than current deposit balance due to slashing,
            // replace it with current deposit balance
            amount = depositBalances[searcher];
        }
        _decDeposit(searcher, amount);

        withdrawReservations[searcher].at = 0;
        withdrawReservations[searcher].amount = 0;

        (bool success, ) = searcher.call{value: amount}("");
        if (!success) revert WithdrawalFailed();

        emit VaultWithdraw(searcher, amount, _getNonce(searcher));
    }

    /// @dev Take a bid
    /// @param searcher The address of the searcher
    /// @param amount The amount of KAIA to take
    /// @return success Whether the bid was taken
    function takeBid(
        address searcher,
        uint256 amount
    ) external override onlyEntryPoint returns (bool) {
        uint256 depositAmount = depositBalances[searcher];
        if (amount > depositAmount) {
            // Return false if searcher has insufficient deposit
            emit InsufficientBalance(searcher, depositAmount, amount);
            return false;
        }

        return _sendBid(searcher, _decDeposit(searcher, amount));
    }

    /// @dev Take a gas reimbursement
    /// @param searcher The address of the searcher
    /// @param gasUsed The amount of gas used
    /// @return success Whether the gas reimbursement was taken
    function takeGas(
        address searcher,
        uint256 gasUsed
    ) external override onlyEntryPoint returns (bool) {
        uint256 gasAmount = tx.gasprice * gasUsed;
        uint256 depositAmount = depositBalances[searcher];
        if (gasAmount > depositAmount) {
            emit InsufficientBalance(searcher, depositAmount, gasAmount);
            return false;
        }

        return _sendGas(searcher, _decDeposit(searcher, gasAmount));
    }

    /* ========== VAULT MANAGEMENT ========== */

    /// @dev Change the auction fee vault
    /// @param newAuctionFeeVault The new auction fee vault
    function changeAuctionFeeVault(
        address newAuctionFeeVault
    ) external override onlyOwner notNull(newAuctionFeeVault) {
        emit ChangeAuctionFeeVault(auctionFeeVault, newAuctionFeeVault);
        auctionFeeVault = newAuctionFeeVault;
    }

    /// @dev Change the minimum deposit amount
    /// @param newMinAmount The new minimum deposit amount
    function changeMinDepositAmount(
        uint256 newMinAmount
    ) external override onlyOwner {
        emit ChangeMinDepositAmount(minDepositAmount, newMinAmount);
        minDepositAmount = newMinAmount;
    }

    /// @dev Change the minimum withdraw lock time
    /// @param newLocktime The new minimum withdraw lock time
    function changeMinWithdrawLocktime(
        uint256 newLocktime
    ) external override onlyOwner {
        emit ChangeMinWithdrawLocktime(minWithdrawLockTime, newLocktime);
        minWithdrawLockTime = newLocktime;
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _deposit(address searcher, uint256 amount) internal {
        if (amount == 0) revert ZeroDepositAmount();
        if (amount + depositBalances[searcher] < minDepositAmount)
            revert MinDepositNotOver();

        uint256 newBalance = _incDeposit(searcher, amount);

        emit VaultDeposit(searcher, amount, newBalance, _getNonce(searcher));
    }

    function _incDeposit(
        address searcher,
        uint256 amount
    ) internal returns (uint256) {
        uint256 newBalance = depositBalances[searcher] + amount;
        depositBalances[searcher] = newBalance;
        _depositAddrs.add(searcher);
        return newBalance;
    }

    function _decDeposit(
        address searcher,
        uint256 amount
    ) internal returns (uint256) {
        uint256 spent = depositBalances[searcher] > amount
            ? amount
            : depositBalances[searcher];
        depositBalances[searcher] -= spent;
        if (depositBalances[searcher] == 0) {
            _depositAddrs.remove(searcher);
        }
        return spent;
    }

    function _sendBid(
        address searcher,
        uint256 amount
    ) internal returns (bool success) {
        bytes memory data = abi.encodeWithSelector(
            IAuctionFeeVault.takeBid.selector,
            searcher
        );
        (success, ) = auctionFeeVault.call{value: amount}(data);
        if (!success) {
            _incDeposit(searcher, amount);
            emit TakenBidFailed(searcher, amount);
        } else {
            emit TakenBid(searcher, amount);
        }
    }

    function _sendGas(
        address searcher,
        uint256 amount
    ) internal returns (bool success) {
        (success, ) = block.coinbase.call{value: amount}("");
        if (!success) {
            _incDeposit(searcher, amount);
            emit TakenGasFailed(searcher, amount);
        } else {
            emit TakenGas(searcher, amount);
        }
    }

    function _getAuctionEntryPointAddress() internal view returns (address) {
        return IRegistry(REGISTRY_ADDRESS).getActiveAddr("AuctionEntryPoint");
    }

    function _getNonce(address searcher) internal view returns (uint256) {
        return INonce(_getAuctionEntryPointAddress()).nonces(searcher);
    }

    /* ========== GETTERS ========== */

    function getDepositAddrsLength() external view override returns (uint256) {
        return _depositAddrs.length();
    }

    function getDepositAddrs(
        uint256 start,
        uint256 limit
    ) external view override returns (address[] memory searchers) {
        uint256 totalAddresses = _depositAddrs.length();
        uint256 end;
        if (limit == 0) {
            end = totalAddresses;
        } else {
            end = start + limit > totalAddresses
                ? totalAddresses
                : start + limit;
        }

        searchers = new address[](end - start);
        for (uint256 i = start; i < end; i++) {
            searchers[i - start] = _depositAddrs.at(i);
        }
    }

    function isMinDepositOver(
        address searcher
    ) public view override returns (bool) {
        return depositBalances[searcher] >= minDepositAmount;
    }

    function getAllAddrsOverMinDeposit(
        uint256 start,
        uint256 limit
    )
        public
        view
        override
        returns (
            address[] memory searchers,
            uint256[] memory depositAmounts,
            uint256[] memory nonces
        )
    {
        uint256 totalAddresses = _depositAddrs.length();
        uint256 end;
        if (limit == 0) {
            end = totalAddresses;
        } else {
            end = start + limit > totalAddresses
                ? totalAddresses
                : start + limit;
        }

        INonce entryPoint = INonce(_getAuctionEntryPointAddress());

        searchers = new address[](end - start);
        depositAmounts = new uint256[](end - start);
        nonces = new uint256[](end - start);
        uint256 cnt = 0;

        for (uint256 i = start; i < end; i++) {
            address searcher = _depositAddrs.at(i);
            if (isMinDepositOver(searcher)) {
                searchers[cnt] = searcher;
                depositAmounts[cnt] = depositBalances[searcher];
                nonces[cnt] = entryPoint.nonces(searcher);
                cnt++;
            }
        }

        assembly {
            mstore(searchers, cnt)
            mstore(depositAmounts, cnt)
            mstore(nonces, cnt)
        }
    }
}

interface INonce {
    function nonces(address searcher) external view returns (uint256);
}

interface IRegistry {
    function getActiveAddr(string memory name) external view returns (address);
}
