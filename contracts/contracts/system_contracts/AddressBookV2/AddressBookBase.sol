// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {ICnStaking} from "../CnStaking/interfaces/ICnStaking.sol";
import {NodeState, BlsPublicKeyInfo, NodeInfo, Profile} from "../types/Node.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {SystemCallable} from "../system/SystemCallable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title AddressBookBase
/// @notice Abstract base for AddressBookV2. Contains both ERC-7201 storage namespaces,
///         constants, access control modifiers, internal helpers, and internal mutators.
/// @dev Two ERC-7201 namespaced storage slots:
///      - NodeStorage at 0x2a6c4... (node info, sets, state counts)
///      - OperationStorage at 0x9114b... (managers, scores, config)
abstract contract AddressBookBase is
    IAddressBookV2,
    Initializable,
    OwnableUpgradeable,
    UUPSUpgradeable,
    SystemCallable
{
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    /// @notice Minimum staking amount required for candidate activation and validator readiness (5M KAIA)
    uint256 public constant MIN_STAKE = 5_000_000 ether;

    /* ========== NODE STORAGE ========== */

    /// @custom:storage-location erc7201:addressbookv2.storage.NodeStorage
    struct NodeStorage {
        /// @notice Node info mapping: nodeId => full NodeInfo
        mapping(address => NodeInfo) nodeInfo;
        /// @notice Active set: all nodes in states other than CandInactive (CandReady through ValExiting)
        EnumerableSet.AddressSet activeSet;
        /// @notice CandInactive set: nodes in CandInactive state only
        EnumerableSet.AddressSet candInactiveSet;
        /// @notice Suspended validators (excluded from effective committee by owner)
        EnumerableSet.AddressSet suspendedSet;
        /// @notice The last assigned governance council (gc) id, starting from 1
        uint256 lastAssignedGCId;
        /// @notice Per-state node count, atomically updated with every state change
        mapping(NodeState => uint256) stateCount;
    }

    // keccak256(abi.encode(uint256(keccak256("addressbookv2.storage.NodeStorage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant NODE_STORAGE_LOCATION = 0x2a6c4f9950bdac880731815408ac22648e61e86a4c1659ca1864add32c0e2900;

    function _getNodeStorage() internal pure returns (NodeStorage storage $) {
        assembly {
            $.slot := NODE_STORAGE_LOCATION
        }
    }

    /* ========== OPERATION STORAGE ========== */

    /// @custom:storage-location erc7201:addressbookv2.storage.OperationStorage
    struct OperationStorage {
        /// @notice nodeId → manager address (access control for node operations)
        mapping(address => address) managers;
        /// @notice Tracks all registered addresses (nodeId, stakingContract, rewardAddress) for uniqueness
        mapping(address => bool) registeredAddresses;
        /// @notice Epoch → nodeId → score
        mapping(uint256 => mapping(address => uint256)) scores;
        /// @notice Number of consecutive proposal failures that triggers forced exit
        uint256 exitThreshold;
        /// @notice Duration (seconds) a paused validator remains in ValPaused before expiry
        uint256 pauseTimeout;
        /// @notice Duration (seconds) an idle validator remains in ValInactive/ValReady before expiry
        uint256 idleTimeout;
        /// @notice Maximum number of validators allowed in the pipeline (activeSet)
        uint256 maxValidatorCount;
        /// @notice Maximum number of candidates in CandReady state
        uint256 maxReadyCandidateCount;
        /// @notice Snapshot of (valActive + valPaused) at epoch start, used for slot math
        uint256 epochValCount;
        /// @notice Whether genesis validators have been initialized (one-time flag)
        bool validatorsInitialized;
    }

    // keccak256(abi.encode(uint256(keccak256("addressbookv2.storage.OperationStorage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant OPERATION_STORAGE_LOCATION =
        0x9114bff6fc2186e0cb5fe9e5ced5765c61f7dbcd47451c215ce8da2cac3ed800;

    function _getOperationStorage() internal pure returns (OperationStorage storage $) {
        assembly {
            $.slot := OPERATION_STORAGE_LOCATION
        }
    }

    /* ========== MODIFIERS ========== */

    /// @notice Restricts a function to the registered manager of the given node
    modifier onlyManager(address nodeId) {
        _onlyManager(nodeId);
        _;
    }

    function _onlyManager(address nodeId) private view {
        if (msg.sender != _getOperationStorage().managers[nodeId]) revert OnlyManager();
    }

    /* ========== INTERNAL MUTATORS ========== */

    /// @notice Stores a new node as CandInactive
    /// @dev Assigns gcId, stores nodeInfo, adds to candInactiveSet.
    function _createNode(address nodeId, NodeInfo memory info) internal {
        NodeStorage storage $ = _getNodeStorage();

        if ($.nodeInfo[nodeId].state != NodeState.Unknown) revert NodeAlreadyExists();

        uint256 gcId = ++$.lastAssignedGCId;

        $.nodeInfo[nodeId] = info;
        $.nodeInfo[nodeId].gcId = gcId;

        $.candInactiveSet.add(nodeId);
        $.stateCount[NodeState.CandInactive]++;

        emit NodeCreated(nodeId, gcId);
    }

    /// @notice Removes a CandInactive node entirely
    function _deleteNode(address nodeId) internal {
        NodeStorage storage $ = _getNodeStorage();

        // Defense-in-depth: caller already checks CandInactive state; this guards against bugs
        if ($.nodeInfo[nodeId].state != NodeState.CandInactive) revert InvalidState();

        $.candInactiveSet.remove(nodeId);
        $.stateCount[NodeState.CandInactive]--;

        delete $.nodeInfo[nodeId];

        emit NodeDeleted(nodeId);
    }

    /// @notice Batch-creates validators directly as ValActive in the active set.
    /// @dev Used by initializeValidators for genesis setup.
    function _createActiveValidators(address[] calldata nodeIds, NodeInfo[] calldata infos) internal {
        NodeStorage storage $ = _getNodeStorage();
        uint256 len = nodeIds.length;

        for (uint256 i; i < len; ) {
            if ($.nodeInfo[nodeIds[i]].state != NodeState.Unknown) revert NodeAlreadyExists();

            uint256 gcId = ++$.lastAssignedGCId;

            $.nodeInfo[nodeIds[i]] = infos[i];
            $.nodeInfo[nodeIds[i]].gcId = gcId;
            $.nodeInfo[nodeIds[i]].state = NodeState.ValActive;
            $.nodeInfo[nodeIds[i]].timeoutAt = 0;

            $.activeSet.add(nodeIds[i]);

            emit NodeCreated(nodeIds[i], gcId);

            unchecked {
                ++i;
            }
        }

        $.stateCount[NodeState.ValActive] += len;
    }

    /// @notice Batch state transition for system processing
    function _batchTransition(
        address[] calldata nodeIds,
        NodeState[] calldata newStates,
        uint256[] calldata timeouts
    ) internal {
        uint256 len = nodeIds.length;
        if (len != newStates.length || len != timeouts.length) revert InvalidInput();

        for (uint256 i; i < len; ) {
            _transition(nodeIds[i], newStates[i], timeouts[i]);
            unchecked {
                ++i;
            }
        }
    }

    /// @notice Atomic state transition: updates set boundary, state, and timeout
    /// @dev No state validation — action functions are responsible for all pre-checks.
    ///      Handles CandInactive <-> activeSet boundary crossing automatically.
    function _transition(address nodeId, NodeState newState, uint256 timeoutAt) internal {
        NodeStorage storage $ = _getNodeStorage();
        NodeInfo storage info = $.nodeInfo[nodeId];
        NodeState oldState = info.state;

        if (oldState == NodeState.Unknown || newState == NodeState.Unknown || oldState == newState) return;

        $.stateCount[oldState]--;
        $.stateCount[newState]++;

        if (oldState == NodeState.CandInactive && newState != NodeState.CandInactive) {
            $.candInactiveSet.remove(nodeId);
            $.activeSet.add(nodeId);
        } else if (oldState != NodeState.CandInactive && newState == NodeState.CandInactive) {
            $.activeSet.remove(nodeId);
            $.candInactiveSet.add(nodeId);
        }

        info.state = newState;
        info.timeoutAt = timeoutAt;

        emit StateChanged(nodeId, oldState, newState);
    }

    /* ========== INTERNAL HELPERS ========== */

    function _isNodeAtState(address nodeId, NodeState state) internal view returns (bool) {
        return _getNodeState(nodeId) == state;
    }

    /// @notice Checks whether a node's effective stake meets MIN_STAKE
    function _isNodeOverMinStake(address nodeId) internal view returns (bool) {
        address stakingContract = _getStakingContract(nodeId);
        uint256 effectiveStake = ICnStaking(stakingContract).staking() - ICnStaking(stakingContract).unstaking();
        return effectiveStake >= MIN_STAKE;
    }

    /// @notice Reverts if the node's timeout has expired.
    /// @dev Defensive guard against the timing window between timeout expiry and
    ///      processSystemTransition execution. The core client will eventually demote
    ///      the node, but without this check a user could call resume() or exit()
    ///      in that window to escape the pending demotion.
    function _revertIfTimeoutExpired(address nodeId) internal view {
        uint256 timeoutAt = _getTimeoutAt(nodeId);
        if (timeoutAt != 0 && block.timestamp >= timeoutAt) revert TimeoutExpired();
    }

    /* ========== INTERNAL GETTERS ========== */

    function _getNodeState(address nodeId) internal view returns (NodeState) {
        return _getNodeStorage().nodeInfo[nodeId].state;
    }

    function _getStakingContract(address nodeId) internal view returns (address) {
        return _getNodeStorage().nodeInfo[nodeId].stakingContract;
    }

    function _getTimeoutAt(address nodeId) internal view returns (uint256) {
        return _getNodeStorage().nodeInfo[nodeId].timeoutAt;
    }

    function _getStateCount(NodeState state) internal view returns (uint256) {
        return _getNodeStorage().stateCount[state];
    }

    function _getActiveSetLength() internal view returns (uint256) {
        return _getNodeStorage().activeSet.length();
    }

    function _getCandInactiveSetLength() internal view returns (uint256) {
        return _getNodeStorage().candInactiveSet.length();
    }

    function _getNodeInfo(address nodeId) internal view returns (NodeInfo storage) {
        return _getNodeStorage().nodeInfo[nodeId];
    }
}
