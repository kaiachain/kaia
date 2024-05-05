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

pragma solidity 0.8.25;

import "../system_contracts/consensus/IAddressBook.sol";
import "../system_contracts/consensus/IStakingTracker.sol";

interface ICnStaking {
    function CONTRACT_TYPE() external view returns (string memory);

    function VERSION() external view returns (uint256);

    function nodeId() external view returns (address);
}

library ValidContract {
    /// @dev Checks if a given address is valid CnStaking contract
    /// Note that it might be reverted earlier if the contract doesn't have expected methods
    function _validCnStaking(address _staking, uint256 __version) internal view returns (bool) {
        if (_staking.code.length == 0) return false;

        ICnStaking _cnStaking = ICnStaking(_staking);

        if (
            keccak256(bytes(_cnStaking.CONTRACT_TYPE())) !=
            hex"a2f5d64a9f0bcdeed97e196203f5a8c1a5c8293988b625b7925686d308055082" // keccak256("CnStakingContract")
        ) return false;

        if (_cnStaking.VERSION() != __version) return false;

        address _nodeId = _cnStaking.nodeId();
        (, address _cnInAB, ) = IAddressBook(0x0000000000000000000000000000000000000400).getCnInfo(_nodeId);
        // Just in case. It was already checked in `getCnInfo`.
        return _cnInAB == _staking;
    }

    /// @dev Checks if a given address is valid StakingTracker contract
    /// Note that it might be reverted earlier if the contract doesn't have expected methods
    function _validStakingTracker(address _tracker, uint256 __version) internal view returns (bool) {
        if (_tracker.code.length == 0) return false;

        IStakingTracker _stakingTracker = IStakingTracker(_tracker);

        if (
            keccak256(bytes(_stakingTracker.CONTRACT_TYPE())) !=
            hex"ccfe28814eb3e9d0e6cfd45eb754f27c5eb4399dac6379181362ebd8b6a865c3" // keccak256("StakingTracker")
        ) return false;

        if (_stakingTracker.VERSION() != __version) return false;

        return true;
    }
}
