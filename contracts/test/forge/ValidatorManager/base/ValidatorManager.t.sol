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

import "forge-std/Test.sol";

import {AddressBookBumped} from "../../../../contracts/testing/system_contracts/AddressBookBumped.sol";
import {SimpleBlsRegistry} from "../../../../contracts/system_contracts/kip113/SimpleBlsRegistry.sol";
import {ERC1967Proxy} from "openzeppelin-contracts-5.0/proxy/ERC1967/ERC1967Proxy.sol";
import {IValidatorManager, ValidatorManager} from "../../../../contracts/system_contracts/consensus/ValidatorManager/ValidatorManager.sol";
import {Registry} from "../../../../contracts/system_contracts/kip149/Registry.sol";
import {CnStakingV3MultiSigFactory} from "../../../../contracts/system_contracts/consensus/CnStakingV3MultiSigFactory/CnStakingV3MultiSigFactory.sol";
import {CnStakingV3MultiSigChunk1} from "../../../../contracts/system_contracts/consensus/CnStakingV3MultiSigFactory/CnStakingV3MultiSigChunk1.sol";
import {CnStakingV3MultiSigChunk2} from "../../../../contracts/system_contracts/consensus/CnStakingV3MultiSigFactory/CnStakingV3MultiSigChunk2.sol";
import {CnStakingV2, ICnStakingV2} from "../../../../contracts/system_contracts/consensus/CnV2/CnStakingV2.sol";
import {CnStakingV3MultiSig} from "../../../../contracts/system_contracts/consensus/CnV3/CnStakingV3MultiSig.sol";
import {PublicDelegationFactory, IPublicDelegation} from "../../../../contracts/system_contracts/consensus/PublicDelegation/PublicDelegationFactory.sol";
import {PublicDelegationFactoryV2} from "../../../../contracts/system_contracts/consensus/PublicDelegation/PublicDelegationFactoryV2.sol";
import {StakingTrackerMockReceiver} from "../../../../contracts/testing/system_contracts/StakingTrackerMock.sol";
import {ICnStakingV3MultiSig} from "../../../../contracts/system_contracts/consensus/CnV3/ICnStakingV3MultiSig.sol";
import {IPublicDelegation} from "../../../../contracts/system_contracts/consensus/PublicDelegation/IPublicDelegation.sol";
import {ICnStakingV3} from "../../../../contracts/system_contracts/consensus/CnV3/ICnStakingV3.sol";
import {IKIP113} from "../../../../contracts/system_contracts/kip113/IKIP113.sol";
import {ISimpleBlsRegistry} from "../../../../contracts/system_contracts/consensus/ISimpleBlsRegistry.sol";
import {IAddressBook} from "../../../../contracts/system_contracts/kip113/IAddressBook.sol";
import {MockPoC} from "../../../../contracts/testing/system_contracts/MockPoC.sol";

