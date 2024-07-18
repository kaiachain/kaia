pragma solidity ^0.8.24;

contract Factory {
    function createContract (uint256 num) public {
        new Contract(num);
    }
}

contract Contract {
    constructor(uint256 num) {
        num -= 100;
    }
}