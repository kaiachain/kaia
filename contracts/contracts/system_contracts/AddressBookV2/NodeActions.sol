// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {AddressBookV2Base} from "./AddressBookV2Base.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {ICnStaking} from "../CnStaking/CnStakingV4/interfaces/ICnStaking.sol";
import {State, NodeInfo, BlsPublicKeyInfo} from "../types/Node.sol";
import {SlotMath} from "../libraries/SlotMath.sol";
import {NodeVerifier} from "../libraries/NodeVerifier.sol";

/// @title NodeActions
/// @notice Node lifecycle operations: user-triggered (createNode, pause, exit, ...)
///         and system-triggered (processSystemTransition, updateScores).
abstract contract NodeActions is AddressBookV2Base {
    /* ========== CONSTRUCTOR ========== */

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(uint256 _epochBlockInterval) AddressBookV2Base(_epochBlockInterval) {}

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
        if (_getNodeState(nodeId) != State.Unknown) revert NodeAlreadyExists();

        ABv2Storage storage $ = _getStorage();
        NodeVerifier.registerNode($.registeredAddresses, nodeId, stakingContract, rewardAddress, blsInfo);

        NodeInfo memory info = NodeInfo({
            manager: msg.sender,
            stakingContract: stakingContract,
            rewardAddress: rewardAddress,
            voterAddress: voterAddress,
            timeoutAt: 0,
            gcId: 0,
            blsInfo: blsInfo,
            metadata: metadata,
            state: State.CandInactive
        });

        _createNode(nodeId, info);
    }

    /// @inheritdoc IAddressBookV2
    function deleteNode(address nodeId) external onlyManager(nodeId) {
        if (!_isNodeAtState(nodeId, State.CandInactive)) revert InvalidState();

        ABv2Storage storage $ = _getStorage();
        NodeInfo storage info = _getNodeInfo(nodeId);
        NodeVerifier.unregisterAddresses($.registeredAddresses, nodeId, info.stakingContract, info.rewardAddress);

        _deleteNode(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function updateManager(address nodeId, address newManager) external onlyManager(nodeId) {
        if (newManager == address(0)) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        address oldManager = $.nodeInfo[nodeId].manager;
        $.nodeInfo[nodeId].manager = newManager;
        emit ManagerUpdated(nodeId, oldManager, newManager);
    }

    /// @inheritdoc IAddressBookV2
    function updateRewardAddress(address nodeId, address newRewardAddress) external onlyManager(nodeId) {
        if (newRewardAddress == address(0)) revert InvalidInput();

        ABv2Storage storage $ = _getStorage();
        NodeInfo storage info = $.nodeInfo[nodeId];

        // Reward address is immutable when public delegation is enabled
        if (ICnStaking(payable(info.stakingContract)).publicDelegation() != address(0)) revert PDEnabled();
        if ($.registeredAddresses[newRewardAddress]) revert NodeVerifier.AddressAlreadyRegistered();

        address oldRewardAddress = info.rewardAddress;

        // Swap reward address
        $.registeredAddresses[oldRewardAddress] = false;
        $.registeredAddresses[newRewardAddress] = true;

        info.rewardAddress = newRewardAddress;
        emit RewardAddressUpdated(nodeId, oldRewardAddress, newRewardAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateVoterAddress(address nodeId, address newVoterAddress) external onlyManager(nodeId) {
        ABv2Storage storage $ = _getStorage();
        address oldVoterAddress = $.nodeInfo[nodeId].voterAddress;
        $.nodeInfo[nodeId].voterAddress = newVoterAddress;
        _refreshVoter($.nodeInfo[nodeId].stakingContract);
        emit VoterAddressUpdated(nodeId, oldVoterAddress, newVoterAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateMetadata(address nodeId, string calldata newMetadata) external onlyManager(nodeId) {
        _getStorage().nodeInfo[nodeId].metadata = newMetadata;
        emit MetadataUpdated(nodeId, newMetadata);
    }

    /// @inheritdoc IAddressBookV2
    function readyCandidate(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.CandInactive)) revert InvalidState();
        if (!_isNodeOverMinStake(nodeId)) revert StakingTooLow();

        ABv2Storage storage $ = _getStorage();
        if (_getStateCount(State.CandReady) >= $.maxReadyCandidateCount) revert SlotsFull();
        if (_getActiveSetLength() >= $.maxValidatorCount) revert SlotsFull();

        _transition(nodeId, State.CandReady, 0);
        emit CandidateReadied(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function unreadyCandidate(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.CandReady)) revert InvalidState();

        _transition(nodeId, State.CandInactive, 0);
        emit CandidateUnreadied(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function readyValidator(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.ValInactive)) revert InvalidState();
        if (!_isNodeOverMinStake(nodeId)) revert StakingTooLow();

        _transition(nodeId, State.ValReady, _getTimeoutAt(nodeId));
    }

    /// @inheritdoc IAddressBookV2
    function unreadyValidator(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.ValReady)) revert InvalidState();

        _transition(nodeId, State.ValInactive, _getTimeoutAt(nodeId));
    }

    /// @inheritdoc IAddressBookV2
    function pause(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.ValActive)) revert InvalidState();

        ABv2Storage storage $ = _getStorage();
        if (_getStateCount(State.ValPaused) >= SlotMath.maxSlotAvailable($.epochValCount)) revert SlotsFull();
        if (_getStateCount(State.ValActive) <= SlotMath.minActiveCount($.epochValCount)) revert SlotsFull();

        uint256 timeout = block.timestamp + $.pauseTimeout;
        _transition(nodeId, State.ValPaused, timeout);
    }

    /// @inheritdoc IAddressBookV2
    function resume(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.ValPaused)) revert InvalidState();

        _revertIfTimeoutExpired(nodeId);

        _transition(nodeId, State.ValActive, 0);
    }

    /// @inheritdoc IAddressBookV2
    function exit(address nodeId) external onlyNodeId(nodeId) {
        State state = _getNodeState(nodeId);
        if (state == State.ValPaused) {
            _revertIfTimeoutExpired(nodeId);
        } else if (state != State.ValActive) {
            revert InvalidState();
        }

        ABv2Storage storage $ = _getStorage();
        if (_getStateCount(State.ValExiting) >= SlotMath.maxSlotAvailable($.epochValCount)) revert SlotsFull();
        if (state == State.ValActive) {
            if (_getStateCount(State.ValActive) <= SlotMath.minActiveCount($.epochValCount)) revert SlotsFull();
        }

        _transition(nodeId, State.ValExiting, 0);
    }

    /// @inheritdoc IAddressBookV2
    function offboard(address nodeId) external onlyNodeId(nodeId) {
        if (!_isNodeAtState(nodeId, State.ValInactive)) revert InvalidState();
        _transition(nodeId, State.CandInactive, 0);
    }

    /* ========== SYSTEM ACTIONS ========== */

    /// @inheritdoc IAddressBookV2
    function processSystemTransition(
        address[] calldata nodeIds,
        State[] calldata newStates,
        uint256[] calldata timeoutAts
    ) external onlySystemTx {
        uint256 len = nodeIds.length;
        if (len != newStates.length || len != timeoutAts.length) revert InvalidInput();
        if (len > 0) {
            _batchTransition(nodeIds, newStates, timeoutAts);
        }
        if (_isEpochBlock()) {
            uint256 epochValCount = _getStateCount(State.ValActive) + _getStateCount(State.ValPaused);
            _getStorage().epochValCount = epochValCount;
            emit EpochTransitionProcessed(epochValCount);
        }
        emit SystemTransitionProcessed(nodeIds, newStates);
    }

    /// @inheritdoc IAddressBookV2
    function updateScores(address[] calldata nodeIds, uint256[] calldata scores) external onlySystemTx onlyEpochBlock {
        if (nodeIds.length != scores.length) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 epoch = _currentEpoch();
        for (uint256 i; i < nodeIds.length; ++i) {
            $.scores[epoch][nodeIds[i]] = scores[i];
        }
        emit ScoresUpdated(epoch, nodeIds, scores);
    }
}
