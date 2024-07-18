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
pragma solidity 0.8.24;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "./IGuardian.sol";

contract Guardian is Initializable, ReentrancyGuardUpgradeable, UUPSUpgradeable, IERC165, IGuardian {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() { _disableInitializers(); }

    function initialize(
        address[] calldata initGuardians,
        uint8 _minGaurdianRequiredConfirm
    )
        public
        initializer
    {

        for (uint8 i=0; i<initGuardians.length; i++) {
            guardians.push(initGuardians[i]);
            isGuardian[initGuardians[i]] = true;
        }
        minGuardianRequiredConfirm = _minGaurdianRequiredConfirm;

        // Fill the first index(0) with dummy tx
        addTransaction(address(0), "", 0);

        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    function _authorizeUpgrade(address newImplementation) internal virtual override onlyGuardian {}

    function supportsInterface(bytes4 interfaceId) external override pure returns (bool) {
        return interfaceId == type(IGuardian).interfaceId;
    }

    /// @dev See {IGuardian-addGuardian}
    function addGuardian(address guardian)
        public
        override
        onlyGuardian
        guardianDoesNotExist(guardian)
        notNull(guardian)
    {
        isGuardian[guardian] = true;
        guardians.push(guardian);
        emit GuardianAddition(guardian);
    }

    /// @dev See {IGuardian-removeGuardian}
    function removeGuardian(address guardian)
        public
        override
        onlyGuardian
        guardianExists(guardian)
    {
        require(guardians.length > 1, "KAIA::Guardian: Guardian size must be greater than one to remove a guardian");
        isGuardian[guardian] = false;
        for (uint256 i=0; i<guardians.length - 1; i++)
            if (guardians[i] == guardian) {
                guardians[i] = guardians[guardians.length - 1];
                break;
            }
        guardians.pop();
        if (minGuardianRequiredConfirm > guardians.length) {
            changeRequirement(uint8(guardians.length));
        }
        emit GuardianRemoval(guardian);
    }

    /// @dev See {IGuardian-replaceGuardian}
    function replaceGuardian(address guardian, address newGuardian)
        public
        override
        onlyGuardian
        guardianExists(guardian)
        guardianDoesNotExist(newGuardian)
        notNull(newGuardian)
    {
        for (uint256 i=0; i<guardians.length; i++) {
            if (guardians[i] == guardian) {
                guardians[i] = newGuardian;
                break;
            }
        }
        isGuardian[guardian] = false;
        isGuardian[newGuardian] = true;
        emit GuardianRemoval(guardian);
        emit GuardianAddition(newGuardian);
    }

    /// @dev See {IGuardian-changeRequirement}
    function changeRequirement(uint8 _minGuardianRequiredConfirm)
        public
        override
        onlyGuardian
        validRequirement(guardians.length, _minGuardianRequiredConfirm)
    {
        minGuardianRequiredConfirm = _minGuardianRequiredConfirm;
        emit RequirementChange(_minGuardianRequiredConfirm);
    }

    /// @dev See {IGuardian-submitTransaction}
    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        public
        override
        returns (uint256)
    {
        uint256 txID = addTransaction(to, data, uniqUserTxIndex);
        confirmTransaction(txID);
        return txID;
    }

    /// @dev See {IGuardian-confirmTransaction}
    function confirmTransaction(uint256 txID)
        public
        override
        txExists(txID)
        notConfirmed(txID, msg.sender)
    {
        confirmations[txID][msg.sender] = true;
        emit Confirmation(msg.sender, txID);
        executeTransaction(txID);
    }

    /// @dev See {IGuardian-revokeConfirmation}
    function revokeConfirmation(uint256 txID)
        public
        override
        guardianExists(msg.sender)
        confirmed(txID, msg.sender)
        notExecuted(txID)
    {
        confirmations[txID][msg.sender] = false;
        emit Revocation(msg.sender, txID);
    }

    /// @dev See {IGuardian-executeTransaction}
    function executeTransaction(uint256 txID)
        public
        override
        guardianExists(msg.sender)
        confirmed(txID, msg.sender)
    {
        // if transaction was already executed, silently return without revert
        if (transactions[txID].executed) {
            return;
        }
        if (isConfirmed(txID)) {
            Transaction storage targetTx = transactions[txID];
            if (predefinedExecute(targetTx.to, targetTx.data)) {
                targetTx.executed = true;
                emit Execution(txID);
            }
        }
    }

    /// @dev Calls predefined function
    /// @param data Calldata
    function predefinedExecute(address to, bytes memory data)
        private
        nonReentrant
        returns (bool)
    {
        // 1. execute transaction
        (bool success, bytes memory res) = to.call(data);
        if (!success) {
            if (res.length > 68) {
                assembly {
                    res := add(res, 0x04)
                }
                revert(abi.decode(res, (string)));
            }
            revert(string(abi.encodePacked(res)));
        }
        return success;
    }

    /// @dev Adds a new transaction to the transaction mapping, if transaction does not exist yet.
    /// @param data Transaction data payload.
    /// @return Returns transaction ID.
    function addTransaction(address to, bytes memory data, uint256 uniqUserTxIndex)
        internal
        returns (uint256)
    {
        transactions.push(Transaction({
            to: to,
            data: data,
            executed: false
        }));
        uint64 txID = uint64(transactions.length - 1);
        emit Submission(txID);

        if (uniqUserTxIndex != 0) {
            require(userIdx2TxID[uniqUserTxIndex] == 0, "KAIA::Guardian: Submission to txID exists");
            userIdx2TxID[uniqUserTxIndex] = txID;
        }
        return txID;
    }

    /// @dev See {IGuardian-isConfirmed}
    function isConfirmed(uint256 txID)
        public
        override
        view
        returns (bool)
    {
        uint256 count = 0;
        for (uint256 i=0; i<guardians.length; i++) {
            if (confirmations[txID][guardians[i]]) {
                count += 1;
            }
            if (count == minGuardianRequiredConfirm) {
                return true;
            }
        }
        return false;
    }

    /// @dev See {IGuardian-getConfirmationCount}
    function getConfirmationCount(uint256 txID) public override view returns (uint256) {
        uint256 count = 0;
        for (uint256 i=0; i<guardians.length; i++) {
            if (confirmations[txID][guardians[i]]) {
                count += 1;
            }
        }
        return count;
    }

    /// @dev See {IGuardian-getTransactionCount}
    function getTransactionCount(bool pending, bool executed) public override view returns (uint256, uint256) {
        uint256 pendingCnt = 0;
        uint256 executedCnt = 0;
        // Ignore the first dummy transaction
        for (uint i=1; i<transactions.length; i++) {
            if (pending && !transactions[i].executed) {
                pendingCnt++;
            }
            if (executed && transactions[i].executed) {
                executedCnt++;
            }
        }
        return (pendingCnt, executedCnt);
    }

    /// @dev See {IGuardian-getGuardians}
    function getGuardians() public override view returns (address[] memory) {
        return guardians;
    }

    /// @dev See {IGuardian-getConfirmations}
    function getConfirmations(uint256 txID) public override view returns (address[] memory) {
        address[] memory confirmationsTemp = new address[](guardians.length);
        uint256 count = 0;
        for (uint256 i=0; i<guardians.length; i++) {
            if (confirmations[txID][guardians[i]]) {
                confirmationsTemp[count++] = guardians[i];
            }
        }
        address[] memory _confirmations = new address[](count);
        for (uint256 i=0; i<count; i++) {
            _confirmations[i] = confirmationsTemp[i];
        }
        return _confirmations;
    }

    /// @dev See {IGuardian-getTransactionIds}
    function getTransactionIds(uint256 from, uint256 to, bool pending, bool executed)
        public
        override
        view
        returns (uint256[] memory, uint256[] memory)
    {
        require(to > from, "KAIA::Guardian: Invalid from and to");
        // Ignore the first dummy transaction
        if (from == 0) {
            from = 1;
        }
        if (to > transactions.length) {
            to = transactions.length;
        }

        uint256 n = to - from;
        uint256[] memory _pendingTxs = new uint[](n);
        uint256[] memory _executedTxs = new uint[](n);
        uint256 pendingCnt = 0;
        uint256 executedCnt = 0;

        for (uint256 i=from; i<to; i++) {
            if (pending && !transactions[i].executed) {
                _pendingTxs[pendingCnt++] = i;
            }
            if (executed && transactions[i].executed) {
                _executedTxs[executedCnt++] = i;
            }
        }

        uint256[] memory pendingTxs = new uint[](pendingCnt);
        uint256[] memory executedTxs = new uint[](executedCnt);
        for (uint256 i=0; i<pendingCnt; i++) {
            pendingTxs[i] = _pendingTxs[i];
        }
        for (uint256 i=0; i<executedCnt; i++) {
            executedTxs[i] = _executedTxs[i];
        }
        return (pendingTxs, executedTxs);
    }

    /// @dev Return a contract version
    function getVersion() public pure returns (string memory) {
        return "0.0.1";
    }
}
