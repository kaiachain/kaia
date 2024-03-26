import { setCode } from "@nomicfoundation/hardhat-network-helpers";
import { expect, Assertion } from "chai";
import { ethers } from "hardhat";
import { Transaction, Contract, BytesLike } from "ethers";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers.js";
import _ from "lodash";

export const ABOOK_ADDRESS = "0x0000000000000000000000000000000000000400";

export const FuncID = {
  Unknown: 0,
  AddAdmin: 1,
  DeleteAdmin: 2,
  UpdateRequirement: 3,
  ClearRequest: 4,
  WithdrawLockupStaking: 5,
  ApproveStakingWithdrawal: 6,
  CancelApprovedStakingWithdrawal: 7,
  UpdateRewardAddress: 8,
  UpdateStakingTracker: 9,
  UpdateVoterAddress: 10,
};
export const RequestState = {
  Unknown: 0,
  NotConfirmed: 1,
  Executed: 2,
  ExecutionFailed: 3,
  Canceled: 4,
};
export const WithdrawalState = {
  Unknown: 0,
  Transferred: 1,
  Canceled: 2,
};

export type CnStakingFixture = {
  contractValidator: SignerWithAddress;
  nodeId: SignerWithAddress;
  rewardAddr: SignerWithAddress;
  adminList: SignerWithAddress[];
  requirement: number;
  unLockTimes: number[];
  unLockAmounts: string[];
  stakingTrackerAddr?: string;
  gcId?: number;
};

export const VoteChoice = {
  No: 0,
  Yes: 1,
  Abstain: 2,
};

export async function submitRequest(
  cnStaking: Contract,
  funcId: number,
  proposer: SignerWithAddress,
  args: any[3]
) {
  switch (funcId) {
    case FuncID.AddAdmin:
      return await cnStaking.connect(proposer).submitAddAdmin(args[0]);
    case FuncID.DeleteAdmin:
      return await cnStaking.connect(proposer).submitDeleteAdmin(args[0]);
    case FuncID.UpdateRequirement:
      return await cnStaking.connect(proposer).submitUpdateRequirement(args[0]);
    case FuncID.ClearRequest:
      return await cnStaking.connect(proposer).submitClearRequest();
    case FuncID.WithdrawLockupStaking:
      return await cnStaking
        .connect(proposer)
        .submitWithdrawLockupStaking(args[0], args[1]);
    case FuncID.ApproveStakingWithdrawal:
      return await cnStaking
        .connect(proposer)
        .submitApproveStakingWithdrawal(args[0], args[1]);
    case FuncID.CancelApprovedStakingWithdrawal:
      return await cnStaking
        .connect(proposer)
        .submitCancelApprovedStakingWithdrawal(args[0]);
    case FuncID.UpdateRewardAddress:
      return await cnStaking
        .connect(proposer)
        .submitUpdateRewardAddress(args[0]);
    case FuncID.UpdateStakingTracker:
      return await cnStaking
        .connect(proposer)
        .submitUpdateStakingTracker(args[0]);
    case FuncID.UpdateVoterAddress:
      return await cnStaking
        .connect(proposer)
        .submitUpdateVoterAddress(args[0]);
    default:
      throw new Error("Invalid funcId");
  }
}

export async function confirmRequests(
  cnStaking: Contract,
  confirmers: SignerWithAddress[],
  args: any[5]
) {
  for (let i = 0; i < confirmers.length; i++) {
    await cnStaking
      .connect(confirmers[i])
      .confirmRequest(
        args[0],
        args[1],
        toBytes32(args[2]),
        toBytes32(args[3]),
        toBytes32(args[4])
      );
  }
}

export async function submitAndExecuteRequest(
  cnStaking: Contract,
  adminList: SignerWithAddress[],
  requirement: number,
  funcId: number,
  proposer: SignerWithAddress,
  args: any[5]
) {
  await submitRequest(cnStaking, funcId, proposer, args.slice(2, 5));
  if (requirement === 1) {
    return;
  }

  let cnt = 0;
  let i = 0;
  const confirmers = [] as SignerWithAddress[];
  while (cnt < requirement - 1) {
    if (adminList[i] !== proposer) {
      confirmers.push(adminList[i]);
      cnt++;
    }
    i++;
  }
  await confirmRequests(cnStaking, confirmers, args);
}

