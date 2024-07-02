// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

contract Counter {
    uint256 public number;

    constructor() {
    }

    function increment() public {
        number++;
    }

    function die() public {
        number++;
        require(1 == 0);
    }

    function dieMsg() public {
        number++;
        require(2 == 1, "bad input");
    }
}
