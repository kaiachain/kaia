// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract MaliciousToken is ERC20 {
    constructor(address initialHolder) ERC20("MaliciousToken", "MTKN") {
        _mint(initialHolder, 1000000 * 10 ** 18);
    }

    function approve(address spender, uint256 amount) public virtual override returns (bool) {
        return false; // Always fails
    }
}
