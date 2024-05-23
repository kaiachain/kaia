import { expect } from "chai";

describe("", function () {
  const validAddrs = [
    "link1hpufl3l8g44aaz3qsqw886sjanhhu73ul6tllxuw3pqlhxzq9e4svku69h",
    "link10pgvx8pn5qwgwv066g93jlqux6mnsp0kajg9hp",
  ];
  const invalidAddrs = [
    "link1hpufl3l8g44aaz3qsqw886sjanhhu63ul6tllxuw3pqlhxzq9e4svku69h",
    "link10pgvx8pn5qwgwv066g93jlqux5mnsp0kajg9hp",
  ];
  let bech32;

  beforeEach(async function() {
    const bech32Factory = await ethers.getContractFactory("Bech32");
    bech32 = await upgrades.deployProxy(bech32Factory);
  })

  it("FNSA address validation", async function () {
    for (let validAddr of validAddrs) {
      expect(await (bech32.verifyAddrFNSA(validAddr.toLowerCase(), true))).to.be.equal(true);
      expect(await (bech32.verifyAddrFNSA(validAddr.toLowerCase(), false))).to.be.equal(true);
      expect(await (bech32.verifyAddrFNSA(validAddr.toUpperCase(), true))).to.be.equal(true);
    }

    for (let invalidAddr of invalidAddrs) {
      await expect(bech32.verifyAddrFNSA(invalidAddr.toLowerCase(), true))
        .to.be.revertedWith("Invalid checksum");
      await expect(bech32.verifyAddrFNSA(invalidAddr.toLowerCase(), false))
        .to.be.revertedWith("Invalid checksum");
      await expect(bech32.verifyAddrFNSA(invalidAddr.toUpperCase(), true))
        .to.be.revertedWith("Invalid checksum");
    }

    await expect(bech32.decodeNoLimit("0x1234", false))
      .to.be.revertedWith("invalid bech32 string length")
  })

  it("Input validation", async function () {
    await expect(bech32.normalize("0xAA"))
      .to.be.revertedWith("Not allowed ASCII value");

    await expect(bech32.normalize(Buffer.from("Link10pgvx8pn5qwgwv066g93jlqux6mnsp0kajg9hp")))
      .to.be.revertedWith("string not all lowercase or all uppercase");

    await expect(bech32.copyBytes("0xaa", 1, 0))
      .to.be.revertedWith("from must be less than to");

    await expect(bech32.copyBytes("0xaa", 0, 10))
      .to.be.revertedWith("to must be less than or equal source length");
  })
})
