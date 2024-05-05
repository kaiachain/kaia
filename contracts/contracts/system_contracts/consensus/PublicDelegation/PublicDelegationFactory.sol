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

import "./PublicDelegation.sol";
import "./IPublicDelegationFactory.sol";

// Simple factory contract to deploy PublicDelegation contracts.
contract PublicDelegationFactory is IPublicDelegationFactory {
    /* ========== CONSTANTS ========== */

    uint256 public constant VERSION = 1;

    string public constant CONTRACT_TYPE = "PublicDelegationFactory";

    /* ========== DEPLOYMENT ========== */

    function deployPublicDelegation(
        IPublicDelegation.PDConstructorArgs memory _pdArgs
    ) external override returns (IPublicDelegation publicDelegation) {
        publicDelegation = new PublicDelegation(msg.sender, _pdArgs);
    }
}
