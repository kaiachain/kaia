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

import "@openzeppelin/contracts/utils/structs/EnumerableSet.sol";

library EnumerableSetUint64 {
    using EnumerableSet for EnumerableSet.UintSet;

    function setAdd(EnumerableSet.UintSet storage set, uint64 v) internal returns (bool) {
        return set.add(v);
    }

    function setRemove(EnumerableSet.UintSet storage set, uint64 v) internal returns (bool) {
        return set.remove(v);
    }

    function setContains(EnumerableSet.UintSet storage set, uint64 v) internal view returns (bool) {
        return set.contains(v);
    }

    function setAt(EnumerableSet.UintSet storage set, uint256 index) internal view returns (uint64) {
        require(index < set.length(), "Index out of bounds");
        return uint64(set.at(index));
    }

    function setLength(EnumerableSet.UintSet storage set) internal view returns (uint256) {
        return set.length();
    }

    function setValues(EnumerableSet.UintSet storage set) internal view returns (uint256[] memory) {
        return set.values();
    }

    function getAll(EnumerableSet.UintSet storage set) internal view returns (uint64[] memory) {
        uint256 len = set.length();
        uint64[] memory vs = new uint64[](len);
        for (uint256 i=0; i<len; i++) {
            vs[i] = uint64(set.at(i));
        }
        return vs;
    }

    function getRange(EnumerableSet.UintSet storage set, uint64 range) internal view returns (uint64[] memory) {
        uint256 to = range;
        uint256 sl = set.length();
        if (range > sl) {
            to = sl;
        }
        uint64[] memory vs = new uint64[](to);
        for (uint256 i=0; i<to; i++) {
            vs[i] = uint64(set.at(i));
        }
        return vs;
    }
}