contract BaseTest is Test {
    address admin;

    address constant ADDRESS_BOOK_ADDR = 0x0000000000000000000000000000000000000400;
    address constant REGISTRY_ADDR = 0x0000000000000000000000000000000000000401;

    uint256 gcId = 1;

    AddressBookBumped addressBook;
    SimpleBlsRegistry simpleBlsRegistry;
    ValidatorManager validatorManager;
    Registry registry;
    PublicDelegationFactoryV2 publicDelegationFactory;
    StakingTrackerMockReceiver stakingTrackerMock;
    CnStakingV3MultiSigFactory cnStakingV3MultiSigFactory;
    CnStakingV3MultiSigChunk1 cnStakingV3MultiSigChunk1;
    CnStakingV3MultiSigChunk2 cnStakingV3MultiSigChunk2;

    struct TestNode {
        address manager;
        address consensusNodeId;
        address stakingContract;
        address rewardAddress;
        uint256 gcId;
        bytes publicKey;
        bytes pop;
    }

    // Initial registered nodes
    TestNode node1;
    TestNode node2; // With public delegation

    modifier asAdmin() {
        vm.startPrank(admin);
        _;
        vm.stopPrank();
    }

    function setUp() public virtual {
        admin = makeAddr("admin");
        address[] memory adminList = new address[](1);
        TestNode[] memory initialNodes = new TestNode[](2);
        IValidatorManager.NodeInfo[] memory nodeInfos = new IValidatorManager.NodeInfo[](2);

        {
            cnStakingV3MultiSigChunk1 = new CnStakingV3MultiSigChunk1();
            cnStakingV3MultiSigChunk2 = new CnStakingV3MultiSigChunk2();
            cnStakingV3MultiSigFactory = new CnStakingV3MultiSigFactory(
                address(cnStakingV3MultiSigChunk1),
                address(cnStakingV3MultiSigChunk2)
            );
        }

        {
            publicDelegationFactory = new PublicDelegationFactoryV2();
            stakingTrackerMock = new StakingTrackerMockReceiver();
        }

        {
            node1 = _deployTestNode(address(0), false, true, true); // Without public delegation
            node2 = _deployTestNode(address(0), true, true, true);
            initialNodes[0] = node1;
            initialNodes[1] = node2;

            nodeInfos[0] = IValidatorManager.NodeInfo({
                manager: node1.manager,
                consensusNodeId: node1.consensusNodeId,
                nodeIds: new address[](0)
            });
            nodeInfos[1] = IValidatorManager.NodeInfo({
                manager: node2.manager,
                consensusNodeId: node2.consensusNodeId,
                nodeIds: new address[](0)
            });
        }

        {
            adminList[0] = admin;
            _deployAddressBook(adminList, initialNodes);
            _deploySimpleBlsRegistry(admin, initialNodes);
            _deployRegistry(admin);
        }

        {
            validatorManager = _deployValidatorManager(admin, nodeInfos);
            // Add validator manager as an admin of AddressBook
            vm.startPrank(admin);
            addressBook.submitAddAdmin(address(validatorManager));
            simpleBlsRegistry.transferOwnership(address(validatorManager));
            registry.transferOwnership(address(validatorManager));
            vm.stopPrank();
        }
    }

    function _deployValidatorManager(
        address owner,
        IValidatorManager.NodeInfo[] memory nodeInfos
    ) public returns (ValidatorManager) {
        ValidatorManager implementation = new ValidatorManager();
        vm.prank(owner);
        ERC1967Proxy proxy = new ERC1967Proxy(
            address(implementation),
            abi.encodeWithSelector(ValidatorManager.initialize.selector, owner, nodeInfos)
        );
        return ValidatorManager(address(proxy));
    }

    function _deployTestNode(
        address manager,
        bool isPublicDelegation,
        bool initializing,
        bool isStaking
    ) public returns (TestNode memory) {
        uint256 gcIndex = gcId++;
        manager = manager == address(0) ? makeAddr(string(abi.encodePacked("manager", gcIndex))) : manager;
        address nodeId = vm.randomAddress();
        // Fund the node id for minimum balance
        vm.deal(nodeId, 10 ether);

        (address cnStakingV3, address rewardAddress) = _deployStakingContract(
            manager,
            nodeId,
            vm.randomAddress(),
            gcIndex,
            isPublicDelegation,
            initializing,
            isStaking
        );

        (bytes memory publicKey, bytes memory pop) = _generateBlsPublicKeyPop();

        return
            TestNode({
                manager: manager,
                consensusNodeId: nodeId,
                stakingContract: address(cnStakingV3),
                rewardAddress: rewardAddress,
                gcId: gcIndex,
                publicKey: publicKey,
                pop: pop
            });
    }

    function _deployStakingContractV2(
        address manager,
        address nodeId,
        address rewardAddress,
        uint256 gcIndex,
        bool initializing,
        bool isStaking
    ) internal returns (address, address) {
        address contractValidator = makeAddr("contractValidator");
        address[] memory adminList = new address[](1);
        adminList[0] = manager;

        uint256[] memory unlockTime = new uint256[](1);
        uint256[] memory unlockAmount = new uint256[](1);

        unlockTime[0] = block.timestamp + 1 weeks;
        unlockAmount[0] = 1 ether;

        vm.prank(manager);
        ICnStakingV2 cnStakingV2 = ICnStakingV2(
            payable(new CnStakingV2(contractValidator, nodeId, rewardAddress, adminList, 1, unlockTime, unlockAmount))
        );

        if (initializing) {
            vm.startPrank(contractValidator);
            cnStakingV2.setStakingTracker(address(stakingTrackerMock));
            cnStakingV2.setGCId(gcIndex);
            vm.stopPrank();

            vm.prank(contractValidator);
            cnStakingV2.reviewInitialConditions();
            vm.prank(manager);
            cnStakingV2.reviewInitialConditions();

            vm.deal(manager, 1 ether);
            vm.prank(manager);
            cnStakingV2.depositLockupStakingAndInit{value: 1 ether}();

            if (isStaking) {
                cnStakingV2.stakeKlay{value: 5_000_000 ether}();
            }
        }

        return (address(cnStakingV2), rewardAddress);
    }

    function _deployStakingContract(
        address manager,
        address nodeId,
        address rewardAddress,
        uint256 gcIndex,
        bool isPublicDelegation,
        bool initializing,
        bool isStaking
    ) internal returns (address, address) {
        address contractValidator = makeAddr("contractValidator");
        address[] memory adminList = new address[](1);
        adminList[0] = manager;
        rewardAddress = isPublicDelegation ? address(0) : rewardAddress;

        vm.prank(manager);
        ICnStakingV3 cnStakingV3 = ICnStakingV3(
            payable(
                cnStakingV3MultiSigFactory.deployCnStakingV3MultiSig(
                    contractValidator,
                    nodeId,
                    rewardAddress,
                    adminList,
                    1,
                    new uint256[](0),
                    new uint256[](0)
                )
            )
        );

        if (initializing) {
            vm.startPrank(contractValidator);
            cnStakingV3.setStakingTracker(address(stakingTrackerMock));
            cnStakingV3.setGCId(gcIndex);
            if (isPublicDelegation) {
                IPublicDelegation.PDConstructorArgs memory pdArgs = IPublicDelegation.PDConstructorArgs({
                    owner: manager,
                    commissionTo: manager,
                    commissionRate: 0,
                    gcName: "GC"
                });
                cnStakingV3.setPublicDelegation(address(publicDelegationFactory), abi.encode(pdArgs));
                rewardAddress = address(cnStakingV3.publicDelegation());
            }
            vm.stopPrank();

            vm.prank(contractValidator);
            cnStakingV3.reviewInitialConditions();
            vm.prank(manager);
            cnStakingV3.reviewInitialConditions();

            vm.prank(manager);
            cnStakingV3.depositLockupStakingAndInit();

            if (isStaking) {
                _stake(manager, address(cnStakingV3), 5_000_000 ether);
            }
        }

        return (address(cnStakingV3), rewardAddress);
    }

    function _deployAddressBook(address[] memory adminList, TestNode[] memory initialNodes) public {
        addressBook = new AddressBookBumped();
        vm.etch(ADDRESS_BOOK_ADDR, address(addressBook).code);
        addressBook = AddressBookBumped(ADDRESS_BOOK_ADDR);

        address poc = address(new MockPoC());
        address kir = address(new MockPoC());

        addressBook.constructContract(adminList, 1);

        vm.startPrank(adminList[0]);
        for (uint256 i = 0; i < initialNodes.length; i++) {
            addressBook.submitRegisterCnStakingContract(
                initialNodes[i].consensusNodeId,
                address(initialNodes[i].stakingContract),
                initialNodes[i].rewardAddress
            );
        }

        addressBook.submitUpdatePocContract(poc, 1);
        addressBook.submitUpdateKirContract(kir, 1);

        addressBook.submitActivateAddressBook();
        vm.stopPrank();
    }

    function _deploySimpleBlsRegistry(address owner, TestNode[] memory initialNodes) public {
        simpleBlsRegistry = new SimpleBlsRegistry();
        vm.store(address(simpleBlsRegistry), bytes32(uint256(151)), bytes32(uint256(uint160(owner))));

        vm.startPrank(owner);
        for (uint256 i = 0; i < initialNodes.length; i++) {
            simpleBlsRegistry.register(initialNodes[i].consensusNodeId, initialNodes[i].publicKey, initialNodes[i].pop);
        }
        vm.stopPrank();
    }

    function _generateBlsPublicKeyPop() public view returns (bytes memory publicKey, bytes memory pop) {
        publicKey = vm.randomBytes(48);
        pop = vm.randomBytes(96);
    }

    function _deployRegistry(address owner) public {
        registry = new Registry();
        vm.etch(REGISTRY_ADDR, address(registry).code);
        registry = Registry(REGISTRY_ADDR);
        vm.store(REGISTRY_ADDR, bytes32(uint256(2)), bytes32(uint256(uint160(owner))));

        vm.startPrank(owner);
        registry.register("SimpleBlsRegistry", address(simpleBlsRegistry), block.number + 1);
        registry.register("CnStakingV3MultiSigFactory", address(cnStakingV3MultiSigFactory), block.number + 1);
        registry.register("PublicDelegationFactory", address(publicDelegationFactory), block.number + 1);
        vm.stopPrank();

        vm.roll(block.number + 1);
    }

    function _stake(address manager, address stakingContract, uint256 amount) public {
        vm.deal(manager, amount);

        ICnStakingV3 cnStakingV3 = ICnStakingV3(payable(stakingContract));

        uint256 initialStaking = cnStakingV3.staking();
        bool isPDEnabled = cnStakingV3.isPublicDelegationEnabled();
        address pd = cnStakingV3.publicDelegation();

        vm.prank(manager);
        if (isPDEnabled) {
            pd.call{value: amount}("");
        } else {
            cnStakingV3.delegate{value: amount}();
        }

        assertEq(cnStakingV3.staking(), initialStaking + amount);
    }

    function _prepareOnboardingRequest(TestNode memory node) internal returns (bytes memory) {
        return
            abi.encode(
                IValidatorManager.OnboardingRequest({
                    manager: node.manager,
                    consensusNodeId: node.consensusNodeId,
                    stakingContract: address(node.stakingContract),
                    rewardAddress: address(node.rewardAddress),
                    blsPublicKeyInfo: IKIP113.BlsPublicKeyInfo({publicKey: node.publicKey, pop: node.pop})
                })
            );
    }

    function _prepareCustomOnboardingRequest(
        address manager,
        address consensusNodeId,
        address stakingContract,
        address rewardAddress,
        bytes memory publicKey,
        bytes memory pop
    ) internal returns (bytes memory) {
        return
            abi.encode(
                IValidatorManager.OnboardingRequest({
                    manager: manager,
                    consensusNodeId: consensusNodeId,
                    stakingContract: stakingContract,
                    rewardAddress: rewardAddress,
                    blsPublicKeyInfo: IKIP113.BlsPublicKeyInfo({publicKey: publicKey, pop: pop})
                })
            );
    }

    function _verifyManagingNodeIds(TestNode[] memory expectedNodes) internal view {
        IValidatorManager.NodeInfo[] memory nodeInfos = validatorManager.getAllNodeInfos();
        assertEq(nodeInfos.length, expectedNodes.length);
        for (uint256 i = 0; i < expectedNodes.length; i++) {
            assertEq(nodeInfos[i].manager, expectedNodes[i].manager);
            assertEq(nodeInfos[i].consensusNodeId, expectedNodes[i].consensusNodeId);
        }
    }

    function _encodeInvalidNodeInfoError(uint256 errorCode) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IValidatorManager.InvalidNodeInfo.selector, errorCode);
    }

    function _encodeInvalidBLSKeyInfoError(uint256 errorCode) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IValidatorManager.InvalidBLSKeyInfo.selector, errorCode);
    }

    function _encodeInvalidStakingContractError(uint256 errorCode) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IValidatorManager.InvalidStakingContract.selector, errorCode);
    }

    function _encodeIncompatibleNewStakingContractError(uint256 errorCode) internal pure returns (bytes memory) {
        return abi.encodeWithSelector(IValidatorManager.IncompatibleNewStakingContract.selector, errorCode);
    }

    /* ========== ERRORS ========== */
    error OwnableUnauthorizedAccount(address account);

    /* ========== EVENTS ========== */
    event NodeInfoAdded(address indexed consensusNodeId);
    event NodeInfoRemoved(address indexed consensusNodeId);
    event RequestSubmitted(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        IValidatorManager.RequestType requestType,
        bytes requestData
    );
    event RequestCancelled(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        IValidatorManager.RequestType requestType,
        bytes requestData
    );
    event RequestApproved(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        IValidatorManager.RequestType requestType,
        bytes requestData
    );
    event RequestRejected(
        address indexed consensusNodeId,
        address indexed manager,
        uint256 requestIndex,
        IValidatorManager.RequestType requestType,
        bytes requestData
    );
    event NodeManagerTransferred(
        address indexed consensusNodeId,
        address indexed oldManager,
        address indexed newManager
    );
}
