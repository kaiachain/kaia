const { expect } = require("chai");
const { ethers } = require("hardhat");

const NOT_OWNER = "Ownable: caller is not the owner";
const EMPTYNAME = "GovParam: name cannot be empty";
const EMPTY_VALUE = "GovParam: val must not be empty if exists=true";
const NONEMPTY_VALUE = "GovParam: val must be empty if exists=false";
const ALREADY_PAST = "GovParam: activation must be in the future";

async function getnow() {
    return parseInt(await hre.network.provider.send("eth_blockNumber"));
}

async function mineMoreBlocks(num) {
    mineblock = "0x" + num.toString(16);
    await hre.network.provider.send("hardhat_mine", [mineblock]);
    // due to bug, the next line is required (https://github.com/NomicFoundation/hardhat/issues/2467)
    await hre.network.provider.send("hardhat_setNextBlockBaseFeePerGas", ["0x0"]);
}

describe("GovParam", function() {
    let accounts;
    let nonOwner;
    let gp;

    const name = "istanbul.committeesize";
    const val1 = "0x1234";
    const val2 = "0x5678";
    const emptyVal = "0x";

    const defaultParam = [ethers.BigNumber.from("0"), false, emptyVal];

    beforeEach(async function() {
        accounts = await hre.ethers.getSigners();

        nonOwner = accounts[1];

        const GovParam = await ethers.getContractFactory("GovParam");
        gp = await GovParam.deploy();
    });

    describe("constructor", function() {
        it("Constructor success", async function() {
            const GovParam = await ethers.getContractFactory("GovParam");
            gp = await GovParam.connect(nonOwner).deploy();

            expect(await gp.owner()).to.equal(nonOwner.address);
        });
    });

    describe("setParam", function() {
        it("setParam success", async function() {
            const activation = await getnow() + 5;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);

            before = [false, emptyVal];
            after = [true, val1];

            for (let i = 0; i < 10; i++) {
                await mineMoreBlocks(1);
                p = await gp.getParam(name);
                now = await getnow();
                if (now < activation) {
                    expect(p).to.deep.equal(before);
                }
                else {
                    expect(p).to.deep.equal(after);
                }
            }
        });

        it("setParam overwrite success", async function() {
            const activation = await getnow() + 5;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);

            const newActivation = activation + 5;
            await expect(gp.setParam(name, true, val2, newActivation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val2, newActivation);

            before = [false, emptyVal];
            after = [true, val2];

            for (let i = 0; i < 20; i++) {
                await mineMoreBlocks(1);
                p = await gp.getParam(name);
                now = await getnow();
                if (now < newActivation) {
                    expect(p).to.deep.equal(before);
                }
                else {
                    expect(p).to.deep.equal(after);
                }
            }
        });

        it("setParam fails with empty name", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam("", true, val1, activation))
                .to.be.revertedWith(EMPTYNAME);
        });

        it("setParam fails with past activation", async function() {
            const activation = await getnow();
            await expect(gp.setParam(name, true, val1, activation))
                .to.be.revertedWith(ALREADY_PAST);
        });

        it("setParam fails with empty value except delete", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, emptyVal, activation))
                .to.be.revertedWith(EMPTY_VALUE);
        });

        it("setParam deleting success", async function() {
            let activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);

            await mineMoreBlocks(10);
            activation = await getnow() + 10;
            await expect(gp.setParam(name, false, emptyVal, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, false, emptyVal, activation);
        });

        it("setParam fails with non-empty value when delete", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, false, val1, activation))
                .to.be.revertedWith(NONEMPTY_VALUE);
        });

        it("setParam fails with non-owner", async function() {
            const activation = await getnow() + 10;
            await expect(gp.connect(nonOwner).setParam(name, false, val1, activation))
                .to.be.revertedWith(NOT_OWNER);
        });
    });

    describe("setParamIn", function() {
        it("setParamIn success", async function() {
            const tx = await gp.setParamIn(name, true, val1, 1);
            await expect(tx)
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, tx.blockNumber + 1);
        });

        it("setParamIn now fails", async function() {
            await expect(gp.setParamIn(name, true, val1, 0))
                .to.be.revertedWith(ALREADY_PAST);
        });

        it("setParamIn fails with non-owner", async function() {
            await expect(gp.connect(nonOwner).setParamIn(name, false, val1, 1))
                .to.be.revertedWith(NOT_OWNER);
        });
    });

    describe("paramNames", function() {
        it("success", async function() {
            const activation = await getnow() + 1000;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            expect(await gp.paramNames(0)).to.equal(name);

            const name2 = "test22222"
            await expect(gp.setParam(name2, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name2, true, val1, activation);
            expect(await gp.paramNames(1)).to.equal(name2);

            // overwrite name2
            await expect(gp.setParam(name2, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name2, true, val1, activation);
            expect(await gp.paramNames(1)).to.equal(name2);

            await mineMoreBlocks(1000);

            // add name2
            const newActivation = await getnow() + 1000;
            await expect(gp.setParam(name2, true, val2, newActivation))
                .to.emit(gp, "SetParam")
                .withArgs(name2, true, val2, newActivation);
            expect(await gp.paramNames(1)).to.equal(name2);
        });
    });

    describe("getAllParamNames", function() {
        it("getAllParamNames success", async function() {
            const activation = await getnow() + 1000;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            expect(await gp.getAllParamNames()).to.deep.equal([name]);

            const name2 = "test22222"
            await expect(gp.setParam(name2, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name2, true, val1, activation);
            expect(await gp.getAllParamNames()).to.deep.equal([name, name2]);
        });
    });

    describe("checkpoints", function() {
        it("checkpoints success", async function() {
            const activation = await getnow() + 1000;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            expect(await gp.checkpoints(name)).to.deep.equal([
                defaultParam,
                [ethers.BigNumber.from(activation), true, val1]
            ]);
        });
    });

    describe("getAllCheckpoints", function() {
        it("getAllCheckpoints success", async function() {
            var expected = [];
            const names = [name, "testtest", "test2222"];

            for (const [i, name] of names.entries()) {
                var p = [];
                p.push(defaultParam);
                for (let j = 0; j < 3; j++) {
                    const activation = await getnow() + 10;
                    const val = "0x" + (16 + 3 * i + j).toString(16);
                    await expect(gp.setParam(name, true, val, activation))
                        .to.emit(gp, "SetParam")
                        .withArgs(name, true, val, activation);
                    await mineMoreBlocks(20);
                    p.push([
                        ethers.BigNumber.from(activation),
                        true,
                        val,
                    ]);
                }
                expected.push(p);
            }
            expect(await gp.getAllCheckpoints()).to.deep.equal([
                names,
                expected
            ]);
        });
    });

    describe("getParam", function() {
        it("getParam success when last checkpoint is activated", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            await mineMoreBlocks(50);

            p = await gp.getParam(name);
            expect(p).to.deep.equal([true, val1]);
        });

        it("getParam success when the last checkpoint is not activated", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);

            p = await gp.getParam(name);
            expect(p).to.deep.equal([false, emptyVal]);
        });

        it("getParam returns false when name is not in paramNames", async function() {
            p = await gp.getParam(name);
            expect(p).to.deep.equal([false, emptyVal]);
        });
    });

    describe("getParamAt", function() {
        it("getParamAt success on past blocks", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            await mineMoreBlocks(50);

            before = [false, emptyVal];
            after = [true, val1];

            for (let i = activation - 5; i < activation + 5; i++) {
                p = await gp.getParamAt(name, i);
                if (i < activation) {
                    expect(p).to.deep.equal(before);
                }
                else {
                    expect(p).to.deep.equal(after);
                }
            }
        });

        it("getParamAt success when the last checkpoint is not activated", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);

            before = [false, emptyVal];
            after = [true, val1];

            for (let i = 0; i < 20; i++) {
                await mineMoreBlocks(1);
                now = await getnow();
                p = await gp.getParamAt(name, now);
                if (now < activation) {
                    expect(p).to.deep.equal(before);
                }
                else {
                    expect(p).to.deep.equal(after);
                }
            }
        });

        it("getParamAt success on future blocks", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            expect(await gp.getParamAt(name, activation-1))
                .to.deep.equal([false, "0x"]);
            expect(await gp.getParamAt(name, activation))
                .to.deep.equal([true, val1]);
        });

        it("getParamAt returns false when name is not in paramNames", async function() {
            p = await gp.getParamAt(name, await getnow());
            expect(p).to.deep.equal([false, emptyVal]);
        });

        it("getParamAt returns false when block is 0", async function() {
            const activation = await getnow() + 10;
            await expect(gp.setParam(name, true, val1, activation))
                .to.emit(gp, "SetParam")
                .withArgs(name, true, val1, activation);
            await mineMoreBlocks(50);

            p = await gp.getParamAt(name, 0);
            expect(p).to.deep.equal([false, emptyVal]);
        });
    });

    describe("getAllParams", function() {
        it("getAllParams success", async function() {
            const names = [name, "testtest", "test2222"];
            var activations = [];
            for (let name of names) {
                const activation = await getnow() + 10;
                const val = val1;
                await expect(gp.setParam(name, true, val, activation))
                    .to.emit(gp, "SetParam")
                    .withArgs(name, true, val, activation);
                activations.push(activation);
            }

            for (let i = 0; i < 10; i++) {
                await mineMoreBlocks(1);
                let now = await getnow();
                let expectedNames = [];
                let expectedVals = [];
                for (const [j, name] of names.entries()) {
                    let activation = activations[j];
                    if (now >= activation) {
                        expectedNames.push(name);
                        expectedVals.push(val1);
                    }
                }

                expect(await gp.getAllParams()).to.deep.equal([
                    expectedNames,
                    expectedVals
                ]);
            }
        });
    });

    describe("getAllParamsAt", function() {
        it("getAllParamsAt success on past blocks", async function() {
            const names = [name, "testtest", "test2222"];
            var activations = [];
            for (let name of names) {
                const activation = await getnow() + 10;
                const val = val1;
                await expect(gp.setParam(name, true, val, activation))
                    .to.emit(gp, "SetParam")
                    .withArgs(name, true, val, activation);
                activations.push(activation);
            }
            await mineMoreBlocks(50);

            for (let i = activations[0] - 5; i < activations[activations.length - 1] + 5; i++) {
                let expectedNames = [];
                let expectedVals = [];
                for (let [j, name] of names.entries()) {
                    if (i >= activations[j]) {
                        expectedNames.push(name);
                        expectedVals.push(val1);
                    }
                }
                expect(await gp.getAllParamsAt(i)).to.deep.equal([
                    expectedNames,
                    expectedVals,
                ]);
            }
        });

        it("getAllParamsAt success on future blocks", async function() {
            const names = [name, "testtest", "test2222"];
            var activations = [];
            for (let name of names) {
                const activation = await getnow() + 10;
                const val = val1;
                await expect(gp.setParam(name, true, val, activation))
                    .to.emit(gp, "SetParam")
                    .withArgs(name, true, val, activation);
                activations.push(activation);
            }

            for (let i = activations[0] - 5; i < activations[activations.length - 1] + 5; i++) {
                let expectedNames = [];
                let expectedVals = [];
                for (let [j, name] of names.entries()) {
                    if (i >= activations[j]) {
                        expectedNames.push(name);
                        expectedVals.push(val1);
                    }
                }
                expect(await gp.getAllParamsAt(i)).to.deep.equal([
                    expectedNames,
                    expectedVals,
                ]);
            }
        });
    });
});
