// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/backend/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package backend

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	kaiaApi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	istanbulCore "github.com/kaiachain/kaia/consensus/istanbul/core"
	"github.com/kaiachain/kaia/networks/rpc"
)

// API is a user facing RPC API to dump Istanbul state
type API struct {
	chain    consensus.ChainReader
	istanbul *backend
}

// GetValidators retrieves the list of qualified validators with the given block number.
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.chain, number)
	if err != nil {
		return nil, err
	}

	valSet, err := api.istanbul.GetValidatorSet(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return valSet.Qualified().List(), nil
}

// GetDemotedValidators retrieves the list of authorized, but demoted validators with the given block number.
func (api *API) GetDemotedValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.chain, number)
	if err != nil {
		return nil, err
	}

	valSet, err := api.istanbul.GetValidatorSet(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return valSet.Demoted().List(), nil
}

// GetValidatorsAtHash retrieves the list of authorized validators with the given block hash.
func (api *API) GetValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetValidators(&rpcBlockNumber)
}

// GetDemotedValidatorsAtHash retrieves the list of demoted validators with the given block hash.
func (api *API) GetDemotedValidatorsAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}
	rpcBlockNumber := rpc.BlockNumber(header.Number.Uint64())
	return api.GetDemotedValidators(&rpcBlockNumber)
}

// Candidates returns the current candidates the node tries to uphold and vote on.
func (api *API) Candidates() map[common.Address]bool {
	api.istanbul.candidatesLock.RLock()
	defer api.istanbul.candidatesLock.RUnlock()

	proposals := make(map[common.Address]bool)
	for address, auth := range api.istanbul.candidates {
		proposals[address] = auth
	}
	return proposals
}

// Propose injects a new authorization candidate that the validator will attempt to
// push through.
func (api *API) Propose(address common.Address, auth bool) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	api.istanbul.candidates[address] = auth
}

// Discard drops a currently running candidate, stopping the validator from casting
// further votes (either for or against).
func (api *API) Discard(address common.Address) {
	api.istanbul.candidatesLock.Lock()
	defer api.istanbul.candidatesLock.Unlock()

	delete(api.istanbul.candidates, address)
}

// API extended by Kaia developers
type APIExtension struct {
	chain    consensus.ChainReader
	istanbul *backend
}

var (
	errPendingNotAllowed       = errors.New("pending is not allowed")
	errInternalError           = errors.New("internal error")
	errStartNotPositive        = errors.New("start block number should be positive")
	errEndLargetThanLatest     = errors.New("end block number should be smaller than the latest block number")
	errStartLargerThanEnd      = errors.New("start should be smaller than end")
	errRequestedBlocksTooLarge = errors.New("number of requested blocks should be smaller than 50")
	errRangeNil                = errors.New("range values should not be nil")
	errNoBlockNumber           = errors.New("block number is not assigned")
)

// GetCouncil retrieves the list of authorized validators at the specified block.
func (api *APIExtension) GetCouncil(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.chain, number)
	if err != nil {
		return nil, err
	}

	valSet, err := api.istanbul.GetValidatorSet(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return valSet.Council().List(), err
}

func (api *APIExtension) GetCouncilSize(number *rpc.BlockNumber) (int, error) {
	council, err := api.GetCouncil(number)
	if err != nil {
		return -1, err
	}
	return len(council), nil
}

func (api *APIExtension) GetCommittee(number *rpc.BlockNumber) ([]common.Address, error) {
	header, err := headerByRpcNumber(api.chain, number)
	if err != nil {
		return nil, err
	}
	roundState, err := api.istanbul.GetCommitteeState(header.Number.Uint64())
	if err != nil {
		return nil, err
	}
	return roundState.Committee().List(), nil
}

func (api *APIExtension) GetCommitteeSize(number *rpc.BlockNumber) (int, error) {
	committee, err := api.GetCommittee(number)
	if err != nil {
		return -1, err
	}
	return len(committee), nil
}

