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

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract Bech32 is Initializable {
    bytes constant charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l";
    bytes32 constant linkHash = keccak256("link");

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() { _disableInitializers(); }

    /// @dev bech32Polymod calculates the BCH checksum for a given hrp, values and
    // checksum data. Checksum is optional, and if nil a 0 checksum is assumed.
    // Values and checksum (if provided) MUST be encoded as 5 bits per element (base
    // 32), otherwise the results are undefined.
    // For more details on the polymod calculation, please refer to BIP 173.
    /// @param hrp HRP part of bech32 address
    /// @param values value part of bech32 address
    /// @param checksum checksum part of bech32 address
    function bech32Polymod(bytes memory hrp, bytes memory values, bytes memory checksum) internal pure returns (uint) {
        uint32[5] memory GENERATOR = [0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3];
        uint chk = 1;

	    // Account for the high bits of the HRP in the checksum.
        for (uint i=0; i<hrp.length; i++) {
            uint b = chk >> 25;
            uint hiBits = uint8(hrp[i]) >> 5;
            chk = (chk&0x1ffffff)<<5 ^ hiBits;
            for (uint j=0; j<5; j++) {
                if ((b >> uint(j))&1 == 1) {
                    chk ^= GENERATOR[j];
                }
            }
        }

        // Account for the separator (0) between high and low bits of the HRP.
        // x^0 == x, so we eliminate the redundant xor used in the other rounds.
        uint _b = chk >> 25;
        chk = (chk & 0x1ffffff) << 5;
        for (uint i=0; i<5; i++) {
            if ((_b>>uint(i))&1 == 1) {
                chk ^= GENERATOR[i];
            }
        }

        // Account for the low bits of the HRP.
        for (uint i=0; i<hrp.length; i++) {
            uint b = chk >> 25;
            uint loBits = uint8(hrp[i]) & 31;
            chk = (chk&0x1ffffff)<<5 ^ loBits;
            for (uint j=0; j<5; j++) {
                if ((b>>uint(j))&1 == 1) {
                    chk ^= GENERATOR[j];
                }
            }
        }

	    // Account for the values.
        for (uint p = 0; p < values.length; p++) {
            uint top = chk >> 25;
            chk = (chk & 0x1ffffff) << 5 ^ uint8(values[p]);
            for (uint i = 0; i < 5; i++) {
                if ((top >> i) & 1 == 1) {
                    chk ^= GENERATOR[i];
                }
            }
        }

        if (checksum.length == 0) {
            // A nil checksum is used during encoding, so assume all bytes are zero.
            // x^0 == x, so we eliminate the redundant xor used in the other rounds.
            for (uint v=0; v<6; v++) {
                uint b = chk >> 25;
                chk = (chk & 0x1ffffff) << 5;
                for (uint i=0; i<5; i++) {
                    if ((b>>uint(i))&1 == 1) {
                        chk ^= GENERATOR[i];
                    }
                }
            }
        } else {
            // Checksum is provided during decoding, so use it.
            for (uint i=0; i<checksum.length; i++) {
                uint b = chk >> 25;
                chk = (chk&0x1ffffff)<<5 ^ uint8(checksum[i]);
                for (uint j=0; j<5; j++) {
                    if ((b>>uint(j))&1 == 1) {
                        chk ^= GENERATOR[j];
                    }
                }
            }
        }
        return chk;
    }

    /// @dev VerifyChecksum verifies whether the bech32 string specified by the provided
    // hrp and payload data (encoded as 5 bits per element byte slice) are validated
    // by the given checksum.
    //
    // For more details on the checksum verification, please refer to BIP 173.
    /// @param hrp HRP part of bech32 address
    /// @param values value part of bech32 address
    /// @param checksum checksum part of bech32 address
    function verifyChecksum(string memory hrp, bytes memory values, bytes memory checksum) public pure returns (bool) {
        uint polymod = bech32Polymod(bytes(hrp), values, checksum);
	    return polymod == 1;
    }

    /// @dev Normalize converts the uppercase letters to lowercase in string, because
    // Bech32 standard uses only the lowercase for of string for checksum calculation.
    // If conversion occurs during function call, `true` will be returned.
    //
    // Mixed case is NOT allowed.
    /// @param bech bech32 address
    function normalize(bytes memory bech) public pure returns (bytes memory) {
        // Only	ASCII characters between 33 and 126 are allowed.
        bool hasLower;
        bool hasUpper;
        for (uint i=0; i<bech.length; i++) {
            if (uint8(bech[i]) < 33 || uint8(bech[i]) > 126) {
                revert("Not allowed ASCII value");
            }

            // The characters must be either all lowercase or all uppercase. Testing
            // directly with ascii codes is safe here, given the previous test.
            hasLower = hasLower || (uint8(bech[i]) >= 97 && uint8(bech[i]) <= 122);
            hasUpper = hasUpper || (uint8(bech[i]) >= 65 && uint8(bech[i]) <= 90);
            if (hasLower && hasUpper) {
                revert("string not all lowercase or all uppercase");
            }
        }

        // Bech32 standard uses only the lowercase for of strings for checksum
        // calculation.
        if (hasUpper) {
            return bytesToLowerString(bech);
        }
        return bech;
    }

    /// @dev convert string to lower case string
    /// @param strBytes bytes string
    function bytesToLowerString(bytes memory strBytes) public pure returns (bytes memory) {
        for (uint i=0; i<strBytes.length; i++) {
            // Check if the character is uppercase A-Z
            if (strBytes[i] >= 0x41 && strBytes[i] <= 0x5A) {
                // Convert to lowercase
                strBytes[i] = bytes1(uint8(strBytes[i]) + 32);
            }
        }
        return strBytes;
    }

    /// @dev Function to find the last index of a byte in a string
    /// @param strBytes bytes string
    /// @param b target char
    function findLastIndexByte(bytes memory strBytes, bytes1 b) public pure returns (int) {
        // Iterate over the string from the end to the beginning
        for (int i=int(strBytes.length) - 1; i>=0; i--) {
            if (strBytes[uint(i)] == b) {
                return i; // Return the index if the byte is found
            }
        }
        return -1; // Return -1 if the byte is not found
    }

    /// @dev Function to find the first index of a byte in a string
    /// @param b target char
    function findIndexByte(bytes1 b) public pure returns (uint) {
        // Iterate over the bytes array
        for (uint i=0; i<charset.length; i++) {
            if (charset[i] == b) {
                return i; // Return the index if the byte is found
            }
        }
        revert("Not found the corresponding index");
    }

    /// @dev toBytes converts each character in the string 'chars' to the value of the
    // index of the correspoding character in 'charset'.
    /// @param strBytes bytes string
    function toBytes(bytes memory strBytes) internal pure returns (bytes memory) {
        bytes memory decoded = new bytes(strBytes.length);
        for (uint i=0; i<strBytes.length; i++) {
            uint idx = findIndexByte(strBytes[i]);
            decoded[i] = bytes1(uint8(idx));
        }
        return decoded;
    }

    /// @dev Function to copy bytes from 'src' starting at 'from' for 'to' bytes
    /// @param src bytes string
    /// @param from start number
    /// @param to end number
    function copyBytes(bytes memory src, uint from, uint to) public pure returns (bytes memory) {
        require(to > from, "from must be less than to");
        require(to <= src.length, "to must be less than or equal source length");

        bytes memory dst = new bytes(to - from);
        uint cnt = 0;
        for (uint i=from; i<to; i++) {
            dst[cnt++] = src[i];
        }
        return dst;
    }

    /// @dev DecodeUnsafe decodes a bech32 encoded string, returning the human-readable
    // part, the data part (excluding the checksum) and the checksum.  This function
    // does NOT validate against the BIP-173 maximum length allowed for bech32 strings
    // and is meant for use in custom applications (such as lightning network payment
    // requests), NOT on-chain addresses.  This function assumes the given string
    // includes lowercase letters only, so if not, you should call Normalize first.
    //
    // Note that the returned data is 5-bit (base32) encoded and the human-readable
    // part will be lowercase.
    /// @param bech bech32 address
    function decodeUnsafe(bytes memory bech) public pure returns (string memory, bytes memory, bytes memory) {
        // The string is invalid if the last '1' is non-existent, it is the
        // first character of the string (no human-readable part) or one of the
        // last 6 characters of the string (since checksum cannot contain '1').
        int one = findLastIndexByte(bech, '1');
        if (one < 1 || one+7 > int(bech.length)) {
            revert("invalid separator index");
        }

        // The human-readable part is everything before the last '1'.
        bytes memory hrp = copyBytes(bech, 0, uint(one));
        bytes memory data = copyBytes(bech, uint(one+1), bech.length);

        // Each character corresponds to the byte with value of the index in 'charset'.
        bytes memory decoded = toBytes(data);
        // return (string(hrp), decoded[:decoded.length - 6], decoded[decoded.length - 6:]);
        return (
            string(hrp),
            copyBytes(decoded, 0, decoded.length - 6),
            copyBytes(decoded, decoded.length - 6, decoded.length)
        );
    }

    /// @dev extract bech32 address
    /// @param hrp HRP part of bech32 address
    /// @param data value part of bech32 address
    function writeBech32Checksum(string memory hrp, bytes memory data) public pure returns (bytes memory) {
        bytes memory empty;
        uint polymod = bech32Polymod(bytes(hrp), data, empty) ^ 1;
        bytes memory checksum = new bytes(6);
        for (uint i=0; i<6; i++) {
            uint b = (polymod >> uint(5*(5-i))) & 31;

            // This can't fail, given we explicitly cap the previous b byte by the
            // first 31 bits.
            checksum[i] = charset[b];
        }
        return checksum;
    }

    /// @dev DecodeNoLimit decodes a bech32 encoded string, returning the human-readable
    // part and the data part excluding the checksum.  This function does NOT
    // validate against the BIP-173 maximum length allowed for bech32 strings and
    // is meant for use in custom applications (such as lightning network payment
    // requests), NOT on-chain addresses.
    //
    // Note that the returned data is 5-bit (base32) encoded and the human-readable
    // part will be lowercase.
    /// @param _bech bech32 address
    function decodeNoLimit(bytes memory _bech, bool useNormalize) public pure returns (string memory, bytes memory) {
        // The minimum allowed size of a bech32 string is 8 characters, since it
        // needs a non-empty HRP, a separator, and a 6 character checksum.
        if (_bech.length < 8) {
            revert("invalid bech32 string length");
        }

        bytes memory bech;
        if (useNormalize) {
            bech = normalize(_bech);
        } else {
            bech = _bech;
        }
        (string memory hrp, bytes memory values, bytes memory checksum) = decodeUnsafe(bech);

        // Verify if the checksum (stored inside decoded[:]) is valid, given the
        // previously decoded hrp.
        if (!verifyChecksum(hrp, values, checksum)) {
            // Invalid checksum. Calculate what it should have been, so that the
            // error contains this information.

            // Extract the actual checksum in the string.
            bytes memory actual = copyBytes(bech, bech.length - 6, bech.length);

            // Calculate the expected checksum, given the hrp and payload data.
            bytes memory extractedChecksum = writeBech32Checksum(hrp, values);
            require(keccak256(extractedChecksum) == keccak256(actual), "Invalid checksum");
        }
        // We exclude the last 6 bytes, which is the checksum.
        return (hrp, values);
    }

    /// @dev verify FNSA address which starts `link` prefix.
    /// @param addr FNSA address
    function verifyAddrFNSA(string memory addr, bool useNormalize) public pure returns (bool) {
        (string memory hrp,) = decodeNoLimit(bytes(addr), useNormalize);
        if (keccak256(abi.encodePacked(hrp)) == linkHash) {
            return true;
        } else {
            return false;
        }
    }
}
