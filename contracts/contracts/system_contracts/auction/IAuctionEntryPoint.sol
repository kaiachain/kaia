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
import "./IAuctionDepositVault.sol";

interface IAuctionEntryPoint {
    /* ========== STRUCT ========== */

    struct AuctionTx {
        bytes32 targetTxHash;
        uint256 blockNumber;
        address sender;
        address to;
        uint256 nonce;
        uint256 bid;
        uint256 callGasLimit;
        bytes data;
        bytes searcherSig; // digest = hashTypedData(AuctionTx)
        bytes auctioneerSig; // digest = hashEthSignedMessage(searcherSig)
    }

    /* ========== EVENT ========== */

    event ChangeDepositVault(address oldDepositVault, address newDepositVault);

    event ChangeAuctioneer(address oldAuctioneer, address newAuctioneer);

    event ChangeGasParameters(
        uint256 gasPerByteIntrinsic,
        uint256 gasPerByteEip7623,
        uint256 gasContractExecution,
        uint256 gasBufferEstimate,
        uint256 gasBufferUnmeasured
    );

    event Call(address sender, uint256 nonce);

    event CallFailed(address sender, uint256 nonce);

    event UseNonce(address searcher, uint256 nonce);

    /* ========== FUNCTION INTERFACE  ========== */

    function changeDepositVault(address _depositVault) external;

    function changeAuctioneer(address _auctioneer) external;

    function changeGasParameters(
        uint256 _gasPerByteIntrinsic,
        uint256 _gasPerByteEip7623,
        uint256 _gasContractExecution,
        uint256 _gasBufferEstimate,
        uint256 _gasBufferUnmeasured
    ) external;

    function call(AuctionTx calldata auctionTx) external;

    function auctioneer() external view returns (address);

    function gasBufferEstimate() external view returns (uint256);

    function depositVault() external view returns (IAuctionDepositVault);

    function getAuctionTxHash(
        AuctionTx calldata auctionTx
    ) external view returns (bytes32);

    function getNoncesAndDeposits(
        address[] memory searchers
    )
        external
        view
        returns (uint256[] memory nonces_, uint256[] memory deposits_);
}
