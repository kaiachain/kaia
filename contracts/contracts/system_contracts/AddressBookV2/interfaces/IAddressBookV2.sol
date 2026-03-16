// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import {State, BlsPublicKeyInfo, NodeInfo, Profile} from "../../types/Node.sol";

/// @title IAddressBookV2
/// @notice Unified interface for AddressBookV2 — the sole entry point for all node operations at 0x400.
interface IAddressBookV2 {
    /* ========== ERRORS ========== */

    error NotInitializable();

    /// @notice Thrown when caller is not the node's registered manager
    error OnlyManager();

    /// @notice Thrown when caller is not the node ID itself
    error OnlyNodeId();

    /// @notice Thrown when a node is not in a valid state for the requested operation
    error InvalidState();

    /// @notice Thrown when a slot limit (candidate ready, pause, exit) would be exceeded
    error SlotsFull();

    /// @notice Thrown when an input parameter is invalid (zero address, mismatched array lengths, etc.)
    error InvalidInput();

    /// @notice Thrown when a node's staking contract balance is below MIN_STAKE
    error StakingTooLow();

    /// @notice Thrown when attempting to create a node that already exists
    error NodeAlreadyExists();

    /// @notice Thrown when operating on a node that does not exist (state == Unknown)
    error NodeNotFound();

    /// @notice Thrown when a user action is attempted after the node's timeout has expired
    error TimeoutExpired();

    /// @notice Thrown when attempting to suspend a validator that is already suspended
    error AlreadySuspended();

    /// @notice Thrown when attempting to unsuspend a validator that is not suspended
    error NotSuspended();

    /// @notice Thrown when attempting to update the reward address while public delegation is enabled
    error PDEnabled();

    /* ========== EVENTS ========== */

    /// @notice Emitted when a node's state changes
    /// @param nodeId The address of the node
    /// @param fromState The previous state
    /// @param toState The new state
    event StateChanged(address indexed nodeId, State indexed fromState, State indexed toState);

    /// @notice Emitted when a new node is created
    /// @param nodeId The address of the newly created node
    /// @param gcId The governance council ID assigned to the node
    event NodeCreated(address indexed nodeId, uint256 gcId);

    /// @notice Emitted when a node is deleted
    /// @param nodeId The address of the deleted node
    event NodeDeleted(address indexed nodeId);

    /// @notice Emitted when a validator is suspended by the owner
    /// @param nodeId The address of the suspended node
    event ValidatorSuspended(address indexed nodeId);

    /// @notice Emitted when a validator is unsuspended by the owner
    /// @param nodeId The address of the unsuspended node
    event ValidatorUnsuspended(address indexed nodeId);

    /// @notice Emitted when a node's manager is updated
    /// @param nodeId The address of the node
    /// @param oldManager The old manager address
    /// @param newManager The new manager address
    event ManagerUpdated(address indexed nodeId, address indexed oldManager, address indexed newManager);

    /// @notice Emitted when a node's reward address is updated
    /// @param nodeId The address of the node
    /// @param oldRewardAddress The old reward address
    /// @param newRewardAddress The new reward address
    event RewardAddressUpdated(
        address indexed nodeId,
        address indexed oldRewardAddress,
        address indexed newRewardAddress
    );

    /// @notice Emitted when a node's voter address is updated
    /// @param nodeId The address of the node
    /// @param oldVoterAddress The old voter address
    /// @param newVoterAddress The new voter address
    event VoterAddressUpdated(address indexed nodeId, address indexed oldVoterAddress, address indexed newVoterAddress);

    /// @notice Emitted when a node's metadata is updated
    /// @param nodeId The address of the node
    /// @param newMetadata The new metadata string
    event MetadataUpdated(address indexed nodeId, string newMetadata);

    /// @notice Emitted when a candidate transitions from CandInactive to CandReady
    /// @param nodeId The address of the readied candidate
    event CandidateReadied(address indexed nodeId);

    /// @notice Emitted when a candidate transitions from CandReady to CandInactive
    /// @param nodeId The address of the unreadied candidate
    event CandidateUnreadied(address indexed nodeId);

