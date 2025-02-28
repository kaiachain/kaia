import { ethers } from "hardhat";
import { parseEther } from "ethers/lib/utils";

import factoryArtifact from "@uniswap/v2-core/build/UniswapV2Factory.json";
import routerArtifact from "@uniswap/v2-periphery/build/UniswapV2Router02.json";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";

export async function gaslessSwapRouterFixture() {
  const [deployer, testUser, thirdUser] = await ethers.getSigners();
  const INITIAL_LIQUIDITY = parseEther("1000");

  // Deploy TestToken
  const TestToken = await ethers.getContractFactory("TestToken");
  const testToken = await TestToken.deploy(testUser.address);
  await testToken.deployed();

  // Deploy WKAIA
  const WKAIA = await ethers.getContractFactory("WKAIA");
  const wkaia = await WKAIA.deploy();
  await wkaia.deployed();

  // Deploy UniswapV2Factory
  const Factory = new ethers.ContractFactory(factoryArtifact.abi, factoryArtifact.bytecode, deployer);
  const uniswapFactory = await Factory.deploy(deployer.address);
  await uniswapFactory.deployed();

  // Deploy UniswapV2Router02
  const Router = new ethers.ContractFactory(routerArtifact.abi, routerArtifact.bytecode, deployer);
  const uniswapRouter = await Router.deploy(uniswapFactory.address, wkaia.address);
  await uniswapRouter.deployed();

  // Deploy GaslessSwapRouter
  const GaslessSwapRouter = await ethers.getContractFactory("GaslessSwapRouter");
  const gaslessRouter = await GaslessSwapRouter.deploy(wkaia.address);
  await gaslessRouter.deployed();

  // Create Uniswap pair and add liquidity
  await uniswapFactory.createPair(testToken.address, wkaia.address);

  // Setup liquidity for testing
  await wkaia.connect(testUser).deposit({ value: INITIAL_LIQUIDITY });
  await testToken.connect(testUser).approve(uniswapRouter.address, INITIAL_LIQUIDITY);
  await wkaia.connect(testUser).approve(uniswapRouter.address, INITIAL_LIQUIDITY);
  const currentBlock = await ethers.provider.getBlock("latest");
  const deadline = currentBlock.timestamp + 3600;
  await uniswapRouter
    .connect(testUser)
    .addLiquidity(
      testToken.address,
      wkaia.address,
      INITIAL_LIQUIDITY,
      INITIAL_LIQUIDITY,
      0,
      0,
      testUser.address,
      deadline,
    );

  return {
    INITIAL_LIQUIDITY,
    deployer,
    gaslessRouter,
    testToken,
    testUser,
    thirdUser,
    uniswapFactory,
    uniswapRouter,
    wkaia,
  };
}

export async function gaslessSwapRouterAddTokenFixture() {
  const {
    INITIAL_LIQUIDITY,
    deployer,
    gaslessRouter,
    testToken,
    testUser,
    thirdUser,
    uniswapFactory,
    uniswapRouter,
    wkaia,
  } = await loadFixture(gaslessSwapRouterFixture);

  await gaslessRouter.addToken(testToken.address, uniswapFactory.address, uniswapRouter.address);

  return {
    INITIAL_LIQUIDITY,
    deployer,
    gaslessRouter,
    testToken,
    testUser,
    thirdUser,
    uniswapFactory,
    uniswapRouter,
    wkaia,
  };
}

export async function gaslessSwapRouterMultiTokenFixture() {
  const liquidityAmount = parseEther("1000");
  const [deployer, testUser] = await ethers.getSigners();

  // Deploy WKAIA
  const WKAIA = await ethers.getContractFactory("WKAIA");
  const wkaia = await WKAIA.deploy();
  await wkaia.deployed();

  // Deploy two separate factories
  const Factory = new ethers.ContractFactory(factoryArtifact.abi, factoryArtifact.bytecode, deployer);
  const factoryA = await Factory.deploy(deployer.address);
  const factoryB = await Factory.deploy(deployer.address);
  await factoryA.deployed();
  await factoryB.deployed();

  // Deploy two separate routers
  const Router = new ethers.ContractFactory(routerArtifact.abi, routerArtifact.bytecode, deployer);
  const routerA = await Router.deploy(factoryA.address, wkaia.address);
  const routerB = await Router.deploy(factoryB.address, wkaia.address);
  await routerA.deployed();
  await routerB.deployed();

  // Deploy two test tokens
  const TestToken = await ethers.getContractFactory("TestToken");
  const tokenA = await TestToken.deploy(testUser.address);
  const tokenB = await TestToken.deploy(testUser.address);
  await tokenA.deployed();
  await tokenB.deployed();

  // Create pairs in both factories
  await factoryA.createPair(tokenA.address, wkaia.address);
  await factoryB.createPair(tokenB.address, wkaia.address);

  // Setup liquidity for both pairs
  await wkaia.connect(testUser).deposit({ value: liquidityAmount.mul(2) });
  await tokenA.connect(testUser).approve(routerA.address, liquidityAmount);
  await tokenB.connect(testUser).approve(routerB.address, liquidityAmount);
  await wkaia.connect(testUser).approve(routerA.address, liquidityAmount);
  await wkaia.connect(testUser).approve(routerB.address, liquidityAmount);

  const currentBlock = await ethers.provider.getBlock("latest");
  const deadline = currentBlock.timestamp + 3600;

  // Add liquidity to pair in Factory A
  await routerA
    .connect(testUser)
    .addLiquidity(tokenA.address, wkaia.address, liquidityAmount, liquidityAmount, 0, 0, testUser.address, deadline);

  // Add liquidity to pair in Factory B
  await routerB
    .connect(testUser)
    .addLiquidity(tokenB.address, wkaia.address, liquidityAmount, liquidityAmount, 0, 0, testUser.address, deadline);

  // Deploy GaslessSwapRouter
  const GaslessSwapRouter = await ethers.getContractFactory("GaslessSwapRouter");
  const gaslessRouter = await GaslessSwapRouter.deploy(wkaia.address);
  await gaslessRouter.deployed();

  return {
    deployer,
    factoryA,
    factoryB,
    gaslessRouter,
    liquidityAmount,
    routerA,
    routerB,
    testUser,
    tokenA,
    tokenB,
    wkaia,
  };
}
