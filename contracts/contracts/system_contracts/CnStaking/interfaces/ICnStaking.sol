// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface ICnStaking {
    function nodeId() external view returns (address);
    function staking() external view returns (uint256);
    function unstaking() external view returns (uint256);
}
