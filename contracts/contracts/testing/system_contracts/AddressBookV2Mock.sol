// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.18;

import {Profile, State} from "../../system_contracts/types/Node.sol";

/// @title AddressBookV2Mock
/// @notice Minimal mock for testing multiCallStakingInfoPermissionless
contract AddressBookV2Mock {
    Profile[] private profiles;
    address public kefAddress;
    address public kifAddress;
    address public kpfAddress;

    function addProfile(
        address nodeId,
        address stakingContract,
        address rewardAddress,
        uint256 timeoutAt,
        State state
    ) external {
        profiles.push(Profile(nodeId, stakingContract, rewardAddress, timeoutAt, state));
    }

    function setFundAddresses(address kef, address kif, address kpf) external {
        kefAddress = kef;
        kifAddress = kif;
        kpfAddress = kpf;
    }

    function getAllProfiles() external view returns (Profile[] memory) {
        return profiles;
    }

    function getFundAddresses() external view returns (address, address, address) {
        return (kefAddress, kifAddress, kpfAddress);
    }
}