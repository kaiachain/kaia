// Copyright 2024 The klaytn Authors
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
pragma solidity 0.8.25;

import "../../../libs/ValidContract.sol";
import "../IAddressBook.sol";
import "../PublicDelegation/IKIP163.sol";
import "../PublicDelegation/IPublicDelegationFactory.sol";
import "./ICnStakingV3.sol";
import "./CnStakingV3Storage.sol";
import "openzeppelin-contracts-5.0/access/IAccessControl.sol";
import "openzeppelin-contracts-5.0/access/extensions/AccessControlEnumerable.sol";

// Features
// Note that extra multisig operations are implemented in the `CnStakingV3MultiSig` contract.
// 1. Initialization
//   - Initialize the contract with the given parameters.
//   - Functions
//     - setStakingTracker: Set the initial stakingTracker address
//     - setGCId: Set the gcId
//     - setPublicDelegation: Set the public delegation contract and initialize the contract
//
// 2. Lockup stakes (Initial lockup)
//   - Initial lockup is a set of long-term fixed lockups.
//   - Admin must agree to the conditions for this contract to initialize.
//   - In CnStakingV3MultiSig, the all admins must agree to the conditions.
//   - KLAYs must be deposited for this contract to initialize.
//   - Functions
//     - reviewInitialConditions(): Agree to the initial lockup conditions
//     - depositLockupStakingAndInit(): Deposit required amount
//     - WithdrawLockupStaking: Withdraw unlocked amount
//
// 3. Non-lockup stakes (Delegated stake)
//   - Delegated stakes can be added or removed at any time.
//   - Delegated stakes can be added either by calling delegate() or sending
//     a transaction to this contract with nonzero KAIA (via fallback).
//   - If public delegation is enabled, only the public delegation contract can call staking functions.
//   - It takes STAKE_LOCKUP after withdrawal request to actually take out the KAIA.
//   - Functions
//     - delegate/receive: Stake KAIA to this contract
//     - ApproveStakingWithdrawal: Schedule a withdrawal
//     - CancelApprovedStakingWithdrawal: Cancel a withdrawal request
//     - withdrawApprovedStaking(): Take out the KAIA or cancel an expired withdrawal request.
//
// 4. Re-delegation
//   - Re-delegation is a process to move stakes from one CnStakingV3 contract to another.
//   - It can be initiated only by the public delegation contract.
//   - It's inactive by default and can be toggled by the owner. (isRedelegationEnabled)
//   - To prevent re-delegation hopping, the last re-delegation time is recorded.
//     Can't redelegate within STAKE_LOCKUP after the last re-delegation.
//     For example, let's assume an user redelegated from CnStakingV3 A to CnStakingV3 B.
//     Then, user cannot redelegate from CnStakingV3 B within STAKE_LOCKUP.
//   - In `handleRedelegation`, target CnStakingV3 will stake the received KAIA on behalf of the user using `stakeFor`,
//     and check if the balance is not changed after the call.
//   - Functions
//     - toggleRedelegation: Enable or disable re-delegation, can be called only when public delegation is enabled.
//     - redelegate: Move stakes to the target CnStakingV3 contract
//     - handleRedelegation: Handle re-delegation from departure CnStakingV3 contract
//
// 5. External accounts
//   - Several addresses constitute the identity of this CN.
//   - Among them, RewardAddress can be modified via CnStaking contract.
//   - Functions
//     - UpdateRewardAddress: Setup pendingRewardAddress
//     - acceptRewardAddress(): Request AddressBook to change reward address.
//     - UpdateStakingTracker: Change the StakingTracker contract to report stakes.
//     - UpdateVoterAddress: Change the Voter account and notify to StakingTracker.

// Code organization
// - Modifiers
// - Mutators
//   - Constructor and initializers
//   - Managing operations
//   - Generic facility
//   - Private helpers
//   - Other public functions
// - Getters

