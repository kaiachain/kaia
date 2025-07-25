// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/account"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/contracts/contracts/system_contracts/misc"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

var logger = log.NewModuleLogger(log.API)

// Error variables for consensus info APIs
var (
	errInternalError           = errors.New("internal error")
	errNoBlockNumber           = errors.New("block number is not assigned")
	errPendingNotAllowed       = errors.New("pending is not allowed")
	errRangeNil                = errors.New("range values should not be nil")
	errStartNotPositive        = errors.New("start block number should be positive")
	errEndLargetThanLatest     = errors.New("end block number should be smaller than the latest block number")
	errStartLargerThanEnd      = errors.New("start should be smaller than end")
	errRequestedBlocksTooLarge = errors.New("number of requested blocks should be smaller than 50")
)

// CreditOutput represents the output structure for GetCypressCredit
type CreditOutput struct {
	Photo       string `json:"photo"`
	Names       string `json:"names"`
	EndingPhoto string `json:"endingPhoto"`
	EndingNames string `json:"endingNames"`
}

// BlockChainAPI provides an API to access the Kaia blockchain.
// It offers only methods that operate on public data that is freely available to anyone.
type KaiaBlockChainAPI struct {
	b Backend
}

// NewKaiaBlockChainAPI creates a new Kaia blockchain API.
func NewKaiaBlockChainAPI(b Backend) *KaiaBlockChainAPI {
	return &KaiaBlockChainAPI{b}
}

// BlockNumber returns the block number of the chain head.
func (s *KaiaBlockChainAPI) BlockNumber() hexutil.Uint64 {
	header, _ := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber) // latest header should always be available
	return hexutil.Uint64(header.Number.Uint64())
}

// ChainID returns the chain ID of the chain from genesis file.
func (s *KaiaBlockChainAPI) ChainID() *hexutil.Big {
	return s.ChainId()
}

// ChainId returns the chain ID of the chain from genesis file.
// This is for compatibility with ethereum client
func (s *KaiaBlockChainAPI) ChainId() *hexutil.Big {
	if s.b.ChainConfig() != nil {
		return (*hexutil.Big)(s.b.ChainConfig().ChainID)
	}
	return nil
}

// IsContractAccount returns true if the account associated with addr has a non-empty codeHash.
// It returns false otherwise.
func (s *KaiaBlockChainAPI) IsContractAccount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return false, err
	}
	return state.IsContractAccount(address), state.Error()
}

// IsHumanReadable returns true if the account associated with addr is a human-readable account.
// It returns false otherwise.
// func (s *KaiaBlockChainAPI) IsHumanReadable(ctx context.Context, address common.Address, blockNr rpc.BlockNumber) (bool, error) {
//	state, _, err := s.b.StateAndHeaderByNumber(ctx, blockNr)
//	if err != nil {
//		return false, err
//	}
//	return state.IsHumanReadable(address), state.Error()
// }

// GetBlockReceipts returns the receipts of all transactions in the block identified by number or hash.
func (s *KaiaBlockChainAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	block, err := s.b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	blockHash := block.Hash()
	receipts := s.b.GetBlockReceipts(ctx, blockHash)
	txs := block.Transactions()
	if receipts.Len() != txs.Len() {
		return nil, fmt.Errorf("the size of transactions and receipts is different in the block (%s)", blockHash.String())
	}
	fieldsList := make([]map[string]interface{}, 0, len(receipts))
	for index, receipt := range receipts {
		fields := RpcOutputReceipt(block.Header(), txs[index], blockHash, block.NumberU64(), uint64(index), receipt, s.b.ChainConfig())
		fieldsList = append(fieldsList, fields)
	}
	return fieldsList, nil
}

// GetBalance returns the amount of kei for the given address in the state of the
// given block number or hash. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers and hash are also allowed.
func (s *KaiaBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}

// AccountCreated returns true if the account associated with the address is created.
// It returns false otherwise.
func (s *KaiaBlockChainAPI) AccountCreated(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (bool, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return false, err
	}
	return state.Exist(address), state.Error()
}

