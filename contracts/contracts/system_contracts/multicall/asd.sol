


// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.4.24;

contract C {
    function A() public {
        uint256 MAX = 100;
        uint256[] memory nums = new uint256[](MAX);
        for (uint i=0; i<MAX; i++) {
            nums[i] = i;
        }
    }

    function B(uint MAX) public returns (uint256[] memory){
        uint256[] memory nums = new uint256[](MAX);
        for (uint i=0; i<MAX; i++) {
            nums[i] = i;
        }
        return nums;
    }

    function D(uint MAX) public returns (uint256[] memory){
        uint256[] memory nums = new uint256[](MAX);
        // for (uint i=0; i<MAX; i++) {
        //     nums[i] = i;
        // }
        return nums;
    }
}
