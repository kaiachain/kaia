import { ethers } from "hardhat";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { Bridge, ServiceChainNFT, ServiceChainToken } from "../../typechain-types";
import { ContractTransaction, Transaction } from "ethers";

const VoteType = {
  ValueTransfer: 0,
  Configuration: 1,
};

const TokenType = {
  KLAY: 0,
  ERC20: 1,
  ERC721: 2,
}

const AddressZero = ethers.constants.AddressZero;
const parseEther = ethers.utils.parseEther;
const txhash = "0x1111111111111111111111111111111111111111111111111111111111111111";
const blockNum = 2222;

interface BridgeFixture {
  bridge: Bridge;
  token: ServiceChainToken;
  nft: ServiceChainNFT;
  owner: SignerWithAddress;
  op1: SignerWithAddress;
  op2: SignerWithAddress;
  op3: SignerWithAddress;
  feeReceiver: SignerWithAddress;
  user1: SignerWithAddress;
  user2: SignerWithAddress;
}

async function deployBridgeFixture(modeMintBurn: boolean = false) {
  const [owner, op1, op2, op3, cBridge, cToken, cNFT, feeReceiver, user1, user2] = await ethers.getSigners();

  const BridgeFactory = await ethers.getContractFactory("Bridge");
  const bridge = await BridgeFactory.connect(owner).deploy(modeMintBurn);

  const TokenFactory = await ethers.getContractFactory("ServiceChainToken");
  const token = await TokenFactory.deploy(bridge.address);

  const NFTFactory = await ethers.getContractFactory("ServiceChainNFT");
  const nft = await NFTFactory.deploy(bridge.address);

  return { bridge, token, nft, owner, op1, op2, op3, cBridge, cToken, cNFT, feeReceiver, user1, user2 };
}

async function deployBridgeConfiguredFixture(modeMintBurn: boolean = false) {
  const fixture = await deployBridgeFixture(modeMintBurn);
  const { bridge, token, nft, owner, op1, op2, op3, cBridge, cToken, cNFT, feeReceiver, user1, user2 } = fixture;

  await bridge.connect(owner).setCounterPartBridge(cBridge.address);
  await bridge.connect(owner).registerToken(token.address, cToken.address);
  await bridge.connect(owner).registerToken(nft.address, cNFT.address);

  await bridge.connect(owner).setKLAYFee(parseEther("0.4"), 0);
  await bridge.connect(owner).setERC20Fee(token.address, parseEther("0.3"), 1);
  await bridge.connect(owner).setFeeReceiver(feeReceiver.address);

  await bridge.connect(owner).registerOperator(op1.address);
  await bridge.connect(owner).registerOperator(op2.address);
  await bridge.connect(owner).registerOperator(op3.address);
  await bridge.connect(owner).deregisterOperator(owner.address);
  await bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 2);
  await bridge.connect(owner).setOperatorThreshold(VoteType.Configuration, 2);

  await ethers.provider.send("hardhat_setBalance", [user1.address, parseEther("5.0").toHexString()]);
  await ethers.provider.send("hardhat_setBalance", [user2.address, "0x0"]);
  await ethers.provider.send("hardhat_setBalance", [feeReceiver.address, "0x0"]);
  await token.connect(owner).transfer(user1.address, parseEther("5.0"));
  await nft.connect(owner).mintWithTokenURI(user1.address, 42, "example.com");

  return fixture;
}

async function deployBridgeConfiguredMintBurnFixture() {
  const fixture = await deployBridgeConfiguredFixture(true);
  const { bridge, owner, token, nft } = fixture;

  token.connect(owner).addMinter(bridge.address);
  nft.connect(owner).addMinter(bridge.address);

  return fixture;
}

