import { ethers } from "hardhat";
import {
  CnStakingV2__factory,
  CnStakingV3__factory,
  PublicDelegation__factory,
  StakingTrackerMockReceiver__factory,
  PublicDelegationFactory__factory,
  MultiCallContract__factory,
} from "../../typechain-types";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { deployAddressBook, nowTime, toPeb } from "./helper";

export const gcId = [1, 2];
export const unlockTime = [100000000000, 200000000000];
export const unlockAmount = [100, 200];

// Roles
export const ROLES = {
  OPERATOR_ROLE: ethers.utils.keccak256(
    ethers.utils.toUtf8Bytes("OPERATOR_ROLE")
  ),
  ADMIN_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE")),
  STAKER_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("STAKER_ROLE")),
  UNSTAKING_APPROVER_ROLE: ethers.utils.keccak256(
    ethers.utils.toUtf8Bytes("UNSTAKING_APPROVER_ROLE")
  ),
  UNSTAKING_CLAIMER_ROLE: ethers.utils.keccak256(
    ethers.utils.toUtf8Bytes("UNSTAKING_CLAIMER_ROLE")
  ),
};

export async function addressBookFixture() {
  return deployAddressBook();
}

// TODO: Tidy up the fixtures
export async function cnV2UnitTestFixture() {
  // Load fixture for address book contract
  const addressBook = await loadFixture(addressBookFixture);

  // Prepare parameters for deploying contracts
  const now = await nowTime();
  const accounts = await ethers.getSigners();
  const [
    contractValidator,
    admin1,
    admin2,
    admin3,
    other1,
    other2,
    nodeId,
    rewardAddr,
  ] = accounts.slice(0, 8);
  const adminList = [admin1, admin2, admin3];
  const unLockTimes = [now + 100, now + 200];
  const unLockAmounts = [200n, 400n].map((x) => toPeb(x));
  const requirement = 2;
  const gcId = 700;

  // Deploy contracts
  const cnStakingV2 = await ethers.deployContract("CnStakingV2", [
    contractValidator.address,
    nodeId.address,
    rewardAddr.address,
    adminList.map((x) => x.address),
    requirement,
    unLockTimes,
    unLockAmounts,
  ]);

  await addressBook.constructContract([contractValidator.address], 1);

  // Register CnStakingV2 contract to AddressBook to use reviseRewardAddress
  await addressBook.registerCnStakingContract(
    nodeId.address,
    cnStakingV2.address,
    rewardAddr.address
  );

  const stakingTrackerMockReceiver = await ethers.deployContract(
    "StakingTrackerMockReceiver"
  );
  const stakingTrackerMockWrong = await ethers.deployContract(
    "StakingTrackerMockWrong"
  );
  const stakingTrackerMockActive = await ethers.deployContract(
    "StakingTrackerMockActive"
  );

  return {
    contractValidator,
    adminList,
    nodeId,
    rewardAddr,
    other1,
    other2,
    unLockTimes,
    unLockAmounts,
    requirement,
    cnStakingV2,
    addressBook,
    stakingTrackerMockReceiver,
    stakingTrackerMockWrong,
    stakingTrackerMockActive,
    gcId,
  };
}

