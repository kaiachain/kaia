// Copyright 2025 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {Initializable} from "openzeppelin-contracts-upgradeable-5.0/proxy/utils/Initializable.sol";
import {UUPSUpgradeable} from "openzeppelin-contracts-upgradeable-5.0/proxy/utils/UUPSUpgradeable.sol";
import {OwnableUpgradeable} from "openzeppelin-contracts-upgradeable-5.0/access/OwnableUpgradeable.sol";
import {Address} from "openzeppelin-contracts-5.0/utils/Address.sol";
import {EnumerableSet} from "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";
import {ArrayUtils} from "../../../libs/ArrayUtils.sol";
import {ICnStakingV3MultiSigFactory} from "../CnStakingV3MultiSigFactory/ICnStakingV3MultiSigFactory.sol";
import {IPublicDelegationFactoryV2} from "../PublicDelegation/IPublicDelegationFactoryV2.sol";
import {ISimpleBlsRegistry} from "../ISimpleBlsRegistry.sol";
import {IAddressBook} from "../IAddressBook.sol";
import {IValidatorManager, ICnStaking, IRegistry} from "./IValidatorManager.sol";

/// @title ValidatorManager
/// @notice The contract that enables the validator manager to manage the validator nodes
/// @dev This contract assumes that after the initial setup, all validator node info in
///      this contract are also correctly registered in AddressBook and SimpleBlsRegistry.
///      If there's any inconsistency occurred, it should be recovered by the owner.
contract ValidatorManager is Initializable, UUPSUpgradeable, OwnableUpgradeable, IValidatorManager {
    using Address for address;
    using ArrayUtils for address[];
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */
    IAddressBook public constant ABOOK = IAddressBook(0x0000000000000000000000000000000000000400);
    IRegistry public constant REGISTRY = IRegistry(0x0000000000000000000000000000000000000401);

    uint256 public constant MIN_CONSENSUS_NODE_BALANCE = 10 ether;

    bytes32 public constant ZERO48HASH = 0xc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293; // keccak256(hex"00"*48)
    bytes32 public constant ZERO96HASH = 0x46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21; // keccak256(hex"00"*96)

    /* ========== STATE VARIABLES ========== */
    mapping(address => NodeInfo) private _nodeInfo; // consensus node id => node info
    mapping(address => Request) private _requests; // consensus node id => request
    uint256 private _requestCount;

    /// Set of consensus node ids for lookup
    EnumerableSet.AddressSet private _consensusNodeIds;
    EnumerableSet.AddressSet private _consensusNodeIdsWithPendingRequest;

    /* ========== MODIFIERS ========== */
    /// Check if this contract has the authority of AddressBook and SimpleBlsRegistry
    modifier onlyWhenAuthorized() {
        _onlyWhenAuthorized();
        _;
    }

    /// Check if the caller is a manager of existing validator
    modifier onlyNodeManager(address consensusNodeId) {
        _onlyNodeManager(msg.sender, consensusNodeId);
        _;
    }

    /* ========== CONSTRUCTOR ========== */
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(address initialOwner, NodeInfo[] memory nodeInfos) public initializer {
        __Ownable_init(initialOwner);
        __UUPSUpgradeable_init();

        _requestCount = 1;

        for (uint256 i = 0; i < nodeInfos.length; i++) {
            NodeInfo memory nodeInfo = nodeInfos[i];
            _checkRegisteredNode(nodeInfo);
            _addNodeInfo(nodeInfo);
        }
    }

    /* ========== OPERATIONS ========== */
    /// @inheritdoc IValidatorManager
    function manage(address target, bytes memory data) external override onlyOwner returns (bytes memory result) {
        result = target.functionCall(data);
    }

    /// @inheritdoc IValidatorManager
    function addNodeInfo(NodeInfo memory nodeInfo) external override onlyOwner {
        _checkRegisteredNode(nodeInfo);

        // If there's a existing consensus node id, it will override the existing one
        _addNodeInfo(nodeInfo);
        _deleteRequest(nodeInfo.consensusNodeId);

        emit NodeInfoAdded(nodeInfo.consensusNodeId);
    }

    /// @inheritdoc IValidatorManager
    function removeNodeInfo(address consensusNodeId) external override onlyOwner {
        if (!_consensusNodeIds.contains(consensusNodeId)) revert NodeNotFound();

        _deleteNodeInfo(consensusNodeId);
        _deleteRequest(consensusNodeId);

        emit NodeInfoRemoved(consensusNodeId);
    }

    /* ========== REQUEST: SUBMISSION & CANCELLATION ========== */
    /// @inheritdoc IValidatorManager
    function request(address consensusNodeId, RequestType requestType, bytes memory requestData) external override {
        address manager = msg.sender;
        if (requestType != RequestType.Onboarding) {
            // For non-onboarding requests, the caller should be the node manager
            _onlyNodeManager(manager, consensusNodeId);
        }
        if (_requests[consensusNodeId].requestType != RequestType.NoRequest) revert RequestExists();

        _checkRequest(consensusNodeId, requestType, requestData);

        uint256 nextRequestIndex = _requestCount++;
        _addRequest(consensusNodeId, manager, nextRequestIndex, requestType, requestData);

        emit RequestSubmitted(consensusNodeId, manager, nextRequestIndex, requestType, requestData);
    }

    /// @inheritdoc IValidatorManager
    function cancelRequest(address consensusNodeId) external override {
        Request memory targetRequest = _requests[consensusNodeId];
        address manager = targetRequest.manager;
        RequestType requestType = targetRequest.requestType;
        // Note that we can't use `onlyNodeManager` since the `_nodeInfo` is not set yet if the request is onboarding.
        if (requestType == RequestType.NoRequest) revert RequestNotExists();
        if (manager != msg.sender) revert RequestNotOwned();

        _deleteRequest(consensusNodeId);

        emit RequestCancelled(
            consensusNodeId,
            manager,
            targetRequest.requestIndex,
            requestType,
            targetRequest.requestData
        );
    }

    /// @inheritdoc IValidatorManager
    function transferManager(
        address consensusNodeId,
        address newManager
    ) external override onlyNodeManager(consensusNodeId) {
        address currentManager = _nodeInfo[consensusNodeId].manager;
        if (_requests[consensusNodeId].requestType != RequestType.NoRequest) revert RequestExists();
        if (newManager == address(0) || newManager == currentManager) revert InvalidNewManager();

        // Update the node manager
        _nodeInfo[consensusNodeId].manager = newManager;

        emit NodeManagerTransferred(consensusNodeId, currentManager, newManager);
    }

    /* ========== REQUEST: APPROVAL & REJECTION ========== */
    /// @inheritdoc IValidatorManager
    function approveRequest(
        address consensusNodeId,
        uint256 requestIndex
    ) external override onlyOwner onlyWhenAuthorized {
        _onlyExistingRequest(consensusNodeId, requestIndex);
        _approveRequest(consensusNodeId, requestIndex);
    }

    /// @inheritdoc IValidatorManager
    function approveBatchRequest(
        address[] memory consensusNodeIds,
        uint256[] memory requestIndices
    ) external override onlyOwner onlyWhenAuthorized {
        if (consensusNodeIds.length != requestIndices.length) revert InvalidBatchRequestLength();

        for (uint256 i = 0; i < consensusNodeIds.length; i++) {
            _onlyExistingRequest(consensusNodeIds[i], requestIndices[i]);
            _approveRequest(consensusNodeIds[i], requestIndices[i]);
        }
    }

    function _approveRequest(address consensusNodeId, uint256 requestIndex) private {
        Request memory targetRequest = _requests[consensusNodeId];
        RequestType requestType = targetRequest.requestType;
        if (requestType == RequestType.Onboarding) _executeOnboarding(consensusNodeId);
        else if (requestType == RequestType.Offboarding) _executeOffboarding(consensusNodeId);
        else if (requestType == RequestType.AddStakingContract) _executeAddStakingContract(consensusNodeId);
        else if (requestType == RequestType.RemoveStakingContract) _executeRemoveStakingContract(consensusNodeId);
        else revert InvalidRequestType();

        _deleteRequest(consensusNodeId);

        emit RequestApproved(
            consensusNodeId,
            targetRequest.manager,
            requestIndex,
            requestType,
            targetRequest.requestData
        );
    }

    /// @inheritdoc IValidatorManager
    function rejectRequest(address consensusNodeId, uint256 requestIndex) external override onlyOwner {
        _onlyExistingRequest(consensusNodeId, requestIndex);
        _rejectRequest(consensusNodeId, requestIndex);
    }

    function rejectBatchRequest(
        address[] memory consensusNodeIds,
        uint256[] memory requestIndices
    ) external override onlyOwner {
        if (consensusNodeIds.length != requestIndices.length) revert InvalidBatchRequestLength();

        for (uint256 i = 0; i < consensusNodeIds.length; i++) {
            _onlyExistingRequest(consensusNodeIds[i], requestIndices[i]);
            _rejectRequest(consensusNodeIds[i], requestIndices[i]);
        }
    }

    function _rejectRequest(address consensusNodeId, uint256 requestIndex) private {
        Request memory targetRequest = _requests[consensusNodeId];

        _deleteRequest(consensusNodeId);

        emit RequestRejected(
            consensusNodeId,
            targetRequest.manager,
            requestIndex,
            targetRequest.requestType,
            targetRequest.requestData
        );
    }

    /// Register the new entry in AddressBook and SimpleBlsRegistry
    function _executeOnboarding(address consensusNodeId) private {
        OnboardingRequest memory onboardingRequest = abi.decode(
            _requests[consensusNodeId].requestData,
            (OnboardingRequest)
        );

        // Last check for request data
        _checkOnboardingRequest(onboardingRequest);

        _registerCnStakingContract(
            onboardingRequest.consensusNodeId,
            onboardingRequest.stakingContract,
            onboardingRequest.rewardAddress
        );
        _registerSimpleBlsRegistry(onboardingRequest);

        _addNodeInfo(NodeInfo(onboardingRequest.manager, onboardingRequest.consensusNodeId, new address[](0)));
    }

    /// In SimpleBlsRegistry, it only holds one record per validator, mapped with the consensus node id.
    /// Since SBR checks the existence of the consensus node id in the AddressBook, we need to unregister the consensus node id first from AddressBook.
    function _executeOffboarding(address consensusNodeId) private {
        // Last check for request data
        _checkOffboardingRequest(consensusNodeId);

        address[] memory nodeIds = _nodeInfo[consensusNodeId].nodeIds;
        for (uint256 i = 0; i < nodeIds.length; i++) {
            _unregisterCnStakingContract(nodeIds[i]);
        }
        _unregisterCnStakingContract(consensusNodeId);
        _unregisterSimpleBlsRegistry(consensusNodeId);

        _deleteNodeInfo(consensusNodeId);
    }

    /// Execute adding a new staking contract
    function _executeAddStakingContract(address consensusNodeId) private {
        (address nodeId, address stakingContract, address rewardAddress) = abi.decode(
            _requests[consensusNodeId].requestData,
            (address, address, address)
        );

        // Last check for request data
        _checkAddStakingContractRequest(consensusNodeId, nodeId, stakingContract, rewardAddress);

        _registerCnStakingContract(nodeId, stakingContract, rewardAddress);

        _nodeInfo[consensusNodeId].nodeIds.push(nodeId);
    }

    /// Execute removing a staking contract
    function _executeRemoveStakingContract(address consensusNodeId) private {
        address nodeId = abi.decode(_requests[consensusNodeId].requestData, (address));

        // Last check for request data
        _checkRemoveStakingContractRequest(consensusNodeId, nodeId);

        _unregisterCnStakingContract(nodeId);

        _nodeInfo[consensusNodeId].nodeIds.remove(nodeId);
    }

    /* ========== VALIDATION ========== */

    // ========== Request Dispatcher ==========
    function _checkRequest(address consensusNodeId, RequestType requestType, bytes memory requestData) private view {
        if (requestType == RequestType.Onboarding) {
            OnboardingRequest memory onboardingRequest = abi.decode(requestData, (OnboardingRequest));
            // Check if the requested manager is the same as the caller
            if (onboardingRequest.manager != msg.sender) revert InvalidNodeManager();
            _checkOnboardingRequest(onboardingRequest);
        } else if (requestType == RequestType.Offboarding) {
            _checkOffboardingRequest(consensusNodeId);
        } else if (requestType == RequestType.AddStakingContract) {
            (address nodeId, address stakingContract, address rewardAddress) = abi.decode(
                requestData,
                (address, address, address)
            );
            _checkAddStakingContractRequest(consensusNodeId, nodeId, stakingContract, rewardAddress);
        } else if (requestType == RequestType.RemoveStakingContract) {
            address nodeId = abi.decode(requestData, (address));
            _checkRemoveStakingContractRequest(consensusNodeId, nodeId);
        } else {
            revert InvalidRequestType();
        }
    }

    // ========== Request-Specific Validations ==========
    function _checkOnboardingRequest(OnboardingRequest memory onboardingRequest) private view {
        address consensusNodeId = onboardingRequest.consensusNodeId;
        // 1. Consensus node id MUST NOT exist
        _onlyNewValidator(consensusNodeId);

        // We already verified `msg.sender == manager` in `_checkRequest`
        address manager = onboardingRequest.manager;

        // 2. BLS key info MUST pass the basic validation
        _checkBlsKeyInfo(onboardingRequest.blsPublicKeyInfo.publicKey, onboardingRequest.blsPublicKeyInfo.pop);

        // 3. Consensus node id MUST have enough balance
        if (consensusNodeId.balance < MIN_CONSENSUS_NODE_BALANCE) revert BelowMinNodeIdBalance();

        // 4. Staking contract MUST be valid
        _checkNewStakingContract(
            manager,
            consensusNodeId,
            onboardingRequest.stakingContract,
            onboardingRequest.rewardAddress
        );
    }

    function _checkOffboardingRequest(address consensusNodeId) private view {
        // 1. Consensus node id MUST exist
        _onlyExistingValidator(consensusNodeId);
    }

    function _checkAddStakingContractRequest(
        address consensusNodeId,
        address nodeId,
        address stakingContract,
        address rewardAddress
    ) private view {
        // 1. Consensus node id MUST exist
        _onlyExistingValidator(consensusNodeId);

        // 2. Staking contract MUST be valid
        // If the new node id is duplicated with the currently managed node id (including the consensus node id), it will be caught by `_isNodeInABook`.
        _checkNewStakingContract(_nodeInfo[consensusNodeId].manager, nodeId, stakingContract, rewardAddress);

        // 3. Staking contract MUST be compatible with the consensus node staking contract
        _checkCompatibleStakingContract(consensusNodeId, stakingContract);
    }

    function _checkRemoveStakingContractRequest(address consensusNodeId, address nodeId) private view {
        // 1. Consensus node id MUST exist
        _onlyExistingValidator(consensusNodeId);

        // 2. New nodeId MUST NOT be the consensus node id
        if (_nodeInfo[consensusNodeId].consensusNodeId == nodeId) revert InvalidTargetNodeId();

        // 3. New nodeId MUST be in the current node info
        if (!_nodeInfo[consensusNodeId].nodeIds.contains(nodeId)) revert InvalidTargetNodeId();
    }

    // ========== Node Validations ==========
    /// Check if the consensus node id doesn't exist
    function _onlyNewValidator(address consensusNodeId) private view {
        if (_consensusNodeIds.contains(consensusNodeId)) revert NodeAlreadyExists();
    }

    /// Check if the consensus node id exists
    function _onlyExistingValidator(address consensusNodeId) private view {
        if (!_consensusNodeIds.contains(consensusNodeId)) revert NodeNotFound();
    }

    /// Check if the node info is valid
    /// It assumes that the `nodeInfo` is already registered in AddressBook and SimpleBlsRegistry
    /// Only used for initial node info registration and emergency recovery
    function _checkRegisteredNode(NodeInfo memory nodeInfo) private view {
        // 1. Addresses MUST NOT be zero address
        if (nodeInfo.manager == address(0)) revert InvalidNodeInfo(0);
        if (nodeInfo.consensusNodeId == address(0)) revert InvalidNodeInfo(1);
        if (nodeInfo.nodeIds.contains(address(0))) revert InvalidNodeInfo(2);

        // 2. Node ids MUST be distinct
        if (nodeInfo.nodeIds.contains(nodeInfo.consensusNodeId)) revert InvalidNodeInfo(3);
        if (!nodeInfo.nodeIds.isDistinct()) revert InvalidNodeInfo(4);

        address consensusNodeId = nodeInfo.consensusNodeId;
        address[] memory nodeIds = nodeInfo.nodeIds;

        // 3. Consensus node id MUST have enough balance
        if (consensusNodeId.balance < MIN_CONSENSUS_NODE_BALANCE) revert InvalidNodeInfo(5);

        // 4. Consensus node id MUST be in AddressBook and SimpleBlsRegistry
        if (!_isNodeInABook(consensusNodeId)) revert InvalidNodeInfo(6);
        if (!_isNodeInSBR(consensusNodeId)) revert InvalidNodeInfo(7);

        // 5. All node ids MUST be in AddressBook and compatible with the consensus node staking contract
        for (uint256 i = 0; i < nodeIds.length; i++) {
            if (!_isNodeInABook(nodeIds[i])) revert InvalidNodeInfo(8);

            (, address stakingContract, ) = ABOOK.getCnInfo(nodeIds[i]);
            _checkCompatibleStakingContract(consensusNodeId, stakingContract);
        }
    }

    // ========== Staking Contract Validations ==========
    /// Check if the new staking contract is valid
    function _checkNewStakingContract(
        address manager,
        address newNodeId,
        address newStakingContract,
        address newRewardAddress
    ) private view {
        // a. Node id MUST NOT be in AddressBook
        //    Don't need to check SimpleBlsRegistry since it can't be registered without AddressBook registration.
        if (_isNodeInABook(newNodeId)) revert InvalidStakingContract(0);

        ICnStaking cnStaking = ICnStaking(newStakingContract);
        // b. Staking contract MUST be deployed by the manager through registered factory
        if (!_getCnStakingV3MultiSigFactory().isDeployedBy(manager, newStakingContract))
            revert InvalidStakingContract(1);

        // c. Staking contract MUST be initialized
        if (!cnStaking.isInitialized()) revert InvalidStakingContract(2);

        // d. Staking contract MUST have the same node id
        if (cnStaking.nodeId() != newNodeId) revert InvalidStakingContract(3);

        // e. Staking contract MUST have the same reward address
        if (cnStaking.rewardAddress() != newRewardAddress) revert InvalidStakingContract(4);

        // f. Public delegation MUST be deployed by the staking contract through registered factory
        // AddStakingContractRequest, which doesn't allow public delegation, will be lazily checked in `_isCompatibleStakingContract`
        if (_isPublicDelegationEnabled(cnStaking)) {
            if (!_getPublicDelegationFactory().isDeployedBy(newStakingContract, newRewardAddress))
                revert InvalidStakingContract(5);
        }
    }

    /// Check if the new staking contract is compatible with the consensus node staking contract
    /// It assumes that the new staking contract is valid
    function _checkCompatibleStakingContract(address consensusNodeId, address newStakingContract) private view {
        (, address consensusStakingContract, ) = ABOOK.getCnInfo(consensusNodeId);
        ICnStaking consensusCnStaking = ICnStaking(consensusStakingContract);
        ICnStaking newCnStaking = ICnStaking(newStakingContract);

        // a. Public delegation MUST be disabled
        if (_isPublicDelegationEnabled(newCnStaking)) revert IncompatibleNewStakingContract(0);

        // b. GC ID MUST be the same
        if (consensusCnStaking.gcId() != newCnStaking.gcId()) revert IncompatibleNewStakingContract(1);

        // c. Reward address MUST be the same
        if (consensusCnStaking.rewardAddress() != newCnStaking.rewardAddress())
            revert IncompatibleNewStakingContract(2);
    }

    // ========== BLS Key Validations ==========
    function _checkBlsKeyInfo(bytes memory publicKey, bytes memory pop) private pure {
        // a. BLS key MUST have correct length
        if (publicKey.length != 48 || pop.length != 96) revert InvalidBLSKeyInfo(0);

        // b. BLS key MUST NOT be zero
        if (keccak256(publicKey) == ZERO48HASH || keccak256(pop) == ZERO96HASH) revert InvalidBLSKeyInfo(1);
    }

    // ========== Registry Helper Checks ==========
    function _isPublicDelegationEnabled(ICnStaking cnStaking) private view returns (bool) {
        try cnStaking.isPublicDelegationEnabled() returns (bool isPublicDelegationEnabled) {
            return isPublicDelegationEnabled;
        } catch {
            // It means CnStakingContract under VERSION < 3, which doesn't support public delegation
            return false;
        }
    }

    function _isNodeInABook(address nodeId) private view returns (bool) {
        try ABOOK.getCnInfo(nodeId) {
            return true;
        } catch {
            return false;
        }
    }

    function _isNodeInSBR(address nodeId) private view returns (bool) {
        (bytes memory publicKey, ) = _getSimpleBlsRegistry().record(nodeId);
        return publicKey.length > 0;
    }

    /* ========== ACTIONS ========== */
    function _registerCnStakingContract(address nodeId, address stakingContract, address rewardAddress) private {
        ABOOK.submitRegisterCnStakingContract(nodeId, stakingContract, rewardAddress);
    }

    function _unregisterCnStakingContract(address nodeId) private {
        ABOOK.submitUnregisterCnStakingContract(nodeId);
    }

    function _registerSimpleBlsRegistry(OnboardingRequest memory onboardingRequest) private {
        _getSimpleBlsRegistry().register(
            onboardingRequest.consensusNodeId,
            onboardingRequest.blsPublicKeyInfo.publicKey,
            onboardingRequest.blsPublicKeyInfo.pop
        );
    }

    function _unregisterSimpleBlsRegistry(address nodeId) private {
        _getSimpleBlsRegistry().unregister(nodeId);
    }

    /* ========== HELPERS & UTILS ========== */
    function _addNodeInfo(NodeInfo memory nodeInfo) private {
        address consensusNodeId = nodeInfo.consensusNodeId;
        _nodeInfo[consensusNodeId] = nodeInfo;
        _consensusNodeIds.add(consensusNodeId);
    }

    function _deleteNodeInfo(address consensusNodeId) private {
        delete _nodeInfo[consensusNodeId];
        _consensusNodeIds.remove(consensusNodeId);
    }

    function _addRequest(
        address consensusNodeId,
        address manager,
        uint256 requestIndex,
        RequestType requestType,
        bytes memory requestData
    ) private {
        _requests[consensusNodeId] = Request(manager, requestIndex, requestType, requestData);
        _consensusNodeIdsWithPendingRequest.add(consensusNodeId);
    }

    function _deleteRequest(address consensusNodeId) private {
        delete _requests[consensusNodeId];
        _consensusNodeIdsWithPendingRequest.remove(consensusNodeId);
    }

    function _onlyWhenAuthorized() private view {
        (address[] memory adminList, uint256 requirement) = ABOOK.getState();
        if (!adminList.contains(address(this))) revert ABookNotExecutable();
        if (requirement != 1) revert ABookNotExecutable();

        if (_getSimpleBlsRegistry().owner() != address(this)) revert SimpleBlsRegistryNotExecutable();
    }

    function _onlyNodeManager(address manager, address consensusNodeId) private view {
        // If node not existed, it will be anyway reverted with `InvalidNodeManager`, but for better error message, we check it explicitly here
        _onlyExistingValidator(consensusNodeId);
        if (_nodeInfo[consensusNodeId].manager != manager) revert InvalidNodeManager();
    }

    function _onlyExistingRequest(address consensusNodeId, uint256 requestIndex) private view {
        if (_requests[consensusNodeId].requestType == RequestType.NoRequest) revert RequestNotExists();
        if (_requests[consensusNodeId].requestIndex != requestIndex) revert RequestNotExists();
    }

    function _getSimpleBlsRegistry() private view returns (ISimpleBlsRegistry) {
        address simpleBlsRegistry = REGISTRY.getActiveAddr("KIP113");
        if (simpleBlsRegistry == address(0)) revert SimpleBlsRegistryNotFound();

        return ISimpleBlsRegistry(simpleBlsRegistry);
    }

    function _getCnStakingV3MultiSigFactory() private view returns (ICnStakingV3MultiSigFactory) {
        address cnStakingV3MultiSigFactory = REGISTRY.getActiveAddr("CnStakingV3MultiSigFactory");
        if (cnStakingV3MultiSigFactory == address(0)) revert CnStakingV3MultiSigFactoryNotFound();

        return ICnStakingV3MultiSigFactory(cnStakingV3MultiSigFactory);
    }

    function _getPublicDelegationFactory() private view returns (IPublicDelegationFactoryV2) {
        address publicDelegationFactory = REGISTRY.getActiveAddr("PublicDelegationFactory");
        if (publicDelegationFactory == address(0)) revert PublicDelegationFactoryNotFound();

        return IPublicDelegationFactoryV2(publicDelegationFactory);
    }

    /* ========== UPGRADE ========== */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}

    /* ========== GETTERS ========== */
    /// @inheritdoc IValidatorManager
    function checkRequest(
        address consensusNodeId,
        RequestType requestType,
        bytes memory requestData
    ) public view override {
        _checkRequest(consensusNodeId, requestType, requestData);
    }

    /// @inheritdoc IValidatorManager
    function getConsensusNodeIds() external view override returns (address[] memory) {
        return _consensusNodeIds.values();
    }

    /// @inheritdoc IValidatorManager
    function getNodeInfo(address consensusNodeId) external view override returns (NodeInfo memory) {
        return _nodeInfo[consensusNodeId];
    }

    /// @inheritdoc IValidatorManager
    function getRequest(address consensusNodeId) external view override returns (Request memory) {
        return _requests[consensusNodeId];
    }

    /// @inheritdoc IValidatorManager
    function getAllNodeInfos() external view override returns (NodeInfo[] memory nodeInfos) {
        uint256 nodeInfoCnt = _consensusNodeIds.length();
        nodeInfos = new NodeInfo[](nodeInfoCnt);
        for (uint256 i = 0; i < nodeInfoCnt; i++) {
            address consensusNodeId = _consensusNodeIds.at(i);
            nodeInfos[i] = _nodeInfo[consensusNodeId];
        }
    }

    /// @inheritdoc IValidatorManager
    function getPendingRequests()
        external
        view
        override
        returns (address[] memory consensusNodeIds, Request[] memory requests)
    {
        uint256 pendingRequestCnt = _consensusNodeIdsWithPendingRequest.length();
        consensusNodeIds = new address[](pendingRequestCnt);
        requests = new Request[](pendingRequestCnt);
        for (uint256 i = 0; i < pendingRequestCnt; i++) {
            address consensusNodeId = _consensusNodeIdsWithPendingRequest.at(i);
            consensusNodeIds[i] = consensusNodeId;
            requests[i] = _requests[consensusNodeId];
        }
    }

    /// @inheritdoc IValidatorManager
    function getDuplicatedNode(address nodeId) external view override returns (address[] memory consensusNodeIds) {
        uint256 cnt = 0;
        uint256 consensusNodeLen = _consensusNodeIds.length();
        consensusNodeIds = new address[](consensusNodeLen);

        if (_consensusNodeIds.contains(nodeId)) {
            consensusNodeIds[cnt++] = nodeId;
        }

        for (uint256 i = 0; i < _consensusNodeIds.length(); i++) {
            address consensusNodeId = _consensusNodeIds.at(i);
            if (_nodeInfo[consensusNodeId].nodeIds.contains(nodeId)) {
                consensusNodeIds[cnt++] = consensusNodeId;
            }
        }

        assembly {
            mstore(consensusNodeIds, cnt)
        }
    }
}
