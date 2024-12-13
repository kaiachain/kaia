import { loadFixture, setBalance } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { airdropTestFixture } from "../materials";
import { Airdrop__factory } from "../../typechain-types";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("Airdrop", function () {
  let fixture: UnPromisify<ReturnType<typeof airdropTestFixture>>;
  beforeEach(async function () {
    fixture = await loadFixture(airdropTestFixture);
    const { airdrop, totalAirdropAmount } = fixture;

    await setBalance(airdrop.address, totalAirdropAmount);
  });
  describe("Set airdrop list", function () {
    it("#addClaim/addBatchClaims: only owner can add claim", async function () {
      const { airdrop, notClaimer, claimInfo } = fixture;

      await expect(
        airdrop.connect(notClaimer).addClaim(claimInfo[0].claimer, claimInfo[0].amount)
      ).to.be.revertedWithCustomError(airdrop, "OwnableUnauthorizedAccount");

      await expect(
        airdrop.connect(notClaimer).addBatchClaims([claimInfo[0].claimer], [claimInfo[0].amount])
      ).to.be.revertedWithCustomError(airdrop, "OwnableUnauthorizedAccount");
    });
    it("#addClaim", async function () {
      const { airdrop, claimInfo } = fixture;

      await airdrop.addClaim(claimInfo[0].claimer, claimInfo[0].amount);

      expect(await airdrop.claims(claimInfo[0].claimer)).to.equal(claimInfo[0].amount);
    });
    it("#addClaim: override existing beneficiary", async function () {
      const { airdrop, claimInfo } = fixture;

      await airdrop.addClaim(claimInfo[0].claimer, claimInfo[0].amount);

      expect(await airdrop.claims(claimInfo[0].claimer)).to.equal(claimInfo[0].amount);

      await airdrop.addClaim(claimInfo[0].claimer, BigInt(claimInfo[0].amount) * 2n);

      expect(await airdrop.claims(claimInfo[0].claimer)).to.equal(BigInt(claimInfo[0].amount) * 2n);
    });
    it("#addBatchClaims", async function () {
      const { airdrop, claimInfo } = fixture;

      await airdrop.addBatchClaims(
        claimInfo.map((claim) => claim.claimer),
        claimInfo.map((claim) => claim.amount)
      );

      for (const claim of claimInfo) {
        expect(await airdrop.claims(claim.claimer)).to.equal(claim.amount);
      }
    });
  });
  describe("Claim airdrop", function () {
    this.beforeEach(async function () {
      const { airdrop, claimInfo } = fixture;

      await airdrop.addBatchClaims(
        claimInfo.map((claim) => claim.claimer),
        claimInfo.map((claim) => claim.amount)
      );

      await airdrop.toggleClaimAllowed();
    });
    it("toggleClaimAllowed: can't claim if claim not allowed", async function () {
      const { airdrop, claimers } = fixture;

      await airdrop.toggleClaimAllowed();

      expect(await airdrop.claimAllowed()).to.be.false;

      await expect(airdrop.connect(claimers[0]).claim()).to.be.revertedWith("Airdrop: claim not allowed");

      await expect(airdrop.connect(claimers[0]).claimFor(claimers[0].address)).to.be.revertedWith(
        "Airdrop: claim not allowed"
      );
    });
    it("#claim/claimFor: can't claim if not in the list", async function () {
      const { airdrop, notClaimer, claimers } = fixture;

      await expect(airdrop.connect(notClaimer).claim()).to.be.revertedWith("Airdrop: no claimable amount");

      await expect(airdrop.connect(claimers[0]).claimFor(notClaimer.address)).to.be.revertedWith(
        "Airdrop: no claimable amount"
      );
    });
    it("#claim: successfully get airdrop", async function () {
      const { airdrop, claimers, claimInfo } = fixture;

      const beforeBalance = await hre.ethers.provider.getBalance(claimers[0].address);

      await expect(airdrop.connect(claimers[0]).claim())
        .to.emit(airdrop, "Claimed")
        .withArgs(claimers[0].address, claimInfo[0].amount);

      const afterBalance = await hre.ethers.provider.getBalance(claimers[0].address);
      // 0.0001 is a transaction fee
      expect(afterBalance.sub(beforeBalance)).to.be.closeTo(claimInfo[0].amount, hre.ethers.utils.parseEther("0.0001"));
      expect(await airdrop.claimed(claimers[0].address)).to.be.true;
    });
    it("#claimFor: successfully get airdrop", async function () {
      const { airdrop, notClaimer, claimers, claimInfo } = fixture;

      const beforeBalance = await hre.ethers.provider.getBalance(claimers[0].address);

      await expect(airdrop.connect(notClaimer).claimFor(claimers[0].address))
        .to.emit(airdrop, "Claimed")
        .withArgs(claimers[0].address, claimInfo[0].amount);

      const afterBalance = await hre.ethers.provider.getBalance(claimers[0].address);
      expect(afterBalance.sub(beforeBalance)).to.equal(claimInfo[0].amount);
      expect(await airdrop.claimed(claimers[0].address)).to.be.true;
    });
    it("#claim/claimFor: can't get twice", async function () {
      const { airdrop, notClaimer, claimers, claimInfo } = fixture;

      await expect(airdrop.connect(claimers[0]).claim())
        .to.emit(airdrop, "Claimed")
        .withArgs(claimers[0].address, claimInfo[0].amount);

      await expect(airdrop.connect(claimers[0]).claim()).to.be.revertedWith("Airdrop: already claimed");

      await expect(airdrop.connect(notClaimer).claimFor(claimers[0].address)).to.be.revertedWith(
        "Airdrop: already claimed"
      );
    });
    it("#claimBatch: successfully get airdrop", async function () {
      const { airdrop, notClaimer, claimers, claimInfo } = fixture;

      const beforeBalances = await Promise.all(
        claimers.map((claimer) => hre.ethers.provider.getBalance(claimer.address))
      );

      await airdrop.connect(notClaimer).claimBatch(claimers.map((claimer) => claimer.address));

      //   const tx = await airdrop.connect(notClaimer).claimBatch(claimers.map((claimer) => claimer.address));
      //   console.log(
      //     "Gas used for claimBatch is",
      //     (await tx.wait()).gasUsed.toString() + " when # of claimers is " + claimers.length,
      //   );

      const afterBalances = await Promise.all(
        claimers.map((claimer) => hre.ethers.provider.getBalance(claimer.address))
      );
      for (let i = 0; i < claimers.length; i++) {
        expect(afterBalances[i].sub(beforeBalances[i])).to.equal(claimInfo[i].amount);
        expect(await airdrop.claimed(claimers[i].address)).to.be.true;
      }
    });
    it("#claimBatch: failed claim", async function () {
      const { airdrop, notClaimer, claimInfo, noReceiverContract } = fixture;

      await airdrop.addClaim(noReceiverContract.address, claimInfo[0].amount);

      const beforeBalance = await hre.ethers.provider.getBalance(airdrop.address);

      await expect(airdrop.connect(notClaimer).claimFor(noReceiverContract.address)).to.be.revertedWith(
        "Airdrop: claim failed"
      );
      expect(await airdrop.claimed(noReceiverContract.address)).to.be.false;
      expect(await hre.ethers.provider.getBalance(airdrop.address)).to.equal(beforeBalance);
    });
  });
  describe("Divide airdrop contract", function () {
    it("#register airdrop contract as claimer", async function () {
      const { deployer, airdrop } = fixture;

      const newAirdrop = await new Airdrop__factory(deployer).deploy();

      await airdrop.addClaim(newAirdrop.address, hre.ethers.utils.parseEther("100"));

      await airdrop.toggleClaimAllowed();

      await airdrop.claimFor(newAirdrop.address);

      expect(await airdrop.claimed(newAirdrop.address)).to.be.true;
      expect(await hre.ethers.provider.getBalance(newAirdrop.address)).to.equal(hre.ethers.utils.parseEther("100"));
    });
  });
  describe("Check view functions", function () {
    this.beforeEach(async function () {
      const { airdrop, claimInfo } = fixture;

      for (const claim of claimInfo) {
        await airdrop.addClaim(claim.claimer, claim.amount);
      }
    });
    it("#getBeneficiariesLength", async function () {
      const { airdrop, claimInfo } = fixture;

      expect(await airdrop.getBeneficiariesLength()).to.equal(claimInfo.length);
    });
    it("#getBeneficiaryAt", async function () {
      const { airdrop, claimInfo } = fixture;

      expect(await airdrop.getBeneficiaryAt(0)).to.equal(claimInfo[0].claimer);
    });
    it("#getBeneficiaries: successfully return beneficiaries", async function () {
      const { airdrop, claimers } = fixture;

      const beneficiaries = await airdrop.getBeneficiaries(0, claimers.length);
      for (let i = 0; i < claimers.length; i++) {
        expect(beneficiaries[i]).to.equal(claimers[i].address);
      }
    });
    it("#getBeneficiaries: end > length", async function () {
      const { airdrop, claimers } = fixture;

      const beneficiaries = await airdrop.getBeneficiaries(0, claimers.length + 5);
      for (let i = 0; i < claimers.length; i++) {
        expect(beneficiaries[i]).to.equal(claimers[i].address);
      }
    });
    it("#getBeneficiaries: start > end", async function () {
      const { airdrop } = fixture;

      const beneficiaries = await airdrop.getBeneficiaries(5, 4);
      expect(beneficiaries.length).to.equal(0);
    });
  });
});
