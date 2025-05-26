import { ethers } from "hardhat";
import { smock } from "@defi-wonderland/smock";
import { SenderTest1, SenderTest1__factory } from "../../../typechain-types";

export const executionBlock = 200;
export const value = hre.ethers.utils.parseEther("20");
export const memo =
  '{ "retirees": [ { "zeroed": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "balance": 20000000000000000000 }, { "zeroed": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "balance": 20000000000000000000 } ], "newbies": [ { "newbie": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "fundAllocated": 10000000000000000000 }, { "newbie": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "fundAllocated": 10000000000000000000 } ], "burnt": 7.2e+37, "success": true }';

export async function rebalanceTestFixture() {
  // Prepare parameters for deploying contracts
  const [rebalance_manager, eoaZeroed, allocated1, allocated2, ...zeroedAdmins] = (await ethers.getSigners()).slice(
    0,
    8
  );

  // Deploy and fund KIF and KEF which will be zeroed after rebalancing.
  const KIF = await ethers.getContractFactory("SenderTest1");
  const zeroed1 = await KIF.deploy();
  await rebalance_manager.sendTransaction({ to: zeroedAdmins[0].address, value: value });
  await rebalance_manager.sendTransaction({ to: zeroedAdmins[1].address, value: value });
  await rebalance_manager.sendTransaction({ to: zeroed1.address, value: value });

  const KEF = await ethers.getContractFactory("SenderTest1");
  const zeroed2 = await KEF.deploy();
  await rebalance_manager.sendTransaction({ to: zeroedAdmins[2].address, value: value });
  await rebalance_manager.sendTransaction({ to: zeroedAdmins[3].address, value: value });
  await rebalance_manager.sendTransaction({ to: zeroed2.address, value: value });

  const mockZeroed3 = await smock.fake<SenderTest1>(SenderTest1__factory.abi, {
    address: "0x38138d89c321b3b5f421e9452b69cf29e4380bae",
  });

  const REBALANCE = await ethers.getContractFactory("TreasuryRebalanceV2");
  const trV2 = await REBALANCE.deploy(executionBlock);

  return {
    rebalance_manager,
    eoaZeroed,
    zeroedAdmins,
    zeroed1,
    zeroed2,
    mockZeroed3,
    allocated1,
    allocated2,
    trV2,
  };
}
