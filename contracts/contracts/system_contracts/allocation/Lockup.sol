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

import "./ILockup.sol";
import "../consensus/PublicDelegation/IPublicDelegation.sol";
import "openzeppelin-contracts-5.0/access/extensions/AccessControlEnumerable.sol";

contract Lockup is AccessControlEnumerable, ILockup {
    /* ========== CONSTANTS ========== */

    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");

    bytes32 public constant SECRETARY_ROLE = keccak256("SECRETARY_ROLE");

    /* ========== STATE VARIABLES ========== */

    // Acquisition request list
    Acquisition[] private _acquisitions;

    // Delegated transfer request list
    // Note that the `DelegatedTransfer` is used when stake the lockup amount, not the `Acquisition`.
    // If the `DelegatedTransfer.to` is not the staking related contract, it won't be confirmed by the secretary.
    DelegatedTransfer[] private _delegatedTransfers;

    uint256 public totalDelegatedAmount;

    bool public isInitialized;

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

    /* ========== CONSTRUCTOR ========== */

    constructor(address _admin, address _secretary) notNull(_admin) notNull(_secretary) {
        _grantRole(ADMIN_ROLE, _admin);
        _grantRole(SECRETARY_ROLE, _secretary);
    }

    // Since the KAIA will directly set to the contract, the admin needs to manually refresh the delegated amount.
    function refreshDelegated() external override onlyRole(ADMIN_ROLE) beforeInit {
        require(address(this).balance > 0, "Lockup: no claimable amount");

        totalDelegatedAmount = address(this).balance;

        isInitialized = true;

        emit RefreshedDelegated(totalDelegatedAmount);
    }

    /* ========== ROLE MANAGEMENT FUNCTIONS ========== */

    function transferAdmin(address _newAdmin) external override onlyRole(ADMIN_ROLE) notNull(_newAdmin) {
        _grantRole(ADMIN_ROLE, _newAdmin);
        _revokeRole(ADMIN_ROLE, msg.sender);
    }

    function transferSecretary(
        address _newSecretary
    ) external override onlyRole(SECRETARY_ROLE) notNull(_newSecretary) {
        _grantRole(SECRETARY_ROLE, _newSecretary);
        _revokeRole(SECRETARY_ROLE, msg.sender);
    }

    /* ========== ADMIN REQUEST FUNCTIONS ========== */

    function proposeAcquisition(uint256 _amount) external override onlyRole(ADMIN_ROLE) afterInit {
        require(_amount > 0 && _amount <= address(this).balance - _totalPendingAmount(), "Lockup: invalid amount");

        uint256 _acReqId = nextAcReqId();

        _acquisitions.push(Acquisition(_acquisitions.length, _amount, AcquisitionStatus.PROPOSED));

        emit ProposeAcquisition(_acReqId, _amount);
    }

    function requestDelegatedTransfer(
        uint256 _amount,
        address _to
    ) external override onlyRole(ADMIN_ROLE) afterInit notNull(_to) {
        require(_amount > 0 && _amount <= address(this).balance - _totalPendingAmount(), "Lockup: invalid amount");

        uint256 _delegatedTransferId = nextDelegatedTransferId();

        _delegatedTransfers.push(
            DelegatedTransfer(_delegatedTransfers.length, _amount, _to, AcquisitionStatus.PROPOSED)
        );

        emit RequestDelegatedTransfer(_delegatedTransferId, _amount, _to);
    }

    /* ========== ADMIN WITHDRAW FUNCTIONS ========== */

    function withdrawAcquisition(uint256 _acReqId) external override onlyRole(ADMIN_ROLE) afterInit {
        require(_acquisitions[_acReqId].status == AcquisitionStatus.CONFIRMED, "Lockup: invalid status");

        _acquisitions[_acReqId].status = AcquisitionStatus.WITHDRAWN;

        (bool success, ) = msg.sender.call{value: _acquisitions[_acReqId].amount}("");
        require(success, "Lockup: transfer failed");

        emit WithdrawAcquisition(_acReqId);
    }

    function withdrawDelegatedTransfer(uint256 _delegatedTransferId) external override onlyRole(ADMIN_ROLE) afterInit {
        require(
            _delegatedTransfers[_delegatedTransferId].status == AcquisitionStatus.CONFIRMED,
            "Lockup: invalid status"
        );

        _delegatedTransfers[_delegatedTransferId].status = AcquisitionStatus.WITHDRAWN;

        (bool success, ) = _delegatedTransfers[_delegatedTransferId].to.call{
            value: _delegatedTransfers[_delegatedTransferId].amount
        }("");
        require(success, "Lockup: transfer failed");

        emit WithdrawDelegatedTransfer(_delegatedTransferId);
    }

    /* ========== ADMIN PUBLIC DELEGATION INTERACTIONS ========== */

    function withdrawStakingAmounts(address _pdKaia, uint256 _shares) external override onlyRole(ADMIN_ROLE) afterInit {
        IPublicDelegation pdKaia = IPublicDelegation(payable(_pdKaia));
        require(_shares > 0 && pdKaia.maxRedeem(address(this)) >= _shares, "Lockup: invalid shares");

        pdKaia.redeem(address(this), _shares);

        emit WithdrawStakingAmounts(_pdKaia, _shares);
    }

    function claimStakingAmounts(address _pdKaia, uint256 _requestId) external override onlyRole(ADMIN_ROLE) afterInit {
        IPublicDelegation pdKaia = IPublicDelegation(payable(_pdKaia));
        require(pdKaia.requestIdToOwner(_requestId) == address(this), "Lockup: invalid request owner");

        pdKaia.claim(_requestId);

        emit ClaimStakingAmounts(_pdKaia, _requestId);
    }

    /* ========== SECRETARY FUNCTIONS ========== */

    function confirmAcquisition(uint256 _acReqId) external override onlyRole(SECRETARY_ROLE) afterInit {
        require(_acquisitions[_acReqId].status == AcquisitionStatus.PROPOSED, "Lockup: invalid status");

        _acquisitions[_acReqId].status = AcquisitionStatus.CONFIRMED;

        emit ConfirmAcquisition(_acReqId);
    }

    function rejectAcquisition(uint256 _acReqId) external override onlyRole(SECRETARY_ROLE) afterInit {
        require(_acquisitions[_acReqId].status == AcquisitionStatus.PROPOSED, "Lockup: invalid status");

        _acquisitions[_acReqId].status = AcquisitionStatus.REJECTED;

        emit RejectAcquisition(_acReqId);
    }

    function confirmDelegatedTransfer(
        uint256 _delegatedTransferId
    ) external override onlyRole(SECRETARY_ROLE) afterInit {
        require(
            _delegatedTransfers[_delegatedTransferId].status == AcquisitionStatus.PROPOSED,
            "Lockup: invalid status"
        );

        _delegatedTransfers[_delegatedTransferId].status = AcquisitionStatus.CONFIRMED;

        emit ConfirmDelegatedTransfer(_delegatedTransferId);
    }

    function rejectDelegatedTransfer(
        uint256 _delegatedTransferId
    ) external override onlyRole(SECRETARY_ROLE) afterInit {
        require(
            _delegatedTransfers[_delegatedTransferId].status == AcquisitionStatus.PROPOSED,
            "Lockup: invalid status"
        );

        _delegatedTransfers[_delegatedTransferId].status = AcquisitionStatus.REJECTED;

        emit RejectDelegatedTransfer(_delegatedTransferId);
    }

    /* ========== PUBLIC FUNCTIONS ========== */

    // `receive` is used when depositing the KAIA that withdrew for staking, which is `DelegatedTransfer`.
    receive() external payable override {}

    /* ========== INTERNAL UTILS ========== */

    function _checkNull(address _address) private pure {
        require(_address != address(0), "Address is null.");
    }

    function _checkInit() private view {
        require(isInitialized, "Contract is not initialized.");
    }

    function _checkNotInit() private view {
        require(!isInitialized, "Contract has been initialized.");
    }

    function _totalPendingAmount() private view returns (uint256 amount) {
        return _pendingAcquisitionAmount() + _pendingDelegatedTransferAmount();
    }

    function _pendingAcquisitionAmount() private view returns (uint256 amount) {
        for (uint256 i = 0; i < _acquisitions.length; i++) {
            if (
                _acquisitions[i].status == AcquisitionStatus.PROPOSED ||
                _acquisitions[i].status == AcquisitionStatus.CONFIRMED
            ) {
                amount += _acquisitions[i].amount;
            }
        }

        return amount;
    }

    function _pendingDelegatedTransferAmount() private view returns (uint256 amount) {
        for (uint256 i = 0; i < _delegatedTransfers.length; i++) {
            if (
                _delegatedTransfers[i].status == AcquisitionStatus.PROPOSED ||
                _delegatedTransfers[i].status == AcquisitionStatus.CONFIRMED
            ) {
                amount += _delegatedTransfers[i].amount;
            }
        }

        return amount;
    }

    /* ========== GETTERS ========== */

    function nextAcReqId() public view override returns (uint256) {
        return _acquisitions.length;
    }

    function nextDelegatedTransferId() public view override returns (uint256) {
        return _delegatedTransfers.length;
    }

    function getAllAcquisitions() external view override returns (Acquisition[] memory) {
        return _acquisitions;
    }

    function getAllDelegatedTransfers() external view override returns (DelegatedTransfer[] memory) {
        return _delegatedTransfers;
    }

    function getAcquisition(uint256 _acReqId) external view override returns (Acquisition memory) {
        return _acquisitions[_acReqId];
    }

    function getDelegatedTransfer(
        uint256 _delegatedTransferId
    ) external view override returns (DelegatedTransfer memory) {
        return _delegatedTransfers[_delegatedTransferId];
    }

    function getAcquisitionAtStatus(
        AcquisitionStatus _status
    ) external view override returns (Acquisition[] memory acquisitions) {
        acquisitions = new Acquisition[](_acquisitions.length);
        uint256 count = 0;

        for (uint256 i = 0; i < _acquisitions.length; i++) {
            if (_acquisitions[i].status == _status) {
                acquisitions[count] = _acquisitions[i];
                count++;
            }
        }

        assembly {
            mstore(acquisitions, count)
        }
    }

    function getDelegatedTransferAtStatus(
        AcquisitionStatus _status
    ) external view override returns (DelegatedTransfer[] memory delegatedTransfers) {
        delegatedTransfers = new DelegatedTransfer[](_delegatedTransfers.length);
        uint256 count = 0;

        for (uint256 i = 0; i < _delegatedTransfers.length; i++) {
            if (_delegatedTransfers[i].status == _status) {
                delegatedTransfers[count] = _delegatedTransfers[i];
                count++;
            }
        }

        assembly {
            mstore(delegatedTransfers, count)
        }
    }
}
