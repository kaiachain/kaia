import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { clRegistryTestFixture } from "../common/fixtures";
import { expect } from "chai";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("CLRegistry add/remove/update", function () {
  let fixture: UnPromisify<ReturnType<typeof clRegistryTestFixture>>;
  const gcId = 1;
  const nodeId = "0x0000000000000000000000000000000000000001";
  const clPool = "0x0000000000000000000000000000000000000002";
  const clStaking = "0x0000000000000000000000000000000000000003";

  beforeEach(async function () {
    fixture = await loadFixture(clRegistryTestFixture);
  });

  it("Add pair", async function () {
    const { clRegistry } = fixture;

    await expect(clRegistry.addCLPair([[nodeId, gcId, clPool, clStaking]])).to.emit(clRegistry, "RegisterPair");

    // Invalid pair input
    await expect(clRegistry.addCLPair([[nodeId, 0, clPool, clStaking]])).to.be.revertedWith(
      "CLRegistry::addCLPair: Invalid pair input",
    );

    // Duplicated registration
    await expect(clRegistry.addCLPair([[nodeId, gcId, clPool, clStaking]])).to.be.revertedWith(
      "CLRegistry::addCLPair: GC ID does exist",
    );

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(
      clRegistry.connect(unknown).addCLPair([[nodeId, gcId, clPool, clStaking]]),
    ).to.be.revertedWithCustomError(clRegistry, "OwnableUnauthorizedAccount");
  });

  it("Add pair(s)", async function () {
    const { clRegistry } = fixture;

    const nodeId2 = "0x000000000000000000000000000000000000000A";
    const clPool2 = "0x000000000000000000000000000000000000000B";
    const clStaking2 = "0x000000000000000000000000000000000000000C";
    const gcId2 = 2;
    const nodeId3 = "0x000000000000000000000000000000000000000D";
    const clPool3 = "0x000000000000000000000000000000000000000E";
    const clStaking3 = "0x000000000000000000000000000000000000000F";
    const gcId3 = 3;

    await expect(
      clRegistry.addCLPair([
        [nodeId, gcId, clPool, clStaking],
        [nodeId2, gcId2, clPool2, clStaking2],
        [nodeId3, gcId3, clPool3, clStaking3],
      ]),
    ).to.emit(clRegistry, "RegisterPair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [clPool, clPool2, clPool3],
      [clStaking, clStaking2, clStaking3],
    ]);

    expect(await clRegistry.getAllGCIds()).to.deep.equal([gcId, gcId2, gcId3]);
  });

  it("Remove pair", async function () {
    const { clRegistry } = fixture;

    // Register a pair
    await expect(clRegistry.addCLPair([[nodeId, gcId, clPool, clStaking]])).to.emit(clRegistry, "RegisterPair");

    // Remoeve a pair registered before
    await expect(clRegistry.removeCLPair(gcId)).to.emit(clRegistry, "RetirePair");

    // Try to remove a not registered pair
    await expect(clRegistry.removeCLPair(gcId)).to.be.revertedWith("CLRegistry::removeCLPair: GC ID does not exist");

    // Invalid GC ID
    await expect(clRegistry.removeCLPair(0)).to.be.revertedWith("CLRegistry::removeCLPair: Invalid GC ID");

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(clRegistry.connect(unknown).removeCLPair(gcId)).to.be.revertedWithCustomError(
      clRegistry,
      "OwnableUnauthorizedAccount",
    );
  });

  it("Update pair", async function () {
    const { clRegistry } = fixture;

    // Register a pair
    await expect(clRegistry.addCLPair([[nodeId, gcId, clPool, clStaking]])).to.emit(clRegistry, "RegisterPair");

    // Update a pair
    const newCLPool = "0x000000000000000000000000000000000000001A";
    await expect(clRegistry.updateCLPair([[nodeId, gcId, newCLPool, clStaking]])).to.emit(clRegistry, "UpdatePair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([[nodeId], [gcId], [newCLPool], [clStaking]]);

    // Invalid pair input
    await expect(clRegistry.updateCLPair([[nodeId, 0, clPool, clStaking]])).to.be.revertedWith(
      "CLRegistry::updateCLPair: Invalid pair input",
    );

    // Try to update not existing pair
    await expect(clRegistry.updateCLPair([[nodeId, gcId + 1, clPool, clStaking]])).to.be.revertedWith(
      "CLRegistry::updateCLPair: GC ID does not exist",
    );

    // ownership check
    const [, unknown] = await ethers.getSigners();
    await expect(
      clRegistry.connect(unknown).updateCLPair([[nodeId, gcId, clPool, clStaking]]),
    ).to.be.revertedWithCustomError(clRegistry, "OwnableUnauthorizedAccount");
  });

  it("Update pair(s)", async function () {
    const { clRegistry } = fixture;

    const nodeId2 = "0x000000000000000000000000000000000000000A";
    const clPool2 = "0x000000000000000000000000000000000000000B";
    const clStaking2 = "0x000000000000000000000000000000000000000C";
    const gcId2 = 2;
    const nodeId3 = "0x000000000000000000000000000000000000000D";
    const clPool3 = "0x000000000000000000000000000000000000000E";
    const clStaking3 = "0x000000000000000000000000000000000000000F";
    const gcId3 = 3;

    const newCLPool1 = "0x000000000000000000000000000000000000001A";
    const newCLPool2 = "0x000000000000000000000000000000000000002A";
    const newCLPool3 = "0x000000000000000000000000000000000000003A";

    await expect(
      clRegistry.addCLPair([
        [nodeId, gcId, clPool, clStaking],
        [nodeId2, gcId2, clPool2, clStaking2],
        [nodeId3, gcId3, clPool3, clStaking3],
      ]),
    ).to.emit(clRegistry, "RegisterPair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [clPool, clPool2, clPool3],
      [clStaking, clStaking2, clStaking3],
    ]);

    await expect(
      clRegistry.updateCLPair([
        [nodeId, gcId, newCLPool1, clStaking],
        [nodeId2, gcId2, newCLPool2, clStaking2],
        [nodeId3, gcId3, newCLPool3, clStaking3],
      ]),
    ).to.emit(clRegistry, "UpdatePair");

    expect(await clRegistry.getAllCLs()).to.deep.equal([
      [nodeId, nodeId2, nodeId3],
      [gcId, gcId2, gcId3],
      [newCLPool1, newCLPool2, newCLPool3],
      [clStaking, clStaking2, clStaking3],
    ]);
  });
});
