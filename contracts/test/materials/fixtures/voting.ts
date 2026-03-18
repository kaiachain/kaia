import { ethers } from "hardhat";
import { loadFixture, setBalance, setCode } from "@nomicfoundation/hardhat-network-helpers";
import {
  CnStakingFixture,
  DAY,
  deployAndInitStakingContract,
  nowTime,
  registerVoter,
  toPeb,
  tx,
  WEEK,
} from "../../common/helper";
import { specialContractsFixture } from "./common";

export async function votingTestFixture() {
  const { addressBook } = await loadFixture(specialContractsFixture);
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

  const RAND_ADDR = "0xe3B0C44298FC1C149AfBF4C8996fb92427aE41E4";
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

  // Set staking tracker address as zero address so that we can use voting contract as owner of staking tracker contract
  const voting = await ethers.deployContract("Voting", [ethers.constants.AddressZero, secretary.address]);
  // Give voting contract KLAY for testing
  await setBalance(voting.address, "0x52B7D2DCC80CD2E4000000");

  // Set stakingTracker2 for testing updateStakingTracker
  const stakingTrackerAddr = await voting.stakingTracker();
  const stakingTracker = await ethers.getContractAt("StakingTracker", stakingTrackerAddr);

  const stakingTrackerByteCode = await ethers.provider.getCode(stakingTracker.address);
  await setCode(RAND_ADDR, stakingTrackerByteCode);
  const stakingTracker2 = await ethers.getContractAt("StakingTracker", RAND_ADDR);

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
