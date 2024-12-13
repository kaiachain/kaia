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

pragma solidity ^0.4.24;

import "./AddressBookMock.sol";

/**
 * @title CLRegistryMockTwoCL
 */

contract CLRegistryMockThreeCL is MockValues {
    function getAllCLs()
        external
        view
        returns (address[] memory, uint256[] memory, address[] memory, address[] memory)
    {
        address[] memory nodeIds = new address[](3);
        uint256[] memory gcIds = new uint256[](3);
        address[] memory clPools = new address[](3);
        address[] memory clStakings = new address[](3);

        nodeIds[0] = nodeId0;
        nodeIds[1] = nodeId1;
        nodeIds[2] = nodeId2; // Doesn't exist in AddressBookMockTwoCN

        gcIds[0] = 1;
        gcIds[1] = 2;
        gcIds[2] = 3;

        clPools[0] = 0x0000000000000000000000000000000000000e00;
        clPools[1] = 0x0000000000000000000000000000000000000e01;
        clPools[2] = 0x0000000000000000000000000000000000000e02;

        clStakings[0] = 0x0000000000000000000000000000000000000e03;
        clStakings[1] = 0x0000000000000000000000000000000000000e04;
        clStakings[2] = 0x0000000000000000000000000000000000000e05;

        return (nodeIds, gcIds, clPools, clStakings);
    }
}

contract RegistryMockForCL {
    // This is a mock implementation of the Registry contract
    // It returns a fixed address for the CLRegistryMockThreeCL address
    function getActiveAddr(string name) external view returns (address) {
        address clRegistryAddr = 0x0000000000000000000000000000000000000Ff0;
        address wrappedKaiaAddr = 0x0000000000000000000000000000000000000Ff1;

        if (keccak256(name) == keccak256("CLRegistry")) {
            return clRegistryAddr;
        } else if (keccak256(name) == keccak256("WrappedKaia")) {
            return wrappedKaiaAddr;
        }
        return address(0);
    }
}

contract RegistryMockZero {
    // This is a mock implementation of the Registry contract
    // It returns a fixed address for the CLRegistryMockThreeCL address
    function getActiveAddr(string name) external view returns (address) {
        return address(0);
    }
}

contract WrappedKaiaMock {
    function balanceOf(address account) external view returns (uint256) {   
        return account.balance;
    }
}
