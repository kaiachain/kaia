import { setCode } from "@nomicfoundation/hardhat-network-helpers";
import { ABOOK_ADDRESS, padUtils } from "../../common/helper";
import { ethers, upgrades } from "hardhat";

// TODO: Manage addresses for system contracts in a single file
export const REGISTRY_ADDRESS = "0x0000000000000000000000000000000000000401";
export const KIP103_ADDRESS = "0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95";
export const GOVERNANCE_ADDRESS = "0xcA4Ef926634A530f12e55A0aEE87F195A7B22Aa3";
export const STAKING_TRACKER_ADDRESS = "0x9b8688d616D3D5180d29520c6a0E28582E82BF4d";
export const RAND_ADDR = "0xe3B0C44298FC1C149AfBF4C8996fb92427aE41E4"; // non-null placeholder

export const preDeployedAddresses = [ABOOK_ADDRESS, KIP103_ADDRESS, GOVERNANCE_ADDRESS, STAKING_TRACKER_ADDRESS];
export const preDeployedActivation = ["0x", "0x", "0x", "0x"]; //0, 0, 0, 0
export const preDeployedName = ["AddressBook", "TreasuryRebalance", "Voting", "StakingTracker"];

export async function registryFixture() {
  // Prepare parameters for deploying contracts
  const accounts = await ethers.getSigners();
  const [a1, a2, a3, a4, a5] = accounts;

  // Deploy registry
  const registryFactory = await ethers.getContractFactory("Registry");
  let registry = await registryFactory.deploy();
  await registry.deployed();

  const registryByteCode = await ethers.provider.getCode(registry.address);
  await setCode(REGISTRY_ADDRESS, registryByteCode);
  registry = await ethers.getContractAt("Registry", REGISTRY_ADDRESS);

  // Inject pre-deployed system contracts by setStorageAt:
  // 1. mapping(string => Records[]) public records -> slot 0
  //   1. MappingKey: keccak256(abi.encode(key, slot))
  //   2. Slot: keccak256(abi.encode(MappingKey))
  //   3. ArraySlot: bytes32(uint256(keccak256(abi.encode(Slot))) + i)
  // 2. string[] public names -> slot 1
  //   1. Array length: keccak256(abi.encode(slot))
  //   2. Array element: bytes32(uint256(keccak256(abi.encode(slot))) + i) (Since there's no string more than slot size for pre-deployed contracts)
  // 3. address private _owner -> slot 2
  //   - owner: at slot 2
  // Reference: https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html

  const baseByte = ethers.utils.hexlify(ethers.utils.zeroPad("0x00", 32));
  const paddedSlot = ethers.utils.zeroPad("0x01", 32);
  const elemSlot = ethers.utils.solidityKeccak256(["bytes"], [paddedSlot]);

  // Total 4 pre-deployed system contracts will be registered
  await hre.network.provider.request({
    method: "hardhat_setStorageAt",
    params: [registry.address, "0x" + Number(1).toString(16), padUtils("0x04", 32)],
  });

  for (let i = 0; i < preDeployedAddresses.length; i++) {
    // For mapping(string => Record[]) public records
    const byteKey = ethers.utils.hexlify(ethers.utils.toUtf8Bytes(preDeployedName[i]));
    const mappingKey = byteKey + baseByte.slice(2); // String key doesn't need to be padded
    const slot = ethers.utils.solidityKeccak256(["bytes"], [mappingKey]);
    const arraySlot = ethers.utils.solidityKeccak256(["bytes"], [slot]);

    const arrayElemSlotAddr = ethers.BigNumber.from(arraySlot).toHexString();
    const arrayElemSlotActivation = ethers.BigNumber.from(arraySlot).add(1).toHexString();

    // Length of array
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [registry.address, slot, padUtils("0x01", 32)],
    });

    // Set Record.addr
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [registry.address, arrayElemSlotAddr, padUtils(preDeployedAddresses[i], 32)],
    });

    // Set Record.activation
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [registry.address, arrayElemSlotActivation, padUtils(preDeployedActivation[i], 32)],
    });

    // For string[] public names
    const elemSlotAt = ethers.BigNumber.from(elemSlot).add(i).toHexString();

    const byteArray = ethers.utils.toUtf8Bytes(preDeployedName[i]);
    const hexString = ethers.utils.hexlify(byteArray);
    const length = byteArray.length * 2;

    const lengthBytes = ethers.utils.hexlify(ethers.utils.zeroPad(ethers.utils.hexlify(length), 32));
    const storedString = hexString + lengthBytes.slice(2 + length);

    // Set names[i]
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [registry.address, elemSlotAt, storedString],
    });
  }

  // For address private _owner
  await hre.network.provider.request({
    method: "hardhat_setStorageAt",
    params: [registry.address, "0x" + Number(2).toString(16), padUtils(a1.address, 32)],
  });

  // Deploy mock upgradeable system contract
  const MockUpgradeableSystemContract = await ethers.getContractFactory("MockUpgradeableSystemContract");
  const mockProxy = await upgrades.deployProxy(MockUpgradeableSystemContract, [3], {
    initializer: "initialize",
    kind: "uups",
  });

  return {
    a1,
    a2,
    a3,
    a4,
    a5,
    registry,
    mockProxy,
  };
}