    /// @notice Emitted when scores are updated for an epoch
    /// @param epoch The epoch for which scores were updated
    /// @param nodeIds The node addresses whose scores were updated
    /// @param scores The new score values
    event ScoresUpdated(uint256 indexed epoch, address[] nodeIds, uint256[] scores);

    /// @notice Emitted when system transitions are processed
    /// @param nodeIds The addresses of the transitioned nodes
    /// @param newStates The new states for each node
    event SystemTransitionProcessed(address[] nodeIds, State[] newStates);

    /// @notice Emitted when the epoch transition is processed
    /// @param epochValCount The epoch validator count
    event EpochTransitionProcessed(uint256 epochValCount);

    /// @notice Emitted when the pause timeout is updated
    /// @param oldPauseTimeout The previous timeout value
    /// @param newPauseTimeout The new timeout value
    event PauseTimeoutUpdated(uint256 oldPauseTimeout, uint256 newPauseTimeout);

    /// @notice Emitted when the idle timeout is updated
    /// @param oldIdleTimeout The previous timeout value
    /// @param newIdleTimeout The new timeout value
    event IdleTimeoutUpdated(uint256 oldIdleTimeout, uint256 newIdleTimeout);

    /// @notice Emitted when the exit threshold is updated
    /// @param oldThreshold The previous threshold value
    /// @param newThreshold The new threshold value
    event ExitThresholdUpdated(uint256 oldThreshold, uint256 newThreshold);

    /// @notice Emitted when the maximum validator count is updated
    /// @param oldCount The previous count
    /// @param newCount The new count
    event MaxValidatorCountUpdated(uint256 oldCount, uint256 newCount);

    /// @notice Emitted when the maximum ready candidate count is updated
    /// @param oldCount The previous count
    /// @param newCount The new count
    event MaxReadyCandidateCountUpdated(uint256 oldCount, uint256 newCount);

    /// @notice Emitted when the KEF address is updated
    /// @param oldAddress The previous address
    /// @param newAddress The new address
    event KefAddressUpdated(address oldAddress, address newAddress);

    /// @notice Emitted when the KIF address is updated
    /// @param oldAddress The previous address
    /// @param newAddress The new address
    event KifAddressUpdated(address oldAddress, address newAddress);

    /// @notice Emitted when the KPF address is updated
    /// @param oldAddress The previous address
    /// @param newAddress The new address
    event KpfAddressUpdated(address oldAddress, address newAddress);

    /// @notice Emitted when genesis validators are initialized
    /// @param nodeIds The addresses of the initialized validators
    event ValidatorsInitialized(address[] nodeIds);

    /* ========== USER FUNCTIONS ========== */

    /// @notice Registers a new node as CandInactive
    /// @param nodeId The address of the node to register
    /// @param stakingContract The staking contract address
    /// @param rewardAddress The reward address
    /// @param voterAddress The voter address for on-chain governance (optional)
    /// @param blsInfo The BLS public key and proof-of-possession
    /// @param metadata The node metadata in JSON format
    function createNode(
        address nodeId,
        address stakingContract,
        address rewardAddress,
        address voterAddress,
        BlsPublicKeyInfo memory blsInfo,
        string memory metadata
    ) external;

    /// @notice Removes a CandInactive node entirely
    /// @param nodeId The address of the node to delete
    function deleteNode(address nodeId) external;

    /// @notice Updates the manager of a node
    /// @param nodeId The address of the node
    /// @param newManager The new manager address
    function updateManager(address nodeId, address newManager) external;

    /// @notice Updates the reward address of a node (only when public delegation is disabled)
    /// @param nodeId The address of the node
    /// @param newRewardAddress The new reward address
    function updateRewardAddress(address nodeId, address newRewardAddress) external;

    /// @notice Updates the voter address of a node
    /// @param nodeId The address of the node
    /// @param newVoterAddress The new voter address
    function updateVoterAddress(address nodeId, address newVoterAddress) external;

    /// @notice Updates the metadata of a node
    /// @param nodeId The address of the node
    /// @param newMetadata The new metadata string
    function updateMetadata(address nodeId, string calldata newMetadata) external;

