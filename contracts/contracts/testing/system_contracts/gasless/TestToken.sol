// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.24;

// Uncomment this line to use console.log
// import "hardhat/console.sol";

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract TestToken is ERC20 {
    constructor(address initialHolder) ERC20("TestToken", "TT") {
        _mint(initialHolder, 1000000 * 1e18);
    }
}