func (api *APIExtension) makeRPCBlockOutput(b *types.Block,
	cInfo consensus.ConsensusInfo, transactions types.Transactions, receipts types.Receipts,
) map[string]interface{} {
	head := b.Header() // copies the header once
	hash := head.Hash()

	r, err := kaiaApi.RpcOutputBlock(b, false, false, api.chain.Config())
	if err != nil {
		logger.Error("failed to RpcOutputBlock", "err", err)
		return nil
	}

	// make transactions
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		if len(receipts) == len(transactions) {
			rpcTransactions[i] = kaiaApi.RpcOutputReceipt(head, tx, hash, head.Number.Uint64(), uint64(i), receipts[i], api.chain.Config())
		} else {
			// fill the transaction output if receipt is not found
			rpcTransactions[i] = kaiaApi.NewRPCTransaction(head, tx, hash, head.Number.Uint64(), uint64(i), api.chain.Config())
		}
	}

	r["committee"] = cInfo.Committee
	r["committers"] = cInfo.Committers
	r["sigHash"] = cInfo.SigHash
	r["proposer"] = cInfo.Proposer
	r["round"] = cInfo.Round
	r["originProposer"] = cInfo.OriginProposer
	r["transactions"] = rpcTransactions
	return r
}

func RecoverCommittedSeals(extra *types.IstanbulExtra, headerHash common.Hash) ([]common.Address, error) {
	committers := make([]common.Address, len(extra.CommittedSeal))
	for idx, cs := range extra.CommittedSeal {
		committer, err := istanbul.GetSignatureAddress(istanbulCore.PrepareCommittedSeal(headerHash), cs)
		if err != nil {
			return nil, err
		}
		committers[idx] = committer
	}
	return committers, nil
}

