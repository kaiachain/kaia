// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {BeaconProxy} from "@openzeppelin/contracts/proxy/beacon/BeaconProxy.sol";
import {ICnStaking} from "../CnStakingV4/interfaces/ICnStaking.sol";
import {IPublicDelegation} from "../../PublicDelegation/interfaces/IPublicDelegation.sol";
import {ICnStakingV4Factory} from "./interfaces/ICnStakingV4Factory.sol";

/// @title CnStakingV4Factory
/// @notice Permissionless factory for deploying CnStaking and PublicDelegation beacon proxies.
///
/// Deploys BeaconProxy instances pointing to shared UpgradeableBeacons.
/// Two deployment modes:
/// 1. deployCnStaking: CnStaking proxy without PublicDelegation.
/// 2. deployCnStakingWithPD: CnStaking + PublicDelegation proxies deployed atomically.
///
/// No owner — fully permissionless. Beacon addresses are immutable.
/// Tracks all deployed CnStaking proxies and their deployers for factory-based validation.
contract CnStakingV4Factory is ICnStakingV4Factory {
    /* ========== IMMUTABLES AND CONSTANTS ========== */

    /// @inheritdoc ICnStakingV4Factory
    address public immutable cnStakingBeacon;

    /// @inheritdoc ICnStakingV4Factory
    address public immutable pdBeacon;

    /// @inheritdoc ICnStakingV4Factory
    uint256 public constant INITIAL_LOCKUP = 1e9;

    /// @inheritdoc ICnStakingV4Factory
    address public constant DEAD_ADDRESS = address(0xdead);

    /* ========== TRACKING STATE ========== */

    /// @notice Tracks all factory-deployed CnStaking proxy addresses.
    mapping(address => bool) private _deployedCnStaking;

    /// @notice Tracks all factory-deployed PublicDelegation proxy addresses.
    mapping(address => bool) private _deployedPublicDelegation;

    /// @notice CnStaking or PD proxy address → deployer (msg.sender who called deploy).
    mapping(address => address) private _deployer;

    /* ========== CONSTRUCTOR ========== */

    /// @param _cnStakingBeacon The UpgradeableBeacon for CnStaking.
    /// @param _pdBeacon        The UpgradeableBeacon for PublicDelegation.
    constructor(address _cnStakingBeacon, address _pdBeacon) {
        if (_cnStakingBeacon == address(0)) revert ZeroAddress();
        if (_pdBeacon == address(0)) revert ZeroAddress();

        cnStakingBeacon = _cnStakingBeacon;
        pdBeacon = _pdBeacon;
    }

    /* ========== DEPLOYMENT FUNCTIONS ========== */

    /// @inheritdoc ICnStakingV4Factory
    function deployCnStaking(address _owner) external returns (address proxy) {
        bytes32 salt = keccak256(abi.encode(msg.sender, _owner));
        proxy = address(new BeaconProxy{salt: salt}(cnStakingBeacon, ""));
        ICnStaking(payable(proxy)).initialize(_owner);

        _deployedCnStaking[proxy] = true;
        _deployer[proxy] = msg.sender;

        emit DeployCnStakingV4(proxy, _owner);
    }

    /// @inheritdoc ICnStakingV4Factory
    function deployCnStakingWithPD(
        address _owner,
        IPublicDelegation.PDConstructorArgs memory _pdArgs
    ) external payable returns (address proxy, address publicDelegation) {
        if (msg.value < INITIAL_LOCKUP) revert InsufficientInitialStake();

        bytes32 cnSalt = keccak256(abi.encode(msg.sender, _owner));
        bytes32 pdSalt = keccak256(abi.encode(msg.sender, _owner, _pdArgs));

        // Deploy both beacon proxies via CREATE2
        proxy = address(new BeaconProxy{salt: cnSalt}(cnStakingBeacon, ""));
        publicDelegation = address(new BeaconProxy{salt: pdSalt}(pdBeacon, ""));

        // Initialize PD first (needs CnStaking address)
        IPublicDelegation(payable(publicDelegation)).initialize(proxy, _pdArgs);

        // Initialize CnStaking (needs PD address)
        ICnStaking(payable(proxy)).initializeWithPD(_owner, publicDelegation);

        // Mint dead shares to prevent inflation attack
        IPublicDelegation(payable(publicDelegation)).stakeFor{value: msg.value}(DEAD_ADDRESS);

        _deployedCnStaking[proxy] = true;
        _deployedPublicDelegation[publicDelegation] = true;
        _deployer[proxy] = msg.sender;
        _deployer[publicDelegation] = msg.sender;

        emit DeployCnStakingV4WithPD(proxy, _owner, publicDelegation);
    }

    /* ========== PUBLIC GETTERS ========== */

    /// @inheritdoc ICnStakingV4Factory
    function isDeployedCnStaking(address _addr) external view returns (bool) {
        return _deployedCnStaking[_addr];
    }

    /// @inheritdoc ICnStakingV4Factory
    function isDeployedPublicDelegation(address _addr) external view returns (bool) {
        return _deployedPublicDelegation[_addr];
    }

    /// @inheritdoc ICnStakingV4Factory
    function getDeployer(address _addr) external view returns (address) {
        return _deployer[_addr];
    }
}
