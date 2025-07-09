package impl

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"

	kaiaApi "github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/networks/rpc"
)

// APIs returns the RPC APIs this valset module provides.
func (v *ValsetModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewValsetAPI(v),
			Public:    true,
		},
	}
}

type ValsetAPI struct {
	vs *ValsetModule
}

func NewValsetAPI(vs *ValsetModule) *ValsetAPI {
	return &ValsetAPI{vs: vs}
}

// GetCouncil retrieves the list of authorized validators at the specified block.
func (api *ValsetAPI) GetCouncil(number *rpc.BlockNumber) ([]common.Address, error) {
	num, err := api.vs.ResolveRpcNumber(number, true)
	if err != nil {
		return nil, err
	}

	return api.vs.GetCouncil(num)
}

func (api *ValsetAPI) GetCouncilSize(number *rpc.BlockNumber) (int, error) {
	council, err := api.GetCouncil(number)
	if err != nil {
		return -1, err
	}
	return len(council), nil
}

func (api *ValsetAPI) GetCommittee(number *rpc.BlockNumber) ([]common.Address, error) {
	// cannot determine the committee of not-yet finalized block because it depends on the round.
	num, err := api.vs.ResolveRpcNumber(number, false)
	if err != nil {
		return nil, err
	}
	header := api.vs.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, errUnknownBlock
	}
	round := uint64(header.Round())

	return api.vs.GetCommittee(num, round)
}

func (api *ValsetAPI) GetCommitteeSize(number *rpc.BlockNumber) (int, error) {
	committee, err := api.GetCommittee(number)
	if err != nil {
		return -1, err
	}
	return len(committee), nil
}

func (api *ValsetAPI) makeRPCBlockOutput(b *types.Block,
	cInfo consensus.ConsensusInfo, transactions types.Transactions, receipts types.Receipts,
) map[string]interface{} {
	head := b.Header() // copies the header once
	hash := head.Hash()

	r, err := kaiaApi.RpcOutputBlock(b, false, false, api.vs.Chain.Config())
	if err != nil {
		logger.Error("failed to RpcOutputBlock", "err", err)
		return nil
	}

	// make transactions
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		if len(receipts) == len(transactions) {
			rpcTransactions[i] = kaiaApi.RpcOutputReceipt(head, tx, hash, head.Number.Uint64(), uint64(i), receipts[i], api.vs.Chain.Config())
		} else {
			// fill the transaction output if receipt is not found
			rpcTransactions[i] = kaiaApi.NewRPCTransaction(head, tx, hash, head.Number.Uint64(), uint64(i), api.vs.Chain.Config())
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

func (api *ValsetAPI) GetBlockWithConsensusInfoByNumber(number *rpc.BlockNumber) (map[string]interface{}, error) {
	b, ok := api.vs.Chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.BlockChain", "type", reflect.TypeOf(api.vs.Chain))
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

	cInfo, err := api.vs.GetConsensusInfo(block)
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

func (api *ValsetAPI) GetBlockWithConsensusInfoByNumberRange(start *rpc.BlockNumber, end *rpc.BlockNumber) (map[string]interface{}, error) {
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

	eChain := api.vs.CurrentHeader().Number.Int64()
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

func (api *ValsetAPI) GetBlockWithConsensusInfoByHash(blockHash common.Hash) (map[string]interface{}, error) {
	b, ok := api.vs.Chain.(*blockchain.BlockChain)
	if !ok {
		logger.Error("chain is not a type of blockchain.Blockchain, returning...", "type", reflect.TypeOf(api.vs.Chain))
		return nil, errInternalError
	}

	block := b.GetBlockByHash(blockHash)
	if block == nil {
		logger.Trace("Finding a block failed.", "blockHash", blockHash)
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", blockHash.String())
	}

	cInfo, err := api.vs.GetConsensusInfo(block)
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

func (api *ValsetAPI) GetAllRecordsFromRegistry(name string, number rpc.BlockNumber) ([]interface{}, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.vs.Chain.CurrentBlock().Number().Int64())
	}

	if api.vs.Chain.Config().IsRandaoForkEnabled(bn) {
		backend := api.vs.NewBlockchainContractBackend()
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

func (api *ValsetAPI) GetActiveAddressFromRegistry(name string, number rpc.BlockNumber) (common.Address, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.vs.Chain.CurrentBlock().Number().Int64())
	}

	if api.vs.Chain.Config().IsRandaoForkEnabled(bn) {
		backend := api.vs.NewBlockchainContractBackend()
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