// GetAccount returns account information of an input address.
func (s *KaiaBlockChainAPI) GetAccount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*account.AccountSerializer, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return &account.AccountSerializer{}, err
	}
	acc := state.GetAccount(address)
	if acc == nil {
		return &account.AccountSerializer{}, err
	}
	serAcc := account.NewAccountSerializerWithAccount(acc)
	return serAcc, state.Error()
}

func (s *KaiaAPI) ForkStatus(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	blkNum := block.Number()
	rules := s.b.ChainConfig().Rules(blkNum)
	status := make(map[string]interface{})

	rulesVal := reflect.ValueOf(rules)
	for i := 0; i < rulesVal.NumField(); i++ {
		val := rulesVal.Field(i)
		typ := rulesVal.Type().Field(i)
		status[typ.Name] = val.Interface()
	}
	// `IsKIP103` is not defined in the `Rules` struct. Exceptionally, we manually add it
	status["IsKIP103"] = s.b.ChainConfig().IsKIP103ForkBlock(blkNum)
	return status, nil
}

// rpcMarshalHeader converts the given header to the RPC output.
func (s *KaiaBlockChainAPI) rpcMarshalHeader(header *types.Header) map[string]interface{} {
	fields := RPCMarshalHeader(header, s.b.ChainConfig().Rules(header.Number))
	return fields
}

// GetHeaderByNumber returns the requested canonical block header.
func (s *KaiaBlockChainAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	header, err := s.b.HeaderByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	return s.rpcMarshalHeader(header), nil
}

// GetHeaderByHash returns the requested header by hash.
func (s *KaiaBlockChainAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	header, err := s.b.HeaderByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return s.rpcMarshalHeader(header), nil
}

// GetBlockByNumber returns the requested block. When blockNr is -1 the chain head is returned. When fullTx is true all
// transactions in the block are returned in full detail, otherwise only the transaction hash is returned.
func (s *KaiaBlockChainAPI) GetBlockByNumber(ctx context.Context, blockNr rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByNumber(ctx, blockNr)
	if block != nil && err == nil {
		response, err := s.rpcOutputBlock(block, true, fullTx)
		if err == nil && blockNr == rpc.PendingBlockNumber {
			// Pending blocks need to nil out a few fields
			for _, field := range []string{"hash", "nonce", "miner"} {
				response[field] = nil
			}
		}
		return response, err
	}
	return nil, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (s *KaiaBlockChainAPI) GetBlockByHash(ctx context.Context, blockHash common.Hash, fullTx bool) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, blockHash)
	if err != nil {
		return nil, err
	}
	return s.rpcOutputBlock(block, true, fullTx)
}

