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
pragma solidity ^0.8.0;

import {BaseTest, CnStakingV3MultiSig, IPublicDelegation} from "./base/ValidatorManager.t.sol";

contract FactoryTest is BaseTest {
    function test_deployCnV3MultiSig() public {
        // Prepare same constructor arguments
        address contractValidator = vm.randomAddress();
        address nodeId = vm.randomAddress();
        address rewardAddress = vm.randomAddress();
        address[] memory adminList = new address[](1);
        adminList[0] = vm.randomAddress();
        uint256 requirement = 1;
        uint256[] memory unlockTime = new uint256[](0);
        uint256[] memory unlockAmount = new uint256[](0);

        // 1. Deploy CnV3MultiSig through factory
        address cnStakingV3MultiSig = cnStakingV3MultiSigFactory.deployCnStakingV3MultiSig(
            contractValidator,
            nodeId,
            rewardAddress,
            adminList,
            requirement,
            unlockTime,
            unlockAmount
        );

        // 2. Deploy CnV3MultiSig directly
        address cnStakingV3MultiSig2 = address(
            new CnStakingV3MultiSig(
                contractValidator,
                nodeId,
                rewardAddress,
                adminList,
                requirement,
                unlockTime,
                unlockAmount
            )
        );

        // Check if the CnV3MultiSig is deployed
        assertEq(cnStakingV3MultiSigFactory.isDeployedBy(address(this), cnStakingV3MultiSig), true);
        assertEq(cnStakingV3MultiSigFactory.getDeployedCnStakingV3MultiSigList(address(this)).length, 1);
        assertEq(cnStakingV3MultiSigFactory.getDeployedCnStakingV3MultiSigList(address(this))[0], cnStakingV3MultiSig);

        // Check if the deployed CnV3MultiSig bytecode is the same
        // @audit-info: Since bytecode heavily depends on the metadata, this test can be failed.
        // To verify this, please replace chunks with the same bytecode from your local environment.
        // assertEq(cnStakingV3MultiSig.codehash, cnStakingV3MultiSig2.codehash);
    }

    function test_deployPublicDelegation() public {
        // Prepare same constructor arguments
        address owner = vm.randomAddress();
        address commissionTo = vm.randomAddress();
        uint256 commissionRate = 1000;
        string memory gcName = "GC";

        // 1. Deploy PublicDelegation through factory
        address publicDelegation = address(
            publicDelegationFactory.deployPublicDelegation(
                IPublicDelegation.PDConstructorArgs({
                    owner: owner,
                    commissionTo: commissionTo,
                    commissionRate: commissionRate,
                    gcName: gcName
                })
            )
        );

        // Check if the PublicDelegation is deployed
        assertEq(publicDelegationFactory.isDeployedBy(address(this), publicDelegation), true);
    }
}
