// Copyright 2025 The kaia Authors
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
pragma solidity ^0.8.24;

import {FnsaVerify} from "../../system_contracts/kaiabridge/FnsaVerify.sol";


contract FnsaVerifyHarness {
    function computeFnsaAddr(bytes calldata publicKey) external pure returns (string memory) {
        return FnsaVerify.computeFnsaAddr(publicKey);
    }

    function computeValoperAddr(bytes calldata publicKey) external pure returns (string memory) {
        return FnsaVerify.computeValoperAddr(publicKey);
    }

    function computeEthAddr(bytes calldata publicKey) external pure returns (address) {
        return FnsaVerify.computeEthAddr(publicKey);
    }

    function recoverEthAddr(bytes32 messageHash, bytes calldata signature) external pure returns (address) {
        return FnsaVerify.recoverEthAddr(messageHash, signature);
    }

    function verify(
        bytes calldata publicKey,
        string calldata fnsaAddress,
        bytes32 messageHash,
        bytes calldata signature
    ) external pure returns (address) {
        return FnsaVerify.verify(publicKey, fnsaAddress, messageHash, signature);
    }
}
