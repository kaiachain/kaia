import { loadFixture, setCode } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import {
  augmentChai,
  nowTime,
  jumpBlock,
  ABOOK_ADDRESS,
  nowBlock,
  padUtils,
} from "../common/helper";
import { ethers, upgrades } from "hardhat";

// TODO: Manage addresses for system contracts in a single file
const REGISTRY_ADDRESS = "0x0000000000000000000000000000000000000401";
const KIP103_ADDRESS = "0xD5ad6D61Dd87EdabE2332607C328f5cc96aeCB95";
const GOVERNANCE_ADDRESS = "0xcA4Ef926634A530f12e55A0aEE87F195A7B22Aa3";
const STAKING_TRACKER_ADDRESS = "0x9b8688d616D3D5180d29520c6a0E28582E82BF4d";
const RAND_ADDR = "0xe3B0C44298FC1C149AfBF4C8996fb92427aE41E4"; // non-null placeholder

const preDeployedAddresses = [
  ABOOK_ADDRESS,
  KIP103_ADDRESS,
  GOVERNANCE_ADDRESS,
  STAKING_TRACKER_ADDRESS,
];
const preDeployedActivation = ["0x", "0x", "0x", "0x"]; //0, 0, 0, 0
const preDeployedName = [
  "AddressBook",
  "TreasuryRebalance",
  "Voting",
  "StakingTracker",
];

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
  async function buildFixture() {
    // Prepare parameters for deploying contracts
    const accounts = await ethers.getSigners();
    const [a1, a2, a3, a4, a5] = accounts;

    // Deploy registry
    const registryFactory = await ethers.getContractFactory("Registry");
    let registry = await registryFactory.deploy();
    await registry.deployed();

    const registryByteCode = await ethers.provider.getCode(registry.address);
    await setCode(REGISTRY_ADDRESS, registryByteCode);
    registry = await ethers.getContractAt("Registry", REGISTRY_ADDRESS);

    // Inject pre-deployed system contracts by setStorageAt:
    // 1. mapping(string => Records[]) public records -> slot 0
    //   1. MappingKey: keccak256(abi.encode(key, slot))
    //   2. Slot: keccak256(abi.encode(MappingKey))
    //   3. ArraySlot: bytes32(uint256(keccak256(abi.encode(Slot))) + i)
    // 2. string[] public names -> slot 1
    //   1. Array length: keccak256(abi.encode(slot))
    //   2. Array element: bytes32(uint256(keccak256(abi.encode(slot))) + i) (Since there's no string more than slot size for pre-deployed contracts)
    // 3. address private _owner -> slot 2
    //   - owner: at slot 2
    // Reference: https://docs.soliditylang.org/en/v0.8.20/internals/layout_in_storage.html

    const baseByte = ethers.utils.hexlify(ethers.utils.zeroPad("0x00", 32));
    const paddedSlot = ethers.utils.zeroPad("0x01", 32);
    const elemSlot = ethers.utils.solidityKeccak256(["bytes"], [paddedSlot]);

    // Total 4 pre-deployed system contracts will be registered
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [
        registry.address,
        "0x" + Number(1).toString(16),
        padUtils("0x04", 32),
      ],
    });

    for (let i = 0; i < preDeployedAddresses.length; i++) {
      // For mapping(string => Record[]) public records
      const byteKey = ethers.utils.hexlify(
        ethers.utils.toUtf8Bytes(preDeployedName[i])
      );
      const mappingKey = byteKey + baseByte.slice(2); // String key doesn't need to be padded
      const slot = ethers.utils.solidityKeccak256(["bytes"], [mappingKey]);
      const arraySlot = ethers.utils.solidityKeccak256(["bytes"], [slot]);

      const arrayElemSlotAddr = ethers.BigNumber.from(arraySlot).toHexString();
      const arrayElemSlotActivation = ethers.BigNumber.from(arraySlot)
        .add(1)
        .toHexString();

      // Length of array
      await hre.network.provider.request({
        method: "hardhat_setStorageAt",
        params: [registry.address, slot, padUtils("0x01", 32)],
      });

      // Set Record.addr
      await hre.network.provider.request({
        method: "hardhat_setStorageAt",
        params: [
          registry.address,
          arrayElemSlotAddr,
          padUtils(preDeployedAddresses[i], 32),
        ],
      });

      // Set Record.activation
      await hre.network.provider.request({
        method: "hardhat_setStorageAt",
        params: [
          registry.address,
          arrayElemSlotActivation,
          padUtils(preDeployedActivation[i], 32),
        ],
      });

      // For string[] public names
      const elemSlotAt = ethers.BigNumber.from(elemSlot).add(i).toHexString();

      const byteArray = ethers.utils.toUtf8Bytes(preDeployedName[i]);
      const hexString = ethers.utils.hexlify(byteArray);
      const length = byteArray.length * 2;

      const lengthBytes = ethers.utils.hexlify(
        ethers.utils.zeroPad(ethers.utils.hexlify(length), 32)
      );
      const storedString = hexString + lengthBytes.slice(2 + length);

      // Set names[i]
      await hre.network.provider.request({
        method: "hardhat_setStorageAt",
        params: [registry.address, elemSlotAt, storedString],
      });
    }

    // For address private _owner
    await hre.network.provider.request({
      method: "hardhat_setStorageAt",
      params: [
        registry.address,
        "0x" + Number(2).toString(16),
        padUtils(a1.address, 32),
      ],
    });

    // Deploy mock upgradeable system contract
    const MockUpgradeableSystemContract = await ethers.getContractFactory(
      "MockUpgradeableSystemContract"
    );
    const mockProxy = await upgrades.deployProxy(
      MockUpgradeableSystemContract,
      [3],
      {
        initializer: "initialize",
        kind: "uups",
      }
    );

    return {
      a1,
      a2,
      a3,
      a4,
      a5,
      registry,
      mockProxy,
    };
  }
  let fixture: UnPromisify<ReturnType<typeof buildFixture>>;
  beforeEach(async function () {
    augmentChai();
    fixture = await loadFixture(buildFixture);
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
        expect(await registry.getActiveAddr(preDeployedName[i])).to.equal(
          preDeployedAddresses[i]
        );
        expect(await registry.names(i)).to.equal(preDeployedName[i]);
      }
    });
  });

  describe("Transfer owner of registry", function () {
    it("#transferOwnership: Can't transfer ownership to zero address", async function () {
      const { registry } = fixture;

      await expect(
        registry.transferOwnership(ethers.constants.AddressZero)
      ).to.be.revertedWith("Zero address");
    });
    it("#transferOwnership: Successfully transfer ownership", async function () {
      const { registry, a2 } = fixture;

      await expect(
        registry.connect(a2).transferOwnership(a2.address)
      ).to.be.revertedWith("Not owner");

      await expect(registry.transferOwnership(a2.address)).to.emit(
        registry,
        "OwnershipTransferred"
      );

      expect(await registry.owner()).to.equal(a2.address);
    });
  });

  describe("Check register function", function () {
    it("#register: Can't register empty name", async function () {
      const { registry } = fixture;

      await expect(
        registry.register("", RAND_ADDR, activeBlockForMockUpgradeable)
      ).to.be.revertedWith("Empty string");
    });
    it("#register: Can't register contract from past", async function () {
      const { registry } = fixture;

      const now = await nowBlock();
      await expect(
        registry.register("MockUpgradeable", RAND_ADDR, now - 1)
      ).to.be.revertedWith("Can't register contract from past");
    });
    it("#register: Successfully register system contract", async function () {
      const { registry, mockProxy } = fixture;

      await expect(
        registry.register(
          "MockUpgradeable",
          mockProxy.address,
          activeBlockForMockUpgradeable
        )
      ).to.emit(registry, "Registered");

      // Not active yet
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        ethers.constants.AddressZero
      );

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        mockProxy.address
      );

      // Check names
      expect(await registry.names(4)).to.equal("MockUpgradeable");
    });
    it("#register: Replace contract by registration", async function () {
      const { registry, mockProxy } = fixture;

      await expect(
        registry.register(
          "MockUpgradeable",
          mockProxy.address,
          activeBlockForMockUpgradeable
        )
      ).to.emit(registry, "Registered");

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        mockProxy.address
      );

      // Replace contract
      await expect(
        registry.register(
          "MockUpgradeable",
          RAND_ADDR,
          activeBlockForMockUpgradeable * 2
        )
      ).to.emit(registry, "Registered");

      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        mockProxy.address
      );

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's replaced to new contract
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        RAND_ADDR
      );
    });
    it("#register: Can overwrite the most recently registered contract", async function () {
      const { registry, mockProxy } = fixture;

      await expect(
        registry.register(
          "MockUpgradeable",
          mockProxy.address,
          activeBlockForMockUpgradeable
        )
      ).to.emit(registry, "Registered");

      // Not active yet
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        ethers.constants.AddressZero
      );

      // It overwrites the previous contract since it is not active yet
      await expect(
        registry.register(
          "MockUpgradeable",
          RAND_ADDR,
          activeBlockForMockUpgradeable
        )
      ).to.emit(registry, "Registered");

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's active
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        RAND_ADDR
      );
    });
    it("#register: Register zero address means deprecation", async function () {
      const { registry, mockProxy } = fixture;

      await expect(
        registry.register(
          "MockUpgradeable",
          mockProxy.address,
          activeBlockForMockUpgradeable
        )
      ).to.emit(registry, "Registered");

      await jumpBlock(activeBlockForMockUpgradeable);

      // It overwrites the previous contract since it is not active yet
      await expect(
        registry.register(
          "MockUpgradeable",
          ethers.constants.AddressZero,
          activeBlockForMockUpgradeable * 2
        )
      ).to.emit(registry, "Registered");

      await jumpBlock(activeBlockForMockUpgradeable);

      // Now it's considered as deprecated
      expect(await registry.getActiveAddr("MockUpgradeable")).to.equal(
        ethers.constants.AddressZero
      );
    });
    describe("Pre-deployed system contracts", function () {
      it("Register new contract to replace pre-deployed contract: AddressBook", async function () {
        const { registry } = fixture;

        // AddressBook
        await expect(
          registry.register(
            "AddressBook",
            RAND_ADDR,
            activeBlockForMockUpgradeable
          )
        ).to.emit(registry, "Registered");

        const records = await registry.getAllRecords("AddressBook");
        expect(records.length).to.equal(2);

        // Not active yet
        expect(await registry.getActiveAddr("AddressBook")).to.equal(
          ABOOK_ADDRESS
        );

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
        await expect(
          registry.register(
            "StakingTracker",
            RAND_ADDR,
            activeBlockForMockUpgradeable
          )
        ).to.emit(registry, "Registered");

        const records = await registry.getAllRecords("StakingTracker");
        expect(records.length).to.equal(2);

        // Not active yet
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          STAKING_TRACKER_ADDRESS
        );

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          RAND_ADDR
        );

        // It shouldn't append "StakingTracker" to names
        const names = await registry.getAllNames();
        expect(names.length).to.equal(4);
      });
      it("Deprecate pre-deployed contract: AddressBook", async function () {
        const { registry } = fixture;

        // AddressBook
        await expect(
          registry.register(
            "AddressBook",
            RAND_ADDR,
            activeBlockForMockUpgradeable
          )
        ).to.emit(registry, "Registered");

        // Not active yet
        expect(await registry.getActiveAddr("AddressBook")).to.equal(
          ABOOK_ADDRESS
        );

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("AddressBook")).to.equal(RAND_ADDR);

        // Deprecate AddressBook by registering zero address.
        await expect(
          registry.register(
            "AddressBook",
            ethers.constants.AddressZero,
            activeBlockForMockUpgradeable * 2
          )
        ).to.emit(registry, "Registered");

        expect(await registry.getActiveAddr("AddressBook")).to.equal(RAND_ADDR);

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's considered as deprecated
        expect(await registry.getActiveAddr("AddressBook")).to.equal(
          ethers.constants.AddressZero
        );
      });
      it("Deprecate pre-deployed contract: StakingTracker", async function () {
        const { registry } = fixture;

        // StakingTracker
        await expect(
          registry.register(
            "StakingTracker",
            RAND_ADDR,
            activeBlockForMockUpgradeable
          )
        ).to.emit(registry, "Registered");

        // Not active yet
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          STAKING_TRACKER_ADDRESS
        );

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          RAND_ADDR
        );

        // Deprecate StakingTracker by registering zero address.
        await expect(
          registry.register(
            "StakingTracker",
            ethers.constants.AddressZero,
            activeBlockForMockUpgradeable * 2
          )
        ).to.emit(registry, "Registered");

        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          RAND_ADDR
        );

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's considered as deprecated
        expect(await registry.getActiveAddr("StakingTracker")).to.equal(
          ethers.constants.AddressZero
        );
      });
    });

    describe("Contract replacement test", function () {
      this.beforeEach(async function () {
        const { registry, mockProxy } = fixture;

        await expect(
          registry.register(
            "MockUpgradeable",
            mockProxy.address,
            activeBlockForMockUpgradeable
          )
        ).to.emit(registry, "Registered");
      });
      it("Update mock upgradeable system contracts by proxy", async function () {
        const { mockProxy } = fixture;

        // Before upgrade, it doesn't have setTime method
        expect(await mockProxy.number()).to.equal(3);
        expect(await mockProxy.setTime).to.be.undefined;

        // 2. Upgrade logic contract of proxy
        const MockUpgradeableSystemContractV2 = await ethers.getContractFactory(
          "MockUpgradeableSystemContractV2"
        );
        let mockProxyWithV2 = await upgrades.upgradeProxy(
          mockProxy.address,
          MockUpgradeableSystemContractV2
        );

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
        await expect(registry.transferOwnership(a3.address)).to.emit(
          registry,
          "OwnershipTransferred"
        );

        expect(await registry.owner()).to.equal(a3.address);

        // 3. Register new voting (governance) contract to Registry
        await expect(
          registry
            .connect(a3)
            .register("Voting", a3.address, activeBlockForMockUpgradeable)
        ).to.emit(registry, "Registered");

        // Not active yet
        expect(await registry.getActiveAddr("Voting")).to.equal(
          GOVERNANCE_ADDRESS
        );

        await jumpBlock(activeBlockForMockUpgradeable);

        // Now it's active
        expect(await registry.getActiveAddr("Voting")).to.equal(a3.address);
      });
    });
  });
});