    /// @notice Readies a candidate: CandInactive → CandReady
    /// @param nodeId The address of the candidate to ready
    function readyCandidate(address nodeId) external;

    /// @notice Unreadies a candidate: CandReady → CandInactive
    /// @param nodeId The address of the candidate to unready
    function unreadyCandidate(address nodeId) external;

    /// @notice Readies a validator: ValInactive → ValReady
    /// @param nodeId The address of the validator
    function readyValidator(address nodeId) external;

    /// @notice Unreadies a validator: ValReady → ValInactive
    /// @param nodeId The address of the validator
    function unreadyValidator(address nodeId) external;

    /// @notice Pauses a validator: ValActive → ValPaused
    /// @param nodeId The address of the validator to pause
    function pause(address nodeId) external;

    /// @notice Resumes a paused validator: ValPaused → ValActive
    /// @param nodeId The address of the validator to resume
    function resume(address nodeId) external;

    /// @notice Exits a validator: ValActive/ValPaused → ValExiting
    /// @param nodeId The address of the validator to exit
    function exit(address nodeId) external;

    /// @notice Offboards a validator: ValInactive → CandInactive
    /// @param nodeId The address of the validator to offboard
    function offboard(address nodeId) external;

    /* ========== SYSTEM FUNCTIONS ========== */

    /// @notice Processes system-triggered state transitions for nodes.
    /// @dev Core client computes all timeout/violation/epoch logic; AddressBookV2 unconditionally records results.
    ///      Updates epochValCount only at epoch boundaries (block.number % EPOCH_BLOCK_INTERVAL == 0).
    /// @param nodeIds Array of node addresses to transition
    /// @param newStates Array of target states for each node
    /// @param timeoutAts Array of timeout timestamps for each node (0 = no timeout)
    function processSystemTransition(
        address[] calldata nodeIds,
        State[] calldata newStates,
        uint256[] calldata timeoutAts
    ) external;

    /// @notice Updates scores for nodes at epoch boundaries
    /// @param nodeIds Array of node addresses
    /// @param scores Array of score values
    function updateScores(address[] calldata nodeIds, uint256[] calldata scores) external;

    /* ========== ADMIN FUNCTIONS ========== */

    /// @notice Suspends a validator (owner emergency action)
    /// @param nodeId The address of the validator to suspend
    function suspendValidator(address nodeId) external;

    /// @notice Unsuspends a validator
    /// @param nodeId The address of the validator to unsuspend
    function unsuspendValidator(address nodeId) external;

    /// @notice Updates the pause timeout duration
    /// @param newPauseTimeout The new timeout in seconds
    function updatePauseTimeout(uint256 newPauseTimeout) external;

    /// @notice Updates the idle timeout duration
    /// @param newIdleTimeout The new timeout in seconds
    function updateIdleTimeout(uint256 newIdleTimeout) external;

    /// @notice Updates the maximum validator count
    /// @param newMaxValidatorCount The new maximum
    function updateMaxValidatorCount(uint256 newMaxValidatorCount) external;

    /// @notice Updates the maximum ready candidate count
    /// @param newMaxReadyCandidateCount The new maximum
    function updateMaxReadyCandidateCount(uint256 newMaxReadyCandidateCount) external;

    /// @notice Updates the proposal failure exit threshold
    /// @param newExitThreshold The new threshold
    function updateExitThreshold(uint256 newExitThreshold) external;

    /// @notice Updates the KEF address
    /// @param newKefAddress The new KEF address
    function updateKefAddress(address newKefAddress) external;

    /// @notice Updates the KIF address
    /// @param newKifAddress The new KIF address
    function updateKifAddress(address newKifAddress) external;

    /// @notice Updates the KPF address
    /// @param newKpfAddress The new KPF address
    function updateKpfAddress(address newKpfAddress) external;

    /* ========== GETTERS (node management) ========== */

    /// @notice Returns the manager address of a node
    /// @param nodeId The address of the node
    /// @return The manager address
    function getManager(address nodeId) external view returns (address);

    /// @notice Checks if an address is registered (nodeId, stakingContract, or rewardAddress)
    /// @param addr The address to check
    /// @return True if the address is registered
    function isRegistered(address addr) external view returns (bool);

