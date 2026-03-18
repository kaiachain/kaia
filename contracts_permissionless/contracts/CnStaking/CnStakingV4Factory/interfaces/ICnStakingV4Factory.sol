// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {IPublicDelegation} from "../../../PublicDelegation/interfaces/IPublicDelegation.sol";

/// @title ICnStakingV4Factory
/// @notice Interface for the permissionless CnStaking and PublicDelegation beacon proxy factory.
interface ICnStakingV4Factory {
    /* ========== EVENTS ========== */

    event DeployCnStakingV4(address indexed proxy, address indexed owner);
    event DeployCnStakingV4WithPD(address indexed proxy, address indexed owner, address publicDelegation);

    /* ========== ERRORS ========== */

    error ZeroAddress();
    error InsufficientInitialStake();

    /* ========== DEPLOYMENT FUNCTIONS ========== */

    /// @notice Deploy a CnStaking proxy without PublicDelegation.
    /// @dev Uses CREATE2 with msg.sender in salt to prevent front-running.
    ///      Same (deployer, owner) combo cannot deploy twice (CREATE2 collision = revert).
    /// @param _owner  The owner of the CnStaking contract.
    /// @return proxy  The deployed BeaconProxy address.
    function deployCnStaking(address _owner) external returns (address proxy);

    /// @notice Deploy CnStaking + PublicDelegation proxies atomically.
    /// @dev Deploys both beacon proxies via CREATE2, initializes PD with CnStaking address,
    ///      then initializes CnStaking with PD address.
    ///      Salt includes msg.sender to prevent front-running.
    /// @param _owner  The owner of both contracts.
    /// @param _pdArgs PublicDelegation constructor arguments.
    /// @return proxy            The deployed CnStaking BeaconProxy address.
    /// @return publicDelegation The deployed PublicDelegation BeaconProxy address.
    function deployCnStakingWithPD(
        address _owner,
        IPublicDelegation.PDConstructorArgs memory _pdArgs
    ) external payable returns (address proxy, address publicDelegation);

    /* ========== GETTERS ========== */

    /// @notice The UpgradeableBeacon for CnStaking implementation.
    function cnStakingBeacon() external view returns (address);

    /// @notice The UpgradeableBeacon for PublicDelegation implementation.
    function pdBeacon() external view returns (address);

    /// @notice Minimum initial stake to prevent inflation attack (minted as dead shares).
    function INITIAL_LOCKUP() external pure returns (uint256);

    /// @notice Dead address that receives initial shares to prevent inflation attack.
    function DEAD_ADDRESS() external pure returns (address);

    /// @notice Returns true if the address is a factory-deployed CnStaking proxy.
    /// @param _addr The address to check.
    function isDeployedCnStaking(address _addr) external view returns (bool);

    /// @notice Returns true if the address is a factory-deployed PublicDelegation proxy.
    /// @param _addr The address to check.
    function isDeployedPublicDelegation(address _addr) external view returns (bool);

    /// @notice Returns the deployer of a factory-deployed CN or PD proxy. address(0) if not deployed by factory.
    /// @param _addr The proxy address to query.
    function getDeployer(address _addr) external view returns (address);
}
