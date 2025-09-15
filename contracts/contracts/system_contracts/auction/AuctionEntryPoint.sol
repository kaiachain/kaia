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

import "openzeppelin-contracts-5.0/access/Ownable.sol";
import "openzeppelin-contracts-5.0/utils/cryptography/EIP712.sol";
import "openzeppelin-contracts-5.0/utils/cryptography/ECDSA.sol";
import "openzeppelin-contracts-5.0/utils/cryptography/MessageHashUtils.sol";
import "openzeppelin-contracts-5.0/utils/Nonces.sol";
import "./AuctionError.sol";
import "./IAuctionEntryPoint.sol";
import "./IAuctionDepositVault.sol";

contract AuctionEntryPoint is
    IAuctionEntryPoint,
    AuctionError,
    Nonces,
    EIP712,
    Ownable
{
    using ECDSA for bytes32;

    /* ========== CONSTANTS ========== */

    bytes32 private constant _AUCTIONTX_TYPEHASH =
        keccak256(
            "AuctionTx(bytes32 targetTxHash,uint256 blockNumber,address sender,address to,uint256 nonce,uint256 bid,uint256 callGasLimit,bytes data)"
        );

    uint256 public constant GAS_PER_BYTE_INTRINSIC = 16; // Base gas cost per byte of msg.data (approximated from 16 gas per non-zero byte + 4 gas per zero byte)
    uint256 public constant GAS_PER_BYTE_EIP_7623 = 40; // Minimum gas cost per byte of msg.data under EIP-7623 (approximated from 40 gas per non-zero byte + 10 gas per zero byte)
    uint256 public constant GAS_BUFFER_ESTIMATE = 180_000; // Buffer for gas calculation except for the main call
    uint256 public constant GAS_BUFFER_UNMEASURED = 20_000; // Buffer for gas calculation that `gasleft()` can't capture after the main call
    uint256 public constant GAS_CONTRACT_EXECUTION = 21_000; // Default transaction gas cost

    string public constant AUCTION_NAME = "KAIA_AUCTION";
    string public constant AUCTION_VERSION = "0.0.1";

    /* ========== STATE VARIABLES ========== */

    IAuctionDepositVault public depositVault;
    address public auctioneer;

    /* ========== MODIFIER ========== */

    modifier depositVaultNotEmpty() {
        if (address(depositVault) == address(0)) revert EmptyDepositVault();
        _;
    }

    modifier onlyProposer() {
        if (msg.sender != block.coinbase) revert OnlyProposer();
        _;
    }

    /* ========== CONSTRUCTOR ========== */

    constructor(
        address initialOwner,
        address _depositVault,
        address _auctioneer
    )
        EIP712(AUCTION_NAME, AUCTION_VERSION)
        Ownable(initialOwner)
        notNull(_depositVault)
        notNull(_auctioneer)
    {
        depositVault = IAuctionDepositVault(_depositVault);
        auctioneer = _auctioneer;
    }

    /* ========== ENTRYPOINT IMPLEMENTATION ========== */

    /// @dev Call the entrypoint
    /// @notice This function is only callable by the proposer, with the bundling mechanism.
    /// @notice This transaction will be discarded from a block at all if reverted, so the proposer won't pay for the gas.
    /// @param auctionTx The auction transaction
    function call(
        AuctionTx calldata auctionTx
    ) external override onlyProposer depositVaultNotEmpty {
        uint256 initialGas = gasleft();
        address searcher = auctionTx.sender;
        uint256 callGasLimit = auctionTx.callGasLimit;

        // 1. Verify input integrity
        if (!_verifyInputIntegrity(auctionTx)) revert();

        // 2. Take bid first
        if (!_checkAndTakeBid(searcher, auctionTx.bid, callGasLimit)) revert();

        // 3. Execute call and refund execution gas
        uint256 nonce = _useNonce(searcher);
        (bool success, ) = auctionTx.to.call{gas: callGasLimit}(auctionTx.data);
        if (success) {
            emit Call(searcher, nonce);
        } else {
            emit CallFailed(searcher, nonce);
        }

        // 4. Refund gas to the proposer
        if (!_payGas(searcher, initialGas)) revert();
    }

    /* ========== ENTRYPOINT MANAGEMENT ========== */

    /// @dev Change the deposit vault
    /// @param _depositVault The new deposit vault
    function changeDepositVault(
        address _depositVault
    ) external override onlyOwner notNull(_depositVault) {
        emit ChangeDepositVault(address(depositVault), _depositVault);
        depositVault = IAuctionDepositVault(_depositVault);
    }

    /// @dev Change the auctioneer
    /// @param _auctioneer The new auctioneer
    function changeAuctioneer(
        address _auctioneer
    ) external override onlyOwner notNull(_auctioneer) {
        emit ChangeAuctioneer(auctioneer, _auctioneer);
        auctioneer = _auctioneer;
    }

    /* ========== INTERNAL FUNCTIONS ========== */

    function _useNonce(address searcher) internal override returns (uint256) {
        uint256 nonce = super._useNonce(searcher);
        // Emit the next nonce for the searcher
        emit UseNonce(searcher, nonce + 1);
        return nonce;
    }

    function _checkAndTakeBid(
        address searcher,
        uint256 bidAmount,
        uint256 callGasLimit
    ) internal returns (bool) {
        uint256 expectedGas = _getMaximumGas(
            callGasLimit + GAS_BUFFER_ESTIMATE
        );
        uint256 expectedSpent = bidAmount + expectedGas * tx.gasprice;

        if (expectedSpent > depositVault.depositBalances(searcher)) {
            return false;
        }

        return depositVault.takeBid(searcher, bidAmount);
    }

    function _payGas(
        address searcher,
        uint256 initialGas
    ) internal returns (bool) {
        uint256 _gasUsed = _getMaximumGas(
            initialGas - gasleft() + GAS_BUFFER_UNMEASURED
        );
        /// @dev The tx.gasprice will be multiplied by the gasUsed in the depositVault
        return depositVault.takeGas(searcher, _gasUsed);
    }

    function _getMaximumGas(
        uint256 executionGas
    ) internal pure returns (uint256) {
        uint256 legacyGas = executionGas + _defaultGas(GAS_PER_BYTE_INTRINSIC);
        uint256 floorDataGas = _defaultGas(GAS_PER_BYTE_EIP_7623);

        return legacyGas > floorDataGas ? legacyGas : floorDataGas;
    }

    function _verifyInputIntegrity(
        AuctionTx calldata auctionTx
    ) internal view returns (bool) {
        /// 1. Check if the block number is correct
        if (auctionTx.blockNumber != block.number) {
            return false;
        }

        /// 2. Check if the bid is greater than 0
        if (auctionTx.bid <= 0) {
            return false;
        }

        /// 3. Check if the auctioneer signature is valid
        bytes32 digest = MessageHashUtils.toEthSignedMessageHash(
            auctionTx.searcherSig
        );
        (address recoveredSigner, , ) = digest.tryRecover(
            auctionTx.auctioneerSig
        );
        if (recoveredSigner != auctioneer) {
            return false;
        }

        /// 4. Check if the searcher signature is valid
        bytes32 structHash = _getAuctionTxHash(auctionTx);
        // Compute the final digest
        digest = _hashTypedDataV4(structHash);
        // Recover the signer from the signature
        (recoveredSigner, , ) = digest.tryRecover(auctionTx.searcherSig);

        if (recoveredSigner != auctionTx.sender) {
            return false;
        }

        return auctionTx.nonce == nonces(auctionTx.sender);
    }

    function _defaultGas(uint256 gasPerByte) internal pure returns (uint256) {
        return msg.data.length * gasPerByte + GAS_CONTRACT_EXECUTION;
    }

    function _getAuctionTxHash(
        AuctionTx calldata auctionTx
    ) internal pure returns (bytes32 auctionTxHash) {
        auctionTxHash = keccak256(
            abi.encode(
                _AUCTIONTX_TYPEHASH,
                auctionTx.targetTxHash,
                auctionTx.blockNumber,
                auctionTx.sender,
                auctionTx.to,
                auctionTx.nonce,
                auctionTx.bid,
                auctionTx.callGasLimit,
                keccak256(auctionTx.data)
            )
        );
    }

    /* ========== GETTERS ========== */

    function getAuctionTxHash(
        AuctionTx calldata auctionTx
    ) external pure override returns (bytes32) {
        return _getAuctionTxHash(auctionTx);
    }

    function getNoncesAndDeposits(
        address[] memory searchers
    )
        external
        view
        override
        returns (uint256[] memory nonces_, uint256[] memory deposits_)
    {
        nonces_ = new uint256[](searchers.length);
        deposits_ = new uint256[](searchers.length);

        for (uint256 i = 0; i < searchers.length; i++) {
            address searcher = searchers[i];
            nonces_[i] = nonces(searcher);
            deposits_[i] = depositVault.depositBalances(searcher);
        }
    }
}
