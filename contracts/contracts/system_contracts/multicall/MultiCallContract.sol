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
pragma solidity 0.8.19;

interface IAddressBook {
    function isActivated() external view returns (bool);

    function getAllAddress() external view returns (uint8[] memory typeList, address[] memory addressList);
}

interface IRegistry {
    function getActiveAddr(string memory name) external view returns (address);
}

interface ICLRegistry {
    function getAllCLs() external view returns (address[] memory, uint256[] memory, address[] memory, address[] memory);
}

interface ICnStaking {
    function VERSION() external view returns (uint256);

    function staking() external view returns (uint256);

    function unstaking() external view returns (uint256);
}

interface IERC20 {
    function balanceOf(address account) external view returns (uint256);
}

// MultiCallContract provides a function to retrieve the any information needed for the Kaia client.
// It will be temporarily injected into state to be used by the Kaia client.
// After retrieving the information, the contract will be removed from the state.
contract MultiCallContract {
    address private constant ADDRESS_BOOK_ADDRESS = 0x0000000000000000000000000000000000000400;
    address private constant REGISTRY_ADDRESS = 0x0000000000000000000000000000000000000401;

    /* ========== STAKING INFORMATION ========== */

    // multiCallStakingInfo returns the staking information of all CNs.
    function multiCallStakingInfo()
        external
        view
        returns (uint8[] memory typeList, address[] memory addressList, uint256[] memory stakingAmounts)
    {
        IAddressBook addressBook = IAddressBook(ADDRESS_BOOK_ADDRESS);
        (typeList, addressList) = addressBook.getAllAddress();

        // Return early if AddressBook hasn't been activated yet or there are no CNs.
        if (addressList.length < 5) {
            return (typeList, addressList, stakingAmounts);
        }

        uint256 lenCnAddress = addressList.length - 2;
        // Just in case.
        if (lenCnAddress % 3 != 0) {
            return (typeList, addressList, stakingAmounts);
        }
        stakingAmounts = new uint256[](lenCnAddress / 3);

        for (uint256 i = 0; i < lenCnAddress; i += 3) {
            stakingAmounts[i / 3] = _getCnStakingAmounts(addressList[i + 1]);
        }

        return (typeList, addressList, stakingAmounts);
    }

    function _getCnStakingAmounts(address cnStaking) private view returns (uint256) {
        return cnStaking.balance;
    }

    function multiCallDPStakingInfo()
        external
        view
        returns (
            address[] memory nodeIds,
            address[] memory clPools,
            address[] memory clStakings,
            uint256[] memory stakingAmounts
        )
    {
        address clreg = IRegistry(REGISTRY_ADDRESS).getActiveAddr("CLRegistry");
        address wKaia = IRegistry(REGISTRY_ADDRESS).getActiveAddr("WrappedKaia");

        if (clreg == address(0)) {
            return (nodeIds, clPools, clStakings, stakingAmounts);
        }

        (nodeIds, , clPools, clStakings) = ICLRegistry(clreg).getAllCLs();
        uint256 poolLength = clPools.length;
        stakingAmounts = new uint256[](poolLength);

        if (wKaia != address(0)) {
            IERC20 wKaiaToken = IERC20(wKaia);
            for (uint256 i = 0; i < poolLength; i++) {
                stakingAmounts[i] = wKaiaToken.balanceOf(clPools[i]);
            }
        }
    }

    /* ========== MORE FUNCTIONS TBA ========== */
}
