import { expect } from "chai";
import { auctionTestFixture } from "../materials";
import hre from "hardhat";
import { getMiner } from "../common/helper";

describe("AuctionFeeVault", () => {
  it("should initialize correctly", async () => {
    const { auctionFeeVault, deployer } = await auctionTestFixture();

    expect(await auctionFeeVault.owner()).to.equal(deployer.address);
    expect(await auctionFeeVault.searcherPaybackRate()).to.equal(0);
    expect(await auctionFeeVault.validatorPaybackRate()).to.equal(1000);
  });
  it("only owner", async () => {
    const { auctionFeeVault, user1 } = await auctionTestFixture();

    await expect(
      auctionFeeVault.connect(user1).setSearcherPaybackRate(100)
    ).to.be.revertedWithCustomError(
      auctionFeeVault,
      "OwnableUnauthorizedAccount"
    );

    await expect(
      auctionFeeVault.connect(user1).setValidatorPaybackRate(100)
    ).to.be.revertedWithCustomError(
      auctionFeeVault,
      "OwnableUnauthorizedAccount"
    );

    await expect(
      auctionFeeVault
        .connect(user1)
        .registerRewardAddresses([user1.address], [user1.address])
    ).to.be.revertedWithCustomError(
      auctionFeeVault,
      "OwnableUnauthorizedAccount"
    );
  });
  it("#setPaybackPercentage: invalid input", async () => {
    const { auctionFeeVault, deployer } = await auctionTestFixture();

    // 10% payback for validator is already set in the constructor
    await expect(
      auctionFeeVault.connect(deployer).setSearcherPaybackRate(9001)
    ).to.be.revertedWithCustomError(auctionFeeVault, "InvalidInput");
  });
  it("#setValidatorPaybackPercentage: invalid input", async () => {
    const { auctionFeeVault, deployer } = await auctionTestFixture();

    await expect(
      auctionFeeVault.connect(deployer).setValidatorPaybackRate(10001)
    ).to.be.revertedWithCustomError(auctionFeeVault, "InvalidInput");
  });
  it("#registerRewardAddresses: invalid input", async () => {
    const { auctionFeeVault, deployer } = await auctionTestFixture();

    await expect(
      auctionFeeVault
        .connect(deployer)
        .registerRewardAddresses(
          [deployer.address],
          [deployer.address, deployer.address]
        )
    ).to.be.revertedWithCustomError(auctionFeeVault, "InvalidInput");
  });
  it("registerRewardAddresses: success", async () => {
    const { auctionFeeVault, deployer, user1, user2 } =
      await auctionTestFixture();

    const nodeIds = [deployer.address, user1.address, user2.address];
    const rewardAddrs = [user1.address, user2.address, deployer.address];

    await auctionFeeVault
      .connect(deployer)
      .registerRewardAddresses(nodeIds, rewardAddrs);

    for (let i = 0; i < nodeIds.length; i++) {
      expect(await auctionFeeVault.getRewardAddr(nodeIds[i])).to.equal(
        rewardAddrs[i]
      );
    }
  });
  it("#registerRewardAddress: only valid validator", async () => {
    const { auctionFeeVault, deployer, user1, addressBook, cnStaking } =
      await auctionTestFixture();

    // not registered in address book
    await expect(
      auctionFeeVault
        .connect(deployer)
        .registerRewardAddress(deployer.address, deployer.address)
    ).to.be.revertedWithoutReason();

    await addressBook.registerCnStakingContract(
      deployer.address,
      cnStaking.address,
      deployer.address
    );

    // user1 is not admin of cnStaking
    await expect(
      auctionFeeVault
        .connect(user1)
        .registerRewardAddress(deployer.address, deployer.address)
    ).to.be.revertedWithCustomError(auctionFeeVault, "OnlyStakingAdmin");

    // invalid CN node ID
    await expect(
      auctionFeeVault
        .connect(deployer)
        .registerRewardAddress(user1.address, deployer.address)
    ).to.be.revertedWith("Invalid CN node ID.");

    // success
    await expect(
      auctionFeeVault
        .connect(deployer)
        .registerRewardAddress(deployer.address, deployer.address)
    )
      .to.emit(auctionFeeVault, "RewardAddressRegistered")
      .withArgs(deployer.address, deployer.address);
  });
  it("#takeBid: should receive KAIA", async () => {
    const { auctionFeeVault, deployer } = await auctionTestFixture();

    await expect(
      auctionFeeVault
        .connect(deployer)
        .takeBid(deployer.address, { value: 100 })
    )
      .to.emit(auctionFeeVault, "FeeDeposit")
      .withArgs(await getMiner(), 100, 0, 0);

    expect(await auctionFeeVault.accumulatedBids()).to.equal(100);

    await auctionFeeVault.connect(deployer).setSearcherPaybackRate(1000);
    await auctionFeeVault.connect(deployer).setValidatorPaybackRate(1000);

    await auctionFeeVault
      .connect(deployer)
      .takeBid(deployer.address, { value: 100 });

    // accumulatedBids is not affected by paybackPercentage and validatorPaybackPercentage
    expect(await auctionFeeVault.accumulatedBids()).to.equal(200);
  });
  it("#takeBid: do not ceil the payback amount", async () => {
    const { auctionFeeVault, deployer, user1, user2, addressBook, cnStaking } =
      await auctionTestFixture();

    // 90% payback to searcher, and 10% payback to validator
    await auctionFeeVault.connect(deployer).setSearcherPaybackRate(9000);

    const miner = await getMiner();
    await addressBook.registerCnStakingContract(
      miner,
      cnStaking.address,
      deployer.address
    );

    await auctionFeeVault
      .connect(deployer)
      .registerRewardAddress(miner, user2.address);

    const beforeBalance = await hre.ethers.provider.getBalance(user1.address);
    const beforeRewardBalance = await hre.ethers.provider.getBalance(
      user2.address
    );

    await expect(
      auctionFeeVault.connect(deployer).takeBid(user1.address, { value: 1245n })
    )
      .to.emit(auctionFeeVault, "FeeDeposit")
      .withArgs(await getMiner(), 1245n, 1120n, 124n);

    expect(await hre.ethers.provider.getBalance(user1.address)).to.equal(
      beforeBalance.add(1120)
    );
    expect(await hre.ethers.provider.getBalance(user2.address)).to.equal(
      beforeRewardBalance.add(124)
    );
  });
  it("#takeBid: should payback to searcher", async () => {
    const { auctionFeeVault, deployer, user1 } = await auctionTestFixture();

    // 10% payback to searcher
    await auctionFeeVault.connect(deployer).setSearcherPaybackRate(1000);

    const beforeBalance = await hre.ethers.provider.getBalance(user1.address);

    await expect(
      auctionFeeVault.connect(deployer).takeBid(user1.address, { value: 100 })
    )
      .to.emit(auctionFeeVault, "FeeDeposit")
      .withArgs(await getMiner(), 100, 10, 0);

    expect(await hre.ethers.provider.getBalance(user1.address)).to.equal(
      beforeBalance.add(10)
    );
  });
  it("#takeBid: should payback to validator", async () => {
    const { auctionFeeVault, deployer, addressBook, cnStaking, user1 } =
      await auctionTestFixture();

    const miner = await getMiner();
    await addressBook.registerCnStakingContract(
      miner,
      cnStaking.address,
      deployer.address
    );

    await auctionFeeVault
      .connect(deployer)
      .registerRewardAddress(miner, user1.address);

    const beforeBalance = await hre.ethers.provider.getBalance(user1.address);

    await expect(
      auctionFeeVault.connect(deployer).takeBid(miner, { value: 100 })
    )
      .to.emit(auctionFeeVault, "FeeDeposit")
      .withArgs(miner, 100, 0, 10);

    expect(await hre.ethers.provider.getBalance(user1.address)).to.equal(
      beforeBalance.add(10)
    );
  });
  it("#takeBid: should payback to validator by registerRewardAddresses", async () => {
    const { auctionFeeVault, deployer, addressBook, cnStaking, user1 } =
      await auctionTestFixture();

    const miner = await getMiner();

    await auctionFeeVault
      .connect(deployer)
      .registerRewardAddresses([miner], [user1.address]);

    const beforeBalance = await hre.ethers.provider.getBalance(user1.address);

    await expect(
      auctionFeeVault.connect(deployer).takeBid(miner, { value: 100 })
    )
      .to.emit(auctionFeeVault, "FeeDeposit")
      .withArgs(miner, 100, 0, 10);

    expect(await hre.ethers.provider.getBalance(user1.address)).to.equal(
      beforeBalance.add(10)
    );
  });
  it("#withdraw: failed to send KAIA", async () => {
    const { auctionFeeVault, deployer, auctionEntryPoint } =
      await auctionTestFixture();

    await auctionFeeVault
      .connect(deployer)
      .takeBid(deployer.address, { value: 100 });

    // AuctionEntryPoint has no receive function, so it will revert
    await expect(
      auctionFeeVault.connect(deployer).withdraw(auctionEntryPoint.address)
    ).to.be.revertedWithCustomError(auctionFeeVault, "WithdrawalFailed");
  });
  it("#withdraw: success", async () => {
    const { auctionFeeVault, deployer, user1 } = await auctionTestFixture();

    await auctionFeeVault
      .connect(deployer)
      .takeBid(deployer.address, { value: 100 });

    const beforeBalance1 = await hre.ethers.provider.getBalance(user1.address);

    await expect(auctionFeeVault.connect(deployer).withdraw(user1.address))
      .to.emit(auctionFeeVault, "FeeWithdrawal")
      .withArgs(100);

    expect(await hre.ethers.provider.getBalance(user1.address)).to.equal(
      beforeBalance1.add(100)
    );
  });
});
