// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

contract MockCnStaking {
    uint256 public staking;
    uint256 public unstaking;
    address public publicDelegation;

    function mockSetStaking(uint256 _staking) external {
        staking = _staking;
    }

    function mockSetUnstaking(uint256 _unstaking) external {
        unstaking = _unstaking;
    }

    function mockSetPublicDelegation(address _pd) external {
        publicDelegation = _pd;
    }
}
