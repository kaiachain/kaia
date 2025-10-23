import { ethers } from "hardhat";
import {
  AuctionEntryPoint__factory,
  AuctionDepositVault__factory,
  AuctionFeeVault__factory,
  TestReceiver__factory,
  MockCnStakingOverV2__factory,
} from "../../../typechain-types";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import {
  registerContract,
  setRegistryOwner,
  specialContractsFixture,
} from "./common";

export const auctionTestFixture = async () => {
  const { addressBook, registry } = await loadFixture(specialContractsFixture);

  const [deployer, user1, user2, user3] = await ethers.getSigners();

  const auctionFeeVault = await new AuctionFeeVault__factory(deployer).deploy(
    deployer.address,
    0,
    1000
  );
  const auctionDepositVault = await new AuctionDepositVault__factory(
    deployer
  ).deploy(deployer.address, auctionFeeVault.address);
  const auctionEntryPoint = await new AuctionEntryPoint__factory(
    deployer
  ).deploy(deployer.address, auctionDepositVault.address, deployer.address);
  const testReceiver = await new TestReceiver__factory(deployer).deploy();
  const cnStaking = await new MockCnStakingOverV2__factory(deployer).deploy();
  await cnStaking.mockSetAdmin(deployer.address);

  await setRegistryOwner(deployer.address);
  await registerContract(
    registry,
    "AuctionEntryPoint",
    auctionEntryPoint.address
  );

  await addressBook.constructContract([deployer.address], 1);

  return {
    auctionEntryPoint,
    auctionDepositVault,
    auctionFeeVault,
    deployer,
    user1,
    user2,
    user3,
    testReceiver,
    cnStaking,
    addressBook,
  };
};
