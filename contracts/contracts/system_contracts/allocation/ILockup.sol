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

interface ILockup {
    /* ========== DATA STRUCTURE ========== */

    enum AcquisitionStatus {
        UNDEFINED,
        PROPOSED,
        CONFIRMED,
        WITHDRAWN,
        REJECTED
    }

    struct Acquisition {
        uint256 acReqId;
        uint256 amount;
        AcquisitionStatus status;
    }

    struct DelegatedTransfer {
        uint256 delegatedTransferId;
        uint256 amount;
        address to;
        AcquisitionStatus status; // Reuse the AcquisitionStatus enum
    }

    /* ========== EVENTS ========== */

    event RefreshedDelegated(uint256 totalDelegatedAmount);

    event ProposeAcquisition(uint256 acReqId, uint256 amount);

    event RequestDelegatedTransfer(uint256 delegatedTransferId, uint256 amount, address to);

    event WithdrawAcquisition(uint256 acReqId);

    event WithdrawDelegatedTransfer(uint256 delegatedTransferId);

    event WithdrawStakingAmounts(address pdKaia, uint256 shares);

    event ClaimStakingAmounts(address pdKaia, uint256 requestId);

    event ConfirmAcquisition(uint256 acReqId);

    event RejectAcquisition(uint256 acReqId);

    event ConfirmDelegatedTransfer(uint256 delegatedTransferId);

    event RejectDelegatedTransfer(uint256 delegatedTransferId);

    /* ========== VIEWS ========== */

    function ADMIN_ROLE() external view returns (bytes32);

    function SECRETARY_ROLE() external view returns (bytes32);

    function nextAcReqId() external view returns (uint256);

    function nextDelegatedTransferId() external view returns (uint256);

    function getAllAcquisitions() external view returns (Acquisition[] memory);

    function getAllDelegatedTransfers() external view returns (DelegatedTransfer[] memory);

    function getAcquisition(uint256 acReqId) external view returns (Acquisition memory);

    function getDelegatedTransfer(uint256 delegatedTransferId) external view returns (DelegatedTransfer memory);

    function getAcquisitionAtStatus(AcquisitionStatus status) external view returns (Acquisition[] memory);

    function getDelegatedTransferAtStatus(AcquisitionStatus status) external view returns (DelegatedTransfer[] memory);

    /* ========== MUTATIVE FUNCTIONS ========== */

    receive() external payable;

    function transferAdmin(address newAdmin) external;

    function transferSecretary(address newSecretary) external;

    function refreshDelegated() external;

    function proposeAcquisition(uint256 amount) external;

    function requestDelegatedTransfer(uint256 amount, address to) external;

    function withdrawAcquisition(uint256 acReqId) external;

    function withdrawDelegatedTransfer(uint256 delegatedTransferId) external;

    function withdrawStakingAmounts(address pdKaia, uint256 shares) external;

    function claimStakingAmounts(address pdKaia, uint256 requestId) external;

    function confirmAcquisition(uint256 acReqId) external;

    function rejectAcquisition(uint256 acReqId) external;

    function confirmDelegatedTransfer(uint256 delegatedTransferId) external;

    function rejectDelegatedTransfer(uint256 delegatedTransferId) external;
}
