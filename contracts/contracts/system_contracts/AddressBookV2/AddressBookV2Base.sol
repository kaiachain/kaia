// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {IAddressBookV2} from "./interfaces/IAddressBookV2.sol";
import {ICnStaking} from "../CnStaking/interfaces/ICnStaking.sol";
import {State, BlsPublicKeyInfo, NodeInfo, Profile} from "../types/Node.sol";
import {EnumerableSet} from "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";
import {SystemCallable} from "../system/SystemCallable.sol";
import {Initializable} from "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import {OwnableUpgradeable} from "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import {UUPSUpgradeable} from "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/// @title AddressBookV2Base
/// @notice Abstract base for AddressBookV2. Contains ERC-7201 storage namespace,
///         constants, access control modifiers, internal helpers, and internal mutators.
abstract contract AddressBookV2Base is
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

    /* ========== STORAGE ========== */

    /// @custom:storage-location erc7201:addressbookv2.storage.ABv2Storage
    struct ABv2Storage {
        // Node data
        mapping(address => NodeInfo) nodeInfo;
        EnumerableSet.AddressSet activeSet;
        EnumerableSet.AddressSet candInactiveSet;
        EnumerableSet.AddressSet suspendedSet;
        uint256 lastAssignedGCId;
        mapping(State => uint256) stateCount;
        // Operations
        mapping(address => bool) registeredAddresses;
        mapping(uint256 => mapping(address => uint256)) scores;
        uint256 exitThreshold;
        uint256 pauseTimeout;
        uint256 idleTimeout;
        uint256 maxValidatorCount;
        uint256 maxReadyCandidateCount;
        uint256 epochValCount;
        address kefAddress;
        address kifAddress;
        address kpfAddress;
    }

    // keccak256(abi.encode(uint256(keccak256("addressbookv2.storage.ABv2Storage")) - 1)) & ~bytes32(uint256(0xff))
    bytes32 private constant STORAGE_LOCATION = 0x1b0484cbd0fba815b5886ffd853c75e18f4b5720362abf431d1f348b59d4ff00;

    function _getStorage() internal pure returns (ABv2Storage storage $) {
        assembly {
            $.slot := STORAGE_LOCATION
        }
    }

    /* ========== MODIFIERS ========== */

    /// @notice Restricts a function to the registered manager of the given node
    modifier onlyManager(address nodeId) {
        _onlyManager(nodeId);
        _;
    }

    function _onlyManager(address nodeId) private view {
        if (msg.sender != _getStorage().nodeInfo[nodeId].manager) revert OnlyManager();
    }

    /* ========== INTERNAL MUTATORS ========== */

    /// @notice Stores a new node as CandInactive
    /// @dev Assigns gcId, stores nodeInfo, adds to candInactiveSet.
    function _createNode(address nodeId, NodeInfo memory info) internal {
        ABv2Storage storage $ = _getStorage();

        uint256 gcId = ++$.lastAssignedGCId;

        $.nodeInfo[nodeId] = info;
        $.nodeInfo[nodeId].gcId = gcId;

        $.candInactiveSet.add(nodeId);
        $.stateCount[State.CandInactive]++;

        emit NodeCreated(nodeId, gcId);
    }

    /// @notice Removes a CandInactive node entirely
    function _deleteNode(address nodeId) internal {
        ABv2Storage storage $ = _getStorage();

        // Defense-in-depth: caller already checks CandInactive state; this guards against bugs
        if ($.nodeInfo[nodeId].state != State.CandInactive) revert InvalidState();

        $.candInactiveSet.remove(nodeId);
        $.stateCount[State.CandInactive]--;

        delete $.nodeInfo[nodeId];

        emit NodeDeleted(nodeId);
    }

    /// @notice Batch-sets validators directly as ValActive in the active set.
    /// @dev Used by initialize() for genesis setup via ABv2DataContract.
    function _setInitialActiveValidators(address[] memory nodeIds, NodeInfo[] memory infos) internal {
        ABv2Storage storage $ = _getStorage();
        uint256 len = nodeIds.length;

        for (uint256 i; i < len; ) {
            uint256 gcId = ++$.lastAssignedGCId;

            $.nodeInfo[nodeIds[i]] = infos[i];
            // TODO-permissionless: Set gcID from data contract
            $.nodeInfo[nodeIds[i]].gcId = gcId;
            $.nodeInfo[nodeIds[i]].state = State.ValActive;
            $.nodeInfo[nodeIds[i]].timeoutAt = 0;

            $.activeSet.add(nodeIds[i]);

            unchecked {
                ++i;
            }
        }

        $.stateCount[State.ValActive] = len;
        $.epochValCount = len;
    }

    /// @notice Batch state transition for system processing
    function _batchTransition(
        address[] calldata nodeIds,
        State[] calldata newStates,
        uint256[] calldata timeouts
    ) internal {
        uint256 len = nodeIds.length;
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
    function _transition(address nodeId, State newState, uint256 timeoutAt) internal {
        ABv2Storage storage $ = _getStorage();
        NodeInfo storage info = $.nodeInfo[nodeId];
        State oldState = info.state;

        if (oldState == State.Unknown || newState == State.Unknown || oldState == newState) return;

        $.stateCount[oldState]--;
        $.stateCount[newState]++;

        if (oldState == State.CandInactive && newState != State.CandInactive) {
            $.candInactiveSet.remove(nodeId);
            $.activeSet.add(nodeId);
        } else if (oldState != State.CandInactive && newState == State.CandInactive) {
            $.activeSet.remove(nodeId);
            $.candInactiveSet.add(nodeId);
        }

        info.state = newState;
        info.timeoutAt = timeoutAt;

        emit StateChanged(nodeId, oldState, newState);
    }

    /* ========== INTERNAL HELPERS ========== */

    function _isNodeAtState(address nodeId, State state) internal view returns (bool) {
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

    function _getNodeState(address nodeId) internal view returns (State) {
        return _getStorage().nodeInfo[nodeId].state;
    }

    function _getStakingContract(address nodeId) internal view returns (address) {
        return _getStorage().nodeInfo[nodeId].stakingContract;
    }

    function _getTimeoutAt(address nodeId) internal view returns (uint256) {
        return _getStorage().nodeInfo[nodeId].timeoutAt;
    }

    function _getStateCount(State state) internal view returns (uint256) {
        return _getStorage().stateCount[state];
    }

    function _getActiveSetLength() internal view returns (uint256) {
        return _getStorage().activeSet.length();
    }

    function _getCandInactiveSetLength() internal view returns (uint256) {
        return _getStorage().candInactiveSet.length();
    }

    function _getNodeInfo(address nodeId) internal view returns (NodeInfo storage) {
        return _getStorage().nodeInfo[nodeId];
    }
}
