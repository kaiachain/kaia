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

import "./ICnStakingV2.sol";
import "../IAddressBook.sol";

// Features
// 1. Administration
//   - Manage multisig admins.
//   - Every multisig operations must be approved (either submit or confirm)
//     by exactly `requirement` number of admins.
//   - Additionally a `contractValidator` takes part in contract initialization.
//   - Functions
//     - multisig AddAdmin: add an admin address
//     - multisig DeleteAdmin: delete an admin address
//     - multisig UpdateRequirement: change multisig threshold
//     - multisig ClearRequest: cancel all pending (NotConfirmed) multisig requests
//
// 2. Lockup stakes (Initial lockup)
//   - Initial lockup is a set of long-term fixed lockups.
//   - Every admins and the contractValidator must agree to the conditions
//     for this contract to initialize.
//   - KLAYs must be deposited for this contract to initialize.
//   - Functions
//     - reviewInitialConditions(): Agree to the initlal lockup conditions
//     - depositLockupStakingAndInit(): Deposit requried amount
//     - multisig WithdrawLockupStaking: Withdraw unlocked amount
//
// 3. Non-lockup stakes (Free stake)
//   - Free stakes can be added or removed at any time on admins' discretion.
//   - Free stakes can be added either by calling stakeKlay() or sending
//     a transaction to this contract with nonzero KLAY (via fallback).
//   - It takes STAKE_LOCKUP after withdrawal request to actually take out the KLAY.
//   - Functions
//     - multisig ApproveStakingWithdrawal: Schedule a withdrawal
//     - multisig CancelApprovedStakingWithdrawal: Cancel a withdrawal request
//     - withdrawApprovedStaking(): Take out the KLAY or cancel an expired withdrawal request.
//
// 4. External accounts
//   - Several addresses constitute the identity of this CN.
//   - Among them, RewardAddress can be modified via CnStaking contract.
//   - Functions
//     - multisig UpdateRewardAddress: Setup pendingRewardAddress
//     - acceptRewardAddress(): Request AddressBook to change reward address.
//     - multisig UpdateStakingTracker: Change the StakingTracker contract to report stakes.
//     - multisig UpdateVoterAddress: Change the Voter account and notify to StakingTracker.

// Code organization
// - Constants
// - States
// - Modifiers
// - Mutators
//   - Constructor and initializers
//   - Specific multisig operations
//   - Generic multisig facility
//   - Private helpers
//   - Other public functions
// - Getters

