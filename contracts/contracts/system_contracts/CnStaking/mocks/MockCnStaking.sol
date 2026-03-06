// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

import {ICnStaking} from "../interfaces/ICnStaking.sol";

contract MockCnStaking is ICnStaking {
    uint256 public staking;
    uint256 public unstaking;
    address public nodeId;

    function mockSetNodeId(address _nodeId) external {
        nodeId = _nodeId;
    }

    function mockSetStaking(uint256 _staking) external {
        staking = _staking;
    }

    function mockSetUnstaking(uint256 _unstaking) external {
        unstaking = _unstaking;
    }
}
