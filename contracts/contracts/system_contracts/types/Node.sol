// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @notice A list of possible states for a node
/// @dev Please refer to KIP-286 for more details
enum State {
    Unknown,      // 0
    CandInactive, // 1
    CandReady,    // 2
    CandTesting,  // 3
    ValInactive,  // 4
    ValReady,     // 5
    ValActive,    // 6
    ValPaused,    // 7
    ValExiting    // 8
}

/// @notice A struct to store the BLS public key and proof-of-possession
/// @dev Please refer to KIP-113 for more details
struct BlsPublicKeyInfo {
    /// @dev compressed BLS12-381 public key (48 bytes)
    bytes publicKey;
    /// @dev proof-of-possession (96 bytes)
    ///  must be a result of PopProve algorithm as per
    ///  draft-irtf-cfrg-bls-signature-05 section 3.3.3.
    ///  with ciphersuite "BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_"
    bytes pop;
}

/// @notice A struct to store the node info
struct NodeInfo {
    /// @dev The manager address that controls this node
    address manager;
    /// @dev Used to stake and unstake, immutable
    address stakingContract;
    /// @dev Used to receive rewards, immutable if PD=on
    address rewardAddress;
    /// @dev Used to vote on-chain governance, mutable
    address voterAddress;
    /// @dev Used to store the pause or idle timeout timestamp, 0 if not paused or idle
    uint256 timeoutAt;
    /// @dev Used to identify the governance council, immutable
    uint256 gcId;
    /// @dev Used to verify signatures, immutable
    BlsPublicKeyInfo blsInfo;
    /// @dev Used to store the metadata of the validator, in JSON format
    string metadata;
    /// @dev Used to store the current state of the node
    State state;
}

/// @notice Lightweight struct for getAllProfiles() — no BLS data
struct Profile {
    address nodeId;
    address stakingContract;
    address rewardAddress;
    uint256 timeoutAt;
    State state;
}
