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
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/rlp"
)

func getProposerAndValidatorsFromBlock(block *types.Block) (common.Address, []common.Address, error) {
	blockNumber := block.NumberU64()
	if blockNumber == 0 {
		return common.Address{}, []common.Address{}, nil
	}
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(block.Header())
	if err != nil {
		return common.Address{}, []common.Address{}, err
	}

	sigHash, err := sigHash(block.Header())
	if err != nil {
		return common.Address{}, []common.Address{}, err
	}
	proposerAddr, err := istanbul.GetSignatureAddress(sigHash.Bytes(), istanbulExtra.Seal)
	if err != nil {
		return common.Address{}, []common.Address{}, err
	}

	return proposerAddr, istanbulExtra.Validators, nil
}

func sigHash(header *types.Header) (hash common.Hash, err error) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	if err := rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false)); err != nil {
		logger.Error("fail to encode", "err", err)
		return common.Hash{}, err
	}
	hasher.Sum(hash[:0])
	return hash, nil
}

func makeBlockGroupOutput(blockchain *blockchain.BlockChain, block *types.Block, cInfo consensus.ConsensusInfo, receipts types.Receipts) map[string]interface{} {
	head := block.Header() // copies the header once
	hash := head.Hash()

	td := blockchain.GetTd(hash, block.NumberU64())
	r, _ := kaiaApi.RpcOutputBlock(block, td, false, false, blockchain.Config())

	// make transactions
	transactions := block.Transactions()
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		rpcTransactions[i] = kaiaApi.RpcOutputReceipt(head, tx, hash, head.Number.Uint64(), uint64(i), receipts[i], blockchain.Config())
	}

	r["committee"] = cInfo.Committee
	r["proposer"] = cInfo.Proposer
	r["round"] = cInfo.Round
	r["originProposer"] = cInfo.OriginProposer
	r["transactions"] = rpcTransactions
	return r
}
