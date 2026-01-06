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

contract RemoveStakingTest is BaseTest {
    struct NewStakingEntry {
        address nodeId;
        address stakingContract;
        address rewardAddress;
    }

    address manager;
    NewStakingEntry newStakingEntry;

    modifier submitRequest() {
        bytes memory requestData = _prepareRequestData();
        vm.prank(manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );
        _;
    }

    function _prepareRequestData() internal view returns (bytes memory) {
        return abi.encode(newStakingEntry.nodeId);
    }

    function setUp() public override {
        super.setUp();

        manager = node1.manager;

        newStakingEntry.nodeId = vm.randomAddress();
        (newStakingEntry.stakingContract, newStakingEntry.rewardAddress) = _deployStakingContract(
            manager,
            newStakingEntry.nodeId,
            node1.rewardAddress,
            node1.gcId,
            false,
            true,
            true
        );

        vm.prank(manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.AddStakingContract,
            abi.encode(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress)
        );

        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 1);
    }

    /* ========== REQUEST: SUBMISSION ========== */
    function test_requestRemoveStakingContract_success() public {
        bytes memory requestData = abi.encode(newStakingEntry.nodeId);

        // Check the request
        validatorManager.checkRequest(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );

        vm.expectEmit(true, true, true, true);
        emit RequestSubmitted(
            node1.consensusNodeId,
            manager,
            2,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );

        vm.prank(manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.RemoveStakingContract));
        assertEq(request.requestData, requestData);
    }

    function test_requestRemoveStakingContract_revert_nodeNotFound() public {
        bytes memory requestData = abi.encode(newStakingEntry.nodeId);

        // Use unmanaged consensus node id
        vm.expectRevert(IValidatorManager.NodeNotFound.selector);
        vm.prank(manager);
        validatorManager.request(vm.randomAddress(), IValidatorManager.RequestType.RemoveStakingContract, requestData);
    }

    function test_requestRemoveStakingContract_revert_invalidNodeManager() public {
        bytes memory requestData = abi.encode(newStakingEntry.nodeId);

        // msg.sender is not the manager
        vm.expectRevert(IValidatorManager.InvalidNodeManager.selector);
        vm.prank(vm.randomAddress());
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );
    }

    function test_requestRemoveStakingContract_revert_invalidTargetNodeId() public {
        // 1. Cannot remove consensus node id
        bytes memory requestData = abi.encode(node1.consensusNodeId);
        vm.expectRevert(IValidatorManager.InvalidTargetNodeId.selector);
        vm.prank(manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );

        // 2. Cannot remove unmanaged node id
        requestData = abi.encode(vm.randomAddress());
        vm.expectRevert(IValidatorManager.InvalidTargetNodeId.selector);
        vm.prank(manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.RemoveStakingContract,
            requestData
        );
    }

    /* ========== REQUEST: APPROVAL ========== */
    function test_approveRequestRemoveStakingContract_success() public submitRequest {
        vm.expectEmit(true, true, true, true);
        emit RequestApproved(
            node1.consensusNodeId,
            manager,
            2,
            IValidatorManager.RequestType.RemoveStakingContract,
            _prepareRequestData()
        );

        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 2);

        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
        assertEq(nodeInfo.consensusNodeId, node1.consensusNodeId);
        assertEq(nodeInfo.nodeIds.length, 0);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));

        // Check AddressBook
        vm.expectRevert(abi.encodePacked("Invalid CN node ID."));
        addressBook.getCnInfo(newStakingEntry.nodeId);
    }

    function test_approveRequestRemoveStakingContract_revert_requestNotFound() public submitRequest asAdmin {
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(vm.randomAddress(), 2);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(vm.randomAddress(), 3);

        validatorManager.approveRequest(node1.consensusNodeId, 2);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 2);
    }
}
