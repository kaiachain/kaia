pragma solidity ^0.8.0;

contract MockInvalidFactory {
    function getPair(address, address) external pure returns (address) {
        revert("No pair exists");
    }
}
