import {
  loadFixture,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import {
  CnStakingV4,
  CnStakingV4__factory,
} from "../../typechain-types";
import { toPeb } from "../common/helper";
import { FakeContract, smock } from "@defi-wonderland/smock";
import { ethers } from "hardhat";

// NOTE: Permissioned multicall tests are in contracts/test/MultiCall/multicall.test.ts

describe("Multicall permissionless", function () {
  const ABOOK_ADDRESS = "0x0000000000000000000000000000000000000400";

  async function deployMulticallFixture() {
    const [deployer] = await ethers.getSigners();

    // Deploy ABv2Mock at 0x400
    const abv2Factory = await ethers.getContractFactory("AddressBookV2Mock");
    const abv2Tmp = await abv2Factory.deploy();
    const abv2Code = await ethers.provider.getCode(abv2Tmp.address);
    await hre.network.provider.request({
      method: "hardhat_setCode",
      params: [ABOOK_ADDRESS, abv2Code],
    });
    const abv2 = await ethers.getContractAt("AddressBookV2Mock", ABOOK_ADDRESS);

    // Create CnStaking fakes
    const stakingData = [
      { staking: 5000n, unstaking: 500n },
      { staking: 10000n, unstaking: 1000n },
      { staking: 15000n, unstaking: 1500n },
    ];
    const expectedAmounts = stakingData.map((d) => toPeb(d.staking - d.unstaking));

    const cnStakingAddrs: string[] = [];
    for (let i = 0; i < stakingData.length; i++) {
      const fake = await smock.fake<CnStakingV4>(CnStakingV4__factory.abi);
      fake.staking.returns(toPeb(stakingData[i].staking));
      fake.unstaking.returns(toPeb(stakingData[i].unstaking));

      await abv2.addProfile(
        deployer.address,
        fake.address,
        deployer.address,
        0,
        6 // State.ValActive
      );
      cnStakingAddrs.push(fake.address);
    }

    // Deploy MultiCallContract
    const multiCallFactory = await ethers.getContractFactory("MultiCallContract");
    const multiCall = await multiCallFactory.deploy();

    return { deployer, abv2, multiCall, cnStakingAddrs, expectedAmounts };
  }

  it("multiCallStakingInfoPermissionless returns profiles, amounts, kef, kif", async function () {
    const { deployer, abv2, multiCall, cnStakingAddrs, expectedAmounts } = await loadFixture(deployMulticallFixture);

    // Set fund addresses
    await abv2.setFundAddresses(
      "0x000000000000000000000000000000000000aaa1",
      "0x000000000000000000000000000000000000aaa2",
      "0x000000000000000000000000000000000000aaa3"
    );

    const result = await multiCall.multiCallStakingInfoPermissionless();
    const [profiles, stakingAmounts, retKef, retKif, retKpf] = result;

    // Verify profiles length
    expect(profiles.length).to.equal(3);

    // Verify profiles and staking amounts
    for (let i = 0; i < 3; i++) {
      expect(profiles[i].nodeId).to.equal(deployer.address);
      expect(profiles[i].stakingContract).to.equal(cnStakingAddrs[i]);
      expect(profiles[i].rewardAddress).to.equal(deployer.address);
      expect(profiles[i].timeoutAt).to.equal(0);
      expect(profiles[i].state).to.equal(6);
      expect(stakingAmounts[i]).to.equal(expectedAmounts[i]);
    }

    expect(retKef.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa1");
    expect(retKif.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa2");
    expect(retKpf.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa3");
  });

  it("multiCallNodeStatesPermissionless returns profiles, amounts, timeouts, maxCounts, thresholds, epochVACount, slotLimits", async function () {
    const { abv2, multiCall, cnStakingAddrs, expectedAmounts } = await loadFixture(deployMulticallFixture);

    // Set timeouts, max counts, thresholds, and epoch VA count
    await abv2.setTimeouts(28800, 2592000); // 8h, 30d in seconds
    await abv2.setMaxCounts(50, 100);
    await abv2.setThresholds(2, 300); // pfsThreshold=2, cfsThreshold=300
    await abv2.setEpochVACount(10); // epochVACount=10 → minActive=ceil(20/3)=7, maxSlot=ceil(3/2)=2

    const result = await multiCall.multiCallNodeStatesPermissionless();
    const [
      profiles, stakingAmounts,
      retPauseTimeout, retIdleTimeout,
      retMaxValCount, retMaxReadyCandCount,
      retPfsThreshold, retCfsThreshold,
      retEpochVACount,
      retMaxSlotAvailable, retMinActiveCount,
    ] = result;

    expect(profiles.length).to.equal(3);
    for (let i = 0; i < 3; i++) {
      expect(profiles[i].stakingContract).to.equal(cnStakingAddrs[i]);
      expect(stakingAmounts[i]).to.equal(expectedAmounts[i]);
    }

    expect(retPauseTimeout).to.equal(28800);
    expect(retIdleTimeout).to.equal(2592000);
    expect(retMaxValCount).to.equal(50);
    expect(retMaxReadyCandCount).to.equal(100);
    expect(retPfsThreshold).to.equal(2);
    expect(retCfsThreshold).to.equal(300);
    expect(retEpochVACount).to.equal(10);
    // SlotMath for n=10: minActive=ceil(20/3)=7, maxSlot=ceil(floor(10/3)/2)=ceil(3/2)=2
    expect(retMaxSlotAvailable).to.equal(2);
    expect(retMinActiveCount).to.equal(7);
  });
});
