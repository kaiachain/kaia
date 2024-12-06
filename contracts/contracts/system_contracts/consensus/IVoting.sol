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

interface IVoting {
    // Types

    enum ProposalState {
        Pending,
        Active,
        Canceled,
        Failed,
        Passed,
        Queued,
        Expired,
        Executed
    }

    enum VoteChoice {
        No,
        Yes,
        Abstain
    }

    struct Receipt {
        bool hasVoted;
        uint8 choice;
        uint256 votes;
    }

    // Events

    /// @dev Emitted when a proposal is created
    /// @param signatures  Array of empty strings; for compatibility with OpenZeppelin
    event ProposalCreated(
        uint256 proposalId,
        address proposer,
        address[] targets,
        uint256[] values,
        string[] signatures,
        bytes[] calldatas,
        uint256 voteStart,
        uint256 voteEnd,
        string description
    );

    /// @dev Emitted when a proposal is canceled
    event ProposalCanceled(uint256 proposalId);

    /// @dev Emitted when a proposal is queued
    /// @param eta  The block number where transaction becomes executable.
    event ProposalQueued(uint256 proposalId, uint256 eta);

    /// @dev Emitted when a proposal is executed
    event ProposalExecuted(uint256 proposalId);

    /// @dev Emitted when a vote is cast
    /// @param reason  An empty string; for compatibility with OpenZeppelin
    event VoteCast(address indexed voter, uint256 proposalId, uint8 choice, uint256 votes, string reason);

    /// @dev Emitted when the StakingTracker is changed
    event UpdateStakingTracker(address oldAddr, address newAddr);

    /// @dev Emitted when the secretary is changed
    event UpdateSecretary(address oldAddr, address newAddr);

    /// @dev Emitted when the AccessRule is changed
    event UpdateAccessRule(bool secretaryPropose, bool voterPropose, bool secretaryExecute, bool voterExecute);

    /// @dev Emitted when the TimingRule is changed
    event UpdateTimingRule(
        uint256 minVotingDelay,
        uint256 maxVotingDelay,
        uint256 minVotingPeriod,
        uint256 maxVotingPeriod
    );

    // Mutators

    function propose(
        string memory description,
        address[] memory targets,
        uint256[] memory values,
        bytes[] memory calldatas,
        uint256 votingDelay,
        uint256 votingPeriod
    ) external returns (uint256 proposalId);

    function cancel(uint256 proposalId) external;

    function castVote(uint256 proposalId, uint8 choice) external;

    function queue(uint256 proposalId) external;

    function execute(uint256 proposalId) external payable;

    function updateStakingTracker(address newAddr) external;

    function updateSecretary(address newAddr) external;

    function updateAccessRule(
        bool secretaryPropose,
        bool voterPropose,
        bool secretaryExecute,
        bool voterExecute
    ) external;

    function updateTimingRule(
        uint256 minVotingDelay,
        uint256 maxVotingDelay,
        uint256 minVotingPeriod,
        uint256 maxVotingPeriod
    ) external;

    // Getters

    function stakingTracker() external view returns (address);

    function secretary() external view returns (address);

    function accessRule()
        external
        view
        returns (bool secretaryPropose, bool voterPropose, bool secretaryExecute, bool voterExecute);

    function timingRule()
        external
        view
        returns (uint256 minVotingDelay, uint256 maxVotingDelay, uint256 minVotingPeriod, uint256 maxVotingPeriod);

    function queueTimeout() external view returns (uint256);

    function execDelay() external view returns (uint256);

    function execTimeout() external view returns (uint256);

    function lastProposalId() external view returns (uint256);

    function state(uint256 proposalId) external view returns (ProposalState);

    function checkQuorum(uint256 proposalId) external view returns (bool);

    function getVotes(uint256 proposalId, address voter) external view returns (uint256, uint256);

    function getProposalContent(
        uint256 proposalId
    ) external view returns (uint256 id, address proposer, string memory description);

    function getActions(
        uint256 proposalId
    )
        external
        view
        returns (
            address[] memory targets,
            uint256[] memory values,
            string[] memory signatures,
            bytes[] memory calldatas
        );

    function getProposalSchedule(
        uint256 proposalId
    )
        external
        view
        returns (
            uint256 voteStart,
            uint256 voteEnd,
            uint256 queueDeadline,
            uint256 eta,
            uint256 execDeadline,
            bool canceled,
            bool queued,
            bool executed
        );

    function getProposalTally(
        uint256 proposalId
    )
        external
        view
        returns (
            uint256 totalYes,
            uint256 totalNo,
            uint256 totalAbstain,
            uint256 quorumCount,
            uint256 quorumPower,
            uint256[] memory voters
        );

    function getReceipt(
        uint256 proposalId,
        uint256 voter
    ) external view returns (bool hasVoted, uint8 choice, uint256 votes);

    function getTrackerSummary(
        uint256 proposalId
    )
        external
        view
        returns (uint256 trackStart, uint256 trackEnd, uint256 numGCs, uint256 totalVotes, uint256 numEligible);

    function getAllTrackedGCs(
        uint256 proposalId
    ) external view returns (uint256[] memory gcIds, uint256[] memory gcBalances, uint256[] memory gcVotes);

    function voterToGCId(address voter) external view returns (uint256 gcId);

    function gcIdToVoter(uint256 gcId) external view returns (address voter);
}
