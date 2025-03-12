import { ethers } from "hardhat";
import { StakingTrackerV2, StakingTracker } from "../../../typechain-types";
import { loadFixture, setBalance } from "@nomicfoundation/hardhat-network-helpers";
import { setRegistryOwner, specialContractsFixture } from "./common";
import { CnStakingFixture, deployAndInitStakingContract, nowBlock, nowTime, toPeb } from "../../common/helper";

export async function stakingTrackerV1TestFixture() {
  const { addressBook } = await loadFixture(specialContractsFixture);
  const ret = await prepareStakingTrackerTestFixture(addressBook, "StakingTracker");
  return {
    ...ret,
    stakingTracker: ret.stakingTracker as StakingTracker,
  };
}

export async function stakingTrackerV2TestFixture() {
  const { addressBook, registry } = await loadFixture(specialContractsFixture);
  const ret = await prepareStakingTrackerTestFixture(addressBook, "StakingTrackerV2");
  return {
    ...ret,
    stakingTracker: ret.stakingTracker as StakingTrackerV2,
    registry,
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

  await setRegistryOwner(contractValidator.address);

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
    [rewardAddrA.address, rewardAddrB.address, rewardAddrC.address, rewardAddrD.address]
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
