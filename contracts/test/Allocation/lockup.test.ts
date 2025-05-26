import { loadFixture, setBalance } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";

import { LockupTestFixture } from "../materials";
import { PublicDelegation, PublicDelegation__factory } from "../../typechain-types";
import { FakeContract, smock } from "@defi-wonderland/smock";

const ADMIN_ROLE = hre.ethers.utils.keccak256(hre.ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const SECRETARY_ROLE = hre.ethers.utils.keccak256(hre.ethers.utils.toUtf8Bytes("SECRETARY_ROLE"));

enum AcquisitionStatus {
  UNDEFINED,
  PROPOSED,
  CONFIRMED,
  WITHDRAWN,
  REJECTED,
}

type UnPromisify<T> = T extends Promise<infer U> ? U : T;
describe("Airdrop", function () {
  let fixture: UnPromisify<ReturnType<typeof LockupTestFixture>>;
  beforeEach(async function () {
    fixture = await loadFixture(LockupTestFixture);
    const { lockup, totalDelegatedAmount } = fixture;

    await setBalance(lockup.address, totalDelegatedAmount);
  });
  describe("Setup", function () {
    it("Initial role setup", async function () {
      const { lockup, deployer, admin } = fixture;

      expect(await lockup.ADMIN_ROLE()).to.equal(ADMIN_ROLE);
      expect(await lockup.SECRETARY_ROLE()).to.equal(SECRETARY_ROLE);

      expect(await lockup.hasRole(ADMIN_ROLE, admin.address)).to.be.true;
      expect(await lockup.hasRole(SECRETARY_ROLE, deployer.address)).to.be.true;
    });
    it("Role transfer", async function () {
      const { lockup, deployer, admin, user } = fixture;

      await lockup.connect(admin).transferAdmin(user.address);
      await lockup.connect(deployer).transferSecretary(user.address);

      expect(await lockup.hasRole(ADMIN_ROLE, user.address)).to.be.true;
      expect(await lockup.hasRole(SECRETARY_ROLE, user.address)).to.be.true;
    });
    it("#refreshDelegated: refresh delegated amount", async function () {
      const { lockup, admin, totalDelegatedAmount } = fixture;

      await lockup.connect(admin).refreshDelegated();

      expect(await lockup.totalDelegatedAmount()).to.equal(totalDelegatedAmount);
      expect(await lockup.isInitialized()).to.be.true;
    });
  });
  describe("only role check", function () {
    it("check all functions", async function () {
      const { lockup, user } = fixture;

      await expect(lockup.connect(user).transferAdmin(user.address)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).transferSecretary(user.address)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).refreshDelegated()).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).proposeAcquisition(1)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).requestDelegatedTransfer(1, user.address)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).withdrawAcquisition(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).withdrawDelegatedTransfer(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).withdrawStakingAmounts(user.address, 1)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).claimStakingAmounts(user.address, 0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).confirmAcquisition(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).rejectAcquisition(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).confirmDelegatedTransfer(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
      await expect(lockup.connect(user).rejectDelegatedTransfer(0)).to.be.revertedWithCustomError(
        lockup,
        "AccessControlUnauthorizedAccount"
      );
    });
  });
  describe("Acquisition", function () {
    this.beforeEach(async function () {
      const { lockup, admin } = fixture;

      await lockup.connect(admin).refreshDelegated();
    });
    describe("Propose Acquisition", function () {
      it("#proposeAcquisition: successfully propose acqusition request", async function () {
        const { lockup, admin } = fixture;

        await expect(lockup.connect(admin).proposeAcquisition(1000))
          .to.emit(lockup, "ProposeAcquisition")
          .withArgs(0, 1000);

        const acquisition = await lockup.getAcquisition(0);
        expect(acquisition.amount).to.equal(1000);
        expect(acquisition.status).to.equal(AcquisitionStatus.PROPOSED);
        expect(await lockup.nextAcReqId()).to.equal(1);
      });
      it("#proposeAcquisition: failed to propose acqusition request with 0 amount", async function () {
        const { lockup, admin } = fixture;

        await expect(lockup.connect(admin).proposeAcquisition(0)).to.be.revertedWith("Lockup: invalid amount");
      });
      it("#proposeAcquisition: failed to propose acquisition request more than total delegated amount", async function () {
        const { lockup, admin, totalDelegatedAmount } = fixture;

        await expect(lockup.connect(admin).proposeAcquisition(totalDelegatedAmount + 1n)).to.be.revertedWith(
          "Lockup: invalid amount"
        );

        // or there's pending acquisition
        await lockup.connect(admin).proposeAcquisition(totalDelegatedAmount);

        await expect(lockup.connect(admin).proposeAcquisition(1)).to.be.revertedWith("Lockup: invalid amount");
      });
    });
    describe("Confirm/Reject Acquisition", function () {
      this.beforeEach(async function () {
        const { lockup, admin } = fixture;

        await lockup.connect(admin).proposeAcquisition(1000);
      });
      it("#confirmAcquisition: successfully confirm acquisition request", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmAcquisition(0)).to.emit(lockup, "ConfirmAcquisition").withArgs(0);

        const acquisition = await lockup.getAcquisition(0);
        expect(acquisition.status).to.equal(AcquisitionStatus.CONFIRMED);
      });
      it("#confirmAcquisition: failed to confirm acquisition request not proposed status", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmAcquisition(0)).to.emit(lockup, "ConfirmAcquisition").withArgs(0);
        // Already confirmed
        await expect(lockup.connect(deployer).confirmAcquisition(0)).to.be.revertedWith("Lockup: invalid status");
      });
      it("#rejectAcquisition: successfully reject acquisition request", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).rejectAcquisition(0)).to.emit(lockup, "RejectAcquisition").withArgs(0);

        const acquisition = await lockup.getAcquisition(0);
        expect(acquisition.status).to.equal(AcquisitionStatus.REJECTED);
      });
      it("#rejectAcquisition: failed to reject acquisition request not proposed status", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmAcquisition(0)).to.emit(lockup, "ConfirmAcquisition").withArgs(0);
        // Already confirmed
        await expect(lockup.connect(deployer).rejectAcquisition(0)).to.be.revertedWith("Lockup: invalid status");
      });
    });
    describe("#withdrawAcquisition", function () {
      const requestAmount = hre.ethers.utils.parseEther("1000");
      this.beforeEach(async function () {
        const { lockup, admin } = fixture;

        await lockup.connect(admin).proposeAcquisition(requestAmount);
      });
      it("successfully withdraw acquisition request", async function () {
        const { lockup, deployer, admin } = fixture;

        const beforeBalance = await hre.ethers.provider.getBalance(admin.address);

        await lockup.connect(deployer).confirmAcquisition(0);
        await expect(lockup.connect(admin).withdrawAcquisition(0)).to.emit(lockup, "WithdrawAcquisition").withArgs(0);

        const afterBalance = await hre.ethers.provider.getBalance(admin.address);
        // 0.0001 is a transaction fee
        expect(afterBalance).to.be.closeTo(beforeBalance.add(requestAmount), hre.ethers.utils.parseEther("0.0001"));
        const acquisition = await lockup.getAcquisition(0);
        expect(acquisition.status).to.equal(AcquisitionStatus.WITHDRAWN);
      });
      it("failed to withdraw acquisition request not confirmed status", async function () {
        const { lockup, deployer, admin } = fixture;

        await expect(lockup.connect(admin).withdrawAcquisition(0)).to.be.revertedWith("Lockup: invalid status");

        await lockup.connect(deployer).confirmAcquisition(0);
        await expect(lockup.connect(admin).withdrawAcquisition(0)).to.emit(lockup, "WithdrawAcquisition").withArgs(0);

        await expect(lockup.connect(admin).withdrawAcquisition(0)).to.be.revertedWith("Lockup: invalid status");
      });
    });
  });
  describe("Delegated Transfer", function () {
    this.beforeEach(async function () {
      const { lockup, admin } = fixture;

      await lockup.connect(admin).refreshDelegated();
    });
    describe("Request Delegated Transfer", function () {
      it("#requestDelegatedTransfer: successfully request delegated transfer", async function () {
        const { lockup, admin, user } = fixture;

        await expect(lockup.connect(admin).requestDelegatedTransfer(1000, user.address))
          .to.emit(lockup, "RequestDelegatedTransfer")
          .withArgs(0, 1000, user.address);

        const transfer = await lockup.getDelegatedTransfer(0);
        expect(transfer.amount).to.equal(1000);
        expect(transfer.to).to.equal(user.address);
        expect(transfer.status).to.equal(AcquisitionStatus.PROPOSED);
        expect(await lockup.nextDelegatedTransferId()).to.equal(1);
      });
      it("#requestDelegatedTransfer: failed to request delegated transfer with 0 amount", async function () {
        const { lockup, admin, user } = fixture;

        await expect(lockup.connect(admin).requestDelegatedTransfer(0, user.address)).to.be.revertedWith(
          "Lockup: invalid amount"
        );
      });
      it("#requestDelegatedTransfer: failed to request delegated transfer more than total delegated amount", async function () {
        const { lockup, admin, user, totalDelegatedAmount } = fixture;

        await expect(
          lockup.connect(admin).requestDelegatedTransfer(totalDelegatedAmount + 1n, user.address)
        ).to.be.revertedWith("Lockup: invalid amount");

        // or there's pending transfer
        await lockup.connect(admin).requestDelegatedTransfer(totalDelegatedAmount, user.address);

        await expect(lockup.connect(admin).requestDelegatedTransfer(1, user.address)).to.be.revertedWith(
          "Lockup: invalid amount"
        );
      });
    });
    describe("Confirm/Reject Delegated Transfer", function () {
      this.beforeEach(async function () {
        const { lockup, admin, user } = fixture;

        await lockup.connect(admin).requestDelegatedTransfer(1000, user.address);
      });
      it("#confirmDelegatedTransfer: successfully confirm delegated transfer", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmDelegatedTransfer(0))
          .to.emit(lockup, "ConfirmDelegatedTransfer")
          .withArgs(0);

        const transfer = await lockup.getDelegatedTransfer(0);
        expect(transfer.status).to.equal(AcquisitionStatus.CONFIRMED);
      });
      it("#confirmDelegatedTransfer: failed to confirm delegated transfer not proposed status", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmDelegatedTransfer(0))
          .to.emit(lockup, "ConfirmDelegatedTransfer")
          .withArgs(0);
        // Already confirmed
        await expect(lockup.connect(deployer).confirmDelegatedTransfer(0)).to.be.revertedWith("Lockup: invalid status");
      });
      it("#rejectDelegatedTransfer: successfully reject delegated transfer", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).rejectDelegatedTransfer(0))
          .to.emit(lockup, "RejectDelegatedTransfer")
          .withArgs(0);

        const transfer = await lockup.getDelegatedTransfer(0);
        expect(transfer.status).to.equal(AcquisitionStatus.REJECTED);
      });
      it("#rejectDelegatedTransfer: failed to reject delegated transfer not proposed status", async function () {
        const { lockup, deployer } = fixture;

        await expect(lockup.connect(deployer).confirmDelegatedTransfer(0))
          .to.emit(lockup, "ConfirmDelegatedTransfer")
          .withArgs(0);
        // Already confirmed
        await expect(lockup.connect(deployer).rejectDelegatedTransfer(0)).to.be.revertedWith("Lockup: invalid status");
      });
    });
    describe("#withdrawDelegatedTransfer", function () {
      const requestAmount = hre.ethers.utils.parseEther("1000");
      this.beforeEach(async function () {
        const { lockup, admin, user } = fixture;

        await lockup.connect(admin).requestDelegatedTransfer(requestAmount, user.address);
      });
      it("successfully withdraw delegated transfer", async function () {
        const { lockup, deployer, admin, user } = fixture;

        const beforeBalance = await hre.ethers.provider.getBalance(user.address);

        await lockup.connect(deployer).confirmDelegatedTransfer(0);
        await expect(lockup.connect(admin).withdrawDelegatedTransfer(0))
          .to.emit(lockup, "WithdrawDelegatedTransfer")
          .withArgs(0);

        const afterBalance = await hre.ethers.provider.getBalance(user.address);
        expect(afterBalance).to.be.equal(beforeBalance.add(requestAmount));

        const transfer = await lockup.getDelegatedTransfer(0);
        expect(transfer.status).to.equal(AcquisitionStatus.WITHDRAWN);
      });
      it("failed to withdraw delegated transfer not confirmed status", async function () {
        const { lockup, deployer, admin } = fixture;

        await expect(lockup.connect(admin).withdrawDelegatedTransfer(0)).to.be.revertedWith("Lockup: invalid status");

        await lockup.connect(deployer).confirmDelegatedTransfer(0);
        await expect(lockup.connect(admin).withdrawDelegatedTransfer(0))
          .to.emit(lockup, "WithdrawDelegatedTransfer")
          .withArgs(0);

        await expect(lockup.connect(admin).withdrawDelegatedTransfer(0)).to.be.revertedWith("Lockup: invalid status");
      });
    });
  });
  describe("Interact with public delegation", function () {
    let mockPD: FakeContract<PublicDelegation>;
    this.beforeEach(async function () {
      const { lockup, admin } = fixture;

      await lockup.connect(admin).refreshDelegated();

      mockPD = await smock.fake<PublicDelegation>(PublicDelegation__factory.abi);
      mockPD.maxRedeem.returns(1000);

      mockPD.requestIdToOwner.returns(lockup.address);
    });
    it("#withdrawStakingAmounts: successfully withdraw staking amounts", async function () {
      const { lockup, admin } = fixture;

      await expect(lockup.connect(admin).withdrawStakingAmounts(mockPD.address, 1000))
        .to.emit(lockup, "WithdrawStakingAmounts")
        .withArgs(mockPD.address, 1000);
    });
    it("#withdrawStakingAmounts: failed to withdraw more than max redeemable shares", async function () {
      const { lockup, admin } = fixture;

      await expect(lockup.connect(admin).withdrawStakingAmounts(mockPD.address, 10_000)).to.be.revertedWith(
        "Lockup: invalid shares"
      );
    });
    it("#claimStakingAmounts: successfully claim staking amounts", async function () {
      const { lockup, admin } = fixture;

      await expect(lockup.connect(admin).claimStakingAmounts(mockPD.address, 0))
        .to.emit(lockup, "ClaimStakingAmounts")
        .withArgs(mockPD.address, 0);
    });
    it("#claimStakingAmounts: failed to claim staking of not owned public delegation", async function () {
      const { lockup, admin, user } = fixture;

      mockPD.requestIdToOwner.returns(user.address); // Not owned

      await expect(lockup.connect(admin).claimStakingAmounts(mockPD.address, 0)).to.be.revertedWith(
        "Lockup: invalid request owner"
      );
    });
  });
  describe("Check view functions", function () {
    let acList: number[];
    let dtList: { amount: number; to: string }[];
    this.beforeEach(async function () {
      const { lockup, admin } = fixture;

      await lockup.connect(admin).refreshDelegated();

      acList = [1000, 2000, 3000];
      dtList = [
        { amount: 1000, to: admin.address },
        { amount: 2000, to: admin.address },
        { amount: 3000, to: admin.address },
      ];

      for (const amount of acList) {
        await lockup.connect(admin).proposeAcquisition(amount);
      }
      for (const { amount, to } of dtList) {
        await lockup.connect(admin).requestDelegatedTransfer(amount, to);
      }
    });
    it("#getAllAcquisitions", async function () {
      const { lockup } = fixture;

      const ret = await lockup.getAllAcquisitions();
      expect(ret.length).to.equal(acList.length);
      for (let i = 0; i < acList.length; i++) {
        expect(ret[i].acReqId).to.equal(i);
        expect(ret[i].amount).to.equal(acList[i]);
      }
    });
    it("#getAllDelegatedTransfers", async function () {
      const { lockup } = fixture;

      const ret = await lockup.getAllDelegatedTransfers();
      expect(ret.length).to.equal(dtList.length);
      for (let i = 0; i < dtList.length; i++) {
        expect(ret[i].delegatedTransferId).to.equal(i);
        expect(ret[i].amount).to.equal(dtList[i].amount);
        expect(ret[i].to).to.equal(dtList[i].to);
      }
    });
    it("#getAcquisitionAtStatus: 1 propose, 2 confirm", async function () {
      const { lockup, deployer } = fixture;

      await lockup.connect(deployer).confirmAcquisition(0);
      await lockup.connect(deployer).confirmAcquisition(1);

      const ret = await lockup.getAcquisitionAtStatus(AcquisitionStatus.CONFIRMED);
      expect(ret.length).to.equal(2);
      expect(ret[0].acReqId).to.equal(0);
      expect(ret[1].acReqId).to.equal(1);

      const ret2 = await lockup.getAcquisitionAtStatus(AcquisitionStatus.PROPOSED);
      expect(ret2.length).to.equal(1);
      expect(ret2[0].acReqId).to.equal(2);
    });
    it("#getDelegatedTransferAtStatus: 1 propose, 2 withdrawn", async function () {
      const { lockup, deployer, admin } = fixture;

      await lockup.connect(deployer).confirmDelegatedTransfer(0);
      await lockup.connect(deployer).confirmDelegatedTransfer(1);

      await lockup.connect(admin).withdrawDelegatedTransfer(0);
      await lockup.connect(admin).withdrawDelegatedTransfer(1);

      const ret = await lockup.getDelegatedTransferAtStatus(AcquisitionStatus.WITHDRAWN);
      expect(ret.length).to.equal(2);
      expect(ret[0].delegatedTransferId).to.equal(0);
      expect(ret[1].delegatedTransferId).to.equal(1);

      const ret2 = await lockup.getDelegatedTransferAtStatus(AcquisitionStatus.PROPOSED);
      expect(ret2.length).to.equal(1);
      expect(ret2[0].delegatedTransferId).to.equal(2);
    });
  });
});