export async function cnV3MultiSigUnitTestFixture() {
  // Load fixture for address book contract
  const addressBook = await loadFixture(addressBookFixture);

  // Prepare parameters for deploying contracts
  const now = await nowTime();
  const accounts = await ethers.getSigners();
  const [
    contractValidator,
    admin1,
    admin2,
    admin3,
    other1,
    other2,
    nodeId,
    rewardAddr,
  ] = accounts.slice(0, 8);
  const adminList = [admin1, admin2, admin3];
  const unLockTimes = [now + 100, now + 200];
  const unLockAmounts = [200n, 400n].map((x) => toPeb(x));
  const requirement = 2;
  const gcId = 700;

  // Deploy contracts
  const cnStakingV3 = await ethers.deployContract("CnStakingV3MultiSig", [
    contractValidator.address,
    nodeId.address,
    rewardAddr.address,
    adminList.map((x) => x.address),
    requirement,
    unLockTimes,
    unLockAmounts,
  ]);

  await addressBook.constructContract([contractValidator.address], 1);

  // Register CnStakingV2 contract to AddressBook to use reviseRewardAddress
  await addressBook.registerCnStakingContract(
    nodeId.address,
    cnStakingV3.address,
    rewardAddr.address
  );

  const stakingTrackerMockReceiver = await ethers.deployContract(
    "StakingTrackerMockReceiver"
  );
  const stakingTrackerMockWrong = await ethers.deployContract(
    "StakingTrackerMockWrong"
  );
  const stakingTrackerMockActive = await ethers.deployContract(
    "StakingTrackerMockActive"
  );

  return {
    contractValidator,
    adminList,
    nodeId,
    rewardAddr,
    other1,
    other2,
    unLockTimes,
    unLockAmounts,
    requirement,
    cnStakingV3,
    addressBook,
    stakingTrackerMockReceiver,
    stakingTrackerMockWrong,
    stakingTrackerMockActive,
    gcId,
  };
}

export async function cnV3MultiSigPublicDelegationTestFixture() {
  const addressBook = await loadFixture(addressBookFixture);

  // Prepare parameters for deploying contracts
  const accounts = await ethers.getSigners();
  const [
    deployer,
    contractValidator,
    admin1,
    admin2,
    other1,
    other2,
    nodeId,
    rewardAddr,
  ] = accounts.slice(0, 8);
  const adminList = [deployer, admin1, admin2];
  const requirement = 2;
  const gcId = 700;

  const pdFactory = await new PublicDelegationFactory__factory(
    deployer
  ).deploy();
  const pdParam = new ethers.utils.AbiCoder().encode(
    ["tuple(address, address,  uint256, string)"],
    [[deployer.address, deployer.address, 0, `GC1`]]
  );
  const stakingTracker = await new StakingTrackerMockReceiver__factory(
    deployer
  ).deploy();

  // Deploy contracts
  const cnStakingV3 = await ethers.deployContract("CnStakingV3MultiSig", [
    contractValidator.address,
    nodeId.address,
    ethers.constants.AddressZero,
    adminList.map((x) => x.address),
    requirement,
    [],
    [],
  ]);

  await cnStakingV3.setStakingTracker(stakingTracker.address);
  await cnStakingV3.setGCId(gcId);

  await cnStakingV3.setPublicDelegation(pdFactory.address, pdParam);
  await cnStakingV3.connect(contractValidator).reviewInitialConditions();
  for (let i = 0; i < adminList.length; i++) {
    await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
  }
  await cnStakingV3.depositLockupStakingAndInit();

  await addressBook.constructContract([contractValidator.address], 1);

  // Register CnStakingV2 contract to AddressBook to use reviseRewardAddress
  await addressBook.registerCnStakingContract(
    nodeId.address,
    cnStakingV3.address,
    rewardAddr.address
  );

  const pd1 = PublicDelegation__factory.connect(
    await cnStakingV3.publicDelegation(),
    deployer
  );

  return {
    deployer,
    contractValidator,
    adminList,
    nodeId,
    rewardAddr,
    other1,
    other2,
    requirement,
    cnStakingV3,
    pd1,
  };
}

