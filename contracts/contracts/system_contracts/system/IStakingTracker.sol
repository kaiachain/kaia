// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

/// @title IStakingTracker
/// @notice Minimal interface for the StakingTracker contract.
/// @dev Used by AddressBookV2 (refreshVoter) and CnStakingV4 (refreshStake).
/// TODO-permissionless: Replace this with official staking tracker interface.
interface IStakingTracker {
    function refreshVoter(address staking) external;
    function refreshStake(address staking) external;
}
