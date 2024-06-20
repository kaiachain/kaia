pragma solidity ^0.5.6;

contract Caller {
    address callee;
    constructor(address _callee) public {
        callee = _callee;
    }

    function callHelloWorld() public payable {
        callee.delegatecall(abi.encodeWithSignature("helloWorld()"));
    }
}

contract Callee {
    function helloWorld() public payable {}
}
