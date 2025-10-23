import { expect } from "chai";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { auctionTestFixture } from "../materials";
import {
  loadFixture,
  impersonateAccount,
  stopImpersonatingAccount,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { ethers } from "hardhat";
import { fillTypeDataArgs, getMiner, nowBlock, toPeb } from "../common/helper";
import { Signer, Wallet } from "ethers";
import { arrayify, hexlify, parseEther, randomBytes } from "ethers/lib/utils";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

async function generateAuctionTx(
  verifyingContract: string,
  sender: SignerWithAddress,
  auctioneer: Signer,
  to: string,
  data: string,
  blockNumber: number,
  callGasLimit = 10_000_000,
  bid = toPeb(10n),
  nonce = 0
) {
  const targetTxHash = hexlify(randomBytes(32));
  const typeData = fillTypeDataArgs({
    verifyingContract: verifyingContract,
    targetTxHash: targetTxHash,
    blockNumber: blockNumber,
    sender: await sender.getAddress(),
    to: to,
    nonce: nonce,
    bid: bid,
    callGasLimit: callGasLimit,
    data: data,
  });

  const signature = await sender._signTypedData(
    typeData.domain,
    typeData.types,
    typeData.message
  );

  const aucSignature = await auctioneer.signMessage(
    ethers.utils.arrayify(signature)
  );

  const auctionTx = {
    targetTxHash: targetTxHash,
    blockNumber: blockNumber,
    sender: await sender.getAddress(),
    to: to,
    nonce: nonce,
    callGasLimit: callGasLimit,
    bid: bid,
    data: data,
    searcherSig: signature,
    auctioneerSig: aucSignature,
  };

  return {
    auctionTx,
    typeData,
    signature,
    aucSignature,
  };
}

describe("AuctionEntryPoint", () => {
  describe("Check initialize", () => {
    it("Check constants", async () => {
      const { auctionEntryPoint } = await loadFixture(auctionTestFixture);

      expect(await auctionEntryPoint.gasPerByteIntrinsic()).to.equal(16);
      expect(await auctionEntryPoint.gasPerByteEip7623()).to.equal(40);
      expect(await auctionEntryPoint.gasContractExecution()).to.equal(21_000);
      expect(await auctionEntryPoint.gasBufferEstimate()).to.equal(200_000);
      expect(await auctionEntryPoint.gasBufferUnmeasured()).to.equal(35_000);

      expect(await auctionEntryPoint.MAX_DATA_SIZE()).to.equal(64 * 1024);
    });
    it("Check initialize", async () => {
      const { auctionEntryPoint, deployer, auctionDepositVault } =
        await auctionTestFixture();

      expect(await auctionEntryPoint.depositVault()).to.equal(
        auctionDepositVault.address
      );
      expect(await auctionEntryPoint.auctioneer()).to.equal(deployer.address);
    });
  });
  describe("Check configuration", () => {
    it("Only owner", async () => {
      const { auctionEntryPoint, user1 } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionEntryPoint.connect(user1).changeAuctioneer(user1.address)
      ).to.be.revertedWithCustomError(
        auctionEntryPoint,
        "OwnableUnauthorizedAccount"
      );
      await expect(
        auctionEntryPoint.connect(user1).changeDepositVault(user1.address)
      ).to.be.revertedWithCustomError(
        auctionEntryPoint,
        "OwnableUnauthorizedAccount"
      );
      await expect(
        auctionEntryPoint
          .connect(user1)
          .changeGasParameters(16, 40, 21_000, 180_000, 30_000)
      ).to.be.revertedWithCustomError(
        auctionEntryPoint,
        "OwnableUnauthorizedAccount"
      );
    });
    it("Not null address", async () => {
      const { auctionEntryPoint, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionEntryPoint
          .connect(deployer)
          .changeAuctioneer(ethers.constants.AddressZero)
      ).to.be.revertedWithCustomError(auctionEntryPoint, "ZeroAddress");
      await expect(
        auctionEntryPoint
          .connect(deployer)
          .changeDepositVault(ethers.constants.AddressZero)
      ).to.be.revertedWithCustomError(auctionEntryPoint, "ZeroAddress");
    });
    it("Change auctioneer", async () => {
      const { auctionEntryPoint, deployer, user1 } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionEntryPoint.connect(deployer).changeAuctioneer(user1.address)
      )
        .to.emit(auctionEntryPoint, "ChangeAuctioneer")
        .withArgs(deployer.address, user1.address);
      expect(await auctionEntryPoint.auctioneer()).to.equal(user1.address);
    });
    it("Change deposit vault", async () => {
      const { auctionEntryPoint, auctionDepositVault, deployer, user1 } =
        await loadFixture(auctionTestFixture);

      await expect(
        auctionEntryPoint.connect(deployer).changeDepositVault(user1.address)
      )
        .to.emit(auctionEntryPoint, "ChangeDepositVault")
        .withArgs(auctionDepositVault.address, user1.address);
      expect(await auctionEntryPoint.depositVault()).to.equal(user1.address);
    });
    it("Change gas parameters", async () => {
      const { auctionEntryPoint, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionEntryPoint
          .connect(deployer)
          .changeGasParameters(100, 100, 100, 100, 100)
      )
        .to.emit(auctionEntryPoint, "ChangeGasParameters")
        .withArgs(100, 100, 100, 100, 100);

      expect(await auctionEntryPoint.gasPerByteIntrinsic()).to.equal(100);
      expect(await auctionEntryPoint.gasPerByteEip7623()).to.equal(100);
      expect(await auctionEntryPoint.gasContractExecution()).to.equal(100);
      expect(await auctionEntryPoint.gasBufferEstimate()).to.equal(100);
      expect(await auctionEntryPoint.gasBufferUnmeasured()).to.equal(100);
    });
  });
  describe("Check call", () => {
    let fixture: UnPromisify<ReturnType<typeof auctionTestFixture>>;
    let coinbase: Signer;
    const user1Deposit = parseEther("100");
    const user2Deposit = parseEther("100");
    const user3Deposit = parseEther("100");
    beforeEach(async () => {
      fixture = await loadFixture(auctionTestFixture);

      const { user1, user2, user3, auctionDepositVault } = fixture;

      await auctionDepositVault.connect(user1).deposit({ value: user1Deposit });
      await auctionDepositVault.connect(user2).deposit({ value: user2Deposit });
      await auctionDepositVault.connect(user3).deposit({ value: user3Deposit });

      const miner = await ethers.provider
        .getBlock("latest")
        .then((block) => block?.miner);
      if (!miner) {
        throw new Error("Miner not found");
      }
      await impersonateAccount(miner);

      await setBalance(miner, parseEther("100"));
      coinbase = await ethers.getSigner(miner);
    });
    afterEach(async () => {
      await stopImpersonatingAccount(await coinbase.getAddress());
    });
    it("#call: Only proposer", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        await nowBlock()
      );

      await expect(
        auctionEntryPoint.connect(user1).call(auctionTx)
      ).to.be.revertedWithCustomError(auctionEntryPoint, "OnlyProposer");
    });
    it("#call: reverted if invalid signature", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1
      );

      auctionTx.targetTxHash = hexlify(randomBytes(32));

      // It silently returns without reverting
      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;

      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(0);
    });
    it("#call: reverted if invalid block number", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) - 1
      );

      // It silently returns without reverting
      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;

      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(0);
    });
    it("#call: reverted if invalid auctioneer signature", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1
      );

      auctionTx.auctioneerSig =
        "0x" + Buffer.from(randomBytes(65)).toString("hex");

      // It silently returns without reverting
      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;

      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(0);
    });
    it("#call: reverted if lower than minimum bid", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(0n)
      );

      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;

      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(0);
    });
    it("#call: reverted if insufficient deposit balance", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(100n)
      );

      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;
    });
    it("#call: reverted if failed to verify the caller integrity", async () => {
      const { auctionEntryPoint, user1, deployer, testReceiver } = fixture;

      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(10n),
        1 // nonce is 1
      );

      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx)).to.be
        .reverted;

      const { auctionTx: auctionTx2 } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x",
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(10n),
        0
      );

      // Make the signature invalid
      auctionTx2.searcherSig = "0x";

      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx2)).to.be
        .reverted;

      const { auctionTx: auctionTx3 } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        "0x" + "ff".repeat(64 * 1024) + "ff", // 64KB + 1 byte
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(10n),
        0
      );

      await expect(auctionEntryPoint.connect(coinbase).call(auctionTx3)).to.be
        .reverted;
    });
    it("#call: should take bid and gas reimbursement if failed to call the target", async () => {
      const {
        auctionEntryPoint,
        auctionDepositVault,
        auctionFeeVault,
        user1,
        deployer,
        testReceiver,
      } = fixture;

      const calldata = testReceiver.interface.encodeFunctionData("increment");
      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        auctionFeeVault.address, // target isn't test receiver
        calldata,
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(10n),
        0
      );

      const tx = await auctionEntryPoint.connect(coinbase).call(auctionTx);
      const receipt = await tx.wait();
      if (!receipt) {
        throw new Error("Receipt is null");
      }

      await expect(tx)
        .to.emit(auctionEntryPoint, "CallFailed")
        .withArgs(user1.address, 0)
        .to.emit(auctionEntryPoint, "UseNonce")
        .withArgs(user1.address, 1)
        .to.emit(auctionDepositVault, "TakenBid")
        .withArgs(user1.address, toPeb(10n))
        .to.emit(auctionDepositVault, "TakenGas")
        .to.emit(auctionFeeVault, "FeeDeposit")
        .withArgs(await getMiner(), toPeb(10n), 0n, 0n);

      const actualGas = receipt.gasUsed.mul(receipt.effectiveGasPrice);
      const spentGas = user1Deposit
        .sub(BigInt(toPeb(10n)))
        .sub(await auctionDepositVault.depositBalances(user1.address));

      expect(actualGas).to.lessThan(spentGas); // The proposer should get reimbursed for the gas
      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(1);
      expect(await testReceiver.count()).to.equal(0);
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(10n));
    });
    it("#call: should take bid and gas reimbursement if insufficient callGasLimit", async () => {
      const {
        auctionEntryPoint,
        auctionDepositVault,
        auctionFeeVault,
        user1,
        deployer,
        testReceiver,
      } = fixture;

      const calldata = testReceiver.interface.encodeFunctionData("increment");
      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        calldata,
        (await nowBlock()) + 1,
        10_000, // insufficient callGasLimit
        toPeb(10n),
        0
      );

      const tx = await auctionEntryPoint.connect(coinbase).call(auctionTx);
      const receipt = await tx.wait();
      if (!receipt) {
        throw new Error("Receipt is null");
      }

      await expect(tx)
        .to.emit(auctionEntryPoint, "CallFailed")
        .withArgs(user1.address, 0)
        .to.emit(auctionEntryPoint, "UseNonce")
        .withArgs(user1.address, 1)
        .to.emit(auctionDepositVault, "TakenBid")
        .withArgs(user1.address, toPeb(10n))
        .to.emit(auctionDepositVault, "TakenGas")
        .to.emit(auctionFeeVault, "FeeDeposit")
        .withArgs(await getMiner(), toPeb(10n), 0n, 0n);

      const actualGas = receipt.gasUsed.mul(receipt.effectiveGasPrice);
      const spentGas = user1Deposit
        .sub(BigInt(toPeb(10n)))
        .sub(await auctionDepositVault.depositBalances(user1.address));

      expect(actualGas).to.lessThan(spentGas); // The proposer should get reimbursed for the gas
      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(1);
      expect(await testReceiver.count()).to.equal(0);
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(10n));
    });
    it("#call: success", async () => {
      const {
        auctionEntryPoint,
        auctionDepositVault,
        auctionFeeVault,
        user1,
        deployer,
        testReceiver,
      } = fixture;

      const calldata = testReceiver.interface.encodeFunctionData("increment");
      const { auctionTx } = await generateAuctionTx(
        auctionEntryPoint.address,
        user1,
        deployer,
        testReceiver.address,
        calldata,
        (await nowBlock()) + 1,
        10_000_000,
        toPeb(10n),
        0
      );

      const tx = await auctionEntryPoint.connect(coinbase).call(auctionTx);
      const receipt = await tx.wait();
      if (!receipt) {
        throw new Error("Receipt is null");
      }

      // const intrinsicGas = ethers.utils.arrayify(tx.data).reduce((acc, curr) => {
      //   return acc + (curr === 0 ? 4 : 16);
      // }, 0);

      await expect(tx)
        .to.emit(auctionEntryPoint, "Call")
        .withArgs(user1.address, 0)
        .to.emit(auctionEntryPoint, "UseNonce")
        .withArgs(user1.address, 1)
        .to.emit(auctionDepositVault, "TakenBid")
        .withArgs(user1.address, toPeb(10n))
        .to.emit(auctionDepositVault, "TakenGas")
        .to.emit(auctionFeeVault, "FeeDeposit")
        .withArgs(await getMiner(), toPeb(10n), 0n, 0n);

      const actualGas = receipt.gasUsed.mul(receipt.effectiveGasPrice);
      const spentGas = user1Deposit
        .sub(toPeb(10n))
        .sub(await auctionDepositVault.depositBalances(user1.address));

      expect(actualGas).to.lessThan(spentGas); // The proposer should get reimbursed for the gas
      expect(await auctionEntryPoint.nonces(user1.address)).to.equal(1);
      expect(await testReceiver.count()).to.equal(1);
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(10n));
    });
  });
});
