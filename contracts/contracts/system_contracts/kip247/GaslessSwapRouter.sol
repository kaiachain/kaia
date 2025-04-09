// SPDX-License-Identifier: MIT
pragma solidity 0.8.24;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@uniswap/v2-periphery/contracts/interfaces/IUniswapV2Router02.sol";
import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Factory.sol";
import "@uniswap/v2-core/contracts/interfaces/IUniswapV2Pair.sol";
import "./IKIP247.sol";
import "./IWKAIA.sol";

/**
 * @title GaslessSwapRouter
 * @dev Implements KIP-247 gasless transaction functionality
 * This contract allows users to swap ERC20 tokens for KAIA without possessing any KAIA.
 * 
 * LIMITATIONS:
 * - This contract does not support Fee-on-transfer (FoT) tokens
 * - Using FoT tokens may result in transaction failures or incorrect amounts
 */
contract GaslessSwapRouter is IKIP247, Ownable {
    using SafeERC20 for IERC20;
    IWKAIA public immutable WKAIA;

    mapping(address => DEXInfo) private _dexInfos;
    address[] private _supportedTokens;

    uint256 public commissionRate; // 10000 = 100%

    event SwappedForGas(
        address indexed proposer,
        uint256 amountRepaid,
        address indexed user,
        uint256 finalUserAmount,
        uint256 commission
    );
    event TokenAdded(address indexed token, address indexed factory, address router);
    event TokenRemoved(address indexed token);
    event CommissionRateUpdated(uint256 oldRate, uint256 newRate);
    event CommissionClaimed(uint256 amount);

    constructor(address _wkaia) {
        require(_wkaia != address(0), "Zero address is not allowed");
        _transferOwnership(msg.sender);
        WKAIA = IWKAIA(_wkaia);
        commissionRate = 0;
    }

    /**
    * @notice Adds a token to the list of supported tokens
    * @dev IMPORTANT: This contract does not support Fee-on-transfer (FoT) tokens.
    * Such tokens will not function correctly with this contract and should not be added.
    */
    function addToken(address token, address factory, address router) external override onlyOwner {
        require(token != address(0), "Invalid token address");
        require(factory != address(0), "Invalid factory address");
        require(router != address(0), "Invalid router address");
        require(_dexInfos[token].factory == address(0), "TokenAlreadySupported");

        address pair;
        bool success;
        try IUniswapV2Factory(factory).getPair(token, address(WKAIA)) returns (address pairAddress) {
            pair = pairAddress;
            success = true;
        } catch {
            success = false;
        }

        require(success, "InvalidDEXAddress");
        require(pair != address(0), "PairDoesNotExist");

        (uint112 reserve0, uint112 reserve1, ) = IUniswapV2Pair(pair).getReserves();
        require(reserve0 > 0 && reserve1 > 0, "NoLiquidity");

        _dexInfos[token] = DEXInfo({factory: factory, router: router});

        _supportedTokens.push(token);

        emit TokenAdded(token, factory, router);
    }

    function removeToken(address token) external override onlyOwner {
        require(_dexInfos[token].factory != address(0), "TokenNotSupported");

        delete _dexInfos[token];

        for (uint i = 0; i < _supportedTokens.length; i++) {
            if (_supportedTokens[i] == token) {
                _supportedTokens[i] = _supportedTokens[_supportedTokens.length - 1];
                _supportedTokens.pop();
                break;
            }
        }

        emit TokenRemoved(token);
    }

    function dexAddress(address token) external view override returns (address) {
        require(_dexInfos[token].factory != address(0), "TokenNotSupported");
        return _dexInfos[token].factory;
    }

    function getDEXInfo(address token) external view returns (address factory, address router) {
        require(isTokenSupported(token), "TokenNotSupported");

        DEXInfo memory info = _dexInfos[token];
        return (info.factory, info.router);
    }

    function claimCommission() external override onlyOwner {
        uint256 amount = address(this).balance;
        require(amount > 0, "NoCommissionToWithdraw");

        (bool success, ) = owner().call{value: amount}("");
        require(success, "CommissionClaimFailed");

        emit CommissionClaimed(amount);
    }

    function updateCommissionRate(uint256 _commissionRate) external override onlyOwner {
        require(_commissionRate <= 10000, "InvalidCommissionRate");

        uint256 oldRate = commissionRate;
        commissionRate = _commissionRate;

        emit CommissionRateUpdated(oldRate, _commissionRate);
    }

    function swapForGas(address token, uint256 amountIn, uint256 minAmountOut, uint256 amountRepay) external override {
        // R2: Token is whitelisted and has a corresponding DEX info
        require(isTokenSupported(token), "TokenNotSupported");

        DEXInfo memory dexInfo = _dexInfos[token];

        // R1: Sender has enough tokens
        require(IERC20(token).balanceOf(msg.sender) >= amountIn, "Insufficient token balance");

        // Part of R3: Check minAmountOut >= amountRepay before swap
        require(minAmountOut >= amountRepay, "InsufficientSwapOutput");

        // Get tokens from user
        IERC20(token).safeTransferFrom(msg.sender, address(this), amountIn);

        // Approve tokens for router
        IERC20(token).safeApprove(dexInfo.router, amountIn);

        // Set up path for swap
        address[] memory path = new address[](2);
        path[0] = token;
        path[1] = address(WKAIA);

        // Execute swap using token-specific router
        IUniswapV2Router02 router = IUniswapV2Router02(dexInfo.router);
        uint256[] memory amounts = router.swapExactTokensForETH(
            amountIn,
            minAmountOut,
            path,
            address(this),
            block.timestamp + 300
        );

        uint256 receivedAmount = amounts[1];

        // Pay the block proposer
        (bool success, ) = block.coinbase.call{value: amountRepay}("");
        require(success, "Failed to send KAIA to proposer");

        // Calculate commission
        uint256 userAmount = receivedAmount - amountRepay;
        uint256 commission = (userAmount * commissionRate) / 10000;
        uint256 finalUserAmount = userAmount - commission;

        // Send remaining KAIA to user
        (bool userTransferSuccess, ) = msg.sender.call{value: finalUserAmount}("");
        require(userTransferSuccess, "FailedToSendKAIA");

        emit SwappedForGas(block.coinbase, amountRepay, msg.sender, finalUserAmount, commission);
    }

    function getAmountIn(address token, uint256 amountOut) external view override returns (uint256 amountIn) {
        require(isTokenSupported(token), "TokenNotSupported");

        DEXInfo memory dexInfo = _dexInfos[token];

        address[] memory path = new address[](2);
        path[0] = token;
        path[1] = address(WKAIA);

        // Use token-specific router
        IUniswapV2Router02 router = IUniswapV2Router02(dexInfo.router);
        uint256[] memory amounts = router.getAmountsIn(amountOut, path);
        return amounts[0];
    }

    function isTokenSupported(address token) public view returns (bool) {
        return _dexInfos[token].factory != address(0);
    }

    function getSupportedTokens() external view returns (address[] memory) {
        return _supportedTokens;
    }

    receive() external payable {}
}
