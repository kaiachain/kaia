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

abstract contract IJudge {
    //////////////////// Struct ////////////////////
    struct Transaction {
        address to;
        bytes data;
        bool executed;
    }

    //////////////////// Modifier ////////////////////
    modifier onlyGuardian() {
        require(msg.sender == guardian, "KAIA::Judge: Sender is not guardian wallet");
        _;
    }

    modifier judgeDoesNotExist(address judge) {
        require(!isJudge[judge], "KAIA::Judge: The address must not be judge");
        _;
    }

    modifier judgeExists(address judge) {
        require(isJudge[judge], "KAIA::Judge: Not an judge");
        _;
    }

    modifier notNull(address addr) {
        require(addr != address(0), "KAIA::Judge: A zero address is not allowed");
        _;
    }

    modifier txExists(uint256 txIdx) {
        require(txIdx < transactions.length, "KAIA::Judge: Transaction does not exist");
        _;
    }

    modifier confirmed(uint256 txID, address judge) {
        require(confirmations[txID][judge], "KAIA::Judge: No confirmation was committed yet");
        _;
    }

    modifier notExecuted(uint256 txID) {
        require(!transactions[txID].executed, "KAIA::Judge: Transaction was already executed");
        _;
    }

    modifier notConfirmed(uint256 txID, address judge) {
        require(!confirmations[txID][judge], "KAIA::Judge: Transaction was already confirmed");
        _;
    }

    modifier validRequirement(uint256 judgeCount, uint256 _minJudgeRequiredConfirm) {
        require(_minJudgeRequiredConfirm <= judgeCount
            && _minJudgeRequiredConfirm != 0
            && judgeCount != 0);
        _;
    }

    //////////////////// Event ////////////////////
    event Confirmation(address indexed sender, uint256 indexed transactionId);
    event Revocation(address indexed sender, uint256 indexed transactionId);
    event Submission(uint256 indexed transactionId);
    event Execution(uint256 indexed transactionId);
    event JudgeAddition(address indexed judge);
    event JudgeRemoval(address indexed judge);
    event ChangeGuardian(address indexed beforeGuardian, address indexed newGuardian);
    event RequirementChange(uint256 indexed required);

    //////////////////// Exported functions ////////////////////
    /// @dev Allows to add a new judge. Transaction has to be sent by contract.
    /// @param judge Address of new judge.
    function addJudge(address judge) external virtual;

    /// @dev Allows to remove an judge. Transaction has to be sent by wallet.
    /// @param judge Address of judge.
    function removeJudge(address judge) external virtual;

    /// @dev Allows to replace an judge with a new judge. Transaction has to be sent by wallet.
    /// @param judge Address of judge to be replaced.
    /// @param newJudge Address of new judge.
    function replaceJudge(address judge, address newJudge) external virtual;

    /// @dev Allows to replace an guardian with a new guardian. Transaction has to be sent by guardian.
    /// @param newGuardian Guardian address
    function changeGuardian(address newGuardian) external virtual;

    /// @dev Allows to change the number of required confirmations. Transaction has to be sent by wallet.
    /// @param _minJudgeRequiredConfirm Number of required confirmations.
    function changeRequirement(uint8 _minJudgeRequiredConfirm) external virtual;

    /// @dev Allows an judge to submit and confirm a transaction.
    /// @param to target contract address
    /// @param data Transaction data payload.
    /// @param uniqUserTxIndex unique user input number to be mapped with tx ID
    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        external
        virtual
        returns (uint256);

    /// @dev Allows an judge to confirm a transaction.
    /// @param txID Transaction ID.
    function confirmTransaction(uint256 txID) external virtual;

    /// @dev Allows an judge to revoke a confirmation for a transaction.
    /// @param txID Transaction ID.
    function revokeConfirmation(uint256 txID) external virtual;

    /// @dev Allows anyone to execute a confirmed transaction.
    /// @param txID Transaction ID.
    function executeTransaction(uint256 txID) external virtual;

    /// @dev Returns the confirmation status of a transaction.
    /// @param txID Transaction ID.
    /// @return Confirmation status.
    function isConfirmed(uint256 txID) external virtual view returns (bool);

    /// @dev Returns number of confirmations of a transaction.
    /// @param txID Transaction ID.
    /// @return Number of confirmations.
    function getConfirmationCount(uint256 txID) external virtual view returns (uint256);

    /// @dev Returns total number of transactions after filers are applied.
    /// @param pending Include pending transactions.
    /// @param executed Include executed transactions.
    /// @return Total number of transactions after filters are applied.
    function getTransactionCount(bool pending, bool executed) external virtual view returns (uint256, uint256);

    /// @dev Returns list of judges.
    /// @return List of judge addresses.
    function getJudges() external virtual view returns (address[] memory);

    /// @dev Returns array with judge addresses, which confirmed transaction.
    /// @param txID Transaction ID.
    /// @return Returns array of judge addresses.
    function getConfirmations(uint256 txID) external virtual view returns (address[] memory);

    /// @dev Returns list of transaction IDs in defined range.
    /// @param from Index start position of transaction array.
    /// @param to Index end position of transaction array.
    /// @param pending Include pending transactions.
    /// @param executed Include executed transactions.
    /// @return Returns array of transaction IDs.
    function getTransactionIds(uint256 from, uint256 to, bool pending, bool executed)
        external
        virtual
        view
        returns (uint256[] memory, uint256[] memory);


    //////////////////// Storage variables ////////////////////
    address[] public judges;
    address public guardian;
    uint8 public minJudgeRequiredConfirm;
    mapping (address => bool) public isJudge;
    Transaction[] public transactions;
    mapping (uint256 => mapping (address => bool)) public confirmations;
    mapping (uint256 => uint64) public submission2TxID; // <unique user input number, tx ID>

    uint256[100] __gap;
}
