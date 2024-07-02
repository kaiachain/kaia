pragma solidity ^0.8.24;

contract Factory {
    event Deploy(address a);
    function deploy(bytes32 salt, uint256 arg) public {
        Contract a = new Contract{salt: salt}(arg);
        emit Deploy(address(a));
    }
}

contract Contract {
    constructor(uint256 arg) {
        arg -= 100; // reverts if arg < 100
    }
}
