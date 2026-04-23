// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.18;

import {Profile, State} from "../types/Node.sol";
import {SlotMath} from "../libraries/SlotMath.sol";

/// @title AddressBookV2Mock
/// @notice Minimal mock for testing multiCallStakingInfoPermissionless and multiCallNodeStatesPermissionless
contract AddressBookV2Mock {
    Profile[] private profiles;
    address public kefAddress;
    address public kifAddress;
    address public kpfAddress;
    uint256 public pauseTimeout;
    uint256 public idleTimeout;
    uint256 public maxNodeCount;
    uint256 public maxCandReadyCount;
    uint256 public pfsThreshold;
    uint256 public cfsThreshold;
    uint256 public epochVACount;

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

    function setMaxCounts(uint256 _maxNodeCount, uint256 _maxCandReadyCount) external {
        maxNodeCount = _maxNodeCount;
        maxCandReadyCount = _maxCandReadyCount;
    }

    function getTimeouts() external view returns (uint256, uint256) {
        return (pauseTimeout, idleTimeout);
    }

    function getMaxCounts() external view returns (uint256, uint256) {
        return (maxNodeCount, maxCandReadyCount);
    }

    function setThresholds(uint256 _pfsThreshold, uint256 _cfsThreshold) external {
        pfsThreshold = _pfsThreshold;
        cfsThreshold = _cfsThreshold;
    }

    function setEpochVACount(uint256 _epochVACount) external {
        epochVACount = _epochVACount;
    }

    function getPfsThreshold() external view returns (uint256) {
        return pfsThreshold;
    }

    function getCfsThreshold() external view returns (uint256) {
        return cfsThreshold;
    }

    function getEpochVACount() external view returns (uint256) {
        return epochVACount;
    }

    function getSlotLimits() external view returns (uint256 maxSlotAvailable, uint256 minActiveCount) {
        return (SlotMath.maxSlotAvailable(epochVACount), SlotMath.minActiveCount(epochVACount));
    }

    function getMaxValActivePausedCount() external pure returns (uint256) {
        return 0;
    }

    function getSuspendedValidators() external pure returns (address[] memory) {
        return new address[](0);
    }
}
