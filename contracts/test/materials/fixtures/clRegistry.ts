import { ethers } from "hardhat";
import { CLRegistry__factory } from "../../../typechain-types";

export async function clRegistryTestFixture() {
  const [deployer] = await ethers.getSigners();
  const clRegistry = await new CLRegistry__factory(deployer).deploy(deployer.address);
  return { clRegistry };
}
