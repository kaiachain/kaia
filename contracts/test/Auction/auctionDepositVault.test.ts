import { expect } from "chai";
import {
  loadFixture,
  impersonateAccount,
  stopImpersonatingAccount,
  setBalance,
} from "@nomicfoundation/hardhat-network-helpers";
import { auctionTestFixture } from "../materials";
import { toPeb, nowTime, setTime, getMiner } from "../common/helper";
import { ethers } from "ethers";
import hre from "hardhat";
import {
  arrayify,
  Fragment,
  id,
  keccak256,
  parseEther,
} from "ethers/lib/utils";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

const MIN_WITHDRAW_LOCK_TIME = 60;

describe("AuctionDepositVault", () => {
  describe("Check initialize", () => {
    it("Check initial values", async () => {
      const { auctionDepositVault, auctionFeeVault, deployer } =
        await loadFixture(auctionTestFixture);

      expect(await auctionDepositVault.owner()).to.equal(deployer.address);
      expect(await auctionDepositVault.auctionFeeVault()).to.equal(
        auctionFeeVault.address
      );

      expect(await auctionDepositVault.minDepositAmount()).to.equal(toPeb(10));
      expect(await auctionDepositVault.minWithdrawLockTime()).to.equal(
        MIN_WITHDRAW_LOCK_TIME
      );
    });
  });
  describe("Check configuration", () => {
    it("Only owner", async () => {
      const { auctionDepositVault, user1 } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault.connect(user1).changeAuctionFeeVault(user1.address)
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "OwnableUnauthorizedAccount"
      );
      await expect(
        auctionDepositVault.connect(user1).changeMinDepositAmount(toPeb(100n))
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "OwnableUnauthorizedAccount"
      );
      await expect(
        auctionDepositVault.connect(user1).changeMinWithdrawLocktime(100)
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "OwnableUnauthorizedAccount"
      );
    });
    it("Not null address", async () => {
      const { auctionDepositVault, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault
          .connect(deployer)
          .changeAuctionFeeVault(ethers.constants.AddressZero)
      ).to.be.revertedWithCustomError(auctionDepositVault, "ZeroAddress");
    });
    it("Change auction fee vault", async () => {
      const { auctionDepositVault, auctionFeeVault, deployer } =
        await loadFixture(auctionTestFixture);

      await expect(
        auctionDepositVault
          .connect(deployer)
          .changeAuctionFeeVault(deployer.address)
      )
        .to.emit(auctionDepositVault, "ChangeAuctionFeeVault")
        .withArgs(auctionFeeVault.address, deployer.address);

      expect(await auctionDepositVault.auctionFeeVault()).to.equal(
        deployer.address
      );
    });
    it("Change min deposit amount", async () => {
      const { auctionDepositVault, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault
          .connect(deployer)
          .changeMinDepositAmount(toPeb(100n))
      )
        .to.emit(auctionDepositVault, "ChangeMinDepositAmount")
        .withArgs(toPeb(10), toPeb(100));

      expect(await auctionDepositVault.minDepositAmount()).to.equal(
        toPeb(100n)
      );
    });
    it("Change min withdraw lock time", async () => {
      const { auctionDepositVault, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault.connect(deployer).changeMinWithdrawLocktime(100)
      )
        .to.emit(auctionDepositVault, "ChangeMinWithdrawLocktime")
        .withArgs(MIN_WITHDRAW_LOCK_TIME, 100);

      expect(await auctionDepositVault.minWithdrawLockTime()).to.equal(100);
    });
  });
  describe("Check deposit", () => {
    it("#deposit(depositFor): zero deposit amount", async () => {
      const { auctionDepositVault, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault.connect(deployer).deposit({ value: 0 })
      ).to.be.revertedWithCustomError(auctionDepositVault, "ZeroDepositAmount");

      await expect(
        auctionDepositVault.connect(deployer).depositFor(deployer.address)
      ).to.be.revertedWithCustomError(auctionDepositVault, "ZeroDepositAmount");
    });
    it("#deposit(depositFor): can't deposit under min deposit amount", async () => {
      const { auctionDepositVault, deployer } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault.connect(deployer).deposit({ value: toPeb(9) })
      ).to.be.revertedWithCustomError(auctionDepositVault, "MinDepositNotOver");

      await expect(
        auctionDepositVault
          .connect(deployer)
          .depositFor(deployer.address, { value: toPeb(9) })
      ).to.be.revertedWithCustomError(auctionDepositVault, "MinDepositNotOver");
    });
    it("#deposit(depositFor): success", async () => {
      const { auctionDepositVault, deployer, user1 } = await loadFixture(
        auctionTestFixture
      );

      await expect(
        auctionDepositVault.connect(deployer).deposit({ value: toPeb(10) })
      )
        .to.emit(auctionDepositVault, "VaultDeposit")
        .withArgs(deployer.address, toPeb(10), toPeb(10), 0);

      expect(
        await auctionDepositVault.depositBalances(deployer.address)
      ).to.equal(toPeb(10));
      expect(await auctionDepositVault.getDepositAddrs(0, 1)).to.deep.equal([
        deployer.address,
      ]);

      await expect(
        auctionDepositVault
          .connect(user1)
          .depositFor(deployer.address, { value: toPeb(10) })
      )
        .to.emit(auctionDepositVault, "VaultDeposit")
        .withArgs(deployer.address, toPeb(10), toPeb(20), 0);

      expect(
        await auctionDepositVault.depositBalances(deployer.address)
      ).to.equal(toPeb(20));
      expect(await auctionDepositVault.getDepositAddrs(0, 1)).to.deep.equal([
        deployer.address,
      ]);
    });
  });
  describe("Check reserve withdraw and withdraw", () => {
    let fixture: UnPromisify<ReturnType<typeof auctionTestFixture>>;
    beforeEach(async () => {
      fixture = await loadFixture(auctionTestFixture);

      const { auctionDepositVault, user1, user2, deployer } = fixture;

      // Set min withdraw lock time to 100
      await auctionDepositVault
        .connect(deployer)
        .changeMinWithdrawLocktime(100);

      await auctionDepositVault.connect(user1).deposit({ value: toPeb(20) });
      await auctionDepositVault.connect(user2).deposit({ value: toPeb(30) });
    });
    it("#reserveWithdraw: success", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await expect(auctionDepositVault.connect(user1).reserveWithdraw())
        .to.emit(auctionDepositVault, "VaultReserveWithdraw")
        .withArgs(user1.address, toPeb(20), 0);

      const [at, amount] = await auctionDepositVault.withdrawReservations(
        user1.address
      );

      expect(at).to.equal((await nowTime()) + 100);
      expect(amount).to.equal(toPeb(20));
    });
    it("#reserveWithdraw: can't reserve withdraw twice", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await expect(
        auctionDepositVault.connect(user1).reserveWithdraw()
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "WithdrawReservationExists"
      );
    });
    it("#reserveWithdraw: can't reserve zero amount", async () => {
      const { auctionDepositVault, user3 } = fixture;

      await expect(
        auctionDepositVault.connect(user3).reserveWithdraw()
      ).to.be.revertedWithCustomError(auctionDepositVault, "ZeroDepositAmount");
    });
    it("#reserveWithdraw: can't deposit when withdraw reservation exists", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await expect(
        auctionDepositVault.connect(user1).deposit({ value: toPeb(10) })
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "WithdrawReservationExists"
      );

      await expect(
        auctionDepositVault.depositFor(user1.address, {
          value: toPeb(10),
        })
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "WithdrawReservationExists"
      );
    });
    it("#withdraw: success", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await setTime((await nowTime()) + 100);

      const beforeBalance = await hre.ethers.provider.getBalance(user1.address);

      await expect(auctionDepositVault.connect(user1).withdraw())
        .to.emit(auctionDepositVault, "VaultWithdraw")
        .withArgs(user1.address, toPeb(20), 0);

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        0
      );
      expect(await hre.ethers.provider.getBalance(user1.address)).to.closeTo(
        beforeBalance.add(BigInt(toPeb(20))),
        toPeb(0.0001)
      );

      const [at, amount] = await auctionDepositVault.withdrawReservations(
        user1.address
      );
      expect(at).to.equal(0);
      expect(amount).to.equal(0);
    });
    it("#withdraw: no reservation", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await expect(
        auctionDepositVault.connect(user1).withdraw()
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "WithdrawalNotAllowedYet"
      );
    });
    it("#withdraw: not allowed yet", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await setTime((await nowTime()) + 50);

      await expect(
        auctionDepositVault.connect(user1).withdraw()
      ).to.be.revertedWithCustomError(
        auctionDepositVault,
        "WithdrawalNotAllowedYet"
      );
    });
    it("#withdraw: can deposit after withdrawal", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await setTime((await nowTime()) + 100);

      await auctionDepositVault.connect(user1).withdraw();

      // Deposit again
      await auctionDepositVault.connect(user1).deposit({ value: toPeb(10) });

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        toPeb(10)
      );
    });
  });
  describe("Check take bid and gas reimbursement", () => {
    let fixture: UnPromisify<ReturnType<typeof auctionTestFixture>>;
    let entryPointSigner: ethers.Signer;
    beforeEach(async () => {
      fixture = await loadFixture(auctionTestFixture);

      const { auctionEntryPoint, auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).deposit({ value: toPeb(30) });

      await impersonateAccount(auctionEntryPoint.address);
      entryPointSigner = await hre.ethers.getSigner(auctionEntryPoint.address);
      await setBalance(auctionEntryPoint.address, parseEther("100"));
    });
    afterEach(async () => {
      await stopImpersonatingAccount(fixture.auctionEntryPoint.address);
    });
    it("#takeBid: success", async () => {
      const { auctionDepositVault, auctionFeeVault, user1 } = fixture;

      await expect(
        auctionDepositVault
          .connect(entryPointSigner)
          .takeBid(user1.address, toPeb(30))
      )
        .to.emit(auctionDepositVault, "TakenBid")
        .withArgs(user1.address, toPeb(30))
        .to.emit(auctionFeeVault, "FeeDeposit")
        .withArgs(await getMiner(), toPeb(30), toPeb(0), toPeb(0));

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        toPeb(0)
      );
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(30));

      // No deposit left so user1 will be removed from deposit addresses
      expect(await auctionDepositVault.getDepositAddrs(0, 1)).to.deep.equal([]);
    });
    it("#takeBid: can take bid during the withdrawal", async () => {
      const { auctionDepositVault, auctionFeeVault, user1 } = fixture;

      await auctionDepositVault.connect(user1).reserveWithdraw();

      await auctionDepositVault
        .connect(entryPointSigner)
        .takeBid(user1.address, toPeb(15));

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        toPeb(15)
      );
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(15));

      await setTime((await nowTime()) + 100);

      await expect(auctionDepositVault.connect(user1).withdraw())
        .to.emit(auctionDepositVault, "VaultWithdraw")
        .withArgs(user1.address, toPeb(15), 0);

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        0
      );
    });
    it("#takeBid: only entry point", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await expect(
        auctionDepositVault.connect(user1).takeBid(user1.address, toPeb(30))
      ).to.be.revertedWithCustomError(auctionDepositVault, "OnlyEntryPoint");
    });
    it("#takeBid: didn't take bid when insufficient balance", async () => {
      const { auctionDepositVault, auctionFeeVault, user1 } = fixture;

      await expect(
        auctionDepositVault
          .connect(entryPointSigner)
          .takeBid(user1.address, toPeb(40))
      )
        .to.emit(auctionDepositVault, "InsufficientBalance")
        .withArgs(user1.address, toPeb(30), toPeb(40));

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        toPeb(30)
      );
      expect(await auctionFeeVault.accumulatedBids()).to.equal(toPeb(0));
    });
    it("#takeBid: increase deposit if failed to send bid", async () => {
      const { auctionDepositVault, user1 } = fixture;

      await auctionDepositVault.changeAuctionFeeVault(
        auctionDepositVault.address
      );

      await expect(
        auctionDepositVault
          .connect(entryPointSigner)
          .takeBid(user1.address, toPeb(30))
      )
        .to.emit(auctionDepositVault, "TakenBidFailed")
        .withArgs(user1.address, toPeb(30));

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        toPeb(30)
      );
      expect(await auctionDepositVault.getDepositAddrs(0, 1)).to.deep.equal([
        user1.address,
      ]);
    });
    it("#takeGas: success", async () => {
      const { auctionDepositVault, user1 } = fixture;

      const gasUsed = 1000000n;

      const tx = await auctionDepositVault
        .connect(entryPointSigner)
        .takeGas(user1.address, gasUsed);

      const receipt = await tx.wait();
      if (!receipt) {
        throw new Error("Transaction failed");
      }

      const gasPrice = receipt.effectiveGasPrice;
      const gasAmount = gasPrice.mul(gasUsed);

      expect(receipt?.logs[0].address).to.equal(auctionDepositVault.address);
      expect(receipt?.logs[0].topics[0]).to.equal(
        id(auctionDepositVault.interface.getEvent("TakenGas").format("sighash"))
      );
      expect(receipt?.logs[0].data).to.equal(
        auctionDepositVault.interface.encodeEventLog("TakenGas", [
          user1.address,
          gasAmount,
        ]).data
      );

      expect(await auctionDepositVault.depositBalances(user1.address)).to.equal(
        parseEther("30").sub(gasAmount)
      );
    });
  });
  describe("Check getters", () => {
    let fixture: UnPromisify<ReturnType<typeof auctionTestFixture>>;
    beforeEach(async () => {
      fixture = await loadFixture(auctionTestFixture);

      const { auctionDepositVault, user1, user2, user3 } = fixture;

      await auctionDepositVault.connect(user1).deposit({ value: toPeb(20) });
      await auctionDepositVault.connect(user2).deposit({ value: toPeb(30) });
      await auctionDepositVault.connect(user3).deposit({ value: toPeb(40) });
    });
    it("#getDepositAddrsLength: success", async () => {
      const { auctionDepositVault } = fixture;

      expect(await auctionDepositVault.getDepositAddrsLength()).to.equal(3);
    });
    it("#getDepositAddrs: success", async () => {
      const { auctionDepositVault } = fixture;

      expect(await auctionDepositVault.getDepositAddrs(0, 0)).to.deep.equal([
        fixture.user1.address,
        fixture.user2.address,
        fixture.user3.address,
      ]);
    });
    it("#getDepositAddrs: invalid range", async () => {
      const { auctionDepositVault } = fixture;

      await expect(
        auctionDepositVault.getDepositAddrs(3, 0)
      ).to.be.revertedWithCustomError(auctionDepositVault, "InvalidRange");
    });
    it("#isMinDepositOver: success", async () => {
      const { auctionDepositVault, user1 } = fixture;

      expect(
        await auctionDepositVault.isMinDepositOver(fixture.user1.address)
      ).to.equal(true);
      expect(
        await auctionDepositVault.isMinDepositOver(fixture.user2.address)
      ).to.equal(true);
      expect(
        await auctionDepositVault.isMinDepositOver(fixture.user3.address)
      ).to.equal(true);

      await auctionDepositVault.connect(user1).reserveWithdraw();

      expect(
        await auctionDepositVault.isMinDepositOver(user1.address)
      ).to.equal(true);

      await setTime((await nowTime()) + MIN_WITHDRAW_LOCK_TIME);

      await auctionDepositVault.connect(user1).withdraw();

      expect(
        await auctionDepositVault.isMinDepositOver(user1.address)
      ).to.equal(false);
    });
    it("#getAllAddrsOverMinDeposit: success", async () => {
      const { auctionDepositVault, user1, user2, user3 } = fixture;

      expect(
        await auctionDepositVault.getAllAddrsOverMinDeposit(0, 0)
      ).to.deep.equal([
        [user1.address, user2.address, user3.address],
        [toPeb(20), toPeb(30), toPeb(40)],
        [0, 0, 0],
      ]);

      await auctionDepositVault.connect(user1).reserveWithdraw();
      await setTime((await nowTime()) + MIN_WITHDRAW_LOCK_TIME);
      await auctionDepositVault.connect(user1).withdraw();

      expect(
        await auctionDepositVault.getAllAddrsOverMinDeposit(0, 0)
      ).to.deep.equal([
        [user3.address, user2.address],
        [toPeb(40), toPeb(30)],
        [0, 0],
      ]);
    });
    it("#getAllAddrsOverMinDeposit: invalid range", async () => {
      const { auctionDepositVault } = fixture;

      await expect(
        auctionDepositVault.getAllAddrsOverMinDeposit(3, 0)
      ).to.be.revertedWithCustomError(auctionDepositVault, "InvalidRange");
    });
  });
});
