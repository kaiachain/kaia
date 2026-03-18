import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import {
  FuncID,
  augmentChai,
  toPeb,
  jumpTime,
  submitAndExecuteRequest,
  ABOOK_ADDRESS,
  setBlock,
  WEEK,
} from "../common/helper";
import { ethers } from "hardhat";
import { stakingTrackerV1TestFixture } from "../materials";

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

/**
 * @dev This unit & scenario test is for StakingTracker.sol
 * 1. Test initializing process
 *    - Check constants
 * 2. Test create/expire trackers
 *    - createTracker
 *    - refreshStake - expire outdated tracker if needed
 * 3. Test functions by calling directly
 *    - refreshStake when all trackers are live
 *    - refreshStake when there's expired tracker
 *    - refreshVoter
 * 4. Test functions by calling through CnStakingV2.sol
 *    - refreshStake
 * *    #stakeKlay
 *      #withdrawLockupStaking
 *      #approveStakingWithdrawal
 *      #cancelApprovedStakingWithdrawal
 *      #withdrawApprovedStaking
 *    - refreshVoter
 *      #updateVoterAddress
 */
describe("StakingTracker.sol", function () {
  let fixture: UnPromisify<ReturnType<typeof stakingTrackerV1TestFixture>>;
  beforeEach(async function () {
    augmentChai();
    fixture = await loadFixture(stakingTrackerV1TestFixture);
  });

  // 1. Test initializing process
  describe("StakingTracker Initialize", function () {
    it("Check staking tracker constants", async function () {
      const { stakingTracker } = fixture;

      const contractType = await stakingTracker.CONTRACT_TYPE();
      const version = await stakingTracker.VERSION();
      const aBookAddr = await stakingTracker.ADDRESS_BOOK_ADDRESS();
      const minStake = await stakingTracker.MIN_STAKE();

      expect(contractType).to.equal("StakingTracker");
      expect(version).to.equal(1);
      expect(aBookAddr).to.equal(ABOOK_ADDRESS);
      expect(minStake).to.equal(BigInt(toPeb(5_000_000n)));
    });
    it("Check AddressBook setting", async function () {
      const {
        addressBook,
        cnStakingV1A,
        cnStakingV2B,
        cnStakingV2C,
        cnStakingV2D,
        argsForContractA,
        argsForContractB,
        argsForContractC,
        argsForContractD,
      } = fixture;

      const allAddressInfo = await addressBook.getAllAddressInfo();
      expect(allAddressInfo[0]).to.equalAddrList([
        argsForContractA.nodeId.address,
        argsForContractB.nodeId.address,
        argsForContractC.nodeId.address,
        argsForContractD.nodeId.address,
      ]);
      expect(allAddressInfo[1]).to.equalAddrList([
        cnStakingV1A.address,
        cnStakingV2B.address,
        cnStakingV2C.address,
        cnStakingV2D.address,
      ]);
      expect(allAddressInfo[2]).to.equalAddrList([
        argsForContractA.rewardAddr.address,
        argsForContractB.rewardAddr.address,
        argsForContractC.rewardAddr.address,
        argsForContractD.rewardAddr.address,
      ]);
    });
    it("Check CnStaking setting", async function () {
      const { cnStakingV1A, cnStakingV2B, cnStakingV2C, cnStakingV2D, admin1, admin2, admin3, admin4 } = fixture;

      const stateOfA = await cnStakingV1A.getState();
      const adminB = await cnStakingV2B.adminList(0);
      const adminC = await cnStakingV2C.adminList(0);
      const adminD = await cnStakingV2D.adminList(0);

      expect(stateOfA[3][0]).to.equal(admin1.address);
      expect(adminB).to.equal(admin2.address);
      expect(adminC).to.equal(admin3.address);
      expect(adminD).to.equal(admin4.address);
    });
  });

  // 2. Test create/expire trackers
  describe("Test create/expire trackers", function () {
    it("Only admin can create tracker", async function () {
      const { stakingTracker, trackStart, trackEnd, other1 } = fixture;

      // Note that contractValidator(= accounts[0]) is an owner of contracts by hardhat default setting
      await expect(stakingTracker.connect(other1).createTracker(trackStart, trackEnd)).to.be.revertedWith(
        "Ownable: caller is not the owner",
      );
    });
    it("Create and check a single tracker", async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;

      await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

      // Checkpoint:
      // 1. LastTrackerId
      // 2. allTrackerIds
      // 3. liveTrackerIds
      // 4. trackerInfo

      const lastTrackerId = await stakingTracker.getLastTrackerId();
      const allTrackerIds = await stakingTracker.getAllTrackerIds();
      const liveTrackerIds = await stakingTracker.getLiveTrackerIds();
      const trackerSummary = await stakingTracker.getTrackerSummary(lastTrackerId);

      expect(lastTrackerId).to.equal(1);
      expect(allTrackerIds).to.equalNumberList([1]);
      expect(liveTrackerIds).to.equalNumberList([1]);

      // Check tracker summary
      expect(trackerSummary[0]).to.equal(trackStart);
      expect(trackerSummary[1]).to.equal(trackEnd);
      // Since V1 contract isn't tracked
      expect(trackerSummary[2]).to.equal(3);
      // Both V2 contracts have over than MIN_STAKE
      expect(trackerSummary[3]).to.equal(3);
      expect(trackerSummary[4]).to.equal(3);
    });
    it("Create and expire a single tracker", async function () {
      const { stakingTracker, trackStart, trackEnd, cnStakingV2B } = fixture;

      // Note that contractValidator(= accounts[0]) is an owner of contracts by hardhat default setting
      await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

      await setBlock(trackEnd + 1);

      // Call refreshStake to expire the outdated tracker
      await expect(stakingTracker.refreshStake(cnStakingV2B.address)).to.emit(stakingTracker, "RetireTracker");

      const liveTrackerIds = await stakingTracker.getLiveTrackerIds();

      // Now we don't have live tracker
      expect(liveTrackerIds).to.equalNumberList([]);
    });
    it("Create trackers and check getter", async function () {
      const { stakingTracker, trackStart, trackEnd, argsForContractB, argsForContractC, argsForContractD } = fixture;

      await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

      const trackedGCB = await stakingTracker.getTrackedGC(1, argsForContractB.gcId!);
      const trackedGCC = await stakingTracker.getTrackedGC(1, argsForContractC.gcId!);
      const trackedGCD = await stakingTracker.getTrackedGC(1, argsForContractD.gcId!);

      const allTrackedGCs = await stakingTracker.getAllTrackedGCs(1);

      expect(trackedGCB[0]).to.equal(toPeb(6_000_000n));
      expect(trackedGCB[1]).to.equal(1);
      expect(trackedGCC[0]).to.equal(toPeb(6_000_000n));
      expect(trackedGCC[1]).to.equal(1);
      expect(trackedGCD[0]).to.equal(toPeb(6_000_000n));
      expect(trackedGCD[1]).to.equal(1);

      if (argsForContractB.gcId && argsForContractC.gcId && argsForContractD.gcId) {
        expect(allTrackedGCs[0]).to.equalNumberList([
          argsForContractB.gcId,
          argsForContractC.gcId,
          argsForContractD.gcId,
        ]);
      }

      expect(allTrackedGCs[1]).to.equalNumberList([toPeb(6_000_000n), toPeb(6_000_000n), toPeb(6_000_000n)]);

      expect(allTrackedGCs[2]).to.equalNumberList([1, 1, 1]);
    });
    it("Create and expire multiple trackers", async function () {
      const {
        stakingTracker,
        admin2,
        admin3,
        other1,
        unLockAmounts,
        trackInterval,
        trackStart,
        trackEnd,
        cnStakingV2B,
        cnStakingV2C,
      } = fixture;

      // Test scenario:
      // 1. Create tracker 1
      // 2. cnStakingV2B adds 4m KLAY stakes
      // Check tracker summary
      // setBlock(trackEnd)
      // jumpTime(105)
      // 3. Create tracker 2
      // 4. CnStakingV2C removes 2m KLAY stakes: expire tracker 1 and CnStakingV2C loses its vote
      // Check tracker states
      // setBlock(trackEnd + trackInterval)
      // 5. Create tracker 3
      // 6. CnStakingV2C adds 2m KLAY free-stakes: expire tracker 2 and CnStakingV2C recovers its vote
      // Check tracker states

      // 1. Create tracker 1
      await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

      // 2. cnStakingV2B adds 4m KLAY stakes
      await expect(cnStakingV2B.connect(admin2).stakeKlay({ value: toPeb(4_000_000n) }))
        .to.emit(cnStakingV2B, "StakeKlay")
        .to.emit(stakingTracker, "RefreshStake");

      // Check tracker summary
      let trackerSummary = await stakingTracker.getTrackerSummary(1);
      // Now cnStakingV2B has 2 votes
      expect(trackerSummary[3]).to.equal(4);
      expect(trackerSummary[4]).to.equal(3);

      await setBlock(trackEnd);

      await jumpTime(105);

      // 3. Create tracker 2
      await expect(stakingTracker.createTracker(trackStart + trackInterval, trackEnd + trackInterval)).to.emit(
        stakingTracker,
        "CreateTracker",
      );

      // 4. CnStakingV2C removes 2m KLAY stakes: expire tracker 1 and CnStakingV2C loses its vote
      await submitAndExecuteRequest(cnStakingV2C, [admin3], 1, FuncID.WithdrawLockupStaking, admin3, [
        0,
        FuncID.WithdrawLockupStaking,
        other1.address,
        BigInt(unLockAmounts[0]),
        0,
      ]);

      // Check tracker states
      let lastTrackerId = await stakingTracker.getLastTrackerId();
      let allTrackerIds = await stakingTracker.getAllTrackerIds();
      let liveTrackerIds = await stakingTracker.getLiveTrackerIds();
      trackerSummary = await stakingTracker.getTrackerSummary(lastTrackerId);

      expect(lastTrackerId).to.equal(2);
      expect(allTrackerIds).to.equalNumberList([1, 2]);
      expect(liveTrackerIds).to.equalNumberList([2]);

      expect(trackerSummary[0]).to.equal(trackStart + trackInterval);
      expect(trackerSummary[1]).to.equal(trackEnd + trackInterval);
      expect(trackerSummary[2]).to.equal(3);
      // Since CnV2C lost its vote
      expect(trackerSummary[3]).to.equal(2);
      expect(trackerSummary[4]).to.equal(2);

      await setBlock(trackEnd + trackInterval);

      // 5. Create tracker 3
      await expect(stakingTracker.createTracker(trackStart + trackInterval * 2, trackEnd + trackInterval * 2)).to.emit(
        stakingTracker,
        "CreateTracker",
      );

      // 6. CnStakingV2C adds 2m KLAY free-stakes: expire tracker 2 and CnStakingV2C recovers its vote
      await expect(cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(2_000_000n) }))
        .to.emit(cnStakingV2C, "StakeKlay")
        .to.emit(stakingTracker, "RefreshStake");

      // Check tracker states
      lastTrackerId = await stakingTracker.getLastTrackerId();
      allTrackerIds = await stakingTracker.getAllTrackerIds();
      liveTrackerIds = await stakingTracker.getLiveTrackerIds();
      trackerSummary = await stakingTracker.getTrackerSummary(lastTrackerId);

      expect(lastTrackerId).to.equal(3);
      expect(allTrackerIds).to.equalNumberList([1, 2, 3]);
      expect(liveTrackerIds).to.equalNumberList([3]);

      expect(trackerSummary[0]).to.equal(trackStart + trackInterval * 2);
      expect(trackerSummary[1]).to.equal(trackEnd + trackInterval * 2);
      expect(trackerSummary[2]).to.equal(3);
      // Since CnV2C recovered its vote
      expect(trackerSummary[3]).to.equal(4);
      expect(trackerSummary[4]).to.equal(3);
    });
  });

  // 3. Test functions by calling directly
  // Since it's redundant to call refreshStake / refreshVote directly,
  // we just need to check it emits events
  describe("Test functions by calling directly", function () {
    describe("#refreshStake", function () {
      it("#refreshStake: When all trackers are live", async function () {
        const { stakingTracker, trackStart, trackEnd, cnStakingV2B } = fixture;

        // 1. Create tracker 1 with (trackStart, trackEnd)
        // 2. Create tracker 2 with (trackStart + 10, trackEnd + 10)
        // 3. Create tracker 3 with (trackStart + 20, trackEnd + 20)
        // setBlock(trackStart + 30)
        // 4. Call refreshStake
        // Check tracker states

        // 1. Create tracker 1
        await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

        // 2. Create tracker 2
        await expect(stakingTracker.createTracker(trackStart + 10, trackEnd + 10)).to.emit(
          stakingTracker,
          "CreateTracker",
        );

        // 3. Create tracker 3
        await expect(stakingTracker.createTracker(trackStart + 20, trackEnd + 20)).to.emit(
          stakingTracker,
          "CreateTracker",
        );

        await setBlock(trackStart + 30);

        // 4. Call refreshStake
        await expect(stakingTracker.refreshStake(cnStakingV2B.address)).to.emit(stakingTracker, "RefreshStake");

        // Check tracker states
        const lastTrackerId = await stakingTracker.getLastTrackerId();
        const allTrackerIds = await stakingTracker.getAllTrackerIds();
        const liveTrackerIds = await stakingTracker.getLiveTrackerIds();

        expect(lastTrackerId).to.equal(3);
        expect(allTrackerIds).to.equalNumberList([1, 2, 3]);
        // refreshStake doesn't expire any tracker since they're all live
        expect(liveTrackerIds).to.equalNumberList([1, 2, 3]);
      });
      it("#refreshStake: When there's expired tracker", async function () {
        const { stakingTracker, trackInterval, trackStart, trackEnd, cnStakingV2B } = fixture;

        // 1. Create tracker 1 with (trackStart, trackEnd)
        // 2. Create tracker 2 with (trackEnd, trackEnd + trackInterval)
        // 3. Create tracker 3 with (trackEnd + 10, trackEnd + trackInterval + 10)
        // setBlock(trackStart + 30): all trackers are live
        // Check tracker state
        // setBlock(trackEnd)
        // 4. Call refreshStake: #2 is live (Need to check if it's intended)

        // 1. Create tracker 1
        await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

        // 2. Create tracker 2
        await expect(stakingTracker.createTracker(trackEnd, trackEnd + trackInterval)).to.emit(
          stakingTracker,
          "CreateTracker",
        );

        // 3. Create tracker 3
        await expect(stakingTracker.createTracker(trackEnd + 10, trackEnd + trackInterval + 10)).to.emit(
          stakingTracker,
          "CreateTracker",
        );

        await setBlock(trackStart + 30);

        // Check tracker states
        let lastTrackerId = await stakingTracker.getLastTrackerId();
        let allTrackerIds = await stakingTracker.getAllTrackerIds();
        let liveTrackerIds = await stakingTracker.getLiveTrackerIds();

        expect(lastTrackerId).to.equal(3);
        expect(allTrackerIds).to.equalNumberList([1, 2, 3]);
        expect(liveTrackerIds).to.equalNumberList([1, 2, 3]);

        await setBlock(trackEnd);

        // 4. Call refreshStake
        await expect(stakingTracker.refreshStake(cnStakingV2B.address))
          .to.emit(stakingTracker, "RefreshStake")
          .to.emit(stakingTracker, "RetireTracker");

        // Check tracker states
        lastTrackerId = await stakingTracker.getLastTrackerId();
        allTrackerIds = await stakingTracker.getAllTrackerIds();
        liveTrackerIds = await stakingTracker.getLiveTrackerIds();

        expect(lastTrackerId).to.equal(3);
        expect(allTrackerIds).to.equalNumberList([1, 2, 3]);
        // refreshStake expires tracker 1
        expect(liveTrackerIds).to.equalNumberList([2]); // TODO: check the standard of liveTrackerIds
      });
    });
    describe("#refreshVoter", function () {
      it("#refreshVoter: Wrong staking contract", async function () {
        const { stakingTracker, trackStart, trackEnd, cnStakingV1A } = fixture;

        await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

        // 1. V1 contract can't update voter
        await expect(stakingTracker.refreshVoter(cnStakingV1A.address)).to.be.revertedWith(
          "Invalid CnStaking contract",
        );

        // 2. Not a staking contract
        await expect(stakingTracker.refreshVoter(stakingTracker.address)).to.be.revertedWith("Not a staking contract");
      });
      it("#refreshVoter: Emits RefreshVoter event", async function () {
        const { stakingTracker, trackStart, trackEnd, cnStakingV2B } = fixture;

        await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");

        await expect(stakingTracker.refreshVoter(cnStakingV2B.address)).to.emit(stakingTracker, "RefreshVoter");
      });
    });
  });

  // 4. Test functions by calling through CnStakingV2.sol
  // Check V1 caller case only in stakeKlay test
  describe("4. Test functions by calling through CnStakingV2.sol", function () {
    this.beforeEach(async function () {
      const { stakingTracker, trackStart, trackEnd } = fixture;

      await expect(stakingTracker.createTracker(trackStart, trackEnd)).to.emit(stakingTracker, "CreateTracker");
    });
    describe("#refreshStake", async function () {
      it("#stakeKlay: By V1 contract", async function () {
        const { stakingTracker, admin1, cnStakingV1A } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // CnStakingV1A stake 4m KLAY
        await expect(cnStakingV1A.connect(admin1).stakeKlay({ value: toPeb(4_000_000n) })).to.emit(
          cnStakingV1A,
          "StakeKlay",
        );

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // Doesn't affect to StakingTracker
        expect(trackerSummary[3]).to.equal(3);
      });
      it("#stakeKlay: By V2 contract", async function () {
        const { stakingTracker, admin2, cnStakingV2B } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // CnStakingV2B stake 4m KLAY
        await expect(cnStakingV2B.connect(admin2).stakeKlay({ value: toPeb(4_000_000n) })).to.emit(
          cnStakingV2B,
          "StakeKlay",
        );

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 4 since CnStakingV2B now has 10m staked KLAY
        expect(trackerSummary[3]).to.equal(4);
      });
      it("#withdrawLockupStaking", async function () {
        const { stakingTracker, admin3, requirement, cnStakingV2C, other1 } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        await jumpTime(105);

        // CnStakingV2C withdraws 2m KLAY
        await submitAndExecuteRequest(cnStakingV2C, [admin3], requirement, FuncID.WithdrawLockupStaking, admin3, [
          0,
          FuncID.WithdrawLockupStaking,
          other1.address,
          BigInt(toPeb(2_000_000n)),
          0,
        ]);

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 2 since CnStakingV2C now has no vote
        expect(trackerSummary[3]).to.equal(2);
      });
      it("#approveStakingWithdrawal", async function () {
        const { stakingTracker, admin3, requirement, cnStakingV2C, other1 } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // stakeKlay 4m KLAY
        await expect(cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(4_000_000n) }))
          .to.emit(cnStakingV2C, "StakeKlay")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(4);

        jumpTime(WEEK);

        // CnStakingV2C withdraws 2m KLAY
        await submitAndExecuteRequest(cnStakingV2C, [admin3], requirement, FuncID.ApproveStakingWithdrawal, admin3, [
          0,
          FuncID.ApproveStakingWithdrawal,
          other1.address,
          BigInt(toPeb(4_000_000n)),
          0,
        ]);

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 3 since CnStakingV2C's unstaking is 4m KLAY
        expect(trackerSummary[3]).to.equal(3);
      });
      it("#cancelApprovedStakingWithdrawal", async function () {
        const { stakingTracker, admin3, requirement, cnStakingV2C, other1 } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // stakeKlay 4m KLAY
        await expect(cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(4_000_000n) }))
          .to.emit(cnStakingV2C, "StakeKlay")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(4);

        jumpTime(WEEK);

        // CnStakingV2C withdraws 4m KLAY
        await submitAndExecuteRequest(cnStakingV2C, [admin3], requirement, FuncID.ApproveStakingWithdrawal, admin3, [
          0,
          FuncID.ApproveStakingWithdrawal,
          other1.address,
          BigInt(toPeb(4_000_000n)),
          0,
        ]);

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 3 since CnStakingV2C's unstaking is 4m KLAY
        expect(trackerSummary[3]).to.equal(3);

        // CnStakingV2C cancels withdrawal
        await submitAndExecuteRequest(
          cnStakingV2C,
          [admin3],
          requirement,
          FuncID.CancelApprovedStakingWithdrawal,
          admin3,
          [1, FuncID.CancelApprovedStakingWithdrawal, 0, 0, 0],
        );

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 4 since CnStakingV2C cancels 4m KLAY unstaking
        expect(trackerSummary[3]).to.equal(4);
      });
      it("#withdrawApprovedStaking: Successfully withdraw free stakes", async function () {
        const { stakingTracker, admin3, requirement, cnStakingV2C, other1 } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // stakeKlay 4m KLAY
        await expect(cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(4_000_000n) }))
          .to.emit(cnStakingV2C, "StakeKlay")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(4);

        // CnStakingV2C withdraws 4m KLAY
        await submitAndExecuteRequest(cnStakingV2C, [admin3], requirement, FuncID.ApproveStakingWithdrawal, admin3, [
          0,
          FuncID.ApproveStakingWithdrawal,
          other1.address,
          BigInt(toPeb(4_000_000n)),
          0,
        ]);

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 3 since CnStakingV2C's unstaking is 4m KLAY
        expect(trackerSummary[3]).to.equal(3);

        jumpTime(WEEK);

        // CnStakingV2C withdraws 4m free-stakes
        await expect(cnStakingV2C.connect(admin3).withdrawApprovedStaking(0))
          .to.emit(cnStakingV2C, "WithdrawApprovedStaking")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 3 since CnStakingV2C withdraws 4m KLAY
        expect(trackerSummary[3]).to.equal(3);
      });
      it("#withdrawApprovedStaking: Can't withdraw it since it's over 2 weeks", async function () {
        const { stakingTracker, admin3, requirement, cnStakingV2C, other1 } = fixture;

        let trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(3);

        // stakeKlay 4m KLAY
        await expect(cnStakingV2C.connect(admin3).stakeKlay({ value: toPeb(4_000_000n) }))
          .to.emit(cnStakingV2C, "StakeKlay")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        expect(trackerSummary[3]).to.equal(4);

        // CnStakingV2C withdraws 4m KLAY
        await submitAndExecuteRequest(cnStakingV2C, [admin3], requirement, FuncID.ApproveStakingWithdrawal, admin3, [
          0,
          FuncID.ApproveStakingWithdrawal,
          other1.address,
          BigInt(toPeb(4_000_000n)),
          0,
        ]);

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 3 since CnStakingV2C's unstaking is 4m KLAY
        expect(trackerSummary[3]).to.equal(3);

        jumpTime(2 * WEEK);

        // CnStakingV2C withdraws 4m free-stakes
        await expect(cnStakingV2C.connect(admin3).withdrawApprovedStaking(0))
          .to.emit(cnStakingV2C, "CancelApprovedStakingWithdrawal")
          .to.emit(stakingTracker, "RefreshStake");

        trackerSummary = await stakingTracker.getTrackerSummary(1);

        // totalVotes will be 4 since CnStakingV2C's withdrawal is canceled
        expect(trackerSummary[3]).to.equal(4);
      });
    });
    describe("#refreshVoter", async function () {
      it("#updateVoterAddress", async function () {
        const { stakingTracker, admin4, requirement, cnStakingV2D, other1 } = fixture;

        // Check voter before updateVoterAddress
        let cnStakingState = await stakingTracker.readCnStaking(cnStakingV2D.address);

        // We haven't set voter address yet
        expect(cnStakingState[4]).to.equal(ethers.constants.AddressZero);

        // submit and execute updateVoterAddress request
        await submitAndExecuteRequest(cnStakingV2D, [admin4], requirement, FuncID.UpdateVoterAddress, admin4, [
          0,
          FuncID.UpdateVoterAddress,
          other1.address,
          0,
          0,
        ]);

        cnStakingState = await stakingTracker.readCnStaking(cnStakingV2D.address);

        // Now voter address should be other1.address
        expect(cnStakingState[4]).to.equal(other1.address);

        await expect(stakingTracker.refreshVoter(cnStakingV2D.address)).to.emit(stakingTracker, "RefreshVoter");
        expect(await stakingTracker.voterToGCId(other1.address)).to.equal(702);
      });
    });
  });
});
