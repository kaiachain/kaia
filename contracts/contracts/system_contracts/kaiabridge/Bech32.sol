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
pragma solidity 0.8.24;

/// @notice Bech32 helper functions ported from
library Bech32 {
    bytes internal constant CHARSET = "qpzry9x8gf2tvdw0s3jn54khce6mua7l";
    bytes32 internal constant LINK_HASH = keccak256("link");

    //-------------------------------------------------------------------------
    // Core checksum helpers
    //-------------------------------------------------------------------------

    function _bech32Polymod(bytes memory hrp, bytes memory values, bytes memory checksum) internal pure returns (uint) {
        uint32[5] memory GENERATOR = [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3];
        uint chk = 1;

        // high bits of HRP
        for (uint i = 0; i < hrp.length; i++) {
            uint b = chk >> 25;
            uint hiBits = uint8(hrp[i]) >> 5;
            chk = ((chk & 0x1ffffff) << 5) ^ hiBits;
            for (uint j = 0; j < 5; j++) {
                if ((b >> j) & 1 == 1) {
                    chk ^= GENERATOR[j];
                }
            }
        }

        // separator 0
        uint _b = chk >> 25;
        chk = (chk & 0x1ffffff) << 5;
        for (uint i = 0; i < 5; i++) {
            if ((_b >> i) & 1 == 1) {
                chk ^= GENERATOR[i];
            }
        }

        // low bits of HRP
        for (uint i = 0; i < hrp.length; i++) {
            uint b = chk >> 25;
            uint loBits = uint8(hrp[i]) & 31;
            chk = ((chk & 0x1ffffff) << 5) ^ loBits;
            for (uint j = 0; j < 5; j++) {
                if ((b >> j) & 1 == 1) {
                    chk ^= GENERATOR[j];
                }
            }
        }

        // values
        for (uint p = 0; p < values.length; p++) {
            uint top = chk >> 25;
            chk = ((chk & 0x1ffffff) << 5) ^ uint8(values[p]);
            for (uint i = 0; i < 5; i++) {
                if ((top >> i) & 1 == 1) {
                    chk ^= GENERATOR[i];
                }
            }
        }

        if (checksum.length == 0) {
            // assume zero checksum during encoding
            for (uint v = 0; v < 6; v++) {
                uint b = chk >> 25;
                chk = (chk & 0x1ffffff) << 5;
                for (uint i = 0; i < 5; i++) {
                    if ((b >> i) & 1 == 1) {
                        chk ^= GENERATOR[i];
                    }
                }
            }
        } else {
            for (uint i = 0; i < checksum.length; i++) {
                uint b = chk >> 25;
                chk = ((chk & 0x1ffffff) << 5) ^ uint8(checksum[i]);
                for (uint j = 0; j < 5; j++) {
                    if ((b >> j) & 1 == 1) {
                        chk ^= GENERATOR[j];
                    }
                }
            }
        }

        return chk;
    }

    function _writeBech32Checksum(string memory hrp, bytes memory data) internal pure returns (bytes memory) {
        bytes memory empty;
        uint polymod = _bech32Polymod(bytes(hrp), data, empty) ^ 1;
        bytes memory checksum = new bytes(6);
        for (uint i = 0; i < 6; i++) {
            uint b = (polymod >> (5 * (5 - i))) & 31;
            checksum[i] = CHARSET[b];
        }
        return checksum;
    }

    //-------------------------------------------------------------------------
    // Case handling utilities
    //-------------------------------------------------------------------------

    function _bytesToLowerString(bytes memory strBytes) internal pure returns (bytes memory) {
        for (uint i = 0; i < strBytes.length; i++) {
            if (strBytes[i] >= 0x41 && strBytes[i] <= 0x5A) {
                strBytes[i] = bytes1(uint8(strBytes[i]) + 32);
            }
        }
        return strBytes;
    }

    function normalize(bytes memory bech) internal pure returns (bytes memory) {
        bool hasLower;
        bool hasUpper;
        for (uint i = 0; i < bech.length; i++) {
            uint8 c = uint8(bech[i]);
            if (c < 33 || c > 126) revert("Not allowed ASCII value");
            if (c >= 97 && c <= 122) hasLower = true;
            if (c >= 65 && c <= 90) hasUpper = true;
            if (hasLower && hasUpper) revert("string not all lowercase or all uppercase");
        }
        return hasUpper ? _bytesToLowerString(bech) : bech;
    }

    function _findLastIndexByte(bytes memory strBytes, bytes1 b) internal pure returns (int) {
        for (int i = int(strBytes.length) - 1; i >= 0; i--) {
            if (strBytes[uint(i)] == b) return i;
        }
        return -1;
    }

    function _findIndexByte(bytes1 b) internal pure returns (uint) {
        for (uint i = 0; i < CHARSET.length; i++) {
            if (CHARSET[i] == b) return i;
        }
        revert("Not found the corresponding index");
    }

    function _toBase32Bytes(bytes memory strBytes) internal pure returns (bytes memory) {
        bytes memory decoded = new bytes(strBytes.length);
        for (uint i = 0; i < strBytes.length; i++) {
            decoded[i] = bytes1(uint8(_findIndexByte(strBytes[i])));
        }
        return decoded;
    }

    function _copyBytes(bytes memory src, uint from, uint to) internal pure returns (bytes memory) {
        require(to > from, "from must be less than to");
        require(to <= src.length, "to must be less than or equal source length");
        bytes memory dst = new bytes(to - from);
        for (uint i = from; i < to; i++) {
            dst[i - from] = src[i];
        }
        return dst;
    }

    //-------------------------------------------------------------------------
    // Decoding helpers
    //-------------------------------------------------------------------------

    function decodeUnsafe(bytes memory bech) internal pure returns (string memory, bytes memory, bytes memory) {
        int one = _findLastIndexByte(bech, "1");
        if (one < 1 || one + 7 > int(bech.length)) revert("invalid separator index");

        bytes memory hrp = _copyBytes(bech, 0, uint(one));
        bytes memory data = _copyBytes(bech, uint(one + 1), bech.length);
        bytes memory decoded = _toBase32Bytes(data);

        return (
            string(hrp),
            _copyBytes(decoded, 0, decoded.length - 6),
            _copyBytes(decoded, decoded.length - 6, decoded.length)
        );
    }

    function verifyChecksum(
        string memory hrp,
        bytes memory values,
        bytes memory checksum
    ) internal pure returns (bool) {
        return _bech32Polymod(bytes(hrp), values, checksum) == 1;
    }

    function decodeNoLimit(bytes memory bech, bool useNormalize) internal pure returns (string memory, bytes memory) {
        if (bech.length < 8) revert("invalid bech32 string length");
        bytes memory normalized = useNormalize ? normalize(bech) : bech;
        (string memory hrp, bytes memory values, bytes memory checksum) = decodeUnsafe(normalized);
        if (!verifyChecksum(hrp, values, checksum)) {
            bytes memory actual = _copyBytes(normalized, normalized.length - 6, normalized.length);
            bytes memory expected = _writeBech32Checksum(hrp, values);
            require(keccak256(expected) == keccak256(actual), "Invalid checksum");
        }
        return (hrp, values);
    }

    function verifyAddrFNSA(string memory addr, bool useNormalize) internal pure returns (bool) {
        (string memory hrp, ) = decodeNoLimit(bytes(addr), useNormalize);
        return keccak256(abi.encodePacked(hrp)) == LINK_HASH;
    }

    //-------------------------------------------------------------------------
    // Encoding helpers
    //-------------------------------------------------------------------------

    function convert(uint[] memory data, uint inBits, uint outBits) internal pure returns (uint[] memory) {
        uint value;
        uint bits;
        uint maxV = (1 << outBits) - 1;
        uint[] memory tmp = new uint[]((data.length * inBits + outBits - 1) / outBits);
        uint j;
        for (uint i = 0; i < data.length; i++) {
            value = (value << inBits) | data[i];
            bits += inBits;
            while (bits >= outBits) {
                bits -= outBits;
                tmp[j] = (value >> bits) & maxV;
                j++;
            }
        }
        uint[] memory ret = new uint[](j);
        for (uint i = 0; i < j; i++) {
            ret[i] = tmp[i];
        }
        return ret;
    }

    function encode(uint[] memory hrp, uint[] memory data) internal pure returns (bytes memory) {
        bytes memory hrpBytes = new bytes(hrp.length);
        for (uint i = 0; i < hrp.length; i++) {
            hrpBytes[i] = bytes1(uint8(hrp[i]));
        }

        bytes memory dataBytes = new bytes(data.length);
        for (uint i = 0; i < data.length; i++) {
            dataBytes[i] = bytes1(uint8(data[i]));
        }

        bytes memory checksum = _writeBech32Checksum(string(hrpBytes), dataBytes);
        bytes memory ret = new bytes(data.length + checksum.length);
        for (uint i = 0; i < data.length; i++) {
            ret[i] = CHARSET[data[i]];
        }
        for (uint i = 0; i < checksum.length; i++) {
            ret[data.length + i] = checksum[i];
        }
        return ret;
    }
}