contract CnStakingV3 is ICnStakingV3, CnStakingV3Storage, AccessControlEnumerable {
    using EnumerableSet for EnumerableSet.AddressSet;
    using ValidContract for address;

    /* ========== MODIFIERS ========== */

    modifier notNull(address _address) {
        _checkNull(_address);
        _;
    }

    modifier beforeInit() {
        _checkNotInit();
        _;
    }

    modifier afterInit() {
        _checkInit();
        _;
    }

    /* ========== INITIALIZE FUNCTIONS ========== */

    /// @dev Fill in initial values for the contract
    /// Emits a DeployContract event.
    /// @param _owner              The owner of this CN
    /// @param _nodeId             The NodeID of this CN
    /// @param _rewardAddress      The RewardBase of this CN
    /// @param _unlockTime         List of initial lockup deadlines in block timestamp
    /// @param _unlockAmount       List of initial lockup amounts in peb
    constructor(
        address _owner,
        address _nodeId,
        address _rewardAddress,
        uint256[] memory _unlockTime,
        uint256[] memory _unlockAmount
    ) notNull(_nodeId) {
        /// @dev Sanitize initial conditions
        _validInitialConditions(_rewardAddress, _unlockTime, _unlockAmount);

        _grantRole(OPERATOR_ROLE, _owner);
        _grantRole(ADMIN_ROLE, _owner);

        isPublicDelegationEnabled = _rewardAddress == address(0);
        lockupConditions.unlockTime = _unlockTime;
        lockupConditions.unlockAmount = _unlockAmount;

        nodeId = _nodeId;

        rewardAddress = _rewardAddress;

        emit DeployCnStakingV3(CONTRACT_TYPE, _nodeId, _rewardAddress, _unlockTime, _unlockAmount);
    }

    /// @dev Set the initial stakingTracker address
    /// Emits a UpdateStakingTracker event.
    function setStakingTracker(address _tracker) external override beforeInit onlyRole(ADMIN_ROLE) {
        require(_tracker._validStakingTracker(1), "Invalid StakingTracker.");
        stakingTracker = _tracker;
        emit UpdateStakingTracker(_tracker);
    }

    /// @dev Set the gcId
    /// The gcId never changes once initialized.
    /// Emits a UpdateGCId event.
    function setGCId(uint256 _gcId) external override beforeInit onlyRole(ADMIN_ROLE) {
        require(_gcId != 0, "GC ID cannot be zero.");
        gcId = _gcId;
        emit UpdateGCId(_gcId);
    }

    /// @dev Set the public delegation contract.
    /// The public delegation contract must be deployed by the PublicDelegationFactory.
    /// The public delegation address never changes once initialized.
    /// Since it's a public delegation, the reward address also never changes.
    /// Emits a SetPublicDelegation event.
    function setPublicDelegation(
        address _pdFactory,
        bytes memory _pdArgs
    ) external override beforeInit notNull(_pdFactory) onlyRole(ADMIN_ROLE) {
        require(isPublicDelegationEnabled, "Public delegation disabled.");
        require(_pdArgs.length > 0, "Invalid args.");

        IPublicDelegation.PDConstructorArgs memory _decodedPDArgs = abi.decode(
            _pdArgs,
            (IPublicDelegation.PDConstructorArgs)
        );
        publicDelegation = address(IPublicDelegationFactory(_pdFactory).deployPublicDelegation(_decodedPDArgs));
        rewardAddress = publicDelegation;

        emit SetPublicDelegation(_msgSender(), publicDelegation, rewardAddress);
    }

    /// @dev Agree on the initial lockup conditions.
    /// Emits a ReviewInitialConditions event.
    /// Emits a CompleteReviewInitialConditions.
    function reviewInitialConditions() external override beforeInit onlyRole(ADMIN_ROLE) {
        address _caller = _msgSender();
        require(lockupConditions.reviewedAdmin[_caller] == false, "Msg.sender already reviewed.");
        lockupConditions.reviewedAdmin[_caller] = true;
        unchecked {
            lockupConditions.reviewedCount++;
        }
        emit ReviewInitialConditions(_caller);

        if (lockupConditions.reviewedCount == getRoleMemberCount(ADMIN_ROLE)) {
            lockupConditions.allReviewed = true;
            emit CompleteReviewInitialConditions();
        }
    }

    /// @dev Completes the contract initialization by depositing initial lockup amounts.
    /// The transaction must send exactly the initial lockup amount of KAIA.
    /// Emits a DepositLockupStakingAndInit event.
    function depositLockupStakingAndInit() public payable virtual override beforeInit {
        require(
            gcId != 0 &&
                stakingTracker != address(0) &&
                (!isPublicDelegationEnabled ||
                    (isPublicDelegationEnabled && publicDelegation != address(0) && rewardAddress != address(0))),
            "Not set up properly."
        );
        require(lockupConditions.allReviewed == true, "Reviews not finished.");

        uint256 requiredStakingAmount;
        for (uint256 i = 0; i < lockupConditions.unlockAmount.length; i++) {
            unchecked {
                requiredStakingAmount += lockupConditions.unlockAmount[i];
            }
        }
        require(msg.value == requiredStakingAmount, "Value does not match.");
        initialLockupStaking = requiredStakingAmount;
        remainingLockupStaking = requiredStakingAmount;

        _grantStakingRoles();

        isInitialized = true;
        emit DepositLockupStakingAndInit(_msgSender(), msg.value);
    }

    /// @dev Grants staking related roles to the appropriate addresses.
    /// It will be overridden in the CnStakingV3MultiSig contract to handle multisig admins.
    function _grantStakingRoles() internal virtual {
        if (!isPublicDelegationEnabled) {
            address _owner = getRoleMember(OPERATOR_ROLE, 0);
            _grantRole(UNSTAKING_APPROVER_ROLE, _owner);
            _grantRole(UNSTAKING_CLAIMER_ROLE, _owner);
        } else {
            _grantRole(STAKER_ROLE, publicDelegation);
            _grantRole(UNSTAKING_APPROVER_ROLE, publicDelegation);
            _grantRole(UNSTAKING_CLAIMER_ROLE, publicDelegation);
        }
    }

    /* ========== OPERATOR FUNCTIONS ========== */
    /// Note that all OPERATOR_ROLE will be multi-sig controlled by the `CnStakingV3MultiSig` contract
    /// by setting the `owner` to the `address(this)`.

    /// @dev Withdraw a part of initial lockup stakes
    /// Emits a WithdrawLockupStaking event.
    function withdrawLockupStaking(
        address payable _to,
        uint256 _value
    ) external override onlyRole(OPERATOR_ROLE) notNull(_to) {
        (, , , , uint256 withdrawableAmount) = getLockupStakingInfo();
        require(_value > 0 && _value <= withdrawableAmount, "Value is not withdrawable.");

        unchecked {
            remainingLockupStaking -= _value;
        }

        (bool success, ) = _to.call{value: _value}("");
        require(success, "Transfer failed.");

        _refreshStake();
        emit WithdrawLockupStaking(_to, _value);
    }

    /// @dev Update the reward address in the AddressBook
    /// Emits an UpdateRewardAddress event.
    /// Need to call acceptRewardAddress() to reflect the change to AddressBook.
    /// The address can be null, which cancels the reward address update attempt.
    function updateRewardAddress(address _addr) external override onlyRole(OPERATOR_ROLE) {
        require(!isPublicDelegationEnabled, "Public delegation enabled.");
        pendingRewardAddress = _addr;
        emit UpdateRewardAddress(_addr);
    }

    /// @dev Finish updating the reward address
    /// Must be called from either the pendingRewardAddress, or one of the AddressBook admins.
    /// This step guarantees that the rewardAddress is owned by the current CN.
    ///
    /// Emits an AcceptRewardAddress event.
    /// Also emits a ReviseRewardAddress event from the AddressBook.
    function acceptRewardAddress(address _addr) external override {
        require(_canAcceptRewardAddress(), "Unauthorized to accept reward address.");
        require(_addr == pendingRewardAddress, "Given address does not match the pending.");

        IAddressBook(ADDRESS_BOOK_ADDRESS).reviseRewardAddress(pendingRewardAddress);
        rewardAddress = pendingRewardAddress;
        pendingRewardAddress = address(0);

        emit AcceptRewardAddress(rewardAddress);
    }

    /// @dev Update the staking tracker
    /// Emits an UpdateStakingTracker event.
    /// Should not be called if there is an active proposal
    function updateStakingTracker(address _tracker) external override onlyRole(OPERATOR_ROLE) {
        require(_tracker._validStakingTracker(1), "Invalid StakingTracker.");
        if (stakingTracker != address(0)) {
            IStakingTracker(stakingTracker).refreshStake(address(this));
            require(
                IStakingTracker(stakingTracker).getLiveTrackerIds().length == 0,
                "Cannot update tracker when there is an active tracker."
            );
        }

        stakingTracker = _tracker;
        emit UpdateStakingTracker(_tracker);
    }

    /// @dev Update the voter address of this CN
    /// Emits an UpdateVoterAddress event.
    function updateVoterAddress(address _addr) external override onlyRole(OPERATOR_ROLE) {
        voterAddress = _addr;
        if (stakingTracker != address(0)) {
            if (_addr != address(0)) {
                require(IStakingTracker(stakingTracker).voterToGCId(_addr) == 0, "Voter already taken.");
            }
            IStakingTracker(stakingTracker).refreshVoter(address(this));
        }
        emit UpdateVoterAddress(_addr);
    }

    /// @dev Toggle redelegation flag
    function toggleRedelegation() external override onlyRole(OPERATOR_ROLE) {
        require(isPublicDelegationEnabled, "Public delegation disabled.");
        isRedelegationEnabled = !isRedelegationEnabled;
        emit ToggleRedelegation(isRedelegationEnabled);
    }

    /* ========== PUBLIC DELEGATION FUNCTIONS ========== */

    /// @dev Delegate KAIA to this contract
    ///
    /// Emits a DelegateKaia event.
    function delegate() external payable override afterInit onlyRole(STAKER_ROLE) {
        _delegateKaia();
    }

    /// @dev The fallback which add more delegated stakes
    receive() external payable override afterInit onlyRole(STAKER_ROLE) {
        _delegateKaia();
    }

    /// @dev Redelegate stakes to another CnStakingV3 contract
    ///
    /// To prevent redelegation hopping, the last redelegation time is recorded.
    /// @param _user        The user of the redelegation (managed by the public delegation contract)
    /// @param _targetCnV3  The target CnStakingV3 contract
    /// @param _value       The amount of KAIA to redelegate
    ///
    /// Emits a Redelegate event.
    function redelegate(address _user, address _targetCnV3, uint256 _value) external override afterInit notNull(_user) {
        require(isRedelegationEnabled && _msgSender() == publicDelegation, "Redelegation disabled.");
        require(_targetCnV3 != address(this), "Target can't be self.");
        require(_targetCnV3._validCnStaking(3), "Invalid CnStakingV3.");
        require(_value > 0 && unstaking + _value <= staking, "Invalid value.");
        require(
            lastRedelegation[_user] == 0 || lastRedelegation[_user] + STAKE_LOCKUP <= block.timestamp,
            "Can't redelegate yet."
        );

        unchecked {
            staking -= _value;
        }

        // Directly transfer the _value to the target contract
        ICnStakingV3(payable(_targetCnV3)).handleRedelegation{value: _value}(_user);

        _refreshStake();

        emit Redelegation(_user, _targetCnV3, _value);
    }

    /// @dev Handle re-delegation from another CnStakingV3 contract
    ///
    /// It stakes KAIA on behalf of the `_user` using `stakeFor` function of the public delegation contract.
    /// Since public delegation will stake the received KAIA back to this contract, the balance should't be changed after the call.
    ///
    /// Emits a HandleRedelegation event.
    function handleRedelegation(address _user) external payable override afterInit notNull(_user) {
        // Early exit if the Redelegation disabled.
        // Note that isRedelegationEnabled can be true only if the public delegation is enabled.
        require(isRedelegationEnabled, "Redelegation disabled.");
        require(_msgSender()._validCnStaking(3), "Invalid CnStakingV3.");

        lastRedelegation[_user] = block.timestamp;

        IKIP163 _publicDelegation = IKIP163(payable(publicDelegation));

        uint256 expected;
        unchecked {
            expected = address(this).balance + _publicDelegation.reward();
        }
        /// @dev This will process new delegated stakes for `_user` in the public delegation contract.
        ///
        /// `stakeFor` function will receive KAIA and process the internal logic, and stake it back to this contract.
        /// If the balance changed after the call, it means the public delegation contract is malfunctioning.
        _publicDelegation.stakeFor{value: msg.value}(_user);
        require(expected == address(this).balance, "Invalid stakeFor.");

        emit HandleRedelegation(_user, _msgSender(), address(this), msg.value);
    }

    /// @dev Withdraw a part of delegated stakes.
    /// Emits a ApproveStakingWithdrawal event.
    function approveStakingWithdrawal(
        address _to,
        uint256 _value
    ) external override notNull(_to) onlyRole(UNSTAKING_APPROVER_ROLE) returns (uint256 id) {
        require(_value > 0 && unstaking + _value <= staking, "Invalid value.");
        id = withdrawalRequestCount;
        uint256 time;
        unchecked {
            withdrawalRequestCount++;

            time = block.timestamp + STAKE_LOCKUP;
            withdrawalRequestMap[id] = WithdrawalRequest({
                to: _to,
                value: _value,
                withdrawableFrom: time,
                state: WithdrawalStakingState.Unknown
            });

            unstaking += _value;
        }
        _refreshStake();
        emit ApproveStakingWithdrawal(id, _to, _value, time);
    }

    /// @dev cancel a withdrawal request
    /// Emits a CancelApprovedStakingWithdrawal event.
    function cancelApprovedStakingWithdrawal(uint256 _id) external override onlyRole(UNSTAKING_APPROVER_ROLE) {
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");

        request.state = WithdrawalStakingState.Canceled;
        unchecked {
            unstaking -= request.value;
        }
        _refreshStake();
        emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
    }

    /// @dev Take out an approved withdrawal amounts.
    ///
    /// If STAKE_LOCKUP has passed since WithdrawalRequest was created,
    /// an admin can call this function to execute the withdrawal.
    ///
    /// If 2*STAKE_LOCKUP has passed since WithdrawalRequest was created,
    /// the withdrawal is canceled by calling this function.
    ///
    /// Either way, unstaking amount decreases.
    ///
    /// The withdrawal request ID can be obtained from ApproveStakingWithdrawal event
    /// or getApprovedStakingWithdrawalIds().
    function withdrawApprovedStaking(uint256 _id) external override onlyRole(UNSTAKING_CLAIMER_ROLE) {
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");
        require(request.value <= staking, "Value is not withdrawable.");
        require(request.withdrawableFrom <= block.timestamp, "Not withdrawable yet.");

        uint256 withdrawableUntil;
        unchecked {
            withdrawableUntil = request.withdrawableFrom + STAKE_LOCKUP;
        }
        if (withdrawableUntil <= block.timestamp) {
            request.state = WithdrawalStakingState.Canceled;
            unchecked {
                unstaking -= request.value;
            }

            _refreshStake();
            emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
        } else {
            request.state = WithdrawalStakingState.Transferred;
            unchecked {
                staking -= request.value;
                unstaking -= request.value;
            }

            (bool success, ) = request.to.call{value: request.value}("");
            require(success, "Transfer failed.");

            _refreshStake();
            emit WithdrawApprovedStaking(_id, request.to, request.value);
        }
    }

    /* ========== PRIVATE FUNCTIONS ========== */

    /// @dev Add more delegated stakes
    /// Emits a DelegateKaia event.
    function _delegateKaia() private {
        require(msg.value > 0, "Invalid amount.");

        unchecked {
            staking += msg.value;
        }

        _refreshStake();
        emit DelegateKaia(_msgSender(), msg.value);
    }

    /// @dev Validate initial conditions
    /// 1. If the reward address is null, the initial lockup is disabled.
    /// 2. Initial lockup conditions must be valid.
    function _validInitialConditions(
        address _ra,
        uint256[] memory _unlockTime,
        uint256[] memory _unlockAmount
    ) private view {
        if (_ra == address(0)) {
            require(_unlockTime.length == 0 && _unlockAmount.length == 0, "Initial lockup disabled.");
        } else {
            require(_unlockTime.length == _unlockAmount.length, "Invalid initial conditions.");
            if (_unlockTime.length > 0) {
                uint256 unlockTime = block.timestamp;
                for (uint256 i = 0; i < _unlockAmount.length; i++) {
                    require(unlockTime < _unlockTime[i], "Unlock time is not in ascending order.");
                    require(_unlockAmount[i] > 0, "Amount is not positive number.");
                    unlockTime = _unlockTime[i];
                }
            }
        }
    }

    /// @dev Refresh the balance of this contract recorded in StakingTracker
    function _refreshStake() private {
        (bool success, ) = address(stakingTracker).call(
            abi.encodeWithSignature("refreshStake(address)", address(this))
        );
        require(success, "StakingTracker call failed.");
    }

    /// @dev Check if the caller can accept the reward address
    function _canAcceptRewardAddress() private view returns (bool) {
        if (_msgSender() == pendingRewardAddress) {
            return true;
        }
        (address[] memory abookAdminList, ) = IAddressBook(ADDRESS_BOOK_ADDRESS).getState();
        for (uint256 i = 0; i < abookAdminList.length; i++) {
            if (_msgSender() == abookAdminList[i]) {
                return true;
            }
        }
        return false;
    }

    function _checkNull(address _address) private pure {
        require(_address != address(0), "Address is null.");
    }

    function _checkInit() private view {
        require(isInitialized, "Contract is not initialized.");
    }

    function _checkNotInit() private view {
        require(!isInitialized, "Contract has been initialized.");
    }

    /* ========== PUBLIC GETTERS ========== */

    /// @dev Returns `true` if `account` has been granted `role`.
    /// If public delegation is disabled, the STAKER_ROLE is granted to everyone.
    function hasRole(bytes32 role, address account) public view override(IAccessControl, AccessControl) returns (bool) {
        if (!isPublicDelegationEnabled && role == STAKER_ROLE) {
            return true;
        }
        return super.hasRole(role, account);
    }

    /// @dev Query initial lockup status
    /// @return unlockTime    List of unlocking times in timestamp
    /// @return unlockAmount  List of unlocking amounts
    /// @return initial       Initial lockup amount
    /// @return remaining     Remaining lockup amount = (initial - withdrawn)
    /// @return withdrawable  Max withdrawable amount = (unlocked - withdrawn)
    function getLockupStakingInfo()
        public
        view
        override
        afterInit
        returns (
            uint256[] memory unlockTime,
            uint256[] memory unlockAmount,
            uint256 initial,
            uint256 remaining,
            uint256 withdrawable
        )
    {
        uint256 unlockedAmount = 0;

        for (uint256 i = 0; i < lockupConditions.unlockTime.length; i++) {
            if (block.timestamp > lockupConditions.unlockTime[i]) {
                unchecked {
                    unlockedAmount += lockupConditions.unlockAmount[i];
                }
            }

            // withdrawable = unlocked - withdrawn
            unchecked {
                withdrawable = unlockedAmount - (initialLockupStaking - remainingLockupStaking);
            }
        }

        return (
            lockupConditions.unlockTime,
            lockupConditions.unlockAmount,
            initialLockupStaking,
            remainingLockupStaking,
            withdrawable
        );
    }

    /// @dev Query withdrawal IDs that matches given state.
    ///
    /// For efficiency, only IDs in range (_from <= id < _to) are searched.
    /// If _to == 0 or _to >= requestCount, then the search range is (_from <= id < requestCount).
    ///
    /// @param _from   search begin index
    /// @param _to     search end index; but search till the end if _to == 0 or _to >= requestCount.
    /// @param _state  withdrawal state
    /// @return ids    withdrawal IDs satisfying the conditions
    function getApprovedStakingWithdrawalIds(
        uint256 _from,
        uint256 _to,
        WithdrawalStakingState _state
    ) external view override returns (uint256[] memory ids) {
        uint256 end = (_to == 0 || _to >= withdrawalRequestCount) ? withdrawalRequestCount : _to;

        ids = new uint256[](end - _from);
        uint256 cnt = 0;

        for (uint256 i = _from; i < end; i++) {
            if (withdrawalRequestMap[i].state == _state) {
                unchecked {
                    ids[cnt++] = i;
                }
            }
        }

        assembly {
            mstore(ids, cnt)
        }
    }

    /// @dev Query a withdrawal request details
    /// @param _index  withdrawal request ID
    /// @return to                recipient
    /// @return value             withdrawing amount
    /// @return withdrawableFrom  withdrawable timestamp
    /// @return state             the request state
    function getApprovedStakingWithdrawalInfo(
        uint256 _index
    )
        external
        view
        override
        returns (address to, uint256 value, uint256 withdrawableFrom, WithdrawalStakingState state)
    {
        return (
            withdrawalRequestMap[_index].to,
            withdrawalRequestMap[_index].value,
            withdrawalRequestMap[_index].withdrawableFrom,
            withdrawalRequestMap[_index].state
        );
    }
}