func (s *KaiaBlockChainAPI) GetBlockWithConsensusInfoByNumber(ctx context.Context, number *rpc.BlockNumber) (map[string]interface{}, error) {
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
		block = s.b.CurrentBlock()
		blockNumber = block.NumberU64()
	} else {
		// rpc.EarliestBlockNumber == 0, no need to treat it as a special case.
		blockNumber = uint64(number.Int64())
		block, _ = s.b.BlockByNumber(ctx, rpc.BlockNumber(blockNumber))
	}

	if block == nil {
		logger.Trace("Finding a block by number failed.", "blockNum", blockNumber)
		return nil, fmt.Errorf("the block does not exist (block number: %d)", blockNumber)
	}
	blockHash := block.Hash()

	cInfo, err := s.b.Engine().GetConsensusInfo(block)
	if err != nil {
		logger.Error("Getting consensus information failed", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := s.b.GetBlockReceipts(ctx, blockHash)
	if receipts == nil {
		receipts = s.b.GetBlockReceipts(ctx, blockHash)
	}

	return s.makeRPCBlockOutputWithConsensusInfo(block, cInfo, block.Transactions(), receipts), nil
}

func (s *KaiaBlockChainAPI) GetBlockWithConsensusInfoByNumberRange(ctx context.Context, start *rpc.BlockNumber, end *rpc.BlockNumber) (map[string]interface{}, error) {
	blocks := make(map[string]interface{})

	if start == nil || end == nil {
		logger.Trace("the range values should not be nil.", "start", start, "end", end)
		return nil, errRangeNil
	}

	// check error status.
	startNum := start.Int64()
	endNum := end.Int64()
	if startNum < 0 {
		logger.Trace("start should be positive", "start", startNum)
		return nil, errStartNotPositive
	}

	eChain := s.b.CurrentBlock().Number().Int64()
	if endNum > eChain {
		logger.Trace("end should be smaller than the latest block number", "end", end, "eChain", eChain)
		return nil, errEndLargetThanLatest
	}

	if startNum > endNum {
		logger.Trace("start should be smaller than end", "start", startNum, "end", endNum)
		return nil, errStartLargerThanEnd
	}

	if (endNum - startNum) > 50 {
		logger.Trace("number of requested blocks should be smaller than 50", "start", startNum, "end", endNum)
		return nil, errRequestedBlocksTooLarge
	}

	// gather start~end blocks
	for i := startNum; i <= endNum; i++ {
		strIdx := fmt.Sprintf("0x%x", i)

		blockNum := rpc.BlockNumber(i)
		b, err := s.GetBlockWithConsensusInfoByNumber(ctx, &blockNum)
		if err != nil {
			logger.Error("error on GetBlockWithConsensusInfoByNumber", "err", err)
			blocks[strIdx] = nil
		} else {
			blocks[strIdx] = b
		}
	}

	return blocks, nil
}

func (s *KaiaBlockChainAPI) GetBlockWithConsensusInfoByHash(ctx context.Context, blockHash common.Hash) (map[string]interface{}, error) {
	block, _ := s.b.BlockByHash(ctx, blockHash)
	if block == nil {
		logger.Trace("Finding a block failed.", "blockHash", blockHash)
		return nil, fmt.Errorf("the block does not exist (block hash: %s)", blockHash.String())
	}

	cInfo, err := s.b.Engine().GetConsensusInfo(block)
	if err != nil {
		logger.Error("Getting consensus information failed", "blockHash", blockHash, "err", err)
		return nil, errInternalError
	}

	receipts := s.b.GetBlockReceipts(ctx, blockHash)
	if receipts == nil {
		receipts = s.b.GetBlockReceipts(ctx, blockHash)
	}

	return s.makeRPCBlockOutputWithConsensusInfo(block, cInfo, block.Transactions(), receipts), nil
}

// GetCode returns the code stored at the given address in the state for the given block number or hash.
func (s *KaiaBlockChainAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	code := state.GetCode(address)
	return code, state.Error()
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers and hash are also allowed.
func (s *KaiaBlockChainAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}
	res := state.GetState(address, common.HexToHash(key))
	return res[:], state.Error()
}

// GetAccountKey returns the account key of EOA at a given address.
// If the account of the given address is a Legacy Account or a Smart Contract Account, it will return nil.
func (s *KaiaBlockChainAPI) GetAccountKey(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*accountkey.AccountKeySerializer, error) {
	state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return &accountkey.AccountKeySerializer{}, err
	}
	if state.Exist(address) == false {
		return nil, nil
	}
	accountKey := state.GetKey(address)
	serAccKey := accountkey.NewAccountKeySerializerWithAccountKey(accountKey)
	return serAccKey, state.Error()
}

// IsParallelDBWrite returns if parallel write is enabled or not.
// If enabled, data written in WriteBlockWithState is being written in parallel manner.
func (s *KaiaBlockChainAPI) IsParallelDBWrite() bool {
	return s.b.IsParallelDBWrite()
}

// IsSenderTxHashIndexingEnabled returns if senderTxHash to txHash mapping information
// indexing is enabled or not.
func (s *KaiaBlockChainAPI) IsSenderTxHashIndexingEnabled() bool {
	return s.b.IsSenderTxHashIndexingEnabled()
}

// IsConsoleLogEnabled returns if console log is enabled or not.
func (s *KaiaBlockChainAPI) IsConsoleLogEnabled() bool {
	return s.b.IsConsoleLogEnabled()
}

