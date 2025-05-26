import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import { augmentChai, nowTime, jumpBlock, ABOOK_ADDRESS, nowBlock } from "../common/helper";
import { ethers, upgrades } from "hardhat";
import {
  GOVERNANCE_ADDRESS,
  preDeployedAddresses,
  preDeployedName,
  RAND_ADDR,
  registryFixture,
  STAKING_TRACKER_ADDRESS,
} from "../materials";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

/**
 * @dev This test is for Registry.sol
 * 1. Test initializing process
 *    - Check injected pre-deployed system contracts
 * 2. Test functions of Registry.sol
 *    - transferOwnership
 *    - register
 * 3. Test pre-deployed system contract
 * 4. Test contract replacement
 *    - Contract upgrade test
 *    - Governance replacement test
 */
describe("Registry.sol", function () {
  let activeBlockForMockUpgradeable: number;

  let fixture: UnPromisify<ReturnType<typeof registryFixture>>;
  beforeEach(async function () {
    augmentChai();
    fixture = await loadFixture(registryFixture);
    activeBlockForMockUpgradeable = (await nowBlock()) + 50;
  });

  describe("Registry Initialize", function () {
    it("Check constructor", async function () {
      const { registry, a1 } = fixture;

      expect(await registry.owner()).to.equal(a1.address);
    });

    it("Check pre-deployed system contracts", async function () {
      const { registry } = fixture;

      for (let i = 0; i < preDeployedAddresses.length; i++) {
        expect(await registry.getActiveAddr(preDeployedName[i])).to.equal(preDeployedAddresses[i]);
        expect(await registry.names(i)).to.equal(preDeployedName[i]);
      }
    });
  });

  describe("Transfer owner of registry", function () {
    it("#transferOwnership: Can't transfer ownership to zero address", async function () {
      const { registry } = fixture;

      await expect(registry.transferOwnership(ethers.constants.AddressZero)).to.be.revertedWith("Zero address");
    });
    it("#transferOwnership: Successfully transfer ownership", async function () {
      const { registry, a2 } = fixture;

      await expect(registry.connect(a2).transferOwnership(a2.address)).to.be.revertedWith("Not owner");

      await expect(registry.transferOwnership(a2.address)).to.emit(registry, "OwnershipTransferred");

      expect(await registry.owner()).to.equal(a2.address);
    });
  });

  describe("Check register function", function () {
    it("#register: Can't register empty name", async function () {
      const { registry } = fixture;

      await expect(registry.register("", RAND_ADDR, activeBlockForMockUpgradeable)).to.be.revertedWith("Empty string");
    });
    it("#register: Can't register contract from past", async function () {
      const { registry } = fixture;

      const now = await nowBlock();
      await expect(registry.register("MockUpgradeable", RAND_ADDR, now - 1)).to.be.revertedWith(
        "Can't register contract from past",
      );
    });
    it("#register: Successfully register system contract", async function () {
      const { registry, mockProxy } = fixture;

      await expect(registry.register("MockUpgradeable", mockProxy.address, activeBlockForMockUpgradeable)).to.emit(
        registry,
        "Registered",
      );

      // Not active yet
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(ethers.constants.AddressZero);

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(mockProxy.address);

      // Check names
      expect(await registry.names(4)).to.equal("MockUpgradeable");
    });
    it("#register: Replace contract by registration", async function () {
      const { registry, mockProxy } = fixture;

      await expect(registry.register("MockUpgradeable", mockProxy.address, activeBlockForMockUpgradeable)).to.emit(
        registry,
        "Registered",
      );

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(mockProxy.address);

      // Replace contract
      await expect(registry.register("MockUpgradeable", RAND_ADDR, activeBlockForMockUpgradeable * 2)).to.emit(
        registry,
        "Registered",
      );

      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(mockProxy.address);

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's replaced to new contract
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(RAND_ADDR);
    });
    it("#register: Can overwrite the most recently registered contract", async function () {
      const { registry, mockProxy } = fixture;

      await expect(registry.register("MockUpgradeable", mockProxy.address, activeBlockForMockUpgradeable)).to.emit(
        registry,
        "Registered",
      );

      // Not active yet
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(ethers.constants.AddressZero);

      // It overwrites the previous contract since it is not active yet
      await expect(registry.register("MockUpgradeable", RAND_ADDR, activeBlockForMockUpgradeable)).to.emit(
        registry,
        "Registered",
      );

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(RAND_ADDR);
    });
    it("#register: Register zero address means deprecation", async function () {
      const { registry, mockProxy } = fixture;

      await expect(registry.register("MockUpgradeable", mockProxy.address, activeBlockForMockUpgradeable)).to.emit(
        registry,
        "Registered",
      );

      await jumpBlock(activeBlockForMockUpgradeable);

      // It overwrites the previous contract since it is not active yet
      await expect(
        registry.register("MockUpgradeable", ethers.constants.AddressZero, activeBlockForMockUpgradeable * 2),
      ).to.emit(registry, "Registered");

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's considered as deprecated
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(ethers.constants.AddressZero);
    });
    describe("Pre-deployed system contracts", function () {
      it("Register new contract to replace pre-deployed contract: AddressBook", async function () {
        const { registry } = fixture;

        // AddressBook
        await expect(registry.register("AddressBook", RAND_ADDR, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );

        const records = await registry.getAllRecords("AddressBook");
        expect(records.length).to.equal(2);

        // Not active yet
        expect(await registry.getActiveAddr("AddressBook")).to.equal(ABOOK_ADDRESS);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("AddressBook")).to.equal(RAND_ADDR);

        // It shouldn't append "AddressBook" to names
        const names = await registry.getAllNames();
        expect(names.length).to.equal(4);
      });
      it("Register new contract to replace pre-deployed contract: StakingTracker", async function () {
        const { registry } = fixture;

        // StakingTracker
        await expect(registry.register("StakingTracker", RAND_ADDR, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );

        const records = await registry.getAllRecords("StakingTracker");
        expect(records.length).to.equal(2);

        // Not active yet
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(STAKING_TRACKER_ADDRESS);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(RAND_ADDR);

        // It shouldn't append "StakingTracker" to names
        const names = await registry.getAllNames();
        expect(names.length).to.equal(4);
      });
      it("Deprecate pre-deployed contract: AddressBook", async function () {
        const { registry } = fixture;

        // AddressBook
        await expect(registry.register("AddressBook", RAND_ADDR, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );

        // Not active yet
        expect(await registry.getActiveAddr("AddressBook")).to.equal(ABOOK_ADDRESS);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("AddressBook")).to.equal(RAND_ADDR);

        // Deprecate AddressBook by registering zero address.
        await expect(
          registry.register("AddressBook", ethers.constants.AddressZero, activeBlockForMockUpgradeable * 2),
        ).to.emit(registry, "Registered");

        expect(await registry.getActiveAddr("AddressBook")).to.equal(RAND_ADDR);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's considered as deprecated
        expect(await registry.getActiveAddr("AddressBook")).to.equal(ethers.constants.AddressZero);
      });
      it("Deprecate pre-deployed contract: StakingTracker", async function () {
        const { registry } = fixture;

        // StakingTracker
        await expect(registry.register("StakingTracker", RAND_ADDR, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );

        // Not active yet
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(STAKING_TRACKER_ADDRESS);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(RAND_ADDR);

        // Deprecate StakingTracker by registering zero address.
        await expect(
          registry.register("StakingTracker", ethers.constants.AddressZero, activeBlockForMockUpgradeable * 2),
        ).to.emit(registry, "Registered");

        expect(await registry.getActiveAddr("StakingTracker")).to.equal(RAND_ADDR);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's considered as deprecated
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(ethers.constants.AddressZero);
      });
    });

    describe("Contract replacement test", function () {
      this.beforeEach(async function () {
        const { registry, mockProxy } = fixture;

        await expect(registry.register("MockUpgradeable", mockProxy.address, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );
      });
      it("Update mock upgradeable system contracts by proxy", async function () {
        const { mockProxy } = fixture;

        // Before upgrade, it doesn't have setTime method
        expect(await mockProxy.number()).to.equal(3);
        expect(await mockProxy.setTime).to.be.undefined;

        // 2. Upgrade logic contract of proxy
        const MockUpgradeableSystemContractV2 = await ethers.getContractFactory("MockUpgradeableSystemContractV2");
        const mockProxyWithV2 = await upgrades.upgradeProxy(mockProxy.address, MockUpgradeableSystemContractV2);

        // 3. Check update process is done successfully
        await mockProxyWithV2.increment();

        expect(await mockProxyWithV2.number()).to.equal(4);
        // Now we have `setTime` variable in system contract
        expect(await mockProxyWithV2.setTime()).to.equal(await nowTime());
      });
    });
    describe("Governance replacement test", function () {
      it("Replace governance contract", async function () {
        const { registry, a3 } = fixture;

        // 1. Update governance contract (skip now)
        // 2. transferOwnership to new governance contract
        // 3. Register new voting (governance) contract to Registry

        // Assume a3.address is new governance contract
        // 2. transferOwnership to new governance contract
        await expect(registry.transferOwnership(a3.address)).to.emit(registry, "OwnershipTransferred");

        expect(await registry.owner()).to.equal(a3.address);

        // 3. Register new voting (governance) contract to Registry
        await expect(registry.connect(a3).register("Voting", a3.address, activeBlockForMockUpgradeable)).to.emit(
          registry,
          "Registered",
        );

        // Not active yet
        expect(await registry.getActiveAddr("Voting")).to.equal(GOVERNANCE_ADDRESS);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("Voting")).to.equal(a3.address);
      });
    });
  });
});
