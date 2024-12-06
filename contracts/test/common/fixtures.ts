import { ethers } from "hardhat";
import {
  CnStakingV2__factory,
  CnStakingV3__factory,
  PublicDelegation__factory,
  StakingTrackerMockReceiver__factory,
  PublicDelegationFactory__factory,
  MultiCallContract__factory,
  Airdrop__factory,
  Lockup__factory,
  PocContract__factory,
  KirContract__factory,
  TestReceiver__factory,
  Delegator__factory,
  GovParam__factory,
  CLRegistry__factory,
  StakingTrackerV2,
  StakingTracker,
} from "../../typechain-types";
import { loadFixture, setBalance, setCode } from "@nomicfoundation/hardhat-network-helpers";
import {
  CnStakingFixture,
  DAY,
  deployAddressBook,
  deployAndInitStakingContract,
  nowBlock,
  nowTime,
  registerVoter,
  toPeb,
  tx,
  WEEK,
} from "./helper";

export const gcId = [1, 2];
export const unlockTime = [100000000000, 200000000000];
export const unlockAmount = [100, 200];

// Roles
export const ROLES = {
  OPERATOR_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("OPERATOR_ROLE")),
  ADMIN_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE")),
  STAKER_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("STAKER_ROLE")),
  UNSTAKING_APPROVER_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("UNSTAKING_APPROVER_ROLE")),
  UNSTAKING_CLAIMER_ROLE: ethers.utils.keccak256(ethers.utils.toUtf8Bytes("UNSTAKING_CLAIMER_ROLE")),
};

export async function addressBookFixture() {
  return deployAddressBook();
}

export async function stakingTrackerV1TestFixture() {
  const addressBook = await loadFixture(addressBookFixture);
  const ret = await prepareStakingTrackerTestFixture(addressBook, "StakingTracker");
  return {
    ...ret,
    stakingTracker: ret.stakingTracker as StakingTracker,
  };
}

export async function stakingTrackerV2TestFixture() {
  const addressBook = await loadFixture(addressBookFixture);
  const ret = await prepareStakingTrackerTestFixture(addressBook, "StakingTrackerV2");
  return {
    ...ret,
    stakingTracker: ret.stakingTracker as StakingTrackerV2,
  };
}

async function prepareStakingTrackerTestFixture(addressBook: any, stakingTrackerName: string) {
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
    other2,
    nodeIdA,
    nodeIdB,
    nodeIdC,
    nodeIdD,
    rewardAddrA,
    rewardAddrB,
    rewardAddrC,
    rewardAddrD,
  ] = accounts;

  // Set 100,000,000 KLAY(= 0x52B7D2DCC80CD2E4000000 in hex) for accounts
  await setBalance(admin1.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin2.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin3.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin4.address, "0x52B7D2DCC80CD2E4000000");

  const trackInterval = 100;
  const trackStart = await nowBlock();
  const trackEnd = (await nowBlock()) + trackInterval;
  const adminListA = [admin1];
  const adminListB = [admin2];
  const adminListC = [admin3];
  const adminListD = [admin4];
  const unLockTimes = [now + 100, now + 200];
  const unLockAmounts = [2_000_000n, 4_000_000n].map((x) => toPeb(x));
  const requirement = 1;
  // const gcIdA = 699; // CnStakingV1 doesn't have gcId
  const gcIdB = 700;
  const gcIdC = 701;
  const gcIdD = 702;

  // Deploy contracts
  let stakingTracker;
  switch (stakingTrackerName) {
    case "StakingTracker":
      stakingTracker = await ethers.deployContract("StakingTracker");
      break;
    case "StakingTrackerV2":
      stakingTracker = await ethers.deployContract("StakingTrackerV2", [contractValidator.address]);
      break;
    default:
      throw new Error("Invalid staking tracker name");
  }
  const stakingTrackerAddr = stakingTracker.address;

  // Prepare args for deploying staking contracts
  const argsForContractA: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdA,
    rewardAddr: rewardAddrA,
    adminList: adminListA,
    requirement,
    unLockTimes,
    unLockAmounts,
  };
  const argsForContractB: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdB,
    rewardAddr: rewardAddrB,
    adminList: adminListB,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdB,
  };
  const argsForContractC: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdC,
    rewardAddr: rewardAddrC,
    adminList: adminListC,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdC,
  };
  const argsForContractD: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdD,
    rewardAddr: rewardAddrD,
    adminList: adminListD,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdD,
  };

  // We'll use 1 V1 contract (cnStakingV1A) and 3 V2 contracts (cnStakingV2B, cnStakingV2C, cnStakingV2D) as staking contracts
  // and set admin and requirement to 1 for testing
  const cnStakingV1A = await deployAndInitStakingContract(1, argsForContractA);
  const cnStakingV2B = await deployAndInitStakingContract(2, argsForContractB);
  const cnStakingV2C = await deployAndInitStakingContract(2, argsForContractC);
  const cnStakingV2D = await deployAndInitStakingContract(2, argsForContractD);

  await addressBook.constructContract([contractValidator.address], 1);

  // Register CnStakingV1/V2 contract to AddressBook
  await addressBook.mockRegisterCnStakingContracts(
    [nodeIdA.address, nodeIdB.address, nodeIdC.address, nodeIdD.address],
    [cnStakingV1A.address, cnStakingV2B.address, cnStakingV2C.address, cnStakingV2D.address],
    [rewardAddrA.address, rewardAddrB.address, rewardAddrC.address, rewardAddrD.address],
  );

  return {
    trackInterval,
    trackStart,
    trackEnd,
    contractValidator,
    admin1,
    admin2,
    admin3,
    admin4,
    cnStakingV1A,
    cnStakingV2B,
    cnStakingV2C,
    cnStakingV2D,
    argsForContractA,
    argsForContractB,
    argsForContractC,
    argsForContractD,
    addressBook,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTracker,
    other1,
    other2,
  };
}

