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

    function KAIA_UNIT() external view returns (uint256);

    function TOTAL_AIRDROP_AMOUNT() external view returns (uint256);

    function claims(address) external view returns (uint256);

    function claimed(address) external view returns (bool);

    /* ========== MUTATIVE FUNCTIONS ========== */

    function addClaim(address _beneficiary, uint256 _amount) external;

    function addBatchClaims(address[] calldata _beneficiaries, uint256[] calldata _amounts) external;

    function claim() external;

    function claimFor(address _beneficiary) external;

    function claimBatch(address[] calldata _beneficiaries) external;
}
