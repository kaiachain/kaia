// Copyright 2025 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

interface IAddressBook {
    function getCnInfo(
        address _cnNodeId
    ) external view returns (address cnNodeId, address cnStakingcontract, address cnRewardAddress);
}

interface ICnStaking {
    function staking() external view returns (uint256);
    function unstaking() external view returns (uint256);
}

interface IValidatorState {
    enum State {
        Unknown,
        CandInactive,
        CandReady,
        CandTesting,
        ValInactive,
        ValPaused,
        ValExiting,
        ValReady,
        ValActive
    }

    struct ValidatorState {
        address addr;
        State state;
        uint256 idleTimeout;
        uint256 pausedTimeout;
    }

    event StateTranstiion(State oldStatw, State newState);
    event SetIdleTimeout(address addr, uint256 timestamp);
    event SetPausedTimeout(address addr, uint256 timestamp);

    function setValidatorStates(ValidatorState[] memory) external;
    function getAllValidators() external view returns (ValidatorState[] memory);
}
