// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "forge-std/Script.sol";
import "../contracts/system_contracts/vrank/ValidatorState.sol";

contract DeployScript is Script {
    function run() external {
        vm.startBroadcast();

        ValidatorState.ValidatorState[] memory valStates = new ValidatorState.ValidatorState[](4);

        valStates[0].addr = 0x3c44CdDDb6a900fa2b585dd299e03d12fA4293BC;
        valStates[0].state = ValidatorState.State(7);

        valStates[1].addr = 0x70997970C51812dc3A010C7d01b50e0d17dc79C8;
        valStates[1].state = ValidatorState.State(7);

        valStates[2].addr = 0x90F79bf6EB2c4f870365E785982E1f101E93b906;
        valStates[2].state = ValidatorState.State(7);

        valStates[3].addr = 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266;
        valStates[3].state = ValidatorState.State(7);

        new ValidatorState(valStates);

        vm.stopBroadcast();
    }
}
