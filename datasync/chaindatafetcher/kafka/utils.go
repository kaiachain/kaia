// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package kafka

import (
	kaiaApi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
)

func makeBlockGroupOutput(blockchain *blockchain.BlockChain, block *types.Block, receipts types.Receipts) (map[string]interface{}, error) {
	head := block.Header() // copies the header once
	hash := head.Hash()

	r, _ := kaiaApi.RpcOutputBlock(block, false, false, blockchain.Config())

	// make transactions
	transactions := block.Transactions()
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		rpcTransactions[i] = kaiaApi.RpcOutputReceipt(head, tx, hash, head.Number.Uint64(), uint64(i), receipts[i], blockchain.Config())
	}

	cInfo, err := kaiaApi.GetConsensusInfo(block, blockchain.Validator().ValsetModule(), blockchain.Sealer())
	if err != nil {
		return nil, err
	}

	r["committee"] = cInfo.Committee
	r["proposer"] = cInfo.Proposer
	r["round"] = cInfo.Round
	r["originProposer"] = cInfo.OriginProposer
	r["transactions"] = rpcTransactions
	return r, nil
}
