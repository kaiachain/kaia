import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { clRegistryTestFixture } from "../materials";
import { expect } from "chai";
import { ethers } from "hardhat";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("CLRegistry add/remove/update", function () {
  let fixture: UnPromisify<ReturnType<typeof clRegistryTestFixture>>;
  const gcId = 1;
  const nodeId = "0x0000000000000000000000000000000000000001";
  const clPool = "0x0000000000000000000000000000000000000002";

  beforeEach(async function () {
    fixture = await loadFixture(clRegistryTestFixture);
  });

  it("Add pair", async function () {
    const { clRegistry } = fixture;

    await expect(clRegistry.addCLPair([{ nodeId, gcId, clPool }])).to.emit(
      clRegistry,
      "RegisterPair"
    );

    // Invalid pair input
    await expect(
      clRegistry.addCLPair([{ nodeId, gcId: 0, clPool }])
    ).to.be.revertedWith("CLRegistry::addCLPair: Invalid pair input");

    // Duplicated registration
    await expect(
      clRegistry.addCLPair([{ nodeId, gcId, clPool }])
    ).to.be.revertedWith("CLRegistry::addCLPair: GC ID does exist");

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(
      clRegistry.connect(unknown).addCLPair([{ nodeId, gcId, clPool }])
    ).to.be.revertedWithCustomError(clRegistry, "OwnableUnauthorizedAccount");
  });

  it("Add pair(s)", async function () {
    const { clRegistry } = fixture;

    const nodeId2 = "0x000000000000000000000000000000000000000A";
    const clPool2 = "0x000000000000000000000000000000000000000B";
    const gcId2 = 2;
    const nodeId3 = "0x000000000000000000000000000000000000000D";
    const clPool3 = "0x000000000000000000000000000000000000000E";
    const gcId3 = 3;

    await expect(
      clRegistry.addCLPair([
        { nodeId, gcId, clPool },
        { nodeId: nodeId2, gcId: gcId2, clPool: clPool2 },
        { nodeId: nodeId3, gcId: gcId3, clPool: clPool3 },
      ])
    ).to.emit(clRegistry, "RegisterPair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [clPool, clPool2, clPool3],
    ]);

    expect(await clRegistry.getAllGCIds()).to.deep.equal([gcId, gcId2, gcId3]);
  });

  it("Remove pair", async function () {
    const { clRegistry } = fixture;

    // Register a pair
    await expect(clRegistry.addCLPair([{ nodeId, gcId, clPool }])).to.emit(
      clRegistry,
      "RegisterPair"
    );

    // Remoeve a pair registered before
    await expect(clRegistry.removeCLPair(gcId)).to.emit(
      clRegistry,
      "RetirePair"
    );

    // Try to remove a not registered pair
    await expect(clRegistry.removeCLPair(gcId)).to.be.revertedWith(
      "CLRegistry::removeCLPair: GC ID does not exist"
    );

    // Invalid GC ID
    await expect(clRegistry.removeCLPair(0)).to.be.revertedWith(
      "CLRegistry::removeCLPair: Invalid GC ID"
    );

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(
      clRegistry.connect(unknown).removeCLPair(gcId)
    ).to.be.revertedWithCustomError(clRegistry, "OwnableUnauthorizedAccount");
  });

  it("Update pair", async function () {
    const { clRegistry } = fixture;

    // Register a pair
    await expect(clRegistry.addCLPair([{ nodeId, gcId, clPool }])).to.emit(
      clRegistry,
      "RegisterPair"
    );

    // Update a pair
    const newCLPool = "0x000000000000000000000000000000000000001A";
    await expect(
      clRegistry.updateCLPair([{ nodeId, gcId, clPool: newCLPool }])
    ).to.emit(clRegistry, "UpdatePair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId],
      [gcId],
      [newCLPool],
    ]);

    // Invalid pair input
    await expect(
      clRegistry.updateCLPair([{ nodeId, gcId: 0, clPool }])
    ).to.be.revertedWith("CLRegistry::updateCLPair: Invalid pair input");

    // Try to update not existing pair
    await expect(
      clRegistry.updateCLPair([{ nodeId, gcId: gcId + 1, clPool }])
    ).to.be.revertedWith("CLRegistry::updateCLPair: GC ID does not exist");

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(
      clRegistry.connect(unknown).updateCLPair([{ nodeId, gcId, clPool }])
    ).to.be.revertedWithCustomError(clRegistry, "OwnableUnauthorizedAccount");
  });

  it("Update pair(s)", async function () {
    const { clRegistry } = fixture;

    const nodeId2 = "0x000000000000000000000000000000000000000A";
    const clPool2 = "0x000000000000000000000000000000000000000B";
    const gcId2 = 2;
    const nodeId3 = "0x000000000000000000000000000000000000000D";
    const clPool3 = "0x000000000000000000000000000000000000000E";
    const gcId3 = 3;

    const newCLPool1 = "0x000000000000000000000000000000000000001A";
    const newCLPool2 = "0x000000000000000000000000000000000000002A";
    const newCLPool3 = "0x000000000000000000000000000000000000003A";

    await expect(
      clRegistry.addCLPair([
        { nodeId, gcId, clPool },
        { nodeId: nodeId2, gcId: gcId2, clPool: clPool2 },
        { nodeId: nodeId3, gcId: gcId3, clPool: clPool3 },
      ])
    ).to.emit(clRegistry, "RegisterPair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [clPool, clPool2, clPool3],
    ]);

    await expect(
      clRegistry.updateCLPair([
        { nodeId, gcId, clPool: newCLPool1 },
        { nodeId: nodeId2, gcId: gcId2, clPool: newCLPool2 },
        { nodeId: nodeId3, gcId: gcId3, clPool: newCLPool3 },
      ])
    ).to.emit(clRegistry, "UpdatePair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [newCLPool1, newCLPool2, newCLPool3],
    ]);
  });
});
