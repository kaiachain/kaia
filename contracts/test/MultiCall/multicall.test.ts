import {
  loadFixture,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

import {
  CnStakingContract,
  CnStakingContract__factory,
  CnStakingV2,
  CnStakingV2__factory,
  CnStakingV3MultiSig,
  CnStakingV3MultiSig__factory,
} from "../../typechain-types";

import { toPeb } from "../common/helper";
import { smock } from "@defi-wonderland/smock";
import { multiCallTestFixture } from "../common/fixtures";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("Multicall", function () {
  let fixture: UnPromisify<ReturnType<typeof multiCallTestFixture>>;
  const expectedStakingAmounts = [
    toPeb(3000n),
    toPeb(3000n),
    toPeb(5500n),
    toPeb(8000n),
    toPeb(4500n),
    toPeb(9000n),
    toPeb(13500n),
  ];
  beforeEach(async function () {
    fixture = await loadFixture(multiCallTestFixture);

    // Assume that initialization has been done
    const { AB, deployer } = fixture;

    const cn = [];
    const nodeIds = [];
    const rewardAddresses = [];

    // Prepare following CNStaking contracts
    // V1: 1 - 3000 KAIA
    // V2: 3 - 3000 KAIA, 6000 KAIA, 9000 KAIA
    //            0 KAIA,  500 KAIA, 1000 KAIA
    // V3: 3 - 5000 KAIA, 10000 KAIA, 15000 KAIA
    //          500 KAIA,  1000 KAIA,  1500 KAIA

    const setupFunction = (contract: any, version: number, address: string) => {
      contract.VERSION.returns(version);
      contract.nodeId.returns(address);
      contract.rewardAddress.returns(address);
    };

    const cnV1 = await smock.fake<CnStakingContract>(
      CnStakingContract__factory.abi
    );
    setupFunction(cnV1, 1, cnV1.address);
    await setBalance(
      cnV1.address,
      hre.ethers.utils.parseEther(3000n.toString())
    );
    cnV1.staking.returns(toPeb(3000n));
    cn.push(cnV1.address);
    nodeIds.push(cnV1.address);
    rewardAddresses.push(cnV1.address);

    for (let i = 0; i < 3; i++) {
      const cnV2 = await smock.fake<CnStakingV2>(CnStakingV2__factory.abi);
      setupFunction(cnV2, 2, cnV2.address);
      await setBalance(
        cnV2.address,
        hre.ethers.utils.parseEther((3000n * (BigInt(i) + 1n)).toString())
      );
      cnV2.staking.returns(toPeb(3000n * (BigInt(i) + 1n)));
      cnV2.unstaking.returns(toPeb(500n * BigInt(i)));
      cn.push(cnV2.address);
      nodeIds.push(cnV2.address);
      rewardAddresses.push(cnV2.address);
    }

    for (let i = 0; i < 3; i++) {
      const cnV3 = await smock.fake<CnStakingV3MultiSig>(
        CnStakingV3MultiSig__factory.abi
      );
      setupFunction(cnV3, 3, cnV3.address);
      await setBalance(
        cnV3.address,
        hre.ethers.utils.parseEther((5000n * (BigInt(i) + 1n)).toString())
      );
      cnV3.staking.returns(toPeb(5000n * (BigInt(i) + 1n)));
      cnV3.unstaking.returns(toPeb(500n * (BigInt(i) + 1n)));
      cn.push(cnV3.address);
      nodeIds.push(cnV3.address);
      rewardAddresses.push(cnV3.address);
    }

    await AB.mockRegisterCnStakingContracts(nodeIds, cn, rewardAddresses);
    await AB.submitUpdatePocContract(deployer.address, 1);
    await AB.submitUpdateKirContract(deployer.address, 1);
  });
  it("Multicall returns staking info", async function () {
    const { AB, multiCall } = fixture;

    await AB.activateAddressBook();

    const stakingInfo = await multiCall.multiCallStakingInfo();

    const stakingAmounts = stakingInfo[2];

    for (let i = 0; i < 7; i++) {
      expect(stakingAmounts[i]).to.equal(expectedStakingAmounts[i]);
    }
  });
  it("Mutlcall returns early if AB not activated", async function () {
    const { multiCall } = fixture;

    const stakingInfo = await multiCall.multiCallStakingInfo();

    const stakingAmounts = stakingInfo[2];

    expect(stakingAmounts).to.have.lengthOf(0);
  });
});
