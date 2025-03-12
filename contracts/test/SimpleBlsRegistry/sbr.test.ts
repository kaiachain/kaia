import { ethers } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { deploySbrFixture, genBLSPubkeyPop, unregisterFixture, getActors } from "../materials";

describe("SimpleBlsRegistry", function () {
  let abookAdmin: SignerWithAddress, cn0: SignerWithAddress, cn1: SignerWithAddress, misc: SignerWithAddress;

  beforeEach(async function () {
    ({ abookAdmin, cn0, cn1, misc } = await getActors());
  });

  describe("constants", function () {
    it("constants are properly set", async function () {
      const { sbr } = await loadFixture(deploySbrFixture);

      expect(await sbr.abook()).to.equal("0x0000000000000000000000000000000000000400");
      expect(await sbr.ZERO48HASH()).to.equal(ethers.utils.keccak256("0x" + Buffer.alloc(48).toString("hex")));
      expect(await sbr.ZERO96HASH()).to.equal(ethers.utils.keccak256("0x" + Buffer.alloc(96).toString("hex")));
    });
  });

  describe("register", function () {
    it("success: register", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      await expect(sbr.register(cn0.address, pk, pop)).to.emit(sbr, "Registered").withArgs(cn0.address, pk, pop);
      expect(await sbr.record(cn0.address)).to.deep.equal([pk, pop]);
      expect(await sbr.getAllBlsInfo()).to.deep.equal([[cn0.address], [[pk, pop]]]);
    });

    it("success: change public key of a registered", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);
      const [newPk, newPop] = await genBLSPubkeyPop(1);

      await expect(sbr.register(cn0.address, pk, pop)).to.emit(sbr, "Registered").withArgs(cn0.address, pk, pop);
      await expect(sbr.register(cn0.address, newPk, newPop))
        .to.emit(sbr, "Registered")
        .withArgs(cn0.address, newPk, newPop);
      expect(await sbr.record(cn0.address)).to.deep.equal([newPk, newPop]);
      expect(await sbr.getAllBlsInfo()).to.deep.equal([[cn0.address], [[newPk, newPop]]]);
    });

    it("success: registering existing public key", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      await expect(sbr.register(cn0.address, pk, pop)).to.emit(sbr, "Registered").withArgs(cn0.address, pk, pop);
      await expect(sbr.register(cn1.address, pk, pop)).to.emit(sbr, "Registered").withArgs(cn1.address, pk, pop);
    });

    it("revert: registering a public key of CN not in AddressBook", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      await expect(sbr.register(misc.address, pk, pop)).to.be.revertedWith("cnNodeId is not in AddressBook");
    });

    it("revert: onlyValidPublicKey fails due to invalid length of public key", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      // 47B pubkey
      await expect(sbr.register(cn0.address, pk.slice(0, -2), pop)).to.be.revertedWith("Public key must be 48 bytes");
      // 49B pubkey
      await expect(sbr.register(cn0.address, pk + "00", pop)).to.be.revertedWith("Public key must be 48 bytes");
    });

    it("revert: onlyValidPublicKey fails when public key is zero", async function () {
      const { sbr, pop } = await loadFixture(deploySbrFixture);
      const pk = Buffer.alloc(48);

      await expect(sbr.register(cn0.address, pk, pop)).to.be.revertedWith("Public key cannot be zero");
    });

    it("revert: onlyValidPop fails due to pop of invalid length", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      // 95B pop
      await expect(sbr.register(cn0.address, pk, pop.slice(0, -2))).to.be.revertedWith("Pop must be 96 bytes");
      // 97B pop
      await expect(sbr.register(cn0.address, pk, pop + "00")).to.be.revertedWith("Pop must be 96 bytes");
    });

    it("revert: onlyValidPop fails when pop is zero", async function () {
      const { sbr, pk } = await loadFixture(deploySbrFixture);
      const pop = Buffer.alloc(96);

      await expect(sbr.register(cn0.address, pk, pop)).to.be.revertedWith("Pop cannot be zero");
    });

    it("revert: msg.sender is not the owner", async function () {
      const { sbr, pk, pop } = await loadFixture(deploySbrFixture);

      await expect(sbr.connect(abookAdmin).register(cn0.address, pk, pop)).to.be.revertedWith(
        "Ownable: caller is not the owner",
      );
    });
  });

  describe("unregister", function () {
    it("success: initial fixture setup", async function () {
      // cn0-2 are registered
      const { sbr, cnList, pkList, popList, blsPubkeyInfoList } = await loadFixture(unregisterFixture);

      for (const [i, cn] of cnList.entries()) {
        const [pk, pop] = [pkList[i], popList[i]];
        expect(await sbr.record(cn.address)).to.deep.equal([pk, pop]);
      }
      expect(await sbr.getAllBlsInfo()).to.deep.equal([cnList.map((x) => x.address), blsPubkeyInfoList]);
    });

    it("success: unregister", async function () {
      const { sbr, cnList, pkList, popList, blsPubkeyInfoList } = await loadFixture(unregisterFixture);
      // unregister cn1 = cnList[1]
      const [cn, pk, pop] = [cnList[1], pkList[1], popList[1]];
      // remove cnList[1] and blsPubkeyInfoList[1]
      cnList.splice(1, 1);
      blsPubkeyInfoList.splice(1, 1);

      await expect(sbr.unregister(cn1.address)).to.emit(sbr, "Unregistered").withArgs(cn.address, pk, pop);
      expect(await sbr.record(cn.address)).to.deep.equal(["0x", "0x"]);
      expect(await sbr.getAllBlsInfo()).to.deep.equal([cnList.map((x) => x.address), blsPubkeyInfoList]);
    });

    it("revert: msg.sender is not the owner", async function () {
      const { sbr } = await loadFixture(unregisterFixture);

      await expect(sbr.connect(abookAdmin).unregister(misc.address)).to.be.revertedWith(
        "Ownable: caller is not the owner",
      );
    });

    it("revert: node not removed from AddressBook", async function () {
      const { sbr } = await loadFixture(unregisterFixture);

      await expect(sbr.unregister(cn0.address)).to.be.revertedWith("CN is still in AddressBook");
    });

    it("revert: unregistering non-existing node", async function () {
      const { sbr } = await loadFixture(unregisterFixture);

      await expect(sbr.unregister(misc.address)).to.be.revertedWith("CN is not registered");
    });
  });
});
