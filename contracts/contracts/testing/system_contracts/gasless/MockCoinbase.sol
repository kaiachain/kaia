pragma solidity ^0.8.24;

contract MockCoinbase {
    receive() external payable {
        revert("Coinbase payment rejected");
    }
}
