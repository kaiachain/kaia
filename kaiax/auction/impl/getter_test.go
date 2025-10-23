// Copyright 2025 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package impl

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/log"
	"github.com/stretchr/testify/require"
)

func genBid(targetTxHash common.Hash) *auction.Bid {
	bid := new(auction.Bid)

	bid.TargetTxHash = targetTxHash
	bid.BlockNumber = 3
	bid.To = common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875701")
	bid.Nonce = uint64(0)
	bid.Bid = new(big.Int).SetBytes(common.Hex2Bytes("8ac7230489e80000"))
	bid.CallGasLimit = uint64(100)
	bid.Data = common.Hex2Bytes("d09de08a")

	return bid
}

func TestGetBidTxGenerator(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	module := prep(t)

	module.Start()
	defer module.Stop()

	// Arbitrary target tx
	tx := types.NewTransaction(0, common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875701"), big.NewInt(0), 1000000, big.NewInt(100), []byte("d09de08a"))

	bid := genBid(tx.Hash())
	txOrGen := module.GetBidTxGenerator(tx, bid)
	require.NotNil(t, txOrGen)

	gasLimit, err := module.bidPool.getBidTxGasLimit(bid)
	require.NoError(t, err)

	// Generate transaction from the generator function
	generatedTx, err := txOrGen.GetTx(0)
	require.NoError(t, err)
	require.NotNil(t, generatedTx)

	// Verify transaction properties
	require.Equal(t, uint16(generatedTx.Type()), uint16(0x7802))
	require.Equal(t, uint64(0), generatedTx.Nonce())
	require.Equal(t, module.bidPool.auctionEntryPoint, *generatedTx.To())
	require.Equal(t, common.Big0, generatedTx.Value())
	require.Equal(t, gasLimit, generatedTx.Gas())
	require.Equal(t, tx.GasFeeCap(), generatedTx.GasFeeCap())
	require.Equal(t, tx.GasTipCap(), generatedTx.GasTipCap())
	require.Equal(t, module.bidPool.ChainConfig.ChainID, generatedTx.ChainId())

	// Verify the transaction is properly signed by the auctioneer
	signer := types.LatestSignerForChainID(module.bidPool.ChainConfig.ChainID)
	sender, err := signer.Sender(generatedTx)
	require.NoError(t, err)
	require.Equal(t, crypto.PubkeyToAddress(module.InitOpts.NodeKey.PublicKey), sender)
}