export async function votingTestFixture() {
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
    other2,
    voter1,
    voter2,
    voter3,
    voter4,
    nodeIdA,
    nodeIdB,
    nodeIdC,
    nodeIdD,
    rewardAddrA,
    rewardAddrB,
    rewardAddrC,
    rewardAddrD,
    secretary,
  ] = accounts;

  // Set 100,000,000 KLAY(= 0x52B7D2DCC80CD2E4000000 in hex) for accounts
  await setBalance(admin1.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin2.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin3.address, "0x52B7D2DCC80CD2E4000000");
  await setBalance(admin4.address, "0x52B7D2DCC80CD2E4000000");

  const adminListA = [admin1];
  const adminListB = [admin2];
  const adminListC = [admin3];
  const adminListD = [admin4];
  const unLockTimes = [now + 100, now + 200];
  const unLockAmounts = [2_000_000n, 4_000_000n].map((x) => toPeb(x));
  const requirement = 1;
  const gcIdA = 700;
  const gcIdB = 701;
  const gcIdC = 702;
  const gcIdD = 703;

  // Deploy contracts
  const stakingTracker = await ethers.deployContract("StakingTracker");
  const voting = await ethers.deployContract("Voting", [stakingTracker.address, secretary.address]);
  const stakingTrackerAddr = stakingTracker.address;
  await stakingTracker.transferOwnership(voting.address);
  // Give voting contract KLAY for testing
  await setBalance(voting.address, "0x52B7D2DCC80CD2E4000000");

  const stakingTracker2 = await ethers.deployContract("StakingTrackerV2", [voting.address]);

  // Prepare args for deploying staking contracts
  const argsForContractA: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdA,
    rewardAddr: rewardAddrA,
    adminList: adminListA,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdA,
  };
  const argsForContractB: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdB,
    rewardAddr: rewardAddrB,
    adminList: adminListB,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdB,
  };
  const argsForContractC: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdC,
    rewardAddr: rewardAddrC,
    adminList: adminListC,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdC,
  };
  const argsForContractD: CnStakingFixture = {
    contractValidator,
    nodeId: nodeIdD,
    rewardAddr: rewardAddrD,
    adminList: adminListD,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTrackerAddr,
    gcId: gcIdD,
  };

  // We'll use 4 V2 staking contracts
  // and set admin and requirement to 1 for testing
  const cnStakingV2A = await deployAndInitStakingContract(2, argsForContractA);
  const cnStakingV2B = await deployAndInitStakingContract(2, argsForContractB);
  const cnStakingV2C = await deployAndInitStakingContract(2, argsForContractC);
  const cnStakingV2D = await deployAndInitStakingContract(2, argsForContractD);

  // Give staking contract free-stakes for testing
  // V2A: 6,000,000 KLAY
  // V2B: 12,000,000 KLAY
  // V2C: 9,000,000 KLAY
  // V2D: 16,000,000 KLAY
  await cnStakingV2B.connect(admin2).stakeKlay({ value: toPeb(6_000_000n) });
  await cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(3_000_000n) });
  await cnStakingV2D.connect(admin4).stakeKlay({ value: toPeb(10_000_000n) });

  // Prepare transactions for execution
  const txSimpleTransfer100KLAY: tx = {
    target: other1.address,
    value: toPeb(100n),
    calldata: "0x",
  };
  const txSimpleTransfer200KLAY: tx = {
    target: other1.address,
    value: toPeb(200n),
    calldata: "0x",
  };

  await addressBook.constructContract([contractValidator.address], 1);

  await addressBook.mockRegisterCnStakingContracts(
    [nodeIdA.address, nodeIdB.address, nodeIdC.address, nodeIdD.address],
    [cnStakingV2A.address, cnStakingV2B.address, cnStakingV2C.address, cnStakingV2D.address],
    [rewardAddrA.address, rewardAddrB.address, rewardAddrC.address, rewardAddrD.address],
  );

  // Register voters to staking contracts
  await registerVoter(cnStakingV2A, admin1, voter1);
  await registerVoter(cnStakingV2B, admin2, voter2);
  await registerVoter(cnStakingV2C, admin3, voter3);
  await registerVoter(cnStakingV2D, admin4, voter4);

  const ABIForVotingUpdate = [
    "function updateStakingTracker(address newAddr)",
    "function updateSecretary(address newAddr)",
    "function updateAccessRule(bool secretaryPropose, bool voterPropose, bool secretaryExecute, bool voterExecute)",
    "function updateTimingRule(uint256 minVotingDelay,uint256 maxVotingDelay, uint256 minVotingPeriod, uint256 maxVotingPeriod)",
  ];

  const iface = new ethers.utils.Interface(ABIForVotingUpdate);

  // Allocate transactions
  const txUpdateStakingTracker: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateStakingTracker", [stakingTracker2.address]),
  };
  const txUpdateSecretary: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateSecretary", [other1.address]),
  };
  const txUpdateAccessRule: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateAccessRule", [true, true, true, true]),
  };
  const txUpdateTimingRule: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateTimingRule", [DAY, WEEK, DAY, WEEK]),
  };
  const txUpdateAccessRuleWrong: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateAccessRule", [false, false, true, true]),
  };
  const txUpdateTimingRuleWrong: tx = {
    target: voting.address,
    value: "0",
    calldata: iface.encodeFunctionData("updateTimingRule", [DAY - 1, WEEK, DAY, WEEK]),
  };

  return {
    contractValidator,
    admin1,
    admin2,
    admin3,
    admin4,
    cnStakingV2A,
    cnStakingV2B,
    cnStakingV2C,
    cnStakingV2D,
    addressBook,
    requirement,
    unLockTimes,
    unLockAmounts,
    stakingTracker,
    stakingTracker2,
    voting,
    other1,
    other2,
    voter1,
    voter2,
    voter3,
    voter4,
    secretary,
    txSimpleTransfer100KLAY,
    txSimpleTransfer200KLAY,
    txUpdateStakingTracker,
    txUpdateSecretary,
    txUpdateAccessRule,
    txUpdateAccessRuleWrong,
    txUpdateTimingRule,
    txUpdateTimingRuleWrong,
  };
}

