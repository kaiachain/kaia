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
pragma solidity 0.8.19;

/// @title KIP-113 BLS public key registry
/// @dev See https://github.com/klaytn/kips/issues/113
interface IKIP113 {
    struct BlsPublicKeyInfo {
        /// @dev uncompressed BLS12-381 G1 public key in EIP-2537 format (128 bytes)
        ///  Encoding: 16 zero bytes || X (48 bytes) || 16 zero bytes || Y (48 bytes)
        bytes publicKey;
        /// @dev proof-of-possession in EIP-2537 G2 format (256 bytes)
        ///  must be a result of PopProve algorithm as per
        ///  draft-irtf-cfrg-bls-signature-05 section 3.3.3.
        ///  with ciphersuite "BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_"
        ///  Encoding: 16 zero bytes || X0 (48 bytes) || 16 zero bytes || X1 (48 bytes)
        ///           || 16 zero bytes || Y0 (48 bytes) || 16 zero bytes || Y1 (48 bytes)
        bytes pop;
    }

    /// @dev Returns all the stored addresses, public keys, and proof-of-possessions at once.
    function getAllBlsInfo() external view returns (address[] memory nodeIdList, BlsPublicKeyInfo[] memory pubkeyList);
}
