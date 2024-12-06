import { loadFixture } from "@nomicfoundation/hardhat-network-helpers";
import { expect } from "chai";
import {
  FuncID,
  augmentChai,
  toPeb,
  jumpTime,
  jumpBlock,
  submitAndExecuteRequest,
  VoteChoice,
  DAY,
  WEEK,
  Actions,
} from "./common/helper";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers.js";
import { ethers } from "hardhat";
import { BytesLike } from "ethers";
import { votingTestFixture } from "./common/fixtures";

const minVotingDelay = 1 * DAY;
const maxVotingDelay = 28 * DAY;
const minVotingPeriod = 1 * DAY;
const maxVotingPeriod = 28 * DAY;

type UnPromisify<T> = T extends Promise<infer U> ? U : T;

/**
 * @dev This test is for Voting.sol
 * 1. Test initializing process
 *    - Check constructor
 * 2. Test functions of Voting.sol
 *    - propose
 *    - cancel
 *    - castVote
 *    - queue
 *    - execute
 *    - updateStakingTracker
 *    - updateAccessRule
 *    - updateTimingRule
 * 3. Test votes count and power
 */
describe("Voting.sol", function () {
  function buildActions(actionList: Actions[]): [string[], string[], BytesLike[]] {
    const {
      txSimpleTransfer100KLAY,
      txSimpleTransfer200KLAY,
      txUpdateStakingTracker,
      txUpdateSecretary,
      txUpdateAccessRule,
      txUpdateTimingRule,
      txUpdateAccessRuleWrong,
      txUpdateTimingRuleWrong,
    } = fixture;

    const allTargets = [
      txSimpleTransfer100KLAY.target,
      txSimpleTransfer200KLAY.target,
      txUpdateStakingTracker.target,
      txUpdateSecretary.target,
      txUpdateAccessRule.target,
      txUpdateTimingRule.target,
      txUpdateAccessRuleWrong.target,
      txUpdateTimingRuleWrong.target,
    ];
    const allValues = [
      txSimpleTransfer100KLAY.value,
      txSimpleTransfer200KLAY.value,
      txUpdateStakingTracker.value,
      txUpdateSecretary.value,
      txUpdateAccessRule.value,
      txUpdateTimingRule.value,
      txUpdateAccessRuleWrong.value,
      txUpdateTimingRuleWrong.value,
    ];
    const allCalldata = [
      txSimpleTransfer100KLAY.calldata,
      txSimpleTransfer200KLAY.calldata,
      txUpdateStakingTracker.calldata,
      txUpdateSecretary.calldata,
      txUpdateAccessRule.calldata,
      txUpdateTimingRule.calldata,
      txUpdateAccessRuleWrong.calldata,
      txUpdateTimingRuleWrong.calldata,
    ];

    if (actionList.length === 0) {
      return [allTargets, allValues, allCalldata];
    }

    const targets = [] as string[];
    const values = [] as string[];
    const calldata = [] as BytesLike[];

    for (let i = 0; i < actionList.length; i++) {
      const action = actionList[i];
      targets[i] = allTargets[action];
      values[i] = allValues[action];
      calldata[i] = allCalldata[action];
    }

    return [targets, values, calldata];
  }
  let fixture: UnPromisify<ReturnType<typeof votingTestFixture>>;
  beforeEach(async function () {
    augmentChai();
    fixture = await loadFixture(votingTestFixture);
  });

  // Test initializing process
  describe("Voting Initialize", function () {
    it("Check constructor", async function () {
      const { voting, stakingTracker, secretary } = fixture;

      const stakingTrackerAddr = await voting.stakingTracker();
      const secretaryAddr = await voting.secretary();

      const nextProposalId = await voting.nextProposalId();

      const accessRule = await voting.accessRule();
      const timingRule = await voting.timingRule();

      expect(stakingTrackerAddr).to.equal(stakingTracker.address);
      expect(secretaryAddr).to.equal(secretary.address);

      expect(nextProposalId).to.equal(1);

      expect(accessRule.slice(0, 4)).to.equalBooleanList([true, false, true, false]);
      expect(timingRule).to.equalNumberList([minVotingDelay, maxVotingDelay, minVotingPeriod, maxVotingPeriod]);
    });
  });

  // 2. Test functions of Voting.sol
  describe("Test functions of Voting.sol", function () {
    describe("#propose", function () {
      it("#propose: targets, values, calldata should have same length", async function () {
        const { voting, txSimpleTransfer100KLAY, txSimpleTransfer200KLAY } = fixture;

        // Target length is 1, while values and calldata are 2
        const targets = [txSimpleTransfer100KLAY.target];
        const values = [txSimpleTransfer100KLAY.value, txSimpleTransfer200KLAY.value];
        const calldata = [txSimpleTransfer100KLAY.calldata, txSimpleTransfer200KLAY.calldata];

        await expect(
          voting.propose("Transfer KLAY to receiver", targets, values, calldata, DAY, DAY),
        ).to.be.revertedWith("Invalid actions");
      });
      it("#propose: Should satisfy timing rule ", async function () {
        const { voting } = fixture;

        const [targets, values, calldata] = buildActions([Actions.txSimpleTransfer100KLAY]);

        await expect(
          voting.propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay - 1, WEEK),
        ).to.be.revertedWith("Invalid votingDelay");

        await expect(
          voting.propose("Transfer KLAY to receiver", targets, values, calldata, maxVotingDelay + 1, WEEK),
        ).to.be.revertedWith("Invalid votingDelay");

        await expect(
          voting.propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod - 1),
        ).to.be.revertedWith("Invalid votingPeriod");

        await expect(
          voting.propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, maxVotingPeriod + 1),
        ).to.be.revertedWith("Invalid votingPeriod");
      });
      it("#propose: Should satisfy access rule ", async function () {
        const { voting, secretary, stakingTracker, voter1 } = fixture;

        const [targets, values, calldata] = buildActions([Actions.txSimpleTransfer100KLAY]);

        // Current access rule:
        // Secretary: Can propose
        // Voter: Can't propose

        // 1. Secretary can propose
        await expect(
          voting
            .connect(secretary)
            .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod),
        )
          .to.emit(voting, "ProposalCreated")
          .to.emit(stakingTracker, "CreateTracker");

        // 2. Voter can't propose
        await expect(
          voting
            .connect(voter1)
            .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod),
        ).to.be.revertedWith("Not the secretary");

        // Update access rule
        await expect(voting.connect(secretary).updateAccessRule(true, true, true, false)).to.emit(
          voting,
          "UpdateAccessRule",
        );

        // 3. Now voter can propose
        await expect(
          voting
            .connect(voter1)
            .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod),
        )
          .to.emit(voting, "ProposalCreated")
          .to.emit(stakingTracker, "CreateTracker");
      });
      it("#propose: Successfully propose proposal", async function () {
        const { voting, secretary, stakingTracker } = fixture;

        const [targets, values, calldata] = buildActions([Actions.txSimpleTransfer100KLAY]);

        // 1. Secretary can propose
        await expect(
          voting
            .connect(secretary)
            .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod),
        )
          .to.emit(voting, "ProposalCreated")
          .to.emit(stakingTracker, "CreateTracker");

        // Check proposal state
        const proposalActions = await voting.getActions(1);

        expect(proposalActions[0]).to.equalAddrList(targets);
        expect(proposalActions[1]).to.equalNumberList(values);
      });
    });
    describe("Tests which need to be done after proposal", function () {
      this.beforeEach(async function () {
        const { voting, secretary } = fixture;

        const [targets, values, calldata] = buildActions([
          Actions.txSimpleTransfer100KLAY,
          Actions.txSimpleTransfer200KLAY,
        ]);

        // Propose #1
        await voting
          .connect(secretary)
          .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod);
      });
      describe("#cancel", function () {
        it("#cancel: Only proposer can cancel proposal", async function () {
          const { voting, voter2 } = fixture;

          // Voter2 is not a proposer
          await expect(voting.connect(voter2).cancel(1)).to.be.revertedWith("Not the proposer");
        });
        it("#cancel: Only pending proposal can be canceled", async function () {
          const { voting, secretary } = fixture;

          await jumpBlock(DAY);

          // Proposal #1 is currently active state
          await expect(voting.connect(secretary).cancel(1)).to.be.revertedWith("Not allowed in current state");
        });
        it("#cancel: Successfully cancel proposal", async function () {
          const { voting, secretary } = fixture;

          await jumpBlock(DAY / 2);

          // Proposal #1 is currently pending state
          await expect(voting.connect(secretary).cancel(1)).to.emit(voting, "ProposalCanceled");
        });
      });
      describe("#castVote", function () {
        it("#castVote: Can't vote for empty proposal", async function () {
          const { voting, voter1 } = fixture;

          // There's no proposal #2
          await expect(voting.connect(voter1).castVote(2, VoteChoice.Yes)).to.be.revertedWith("No such proposal");
        });
        it("#castVote: Can't vote before voting starts", async function () {
          const { voting, voter1 } = fixture;

          await jumpTime(DAY / 2);

          // Not active proposal yet
          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.be.revertedWith(
            "Not allowed in current state",
          );
        });
        it("#castVote: Can't vote for canceled proposal", async function () {
          const { voting, secretary, voter1 } = fixture;

          await jumpTime(DAY / 2);

          // Not active proposal yet
          await expect(voting.connect(secretary).cancel(1)).to.emit(voting, "ProposalCanceled");

          await jumpBlock(DAY / 2);

          // Proposal #1 is currently canceled state
          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.be.revertedWith(
            "Not allowed in current state",
          );
        });
        it("#castVote: Only registered voter can vote", async function () {
          const { voting, other1 } = fixture;

          await jumpBlock(DAY + 1);

          // Not a registered voter
          await expect(voting.connect(other1).castVote(1, VoteChoice.Yes)).to.be.revertedWith("Not a registered voter");
        });
        it("#castVote: Need to have votes", async function () {
          const { voting, cnStakingV2A, admin1, voter1, other1 } = fixture;

          await jumpTime(105);

          // CnStakingV2A withdraws KLAY before it starts, becoming a voter with 0 votes
          await submitAndExecuteRequest(cnStakingV2A, [admin1], 1, FuncID.WithdrawLockupStaking, admin1, [
            0,
            FuncID.WithdrawLockupStaking,
            other1.address,
            BigInt(toPeb(2_000_000n)),
            0,
          ]);

          await jumpBlock(DAY);

          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.be.revertedWith("Not eligible to vote");
        });
        it("#castVote: After pending, tracker doesn't track staking", async function () {
          const { voting, cnStakingV2A, admin1, voter1, other1 } = fixture;

          await jumpTime(105);

          await jumpBlock(DAY);

          // CnStakingV2A withdraws KLAY after pending state, which means it has votes for proposal #1
          await submitAndExecuteRequest(cnStakingV2A, [admin1], 1, FuncID.WithdrawLockupStaking, admin1, [
            0,
            FuncID.WithdrawLockupStaking,
            other1.address,
            BigInt(toPeb(2_000_000n)),
            0,
          ]);

          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");
        });
        it("#castVote: Not a valid choice", async function () {
          const { voting, voter1 } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter1).castVote(1, VoteChoice.Abstain + 1)).to.be.revertedWith(
            "Not a valid choice",
          );
        });
        it("#castVote: Can't vote twice", async function () {
          const { voting, voter1 } = fixture;

          await jumpBlock(DAY);

          // Vote for proposal #1
          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          // Can't vote twice
          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.be.revertedWith("Already voted");
        });
        it("#castVote: Successfully vote and check proposal", async function () {
          const { voting, voter1, voter2 } = fixture;

          await jumpBlock(DAY);

          // Vote for proposal #1
          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          // Vote for proposal #1
          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          const proposalTally = await voting.getProposalTally(1);

          expect(proposalTally[0]).to.equal(3);
          expect(proposalTally[1]).to.equal(0);
          expect(proposalTally[2]).to.equal(0);
          // 4 GCs, 4 eligible count, (4 + 2) / 3 = 2
          expect(proposalTally[3]).to.equal(2);
          // 7 votes power, (7 + 2) / 3 = 3
          expect(proposalTally[4]).to.equal(3);
          expect(proposalTally[5]).to.equalNumberList([700, 701]);
        });
      });
      describe("#queue", function () {
        it("#queue: Can't queue failed proposal (Not enough participation)", async function () {
          const { voting, voter1 } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          // Only voter1 votes for proposal #1, which can't satisfy quorum
          await expect(voting.connect(voter1).queue(1)).to.be.revertedWith("Not allowed in current state");
        });
        it("#queue: Can't queue failed proposal (Not enough YES votes)", async function () {
          const { voting, voter1, voter2, voter4 } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter1).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");
          await expect(voting.connect(voter4).castVote(1, VoteChoice.No)).to.emit(voting, "VoteCast");
          await expect(voting.connect(voter2).castVote(1, VoteChoice.No)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          // Now it satisfies quorum, but not passed
          await expect(voting.connect(voter1).queue(1)).to.be.revertedWith("Not allowed in current state");
        });
        it("#queue: Only executor can queue proposal", async function () {
          const { voting, voter2, voter4 } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(voter4).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          // Currently, voter is not an executor
          await expect(voting.connect(voter2).queue(1)).to.be.revertedWith("Not the secretary");
        });
        it("#queue: Successfully queue proposal", async function () {
          const { voting, voter2, voter4, secretary } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(voter4).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          await expect(voting.connect(secretary).queue(1)).to.emit(voting, "ProposalQueued");

          const proposalSchedule = await voting.getProposalSchedule(1);

          expect(proposalSchedule[6]).to.equal(true);
        });
      });
      describe("#execute", function () {
        it("#execute: Can't execute proposal in not proper timing", async function () {
          const { voting, voter2, voter4, secretary } = fixture;

          // 1. In active state
          await jumpBlock(DAY);

          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(voter4).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("Not allowed in current state");

          // 2. In active < now < queue state
          await jumpBlock(DAY);

          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("Not allowed in current state");

          // 3. In queue< now < eta state
          await expect(voting.connect(secretary).queue(1)).to.emit(voting, "ProposalQueued");

          await jumpBlock(DAY);

          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("Not yet executable");

          // 4 In now > execute timeout state
          await jumpBlock(3 * WEEK);
          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("Not allowed in current state");
        });
        it("#execute: Only executor can execute proposal", async function () {
          const { voting, voter2, voter4, secretary } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(voter4).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          await expect(voting.connect(secretary).queue(1)).to.emit(voting, "ProposalQueued");

          // After eta
          await jumpBlock(2 * DAY);

          await expect(voting.connect(voter2).execute(1)).to.be.revertedWith("Not the secretary");
        });
        it("#execute: Successfully execute proposal", async function () {
          const { voting, voter2, voter4, secretary, other1 } = fixture;

          await jumpBlock(DAY);

          await expect(voting.connect(voter2).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await expect(voting.connect(voter4).castVote(1, VoteChoice.Yes)).to.emit(voting, "VoteCast");

          await jumpBlock(DAY);

          await expect(voting.connect(secretary).queue(1)).to.emit(voting, "ProposalQueued");

          // After eta
          await jumpBlock(2 * DAY);

          const other1Before = await ethers.provider.getBalance(other1.address);

          await expect(voting.connect(secretary).execute(1)).to.emit(voting, "ProposalExecuted");

          const other1After = await ethers.provider.getBalance(other1.address);

          // Our actions are transfer 100, 200 KLAY to other1.address
          expect(other1After.sub(other1Before)).to.equal(toPeb(300n));
        });
      });
    });
    describe("Tests for updating contract parameters", function () {
      async function goToQueueState(voter: SignerWithAddress[], choice: number) {
        const { voting, secretary } = fixture;

        await jumpBlock(DAY);

        for (let i = 0; i < voter.length; i++) {
          await expect(voting.connect(voter[i]).castVote(1, choice)).to.emit(voting, "VoteCast");
        }

        await jumpBlock(DAY);

        await expect(voting.connect(secretary).queue(1)).to.emit(voting, "ProposalQueued");
      }
      describe("#updateStakingTracker", function () {
        it("#updateStakingTracker: Can't update staking tracker while there's an active tracker", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateStakingTracker]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update staking tracker address", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          // Propose #2
          await voting
            .connect(secretary)
            .propose("Update staking tracker address", targets, values, calldata, minVotingDelay, minVotingPeriod);

          // Proposal #2 is currently active state
          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith(
            "Cannot update tracker when there is an active tracker",
          );
        });
        it("#updateStakingTracker: Successfully update staking tracker", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateStakingTracker]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update staking tracker address", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.emit(voting, "UpdateStakingTracker");
        });
      });
      describe("#updateSecretary", function () {
        it("#updateSecretary: Successfully update secretary address", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateSecretary]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update secretary address", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.emit(voting, "UpdateSecretary");
        });
      });
      describe("#updateAccessRule", function () {
        it("#updateAccessRule: Invalid access rule", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [target, value, calldata] = buildActions([Actions.txUpdateAccessRuleWrong]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update access rule", target, value, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("No propose access");
        });
        it("#updateAccessRule: Successfully update access rule through secretary", async function () {
          const { voting, secretary, voter2 } = fixture;

          // Only secretary can update access rule
          await expect(voting.connect(voter2).updateAccessRule(true, true, true, true)).to.be.revertedWith(
            "Not a governance transaction or secretary",
          );

          await expect(voting.connect(secretary).updateAccessRule(true, true, true, true)).to.emit(
            voting,
            "UpdateAccessRule",
          );

          expect(await voting.accessRule()).to.equalBooleanList([true, true, true, true]);
        });
        it("#updateAccessRule: Successfully update access rule through governance", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateAccessRule]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update access rule", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.emit(voting, "UpdateAccessRule");

          expect(await voting.accessRule()).to.equalBooleanList([true, true, true, true]);
        });
      });
      describe("#updateTimingRule", function () {
        it("#updateTimingRule: Invalid timing rule", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateTimingRuleWrong]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update timing rule", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.be.revertedWith("Invalid timing");
        });
        it("#updateTimingRule: Successfully update timing rule through secretary", async function () {
          const { voting, secretary, voter2 } = fixture;

          // Only secretary can update access rule
          await expect(voting.connect(voter2).updateTimingRule(DAY, WEEK, DAY, WEEK)).to.be.revertedWith(
            "Not a governance transaction or secretary",
          );

          await expect(voting.connect(secretary).updateTimingRule(DAY, WEEK, DAY, WEEK)).to.emit(
            voting,
            "UpdateTimingRule",
          );

          expect(await voting.timingRule()).to.equalNumberList([DAY, WEEK, DAY, WEEK]);
        });
        it("#updateAccessRule: Successfully update timing rule through governance", async function () {
          const { voting, secretary, voter2, voter4 } = fixture;

          const [targets, values, calldata] = buildActions([Actions.txUpdateTimingRule]);

          // Propose #1
          await voting
            .connect(secretary)
            .propose("Update access rule", targets, values, calldata, minVotingDelay, minVotingPeriod);

          await goToQueueState([voter2, voter4], VoteChoice.Yes);

          await jumpBlock(2 * DAY);

          await expect(voting.connect(secretary).execute(1)).to.emit(voting, "UpdateTimingRule");

          expect(await voting.timingRule()).to.equalNumberList([DAY, WEEK, DAY, WEEK]);
        });
      });
    });
  });

  // 3. Test votes count and power
  describe("Test votes count and power", function () {
    this.beforeEach(async function () {
      const { voting, secretary } = fixture;

      const [targets, values, calldata] = buildActions([Actions.txSimpleTransfer100KLAY]);

      // Propose #1
      await voting
        .connect(secretary)
        .propose("Transfer KLAY to receiver", targets, values, calldata, minVotingDelay, minVotingPeriod);
    });
    it("Withdraw stakes in pending state", async function () {
      const { voting, cnStakingV2A, admin1, secretary, voter1, voter4, other1 } = fixture;

      const [targets, values, calldata] = buildActions([Actions.txSimpleTransfer200KLAY]);
      // Current eligible voter and votes power:
      // A: 1
      // B: 2
      // C: 1
      // D: 3

      // If A withdraws 2m KLAY, then it also affects to D's votes power since
      // votes power is capped by (numEligibles - 1), which is now 2.

      await jumpTime(105);

      // Now A has no votes
      await submitAndExecuteRequest(cnStakingV2A, [admin1], 1, FuncID.WithdrawLockupStaking, admin1, [
        0,
        FuncID.WithdrawLockupStaking,
        other1.address,
        BigInt(toPeb(2_000_000n)),
        0,
      ]);

      let votesInfoA = await voting.getVotes(1, voter1.address);

      expect(votesInfoA[0]).to.equal(700);
      expect(votesInfoA[1]).to.equal(0);

      let votesInfoD = await voting.getVotes(1, voter4.address);

      expect(votesInfoD[0]).to.equal(703);
      expect(votesInfoD[1]).to.equal(2);

      await jumpBlock(DAY);

      // Now eligible voter and votes power:
      // A: 0
      // B: 2
      // C: 1
      // D: 2

      // Propose #2
      await voting
        .connect(secretary)
        .propose("Transfer 200 KLAY", targets, values, calldata, minVotingDelay, minVotingPeriod);

      // Since now proposal #2 is in pending state, A has 2 votes power and D also recovers its votes power
      await cnStakingV2A.connect(admin1).stakeKlay({ value: toPeb(6_000_000n) });

      votesInfoA = await voting.getVotes(2, voter1.address);

      expect(votesInfoA[0]).to.equal(700);
      expect(votesInfoA[1]).to.equal(2);

      votesInfoD = await voting.getVotes(2, voter4.address);

      expect(votesInfoD[0]).to.equal(703);
      expect(votesInfoD[1]).to.equal(3);
    });
  });
});
