// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {NodeActions} from "./NodeActions.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {IABv2DataContract} from "./interfaces/IABv2DataContract.sol";
import {IRegistry} from "../system/IRegistry.sol";
import {State, BlsPublicKeyInfo, NodeInfo, Profile} from "../types/Node.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @title AddressBookV2
/// @notice The sole entry point for all node operations, deployed at 0x400.
///         Inherits all node actions (user, system, admin) from NodeActions and
///         exposes all public getters in a single UUPS-upgradeable contract.
contract AddressBookV2 is NodeActions {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    string public constant CONTRACT_TYPE = "AddressBook";
    uint256 public constant VERSION = 2;

    /// @notice System registry address
    address internal constant REGISTRY = address(0x401);

    /* ========== CONSTRUCTOR ========== */

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /* ========== INITIALIZER ========== */

    /// @notice Initializes ABv2 by reading all genesis data from ABv2DataContract.
    /// @dev Resolves the data contract via Registry(0x401).getActiveAddr("ABv2DataContract").
    ///      It assumes that the data contract is already initialized and executable.
    ///      All validation is done at data contract construction time — this function trusts the data.
    function initialize() public initializer {
        // Resolve data contract from system registry
        address dataContract = IRegistry(REGISTRY).getActiveAddr("ABv2DataContract");
        if (dataContract == address(0)) revert NotInitializable();
        IABv2DataContract.InitData memory d = IABv2DataContract(dataContract).getData();
        ABv2Storage storage $ = _getStorage();

        // Set owner and config
        __Ownable_init(d.initialOwner);
        $.exitThreshold = d.exitThreshold;
        $.pauseTimeout = d.pauseTimeout;
        $.idleTimeout = d.idleTimeout;
        $.maxValidatorCount = d.maxValidatorCount;
        $.maxReadyCandidateCount = d.maxReadyCandidateCount;
        $.kefAddress = d.kefAddress;
        $.kifAddress = d.kifAddress;
        $.kpfAddress = d.kpfAddress;

        // Register addresses
        uint256 len = d.nodeIds.length;
        for (uint256 i; i < len; ) {
            $.registeredAddresses[d.nodeIds[i]] = true;
            $.registeredAddresses[d.infos[i].stakingContract] = true;
            $.registeredAddresses[d.infos[i].rewardAddress] = true;
            unchecked {
                ++i;
            }
        }

        // Set initial active validators and set epoch count
        _setInitialActiveValidators(d.nodeIds, d.infos);

        emit ValidatorsInitialized(d.nodeIds);
    }

    /* ========== CONFIGURATIONS  ========== */

    /// @inheritdoc IAddressBookV2
    function suspendValidator(address nodeId) external onlyOwner {
        ABv2Storage storage $ = _getStorage();
        if (!$.suspendedSet.add(nodeId)) revert AlreadySuspended();
        emit ValidatorSuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function unsuspendValidator(address nodeId) external onlyOwner {
        ABv2Storage storage $ = _getStorage();
        if (!$.suspendedSet.remove(nodeId)) revert NotSuspended();
        emit ValidatorUnsuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function updatePauseTimeout(uint256 newPauseTimeout) external onlyOwner {
        if (newPauseTimeout == 0) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 oldPauseTimeout = $.pauseTimeout;
        $.pauseTimeout = newPauseTimeout;
        emit PauseTimeoutUpdated(oldPauseTimeout, newPauseTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateIdleTimeout(uint256 newIdleTimeout) external onlyOwner {
        if (newIdleTimeout == 0) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 oldIdleTimeout = $.idleTimeout;
        $.idleTimeout = newIdleTimeout;
        emit IdleTimeoutUpdated(oldIdleTimeout, newIdleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxValidatorCount(uint256 newMaxValidatorCount) external onlyOwner {
        if (newMaxValidatorCount == 0) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 oldCount = $.maxValidatorCount;
        $.maxValidatorCount = newMaxValidatorCount;
        emit MaxValidatorCountUpdated(oldCount, newMaxValidatorCount);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxReadyCandidateCount(uint256 newMaxReadyCandidateCount) external onlyOwner {
        if (newMaxReadyCandidateCount == 0) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 oldCount = $.maxReadyCandidateCount;
        $.maxReadyCandidateCount = newMaxReadyCandidateCount;
        emit MaxReadyCandidateCountUpdated(oldCount, newMaxReadyCandidateCount);
    }

    /// @inheritdoc IAddressBookV2
    function updateExitThreshold(uint256 newExitThreshold) external onlyOwner {
        if (newExitThreshold == 0) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        uint256 oldThreshold = $.exitThreshold;
        $.exitThreshold = newExitThreshold;
        emit ExitThresholdUpdated(oldThreshold, newExitThreshold);
    }

    /// @inheritdoc IAddressBookV2
    function updateKefAddress(address newKefAddress) external onlyOwner {
        if (newKefAddress == address(0)) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        address oldAddress = $.kefAddress;
        $.kefAddress = newKefAddress;
        emit KefAddressUpdated(oldAddress, newKefAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKifAddress(address newKifAddress) external onlyOwner {
        if (newKifAddress == address(0)) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        address oldAddress = $.kifAddress;
        $.kifAddress = newKifAddress;
        emit KifAddressUpdated(oldAddress, newKifAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKpfAddress(address newKpfAddress) external onlyOwner {
        if (newKpfAddress == address(0)) revert InvalidInput();
        ABv2Storage storage $ = _getStorage();
        address oldAddress = $.kpfAddress;
        $.kpfAddress = newKpfAddress;
        emit KpfAddressUpdated(oldAddress, newKpfAddress);
    }

    /* ========== GETTERS ========== */

    /// @inheritdoc IAddressBookV2
    function getTimeouts() external view returns (uint256, uint256) {
        ABv2Storage storage $ = _getStorage();
        return ($.pauseTimeout, $.idleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function getMaxCounts() external view returns (uint256, uint256) {
        ABv2Storage storage $ = _getStorage();
        return ($.maxValidatorCount, $.maxReadyCandidateCount);
    }

    /// @inheritdoc IAddressBookV2
    function getExitThreshold() external view returns (uint256) {
        return _getStorage().exitThreshold;
    }

    /// @inheritdoc IAddressBookV2
    function getFundAddresses() external view returns (address, address, address) {
        ABv2Storage storage $ = _getStorage();
        return ($.kefAddress, $.kifAddress, $.kpfAddress);
    }

    /// @inheritdoc IAddressBookV2
    function getManager(address nodeId) external view returns (address) {
        return _getStorage().nodeInfo[nodeId].manager;
    }

    /// @inheritdoc IAddressBookV2
    function isRegistered(address addr) external view returns (bool) {
        return _getStorage().registeredAddresses[addr];
    }

    /// @inheritdoc IAddressBookV2
    function getScore(uint256 epoch, address nodeId) external view returns (uint256) {
        return _getStorage().scores[epoch][nodeId];
    }

    /// @inheritdoc IAddressBookV2
    function currentEpoch() public view returns (uint256) {
        return _currentEpoch();
    }

    /// @inheritdoc IAddressBookV2
    function getEpochValCount() external view returns (uint256) {
        return _getStorage().epochValCount;
    }

    /// @inheritdoc IAddressBookV2
    function getNodeInfo(address nodeId) external view override returns (NodeInfo memory) {
        ABv2Storage storage $ = _getStorage();
        if ($.nodeInfo[nodeId].state == State.Unknown) revert NodeNotFound();
        return $.nodeInfo[nodeId];
    }

    /// @inheritdoc IAddressBookV2
    function getNodeInfos(address[] calldata nodeIds) external view override returns (NodeInfo[] memory infos) {
        ABv2Storage storage $ = _getStorage();
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
        ABv2Storage storage $ = _getStorage();
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
        ABv2Storage storage $ = _getStorage();
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
    function getNodeState(address nodeId) external view override returns (State) {
        return _getStorage().nodeInfo[nodeId].state;
    }

    /// @inheritdoc IAddressBookV2
    function isCandInactive(address nodeId) external view override returns (bool) {
        return _getStorage().candInactiveSet.contains(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function isInActiveSet(address nodeId) external view override returns (bool) {
        return _getStorage().activeSet.contains(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function getSuspendedValidators() external view override returns (address[] memory) {
        return _getStorage().suspendedSet.values();
    }

    /// @inheritdoc IAddressBookV2
    function getStakingContract(address nodeId) external view override returns (address) {
        return _getStorage().nodeInfo[nodeId].stakingContract;
    }

    /// @inheritdoc IAddressBookV2
    function getTimeoutAt(address nodeId) external view override returns (uint256) {
        return _getStorage().nodeInfo[nodeId].timeoutAt;
    }

    /// @inheritdoc IAddressBookV2
    function getActiveSetLength() external view override returns (uint256) {
        return _getStorage().activeSet.length();
    }

    /// @inheritdoc IAddressBookV2
    function getCandInactiveSetLength() external view override returns (uint256) {
        return _getStorage().candInactiveSet.length();
    }

    /// @inheritdoc IAddressBookV2
    function getStateCount(State state) external view override returns (uint256) {
        return _getStorage().stateCount[state];
    }

    /* ========== UPGRADE FUNCTIONS ========== */

    function _authorizeUpgrade(address) internal override onlyOwner {}
}
