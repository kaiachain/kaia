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
import "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";
import "openzeppelin-contracts-5.0/access/Ownable.sol";

contract Airdrop is Ownable, IAirdrop {
    using EnumerableSet for EnumerableSet.AddressSet;

    /* ========== CONSTANTS ========== */

    uint256 public constant KAIA_UNIT = 1e18;

    uint256 public constant TOTAL_AIRDROP_AMOUNT = 80_000_000 * KAIA_UNIT;

    /* ========== STATE VARIABLES ========== */

    EnumerableSet.AddressSet private _beneficiaries;

    mapping(address => uint256) public claims;

    mapping(address => bool) public claimed;

    /* ========== CONSTRUCTOR ========== */

    constructor() Ownable(msg.sender) {}

    /* ========== OPERATOR FUNCTIONS ========== */

    function addClaim(address beneficiary, uint256 amount) external override onlyOwner {
        _addClaim(beneficiary, amount);
    }

    function addBatchClaims(address[] calldata beneficiaries, uint256[] calldata amounts) external override onlyOwner {
        require(beneficiaries.length == amounts.length, "Airdrop: invalid input length");

        for (uint256 i = 0; i < beneficiaries.length; i++) {
            _addClaim(beneficiaries[i], amounts[i]);
        }
    }

    /* ========== PUBLIC FUNCTIONS ========== */

    function claim() external override {
        _claim(msg.sender);
    }

    function claimFor(address beneficiary) public override {
        _claim(beneficiary);
    }

    function claimBatch(address[] calldata beneficiary) external override {
        for (uint256 i = 0; i < beneficiary.length; i++) {
            _claim(beneficiary[i]);
        }
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _addClaim(address beneficiary, uint256 amount) internal {
        require(_beneficiaries.add(beneficiary), "Airdrop: beneficiary already exists");
        claims[beneficiary] = amount;
    }

    function _claim(address beneficiary) internal {
        require(!claimed[beneficiary], "Airdrop: already claimed");

        uint256 _amount = claims[beneficiary];
        require(_amount > 0 && _amount <= address(this).balance, "Airdrop: no claimable amount");

        claimed[beneficiary] = true;

        (bool success, ) = beneficiary.call{value: _amount}("");
        require(success, "Transfer failed.");

        emit Claimed(beneficiary, _amount);
    }

    /* ========== GETTERS ========== */

    function getBeneficiaries(uint256 start, uint256 end) external override view returns (address[] memory result) {
        end = end > _beneficiaries.length() ? _beneficiaries.length() : end;
        if (start >= end) {
            return new address[](0);
        }

        result = new address[](end - start);
        for (uint256 i = start; i < end; i++) { 
            unchecked {
                result[i - start] = _beneficiaries.at(i);
            }
        }
    }

    function getBeneficiaryAt(uint256 index) external override view returns (address) {
        return _beneficiaries.at(index);
    }

    function getBeneficiariesLength() external override view returns (uint256) {
        return _beneficiaries.length();
    }
}
