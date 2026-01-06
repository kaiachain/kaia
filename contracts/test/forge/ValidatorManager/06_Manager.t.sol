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

contract TransferManagerTest is BaseTest {
    modifier submitRequest() {
        vm.prank(node1.manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
        _;
    }

    function test_cancelRequest_success() public submitRequest {
        vm.expectEmit(true, true, true, true);
        emit RequestCancelled(
            node1.consensusNodeId,
            node1.manager,
            1,
            IValidatorManager.RequestType.Offboarding,
            new bytes(0)
        );

        vm.prank(node1.manager);
        validatorManager.cancelRequest(node1.consensusNodeId);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));
    }

    function test_cancelRequest_revert_requestNotExists() public submitRequest {
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        vm.prank(node1.manager);
        validatorManager.cancelRequest(vm.randomAddress());
    }

    function test_cancelRequest_revert_requestNotOwned() public submitRequest {
        vm.expectRevert(IValidatorManager.RequestNotOwned.selector);
        vm.prank(vm.randomAddress());
        validatorManager.cancelRequest(node1.consensusNodeId);
    }

    function test_cancelRequest_revert_requestNotFound() public submitRequest {
        vm.prank(node1.manager);
        validatorManager.cancelRequest(node1.consensusNodeId);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        vm.prank(node1.manager);
        validatorManager.cancelRequest(node1.consensusNodeId);
    }

    function test_transferManager_success() public {
        address newManager = vm.randomAddress();

        vm.expectEmit(true, true, true, true);
        emit NodeManagerTransferred(node1.consensusNodeId, node1.manager, newManager);

        vm.prank(node1.manager);
        validatorManager.transferManager(node1.consensusNodeId, newManager);

        // New manager should be able to manage the nodes
        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
        assertEq(nodeInfo.manager, newManager);
        assertEq(nodeInfo.consensusNodeId, node1.consensusNodeId);
        assertEq(nodeInfo.nodeIds.length, 0);
    }

    function test_transferManager_revert_requestExists() public submitRequest {
        address newManager = vm.randomAddress();

        vm.expectRevert(IValidatorManager.RequestExists.selector);
        vm.prank(node1.manager);
        validatorManager.transferManager(node1.consensusNodeId, newManager);
    }

    function test_transferManager_revert_nodeNotFound() public {
        vm.expectRevert(IValidatorManager.NodeNotFound.selector);
        vm.prank(node1.manager);
        validatorManager.transferManager(vm.randomAddress(), vm.randomAddress());
    }

    function test_transferManager_revert_invalidNodeManager() public {
        vm.expectRevert(IValidatorManager.InvalidNodeManager.selector);
        vm.prank(vm.randomAddress());
        validatorManager.transferManager(node1.consensusNodeId, vm.randomAddress());
    }

    function test_transferManager_revert_invalidNewManager() public {
        vm.startPrank(node1.manager);
        // 1. Zero address
        vm.expectRevert(IValidatorManager.InvalidNewManager.selector);
        validatorManager.transferManager(node1.consensusNodeId, address(0));

        // 2. Current manager
        vm.expectRevert(IValidatorManager.InvalidNewManager.selector);
        validatorManager.transferManager(node1.consensusNodeId, node1.manager);
    }
}
