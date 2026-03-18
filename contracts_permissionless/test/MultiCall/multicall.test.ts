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

  it("multiCallStakingInfoPermissionless returns profiles, amounts, kef, kif", async function () {
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

    // Set fund addresses
    await abv2.setFundAddresses(
      "0x000000000000000000000000000000000000aaa1",
      "0x000000000000000000000000000000000000aaa2",
      "0x000000000000000000000000000000000000aaa3"
    );

    // Create CnStaking fakes with staking/unstaking values
    // Node 0: staking 5000, unstaking 500 → effective 4500
    // Node 1: staking 10000, unstaking 1000 → effective 9000
    // Node 2: staking 15000, unstaking 1500 → effective 13500
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

      // Add profile: nodeId = deployer-derived address, stakingContract = fake, rewardAddress = deployer
      await abv2.addProfile(
        deployer.address, // nodeId (doesn't matter for this test)
        fake.address,     // stakingContract
        deployer.address, // rewardAddress
        0,                // timeoutAt
        6                 // State.ValActive
      );
      cnStakingAddrs.push(fake.address);
    }

    // Deploy MultiCallContract
    const multiCallFactory = await ethers.getContractFactory("MultiCallContract");
    const multiCall = await multiCallFactory.deploy();

    // Call multiCallStakingInfoPermissionless
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
      expect(profiles[i].state).to.equal(6); // ValActive
      expect(stakingAmounts[i]).to.equal(expectedAmounts[i]);
    }

    // Verify fund addresses
    expect(retKef.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa1");
    expect(retKif.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa2");
    expect(retKpf.toLowerCase()).to.equal("0x000000000000000000000000000000000000aaa3");
  });
});
