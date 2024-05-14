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

abstract contract IGuardian {
    //////////////////// Struct ////////////////////
    struct Transaction {
        address to;
        bytes data;
        bool executed;
    }

    //////////////////// Modifier ////////////////////
    modifier onlyGuardian() {
        require(msg.sender == address(this), "KAIA::Guardian: Sender is not guardian wallet");
        _;
    }

    modifier guardianDoesNotExist(address guardian) {
        require(!isGuardian[guardian], "KAIA::Guardian: The address must not be guardian");
        _;
    }

    modifier guardianExists(address guardian) {
        require(isGuardian[guardian], "KAIA::Guardian: Not an guardian");
        _;
    }

    modifier notNull(address addr) {
        require(addr != address(0), "KAIA::Guardian: A zero address is not allowed");
        _;
    }

    modifier txExists(uint256 txIdx) {
        require(txIdx < transactions.length, "KAIA::Guardian: Transaction does not exist");
        _;
    }

    modifier confirmed(uint256 txID, address guardian) {
        require(confirmations[txID][guardian], "KAIA::Guardian: No confirmation was committed yet");
        _;
    }

    modifier notExecuted(uint256 txID) {
        require(!transactions[txID].executed, "KAIA::Guardian: Transaction was already executed");
        _;
    }

    modifier notConfirmed(uint256 txID, address guardian) {
        require(!confirmations[txID][guardian], "KAIA::Guardian: Transaction was already confirmed");
        _;
    }

    modifier validRequirement(uint256 guardianCount, uint256 _minGuardianRequiredConfirm) {
        require(_minGuardianRequiredConfirm <= guardianCount
            && _minGuardianRequiredConfirm != 0
            && guardianCount != 0);
        _;
    }

    //////////////////// Event ////////////////////
    event Confirmation(address indexed sender, uint256 indexed transactionId);
    event Revocation(address indexed sender, uint256 indexed transactionId);
    event Submission(uint256 indexed transactionId);
    event Execution(uint256 indexed transactionId);
    event GuardianAddition(address indexed guardian);
    event GuardianRemoval(address indexed guardian);
    event RequirementChange(uint256 indexed required);

    //////////////////// Exported functions ////////////////////
    /// @dev Allows to add a new guardian. Transaction has to be sent by contract.
    /// @param guardian Address of new guardian.
    function addGuardian(address guardian) external virtual;

    /// @dev Allows to remove an guardian. Transaction has to be sent by wallet.
    /// @param guardian Address of guardian.
    function removeGuardian(address guardian) external virtual;

    /// @dev Allows to replace an guardian with a new guardian. Transaction has to be sent by wallet.
    /// @param guardian Address of guardian to be replaced.
    /// @param newGuardian Address of new guardian.
    function replaceGuardian(address guardian, address newGuardian) external virtual;

    /// @dev Allows to change the number of required confirmations. Transaction has to be sent by wallet.
    /// @param _minGuardianRequiredConfirm Number of required confirmations.
    function changeRequirement(uint8 _minGuardianRequiredConfirm) external virtual;

    /// @dev Allows an guardian to submit and confirm a transaction.
    /// @param to target contract address
    /// @param data Transaction data payload.
    /// @param uniqUserTxIndex unique user input number to be mapped with tx ID
    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        external
        virtual
        returns (uint256);

    /// @dev Allows an guardian to confirm a transaction.
    /// @param txID Transaction ID.
    function confirmTransaction(uint256 txID) external virtual;

    /// @dev Allows an guardian to revoke a confirmation for a transaction.
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

    /// @dev Returns list of guardians.
    /// @return List of guardian addresses.
    function getGuardians() external virtual view returns (address[] memory);

    /// @dev Returns array with guardian addresses, which confirmed transaction.
    /// @param txID Transaction ID.
    /// @return Returns array of guardian addresses.
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
    address[] public guardians;
    uint8 public minGuardianRequiredConfirm;
    mapping (address => bool) public isGuardian;
    Transaction[] public transactions;
    mapping (uint256 => mapping (address => bool)) public confirmations;
    mapping (uint256 => uint64) public submission2TxID; // <unique user input number, tx ID>

    uint256[100] __gap;
}
