import { ethers } from "hardhat";
import { MultiCallContract__factory } from "../../../typechain-types";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { setRegistryOwner, specialContractsFixture } from "./common";

export async function multiCallTestFixture() {
  const { addressBook, registry } = await loadFixture(specialContractsFixture);

  const [deployer] = await ethers.getSigners();

  await addressBook.constructContract([deployer.address], 1);

  await setRegistryOwner(deployer.address);
  const multiCall = await new MultiCallContract__factory(deployer).deploy();

  return {
    addressBook,
    registry,
    multiCall,
    deployer,
  };
}