// TODO-Kaia: This API functions should be managed with API functions with namespace "kaia"
func (api *APIExtension) GetBlockWithConsensusInfoByNumber(number *rpc.BlockNumber) (map[string]interface{}, error) {
	b, ok := api.chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.BlockChain", "type", reflect.TypeOf(api.chain))
		return nil, errInternalError
	}
	var block *types.Block
	var blockNumber uint64

	if number == nil {
		logger.Trace("block number is not assigned")
		return nil, errNoBlockNumber
	}

	if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get consensus information of the PendingBlock.")
		return nil, errPendingNotAllowed
	}

	if *number == rpc.LatestBlockNumber {
		block = b.CurrentBlock()
		blockNumber = block.NumberU64()
	} else {
		// rpc.EarliestBlockNumber == 0, no need to treat it as a special case.
		blockNumber = uint64(number.Int64())
		block = b.GetBlockByNumber(blockNumber)
	}

	if block == nil {
		logger.Trace("Finding a block by number failed.", "blockNum", blockNumber)
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}
	blockHash := block.Hash()

	cInfo, err := api.istanbul.GetConsensusInfo(block)
	if err != nil {
		logger.Error("Getting the proposer and validators failed.", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := b.GetBlockReceiptsInCache(blockHash)
	if receipts == nil {
		receipts = b.GetReceiptsByBlockHash(blockHash)
	}

	return api.makeRPCBlockOutput(block, cInfo, block.Transactions(), receipts), nil
}

func (api *APIExtension) GetBlockWithConsensusInfoByNumberRange(start *rpc.BlockNumber, end *rpc.BlockNumber) (map[string]interface{}, error) {
	blocks := make(map[string]interface{})

	if start == nil || end == nil {
		logger.Trace("the range values should not be nil.", "start", start, "end", end)
		return nil, errRangeNil
	}

	// check error status.
	s := start.Int64()
	e := end.Int64()
	if s < 0 {
		logger.Trace("start should be positive", "start", s)
		return nil, errStartNotPositive
	}

	eChain := api.chain.CurrentHeader().Number.Int64()
	if e > eChain {
		logger.Trace("end should be smaller than the lastest block number", "end", end, "eChain", eChain)
		return nil, errEndLargetThanLatest
	}

	if s > e {
		logger.Trace("start should be smaller than end", "start", s, "end", e)
		return nil, errStartLargerThanEnd
	}

	if (e - s) > 50 {
		logger.Trace("number of requested blocks should be smaller than 50", "start", s, "end", e)
		return nil, errRequestedBlocksTooLarge
	}

	// gather s~e blocks
	for i := s; i <= e; i++ {
		strIdx := fmt.Sprintf("0x%x", i)

		blockNum := rpc.BlockNumber(i)
		b, err := api.GetBlockWithConsensusInfoByNumber(&blockNum)
		if err != nil {
			logger.Error("error on GetBlockWithConsensusInfoByNumber", "err", err)
			blocks[strIdx] = nil
		} else {
			blocks[strIdx] = b
		}
	}

	return blocks, nil
}

func (api *APIExtension) GetBlockWithConsensusInfoByHash(blockHash common.Hash) (map[string]interface{}, error) {
	b, ok := api.chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.Blockchain, returning...", "type", reflect.TypeOf(api.chain))
		return nil, errInternalError
	}

	block := b.GetBlockByHash(blockHash)
	if block == nil {
		logger.Trace("Finding a block failed.", "blockHash", blockHash)
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", blockHash.String())
	}

	cInfo, err := api.istanbul.GetConsensusInfo(block)
	if err != nil {
		logger.Error("Getting the proposer and validators failed.", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := b.GetBlockReceiptsInCache(blockHash)
	if receipts == nil {
		receipts = b.GetReceiptsByBlockHash(blockHash)
	}

	return api.makeRPCBlockOutput(block, cInfo, block.Transactions(), receipts), nil
}

func (api *APIExtension) GetAllRecordsFromRegistry(name string, number rpc.BlockNumber) ([]interface{}, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.chain.CurrentBlock().Number().Int64())
	}

	if api.chain.Config().IsRandaoForkEnabled(bn) {
		backend := backends.NewBlockchainContractBackend(api.chain, nil, nil)
		records, err := system.ReadAllRecordsFromRegistry(backend, name, bn)
		if err != nil {
			return nil, err
		}

		if len(records) == 0 {
			return nil, fmt.Errorf("%s has not been registered", name)
		}

		recordsList := make([]interface{}, len(records))
		for i, record := range records {
			recordsList[i] = map[string]interface{}{"addr": record.Addr, "activation": record.Activation}
		}
		return recordsList, nil
	} else {
		return nil, errors.New("Randao fork is not enabled")
	}
}

func (api *APIExtension) GetActiveAddressFromRegistry(name string, number rpc.BlockNumber) (common.Address, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.chain.CurrentBlock().Number().Int64())
	}

	if api.chain.Config().IsRandaoForkEnabled(bn) {
		backend := backends.NewBlockchainContractBackend(api.chain, nil, nil)
		addr, err := system.ReadActiveAddressFromRegistry(backend, name, bn)
		if err != nil {
			return common.Address{}, err
		}

		if addr == (common.Address{}) {
			return common.Address{}, errors.New("no active address for " + name)
		}
		return addr, nil
	} else {
		return common.Address{}, errors.New("Randao fork is not enabled")
	}
}

func (api *API) GetTimeout() uint64 {
	return istanbul.DefaultConfig.Timeout
}

// Retrieve the header at requested block number
func headerByRpcNumber(chain consensus.ChainReader, number *rpc.BlockNumber) (*types.Header, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = chain.CurrentHeader()
	} else if *number == rpc.PendingBlockNumber {
		logger.Trace("Cannot get snapshot of the pending block.", "number", number)
		return nil, errPendingNotAllowed
	} else {
		header = chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}
	return header, nil
}