export async function cnV2ScenarioTestFixture() {
  // Load fixture for address book contract
  const addressBook = await loadFixture(addressBookFixture);

  // Prepare parameters for deploying contracts
  const now = await nowTime();
  const accounts = await ethers.getSigners();
  const [
    contractValidator,
    admin1,
    admin2,
    admin3,
    admin4,
    other1,
    nodeId,
    rewardAddr,
  ] = accounts.slice(0, 9);
  const adminList = [admin1, admin2, admin3, admin4];
  const unLockTimes = [now + 100, now + 200];
  const unLockAmounts = [200n, 400n].map((x) => toPeb(x));
  const requirement = 3;
  const gcId = 700;

  // Deploy contracts
  const cnStakingV2 = await ethers.deployContract("CnStakingV2", [
    contractValidator.address,
    nodeId.address,
    rewardAddr.address,
    adminList.map((x) => x.address),
    requirement,
    unLockTimes,
    unLockAmounts,
  ]);

  await addressBook.constructContract([contractValidator.address], 1);

  // Register CnStakingV2 contract to AddressBook to use reviseRewardAddress
  await addressBook.registerCnStakingContract(
    nodeId.address,
    cnStakingV2.address,
    rewardAddr.address
  );

  const stakingTrackerMockReceiver = await ethers.deployContract(
    "StakingTrackerMockReceiver"
  );
  const stakingTrackerMockWrong = await ethers.deployContract(
    "StakingTrackerMockWrong"
  );
  const stakingTrackerMockActive = await ethers.deployContract(
    "StakingTrackerMockActive"
  );

  return {
    contractValidator,
    adminList,
    nodeId,
    rewardAddr,
    other1,
    unLockTimes,
    unLockAmounts,
    requirement,
    cnStakingV2,
    addressBook,
    stakingTrackerMockReceiver,
    stakingTrackerMockWrong,
    stakingTrackerMockActive,
    gcId,
  };
}

export async function cnV3InitialLockupNotDepositedFixture() {
  const AB = await loadFixture(addressBookFixture);
  const [
    deployer,
    contractValidator,
    nodeId,
    rewardAddress,
    user,
    voterAddress,
  ] = await ethers.getSigners();

  const stakingTracker = await new StakingTrackerMockReceiver__factory(
    deployer
  ).deploy();
  const cnV3 = await new CnStakingV3__factory(deployer).deploy(
    deployer.address,
    nodeId.address,
    rewardAddress.address,
    unlockTime,
    unlockAmount
  );

  await cnV3.setStakingTracker(stakingTracker.address);
  await cnV3.setGCId(1);
  await cnV3.reviewInitialConditions();

  await AB.constructContract([deployer.address], 1);
  await AB.registerCnStakingContract(
    nodeId.address,
    cnV3.address,
    rewardAddress.address
  );

  return {
    cnV3,
    stakingTracker,
    AB,
    deployer,
    contractValidator,
    nodeId: nodeId.address,
    rewardAddress: rewardAddress.address,
    user,
    voterAddress,
  };
}

export async function cnV3InitialLockupFixture() {
  const fixture = await loadFixture(cnV3InitialLockupNotDepositedFixture);

  await fixture.cnV3.depositLockupStakingAndInit({ value: 300 });

  await fixture.cnV3.updateVoterAddress(fixture.voterAddress.address);

  return fixture;
}

export async function cnV3PublicDelegationNotRegisteredFixture() {
  const AB = await loadFixture(addressBookFixture);
  const [
    deployer,
    contractValidator,
    commission1,
    commission2,
    nodeIdForV2,
    nodeIdForV3InitialLockup,
    rewardAddressForCnV3,
    nodeId,
    nodeId2,
    nodeIdMock,
    user1,
    user2,
    user3,
    voterAddress1,
    voterAddress2,
    voterAddressMock,
  ] = await ethers.getSigners();

  const pdFactory = await new PublicDelegationFactory__factory(
    deployer
  ).deploy();
  const stakingTracker = await new StakingTrackerMockReceiver__factory(
    deployer
  ).deploy();

  const commissionTo = [commission1.address, commission2.address];
  const node = [nodeId.address, nodeId2.address];
  const voterAddress = [voterAddress1.address, voterAddress2.address];
  const cnV3s = [];

  for (let i = 0; i < gcId.length; i++) {
    cnV3s.push(
      await new CnStakingV3__factory(deployer).deploy(
        deployer.address,
        node[i],
        ethers.constants.AddressZero,
        [],
        []
      )
    );
    await cnV3s[i].setStakingTracker(stakingTracker.address);
    await cnV3s[i].setGCId(gcId[i]);
  }

  const testingPsParam = new ethers.utils.AbiCoder().encode(
    ["tuple(address, address,  uint256, string)"],
    [[deployer.address, commissionTo[0], 0, `GC1`]]
  );

  return {
    AB,
    cnV3s,
    stakingTracker,
    deployer,
    contractValidator,
    node,
    nodeIdMock,
    commissionTo,
    user1,
    user2,
    user3,
    voterAddress,
    voterAddressMock,
    nodeIdForV2,
    nodeIdForV3InitialLockup,
    rewardAddressForCnV3,
    pdFactory,
    testingPsParam,
  };
}

