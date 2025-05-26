// Copyright 2022 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/Strings.sol";
import "./IVoting.sol";
import "./StakingTracker.sol";

contract Voting is IVoting {
    // Types

    struct Proposal {
        // Contents
        address proposer;
        string description;
        address[] targets; // Transaction 'to' addresses
        uint256[] values; // Transaction 'value' amounts
        bytes[] calldatas; // Transaction 'input' data
        // Schedule
        uint256 voteStart; // propose()d block + votingDelay
        uint256 voteEnd; // voteStart + votingPeriod
        uint256 queueDeadline; // voteEnd + queueTimeout
        uint256 eta; // queue()d block + execDelay
        uint256 execDeadline; // queue()d block + execDelay + execTimeout
        bool canceled; // true if successfully cancel()ed
        bool queued; // true if successfully queue()d
        bool executed; // true if successfully execute()d
        // Vote counting
        address stakingTracker;
        uint256 trackerId;
        uint256 totalYes;
        uint256 totalNo;
        uint256 totalAbstain;
        uint256 quorumCount;
        uint256 quorumPower;
        uint256[] voters;
        mapping(uint256 => Receipt) receipts;
    }

    struct AccessRule {
        // True if the secretary can propose()
        bool secretaryPropose;
        // True if any eligible voter at the time of the submission can propose() a proposal
        bool voterPropose;
        // True if the secretary can queue() and execute()
        bool secretaryExecute;
        // True if any eligible voter of a given proposal can queue() and execute() the proposal.
        bool voterExecute;
    }

    struct TimingRule {
        uint256 minVotingDelay;
        uint256 maxVotingDelay;
        uint256 minVotingPeriod;
        uint256 maxVotingPeriod;
    }

    // States

    mapping(uint256 => Proposal) private proposals;
    uint256 public nextProposalId;

    /// @dev The address of StakingTracker.
    /// Intended for internal use only, but is public for debugging purposes.
    /// This address is used by newly created proposals.
    address public override stakingTracker;

    /// @dev The address of the Voting secretary.
    /// The secretary can be zero address to signify the absence of the secretary.
    address public override secretary;

    /// @dev The access control rule of some important functions.
    AccessRule public override accessRule;
    /// @dev The timing rules of proposal schedule
    TimingRule public override timingRule;

    uint256 public constant DAY = 86400;
    /// @dev Grace period to queue() passed proposals in block numbers
    uint256 public override queueTimeout = 14 * DAY;
    /// @dev A minimum delay before a queued transaction can be executed in block numbers
    uint256 public override execDelay = 2 * DAY;
    /// @dev Grace period to execute() queued proposals since `execDelay` in block numbers
    uint256 public override execTimeout = 14 * DAY;

    constructor(address _tracker, address _secretary) {
        if (_tracker != address(0)) {
            stakingTracker = _tracker;
        } else {
            // This contract becomes the owner
            stakingTracker = address(new StakingTracker());
        }

        secretary = _secretary;

        nextProposalId = 1;

        // Initial rules
        accessRule.secretaryPropose = true;
        accessRule.voterPropose = false;
        accessRule.secretaryExecute = true;
        accessRule.voterExecute = false;
        validateAccessRule();

        timingRule.minVotingDelay = 1 * DAY;
        timingRule.maxVotingDelay = 28 * DAY;
        timingRule.minVotingPeriod = 1 * DAY;
        timingRule.maxVotingPeriod = 28 * DAY;
        validateTimingRule();
    }

    /// @dev Check for propose() access permission
    function checkProposeAccess(uint256 proposalId) internal view {
        checkAccess(proposalId, accessRule.secretaryPropose, accessRule.voterPropose);
    }

    /// @dev Check for queue() and execute() access permission
    function checkExecuteAccess(uint256 proposalId) internal view {
        checkAccess(proposalId, accessRule.secretaryExecute, accessRule.voterExecute);
    }

    /// @dev Check that sender has access to a certain operation for the given proposal.
    ///
    /// @param proposalId       The proposal ID which the operation changes
    /// @param secretaryAccess  True if the operation is allowed to the secretary
    /// @param voterAccess      True if the operation is allowed to any voter of the proposal
    function checkAccess(uint256 proposalId, bool secretaryAccess, bool voterAccess) internal view {
        // if ( sA &&  vA), msg.sender must be the secretary or a voter.
        //   Note that in this case, the revert message would be
        //   "Not a registered voter" or "Not eligible to vote".
        // if ( sA && !vA), msg.sender must be the secretary.
        // if (!sa &&  vA), msg.sender must be a voter.
        if (secretaryAccess && msg.sender == secretary) {
            return;
        } else if (voterAccess) {
            // check that the sender is an eligible voter of the given proposal.
            (uint256 gcId, uint256 votes) = getVotes(proposalId, msg.sender);
            require(gcId != 0, "Not a registered voter");
            require(votes > 0, "Not eligible to vote");
        } else {
            revert("Not the secretary");
        }
    }

    // Modifiers

    /// @dev Sender must have execute permission of the proposal
    modifier onlyExecutor(uint256 proposalId) {
        checkExecuteAccess(proposalId);
        _;
    }

    /// @dev The proposal must exist and is in the speciefied state
    modifier onlyState(uint256 proposalId, ProposalState s) {
        require(proposals[proposalId].proposer != address(0), "No such proposal");
        require(state(proposalId) == s, "Not allowed in current state");
        _;
    }

    /// @dev Sender must be this contract, i.e. executed via governance proposal
    modifier onlyGovernance() {
        require(address(this) == msg.sender, "Not a governance transaction");
        _;
    }

    /// @dev Sender must be this contract or the secretary.
    modifier onlyGovernanceOrSecretary() {
        require(msg.sender == address(this) || msg.sender == secretary, "Not a governance transaction or secretary");
        _;
    }

    // Mutators

    /// @dev Create a Proposal
    /// @param description   Proposal text
    /// @param targets       List of transaction target addresses
    /// @param values        List of KLAY values to send along with transactions
    /// @param calldatas     List of transaction calldatas
    /// @param votingDelay   Delay from proposal submission to voting start in block numbers
    /// @param votingPeriod  Duration of the voting in block numbers
    function propose(
        string memory description,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        uint256 votingDelay,
        uint256 votingPeriod
    ) external override returns (uint256 proposalId) {
        require(targets.length == values.length && targets.length == calldatas.length, "Invalid actions");
        require(
            timingRule.minVotingDelay <= votingDelay && votingDelay <= timingRule.maxVotingDelay,
            "Invalid votingDelay"
        );
        require(
            timingRule.minVotingPeriod <= votingPeriod && votingPeriod <= timingRule.maxVotingPeriod,
            "Invalid votingPeriod"
        );

        proposalId = nextProposalId;
        nextProposalId++;
        Proposal storage p = proposals[proposalId];

        p.proposer = msg.sender;
        p.description = description;
        p.targets = targets;
        p.values = values;
        p.calldatas = calldatas;

        p.voteStart = block.number + votingDelay;
        p.voteEnd = p.voteStart + votingPeriod;
        p.queueDeadline = p.voteEnd + queueTimeout;

        // Finalize voter list and track balance changes during the preparation period
        p.stakingTracker = stakingTracker;
        p.trackerId = IStakingTracker(p.stakingTracker).createTracker(block.number, p.voteStart);

        // Permission check must be done here since it requires trackerId.
        checkProposeAccess(proposalId);

        emit ProposalCreated(
            proposalId,
            p.proposer,
            p.targets,
            p.values,
            new string[](p.targets.length),
            p.calldatas,
            p.voteStart,
            p.voteEnd,
            p.description
        );
    }

    /// @dev Cancel a proposal
    /// The proposal must be in Pending state
    /// Only the proposer of the proposal can cancel the proposal.
    function cancel(uint256 proposalId) external override onlyState(proposalId, ProposalState.Pending) {
        Proposal storage p = proposals[proposalId];
        require(p.proposer == msg.sender, "Not the proposer");

        p.canceled = true;
        emit ProposalCanceled(proposalId);
    }

    /// @dev Cast a vote to a proposal
    /// The proposal must be in Active state
    /// A node can only vote once for a proposal
    /// choice must be one of VoteChoice.
    function castVote(uint256 proposalId, uint8 choice) external override onlyState(proposalId, ProposalState.Active) {
        Proposal storage p = proposals[proposalId];

        // cache quorums to (1) save gas for checkQuorum,
        // (2) prevent any unintended outcome of updating stakingTracker address.
        if (p.quorumCount == 0) {
            (uint256 quorumCount, uint256 quorumPower) = getQuorum(proposalId);
            p.quorumCount = quorumCount;
            p.quorumPower = quorumPower;
        }

        (uint256 gcId, uint256 votes) = getVotes(proposalId, msg.sender);
        require(gcId != 0, "Not a registered voter");
        require(votes > 0, "Not eligible to vote");

        require(
            choice == uint8(VoteChoice.Yes) || choice == uint8(VoteChoice.No) || choice == uint8(VoteChoice.Abstain),
            "Not a valid choice"
        );

        require(!p.receipts[gcId].hasVoted, "Already voted");
        p.receipts[gcId].hasVoted = true;
        p.receipts[gcId].choice = choice;
        p.receipts[gcId].votes = votes;

        incrementTally(proposalId, choice, votes);
        p.voters.push(gcId);

        emit VoteCast(msg.sender, proposalId, choice, votes, Strings.toHexString(gcId, 32));
    }

    function incrementTally(uint256 proposalId, uint8 choice, uint256 votes) private {
        Proposal storage p = proposals[proposalId];
        if (choice == uint8(VoteChoice.Yes)) {
            p.totalYes += votes;
        } else if (choice == uint8(VoteChoice.No)) {
            p.totalNo += votes;
        } else if (choice == uint8(VoteChoice.Abstain)) {
            p.totalAbstain += votes;
        }
    }

    /// @dev Queue a passed proposal
    /// The proposal must be in Passed state
    /// Current block must be before `queueDeadline` of this proposal
    /// If secretary is null, any GC with at least 1 vote can queue.
    /// Otherwise only secretary can queue.
    function queue(
        uint256 proposalId
    ) external override onlyState(proposalId, ProposalState.Passed) onlyExecutor(proposalId) {
        Proposal storage p = proposals[proposalId];
        require(p.targets.length > 0, "Proposal has no action");

        p.eta = block.number + execDelay;
        p.execDeadline = p.eta + execTimeout;
        p.queued = true;

        emit ProposalQueued(proposalId, p.eta);
    }

    /// @dev Execute a queued proposal
    /// The proposal must be in Queued state
    /// Current block must be after `eta` and before `execDeadline` of this proposal
    /// If secretary is null, any GC with at least 1 vote can execute.
    /// Otherwise only secretary can execute.
    function execute(
        uint256 proposalId
    ) external payable override onlyState(proposalId, ProposalState.Queued) onlyExecutor(proposalId) {
        Proposal storage p = proposals[proposalId];
        require(block.number >= p.eta, "Not yet executable");

        for (uint256 i = 0; i < p.targets.length; i++) {
            (bool success, bytes memory result) = p.targets[i].call{value: p.values[i]}(p.calldatas[i]);
            handleCallResult(success, result);
        }

        p.executed = true;

        emit ProposalExecuted(proposalId);
    }

    function handleCallResult(bool success, bytes memory result) private pure {
        if (success) {
            return;
        }

        if (result.length == 0) {
            // Call failed without message.
            revert("Transaction failed");
        } else {
            // https://github.com/OpenZeppelin/openzeppelin-contracts/blob/v4.7.3/contracts/utils/Address.sol
            // Toss the result, which would contain error instances.
            assembly {
                let result_size := mload(result)
                revert(add(32, result), result_size)
            }
        }
    }

    // Governance functions

    /// @dev Update the StakingTracker address
    /// Should not be called if there is an active proposal
    function updateStakingTracker(address newAddr) public override onlyGovernance {
        // Retire expired trackers
        require(newAddr != address(0), "Address is null");
        IStakingTracker(stakingTracker).refreshStake(address(0));
        require(
            IStakingTracker(stakingTracker).getLiveTrackerIds().length == 0,
            "Cannot update tracker when there is an active tracker"
        );
        address oldAddr = stakingTracker;
        stakingTracker = newAddr;
        emit UpdateStakingTracker(oldAddr, newAddr);
    }

    /// @dev Update the secretary account
    /// Must be called by address(this), i.e. via governance proposal.
    function updateSecretary(address newAddr) public override onlyGovernance {
        address oldAddr = secretary;
        secretary = newAddr;
        validateAccessRule();
        emit UpdateSecretary(oldAddr, newAddr);
    }

    /// @dev Update the access rule
    function updateAccessRule(
        bool secretaryPropose,
        bool voterPropose,
        bool secretaryExecute,
        bool voterExecute
    ) public override onlyGovernanceOrSecretary {
        AccessRule storage ar = accessRule;
        ar.secretaryPropose = secretaryPropose;
        ar.voterPropose = voterPropose;
        ar.secretaryExecute = secretaryExecute;
        ar.voterExecute = voterExecute;

        validateAccessRule();

        emit UpdateAccessRule(ar.secretaryPropose, ar.voterPropose, ar.secretaryExecute, ar.voterExecute);
    }

    function validateAccessRule() internal view {
        AccessRule storage ar = accessRule;
        require((ar.secretaryPropose && secretary != address(0)) || ar.voterPropose, "No propose access");
        require((ar.secretaryExecute && secretary != address(0)) || ar.voterExecute, "No execute access");
    }

    /// @dev Update the timing rule
    function updateTimingRule(
        uint256 minVotingDelay,
        uint256 maxVotingDelay,
        uint256 minVotingPeriod,
        uint256 maxVotingPeriod
    ) public override onlyGovernanceOrSecretary {
        TimingRule storage tr = timingRule;
        tr.minVotingDelay = minVotingDelay;
        tr.maxVotingDelay = maxVotingDelay;
        tr.minVotingPeriod = minVotingPeriod;
        tr.maxVotingPeriod = maxVotingPeriod;

        validateTimingRule();

        emit UpdateTimingRule(tr.minVotingDelay, tr.maxVotingDelay, tr.minVotingPeriod, tr.maxVotingPeriod);
    }

    function validateTimingRule() internal view {
        TimingRule storage tr = timingRule;
        require(tr.minVotingDelay >= 1 * DAY, "Invalid timing");
        require(tr.minVotingPeriod >= 1 * DAY, "Invalid timing");
        require(tr.minVotingDelay <= tr.maxVotingDelay, "Invalid timing");
        require(tr.minVotingPeriod <= tr.maxVotingPeriod, "Invalid timing");
    }

    // Getters

    /// @dev The id of the last created proposal
    /// Retrurns 0 if there is no proposal.
    function lastProposalId() external view override returns (uint256) {
        return nextProposalId - 1;
    }

    /// @dev State of a proposal
    function state(uint256 proposalId) public view override returns (ProposalState) {
        Proposal storage p = proposals[proposalId];

        if (p.executed) {
            return ProposalState.Executed;
        } else if (p.canceled) {
            return ProposalState.Canceled;
        } else if (block.number < p.voteStart) {
            return ProposalState.Pending;
        } else if (block.number <= p.voteEnd) {
            return ProposalState.Active;
        } else if (!checkQuorum(proposalId)) {
            return ProposalState.Failed;
        }

        if (!p.queued) {
            if (block.number <= p.queueDeadline || p.targets.length == 0) {
                return ProposalState.Passed;
            } else {
                return ProposalState.Expired;
            }
        } else {
            if (block.number <= p.execDeadline) {
                return ProposalState.Queued;
            } else {
                return ProposalState.Expired;
            }
        }
    }

    /// @dev Check if a proposal is passed
    /// Note that its return value represents the current voting status,
    /// and is subject to change until the voting ends.
    function checkQuorum(uint256 proposalId) public view override returns (bool) {
        Proposal storage p = proposals[proposalId];

        (uint256 quorumCount, uint256 quorumPower) = getQuorum(proposalId);
        uint256 totalVotes = p.totalYes + p.totalNo + p.totalAbstain;
        uint256 quorumYes = p.totalNo + p.totalAbstain + 1; // more than half of all votes

        bool countPass = (p.voters.length >= quorumCount);
        bool powerPass = (totalVotes >= quorumPower);
        bool approval = (p.totalYes >= quorumYes);

        return ((countPass || powerPass) && approval);
    }

    /// @dev Calculate count and power quorums for a proposal
    function getQuorum(uint256 proposalId) private view returns (uint256 quorumCount, uint256 quorumPower) {
        Proposal storage p = proposals[proposalId];
        if (p.quorumCount != 0) {
            // return cached numbers
            return (p.quorumCount, p.quorumPower);
        }

        (, , , uint256 totalVotes, uint256 numEligible) = IStakingTracker(p.stakingTracker).getTrackerSummary(
            p.trackerId
        );

        quorumCount = (numEligible + 2) / 3; // more than or equal to 1/3 of all GC members
        quorumPower = (totalVotes + 2) / 3; // more than or equal to 1/3 of all voting powers
        return (quorumCount, quorumPower);
    }

    /// @dev Resolve the voter account into its gcId and voting powers
    /// Returns the currently assigned gcId. Returns the voting powers
    /// effective at the given proposal. Returns zero gcId and 0 votes
    /// if the voter account is not assigned to any eligible GC.
    ///
    /// @param proposalId  The proposal id
    /// @return gcId    The gcId assigned to this voter account
    /// @return votes   The amount of voting powers the voter account represents
    function getVotes(uint256 proposalId, address voter) public view override returns (uint256 gcId, uint256 votes) {
        Proposal storage p = proposals[proposalId];

        gcId = IStakingTracker(p.stakingTracker).voterToGCId(voter);
        (, votes) = IStakingTracker(p.stakingTracker).getTrackedGC(p.trackerId, gcId);
    }

    /// @dev General contents of a proposal
    function getProposalContent(
        uint256 proposalId
    ) external view override returns (uint256 id, address proposer, string memory description) {
        Proposal storage p = proposals[proposalId];
        return (proposalId, p.proposer, p.description);
    }

    /// @dev Transactions in a proposal
    /// signatures is Array of empty strings; for compatibility with OpenZeppelin
    function getActions(
        uint256 proposalId
    )
        external
        view
        override
        returns (
            address[] memory targets,
            uint256[] memory values,
            string[] memory signatures,
            bytes[] memory calldatas
        )
    {
        Proposal storage p = proposals[proposalId];
        return (p.targets, p.values, new string[](p.targets.length), p.calldatas);
    }

    /// @dev Timing and state related properties of a proposal
    function getProposalSchedule(
        uint256 proposalId
    )
        external
        view
        override
        returns (
            uint256 voteStart,
            uint256 voteEnd,
            uint256 queueDeadline,
            uint256 eta,
            uint256 execDeadline,
            bool canceled,
            bool queued,
            bool executed
        )
    {
        Proposal storage p = proposals[proposalId];
        return (p.voteStart, p.voteEnd, p.queueDeadline, p.eta, p.execDeadline, p.canceled, p.queued, p.executed);
    }

    /// @dev Vote counting related properties of a proposal
    function getProposalTally(
        uint256 proposalId
    )
        external
        view
        override
        returns (
            uint256 totalYes,
            uint256 totalNo,
            uint256 totalAbstain,
            uint256 quorumCount,
            uint256 quorumPower,
            uint256[] memory voters
        )
    {
        Proposal storage p = proposals[proposalId];
        (quorumCount, quorumPower) = getQuorum(proposalId);
        return (p.totalYes, p.totalNo, p.totalAbstain, quorumCount, quorumPower, p.voters);
    }

    /// @dev Individual vote receipt
    function getReceipt(
        uint256 proposalId,
        uint256 gcId
    ) external view override returns (bool hasVoted, uint8 choice, uint256 votes) {
        Proposal storage p = proposals[proposalId];
        Receipt storage r = p.receipts[gcId];
        return (r.hasVoted, r.choice, r.votes);
    }

    function getTrackerSummary(
        uint256 proposalId
    )
        external
        view
        override
        returns (uint256 trackStart, uint256 trackEnd, uint256 numGCs, uint256 totalVotes, uint256 numEligible)
    {
        Proposal storage p = proposals[proposalId];
        return IStakingTracker(p.stakingTracker).getTrackerSummary(p.trackerId);
    }

    function getAllTrackedGCs(
        uint256 proposalId
    ) external view override returns (uint256[] memory gcIds, uint256[] memory gcBalances, uint256[] memory gcVotes) {
        Proposal storage p = proposals[proposalId];
        return IStakingTracker(p.stakingTracker).getAllTrackedGCs(p.trackerId);
    }

    function voterToGCId(address voter) external view override returns (uint256 gcId) {
        return IStakingTracker(stakingTracker).voterToGCId(voter);
    }

    function gcIdToVoter(uint256 gcId) external view override returns (address voter) {
        return IStakingTracker(stakingTracker).gcIdToVoter(gcId);
    }
}
