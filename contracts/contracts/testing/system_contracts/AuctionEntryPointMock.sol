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

pragma solidity ^0.8.18;

contract AuctionEntryPointMock {
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

    address public auctioneer;
    uint256 public count;

    modifier onlyProposer() {
        if (msg.sender != block.coinbase) revert();
        _;
    }

    function setAuctioneer(address _auctioneer) public {
        auctioneer = _auctioneer;
    }

    function call(AuctionTx calldata auctionTx) external onlyProposer {
        // 1. Verify input integrity
        if (!_verifyInputIntegrity(auctionTx)) revert();

        // // 2. Take bid first
        // if (!_checkAndTakeBid(searcher, auctionTx.bid, callGasLimit)) revert();

        // // 3. Execute call and refund execution gas
        // uint256 nonce = _useNonce(searcher);
        // (bool success, ) = auctionTx.to.call{gas: callGasLimit}(auctionTx.data);
        // if (success) {
        //     emit Call(searcher, nonce);
        // } else {
        //     emit CallFailed(searcher, nonce);
        // }

        // // 4. Refund gas to the proposer
        // if (!_payGas(searcher, initialGas)) revert();
        count++;
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

        // /// 3. Check if the auctioneer signature is valid
        // bytes32 digest = MessageHashUtils.toEthSignedMessageHash(auctionTx.searcherSig);
        // (address recoveredSigner, , ) = digest.tryRecover(auctionTx.auctioneerSig);
        // if (recoveredSigner != auctioneer) {
        //     return false;
        // }

        // /// 4. Check if the searcher signature is valid
        // bytes32 structHash = _getAuctionTxHash(auctionTx);
        // // Compute the final digest
        // digest = _hashTypedDataV4(structHash);
        // // Recover the signer from the signature
        // (recoveredSigner, , ) = digest.tryRecover(auctionTx.searcherSig);

        // if (recoveredSigner != auctionTx.sender) {
        //     return false;
        // }

        // return auctionTx.nonce == nonces(auctionTx.sender);
        return true;
    }
}
