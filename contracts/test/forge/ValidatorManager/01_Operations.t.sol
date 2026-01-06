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

import {BaseTest, IValidatorManager, IKIP113, ISimpleBlsRegistry, IAddressBook} from "./base/ValidatorManager.t.sol";

contract OperationTest is BaseTest {
    struct NewStakingEntry {
        address nodeId;
        address stakingContract;
        address rewardAddress;
    }

    TestNode onboardingNode;
    NewStakingEntry newStakingEntry;
    NewStakingEntry newStakingEntryWithPD;

    modifier submitRequest() {
        // Request add staking contract from node 1, offboarding request from node 2
        vm.prank(node1.manager);
        validatorManager.request(
            node1.consensusNodeId,
            IValidatorManager.RequestType.AddStakingContract,
            abi.encode(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress)
        );
        vm.prank(node2.manager);
        validatorManager.request(node2.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));
        _;
    }

    function setUp() public override {
        super.setUp();

        onboardingNode = _deployTestNode(address(0), false, true, false);

        newStakingEntry.nodeId = vm.randomAddress();
        (newStakingEntry.stakingContract, newStakingEntry.rewardAddress) = _deployStakingContract(
            node1.manager,
            newStakingEntry.nodeId,
            node1.rewardAddress,
            node1.gcId,
            false,
            true,
            true
        );

        newStakingEntryWithPD.nodeId = vm.randomAddress();
        (newStakingEntryWithPD.stakingContract, newStakingEntryWithPD.rewardAddress) = _deployStakingContract(
            node1.manager,
            newStakingEntryWithPD.nodeId,
            node1.rewardAddress,
            node1.gcId,
            true, // isPublicDelegation
            true,
            true
        );
    }

    function _registerInABook(address nodeId, address stakingContract, address rewardAddress) internal {
        vm.prank(admin);
        addressBook.submitRegisterCnStakingContract(nodeId, stakingContract, rewardAddress);
    }

    function _registerInSBR(address nodeId, bytes memory publicKey, bytes memory pop) internal {
        vm.prank(address(validatorManager));
        simpleBlsRegistry.register(nodeId, publicKey, pop);
    }

    /* ========== TEST: MANAGE ========== */
    function test_manageAddressBook_success() public asAdmin {
        address newAdmin = vm.randomAddress();
        bytes memory data = abi.encodeWithSelector(IAddressBook.submitAddAdmin.selector, newAdmin);

        validatorManager.manage(ADDRESS_BOOK_ADDR, data);

        // Check if new admin is added
        (address[] memory adminList, ) = addressBook.getState();
        assertEq(adminList.length, 3);
        assertEq(adminList[2], newAdmin);
    }

    function test_manageSimpleBlsRegistry_success() public asAdmin {
        address newAdmin = vm.randomAddress();
        bytes memory data = abi.encodeWithSignature("transferOwnership(address)", newAdmin);

        validatorManager.manage(address(simpleBlsRegistry), data);

        // Check if new admin is added
        assertEq(simpleBlsRegistry.owner(), newAdmin);
    }

    function test_manage_revert_onlyOwner() public {
        vm.expectRevert(abi.encodeWithSelector(OwnableUnauthorizedAccount.selector, vm.addr(1)));
        vm.prank(vm.addr(1));
        validatorManager.manage(ADDRESS_BOOK_ADDR, new bytes(0));
    }

    /* ========== TEST: ADD NODE INFO ========== */
    function test_addNodeInfo_success() public {
        _registerInABook(onboardingNode.consensusNodeId, onboardingNode.stakingContract, onboardingNode.rewardAddress);
        _registerInSBR(onboardingNode.consensusNodeId, onboardingNode.publicKey, onboardingNode.pop);

        address consensusNodeId = onboardingNode.consensusNodeId;
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: onboardingNode.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: new address[](0)
        });

        vm.expectEmit(true, true, true, true);
        emit NodeInfoAdded(consensusNodeId);
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);

        IValidatorManager.NodeInfo memory added = validatorManager.getNodeInfo(consensusNodeId);
        assertEq(added.manager, onboardingNode.manager);
        assertEq(added.consensusNodeId, consensusNodeId);
        assertEq(added.nodeIds.length, 0);
    }

    function test_addNodeInfo_override_success() public {
        _registerInABook(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress);

        address[] memory nodeIds = new address[](1);
        nodeIds[0] = newStakingEntry.nodeId;
        address consensusNodeId = node1.consensusNodeId;
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: node1.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: nodeIds
        });

        // Submit an offboarding request, will be deleted
        vm.prank(node1.manager);
        validatorManager.request(consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        vm.expectEmit(true, true, true, true);
        emit NodeInfoAdded(consensusNodeId);
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);

        IValidatorManager.NodeInfo memory added = validatorManager.getNodeInfo(consensusNodeId);
        assertEq(added.manager, node1.manager);
        assertEq(added.consensusNodeId, consensusNodeId);
        assertEq(added.nodeIds.length, nodeIds.length);
        for (uint256 i = 0; i < added.nodeIds.length; i++) {
            assertEq(added.nodeIds[i], nodeIds[i]);
        }

        IValidatorManager.Request memory request = validatorManager.getRequest(consensusNodeId);
        assertEq(uint256(request.requestType), uint256(IValidatorManager.RequestType.NoRequest));
    }

    function test_addNodeInfo_revert_zeroManager() public asAdmin {
        IValidatorManager.NodeInfo memory invalidNodeInfo = IValidatorManager.NodeInfo({
            manager: address(0),
            consensusNodeId: onboardingNode.consensusNodeId,
            nodeIds: new address[](0)
        });
        vm.expectRevert(_encodeInvalidNodeInfoError(0));
        validatorManager.addNodeInfo(invalidNodeInfo);
    }

    function test_addNodeInfo_revert_invalidNodeIds() public asAdmin {
        address[] memory zeroNodeIds = new address[](1);
        address[] memory consensusNodeIdNodeIds = new address[](1);
        consensusNodeIdNodeIds[0] = node1.consensusNodeId;
        address[] memory duplicateNodeIds = new address[](2);
        duplicateNodeIds[0] = vm.randomAddress();
        duplicateNodeIds[1] = duplicateNodeIds[0];
        address[] memory validNodeIds = new address[](1);
        validNodeIds[0] = vm.randomAddress();

        // 1. Consensus node id is zero address
        IValidatorManager.NodeInfo memory invalidNodeInfo3 = IValidatorManager.NodeInfo({
            manager: admin,
            consensusNodeId: address(0),
            nodeIds: validNodeIds
        });
        vm.expectRevert(_encodeInvalidNodeInfoError(1));
        validatorManager.addNodeInfo(invalidNodeInfo3);

        // 2. Zero address is in node ids
        IValidatorManager.NodeInfo memory invalidNodeInfo1 = IValidatorManager.NodeInfo({
            manager: admin,
            consensusNodeId: node1.consensusNodeId,
            nodeIds: zeroNodeIds
        });
        vm.expectRevert(_encodeInvalidNodeInfoError(2));
        validatorManager.addNodeInfo(invalidNodeInfo1);

        // 3. Consensus node id is in node ids
        IValidatorManager.NodeInfo memory invalidNodeInfo2 = IValidatorManager.NodeInfo({
            manager: admin,
            consensusNodeId: node1.consensusNodeId,
            nodeIds: consensusNodeIdNodeIds
        });
        vm.expectRevert(_encodeInvalidNodeInfoError(3));
        validatorManager.addNodeInfo(invalidNodeInfo2);

        // 4. Node ids are not distinct
        IValidatorManager.NodeInfo memory invalidNodeInfo4 = IValidatorManager.NodeInfo({
            manager: admin,
            consensusNodeId: node1.consensusNodeId,
            nodeIds: duplicateNodeIds
        });
        vm.expectRevert(_encodeInvalidNodeInfoError(4));
        validatorManager.addNodeInfo(invalidNodeInfo4);
    }

    function test_addNodeInfo_revert_notEnoughBalance() public {
        address consensusNodeId = onboardingNode.consensusNodeId;
        // Set the balance to 0
        vm.deal(consensusNodeId, 0);
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: onboardingNode.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: new address[](0)
        });

        vm.expectRevert(_encodeInvalidNodeInfoError(5));
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);
    }

    function test_addNodeInfo_revert_notInABook() public {
        address consensusNodeId = onboardingNode.consensusNodeId;
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: onboardingNode.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: new address[](0)
        });

        vm.expectRevert(_encodeInvalidNodeInfoError(6));
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);
    }

    function test_addNodeInfo_revert_notInSBR() public {
        _registerInABook(onboardingNode.consensusNodeId, onboardingNode.stakingContract, onboardingNode.rewardAddress);

        address consensusNodeId = onboardingNode.consensusNodeId;
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: onboardingNode.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: new address[](0)
        });

        vm.expectRevert(_encodeInvalidNodeInfoError(7));
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);
    }

    function test_addNodeInfo_revert_notCompatibleStakingContract() public {
        _registerInABook(
            newStakingEntryWithPD.nodeId,
            newStakingEntryWithPD.stakingContract,
            newStakingEntryWithPD.rewardAddress
        );

        address[] memory nodeIds = new address[](1);
        nodeIds[0] = newStakingEntryWithPD.nodeId;
        address consensusNodeId = node1.consensusNodeId;
        IValidatorManager.NodeInfo memory newNodeInfo = IValidatorManager.NodeInfo({
            manager: node1.manager,
            consensusNodeId: consensusNodeId,
            nodeIds: nodeIds
        });

        vm.expectRevert(_encodeIncompatibleNewStakingContractError(0));
        vm.prank(admin);
        validatorManager.addNodeInfo(newNodeInfo);
    }

    /* ========== TEST: REMOVE NODE INFO ========== */
    function test_removeNodeInfo_success() public {
        // Submit an onboarding request, will be deleted
        vm.prank(node1.manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Offboarding, new bytes(0));

        vm.expectEmit(true, true, true, true);
        emit NodeInfoRemoved(node1.consensusNodeId);
        vm.prank(admin);
        validatorManager.removeNodeInfo(node1.consensusNodeId);

        IValidatorManager.NodeInfo memory removed = validatorManager.getNodeInfo(node1.consensusNodeId);
        assertEq(removed.manager, address(0));
        assertEq(removed.consensusNodeId, address(0));
        assertEq(removed.nodeIds.length, 0);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint256(request.requestType), uint256(IValidatorManager.RequestType.NoRequest));
    }

    function test_removeNodeInfo_revert_nodeNotFound() public {
        vm.expectRevert(IValidatorManager.NodeNotFound.selector);
        vm.prank(admin);
        validatorManager.removeNodeInfo(vm.randomAddress());
    }

    /* ========== TEST: BATCH OPERATIONS ========== */
    function test_approveBatchRequest_success() public submitRequest {
        address[] memory consensusNodeIds = new address[](2);
        consensusNodeIds[0] = node1.consensusNodeId;
        consensusNodeIds[1] = node2.consensusNodeId;
        uint256[] memory requestIndices = new uint256[](2);
        requestIndices[0] = 1;
        requestIndices[1] = 2;

        vm.expectEmit(true, true, true, true);
        emit RequestApproved(
            node1.consensusNodeId,
            node1.manager,
            1,
            IValidatorManager.RequestType.AddStakingContract,
            abi.encode(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress)
        );
        emit RequestApproved(
            node2.consensusNodeId,
            node2.manager,
            2,
            IValidatorManager.RequestType.Offboarding,
            new bytes(0)
        );
        vm.prank(admin);
        validatorManager.approveBatchRequest(consensusNodeIds, requestIndices);

        // All nodes are offboarded
        TestNode[] memory expectedNodes = new TestNode[](1);
        expectedNodes[0] = node1;
        _verifyManagingNodeIds(expectedNodes);

        // No pending requests
        (
            address[] memory pendingConsensusNodeIds,
            IValidatorManager.Request[] memory pendingRequests
        ) = validatorManager.getPendingRequests();
        assertEq(pendingConsensusNodeIds.length, 0);
        assertEq(pendingRequests.length, 0);
    }

    function test_approveBatchRequest_revert_invalidBatchRequestLength() public submitRequest {
        address[] memory consensusNodeIds = new address[](2);
        consensusNodeIds[0] = node1.consensusNodeId;
        consensusNodeIds[1] = node2.consensusNodeId;
        uint256[] memory requestIndices = new uint256[](1);
        requestIndices[0] = 1;
        vm.expectRevert(IValidatorManager.InvalidBatchRequestLength.selector);
        vm.prank(admin);
        validatorManager.approveBatchRequest(consensusNodeIds, requestIndices);
    }

    function test_approveBatchRequest_revert_requestNotExists() public submitRequest {
        address[] memory consensusNodeIds = new address[](1);
        consensusNodeIds[0] = vm.randomAddress();
        uint256[] memory requestIndices = new uint256[](1);
        requestIndices[0] = 1;
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        vm.prank(admin);
        validatorManager.approveBatchRequest(consensusNodeIds, requestIndices);
    }

    function test_rejectBatchRequest_success() public submitRequest {
        address[] memory consensusNodeIds = new address[](2);
        consensusNodeIds[0] = node1.consensusNodeId;
        consensusNodeIds[1] = node2.consensusNodeId;
        uint256[] memory requestIndices = new uint256[](2);
        requestIndices[0] = 1;
        requestIndices[1] = 2;

        vm.expectEmit(true, true, true, true);
        emit RequestRejected(
            node1.consensusNodeId,
            node1.manager,
            1,
            IValidatorManager.RequestType.AddStakingContract,
            abi.encode(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress)
        );
        emit RequestRejected(
            node2.consensusNodeId,
            node2.manager,
            2,
            IValidatorManager.RequestType.Offboarding,
            new bytes(0)
        );
        vm.prank(admin);
        validatorManager.rejectBatchRequest(consensusNodeIds, requestIndices);

        // No pending requests
        (
            address[] memory pendingConsensusNodeIds,
            IValidatorManager.Request[] memory pendingRequests
        ) = validatorManager.getPendingRequests();
        assertEq(pendingConsensusNodeIds.length, 0);
    }

    function test_rejectBatchRequest_revert_invalidBatchRequestLength() public submitRequest {
        address[] memory consensusNodeIds = new address[](2);
        consensusNodeIds[0] = node1.consensusNodeId;
        consensusNodeIds[1] = node2.consensusNodeId;
        uint256[] memory requestIndices = new uint256[](1);
        requestIndices[0] = 1;
        vm.expectRevert(IValidatorManager.InvalidBatchRequestLength.selector);
        vm.prank(admin);
        validatorManager.rejectBatchRequest(consensusNodeIds, requestIndices);
    }

    function test_rejectBatchRequest_revert_requestNotExists() public submitRequest {
        address[] memory consensusNodeIds = new address[](1);
        consensusNodeIds[0] = vm.randomAddress();
        uint256[] memory requestIndices = new uint256[](1);
        requestIndices[0] = 1;
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        vm.prank(admin);
        validatorManager.rejectBatchRequest(consensusNodeIds, requestIndices);
    }
}
