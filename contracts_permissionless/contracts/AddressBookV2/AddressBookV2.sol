// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {NodeActions} from "./NodeActions.sol";
import {AddressBookLegacy} from "./AddressBookLegacy.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {IABv2DataContract} from "./interfaces/IABv2DataContract.sol";
import {IRegistry} from "../system/IRegistry.sol";
import {ABv2ConfigLib} from "../libraries/ABv2ConfigLib.sol";
import {State, BlsPublicKeyInfo, NodeInfo, Profile} from "../types/Node.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

/// @title AddressBookV2
/// @notice The sole entry point for all node operations, deployed at 0x400.
///         Inherits all node actions from NodeActions and exposes all public getters.
contract AddressBookV2 is NodeActions, AddressBookLegacy {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTRUCTOR ========== */

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor(uint256 _epochBlockInterval) NodeActions(_epochBlockInterval) {}

    /* ========== INITIALIZER ========== */

    /// @notice Initializes ABv2 by reading all genesis data from ABv2DataContract.
    /// @dev Resolves the data contract via Registry(0x401).getActiveAddr("ABv2DataContract").
    ///      It assumes that the data contract is already initialized and executable.
    ///      All validation is done at data contract construction time — this function trusts the data.
    function initialize() public initializer {
        // Resolve data contract from system registry
        address dataContract = IRegistry(REGISTRY_ADDRESS).getActiveAddr("ABv2DataContract");
        if (dataContract == address(0)) revert NotInitializable();
        IABv2DataContract.InitData memory d = IABv2DataContract(dataContract).getInitData();
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
        for (uint256 i; i < len; ++i) {
            $.registeredAddresses[d.nodeIds[i]] = true;
            $.registeredAddresses[d.infos[i].stakingContract] = true;
            $.registeredAddresses[d.infos[i].rewardAddress] = true;
        }

        // Set initial active validators and set epoch count
        _setInitialActiveValidators(d.nodeIds, d.infos);

        // GC IDs for genesis validators come from data contract.
        // Start counter at 100 so post-init nodes get gcId >= 101.
        $.lastAssignedGCId = 100;

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
        emit PauseTimeoutUpdated(ABv2ConfigLib.PAUSE_TIMEOUT.updateUint(newPauseTimeout), newPauseTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateIdleTimeout(uint256 newIdleTimeout) external onlyOwner {
        emit IdleTimeoutUpdated(ABv2ConfigLib.IDLE_TIMEOUT.updateUint(newIdleTimeout), newIdleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxValidatorCount(uint256 newMaxValidatorCount) external onlyOwner {
        emit MaxValidatorCountUpdated(
            ABv2ConfigLib.MAX_VALIDATOR_COUNT.updateUint(newMaxValidatorCount),
            newMaxValidatorCount
        );
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxReadyCandidateCount(uint256 newMaxReadyCandidateCount) external onlyOwner {
        emit MaxReadyCandidateCountUpdated(
            ABv2ConfigLib.MAX_READY_CANDIDATE_COUNT.updateUint(newMaxReadyCandidateCount),
            newMaxReadyCandidateCount
        );
    }

    /// @inheritdoc IAddressBookV2
    function updateExitThreshold(uint256 newExitThreshold) external onlyOwner {
        emit ExitThresholdUpdated(ABv2ConfigLib.EXIT_THRESHOLD.updateUint(newExitThreshold), newExitThreshold);
    }

    /// @inheritdoc IAddressBookV2
    function updateKefAddress(address newKefAddress) external onlyOwner {
        emit KefAddressUpdated(ABv2ConfigLib.KEF_ADDRESS.updateAddress(newKefAddress), newKefAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKifAddress(address newKifAddress) external onlyOwner {
        emit KifAddressUpdated(ABv2ConfigLib.KIF_ADDRESS.updateAddress(newKifAddress), newKifAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKpfAddress(address newKpfAddress) external onlyOwner {
        emit KpfAddressUpdated(ABv2ConfigLib.KPF_ADDRESS.updateAddress(newKpfAddress), newKpfAddress);
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
        NodeInfo storage info = _getStorage().nodeInfo[nodeId];
        if (info.state == State.Unknown) revert NodeNotFound();
        return info;
    }

    /// @inheritdoc IAddressBookV2
    function getNodeInfos(address[] calldata nodeIds) external view override returns (NodeInfo[] memory infos) {
        ABv2Storage storage $ = _getStorage();
        uint256 len = nodeIds.length;
        infos = new NodeInfo[](len);
        for (uint256 i; i < len; ++i) {
            infos[i] = $.nodeInfo[nodeIds[i]];
        }
    }

    /// @inheritdoc IAddressBookV2
    function getAllProfiles() external view override returns (Profile[] memory profiles) {
        ABv2Storage storage $ = _getStorage();
        address[] memory nodeIds = _getNonSuspendedNodeIds();
        uint256 len = nodeIds.length;
        profiles = new Profile[](len);
        for (uint256 i; i < len; ++i) {
            address nid = nodeIds[i];
            NodeInfo storage info = $.nodeInfo[nid];
            profiles[i] = Profile(nid, info.stakingContract, info.rewardAddress, info.timeoutAt, info.state);
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
        nodeIdList = _getNonSuspendedNodeIds();
        uint256 len = nodeIdList.length;
        pubkeyList = new BlsPublicKeyInfo[](len);
        for (uint256 i; i < len; ++i) {
            pubkeyList[i] = $.nodeInfo[nodeIdList[i]].blsInfo;
        }
    }

    /// @inheritdoc IAddressBookV2
    function getNodeState(address nodeId) external view override returns (State) {
        return _getStorage().nodeInfo[nodeId].state;
    }

    /// @inheritdoc IAddressBookV2
    function getStateCount(State state) external view override returns (uint256) {
        return _getStorage().stateCount[state];
    }

    /// @inheritdoc IAddressBookV2
    function isInActiveSet(address nodeId) external view override returns (bool) {
        return _getStorage().activeSet.contains(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function getActiveSetLength() external view override returns (uint256) {
        return _getStorage().activeSet.length();
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

    /* ========== LEGACY COMPATIBILITY ========== */

    /// @notice Legacy getState() — returns admin list and requirement for backward compatibility.
    /// @dev Different selector than getState(address) so both coexist.
    ///      Returns ([owner], 1) to reflect single-owner governance.
    function getState() external view returns (address[] memory, uint256) {
        address[] memory admins = new address[](1);
        admins[0] = owner();
        return (admins, 1);
    }

    /// @inheritdoc AddressBookLegacy
    function _getAllAddressData()
        internal
        view
        override
        returns (
            address[] memory nodeIds,
            address[] memory stakingContracts,
            address[] memory rewardAddresses,
            address kifAddress,
            address kefAddress
        )
    {
        ABv2Storage storage $ = _getStorage();
        uint256 activeLen = $.activeSet.length();

        nodeIds = new address[](activeLen);
        stakingContracts = new address[](activeLen);
        rewardAddresses = new address[](activeLen);

        for (uint256 i; i < activeLen; ++i) {
            address nid = $.activeSet.at(i);
            NodeInfo storage info = $.nodeInfo[nid];
            nodeIds[i] = nid;
            stakingContracts[i] = info.stakingContract;
            rewardAddresses[i] = info.rewardAddress;
        }

        kifAddress = $.kifAddress;
        kefAddress = $.kefAddress;
    }

    /// @inheritdoc AddressBookLegacy
    function _getNodeData(
        address nodeId
    ) internal view override returns (address stakingContract, address rewardAddress, bool exists) {
        ABv2Storage storage $ = _getStorage();
        NodeInfo storage info = $.nodeInfo[nodeId];
        if (info.state == State.Unknown) {
            return (address(0), address(0), false);
        }
        return (info.stakingContract, info.rewardAddress, true);
    }

    /* ========== INTERNAL HELPERS ========== */

    /// @notice Returns all node IDs from activeSet, excluding suspended nodes.
    function _getNonSuspendedNodeIds() private view returns (address[] memory nodeIds) {
        ABv2Storage storage $ = _getStorage();
        uint256 activeLen = $.activeSet.length();

        nodeIds = new address[](activeLen);
        uint256 idx;

        for (uint256 i; i < activeLen; ++i) {
            address nid = $.activeSet.at(i);
            if (!$.suspendedSet.contains(nid)) {
                nodeIds[idx] = nid;
                unchecked {
                    ++idx;
                }
            }
        }

        assembly {
            mstore(nodeIds, idx)
        }
    }

    /* ========== UPGRADE FUNCTIONS ========== */

    function _authorizeUpgrade(address) internal override onlyOwner {}
}
