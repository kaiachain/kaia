import { ethers } from "hardhat";
import { StakingTrackerMockReceiver__factory, Airdrop__factory, Lockup__factory } from "../../../typechain-types";
import { toPeb } from "../../common/helper";

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
