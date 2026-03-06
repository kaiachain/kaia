// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {NodeActions} from "./NodeActions.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {NodeState, BlsPublicKeyInfo, NodeInfo, Profile} from "../types/Node.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @title AddressBookV2
/// @notice The sole entry point for all node operations, deployed at 0x400.
///         Inherits all node actions (user, system, admin) from NodeActions and
///         exposes all public getters in a single UUPS-upgradeable contract.
/// @dev Uses two ERC-7201 namespaced storage slots:
///      - NodeStorage at 0x2a6c4... (node info, sets, state counts)
///      - OperationStorage at 0x9114b... (managers, scores, config)
contract AddressBookV2 is NodeActions {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "AddressBook";
    uint256 public constant VERSION = 2;

    /* ========== CONSTRUCTOR ========== */

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /* ========== INITIALIZER ========== */

    /// @notice Initializes the contract with owner and all configuration parameters
    /// @param initialOwner The owner address (can upgrade, suspend, update config)
    /// @param _exitThreshold Proposal failure count that triggers forced exit
    /// @param _pauseTimeout Duration (seconds) for ValPaused timeout
    /// @param _idleTimeout Duration (seconds) for ValInactive/ValReady timeout
    /// @param _maxValidatorCount Maximum validators in the pipeline (activeSet)
    /// @param _maxReadyCandidateCount Maximum candidates in CandReady state
    function initialize(
        address initialOwner,
        uint256 _exitThreshold,
        uint256 _pauseTimeout,
        uint256 _idleTimeout,
        uint256 _maxValidatorCount,
        uint256 _maxReadyCandidateCount
    ) public initializer {
        __Ownable_init();
        _transferOwnership(initialOwner);
        if (_exitThreshold == 0 || _pauseTimeout == 0 || _idleTimeout == 0) revert InvalidInput();
        if (_maxValidatorCount == 0 || _maxReadyCandidateCount == 0) revert InvalidInput();

        OperationStorage storage $ = _getOperationStorage();
        $.exitThreshold = _exitThreshold;
        $.pauseTimeout = _pauseTimeout;
        $.idleTimeout = _idleTimeout;
        $.maxValidatorCount = _maxValidatorCount;
        $.maxReadyCandidateCount = _maxReadyCandidateCount;
    }

    /* ========== CONFIGURATIONS  ========== */

    /// @inheritdoc IAddressBookV2
    function suspendValidator(address nodeId) external onlyOwner {
        NodeStorage storage $ = _getNodeStorage();
        if (!$.suspendedSet.add(nodeId)) revert AlreadySuspended();
        emit ValidatorSuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function unsuspendValidator(address nodeId) external onlyOwner {
        NodeStorage storage $ = _getNodeStorage();
        if (!$.suspendedSet.remove(nodeId)) revert NotSuspended();
        emit ValidatorUnsuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function updatePauseTimeout(uint256 newPauseTimeout) external onlyOwner {
        if (newPauseTimeout == 0) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 oldPauseTimeout = $.pauseTimeout;
        $.pauseTimeout = newPauseTimeout;
        emit PauseTimeoutUpdated(oldPauseTimeout, newPauseTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateIdleTimeout(uint256 newIdleTimeout) external onlyOwner {
        if (newIdleTimeout == 0) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 oldIdleTimeout = $.idleTimeout;
        $.idleTimeout = newIdleTimeout;
        emit IdleTimeoutUpdated(oldIdleTimeout, newIdleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxValidatorCount(uint256 newMaxValidatorCount) external onlyOwner {
        if (newMaxValidatorCount == 0) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 oldCount = $.maxValidatorCount;
        $.maxValidatorCount = newMaxValidatorCount;
        emit MaxValidatorCountUpdated(oldCount, newMaxValidatorCount);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxReadyCandidateCount(uint256 newMaxReadyCandidateCount) external onlyOwner {
        if (newMaxReadyCandidateCount == 0) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 oldCount = $.maxReadyCandidateCount;
        $.maxReadyCandidateCount = newMaxReadyCandidateCount;
        emit MaxReadyCandidateCountUpdated(oldCount, newMaxReadyCandidateCount);
    }

    /// @inheritdoc IAddressBookV2
    function updateExitThreshold(uint256 newExitThreshold) external onlyOwner {
        if (newExitThreshold == 0) revert InvalidInput();
        OperationStorage storage $ = _getOperationStorage();
        uint256 oldThreshold = $.exitThreshold;
        $.exitThreshold = newExitThreshold;
        emit ExitThresholdUpdated(oldThreshold, newExitThreshold);
    }

    /* ========== GETTERS ========== */

    /// @inheritdoc IAddressBookV2
    function getTimeouts() external view returns (uint256, uint256) {
        OperationStorage storage $ = _getOperationStorage();
        return ($.pauseTimeout, $.idleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function getMaxCounts() external view returns (uint256, uint256) {
        OperationStorage storage $ = _getOperationStorage();
        return ($.maxValidatorCount, $.maxReadyCandidateCount);
    }

    /// @inheritdoc IAddressBookV2
    function getExitThreshold() external view returns (uint256) {
        return _getOperationStorage().exitThreshold;
    }

    /// @inheritdoc IAddressBookV2
    function getManager(address nodeId) external view returns (address) {
        return _getOperationStorage().managers[nodeId];
    }

    /// @inheritdoc IAddressBookV2
    function isRegistered(address addr) external view returns (bool) {
        return _getOperationStorage().registeredAddresses[addr];
    }

    /// @inheritdoc IAddressBookV2
    function getScore(uint256 epoch, address nodeId) external view returns (uint256) {
        return _getOperationStorage().scores[epoch][nodeId];
    }

    /// @inheritdoc IAddressBookV2
    function currentEpoch() public view returns (uint256) {
        return _currentEpoch();
    }

    /// @inheritdoc IAddressBookV2
    function getEpochValCount() external view returns (uint256) {
        return _getOperationStorage().epochValCount;
    }

    /// @inheritdoc IAddressBookV2
    function isValidatorsInitialized() external view returns (bool) {
        return _getOperationStorage().validatorsInitialized;
    }

    /// @inheritdoc IAddressBookV2
    function getNodeInfo(address nodeId) external view override returns (NodeInfo memory) {
        NodeStorage storage $ = _getNodeStorage();
        if ($.nodeInfo[nodeId].state == NodeState.Unknown) revert NodeNotFound();
        return $.nodeInfo[nodeId];
    }

    /// @inheritdoc IAddressBookV2
    function getNodeInfos(address[] calldata nodeIds) external view override returns (NodeInfo[] memory infos) {
        NodeStorage storage $ = _getNodeStorage();
        uint256 len = nodeIds.length;
        infos = new NodeInfo[](len);
        for (uint256 i; i < len; ) {
            infos[i] = $.nodeInfo[nodeIds[i]];
            unchecked {
                ++i;
            }
        }
    }

    /// @inheritdoc IAddressBookV2
    function getAllProfiles() external view override returns (Profile[] memory profiles) {
        NodeStorage storage $ = _getNodeStorage();
        uint256 activeLen = $.activeSet.length();
        uint256 inactiveLen = $.candInactiveSet.length();

        profiles = new Profile[](activeLen + inactiveLen);
        uint256 idx;

        for (uint256 i; i < activeLen; ) {
            address nid = $.activeSet.at(i);
            if (!$.suspendedSet.contains(nid)) {
                NodeInfo storage info = $.nodeInfo[nid];
                profiles[idx] = Profile(nid, info.stakingContract, info.rewardAddress, info.timeoutAt, info.state);
                unchecked {
                    ++idx;
                }
            }
            unchecked {
                ++i;
            }
        }

        for (uint256 i; i < inactiveLen; ) {
            address nid = $.candInactiveSet.at(i);
            if (!$.suspendedSet.contains(nid)) {
                NodeInfo storage info = $.nodeInfo[nid];
                profiles[idx] = Profile(nid, info.stakingContract, info.rewardAddress, info.timeoutAt, info.state);
                unchecked {
                    ++idx;
                }
            }
            unchecked {
                ++i;
            }
        }

        assembly {
            mstore(profiles, idx)
        }
    }

    /// @inheritdoc IAddressBookV2
    function getAllBlsInfo()
        external
        view
        override
        returns (address[] memory nodeIdList, BlsPublicKeyInfo[] memory pubkeyList)
    {
        NodeStorage storage $ = _getNodeStorage();
        uint256 activeLen = $.activeSet.length();
        uint256 inactiveLen = $.candInactiveSet.length();

        nodeIdList = new address[](activeLen + inactiveLen);
        pubkeyList = new BlsPublicKeyInfo[](activeLen + inactiveLen);
        uint256 idx;

        for (uint256 i; i < activeLen; ) {
            address nid = $.activeSet.at(i);
            if (!$.suspendedSet.contains(nid)) {
                nodeIdList[idx] = nid;
                pubkeyList[idx] = $.nodeInfo[nid].blsInfo;
                unchecked {
                    ++idx;
                }
            }
            unchecked {
                ++i;
            }
        }

        for (uint256 i; i < inactiveLen; ) {
            address nid = $.candInactiveSet.at(i);
            if (!$.suspendedSet.contains(nid)) {
                nodeIdList[idx] = nid;
                pubkeyList[idx] = $.nodeInfo[nid].blsInfo;
                unchecked {
                    ++idx;
                }
            }
            unchecked {
                ++i;
            }
        }

        assembly {
            mstore(nodeIdList, idx)
            mstore(pubkeyList, idx)
        }
    }

    /// @inheritdoc IAddressBookV2
    function getNodeState(address nodeId) external view override returns (NodeState) {
        return _getNodeStorage().nodeInfo[nodeId].state;
    }

    /// @inheritdoc IAddressBookV2
    function isCandInactive(address nodeId) external view override returns (bool) {
        return _getNodeStorage().candInactiveSet.contains(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function isInActiveSet(address nodeId) external view override returns (bool) {
        return _getNodeStorage().activeSet.contains(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function getSuspendedValidators() external view override returns (address[] memory) {
        return _getNodeStorage().suspendedSet.values();
    }

    /// @inheritdoc IAddressBookV2
    function getStakingContract(address nodeId) external view override returns (address) {
        return _getNodeStorage().nodeInfo[nodeId].stakingContract;
    }

    /// @inheritdoc IAddressBookV2
    function getTimeoutAt(address nodeId) external view override returns (uint256) {
        return _getNodeStorage().nodeInfo[nodeId].timeoutAt;
    }

    /// @inheritdoc IAddressBookV2
    function getActiveSetLength() external view override returns (uint256) {
        return _getNodeStorage().activeSet.length();
    }

    /// @inheritdoc IAddressBookV2
    function getCandInactiveSetLength() external view override returns (uint256) {
        return _getNodeStorage().candInactiveSet.length();
    }

    /// @inheritdoc IAddressBookV2
    function getStateCount(NodeState state) external view override returns (uint256) {
        return _getNodeStorage().stateCount[state];
    }

    /* ========== UPGRADE FUNCTIONS ========== */

    function _authorizeUpgrade(address) internal override onlyOwner {}
}
