// Copyright 2024 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import "openzeppelin-contracts-5.0/access/Ownable.sol";
import "./IStakingTrackerV2.sol";

contract StakingTrackerV2 is IStakingTrackerV2, Ownable {
    /* ========== STATE VARIABLES ========== */

    struct Tracker {
        // Block range for tracking.
        // Updates only occur when: trackStart <= block.number < trackEnd
        uint256 trackStart;
        uint256 trackEnd;
        // Governance Council (GC) tracking
        // Set at createTracker() and remains immutable
        uint256[] gcIds;
        mapping(uint256 => bool) gcExists;
        mapping(address => uint256) stakingToGCId;
        // Balance tracking
        // Initially set at createTracker(), updated via refreshStake() until trackEnd
        mapping(address => uint256) stakingBalances; // All staking contract balances (CnStaking + CLPool)
        mapping(uint256 => uint256) cnStakingBalances; // CnStaking-specific balances per GC
        mapping(address => bool) isCLPool; // Identifies CnStaking vs CLPool contracts. Mark CLPool to save gas
        mapping(uint256 => uint256) gcBalances; // Total GC balances (CnStaking + CLPool combined)
        // Voting power tracking
        // Updated via refreshStake() until trackEnd
        mapping(uint256 => uint256) gcVotes; // Voting power per GC
        uint256 totalVotes; // Sum of all GC voting power
        uint256 numEligible; // Count of GCs meeting minimum stake requirement
    }

    // Store tracker objects
    mapping(uint256 => Tracker) internal trackers; // indexed by trackerId
    uint256[] internal allTrackerIds; // append-only list of trackerIds
    uint256[] internal liveTrackerIds; // trackerIds with block.number < trackEnd. Not in order.

    // 1-to-1 mapping between gcId and voter account
    mapping(uint256 => address) public override gcIdToVoter;
    mapping(address => uint256) public override voterToGCId;

    /* ========== CONSTANTS ========== */

    function CONTRACT_TYPE()
        external
        view
        virtual
        override
        returns (string memory)
    {
        return "StakingTracker";
    }

    /// @dev It remains 1 for the compatibility with the previous version.
    function VERSION() external view virtual override returns (uint256) {
        return 1;
    }

    function ADDRESS_BOOK_ADDRESS()
        public
        view
        virtual
        override
        returns (address)
    {
        return 0x0000000000000000000000000000000000000400;
    }

    function REGISTRY_ADDRESS() public view virtual override returns (address) {
        return 0x0000000000000000000000000000000000000401;
    }

    function MIN_STAKE() public view virtual override returns (uint256) {
        return 5000000 ether;
    }

    /* ========== CONSTRUCTOR ========== */

    constructor(address initialOwner) Ownable(initialOwner) {}

    /* ========== INITIALIZER ========== */

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

        _populateTracker(trackerId);

        emit CreateTracker(trackerId, trackStart, trackEnd, tracker.gcIds);
        return trackerId;
    }

    function _populateTracker(uint256 trackerId) private {
        _populateFromAddressBook(trackerId);
        _populateFromCLRegistry(trackerId);

        _calcAllVotes(trackerId);
    }

    /// @dev Populate a tracker with staking balances from AddressBook
    function _populateFromAddressBook(uint256 trackerId) private {
        Tracker storage tracker = trackers[trackerId];

        (, address[] memory stakingContracts, ) = _getAddressBookLists();

        for (uint256 i = 0; i < stakingContracts.length; i++) {
            address staking = stakingContracts[i];
            (
                bool isV2,
                uint256 balance,
                uint256 gcId,
                address stakingTracker,

            ) = readCnStaking(staking);

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
            tracker.cnStakingBalances[gcId] += balance;
            tracker.gcBalances[gcId] += balance;
        }
    }

    /// @dev Populate a tracker with staking balances from CLRegistry
    function _populateFromCLRegistry(uint256 trackerId) private {
        Tracker storage tracker = trackers[trackerId];

        (
            uint256[] memory gcIds,
            address[] memory pools
        ) = _getCLRegistryLists();

        for (uint256 i = 0; i < pools.length; i++) {
            (uint256 gcId, address pool) = (gcIds[i], pools[i]);
            (address stakingTracker, uint256 balance) = readCLPool(pool);

            if (!tracker.gcExists[gcId]) {
                // Ignore non-existent GCs
                continue;
            }
            if (stakingTracker != address(this)) {
                // Ignore CLPool that does not point to this StakingTracker.
                continue;
            }

            tracker.isCLPool[pool] = true;
            tracker.stakingToGCId[pool] = gcId;
            tracker.stakingBalances[pool] = balance;
            tracker.gcBalances[gcId] += balance;
        }
    }

    /// @dev Populate a tracker with voting powers
    function _calcAllVotes(uint256 trackerId) private {
        Tracker storage tracker = trackers[trackerId];

        uint256 numEligible = 0;
        uint256 totalVotes = 0;

        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            if (tracker.cnStakingBalances[gcId] >= MIN_STAKE()) {
                numEligible++;
            }
        }
        tracker.numEligible = numEligible;

        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            uint256 votes = _calcVotes(trackerId, gcId);
            tracker.gcVotes[gcId] = votes;
            totalVotes += votes;
        }

        tracker.totalVotes = totalVotes; // only write final result to save gas
    }

    /* ========== UPDATE TRACKER ========== */

    /// @dev Re-evaluate Tracker contents related to the staking contract
    /// Anyone can call this function, but `staking` must be a staking contract
    /// registered in tracker.
    function refreshStake(address staking) external virtual override {
        uint256 i = 0;
        while (i < liveTrackerIds.length) {
            uint256 currId = liveTrackerIds[i];

            // Remove expired tracker as soon as we discover it
            if (!_isTrackerLive(currId)) {
                uint256 lastId = liveTrackerIds[liveTrackerIds.length - 1];
                liveTrackerIds[i] = lastId;
                liveTrackerIds.pop();
                emit RetireTracker(currId);
                continue;
            }

            _updateTracker(currId, staking);
            i++;
        }
    }

    /// @dev Re-evaluate balances and subsequently voting power
    function _updateTracker(uint256 trackerId, address staking) private {
        Tracker storage tracker = trackers[trackerId];

        uint256 gcId = tracker.stakingToGCId[staking];
        if (gcId == 0) {
            return;
        }

        // Update general staking balances and GC balances and determine if recalculating all votes is necessary
        (uint256 newBalance, bool recalcAllVotes) = _updateBalances(
            trackerId,
            gcId,
            staking
        );

        if (recalcAllVotes) {
            _updateAllVotes(trackerId, gcId);
        } else {
            _updateVotes(trackerId, gcId);
        }

        emit RefreshStake(
            trackerId,
            gcId,
            staking,
            newBalance,
            tracker.gcBalances[gcId],
            tracker.gcVotes[gcId],
            tracker.totalVotes
        );
    }

    /// @dev Update general staking balances and GC balances and determine if recalculating all votes is necessary
    function _updateBalances(
        uint256 trackerId,
        uint256 gcId,
        address cnStakingOrCLPool
    ) private returns (uint256 newBalance, bool recalcAllVotes) {
        Tracker storage tracker = trackers[trackerId];
        bool isCL = tracker.isCLPool[cnStakingOrCLPool];

        uint256 oldBalance = tracker.stakingBalances[cnStakingOrCLPool];
        newBalance = _readStakingBalance(isCL, cnStakingOrCLPool);

        tracker.stakingBalances[cnStakingOrCLPool] = newBalance;
        tracker.gcBalances[gcId] =
            tracker.gcBalances[gcId] -
            oldBalance +
            newBalance;

        if (!isCL) {
            uint256 oldCnStakingBalance = tracker.cnStakingBalances[gcId];
            uint256 newCnStakingBalance = oldCnStakingBalance -
                oldBalance +
                newBalance;
            tracker.cnStakingBalances[gcId] = newCnStakingBalance;
            recalcAllVotes = _shouldRecalcAllVotes(
                oldCnStakingBalance,
                newCnStakingBalance
            );
        }

        return (newBalance, recalcAllVotes);
    }

    /// @dev Update the voting power of every GC
    function _updateAllVotes(uint256 trackerId, uint256 gcId) private {
        Tracker storage tracker = trackers[trackerId];

        bool wasEligible = tracker.cnStakingBalances[gcId] < MIN_STAKE();
        if (wasEligible) {
            // eligible -> not eligible
            tracker.numEligible -= 1;
        } else {
            // not eligible -> eligible
            tracker.numEligible += 1;
        }

        _recalcAllVotes(trackerId);
    }

    /// @dev Update votes
    function _updateVotes(uint256 trackerId, uint256 gcId) private {
        Tracker storage tracker = trackers[trackerId];

        uint256 oldVotes = tracker.gcVotes[gcId];
        uint256 newVotes = _calcVotes(trackerId, gcId);

        tracker.gcVotes[gcId] = newVotes;
        tracker.totalVotes = tracker.totalVotes - oldVotes + newVotes;
    }

    /// @dev Determine if recalculating all votes is necessary
    function _shouldRecalcAllVotes(
        uint256 oldCnStakingBalance,
        uint256 newCnStakingBalance
    ) private view returns (bool) {
        bool wasEligible = oldCnStakingBalance >= MIN_STAKE();
        bool isEligible = newCnStakingBalance >= MIN_STAKE();
        return wasEligible != isEligible;
    }

    /// @dev Recalculate votes with new numEligible
    /// Assume numEligible is already updated at tracker.numEligible
    function _recalcAllVotes(uint256 trackerId) private {
        Tracker storage tracker = trackers[trackerId];

        uint256 totalVotes = tracker.totalVotes;
        for (uint256 i = 0; i < tracker.gcIds.length; i++) {
            uint256 gcId = tracker.gcIds[i];
            uint256 oldVotes = tracker.gcVotes[gcId];
            uint256 newVotes = _calcVotes(trackerId, gcId);

            if (oldVotes != newVotes) {
                tracker.gcVotes[gcId] = newVotes;
                totalVotes -= oldVotes;
                totalVotes += newVotes;
            }
        }

        tracker.totalVotes = totalVotes; // only write final result to save gas
    }

    /* ========== UPDATE VOTER ========== */

    /// @dev Re-evaluate voter account mapping related to the staking contract
    /// Anyone can call this function, but `staking` must be a staking contract
    /// registered to the current AddressBook.
    ///
    /// Updates the voter account of the GC of the `staking` with respect to
    /// the current AddressBook.
    ///
    /// If the GC already had a voter account, the account will be unregistered.
    /// If the new voter account is already appointed for another GC,
    /// this function reverts.
    function refreshVoter(address staking) external virtual override {
        (, address[] memory stakingContracts, ) = _getAddressBookLists();
        bool stakingInAddressBook = false;
        for (uint256 i = 0; i < stakingContracts.length; i++) {
            if (stakingContracts[i] == staking) {
                stakingInAddressBook = true;
                break;
            }
        }
        require(stakingInAddressBook, "Not a staking contract");

        (bool isV2, , uint256 gcId, , address newVoter) = readCnStaking(
            staking
        );
        require(isV2, "Invalid CnStaking contract");

        _updateVoter(gcId, newVoter);

        emit RefreshVoter(gcId, staking, newVoter);
    }

    function _updateVoter(uint256 gcId, address newVoter) private {
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

    /* ========== HELPER FUNCTIONS ========== */

    /// @dev Calculate voting power from staking amounts.
    /// One integer vote is granted for each MIN_STAKE() balance. But the number of votes
    /// is at most ([number of eligible GCs] - 1).
    function _calcVotes(
        uint256 trackerId,
        uint256 gcId
    ) internal view returns (uint256) {
        Tracker storage tracker = trackers[trackerId];
        uint256 cnStakingBalance = tracker.cnStakingBalances[gcId];
        uint256 balance = tracker.gcBalances[gcId];
        uint256 numEligible = tracker.numEligible;

        if (cnStakingBalance < MIN_STAKE()) {
            return 0;
        }

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

    /// @dev Query the 3-tuples (node, staking, reward) from AddressBook
    function _getAddressBookLists()
        internal
        view
        returns (
            address[] memory nodeIds,
            address[] memory stakingContracts,
            address[] memory rewardAddrs
        )
    {
        (
            nodeIds,
            stakingContracts,
            rewardAddrs /* kgf */ /* kir */,
            ,

        ) = IAddressBook(ADDRESS_BOOK_ADDRESS()).getAllAddressInfo();
        require(
            nodeIds.length == stakingContracts.length &&
                nodeIds.length == rewardAddrs.length,
            "Invalid data"
        );
    }

    /// @dev Get the address of the CLRegistry
    function _getCLRegistryAddress() internal view returns (address) {
        return IRegistry(REGISTRY_ADDRESS()).getActiveAddr("CLRegistry");
    }

    /// @dev Get the address of the Wrapped Kaia
    function _getWKaiaAddress() internal view returns (address) {
        return IRegistry(REGISTRY_ADDRESS()).getActiveAddr("WrappedKaia");
    }

    /// @dev Query the 2-tuples (gcId, pool) from CLRegistry
    /// Returns empty arrays if CLRegistry is not found
    function _getCLRegistryLists()
        internal
        view
        returns (uint256[] memory gcIds, address[] memory pools)
    {
        address clRegistry = _getCLRegistryAddress();
        if (clRegistry == address(0)) {
            return (gcIds, pools);
        }
        (, gcIds, pools) = ICLRegistry(clRegistry).getAllCLs();
        require(gcIds.length == pools.length, "Invalid data");
    }

    /// @dev Read staking balance from staking contract
    function _readStakingBalance(
        bool isCL,
        address staking
    ) internal view returns (uint256 balance) {
        if (isCL) {
            (, balance) = readCLPool(staking);
        } else {
            (, balance, , , ) = readCnStaking(staking);
        }
    }

    /// @dev Determine if given tracker is updatable with respect to current block.
    function _isTrackerLive(uint256 trackerId) internal view returns (bool) {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.trackStart <= block.number &&
            block.number < tracker.trackEnd);
    }

    /// @dev Test if the given contract is a CnStakingV2 instance
    /// Does not check if the contract is registered in AddressBook.
    function isCnStakingV2(address staking) public view returns (bool) {
        bool ok;
        bytes memory out;

        (ok, out) = staking.staticcall(
            abi.encodeWithSignature("CONTRACT_TYPE()")
        );
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
        returns (
            bool isV2,
            uint256 effectiveBalance,
            uint256 gcId,
            address stakingTracker,
            address voterAddress
        )
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

    /// @dev Read balance from CLPool
    function readCLPool(
        address pool
    ) public view virtual returns (address stakingTracker, uint256 balance) {
        address wKaia = _getWKaiaAddress();
        if (wKaia != address(0)) {
            balance = IERC20(wKaia).balanceOf(pool);
        }
        stakingTracker = ICLPool(pool).stakingTracker();
    }

    /* ========== GETTERS ========== */

    function getLastTrackerId() public view override returns (uint256) {
        return allTrackerIds.length;
    }

    function getAllTrackerIds()
        external
        view
        override
        returns (uint256[] memory)
    {
        return allTrackerIds;
    }

    function getLiveTrackerIds()
        external
        view
        override
        returns (uint256[] memory)
    {
        return liveTrackerIds;
    }

    function getTrackerSummary(
        uint256 trackerId
    )
        public
        view
        override
        returns (
            uint256 trackStart,
            uint256 trackEnd,
            uint256 numGCs,
            uint256 totalVotes,
            uint256 numEligible
        )
    {
        Tracker storage tracker = trackers[trackerId];
        return (
            tracker.trackStart,
            tracker.trackEnd,
            tracker.gcIds.length,
            tracker.totalVotes,
            tracker.numEligible
        );
    }

    function getTrackedGC(
        uint256 trackerId,
        uint256 gcId
    ) external view override returns (uint256 gcBalance, uint256 gcVotes) {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.gcBalances[gcId], tracker.gcVotes[gcId]);
    }

    function getTrackedGCBalance(
        uint256 trackerId,
        uint256 gcId
    ) external view override returns (uint256, uint256) {
        Tracker storage tracker = trackers[trackerId];
        return (tracker.cnStakingBalances[gcId], tracker.gcBalances[gcId]);
    }

    function getAllTrackedGCs(
        uint256 trackerId
    )
        public
        view
        override
        returns (
            uint256[] memory gcIds,
            uint256[] memory gcBalances,
            uint256[] memory gcVotes
        )
    {
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

    function stakingToGCId(
        uint256 trackerId,
        address staking
    ) external view override returns (uint256) {
        Tracker storage tracker = trackers[trackerId];
        return tracker.stakingToGCId[staking];
    }

    function isCLPool(
        uint256 trackerId,
        address staking
    ) external view override returns (bool) {
        Tracker storage tracker = trackers[trackerId];
        return tracker.isCLPool[staking];
    }
}

interface IAddressBook {
    function getAllAddressInfo()
        external
        view
        returns (
            address[] memory,
            address[] memory,
            address[] memory,
            address,
            address
        );
}

interface ICnStakingV2 {
    function VERSION() external view returns (uint256);

    function rewardAddress() external view returns (address);

    function stakingTracker() external view returns (address);

    function voterAddress() external view returns (address);

    function gcId() external view returns (uint256);

    function unstaking() external view returns (uint256);
}

interface IRegistry {
    function getActiveAddr(string memory name) external view returns (address);
}

interface ICLRegistry {
    function getAllCLs()
        external
        view
        returns (address[] memory, uint256[] memory, address[] memory);
}

interface ICLPool {
    function stakingTracker() external view returns (address);
}

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);
}
