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

import {BaseTest, PublicDelegationFactoryV2, ICnStakingV3MultiSig, IValidatorManager, IKIP113} from "./base/ValidatorManager.t.sol";

contract OnboardingTest is BaseTest {
    TestNode onboardingNode;
    TestNode uninitializedNode;
    TestNode customPublicDelegationNode;
    address manager;

    modifier submitRequest() {
        _stake(manager, address(onboardingNode.stakingContract), 5_000_000 ether);
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);
        vm.prank(manager);
        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
        _;
    }

    function setUp() public override {
        super.setUp();

        onboardingNode = _deployTestNode(address(0), false, true, false);
        uninitializedNode = _deployTestNode(address(0), false, false, true);

        // Temporarily replace the public delegation factory with the custom one
        PublicDelegationFactoryV2 customPublicDelegationFactory = new PublicDelegationFactoryV2();
        PublicDelegationFactoryV2 originalPublicDelegationFactory = publicDelegationFactory;
        publicDelegationFactory = customPublicDelegationFactory;

        customPublicDelegationNode = _deployTestNode(address(0), true, true, true);

        publicDelegationFactory = originalPublicDelegationFactory;

        manager = onboardingNode.manager;
    }

    function _sendOnboardingRequest(address consensusNodeId, bytes memory requestData) internal {
        vm.prank(manager);
        validatorManager.request(consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
    }

    /* ========== REQUEST: SUBMISSION ========== */
    function test_requestOnboarding_success() public {
        // Prepare onboarding request
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);

        vm.startPrank(manager);

        // Check the request
        validatorManager.checkRequest(
            onboardingNode.consensusNodeId,
            IValidatorManager.RequestType.Onboarding,
            requestData
        );

        vm.expectEmit(true, true, true, true);
        emit RequestSubmitted(
            onboardingNode.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Onboarding,
            requestData
        );

        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);

        IValidatorManager.Request memory request = validatorManager.getRequest(onboardingNode.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.Onboarding));
        assertEq(request.requestData, requestData);

        vm.stopPrank();
    }

    function test_cancelRequestOnboarding_success() public submitRequest {
        // This test is required since during onboarding, the _nodeInfo is not set yet.
        vm.expectEmit(true, true, true, true);
        emit RequestCancelled(
            onboardingNode.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Onboarding,
            _prepareOnboardingRequest(onboardingNode)
        );

        vm.prank(manager);
        validatorManager.cancelRequest(onboardingNode.consensusNodeId);

        IValidatorManager.Request memory request = validatorManager.getRequest(onboardingNode.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));
    }

    function test_requestOnboarding_revert_wrongRequestData() public {
        // Wrong request data
        bytes memory requestData = abi.encode(vm.randomBytes(32));
        vm.expectRevert();
        vm.prank(manager);
        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
    }

    function test_requestOnboarding_revert_invalidNodeManager() public {
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);
        vm.expectRevert(IValidatorManager.InvalidNodeManager.selector);
        // msg.sender is not the manager
        vm.prank(vm.randomAddress());
        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
    }

    function test_requestOnboarding_revert_existingRequest() public {
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);

        vm.prank(manager);
        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);

        // Second request from same manager
        vm.expectRevert(IValidatorManager.RequestExists.selector);
        vm.prank(manager);
        validatorManager.request(onboardingNode.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
    }

    function test_requestOnboarding_revert_invalidBLSKey() public {
        bytes memory requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            new bytes(47), // Invalid length
            onboardingNode.pop
        );
        vm.expectRevert(_encodeInvalidBLSKeyInfoError(0));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);

        requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            onboardingNode.publicKey,
            new bytes(95) // Invalid length
        );
        vm.expectRevert(_encodeInvalidBLSKeyInfoError(0));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);

        requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            new bytes(48), // Zero public key
            onboardingNode.pop
        );
        vm.expectRevert(_encodeInvalidBLSKeyInfoError(1));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);

        requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            onboardingNode.publicKey,
            new bytes(96) // Zero pop
        );
        vm.expectRevert(_encodeInvalidBLSKeyInfoError(1));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);
    }

    function test_requestOnboarding_revert_nodeAlreadyExists() public {
        // Use node1.consensusNodeId, which is already onboarded
        bytes memory requestData = _prepareCustomOnboardingRequest(
            node1.manager,
            node1.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            onboardingNode.publicKey,
            onboardingNode.pop
        );

        vm.expectRevert(IValidatorManager.NodeAlreadyExists.selector);
        vm.prank(node1.manager);
        validatorManager.request(node1.consensusNodeId, IValidatorManager.RequestType.Onboarding, requestData);
    }

    function test_requestOnboarding_revert_insufficientConsensusNodeIdBalance() public {
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);
        // Set the balance to 0
        vm.deal(onboardingNode.consensusNodeId, 0);

        vm.expectRevert(IValidatorManager.BelowMinNodeIdBalance.selector);
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);
    }

    function test_requestOnboarding_revert_nodeAlreadyInABook() public {
        bytes memory requestData = _prepareOnboardingRequest(onboardingNode);

        // Register the consensus node id in AddressBook
        vm.prank(admin);
        addressBook.submitRegisterCnStakingContract(
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress)
        );

        vm.expectRevert(_encodeInvalidStakingContractError(0));
        _sendOnboardingRequest(node1.consensusNodeId, requestData);
    }

    function test_requestOnboarding_revert_notDeployedByManager() public {
        // Use uninitializedNode.stakingContract, which is not deployed by the manager
        bytes memory requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(uninitializedNode.stakingContract),
            address(onboardingNode.rewardAddress),
            onboardingNode.publicKey,
            onboardingNode.pop
        );
        vm.expectRevert(_encodeInvalidStakingContractError(1));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);
    }

    function test_requestOnboarding_revert_notInitialized() public {
        bytes memory requestData = _prepareCustomOnboardingRequest(
            uninitializedNode.manager,
            uninitializedNode.consensusNodeId,
            address(uninitializedNode.stakingContract),
            address(uninitializedNode.rewardAddress),
            uninitializedNode.publicKey,
            uninitializedNode.pop
        );
        vm.expectRevert(_encodeInvalidStakingContractError(2));
        vm.prank(uninitializedNode.manager);
        validatorManager.request(
            uninitializedNode.consensusNodeId,
            IValidatorManager.RequestType.Onboarding,
            requestData
        );
    }

    function test_requestOnboarding_revert_unmatchedConsensusNodeId() public {
        address randomNodeId = vm.randomAddress();
        vm.deal(randomNodeId, 10 ether);

        bytes memory requestData = _prepareCustomOnboardingRequest(
            manager,
            randomNodeId,
            address(onboardingNode.stakingContract),
            address(onboardingNode.rewardAddress),
            onboardingNode.publicKey,
            onboardingNode.pop
        );
        vm.expectRevert(_encodeInvalidStakingContractError(3));
        _sendOnboardingRequest(randomNodeId, requestData);
    }

    function test_requestOnboarding_revert_unmatchedRewardAddress() public {
        // Onboarding node has a different reward address with node1
        bytes memory requestData = _prepareCustomOnboardingRequest(
            manager,
            onboardingNode.consensusNodeId,
            address(onboardingNode.stakingContract),
            address(node1.rewardAddress),
            onboardingNode.publicKey,
            onboardingNode.pop
        );
        vm.expectRevert(_encodeInvalidStakingContractError(4));
        _sendOnboardingRequest(onboardingNode.consensusNodeId, requestData);
    }

    function test_requestOnboarding_revert_notDeployedByPublicDelegationFactory() public {
        bytes memory requestData = _prepareCustomOnboardingRequest(
            customPublicDelegationNode.manager,
            customPublicDelegationNode.consensusNodeId,
            address(customPublicDelegationNode.stakingContract),
            address(customPublicDelegationNode.rewardAddress),
            customPublicDelegationNode.publicKey,
            customPublicDelegationNode.pop
        );
        vm.expectRevert(_encodeInvalidStakingContractError(5));
        vm.prank(customPublicDelegationNode.manager);
        validatorManager.request(
            customPublicDelegationNode.consensusNodeId,
            IValidatorManager.RequestType.Onboarding,
            requestData
        );
    }

    /* ========== REQUEST: APPROVAL & REJECTION ========== */
    function test_approveRequestOnboarding_success() public submitRequest {
        vm.expectEmit(true, true, true, true);
        emit RequestApproved(
            onboardingNode.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Onboarding,
            _prepareOnboardingRequest(onboardingNode)
        );

        vm.prank(admin);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);

        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(onboardingNode.consensusNodeId);
        assertEq(nodeInfo.manager, manager);
        assertEq(nodeInfo.consensusNodeId, onboardingNode.consensusNodeId);
        assertEq(nodeInfo.nodeIds.length, 0);

        TestNode[] memory expectedNodes = new TestNode[](3);
        expectedNodes[0] = node1;
        expectedNodes[1] = node2;
        expectedNodes[2] = onboardingNode;
        _verifyManagingNodeIds(expectedNodes);

        IValidatorManager.Request memory request = validatorManager.getRequest(onboardingNode.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));

        // Check AddressBook and SimpleBlsRegistry
        (address node, address stakingContract, address rewardAddress) = addressBook.getCnInfo(
            onboardingNode.consensusNodeId
        );
        assertEq(node, onboardingNode.consensusNodeId);
        assertEq(stakingContract, address(onboardingNode.stakingContract));
        assertEq(rewardAddress, address(onboardingNode.rewardAddress));

        (bytes memory publicKey, bytes memory pop) = simpleBlsRegistry.record(onboardingNode.consensusNodeId);
        assertEq(publicKey, onboardingNode.publicKey);
        assertEq(pop, onboardingNode.pop);
    }

    function test_approveRequestOnboarding_revert_onlyOwner() public submitRequest {
        vm.expectRevert(abi.encodeWithSelector(OwnableUnauthorizedAccount.selector, vm.addr(1)));
        vm.prank(vm.addr(1));
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);
    }

    function test_approveRequestOnboarding_revert_ABookNotExecutable() public submitRequest asAdmin {
        // Delete admin from AB
        addressBook.submitDeleteAdmin(address(validatorManager));
        vm.expectRevert(IValidatorManager.ABookNotExecutable.selector);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);

        addressBook.submitAddAdmin(address(validatorManager));

        // Set AB requirement to 2
        addressBook.submitUpdateRequirement(2);
        vm.expectRevert(IValidatorManager.ABookNotExecutable.selector);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);
    }

    function test_approveRequestOnboarding_revert_SimpleBlsRegistryNotExecutable() public submitRequest {
        // Transfer ownership to admin
        vm.prank(address(validatorManager));
        simpleBlsRegistry.transferOwnership(admin);

        vm.expectRevert(IValidatorManager.SimpleBlsRegistryNotExecutable.selector);
        vm.prank(admin);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);
    }

    function test_approveRequestOnboarding_revert_requestNotFound() public submitRequest asAdmin {
        // No such manager
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(vm.randomAddress(), 1);

        // Not matched request index
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 2);

        // Approve the request
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);

        // Already approved
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.approveRequest(onboardingNode.consensusNodeId, 1);
    }

    function test_rejectRequestOnboarding_success() public submitRequest asAdmin {
        vm.expectEmit(true, true, true, true);
        emit RequestRejected(
            onboardingNode.consensusNodeId,
            manager,
            1,
            IValidatorManager.RequestType.Onboarding,
            _prepareOnboardingRequest(onboardingNode)
        );

        validatorManager.rejectRequest(onboardingNode.consensusNodeId, 1);

        IValidatorManager.Request memory request = validatorManager.getRequest(onboardingNode.consensusNodeId);
        assertEq(uint(request.requestType), uint(IValidatorManager.RequestType.NoRequest));

        IValidatorManager.NodeInfo memory nodeInfo = validatorManager.getNodeInfo(onboardingNode.consensusNodeId);
        assertEq(nodeInfo.manager, address(0));
        assertEq(nodeInfo.consensusNodeId, address(0));
        assertEq(nodeInfo.nodeIds.length, 0);
    }

    function test_rejectRequestOnboarding_revert_requestNotFound() public submitRequest asAdmin {
        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.rejectRequest(vm.randomAddress(), 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.rejectRequest(onboardingNode.consensusNodeId, 2);

        validatorManager.rejectRequest(onboardingNode.consensusNodeId, 1);

        vm.expectRevert(IValidatorManager.RequestNotExists.selector);
        validatorManager.rejectRequest(onboardingNode.consensusNodeId, 1);
    }
}