// CallArgs represents the arguments for a call.
type CallArgs struct {
	From                 common.Address  `json:"from"`
	To                   *common.Address `json:"to"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	Value                hexutil.Big     `json:"value"`
	Data                 hexutil.Bytes   `json:"data"`
	Input                hexutil.Bytes   `json:"input"`

	// Introduced by AccessListTxType transaction.
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

func (args *CallArgs) InputData() []byte {
	if args.Input != nil {
		return args.Input
	}
	if args.Data != nil {
		return args.Data
	}
	return nil
}

func (args *CallArgs) GetAccessList() types.AccessList {
	if args.AccessList != nil {
		return *args.AccessList
	} else {
		return nil
	}
}

func DoCall(ctx context.Context, b Backend, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, vmCfg vm.Config, timeout time.Duration, globalGasCap *big.Int) (*blockchain.ExecutionResult, uint64, error) {
	defer func(start time.Time) { logger.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, 0, err
	}
	// Setup context so it may be cancelled the call has completed
	// or, in case of unmetered gas, setup a context with a timeout.
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}
	// Make sure the context is cancelled when the call has completed
	// this makes sure resources are cleaned up.
	defer cancel()

	intrinsicGas, err := types.IntrinsicGas(args.InputData(), args.GetAccessList(), nil, args.To == nil, b.ChainConfig().Rules(header.Number))
	if err != nil {
		return nil, 0, err
	}

	// header.BaseFee != nil means magma hardforked
	var baseFee *big.Int
	if header.BaseFee != nil {
		baseFee = header.BaseFee
	} else {
		baseFee = new(big.Int).SetUint64(params.ZeroBaseFee)
	}
	msg, err := args.ToMessage(globalGasCap.Uint64(), baseFee, intrinsicGas)
	if err != nil {
		return nil, 0, err
	}

	// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
	state.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.EffectiveGasPrice(header, b.ChainConfig())))

	// The intrinsicGas is checked again later in the blockchain.ApplyMessage function,
	// but we check in advance here in order to keep StateTransition.TransactionDb method as unchanged as possible
	// and to clarify error reason correctly to serve eth namespace APIs.
	// This case is handled by DoEstimateGas function.
	if msg.Gas() < intrinsicGas {
		return nil, 0, fmt.Errorf("%w: msg.gas %d, want %d", blockchain.ErrIntrinsicGas, msg.Gas(), intrinsicGas)
	}
	evm, vmError, err := b.GetEVM(ctx, msg, state, header, vmCfg)
	if err != nil {
		return nil, 0, err
	}
	// Wait for the context to be done and cancel the evm. Even if the
	// EVM has finished, cancelling may be done (repeatedly)
	go func() {
		<-ctx.Done()
		evm.Cancel(vm.CancelByCtxDone)
	}()

	// Execute the message.
	result, err := blockchain.ApplyMessage(evm, msg)
	if err := vmError(); err != nil {
		return nil, 0, err
	}
	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, 0, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	if err != nil {
		return result, 0, fmt.Errorf("err: %w (supplied gas %d)", err, msg.Gas())
	}
	return result, evm.GetOpCodeComputationCost(), nil
}

// Call executes the given transaction on the state for the given block number or hash.
// It doesn't make and changes in the state/blockchain and is useful to execute and retrieve values.
func (s *KaiaBlockChainAPI) Call(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	gasCap := big.NewInt(0)
	if rpcGasCap := s.b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap
	}
	result, _, err := DoCall(ctx, s.b, args, blockNrOrHash, vm.Config{ComputationCostLimit: params.OpcodeComputationCostLimitInfinite, UseConsoleLog: s.b.IsConsoleLogEnabled()}, s.b.RPCEVMTimeout(), gasCap)
	if err != nil {
		return nil, err
	}

	if len(result.Revert()) > 0 {
		return nil, blockchain.NewRevertError(result)
	}
	return result.Return(), result.Unwrap()
}

func (s *KaiaBlockChainAPI) EstimateComputationCost(ctx context.Context, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	gasCap := big.NewInt(0)
	if rpcGasCap := s.b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap
	}
	_, computationCost, err := DoCall(ctx, s.b, args, blockNrOrHash, vm.Config{ComputationCostLimit: params.OpcodeComputationCostLimitInfinite}, s.b.RPCEVMTimeout(), gasCap)
	return (hexutil.Uint64)(computationCost), err
}

// EstimateGas returns an estimate of the amount of gas needed to execute the given transaction against the latest block.
func (s *KaiaBlockChainAPI) EstimateGas(ctx context.Context, args CallArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *EthStateOverride) (hexutil.Uint64, error) {
	gasCap := uint64(0)
	if rpcGasCap := s.b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	bNrOrHash := rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	return DoEstimateGas(ctx, s.b, args, bNrOrHash, overrides, s.b.RPCEVMTimeout(), new(big.Int).SetUint64(gasCap))
}

func DoEstimateGas(ctx context.Context, b Backend, args CallArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride, timeout time.Duration, gasCap *big.Int) (hexutil.Uint64, error) {
	var gasLimit uint64 = 0
	if args.Gas != nil {
		gasLimit = uint64(*args.Gas)
	}

	var feeCap *big.Int
	if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else {
		feeCap = common.Big0
	}

	state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return 0, err
	}
	if err := overrides.Apply(state); err != nil {
		return 0, err
	}
	balance := state.GetBalance(args.From) // from can't be nil

	// Create a helper to check if a gas allowance results in an executable transaction
	executable := func(gas uint64) (bool, *blockchain.ExecutionResult, error) {
		args.Gas = (*hexutil.Uint64)(&gas)
		result, _, err := DoCall(ctx, b, args, blockNrOrHash, vm.Config{ComputationCostLimit: params.OpcodeComputationCostLimitInfinite}, timeout, gasCap)
		if err != nil {
			if errors.Is(err, blockchain.ErrIntrinsicGas) || errors.Is(err, blockchain.ErrFloorDataGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			return true, nil, err // Bail out
		}
		return result.Failed(), result, nil
	}

	return blockchain.DoEstimateGas(ctx, gasLimit, gasCap.Uint64(), args.Value.ToInt(), feeCap, balance, executable)
}

// ExecutionResult groups all structured logs emitted by the EVM
// while replaying a transaction in debug mode as well as transaction
// execution status, the amount of gas used and the return value
type ExecutionResult struct {
	Gas         uint64         `json:"gas"`
	Failed      bool           `json:"failed"`
	ReturnValue string         `json:"returnValue"`
	StructLogs  []StructLogRes `json:"structLogs"`
}

// AccessListResult returns an optional accesslist
// Its the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type AccessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

// CreateAccessList creates a EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (s *KaiaBlockChainAPI) CreateAccessList(ctx context.Context, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (interface{}, error) {
	return doCreateAccessList(ctx, s.b, args, blockNrOrHash)
}

func (s *KaiaBlockChainAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*EthAccountResult, error) {
	return doGetProof(ctx, s.b, address, storageKeys, blockNrOrHash)
}

// StructLogRes stores a structured log emitted by the EVM while replaying a
// transaction in debug mode
type StructLogRes struct {
	Pc              uint64             `json:"pc"`
	Op              string             `json:"op"`
	Gas             uint64             `json:"gas"`
	GasCost         uint64             `json:"gasCost"`
	Depth           int                `json:"depth"`
	Error           error              `json:"error,omitempty"`
	Stack           *[]string          `json:"stack,omitempty"`
	Memory          *[]string          `json:"memory,omitempty"`
	Storage         *map[string]string `json:"storage,omitempty"`
	Computation     uint64             `json:"computation"`
	ComputationCost uint64             `json:"computationCost"`
}

// formatLogs formats EVM returned structured logs for json output
func FormatLogs(timeout time.Duration, logs []vm.StructLog) ([]StructLogRes, error) {
	logTimeout := false
	deadlineCtx, cancel := context.WithTimeout(context.Background(), timeout)
	go func() {
		<-deadlineCtx.Done()
		logger.Debug("trace logger timeout", "timeout", timeout, "err", deadlineCtx.Err())
		logTimeout = true
	}()
	defer cancel()

	formatted := make([]StructLogRes, len(logs))
	for index, trace := range logs {
		if logTimeout {
			return nil, fmt.Errorf("trace logger timeout")
		}
		formatted[index] = StructLogRes{
			Pc:              trace.Pc,
			Op:              trace.Op.String(),
			Gas:             trace.Gas,
			GasCost:         trace.GasCost,
			Depth:           trace.Depth,
			Error:           trace.Err,
			Computation:     trace.Computation,
			ComputationCost: trace.ComputationCost,
		}
		if trace.Stack != nil {
			stack := make([]string, len(trace.Stack))
			for i, stackValue := range trace.Stack {
				stack[i] = fmt.Sprintf("%x", math.PaddedBigBytes(stackValue, 32))
			}
			formatted[index].Stack = &stack
		}
		if trace.Memory != nil {
			memory := make([]string, 0, (len(trace.Memory)+31)/32)
			for i := 0; i+32 <= len(trace.Memory); i += 32 {
				memory = append(memory, fmt.Sprintf("%x", trace.Memory[i:i+32]))
			}
			formatted[index].Memory = &memory
		}
		if trace.Storage != nil {
			storage := make(map[string]string)
			for i, storageValue := range trace.Storage {
				storage[fmt.Sprintf("%x", i)] = fmt.Sprintf("%x", storageValue)
			}
			formatted[index].Storage = &storage
		}
	}
	return formatted, nil
}

// For kaia_getBlockByNumber, kaia_getBlockByHash, kaia_getBlockWithconsensusInfoByNumber, kaia_getBlockWithconsensusInfoByHash APIs
// and Kafka chaindatafetcher.
func RpcOutputBlock(b *types.Block, inclTx bool, fullTx bool, config *params.ChainConfig) (map[string]interface{}, error) {
	head := b.Header() // copies the header once
	fields := map[string]interface{}{
		"number":           (*hexutil.Big)(head.Number),
		"hash":             b.Hash(),
		"parentHash":       head.ParentHash,
		"logsBloom":        head.Bloom,
		"stateRoot":        head.Root,
		"reward":           head.Rewardbase,
		"blockScore":       (*hexutil.Big)(head.BlockScore),
		"extraData":        hexutil.Bytes(head.Extra),
		"governanceData":   hexutil.Bytes(head.Governance),
		"voteData":         hexutil.Bytes(head.Vote),
		"size":             hexutil.Uint64(b.Size()),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        (*hexutil.Big)(head.Time),
		"timestampFoS":     (hexutil.Uint)(head.TimeFoS),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}

		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newRPCTransactionFromBlockHash(b, tx.Hash(), config), nil
			}
		}

		txs := b.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range b.Transactions() {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}

	rules := config.Rules(b.Number())
	if rules.IsMagma {
		if head.BaseFee == nil {
			fields["baseFeePerGas"] = (*hexutil.Big)(new(big.Int).SetUint64(params.ZeroBaseFee))
		} else {
			fields["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
		}
	}
	if rules.IsRandao {
		fields["randomReveal"] = hexutil.Bytes(head.RandomReveal)
		fields["mixHash"] = hexutil.Bytes(head.MixHash)
	}

	return fields, nil
}

// rpcOutputBlock converts the given block to the RPC output which depends on fullTx. If inclTx is true transactions are
// returned. When fullTx is true the returned block contains full transaction details, otherwise it will only contain
// transaction hashes.
func (s *KaiaBlockChainAPI) rpcOutputBlock(b *types.Block, inclTx bool, fullTx bool) (map[string]interface{}, error) {
	return RpcOutputBlock(b, inclTx, fullTx, s.b.ChainConfig())
}

func getFrom(tx *types.Transaction) common.Address {
	var from common.Address
	if tx.IsEthereumTransaction() {
		signer := types.LatestSignerForChainID(tx.ChainId())
		from, _ = types.Sender(signer, tx)
	} else {
		from, _ = tx.From()
	}
	return from
}

func NewRPCTransaction(head *types.Header, tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, config *params.ChainConfig) map[string]interface{} {
	return newRPCTransaction(head, tx, blockHash, blockNumber, index, config)
}

// newRPCTransaction returns a transaction that will serialize to the RPC
// representation, with the given location metadata set (if available).
func newRPCTransaction(header *types.Header, tx *types.Transaction, blockHash common.Hash, blockNumber uint64, index uint64, config *params.ChainConfig) map[string]interface{} {
	output := tx.MakeRPCOutput()
	output["senderTxHash"] = tx.SenderTxHashAll()
	output["blockHash"] = blockHash
	output["blockNumber"] = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
	output["from"] = getFrom(tx)
	output["hash"] = tx.Hash()
	output["transactionIndex"] = hexutil.Uint(index)
	if tx.Type() == types.TxTypeEthereumDynamicFee || tx.Type() == types.TxTypeEthereumSetCode {
		if header != nil {
			output["gasPrice"] = (*hexutil.Big)(tx.EffectiveGasPrice(header, config))
		} else {
			// transaction is not processed yet
			output["gasPrice"] = (*hexutil.Big)(tx.EffectiveGasPrice(nil, nil))
		}
	}

	return output
}

// newRPCPendingTransaction returns a pending transaction that will serialize to the RPC representation
func newRPCPendingTransaction(tx *types.Transaction, config *params.ChainConfig) map[string]interface{} {
	return newRPCTransaction(nil, tx, common.Hash{}, 0, 0, config)
}

// newRPCTransactionFromBlockIndex returns a transaction that will serialize to the RPC representation.
// non-null of b(block) is guaranteed
func newRPCTransactionFromBlockIndex(b *types.Block, index uint64, config *params.ChainConfig) map[string]interface{} {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	return newRPCTransaction(b.Header(), txs[index], b.Hash(), b.NumberU64(), index, config)
}

// RPCMarshalHeader converts the given header to the RPC output that includes Klaytn-specific fields.
// For kaia_getHeaderByNumber and kaia_getHeaderByHash APIs.
func RPCMarshalHeader(head *types.Header, rules params.Rules) map[string]interface{} {
	result := map[string]interface{}{
		"parentHash":       head.ParentHash,
		"reward":           head.Rewardbase,
		"stateRoot":        head.Root,
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
		"logsBloom":        head.Bloom,
		"blockScore":       (*hexutil.Big)(head.BlockScore),
		"number":           (*hexutil.Big)(head.Number),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        (*hexutil.Big)(head.Time),
		"timestampFoS":     hexutil.Uint(head.TimeFoS),
		"extraData":        hexutil.Bytes(head.Extra),
		"governanceData":   hexutil.Bytes(head.Governance),
		"voteData":         hexutil.Bytes(head.Vote),
		"hash":             head.Hash(),
	}

	if rules.IsEthTxType {
		if head.BaseFee == nil {
			result["baseFeePerGas"] = (*hexutil.Big)(new(big.Int).SetUint64(params.ZeroBaseFee))
		} else {
			result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
		}
	}
	if rules.IsRandao {
		result["randomReveal"] = hexutil.Bytes(head.RandomReveal)
		result["mixhash"] = hexutil.Bytes(head.MixHash)
	}

	return result
}

// newRPCRawTransactionFromBlockIndex returns the bytes of a transaction given a block and a transaction index.
func newRPCRawTransactionFromBlockIndex(b *types.Block, index uint64) hexutil.Bytes {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		return nil
	}
	blob, _ := rlp.EncodeToBytes(txs[index])
	return blob
}

// newRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newRPCTransactionFromBlockHash(b *types.Block, hash common.Hash, config *params.ChainConfig) map[string]interface{} {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newRPCTransactionFromBlockIndex(b, uint64(idx), config)
		}
	}
	return nil
}

func (args *CallArgs) ToMessage(globalGasCap uint64, baseFee *big.Int, intrinsicGas uint64) (*types.Transaction, error) {
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.MaxFeePerGas != nil && args.MaxPriorityFeePerGas != nil {
		if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return nil, errors.New("MaxPriorityFeePerGas is greater than MaxFeePerGas")
		}
	}

	// Set sender address or use zero address if none specified.
	addr := args.From

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = params.UpperGasLimit
	}
	if args.Gas != nil {
		gas = uint64(*args.Gas)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		logger.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}

	// Do not update gasPrice unless any of args.GasPrice and args.MaxFeePerGas is specified.
	gasPrice := new(big.Int)
	if args.GasPrice != nil {
		gasPrice = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		gasPrice = args.MaxFeePerGas.ToInt()
	}

	if baseFee.Cmp(new(big.Int).SetUint64(params.ZeroBaseFee)) != 0 {
		if args.GasPrice == nil && args.MaxFeePerGas == nil {
			// User specified neither GasPrice nor MaxFeePerGas, use baseFee
			gasPrice = new(big.Int).Mul(baseFee, common.Big2)
		}
	}
	value := new(big.Int)
	if &args.Value != nil {
		value = args.Value.ToInt()
	}

	var accessList types.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}
	return types.NewMessage(addr, args.To, 0, value, gas, gasPrice, nil, nil, args.InputData(), false, intrinsicGas, accessList, nil, nil), nil
}

// MainnetCredit contract is stored in the address zero.
var (
	mainnetCreditContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")
	latestBlockNrOrHash          = rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	errNoCypressCreditContract   = errors.New("no mainnet credit contract")
)

// callCypressCreditGetFunc executes funcName in CypressCreditContract and returns the output.
func (s *KaiaBlockChainAPI) callCypressCreditGetFunc(ctx context.Context, parsed *abi.ABI, funcName string) (*string, error) {
	abiGet, err := parsed.Pack(funcName)
	if err != nil {
		return nil, err
	}

	args := CallArgs{
		To:   &mainnetCreditContractAddress,
		Data: abiGet,
	}
	ret, err := s.Call(ctx, args, latestBlockNrOrHash)
	if err != nil {
		return nil, err
	}

	output := new(string)
	err = parsed.UnpackIntoInterface(output, funcName, ret)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetCypressCredit calls getPhoto and getNames in the CypressCredit contract
// and returns all the results as a struct.
func (s *KaiaBlockChainAPI) GetCypressCredit(ctx context.Context) (*CreditOutput, error) {
	if ok, err := s.IsContractAccount(ctx, mainnetCreditContractAddress, latestBlockNrOrHash); err != nil {
		return nil, err
	} else if !ok {
		return nil, errNoCypressCreditContract
	}

	parsed, err := abi.JSON(strings.NewReader(misc.CypressCreditV2ABI))
	if err != nil {
		return nil, err
	}

	output := new(CreditOutput)

	// getPhoto and getNames must exist from the Cypress genesis.
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getPhoto"); err == nil {
		output.Photo = *str
	} else {
		return nil, err
	}
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getNames"); err == nil {
		output.Names = *str
	} else {
		return nil, err
	}

	// getEndingPhoto and getEndingNames are added at some nonzero block. They are returned if they exist.
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getEndingPhoto"); err == nil {
		output.EndingPhoto = *str
	}
	if str, err := s.callCypressCreditGetFunc(ctx, &parsed, "getEndingNames"); err == nil {
		output.EndingNames = *str
	}

	return output, nil
}

func (s *KaiaBlockChainAPI) makeRPCBlockOutputWithConsensusInfo(b *types.Block,
	cInfo consensus.ConsensusInfo, transactions types.Transactions, receipts types.Receipts,
) map[string]interface{} {
	head := b.Header() // copies the header once
	hash := head.Hash()

	r, err := RpcOutputBlock(b, false, false, s.b.ChainConfig())
	if err != nil {
		logger.Error("failed to RpcOutputBlock", "err", err)
		return nil
	}

	// make transactions
	numTxs := len(transactions)
	rpcTransactions := make([]map[string]interface{}, numTxs)
	for i, tx := range transactions {
		if len(receipts) == len(transactions) {
			rpcTransactions[i] = RpcOutputReceipt(head, tx, hash, head.Number.Uint64(), uint64(i), receipts[i], s.b.ChainConfig())
		} else {
			// fill the transaction output if receipt is not found
			rpcTransactions[i] = NewRPCTransaction(head, tx, hash, head.Number.Uint64(), uint64(i), s.b.ChainConfig())
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
