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
import "../../system_contracts/kaiabridge/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "../../system_contracts/kaiabridge/IGuardian.sol";
import "../../system_contracts/kaiabridge/IJudge.sol";

contract NewJudge is Initializable, ReentrancyGuardUpgradeable, UUPSUpgradeable, IERC165, IJudge {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() { _disableInitializers(); }

    function initialize(
        address[] calldata initJudges,
        address initGuardian,
        uint8 _minJudgeRequiredConfirm
    )
        public
        initializer
    {
        require(IERC165(initGuardian).supportsInterface(type(IGuardian).interfaceId), "KAIA::Judge: Guardian contract address does not implement IGuardian");

        for (uint8 i=0; i<initJudges.length; i++) {
            judges.push(initJudges[i]);
            isJudge[initJudges[i]] = true;
        }
        minJudgeRequiredConfirm = _minJudgeRequiredConfirm;
        guardian = initGuardian;

        // Fill the first index(0) with dummy tx
        addTransaction(address(0), "", 0);

        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    function _authorizeUpgrade(address newImplementation) internal virtual override onlyGuardian {}

    function supportsInterface(bytes4 interfaceId) external override pure returns (bool) {
        return interfaceId == type(IJudge).interfaceId;
    }

    /// @dev See {IJudge-addJudge}
    function addJudge(address judge)
        public
        override
        onlyGuardian
        judgeDoesNotExist(judge)
        notNull(judge)
    {
        isJudge[judge] = true;
        judges.push(judge);
        emit JudgeAddition(judge);
    }

    /// @dev See {IJudge-removeJudge}
    function removeJudge(address judge)
        public
        override
        onlyGuardian
        judgeExists(judge)
    {
        isJudge[judge] = false;
        for (uint256 i=0; i<judges.length - 1; i++)
            if (judges[i] == judge) {
                judges[i] = judges[judges.length - 1];
                break;
            }
        judges.pop();
        if (minJudgeRequiredConfirm > judges.length) {
            changeRequirement(uint8(judges.length));
        }
        emit JudgeRemoval(judge);
    }

    /// @dev See {IJudge-replaceJudge}
    function replaceJudge(address judge, address newJudge)
        public
        override
        onlyGuardian
        judgeExists(judge)
        judgeDoesNotExist(newJudge)
    {
        for (uint256 i=0; i<judges.length; i++) {
            if (judges[i] == judge) {
                judges[i] = newJudge;
                break;
            }
        }
        isJudge[judge] = false;
        isJudge[newJudge] = true;
        emit JudgeRemoval(judge);
        emit JudgeAddition(newJudge);
    }

    /// @dev See {IJudge-changeGuardian}
    function changeGuardian(address newGuardian)
        public
        override
        onlyGuardian
    {
        emit ChangeGuardian(guardian, newGuardian);
        guardian = newGuardian;
    }

    /// @dev See {IJudge-changeRequirement}
    function changeRequirement(uint8 _minJudgeRequiredConfirm)
        public
        override
        onlyGuardian
        validRequirement(judges.length, _minJudgeRequiredConfirm)
    {
        minJudgeRequiredConfirm = _minJudgeRequiredConfirm;
        emit RequirementChange(_minJudgeRequiredConfirm);
    }

    /// @dev See {IJudge-submitTransaction}
    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        public
        override
        returns (uint256)
    {
        uint256 txID = addTransaction(to, data, uniqUserTxIndex);
        confirmTransaction(txID);
        return txID;
    }

    /// @dev See {IJudge-confirmTransaction}
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

    /// @dev See {IJudge-revokeConfirmation}
    function revokeConfirmation(uint256 txID)
        public
        override
        judgeExists(msg.sender)
        confirmed(txID, msg.sender)
        notExecuted(txID)
    {
        confirmations[txID][msg.sender] = false;
        emit Revocation(msg.sender, txID);
    }

    /// @dev See {IJudge-executeTransaction}
    function executeTransaction(uint256 txID)
        public
        override
        judgeExists(msg.sender)
        confirmed(txID, msg.sender)
        notExecuted(txID)
    {
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
            require(submission2TxID[uniqUserTxIndex] == 0, "KAIA::Operator: Submission to txID exists");
            submission2TxID[uniqUserTxIndex] = txID;
        }
        return txID;
    }

    /// @dev See {IJudge-isConfirmed}
    function isConfirmed(uint256 txID)
        public
        override
        view
        returns (bool)
    {
        uint256 count = 0;
        for (uint256 i=0; i<judges.length; i++) {
            if (confirmations[txID][judges[i]]) {
                count += 1;
            }
            if (count == minJudgeRequiredConfirm) {
                return true;
            }
        }
        return false;
    }

    /// @dev See {IJudge-getConfirmationCount}
    function getConfirmationCount(uint256 txID) public override view returns (uint256) {
        uint256 count = 0;
        for (uint256 i=0; i<judges.length; i++) {
            if (confirmations[txID][judges[i]]) {
                count += 1;
            }
        }
        return count;
    }

    /// @dev See {IJudge-getTransactionCount}
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

    /// @dev See {IJudge-getJudges}
    function getJudges() public override view returns (address[] memory) {
        return judges;
    }

    /// @dev See {IJudge-getConfirmations}
    function getConfirmations(uint256 txID) public override view returns (address[] memory) {
        address[] memory confirmationsTemp = new address[](judges.length);
        uint256 count = 0;
        for (uint256 i=0; i<judges.length; i++) {
            if (confirmations[txID][judges[i]]) {
                confirmationsTemp[count++] = judges[i];
            }
        }
        address[] memory _confirmations = new address[](count);
        for (uint256 i=0; i<count; i++) {
            _confirmations[i] = confirmationsTemp[i];
        }
        return _confirmations;
    }

    /// @dev See {IJudge-getTransactionIds}
    function getTransactionIds(uint256 from, uint256 to, bool pending, bool executed)
        public
        override
        view
        returns (uint256[] memory, uint256[] memory)
    {
        require(to > from, "KAIA::Judge: Invalid from and to");
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
        return "0.0.2";
    }
}
