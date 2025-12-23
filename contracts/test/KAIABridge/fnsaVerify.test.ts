import { expect } from "chai";
import { encode , toWords } from "bech32";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { ethers } from "hardhat";

const { Wallet } = ethers;
const { concat, hashMessage, toUtf8Bytes, keccak256 } = ethers.utils;

describe("Fnsa", function () {
  const mnemonic = "arena click issue slot sleep tag access exotic opera pattern code coral"; // randomly generated test wallet
  const path = "m/44'/438'/0'/0/0"; // 438 is FNSA's coin_type. https://github.com/satoshilabs/slips/blob/master/slip-0044.md?plain=1
  const wallet = Wallet.fromMnemonic(mnemonic, path);
  const signingKey = wallet._signingKey();
  const walletPub = signingKey.publicKey; // 65-byte
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
    const message = "kaiabridge" + wallet.address.toLowerCase();
    const ethHash = ethers.utils.hashMessage(message);
    const sig = await wallet.signMessage(message);

    const { fnsa } = await loadFixture(deployFnsaFixture);
    const verifiedAddr = await fnsa.verify(walletPub, expectedFnsaAddr, ethHash, sig);
    expect(verifiedAddr).to.equal(wallet.address);
    const verifiedValoper = await fnsa.verify(walletPub, expectedValoperAddr, ethHash, sig);
    expect(verifiedValoper).to.equal(wallet.address);

    // Detects mismatches
    const otherFnsaAddr = "link1lt90gyw368jj9h547ehe8v9t4cupcsca9g9wc6";
    await expect(fnsa.verify(walletPub, otherFnsaAddr, ethHash, sig)).to.be.revertedWith("Invalid fnsa address");
    const otherValoperAddr = "linkvaloper1lt90gyw368jj9h547ehe8v9t4cupcscm9yywa";
    await expect(fnsa.verify(walletPub, otherValoperAddr, ethHash, sig)).to.be.revertedWith("Invalid fnsa address");

    // Signature mismatch (wrong message)
    const otherMessage = "Goodbye, world!";
    const otherSig = await wallet.signMessage(otherMessage);
    await expect(fnsa.verify(walletPub, expectedFnsaAddr, ethHash, otherSig)).to.be.revertedWith(
      "Invalid signature",
    );

    // Signature mismatch (wrong signer)
    const otherWallet = Wallet.createRandom();
    const otherSig2 = await otherWallet.signMessage(message);
    await expect(fnsa.verify(walletPub, expectedFnsaAddr, ethHash, otherSig2)).to.be.revertedWith(
      "Invalid signature",
    );
  });

  it("on-chain: verify (Klaytn prefix)", async function () {
    const message = "kaiabridge" + wallet.address.toLowerCase();
    const klayPrefix = "\x19Klaytn Signed Message:\n52";
    const klayHash = keccak256(concat([toUtf8Bytes(klayPrefix), toUtf8Bytes(message)]));
    const sig = wallet._signingKey().signDigest(klayHash);
    const sigSerialized = ethers.utils.joinSignature(sig);

    const { fnsa } = await loadFixture(deployFnsaFixture);
    const verifiedAddr = await fnsa.verify(walletPub, expectedFnsaAddr, klayHash, sigSerialized);
    expect(verifiedAddr).to.equal(wallet.address);
    const verifiedValoper = await fnsa.verify(walletPub, expectedValoperAddr, klayHash, sigSerialized);
    expect(verifiedValoper).to.equal(wallet.address);
  });

  it("on-chain: verify with messageHash rejects mismatch", async function () {
    const message = "kaiabridge" + wallet.address.toLowerCase();
    const ethHash = ethers.utils.hashMessage(message);
    const sig = await wallet.signMessage(message);
    const wrongHash = keccak256(toUtf8Bytes("wrong"));

    const { fnsa } = await loadFixture(deployFnsaFixture);
    await expect(fnsa.verify(walletPub, expectedFnsaAddr, wrongHash, sig)).to.be.revertedWith(
      "messageHash mismatch",
    );

    expect(await fnsa.verify(walletPub, expectedFnsaAddr, ethHash, sig)).to.equal(wallet.address);
  });
});
