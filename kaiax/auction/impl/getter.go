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
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/kaiax/builder"
)

const (
	AuctionTxMaxGasLimit = uint64(12_000_000)
)

func (a *AuctionModule) GetBidTxGenerator(tx *types.Transaction, bid *auction.Bid) *builder.TxOrGen {
	gen := func(nonce uint64) (*types.Transaction, error) {
		var (
			chainId           = a.InitOpts.ChainConfig.ChainID
			signer            = types.LatestSignerForChainID(chainId)
			auctionEntryPoint = a.bidPool.GetAuctionEntryPoint()
			key               = a.InitOpts.NodeKey
		)

		data, err := system.EncodeAuctionCallData(bid)
		if err != nil {
			return nil, err
		}

		tx, err := types.NewTransactionWithMap(types.TxTypeEthereumDynamicFee, map[types.TxValueKeyType]interface{}{
			types.TxValueKeyNonce:      nonce,
			types.TxValueKeyTo:         &auctionEntryPoint,
			types.TxValueKeyAmount:     common.Big0,
			types.TxValueKeyData:       data,
			types.TxValueKeyGasLimit:   AuctionTxMaxGasLimit,
			types.TxValueKeyGasFeeCap:  tx.GasFeeCap(),
			types.TxValueKeyGasTipCap:  tx.GasTipCap(),
			types.TxValueKeyAccessList: types.AccessList{},
			types.TxValueKeyChainID:    chainId,
		})
		if err != nil {
			return nil, err
		}

		err = tx.Sign(signer, key)

		return tx, err
	}

	return builder.NewTxOrGenFromGen(gen, bid.Hash())
}
