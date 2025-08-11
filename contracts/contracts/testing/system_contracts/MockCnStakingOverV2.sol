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

contract MockCnStakingOverV2 {
    string public CONTRACT_TYPE = "CnStakingContract";
    uint256 public VERSION = 2;
    address public nodeId;
    address public rewardAddress;
    address public admin;
    uint256 public staking;
    uint256 public unstaking;

    function mockSetVersion(uint256 _version) external {
        VERSION = _version;
    }

    function mockSetNodeId(address _nodeId) external {
        nodeId = _nodeId;
    }

    function mockSetRewardAddress(address _rewardAddress) external {
        rewardAddress = _rewardAddress;
    }

    function mockSetStaking(uint256 _staking) external {
        staking = _staking;
    }

    function mockSetUnstaking(uint256 _unstaking) external {
        unstaking = _unstaking;
    }

    function mockSetAdmin(address _admin) external {
        admin = _admin;
    }

    function isAdmin(address _admin) external view returns (bool) {
        return _admin == admin;
    }
}