export async function cnV3PublicDelegationFixture() {
  const fixture = await loadFixture(cnV3PublicDelegationNotRegisteredFixture);
  const {
    AB,
    cnV3s,
    stakingTracker,
    deployer,
    contractValidator,
    node,
    nodeIdMock,
    commissionTo,
    user1,
    user2,
    user3,
    voterAddress,
    voterAddressMock,
    nodeIdForV2,
    nodeIdForV3InitialLockup,
    rewardAddressForCnV3,
    pdFactory,
  } = fixture;

  // Dummy CnStakingV2
  const cnV2 = await new CnStakingV2__factory(deployer).deploy(
    contractValidator.address,
    nodeIdForV2.address,
    deployer.address,
    [deployer.address],
    1,
    unlockTime,
    unlockAmount
  );

  // Dummy CnV3 with initial lockup
  const cnV3WithInitialLockup = await new CnStakingV3__factory(deployer).deploy(
    deployer.address,
    nodeIdForV3InitialLockup.address,
    rewardAddressForCnV3.address,
    unlockTime,
    unlockAmount
  );

  for (let i = 0; i < gcId.length; i++) {
    const pdParam = new ethers.utils.AbiCoder().encode(
      ["tuple(address, address,  uint256, string)"],
      [[deployer.address, commissionTo[i], 0, `GC${i + 1}`]]
    );

    await cnV3s[i].setPublicDelegation(pdFactory.address, pdParam);

    await cnV3s[i].reviewInitialConditions();
    await cnV3s[i].depositLockupStakingAndInit();
    await cnV3s[i].toggleRedelegation();
    await cnV3s[i].updateVoterAddress(voterAddress[i]);
  }

  const pd1 = PublicDelegation__factory.connect(
    await cnV3s[0].publicDelegation(),
    deployer
  );
  const pd2 = PublicDelegation__factory.connect(
    await cnV3s[1].publicDelegation(),
    deployer
  );

  await AB.constructContract([deployer.address], 1);
  await AB.mockRegisterCnStakingContracts(
    [...node, nodeIdForV3InitialLockup.address, nodeIdForV2.address],
    [
      ...cnV3s.map((c) => c.address),
      cnV3WithInitialLockup.address,
      cnV2.address,
    ],
    [
      ...cnV3s.map(async (c) => await c.rewardAddress()),
      rewardAddressForCnV3.address,
      deployer.address,
    ]
  );

  return {
    AB,
    cnV3s,
    cnV2,
    pd1,
    pd2,
    stakingTracker,
    pdFactory,
    deployer,
    contractValidator,
    node,
    nodeIdMock,
    commissionTo,
    user1,
    user2,
    user3,
    voterAddress,
    voterAddressMock: voterAddressMock.address,
  };
}

export async function multiCallTestFixture() {
  const AB = await loadFixture(addressBookFixture);

  const [deployer] = await ethers.getSigners();

  await AB.constructContract([deployer.address], 1);

  const multiCall = await new MultiCallContract__factory(deployer).deploy();

  return {
    AB,
    multiCall,
    deployer,
  };
}
