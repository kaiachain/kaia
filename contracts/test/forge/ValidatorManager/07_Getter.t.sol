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

import {BaseTest, IValidatorManager} from "./base/ValidatorManager.t.sol";

contract GetterTest is BaseTest {
    function test_getConsensusNodeIds() public {
        address[] memory consensusNodeIds = validatorManager.getConsensusNodeIds();
        assertEq(consensusNodeIds.length, 2);
        assertEq(consensusNodeIds[0], node1.consensusNodeId);
        assertEq(consensusNodeIds[1], node2.consensusNodeId);
    }

    function test_getAllNodeInfos() public {
        IValidatorManager.NodeInfo[] memory nodeInfos = validatorManager.getAllNodeInfos();
        assertEq(nodeInfos.length, 2);
        assertEq(nodeInfos[0].manager, node1.manager);
        assertEq(nodeInfos[0].consensusNodeId, node1.consensusNodeId);
        assertEq(nodeInfos[0].manager, node1.manager);
        assertEq(nodeInfos[0].nodeIds.length, 0);
        assertEq(nodeInfos[1].manager, node2.manager);
        assertEq(nodeInfos[1].consensusNodeId, node2.consensusNodeId);
    }

    function test_getPendingRequests_noRequest() public {
        (address[] memory consensusNodeIds, IValidatorManager.Request[] memory requests) = validatorManager
            .getPendingRequests();
        assertEq(consensusNodeIds.length, 0);
        assertEq(requests.length, 0);
    }

    function test_getPendingRequests_oneRequest() public {
        vm.prank(node1.manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        vm.prank(node2.manager);
        validatorManager.request(node2.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        (address[] memory consensusNodeIds, IValidatorManager.Request[] memory requests) = validatorManager
            .getPendingRequests();
        assertEq(consensusNodeIds.length, 2);
        assertEq(requests.length, 2);

        assertEq(consensusNodeIds[0], node1.consensusNodeId);
        assertEq(requests[0].manager, node1.manager);
        assertEq(requests[0].requestIndex, 1);
        assertEq(uint(requests[0].requestType), uint(IValidatorManager.RequestType.Offboarding));
        assertEq(requests[0].requestData, new bytes(0));

        assertEq(consensusNodeIds[1], node2.consensusNodeId);
        assertEq(requests[1].manager, node2.manager);
        assertEq(requests[1].requestIndex, 2);
        assertEq(uint(requests[1].requestType), uint(IValidatorManager.RequestType.Offboarding));
        assertEq(requests[1].requestData, new bytes(0));
    }

    function test_getDuplicatedNode_consensusNodeId() public {
        address[] memory consensusNodeIds = validatorManager.getDuplicatedNode(vm.randomAddress());
        assertEq(consensusNodeIds.length, 0);

        consensusNodeIds = validatorManager.getDuplicatedNode(node1.consensusNodeId);
        assertEq(consensusNodeIds.length, 1);
        assertEq(consensusNodeIds[0], node1.consensusNodeId);

        consensusNodeIds = validatorManager.getDuplicatedNode(node2.consensusNodeId);
        assertEq(consensusNodeIds.length, 1);
        assertEq(consensusNodeIds[0], node2.consensusNodeId);
    }
}
