import { ethers, upgrades } from "hardhat";
import { deployAddressBook } from "../../common/helper";
import loadBls, { ModuleInstance } from "bls-signatures";

let BLS: ModuleInstance;

// can't use lib due to typechain dependency
export async function genBLSPubkeyPop(idx: number): Promise<string[]> {
  if (BLS == null) {
    BLS = await loadBls();
  }

  // prettier-ignore
  const seed = Uint8Array.from([
    idx, 50, 6, 244, 24, 199, 1, 25, 52, 88, 192, 19, 18, 12, 89, 6, 220, 18, 102, 58, 209, 82, 12, 62, 89, 110, 182, 9, 44, 20, 254, 22,
  ]);
  const sk = BLS.PopSchemeMPL.key_gen(seed);
  const pk = sk.get_g1().serialize();
  const pop = BLS.PopSchemeMPL.pop_prove(sk).serialize();
  // Uint8Array to string
  return [pk, pop].map((buf) => "0x" + Buffer.from(buf).toString("hex"));
}

export async function getActors() {
  const signers = await ethers.getSigners();
  const [deployer, abookAdmin, cn0, cn1, cn2, misc] = signers;
  const cnList = [cn0, cn1, cn2];
  return { deployer, abookAdmin, cn0, cn1, cn2, misc, cnList };
}

export async function deploySbrFixture() {
  const SbrFactory = await ethers.getContractFactory("SimpleBlsRegistry");
  const sbr = await upgrades.deployProxy(SbrFactory, [], { initializer: "initialize", kind: "uups" });
  const [pk, pop] = await genBLSPubkeyPop(0);

  const abook = await deployAddressBook("AddressBookMockThreeCN");

  return { abook, sbr, pk, pop };
}

export async function unregisterFixture() {
  const { sbr } = await deploySbrFixture();
  const pkList: string[] = [];
  const popList: string[] = [];
  const { cnList } = await getActors();
  const blsPubkeyInfoList: [string, string][] = [];

  for (const [i, cnNodeId] of cnList.entries()) {
    const [pk, pop] = await genBLSPubkeyPop(i);
    pkList.push(pk);
    popList.push(pop);
    blsPubkeyInfoList.push([pk, pop]);

    await sbr.register(cnNodeId.address, pk, pop);
  }

  const abook = await deployAddressBook("AddressBookMockOneCN");

  return { abook, sbr, cnList, pkList, popList, blsPubkeyInfoList };
}
