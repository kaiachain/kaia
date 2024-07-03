pragma solidity ^0.5.6;

contract StaticCaller {
    address callee;
    constructor(address _callee) public {
        callee = _callee;
    }

    function callHelloWorld() public payable {
        callee.staticcall(abi.encodeWithSignature("helloWorld()"));
    }
}

contract Callee {
    function helloWorld() public view {}
}
