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
import "./IOperator.sol";
import "./IGuardian.sol";
import "./IBridge.sol";

contract Operator is Initializable, ReentrancyGuardUpgradeable, UUPSUpgradeable, IERC165, IOperator {
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() { _disableInitializers(); }

    function initialize(
        address[] calldata initOperators,
        address initGuardian,
        uint8 _minGaurdianRequiredConfirm
    )
        public
        initializer
    {
        require(IERC165(initGuardian).supportsInterface(type(IGuardian).interfaceId), "KAIA::Operator: Operator contract address does not implement IGuardian");

        for (uint8 i=0; i<initOperators.length; i++) {
            operators.push(initOperators[i]);
            isOperator[initOperators[i]] = true;
        }
        minOperatorRequiredConfirm = _minGaurdianRequiredConfirm;
        guardian = initGuardian;

        // Fill the first index(0) with dummy tx
        addTransaction(address(0), "", 0);

        __UUPSUpgradeable_init();
        __ReentrancyGuard_init();
    }

    function _authorizeUpgrade(address newImplementation) internal virtual override onlyGuardian {}

    function supportsInterface(bytes4 interfaceId) external override pure returns (bool) {
        return interfaceId == type(IOperator).interfaceId;
    }

    /// @dev See {IOperator-addOperator}
    function addOperator(address operator)
        public
        override
        onlyGuardian
        operatorDoesNotExist(operator)
        notNull(operator)
    {
        isOperator[operator] = true;
        operators.push(operator);
        emit OperatorAddition(operator);
    }

    /// @dev See {IOperator-removeOperator}
    function removeOperator(address operator)
        public
        override
        onlyGuardian
        operatorExists(operator)
    {
        isOperator[operator] = false;
        for (uint64 i=0; i<operators.length - 1; i++)
            if (operators[i] == operator) {
                operators[i] = operators[operators.length - 1];
                break;
            }
        operators.pop();
        if (minOperatorRequiredConfirm > operators.length) {
            changeRequirement(uint8(operators.length));
        }
        emit OperatorRemoval(operator);
    }

    /// @dev See {IOperator-replaceOperator}
    function replaceOperator(address operator, address newOperator)
        public
        override
        onlyGuardian
        operatorExists(operator)
        operatorDoesNotExist(newOperator)
    {
        for (uint64 i=0; i<operators.length; i++) {
            if (operators[i] == operator) {
                operators[i] = newOperator;
                break;
            }
        }
        isOperator[operator] = false;
        isOperator[newOperator] = true;
        emit OperatorRemoval(operator);
        emit OperatorAddition(newOperator);
    }

    /// @dev See {IOperator-changeGuardian}
    function changeGuardian(address newGuardian)
        public
        override
        onlyGuardian
    {
        emit ChangeGuardian(guardian, newGuardian);
        guardian = newGuardian;
    }

    /// @dev See {IOperator-changeBridge}
    function changeBridge(address newBridge)
        public
        override
        onlyGuardian
    {
        emit ChangeBridge(bridge, newBridge);
        bridge = newBridge;
    }

    /// @dev See {IOperator-changeRequirement}
    function changeRequirement(uint8 _minOperatorRequiredConfirm)
        public
        override
        onlyGuardian
        validRequirement(operators.length, _minOperatorRequiredConfirm)
    {
        minOperatorRequiredConfirm = _minOperatorRequiredConfirm;
        emit RequirementChange(_minOperatorRequiredConfirm);
    }

    function submitTransaction(address to, bytes calldata data, uint256 uniqUserTxIndex)
        public
        override
        returns (uint64)
    {
        require(data.length >= 4, "Calldata length must be larger than 4bytes");
        uint64 seq = 0;
        uint64 txID = 0;

        try IBridge(bridge).bytes2Provision(data) returns (IBridge.ProvisionData memory provision) {
            require(bridge == to, "KAIA::Operator: Provision transaction must be targeted to known bridge contract address");
            require(provision.seq > 0 , "KAIA::Operator: Provision sequence number must be greater than zero");

            seq = provision.seq;
            bytes32 calldataHash = keccak256(data);
            txID = calldataHashes[calldataHash];

            if (txID != 0) {
                confirmTransaction(txID);
            } else {
                txID = addTransaction(bridge, data, uniqUserTxIndex);
                calldataHashes[calldataHash] = txID;
                provisions[txID] = provision;
                txID2Seq[txID] = seq;
                confirmTransaction(txID);
            }
            return txID;
        } catch {}

        txID = addTransaction(to, data, uniqUserTxIndex);
        confirmTransaction(txID);
        return txID;
    }

    /// @dev See {IOperator-confirmTransaction}
    function confirmTransaction(uint64 txID)
        public
        override
        txExists(txID)
        notConfirmed(txID, msg.sender)
    {
        // `provisions[txID]` is set only if the target contract address is bridge
        IBridge.ProvisionData storage provision = provisions[txID];
        if (provision.seq != 0) {
            EnumerableSetUint64.setAdd(seq2TxID[provision.seq], txID);
            updateGreatestSubmittedSeq(provision.seq);
            updateNextSeq(provision.seq);
            emit IBridge.Provision(IBridge.ProvisionIndividualEvent({
                seq: provision.seq,
                sender: provision.sender,
                receiver: provision.receiver,
                amount: provision.amount,
                txID: txID,
                operator: msg.sender
            }));
        }

        confirmations[txID][msg.sender] = true;
        emit Confirmation(msg.sender, txID);
        executeTransaction(txID);
    }

    /// @dev See {IOperator-revokeConfirmation}
    function revokeConfirmation(uint64 txID)
        public
        override
        operatorExists(msg.sender)
        confirmed(txID, msg.sender)
        notExecuted(txID)
    {
        confirmations[txID][msg.sender] = false;
        emit Revocation(msg.sender, txID);
    }

    /// @dev See {IOperator-executeTransaction}
    function executeTransaction(uint64 txID)
        public
        override
        operatorExists(msg.sender)
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
        returns (uint64)
    {
        transactions.push(Transaction({
            to: to,
            data: data,
            executed: false
        }));
        uint64 txID = uint64(transactions.length - 1);
        emit Submission(txID);

        if (uniqUserTxIndex != 0) {
            require(userIdx2TxID[uniqUserTxIndex] == 0, "KAIA::Operator: Submission to txID exists");
            userIdx2TxID[uniqUserTxIndex] = txID;
        }
        return txID;
    }

    /// @dev Update greatest sequence per operator
    /// @param seq ProvisionData sequence number
    function updateGreatestSubmittedSeq(uint64 seq) internal operatorExists(msg.sender) {
        if (greatestSubmittedSeq[msg.sender] < seq) {
            greatestSubmittedSeq[msg.sender] = seq;
        }
    }

    /// @dev Update next sequence per operator
    /// @param seq ProvisionData sequence number
    function updateNextSeq(uint64 seq) internal operatorExists(msg.sender) {
        if (seq > 0 && nextProvisionSeq[msg.sender] == seq - 1) {
            nextProvisionSeq[msg.sender] = seq;
        }
    }

    /// @dev See {IOperator-isConfirmed}
    function isConfirmed(uint64 txID)
        public
        override
        view
        returns (bool)
    {
        uint64 count = 0;
        for (uint64 i=0; i<operators.length; i++) {
            if (confirmations[txID][operators[i]]) {
                count += 1;
            }
            if (count == minOperatorRequiredConfirm) {
                return true;
            }
        }
        return false;
    }

    /// @dev See {IOperator-getConfirmationCount}
    function getConfirmationCount(uint64 txID) public override view returns (uint64) {
        uint64 count = 0;
        for (uint64 i=0; i<operators.length; i++) {
            if (confirmations[txID][operators[i]]) {
                count += 1;
            }
        }
        return count;
    }

    /// @dev See {IOperator-getTransactionCount}
    function getTransactionCount(bool pending, bool executed) public override view returns (uint64, uint64) {
        uint64 pendingCnt = 0;
        uint64 executedCnt = 0;
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

    /// @dev See {IOperator-getOperators}
    function getOperators() public override view returns (address[] memory) {
        return operators;
    }

    /// @dev See {IOperator-getConfirmations}
    function getConfirmations(uint64 txID) public override view returns (address[] memory) {
        address[] memory confirmationsTemp = new address[](operators.length);
        uint64 count = 0;
        for (uint64 i=0; i<operators.length; i++) {
            if (confirmations[txID][operators[i]]) {
                confirmationsTemp[count++] = operators[i];
            }
        }
        address[] memory _confirmations = new address[](count);
        for (uint64 i=0; i<count; i++) {
            _confirmations[i] = confirmationsTemp[i];
        }
        return _confirmations;
    }

    /// @dev See {IOperator-getTransactionIds}
    function getTransactionIds(uint64 from, uint64 to, bool pending, bool executed)
        public
        override
        view
        returns (uint64[] memory, uint64[] memory)
    {
        require(to > from, "KAIA::Operator: Invalid from and to");
        // Ignore the first dummy transaction
        if (from == 0) {
            from = 1;
        }
        if (to > transactions.length) {
            to = uint64(transactions.length);
        }

        uint64 n = uint64(to - from);
        uint64[] memory _pendingTxs = new uint64[](n);
        uint64[] memory _executedTxs = new uint64[](n);
        uint64 pendingCnt = 0;
        uint64 executedCnt = 0;

        for (uint64 i=from; i<to; i++) {
            if (pending && !transactions[i].executed) {
                _pendingTxs[pendingCnt++] = i;
            }
            if (executed && transactions[i].executed) {
                _executedTxs[executedCnt++] = i;
            }
        }

        uint64[] memory pendingTxs = new uint64[](pendingCnt);
        uint64[] memory executedTxs = new uint64[](executedCnt);
        for (uint64 i=0; i<pendingCnt; i++) {
            pendingTxs[i] = _pendingTxs[i];
        }
        for (uint64 i=0; i<executedCnt; i++) {
            executedTxs[i] = _executedTxs[i];
        }
        return (pendingTxs, executedTxs);
    }

    /// @dev See {IOperator-getUnconfirmedProvisionSeqs}
    function getUnconfirmedProvisionSeqs(address targetOperator, uint64 range)
        operatorExists(targetOperator)
        public
        override
        view
        returns (uint64[] memory)
    {
        uint64 seqFrom = unsubmittedNextSeq[targetOperator];
        return doGetUnconfirmedProvisionSeqs(targetOperator, seqFrom, seqFrom + range);
    }


    /// @dev See {IOperator-doGetUnconfirmedProvisionSeqs}
    function doGetUnconfirmedProvisionSeqs(address targetOperator, uint64 seqFrom, uint64 seqTo)
        operatorExists(targetOperator)
        public
        override
        view
        returns (uint64[] memory)
    {
        require(seqTo > seqFrom, "KAIA::Operator: Invalid from and to");
        // Ignore the first dummy sequence (sequence number starts from 1.)
        if (seqFrom == 0) {
            seqFrom = 1;
        }
        uint64 n = seqTo - seqFrom;
        uint64[] memory unsubmittedProvisionSeqsTemp = new uint64[](n);
        uint unsubmittedCnt = 0;

        for (uint64 seq=seqFrom; seq<seqTo; seq++) {
            uint64[] memory txIDs = EnumerableSetUint64.getAll(seq2TxID[seq]);
            uint64 unconfirmedCnt = 0;
            bool txExecuted = false;

            for (uint i=0; i<txIDs.length; i++) {
                Transaction storage txData = transactions[txIDs[i]];
                if (!confirmations[txIDs[i]][targetOperator]) {
                    unconfirmedCnt++;
                }
                if (txData.executed) {
                    txExecuted = true;
                    break;
                }
            }

            bool revoked = EnumerableSetUint64.setContains(revokedProvisionSeqs, seq);
            if (revoked) {
                // Condition0: If the sequence was revoked
                unsubmittedProvisionSeqsTemp[unsubmittedCnt++] = seq;
            } else if (!txExecuted && unconfirmedCnt == txIDs.length) {
                // Condition1: None of a list of provision txs were not executed
                // Condition2: Operator did not send among the list of provision txs
                unsubmittedProvisionSeqsTemp[unsubmittedCnt++] = seq;
            }
        }

        // fitting
        uint64[] memory unsubmittedProvisionSeqs = new uint64[](unsubmittedCnt);
        for (uint64 i=0; i<unsubmittedCnt; i++) {
            unsubmittedProvisionSeqs[i] = unsubmittedProvisionSeqsTemp[i];
        }
        return unsubmittedProvisionSeqs;
    }

    /// @dev See {IOperator-unmarkRevokeSeq}
    function unmarkRevokeSeq(uint64 seq) public override onlyBridge {
        EnumerableSetUint64.setRemove(revokedProvisionSeqs, seq);
    }

    /// @dev See {IOperator-markRevokeSeq}
    function markRevokeSeq(uint64 seq) public override onlyBridge {
        EnumerableSetUint64.setAdd(revokedProvisionSeqs, seq);
        emit RevokedProvision(seq);
    }

    /// @dev See {IOperator-updateNextUnsubmittedSeq}
    function updateNextUnsubmittedSeq(uint64 nextUnprovisionedSeq)
        public
        override
        operatorExists(msg.sender)
    {
        require(nextUnprovisionedSeq != 0, "KAIA::Operator: sequence number must not be zero");
        emit UnsubmittedNextSeqUpdate(unsubmittedNextSeq[msg.sender], nextUnprovisionedSeq);
        unsubmittedNextSeq[msg.sender] = nextUnprovisionedSeq;
    }

    /// @dev See {IOperator-getSeq2TxIDs}
    function getSeq2TxIDs(uint64 seq) public override view returns (uint64[] memory) {
        return EnumerableSetUint64.getAll(seq2TxID[seq]);
    }

    /// @dev See {IOperator-checkProvisionShouldSubmit}
    function checkProvisionShouldSubmit(bytes32 hashedData, address operator) public override view returns (bool) {
        uint64 txID = calldataHashes[hashedData];
        if (txID > 0 && txID < transactions.length) {
            bool executed = transactions[txID].executed;
            bool confirmed = confirmations[txID][operator];
            return !confirmed && !executed;
        }
        if (txID == 0) { // not submitted before for this payload
            return true;
        }
        return false;
    }

    /// @dev Return a contract version
    function getVersion() public pure returns (string memory) {
        return "0.0.1";
    }
}