    /// @notice Returns the pause and idle timeout durations
    /// @return pauseTimeout The pause timeout in seconds
    /// @return idleTimeout The idle timeout in seconds
    function getTimeouts() external view returns (uint256 pauseTimeout, uint256 idleTimeout);

    /// @notice Returns the maximum validator and ready candidate counts
    /// @return maxValidatorCount The maximum validator count
    /// @return maxReadyCandidateCount The maximum ready candidate count
    function getMaxCounts() external view returns (uint256 maxValidatorCount, uint256 maxReadyCandidateCount);

    /// @notice Returns the proposal failure exit threshold
    /// @return The exit threshold
    function getExitThreshold() external view returns (uint256);

    /// @notice Returns the fund addresses (KEF, KIF, KPF)
    /// @return kefAddress The Kaia Ecosystem Fund address
    /// @return kifAddress The Kaia Infrastructure Fund address
    /// @return kpfAddress The Kaia Protocol Fund address
    function getFundAddresses() external view returns (address kefAddress, address kifAddress, address kpfAddress);

    /// @notice Returns the score for a node in a specific epoch
    /// @param epoch The epoch number
    /// @param nodeId The address of the node
    /// @return The score
    function getScore(uint256 epoch, address nodeId) external view returns (uint256);

    /// @notice Returns the current epoch number
    /// @return The current epoch
    function currentEpoch() external view returns (uint256);

    /// @notice Returns the epoch validator count snapshot (ValActive + ValPaused at last epoch)
    /// @return The epoch validator count used for mid-epoch slot math
    function getEpochValCount() external view returns (uint256);

    /* ========== GETTERS (node info) ========== */

    /// @notice Returns the full NodeInfo for a specific node
    /// @dev Reverts NodeNotFound if the node's state is Unknown.
    /// @param nodeId The address of the node
    /// @return The NodeInfo struct
    function getNodeInfo(address nodeId) external view returns (NodeInfo memory);

    /// @notice Returns NodeInfo structs for multiple nodes
    /// @param nodeIds Array of node addresses to fetch
    /// @return Array of NodeInfo structs in the same order as nodeIds
    function getNodeInfos(address[] calldata nodeIds) external view returns (NodeInfo[] memory);

    /// @notice Returns lightweight profiles for all registered nodes, excluding suspended validators
    /// @dev Iterates activeSet, skips suspended nodes.
    /// @return Array of Profile structs
    function getAllProfiles() external view returns (Profile[] memory);

    /// @notice Returns BLS public key info for all registered nodes
    /// @dev Iterates activeSet. Matches SimpleBlsRegistry interface.
    /// @return nodeIdList Array of node addresses
    /// @return pubkeyList Array of BlsPublicKeyInfo structs in the same order
    function getAllBlsInfo() external view returns (address[] memory nodeIdList, BlsPublicKeyInfo[] memory pubkeyList);

    /// @notice Returns the current state of a node
    /// @param nodeId The address of the node
    /// @return The node's current State
    function getNodeState(address nodeId) external view returns (State);

    /// @notice Checks if a node is in the active set (all states except CandInactive)
    /// @param nodeId The address of the node
    /// @return True if the node is in activeSet
    function isInActiveSet(address nodeId) external view returns (bool);

    /// @notice Returns all currently suspended validators
    /// @return Array of suspended node addresses
    function getSuspendedValidators() external view returns (address[] memory);

    /// @notice Returns the staking contract address of a node
    /// @param nodeId The address of the node
    /// @return The staking contract address
    function getStakingContract(address nodeId) external view returns (address);

    /// @notice Returns the timeout timestamp of a node
    /// @param nodeId The address of the node
    /// @return The timeout timestamp (0 if no timeout)
    function getTimeoutAt(address nodeId) external view returns (uint256);

    /// @notice Returns the number of nodes in the active set
    /// @return The active set length
    function getActiveSetLength() external view returns (uint256);

    /// @notice Returns the number of nodes in a given state
    /// @dev Atomically maintained inside _transition(), _createNode(), _deleteNode().
    /// @param state The State to query
    /// @return The count of nodes in that state
    function getStateCount(State state) external view returns (uint256);
}
