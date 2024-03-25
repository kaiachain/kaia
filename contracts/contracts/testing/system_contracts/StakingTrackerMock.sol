// Copyright 2022 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

contract StakingTrackerMockReceiver {
    event RefreshStake();
    event RefreshVoter();

    function refreshStake(address) external {
        emit RefreshStake();
    }

    function refreshVoter(address) external {
        emit RefreshVoter();
    }

    function CONTRACT_TYPE() external pure returns (string memory) {
        return "StakingTracker";
    }

    function VERSION() external pure returns (uint256) {
        return 1;
    }

    function voterToGCId(address) external pure returns (uint256) {
        return 0;
    }

    function getLiveTrackerIds() external pure returns (uint256[] memory) {
        return new uint256[](0);
    }
}

contract StakingTrackerMockActive {
    event RefreshStake();

    function CONTRACT_TYPE() external pure returns (string memory) {
        return "StakingTracker";
    }

    function VERSION() external pure returns (uint256) {
        return 1;
    }

    function refreshStake(address) external {
        emit RefreshStake();
    }

    function getLiveTrackerIds() external pure returns (uint256[] memory) {
        return new uint256[](1);
    }
}

contract StakingTrackerMockWrong {
    function CONTRACT_TYPE() external pure returns (string memory) {
        return "Wrong";
    }

    function VERSION() external pure returns (uint256) {
        return 1;
    }
}

contract StakingTrackerMockInvalid {
    function CONTRACT_TYPE() external pure returns (string memory) {
        return "";
    }
    // no VERSION() function
}
