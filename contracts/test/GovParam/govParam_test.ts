import { ethers } from "hardhat";
import { expect } from "chai";
import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { nowBlock, jumpBlock } from "../common/helper";
import { govParamTestFixture } from "../materials";

const NOT_OWNER = "Ownable: caller is not the owner";
const EMPTYNAME = "GovParam: name cannot be empty";
const EMPTY_VALUE = "GovParam: val must not be empty if exists=true";
const NONEMPTY_VALUE = "GovParam: val must be empty if exists=false";
const ALREADY_PAST = "GovParam: activation must be in the future";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

describe("GovParam", function () {
  let fixture: UnPromisify<ReturnType<typeof govParamTestFixture>>;
  beforeEach(async function () {
    fixture = await loadFixture(govParamTestFixture);
  });

  const name = "istanbul.committeesize";
  const val1 = "0x1234";
  const val2 = "0x5678";
  const emptyVal = "0x";

  const defaultParam = [ethers.BigNumber.from("0"), false, emptyVal];

  describe("constructor", function () {
    it("Constructor success", async function () {
      const { gp, deployer } = fixture;

      expect(await gp.owner()).to.equal(deployer.address);
    });
  });

  describe("setParam", function () {
    it("setParam success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 5;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);

      const before = [false, emptyVal];
      const after = [true, val1];

      for (let i = 0; i < 10; i++) {
        await jumpBlock(1);
        const p = await gp.getParam(name);
        const now = await nowBlock();
        if (now < activation) {
          expect(p).to.deep.equal(before);
        } else {
          expect(p).to.deep.equal(after);
        }
      }
    });

    it("setParam overwrite success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 5;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);

      const newActivation = activation + 5;
      await expect(gp.setParam(name, true, val2, newActivation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val2, newActivation);

      const before = [false, emptyVal];
      const after = [true, val2];

      for (let i = 0; i < 20; i++) {
        await jumpBlock(1);
        const p = await gp.getParam(name);
        const now = await nowBlock();
        if (now < newActivation) {
          expect(p).to.deep.equal(before);
        } else {
          expect(p).to.deep.equal(after);
        }
      }
    });

    it("setParam fails with empty name", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam("", true, val1, activation)).to.be.revertedWith(EMPTYNAME);
    });

    it("setParam fails with past activation", async function () {
      const { gp } = fixture;
      const activation = await nowBlock();
      await expect(gp.setParam(name, true, val1, activation)).to.be.revertedWith(ALREADY_PAST);
    });

    it("setParam fails with empty value except delete", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, emptyVal, activation)).to.be.revertedWith(EMPTY_VALUE);
    });

    it("setParam deleting success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);

      await jumpBlock(10);
      const newActivation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, false, emptyVal, newActivation))
        .to.emit(gp, "SetParam")
        .withArgs(name, false, emptyVal, newActivation);
    });

    it("setParam fails with non-empty value when delete", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, false, val1, activation)).to.be.revertedWith(NONEMPTY_VALUE);
    });

    it("setParam fails with non-owner", async function () {
      const { gp, nonOwner } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.connect(nonOwner).setParam(name, false, val1, activation)).to.be.revertedWith(NOT_OWNER);
    });
  });

  describe("setParamIn", function () {
    it("setParamIn success", async function () {
      const { gp } = fixture;
      const tx = await gp.setParamIn(name, true, val1, 1);
      await tx.wait();
      await expect(tx)
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, tx.blockNumber! + 1);
    });

    it("setParamIn now fails", async function () {
      const { gp } = fixture;
      await expect(gp.setParamIn(name, true, val1, 0)).to.be.revertedWith(ALREADY_PAST);
    });

    it("setParamIn fails with non-owner", async function () {
      const { gp, nonOwner } = fixture;
      await expect(gp.connect(nonOwner).setParamIn(name, false, val1, 1)).to.be.revertedWith(NOT_OWNER);
    });
  });

  describe("paramNames", function () {
    it("success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 1000;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      expect(await gp.paramNames(0)).to.equal(name);

      const name2 = "test22222";
      await expect(gp.setParam(name2, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name2, true, val1, activation);
      expect(await gp.paramNames(1)).to.equal(name2);

      // overwrite name2
      await expect(gp.setParam(name2, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name2, true, val1, activation);
      expect(await gp.paramNames(1)).to.equal(name2);

      await jumpBlock(1000);

      // add name2
      const newActivation = (await nowBlock()) + 1000;
      await expect(gp.setParam(name2, true, val2, newActivation))
        .to.emit(gp, "SetParam")
        .withArgs(name2, true, val2, newActivation);
      expect(await gp.paramNames(1)).to.equal(name2);
    });
  });

  describe("getAllParamNames", function () {
    it("getAllParamNames success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 1000;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      expect(await gp.getAllParamNames()).to.deep.equal([name]);

      const name2 = "test22222";
      await expect(gp.setParam(name2, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name2, true, val1, activation);
      expect(await gp.getAllParamNames()).to.deep.equal([name, name2]);
    });
  });

  describe("checkpoints", function () {
    it("checkpoints success", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 1000;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      expect(await gp.checkpoints(name)).to.deep.equal([defaultParam, [ethers.BigNumber.from(activation), true, val1]]);
    });
  });

  describe("getAllCheckpoints", function () {
    it("getAllCheckpoints success", async function () {
      const { gp } = fixture;
      const expected = [];
      const names = [name, "testtest", "test2222"];

      for (const [i, name] of names.entries()) {
        const p = [];
        p.push(defaultParam);
        for (let j = 0; j < 3; j++) {
          const activation = (await nowBlock()) + 10;
          const val = "0x" + (16 + 3 * i + j).toString(16);
          await expect(gp.setParam(name, true, val, activation))
            .to.emit(gp, "SetParam")
            .withArgs(name, true, val, activation);
          await jumpBlock(20);
          p.push([ethers.BigNumber.from(activation), true, val]);
        }
        expected.push(p);
      }
      expect(await gp.getAllCheckpoints()).to.deep.equal([names, expected]);
    });
  });

  describe("getParam", function () {
    it("getParam success when last checkpoint is activated", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      await jumpBlock(50);

      const p = await gp.getParam(name);
      expect(p).to.deep.equal([true, val1]);
    });

    it("getParam success when the last checkpoint is not activated", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);

      const p = await gp.getParam(name);
      expect(p).to.deep.equal([false, emptyVal]);
    });

    it("getParam returns false when name is not in paramNames", async function () {
      const { gp } = fixture;
      const p = await gp.getParam(name);
      expect(p).to.deep.equal([false, emptyVal]);
    });
  });

  describe("getParamAt", function () {
    it("getParamAt success on past blocks", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      await jumpBlock(50);

      const before = [false, emptyVal];
      const after = [true, val1];

      for (let i = activation - 5; i < activation + 5; i++) {
        const p = await gp.getParamAt(name, i);
        if (i < activation) {
          expect(p).to.deep.equal(before);
        } else {
          expect(p).to.deep.equal(after);
        }
      }
    });

    it("getParamAt success when the last checkpoint is not activated", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);

      const before = [false, emptyVal];
      const after = [true, val1];

      for (let i = 0; i < 20; i++) {
        await jumpBlock(1);
        const now = await nowBlock();
        const p = await gp.getParamAt(name, now);
        if (now < activation) {
          expect(p).to.deep.equal(before);
        } else {
          expect(p).to.deep.equal(after);
        }
      }
    });

    it("getParamAt success on future blocks", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      expect(await gp.getParamAt(name, activation - 1)).to.deep.equal([false, "0x"]);
      expect(await gp.getParamAt(name, activation)).to.deep.equal([true, val1]);
    });

    it("getParamAt returns false when name is not in paramNames", async function () {
      const { gp } = fixture;
      const p = await gp.getParamAt(name, await nowBlock());
      expect(p).to.deep.equal([false, emptyVal]);
    });

    it("getParamAt returns false when block is 0", async function () {
      const { gp } = fixture;
      const activation = (await nowBlock()) + 10;
      await expect(gp.setParam(name, true, val1, activation))
        .to.emit(gp, "SetParam")
        .withArgs(name, true, val1, activation);
      await jumpBlock(50);

      const p = await gp.getParamAt(name, 0);
      expect(p).to.deep.equal([false, emptyVal]);
    });
  });

  describe("getAllParams", function () {
    it("getAllParams success", async function () {
      const { gp } = fixture;
      const names = [name, "testtest", "test2222"];
      const activations = [];
      for (const name of names) {
        const activation = (await nowBlock()) + 10;
        const val = val1;
        await expect(gp.setParam(name, true, val, activation))
          .to.emit(gp, "SetParam")
          .withArgs(name, true, val, activation);
        activations.push(activation);
      }

      for (let i = 0; i < 10; i++) {
        await jumpBlock(1);
        const now = await nowBlock();
        const expectedNames = [];
        const expectedVals = [];
        for (const [j, name] of names.entries()) {
          const activation = activations[j];
          if (now >= activation) {
            expectedNames.push(name);
            expectedVals.push(val1);
          }
        }

        expect(await gp.getAllParams()).to.deep.equal([expectedNames, expectedVals]);
      }
    });
  });

  describe("getAllParamsAt", function () {
    it("getAllParamsAt success on past blocks", async function () {
      const { gp } = fixture;
      const names = [name, "testtest", "test2222"];
      const activations = [];
      for (const name of names) {
        const activation = (await nowBlock()) + 10;
        const val = val1;
        await expect(gp.setParam(name, true, val, activation))
          .to.emit(gp, "SetParam")
          .withArgs(name, true, val, activation);
        activations.push(activation);
      }
      await jumpBlock(50);

      for (let i = activations[0] - 5; i < activations[activations.length - 1] + 5; i++) {
        const expectedNames = [];
        const expectedVals = [];
        for (const [j, name] of names.entries()) {
          if (i >= activations[j]) {
            expectedNames.push(name);
            expectedVals.push(val1);
          }
        }
        expect(await gp.getAllParamsAt(i)).to.deep.equal([expectedNames, expectedVals]);
      }
    });

    it("getAllParamsAt success on future blocks", async function () {
      const { gp } = fixture;
      const names = [name, "testtest", "test2222"];
      const activations = [];
      for (const name of names) {
        const activation = (await nowBlock()) + 10;
        const val = val1;
        await expect(gp.setParam(name, true, val, activation))
          .to.emit(gp, "SetParam")
          .withArgs(name, true, val, activation);
        activations.push(activation);
      }

      for (let i = activations[0] - 5; i < activations[activations.length - 1] + 5; i++) {
        const expectedNames = [];
        const expectedVals = [];
        for (const [j, name] of names.entries()) {
          if (i >= activations[j]) {
            expectedNames.push(name);
            expectedVals.push(val1);
          }
        }
        expect(await gp.getAllParamsAt(i)).to.deep.equal([expectedNames, expectedVals]);
      }
    });
  });
});
