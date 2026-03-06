// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {AddressBookBase} from "./AddressBookBase.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {NodeState, NodeInfo, BlsPublicKeyInfo} from "../types/Node.sol";
import {SlotMath} from "../libraries/SlotMath.sol";
import {NodeVerifier} from "../libraries/NodeVerifier.sol";

/// @title NodeActions
/// @notice Node lifecycle operations: user-triggered (createNode, pause, exit, ...)
///         and system-triggered (processSystemTransition, updateScores, initializeValidators).
abstract contract NodeActions is AddressBookBase {
    /* ========== USER ACTIONS ========== */

    /// @inheritdoc IAddressBookV2
    function createNode(
        address nodeId,
        address stakingContract,
        address rewardAddress,
        address voterAddress,
        BlsPublicKeyInfo memory blsInfo,
        string memory metadata
    ) external {
        if (_getNodeState(nodeId) != NodeState.Unknown) revert NodeAlreadyExists();

        OperationStorage storage $ = _getOperationStorage();
        NodeVerifier.registerNode($.registeredAddresses, nodeId, stakingContract, rewardAddress, blsInfo);

        $.managers[nodeId] = msg.sender;

        NodeInfo memory info = NodeInfo({
            stakingContract: stakingContract,
            rewardAddress: rewardAddress,
            voterAddress: voterAddress,
            timeoutAt: 0,
            gcId: 0,
            blsInfo: blsInfo,
            metadata: metadata,
            state: NodeState.CandInactive
        });

        _createNode(nodeId, info);
    }

    /// @inheritdoc IAddressBookV2
    function deleteNode(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.CandInactive)) revert InvalidState();

        OperationStorage storage $ = _getOperationStorage();
        NodeInfo storage info = _getNodeInfo(nodeId);
        NodeVerifier.unregisterAddresses($.registeredAddresses, nodeId, info.stakingContract, info.rewardAddress);
        delete $.managers[nodeId];

        _deleteNode(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function updateManager(address nodeId, address newManager) external onlyManager(nodeId) {
        if (newManager == address(0)) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        address oldManager = $.managers[nodeId];
        $.managers[nodeId] = newManager;
        emit ManagerUpdated(nodeId, oldManager, newManager);
    }

    /// @inheritdoc IAddressBookV2
    function activateCandidate(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.CandInactive)) revert InvalidState();
        if (!_isNodeOverMinStake(nodeId)) revert StakingTooLow();

        OperationStorage storage $ = _getOperationStorage();
        if (_getStateCount(NodeState.CandReady) >= $.maxReadyCandidateCount) revert SlotsFull();
        if (_getActiveSetLength() >= $.maxValidatorCount) revert SlotsFull();

        _transition(nodeId, NodeState.CandReady, 0);
        emit CandidateActivated(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function deactivateCandidate(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.CandReady)) revert InvalidState();

        _transition(nodeId, NodeState.CandInactive, 0);
        emit CandidateDeactivated(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function readyValidator(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.ValInactive)) revert InvalidState();
        if (!_isNodeOverMinStake(nodeId)) revert StakingTooLow();

        _transition(nodeId, NodeState.ValReady, _getTimeoutAt(nodeId));
    }

    /// @inheritdoc IAddressBookV2
    function unreadyValidator(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.ValReady)) revert InvalidState();

        _transition(nodeId, NodeState.ValInactive, _getTimeoutAt(nodeId));
    }

    /// @inheritdoc IAddressBookV2
    function pause(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.ValActive)) revert InvalidState();

        OperationStorage storage $ = _getOperationStorage();
        if (_getStateCount(NodeState.ValPaused) >= SlotMath.maxSlotAvailable($.epochValCount)) revert SlotsFull();
        if (_getStateCount(NodeState.ValActive) <= SlotMath.minActiveCount($.epochValCount)) revert SlotsFull();

        uint256 timeout = block.timestamp + $.pauseTimeout;
        _transition(nodeId, NodeState.ValPaused, timeout);
    }

    /// @inheritdoc IAddressBookV2
    function resume(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.ValPaused)) revert InvalidState();

        _revertIfTimeoutExpired(nodeId);

        _transition(nodeId, NodeState.ValActive, 0);
    }

    /// @inheritdoc IAddressBookV2
    function exit(address nodeId) external onlyManager(nodeId) {
        NodeState state = _getNodeState(nodeId);
        if (state == NodeState.ValPaused) {
            _revertIfTimeoutExpired(nodeId);
        } else if (state != NodeState.ValActive) {
            revert InvalidState();
        }

        OperationStorage storage $ = _getOperationStorage();
        if (_getStateCount(NodeState.ValExiting) >= SlotMath.maxSlotAvailable($.epochValCount)) revert SlotsFull();
        if (state == NodeState.ValActive) {
            if (_getStateCount(NodeState.ValActive) <= SlotMath.minActiveCount($.epochValCount)) revert SlotsFull();
        }

        _transition(nodeId, NodeState.ValExiting, 0);
    }

    /// @inheritdoc IAddressBookV2
    function offboard(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, NodeState.ValInactive)) revert InvalidState();
        _transition(nodeId, NodeState.CandInactive, 0);
    }

    /* ========== SYSTEM ACTIONS ========== */

    /// @inheritdoc IAddressBookV2
    function initializeValidators(address[] calldata nodeIds, NodeInfo[] calldata infos) external onlySystemTx {
        OperationStorage storage $ = _getOperationStorage();
        if ($.validatorsInitialized) revert AlreadyInitialized();
        $.validatorsInitialized = true;

        uint256 len = nodeIds.length;
        if (len != infos.length) revert InvalidInput();

        for (uint256 i; i < len; ) {
            NodeVerifier.registerNode(
                $.registeredAddresses,
                nodeIds[i],
                infos[i].stakingContract,
                infos[i].rewardAddress,
                infos[i].blsInfo
            );

            $.managers[nodeIds[i]] = nodeIds[i];

            unchecked {
                ++i;
            }
        }

        _createActiveValidators(nodeIds, infos);

        $.epochValCount += len;

        emit ValidatorsInitialized(nodeIds);
    }

    /// @inheritdoc IAddressBookV2
    function processSystemTransition(
        address[] calldata nodeIds,
        NodeState[] calldata newStates,
        uint256[] calldata timeoutAts
    ) external onlySystemTx {
        uint256 len = nodeIds.length;
        if (len != newStates.length || len != timeoutAts.length) revert InvalidInput();
        if (len > 0) {
            _batchTransition(nodeIds, newStates, timeoutAts);
        }
        if (_isEpochBlock()) {
            uint256 epochValCount = _getStateCount(NodeState.ValActive) + _getStateCount(NodeState.ValPaused);
            _getOperationStorage().epochValCount = epochValCount;
            emit EpochTransitionProcessed(epochValCount);
        }
        emit SystemTransitionProcessed(nodeIds, newStates);
    }

    /// @inheritdoc IAddressBookV2
    function updateScores(address[] calldata nodeIds, uint256[] calldata scores) external onlySystemTx onlyEpochBlock {
        if (nodeIds.length != scores.length) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 epoch = _currentEpoch();
        for (uint256 i; i < nodeIds.length; ) {
            $.scores[epoch][nodeIds[i]] = scores[i];
            unchecked {
                ++i;
            }
        }
        emit ScoresUpdated(epoch, nodeIds, scores);
    }
}
