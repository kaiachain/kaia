// Copyright 2023 The klaytn Authors
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
pragma solidity ^0.8.18;

contract MultiCallContractMock {
    function multiCallStakingInfo()
        external
        view
        returns (uint8[] memory typeList, address[] memory addressList, uint256[] memory stakingAmounts)
    {
        typeList = new uint8[](5);
        addressList = new address[](5);
        stakingAmounts = new uint256[](1);

        typeList[0] = 0; // Node address
        typeList[1] = 1; // Staking address
        typeList[2] = 2; // Reward address
        typeList[3] = 3; // POC address
        typeList[4] = 4; // KIR address

        addressList[0] = 0x0000000000000000000000000000000000000F00;
        addressList[1] = 0x0000000000000000000000000000000000000F01;
        addressList[2] = 0x0000000000000000000000000000000000000f02;
        addressList[3] = 0x0000000000000000000000000000000000000F03;
        addressList[4] = 0x0000000000000000000000000000000000000f04;

        stakingAmounts[0] = 7_000_000 ether;
    }
}
