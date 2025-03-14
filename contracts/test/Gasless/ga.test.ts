import { expect } from "chai";
import { ethers, network } from "hardhat";
import { Contract } from "ethers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { parseEther } from "@ethersproject/units";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";

import {
  gaslessSwapRouterFixture,
  gaslessSwapRouterAddTokenFixture,
  gaslessSwapRouterMultiTokenFixture,
} from "../materials";

describe("GaslessSwapRouter", function () {
  describe("Token Management", function () {
    it("should add token successfully", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      expect(await gaslessRouter.isTokenSupported(testToken.address)).to.be
        .true;
      expect(await gaslessRouter.dexAddress(testToken.address)).to.equal(
        uniswapFactory.address
      );
    });

    it("should fail to add token if already supported", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await expect(
        gaslessRouter.addToken(
          testToken.address,
          uniswapFactory.address,
          uniswapRouter.address
        )
      ).to.be.revertedWith("TokenAlreadySupported");
    });

    it("should remove token successfully", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await gaslessRouter.removeToken(testToken.address);
      expect(await gaslessRouter.isTokenSupported(testToken.address)).to.be
        .false;
    });

    it("should fail to remove token if not supported", async function () {
      const { gaslessRouter, testToken } = await loadFixture(
        gaslessSwapRouterFixture
      );

      await expect(
        gaslessRouter.removeToken(testToken.address)
      ).to.be.revertedWith("TokenNotSupported");
    });

    it("should get all supported tokens", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      const supportedTokens = await gaslessRouter.getSupportedTokens();
      expect(supportedTokens).to.include(testToken.address);
      expect(supportedTokens.length).to.equal(1);
    });

    it("should add and remove multiple tokens in various orders", async function () {
      const {
        gaslessRouter,
        testToken,
        testUser,
        wkaia,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      // Deploy multiple tokens
      const TokenFactory = await ethers.getContractFactory("TestToken");
      const token1 = await TokenFactory.deploy(testUser.address);
      const token2 = await TokenFactory.deploy(testUser.address);
      const token3 = await TokenFactory.deploy(testUser.address);
      await token1.deployed();
      await token2.deployed();
      await token3.deployed();

      // Create pairs for all tokens
      await uniswapFactory.createPair(token1.address, wkaia.address);
      await uniswapFactory.createPair(token2.address, wkaia.address);
      await uniswapFactory.createPair(token3.address, wkaia.address);

      // Add all tokens to the router
      await gaslessRouter.addToken(
        token1.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await gaslessRouter.addToken(
        token2.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await gaslessRouter.addToken(
        token3.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      // Verify all tokens are supported
      let supportedTokens = await gaslessRouter.getSupportedTokens();
      expect(supportedTokens.length).to.equal(4);
      expect(supportedTokens).to.include(token1.address);
      expect(supportedTokens).to.include(token2.address);
      expect(supportedTokens).to.include(token3.address);
      expect(supportedTokens).to.include(testToken.address);

      // Remove tokens in different order
      await gaslessRouter.removeToken(token2.address);
      supportedTokens = await gaslessRouter.getSupportedTokens();
      expect(supportedTokens.length).to.equal(3);
      expect(supportedTokens).to.include(token1.address);
      expect(supportedTokens).to.include(token3.address);
      expect(supportedTokens).to.include(testToken.address);
      expect(supportedTokens).to.not.include(token2.address);

      // Remove token from the middle of the array
      await gaslessRouter.removeToken(token1.address);
      supportedTokens = await gaslessRouter.getSupportedTokens();
      expect(supportedTokens.length).to.equal(2);
      expect(supportedTokens).to.not.include(token1.address);
      expect(supportedTokens).to.include(token3.address);
      expect(supportedTokens).to.include(testToken.address);

      // Try to get DEX info for removed token
      await expect(gaslessRouter.getDEXInfo(token1.address)).to.be.revertedWith(
        "TokenNotSupported"
      );

      // Verify remaining tokens still work
      const dexInfo = await gaslessRouter.getDEXInfo(token3.address);
      expect(dexInfo.factory).to.equal(uniswapFactory.address);
      expect(dexInfo.router).to.equal(uniswapRouter.address);
    });

    it("should fail to add zero address as token", async function () {
      const { gaslessRouter, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      await expect(
        gaslessRouter.addToken(
          ethers.constants.AddressZero,
          uniswapFactory.address,
          uniswapRouter.address
        )
      ).to.be.revertedWith("Invalid token address");
    });

    it("should fail when adding token with zero address router or factory", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterFixture);

      // Test with zero router
      await expect(
        gaslessRouter.addToken(
          testToken.address,
          uniswapFactory.address,
          ethers.constants.AddressZero
        )
      ).to.be.revertedWith("Invalid router address");

      // Test with zero factory
      await expect(
        gaslessRouter.addToken(
          testToken.address,
          ethers.constants.AddressZero,
          uniswapRouter.address
        )
      ).to.be.revertedWith("Invalid factory address");
    });

    it("should fail to add token with invalid DEX address", async function () {
      const { gaslessRouter, testUser, uniswapRouter } = await loadFixture(
        gaslessSwapRouterFixture
      );

      const MockFactory = await ethers.getContractFactory("MockInvalidFactory");
      const mockFactory = await MockFactory.deploy();
      await mockFactory.deployed();

      const anotherToken = await (
        await ethers.getContractFactory("TestToken")
      ).deploy(testUser.address);
      await anotherToken.deployed();

      await expect(
        gaslessRouter.addToken(
          anotherToken.address,
          mockFactory.address,
          uniswapRouter.address
        )
      ).to.be.revertedWith("InvalidDEXAddress");
    });

    it("should handle multiple token additions and removals", async function () {
      const {
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
        wkaia,
      } = await loadFixture(gaslessSwapRouterFixture);

      // First make sure testToken is already added
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      // Deploy additional token
      const additionalToken = await (
        await ethers.getContractFactory("TestToken")
      ).deploy(testUser.address);
      await additionalToken.deployed();

      // Create pair for additional token
      await uniswapFactory.createPair(additionalToken.address, wkaia.address);

      // Add the additional token
      await gaslessRouter.addToken(
        additionalToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      // Get supported tokens and verify both tokens are included
      const supportedTokens = await gaslessRouter.getSupportedTokens();
      expect(supportedTokens).to.include(testToken.address);
      expect(supportedTokens).to.include(additionalToken.address);
      expect(supportedTokens.length).to.equal(2);

      // Remove additional token
      await gaslessRouter.removeToken(additionalToken.address);

      // Verify only testToken remains
      const updatedSupportedTokens = await gaslessRouter.getSupportedTokens();
      expect(updatedSupportedTokens).to.include(testToken.address);
      expect(updatedSupportedTokens).to.not.include(additionalToken.address);
      expect(updatedSupportedTokens.length).to.equal(1);
    });
  });

  describe("Commission Management", function () {
    it("should update commission rate", async function () {
      const { gaslessRouter } = await loadFixture(gaslessSwapRouterFixture);

      const newRate = 500; // 5%
      await gaslessRouter.updateCommissionRate(newRate);
      expect(await gaslessRouter.commissionRate()).to.equal(newRate);
    });

    it("should fail to update commission rate if too high", async function () {
      const { gaslessRouter } = await loadFixture(gaslessSwapRouterFixture);

      const tooHighRate = 11000; // 110%
      await expect(
        gaslessRouter.updateCommissionRate(tooHighRate)
      ).to.be.revertedWith("InvalidCommissionRate");
    });

    it("should fail to claim commission if zero", async function () {
      const { gaslessRouter } = await loadFixture(gaslessSwapRouterFixture);

      await expect(gaslessRouter.claimCommission()).to.be.revertedWith(
        "NoCommissionToWithdraw"
      );
    });

    it("should accumulate and claim commission successfully", async function () {
      const {
        deployer,
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      // Set commission rate to 10%
      await gaslessRouter.updateCommissionRate(1000);

      // Add token and approve for swap
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      const swapAmount = parseEther("10.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      // Calculate gas repayment
      const feeData = await ethers.provider.getFeeData();
      const gasPriceBN = feeData.gasPrice!;
      const amountRepay = gasPriceBN.mul(600000); // Rough estimate

      // Execute swap
      const minAmountOut = amountRepay.add(parseEther("1")); // Ensure enough output
      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      // Check accumulated commission
      const accumulatedCommission = await ethers.provider.getBalance(
        gaslessRouter.address
      );
      expect(accumulatedCommission).to.be.gt(0);

      // Get contract balance
      const contractBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );
      expect(contractBalance).to.be.gt(0);

      // Claim commission
      const initialBalance = await deployer.getBalance();
      await gaslessRouter.claimCommission();
      const finalBalance = await deployer.getBalance();

      expect(finalBalance.sub(initialBalance)).to.be.gt(0);
      expect(await ethers.provider.getBalance(gaslessRouter.address)).to.equal(
        0
      );
    });

    it("should handle commission rate at boundary (10000)", async function () {
      const {
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.updateCommissionRate(10000); // 100%

      const swapAmount = parseEther("1.0");
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("0.1"));

      const initialBalance = await testUser.getBalance();
      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const finalBalance = await testUser.getBalance();
      expect(finalBalance).to.be.lt(initialBalance.add(amountRepay));
    });

    it("should handle zero commission correctly", async function () {
      const {
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      // Ensure commission rate is 0
      await gaslessRouter.updateCommissionRate(0);

      const swapAmount = parseEther("1.0");
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("0.1"));

      const initialBalance = await testUser.getBalance();

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const finalBalance = await testUser.getBalance();
      expect(finalBalance).to.be.gt(initialBalance);
      expect(await ethers.provider.getBalance(gaslessRouter.address)).to.equal(
        0
      );
    });

    it("should calculate commission precisely", async function () {
      const {
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
        wkaia,
      } = await loadFixture(gaslessSwapRouterFixture);

      await gaslessRouter.updateCommissionRate(500); // 5%

      const swapAmount = parseEther("10.0");
      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("1"));

      const path = [testToken.address, wkaia.address];
      const [, expectedOutput] = await uniswapRouter.getAmountsOut(
        swapAmount,
        path
      );
      const userAmount = expectedOutput.sub(amountRepay);

      const commission = userAmount.mul(500).div(10000);

      const preCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const postCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );
      const actualCommission = postCommissionBalance.sub(preCommissionBalance);

      expect(actualCommission).to.equal(commission);
    });
  });

  describe("Ownership Management", function () {
    it("should transfer ownership correctly", async function () {
      const {
        deployer,
        gaslessRouter,
        testUser,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      const newOwner = testUser;
      await gaslessRouter.transferOwnership(newOwner.address);

      await expect(
        gaslessRouter
          .connect(newOwner)
          .addToken(
            ethers.Wallet.createRandom().address,
            uniswapFactory.address,
            uniswapRouter.address
          )
      ).to.not.be.reverted;

      await expect(
        gaslessRouter
          .connect(deployer)
          .addToken(
            ethers.Wallet.createRandom().address,
            uniswapFactory.address,
            uniswapRouter.address
          )
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("should fail to claim commission by non-owner", async function () {
      const { gaslessRouter, testUser } = await loadFixture(
        gaslessSwapRouterFixture
      );

      await expect(
        gaslessRouter.connect(testUser).claimCommission()
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("should comprehensively test ownership functions", async function () {
      const {
        deployer,
        gaslessRouter,
        testToken,
        testUser,
        uniswapFactory,
        uniswapRouter,
      } = await loadFixture(gaslessSwapRouterFixture);

      // Test initial ownership
      expect(await gaslessRouter.owner()).to.equal(deployer.address);

      // Test transferring ownership
      await gaslessRouter.transferOwnership(testUser.address);
      expect(await gaslessRouter.owner()).to.equal(testUser.address);

      // Test function access by previous owner (should fail)
      await expect(
        gaslessRouter
          .connect(deployer)
          .addToken(
            testToken.address,
            uniswapFactory.address,
            uniswapRouter.address
          )
      ).to.be.revertedWith("Ownable: caller is not the owner");

      // Transfer ownership back
      await gaslessRouter.connect(testUser).transferOwnership(deployer.address);
      expect(await gaslessRouter.owner()).to.equal(deployer.address);

      // Test various privileged functions with non-owner
      await expect(
        gaslessRouter.connect(testUser).updateCommissionRate(500)
      ).to.be.revertedWith("Ownable: caller is not the owner");

      await expect(
        gaslessRouter.connect(testUser).claimCommission()
      ).to.be.revertedWith("Ownable: caller is not the owner");

      await expect(
        gaslessRouter
          .connect(testUser)
          .removeToken(ethers.constants.AddressZero)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("Basic Functions", function () {
    it("should receive KAIA successfully", async function () {
      const { gaslessRouter, testUser } = await loadFixture(
        gaslessSwapRouterFixture
      );

      const amount = parseEther("1.0");
      const initialBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );

      await testUser.sendTransaction({
        to: gaslessRouter.address,
        value: amount,
      });

      const finalBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );
      expect(finalBalance.sub(initialBalance)).to.equal(amount);
    });
  });

  describe("DEX Information Management", function () {
    it("should handle all DEX info retrieval edge cases", async function () {
      const { gaslessRouter, testToken, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);

      // Get DEX info for valid token
      const dexInfo = await gaslessRouter.getDEXInfo(testToken.address);
      expect(dexInfo.factory).to.equal(uniswapFactory.address);
      expect(dexInfo.router).to.equal(uniswapRouter.address);

      // Try to get DEX info for non-existent token
      const randomToken = ethers.Wallet.createRandom().address;
      await expect(gaslessRouter.getDEXInfo(randomToken)).to.be.revertedWith(
        "TokenNotSupported"
      );

      // Try to get DEX address for non-existent token
      await expect(gaslessRouter.dexAddress(randomToken)).to.be.revertedWith(
        "TokenNotSupported"
      );

      // Check if token is supported
      expect(await gaslessRouter.isTokenSupported(testToken.address)).to.be
        .true;
      expect(await gaslessRouter.isTokenSupported(randomToken)).to.be.false;
    });

    it("should revert when getting dex address for unsupported token", async function () {
      const { gaslessRouter } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );

      const unknownToken = ethers.Wallet.createRandom().address;
      await expect(gaslessRouter.dexAddress(unknownToken)).to.be.revertedWith(
        "TokenNotSupported"
      );
    });
  });
});

describe("GaslessSwapRouter: Swap Operations", function () {
  describe("Core Swap Functionality", function () {
    it("should calculate amount in correctly", async function () {
      const { gaslessRouter, testToken } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const amountOut = parseEther("1");
      const amountIn = await gaslessRouter.getAmountIn(
        testToken.address,
        amountOut
      );
      expect(amountIn).to.be.gt(0);
    });

    it("should execute swap successfully", async function () {
      const { gaslessRouter, testToken, testUser, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      // Calculate gas repayment
      const feeData = await ethers.provider.getFeeData();
      const gasPriceBN = feeData.gasPrice!;
      const R1 = gasPriceBN.mul(21000);
      const R2 = ethers.BigNumber.from(200000);
      const R3 = gasPriceBN.mul(300000);
      const amountRepay = R1.add(R2).add(R3);

      // Get expected output
      const path = [testToken.address, wkaia.address];
      const [, expectedOutput] = await uniswapRouter.getAmountsOut(
        swapAmount,
        path
      );
      const margin = expectedOutput.mul(1).div(100); // 1% margin
      const minAmountOut = amountRepay.add(margin);

      const initialBalance = await testUser.getBalance();
      const tx = await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);
      await tx.wait();

      const finalBalance = await testUser.getBalance();
      expect(finalBalance).to.be.gt(initialBalance);
    });

    it("should fail swap if token not supported", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const invalidToken = ethers.Wallet.createRandom().address;

      // First approve some tokens to avoid approval revert
      const amount = parseEther("1");
      await testToken.connect(testUser).approve(gaslessRouter.address, amount);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(invalidToken, amount, amount, parseEther("0.1"))
      ).to.be.revertedWith("TokenNotSupported");
    });

    it("should prevent operations on removed token", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      await gaslessRouter.removeToken(testToken.address);

      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            swapAmount,
            parseEther("0.5"),
            amountRepay
          )
      ).to.be.revertedWith("TokenNotSupported");
    });

    it("should handle multiple consecutive swaps", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      await gaslessRouter.updateCommissionRate(500); // 5%

      const swapAmount = parseEther("5.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount.mul(2));

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("1"));

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const initialCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const finalCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );
      expect(finalCommissionBalance).to.be.gt(initialCommissionBalance);
    });

    it("should verify swap balances and commission", async function () {
      const { gaslessRouter, testToken, testUser, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      await gaslessRouter.updateCommissionRate(1000); // 10%

      const swapAmount = parseEther("10.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      const path = [testToken.address, wkaia.address];
      const [, expectedOutput] = await uniswapRouter.getAmountsOut(
        swapAmount,
        path
      );

      const userAmount = expectedOutput.sub(amountRepay);
      const expectedCommission = userAmount.mul(1000).div(10000);
      const expectedUserFinalAmount = userAmount.sub(expectedCommission);

      const initialUserBalance = await testUser.getBalance();
      const initialCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );

      const tx = await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, expectedOutput, amountRepay);
      await tx.wait();

      const finalUserBalance = await testUser.getBalance();
      const finalCommissionBalance = await ethers.provider.getBalance(
        gaslessRouter.address
      );

      expect(finalUserBalance.sub(initialUserBalance)).to.be.closeTo(
        expectedUserFinalAmount,
        parseEther("0.1")
      );
      expect(finalCommissionBalance.sub(initialCommissionBalance)).to.equal(
        expectedCommission
      );
    });
  });

  describe("Swap Parameter Validation", function () {
    it("should fail swap if insufficient output", async function () {
      const { gaslessRouter, testToken, testUser, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      // Get expected output
      const path = [testToken.address, wkaia.address];
      const [, expectedOutput] = await uniswapRouter.getAmountsOut(
        swapAmount,
        path
      );

      // Calculate minimum output that's much higher than possible
      const tooHighMinOutput = expectedOutput.mul(2);
      const amountRepay = parseEther("0.1");

      // Min output should be greater than amountRepay
      expect(tooHighMinOutput).to.be.gt(amountRepay);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            swapAmount,
            tooHighMinOutput,
            amountRepay
          )
      ).to.be.revertedWith("UniswapV2Router: INSUFFICIENT_OUTPUT_AMOUNT");
    });

    it("should revert when minAmountOut is less than amountRepay", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const amountRepay = parseEther("0.1");
      const minAmountOut = parseEther("0.05"); // Less than amountRepay

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay)
      ).to.be.revertedWith("InsufficientSwapOutput");
    });

    it("should revert with insufficient token balance", async function () {
      const {
        deployer,
        gaslessRouter,
        testUser,
        uniswapFactory,
        uniswapRouter,
        wkaia,
      } = await loadFixture(gaslessSwapRouterAddTokenFixture);
      // Deploy a new token with no balance
      const NoBalanceToken = await ethers.getContractFactory("TestToken");
      const noBalanceToken = await NoBalanceToken.deploy(deployer.address);
      await noBalanceToken.deployed();

      // Create pair
      await uniswapFactory.createPair(noBalanceToken.address, wkaia.address);

      // Add token
      await gaslessRouter.addToken(
        noBalanceToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      // Try to swap with no balance
      const swapAmount = parseEther("1.0");
      await noBalanceToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            noBalanceToken.address,
            swapAmount,
            parseEther("0.1"),
            parseEther("0.05")
          )
      ).to.be.revertedWith("Insufficient token balance");
    });

    it("should handle token approval limits", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const maxApproval = ethers.constants.MaxUint256;
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, maxApproval);

      const swapAmount = parseEther("10.0");

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("1"));

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const initialBalance = await testToken.balanceOf(testUser.address);
      expect(initialBalance).to.be.gt(0);
    });

    it("should fail swap if gas price is too high", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      try {
        const swapAmount = parseEther("1.0");
        await testToken
          .connect(testUser)
          .approve(gaslessRouter.address, swapAmount);

        await expect(
          gaslessRouter
            .connect(testUser)
            .swapForGas(
              testToken.address,
              swapAmount,
              parseEther("1000"),
              parseEther("0.1")
            )
        ).to.be.reverted;
      } catch (e) {
        const error = e as Error;
        expect(error.message).to.include("revert");
      }
    });
  });

  describe("Swap Edge Cases", function () {
    it("should handle extreme input in getAmountIn", async function () {
      const { gaslessRouter, testToken } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const largeAmountOut = parseEther("0.00000001");
      const amountIn = await gaslessRouter.getAmountIn(
        testToken.address,
        largeAmountOut
      );
      expect(amountIn).to.be.gt(0);
    });

    it("should handle getAmountIn with extremely large output", async function () {
      const { gaslessRouter, testToken } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const extremeAmountOut = ethers.constants.MaxUint256;

      await expect(
        gaslessRouter.getAmountIn(testToken.address, extremeAmountOut)
      ).to.be.reverted;
    });

    it("should handle getAmountIn with zero output", async function () {
      const { gaslessRouter, testToken } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      await expect(gaslessRouter.getAmountIn(testToken.address, 0)).to.be
        .reverted;
    });

    it("should handle very small swap amount", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const verySmallAmount = parseEther("0.000001");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, verySmallAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            verySmallAmount,
            amountRepay,
            amountRepay.div(2)
          )
      ).to.be.reverted;
    });

    it("should handle very large swap amounts", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      const largeSwapAmount = parseEther("1000000");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, largeSwapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            largeSwapAmount,
            parseEther("0.5"),
            amountRepay
          )
      ).to.be.reverted;
    });

    it("should fail if coin transfer fails", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      // Create a mock coinbase that rejects payments
      const MockCoinbase = await ethers.getContractFactory("MockCoinbase");
      const mockCoinbase = await MockCoinbase.deploy();

      // Modify block.coinbase to return mock address
      await network.provider.send("hardhat_setCoinbase", [
        mockCoinbase.address,
      ]);

      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            swapAmount,
            parseEther("0.5"),
            parseEther("0.1")
          )
      ).to.be.revertedWith("Failed to send KAIA to proposer");

      // Reset coinbase
      await network.provider.send("hardhat_setCoinbase", [
        "0x0000000000000000000000000000000000000000",
      ]);
    });

    it("should fail if transfer to user fails", async function () {
      const { gaslessRouter, testToken, testUser, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      // Deploy a malicious contract that rejects KAIA transfers
      const MaliciousReceiver = await ethers.getContractFactory(
        "MaliciousReceiver"
      );
      const maliciousReceiver = await MaliciousReceiver.deploy();

      // Impersonate the malicious contract
      await network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [maliciousReceiver.address],
      });

      // Fund the malicious account with ETH for gas fees
      await network.provider.send("hardhat_setBalance", [
        maliciousReceiver.address,
        ethers.utils.hexValue(parseEther("10.0")),
      ]);

      const maliciousSigner = await ethers.getSigner(maliciousReceiver.address);

      // Fund the malicious contract with some test tokens
      await testToken
        .connect(testUser)
        .transfer(maliciousReceiver.address, parseEther("10"));
      await testToken
        .connect(maliciousSigner)
        .approve(gaslessRouter.address, parseEther("1"));

      // Get expected output for setting appropriate minAmountOut
      const path = [testToken.address, wkaia.address];
      await uniswapRouter.getAmountsOut(parseEther("1"), path);
      const amountRepay = parseEther("0.1");

      await expect(
        gaslessRouter
          .connect(maliciousSigner)
          .swapForGas(
            testToken.address,
            parseEther("1"),
            amountRepay.add(parseEther("0.1")),
            amountRepay
          )
      ).to.be.revertedWith("FailedToSendKAIA");

      await network.provider.request({
        method: "hardhat_stopImpersonatingAccount",
        params: [maliciousReceiver.address],
      });
    });
  });

  describe("Multiple DEX Support", function () {
    it("should add tokens with different factories", async function () {
      const {
        gaslessRouter,
        tokenA,
        factoryA,
        routerA,
        tokenB,
        factoryB,
        routerB,
      } = await loadFixture(gaslessSwapRouterMultiTokenFixture);
      // Add tokenA with factoryA
      await gaslessRouter.addToken(
        tokenA.address,
        factoryA.address,
        routerA.address
      );

      // Add tokenB with factoryB
      await gaslessRouter.addToken(
        tokenB.address,
        factoryB.address,
        routerB.address
      );

      // Verify both tokens are supported
      expect(await gaslessRouter.isTokenSupported(tokenA.address)).to.be.true;
      expect(await gaslessRouter.isTokenSupported(tokenB.address)).to.be.true;

      // Verify correct factory addresses are stored
      expect(await gaslessRouter.dexAddress(tokenA.address)).to.equal(
        factoryA.address
      );
      expect(await gaslessRouter.dexAddress(tokenB.address)).to.equal(
        factoryB.address
      );
    });

    it("should execute swaps through different factories", async function () {
      const {
        gaslessRouter,
        tokenA,
        factoryA,
        routerA,
        tokenB,
        factoryB,
        routerB,
        testUser,
        wkaia,
      } = await loadFixture(gaslessSwapRouterMultiTokenFixture);
      // Add tokens with their respective factories
      await gaslessRouter.addToken(
        tokenA.address,
        factoryA.address,
        routerA.address
      );
      await gaslessRouter.addToken(
        tokenB.address,
        factoryB.address,
        routerB.address
      );

      // Setup for swaps
      const swapAmount = parseEther("1.0");
      await tokenA.connect(testUser).approve(gaslessRouter.address, swapAmount);
      await tokenB.connect(testUser).approve(gaslessRouter.address, swapAmount);

      // Calculate gas repayment
      const feeData = await ethers.provider.getFeeData();
      const gasPriceBN = feeData.gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      // Calculate minAmountOut for both tokens
      const pathA = [tokenA.address, wkaia.address];
      const [, expectedOutputA] = await routerA.getAmountsOut(
        swapAmount,
        pathA
      );
      const minAmountOutA = amountRepay.add(expectedOutputA.mul(1).div(100)); // Add 1% margin

      const pathB = [tokenB.address, wkaia.address];
      const [, expectedOutputB] = await routerB.getAmountsOut(
        swapAmount,
        pathB
      );
      const minAmountOutB = amountRepay.add(expectedOutputB.mul(1).div(100)); // Add 1% margin

      // Execute swap for tokenA through factoryA
      const initialBalanceA = await testUser.getBalance();
      await gaslessRouter
        .connect(testUser)
        .swapForGas(tokenA.address, swapAmount, minAmountOutA, amountRepay);
      const finalBalanceA = await testUser.getBalance();

      // Verify swap for tokenA was successful
      expect(finalBalanceA).to.be.gt(initialBalanceA);

      // Execute swap for tokenB through factoryB
      const initialBalanceB = await testUser.getBalance();
      await gaslessRouter
        .connect(testUser)
        .swapForGas(tokenB.address, swapAmount, minAmountOutB, amountRepay);
      const finalBalanceB = await testUser.getBalance();

      // Verify swap for tokenB was successful
      expect(finalBalanceB).to.be.gt(initialBalanceB);
    });

    it("should use the correct factory and router for each token", async function () {
      const {
        factoryA,
        factoryB,
        routerA,
        routerB,
        testUser,
        tokenA,
        tokenB,
        wkaia,
      } = await loadFixture(gaslessSwapRouterMultiTokenFixture);
      // Deploy a GaslessSwapRouter with a custom router
      const GaslessSwapRouter = await ethers.getContractFactory(
        "GaslessSwapRouter"
      );
      const customGaslessRouter = await GaslessSwapRouter.deploy(wkaia.address);
      await customGaslessRouter.deployed();

      // Add tokens with specific factory addresses
      await customGaslessRouter.addToken(
        tokenA.address,
        factoryA.address,
        routerA.address
      );
      await customGaslessRouter.addToken(
        tokenB.address,
        factoryB.address,
        routerB.address
      );

      // Setup for swaps
      const swapAmount = parseEther("1.0");
      await tokenA
        .connect(testUser)
        .approve(customGaslessRouter.address, swapAmount);
      await tokenB
        .connect(testUser)
        .approve(customGaslessRouter.address, swapAmount);

      // Get minAmountOut values
      const pathA = [tokenA.address, wkaia.address];
      const [, expectedOutputA] = await routerA.getAmountsOut(
        swapAmount,
        pathA
      );

      const pathB = [tokenB.address, wkaia.address];
      const [, expectedOutputB] = await routerB.getAmountsOut(
        swapAmount,
        pathB
      );

      // Calculate gas repayment
      const feeData = await ethers.provider.getFeeData();
      const gasPriceBN = feeData.gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);

      const minAmountOutA = amountRepay.add(expectedOutputA.mul(1).div(100));
      const minAmountOutB = amountRepay.add(expectedOutputB.mul(1).div(100));

      // Both swaps should work since the router can find paths through
      // the correct factories based on the stored factory addresses
      await customGaslessRouter
        .connect(testUser)
        .swapForGas(tokenA.address, swapAmount, minAmountOutA, amountRepay);
      await customGaslessRouter
        .connect(testUser)
        .swapForGas(tokenB.address, swapAmount, minAmountOutB, amountRepay);

      // Additional check: getAmountIn should use the correct factory for each token
      const amountOutA = parseEther("0.1");
      const amountInA = await customGaslessRouter.getAmountIn(
        tokenA.address,
        amountOutA
      );

      const amountOutB = parseEther("0.1");
      const amountInB = await customGaslessRouter.getAmountIn(
        tokenB.address,
        amountOutB
      );

      // Due to different factory configurations, the amount in may differ
      expect(amountInA).to.not.equal(0);
      expect(amountInB).to.not.equal(0);
    });
  });
});

