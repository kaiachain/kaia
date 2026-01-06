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

interface ICnStakingV3MultiSigFactory {
    event CnStakingV3MultiSigDeployed(address indexed deployer, address indexed cnStakingV3MultiSig);

    error DeploymentFailed();

    function deployCnStakingV3MultiSig(
        address _contractValidator,
        address _nodeId,
        address _rewardAddress,
        address[] memory _cnAdminlist,
        uint256 _requirement,
        uint256[] memory _unlockTime,
        uint256[] memory _unlockAmount
    ) external returns (address cnStakingV3MultiSig);

    function getDeployedCnStakingV3MultiSigList(address _deployer) external view returns (address[] memory);

    function isDeployedBy(address _deployer, address _cnStakingV3MultiSig) external view returns (bool);
}
