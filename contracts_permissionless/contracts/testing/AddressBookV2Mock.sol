// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.18;

import {Profile, State} from "../types/Node.sol";

/// @title AddressBookV2Mock
/// @notice Minimal mock for testing multiCallStakingInfoPermissionless
contract AddressBookV2Mock {
    Profile[] private profiles;
    address public kefAddress;
    address public kifAddress;
    address public kpfAddress;
    uint256 public pauseTimeout;
    uint256 public idleTimeout;
    uint256 public maxValidatorCount;
    uint256 public maxReadyCandidateCount;

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

    function setTimeouts(uint256 _pauseTimeout, uint256 _idleTimeout) external {
        pauseTimeout = _pauseTimeout;
        idleTimeout = _idleTimeout;
    }

    function setMaxCounts(uint256 _maxValidatorCount, uint256 _maxReadyCandidateCount) external {
        maxValidatorCount = _maxValidatorCount;
        maxReadyCandidateCount = _maxReadyCandidateCount;
    }

    function getTimeouts() external view returns (uint256, uint256) {
        return (pauseTimeout, idleTimeout);
    }

    function getMaxCounts() external view returns (uint256, uint256) {
        return (maxValidatorCount, maxReadyCandidateCount);
    }
}