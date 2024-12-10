import { ethers } from "hardhat";
import { CnStakingV2__factory, CnStakingV3__factory, PublicDelegation__factory } from "../../../typechain-types";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { gcId, unlockTime, unlockAmount } from "./common";
import { cnV3PublicDelegationNotRegisteredFixture } from "./cnV3";

export async function cnV3PublicDelegationFixture() {
  const fixture = await loadFixture(cnV3PublicDelegationNotRegisteredFixture);
  const {
    addressBook,
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

  const pd1 = PublicDelegation__factory.connect(await cnV3s[0].publicDelegation(), deployer);
  const pd2 = PublicDelegation__factory.connect(await cnV3s[1].publicDelegation(), deployer);

  await addressBook.constructContract([deployer.address], 1);
  await addressBook.mockRegisterCnStakingContracts(
    [...node, nodeIdForV3InitialLockup.address, nodeIdForV2.address],
    [...cnV3s.map((c) => c.address), cnV3WithInitialLockup.address, cnV2.address],
    [...cnV3s.map(async (c) => await c.rewardAddress()), rewardAddressForCnV3.address, deployer.address]
  );

  return {
    addressBook,
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
