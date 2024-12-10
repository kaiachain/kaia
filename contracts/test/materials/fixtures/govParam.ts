import { ethers } from "hardhat";
import { GovParam__factory } from "../../../typechain-types";

export async function govParamTestFixture() {
  const accounts = await ethers.getSigners();
  const gp = await new GovParam__factory(accounts[0]).deploy();

  return { gp, deployer: accounts[0], nonOwner: accounts[1], accounts };
}
