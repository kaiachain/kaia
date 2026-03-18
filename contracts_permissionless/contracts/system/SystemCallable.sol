// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.25;

/// @title SystemCallable
/// @notice Shared base for contracts called via system transactions
abstract contract SystemCallable {
    address public constant SYSTEM_SENDER = 0xffffFFFfFFffffffffffffffFfFFFfffFFFfFFfE;

    error OnlySystemTx();

    modifier onlySystemTx() {
        _onlySystemTx();
        _;
    }

    function _onlySystemTx() private view {
        if (msg.sender != SYSTEM_SENDER) revert OnlySystemTx();
    }
}
