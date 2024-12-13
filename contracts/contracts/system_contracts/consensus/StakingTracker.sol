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

import "@openzeppelin/contracts/access/Ownable.sol";
import "./IStakingTracker.sol";

contract StakingTracker is IStakingTracker, Ownable {
    struct Tracker {
        // Tracked block range.
        // Balance changes are only updated if trackStart <= block.number < trackEnd.
        uint256 trackStart;
        uint256 trackEnd;
        // List of eligible GCs and their staking addresses.
        // Determined at crateTracker() and does not change.
        uint256[] gcIds;
        mapping(uint256 => bool) gcExists;
        mapping(address => uint256) stakingToGCId;
        // Balances and voting powers.
        // First collected at crateTracker() and updated at refreshStake() until trackEnd.
        mapping(address => uint256) stakingBalances; // staking address balances
        mapping(uint256 => uint256) gcBalances; // consolidated GC balances
        mapping(uint256 => uint256) gcVotes; // GC voting powers
        uint256 totalVotes;
        uint256 numEligible;
    }

    // Store tracker objects
    mapping(uint256 => Tracker) internal trackers; // indexed by trackerId
    uint256[] internal allTrackerIds; // append-only list of trackerIds
    uint256[] internal liveTrackerIds; // trackerIds with block.number < trackEnd. Not in order.

    // 1-to-1 mapping between gcId and voter account
    mapping(uint256 => address) public override gcIdToVoter;
    mapping(address => uint256) public override voterToGCId;

    // Constants
    function CONTRACT_TYPE() external view virtual override returns (string memory) {
        return "StakingTracker";
    }

    function VERSION() external view virtual override returns (uint256) {
        return 1;
    }

    function ADDRESS_BOOK_ADDRESS() public view virtual override returns (address) {
        return 0x0000000000000000000000000000000000000400;
    }

    function MIN_STAKE() public view virtual override returns (uint256) {
        return 5000000 ether;
    }

    // Mutators

    /// @dev Creates a new Tracker and populate initial values from AddressBook
    /// Only allowed to the contract owner.
    function createTracker(
        uint256 trackStart,
        uint256 trackEnd
    ) public virtual override onlyOwner returns (uint256 trackerId) {
        trackerId = getLastTrackerId() + 1;
        allTrackerIds.push(trackerId);
        liveTrackerIds.push(trackerId);

        Tracker storage tracker = trackers[trackerId];
        tracker.trackStart = trackStart;
        tracker.trackEnd = trackEnd;

        populateFromAddressBook(trackerId);
        calcAllVotes(trackerId);

        emit CreateTracker(trackerId, trackStart, trackEnd, tracker.gcIds);
        return trackerId;
    }

    /// @dev Populate a tracker with staking balances from AddressBook
    function populateFromAddressBook(uint256 trackerId) internal {
        Tracker storage tracker = trackers[trackerId];

        (, address[] memory stakingContracts, ) = getAddressBookLists();

        for (uint256 i = 0; i < stakingContracts.length; i++) {
            address staking = stakingContracts[i];

            (bool isV2, uint256 balance, uint256 gcId, address stakingTracker, ) = readCnStaking(staking);
            if (!isV2) {
                // Ignore V1 contract
                continue;
            }
            if (stakingTracker != address(this)) {
                // Ignore CnStaking that does not point to this StakingTracker.
                // Hinders an attack where the CnStaking evades real-time voting
                // power calculation via staking withdrawal.
                continue;
            }

            if (!tracker.gcExists[gcId]) {
                tracker.gcExists[gcId] = true;
                tracker.gcIds.push(gcId);
            }

            tracker.stakingToGCId[staking] = gcId;
            tracker.stakingBalances[staking] = balance;
            tracker.gcBalances[gcId] += balance;
        }
    }

    /// @dev Populate a tracker with voting powers
    function calcAllVotes(uint256 trackerId) internal {
        Tracker storage tracker = trackers[trackerId];
        uint256 numEligible = 0;
        uint256 totalVotes = 0;

        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            if (tracker.gcBalances[gcId] >= MIN_STAKE()) {
                numEligible++;
            }
        }
        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            uint256 balance = tracker.gcBalances[gcId];
            uint256 votes = calcVotes(numEligible, balance);
            tracker.gcVotes[gcId] = votes;
            totalVotes += votes;
        }

        tracker.numEligible = numEligible;
        tracker.totalVotes = totalVotes; // only write final result to save gas
    }

    /// @dev Re-evaluate Tracker contents related to the staking contract
    /// Anyone can call this function, but `staking` must be a staking contract
    /// registered in tracker.
    function refreshStake(address staking) external virtual override {
        uint256 i = 0;
        while (i < liveTrackerIds.length) {
            uint256 currId = liveTrackerIds[i];

            // Remove expired tracker as soon as we discover it
            if (!isTrackerLive(currId)) {
                uint256 lastId = liveTrackerIds[liveTrackerIds.length - 1];
                liveTrackerIds[i] = lastId;
                liveTrackerIds.pop();
                emit RetireTracker(currId);
                continue;
            }

            updateTracker(currId, staking);
            i++;
        }
    }

    /// @dev Re-evalute balances and subsequently voting power
    function updateTracker(uint256 trackerId, address staking) private {
        Tracker storage tracker = trackers[trackerId];

        // Resolve GC
        uint256 gcId = tracker.stakingToGCId[staking];
        if (gcId == 0) {
            return;
        }

        // Update balance
        uint256 oldBalance = tracker.stakingBalances[staking];
        (, uint256 newBalance, , , ) = readCnStaking(staking);
        tracker.stakingBalances[staking] = newBalance;

        uint256 oldGcBalance = tracker.gcBalances[gcId];
        tracker.gcBalances[gcId] -= oldBalance;
        tracker.gcBalances[gcId] += newBalance;
        uint256 newGcBalance = tracker.gcBalances[gcId];

        // Update vote cap if necessary
        recalcAllVotesIfNeeded(trackerId, oldGcBalance, newGcBalance);

        // Update votes
        uint256 oldVotes = tracker.gcVotes[gcId];
        uint256 newVotes = calcVotes(tracker.numEligible, newGcBalance);
        tracker.gcVotes[gcId] = newVotes;
        tracker.totalVotes -= oldVotes;
        tracker.totalVotes += newVotes;

        emit RefreshStake(trackerId, gcId, staking, newBalance, newGcBalance, newVotes, tracker.totalVotes);
    }

    function recalcAllVotesIfNeeded(uint256 trackerId, uint256 oldGcBalance, uint256 newGcBalance) internal {
        Tracker storage tracker = trackers[trackerId];

        bool wasEligible = oldGcBalance >= MIN_STAKE();
        bool isEligible = newGcBalance >= MIN_STAKE();
        if (wasEligible != isEligible) {
            if (wasEligible) {
                // eligible -> not eligible
                tracker.numEligible -= 1;
            } else {
                // not eligible -> eligible
                tracker.numEligible += 1;
            }
            recalcAllVotes(trackerId);
        }
    }

    /// @dev Recalculate votes with new numEligible
    function recalcAllVotes(uint256 trackerId) internal {
        Tracker storage tracker = trackers[trackerId];

        uint256 totalVotes = tracker.totalVotes;
        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            uint256 gcBalance = tracker.gcBalances[gcId];
            uint256 oldVotes = tracker.gcVotes[gcId];
            uint256 newVotes = calcVotes(tracker.numEligible, gcBalance);

            if (oldVotes != newVotes) {
                tracker.gcVotes[gcId] = newVotes;
                totalVotes -= oldVotes;
                totalVotes += newVotes;
            }
        }

        tracker.totalVotes = totalVotes; // only write final result to save gas
    }

    /// @dev Re-evaluate voter account mapping related to the staking contract
    /// Anyone can call this function, but `staking` must be a staking contract
    /// registered to the current AddressBook.
    ///
    /// Updates the voter account of the GC of the `staking` with respect to
    /// the corrent AddressBook.
    ///
    /// If the GC already had a voter account, the account will be unregistered.
    /// If the new voter account is already appointed for another GC,
    /// this function reverts.
    function refreshVoter(address staking) external virtual override {
        (, address[] memory stakingContracts, ) = getAddressBookLists();
        bool stakingInAddressBook = false;
        for (uint256 i = 0; i < stakingContracts.length; i++) {
            if (stakingContracts[i] == staking) {
                stakingInAddressBook = true;
                break;
            }
        }
        require(stakingInAddressBook, "Not a staking contract");

        (bool isV2, , uint256 gcId, , address newVoter) = readCnStaking(staking);
        require(isV2, "Invalid CnStaking contract");

        updateVoter(gcId, newVoter);

        emit RefreshVoter(gcId, staking, newVoter);
    }

    function updateVoter(uint256 gcId, address newVoter) internal {
        // Unlink existing two-way mapping
        address oldVoter = gcIdToVoter[gcId];
        if (oldVoter != address(0)) {
            voterToGCId[oldVoter] = 0;
            gcIdToVoter[gcId] = address(0);
        }

        // Create new mapping
        if (newVoter != address(0)) {
            require(voterToGCId[newVoter] == 0, "Voter address already taken");
            voterToGCId[newVoter] = gcId;
            gcIdToVoter[gcId] = newVoter;
        }
    }

    // Helper fucntions

    /// @dev Query the 3-tuples (node, staking, reward) from AddressBook
    function getAddressBookLists()
        internal
        view
        returns (address[] memory nodeIds, address[] memory stakingContracts, address[] memory rewardAddrs)
    {
        (nodeIds, stakingContracts, rewardAddrs /* kgf */ /* kir */, , ) = IAddressBook(ADDRESS_BOOK_ADDRESS())
            .getAllAddressInfo();
        require(nodeIds.length == stakingContracts.length && nodeIds.length == rewardAddrs.length, "Invalid data");
    }

    /// @dev Test if the given contract is a CnStakingV2 instance
    /// Does not check if the contract is registered in AddressBook.
    function isCnStakingV2(address staking) public view returns (bool) {
        bool ok;
        bytes memory out;

        (ok, out) = staking.staticcall(abi.encodeWithSignature("CONTRACT_TYPE()"));
        if (!ok || out.length == 0) {
            return false;
        }
        string memory _type = abi.decode(out, (string));
        if (keccak256(bytes(_type)) != keccak256(bytes("CnStakingContract"))) {
            return false;
        }

        (ok, out) = staking.staticcall(abi.encodeWithSignature("VERSION()"));
        if (!ok || out.length == 0) {
            return false;
        }
        uint256 _version = abi.decode(out, (uint256));
        if (_version < 2) {
            return false;
        }

        return true;
    }

    /// @dev Read various fields from a CnStaking contract
    function readCnStaking(
        address staking
    )
        public
        view
        virtual
        returns (bool isV2, uint256 effectiveBalance, uint256 gcId, address stakingTracker, address voterAddress)
    {
        if (isCnStakingV2(staking)) {
            return (
                true,
                staking.balance - ICnStakingV2(staking).unstaking(),
                ICnStakingV2(staking).gcId(),
                ICnStakingV2(staking).stakingTracker(),
                ICnStakingV2(staking).voterAddress()
            );
        }
        return (false, 0, 0, address(0), address(0));
    }

    /// @dev Calculate voting power from staking amounts.
    /// One integer vote is granted for each MIN_STAKE() balance. But the number of votes
    /// is at most ([number of eligible GCs] - 1).
    function calcVotes(uint256 numEligible, uint256 balance) private view returns (uint256) {
        uint256 voteCap = 1;
        if (numEligible > 1) {
            voteCap = numEligible - 1;
        }

        uint256 votes = balance / MIN_STAKE();
        if (votes > voteCap) {
            votes = voteCap;
        }
        return votes;
    }

    /// @dev Determine if given tracker is updatable with respect to current block.
    function isTrackerLive(uint256 trackerId) private view returns (bool) {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.trackStart <= block.number && block.number < tracker.trackEnd);
    }

    // Getter functions

    function getLastTrackerId() public view override returns (uint256) {
        return allTrackerIds.length;
    }

    function getAllTrackerIds() external view override returns (uint256[] memory) {
        return allTrackerIds;
    }

    function getLiveTrackerIds() external view override returns (uint256[] memory) {
        return liveTrackerIds;
    }

    function getTrackerSummary(
        uint256 trackerId
    )
        public
        view
        override
        returns (uint256 trackStart, uint256 trackEnd, uint256 numGCs, uint256 totalVotes, uint256 numEligible)
    {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.trackStart, tracker.trackEnd, tracker.gcIds.length, tracker.totalVotes, tracker.numEligible);
    }

    function getTrackedGC(
        uint256 trackerId,
        uint256 gcId
    ) external view override returns (uint256 gcBalance, uint256 gcVotes) {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.gcBalances[gcId], tracker.gcVotes[gcId]);
    }

    function getAllTrackedGCs(
        uint256 trackerId
    ) public view override returns (uint256[] memory gcIds, uint256[] memory gcBalances, uint256[] memory gcVotes) {
        Tracker storage tracker = trackers[trackerId];
        uint256 numGCs = tracker.gcIds.length;
        gcIds = tracker.gcIds;

        gcBalances = new uint256[](numGCs);
        gcVotes = new uint256[](numGCs);
        for (uint256 i = 0; i < numGCs; i++) {
            uint256 gcId = tracker.gcIds[i];
            gcBalances[i] = tracker.gcBalances[gcId];
            gcVotes[i] = tracker.gcVotes[gcId];
        }
    }

    function stakingToGCId(uint256 trackerId, address staking) external view override returns (uint256) {
        Tracker storage tracker = trackers[trackerId];
        return tracker.stakingToGCId[staking];
    }
}

interface IAddressBook {
    function getAllAddressInfo()
        external
        view
        returns (address[] memory, address[] memory, address[] memory, address, address);
}

interface ICnStakingV2 {
    function VERSION() external view returns (uint256);

    function rewardAddress() external view returns (address);

    function stakingTracker() external view returns (address);

    function voterAddress() external view returns (address);

    function gcId() external view returns (uint256);

    function unstaking() external view returns (uint256);
}
