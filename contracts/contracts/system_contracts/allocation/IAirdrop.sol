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

interface IAirdrop {
    /* ========== EVENT ========== */

    event Claimed(address indexed beneficiary, uint256 amount);

    /* ========== VIEWS ========== */

    function claimAllowed() external view returns (bool);

    function claims(address) external view returns (uint256);

    function claimed(address) external view returns (bool);

    function getBeneficiaries(uint256 start, uint256 end) external view returns (address[] memory);

    function getBeneficiaryAt(uint256 index) external view returns (address);

    function getBeneficiariesLength() external view returns (uint256);

    /* ========== MUTATIVE FUNCTIONS ========== */

    receive() external payable;

    function toggleClaimAllowed() external;

    function addClaim(address beneficiary, uint256 amount) external;

    function addBatchClaims(address[] calldata beneficiaries, uint256[] calldata amounts) external;

    function claim() external;

    function claimFor(address beneficiary) external;

    function claimBatch(address[] calldata beneficiaries) external;
}
