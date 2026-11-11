pragma solidity ^0.8.0;

contract TestReceiver {
    uint256 public count;

    receive() external payable {
        count++;
    }

    function increment() public {
        count++;
    }

    function makeRevert() public pure {
        revert();
    }
}
