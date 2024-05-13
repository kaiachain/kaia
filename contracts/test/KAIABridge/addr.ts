import { expect } from "chai";

describe("", function () {
  let bech32;

  beforeEach(async function() {
    const bech32Factory = await ethers.getContractFactory("Bech32");
    bech32 = await upgrades.deployProxy(bech32Factory);
  })

  it("FNSA address validation", async function () {
    const validAddrs = [
      "link1hpufl3l8g44aaz3qsqw886sjanhhu73ul6tllxuw3pqlhxzq9e4svku69h",
      "link10pgvx8pn5qwgwv066g93jlqux6mnsp0kajg9hp",
    ];

    const invalidAddrs = [
      "link1hpufl3l8g44aaz3qsqw886sjanhhu63ul6tllxuw3pqlhxzq9e4svku69h",
      "link10pgvx8pn5qwgwv066g93jlqux5mnsp0kajg9hp",
    ];

    for (let validAddr of validAddrs) {
      expect(await (bech32.verifyAddrFNSA(validAddr, false))).to.be.equal(true);
    }

    for (let invalidAddr of invalidAddrs) {
      await expect(bech32.verifyAddrFNSA(invalidAddr, false))
        .to.be.revertedWith("Invalid checksum");
    }
  })
})
