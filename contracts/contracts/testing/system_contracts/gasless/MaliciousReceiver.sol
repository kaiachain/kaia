// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

contract MaliciousReceiver {
    receive() external payable {
        revert("I reject KAIA");
    }
}
