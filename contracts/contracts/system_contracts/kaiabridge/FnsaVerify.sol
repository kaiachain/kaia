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

import {Bech32} from "../kaiabridge/Bech32.sol";

library FnsaVerify {
    /* ========== CONSTANTS ========== */

    uint256 internal constant PUBLIC_KEY_LENGTH = 65;
    uint256 internal constant SIGNATURE_LENGTH = 65;
    bytes1 internal constant UNCOMPRESSED_PREFIX = 0x04;

    string internal constant BECH32_PREFIX = "link1";
    string internal constant BECH32_VALOPER_PREFIX = "linkvaloper1";

    /* ========== PUBLIC FUNCTIONS ========== */

    function verify(
        bytes calldata publicKey,
        string memory fnsaAddress,
        bytes32 messageHash,
        bytes calldata signature
    ) internal pure returns (address) {
        _validatePublicKey(publicKey);
        _validateSignature(signature);

        string memory computedFnsaAddr = computeFnsaAddr(publicKey);
        string memory computedValoperAddr = computeValoperAddr(publicKey);
        bytes32 providedHash = keccak256(abi.encodePacked(fnsaAddress));
        require(
            keccak256(abi.encodePacked(computedFnsaAddr)) == providedHash ||
                keccak256(abi.encodePacked(computedValoperAddr)) == providedHash,
            "Invalid fnsa address"
        );

        address computedEthAddr = computeEthAddr(publicKey);
        address recoveredEthAddr = recoverEthAddr(messageHash, signature);
        require(recoveredEthAddr == computedEthAddr, "Invalid signature");

        return computedEthAddr;
    }

    function computeFnsaAddr(bytes calldata publicKey) internal pure returns (string memory) {
        _validatePublicKey(publicKey);

        uint[] memory hrp = _bech32Hrp();
        bytes memory encoded = _encodePublicKey(publicKey, hrp);
        return string(abi.encodePacked(BECH32_PREFIX, encoded));
    }

    function computeValoperAddr(bytes calldata publicKey) internal pure returns (string memory) {
        _validatePublicKey(publicKey);

        uint[] memory hrp = _bech32ValoperHrp();
        bytes memory encoded = _encodePublicKey(publicKey, hrp);
        return string(abi.encodePacked(BECH32_VALOPER_PREFIX, encoded));
    }

    function computeEthAddr(bytes calldata publicKey) internal pure returns (address) {
        _validatePublicKey(publicKey);

        return address(uint160(uint256(keccak256(publicKey[1:]))));
    }

    function recoverEthAddr(bytes32 messageHash, bytes calldata signature) internal pure returns (address) {
        _validateSignature(signature);

        bytes32 r = bytes32(signature[:32]);
        bytes32 s = bytes32(signature[32:64]);
        uint8 v = uint8(signature[64]);
        return ecrecover(messageHash, v, r, s);
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _validatePublicKey(bytes calldata publicKey) private pure {
        require(publicKey.length == PUBLIC_KEY_LENGTH, "Invalid public key length");
        require(publicKey[0] == UNCOMPRESSED_PREFIX, "Invalid public key");
    }

    function _validateSignature(bytes calldata signature) private pure {
        require(signature.length == SIGNATURE_LENGTH, "Invalid signature length");
    }

    function _encodePublicKey(bytes calldata publicKey, uint[] memory hrp) private pure returns (bytes memory) {
        bytes32 pubX = bytes32(publicKey[1:33]);
        bytes32 pubY = bytes32(publicKey[33:]);
        uint8 pubPrefix = (uint256(pubY) % 2 == 0) ? 0x02 : 0x03;
        bytes memory compressedPub = abi.encodePacked(pubPrefix, pubX);

        bytes32 sha = sha256(abi.encodePacked(compressedPub));
        bytes20 digest = ripemd160(abi.encodePacked(sha));

        uint[] memory words8 = new uint[](20);
        for (uint i = 0; i < 20; i++) {
            words8[i] = uint8(digest[i]);
        }
        uint[] memory words5 = Bech32.convert(words8, 8, 5);
        return Bech32.encode(hrp, words5);
    }

    function _bech32Hrp() private pure returns (uint[] memory hrp) {
        hrp = new uint[](4);
        hrp[0] = 108;
        hrp[1] = 105;
        hrp[2] = 110;
        hrp[3] = 107;
    }

    function _bech32ValoperHrp() private pure returns (uint[] memory hrp) {
        hrp = new uint[](11);
        hrp[0] = 108;
        hrp[1] = 105;
        hrp[2] = 110;
        hrp[3] = 107;
        hrp[4] = 118;
        hrp[5] = 97;
        hrp[6] = 108;
        hrp[7] = 111;
        hrp[8] = 112;
        hrp[9] = 101;
        hrp[10] = 114;
    }
}
