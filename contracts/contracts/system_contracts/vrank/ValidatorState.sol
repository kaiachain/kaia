// Copyright 2025 The kaia Authors
// This file is part of the kaia library.
//
// The kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the kaia library. If not, see <http://www.gnu.org/licenses/>.

// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

import {EnumerableSet} from "openzeppelin-contracts-5.0/utils/structs/EnumerableSet.sol";
import {IValidatorState, IAddressBook, ICnStaking} from "./IValidatorState.sol";

contract ValidatorState is IValidatorState {
    using EnumerableSet for EnumerableSet.AddressSet;

    address public owner;
    address private _systemAddress = address(0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE);
    IAddressBook public constant ABOOK = IAddressBook(0x0000000000000000000000000000000000000400);
    uint256 public constant ValPausedTimeout = 8 hours;
    uint256 public constant ValIdleTimeout = 30 days;
    mapping(address => ValidatorState) public validatorStates;
    EnumerableSet.AddressSet private _validatorAddrs;

    modifier onlyOwner() {
        require(msg.sender == owner || msg.sender == _systemAddress, "Error: Caller is not the owner");
        _;
    }

    modifier mustBeValidator(address addr) {
        require(_validatorAddrs.contains(addr), "Error: given address is not a vaildator");
        _;
    }

    modifier mustBeState(address addr, State expectedState) {
        require(validatorStates[addr].state == expectedState, "Error: Not an expected state");
        _;
    }

    modifier mustBeEitherState(address addr, State expectedState1, State expectedState2) {
        require(
            validatorStates[addr].state == expectedState1 || validatorStates[addr].state == expectedState2,
            "Error: Not an expected state"
        );
        _;
    }

    constructor(SystemValidatorUpdateRequest[] memory valStates) {
        owner = msg.sender;
        _setValidatorStates(valStates);
    }

    function _setValidatorStates(SystemValidatorUpdateRequest[] memory valStates) internal  {
        for (uint256 i=0; i<valStates.length; i++) {
            address addr = valStates[i].addr;
            State state = valStates[i].state;
            setState(addr, state);
            _validatorAddrs.add(addr);

            if (state == State.ValInactive) {
                setIdleTimeout(addr);
            }
        }
    }

    function setState(address addr, State state) internal {
        emit StateTranstiion(validatorStates[addr].state, state);
        validatorStates[addr].state = state;
    }

    function setIdleTimeout(address addr) internal {
        if (validatorStates[addr].idleTimeout < block.timestamp) {
            // set idle time if and only if the idle time is elapsed
            validatorStates[addr].idleTimeout = block.timestamp + ValIdleTimeout;
        }
        emit SetIdleTimeout(addr, validatorStates[addr].idleTimeout);
    }

    function setPauesdTimeout(address addr) internal {
        validatorStates[addr].pausedTimeout = block.timestamp + ValPausedTimeout;
        emit SetPausedTimeout(addr, validatorStates[addr].pausedTimeout);
    }

    // system tx (setter)
    function setValidatorStates(SystemValidatorUpdateRequest[] memory valStates) external onlyOwner {
        _setValidatorStates(valStates);
    }

    // public getter
    function getAllValidators() external view returns (ValidatorState[] memory) {
        uint256 valLen = _validatorAddrs.length();
        ValidatorState[] memory valStates = new ValidatorState[](valLen);
        for (uint256 i=0; i<valLen; i++) {
            address addr = _validatorAddrs.at(i);
            valStates[i] = ValidatorState(addr, validatorStates[addr].state, validatorStates[addr].idleTimeout, validatorStates[addr].pausedTimeout);
        }
        return valStates;
    }

    // user tx
    function setValReady() mustBeValidator(msg.sender) mustBeState(msg.sender, State.ValInactive) external {
        address addr = msg.sender;
        (, address stakingContract, ) = ABOOK.getCnInfo(addr);
        uint256 stakingAmount = ICnStaking(stakingContract).staking() - ICnStaking(stakingContract).unstaking();
        require(stakingAmount >= 5000000 ether, "Error: insufficient staking amount");

        setState(addr, State.ValReady);
        setIdleTimeout(addr);
    }

    // user tx
    // TODO-Permissionles: Slot constraint is not considered in this mock contract
    function setValPaused() mustBeValidator(msg.sender) mustBeState(msg.sender, State.ValActive) external {
        address addr = msg.sender;
        setState(addr, State.ValPaused);
        setPauesdTimeout(addr);
    }

    // user tx
    function setValActive() mustBeValidator(msg.sender) mustBeState(msg.sender, State.ValPaused) external {
        address addr = msg.sender;
        setState(addr, State.ValActive);
    }

    // user tx
    // TODO-Permissionles: Slot constraint is not considered in this mock contract
    function setValExiting() mustBeValidator(msg.sender) mustBeEitherState(msg.sender, State.ValActive, State.ValPaused) external {
        address addr = msg.sender;
        setState(addr, State.ValExiting);
    }
    
    // user tx
    // TODO-Permissionles: Slot constraint is not considered in this mock contract
    function setCandReady() mustBeEitherState(msg.sender, State.Unknown, State.CandInactive) external {
        address addr = msg.sender;
        (, address stakingContract, ) = ABOOK.getCnInfo(addr);
        uint256 stakingAmount = ICnStaking(stakingContract).staking() - ICnStaking(stakingContract).unstaking();
        require(stakingAmount >= 5000000 ether, "Error: insufficient staking amount");
        setState(addr, State.CandReady);
    }

    // user tx
    function setCandInactive() mustBeEitherState(msg.sender, State.CandReady, State.ValInactive) external {
        address addr = msg.sender;
        setState(addr, State.CandInactive);
    }

    // mock tx
    function setCandInactiveSuper() external onlyOwner {
        address addr = msg.sender;
        setState(addr, State.CandInactive);
    }

    // mock tx
    function setIdleTimeoutSuper(uint256 duration) external mustBeValidator(msg.sender) onlyOwner {
        validatorStates[msg.sender].idleTimeout = block.timestamp + duration;
    }
}
