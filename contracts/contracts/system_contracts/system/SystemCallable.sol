// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

/// @title SystemCallable
/// @notice Shared base for contracts called via system transactions at epoch boundaries
abstract contract SystemCallable {
    uint256 public constant EPOCH_BLOCK_INTERVAL = 86_400;
    address public constant SYSTEM_SENDER = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    error OnlyEpochBlock();
    error OnlySystemTx();

    modifier onlyEpochBlock() {
        _onlyEpochBlock();
        _;
    }

    function _onlyEpochBlock() private view {
        if (block.number % EPOCH_BLOCK_INTERVAL != 0) revert OnlyEpochBlock();
    }

    modifier onlySystemTx() {
        _onlySystemTx();
        _;
    }

    function _onlySystemTx() private view {
        if (msg.sender != SYSTEM_SENDER) revert OnlySystemTx();
    }

    function _isEpochBlock() internal view returns (bool) {
        return block.number % EPOCH_BLOCK_INTERVAL == 0;
    }

    function _currentEpoch() internal view returns (uint256) {
        return block.number / EPOCH_BLOCK_INTERVAL;
    }
}
