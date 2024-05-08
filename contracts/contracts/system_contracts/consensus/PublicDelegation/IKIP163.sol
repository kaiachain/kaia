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

/// @title KIP-163 Public Delegation Interface
/// @dev See https://github.com/klaytn/kips/issues/163
interface IKIP163 {
    /// @dev Stake KAIA for the _recipient.
    /// It is used in CnStakingV3.handleRedelegation to stake KAIA for the _recipient when handling redelegation.
    /// It must stake KAIA to the CnStakingV3 in the same transaction.
    /// See CnStakingV3.handleRedelegation for more details.
    /// @param _recipient The address to stake for
    function stakeFor(address _recipient) external payable;

    /// @dev Returns the current rewards to be automatically compounded during the `stakeFor` function.
    /// It is used in CnStakingV3.handleRedelegation to calculate the expected balance after calling `stakeFor`.
    /// If implemented public delegation doesn't support auto-compounding during the `stakeFor`, it should return 0.
    /// See CnStakingV3.handleRedelegation for more details.
    function reward() external view returns (uint256);
}
