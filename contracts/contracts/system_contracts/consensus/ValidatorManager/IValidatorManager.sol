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
pragma solidity ^0.8.0;

import {IKIP113} from "../../kip113/IKIP113.sol";

interface IValidatorManager {
    /* ========== TYPES ========== */
    enum RequestType {
        NoRequest,
        Onboarding,
        Offboarding,
        AddStakingContract,
        RemoveStakingContract
    }

    struct Request {
        address manager; // the manager that submitted the request
        uint256 requestIndex; // requestIndex is used to prevent request-replacement attack when approving or rejecting a pending request
        RequestType requestType;
        bytes requestData;
    }

    struct OnboardingRequest {
        address manager; // it should be same as the msg.sender at the time of request submission
        address consensusNodeId; // consensus node id that used for consensus
        address stakingContract; // staking contract address
        address rewardAddress; // if staking contract uses public delegation, this should be same as public delegation contract
        IKIP113.BlsPublicKeyInfo blsPublicKeyInfo; // refer to IKIP113.sol
    }

    struct NodeInfo {
        address manager;
        address consensusNodeId;
        address[] nodeIds; // adding staking contract requires separate node id and it won't participate in consensus
    }

    /* ========== ERRORS ========== */
    error RequestExists();
    error RequestNotExists();
    error RequestNotOwned();
    error InvalidNewManager();
    error InvalidRequestType();
    error NodeAlreadyExists();
    error BLSKeyAlreadyExists();
    error BelowMinNodeIdBalance();
    error InvalidTargetNodeId();
    error InvalidBatchRequestLength();
    error ABookNotExecutable();
    error SimpleBlsRegistryNotExecutable();
    error NodeNotFound();
    error InvalidNodeManager();
    error SimpleBlsRegistryNotFound();
    error CnStakingV3MultiSigFactoryNotFound();
    error PublicDelegationFactoryNotFound();

    /// @dev Error codes for InvalidNodeInfo:
    ///      0: Zero manager address
    ///      1: Zero consensus node id
    ///      2: Zero node id in node ids
    ///      3: Consensus node id is in node ids
    ///      4: Node ids are not distinct
    ///      5: Below min node id balance
    ///      6: Consensus node id is not in AddressBook
    ///      7: Consensus node id is not in SimpleBlsRegistry
    ///      8: Node id is not in AddressBook
    error InvalidNodeInfo(uint256 errorCode);

    /// @dev Error codes for InvalidBLSKeyInfo:
    ///      0: Invalid BLS key length
    ///      1: BLS key is zero
    error InvalidBLSKeyInfo(uint256 errorCode);

    /// @dev Error codes for InvalidStakingContract:
    ///      0: Staking contract is already in AddressBook
    ///      1: Staking contract is not deployed by the manager through registered factory
    ///      2: Staking contract is not initialized
    ///      3: Staking contract has the different node id
    ///      4: Staking contract has the different reward address
    ///      5: Public delegation is not deployed by the staking contract through registered factory
    error InvalidStakingContract(uint256 errorCode);

    /// @dev Error codes for IncompatibleNewStakingContract:
    ///      0: Public delegation is enabled
    ///      1: GC ID is different
    ///      2: Reward address is different
    error IncompatibleNewStakingContract(uint256 errorCode);

    /* ========== EVENTS ========== */
    event NodeManagerTransferred(
        address indexed consensusNodeId,
        address indexed oldNodeManager,
        address indexed newNodeManager
    );
    event NodeInfoAdded(address indexed consensusNodeId);
    event NodeInfoRemoved(address indexed consensusNodeId);
    event RequestSubmitted(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        RequestType requestType,
        bytes requestData
    );
    event RequestCancelled(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        RequestType requestType,
        bytes requestData
    );
    event RequestApproved(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        RequestType requestType,
        bytes requestData
    );
    event RequestRejected(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        RequestType requestType,
        bytes requestData
    );

    /* ========== FUNCTIONS ========== */
    /// @notice Manage the target contract, especially used for AddressBook and SimpleBlsRegistry
    /// @dev Callable by owner
    ///      It grants ValidatorManager the permission to call any function of the target contract
    function manage(
        address target,
        bytes memory data
    ) external returns (bytes memory result);

    /// @notice Directly add the node info
    /// @param nodeInfo The node info to be added
    /// @dev Callable by owner
    ///      This function is used to prevent any state inconsistency between AddressBook and ValidatorManager
    ///      and assumes that the node info is already registered in AddressBook
    function addNodeInfo(NodeInfo memory nodeInfo) external;

