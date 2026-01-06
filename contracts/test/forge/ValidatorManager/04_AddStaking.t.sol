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

import {BaseTest, IValidatorManager, ICnStakingV3, ICnStakingV3MultiSig} from "./base/ValidatorManager.t.sol";

contract AddStakingTest is BaseTest {
    struct NewStakingEntry {
        address nodeId;
        address stakingContract;
        address rewardAddress;
    }

    address manager;
    NewStakingEntry newStakingEntry;
    NewStakingEntry newStakingEntryWithPD;
    NewStakingEntry newStakingEntryWithUnmatchedGcId;

    modifier submitRequest() {
        bytes memory requestData = _prepareRequestData();
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
        _;
    }

    function _prepareRequestData() internal view returns (bytes memory) {
        return abi.encode(newStakingEntry.nodeId, newStakingEntry.stakingContract, newStakingEntry.rewardAddress);
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

        newStakingEntryWithPD.nodeId = vm.randomAddress();
        (newStakingEntryWithPD.stakingContract, newStakingEntryWithPD.rewardAddress) = _deployStakingContract(
            manager,
            newStakingEntryWithPD.nodeId,
            node1.rewardAddress,
            node1.gcId,
            true, // isPublicDelegation
            true,
            true
        );

        newStakingEntryWithUnmatchedGcId.nodeId = vm.randomAddress();
        (
            newStakingEntryWithUnmatchedGcId.stakingContract,
            newStakingEntryWithUnmatchedGcId.rewardAddress
        ) = _deployStakingContract(
            manager,
            newStakingEntryWithUnmatchedGcId.nodeId,
            node1.rewardAddress,
            node1.gcId + 1,
            false,
            true,
            true
        );
    }

    /* ========== REQUEST: SUBMISSION ========== */
    function test_requestAddStakingContract_success() public {
        bytes memory requestData = abi.encode(
            newStakingEntry.nodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        // Check the request
        validatorManager.checkRequest(
            node1.consensusNodeId,
            IValidatorManager.RequestType.AddStakingContract,
            requestData
        );

        vm.expectEmit(true, true, true, true);
        emit RequestSubmitted(
            node1.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.AddStakingContract,
            requestData
        );

        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.AddStakingContract));
        assertEq(request.requestData, requestData);
    }

    function test_requestAddStakingContract_revert_nodeNotFound() public {
        bytes memory requestData = abi.encode(
            newStakingEntry.nodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        // Use unmanaged consensus node id
        vm.expectRevert(IValidatorManager.NodeNotFound.selector);
        vm.prank(manager);
        validatorManager.request(vm.randomAddress(), IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_invalidNodeManager() public {
        bytes memory requestData = abi.encode(
            newStakingEntry.nodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        // msg.sender is not the manager
        vm.expectRevert(IValidatorManager.InvalidNodeManager.selector);
        vm.prank(vm.randomAddress());
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_addConsensusNodeId() public {
        // Add consensus node id as a new node id
        bytes memory requestData = abi.encode(
            node1.consensusNodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        vm.expectRevert(_encodeInvalidStakingContractError(0));
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_addAlreadyManagedNodeId() public submitRequest {
        bytes memory requestData = abi.encode(
            newStakingEntry.nodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        // newStakingEntry is registered already
        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 1);

        // Request again with the same data
        vm.expectRevert(_encodeInvalidStakingContractError(0));
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_addPDEnabled() public {
        // Add public delegation enabled staking contract
        bytes memory requestData = abi.encode(
            newStakingEntryWithPD.nodeId,
            newStakingEntryWithPD.stakingContract,
            newStakingEntryWithPD.rewardAddress
        );

        vm.expectRevert(_encodeIncompatibleNewStakingContractError(0));
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_unmatchedGcIdWithConsensusNode() public {
        bytes memory requestData = abi.encode(
            newStakingEntryWithUnmatchedGcId.nodeId,
            newStakingEntryWithUnmatchedGcId.stakingContract,
            newStakingEntryWithUnmatchedGcId.rewardAddress
        );

        vm.expectRevert(_encodeIncompatibleNewStakingContractError(1));
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    function test_requestAddStakingContract_revert_unmatchedRewardAddressWithConsensusNode() public {
        // Change the reward address of consensus node
        address newRewardAddress = vm.randomAddress();
        ICnStakingV3 cnStakingV3 = ICnStakingV3(payable(node1.stakingContract));
        vm.prank(node1.stakingContract);
        cnStakingV3.updateRewardAddress(newRewardAddress);

        vm.prank(newRewardAddress);
        cnStakingV3.acceptRewardAddress(newRewardAddress);

        assertEq(cnStakingV3.rewardAddress(), newRewardAddress);

        bytes memory requestData = abi.encode(
            newStakingEntry.nodeId,
            newStakingEntry.stakingContract,
            newStakingEntry.rewardAddress
        );

        vm.expectRevert(_encodeIncompatibleNewStakingContractError(2));
        vm.prank(manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.AddStakingContract, requestData);
    }

    /* ========== REQUEST: APPROVAL ========== */
    function test_approveRequestAddStakingContract_success() public submitRequest {
        vm.expectEmit(true, true, true, true);
        emit RequestApproved(
            node1.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.AddStakingContract,
            _prepareRequestData()
        );

        vm.prank(admin);
        validatorManager.approveRequest(node1.consensusNodeId, 1);

        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(node1.consensusNodeId);
        assertEq(nodeInfo.nodeIds.length, 1);
        assertEq(nodeInfo.nodeIds[0], newStakingEntry.nodeId);

        IValidatorManager.Request memory request = validatorManager.getRequest(node1.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));

        // Check AddressBook
        (address node, address stakingContract, address rewardAddress) = addressBook.getCnInfo(newStakingEntry.nodeId);
        assertEq(node, newStakingEntry.nodeId);
        assertEq(stakingContract, newStakingEntry.stakingContract);
        assertEq(rewardAddress, newStakingEntry.rewardAddress);
    }

    function test_approveRequestAddStakingContract_revert_requestNotFound() public submitRequest asAdmin {
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(vm.randomAddress(), 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 2);

        validatorManager.approveRequest(node1.consensusNodeId, 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(node1.consensusNodeId, 1);
    }
}
