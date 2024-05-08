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

import "./ICnStakingV3MultiSig.sol";
import "./CnStakingV3MultiSigStorage.sol";
import "./CnStakingV3.sol";

// Feature
// Provide all multisig operations for `CnStakingV3`, ported from `CnStakingV2`.
// The owner of `CnStakingV3` is set to `address(this)`, which means
// `onlyRole(OPERATOR_ROLE)` replaces `onlyMultisigTx` in `CnStakingV2`.
// Multisig operations
//   - Every multisig operations must be approved (either submit or confirm)
//     by exactly `requirement` number of admins.
//   - Additionally a `contractValidator` takes part in contract initialization.
//   - Functions
//     - multisig AddAdmin: add an admin address
//     - multisig DeleteAdmin: delete an admin address
//     - multisig UpdateRequirement: change multisig threshold
//     - multisig ClearRequest: cancel all pending (NotConfirmed) multisig requests
//     - multisig WithdrawLockupStaking: withdraw a part of initial lockup stakes
//     - multisig ApproveStakingWithdrawal: approve a withdrawal request
//       - Only available when public delegation is disabled
//     - multisig CancelApprovedStakingWithdrawal: cancel a withdrawal request
//       - Only available when public delegation is disabled
//     - multisig UpdateRewardAddress: update the reward address
//     - multisig UpdateStakingTracker: update the staking tracker
//     - multisig UpdateVoterAddress: update the voter address
//     - multisig ToggleRedelegation: toggle redelegation

// Code organization
// - Modifiers
// - Mutators
//   - Constructor and initializers
//   - Managing operations
//   - Generic facility
//   - Private helpers
// - Getters