export async function deployAndInitStakingContract(
  version: number,
  args: CnStakingFixture
): Promise<Contract> {
  let stakingContract = {} as Contract;

  if (version === 1) {
    const cnStakingV1 = await ethers.deployContract("CnStakingContract", [
      args.contractValidator.address,
      args.nodeId.address,
      args.rewardAddr.address,
      args.adminList.map((admin) => admin.address),
      args.requirement,
      args.unLockTimes,
      args.unLockAmounts,
    ]);
    stakingContract = cnStakingV1;
  } else if (version === 2) {
    const cnStakingV2 = await ethers.deployContract("CnStakingV2", [
      args.contractValidator.address,
      args.nodeId.address,
      args.rewardAddr.address,
      args.adminList.map((admin) => admin.address),
      args.requirement,
      args.unLockTimes,
      args.unLockAmounts,
    ]);
    stakingContract = cnStakingV2;
  } else {
    throw new Error("Invalid version");
  }
  if (version == 2) {
    await stakingContract
      .connect(args.adminList[0])
      .setStakingTracker(args.stakingTrackerAddr);

    await stakingContract.connect(args.adminList[0]).setGCId(args.gcId);
  }

  await stakingContract
    .connect(args.contractValidator)
    .reviewInitialConditions();
  await stakingContract.connect(args.adminList[0]).reviewInitialConditions();

  await stakingContract
    .connect(args.adminList[0])
    .depositLockupStakingAndInit({ value: toPeb(6_000_000n) });

  return stakingContract;
}

export async function getRuntimecode(name: string): Promise<string> {
  const contract = await ethers.deployContract(name);
  return await ethers.provider.getCode(contract.address);
}

export async function setCodeAt(name: string, addr: string) {
  const runtimeCode = await getRuntimecode(name);
  await setCode(addr, runtimeCode);
}

export async function deployAddressBook(kind = "AddressBookMock") {
  const kinds = [
    "AddressBookMock",
    "AddressBookMockThreeCN",
    "AddressBookMockOneCN",
  ];
  for (const k of kinds) {
    if (k === kind) {
      await setCodeAt(k, ABOOK_ADDRESS);
      return await ethers.getContractAt(k, ABOOK_ADDRESS);
    }
  }
}

// Time related
export async function nowBlock() {
  return parseInt(await hre.network.provider.send("eth_blockNumber"));
}

export async function nowTime() {
  // hardhat simulated node has separate timer that increases every time
  // a new block is mined (which is basically every transaction).
  // Therefore nowTime() != Date.now().
  const block = await hre.network.provider.send("eth_getBlockByNumber", [
    "latest",
    false,
  ]);
  return parseInt(block.timestamp);
}

export async function setBlock(num: number) {
  const now = await nowBlock();
  if (now < num) {
    const blocksToMine = "0x" + (num - now).toString(16);
    await hre.network.provider.send("hardhat_mine", [blocksToMine]);
  }
}

export async function setTime(timestamp: number) {
  // https://ethereum.stackexchange.com/questions/86633/time-dependent-tests-with-hardhat
  await hre.network.provider.send("evm_mine", [timestamp]);
}

// Query chain
export async function getBalance(address: string) {
  const hex = await hre.network.provider.send("eth_getBalance", [address]);
  return ethers.BigNumber.from(hex).toString();
}

// Data conversion
export function toPeb(klay: bigint) {
  return ethers.utils.parseEther(klay.toString()).toString();
}

export function toKlay(peb: bigint) {
  return ethers.utils.formatEther(peb);
}

export function addPebs(a: bigint | string, b: bigint | string) {
  const bigA = ethers.BigNumber.from(a);
  const bigB = ethers.BigNumber.from(b);
  return bigA.add(bigB).toString();
}

export function subPebs(a: bigint | string, b: bigint | string) {
  const bigA = ethers.BigNumber.from(a);
  const bigB = ethers.BigNumber.from(b);
  return bigA.sub(bigB).toString();
}

export function toBytes32(x: any) {
  try {
    return ethers.utils.hexZeroPad(x, 32).toLowerCase();
    // eslint-disable-next-line no-empty
  } catch {}

  return x;
}

