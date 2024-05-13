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

import "./IBridge.sol";
import "./EnumerableSetUint64.sol";

abstract contract IOperator {
    //////////////////// Struct ////////////////////
    struct Transaction {
        address to;
        bytes data;
        bool executed;
    }

    //////////////////// Modifier ////////////////////
    modifier onlyGuardian() {
        require(msg.sender == guardian, "PDT::Operator: Sender is not guardian contract");
        _;
    }

    modifier onlyBridge() {
        require(msg.sender == bridge, "PDT::Operator: Sender is not bridge contract");
        _;
    }

    modifier operatorDoesNotExist(address operator) {
        require(!isOperator[operator], "PDT::Operator: The address must not be operator");
        _;
    }

    modifier operatorExists(address operator) {
        require(isOperator[operator], "PDT::Operator: Not an operator");
        _;
    }

    modifier notNull(address addr) {
        require(addr != address(0), "PDT::Operator: A zero address is not allowed");
        _;
    }

    modifier txExists(uint64 txIdx) {
        require(txIdx < transactions.length, "PDT::Operator: Transaction does not exist");
        _;
    }

    modifier confirmed(uint64 txID, address operator) {
        require(confirmations[txID][operator], "PDT::Operator: No confirmation was committed yet");
        _;
    }

    modifier notExecuted(uint64 txID) {
        require(!transactions[txID].executed, "PDT::Operator: Transaction was already executed");
        _;
    }

    modifier notConfirmed(uint64 txID, address operator) {
        require(!confirmations[txID][operator], "PDT::Operator: Transaction was already confirmed");
        _;
    }

    modifier validRequirement(uint256 operatorCount, uint64 _minOperatorRequiredConfirm) {
        require(_minOperatorRequiredConfirm <= operatorCount
            && _minOperatorRequiredConfirm != 0
            && operatorCount != 0);
        _;
    }

    //////////////////// Event ////////////////////
    event Confirmation(address indexed sender, uint64 indexed transactionId);
    event Revocation(address indexed sender, uint64 indexed transactionId);
    event Submission(uint64 indexed transactionId);
    event Execution(uint64 indexed transactionId);
    event OperatorAddition(address indexed operator);
    event OperatorRemoval(address indexed operator);
    event ChangeGuardian(address indexed beforeGuardian, address indexed newGuardian);
    event ChangeBridge(address indexed beforeBridge, address indexed newBridge);
    event RequirementChange(uint64 indexed required);
    event UnsubmittedNextSeqUpdate(uint64 indexed unprovisionedNextSeq, uint64 indexed newSeq);
    event RevokedProvision(uint64 indexed seq);

    //////////////////// Exported functions ////////////////////
    /// @dev Allows to add a new operator. Transaction has to be sent by guardian.
    /// @param operator Address of new operator.
    function addOperator(address operator) external virtual;

    /// @dev Allows to remove an operator. Transaction has to be sent by guardian.
    /// @param operator Address of operator.
    function removeOperator(address operator) external virtual;

    /// @dev Allows to replace an operator with a new operator. Transaction has to be sent by guardian.
    /// @param operator Address of operator to be replaced.
    /// @param newOperator Address of new operator.
    function replaceOperator(address operator, address newOperator) external virtual;

    /// @dev Allows to replace an guardian with a new guardian. Transaction has to be sent by guardian.
    /// @param newGuardian Guardian address
    function changeGuardian(address newGuardian) external virtual;

    /// @dev Allows to replace an bridge with a new bridge. Transaction has to be sent by guardian.
    /// @param newBridge Bridge address
    function changeBridge(address newBridge) external virtual;

    /// @dev Allows to change the number of required confirmations. Transaction has to be sent by wallet.
    /// @param _minOperatorRequiredConfirm Number of required confirmations.
    function changeRequirement(uint8 _minOperatorRequiredConfirm) external virtual;

    /// @dev Allows an operator to submit and confirm a transaction.
    /// @param to target contract address
    /// @param data Transaction data payload.
    /// @param uniqUserTxIndex unique user input number to be mapped with tx ID
    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        external
        virtual
        returns (uint64);

    /// @dev Allows an operator to confirm a transaction.
    /// @param txID Transaction ID.
    function confirmTransaction(uint64 txID) external virtual;

    /// @dev Allows an operator to revoke a confirmation for a transaction.
    /// @param txID Transaction ID.
    function revokeConfirmation(uint64 txID) external virtual;

    /// @dev Allows anyone to execute a confirmed transaction.
    /// @param txID Transaction ID.
    function executeTransaction(uint64 txID) external virtual;

    /// @dev Returns the confirmation status of a transaction.
    /// @param txID Transaction ID.
    /// @return Confirmation status.
    function isConfirmed(uint64 txID) external virtual view returns (bool);

    /// @dev Returns number of confirmations of a transaction.
    /// @param txID Transaction ID.
    /// @return Number of confirmations.
    function getConfirmationCount(uint64 txID) external virtual view returns (uint64);

    /// @dev Returns total number of transactions after filers are applied.
    /// @param pending Include pending transactions.
    /// @param executed Include executed transactions.
    /// @return Total number of transactions after filters are applied.
    function getTransactionCount(bool pending, bool executed) external virtual view returns (uint64, uint64);

    /// @dev Returns list of operators.
    /// @return List of operator addresses.
    function getOperators() external virtual view returns (address[] memory);

    /// @dev Returns array with operator addresses, which confirmed transaction.
    /// @param txID Transaction ID.
    /// @return Returns array of operator addresses.
    function getConfirmations(uint64 txID) external virtual view returns (address[] memory);

    /// @dev Returns list of transaction IDs in defined range.
    /// @param from Index start position of transaction array.
    /// @param to Index end position of transaction array.
    /// @param pending Include pending transactions.
    /// @param executed Include executed transactions.
    /// @return Returns array of transaction IDs.
    function getTransactionIds(uint64 from, uint64 to, bool pending, bool executed)
        external
        virtual
        view
        returns (uint64[] memory, uint64[] memory);

    /// @dev Returns list of unprovisioned sequence and block number pairs for the given range of sequence
    /// @param targetOperator operator address
    /// @param range Unconfrimed provision set iteration range
    function getUnconfirmedProvisionSeqs(address targetOperator, uint64 range)
        external
        virtual
        view
        returns (uint64[] memory);

    /// @dev The same version of `getUnconfirmedProvisionSeqs()`, but it takes arbitrary sequence start and end number.
    /// @param targetOperator operator address
    /// @param seqFrom start sequence number
    /// @param seqTo end sequence number
    function doGetUnconfirmedProvisionSeqs(address targetOperator, uint64 seqFrom, uint64 seqTo)
        external
        virtual
        view
        returns (uint64[] memory);

    /// @dev update `unprovisionedNextSeq` value
    /// @param nextUnprovisionedSeq value to be replaced
    /// NOTE: This update function should be called once the getter(getUnprovisionedSeqs() or doGetUnprovisionedSeqs()) condition is satisfied.
    /// condition: Completness value (second value in pair) must be non-zero
    function updateNextUnsubmittedSeq(uint64 nextUnprovisionedSeq) external virtual;

    /// @dev Unmark revoke sign to the provision sequence number
    /// @param seq sequence number
    function unmarkRevokeSeq(uint64 seq) external virtual;

    /// @dev Mark revoke sign to the provision sequence number
    /// @param seq sequence number
    function markRevokeSeq(uint64 seq) external virtual;

    /// @dev Return tx ID list, corresponding the input sequence number
    /// @param seq Sequence number
    function getSeq2TxIDs(uint64 seq) external virtual view returns (uint64[] memory);


    //////////////////// Storage variables ////////////////////
    address[] public operators;
    address public guardian;
    address public bridge;
    uint8 public minOperatorRequiredConfirm;
    mapping (address => bool) public isOperator;
    Transaction[] public transactions;
    mapping (uint64 => mapping (address => bool)) public confirmations;

    mapping (uint64 => uint64) public txID2Seq; // <tx id, provision sequence>
    mapping (bytes32 => uint64) public calldataHashes;
    mapping (uint64 => IBridge.ProvisionData) public provisions; // <txID, provision info>
    mapping (address => uint64) public unprovisionedNextSeq;
    mapping (uint64 => EnumerableSet.UintSet) seq2TxID; // <sequence number, all provision transaction ID set>
    mapping (address => uint64) public unsubmittedNextSeq; // <operator address, next unsubmitted sequence number>
    mapping (address => uint64) public greatestSubmittedSeq; // Gratest submitted seq per operator
    EnumerableSet.UintSet revokedProvisionSeqs; // revoked provision sequence list
    mapping (uint256 => uint64) public submission2TxID; // <unique user input number, tx ID>

    uint256[100] __gap;
}