    /// @notice Directly remove the node info
    /// @param consensusNodeId The consensus node id to be removed
    /// @dev Callable by owner
    function removeNodeInfo(address consensusNodeId) external;

    /// @notice Request a new validator onboarding, offboarding, adding a new staking contract, or removing a staking contract
    /// @param consensusNodeId The consensus node id that the request is for
    /// @param requestType The type of request
    /// @param requestData The encoded request data
    /// @dev Callable by validator manager
    /// @dev Note that this only allows to submit one request at a time
    function request(
        address consensusNodeId,
        RequestType requestType,
        bytes memory requestData
    ) external;

    /// @notice Cancel existing request for a consensus node id
    /// @param consensusNodeId The consensus node id that the request is for
    /// @dev Callable by validator manager
    function cancelRequest(address consensusNodeId) external;

    /// @notice Transfer the validator manager to a new address
    /// @dev Callable by validator manager
    ///      Can't transfer manager if there is a pending request
    function transferManager(
        address consensusNodeId,
        address newManager
    ) external;

    /// @notice Approve a request
    /// @param consensusNodeId The consensus node id that has a request to be approved
    /// @param requestIndex The request index to be approved
    /// @dev Callable by owner
    function approveRequest(
        address consensusNodeId,
        uint256 requestIndex
    ) external;

    /// @notice Approve a batch of requests
    /// @param consensusNodeIds The consensus node ids that have requests to be approved
    /// @param requestIndices The request indices to be approved
    /// @dev Callable by owner
    function approveBatchRequest(
        address[] memory consensusNodeIds,
        uint256[] memory requestIndices
    ) external;

    /// @notice Reject request of a validator manager
    /// @param consensusNodeId The consensus node id that has a request to be rejected
    /// @param requestIndex The request index to be rejected
    /// @dev Callable by owner
    function rejectRequest(
        address consensusNodeId,
        uint256 requestIndex
    ) external;

    /// @notice Reject a batch of requests
    /// @param consensusNodeIds The consensus node ids that have requests to be rejected
    /// @param requestIndices The request indices to be rejected
    /// @dev Callable by owner
    function rejectBatchRequest(
        address[] memory consensusNodeIds,
        uint256[] memory requestIndices
    ) external;

    /* ========== GETTERS ========== */
    /// @notice Check if a request is valid
    /// @param consensusNodeId The consensus node id that the request is for
    /// @param requestType The type of request
    /// @param requestData The encoded request data
    /// @dev This function is a wrapper function
    ///      When the request is onboarding, the "from" field should be set to the node manager
    function checkRequest(
        address consensusNodeId,
        RequestType requestType,
        bytes memory requestData
    ) external view;

    /// @notice Get all consensus node ids
    /// @return consensusNodeIds All consensus node ids
    function getConsensusNodeIds() external view returns (address[] memory);

    /// @notice Get the node info of a consensus node id
    /// @param consensusNodeId The consensus node id
    /// @return nodeInfo The node info of the consensus node id
    function getNodeInfo(
        address consensusNodeId
    ) external view returns (NodeInfo memory);

    /// @notice Get the request of a consensus node id
    /// @param consensusNodeId The consensus node id
    /// @return request The request of the consensus node id
    function getRequest(
        address consensusNodeId
    ) external view returns (Request memory);

    /// @notice Get all node infos
    /// @return nodeInfos All node infos
    function getAllNodeInfos()
        external
        view
        returns (NodeInfo[] memory nodeInfos);

    /// @notice Get all pending requests
    /// @return consensusNodeIds All consensus node ids with pending requests
    /// @return requests All pending requests
    function getPendingRequests()
        external
        view
        returns (address[] memory consensusNodeIds, Request[] memory requests);

    /// @notice Get the duplicated node ids
    /// @param nodeId The node id to get the duplicated node ids
    /// @return consensusNodeIds The consensus node ids that have the duplicated node id
    /// @dev If it returns length > 1, it means that the node id is duplicated between multiple nodes, which requires state correction
    function getDuplicatedNode(
        address nodeId
    ) external view returns (address[] memory consensusNodeIds);
}

interface ICnStaking {
    function isInitialized() external view returns (bool);

    function nodeId() external view returns (address);

    function rewardAddress() external view returns (address);

    function gcId() external view returns (uint256);

    function voterAddress() external view returns (address);

    function isPublicDelegationEnabled() external view returns (bool);
}

interface IRegistry {
    function getActiveAddr(string memory name) external view returns (address);
}
