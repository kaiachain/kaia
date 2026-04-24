// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {NodeActions} from "./NodeActions.sol";
import {AddressBookLegacy} from "./AddressBookLegacy.sol";
import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {IABv2DataContract} from "./interfaces/IABv2DataContract.sol";
import {IRegistry} from "../system/IRegistry.sol";
import {ABv2ConfigLib} from "../libraries/ABv2ConfigLib.sol";
import {SlotMath} from "../libraries/SlotMath.sol";
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
        $.pfsThreshold = d.pfsThreshold;
        $.cfsThreshold = d.cfsThreshold;
        $.pauseTimeout = d.pauseTimeout;
        $.idleTimeout = d.idleTimeout;
        $.maxNodeCount = d.maxNodeCount;
        $.maxValActivePausedCount = d.maxValActivePausedCount;
        $.maxCandReadyCount = d.maxCandReadyCount;
        $.kefAddress = d.kefAddress;
        $.kifAddress = d.kifAddress;
        $.kpfAddress = d.kpfAddress;
        $.suspender = d.initialSuspender;
        $.configurator = d.initialConfigurator;

        // Register addresses
        uint256 len = d.nodeIds.length;
        for (uint256 i; i < len; ++i) {
            $.usedAddresses[d.nodeIds[i]] = true;
            $.usedAddresses[d.infos[i].stakingContract] = true;
            $.usedAddresses[d.infos[i].rewardAddress] = true;
        }

        // Set initial active validators and set epoch count
        _setInitialActiveValidators(d.nodeIds, d.infos);

        // Start counter at 100 so post-init nodes get gcId >= 101.
        $.lastAssignedGCId = 100;
    }

    /* ========== CONFIGURATIONS  ========== */

    /// @inheritdoc IAddressBookV2
    function suspendValidator(address nodeId) external onlySuspender {
        ABv2Storage storage $ = _getStorage();
        if (!$.suspendedSet.add(nodeId)) revert AlreadySuspended();
        emit ValidatorSuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function unsuspendValidator(address nodeId) external onlySuspender {
        ABv2Storage storage $ = _getStorage();
        if (!$.suspendedSet.remove(nodeId)) revert NotSuspended();
        emit ValidatorUnsuspended(nodeId);
    }

    /// @inheritdoc IAddressBookV2
    function updateSuspender(address newSuspender) external onlyOwner {
        emit SuspenderUpdated(ABv2ConfigLib.SUSPENDER.updateAddress(newSuspender), newSuspender);
    }

    /// @inheritdoc IAddressBookV2
    function updateConfigurator(address newConfigurator) external onlyOwner {
        emit ConfiguratorUpdated(ABv2ConfigLib.CONFIGURATOR.updateAddress(newConfigurator), newConfigurator);
    }

    /// @inheritdoc IAddressBookV2
    function assignGcId(address nodeId) external onlyConfigurator {
        ABv2Storage storage $ = _getStorage();
        if ($.nodeInfo[nodeId].state == State.Unknown) revert NodeNotFound();
        if ($.nodeInfo[nodeId].gcId != 0) revert GcIdAlreadyAssigned();
        uint256 gcId = ++$.lastAssignedGCId;
        $.nodeInfo[nodeId].gcId = gcId;
    }

    /// @inheritdoc IAddressBookV2
    function revokeGcId(address nodeId) external onlyConfigurator {
        ABv2Storage storage $ = _getStorage();
        uint256 gcId = $.nodeInfo[nodeId].gcId;
        if (gcId == 0) revert GcIdNotAssigned();
        $.nodeInfo[nodeId].gcId = 0;
    }

    /// @inheritdoc IAddressBookV2
    function updatePauseTimeout(uint256 newPauseTimeout) external onlyConfigurator {
        emit PauseTimeoutUpdated(ABv2ConfigLib.PAUSE_TIMEOUT.updateUint(newPauseTimeout), newPauseTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateIdleTimeout(uint256 newIdleTimeout) external onlyConfigurator {
        emit IdleTimeoutUpdated(ABv2ConfigLib.IDLE_TIMEOUT.updateUint(newIdleTimeout), newIdleTimeout);
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxNodeCount(uint256 newMaxNodeCount) external onlyConfigurator {
        emit MaxNodeCountUpdated(
            ABv2ConfigLib.MAX_NODE_COUNT.updateUint(newMaxNodeCount),
            newMaxNodeCount
        );
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxValActivePausedCount(uint256 newMaxValActivePausedCount) external onlyConfigurator {
        emit MaxValActivePausedCountUpdated(
            ABv2ConfigLib.MAX_VAL_ACTIVE_PAUSED_COUNT.updateUint(newMaxValActivePausedCount),
            newMaxValActivePausedCount
        );
    }

    /// @inheritdoc IAddressBookV2
    function updateMaxCandReadyCount(uint256 newMaxCandReadyCount) external onlyConfigurator {
        emit MaxCandReadyCountUpdated(
            ABv2ConfigLib.MAX_CAND_READY_COUNT.updateUint(newMaxCandReadyCount),
            newMaxCandReadyCount
        );
    }

    /// @inheritdoc IAddressBookV2
    function updatePfsThreshold(uint256 newPfsThreshold) external onlyConfigurator {
        emit PfsThresholdUpdated(ABv2ConfigLib.PFS_THRESHOLD.updateUint(newPfsThreshold), newPfsThreshold);
    }

    /// @inheritdoc IAddressBookV2
    function updateCfsThreshold(uint256 newCfsThreshold) external onlyConfigurator {
        emit CfsThresholdUpdated(ABv2ConfigLib.CFS_THRESHOLD.updateUint(newCfsThreshold), newCfsThreshold);
    }

    /// @inheritdoc IAddressBookV2
    function updateKefAddress(address newKefAddress) external onlyConfigurator {
        emit KefAddressUpdated(ABv2ConfigLib.KEF_ADDRESS.updateAddress(newKefAddress), newKefAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKifAddress(address newKifAddress) external onlyConfigurator {
        emit KifAddressUpdated(ABv2ConfigLib.KIF_ADDRESS.updateAddress(newKifAddress), newKifAddress);
    }

    /// @inheritdoc IAddressBookV2
    function updateKpfAddress(address newKpfAddress) external onlyConfigurator {
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
        return ($.maxNodeCount, $.maxCandReadyCount);
    }

    /// @inheritdoc IAddressBookV2
    function getMaxValActivePausedCount() external view returns (uint256) {
        return _getStorage().maxValActivePausedCount;
    }

    /// @inheritdoc IAddressBookV2
    function getPfsThreshold() external view returns (uint256) {
        return _getStorage().pfsThreshold;
    }

    /// @inheritdoc IAddressBookV2
    function getCfsThreshold() external view returns (uint256) {
        return _getStorage().cfsThreshold;
    }

    /// @inheritdoc IAddressBookV2
    function getSlotLimits() external view returns (uint256 maxSlotAvailable, uint256 minActiveCount) {
        return getSlotLimitsFor(_getStorage().epochVACount);
    }

    /// @inheritdoc IAddressBookV2
    function getSlotLimitsFor(uint256 n) public pure returns (uint256 maxSlotAvailable, uint256 minActiveCount) {
        return (SlotMath.maxSlotAvailable(n), SlotMath.minActiveCount(n));
    }

    /// @inheritdoc IAddressBookV2
    function getFundAddresses() external view returns (address, address, address) {
        ABv2Storage storage $ = _getStorage();
        return ($.kefAddress, $.kifAddress, $.kpfAddress);
    }

    /// @inheritdoc IAddressBookV2
    function getSuspender() external view returns (address) {
        return _getStorage().suspender;
    }

    /// @inheritdoc IAddressBookV2
    function getConfigurator() external view returns (address) {
        return _getStorage().configurator;
    }

    /// @inheritdoc IAddressBookV2
    function isUsedAddress(address addr) external view returns (bool) {
        return _getStorage().usedAddresses[addr];
    }

    /// @inheritdoc IAddressBookV2
    function currentEpoch() public view returns (uint256) {
        return _currentEpoch();
    }

    /// @inheritdoc IAddressBookV2
    function getEpochVACount() external view returns (uint256) {
        return _getStorage().epochVACount;
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
        uint256 activeLen = $.allNodes.length();
        profiles = new Profile[](activeLen);
        for (uint256 i; i < activeLen; ++i) {
            address nid = $.allNodes.at(i);
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
        uint256 len = $.allNodes.length();
        nodeIdList = new address[](len);
        pubkeyList = new BlsPublicKeyInfo[](len);
        for (uint256 i; i < len; ++i) {
            nodeIdList[i] = $.allNodes.at(i);
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
    function getAllNodesLength() external view override returns (uint256) {
        return _getStorage().allNodes.length();
    }

    /// @inheritdoc IAddressBookV2
    function getSuspendedValidators() external view override returns (address[] memory) {
        return _getStorage().suspendedSet.values();
    }

    /// @inheritdoc IAddressBookV2
    function getRegisteredNodes() external view override returns (address[] memory) {
        return _getStorage().registeredNodes.values();
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
        uint256 activeLen = $.allNodes.length();

        nodeIds = new address[](activeLen);
        stakingContracts = new address[](activeLen);
        rewardAddresses = new address[](activeLen);

        for (uint256 i; i < activeLen; ++i) {
            address nid = $.allNodes.at(i);
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

    /* ========== UPGRADE FUNCTIONS ========== */

    function _authorizeUpgrade(address) internal override onlyOwner {}
}
