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

import {BaseTest, IValidatorManager, IKIP113} from "./base/ValidatorManager.t.sol";

contract ConstructorTest is BaseTest {
    function test_constants() public view {
        assertEq(address(validatorManager.ABOOK()), ADDRESS_BOOK_ADDR);
        assertEq(address(validatorManager.REGISTRY()), REGISTRY_ADDR);

        assertEq(validatorManager.MIN_CONSENSUS_NODE_BALANCE(), 10 ether);

        assertEq(validatorManager.ZERO48HASH(), keccak256(new bytes(48)));
        assertEq(validatorManager.ZERO96HASH(), keccak256(new bytes(96)));
    }

    function test_constructor() public view {
        assertEq(address(validatorManager.owner()), admin);

        IValidatorManager.NodeInfo memory nodeInfo1 = validatorManager.getNodeInfo(node1.consensusNodeId);
        IValidatorManager.NodeInfo memory nodeInfo2 = validatorManager.getNodeInfo(node2.consensusNodeId);

        assertEq(nodeInfo1.manager, node1.manager);
        assertEq(nodeInfo1.consensusNodeId, node1.consensusNodeId);
        assertEq(nodeInfo1.nodeIds.length, 0);

        assertEq(nodeInfo2.manager, node2.manager);
        assertEq(nodeInfo2.consensusNodeId, node2.consensusNodeId);
        assertEq(nodeInfo2.nodeIds.length, 0);
    }

    function test_checkAddressBook() public view {
        (address nodeId, address stakingContract, address rewardAddress) = addressBook.getCnInfo(node1.consensusNodeId);
        assertEq(nodeId, node1.consensusNodeId);
        assertEq(stakingContract, address(node1.stakingContract));
        assertEq(rewardAddress, node1.rewardAddress);

        (nodeId, stakingContract, rewardAddress) = addressBook.getCnInfo(node2.consensusNodeId);
        assertEq(nodeId, node2.consensusNodeId);
        assertEq(stakingContract, address(node2.stakingContract));
        assertEq(rewardAddress, node2.rewardAddress);
    }

    function test_checkSimpleBlsRegistry() public view {
        (address[] memory nodeIdList, IKIP113.BlsPublicKeyInfo[] memory pubkeyList) = simpleBlsRegistry.getAllBlsInfo();

        assertEq(nodeIdList.length, 2);
        assertEq(pubkeyList.length, 2);
        assertEq(nodeIdList[0], node1.consensusNodeId);
        assertEq(pubkeyList[0].publicKey, node1.publicKey);
        assertEq(pubkeyList[0].pop, node1.pop);
        assertEq(nodeIdList[1], node2.consensusNodeId);
        assertEq(pubkeyList[1].publicKey, node2.publicKey);
        assertEq(pubkeyList[1].pop, node2.pop);
    }
}
