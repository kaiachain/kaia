const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("TreasuryRebalance", function() {
  let TreasuryRebalance,
    treasuryRebalance,
    kgf,
    kir,
    test,
    owner,
    retired1,
    retired2,
    newbie1,
    newbie2;
  let currentBlock, executionBlock, memo;

  beforeEach(async function() {
    [owner, account1, newbie1, newbie2] = await ethers.getSigners();
    currentBlock = await ethers.provider.getBlockNumber();
    executionBlock = currentBlock + 40;

    // Deploy dependancy treasuryRebalances to register as retireds
    // retired1 and retired2 simulates KGF and KIR treasuryRebalances respectively
    const KGF = await ethers.getContractFactory("SenderTest1");
    kgf = await KGF.deploy();
    retired1 = kgf.address;

    const KIR = await ethers.getContractFactory("SenderTest2");
    kir = await KIR.deploy();
    retired2 = kir.address;

    const TEST = await ethers.getContractFactory(
      "contracts/libs/Ownable.sol:Ownable"
    );
    test = await TEST.deploy();

    // Send some funds to retired1 to simulate KFG funds
    await owner.sendTransaction({
      to: retired1,
      value: hre.ethers.utils.parseEther("20"),
    });

    // Send some funds to retired2 to simulate KIR funds
    await owner.sendTransaction({
      to: retired2,
      value: hre.ethers.utils.parseEther("20"),
    });

    // Deploy Treasury Rebalance treasuryRebalance
    treasuryFund = hre.ethers.utils.parseEther("30");
    TreasuryRebalance = await ethers.getContractFactory("TreasuryRebalance");
    treasuryRebalance = await TreasuryRebalance.deploy(executionBlock);

    // memo format
    memo =
      '{ "retirees": [ { "retired": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "balance": 20000000000000000000 }, { "retired": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "balance": 20000000000000000000 } ], "newbies": [ { "newbie": "0x38138d89c321b3b5f421e9452b69cf29e4380bae", "fundAllocated": 10000000000000000000 }, { "newbie": "0x0a33a1b99bd67a7189573dd74de80293afdf969a", "fundAllocated": 10000000000000000000 } ], "burnt": 7.2e+37, "success": true }';
  });

  describe("Deployment", function() {
    it("Should check the correct initial values for dependancy treasuryRebalances", async function() {
      const retired1Balance = await ethers.provider.getBalance(retired1);
      const retired2Balance = await ethers.provider.getBalance(retired1);
      expect(await retired1Balance).to.equal(hre.ethers.utils.parseEther("20"));
      expect(await retired2Balance).to.equal(hre.ethers.utils.parseEther("20"));

      const [adminList] = await kgf.getState();
      expect(adminList[0]).to.equal(owner.address);
    });

    it("Should set the correct initial values for main treasuryRebalance", async function() {
      expect(await treasuryRebalance.status()).to.equal(0);
      expect(await treasuryRebalance.getTreasuryAmount()).to.equal(0);
      expect(await treasuryRebalance.rebalanceBlockNumber()).to.be.greaterThan(
        currentBlock
      );
    });
  });

  describe("registerRetired()", function() {
    it("Should add a retired", async function() {
      await treasuryRebalance.registerRetired(retired1);
      const retired = await treasuryRebalance.getRetired(retired1);
      expect(retired[0]).to.equal(retired1);
      expect(retired[1].length).to.equal(0);
    });

    it("Should emit a RegisterRetired event", async function() {
      await expect(treasuryRebalance.registerRetired(retired1))
        .to.emit(treasuryRebalance, "RetiredRegistered")
        .withArgs(retired1);
    });

    it("Should not allow adding the same retired twice", async function() {
      expect(await treasuryRebalance.retiredExists(retired1)).to.equal(false);
      await treasuryRebalance.registerRetired(retired1);
      await expect(
        treasuryRebalance.registerRetired(retired1)
      ).to.be.revertedWith("Retired address is already registered");
    });

    it("Should not allow non-owner to add a retired", async function() {
      await expect(
        treasuryRebalance.connect(account1).registerRetired(retired2)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should not register retired when contract is not in Initialized state", async function() {
      await treasuryRebalance.finalizeRegistration();
      await expect(
        treasuryRebalance.registerRetired(retired1)
      ).to.be.revertedWith("Not in the designated status");
    });

    it("retireds length should be one", async function() {
      await treasuryRebalance.registerRetired(retired1);
      const length = await treasuryRebalance.getRetiredCount();
      expect(length).to.equal(1);
    });

    it("Should revert if the retired address is zero", async function() {
      await expect(
        treasuryRebalance.registerRetired(ethers.constants.AddressZero)
      ).to.be.revertedWith("Invalid address");
    });
  });

  describe("removeRetired()", function() {
    beforeEach(async function() {
      await treasuryRebalance.registerRetired(retired1);
    });

    it("Should remove a retired", async function() {
      await treasuryRebalance.removeRetired(retired1);
      await expect(treasuryRebalance.getRetired(retired1)).to.be.revertedWith(
        "Retired not registered"
      );
    });

    it("Should emit a RemoveRetired event", async function() {
      await expect(treasuryRebalance.removeRetired(retired1))
        .to.emit(treasuryRebalance, "RetiredRemoved")
        .withArgs(retired1);
    });

    it("Should not allow removing a non-existent retired", async function() {
      await expect(treasuryRebalance.removeRetired(owner.address)).to.be
        .reverted;
    });

    it("Should not allow non-owner to remove a retired", async function() {
      await expect(
        treasuryRebalance.connect(account1).removeRetired(retired1)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should not remove retired when contract is not in Initialized state", async function() {
      await treasuryRebalance.finalizeRegistration();
      await expect(
        treasuryRebalance.removeRetired(retired1)
      ).to.be.revertedWith("Not in the designated status");
    });
  });

  describe("registerNewbie", function() {
    let newbieAddress;

    beforeEach(async function() {
      newbieAddress = newbie1.address;
    });

    it("Should register newbie address and its fund distribution", async function() {
      const amount = hre.ethers.utils.parseEther("20");

      await treasuryRebalance.registerNewbie(newbieAddress, amount);
      const newbie = await treasuryRebalance.getNewbie(newbieAddress);
      expect(newbie[0]).to.equal(newbieAddress);
      expect(newbie[1]).to.equal(amount);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(1);

      const treasuryAmount = await treasuryRebalance.getTreasuryAmount();
      expect(treasuryAmount).to.equal(amount);
    });

    it("Should emit a RegisterNewbie event", async function() {
      const amount = hre.ethers.utils.parseEther("20");
      await expect(treasuryRebalance.registerNewbie(newbieAddress, amount))
        .to.emit(treasuryRebalance, "NewbieRegistered")
        .withArgs(newbieAddress, amount);
    });

    it("Should revert if register newbie twice", async function() {
      const amount1 = hre.ethers.utils.parseEther("20");
      await treasuryRebalance.registerNewbie(newbie1.address, amount1);
      await treasuryRebalance.registerNewbie(newbie2.address, amount1);
      await expect(
        treasuryRebalance.registerNewbie(newbie1.address, amount1)
      ).to.be.revertedWith("Newbie address is already registered");
    });

    it("Should revert if newbie when contract is not in Initialized state", async function() {
      const amount1 = hre.ethers.utils.parseEther("20");
      await treasuryRebalance.finalizeRegistration();
      await expect(
        treasuryRebalance.registerNewbie(newbieAddress, amount1)
      ).to.be.revertedWith("Not in the designated status");
    });

    it("Should not allow non-owner to add a newbie", async function() {
      const amount1 = hre.ethers.utils.parseEther("20");
      await expect(
        treasuryRebalance
          .connect(account1)
          .registerNewbie(newbieAddress, amount1)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should revert if the newbie address is zero", async function() {
      const amount1 = hre.ethers.utils.parseEther("20");
      await expect(
        treasuryRebalance.registerNewbie(ethers.constants.AddressZero, amount1)
      ).to.be.revertedWith("Invalid address");
    });

    it("Should revert if the amount is set to 0", async function() {
      await expect(
        treasuryRebalance.registerNewbie(newbieAddress, 0)
      ).to.be.revertedWith("Amount cannot be set to 0");
    });
  });

  describe("removeNewbie", function() {
    let newbieAddress;
    let amount;

    beforeEach(async function() {
      newbieAddress = newbie1.address;
      amount = hre.ethers.utils.parseEther("20");
      await treasuryRebalance.registerNewbie(newbieAddress, amount);
    });

    it("Should remove newbie", async function() {
      await treasuryRebalance.removeNewbie(newbieAddress);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(0);
      expect(await treasuryRebalance.getTreasuryAmount()).to.equal(0);
      await expect(
        treasuryRebalance.getNewbie(newbieAddress)
      ).to.be.revertedWith("Newbie not registered");
    });

    it("Should emit RemoveNewbie event", async function() {
      await expect(treasuryRebalance.removeNewbie(newbieAddress))
        .to.emit(treasuryRebalance, "NewbieRemoved")
        .withArgs(newbieAddress);
    });

    it("Should not remove unregistered newbie", async function() {
      await expect(treasuryRebalance.removeNewbie(newbie2.address)).to.be
        .reverted;
    });

    it("Should not allow non-owner to remove a newbie", async function() {
      await expect(
        treasuryRebalance.connect(account1).removeNewbie(newbieAddress)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should not remove newbie when contract is not in Initialized state", async function() {
      await treasuryRebalance.finalizeRegistration();
      await expect(
        treasuryRebalance.removeNewbie(newbieAddress)
      ).to.be.revertedWith("Not in the designated status");
    });
  });

  describe("approve", function() {
    beforeEach(async function() {
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerRetired(retired2);
      await treasuryRebalance.registerRetired(owner.address);
      await treasuryRebalance.registerRetired(newbie1.address);
      await treasuryRebalance.registerRetired(test.address);
      await treasuryRebalance.finalizeRegistration();
    });

    it("Should approve retired if msg.sender is admin of retired contract", async function() {
      const retired = await treasuryRebalance.getRetired(retired1);
      expect(retired[1].length).to.equal(0);
      const tx = await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);

      const updatedRetiredDetails = await treasuryRebalance.getRetired(
        retired1
      );
      expect(updatedRetiredDetails[1][0]).to.equal(owner.address);

      await expect(tx)
        .to.emit(treasuryRebalance, "Approved")
        .withArgs(retired1, owner.address, 1);
    });

    it("Should approve retiredAddress is the msg.sender if retired is a EOA", async function() {
      await treasuryRebalance.approve(owner.address);
      const retired = await treasuryRebalance.getRetired(owner.address);
      expect(retired[1][0]).to.equal(owner.address);
    });

    it("Should revert if retired is already approved", async function() {
      await treasuryRebalance.approve(retired1);
      await expect(treasuryRebalance.approve(retired1)).to.be.revertedWith(
        "Already approved"
      );
    });

    it("Should revert if retired is not registered", async function() {
      // try to approve unregistered retired
      await expect(
        treasuryRebalance.approve(newbie2.address)
      ).to.be.revertedWith("retired needs to be registered before approval");
    });

    it("Should revert if retired is a EOA and if msg.sender is not the admin", async function() {
      await expect(
        treasuryRebalance.approve(newbie1.address)
      ).to.be.revertedWith("retiredAddress is not the msg.sender");
    });

    it("Should revert if retired is a contract address but does not have getState() method", async function() {
      await expect(treasuryRebalance.approve(test.address)).to.be.reverted;
    });

    it("Should revert if retired is a contract but adminList is empty", async function() {
      await kgf.emptyAdminList();
      await expect(treasuryRebalance.approve(retired1)).to.be.revertedWith(
        "admin list cannot be empty"
      );
    });

    it("Should not approve if retired is a contract but msg.sender is not the admin", async function() {
      await expect(
        treasuryRebalance.connect(account1).approve(retired1)
      ).to.be.revertedWith("msg.sender is not the admin");
    });
  });

  describe("setStatus", function() {
    let initialStatus;

    beforeEach(async function() {
      initialStatus = await treasuryRebalance.status();
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerRetired(retired2);
    });

    it("Should set status to Registered", async function() {
      expect(initialStatus).to.equal(0);
      await treasuryRebalance.finalizeRegistration();
      expect(await treasuryRebalance.status()).to.equal(1);
    });

    it("Should set status to Approved", async function() {
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
      await treasuryRebalance.finalizeApproval();
      expect(await treasuryRebalance.status()).to.equal(2);
    });

    it("Should not set status to Approved when treasury amount exceeds balance of retirees", async function() {
      const amount = hre.ethers.utils.parseEther("50");
      await treasuryRebalance.registerNewbie(newbie1.address, amount);
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
      await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
        "treasury amount should be less than the sum of all retired address balances"
      );
    });

    it("Should revert if the current status is tried to set again", async function() {
      await treasuryRebalance.finalizeRegistration();
      await expect(treasuryRebalance.finalizeRegistration()).to.be.revertedWith(
        "Not in the designated status"
      );
    });

    it("Should revert if owner tries to set Finalize after Registered", async function() {
      await treasuryRebalance.finalizeRegistration();
      await expect(treasuryRebalance.finalizeContract(memo)).to.be.revertedWith(
        "Not in the designated status"
      );
    });

    it("Should revert if owner tries to set Approved before Registered", async function() {
      await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
        "Not in the designated status"
      );
    });

    it("Should revert if owner tries to set Registered after Approved", async function() {
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
      await treasuryRebalance.finalizeApproval();
      await expect(treasuryRebalance.finalizeRegistration()).to.be.revertedWith(
        "Not in the designated status"
      );
    });

    it("should revert if not called by the owner", async () => {
      await expect(
        treasuryRebalance.connect(account1).finalizeRegistration()
      ).to.be.revertedWith("Ownable: caller is not the owner");
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
      await expect(
        treasuryRebalance.connect(account1).finalizeApproval()
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("Should emit StatusChanged event", async function() {
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
      await expect(treasuryRebalance.finalizeApproval())
        .to.emit(treasuryRebalance, "StatusChanged")
        .withArgs(2);
    });

    describe("Reach Quorom", function() {
      it("Should revert if min required admins does not approve", async function() {
        await treasuryRebalance.finalizeRegistration();
        await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );
      });

      it("Should revert if approved admin change during the contract ", async function() {
        await treasuryRebalance.finalizeRegistration();
        await treasuryRebalance.approve(retired1);
        await treasuryRebalance.approve(retired2);
        await kgf.addAdmin(account1.address);
        await kgf.changeMinReq(2);
        await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );

        await treasuryRebalance.connect(account1).approve(retired1);
        await treasuryRebalance.finalizeApproval();
      });

      it("Should revert if approved admin change during the contract ", async function() {
        await treasuryRebalance.finalizeRegistration();
        await treasuryRebalance.approve(retired1);
        await kgf.changeMinReq(2);
        await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );

        await expect(
          treasuryRebalance.connect(account1).approve(retired1),
          "msg.sender is not the admin"
        );
      });

      it("Should revert if admin list change during the contract ", async function() {
        await treasuryRebalance.finalizeRegistration();
        await treasuryRebalance.approve(retired1);
        await kgf.emptyAdminList();
        await kgf.addAdmin(account1.address);
        await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
          "min required admins should approve"
        );
      });

      it("Should revert if EOA did not approve", async function() {
        await treasuryRebalance.registerRetired(owner.address);
        await treasuryRebalance.registerRetired(account1.address);
        await treasuryRebalance.finalizeRegistration();
        await treasuryRebalance.approve(retired1);
        await treasuryRebalance.approve(retired2);
        await treasuryRebalance.approve(owner.address);
        await expect(treasuryRebalance.finalizeApproval()).to.be.revertedWith(
          "EOA should approve"
        );
      });
    });
  });

  describe("finalize contract", function() {
    beforeEach(async function() {
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerRetired(retired2);
      await treasuryRebalance.finalizeRegistration();
      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.approve(retired2);
    });

    it("should set the memo and status to Finalized", async function() {
      await treasuryRebalance.finalizeApproval();
      await hre.network.provider.send("hardhat_mine", ["0x32"]);
      await treasuryRebalance.finalizeContract(memo);
      expect(await treasuryRebalance.memo()).to.equal(memo);
      expect(await treasuryRebalance.status()).to.equal(3);
    });

    it("Should emit Finalize event", async function() {
      await treasuryRebalance.finalizeApproval();
      await hre.network.provider.send("hardhat_mine", ["0x32"]);
      await expect(treasuryRebalance.finalizeContract(memo))
        .to.emit(treasuryRebalance, "Finalized")
        .withArgs(memo, 3);
    });

    it("should revert if not called by the owner", async () => {
      await treasuryRebalance.finalizeApproval();
      await expect(
        treasuryRebalance.connect(account1).finalizeContract(memo)
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("should revert if its not called at the Approved stage", async () => {
      await expect(treasuryRebalance.finalizeContract(memo)).to.be.revertedWith(
        "Not in the designated status"
      );
    });

    it("should revert before rebalanceBlockNumber", async () => {
      await treasuryRebalance.finalizeApproval();
      await expect(treasuryRebalance.finalizeContract(memo)).to.be.revertedWith(
        "Contract can only finalize after executing rebalancing"
      );
    });
  });

  describe("reset", function() {
    beforeEach(async function() {
      await treasuryRebalance;
    });

    it("should reset all storage values to 0 at Initialize state", async function() {
      await treasuryRebalance.reset();
      expect(await treasuryRebalance.getRetiredCount()).to.equal(0);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(0);
      expect(await treasuryRebalance.getTreasuryAmount()).to.equal(0);
      expect(await treasuryRebalance.memo()).to.equal("");
      expect(await treasuryRebalance.status()).to.equal(0);
      expect(await treasuryRebalance.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("should reset all storage values to 0 at Registered state", async function() {
      const amount = hre.ethers.utils.parseEther("50");
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerRetired(retired2);

      await treasuryRebalance.registerNewbie(newbie1.address, amount);
      await treasuryRebalance.finalizeRegistration();
      expect(await treasuryRebalance.getRetiredCount()).to.equal(2);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(1);

      await treasuryRebalance.reset();
      expect(await treasuryRebalance.getRetiredCount()).to.equal(0);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(0);
      expect(await treasuryRebalance.getTreasuryAmount()).to.equal(0);
      expect(await treasuryRebalance.memo()).to.equal("");
      expect(await treasuryRebalance.status()).to.equal(0);
      expect(await treasuryRebalance.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("should reset all storage values to 0 at Approved state", async function() {
      const amount = hre.ethers.utils.parseEther("10");
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerNewbie(newbie1.address, amount);
      await treasuryRebalance.finalizeRegistration();

      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.finalizeApproval();
      expect(await treasuryRebalance.status()).to.equal(2);

      expect(await treasuryRebalance.getRetiredCount()).to.equal(1);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(1);

      await treasuryRebalance.reset();
      expect(await treasuryRebalance.getRetiredCount()).to.equal(0);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(0);
      expect(await treasuryRebalance.getTreasuryAmount()).to.equal(0);
      expect(await treasuryRebalance.memo()).to.equal("");
      expect(await treasuryRebalance.status()).to.equal(0);
      expect(await treasuryRebalance.rebalanceBlockNumber()).to.not.equal(0);
    });

    it("should revert when tried to reset after finalization", async function() {
      const amount = hre.ethers.utils.parseEther("10");
      await treasuryRebalance.registerRetired(retired1);
      await treasuryRebalance.registerNewbie(newbie1.address, amount);
      await treasuryRebalance.finalizeRegistration();

      await treasuryRebalance.approve(retired1);
      await treasuryRebalance.finalizeApproval();
      await hre.network.provider.send("hardhat_mine", ["0x32"]);
      await treasuryRebalance.finalizeContract(memo);
      expect(await treasuryRebalance.status()).to.equal(3);

      expect(await treasuryRebalance.getRetiredCount()).to.equal(1);
      expect(await treasuryRebalance.getNewbieCount()).to.equal(1);

      await expect(treasuryRebalance.reset()).to.be.revertedWith(
        "Contract is finalized, cannot reset values"
      );
    });

    it("should revert when tried to reset after it passes the execution block", async function() {
      await hre.network.provider.send("hardhat_mine", ["0x32"]);
      await expect(treasuryRebalance.reset()).to.be.revertedWith(
        "Contract is finalized, cannot reset values"
      );
    });

    it("Should not allow non-owner to reset", async function() {
      await expect(
        treasuryRebalance.connect(account1).reset()
      ).to.be.revertedWith("Ownable: caller is not the owner");
    });
  });

  describe("fallback", function() {
    it("should revert if KLAY is sent to the contract address", async function() {
      await expect(
        owner.sendTransaction({
          to: treasuryRebalance.address,
          value: hre.ethers.utils.parseEther("20"),
        })
      ).to.be.revertedWith("This contract does not accept any payments");
    });
  });

  describe("isContract", function() {
    it("should check whether EOA/contract address", async function() {
      const eoa = await treasuryRebalance.isContractAddr(owner.address);
      const contractAddress = await treasuryRebalance.isContractAddr(
        treasuryRebalance.address
      );
      expect(eoa).to.equal(false);
      expect(contractAddress).to.equal(true);
    });
  });
});
