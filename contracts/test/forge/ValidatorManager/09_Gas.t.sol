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

import {BaseTest, ValidatorManager, IValidatorManager, IKIP113} from "./base/ValidatorManager.t.sol";

contract ValidatorManagerGasTest is BaseTest {
    TestNode node;
    uint256 constant NODE_COUNT = 40;

    function test_constructor_gas_with_40_nodes() public {
        // Prepare 40 nodes with:
        // - one consensus node id
        // - two additional node ids
        // Total: 40 * 3 = 120 node ids
        IValidatorManager.NodeInfo[] memory nodeInfos = new IValidatorManager.NodeInfo[](NODE_COUNT);
        for (uint256 i = 0; i < NODE_COUNT; i++) {
            node = _deployTestNode(address(0), false, true, true);

            address nodeId1 = vm.randomAddress();
            (address stakingContract1, address rewardAddress1) = _deployStakingContract(
                node.manager,
                nodeId1,
                node.rewardAddress,
                node.gcId,
                false,
                true,
                true
            );

            address nodeId2 = vm.randomAddress();
            address stakingContract2;
            address rewardAddress2;
            if (i % 2 == 0) {
                (stakingContract2, rewardAddress2) = _deployStakingContract(
                    node.manager,
                    nodeId2,
                    node.rewardAddress,
                    node.gcId,
                    false,
                    true,
                    true
                );
            } else {
                (stakingContract2, rewardAddress2) = _deployStakingContractV2(
                    node.manager,
                    nodeId2,
                    node.rewardAddress,
                    node.gcId,
                    true,
                    true
                );
            }

            address[] memory nodeIds = new address[](2);
            nodeIds[0] = nodeId1;
            nodeIds[1] = nodeId2;

            // Register in AddressBook
            vm.startPrank(admin);
            addressBook.submitRegisterCnStakingContract(node.consensusNodeId, node.stakingContract, node.rewardAddress);
            addressBook.submitRegisterCnStakingContract(nodeId1, stakingContract1, rewardAddress1);
            addressBook.submitRegisterCnStakingContract(nodeId2, stakingContract2, rewardAddress2);
            vm.stopPrank();

            vm.prank(address(validatorManager));
            simpleBlsRegistry.register(node.consensusNodeId, node.publicKey, node.pop);

            nodeInfos[i] = IValidatorManager.NodeInfo({
                manager: node.manager,
                consensusNodeId: node.consensusNodeId,
                nodeIds: nodeIds
            });
        }

        // Deploy fresh ValidatorManager
        vm.startSnapshotGas(string(abi.encodePacked("constructor_gas_with_", vm.toString(NODE_COUNT), "*3_nodes")));
        validatorManager = _deployValidatorManager(admin, nodeInfos);
        vm.stopSnapshotGas();

        // Check registration just in case
        for (uint256 i = 0; i < NODE_COUNT; i++) {
            IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(nodeInfos[i].consensusNodeId);
            assertEq(nodeInfo.manager, nodeInfos[i].manager);
            assertEq(nodeInfo.consensusNodeId, nodeInfos[i].consensusNodeId);
            assertEq(nodeInfo.nodeIds.length, nodeInfos[i].nodeIds.length);
            for (uint256 j = 0; j < nodeInfo.nodeIds.length; j++) {
                assertEq(nodeInfo.nodeIds[j], nodeInfos[i].nodeIds[j]);
            }
        }
    }
}
