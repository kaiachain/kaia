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

import "./IAirdrop.sol";
import "openzeppelin-contracts-5.0/access/Ownable.sol";

contract Airdrop is Ownable, IAirdrop {
    /* ========== CONSTANTS ========== */

    uint256 public constant KAIA_UNIT = 1e18;

    uint256 public constant TOTAL_AIRDROP_AMOUNT = 80_000_000 * KAIA_UNIT;

    /* ========== STATE VARIABLES ========== */

    mapping(address => uint256) public claims;

    mapping(address => bool) public claimed;

    /* ========== CONSTRUCTOR ========== */

    constructor() Ownable(msg.sender) {}

    /* ========== OPERATOR FUNCTIONS ========== */

    function addClaim(address _beneficiary, uint256 _amount) external override onlyOwner {
        claims[_beneficiary] = _amount;
    }

    function addBatchClaims(
        address[] calldata _beneficiaries,
        uint256[] calldata _amounts
    ) external override onlyOwner {
        require(_beneficiaries.length == _amounts.length, "Airdrop: invalid input length");

        for (uint256 i = 0; i < _beneficiaries.length; i++) {
            claims[_beneficiaries[i]] = _amounts[i];
        }
    }

    /* ========== PUBLIC FUNCTIONS ========== */

    function claim() external override {
        _claim(msg.sender);
    }

    function claimFor(address _beneficiary) public override {
        _claim(_beneficiary);
    }

    function claimBatch(address[] calldata _beneficiary) external override {
        for (uint256 i = 0; i < _beneficiary.length; i++) {
            _claim(_beneficiary[i]);
        }
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _claim(address _beneficiary) internal {
        require(!claimed[_beneficiary], "Airdrop: already claimed");

        uint256 _amount = claims[_beneficiary];
        require(_amount > 0 && _amount <= address(this).balance, "Airdrop: no claimable amount");

        claimed[_beneficiary] = true;

        (bool success, ) = _beneficiary.call{value: _amount}("");
        require(success, "Transfer failed.");

        emit Claimed(_beneficiary, _amount);
    }
}