contract CnStakingV2 is ICnStakingV2 {
    // Constants
    // - Constants are defined as virtual functions to allow easier unit tests.
    uint256 public constant ONE_WEEK = 1 weeks;

    function MAX_ADMIN() public view virtual override returns (uint256) {
        return 50;
    }

    function CONTRACT_TYPE() public view virtual override returns (string memory) {
        return "CnStakingContract";
    }

    function VERSION() public view virtual override returns (uint256) {
        return 2;
    }

    function ADDRESS_BOOK_ADDRESS() public view virtual override returns (address) {
        return 0x0000000000000000000000000000000000000400;
    }

    function STAKE_LOCKUP() public view virtual override returns (uint256) {
        return ONE_WEEK;
    }

    // State variables

    // Multisig admin list
    address public contractValidator; // temporary admin only used during initialization
    address[] public adminList; // all persistent admins
    uint256 public requirement; // this number of admins must approve a request
    mapping(address => bool) public isAdmin;

    // Multisig requests
    uint256 public lastClearedId; // For efficient ClearRequest
    uint256 public requestCount;
    mapping(uint256 => Request) private requestMap;
    struct Request {
        Functions functionId;
        bytes32 firstArg;
        bytes32 secondArg;
        bytes32 thirdArg;
        address requestProposer;
        address[] confirmers;
        RequestState state;
    }

    // Initial lockup
    LockupConditions public lockupConditions;
    uint256 public initialLockupStaking;
    uint256 public remainingLockupStaking;
    bool public isInitialized;
    struct LockupConditions {
        uint256[] unlockTime;
        uint256[] unlockAmount;
        bool allReviewed;
        uint256 reviewedCount;
        mapping(address => bool) reviewedAdmin;
    }

    // Free stakes
    uint256 public staking;
    uint256 public unstaking;
    uint256 public withdrawalRequestCount;
    mapping(uint256 => WithdrawalRequest) private withdrawalRequestMap;
    struct WithdrawalRequest {
        address to;
        uint256 value;
        uint256 withdrawableFrom;
        WithdrawalStakingState state;
    }

    // External accounts
    uint256 public override gcId; // used to group staking contracts
    address public override nodeId; // informational
    address public override rewardAddress; // informational
    address public override pendingRewardAddress; // used in updateRewardAddress in progress
    address public override stakingTracker; // used to call refreshStake(), refreshVoter()
    address public override voterAddress; // read by StakingTracker

    modifier onlyMultisigTx() {
        require(msg.sender == address(this), "Not a multisig-transaction.");
        _;
    }

    modifier onlyAdmin(address _admin) {
        require(isAdmin[_admin], "Address is not admin.");
        _;
    }

    modifier adminDoesNotExist(address _admin) {
        require(!isAdmin[_admin], "Admin already exists.");
        _;
    }

    modifier notNull(address _address) {
        require(_address != address(0), "Address is null");
        _;
    }

    modifier notConfirmedRequest(uint256 _id) {
        require(requestMap[_id].state == RequestState.NotConfirmed, "Must be at not-confirmed state.");
        _;
    }

    modifier validRequirement(uint256 _adminCount, uint256 _requirement) {
        require(
            _adminCount <= MAX_ADMIN() && _requirement <= _adminCount && _requirement != 0 && _adminCount != 0,
            "Invalid requirement."
        );
        _;
    }

    modifier beforeInit() {
        require(isInitialized == false, "Contract has been initialized.");
        _;
    }

    modifier afterInit() {
        require(isInitialized == true, "Contract is not initialized.");
        _;
    }

    // Initialization functions

    /// @dev Fill in initial values for the contract
    /// Emits a DeployContract event.
    /// @param _contractValidator  A temporary admin to perform initial condition checks
    /// @param _nodeId             The NodeID of this CN
    /// @param _rewardAddress      The RewardBase of this CN
    /// @param _cnAdminlist        Initial list of admins
    /// @param _requirement        Number of required multisig confirmations
    /// @param _unlockTime         List of initial lockup deadlines in block timestamp
    /// @param _unlockAmount       List of initial lockup amounts in peb
    constructor(
        address _contractValidator,
        address _nodeId,
        address _rewardAddress,
        address[] memory _cnAdminlist,
        uint256 _requirement,
        uint256[] memory _unlockTime,
        uint256[] memory _unlockAmount
    )
        notNull(_contractValidator)
        notNull(_nodeId)
        notNull(_rewardAddress)
        validRequirement(_cnAdminlist.length, _requirement)
    {
        // Sanitize _cnAdminlist
        isAdmin[_contractValidator] = true;
        for (uint256 i = 0; i < _cnAdminlist.length; i++) {
            require(!isAdmin[_cnAdminlist[i]] && _cnAdminlist[i] != address(0), "Address is null or not unique.");
            isAdmin[_cnAdminlist[i]] = true;
        }

        // Sanitize _unlockTime and _unlockAmount
        require(
            _unlockTime.length != 0 && _unlockAmount.length != 0 && _unlockTime.length == _unlockAmount.length,
            "Invalid unlock time and amount."
        );
        uint256 unlockTime = block.timestamp;

        for (uint256 i = 0; i < _unlockAmount.length; i++) {
            require(unlockTime < _unlockTime[i], "Unlock time is not in ascending order.");
            require(_unlockAmount[i] > 0, "Amount is not positive number.");
            unlockTime = _unlockTime[i];
        }

        contractValidator = _contractValidator;
        nodeId = _nodeId;
        rewardAddress = _rewardAddress;

        adminList = _cnAdminlist;
        requirement = _requirement;

        lockupConditions.unlockTime = _unlockTime;
        lockupConditions.unlockAmount = _unlockAmount;
        isInitialized = false;

        emit DeployContract(
            CONTRACT_TYPE(),
            _contractValidator,
            _nodeId,
            _rewardAddress,
            _cnAdminlist,
            _requirement,
            _unlockTime,
            _unlockAmount
        );
    }

    /// @dev Set the initial stakingTracker address
    /// Emits a UpdateStakingTracker event.
    /// This step can be skipped if automatic StakingTracker refresh is not needed.
    function setStakingTracker(address _tracker) external override beforeInit onlyAdmin(msg.sender) notNull(_tracker) {
        require(validStakingTracker(_tracker), "Invalid contract");

        stakingTracker = _tracker;
        emit UpdateStakingTracker(_tracker);
    }

    /// @dev Set the gcId
    /// The gcId never changes once initialized.
    /// Emits a UpdateCouncilId event.
    function setGCId(uint256 _gcId) external override beforeInit onlyAdmin(msg.sender) {
        require(_gcId != 0, "GC ID cannot be zero");
        gcId = _gcId;
        emit UpdateGCId(_gcId);
    }

    /// @dev Agree on the initial lockup conditions
    /// The contractValidator and every initial admins (cnAdminList) must agree
    /// for this contract to initialize.
    /// Emits a ReviewInitialConditions event.
    /// Emits a CompleteReviewInitialConditions if everyone has reviewed.
    function reviewInitialConditions() external override beforeInit onlyAdmin(msg.sender) {
        require(lockupConditions.reviewedAdmin[msg.sender] == false, "Msg.sender already reviewed.");
        lockupConditions.reviewedAdmin[msg.sender] = true;
        lockupConditions.reviewedCount++;
        emit ReviewInitialConditions(msg.sender);

        if (lockupConditions.reviewedCount == adminList.length + 1) {
            lockupConditions.allReviewed = true;
            emit CompleteReviewInitialConditions();
        }
    }

    /// @dev Completes the contract initialization by depositing initial lockup amounts.
    /// Everyone must have agreed on initial lockup conditions,
    /// The transaction must send exactly the initial lockup amount of KLAY.
    /// Emits a DepositLockupStakingAndInit event.
    function depositLockupStakingAndInit() external payable override beforeInit {
        require(gcId != 0, "GC ID cannot be zero");
        require(lockupConditions.allReviewed == true, "Reviewing is not finished.");

        uint256 requiredStakingAmount;
        for (uint256 i = 0; i < lockupConditions.unlockAmount.length; i++) {
            requiredStakingAmount += lockupConditions.unlockAmount[i];
        }
        require(msg.value == requiredStakingAmount, "Value does not match.");
        initialLockupStaking = requiredStakingAmount;
        remainingLockupStaking = requiredStakingAmount;

        // Remove the temporary admin (i.e. contractValidator)
        isAdmin[contractValidator] = false;
        delete contractValidator;

        isInitialized = true;
        emit DepositLockupStakingAndInit(msg.sender, msg.value);
    }

    // Multisig operations

    /// @dev Submit a request to add an admin to adminList
    /// @param _admin  new admin address
    function submitAddAdmin(
        address _admin
    )
        external
        override
        afterInit
        onlyAdmin(msg.sender)
        notNull(_admin)
        adminDoesNotExist(_admin)
        validRequirement(adminList.length + 1, requirement)
    {
        uint256 id = submitRequest(Functions.AddAdmin, toBytes32(_admin), 0, 0);
        confirmRequest(id);
    }

    /// @dev Add an admin to adminList
    /// @param _admin  new admin address
    /// Emits an AddAdmin event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function addAdmin(
        address _admin
    )
        external
        override
        onlyMultisigTx
        notNull(_admin)
        adminDoesNotExist(_admin)
        validRequirement(adminList.length + 1, requirement)
    {
        isAdmin[_admin] = true;
        adminList.push(_admin);
        clearRequest();
        emit AddAdmin(_admin);
    }

    /// @dev Submit a request to delete an admin from adminList
    /// @param _admin  the admin address
    function submitDeleteAdmin(
        address _admin
    )
        external
        override
        afterInit
        onlyAdmin(msg.sender)
        notNull(_admin)
        onlyAdmin(_admin)
        validRequirement(adminList.length - 1, requirement)
    {
        uint256 id = submitRequest(Functions.DeleteAdmin, toBytes32(_admin), 0, 0);
        confirmRequest(id);
    }

    /// @dev Delete an admin from adminList
    /// @param _admin  the admin address
    /// Emits a DeleteAdmin event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function deleteAdmin(
        address _admin
    )
        external
        override
        onlyMultisigTx
        notNull(_admin)
        onlyAdmin(_admin)
        validRequirement(adminList.length - 1, requirement)
    {
        deleteArrayElement(adminList, _admin);
        isAdmin[_admin] = false;
        clearRequest();
        emit DeleteAdmin(_admin);
    }

    /// @dev submit a request to update the confirmation threshold
    /// @param _requirement  new confirmation threshold
    function submitUpdateRequirement(
        uint256 _requirement
    ) external override afterInit onlyAdmin(msg.sender) validRequirement(adminList.length, _requirement) {
        require(_requirement != requirement, "Invalid value");
        uint256 id = submitRequest(Functions.UpdateRequirement, bytes32(_requirement), 0, 0);
        confirmRequest(id);
    }

    /// @dev update the confirmation threshold
    /// @param _requirement  new confirmation threshold
    /// Emits an UpdateRequirement event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function updateRequirement(
        uint256 _requirement
    ) external override onlyMultisigTx validRequirement(adminList.length, _requirement) {
        requirement = _requirement;
        clearRequest();
        emit UpdateRequirement(_requirement);
    }

    /// @dev submit a request to cancel all outstanding (i.e. NotConfirmed) requests
    function submitClearRequest() external override afterInit onlyAdmin(msg.sender) {
        uint256 id = submitRequest(Functions.ClearRequest, 0, 0, 0);
        confirmRequest(id);
    }

    /// @dev cancel all outstanding (i.e. NotConfirmed) requests
    /// Emits a ClearRequest event.
    function clearRequest() public override onlyMultisigTx {
        for (uint256 i = lastClearedId; i < requestCount; i++) {
            if (requestMap[i].state == RequestState.NotConfirmed) {
                requestMap[i].state = RequestState.Canceled;
            }
        }
        lastClearedId = requestCount;
        emit ClearRequest();
    }

    /// @dev Submit a request to withdraw a part of initial lockup stakes
    ///
    /// Max withdrawable amount is (unlocked - withdrawn),
    /// where unlocked = amounts that lockup period has passed,
    /// and withdrawn = (initial - remaining).
    ///
    /// @param _to     The recipient address
    /// @param _value  The amount
    function submitWithdrawLockupStaking(
        address payable _to,
        uint256 _value
    ) external override afterInit onlyAdmin(msg.sender) notNull(_to) {
        (, , , , uint256 withdrawableAmount) = getLockupStakingInfo();
        require(_value > 0 && _value <= withdrawableAmount, "Invalid value.");

        uint256 id = submitRequest(Functions.WithdrawLockupStaking, toBytes32(_to), bytes32(_value), 0);
        confirmRequest(id);
    }

    /// @dev Withdraw a part of initial lockup stakes
    /// Emits a WithdrawLockupStaking event.
    function withdrawLockupStaking(address payable _to, uint256 _value) external override onlyMultisigTx notNull(_to) {
        (, , , , uint256 withdrawableAmount) = getLockupStakingInfo();
        require(_value > 0 && _value <= withdrawableAmount, "Value is not withdrawable.");

        remainingLockupStaking -= _value;

        (bool success, ) = _to.call{value: _value}("");
        require(success, "Transfer failed.");

        safeRefreshStake();
        emit WithdrawLockupStaking(_to, _value);
    }

    /// @dev submit a request to withdraw a part of free stakes.
    ///
    /// Creates a new WithdrawalRequest
    /// The WithdrawalRequest is withdrawable from request creation + STAKE_LOCKUP.
    /// The WithdrawalRequest expires from request creation + 2 * STAKE_LOCKUP.
    ///
    /// Max withdrawable amount is (staked - unstaking).
    /// Once the WithdrawalRequest is created, unstaking amount increases.
    ///
    /// @param _to     The recipient address
    /// @param _value  The amount
    function submitApproveStakingWithdrawal(
        address _to,
        uint256 _value
    ) external override afterInit onlyAdmin(msg.sender) notNull(_to) {
        require(_value > 0 && _value <= staking, "Invalid value.");
        require(unstaking + _value <= staking, "Too much outstanding withdrawal");
        uint256 id = submitRequest(Functions.ApproveStakingWithdrawal, toBytes32(_to), bytes32(_value), 0);
        confirmRequest(id);
    }

    /// @dev Withdraw a part of free stakes.
    /// Emits a ApproveStakingWithdrawal event.
    function approveStakingWithdrawal(address _to, uint256 _value) external override onlyMultisigTx notNull(_to) {
        require(_value > 0 && _value <= staking, "Invalid value.");
        require(unstaking + _value <= staking, "Too much outstanding withdrawal");
        uint256 id = withdrawalRequestCount;
        withdrawalRequestCount++;

        uint256 time = block.timestamp + STAKE_LOCKUP();
        withdrawalRequestMap[id] = WithdrawalRequest({
            to: _to,
            value: _value,
            withdrawableFrom: time,
            state: WithdrawalStakingState.Unknown
        });
        unstaking += _value;
        safeRefreshStake();
        emit ApproveStakingWithdrawal(id, _to, _value, time);
    }

    /// @dev submit a request to cancel a withdrawal request
    /// The withdrawal request ID can be obtained from ApproveStakingWithdrawal event
    /// or getApprovedStakingWithdrawalIds().
    /// Unstaking amount decreases.
    function submitCancelApprovedStakingWithdrawal(uint256 _id) external override afterInit onlyAdmin(msg.sender) {
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");

        uint256 id = submitRequest(Functions.CancelApprovedStakingWithdrawal, bytes32(_id), 0, 0);
        confirmRequest(id);
    }

    /// @dev cancel a withdrawal request
    /// Emits a CancelApprovedStakingWithdrawal event.
    function cancelApprovedStakingWithdrawal(uint256 _id) external override onlyMultisigTx {
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");

        request.state = WithdrawalStakingState.Canceled;
        unstaking -= request.value;
        safeRefreshStake();
        emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
    }

    /// @dev submit a request to update the reward address of this CN
    function submitUpdateRewardAddress(address _addr) external override afterInit onlyAdmin(msg.sender) {
        uint256 id = submitRequest(Functions.UpdateRewardAddress, toBytes32(_addr), 0, 0);
        confirmRequest(id);
    }

    /// @dev Update the reward address in the AddressBook
    /// Emits an UpdateRewardAddress event.
    /// Need to call acceptRewardAddress() to reflect the change to AddressBook.
    /// The address can be null, which cancels the reward address update attempt.
    function updateRewardAddress(address _addr) external override onlyMultisigTx {
        pendingRewardAddress = _addr;
        emit UpdateRewardAddress(_addr);
    }

    /// @dev submit a request to update the staking tracker this CN reports to
    /// Should not be called if there is an active proposal
    function submitUpdateStakingTracker(
        address _tracker
    ) external override afterInit onlyAdmin(msg.sender) notNull(_tracker) {
        require(validStakingTracker(_tracker), "Invalid contract");
        if (stakingTracker != address(0)) {
            IStakingTracker(stakingTracker).refreshStake(address(this));
            require(
                IStakingTracker(stakingTracker).getLiveTrackerIds().length == 0,
                "Cannot update tracker when there is an active tracker"
            );
        }

        uint256 id = submitRequest(Functions.UpdateStakingTracker, toBytes32(_tracker), 0, 0);
        confirmRequest(id);
    }

    /// @dev Update the staking tracker
    /// Emits an UpdateStakingTracker event.
    /// Should not be called if there is an active proposal
    function updateStakingTracker(address _tracker) external override onlyMultisigTx notNull(_tracker) {
        require(validStakingTracker(_tracker), "Invalid contract");
        if (stakingTracker != address(0)) {
            IStakingTracker(stakingTracker).refreshStake(address(this));
            require(
                IStakingTracker(stakingTracker).getLiveTrackerIds().length == 0,
                "Cannot update tracker when there is an active tracker"
            );
        }

        stakingTracker = _tracker;
        emit UpdateStakingTracker(_tracker);
    }

    /// @dev submit a request to update the voter address of this CN
    function submitUpdateVoterAddress(address _addr) external override afterInit onlyAdmin(msg.sender) {
        if (stakingTracker != address(0) && _addr != address(0)) {
            address oldGCId = IStakingTracker(stakingTracker).voterToGCId(_addr);
            require(oldGCId == address(0), "Voter address already taken");
        }
        uint256 id = submitRequest(Functions.UpdateVoterAddress, toBytes32(_addr), 0, 0);
        confirmRequest(id);
    }

    /// @dev Update the voter address of this CN
    /// Emits an UpdateVoterAddress event.
    function updateVoterAddress(address _addr) external override onlyMultisigTx {
        voterAddress = _addr;

        if (stakingTracker != address(0)) {
            IStakingTracker(stakingTracker).refreshVoter(address(this));
        }
        emit UpdateVoterAddress(_addr);
    }

    // Generic multisig facility

    /// @dev Submits a request
    /// Emits a SubmitRequest event.
    /// @return  the request ID
    function submitRequest(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) private returns (uint256) {
        uint256 id = requestCount;
        requestCount++;

        requestMap[id] = Request({
            functionId: _functionId,
            firstArg: _firstArg,
            secondArg: _secondArg,
            thirdArg: _thirdArg,
            requestProposer: msg.sender,
            confirmers: new address[](0),
            state: RequestState.NotConfirmed
        });
        emit SubmitRequest(id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg);
        return id;
    }

    /// @dev Confirm a submitted request by another admin
    /// Note that a submitXYZ() automatically calls confirmRequest().
    /// Therefore an explicit confirmRequest() is only relevant when requirement >= 2.
    ///
    /// Emits a ConfirmRequest event.
    /// The necessary data can be obtained from SubmitRequest event or getRequestInfo().
    ///
    /// @param _id          The request ID
    /// @param _functionId  The function ID in enum Functions
    /// @param _firstArg    The first argument
    /// @param _secondArg   The second argument
    /// @param _thirdArg    The third argument
    function confirmRequest(
        uint256 _id,
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) public override notConfirmedRequest(_id) onlyAdmin(msg.sender) {
        require(!hasConfirmed(_id, msg.sender), "Msg.sender already confirmed.");
        require(
            requestMap[_id].functionId == _functionId &&
                requestMap[_id].firstArg == _firstArg &&
                requestMap[_id].secondArg == _secondArg &&
                requestMap[_id].thirdArg == _thirdArg,
            "Function id and arguments do not match."
        );

        requestMap[_id].confirmers.push(msg.sender);
        emit ConfirmRequest(_id, msg.sender, _functionId, _firstArg, _secondArg, _thirdArg, requestMap[_id].confirmers);

        if (requestMap[_id].confirmers.length >= requirement) {
            executeRequest(_id);
        }
    }

    /// @dev Shortcut of confirmRequest(...)
    /// Used by submitXYZ() functions.
    function confirmRequest(uint256 id) private {
        confirmRequest(
            id,
            requestMap[id].functionId,
            requestMap[id].firstArg,
            requestMap[id].secondArg,
            requestMap[id].thirdArg
        );
    }

    /// @dev Revoke a confirmation to a request
    /// If the sender is the proposer of the request, the request is canceled.
    /// Otherwise, the sender is simply deleted from the confirmers list.
    ///
    /// Emits a CancelRequest or RevokeConfirmation event.
    /// The necessary data can be obtained from SubmitRequest event or getRequestInfo().
    ///
    /// @param _id          The request ID
    /// @param _functionId  The function ID in enum Functions
    /// @param _firstArg    The first argument
    /// @param _secondArg   The second argument
    /// @param _thirdArg    The third argument
    function revokeConfirmation(
        uint256 _id,
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) external override notConfirmedRequest(_id) onlyAdmin(msg.sender) {
        require(hasConfirmed(_id, msg.sender), "Msg.sender has not confirmed.");
        require(
            requestMap[_id].functionId == _functionId &&
                requestMap[_id].firstArg == _firstArg &&
                requestMap[_id].secondArg == _secondArg &&
                requestMap[_id].thirdArg == _thirdArg,
            "Function id and arguments do not match."
        );

        if (requestMap[_id].requestProposer == msg.sender) {
            requestMap[_id].state = RequestState.Canceled;
            emit CancelRequest(
                _id,
                msg.sender,
                requestMap[_id].functionId,
                requestMap[_id].firstArg,
                requestMap[_id].secondArg,
                requestMap[_id].thirdArg
            );
        } else {
            deleteArrayElement(requestMap[_id].confirmers, msg.sender);
            emit RevokeConfirmation(
                _id,
                msg.sender,
                requestMap[_id].functionId,
                requestMap[_id].firstArg,
                requestMap[_id].secondArg,
                requestMap[_id].thirdArg,
                requestMap[_id].confirmers
            );
        }
    }

    /// @dev execute a requested function
    /// Used by confirmRequest when enough confirmations are made.
    /// Emits a ExecuteRequestSuccess or ExecuteRequestFailure event.
    function executeRequest(uint256 _id) private {
        bool ok = false;
        bytes memory out;
        Functions funcId = requestMap[_id].functionId;
        bytes32 a1 = requestMap[_id].firstArg;
        bytes32 a2 = requestMap[_id].secondArg;
        bytes32 a3 = requestMap[_id].thirdArg;

        if (funcId == Functions.AddAdmin) {
            (ok, out) = address(this).call(abi.encodeWithSignature("addAdmin(address)", a1));
        } else if (funcId == Functions.DeleteAdmin) {
            (ok, out) = address(this).call(abi.encodeWithSignature("deleteAdmin(address)", a1));
        } else if (funcId == Functions.UpdateRequirement) {
            (ok, out) = address(this).call(abi.encodeWithSignature("updateRequirement(uint256)", a1));
        } else if (funcId == Functions.ClearRequest) {
            (ok, out) = address(this).call(abi.encodeWithSignature("clearRequest()"));
        } else if (funcId == Functions.WithdrawLockupStaking) {
            (ok, out) = address(this).call(abi.encodeWithSignature("withdrawLockupStaking(address,uint256)", a1, a2));
        } else if (funcId == Functions.ApproveStakingWithdrawal) {
            (ok, out) = address(this).call(
                abi.encodeWithSignature("approveStakingWithdrawal(address,uint256)", a1, a2)
            );
        } else if (funcId == Functions.CancelApprovedStakingWithdrawal) {
            (ok, out) = address(this).call(abi.encodeWithSignature("cancelApprovedStakingWithdrawal(uint256)", a1));
        } else if (funcId == Functions.UpdateRewardAddress) {
            (ok, out) = address(this).call(abi.encodeWithSignature("updateRewardAddress(address)", a1));
        } else if (funcId == Functions.UpdateStakingTracker) {
            (ok, out) = address(this).call(abi.encodeWithSignature("updateStakingTracker(address)", a1));
        } else if (funcId == Functions.UpdateVoterAddress) {
            (ok, out) = address(this).call(abi.encodeWithSignature("updateVoterAddress(address)", a1));
        } else {
            revert("Unsupported function");
        }

        if (ok) {
            requestMap[_id].state = RequestState.Executed;
            emit ExecuteRequestSuccess(_id, msg.sender, funcId, a1, a2, a3);
        } else {
            requestMap[_id].state = RequestState.ExecutionFailed;
            emit ExecuteRequestFailure(_id, msg.sender, funcId, a1, a2, a3);
        }
    }

    // Helper functions

    function toBytes32(address _x) private pure returns (bytes32) {
        return bytes32(uint256(uint160(_x)));
    }

    function hasConfirmed(uint256 _id, address addr) private view returns (bool) {
        for (uint i = 0; i < requestMap[_id].confirmers.length; i++) {
            if (requestMap[_id].confirmers[i] == addr) {
                return true;
            }
        }
        return false;
    }

    function deleteArrayElement(address[] storage array, address target) private {
        for (uint i = 0; i < array.length; i++) {
            if (array[i] == target) {
                if (i != array.length - 1) {
                    array[i] = array[array.length - 1];
                }
                array.pop();
                return;
            }
        }
    }

    /// @dev Checks if a given address is valid StakingTracker contract
    function validStakingTracker(address _tracker) private view returns (bool) {
        string memory _type = IStakingTracker(_tracker).CONTRACT_TYPE();
        uint256 _version = IStakingTracker(_tracker).VERSION();
        return (keccak256(bytes(_type)) == keccak256(bytes("StakingTracker")) && _version == 1);
    }

    // Public functions

    /// @dev Add more free stakes
    /// Emits a StakeKlay event.
    function stakeKlay() public payable override afterInit {
        require(msg.value > 0, "Invalid amount.");
        staking += msg.value;
        safeRefreshStake();
        emit StakeKlay(msg.sender, msg.value);
    }

    /// @dev The fallback which add more free stakes
    ///
    /// Note that This fallback only accept transactions with empty calldata.
    /// contract calls with wrong function signature is reverted despite this fallback.
    receive() external payable override afterInit {
        stakeKlay();
    }

    /// @dev Refresh the balance of this contract recorded in StakingTracker
    /// This function should never revert to allow financial features to work
    /// even if stakingTracker is accidentally malfunctioning.
    function safeRefreshStake() private {
        stakingTracker.call(abi.encodeWithSignature("refreshStake(address)", address(this)));
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
    function withdrawApprovedStaking(uint256 _id) external override onlyAdmin(msg.sender) {
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");
        require(request.value <= staking, "Value is not withdrawable.");
        require(request.withdrawableFrom <= block.timestamp, "Not withdrawable yet.");

        uint256 withdrawableUntil = request.withdrawableFrom + STAKE_LOCKUP();
        if (withdrawableUntil <= block.timestamp) {
            request.state = WithdrawalStakingState.Canceled;
            unstaking -= request.value;

            safeRefreshStake();
            emit CancelApprovedStakingWithdrawal(_id, request.to, request.value);
        } else {
            request.state = WithdrawalStakingState.Transferred;
            staking -= request.value;
            unstaking -= request.value;

            (bool success, ) = request.to.call{value: request.value}("");
            require(success, "Transfer failed.");

            safeRefreshStake();
            emit WithdrawApprovedStaking(_id, request.to, request.value);
        }
    }

    /// @dev Finish updating the reward address
    /// Must be called from either the pendingRewardAddress, or one of the AddressBook admins.
    /// This step guarantees that the rewardAddress is owned by the current CN.
    ///
    /// Emits an AcceptRewardAddress event.
    /// Also emits a ReviseRewardAddress event from the AddressBook.
    function acceptRewardAddress(address _addr) external override {
        require(canAcceptRewardAddress(), "Unauthorized to accept reward address");
        require(_addr == pendingRewardAddress, "Given address does not match the pending");

        IAddressBook(ADDRESS_BOOK_ADDRESS()).reviseRewardAddress(pendingRewardAddress);
        rewardAddress = pendingRewardAddress;
        pendingRewardAddress = address(0);

        emit UpdateRewardAddress(rewardAddress);
    }

    function canAcceptRewardAddress() private returns (bool) {
        if (msg.sender == pendingRewardAddress) {
            return true;
        }
        (address[] memory abookAdminList, ) = IAddressBook(ADDRESS_BOOK_ADDRESS()).getState();
        for (uint256 i = 0; i < abookAdminList.length; i++) {
            if (msg.sender == abookAdminList[i]) {
                return true;
            }
        }
        return false;
    }

    // Public getters

    /// @dev Return the reviewers of the initial lockup conditions
    /// @return  reviewers addresses
    function getReviewers() external view override beforeInit returns (address[] memory) {
        address[] memory reviewers = new address[](lockupConditions.reviewedCount);
        uint256 id = 0;
        if (lockupConditions.reviewedAdmin[contractValidator] == true) {
            reviewers[id] = contractValidator;
            id++;
        }
        for (uint256 i = 0; i < adminList.length; i++) {
            if (lockupConditions.reviewedAdmin[adminList[i]] == true) {
                reviewers[id] = adminList[i];
                id++;
            }
        }
        return reviewers;
    }

    /// @dev Return the overall adminstrative states
    function getState()
        external
        view
        override
        returns (
            address _contractValidator,
            address _nodeId,
            address _rewardAddress,
            address[] memory _adminList,
            uint256 _requirement,
            uint256[] memory _unlockTime,
            uint256[] memory _unlockAmount,
            bool _allReviewed,
            bool _isInitialized
        )
    {
        return (
            contractValidator,
            nodeId,
            rewardAddress,
            adminList,
            requirement,
            lockupConditions.unlockTime,
            lockupConditions.unlockAmount,
            lockupConditions.allReviewed,
            isInitialized
        );
    }

    /// @dev Query request IDs that matches given state.
    ///
    /// For efficiency, only IDs in range (_from <= id < _to) are searched.
    /// If _to == 0 or _to >= requestCount, then the search range is (_from <= id < requestCount).
    ///
    /// @param _from   search begin index
    /// @param _to     search end index; but search till the end if _to == 0 or _to >= requestCount.
    /// @param _state  request state
    /// @return ids    request IDs satisfying the conditions
    function getRequestIds(
        uint256 _from,
        uint256 _to,
        RequestState _state
    ) external view override returns (uint256[] memory ids) {
        uint256 begin = _from;
        uint256 end = _to;
        if (_to == 0 || _to >= requestCount) {
            end = requestCount;
        }

        // Because memory array cannot grow, we must calculate size first.
        uint cnt = 0;
        for (uint i = begin; i < end; i++) {
            if (requestMap[i].state == _state) {
                cnt++;
            }
        }
        ids = new uint256[](cnt);
        cnt = 0;
        for (uint i = begin; i < end; i++) {
            if (requestMap[i].state == _state) {
                ids[cnt] = i;
                cnt++;
            }
        }
        return ids;
    }

    /// @dev Query a request details
    /// @param _id  requestID
    function getRequestInfo(
        uint256 _id
    )
        external
        view
        override
        returns (
            Functions functionId,
            bytes32 firstArg,
            bytes32 secondArg,
            bytes32 thirdArg,
            address proposer,
            address[] memory confirmers,
            RequestState state
        )
    {
        Request storage r = requestMap[_id];
        return (r.functionId, r.firstArg, r.secondArg, r.thirdArg, r.requestProposer, r.confirmers, r.state);
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
                unlockedAmount += lockupConditions.unlockAmount[i];
            }
        }

        uint256 withdrawnAmount = initialLockupStaking - remainingLockupStaking;
        uint256 withdrawableAmount = unlockedAmount - withdrawnAmount;

        return (
            lockupConditions.unlockTime,
            lockupConditions.unlockAmount,
            initialLockupStaking,
            remainingLockupStaking,
            withdrawableAmount
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
        uint256 begin = _from;
        uint256 end = _to;
        if (_to == 0 || _to >= withdrawalRequestCount) {
            end = withdrawalRequestCount;
        }

        // Because memory array cannot grow, we must calculate size first.
        uint cnt = 0;
        for (uint i = begin; i < end; i++) {
            if (withdrawalRequestMap[i].state == _state) {
                cnt += 1;
            }
        }
        ids = new uint256[](cnt);
        cnt = 0;
        for (uint i = begin; i < end; i++) {
            if (withdrawalRequestMap[i].state == _state) {
                ids[cnt] = i;
                cnt++;
            }
        }
        return ids;
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

interface IStakingTracker {
    function refreshStake(address staking) external;

    function refreshVoter(address voter) external;

    function CONTRACT_TYPE() external view returns (string memory);

    function VERSION() external view returns (uint256);

    function voterToGCId(address voter) external view returns (address nodeId);

    function getLiveTrackerIds() external view returns (uint256[] memory);
}