// TODO: Tidy up the fixtures
export async function cnV2UnitTestFixture() {
  // Load fixture for address book contract
  const addressBook = await loadFixture(addressBookFixture);

  // Prepare parameters for deploying contracts
  const now = await nowTime();
  const accounts = await ethers.getSigners();
  const [contractValidator, admin1, admin2, admin3, other1, other2, nodeId, rewardAddr] = accounts.slice(0, 8);
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
  await addressBook.registerCnStakingContract(nodeId.address, cnStakingV2.address, rewardAddr.address);

  const stakingTrackerMockReceiver = await ethers.deployContract("StakingTrackerMockReceiver");
  const stakingTrackerMockWrong = await ethers.deployContract("StakingTrackerMockWrong");
  const stakingTrackerMockActive = await ethers.deployContract("StakingTrackerMockActive");

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
  const [contractValidator, admin1, admin2, admin3, other1, other2, nodeId, rewardAddr] = accounts.slice(0, 8);
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
  await addressBook.registerCnStakingContract(nodeId.address, cnStakingV3.address, rewardAddr.address);

  const stakingTrackerMockReceiver = await ethers.deployContract("StakingTrackerMockReceiver");
  const stakingTrackerMockWrong = await ethers.deployContract("StakingTrackerMockWrong");
  const stakingTrackerMockActive = await ethers.deployContract("StakingTrackerMockActive");

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
  const [deployer, contractValidator, admin1, admin2, other1, other2, nodeId, rewardAddr] = accounts.slice(0, 8);
  const adminList = [deployer, admin1, admin2];
  const requirement = 2;
  const gcId = 700;

  const pdFactory = await new PublicDelegationFactory__factory(deployer).deploy();
  const pdParam = new ethers.utils.AbiCoder().encode(
    ["tuple(address, address,  uint256, string)"],
    [[deployer.address, deployer.address, 0, `GC1`]],
  );
  const stakingTracker = await new StakingTrackerMockReceiver__factory(deployer).deploy();

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
  await addressBook.registerCnStakingContract(nodeId.address, cnStakingV3.address, rewardAddr.address);

  const pd1 = PublicDelegation__factory.connect(await cnStakingV3.publicDelegation(), deployer);

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
  const [contractValidator, admin1, admin2, admin3, admin4, other1, nodeId, rewardAddr] = accounts.slice(0, 9);
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
  await addressBook.registerCnStakingContract(nodeId.address, cnStakingV2.address, rewardAddr.address);

  const stakingTrackerMockReceiver = await ethers.deployContract("StakingTrackerMockReceiver");
  const stakingTrackerMockWrong = await ethers.deployContract("StakingTrackerMockWrong");
  const stakingTrackerMockActive = await ethers.deployContract("StakingTrackerMockActive");

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
  const [deployer, contractValidator, nodeId, rewardAddress, user, voterAddress] = await ethers.getSigners();

  const stakingTracker = await new StakingTrackerMockReceiver__factory(deployer).deploy();
  const cnV3 = await new CnStakingV3__factory(deployer).deploy(
    deployer.address,
    nodeId.address,
    rewardAddress.address,
    unlockTime,
    unlockAmount,
  );

  await cnV3.setStakingTracker(stakingTracker.address);
  await cnV3.setGCId(1);
  await cnV3.reviewInitialConditions();

  await AB.constructContract([deployer.address], 1);
  await AB.registerCnStakingContract(nodeId.address, cnV3.address, rewardAddress.address);

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

  const pdFactory = await new PublicDelegationFactory__factory(deployer).deploy();
  const stakingTracker = await new StakingTrackerMockReceiver__factory(deployer).deploy();

  const commissionTo = [commission1.address, commission2.address];
  const node = [nodeId.address, nodeId2.address];
  const voterAddress = [voterAddress1.address, voterAddress2.address];
  const cnV3s = [];

  for (let i = 0; i < gcId.length; i++) {
    cnV3s.push(
      await new CnStakingV3__factory(deployer).deploy(deployer.address, node[i], ethers.constants.AddressZero, [], []),
    );
    await cnV3s[i].setStakingTracker(stakingTracker.address);
    await cnV3s[i].setGCId(gcId[i]);
  }

  const testingPsParam = new ethers.utils.AbiCoder().encode(
    ["tuple(address, address,  uint256, string)"],
    [[deployer.address, commissionTo[0], 0, `GC1`]],
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
    unlockAmount,
  );

  // Dummy CnV3 with initial lockup
  const cnV3WithInitialLockup = await new CnStakingV3__factory(deployer).deploy(
    deployer.address,
    nodeIdForV3InitialLockup.address,
    rewardAddressForCnV3.address,
    unlockTime,
    unlockAmount,
  );

  for (let i = 0; i < gcId.length; i++) {
    const pdParam = new ethers.utils.AbiCoder().encode(
      ["tuple(address, address,  uint256, string)"],
      [[deployer.address, commissionTo[i], 0, `GC${i + 1}`]],
    );

    await cnV3s[i].setPublicDelegation(pdFactory.address, pdParam);

    await cnV3s[i].reviewInitialConditions();
    await cnV3s[i].depositLockupStakingAndInit();
    await cnV3s[i].toggleRedelegation();
    await cnV3s[i].updateVoterAddress(voterAddress[i]);
  }

  const pd1 = PublicDelegation__factory.connect(await cnV3s[0].publicDelegation(), deployer);
  const pd2 = PublicDelegation__factory.connect(await cnV3s[1].publicDelegation(), deployer);

  await AB.constructContract([deployer.address], 1);
  await AB.mockRegisterCnStakingContracts(
    [...node, nodeIdForV3InitialLockup.address, nodeIdForV2.address],
    [...cnV3s.map((c) => c.address), cnV3WithInitialLockup.address, cnV2.address],
    [...cnV3s.map(async (c) => await c.rewardAddress()), rewardAddressForCnV3.address, deployer.address],
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

export async function clRegistryTestFixture() {
  const [deployer] = await ethers.getSigners();
  const clRegistry = await new CLRegistry__factory(deployer).deploy(deployer.address);
  return { clRegistry };
}

export async function airdropTestFixture() {
  const [deployer, notClaimer, ...claimers] = await ethers.getSigners();

  const claimInfo = [];

  let totalAirdropAmount = 0n;

  for (let i = 0; i < claimers.length; i++) {
    claimInfo.push({
      claimer: claimers[i].address,
      amount: toPeb(100n),
    });
    totalAirdropAmount += BigInt(toPeb(100n));
  }

  const airdrop = await new Airdrop__factory(deployer).deploy();

  // Just mock contract without receiver function to test claim failed case.
  const noReceiverContract = await new StakingTrackerMockReceiver__factory(deployer).deploy();

  return {
    airdrop,
    deployer,
    notClaimer,
    claimers,
    claimInfo,
    totalAirdropAmount,
    noReceiverContract,
  };
}

export async function LockupTestFixture() {
  const [deployer, admin, user] = await ethers.getSigners();

  const lockup = await new Lockup__factory(deployer).deploy(admin.address, deployer.address);

  const totalDelegatedAmount = BigInt(toPeb(1_000n));

  return {
    lockup,
    deployer, // secretary
    admin,
    user,
    totalDelegatedAmount,
  };
}

export async function reserveTestFixture() {
  const [deployer, admin, eoa] = await ethers.getSigners();
  const poc = await new PocContract__factory(deployer).deploy([admin.address], 1);
  const kir = await new KirContract__factory(deployer).deploy([admin.address], 1);
  const sca = await new TestReceiver__factory(deployer).deploy();

  return {
    deployer,
    admin,
    sca,
    eoa,
    poc,
    kir,
  };
}

export async function delegatorTestFixture() {
  const fixture = await loadFixture(cnV3PublicDelegationFixture);
  // [cnV3s, pd1, pd2, deployer, user1, user2, user3]

  const delegatorContract = await new Delegator__factory(fixture.deployer).deploy(
    fixture.deployer.address,
    fixture.contractValidator.address,
    fixture.pd1.address,
  );

  return {
    delegator: fixture.deployer,
    delegatee: fixture.contractValidator,
    cn: fixture.cnV3s[0],
    pd: fixture.pd1,
    dc: delegatorContract,
    user: fixture.user1,
  };
}

export async function govParamTestFixture() {
  const accounts = await ethers.getSigners();
  const gp = await new GovParam__factory(accounts[0]).deploy();

  return { gp, deployer: accounts[0], nonOwner: accounts[1], accounts };
}
