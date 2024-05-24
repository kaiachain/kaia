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

import "./ICnStakingV3MultiSig.sol";

abstract contract CnStakingV3MultiSigStorage is ICnStakingV3MultiSig {
    /* ========== CONSTANTS ========== */

    uint256 public constant MAX_ADMIN = 50;

    /* ========== STATE VARIABLES ========== */

    /// @dev Instead of replacing `enum Functions` with `bytes4`, we use a mapping to store the function selector
    /// for the compatibility with the previous CN staking contracts. (see `confirmRequest`, `revokeRequest`)
    /// It will be initialized in the constructor.
    mapping(Functions => bytes4) internal _fnSelMap;

    // Temporary admin to validate the CnStakingV3 contract
    address public contractValidator;
    uint256 public requirement;
    uint256 public lastClearedId;
    uint256 public requestCount;
    mapping(uint256 => Request) internal _requestMap;
}
