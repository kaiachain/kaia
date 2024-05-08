import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import {
  FuncIDV3,
  RequestState,
  WithdrawalState,
  augmentChai,
  nowTime,
  getBalance,
  setTime,
  toPeb,
  toBytes32,
  jumpTime,
  onlyAccessControlFail,
  notNullFailWithPoint,
  notConfirmedRequestFail,
  beforeInitFail,
  afterInitFail,
  checkRequestInfo,
  addTime,
} from "../common/helper";
import { ethers } from "hardhat";
import { ROLES, cnV3MultiSigPublicDelegationTestFixture, cnV3MultiSigUnitTestFixture } from "../common/fixtures";

const DAY = 24 * 60 * 60;
const WEEK = 7 * DAY;

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

/**
 * @dev Import unit test from cnStakingV2.test.ts to test multisig functions
 * Update:
 * 1. Update revert messages
 * 2. Add `submitToggleRedelegation` test
 * 3. The adminList contains contractValidator before initialization
 */
describe("CnStakingV3MultiSig", function () {
  describe("CnStakingV3MultiSig without public delegation", function () {
    let fixture: UnPromisify<ReturnType<typeof cnV3MultiSigUnitTestFixture>>;
    beforeEach(async function () {
      augmentChai();
      fixture = await loadFixture(cnV3MultiSigUnitTestFixture);
    });
    async function initialize() {
      const { contractValidator, adminList, stakingTrackerMockReceiver, cnStakingV3, gcId } = fixture;

      await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
        cnStakingV3,
        "UpdateStakingTracker",
      );

      await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");

      await cnStakingV3.connect(contractValidator).reviewInitialConditions();
      for (let i = 0; i < adminList.length; i++) {
        await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
      }

      await expect(cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) })).to.emit(
        cnStakingV3,
        "DepositLockupStakingAndInit",
      );
    }

    // 1. Test initializing process
    describe("CnStakingV3 Initialize", function () {
      async function initializeBeforeDeposit(opt: boolean) {
        const { contractValidator, adminList, stakingTrackerMockReceiver, cnStakingV3, gcId } = fixture;
        // Setup initialization
        if (opt) {
          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
            cnStakingV3,
            "UpdateStakingTracker",
          );
        }

        await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");

        await cnStakingV3.connect(contractValidator).reviewInitialConditions();
        for (let i = 0; i < adminList.length; i++) {
          await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
        }
      }
      it("Check constructor", async function () {
        const { contractValidator, adminList, nodeId, rewardAddr, requirement, cnStakingV3 } = fixture;

        expect(await cnStakingV3.contractValidator()).to.equal(contractValidator.address);
        expect(await cnStakingV3.isAdmin(contractValidator.address)).to.equal(true);
        expect(await cnStakingV3.adminList(0)).to.equal(contractValidator.address);
        expect(await cnStakingV3.adminList(1)).to.equal(adminList[0].address);
        expect(await cnStakingV3.adminList(2)).to.equal(adminList[1].address);
        expect(await cnStakingV3.adminList(3)).to.equal(adminList[2].address);
        expect(await cnStakingV3.requirement()).to.equal(requirement);
        expect(await cnStakingV3.nodeId()).to.equal(nodeId.address);
        expect(await cnStakingV3.rewardAddress()).to.equal(rewardAddr.address);
      });
      it("cnAdminList should be unique and not zero", async function () {
        const { contractValidator, adminList, nodeId, rewardAddr, requirement, unLockAmounts, unLockTimes } = fixture;

        await expect(
          ethers.deployContract("CnStakingV3MultiSig", [
            contractValidator.address,
            nodeId.address,
            rewardAddr.address,
            [adminList[0].address, adminList[0].address, adminList[1].address],
            requirement,
            unLockTimes,
            unLockAmounts,
          ]),
        ).to.be.revertedWith("Address is null or not unique.");
        await expect(
          ethers.deployContract("CnStakingV3MultiSig", [
            contractValidator.address,
            nodeId.address,
            rewardAddr.address,
            [adminList[0].address, adminList[1].address, ethers.constants.AddressZero],
            requirement,
            unLockTimes,
            unLockAmounts,
          ]),
        ).to.be.revertedWith("Address is null or not unique.");
      });

      it("Check constants", async function () {
        const { cnStakingV3, addressBook } = fixture;

        expect(await cnStakingV3.MAX_ADMIN()).to.equal(50);
        expect(await cnStakingV3.CONTRACT_TYPE()).to.equal("CnStakingContract");
        expect(await cnStakingV3.VERSION()).to.equal(3);
        expect(await cnStakingV3.ADDRESS_BOOK_ADDRESS()).to.equal(addressBook.address);
        expect(await cnStakingV3.STAKE_LOCKUP()).to.equal(WEEK);
      });

      describe("Check initializing process #setStakingTracker", function () {
        it("#setStakingTracker: Wrong msg.sender", async function () {
          const { other1, cnStakingV3, stakingTrackerMockReceiver } = fixture;

          const setStakingTrackerTx = cnStakingV3.connect(other1).setStakingTracker(stakingTrackerMockReceiver.address);
          await onlyAccessControlFail(setStakingTrackerTx, cnStakingV3);
        });
        it("#setStakingTracker: Tracker address can't be zero address", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(ethers.constants.AddressZero)).to.be
            .reverted;
        });
        it("#setStakingTracker: Tracker address can't be set after initialization", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await initializeBeforeDeposit(true);

          await expect(cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) })).to.emit(
            cnStakingV3,
            "DepositLockupStakingAndInit",
          );

          const setStakingTrackerTx = cnStakingV3.connect(adminList[0]).setStakingTracker(ethers.constants.AddressZero);
          await beforeInitFail(setStakingTrackerTx);
        });
        it("#setStakingTracker: Wrong tracker contract", async function () {
          const { adminList, cnStakingV3, stakingTrackerMockWrong } = fixture;

          await expect(
            cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockWrong.address),
          ).to.be.revertedWith("Invalid StakingTracker.");
        });
        it("#setStakingTracker: Successfully set staking tracker", async function () {
          const { adminList, cnStakingV3, stakingTrackerMockReceiver } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
            cnStakingV3,
            "UpdateStakingTracker",
          );

          expect(await cnStakingV3.stakingTracker()).to.equal(stakingTrackerMockReceiver.address);
        });
      });
      describe("Check initializing process #setGCId", function () {
        it("#setGCId: Wrong msg.sender", async function () {
          const { other1, cnStakingV3, gcId } = fixture;

          const setGCIdTx = cnStakingV3.connect(other1).setGCId(gcId);
          await onlyAccessControlFail(setGCIdTx, cnStakingV3);
        });
        it("#setGCId: GC ID can't be zero", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setGCId(0)).to.be.revertedWith("GC ID cannot be zero.");
        });
        it("#setGCId: GCId can't be set after initialization", async function () {
          const { adminList, cnStakingV3, gcId } = fixture;

          await initializeBeforeDeposit(true);

          await expect(cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) })).to.emit(
            cnStakingV3,
            "DepositLockupStakingAndInit",
          );

          const setGCIdTx = cnStakingV3.connect(adminList[0]).setGCId(gcId);
          await beforeInitFail(setGCIdTx);
        });
        it("#setGCId: Successfully set GCId", async function () {
          const { adminList, cnStakingV3, gcId } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");
          expect(await cnStakingV3.gcId()).to.equal(gcId);
        });
      });
      describe("Check initializing process #reviewInitialConditions", function () {
        it("#reviewInitialConditions: Wrong msg.sender", async function () {
          const { other1, cnStakingV3 } = fixture;

          const reviewInitialConditionsTx = cnStakingV3.connect(other1).reviewInitialConditions();
          await onlyAccessControlFail(reviewInitialConditionsTx, cnStakingV3);
        });
        it("#reviewInitialConditions: Admin can review only once", async function () {
          const { adminList, cnStakingV3 } = fixture;

          // First review
          await expect(cnStakingV3.connect(adminList[0]).reviewInitialConditions()).to.emit(
            cnStakingV3,
            "ReviewInitialConditions",
          );

          // Second review by same admin - fail
          await expect(cnStakingV3.connect(adminList[0]).reviewInitialConditions()).to.be.revertedWith(
            "Msg.sender already reviewed.",
          );
        });
        it("#reviewInitialConditions: Can't review after initialization", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await initializeBeforeDeposit(true);

          await expect(cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) })).to.emit(
            cnStakingV3,
            "DepositLockupStakingAndInit",
          );

          const reviewInitialConditionsTx = cnStakingV3.connect(adminList[0]).reviewInitialConditions();
          await beforeInitFail(reviewInitialConditionsTx);
        });
        it("#reviewInitialConditions: Successfully review initial conditions", async function () {
          const { contractValidator, adminList, cnStakingV3, stakingTrackerMockReceiver, gcId } = fixture;

          // Setup initialization
          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
            cnStakingV3,
            "UpdateStakingTracker",
          );

          await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");

          await expect(cnStakingV3.connect(contractValidator).reviewInitialConditions()).to.emit(
            cnStakingV3,
            "ReviewInitialConditions",
          );
          for (let i = 0; i < adminList.length - 1; i++) {
            await expect(cnStakingV3.connect(adminList[i]).reviewInitialConditions()).to.emit(
              cnStakingV3,
              "ReviewInitialConditions",
            );
          }

          await expect(cnStakingV3.connect(adminList[adminList.length - 1]).reviewInitialConditions()).to.emit(
            cnStakingV3,
            "CompleteReviewInitialConditions",
          );

          expect((await cnStakingV3.lockupConditions()).allReviewed).to.equal(true);
          expect((await cnStakingV3.lockupConditions()).reviewedCount).to.equal(4);

          // Get reviewers
          const reviewers = await cnStakingV3.getReviewers();
          expect(reviewers).to.equalAddrList([contractValidator.address, ...adminList.map((x) => x.address)]);
        });
      });
      describe("Check initializing process #depositLockupStakingAndInit", function () {
        it("#depositLockupStakingAndInit: Can't deposit before setting GC Id", async function () {
          const { contractValidator, adminList, cnStakingV3, stakingTrackerMockReceiver } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
            cnStakingV3,
            "UpdateStakingTracker",
          );

          await cnStakingV3.connect(contractValidator).reviewInitialConditions();
          for (let i = 0; i < adminList.length; i++) {
            await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
          }

          await expect(
            cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) }),
          ).to.be.revertedWith("Not set up properly.");
        });
        it("#depositLockupStakingAndInit: Can't deposit before setting staking tracker", async function () {
          const { contractValidator, adminList, cnStakingV3, gcId } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");

          await cnStakingV3.connect(contractValidator).reviewInitialConditions();
          for (let i = 0; i < adminList.length; i++) {
            await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
          }

          await expect(
            cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) }),
          ).to.be.revertedWith("Not set up properly.");
        });
        it("#depositLockupStakingAndInit: Can't deposit before allReviewed", async function () {
          const { adminList, cnStakingV3, stakingTrackerMockReceiver, gcId } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).setStakingTracker(stakingTrackerMockReceiver.address)).to.emit(
            cnStakingV3,
            "UpdateStakingTracker",
          );
          await expect(cnStakingV3.connect(adminList[0]).setGCId(gcId)).to.emit(cnStakingV3, "UpdateGCId");

          // contractValidator doesn't review the initial condition, which means allReviewed is false
          // await cnStakingV3.connect(contractValidator).reviewInitialConditions();
          for (let i = 0; i < adminList.length; i++) {
            await cnStakingV3.connect(adminList[i]).reviewInitialConditions();
          }

          await expect(
            cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) }),
          ).to.be.revertedWith("Reviews not finished.");
        });
        it("#depositLockupStakingAndInit: Msg.value should be equal to unlockAmount", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await initializeBeforeDeposit(true);

          await expect(
            cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(500n) }),
          ).to.be.revertedWith("Value does not match.");
        });
        it("#depositLockupStakingAndInit: Successfully deposit and initialize contract", async function () {
          const { adminList, unLockTimes, unLockAmounts, cnStakingV3, nodeId, rewardAddr, requirement } = fixture;

          await initializeBeforeDeposit(true);

          await expect(cnStakingV3.connect(adminList[0]).depositLockupStakingAndInit({ value: toPeb(600n) })).to.emit(
            cnStakingV3,
            "DepositLockupStakingAndInit",
          );

          expect(await cnStakingV3.initialLockupStaking()).to.equal(toPeb(600n));
          expect(await cnStakingV3.remainingLockupStaking()).to.equal(toPeb(600n));
          expect(await getBalance(cnStakingV3.address)).to.equal(toPeb(600n));

          const state = await cnStakingV3.getState();
          expect(state[0]).to.equal(ethers.constants.AddressZero);
          expect(state[1]).to.equal(nodeId.address);
          expect(state[2]).to.equal(rewardAddr.address);
          expect([...state[3]].sort()).to.equalAddrList(adminList.map((x) => x.address).sort());
          expect(state[4]).to.equal(requirement);
          expect(state[5]).to.equalNumberList(unLockTimes);
          expect(state[6]).to.equalNumberList(unLockAmounts);
          expect(state[7]).to.equal(true);
          expect(state[8]).to.equal(true);

          // Check roles
          for (let i = 0; i < adminList.length; i++) {
            expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, adminList[i].address)).to.equal(true);
            expect(await cnStakingV3.hasRole(ROLES.ADMIN_ROLE, adminList[i].address)).to.equal(true);
          }
          expect(await cnStakingV3.hasRole(ROLES.OPERATOR_ROLE, cnStakingV3.address)).to.equal(true);
          expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_APPROVER_ROLE, cnStakingV3.address)).to.equal(true);

          expect(await cnStakingV3.getRoleMemberCount(ROLES.OPERATOR_ROLE)).to.equal(1);
          expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_APPROVER_ROLE)).to.equal(1);
          expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_CLAIMER_ROLE)).to.equal(3);
          expect(await cnStakingV3.getRoleMemberCount(ROLES.ADMIN_ROLE)).to.equal(3);
          expect(await cnStakingV3.getRoleMemberCount(ROLES.STAKER_ROLE)).to.equal(0);
        });
      });
    });

    // 2. Test afterInit condition of all multisig functions
    describe("Check afterInit condition of all multisig functions", function () {
      it("#AddAdmin", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitAddAdminTx = cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address);
        await afterInitFail(submitAddAdminTx);
      });
      it("#DeleteAdmin", async function () {
        const { adminList, cnStakingV3 } = fixture;

        const submitDeleteAdminTx = cnStakingV3.connect(adminList[0]).submitDeleteAdmin(adminList[2].address);
        await afterInitFail(submitDeleteAdminTx);
      });
      it("#UpdateRequirement", async function () {
        const { adminList, cnStakingV3 } = fixture;

        const submitUpdateRequirement = cnStakingV3.connect(adminList[0]).submitUpdateRequirement(3);
        await afterInitFail(submitUpdateRequirement);
      });
      it("#ClearRequest", async function () {
        const { adminList, cnStakingV3 } = fixture;

        const submitClearRequestTx = cnStakingV3.connect(adminList[0]).submitClearRequest();
        await afterInitFail(submitClearRequestTx);
      });
      it("#UpdateStakingTracker", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitUpdateStakingTrackerTx = cnStakingV3
          .connect(adminList[0])
          .submitUpdateStakingTracker(other1.address);
        await afterInitFail(submitUpdateStakingTrackerTx);
      });
      it("#UpdateVoterAddress", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitUpdateVoterAddressTx = cnStakingV3.connect(adminList[0]).submitUpdateVoterAddress(other1.address);
        await afterInitFail(submitUpdateVoterAddressTx);
      });
      it("#UpdateRewardAddress", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitUpdateRewardAddressTx = cnStakingV3.connect(adminList[0]).submitUpdateRewardAddress(other1.address);
        await afterInitFail(submitUpdateRewardAddressTx);
      });
      it("#WithdrawLockupStaking", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitWithdrawLockupStakingTx = cnStakingV3
          .connect(adminList[0])
          .submitWithdrawLockupStaking(other1.address, toPeb(100n));
        await afterInitFail(submitWithdrawLockupStakingTx);
      });
      it("#ApproveStakingWithdrawal", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        const submitApproveStakingWithdrawalTx = cnStakingV3
          .connect(adminList[0])
          .submitApproveStakingWithdrawal(other1.address, toPeb(100n));
        await afterInitFail(submitApproveStakingWithdrawalTx);
      });
      it("#CancelApprovedStakingWithdrawal", async function () {
        const { adminList, cnStakingV3 } = fixture;

        const submitCancelApprovedStakingWithdrawalTx = cnStakingV3
          .connect(adminList[0])
          .submitCancelApprovedStakingWithdrawal(0);
        await afterInitFail(submitCancelApprovedStakingWithdrawalTx);
      });
    });

    // 3. Test submit/confirm functions not related to stake
    describe("Check multisig tx not related to stake", function () {
      // Now we can assume that the contract is initialized
      this.beforeEach(async () => {
        await initialize();
      });
      describe("Check addAdmin process", function () {
        it("#submitAddAdmin: Wrong msg.sender", async function () {
          const { cnStakingV3, other1, other2 } = fixture;

          const submitAddAdminTx = cnStakingV3.connect(other1).submitAddAdmin(other2.address);
          await onlyAccessControlFail(submitAddAdminTx, cnStakingV3);
        });
        it("#submitAddAdmin: Admin can't be zero address", async function () {
          const { adminList, cnStakingV3 } = fixture;

          const submitAddAdminTx = cnStakingV3.connect(adminList[0]).submitAddAdmin(ethers.constants.AddressZero);
          await notNullFailWithPoint(submitAddAdminTx);
        });
        it("#submitAddAdmin: Can't add admin that already exist", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(adminList[2].address)).to.be.revertedWith(
            "Admin already exists.",
          );
        });
        it("#submitAddAdmin::revokeConfirmation: Request proposer can cancel request", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Check request state
          const request = await cnStakingV3.getRequestInfo(0);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address],
              RequestState.NotConfirmed,
            ],
            request,
          );

          // Revoke request by proposer
          await expect(
            cnStakingV3
              .connect(adminList[0])
              .revokeConfirmation(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.emit(cnStakingV3, "CancelRequest");

          // Check updated request state
          const updatedRequest = await cnStakingV3.getRequestInfo(0);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address],
              RequestState.Canceled,
            ],
            updatedRequest,
          );
        });
        it("#submitAddAdmin::revokeConfirmation: Can revoke confirmation", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Update requirement to 3
          await cnStakingV3.connect(adminList[0]).submitUpdateRequirement(3);
          await cnStakingV3
            .connect(adminList[1])
            .confirmRequest(0, FuncIDV3.UpdateRequirement, toBytes32(3), toBytes32(0), toBytes32(0));

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(1, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.emit(cnStakingV3, "ConfirmRequest");

          // Revoke request by adminList[1]
          await expect(
            cnStakingV3
              .connect(adminList[1])
              .revokeConfirmation(1, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.emit(cnStakingV3, "RevokeConfirmation");

          // Check updated request state
          const updatedRequest = await cnStakingV3.getRequestInfo(1);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address],
              RequestState.NotConfirmed,
            ],
            updatedRequest,
          );
        });
        it("#addAdmin::revokeConfirmation: Can't cancel executed request", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Submit confirmRequest tx by adminList[1], and it will be confirmed since
          // current requirement is 2
          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // Can't cancel executed request
          const confirmRequestTx = cnStakingV3
            .connect(adminList[0])
            .revokeConfirmation(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0));
          await notConfirmedRequestFail(confirmRequestTx);
        });
        it("#submitAddAdmin::confirmRequest: Can't confirm a same request twice", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Can't confirm a same request twice
          await expect(
            cnStakingV3
              .connect(adminList[0])
              .confirmRequest(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.be.revertedWith("Msg.sender already confirmed.");
        });
        it("#submitAddAdmin::confirmRequest: Can't confirm request with wrong args", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Can't confirm request with wrong args
          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.DeleteAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.be.revertedWith("Function id and arguments do not match.");
        });
        it("#submitAddAdmin::confirmRequest: Can't confirm canceled request", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Revoke request by proposer
          await expect(
            cnStakingV3
              .connect(adminList[0])
              .revokeConfirmation(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          ).to.emit(cnStakingV3, "CancelRequest");

          // Can't confirm canceled request
          const confirmRequestTx = cnStakingV3
            .connect(adminList[1])
            .confirmRequest(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0));

          await notConfirmedRequestFail(confirmRequestTx);
        });
        it("#addAdmin: Successfully add admin", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          // Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Submit confirmRequest tx by adminList[1], and it will be confirmed since
          // current requirement is 2
          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // Check updated adminList
          const updatedAdminList = (await cnStakingV3.getState())[3];
          expect([...updatedAdminList].sort()).to.equalAddrList(
            [...adminList.map((x) => x.address), other1.address].sort(),
          );

          // Check updated UNSTAKING_CLAIMER_ROLE list.
          expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, other1.address)).to.equal(true);

          // Check updated request state
          const updatedRequest = await cnStakingV3.getRequestInfo(0);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address, adminList[1].address],
              RequestState.Executed,
            ],
            updatedRequest,
          );

          // Can't confirm executed request
          const confirmRequestTx = cnStakingV3
            .connect(adminList[2])
            .confirmRequest(0, FuncIDV3.AddAdmin, toBytes32(other1.address), toBytes32(0), toBytes32(0));
          await notConfirmedRequestFail(confirmRequestTx);
        });
      });
      describe("Check deleteAdmin process", function () {
        it("#submitDeleteAdmin: Wrong msg.sender", async function () {
          const { cnStakingV3, other1, other2 } = fixture;

          const submitDeleteAdminTx = cnStakingV3.connect(other1).submitDeleteAdmin(other2.address);
          await onlyAccessControlFail(submitDeleteAdminTx, cnStakingV3);
        });
        it("#submitDeleteAdmin: Admin can't be zero address", async function () {
          const { adminList, cnStakingV3 } = fixture;

          const submitDeleteAdminTx = cnStakingV3.connect(adminList[0]).submitDeleteAdmin(ethers.constants.AddressZero);
          await notNullFailWithPoint(submitDeleteAdminTx);
        });
        it("#submitDeleteAdmin: Can't delete admin that doesn't exist", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitDeleteAdmin(other1.address)).to.be.revertedWith(
            "Admin does not exist.",
          );
        });
        it("#deleteAdmin: Successfully delete admin", async function () {
          const { adminList, cnStakingV3 } = fixture;

          // Submit submitDeleteAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitDeleteAdmin(adminList[2].address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.DeleteAdmin, toBytes32(adminList[2].address), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // Check updated adminList
          const updatedAdminList = (await cnStakingV3.getState())[3] as string[];
          expect([...updatedAdminList].sort()).to.equalAddrList([adminList[0].address, adminList[1].address].sort());

          // Check updated UNSTAKING_CLAIMER_ROLE list.
          expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, adminList[2].address)).to.equal(false);
        });
      });
      describe("Check updateRequirement process", function () {
        it("#submitUpdateRequirement: Wrong msg.sender", async function () {
          const { cnStakingV3, other1 } = fixture;

          const submitUpdateRequirementTx = cnStakingV3.connect(other1).submitUpdateRequirement(3);
          await onlyAccessControlFail(submitUpdateRequirementTx, cnStakingV3);
        });
        it("#submitUpdateRequirement: Can't update requirement to same value", async function () {
          const { adminList, cnStakingV3 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitUpdateRequirement(2)).to.be.revertedWith(
            "Invalid value.",
          );
        });
        it("#updateRequirement: Successfully update requirement", async function () {
          const { adminList, cnStakingV3 } = fixture;

          // Submit submitUpdateRequirement tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitUpdateRequirement(3))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.UpdateRequirement, toBytes32(3), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // Check update requirement
          const updatedRequirement = (await cnStakingV3.getState())[4];
          expect(updatedRequirement).to.equal(3);
        });
      });
      describe("Check clearRequest process", function () {
        async function addRequestsForTest() {
          const { adminList, cnStakingV3, other1 } = fixture;

          // 1. Submit submitAddAdmin tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // 2. Submit updateRequirement tx by adminList[1]
          await expect(cnStakingV3.connect(adminList[1]).submitUpdateRequirement(3))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");
        }
        it("#submitClearRequest: Wrong msg.sender", async function () {
          const { cnStakingV3, other1 } = fixture;

          const clearRequestTx = cnStakingV3.connect(other1).submitClearRequest();
          await onlyAccessControlFail(clearRequestTx, cnStakingV3);
        });
        it("#clearRequest: Automatically clear outdated request when add/delete/requirement request has been executed", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          await addRequestsForTest();

          expect(await cnStakingV3.getRequestState(0)).to.equal(RequestState.NotConfirmed);
          expect(await cnStakingV3.getRequestState(1)).to.equal(RequestState.NotConfirmed);

          // Confirming #2 request results in clearing outdated request since it's about requirement update
          await expect(
            cnStakingV3
              .connect(adminList[0])
              .confirmRequest(1, FuncIDV3.UpdateRequirement, toBytes32(3), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // #1 Request state should be canceled
          const requestFor0 = await cnStakingV3.getRequestInfo(0);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address],
              // Update to Canceled state
              RequestState.Canceled,
            ],
            requestFor0,
          );

          expect(await cnStakingV3.getRequestState(1)).to.equal(RequestState.Executed);
        });
        it("#clearRequest: Manually clear outdated request", async function () {
          const { adminList, cnStakingV3, other1 } = fixture;

          await addRequestsForTest();

          // 3. Submit submitClearRequest tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitClearRequest())
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          // Confirming #3 request results in clearing outdated request
          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(2, FuncIDV3.ClearRequest, toBytes32(0), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");

          // Check request state of #0 and #1
          const requestFor0 = await cnStakingV3.getRequestInfo(0);
          checkRequestInfo(
            [
              FuncIDV3.AddAdmin,
              toBytes32(other1.address),
              toBytes32(0),
              toBytes32(0),
              adminList[0].address,
              [adminList[0].address],
              // Update to Canceled state
              RequestState.Canceled,
            ],
            requestFor0,
          );

          const requestFor1 = await cnStakingV3.getRequestInfo(1);
          checkRequestInfo(
            [
              FuncIDV3.UpdateRequirement,
              toBytes32(3),
              toBytes32(0),
              toBytes32(0),
              adminList[1].address,
              [adminList[1].address],
              // Update to Canceled state
              RequestState.Canceled,
            ],
            requestFor1,
          );
        });
      });
      describe("Check toggleRedelegation process", function () {
        it("#submitToggleRedelegation: Wrong msg.sender", async function () {
          const { cnStakingV3, other1 } = fixture;

          const toggleRedelegationTx = cnStakingV3.connect(other1).submitToggleRedelegation();
          await onlyAccessControlFail(toggleRedelegationTx, cnStakingV3);
        });
        it("#submitToggleRedelegation: can't toggle if not using public delegation", async function () {
          const { cnStakingV3, adminList } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitToggleRedelegation()).to.be.revertedWith(
            "Public delegation disabled.",
          );
        });
      });
      describe("Check updateStakingTracker process", function () {
        it("#submitUpdateStakingTracker: Wrong msg.sender", async function () {
          const { cnStakingV3, other1, stakingTrackerMockReceiver } = fixture;

          const updateStakingTrackerTx = cnStakingV3
            .connect(other1)
            .submitUpdateStakingTracker(stakingTrackerMockReceiver.address);
          await onlyAccessControlFail(updateStakingTrackerTx, cnStakingV3);
        });
        it("#submitUpdateStakingTracker: Staking tracker can't be zero address", async function () {
          const { cnStakingV3, adminList } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitUpdateStakingTracker(ethers.constants.AddressZero)).to.be
            .reverted;
        });
        it("#submitUpdateStakingTracker: Wrong staking tracker contract", async function () {
          const { cnStakingV3, adminList, stakingTrackerMockWrong } = fixture;

          await expect(
            cnStakingV3.connect(adminList[0]).submitUpdateStakingTracker(stakingTrackerMockWrong.address),
          ).to.be.revertedWith("Invalid StakingTracker.");
        });
        it("#updateStakingTracker: Can't update staking tracker if there's an active tracker", async function () {
          const { cnStakingV3, adminList, stakingTrackerMockReceiver, stakingTrackerMockActive } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitUpdateStakingTracker(stakingTrackerMockActive.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(
                0,
                FuncIDV3.UpdateStakingTracker,
                toBytes32(stakingTrackerMockActive.address),
                toBytes32(0),
                toBytes32(0),
              ),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          // Can't update stakingTrackerMockActive since it has live staking tracker
          await expect(
            cnStakingV3.connect(adminList[0]).submitUpdateStakingTracker(stakingTrackerMockReceiver.address),
          ).to.be.revertedWith("Cannot update tracker when there is an active tracker.");
        });
        it("#updateStakingTracker: Successfully update staking tracker", async function () {
          const { cnStakingV3, adminList, stakingTrackerMockReceiver, stakingTrackerMockActive } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitUpdateStakingTracker(stakingTrackerMockActive.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(
                0,
                FuncIDV3.UpdateStakingTracker,
                toBytes32(stakingTrackerMockActive.address),
                toBytes32(0),
                toBytes32(0),
              ),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          // Check updated staking tracker
          expect(await cnStakingV3.stakingTracker()).to.be.equal(stakingTrackerMockActive.address);
        });
      });
      describe("Check updateVoterAddress process", function () {
        it("#submitUpdateVoterAddress: Wrong msg.sender", async function () {
          const { cnStakingV3, other1 } = fixture;

          const updateVoterAddressTx = cnStakingV3.connect(other1).submitUpdateVoterAddress(other1.address);
          await onlyAccessControlFail(updateVoterAddressTx, cnStakingV3);
        });
        it("#updateVoterAddress: Successfully update voter address", async function () {
          const { cnStakingV3, stakingTrackerMockReceiver, adminList, other1 } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).submitUpdateVoterAddress(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.UpdateVoterAddress, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess")
            .to.emit(stakingTrackerMockReceiver, "RefreshVoter");
        });
      });
      describe("Check updateRewardAddress process", function () {
        async function updatePendingRewardAddress() {
          const { cnStakingV3, adminList, other1 } = fixture;

          // 1. Submit submitUpdateRewardAddress tx by adminList[0]
          await expect(cnStakingV3.connect(adminList[0]).submitUpdateRewardAddress(other1.address))
            .to.emit(cnStakingV3, "SubmitRequest")
            .to.emit(cnStakingV3, "ConfirmRequest");

          await expect(
            cnStakingV3
              .connect(adminList[1])
              .confirmRequest(0, FuncIDV3.UpdateRewardAddress, toBytes32(other1.address), toBytes32(0), toBytes32(0)),
          )
            .to.emit(cnStakingV3, "ConfirmRequest")
            .to.emit(cnStakingV3, "ExecuteRequestSuccess");
        }
        it("#submitUpdateRewardAddress: Wrong msg.sender", async function () {
          const { cnStakingV3, other1 } = fixture;

          const updateRewardAddressTx = cnStakingV3.connect(other1).submitUpdateRewardAddress(other1.address);
          await onlyAccessControlFail(updateRewardAddressTx, cnStakingV3);
        });
        it("#submitUpdateRewardAddress: Zero reward address", async function () {
          const { cnStakingV3, adminList } = fixture;

          const updateRewardAddressTx = cnStakingV3
            .connect(adminList[0])
            .submitUpdateRewardAddress(ethers.constants.AddressZero);
          await expect(updateRewardAddressTx).to.be.revertedWith("Address is null.");
        });
        it("#updateRewardAddress: Successfully update pending reward address", async function () {
          const { cnStakingV3, other1 } = fixture;

          await updatePendingRewardAddress();

          // Check pending reward address
          const pendingRewardAddress = await cnStakingV3.pendingRewardAddress();
          expect(pendingRewardAddress).to.be.equal(other1.address);
        });
        it("#acceptRewardAddress: Unauthorized address can't accept reward address", async function () {
          const { cnStakingV3, other1, other2 } = fixture;

          await updatePendingRewardAddress();

          // Accept reward address by non-authorized address
          await expect(cnStakingV3.connect(other2).acceptRewardAddress(other1.address)).to.be.revertedWith(
            "Unauthorized to accept reward address.",
          );
        });
        it("#acceptRewardAddress: Accept Reward address by abook admin", async function () {
          const { contractValidator, cnStakingV3, addressBook, other1 } = fixture;

          await updatePendingRewardAddress();

          // Accept reward address by abook admin
          // Note that contract validator is an abook admin
          await expect(cnStakingV3.connect(contractValidator).acceptRewardAddress(other1.address))
            .to.emit(addressBook, "ReviseRewardAddress")
            .to.emit(cnStakingV3, "AcceptRewardAddress");
        });
        it("#acceptRewardAddress: Accept Reward address by reward address", async function () {
          const { cnStakingV3, addressBook, other1 } = fixture;

          await updatePendingRewardAddress();

          // Accept reward address by reward address
          await expect(cnStakingV3.connect(other1).acceptRewardAddress(other1.address))
            .to.emit(addressBook, "ReviseRewardAddress")
            .to.emit(cnStakingV3, "AcceptRewardAddress");
        });
      });
    });

    // 4. Test lockup stakes
    describe("Check withdrawal process of lockup stakes (Initial lockup)", function () {
      async function withdraw(id: number, amount: bigint) {
        const { adminList, cnStakingV3, stakingTrackerMockReceiver, other1 } = fixture;

        await expect(cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, amount))
          .to.emit(cnStakingV3, "SubmitRequest")
          .to.emit(cnStakingV3, "ConfirmRequest");

        await expect(
          cnStakingV3
            .connect(adminList[1])
            .confirmRequest(
              id,
              FuncIDV3.WithdrawLockupStaking,
              toBytes32(other1.address),
              toBytes32(amount),
              toBytes32(0),
            ),
        )
          .to.emit(cnStakingV3, "ConfirmRequest")
          .to.emit(cnStakingV3, "ExecuteRequestSuccess")
          .to.emit(cnStakingV3, "WithdrawLockupStaking")
          .to.emit(stakingTrackerMockReceiver, "RefreshStake");
      }
      this.beforeEach(async () => {
        await initialize();
      });
      it("#submitWithdrawLockupStaking: Wrong msg.sender", async function () {
        const { cnStakingV3, other1 } = fixture;

        await jumpTime(105);

        const submitWithdrawLockupStakingTx = cnStakingV3
          .connect(other1)
          .submitWithdrawLockupStaking(other1.address, toPeb(100n));
        await onlyAccessControlFail(submitWithdrawLockupStakingTx, cnStakingV3);
      });
      it("#submitWithdrawLockupStaking: Receiver can't be zero address", async function () {
        const { adminList, cnStakingV3 } = fixture;

        await jumpTime(105);

        const submitWithdrawLockupStakingTx = cnStakingV3
          .connect(adminList[0])
          .submitWithdrawLockupStaking(ethers.constants.AddressZero, toPeb(100n));
        await notNullFailWithPoint(submitWithdrawLockupStakingTx);
      });
      it("#submitWithdrawLockupStaking: Value can't be zero", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        await jumpTime(105);

        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(0n)),
        ).to.be.revertedWith("Invalid value.");
      });
      it("#submitWithdrawLockupStaking: Not enough withdrawable amount", async function () {
        const { adminList, cnStakingV3, other1 } = fixture;

        // 1. now < unlockTime[0]: 0
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(100n)),
        ).to.be.revertedWith("Invalid value.");

        await jumpTime(105);

        // 2. unlockTime[0] < now < unlockTime[1]: unlockAmount[0] = 200 KAIA
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(300n)),
        ).to.be.revertedWith("Invalid value.");

        // 2-1. First withdraw: 100 KAIA
        await withdraw(0, BigInt(toPeb(100n)));

        // 2-2. Second withdraw: 150 KAIA => Not enough withdrawable amount
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(150n)),
        ).to.be.revertedWith("Invalid value.");

        await jumpTime(105);

        // 3. unlockTime[1] < now: unlockAmount[0] + unlockAmount[1] - 100 KAIA = 500 KAIA
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(700n)),
        ).to.be.revertedWith("Invalid value.");
      });
      it("#withdrawLockupStaking: Successfully withdraw lockup stakes", async function () {
        const { adminList, cnStakingV3, unLockTimes, unLockAmounts, other1 } = fixture;
        const balanceOfOther1Before = await ethers.provider.getBalance(other1.address);

        let lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[0]).to.equalNumberList(unLockTimes);
        expect(lockupStakingInfo[1]).to.equalNumberList(unLockAmounts);
        expect(lockupStakingInfo[2]).to.equal(BigInt(toPeb(600n)));
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(600n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(0n)));

        // 1. now < unlockTime[0]: 0
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(100n)),
        ).to.be.revertedWith("Invalid value.");

        await jumpTime(105);

        // 2. unlockTime[0] < now < unlockTime[1]: unlockAmount[0] = 200 KAIA
        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(200n)));

        await withdraw(0, BigInt(toPeb(100n)));

        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(500n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(100n)));

        // Remaining withdrawable amount: 100 KAIA
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(200n)),
        ).to.be.revertedWith("Invalid value.");

        // It's possible to withdraw 100 KAIA
        await withdraw(1, BigInt(toPeb(100n)));

        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(400n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(0n)));

        await jumpTime(105);

        // 3. unlockTime[1] < now: unlockAmount[0] + unlockAmount[1] - 200 KAIA = 400 KAIA
        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(400n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(400n)));

        await withdraw(2, BigInt(toPeb(250n)));

        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(150n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(150n)));

        // Remaining withdrawable amount: 150 KAIA
        await expect(
          cnStakingV3.connect(adminList[0]).submitWithdrawLockupStaking(other1.address, toPeb(200n)),
        ).to.be.revertedWith("Invalid value.");

        // It's possible to withdraw 150 KAIA
        await withdraw(3, BigInt(toPeb(150n)));

        lockupStakingInfo = await cnStakingV3.getLockupStakingInfo();
        expect(lockupStakingInfo[3]).to.equal(BigInt(toPeb(0n)));
        expect(lockupStakingInfo[4]).to.equal(BigInt(toPeb(0n)));

        // Check balance of the contract and receiver
        const balanceOfOther1After = await ethers.provider.getBalance(other1.address);

        expect(await ethers.provider.getBalance(cnStakingV3.address)).to.be.equal(0);
        expect(balanceOfOther1After.sub(balanceOfOther1Before)).to.be.equal(BigInt(toPeb(600n)));

        expect(await cnStakingV3.getRequestIds(0, 5, 2)).to.deep.equal([0, 1, 2, 3]);
        expect(await cnStakingV3.getRequestIds(0, 0, 2)).to.deep.equal([0, 1, 2, 3]);
      });
    });

    // 5. Test free stakes
    describe("Check free stakes process", function () {
      async function approveStakingWithdrawal(id: number, amount: bigint) {
        const { cnStakingV3, stakingTrackerMockReceiver, adminList, other1 } = fixture;

        await expect(cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(other1.address, amount))
          .to.emit(cnStakingV3, "SubmitRequest")
          .to.emit(cnStakingV3, "ConfirmRequest");

        const now = await nowTime();
        await setTime(now + 1);

        await expect(
          cnStakingV3
            .connect(adminList[1])
            .confirmRequest(
              id,
              FuncIDV3.ApproveStakingWithdrawal,
              toBytes32(other1.address),
              toBytes32(amount),
              toBytes32(0),
            ),
        )
          .to.emit(cnStakingV3, "ConfirmRequest")
          .to.emit(cnStakingV3, "ExecuteRequestSuccess")
          .to.emit(stakingTrackerMockReceiver, "RefreshStake");

        // Since confirmRequest executed after 1 second, the timestamp of the request is now + 2
        return now + 2;
      }
      this.beforeEach(async () => {
        await initialize();
      });
      describe("Check deposit process", function () {
        it("#delegate: Can't stake 0 KAIA", async function () {
          const { cnStakingV3, adminList } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).delegate({ value: 0 })).to.be.revertedWith("Invalid amount.");
        });
        it("#delegate: Successfully stake KAIA via delegate function", async function () {
          const { cnStakingV3, stakingTrackerMockReceiver, adminList } = fixture;

          await expect(cnStakingV3.connect(adminList[0]).delegate({ value: toPeb(500n) }))
            .to.emit(cnStakingV3, "DelegateKaia")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          const balanceOfContract = await ethers.provider.getBalance(cnStakingV3.address);
          expect(balanceOfContract).to.be.equal(BigInt(toPeb(1100n)));
        });
        it("#fallback: Can't stake 0 KAIA", async function () {
          const { cnStakingV3, adminList } = fixture;

          await expect(adminList[0].sendTransaction({ to: cnStakingV3.address, value: toPeb(0n) })).to.be.revertedWith(
            "Invalid amount.",
          );
        });
        it("#fallback: Successfully stake KAIA via fallback", async function () {
          const { cnStakingV3, stakingTrackerMockReceiver, adminList } = fixture;

          await expect(adminList[0].sendTransaction({ to: cnStakingV3.address, value: toPeb(500n) }))
            .to.emit(cnStakingV3, "DelegateKaia")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");

          const balanceOfContract = await ethers.provider.getBalance(cnStakingV3.address);
          expect(balanceOfContract).to.be.equal(BigInt(toPeb(1100n)));
        });
      });
      describe("Check withdrawal process", function () {
        this.beforeEach(async () => {
          const { adminList, cnStakingV3, stakingTrackerMockReceiver } = fixture;

          // Free stake 500 KAIA
          await expect(cnStakingV3.connect(adminList[0]).delegate({ value: toPeb(500n) }))
            .to.emit(cnStakingV3, "DelegateKaia")
            .to.emit(stakingTrackerMockReceiver, "RefreshStake");
        });
        describe("Check approveStakingWithdrawal process", function () {
          it("#submitApproveStakingWithdrawal: Wrong msg.sender", async function () {
            const { cnStakingV3, other1 } = fixture;

            const submitApproveStakingWithdrawalTx = cnStakingV3
              .connect(other1)
              .submitApproveStakingWithdrawal(other1.address, toPeb(100n));
            await onlyAccessControlFail(submitApproveStakingWithdrawalTx, cnStakingV3);
          });
          it("#submitApproveStakingWithdrawal: Receiver can't be zero address", async function () {
            const { cnStakingV3, adminList } = fixture;

            const submitApproveStakingWithdrawalTx = cnStakingV3
              .connect(adminList[0])
              .submitApproveStakingWithdrawal(ethers.constants.AddressZero, toPeb(100n));
            await notNullFailWithPoint(submitApproveStakingWithdrawalTx);
          });
          it("#submitApproveStakingWithdrawal: Value can't be zero", async function () {
            const { cnStakingV3, adminList, other1 } = fixture;

            await expect(
              cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(other1.address, toPeb(0n)),
            ).to.be.revertedWith("Invalid value.");
          });
          it("#submitApproveStakingWithdrawal: Not enough stakes to withdraw", async function () {
            const { cnStakingV3, adminList, other1 } = fixture;

            await expect(
              cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(other1.address, toPeb(550n)),
            ).to.be.revertedWith("Invalid value.");
          });
          it("#submitApproveStakingWithdrawal: unstaking + value can't exceed staking amount", async function () {
            const { cnStakingV3, adminList, other1 } = fixture;

            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            // Now, we have unstaked 200 KAIA. So, we can't approve more than 300 KAIA
            await expect(
              cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(other1.address, toPeb(350n)),
            ).to.be.revertedWith("Invalid value.");

            // It's possible below 300 KAIA
            await expect(cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(other1.address, toPeb(200n)))
              .to.emit(cnStakingV3, "SubmitRequest")
              .to.emit(cnStakingV3, "ConfirmRequest");
          });
          it("#approveStakingWithdrawal: Not withdrawable before 1 week", async function () {
            const { cnStakingV3, adminList } = fixture;

            // Approve withdraw 300 KAIA
            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(4 * DAY);

            // Not withdrawable before 1 week
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0)).to.be.revertedWith(
              "Not withdrawable yet.",
            );
          });
          it("#approveStakingWithdrawal: Need to withdraw free stakes before 2 weeks", async function () {
            const { cnStakingV3, adminList, stakingTrackerMockReceiver } = fixture;

            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(2 * WEEK);

            await approveStakingWithdrawal(1, BigInt(toPeb(200n)));

            // Admin can't withdraw stakes after withdrawableFrom + WEEK (2 weeks after approve)
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0))
              .to.emit(cnStakingV3, "CancelApprovedStakingWithdrawal")
              .to.emit(stakingTrackerMockReceiver, "RefreshStake");

            // Check staking related state
            const balanceOfContract = await ethers.provider.getBalance(cnStakingV3.address);
            const staking = await cnStakingV3.staking();
            const unstaking = await cnStakingV3.unstaking();

            expect(balanceOfContract).to.equal(BigInt(toPeb(1100n)));
            expect(staking).to.equal(BigInt(toPeb(500n)));
            expect(unstaking).to.equal(BigInt(toPeb(200n)));
          });
          it("#approveStakingWithdrawal: Check approved staking withdrawal info", async function () {
            const { cnStakingV3, other1 } = fixture;

            const now1 = await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(1000);

            const now2 = await approveStakingWithdrawal(1, BigInt(toPeb(200n)));

            // Check approved staking withdrawal info
            let approvedStakingWithdrawalInfo = await cnStakingV3.getApprovedStakingWithdrawalInfo(0);
            expect(approvedStakingWithdrawalInfo[0]).to.be.equal(other1.address);
            expect(approvedStakingWithdrawalInfo[1]).to.be.equal(BigInt(toPeb(300n)));
            expect(Number(approvedStakingWithdrawalInfo[2])).to.be.equal(now1 + WEEK);
            expect(Number(approvedStakingWithdrawalInfo[3])).to.be.equal(RequestState.Unknown);

            approvedStakingWithdrawalInfo = await cnStakingV3.getApprovedStakingWithdrawalInfo(1);
            expect(approvedStakingWithdrawalInfo[0]).to.be.equal(other1.address);
            expect(approvedStakingWithdrawalInfo[1]).to.be.equal(BigInt(toPeb(200n)));
            expect(Number(approvedStakingWithdrawalInfo[2])).to.be.equal(now2 + WEEK);
            expect(Number(approvedStakingWithdrawalInfo[3])).to.be.equal(RequestState.Unknown);
          });
          it("#withdrawApprovedStaking: Can't withdraw twice from a same request", async function () {
            const { cnStakingV3, adminList, stakingTrackerMockReceiver } = fixture;

            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(WEEK);

            // Admin can withdraw first withdraw stake request
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0))
              .to.emit(cnStakingV3, "WithdrawApprovedStaking")
              .to.emit(stakingTrackerMockReceiver, "RefreshStake");

            // Admin already withdrew first withdraw stake request
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0)).to.be.revertedWith(
              "Invalid state.",
            );
          });
          it("#withdrawApprovedStaking: Successfully withdraw all free stakes", async function () {
            const { cnStakingV3, adminList, stakingTrackerMockReceiver } = fixture;

            // Let's withdraw all funds through 2 requests
            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(3 * DAY);

            await approveStakingWithdrawal(1, BigInt(toPeb(200n)));

            await jumpTime(4 * DAY);

            // Admin can withdraw first withdraw stake request
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0))
              .to.emit(cnStakingV3, "WithdrawApprovedStaking")
              .to.emit(stakingTrackerMockReceiver, "RefreshStake");

            await jumpTime(4 * DAY);

            // Admin can withdraw third withdraw stake request
            await expect(cnStakingV3.connect(adminList[2]).withdrawApprovedStaking(1))
              .to.emit(cnStakingV3, "WithdrawApprovedStaking")
              .to.emit(stakingTrackerMockReceiver, "RefreshStake");

            // Check staking related state
            const unstaking = await cnStakingV3.unstaking();
            const staking = await cnStakingV3.staking();

            expect(unstaking).to.equal(BigInt(toPeb(0n)));
            expect(staking).to.equal(BigInt(toPeb(0n)));
          });
        });
        describe("Check cancelApprovedStakingWithdrawal process", function () {
          this.beforeEach(async () => {
            await approveStakingWithdrawal(0, BigInt(toPeb(300n)));

            await jumpTime(1000);

            await approveStakingWithdrawal(1, BigInt(toPeb(200n)));
          });
          it("#submitCancelApprovedStakingWithdrawal: Wrong msg.sender", async function () {
            const { cnStakingV3, other1 } = fixture;

            const submitCancelApprovedStakingWithdrawalTx = cnStakingV3
              .connect(other1)
              .submitCancelApprovedStakingWithdrawal(0);
            await onlyAccessControlFail(submitCancelApprovedStakingWithdrawalTx, cnStakingV3);
          });
          it("#submitCancelApprovedStakingWithdrawal: Can't cancel empty request", async function () {
            const { cnStakingV3, adminList } = fixture;

            await expect(cnStakingV3.connect(adminList[0]).submitCancelApprovedStakingWithdrawal(3)).to.be.revertedWith(
              "Withdrawal request does not exist.",
            );
          });
          it("#submitCancelApprovedStakingWithdrawal: Can't cancel transferred state", async function () {
            const { cnStakingV3, adminList, stakingTrackerMockReceiver } = fixture;

            await jumpTime(WEEK + DAY);

            // Withdraw first request
            await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0))
              .to.emit(cnStakingV3, "WithdrawApprovedStaking")
              .to.emit(stakingTrackerMockReceiver, "RefreshStake");

            // Check balance of contract
            expect(await ethers.provider.getBalance(cnStakingV3.address)).to.be.equal(BigInt(toPeb(800n)));

            // Can't cancel first request since it's already transferred
            await expect(cnStakingV3.connect(adminList[0]).submitCancelApprovedStakingWithdrawal(0)).to.be.revertedWith(
              "Invalid state.",
            );
          });
          it("#submitCancelApprovedStakingWithdrawal: Successfully cancel approved staking withdrawal request", async function () {
            const { cnStakingV3, adminList } = fixture;

            await expect(cnStakingV3.connect(adminList[0]).submitCancelApprovedStakingWithdrawal(0))
              .to.emit(cnStakingV3, "SubmitRequest")
              .to.emit(cnStakingV3, "ConfirmRequest");

            // Cancel first approved staking withdrawal request
            await expect(
              cnStakingV3
                .connect(adminList[1])
                .confirmRequest(2, FuncIDV3.CancelApprovedStakingWithdrawal, toBytes32(0), toBytes32(0), toBytes32(0)),
            )
              .to.emit(cnStakingV3, "ConfirmRequest")
              .to.emit(cnStakingV3, "ExecuteRequestSuccess");

            // Check withdrawal staking state and total unstaking amount
            const approvedStakingWithdrawalInfo = await cnStakingV3.getApprovedStakingWithdrawalInfo(0);
            expect(approvedStakingWithdrawalInfo[3]).to.equal(WithdrawalState.Canceled);

            const unstaking = await cnStakingV3.unstaking();
            expect(unstaking).to.equal(BigInt(toPeb(200n)));
          });
        });
      });
    });
  });
  describe("CnStakingV3MultiSig with public delegation", function () {
    let fixture: UnPromisify<ReturnType<typeof cnV3MultiSigPublicDelegationTestFixture>>;
    beforeEach(async function () {
      augmentChai();
      fixture = await loadFixture(cnV3MultiSigPublicDelegationTestFixture);

      const { pd1, adminList } = fixture;

      await pd1.connect(adminList[0]).stake({ value: toPeb(100n) });
    });
    describe("Check roles", function () {
      it("#check roles", async function () {
        const { cnStakingV3, adminList, pd1 } = fixture;
        // Check roles
        for (let i = 0; i < adminList.length; i++) {
          expect(await cnStakingV3.hasRole(ROLES.ADMIN_ROLE, adminList[i].address)).to.equal(true);
        }
        expect(await cnStakingV3.hasRole(ROLES.OPERATOR_ROLE, cnStakingV3.address)).to.equal(true);
        expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_APPROVER_ROLE, pd1.address)).to.equal(true);
        expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, pd1.address)).to.equal(true);
        expect(await cnStakingV3.hasRole(ROLES.STAKER_ROLE, pd1.address)).to.equal(true);

        expect(await cnStakingV3.getRoleMemberCount(ROLES.OPERATOR_ROLE)).to.equal(1);
        expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_APPROVER_ROLE)).to.equal(1);
        expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_CLAIMER_ROLE)).to.equal(1);
        expect(await cnStakingV3.getRoleMemberCount(ROLES.ADMIN_ROLE)).to.equal(3);
        expect(await cnStakingV3.getRoleMemberCount(ROLES.STAKER_ROLE)).to.equal(1);
      });
    });
    describe("Check withdrawal request", function () {
      it("#submitApprovedStakingWithdrawal: Can't submit approved staking withdrawal request if public delegation is enabled", async function () {
        const { cnStakingV3, adminList } = fixture;

        await expect(
          cnStakingV3.connect(adminList[0]).submitApproveStakingWithdrawal(adminList[0].address, toPeb(100n)),
        ).to.be.revertedWith("Public delegation enabled.");
      });
      it("#withdrawApprovedStaking: Can't withdraw approved staking withdrawal request if public delegation is enabled", async function () {
        const { cnStakingV3, pd1, adminList } = fixture;

        await pd1.connect(adminList[0]).redeem(adminList[0].address, toPeb(100n));

        await addTime(WEEK + 10);

        await expect(cnStakingV3.connect(adminList[0]).withdrawApprovedStaking(0)).to.be.revertedWithCustomError(
          cnStakingV3,
          "AccessControlUnauthorizedAccount",
        );
      });
      it("#cancelApprovedStakingWithdrawal: Can't cancel approved staking withdrawal request if public delegation is enabled", async function () {
        const { cnStakingV3, adminList } = fixture;

        await expect(cnStakingV3.connect(adminList[0]).submitCancelApprovedStakingWithdrawal(0)).to.be.revertedWith(
          "Public delegation enabled.",
        );
      });
    });
    describe("Check add/delete admin process", function () {
      it("#addAdmin: Doesn't add admin to unstaking claimer if public delegation is enabled", async function () {
        const { cnStakingV3, adminList, other1 } = fixture;

        await expect(cnStakingV3.connect(adminList[0]).submitAddAdmin(other1.address))
          .to.emit(cnStakingV3, "SubmitRequest")
          .to.emit(cnStakingV3, "ConfirmRequest");

        expect(await cnStakingV3.hasRole(ROLES.UNSTAKING_CLAIMER_ROLE, other1.address)).to.be.equal(false);
      });
      it("#deleteAdmin: Doesn't delete admin from unstaking claimer if public delegation is enabled", async function () {
        const { cnStakingV3, adminList, pd1 } = fixture;

        expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_CLAIMER_ROLE)).to.equal(1);
        expect(await cnStakingV3.getRoleMember(ROLES.UNSTAKING_CLAIMER_ROLE, 0)).to.equal(pd1.address);

        await expect(cnStakingV3.connect(adminList[0]).submitDeleteAdmin(adminList[0].address))
          .to.emit(cnStakingV3, "SubmitRequest")
          .to.emit(cnStakingV3, "ConfirmRequest");

        // No changes in unstaking claimer
        expect(await cnStakingV3.getRoleMemberCount(ROLES.UNSTAKING_CLAIMER_ROLE)).to.equal(1);
        expect(await cnStakingV3.getRoleMember(ROLES.UNSTAKING_CLAIMER_ROLE, 0)).to.equal(pd1.address);
      });
    });
    describe("Check toggle redelegation process", function () {
      it("#submitToggleRedelegation: Wrong msg.sender", async function () {
        const { cnStakingV3, other1 } = fixture;

        const toggleRedelegationTx = cnStakingV3.connect(other1).submitToggleRedelegation();
        await onlyAccessControlFail(toggleRedelegationTx, cnStakingV3);
      });
      it("#submitToggleRedelegation: Successfully toggle redelegation", async function () {
        const { cnStakingV3, adminList } = fixture;

        await expect(cnStakingV3.connect(adminList[0]).submitToggleRedelegation())
          .to.emit(cnStakingV3, "SubmitRequest")
          .to.emit(cnStakingV3, "ConfirmRequest");

        await expect(
          cnStakingV3
            .connect(adminList[1])
            .confirmRequest(0, FuncIDV3.ToggleRedelegation, toBytes32(0), toBytes32(0), toBytes32(0)),
        )
          .to.emit(cnStakingV3, "ConfirmRequest")
          .to.emit(cnStakingV3, "ExecuteRequestSuccess");

        // Check updated redelegation state
        expect(await cnStakingV3.isRedelegationEnabled()).to.be.equal(true);
      });
    });
  });
});
