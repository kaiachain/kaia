// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.18;

interface IKIP247 {
    function swapForGas(address token, uint256 amountIn, uint256 minAmountOut, uint256 amountRepay) external;

    function addToken(address token, address factory, address router) external;

    function removeToken(address token) external;

    function claimCommission() external;

    function updateCommissionRate(uint256 _commissionRate) external;

    function dexAddress(address token) external view returns (address dex);

    function getAmountIn(address token, uint amountOut) external view returns (uint amountIn);

    function getSupportedTokens() external view returns (address[] memory);
}