export function padUtils(value: string | BytesLike, length: number) {
  return ethers.utils.hexlify(
    ethers.utils.zeroPad(ethers.utils.hexlify(value), length)
  );
}

export function numericAddr(n: number, m: number) {
  // Return a human-friendly address to be used as placeholders.
  // ex. CN #42, second node ID is:
  // numericAddr(42, 2) => 0x4202000000000000000000000000000000000001
  const a = n < 10 ? "0" + n : "" + n;
  const b = m < 10 ? "0" + m : "" + m;
  return "0x" + a + b + "00".repeat(17) + "01";
}

export async function jumpTime(time: number, mine = true) {
  await ethers.provider.send("evm_increaseTime", [time]);
  if (mine) {
    await ethers.provider.send("evm_mine", []);
  }
}

export async function jumpBlock(num: number) {
  if (num > 0) {
    const blocksToMine = "0x" + num.toString(16);
    await hre.network.provider.send("hardhat_mine", [blocksToMine]);
  }
}

/* ========== TESTING HELPERS ========== */
// declare global to add custom properties in typescript
declare global {
  export namespace Chai {
    interface Assertion {
      equalAddrList(arr: string[]): void;
      equalStrList(arr: string[]): void;
      equalNumberList(arr: bigint[] | string[] | number[]): void;
      equalBooleanList(arr: boolean[]): void;
    }
  }
}

// Augment chai expect(..) assertion
// - .to.equal(..) with more generous type check
// - .to.emit(..) for CnStaking specific events
export function augmentChai() {
  Assertion.addMethod("equalAddrList", function (arr) {
    arr = _.map(arr, (elem) => elem.address || elem);
    const expected = _.map(arr, (elem) => elem.toLowerCase());
    const actual = _.map(this._obj, (elem) => elem.toLowerCase());
    return this.assert(
      _.isEqual(expected, actual),
      "expected #{this} to be equal to #{arr}",
      "expected #{this} to not equal to #{arr}",
      expected,
      actual
    );
  });
  Assertion.addMethod("equalStrList", function (arr) {
    const expected = _.map(arr, (elem) => elem.toLowerCase());
    const actual = _.map(this._obj, (elem) => elem.toLowerCase());
    return this.assert(
      _.isEqual(expected, actual),
      "expected #{this} to be equal to #{arr}",
      "expected #{this} to not equal to #{arr}",
      expected,
      actual
    );
  });
  Assertion.addMethod("equalNumberList", function (arr) {
    const expected = _.map(arr, (elem) => elem.toString());
    const actual = _.map(this._obj, (elem) => elem.toString());
    return this.assert(
      _.isEqual(expected, actual),
      "expected #{this} to be equal to #{arr}",
      "expected #{this} to not equal to #{arr}",
      expected,
      actual
    );
  });
  Assertion.addMethod("equalBooleanList", function (arr) {
    const expected = arr;
    const actual = this._obj;
    return this.assert(
      _.isEqual(expected, actual),
      "expected #{this} to be equal to #{arr}",
      "expected #{this} to not equal to #{arr}",
      expected,
      actual
    );
  });
}

// Modifier failure messages
export async function onlyAdminFail(tx: Transaction) {
  await expectRevert(tx, "Address is not admin.");
}
export async function notNullFail(tx: Transaction) {
  await expectRevert(tx, "Address is null");
}
export async function notConfirmedRequestFail(tx: Transaction) {
  await expectRevert(tx, "Must be at not-confirmed state.");
}
export async function beforeInitFail(tx: Transaction) {
  await expectRevert(tx, "Contract has been initialized.");
}
export async function afterInitFail(tx: Transaction) {
  await expectRevert(tx, "Contract is not initialized.");
}

export async function expectRevert(
  expr: Transaction,
  message: string | RegExp
) {
  return expect(expr).to.be.revertedWith(message);
}

export function checkRequestInfo(expected: any, returned: any) {
  expect(returned[0]).to.equal(expected[0]);
  expect(returned[1]).to.equal(expected[1]);
  expect(returned[2]).to.equal(expected[2]);
  expect(returned[3]).to.equal(expected[3]);
  expect(returned[4]).to.equal(expected[4]);
  expect(returned[5]).to.equalAddrList(expected[5]);
  expect(returned[6]).to.equal(expected[6]);
}