describe("GaslessSwapRouter: Error Handling & Mock Contract Tests", function () {
  describe("Error Cases", function () {
    it("should fail if token approval fails", async function () {
      const { gaslessRouter, testUser, uniswapFactory, uniswapRouter } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);

      // Deploy a malicious token that always fails on approve
      const MaliciousToken = await ethers.getContractFactory("MaliciousToken");
      const maliciousToken = await MaliciousToken.deploy(testUser.address);
      await gaslessRouter.addToken(
        maliciousToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            maliciousToken.address,
            parseEther("1"),
            parseEther("1"),
            parseEther("0.1")
          )
      ).to.be.reverted;
    });

    it("should fail if transfer to user fails", async function () {
      const { gaslessRouter, testToken, testUser, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      // Deploy a malicious contract that rejects KAIA transfers
      const MaliciousReceiver = await ethers.getContractFactory(
        "MaliciousReceiver"
      );
      const maliciousReceiver = await MaliciousReceiver.deploy();

      // Impersonate the malicious contract
      await network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [maliciousReceiver.address],
      });

      // Fund the malicious account with ETH for gas fees
      await network.provider.send("hardhat_setBalance", [
        maliciousReceiver.address,
        ethers.utils.hexValue(parseEther("10.0")),
      ]);

      const maliciousSigner = await ethers.getSigner(maliciousReceiver.address);

      // Fund the malicious contract with some test tokens
      await testToken
        .connect(testUser)
        .transfer(maliciousReceiver.address, parseEther("10"));
      await testToken
        .connect(maliciousSigner)
        .approve(gaslessRouter.address, parseEther("1"));

      // Get expected output for setting appropriate minAmountOut
      const path = [testToken.address, wkaia.address];
      await uniswapRouter.getAmountsOut(parseEther("1"), path);
      const amountRepay = parseEther("0.1");

      await expect(
        gaslessRouter
          .connect(maliciousSigner)
          .swapForGas(
            testToken.address,
            parseEther("1"),
            amountRepay.add(parseEther("0.1")),
            amountRepay
          )
      ).to.be.revertedWith("FailedToSendKAIA");

      await network.provider.request({
        method: "hardhat_stopImpersonatingAccount",
        params: [maliciousReceiver.address],
      });
    });

    it("should fail if WKAIA withdraw fails", async function () {
      const { testToken, testUser, uniswapFactory, uniswapRouter, wkaia } =
        await loadFixture(gaslessSwapRouterAddTokenFixture);
      const MaliciousWKAIA = await ethers.getContractFactory("MaliciousWKAIA");
      const maliciousWKAIA = await MaliciousWKAIA.deploy();

      const GaslessSwapRouter = await ethers.getContractFactory(
        "GaslessSwapRouter"
      );
      let gaslessRouter = await GaslessSwapRouter.deploy(
        maliciousWKAIA.address
      );

      await gaslessRouter.addToken(
        testToken.address,
        uniswapFactory.address,
        uniswapRouter.address
      );

      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      // Get expected output for setting appropriate minAmountOut
      const path = [testToken.address, wkaia.address];
      const [, expectedOutput] = await uniswapRouter.getAmountsOut(
        swapAmount,
        path
      );
      const margin = expectedOutput.mul(1).div(100); // 1% margin
      const amountRepay = parseEther("0.05");
      const minAmountOut = amountRepay.add(margin);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay)
      ).to.be.reverted;
    });

    it("should fail if coinbase payment fails", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      // Create a mock coinbase that rejects payments
      const MockCoinbase = await ethers.getContractFactory("MockCoinbase");
      const mockCoinbase = await MockCoinbase.deploy();

      // Modify block.coinbase to return mock address
      await network.provider.send("hardhat_setCoinbase", [
        mockCoinbase.address,
      ]);

      const swapAmount = parseEther("1.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      await expect(
        gaslessRouter
          .connect(testUser)
          .swapForGas(
            testToken.address,
            swapAmount,
            parseEther("0.5"),
            parseEther("0.1")
          )
      ).to.be.revertedWith("Failed to send KAIA to proposer");

      // Reset coinbase
      await network.provider.send("hardhat_setCoinbase", [
        "0x0000000000000000000000000000000000000000",
      ]);
    });

    it("should handle commission claim failure", async function () {
      const { gaslessRouter, testToken, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      await gaslessRouter.updateCommissionRate(1000); // 10%
      const swapAmount = parseEther("10.0");
      await testToken
        .connect(testUser)
        .approve(gaslessRouter.address, swapAmount);

      const gasPriceBN = (await ethers.provider.getFeeData()).gasPrice!;
      const amountRepay = gasPriceBN.mul(600000);
      const minAmountOut = amountRepay.add(parseEther("1"));

      await gaslessRouter
        .connect(testUser)
        .swapForGas(testToken.address, swapAmount, minAmountOut, amountRepay);

      const MaliciousReceiver = await ethers.getContractFactory(
        "MaliciousReceiver"
      );
      const maliciousReceiver = await MaliciousReceiver.deploy();

      await gaslessRouter.transferOwnership(maliciousReceiver.address);

      await expect(
        gaslessRouter.connect(testUser).claimCommission()
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("Mock Contract Tests", function () {
    it("should safely handle MaliciousWKAIA", async function () {
      // Deploy MaliciousWKAIA
      const MaliciousWKAIAFactory = await ethers.getContractFactory(
        "MaliciousWKAIA"
      );
      const maliciousWkaia = await MaliciousWKAIAFactory.deploy();
      await maliciousWkaia.deployed();

      // Test withdraw method
      await expect(maliciousWkaia.withdraw(parseEther("1"))).to.be.revertedWith(
        "Withdrawal failed"
      );
    });

    it("should safely handle MaliciousToken", async function () {
      const { deployer, testUser } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      // Deploy MaliciousToken
      const MaliciousTokenFactory = await ethers.getContractFactory(
        "MaliciousToken"
      );
      const maliciousToken = await MaliciousTokenFactory.deploy(
        deployer.address
      );
      await maliciousToken.deployed();

      // Test transfer method (should work)
      await maliciousToken.transfer(testUser.address, parseEther("1"));
      expect(await maliciousToken.balanceOf(testUser.address)).to.equal(
        parseEther("1")
      );
    });

    it("should test commission claim failure scenarios", async function () {
      const { deployer, gaslessRouter } = await loadFixture(
        gaslessSwapRouterAddTokenFixture
      );
      // Setup a malicious contract that rejects ETH
      const MaliciousReceiver = await ethers.getContractFactory(
        "MaliciousReceiver"
      );
      const maliciousReceiver = await MaliciousReceiver.deploy();
      await maliciousReceiver.deployed();

      // Add some funds to the contract
      await deployer.sendTransaction({
        to: gaslessRouter.address,
        value: parseEther("1"),
      });

      // Transfer ownership to the malicious contract
      await gaslessRouter.transferOwnership(maliciousReceiver.address);

      // Impersonate the malicious contract
      await network.provider.request({
        method: "hardhat_impersonateAccount",
        params: [maliciousReceiver.address],
      });

      // Fund the malicious account for gas
      await network.provider.send("hardhat_setBalance", [
        maliciousReceiver.address,
        ethers.utils.hexValue(parseEther("10.0")),
      ]);

      const maliciousSigner = await ethers.getSigner(maliciousReceiver.address);

      // Try to claim commission - should fail because the malicious contract rejects ETH
      await expect(
        gaslessRouter.connect(maliciousSigner).claimCommission()
      ).to.be.revertedWith("CommissionClaimFailed");

      // Stop impersonating
      await network.provider.request({
        method: "hardhat_stopImpersonatingAccount",
        params: [maliciousReceiver.address],
      });
    });
  });
});

describe("WKAIA Contract Tests", function () {
  let wkaia: Contract;
  let deployer: SignerWithAddress;
  let user1: SignerWithAddress;
  let user2: SignerWithAddress;

  beforeEach(async function () {
    [deployer, user1, user2] = await ethers.getSigners();

    // Deploy WKAIA
    const WKAIA = await ethers.getContractFactory("WKAIA");
    wkaia = await WKAIA.deploy();
    await wkaia.deployed();
  });

  describe("Basic Functionality", function () {
    it("should have correct name, symbol, and decimals", async function () {
      expect(await wkaia.name()).to.equal("Wrapped Kaia");
      expect(await wkaia.symbol()).to.equal("WKAIA");
      expect(await wkaia.decimals()).to.equal(18);
    });

    it("should deposit KAIA via deposit() function", async function () {
      const depositAmount = parseEther("1.0");

      await wkaia.connect(user1).deposit({ value: depositAmount });

      expect(await wkaia.balanceOf(user1.address)).to.equal(depositAmount);
      expect(await ethers.provider.getBalance(wkaia.address)).to.equal(
        depositAmount
      );
    });

    it("should deposit KAIA via receive() function", async function () {
      const depositAmount = parseEther("1.0");

      await user1.sendTransaction({
        to: wkaia.address,
        value: depositAmount,
      });

      expect(await wkaia.balanceOf(user1.address)).to.equal(depositAmount);
      expect(await ethers.provider.getBalance(wkaia.address)).to.equal(
        depositAmount
      );
    });

    it("should withdraw KAIA via withdraw() function", async function () {
      const depositAmount = parseEther("1.0");
      const withdrawAmount = parseEther("0.5");

      // First deposit
      await wkaia.connect(user1).deposit({ value: depositAmount });

      const balanceBefore = await ethers.provider.getBalance(user1.address);

      // Then withdraw
      const tx = await wkaia.connect(user1).withdraw(withdrawAmount);
      const receipt = await tx.wait();
      const gasUsed = receipt.gasUsed.mul(receipt.effectiveGasPrice);

      const balanceAfter = await ethers.provider.getBalance(user1.address);

      // Check balance changes
      expect(balanceAfter.add(gasUsed).sub(balanceBefore)).to.equal(
        withdrawAmount
      );
      expect(await wkaia.balanceOf(user1.address)).to.equal(
        depositAmount.sub(withdrawAmount)
      );
      expect(await ethers.provider.getBalance(wkaia.address)).to.equal(
        depositAmount.sub(withdrawAmount)
      );
    });

    it("should fail to withdraw more than balance", async function () {
      const depositAmount = parseEther("1.0");
      const withdrawAmount = parseEther("2.0");

      await wkaia.connect(user1).deposit({ value: depositAmount });

      await expect(wkaia.connect(user1).withdraw(withdrawAmount)).to.be
        .reverted;
    });

    it("should test totalSupply function", async function () {
      // Initially should be zero
      expect(await wkaia.totalSupply()).to.equal(0);

      // After deposits, should match contract's ETH balance
      await user1.sendTransaction({
        to: wkaia.address,
        value: parseEther("3.0"),
      });
      await user2.sendTransaction({
        to: wkaia.address,
        value: parseEther("2.0"),
      });

      expect(await wkaia.totalSupply()).to.equal(parseEther("5.0"));
      expect(await wkaia.totalSupply()).to.equal(
        await ethers.provider.getBalance(wkaia.address)
      );
    });
  });

  describe("ERC20 Functionality", function () {
    beforeEach(async function () {
      // User1 deposits some WKAIA
      await wkaia.connect(user1).deposit({ value: parseEther("10.0") });
    });

    it("should transfer WKAIA tokens", async function () {
      const transferAmount = parseEther("2.0");

      await wkaia.connect(user1).transfer(user2.address, transferAmount);

      expect(await wkaia.balanceOf(user1.address)).to.equal(parseEther("8.0"));
      expect(await wkaia.balanceOf(user2.address)).to.equal(transferAmount);
    });

    it("should approve and transferFrom WKAIA tokens", async function () {
      const approveAmount = parseEther("5.0");
      const transferAmount = parseEther("3.0");

      // User1 approves User2 to spend tokens
      await wkaia.connect(user1).approve(user2.address, approveAmount);
      expect(await wkaia.allowance(user1.address, user2.address)).to.equal(
        approveAmount
      );

      // User2 transfers tokens from User1 to themselves
      await wkaia
        .connect(user2)
        .transferFrom(user1.address, user2.address, transferAmount);

      // Check balances and allowance
      expect(await wkaia.balanceOf(user1.address)).to.equal(parseEther("7.0"));
      expect(await wkaia.balanceOf(user2.address)).to.equal(transferAmount);
      expect(await wkaia.allowance(user1.address, user2.address)).to.equal(
        approveAmount.sub(transferAmount)
      );
    });

    it("should handle failed transfers due to insufficient balance", async function () {
      const transferAmount = parseEther("15.0"); // More than User1 has

      await expect(wkaia.connect(user1).transfer(user2.address, transferAmount))
        .to.be.reverted;
    });

    it("should handle failed transferFrom due to insufficient allowance", async function () {
      const approveAmount = parseEther("2.0");
      const transferAmount = parseEther("3.0"); // More than approved

      // User1 approves User2 to spend tokens
      await wkaia.connect(user1).approve(user2.address, approveAmount);

      // User2 tries to transfer more than approved
      await expect(
        wkaia
          .connect(user2)
          .transferFrom(user1.address, user2.address, transferAmount)
      ).to.be.reverted;
    });
  });

  describe("Edge Cases", function () {
    it("should handle multiple deposits and withdrawals", async function () {
      // Make multiple deposits
      await wkaia.connect(user1).deposit({ value: parseEther("1.0") });
      await wkaia.connect(user1).deposit({ value: parseEther("2.0") });
      await user1.sendTransaction({
        to: wkaia.address,
        value: parseEther("3.0"),
      });

      expect(await wkaia.balanceOf(user1.address)).to.equal(parseEther("6.0"));

      // Make multiple withdrawals
      await wkaia.connect(user1).withdraw(parseEther("1.0"));
      await wkaia.connect(user1).withdraw(parseEther("2.0"));

      expect(await wkaia.balanceOf(user1.address)).to.equal(parseEther("3.0"));
      expect(await ethers.provider.getBalance(wkaia.address)).to.equal(
        parseEther("3.0")
      );
    });

    it("should handle zero value deposits", async function () {
      const initialBalance = await wkaia.balanceOf(user1.address);

      // Zero value deposit via deposit()
      await wkaia.connect(user1).deposit({ value: 0 });
      expect(await wkaia.balanceOf(user1.address)).to.equal(initialBalance);

      // Zero value deposit via receive()
      await user1.sendTransaction({ to: wkaia.address, value: 0 });
      expect(await wkaia.balanceOf(user1.address)).to.equal(initialBalance);
    });

    it("should handle zero value transfers and approvals", async function () {
      // Zero transfer
      await wkaia.connect(user1).transfer(user2.address, 0);
      expect(await wkaia.balanceOf(user2.address)).to.equal(0);

      // Zero approval
      await wkaia.connect(user1).approve(user2.address, 0);
      expect(await wkaia.allowance(user1.address, user2.address)).to.equal(0);

      // Zero transferFrom (should work even with zero allowance)
      await wkaia.connect(user2).transferFrom(user1.address, user2.address, 0);
      expect(await wkaia.balanceOf(user2.address)).to.equal(0);
    });

    it("should handle transfers to the contract itself", async function () {
      // Deposit some WKAIA first
      await wkaia.connect(user1).deposit({ value: parseEther("5.0") });

      // Transfer to the contract itself
      await wkaia.connect(user1).transfer(wkaia.address, parseEther("1.0"));

      // Check balances
      expect(await wkaia.balanceOf(user1.address)).to.equal(parseEther("4.0"));
      expect(await wkaia.balanceOf(wkaia.address)).to.equal(parseEther("1.0"));
    });

    it("should handle unlimited approvals correctly", async function () {
      // Set unlimited approval (uint256 max value)
      const unlimitedApproval = ethers.constants.MaxUint256;
      await wkaia.connect(user1).approve(user2.address, unlimitedApproval);

      // Deposit some tokens
      await wkaia.connect(user1).deposit({ value: parseEther("10.0") });

      // Transfer multiple times without approval decreasing
      await wkaia
        .connect(user2)
        .transferFrom(user1.address, user2.address, parseEther("2.0"));
      expect(await wkaia.allowance(user1.address, user2.address)).to.equal(
        unlimitedApproval
      );

      await wkaia
        .connect(user2)
        .transferFrom(user1.address, user2.address, parseEther("3.0"));
      expect(await wkaia.allowance(user1.address, user2.address)).to.equal(
        unlimitedApproval
      );
    });
  });

  describe("Events", function () {
    it("should check that approve emits an event", async function () {
      // Test that approve emits correct event
      await expect(
        wkaia.connect(user1).approve(user2.address, parseEther("1.0"))
      )
        .to.emit(wkaia, "Approval")
        .withArgs(user1.address, user2.address, parseEther("1.0"));
    });

    it("should check that deposit emits an event", async function () {
      // Test that deposit emits correct event
      await expect(wkaia.connect(user1).deposit({ value: parseEther("1.0") }))
        .to.emit(wkaia, "Deposit")
        .withArgs(user1.address, parseEther("1.0"));
    });

    it("should check that withdraw emits an event", async function () {
      // First deposit some funds
      await wkaia.connect(user1).deposit({ value: parseEther("5.0") });

      // Test that withdraw emits correct event
      await expect(wkaia.connect(user1).withdraw(parseEther("1.0")))
        .to.emit(wkaia, "Withdrawal")
        .withArgs(user1.address, parseEther("1.0"));
    });

    it("should check that transfer emits an event", async function () {
      // First deposit some funds
      await wkaia.connect(user1).deposit({ value: parseEther("5.0") });

      // Test that transfer emits correct event
      await expect(
        wkaia.connect(user1).transfer(user2.address, parseEther("1.0"))
      )
        .to.emit(wkaia, "Transfer")
        .withArgs(user1.address, user2.address, parseEther("1.0"));
    });
  });
});
