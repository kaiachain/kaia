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
  CnStakingV4,
  CnStakingV4__factory,
  IERC20,
  IERC20__factory,
} from "../../typechain-types";
import { jumpBlock, nowBlock, toPeb } from "../common/helper";
import { FakeContract, smock } from "@defi-wonderland/smock";
import {
  multiCallTestFixture,
  clRegistryTestFixture,
  registerContract,
} from "../materials";
import { ethers } from "hardhat";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("Multicall Permissioned", function () {
  let multiCallFixture: UnPromisify<ReturnType<typeof multiCallTestFixture>>;
  let fakeWKaia: FakeContract<IERC20>;

  const expectedStakingAmounts = [
    toPeb(3000n),
    toPeb(3000n),
    toPeb(6000n),
    toPeb(9000n),
    toPeb(5000n),
    toPeb(10000n),
    toPeb(15000n),
  ];
  // Test params for CL staking
  const gcId1 = 1;
  const nodeId1 = "0x0000000000000000000000000000000000000001";
  const clPool1 = "0x0000000000000000000000000000000000000002";
  const gcId2 = 2;
  const nodeId2 = "0x0000000000000000000000000000000000000004";
  const clPool2 = "0x0000000000000000000000000000000000000005";

  beforeEach(async function () {
    multiCallFixture = await loadFixture(multiCallTestFixture);

    // Assume that initialization has been done
    const { addressBook, deployer } = multiCallFixture;

    const cn = [];
    const nodeIds = [];
    const rewardAddresses = [];

    // Prepare following CNStaking contracts
    // Note that unstaking amount will be ignored
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
    await addressBook.mockRegisterCnStakingContracts(
      nodeIds,
      cn,
      rewardAddresses
    );
    await addressBook.submitUpdatePocContract(deployer.address, 1);
    await addressBook.submitUpdateKirContract(deployer.address, 1);

    fakeWKaia = await smock.fake<IERC20>(IERC20__factory.abi);
  });
  it("Multicall returns staking info", async function () {
    const { addressBook, multiCall } = multiCallFixture;

    await addressBook.activateAddressBook();

    const stakingInfo = await multiCall.multiCallStakingInfo();

    const stakingAmounts = stakingInfo[2];

    for (let i = 0; i < 7; i++) {
      expect(stakingAmounts[i]).to.equal(expectedStakingAmounts[i]);
    }
  });
  it("Mutlcall returns early if AB not activated", async function () {
    const { multiCall } = multiCallFixture;

    const stakingInfo = await multiCall.multiCallStakingInfo();

    const stakingAmounts = stakingInfo[2];

    expect(stakingAmounts).to.have.lengthOf(0);
  });
  it("Multicall returns DP staking info", async function () {
    const { multiCall, registry } = multiCallFixture;
    const { clRegistry } = await clRegistryTestFixture();

    const curBlock = await nowBlock();
    // Enroll a CLRegistry and WrappedKaia contract address
    await registerContract(registry, "CLRegistry", clRegistry.address);
    await registerContract(registry, "WrappedKaia", fakeWKaia.address);

    fakeWKaia.balanceOf.whenCalledWith(clPool1).returns(toPeb(3000n));
    fakeWKaia.balanceOf.whenCalledWith(clPool2).returns(toPeb(10000n));

    // Add a CL pair1
    await expect(
      clRegistry.addCLPair([{ nodeId: nodeId1, gcId: gcId1, clPool: clPool1 }])
    ).to.emit(clRegistry, "RegisterPair");
    // Add a CL pair2
    await expect(
      clRegistry.addCLPair([{ nodeId: nodeId2, gcId: gcId2, clPool: clPool2 }])
    ).to.emit(clRegistry, "RegisterPair");

    await jumpBlock(curBlock + 100);

    expect(await registry.getActiveAddr("CLRegistry")).to.equal(
      clRegistry.address
    );
    expect(await multiCall.multiCallDPStakingInfo()).to.deep.equal([
      [nodeId1, nodeId2],
      [clPool1, clPool2],
      [toPeb(3000n), toPeb(10000n)],
    ]);
  });
  it("Multicall returns DP staking info (no WKaia)", async function () {
    const { multiCall, registry } = multiCallFixture;
    const { clRegistry } = await clRegistryTestFixture();

    const curBlock = await nowBlock();
    // Enroll a CLRegistry contract address but no WrappedKaia
    await registerContract(registry, "CLRegistry", clRegistry.address);

    fakeWKaia.balanceOf.whenCalledWith(clPool1).returns(toPeb(3000n));
    fakeWKaia.balanceOf.whenCalledWith(clPool2).returns(toPeb(10000n));

    // Add a CL pair1
    await expect(
      clRegistry.addCLPair([{ nodeId: nodeId1, gcId: gcId1, clPool: clPool1 }])
    ).to.emit(clRegistry, "RegisterPair");
    // Add a CL pair2
    await expect(
      clRegistry.addCLPair([{ nodeId: nodeId2, gcId: gcId2, clPool: clPool2 }])
    ).to.emit(clRegistry, "RegisterPair");

    await jumpBlock(curBlock + 100);

    expect(await registry.getActiveAddr("CLRegistry")).to.equal(
      clRegistry.address
    );
    expect(await registry.getActiveAddr("WrappedKaia")).to.equal(
      ethers.constants.AddressZero
    );
    expect(await multiCall.multiCallDPStakingInfo()).to.deep.equal([
      [nodeId1, nodeId2],
      [clPool1, clPool2],
      [toPeb(0n), toPeb(0n)], // No WKaia registered in Registry
    ]);
  });
  it("Multicall returns DP staking info (not activated)", async function () {
    const { multiCall, registry } = multiCallFixture;

    expect(await registry.getActiveAddr("CLRegistry")).to.equal(
      ethers.constants.AddressZero
    );
    expect(await multiCall.multiCallDPStakingInfo()).to.deep.equal([
      [],
      [],
      [],
    ]);
  });
});

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
