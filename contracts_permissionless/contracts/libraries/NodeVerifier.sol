// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import {BlsPublicKeyInfo} from "../types/Node.sol";
import {IRegistry} from "../system/IRegistry.sol";
import {ICnStaking} from "../CnStaking/CnStakingV4/interfaces/ICnStaking.sol";
import {ICnStakingV4Factory} from "../CnStaking/CnStakingV4Factory/interfaces/ICnStakingV4Factory.sol";

/// @title NodeVerifier
/// @notice Validates node registration inputs and manages address uniqueness registry.
/// @dev Used by AddressBookV2 (NodeActions) for createNode / deleteNode.
library NodeVerifier {
    error InvalidInput();
    error AddressAlreadyRegistered();
    error FactoryNotFound();
    error StakingDeployerMismatch();

    address private constant REGISTRY_ADDRESS = address(0x401);

    bytes32 private constant ZERO48HASH = 0xc980e59163ce244bb4bb6211f48c7b46f88a4f40943e84eb99bdc41e129bd293; // keccak256(hex"00"*48)
    bytes32 private constant ZERO96HASH = 0x46700b4d40ac5c35af2c22dda2787a91eb567b06c924a8fb8ae9a05b20c08c21; // keccak256(hex"00"*96)

    /// @notice Validates node registration inputs, verifies the staking contract deployer, and registers addresses.
    /// @dev Used by NodeActions.createNode() at runtime. Checks that msg.sender deployed the staking contract.
    /// @param registry The address uniqueness registry
    /// @param nodeId The node address
    /// @param stakingContract The staking contract address
    /// @param rewardAddress The reward address
    /// @param blsInfo The BLS public key and proof-of-possession
    function registerNode(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress,
        BlsPublicKeyInfo memory blsInfo
    ) internal {
        _checkInputs(nodeId, stakingContract, rewardAddress, blsInfo);

        // Deployer check: msg.sender must be the one who deployed the staking contract via factory
        address factory = _getFactory();
        if (factory == address(0)) revert FactoryNotFound();
        if (ICnStakingV4Factory(factory).getDeployer(stakingContract) != msg.sender) revert StakingDeployerMismatch();

        _checkPublicDelegation(stakingContract, rewardAddress);
        _registerAddresses(registry, nodeId, stakingContract, rewardAddress);
    }

    /// @notice Validates node registration inputs without factory check.
    /// @dev Used by ABv2DataContract constructor for genesis data. Genesis data is trusted.
    function registerNodeGenesis(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress,
        BlsPublicKeyInfo memory blsInfo
    ) internal {
        _checkInputs(nodeId, stakingContract, rewardAddress, blsInfo);
        _checkPublicDelegation(stakingContract, rewardAddress);

        _registerAddresses(registry, nodeId, stakingContract, rewardAddress);
    }

    function _checkInputs(
        address nodeId,
        address stakingContract,
        address rewardAddress,
        BlsPublicKeyInfo memory blsInfo
    ) private pure {
        // Zero address checks
        if (nodeId == address(0) || stakingContract == address(0) || rewardAddress == address(0)) revert InvalidInput();
        // Mutual uniqueness among registerable addresses
        if (nodeId == stakingContract || nodeId == rewardAddress || stakingContract == rewardAddress) {
            revert InvalidInput();
        }
        // BLS public key (48 bytes) and proof-of-possession (96 bytes)
        if (blsInfo.publicKey.length != 48 || blsInfo.pop.length != 96) revert InvalidInput();
        if (keccak256(blsInfo.publicKey) == ZERO48HASH || keccak256(blsInfo.pop) == ZERO96HASH) revert InvalidInput();
    }

    function _checkPublicDelegation(address stakingContract, address rewardAddress) private view {
        address pd = ICnStaking(payable(stakingContract)).publicDelegation();
        if (pd != address(0) && pd != rewardAddress) revert InvalidInput();
    }

    function _registerAddresses(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress
    ) private {
        // Registry uniqueness check
        if (registry[nodeId] || registry[stakingContract] || registry[rewardAddress]) revert AddressAlreadyRegistered();

        // Register addresses
        registry[nodeId] = true;
        registry[stakingContract] = true;
        registry[rewardAddress] = true;
    }

    /// @dev Resolves CnStakingFactory from registry
    function _getFactory() private view returns (address) {
        return IRegistry(REGISTRY_ADDRESS).getActiveAddr("CnStakingFactory");
    }

    /// @notice Unregisters previously registered addresses
    function unregisterAddresses(
        mapping(address => bool) storage registry,
        address nodeId,
        address stakingContract,
        address rewardAddress
    ) internal {
        registry[nodeId] = false;
        registry[stakingContract] = false;
        registry[rewardAddress] = false;
    }
}