contract CnStakingV3MultiSig is ICnStakingV3MultiSig, CnStakingV3MultiSigStorage, CnStakingV3 {
    using EnumerableSet for EnumerableSet.AddressSet;
    using ValidContract for address;

    /* ========== MODIFIERS ========== */

    modifier notConfirmedRequest(uint256 _id) {
        _checkNotConfirmed(_id);
        _;
    }

    modifier validRequirement(uint256 _adminCount, uint256 _requirement) {
        _checkValidRequirement(_adminCount, _requirement);
        _;
    }

    /* ========== CONSTRUCTOR ========== */

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
        validRequirement(_cnAdminlist.length, _requirement)
        CnStakingV3(address(this), _nodeId, _rewardAddress, _unlockTime, _unlockAmount)
    {
        // Initialize the function selector map
        _setFnSelMap();

        // Admin role was granted to address(this) in CnStakingV3, so revoke it here.
        _revokeRole(ADMIN_ROLE, address(this));

        // Sanitize _cnAdminlist
        _grantRole(ADMIN_ROLE, _contractValidator);
        for (uint256 i = 0; i < _cnAdminlist.length; i++) {
            require(
                _cnAdminlist[i] != address(0) && _grantRole(ADMIN_ROLE, _cnAdminlist[i]),
                "Address is null or not unique."
            );
        }

        contractValidator = _contractValidator;
        requirement = _requirement;

        emit DeployCnStakingV3MultiSig(CONTRACT_TYPE, _contractValidator, _cnAdminlist, _requirement);
    }

    /// @dev Initialize the _fnSelMap
    function _setFnSelMap() private {
        _fnSelMap[Functions.AddAdmin] = this.addAdmin.selector;
        _fnSelMap[Functions.DeleteAdmin] = this.deleteAdmin.selector;
        _fnSelMap[Functions.UpdateRequirement] = this.updateRequirement.selector;
        _fnSelMap[Functions.ClearRequest] = this.clearRequest.selector;
        _fnSelMap[Functions.WithdrawLockupStaking] = this.withdrawLockupStaking.selector;
        _fnSelMap[Functions.ApproveStakingWithdrawal] = this.approveStakingWithdrawal.selector;
        _fnSelMap[Functions.CancelApprovedStakingWithdrawal] = this.cancelApprovedStakingWithdrawal.selector;
        _fnSelMap[Functions.UpdateRewardAddress] = this.updateRewardAddress.selector;
        _fnSelMap[Functions.UpdateStakingTracker] = this.updateStakingTracker.selector;
        _fnSelMap[Functions.UpdateVoterAddress] = this.updateVoterAddress.selector;
        _fnSelMap[Functions.ToggleRedelegation] = this.toggleRedelegation.selector;
    }

    /// @dev Completes the contract initialization by depositing initial lockup amounts.
    /// The transaction must send exactly the initial lockup amount of KAIA.
    /// Emits a DepositLockupStakingAndInit event.
    function depositLockupStakingAndInit() public payable override beforeInit {
        _revokeRole(ADMIN_ROLE, contractValidator);
        delete contractValidator;

        super.depositLockupStakingAndInit();
    }

    /// @dev Grants staking related roles to the appropriate addresses.
    /// It overrides the function in CnStakingV3 to handle the multisig admins.
    function _grantStakingRoles() internal override {
        if (!isPublicDelegationEnabled) {
            _grantRole(UNSTAKING_APPROVER_ROLE, address(this));
            for (uint256 i = 0; i < getRoleMemberCount(ADMIN_ROLE); i++) {
                address _operator = getRoleMember(ADMIN_ROLE, i);
                _grantRole(UNSTAKING_CLAIMER_ROLE, _operator);
            }
        } else {
            _grantRole(STAKER_ROLE, publicDelegation);
            _grantRole(UNSTAKING_APPROVER_ROLE, publicDelegation);
            _grantRole(UNSTAKING_CLAIMER_ROLE, publicDelegation);
        }
    }

    /* ========== MULTISIG OPERATIONS (SUBMIT - CONFIRM - EXECUTE) ========== */

    /// @dev Submit a request to add an admin to adminList
    /// @param _admin  new admin address
    function submitAddAdmin(
        address _admin
    )
        external
        override
        afterInit
        onlyRole(ADMIN_ROLE)
        notNull(_admin)
        validRequirement(getRoleMemberCount(ADMIN_ROLE) + 1, requirement)
    {
        require(!hasRole(ADMIN_ROLE, _admin), "Admin already exists.");
        uint256 id = _submitRequest(Functions.AddAdmin, _toBytes32(_admin), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Add an admin to adminList
    /// @param _admin  new admin address
    /// Emits an AddAdmin event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function addAdmin(
        address _admin
    ) external override onlyRole(OPERATOR_ROLE) validRequirement(getRoleMemberCount(ADMIN_ROLE) + 1, requirement) {
        require(_grantRole(ADMIN_ROLE, _admin), "Admin already exists.");
        if (!isPublicDelegationEnabled) {
            require(_grantRole(UNSTAKING_CLAIMER_ROLE, _admin), "Admin already exists.");
        }

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
        onlyRole(ADMIN_ROLE)
        notNull(_admin)
        validRequirement(getRoleMemberCount(ADMIN_ROLE) - 1, requirement)
    {
        require(hasRole(ADMIN_ROLE, _admin), "Admin does not exist.");
        uint256 id = _submitRequest(Functions.DeleteAdmin, _toBytes32(_admin), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Delete an admin from adminList
    /// @param _admin  the admin address
    /// Emits a DeleteAdmin event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function deleteAdmin(
        address _admin
    ) external override onlyRole(OPERATOR_ROLE) validRequirement(getRoleMemberCount(ADMIN_ROLE) - 1, requirement) {
        require(_revokeRole(ADMIN_ROLE, _admin), "Admin does not exist.");
        if (!isPublicDelegationEnabled) {
            require(_revokeRole(UNSTAKING_CLAIMER_ROLE, _admin), "Admin does not exist.");
        }

        clearRequest();
        emit DeleteAdmin(_admin);
    }

    /// @dev Submit a request to update the confirmation threshold
    /// @param _requirement  new confirmation threshold
    function submitUpdateRequirement(
        uint256 _requirement
    ) external override afterInit onlyRole(ADMIN_ROLE) validRequirement(getRoleMemberCount(ADMIN_ROLE), _requirement) {
        require(_requirement != requirement, "Invalid value.");
        uint256 id = _submitRequest(Functions.UpdateRequirement, bytes32(_requirement), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Update the confirmation threshold
    /// @param _requirement  new confirmation threshold
    /// Emits an UpdateRequirement event.
    /// All outstanding requests (i.e. NotConfirmed) are canceled.
    function updateRequirement(
        uint256 _requirement
    ) external override onlyRole(OPERATOR_ROLE) validRequirement(getRoleMemberCount(ADMIN_ROLE), _requirement) {
        requirement = _requirement;
        clearRequest();
        emit UpdateRequirement(_requirement);
    }

    /// @dev Submit a request to cancel all outstanding (i.e. NotConfirmed) requests
    function submitClearRequest() external override afterInit onlyRole(ADMIN_ROLE) {
        uint256 id = _submitRequest(Functions.ClearRequest, 0, 0, 0);
        _confirmRequest(id);
    }

    /// @dev Implicitly cancel all outstanding (i.e. NotConfirmed) requests
    /// by setting lastClearedId to requestCount.
    /// See `getRequestState` for the state transition.
    /// Emits a ClearRequest event.
    function clearRequest() public override onlyRole(OPERATOR_ROLE) {
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
    ) external override afterInit onlyRole(ADMIN_ROLE) notNull(_to) {
        (, , , , uint256 withdrawableAmount) = getLockupStakingInfo();
        require(_value > 0 && _value <= withdrawableAmount, "Invalid value.");

        uint256 id = _submitRequest(Functions.WithdrawLockupStaking, _toBytes32(_to), bytes32(_value), 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to withdraw a part of delegated stakes.
    ///
    /// If public delegation is enabled, the request by multisig is not allowed.
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
    ) external override afterInit onlyRole(ADMIN_ROLE) notNull(_to) {
        require(!isPublicDelegationEnabled, "Public delegation enabled.");
        require(_value > 0 && unstaking + _value <= staking, "Invalid value.");

        uint256 id = _submitRequest(Functions.ApproveStakingWithdrawal, _toBytes32(_to), bytes32(_value), 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to cancel a withdrawal request
    ///
    /// If public delegation is enabled, the request by multisig is not allowed.
    ///
    /// The withdrawal request ID can be obtained from ApproveStakingWithdrawal event
    /// or getApprovedStakingWithdrawalIds().
    /// Unstaking amount decreases.
    function submitCancelApprovedStakingWithdrawal(uint256 _id) external override afterInit onlyRole(ADMIN_ROLE) {
        require(!isPublicDelegationEnabled, "Public delegation enabled.");
        WithdrawalRequest storage request = withdrawalRequestMap[_id];
        require(request.to != address(0), "Withdrawal request does not exist.");
        require(request.state == WithdrawalStakingState.Unknown, "Invalid state.");

        uint256 id = _submitRequest(Functions.CancelApprovedStakingWithdrawal, bytes32(_id), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to update the reward address of this CN
    function submitUpdateRewardAddress(address _addr) external override afterInit onlyRole(ADMIN_ROLE) notNull(_addr) {
        require(!isPublicDelegationEnabled, "Public delegation enabled.");
        uint256 id = _submitRequest(Functions.UpdateRewardAddress, _toBytes32(_addr), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to update the staking tracker this CN reports to
    /// Should not be called if there is an active proposal
    function submitUpdateStakingTracker(address _tracker) external override afterInit onlyRole(ADMIN_ROLE) {
        require(_tracker._validStakingTracker(1), "Invalid StakingTracker.");
        if (stakingTracker != address(0)) {
            IStakingTracker(stakingTracker).refreshStake(address(this));
            require(
                IStakingTracker(stakingTracker).getLiveTrackerIds().length == 0,
                "Cannot update tracker when there is an active tracker."
            );
        }

        uint256 id = _submitRequest(Functions.UpdateStakingTracker, _toBytes32(_tracker), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to update the voter address of this CN
    function submitUpdateVoterAddress(address _addr) external override afterInit onlyRole(ADMIN_ROLE) {
        if (stakingTracker != address(0) && _addr != address(0)) {
            require(IStakingTracker(stakingTracker).voterToGCId(_addr) == 0, "Voter already taken.");
        }
        uint256 id = _submitRequest(Functions.UpdateVoterAddress, _toBytes32(_addr), 0, 0);
        _confirmRequest(id);
    }

    /// @dev Submit a request to toggle redelegation
    /// Only available when public delegation is enabled
    function submitToggleRedelegation() external override afterInit onlyRole(ADMIN_ROLE) {
        require(isPublicDelegationEnabled, "Public delegation disabled.");
        uint256 id = _submitRequest(Functions.ToggleRedelegation, 0, 0, 0);
        _confirmRequest(id);
    }

    // Generic multisig facility

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
    ) public override notConfirmedRequest(_id) onlyRole(ADMIN_ROLE) {
        require(!_hasConfirmed(_id, _msgSender()), "Msg.sender already confirmed.");
        Request storage _req = _requestMap[_id];
        require(
            _req.functionId == _functionId &&
                _req.firstArg == _firstArg &&
                _req.secondArg == _secondArg &&
                _req.thirdArg == _thirdArg,
            "Function id and arguments do not match."
        );

        _req.confirmers.add(_msgSender());
        emit ConfirmRequest(_id, _msgSender(), _functionId, _firstArg, _secondArg, _thirdArg, _req.confirmers.values());

        if (_req.confirmers.length() >= requirement) {
            _executeRequest(_id);
        }
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
    ) external override notConfirmedRequest(_id) onlyRole(ADMIN_ROLE) {
        require(_hasConfirmed(_id, _msgSender()), "Msg.sender has not confirmed.");
        Request storage _req = _requestMap[_id];
        require(
            _req.functionId == _functionId &&
                _req.firstArg == _firstArg &&
                _req.secondArg == _secondArg &&
                _req.thirdArg == _thirdArg,
            "Function id and arguments do not match."
        );

        if (_req.requestProposer == _msgSender()) {
            _req.state = RequestState.Canceled;
            emit CancelRequest(_id, _msgSender(), _req.functionId, _req.firstArg, _req.secondArg, _req.thirdArg);
        } else {
            _req.confirmers.remove(_msgSender());
            emit RevokeConfirmation(
                _id,
                _msgSender(),
                _req.functionId,
                _req.firstArg,
                _req.secondArg,
                _req.thirdArg,
                _req.confirmers.values()
            );
        }
    }

    /// @dev Submit a request
    /// Emits a SubmitRequest event.
    /// @return  the request ID
    function _submitRequest(
        Functions _functionId,
        bytes32 _firstArg,
        bytes32 _secondArg,
        bytes32 _thirdArg
    ) private returns (uint256) {
        uint256 id = requestCount;
        unchecked {
            requestCount++;
        }

        Request storage req = _requestMap[id];
        req.functionId = _functionId;
        req.firstArg = _firstArg;
        req.secondArg = _secondArg;
        req.thirdArg = _thirdArg;
        req.requestProposer = _msgSender();
        req.state = RequestState.NotConfirmed;

        emit SubmitRequest(id, _msgSender(), _functionId, _firstArg, _secondArg, _thirdArg);
        return id;
    }

    /// @dev Shortcut of confirmRequest(...)
    /// Used by submitXYZ() functions.
    function _confirmRequest(uint256 id) private {
        confirmRequest(
            id,
            _requestMap[id].functionId,
            _requestMap[id].firstArg,
            _requestMap[id].secondArg,
            _requestMap[id].thirdArg
        );
    }

    /// @dev execute a requested function
    /// Used by confirmRequest when enough confirmations are made.
    /// Emits a ExecuteRequestSuccess or ExecuteRequestFailure event.
    function _executeRequest(uint256 _id) private {
        bool ok = false;
        Request storage req = _requestMap[_id];
        Functions funcId = req.functionId;
        bytes32 a1 = req.firstArg;
        bytes32 a2 = req.secondArg;
        bytes32 a3 = req.thirdArg;

        bytes memory _data = abi.encodeWithSelector(_fnSelMap[funcId], a1, a2, a3);
        (ok, ) = address(this).call(_data);

        if (ok) {
            req.state = RequestState.Executed;
            emit ExecuteRequestSuccess(_id, _msgSender(), funcId, a1, a2, a3);
        } else {
            req.state = RequestState.ExecutionFailed;
            emit ExecuteRequestFailure(_id, _msgSender(), funcId, a1, a2, a3);
        }
    }

    /* ========== PRIVATE HELPERS ========== */

    /// @dev Convert an address to bytes32
    function _toBytes32(address _x) private pure returns (bytes32) {
        return bytes32(uint256(uint160(_x)));
    }

    /// @dev Check if the given admin has confirmed the request
    function _hasConfirmed(uint256 _id, address addr) private view returns (bool) {
        return _requestMap[_id].confirmers.contains(addr);
    }

    /// @dev Check if the given request is at NotConfirmed state
    function _checkNotConfirmed(uint256 _id) private view {
        require(getRequestState(_id) == RequestState.NotConfirmed, "Must be at not-confirmed state.");
    }

    /// @dev Check if the given requirement is valid
    function _checkValidRequirement(uint256 _adminCount, uint256 _requirement) private pure {
        require(
            _adminCount <= MAX_ADMIN && _requirement <= _adminCount && _requirement != 0 && _adminCount != 0,
            "Invalid requirement."
        );
    }

    /// @dev Return the all admin addresses
    function _getAdminList() private view returns (address[] memory admins) {
        uint256 len = getRoleMemberCount(ADMIN_ROLE);
        admins = new address[](len);
        for (uint256 i = 0; i < len; i++) {
            admins[i] = getRoleMember(ADMIN_ROLE, i);
        }
    }

    /* ========== PUBLIC GETTERS ========== */
    /// Note Below `admin` related getters are for the convenience of the frontend to be compatible with CnStakingV2.

    /// @dev Check if the given address is an admins
    function isAdmin(address _admin) external view override returns (bool) {
        return hasRole(ADMIN_ROLE, _admin);
    }

    /// @dev Return an admin address at the given position
    function adminList(uint256 _pos) external view override returns (address) {
        return getRoleMember(ADMIN_ROLE, _pos);
    }

    /// @dev Return the current state of the request
    function getRequestState(uint256 _id) public view override returns (RequestState) {
        if (_id >= requestCount) {
            return RequestState.Unknown;
        }

        if (_requestMap[_id].state != RequestState.NotConfirmed) {
            return _requestMap[_id].state;
        }

        return _id < lastClearedId ? RequestState.Canceled : RequestState.NotConfirmed;
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
            address[] memory _adminListArr,
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
            _getAdminList(),
            requirement,
            lockupConditions.unlockTime,
            lockupConditions.unlockAmount,
            lockupConditions.allReviewed,
            isInitialized
        );
    }

    /// @dev Return the reviewers of the initial lockup conditions
    /// @return  reviewers addresses
    function getReviewers() external view override beforeInit returns (address[] memory) {
        address[] memory reviewers = new address[](lockupConditions.reviewedCount);
        uint256 id = 0;
        for (uint256 i = 0; i < getRoleMemberCount(ADMIN_ROLE); i++) {
            address _admin = getRoleMember(ADMIN_ROLE, i);
            if (lockupConditions.reviewedAdmin[_admin] == true) {
                reviewers[id] = _admin;
                unchecked {
                    id++;
                }
            }
        }

        return reviewers;
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
        uint256 _end = (_to == 0 || _to >= requestCount) ? requestCount : _to;

        ids = new uint256[](_end - _from);
        uint cnt = 0;
        for (uint i = _from; i < _end; i++) {
            if (getRequestState(i) == _state) {
                unchecked {
                    ids[cnt++] = i;
                }
            }
        }

        assembly {
            mstore(ids, cnt)
        }
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
        Request storage r = _requestMap[_id];
        return (
            r.functionId,
            r.firstArg,
            r.secondArg,
            r.thirdArg,
            r.requestProposer,
            r.confirmers.values(),
            getRequestState(_id)
        );
    }
}