describe("Bridge", function () {
  describe("BridgeCounterpart", function () {
    it("set counterpart bridge", async function () {
      const { bridge, owner, cBridge, user1 } = await loadFixture(deployBridgeFixture);
      expect(await bridge.counterpartBridge()).to.equal(AddressZero);
      await bridge.connect(owner).setCounterPartBridge(cBridge.address);
      expect(await bridge.counterpartBridge()).to.equal(cBridge.address);

      expect(bridge.connect(user1).setCounterPartBridge(user1.address)).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("BridgeFee", function () {
    it("set configuration", async function () {
      const { bridge, token, owner, op1, feeReceiver, user1 } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerOperator(op1.address);
      await bridge.connect(owner).setOperatorThreshold(VoteType.Configuration, 2);

      expect(await bridge.feeOfKLAY()).to.equal(0);
      expect(await bridge.feeReceiver()).to.equal(AddressZero);

      await bridge.connect(owner).setKLAYFee(7, 0);
      await bridge.connect(op1).setKLAYFee(7, 0);
      expect(await bridge.feeOfKLAY()).to.equal(7);

      await bridge.connect(owner).setERC20Fee(token.address, 77, 1);
      await bridge.connect(op1).setERC20Fee(token.address, 77, 1);
      expect(await bridge.feeOfERC20(token.address)).to.equal(77);

      await bridge.connect(owner).setFeeReceiver(feeReceiver.address);
      expect(await bridge.feeReceiver()).to.equal(feeReceiver.address);

      expect(bridge.connect(user1).setKLAYFee(token.address, 7)).to.be.revertedWith("msg.sender is not an operator");
      expect(bridge.connect(user1).setERC20Fee(token.address, 77, 1)).to.be.revertedWith("msg.sender is not an operator");
      expect(bridge.connect(op1).setFeeReceiver(feeReceiver.address)).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("BridgeOperator", function () {
    it("initial configuration with the only owner/operator", async function () {
      const { bridge, owner } = await loadFixture(deployBridgeFixture);

      expect(await bridge.owner()).to.equal(owner.address);
      expect(await bridge.MAX_OPERATOR()).to.equal(12);
      expect(await bridge.operators(owner.address)).to.equal(true);
      expect(await bridge.getOperatorList()).to.deep.equal([owner.address]);
      expect(await bridge.operatorThresholds(VoteType.ValueTransfer)).to.equal(1);
      expect(await bridge.operatorThresholds(VoteType.Configuration)).to.equal(1);
    });
    it("register operator", async function () {
      const { bridge, owner, op1, op2 } = await loadFixture(deployBridgeFixture);

      expect(await bridge.operators(op1.address)).to.equal(false);
      await bridge.connect(owner).registerOperator(op1.address);
      expect(await bridge.operators(op1.address)).to.equal(true);
      expect(await bridge.getOperatorList()).to.deep.equal([owner.address, op1.address]);

      expect(bridge.connect(op1).registerOperator(op2.address)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).registerOperator(op1.address)).to.be.revertedWith("exist operator");
    });
    it("deregister operator", async function () {
      const { bridge, owner, op1, op2, op3 } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerOperator(op1.address);
      await bridge.connect(owner).registerOperator(op2.address);
      await bridge.connect(owner).registerOperator(op3.address);
      expect(await bridge.getOperatorList()).to.deep.equal([owner.address, op1.address, op2.address, op3.address]);

      expect(await bridge.operators(op3.address)).to.equal(true);
      await bridge.connect(owner).deregisterOperator(op3.address);
      expect(await bridge.operators(op3.address)).to.equal(false);
      expect(await bridge.getOperatorList()).to.deep.equal([owner.address, op1.address, op2.address]);

      expect(bridge.connect(op1).deregisterOperator(op2.address)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).deregisterOperator(op3.address)).to.be.revertedWith("");
    });
    it("set operator threshold", async function () {
      const { bridge, owner, op1 } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerOperator(op1.address);

      expect(await bridge.operatorThresholds(VoteType.ValueTransfer)).to.equal(1);
      await bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 2);
      expect(await bridge.operatorThresholds(VoteType.ValueTransfer)).to.equal(2);

      expect(bridge.connect(op1).setOperatorThreshold(VoteType.ValueTransfer, 1)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 0)).to.be.revertedWith("zero threshold");
      expect(bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 3)).to.be.revertedWith("bigger than num of operators");
    });
    it("vote configuration", async function () {
      const { bridge, owner, op1, op2, op3, user1 } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerOperator(op1.address);
      await bridge.connect(owner).registerOperator(op2.address);
      await bridge.connect(owner).registerOperator(op3.address);
      await bridge.connect(owner).deregisterOperator(owner.address);
      await bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 2);
      await bridge.connect(owner).setOperatorThreshold(VoteType.Configuration, 2);

      expect(bridge.connect(owner).setKLAYFee(7, 0)).to.be.revertedWith("msg.sender is not an operator");
      expect(bridge.connect(user1).setKLAYFee(7, 0)).to.be.revertedWith("msg.sender is not an operator");
      expect(bridge.connect(op1).setKLAYFee(7, 1)).to.be.revertedWith("nonce mismatch");

      expect(await bridge.configurationNonce()).to.equal(0);
      expect(await bridge.feeOfKLAY()).to.equal(0);

      await bridge.connect(op1).setKLAYFee(7, 0);
      expect(await bridge.configurationNonce()).to.equal(0);
      expect(await bridge.feeOfKLAY()).to.equal(0);

      await bridge.connect(op1).setKLAYFee(9, 0); // Can vote again with different data
      expect(await bridge.configurationNonce()).to.equal(0);
      expect(await bridge.feeOfKLAY()).to.equal(0);

      await bridge.connect(op2).setKLAYFee(9, 0);
      expect(await bridge.configurationNonce()).to.equal(1);
      expect(await bridge.feeOfKLAY()).to.equal(9);

      // The configuration nonce has increased; not votable anymore.
      expect(bridge.connect(op1).setKLAYFee(9, 0)).to.be.revertedWith("nonce mismatch");
    });
    it("vote configuration interleaved", async function () {
      const { bridge, owner, op1, op2, op3, token } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerOperator(op1.address);
      await bridge.connect(owner).registerOperator(op2.address);
      await bridge.connect(owner).registerOperator(op3.address);
      await bridge.connect(owner).deregisterOperator(owner.address);
      await bridge.connect(owner).setOperatorThreshold(VoteType.ValueTransfer, 2);
      await bridge.connect(owner).setOperatorThreshold(VoteType.Configuration, 2);

      // Two operators accidentally voted on same nonce different configuration,
      // temporarily stuck but recovers by amending the votes.

      // Operator 1 votes setKLAYFee
      await bridge.connect(op1).setKLAYFee(7, 0);
      expect(await bridge.configurationNonce()).to.equal(0);
      expect(await bridge.feeOfKLAY()).to.equal(0);

      // Operator 2 votes setERC20Fee
      await bridge.connect(op2).setERC20Fee(token.address, 7, 0);
      expect(await bridge.configurationNonce()).to.equal(0);
      expect(await bridge.feeOfKLAY()).to.equal(0);

      // Operator 2 amends, now votes setKLAYFee, changing the configuration.
      await bridge.connect(op2).setKLAYFee(7, 0);
      expect(await bridge.configurationNonce()).to.equal(1);
      expect(await bridge.feeOfKLAY()).to.equal(7);
    });
  });

  describe("BridgeToken", function () {
    it("register and deregister token", async function () {
      const { bridge, token, owner, cToken, user1 } = await loadFixture(deployBridgeFixture);
      expect(await bridge.registeredTokens(token.address)).to.equal(AddressZero);

      // Test register
      expect(bridge.connect(user1).registerToken(token.address, cToken.address)).to.be.revertedWith("Ownable: caller is not the owner");

      await bridge.connect(owner).registerToken(token.address, cToken.address);
      expect(await bridge.registeredTokens(token.address)).to.equal(cToken.address);
      expect(await bridge.getRegisteredTokenList()).to.deep.equal([token.address]);

      expect(bridge.connect(owner).registerToken(token.address, cToken.address)).to.be.revertedWith("allowed token");

      // Test deregister
      expect(bridge.connect(user1).deregisterToken(token.address)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).deregisterToken(user1.address)).to.be.revertedWith("not allowed token");

      await bridge.connect(owner).deregisterToken(token.address);
      expect(await bridge.registeredTokens(token.address)).to.equal(AddressZero);
      expect(await bridge.getRegisteredTokenList()).to.deep.equal([]);
    });
    it("lock and unlock token", async function () {
      const { bridge, token, owner, cToken, user1 } = await loadFixture(deployBridgeFixture);
      await bridge.connect(owner).registerToken(token.address, cToken.address);
      expect(await bridge.getRegisteredTokenList()).to.deep.equal([token.address]);

      // Test lock
      expect(bridge.connect(user1).lockToken(token.address)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).lockToken(user1.address)).to.be.revertedWith("not allowed token");

      await bridge.connect(owner).lockToken(token.address);
      expect(await bridge.lockedTokens(token.address)).to.equal(true);

      expect(bridge.connect(owner).lockToken(token.address)).to.be.revertedWith("locked token");

      // Test unlock
      expect(bridge.connect(user1).unlockToken(token.address)).to.be.revertedWith("Ownable: caller is not the owner");
      expect(bridge.connect(owner).unlockToken(user1.address)).to.be.revertedWith("not allowed token");

      await bridge.connect(owner).unlockToken(token.address);
      expect(await bridge.lockedTokens(token.address)).to.equal(false);

      expect(bridge.connect(owner).unlockToken(token.address)).to.be.revertedWith("unlocked token");
    });
  });

  describe("BridgeTransfer", function () {
    it("setRunningStatus", async function () {
      const { bridge, owner, user1 } = await loadFixture(deployBridgeFixture);

      expect(bridge.connect(user1).setRunningStatus(false)).to.be.revertedWith("Ownable: caller is not the owner");

      expect(await bridge.isRunning()).to.equal(true);
      await bridge.connect(owner).setRunningStatus(false);
      expect(await bridge.isRunning()).to.equal(false);
    });
    it("start", async function () {
      const { bridge, owner, user1 } = await loadFixture(deployBridgeFixture);

      expect(bridge.connect(user1).start(false)).to.be.revertedWith("Ownable: caller is not the owner");

      expect(await bridge.isRunning()).to.equal(true);
      await bridge.connect(owner).start(false);
      expect(await bridge.isRunning()).to.equal(false);
    });
  });

  describe("BridgeTransferERC20", function () {
    //               before after
    // user          5.0    3.7
    // feeReceiver   0.0    0.3
    // bridge        0.0    1.0 or 0.0 depending on modeMintBurn
    async function expectSendERC20(fixture: BridgeFixture, tx: any) {
      const { bridge, token, feeReceiver, user1, user2 } = fixture;
      const modeMintBurn = await bridge.modeMintBurn();

      const expectTx = expect(await tx)
        .to.emit(bridge, "RequestValueTransfer")
        .withArgs(
          TokenType.ERC20,
          user1.address,
          user2.address,
          token.address,
          parseEther("1.0"),
          0,
          parseEther("0.3"),
          "0x",
        )
        .to.emit(token, "Transfer").withArgs(user1.address, bridge.address, parseEther("1.5")) // Send value + feeLimit
        .to.emit(token, "Transfer").withArgs(user1.address, feeReceiver.address, parseEther("0.3")) // Takeout fee
        .to.emit(token, "Transfer").withArgs(bridge.address, user1.address, parseEther("0.2")); // Refund (feeLimit - fee)
      if (modeMintBurn) {
        expectTx.to.emit(token, "Burn").withArgs(bridge.address, parseEther("1.0")); // Burn value
      } else {
        expectTx.to.not.emit(token, "Burn"); // No burn
      }

      expect(await bridge.requestNonce()).to.equal(1);
      expect(await token.balanceOf(user1.address)).to.equal(parseEther("3.7"));
      expect(await token.balanceOf(feeReceiver.address)).to.equal(parseEther("0.3"));
      if (modeMintBurn) {
        expect(await token.balanceOf(bridge.address)).to.equal(parseEther("0.0")); // Value is burned
      } else {
        expect(await token.balanceOf(bridge.address)).to.equal(parseEther("1.0")); // Bridge holds value
      }
    }
    // 2-step: first approve then request bridge
    async function expectSendERC20_2step(fixture: BridgeFixture) {
      const { bridge, token, user1, user2 } = fixture;

      await token.connect(user1).approve(bridge.address, parseEther("9999.0"));
      await expectSendERC20(fixture, bridge.connect(user1).requestERC20Transfer(
        token.address, user2.address, parseEther("1.0"), parseEther("0.5"), "0x"));
    }
    // 1-step: request directly but requires request feature in ERC20
    async function expectSendERC20_1step(fixture: BridgeFixture) {
      const { token, user1, user2 } = fixture;

      await expectSendERC20(fixture, token.connect(user1).requestValueTransfer(
        parseEther("1.0"), user2.address, parseEther("0.5"), "0x"));
    }

    async function expectReceiveERC20(fixture: BridgeFixture) {
      const { bridge, op1, op2, token, user1, user2 } = fixture;
      const modeMintBurn = await bridge.modeMintBurn();

      await bridge.connect(op1).handleERC20Transfer(txhash, user1.address, user2.address, token.address, parseEther("1.0"), 0, blockNum, "0x");
      const expectTx = expect(await bridge.connect(op2).handleERC20Transfer(txhash, user1.address, user2.address, token.address, parseEther("1.0"), 0, blockNum, "0x"))
        .to.emit(bridge, "HandleValueTransfer")
        .withArgs(
          txhash,
          TokenType.ERC20,
          user1.address,
          user2.address,
          token.address,
          parseEther("1.0"),
          0, // requestNonce
          0, // lowerHandleNonce
          "0x",
        );
      if (modeMintBurn) {
        expectTx.to.emit(token, "Mint").withArgs(user2.address, parseEther("1.0"));
      } else {
        expectTx.to.emit(token, "Transfer").withArgs(bridge.address, user2.address, parseEther("1.0"));
      }

      expect(await token.balanceOf(user2.address)).to.equal(parseEther("1.0"));
    }

    it("send 2-step no burn", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      await expectSendERC20_2step(fixture);
    });
    it("send 2-step burn", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      await expectSendERC20_2step(fixture);
    });
    it("send 1-step no burn", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      await expectSendERC20_1step(fixture);
    });
    it("send 1-step burn", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      await expectSendERC20_1step(fixture);
    });
    it("nonzero fee no refund", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, token, op1, op2, feeReceiver, user1, user2 } = fixture;

      // fee = 0.3, feeLimit = 0.3 -> no refund
      await expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.3"), "0x"))
        .to.emit(token, "Transfer").withArgs(user1.address, bridge.address, parseEther("1.3"));
      expect(await bridge.requestNonce()).to.equal(1);
      expect(await token.balanceOf(user1.address)).to.equal(parseEther("3.7"));
      expect(await token.balanceOf(feeReceiver.address)).to.equal(parseEther("0.3"));
      expect(await token.balanceOf(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("zero fee no refund", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, token, op1, op2, feeReceiver, user1, user2 } = fixture;

      await bridge.connect(op1).setERC20Fee(token.address, 0, 2);
      await bridge.connect(op2).setERC20Fee(token.address, 0, 2);

      // fee = 0, feeLimit = 0 -> no refund
      await expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.0"), "0x"))
        .to.emit(token, "Transfer").withArgs(user1.address, bridge.address, parseEther("1.0"));
      expect(await bridge.requestNonce()).to.equal(1);
      expect(await token.balanceOf(user1.address)).to.equal(parseEther("4.0"));
      expect(await token.balanceOf(feeReceiver.address)).to.equal(parseEther("0.0"));
      expect(await token.balanceOf(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("zero fee full refund", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, token, op1, op2, feeReceiver, user1, user2 } = fixture;

      await bridge.connect(op1).setERC20Fee(token.address, 0, 2);
      await bridge.connect(op2).setERC20Fee(token.address, 0, 2);

      // fee = 0, feeLimit = 0.2 -> refund feeLimit
      await expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.2"), "0x"))
        .to.emit(token, "Transfer").withArgs(user1.address, bridge.address, parseEther("1.2"))
        .to.emit(token, "Transfer").withArgs(bridge.address, user1.address, parseEther("0.2"));
      expect(await bridge.requestNonce()).to.equal(1);
      expect(await token.balanceOf(user1.address)).to.equal(parseEther("4.0"));
      expect(await token.balanceOf(feeReceiver.address)).to.equal(parseEther("0.0"));
      expect(await token.balanceOf(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("receive no mint", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, token } = fixture;

      // Add liquidity
      await token.connect(owner).transfer(bridge.address, parseEther("9999.0"));
      await expectReceiveERC20(fixture);
      expect(await token.balanceOf(bridge.address)).to.equal(parseEther("9998.0"));
    });
    it("receive mint", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      const { bridge, token } = fixture;

      // No liquidity needed
      await expectReceiveERC20(fixture);
      expect(await token.balanceOf(bridge.address)).to.equal(parseEther("0.0"));
    });
    it("insufficient feeLimit", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { token, user1, user2 } = fixture;

      expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.01"), "0x"))
        .to.be.revertedWith("insufficient feeLimit");
    });
    it("unregistered token", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, token, user1, user2 } = fixture;

      await bridge.connect(owner).deregisterToken(token.address);
      expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.5"), "0x"))
        .to.be.revertedWith("not allowed token");
    });
    it("locked token", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, token, user1, user2 } = fixture;

      await bridge.connect(owner).lockToken(token.address);
      expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.5"), "0x"))
        .to.be.revertedWith("locked token");
    });
    it("stopped bridge", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, token, user1, user2 } = fixture;

      await bridge.connect(owner).start(false);
      expect(token.connect(user1).requestValueTransfer(parseEther("1.0"), user2.address, parseEther("0.5"), "0x"))
        .to.be.revertedWith("stopped bridge");
    });
    it("zero value", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, token, user1, user2 } = fixture;

      await bridge.connect(owner).start(false);
      expect(token.connect(user1).requestValueTransfer(0, user2.address, 0, "0x"))
        .to.be.revertedWith("zero msg.value");
    });
  });

  describe("BridgeTransferERC721", function () {
    async function expectSendERC721(fixture: BridgeFixture, tx: any) {
      const { bridge, nft, user1, user2 } = fixture;
      const modeMintBurn = await bridge.modeMintBurn();

      const expectTx = expect(await tx)
        .to.emit(bridge, "RequestValueTransferEncoded")
        .withArgs(
          TokenType.ERC721,
          user1.address,
          user2.address,
          nft.address,
          42,
          0,
          0,
          "0x",
          2,
          "example.com",
        )
        .to.emit(nft, "Transfer").withArgs(user1.address, bridge.address, 42);
      if (modeMintBurn) {
        expectTx.to.emit(nft, "Burn").withArgs(bridge.address, 42);
      } else {
        expectTx.to.not.emit(nft, "Burn");
      }

      expect(await bridge.requestNonce()).to.equal(1);
      if (modeMintBurn) {
        expect(nft.ownerOf(42)).to.be.revertedWith("ERC721: owner query for nonexistent token");
      } else {
        expect(await nft.ownerOf(42)).to.equal(bridge.address);
      }
    }
    // 2-step: first approve then request bridge
    async function expectSendERC721_2step(fixture: BridgeFixture) {
      const { bridge, nft, user1, user2 } = fixture;

      await nft.connect(user1).approve(bridge.address, 42);
      await expectSendERC721(fixture, bridge.connect(user1).requestERC721Transfer(nft.address, user2.address, 42, "0x"));
    }
    // 1-step: request directly but requires request feature in ERC721
    async function expectSendERC721_1step(fixture: BridgeFixture) {
      const { nft, user1, user2 } = fixture;

      await expectSendERC721(fixture, nft.connect(user1).requestValueTransfer(42, user2.address, "0x"));
    }

    async function expectReceiveERC721(fixture: BridgeFixture) {
      const { bridge, op1, op2, nft, user1, user2 } = fixture;
      const modeMintBurn = await bridge.modeMintBurn();

      await bridge.connect(op1).handleERC721Transfer(txhash, user1.address, user2.address, nft.address, 99, 0, blockNum, "example.com", "0x");
      const expectTx = expect(await bridge.connect(op2).handleERC721Transfer(txhash, user1.address, user2.address, nft.address, 99, 0, blockNum, "example.com", "0x"))
        .to.emit(bridge, "HandleValueTransfer")
        .withArgs(
          txhash,
          TokenType.ERC721,
          user1.address,
          user2.address,
          nft.address,
          99,
          0, // requestNonce
          0, // lowerHandleNonce
          "0x",
        )
      if (modeMintBurn) {
        expectTx.to.emit(nft, "Mint").withArgs(user2.address, 99);
      } else {
        expectTx.to.emit(nft, "Transfer").withArgs(bridge.address, user2.address, 99);
      }

      expect(await nft.ownerOf(99)).to.equal(user2.address);
      expect(await nft.tokenURI(99)).to.equal("example.com");
    }

    it("send 2-step no burn", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      await expectSendERC721_2step(fixture);
    });
    it("send 2-step burn", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      await expectSendERC721_2step(fixture);
    });
    it("send 1-step no burn", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      await expectSendERC721_1step(fixture);
    });
    it("send 1-step burn", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      await expectSendERC721_1step(fixture);
    });
    it("receive no mint", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, nft } = fixture;

      // Bridge must have the token beforehand
      await nft.connect(owner).mintWithTokenURI(bridge.address, 99, "example.com");
      await expectReceiveERC721(fixture);
    });
    it("receive mint", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredMintBurnFixture);
      await expectReceiveERC721(fixture);
    });
    it("unregistered token", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, nft, user1, user2 } = fixture;

      await bridge.connect(owner).deregisterToken(nft.address);
      expect(nft.connect(user1).requestValueTransfer(42, user2.address, "0x"))
        .to.be.revertedWith("not allowed token");
    });
    it("locked token", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, nft, user1, user2 } = fixture;

      await bridge.connect(owner).lockToken(nft.address);
      expect(nft.connect(user1).requestValueTransfer(42, user2.address, "0x"))
        .to.be.revertedWith("locked token");
    });
    it("stopped bridge", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, nft, user1, user2 } = fixture;

      await bridge.connect(owner).start(false);
      expect(nft.connect(user1).requestValueTransfer(42, user2.address, "0x"))
        .to.be.revertedWith("stopped bridge");
    });
  });

  describe("BridgeTransferKLAY", function () {
    //               before after
    // user          5.0    3.6 - gasFee
    // feeReceiver   0.0    0.4
    // bridge        0.0    1.0
    async function gasFee(sentTx: ContractTransaction) {
      const receipt = await sentTx.wait();
      return receipt.gasUsed.mul(receipt.effectiveGasPrice);
    }
    async function expectSendKLAY(fixture: BridgeFixture, tx: Promise<ContractTransaction>) {
      const { bridge, feeReceiver, user1 } = fixture;

      const sentTx = await tx;
      expect(sentTx)
        .to.emit(bridge, "RequestValueTransfer")
        .withArgs(
          TokenType.KLAY,
          user1.address,
          user1.address,
          AddressZero,
          parseEther("1.0"), // -- eventArg
          0,
          parseEther("0.4"),
          "0x",
        )
      // When using requestKLAYTransfer, eventArg = msg.value - feeLimit = msg.value - (msg.value - _value) = _value
      // When using fallback, eventArg = msg.value - feeLimit = msg.value - fee.

      expect(await bridge.requestNonce()).to.equal(1);
      expect(await ethers.provider.getBalance(user1.address)).to.equal(parseEther("3.6").sub(await gasFee(sentTx)));
      expect(await ethers.provider.getBalance(feeReceiver.address)).to.equal(parseEther("0.4"));
      expect(await ethers.provider.getBalance(bridge.address)).to.equal(parseEther("1.0"));
    }

    it("send function", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, user1, user2 } = fixture;

      // _value = 1.0, msg.value = 1.5 => feeLimit = 0.5, eventArg = 1.0
      await expectSendKLAY(fixture, bridge.connect(user1).requestKLAYTransfer(
        user2.address,
        parseEther("1.0"),
        "0x",
        {value: parseEther("1.5")})
      );
    });
    it("send fallback", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, user1 } = fixture;

      await expectSendKLAY(fixture, bridge.connect(user1).fallback({value: parseEther("1.4")}));
    });
    it("zero fee no refund", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, op1, op2, feeReceiver, user1 } = fixture;

      // Set fee = 0
      await bridge.connect(op1).setKLAYFee(0, 2);
      await bridge.connect(op2).setKLAYFee(0, 2);

      // fee = 0, feeLimit = 0 -> no refund
      const sentTx = await bridge.connect(user1).fallback({value: parseEther("1.0")});
      expect(await ethers.provider.getBalance(user1.address)).to.equal(parseEther("4.0").sub(await gasFee(sentTx)));
      expect(await ethers.provider.getBalance(feeReceiver.address)).to.equal(parseEther("0.0"));
      expect(await ethers.provider.getBalance(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("zero fee full refund", async function () {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, op1, op2, feeReceiver, user1, user2 } = fixture;

      // Set fee = 0
      await bridge.connect(op1).setKLAYFee(0, 2);
      await bridge.connect(op2).setKLAYFee(0, 2);

      // fee = 0, feeLimit = 0.2 -> refund feeLimit
      const sentTx = await bridge.connect(user1).requestKLAYTransfer(
        user2.address, parseEther("1.0"), "0x", {value: parseEther("1.2")});
      expect(await ethers.provider.getBalance(user1.address)).to.equal(parseEther("4.0").sub(await gasFee(sentTx)));
      expect(await ethers.provider.getBalance(feeReceiver.address)).to.equal(parseEther("0.0"));
      expect(await ethers.provider.getBalance(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("receive", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, op1, op2, user1, user2 } = fixture;

      // Add liquidity
      await bridge.connect(owner).chargeWithoutEvent({value: parseEther("1000.0")});

      expect(bridge.connect(op1).handleKLAYTransfer(txhash, user1.address, user2.address, parseEther("1.0"), 0, blockNum, "0x"))
        .to.be.revertedWith("msg.sender is not an operator");

      await bridge.connect(op1).handleKLAYTransfer(txhash, user1.address, user2.address, parseEther("1.0"), 0, blockNum, "0x");
      expect(await bridge.connect(op2).handleKLAYTransfer(txhash, user1.address, user2.address, parseEther("1.0"), 0, blockNum, "0x"))
        .to.emit(bridge, "HandleValueTransfer")
        .withArgs(
          txhash,
          TokenType.KLAY,
          user1.address,
          user2.address,
          AddressZero,
          parseEther("1.0"),
          0, // requestNonce
          0, // lowerHandleNonce
          "0x",
        );

      expect(await ethers.provider.getBalance(user2.address)).to.equal(parseEther("1.0"));
    });
    it("insufficient feeLimit", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, user1, user2 } = fixture;

      expect(bridge.connect(user1).fallback({value: parseEther("0.01")}))
        .to.be.revertedWith("insufficient feeLimit");
      expect(bridge.connect(user1).requestKLAYTransfer(user2.address, parseEther("1.0"), "0x", {value: parseEther("1.01")}))
        .to.be.revertedWith("insufficient feeLimit");
    });
    it("locked KLAY", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, user1 } = fixture;

      expect(bridge.connect(user1).lockKLAY()).to.be.revertedWith("Ownable: caller is not the owner");
      await bridge.connect(owner).lockKLAY();
      expect(bridge.connect(owner).lockKLAY()).to.be.revertedWith("locked");

      // Transfer not working when locked
      expect(bridge.connect(user1).fallback({value: parseEther("1.0")})).to.be.revertedWith("locked");

      expect(bridge.connect(user1).unlockKLAY()).to.be.revertedWith("Ownable: caller is not the owner");
      await bridge.connect(owner).unlockKLAY();
      expect(bridge.connect(owner).unlockKLAY()).to.be.revertedWith("unlocked");

      // Transfer working when unlocked
      await bridge.connect(user1).fallback({value: parseEther("1.4")});
      expect(await ethers.provider.getBalance(bridge.address)).to.equal(parseEther("1.0"));
    });
    it("stopped bridge", async function() {
      const fixture = await loadFixture(deployBridgeConfiguredFixture);
      const { bridge, owner, user1 } = fixture;

      await bridge.connect(owner).start(false);
      expect(bridge.connect(user1).fallback({value: parseEther("1.0")}))
        .to.be.revertedWith("stopped bridge");
    });
  });
});
