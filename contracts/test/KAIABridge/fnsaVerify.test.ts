import { expect } from "chai";
import { encode , toWords } from "bech32";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { ethers } from "hardhat";

const { Wallet } = ethers;
const { arrayify, concat, sha256, ripemd160, hashMessage, hexlify } = ethers.utils;

describe("Fnsa", function () {
  const mnemonic = "arena click issue slot sleep tag access exotic opera pattern code coral"; // randomly generated test wallet
  const path = "m/44'/438'/0'/0/0"; // 438 is FNSA's coin_type. https://github.com/satoshilabs/slips/blob/master/slip-0044.md?plain=1
  const wallet = Wallet.fromMnemonic(mnemonic, path);
  const signingKey = wallet._signingKey();
  const walletPub = signingKey.publicKey; // 65-byte
  const bech32Prefix = "link";
  const expectedFnsaAddr = "link19sl2wemng3ayh3uwuw2xvj6zsypzacz0e6cahc";
  const expectedValoperAddr = "linkvaloper19sl2wemng3ayh3uwuw2xvj6zsypzacz0tw6qet";

  async function deployFnsaFixture() {
    const Fnsa = await ethers.getContractFactory("FnsaVerifyHarness");
    const fnsa = await Fnsa.deploy();
    return { fnsa };
  }

  it("on-chain: pub => fnsaAddr (computeFnsaAddr)", async function () {
    const { fnsa } = await loadFixture(deployFnsaFixture);
    const fnsaAddr = await fnsa.computeFnsaAddr(walletPub);
    expect(fnsaAddr).to.equal(expectedFnsaAddr);
  });

  it("on-chain: pub => ethAddr (computeEthAddr)", async function () {
    const { fnsa } = await loadFixture(deployFnsaFixture);
    const ethAddr = await fnsa.computeEthAddr(walletPub);
    expect(ethAddr).to.equal(wallet.address);
  });

  it("on-chain: pub => valoperAddr (computeValoperAddr)", async function () {
    const { fnsa } = await loadFixture(deployFnsaFixture);
    const valoperAddr = await fnsa.computeValoperAddr(walletPub);
    expect(valoperAddr).to.equal(expectedValoperAddr);
  });

  it("on-chain: sig => ethAddr (recoverEthAddr)", async function () {
    const message = "Hello, world!";
    const messageHash = hashMessage(message);
    const sig = await wallet.signMessage(message);

    const { fnsa } = await loadFixture(deployFnsaFixture);
    const ethAddr = await fnsa.recoverEthAddr(messageHash, sig);
    expect(ethAddr).to.equal(wallet.address);
  });

  it("on-chain: verify", async function () {
    const message = "Hello, world!";
    const messageHash = hashMessage(message);
    const sig = await wallet.signMessage(message);

    const { fnsa } = await loadFixture(deployFnsaFixture);
    await expect(fnsa.verify(walletPub, expectedFnsaAddr, messageHash, sig)).to.not.be.reverted;
    await expect(fnsa.verify(walletPub, expectedValoperAddr, messageHash, sig)).to.not.be.reverted;

    // Detects mismatches
    const otherFnsaAddr = "link1lt90gyw368jj9h547ehe8v9t4cupcsca9g9wc6";
    await expect(fnsa.verify(walletPub, otherFnsaAddr, messageHash, sig)).to.be.revertedWith("Invalid fnsa address");
    const otherValoperAddr = "linkvaloper1lt90gyw368jj9h547ehe8v9t4cupcscm9yywa";
    await expect(
      fnsa.verify(walletPub, otherValoperAddr, messageHash, sig),
    ).to.be.revertedWith("Invalid fnsa address");

    const otherMessage = "Goodbye, world!";
    const otherMessageHash = hashMessage(otherMessage);
    await expect(
      fnsa.verify(walletPub, expectedFnsaAddr, otherMessageHash, sig),
    ).to.be.revertedWith("Invalid signature");

    const otherWallet = Wallet.createRandom();
    const otherSig = await otherWallet.signMessage(message);
    await expect(
      fnsa.verify(walletPub, expectedFnsaAddr, messageHash, otherSig),
    ).to.be.revertedWith("Invalid signature");
  });
});
