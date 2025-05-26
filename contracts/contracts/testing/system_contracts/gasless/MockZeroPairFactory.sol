// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract MockZeroPairFactory {
    function getPair(address _tokenA, address _tokenB) external pure returns (address) {
        // Always return zero address to simulate a non-existent pair
        return address(0);
    }
}
