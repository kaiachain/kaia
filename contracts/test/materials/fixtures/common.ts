import { ethers } from "hardhat";
import { deployAddressBook, deployRegistry, nowBlock, padUtils } from "../../common/helper";
import { REGISTRY_ADDRESS } from "./registry";

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

export async function specialContractsFixture() {
  const ab = await deployAddressBook();
  const registry = await deployRegistry();
  return {
    addressBook: ab,
    registry,
  };
}

export async function setRegistryOwner(owner: string) {
  await hre.network.provider.request({
    method: "hardhat_setStorageAt",
    params: [REGISTRY_ADDRESS, "0x" + Number(2).toString(16), padUtils(owner, 32)],
  });
}

export async function registerContract(registry: any, name: string, address: string) {
  const now = await nowBlock();
  await registry.register(name, address, now + 2);
}
