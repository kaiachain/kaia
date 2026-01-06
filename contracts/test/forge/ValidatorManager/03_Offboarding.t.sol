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

contract OffboardingTest is BaseTest {
    address manager;

    modifier submitRequest() {
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
        _;
    }

    function setUp() public override {
        super.setUp();
        manager = node1.manager;
    }

    /* ========== REQUEST: SUBMISSION ========== */
    function test_requestOffboarding_success() public {
        // Check the request
        validatorManager.checkRequest(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        vm.expectEmit(true, true, true, true);
        emit RequestSubmitted(
            node1.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Offboarding,
            new bytes(0)
        );

        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.Offboarding));
        assertEq(request.requestData, new bytes(0));
    }

    function test_requestOffboarding_revert_nodeNotFound() public {
        // vm.randomAddress() != node1.consensusNodeId
        vm.expectRevert(IValidatorManager.NodeNotFound.selector);
        vm.prank(manager);
        validatorManager.request(vm.randomAddress(), IValidatorManager.RequestType.Offboarding, new bytes(0));
    }

    function test_requestOffboarding_revert_invalidNodeManager() public {
        // msg.sender is not the manager
        vm.expectRevert(IValidatorManager.InvalidNodeManager.selector);
        vm.prank(vm.randomAddress());
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
    }

    function test_requestOffboarding_revert_requestExists() public {
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        vm.expectRevert(IValidatorManager.RequestExists.selector);
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
    }

    /* ========== REQUEST: APPROVAL ========== */
    function test_approveRequestOffboarding_success() public submitRequest {
        vm.expectEmit(true, true, true, true);
        emit RequestApproved(
            node1.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Offboarding,
            new bytes(0)
        );

        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 1);

        TestNode[] memory expectedNodes = new TestNode[](1);
        expectedNodes[0] = node2;
        _verifyManagingNodeIds(expectedNodes);

        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
        assertEq(nodeInfo.manager, address(0));
        assertEq(nodeInfo.consensusNodeId, address(0));
        assertEq(nodeInfo.nodeIds.length, 0);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));

        // Check AddressBook and SimpleBlsRegistry
        vm.expectRevert(abi.encodePacked("Invalid CN node ID."));
        addressBook.getCnInfo(node1.consensusNodeId);

        (bytes memory publicKey, bytes memory pop) = simpleBlsRegistry.record(node1.consensusNodeId);
        assertEq(publicKey, new bytes(0));
        assertEq(pop, new bytes(0));
    }

    function test_approveRequestOffboarding_success_multipleNodes() public {
        // First register two more staking contracts
        address nodeId3 = vm.randomAddress();
        address nodeId4 = vm.randomAddress();
        address cnStaking3;
        address cnStaking4;
        address rewardAddress3;
        address rewardAddress4;
        (cnStaking3, rewardAddress3) = _deployStakingContract(
            manager,
            nodeId3,
            node1.rewardAddress,
            node1.gcId,
            false,
            true,
            true
        );
        (cnStaking4, rewardAddress4) = _deployStakingContract(
            manager,
            nodeId4,
            node1.rewardAddress,
            node1.gcId,
            false,
            true,
            true
        );

        {
            vm.prank(manager);
            validatorManager.request(
                node1.consensusNodeId,
                IValidatorManager.RequestType.AddStakingContract,
                abi.encode(nodeId3, cnStaking3, rewardAddress3)
            );
            vm.prank(admin);
            validatorManager.approveRequest(node1.consensusNodeId, 1);

            vm.prank(manager);
            validatorManager.request(
                node1.consensusNodeId,
                IValidatorManager.RequestType.AddStakingContract,
                abi.encode(nodeId4, cnStaking4, rewardAddress4)
            );
            vm.prank(admin);
            validatorManager.approveRequest(node1.consensusNodeId, 2);

            // Check node info
            IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
            assertEq(nodeInfo.consensusNodeId, node1.consensusNodeId);
            assertEq(nodeInfo.nodeIds.length, 2);
            assertEq(nodeInfo.nodeIds[0], nodeId3);
            assertEq(nodeInfo.nodeIds[1], nodeId4);

            // Check AddressBook
            addressBook.getCnInfo(nodeId3);
            addressBook.getCnInfo(nodeId4);
        }

        // Offboard node will remove all nodes from AddressBook and SimpleBlsRegistry
        {
            vm.prank(manager);
            validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
            vm.prank(admin);
            validatorManager.approveRequest(node1.consensusNodeId, 3);

            // Check node info
            IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
            assertEq(nodeInfo.consensusNodeId, address(0));
            assertEq(nodeInfo.nodeIds.length, 0);

            // Check AddressBook and SimpleBlsRegistry
            vm.expectRevert(abi.encodePacked("Invalid CN node ID."));
            addressBook.getCnInfo(nodeId3);
            vm.expectRevert(abi.encodePacked("Invalid CN node ID."));
            addressBook.getCnInfo(nodeId4);

            (bytes memory publicKey, bytes memory pop) = simpleBlsRegistry.record(node1.consensusNodeId);
            assertEq(publicKey, new bytes(0));
            assertEq(pop, new bytes(0));
        }
    }

    function test_approveRequestOffboarding_revert_ABookNotExecutable() public submitRequest asAdmin {
        // Delete admin from AB
        addressBook.submitDeleteAdmin(address(validatorManager));
        vm.expectRevert(IValidatorManager.ABookNotExecutable.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 1);

        addressBook.submitAddAdmin(address(validatorManager));

        // Set AB requirement to 2
        addressBook.submitUpdateRequirement(2);
        vm.expectRevert(IValidatorManager.ABookNotExecutable.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 1);
    }

    function test_approveRequestOffboarding_revert_SimpleBlsRegistryNotExecutable() public submitRequest {
        vm.prank(address(validatorManager));
        simpleBlsRegistry.transferOwnership(admin);

        vm.expectRevert(IValidatorManager.SimpleBlsRegistryNotExecutable.selector);
        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 1);
    }

    function test_approveRequestOffboarding_revert_requestNotFound() public submitRequest asAdmin {
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(vm.randomAddress(), 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 2);

        validatorManager.approveRequest(node1.consensusNodeId, 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 1);
    }
}
