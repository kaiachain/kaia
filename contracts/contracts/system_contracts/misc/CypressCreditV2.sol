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

pragma solidity 0.4.24;

import "./CypressCredit.sol";

/// @title Cypress Genesis and Ending Block Information

// CypressCreditV2 keeps the information about Klaytn team when this contract is
// created, and it replaces the CypressCredit contract in the Klaytn Cypress mainnet.
contract CypressCreditV2 is CypressCredit {
    // getPhoto returns a base64-encoded photo of the Klaytn team.
    function getEndingPhoto() public pure returns (string) {
        string memory photo = "EndingPhoto TBU";
        return photo;
    }

    // getNames returns comma-separated names of all the members in the Klaytn team.
    function getEndingNames() public pure returns (string) {
        string memory names = "EndingNames TBU";
        return names;
    }
}
