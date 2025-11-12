import { expect } from "chai";
import { ethers, upgrades } from "hardhat";

const { parseUnits, parseEther, keccak256, toUtf8Bytes, joinSignature } = ethers.utils;
const { BigNumber, Wallet } = ethers;

type SignerWithAddress = Awaited<ReturnType<typeof ethers.getSigners>>[number];

describe("HolderVerifier", function () {
  let holderVerifier: any;
  let owner: SignerWithAddress;
  let user1: SignerWithAddress;
  let nonOwner: SignerWithAddress;
  let bridge: any;
  let operator: any;
  let guardian: any;
  let judge: any;

  const CONV_RATE = BigNumber.from("148079656000000");
  const fnsaAddr1 = "link1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz";
  const fnsaAddr2 = "link1def456ghi789jkl012mno345pqr678stu901vwx234yzabc123";
  const fnsaAddr3 = "link1ghi789jkl012mno345pqr678stu901vwx234yzabc123def456";
  const valoperAddr1 = "linkvaloper1abc123def456ghi789jkl012mno345pqr678stu901vwx234yz";
  const conyBalance1 = parseUnits("1000", 6);
  const conyBalance2 = parseUnits("2000", 6);
  const conyBalance3 = parseUnits("3000", 6);

  const mockSignature = "0x" + "11".repeat(65);
  const mockMessageHash = keccak256(toUtf8Bytes("message"));

  beforeEach(async function () {
    [owner, user1, nonOwner] = await ethers.getSigners();

    const guardianFactory = await ethers.getContractFactory("Guardian", owner);
    guardian = (await upgrades.deployProxy(guardianFactory, [[owner.address], 1])) as unknown as any;
    await guardian.deployed();
    const guardianAddr = guardian.address;

    const operatorFactory = await ethers.getContractFactory("Operator", owner);
    operator = (await upgrades.deployProxy(operatorFactory, [
      [owner.address],
      guardianAddr,
      1,
    ])) as unknown as any;
    await operator.deployed();
    const operatorAddr = operator.address;

    const judgeFactory = await ethers.getContractFactory("Judge", owner);
    judge = (await upgrades.deployProxy(judgeFactory, [
      [owner.address],
      guardianAddr,
      1,
    ])) as unknown as any;
    await judge.deployed();
    const judgeAddr = judge.address;

    const bridgeFactory = await ethers.getContractFactory("KAIABridge", owner);
    bridge = (await upgrades.deployProxy(bridgeFactory, [
      operatorAddr,
      guardianAddr,
      judgeAddr,
      3,
    ])) as unknown as any;
    await bridge.deployed();

    const bridgeAddr = bridge.address;
    const changeBridgeData = operator.interface.encodeFunctionData("changeBridge", [bridgeAddr]);
    await (await guardian.submitTransaction(operatorAddr, changeBridgeData, 0)).wait();

    const timelockData = bridge.interface.encodeFunctionData("changeTransferTimeLock", [0]);
    await (await guardian.submitTransaction(bridgeAddr, timelockData, 0)).wait();

    const holderVerifierFactory = await ethers.getContractFactory("HolderVerifier", owner);
    holderVerifier = (await holderVerifierFactory.deploy(operatorAddr)) as unknown as any;
    await holderVerifier.deployed();

    const holderVerifierAddr = holderVerifier.address;
    const addOperatorData = operator.interface.encodeFunctionData("addOperator", [holderVerifierAddr]);
    await (await guardian.submitTransaction(operatorAddr, addOperatorData, 0)).wait();

    await owner.sendTransaction({ to: bridgeAddr, value: parseEther("1000") });
  });

  describe("Deployment", function () {
    it("sets the right owner", async function () {
      expect(await holderVerifier.owner()).to.equal(owner.address);
    });
  });

  describe("Record Management", function () {
    describe("addRecord", function () {
      it("adds a record for new address", async function () {
        await expect(holderVerifier.addRecord(fnsaAddr1, conyBalance1))
          .to.emit(holderVerifier, "RecordAdded")
          .withArgs(fnsaAddr1, conyBalance1);

        expect(await holderVerifier.conyBalances(fnsaAddr1)).to.equal(conyBalance1);
        expect(await holderVerifier.provisioned(fnsaAddr1)).to.be.false;
        expect(await holderVerifier.allConyBalances()).to.equal(conyBalance1);
        expect(await holderVerifier.provisionedConyBalances()).to.equal(0);
        expect(await holderVerifier.provisionedAccounts()).to.equal(0);
      });

      it("accepts valoper addresses", async function () {
        await expect(holderVerifier.addRecord(valoperAddr1, conyBalance1))
          .to.emit(holderVerifier, "RecordAdded")
          .withArgs(valoperAddr1, conyBalance1);

        expect(await holderVerifier.conyBalances(valoperAddr1)).to.equal(conyBalance1);
      });

      it("reverts when called by non-owner", async function () {
        await expect(holderVerifier.connect(nonOwner).addRecord(fnsaAddr1, conyBalance1)).to.be.revertedWithCustomError(
          holderVerifier,
          "OwnableUnauthorizedAccount",
        );
      });
    });

    describe("addRecords", function () {
      it("adds multiple records", async function () {
        await expect(
          holderVerifier.addRecords([fnsaAddr1, fnsaAddr2, fnsaAddr3], [conyBalance1, conyBalance2, conyBalance3]),
        )
          .to.emit(holderVerifier, "RecordsAdded")
          .withArgs(3);

        expect(await holderVerifier.allConyBalances()).to.equal(
          conyBalance1.add(conyBalance2).add(conyBalance3),
        );
        expect(await holderVerifier.provisionedConyBalances()).to.equal(0);
        expect(await holderVerifier.provisionedAccounts()).to.equal(0);
      });
    });
  });

  describe("Configuration", function () {
    it("allows owner to change operator", async function () {
      const newOperatorFactory = await ethers.getContractFactory("Operator", owner);
      const newOperator = (await upgrades.deployProxy(newOperatorFactory, [
        [owner.address],
        guardian.address,
        1,
      ])) as unknown as any;
      await newOperator.deployed();
      const currentOperatorAddress = operator.address;
      const newOperatorAddress = newOperator.address;

      await expect(holderVerifier.changeOperator(newOperatorAddress))
        .to.emit(holderVerifier, "OperatorChanged")
        .withArgs(currentOperatorAddress, newOperatorAddress);

      expect(await holderVerifier.operator()).to.equal(newOperatorAddress);
    });
  });

  describe("Access control", function () {
    it("allows owner to transfer ownership", async function () {
      await holderVerifier.transferOwnership(user1.address);
      expect(await holderVerifier.owner()).to.equal(user1.address);
    });
  });

  describe("Error cases", function () {
    beforeEach(async function () {
      await holderVerifier.addRecord(fnsaAddr1, conyBalance1);
    });

    it("reverts with invalid public key length", async function () {
      const invalidPubKey = "0x" + "04" + "11".repeat(32);
      await expect(
        holderVerifier.requestProvision(invalidPubKey, fnsaAddr1, mockMessageHash, mockSignature),
      ).to.be.revertedWith("Invalid public key length");
    });
  });

  describe("Provisioning", function () {
    it("emits ProvisionRequested with correct KAIA amount", async function () {
      const harnessFactory = await ethers.getContractFactory("FnsaVerifyHarness");
      const harness = await harnessFactory.deploy();
      await harness.deployed();

      const wallet = Wallet.createRandom();
      const signingKey = wallet._signingKey();
      const publicKey = signingKey.publicKey;
      const fnsaAddr = await harness.computeFnsaAddr(publicKey);

      const conyBalance = parseUnits("1", 6);
      await holderVerifier.addRecord(fnsaAddr, conyBalance);

      const messageHash = keccak256(toUtf8Bytes("kaia-holder-verifier"));
      const signature = joinSignature(signingKey.signDigest(messageHash));

      const expectedKaiaAmount = conyBalance.mul(CONV_RATE);

      const tx = await holderVerifier.connect(user1).requestProvision(publicKey, fnsaAddr, messageHash, signature);
      await expect(tx)
        .to.emit(holderVerifier, "ProvisionRequested")
        .withArgs(fnsaAddr, user1.address, conyBalance, expectedKaiaAmount);
      await tx.wait();

      const txId = await operator.userIdx2TxID(1n);
      expect(txId).to.equal(1n);

      const submission = await operator.transactions(txId);
      expect(submission.to).to.equal(bridge.address);

      const provisionData = await operator.provisions(txId);
      expect(provisionData.seq).to.equal(1n);
      expect(provisionData.sender).to.equal(fnsaAddr);
      expect(provisionData.receiver).to.equal(user1.address);
      expect(provisionData.amount).to.equal(expectedKaiaAmount);

      expect(await holderVerifier.provisioned(fnsaAddr)).to.be.true;
      expect(await holderVerifier.provisionedConyBalances()).to.equal(conyBalance);
      expect(await holderVerifier.provisionedAccounts()).to.equal(1n);
    });
  });
});
