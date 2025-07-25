// Modifications Copyright 2024 The Kaia Authors
// Copyright 2021 The klaytn Authors
// This file is part of the klaytn library.
//
// Modifications Copyright 2024 The Kaia Authors
// Copyright 2021 The klaytn Authors
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

package api

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/statedb"
)

const (
	// EmptySha3Uncles always have value which is the result of
	// `crypto.Keccak256Hash(rlp.EncodeToBytes([]*types.Header(nil)).String())`
	// because there is no uncles in Kaia.
	// Just use const value because we don't have to calculate it everytime which always be same result.
	EmptySha3Uncles = "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347"
	// ZeroHashrate exists for supporting Ethereum compatible data structure.
	// There is no POW mining mechanism in Kaia.
	ZeroHashrate uint64 = 0
	// ZeroUncleCount is always zero because there is no uncle blocks in Kaia.
	ZeroUncleCount uint = 0
)

var (
	errNoMiningWork  = errors.New("no mining work available yet")
	errNotFoundBlock = errors.New("can't find a block in database")
)

// EthAPI provides an API to access the Kaia through the `eth` namespace.
type EthAPI struct {
	kaiaAPI            *KaiaAPI
	kaiaBlockChainAPI  *KaiaBlockChainAPI
	kaiaTransactionAPI *KaiaTransactionAPI
	kaiaAccountAPI     *KaiaAccountAPI
	nodeAddress        common.Address
}

// NewEthAPI creates a new ethereum API.
// EthAPI operates using Kaia's API internally without overriding.
// Therefore, it is necessary to use APIs defined in two different packages(cn and api),
// so those apis will be defined through a setter.
func NewEthAPI(
	kaiaAPI *KaiaAPI,
	kaiaBlockChainAPI *KaiaBlockChainAPI,
	kaiaTransactionAPI *KaiaTransactionAPI,
	kaiaAccountAPI *KaiaAccountAPI,
	nodeAddress common.Address,
) *EthAPI {
	return &EthAPI{
		kaiaAPI,
		kaiaBlockChainAPI,
		kaiaTransactionAPI,
		kaiaAccountAPI,
		nodeAddress,
	}
}

// Etherbase is the address of operating node.
// Unlike Ethereum, it only returns the node address because Kaia does not have a POW mechanism.
func (api *EthAPI) Etherbase() (common.Address, error) {
	return api.nodeAddress, nil
}

// Coinbase is the address of operating node (alias for Etherbase).
func (api *EthAPI) Coinbase() (common.Address, error) {
	return api.Etherbase()
}

// Hashrate returns the POW hashrate.
// Unlike Ethereum, it always returns ZeroHashrate because Kaia does not have a POW mechanism.
func (api *EthAPI) Hashrate() hexutil.Uint64 {
	return hexutil.Uint64(ZeroHashrate)
}

// Mining returns an indication if this node is currently mining.
// Unlike Ethereum, it always returns false because Kaia does not have a POW mechanism,
func (api *EthAPI) Mining() bool {
	return false
}

// GetWork returns an errNoMiningWork because Kaia does not have a POW mechanism.
func (api *EthAPI) GetWork() ([4]string, error) {
	return [4]string{}, errNoMiningWork
}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

// SubmitWork returns false because Kaia does not have a POW mechanism.
func (api *EthAPI) SubmitWork(nonce BlockNonce, hash, digest common.Hash) bool {
	return false
}

// SubmitHashrate returns false because Kaia does not have a POW mechanism.
func (api *EthAPI) SubmitHashrate(rate hexutil.Uint64, id common.Hash) bool {
	return false
}

// GetHashrate returns ZeroHashrate because Kaia does not have a POW mechanism.
func (api *EthAPI) GetHashrate() uint64 {
	return ZeroHashrate
}

// GasPrice returns a suggestion for a gas price.
func (api *EthAPI) GasPrice(ctx context.Context) (*hexutil.Big, error) {
	return api.kaiaAPI.GasPrice(ctx)
}

func (api *EthAPI) UpperBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(api.kaiaAPI.UpperBoundGasPrice(ctx))
}

func (api *EthAPI) LowerBoundGasPrice(ctx context.Context) *hexutil.Big {
	return (*hexutil.Big)(api.kaiaAPI.LowerBoundGasPrice(ctx))
}

// MaxPriorityFeePerGas returns a suggestion for a gas tip cap for dynamic fee transactions.
func (api *EthAPI) MaxPriorityFeePerGas(ctx context.Context) (*hexutil.Big, error) {
	return api.kaiaAPI.MaxPriorityFeePerGas(ctx)
}

// DecimalOrHex unmarshals a non-negative decimal or hex parameter into a uint64.
type DecimalOrHex uint64

// UnmarshalJSON implements json.Unmarshaler.
func (dh *DecimalOrHex) UnmarshalJSON(data []byte) error {
	input := strings.TrimSpace(string(data))
	if len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"' {
		input = input[1 : len(input)-1]
	}

	value, err := strconv.ParseUint(input, 10, 64)
	if err != nil {
		value, err = hexutil.DecodeUint64(input)
	}
	if err != nil {
		return err
	}
	*dh = DecimalOrHex(value)
	return nil
}

func (api *EthAPI) FeeHistory(ctx context.Context, blockCount DecimalOrHex, lastBlock rpc.BlockNumber, rewardPercentiles []float64) (*FeeHistoryResult, error) {
	return api.kaiaAPI.FeeHistory(ctx, blockCount, lastBlock, rewardPercentiles)
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (api *EthAPI) Syncing() (interface{}, error) {
	return api.kaiaAPI.Syncing()
}

// ChainId is the EIP-155 replay-protection chain id for the current ethereum chain config.
func (api *EthAPI) ChainId() (*hexutil.Big, error) {
	return api.kaiaBlockChainAPI.ChainId(), nil
}

// BlockNumber returns the block number of the chain head.
func (api *EthAPI) BlockNumber() hexutil.Uint64 {
	return api.kaiaBlockChainAPI.BlockNumber()
}

// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (api *EthAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
	return api.kaiaBlockChainAPI.GetBalance(ctx, address, blockNrOrHash)
}

// EthAccountResult structs for GetProof
// AccountResult in go-ethereum has been renamed to EthAccountResult.
// AccountResult is defined in go-ethereum's internal package, so AccountResult is redefined here as EthAccountResult.
type EthAccountResult struct {
	Address      common.Address     `json:"address"`
	AccountProof []string           `json:"accountProof"`
	Balance      *hexutil.Big       `json:"balance"`
	CodeHash     common.Hash        `json:"codeHash"`
	Nonce        hexutil.Uint64     `json:"nonce"`
	StorageHash  common.Hash        `json:"storageHash"`
	StorageProof []EthStorageResult `json:"storageProof"`
}

// StorageResult in go-ethereum has been renamed to EthStorageResult.
// StorageResult is defined in go-ethereum's internal package, so StorageResult is redefined here as EthStorageResult.
type EthStorageResult struct {
	Key   string       `json:"key"`
	Value *hexutil.Big `json:"value"`
	Proof []string     `json:"proof"`
}

// proofList implements KeyValueWriter and collects the proofs as
// hex-strings for delivery to rpc-caller.
type proofList []string

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, hexutil.Encode(value))
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

func (n *proofList) WriteMerkleProof(key, value []byte) {
	n.Put(key, value)
}

// decodeHash parses a hex-encoded 32-byte hash. The input may optionally
// be prefixed by 0x and can have a byte length up to 32.
func decodeHash(s string) (h common.Hash, inputLength int, err error) {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if (len(s) & 1) > 0 {
		s = "0" + s
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return common.Hash{}, 0, errors.New("hex string invalid")
	}
	if len(b) > 32 {
		return common.Hash{}, len(b), errors.New("hex string too long, want at most 32 bytes")
	}
	return common.BytesToHash(b), len(b), nil
}

func doGetProof(ctx context.Context, b Backend, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*EthAccountResult, error) {
	var (
		keys         = make([]common.Hash, len(storageKeys))
		keyLengths   = make([]int, len(storageKeys))
		storageProof = make([]EthStorageResult, len(storageKeys))
	)
	// Deserialize all keys. This prevents state access on invalid input.
	for i, hexKey := range storageKeys {
		var err error
		keys[i], keyLengths[i], err = decodeHash(hexKey)
		if err != nil {
			return nil, err
		}
	}
	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	codeHash := state.GetCodeHash(address)

	contractStorageRootExt, err := state.GetContractStorageRoot(address)
	if err != nil {
		return nil, err
	}
	contractStorageRoot := contractStorageRootExt.Unextend()

	// if we have a storageTrie, (which means the account exists), we can update the storagehash
	if len(keys) > 0 {
		storageTrie, err := statedb.NewTrie(contractStorageRoot, state.Database().TrieDB(), nil)
		if err != nil {
			return nil, err
		}
		// Create the proofs for the storageKeys.
		for i, key := range keys {
			// Output key encoding is a bit special: if the input was a 32-byte hash, it is
			// returned as such. Otherwise, we apply the QUANTITY encoding mandated by the
			// JSON-RPC spec for getProof. This behavior exists to preserve backwards
			// compatibility with older client versions.
			var outputKey string
			if keyLengths[i] != 32 {
				outputKey = hexutil.EncodeBig(key.Big())
			} else {
				outputKey = hexutil.Encode(key[:])
			}
			if storageTrie == nil {
				storageProof[i] = EthStorageResult{outputKey, &hexutil.Big{}, []string{}}
				continue
			}
			var proof proofList
			if err := storageTrie.Prove(crypto.Keccak256(key.Bytes()), 0, &proof); err != nil {
				return nil, err
			}
			value := (*hexutil.Big)(state.GetState(address, key).Big())
			storageProof[i] = EthStorageResult{outputKey, value, proof}
		}
	}

	// Create the accountProof.
	trie, err := statedb.NewTrie(header.Root, state.Database().TrieDB(), nil)
	if err != nil {
		return nil, err
	}
	var accountProof proofList
	if err := trie.Prove(crypto.Keccak256(address.Bytes()), 0, &accountProof); err != nil {
		return nil, err
	}

	return &EthAccountResult{
		Address:      address,
		AccountProof: accountProof,
		Balance:      (*hexutil.Big)(state.GetBalance(address)),
		CodeHash:     codeHash,
		Nonce:        hexutil.Uint64(state.GetNonce(address)),
		StorageHash:  contractStorageRoot,
		StorageProof: storageProof,
	}, state.Error()
}

// GetProof returns the Merkle-proof for a given account and optionally some storage keys
func (api *EthAPI) GetProof(ctx context.Context, address common.Address, storageKeys []string, blockNrOrHash rpc.BlockNumberOrHash) (*EthAccountResult, error) {
	return doGetProof(ctx, api.kaiaBlockChainAPI.b, address, storageKeys, blockNrOrHash)
}

// GetHeaderByNumber returns the requested canonical block header.
// * When blockNr is -1 the chain head is returned.
// * When blockNr is -2 the pending chain head is returned.
func (api *EthAPI) GetHeaderByNumber(ctx context.Context, number rpc.BlockNumber) (map[string]interface{}, error) {
	// In Ethereum, err is always nil because the backend of Ethereum always return nil.
	header, err := api.kaiaBlockChainAPI.b.HeaderByNumber(ctx, number)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	inclMiner := number != rpc.PendingBlockNumber
	response, err := api.rpcMarshalHeader(header, inclMiner)
	if err != nil {
		return nil, err
	}
	if number == rpc.PendingBlockNumber {
		// Pending header need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}
	return response, nil
}

// GetHeaderByHash returns the requested header by hash.
func (api *EthAPI) GetHeaderByHash(ctx context.Context, hash common.Hash) map[string]interface{} {
	// In Ethereum, err is always nil because the backend of Ethereum always return nil.
	header, _ := api.kaiaBlockChainAPI.b.HeaderByHash(ctx, hash)
	if header != nil {
		response, err := api.rpcMarshalHeader(header, true)
		if err != nil {
			return nil
		}
		return response
	}
	return nil
}

// GetBlockByNumber returns the requested canonical block.
//   - When blockNr is -1 the chain head is returned.
//   - When blockNr is -2 the pending chain head is returned.
//   - When fullTx is true all transactions in the block are returned, otherwise
//     only the transaction hash is returned.
func (api *EthAPI) GetBlockByNumber(ctx context.Context, number rpc.BlockNumber, fullTx bool) (map[string]interface{}, error) {
	// Kaia backend returns error when there is no matched block but
	// Ethereum returns it as nil without error, so we should return is as nil when there is no matched block.
	block, err := api.kaiaBlockChainAPI.b.BlockByNumber(ctx, number)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}

	inclMiner := number != rpc.PendingBlockNumber
	inclTx := true
	response, err := api.rpcMarshalBlock(block, inclMiner, inclTx, fullTx)
	if err == nil && number == rpc.PendingBlockNumber {
		// Pending blocks need to nil out a few fields
		for _, field := range []string{"hash", "nonce", "miner"} {
			response[field] = nil
		}
	}
	return response, err
}

// GetBlockByHash returns the requested block. When fullTx is true all transactions in the block are returned in full
// detail, otherwise only the transaction hash is returned.
func (api *EthAPI) GetBlockByHash(ctx context.Context, hash common.Hash, fullTx bool) (map[string]interface{}, error) {
	// Kaia backend returns error when there is no matched block but
	// Ethereum returns it as nil without error, so we should return is as nil when there is no matched block.
	block, err := api.kaiaBlockChainAPI.b.BlockByHash(ctx, hash)
	if err != nil {
		if strings.Contains(err.Error(), "does not exist") {
			return nil, nil
		}
		return nil, err
	}
	return api.rpcMarshalBlock(block, true, true, fullTx)
}

// GetUncleByBlockNumberAndIndex returns nil because there is no uncle block in Kaia.
func (api *EthAPI) GetUncleByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) (map[string]interface{}, error) {
	return nil, nil
}

// GetUncleByBlockHashAndIndex returns nil because there is no uncle block in Kaia.
func (api *EthAPI) GetUncleByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) (map[string]interface{}, error) {
	return nil, nil
}

// GetUncleCountByBlockNumber returns 0 when given blockNr exists because there is no uncle block in Kaia.
func (api *EthAPI) GetUncleCountByBlockNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	if block, _ := api.kaiaBlockChainAPI.b.BlockByNumber(ctx, blockNr); block != nil {
		n := hexutil.Uint(ZeroUncleCount)
		return &n
	}
	return nil
}

// GetUncleCountByBlockHash returns 0 when given blockHash exists because there is no uncle block in Kaia.
func (api *EthAPI) GetUncleCountByBlockHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	if block, _ := api.kaiaBlockChainAPI.b.BlockByHash(ctx, blockHash); block != nil {
		n := hexutil.Uint(ZeroUncleCount)
		return &n
	}
	return nil
}

// GetCode returns the code stored at the given address in the state for the given block number.
func (api *EthAPI) GetCode(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return api.kaiaBlockChainAPI.GetCode(ctx, address, blockNrOrHash)
}

// GetStorageAt returns the storage from the state at the given address, key and
// block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta block
// numbers are also allowed.
func (api *EthAPI) GetStorageAt(ctx context.Context, address common.Address, key string, blockNrOrHash rpc.BlockNumberOrHash) (hexutil.Bytes, error) {
	return api.kaiaBlockChainAPI.GetStorageAt(ctx, address, key, blockNrOrHash)
}

// EthOverrideAccount indicates the overriding fields of account during the execution
// of a message call.
// Note, state and stateDiff can't be specified at the same time. If state is
// set, message execution will only use the data in the given state. Otherwise
// if statDiff is set, all diff will be applied first and then execute the call
// message.
// OverrideAccount in go-ethereum has been renamed to EthOverrideAccount.
// OverrideAccount is defined in go-ethereum's internal package, so OverrideAccount is redefined here as EthOverrideAccount.
type EthOverrideAccount struct {
	Nonce     *hexutil.Uint64              `json:"nonce"`
	Code      *hexutil.Bytes               `json:"code"`
	Balance   **hexutil.Big                `json:"balance"`
	State     *map[common.Hash]common.Hash `json:"state"`
	StateDiff *map[common.Hash]common.Hash `json:"stateDiff"`
}

// EthStateOverride is the collection of overridden accounts.
// StateOverride in go-ethereum has been renamed to EthStateOverride.
// StateOverride is defined in go-ethereum's internal package, so StateOverride is redefined here as EthStateOverride.
type EthStateOverride map[common.Address]EthOverrideAccount

func (diff *EthStateOverride) Apply(state *state.StateDB) error {
	if diff == nil {
		return nil
	}
	for addr, account := range *diff {
		// Override account nonce.
		if account.Nonce != nil {
			state.SetNonce(addr, uint64(*account.Nonce))
		}
		// Override account(contract) code.
		if account.Code != nil {
			state.SetCode(addr, *account.Code)
		}
		// Override account balance.
		if account.Balance != nil {
			state.SetBalance(addr, (*big.Int)(*account.Balance))
		}
		if account.State != nil && account.StateDiff != nil {
			return fmt.Errorf("account %s has both 'state' and 'stateDiff'", addr.Hex())
		}
		// Replace entire state if caller requires.
		if account.State != nil {
			state.SetStorage(addr, *account.State)
		}
		// Apply state diff into specified accounts.
		if account.StateDiff != nil {
			for key, value := range *account.StateDiff {
				state.SetState(addr, key, value)
			}
		}
	}
	return nil
}

// Call executes the given transaction on the state for the given block number.
//
// Additionally, the caller can specify a batch of contract for fields overriding.
//
// Note, this function doesn't make and changes in the state/blockchain and is
// useful to execute and retrieve values.
func (api *EthAPI) Call(ctx context.Context, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride) (hexutil.Bytes, error) {
	bcAPI := api.kaiaBlockChainAPI.b
	gasCap := uint64(0)
	if rpcGasCap := bcAPI.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	result, err := EthDoCall(ctx, bcAPI, args, blockNrOrHash, overrides, bcAPI.RPCEVMTimeout(), gasCap)
	if err != nil {
		return nil, err
	}

	if len(result.Revert()) > 0 {
		return nil, blockchain.NewRevertError(result)
	}
	return result.Return(), result.Unwrap()
}

// EstimateGas returns an estimate of the amount of gas needed to execute the
// given transaction against the current pending block.
func (api *EthAPI) EstimateGas(ctx context.Context, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash, overrides *EthStateOverride) (hexutil.Uint64, error) {
	bcAPI := api.kaiaBlockChainAPI.b
	bNrOrHash := rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	gasCap := uint64(0)
	if rpcGasCap := bcAPI.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}
	return EthDoEstimateGas(ctx, bcAPI, args, bNrOrHash, overrides, gasCap)
}

// GetBlockTransactionCountByNumber returns the number of transactions in the block with the given block number.
func (api *EthAPI) GetBlockTransactionCountByNumber(ctx context.Context, blockNr rpc.BlockNumber) *hexutil.Uint {
	return api.kaiaTransactionAPI.GetBlockTransactionCountByNumber(ctx, blockNr)
}

// GetBlockTransactionCountByHash returns the number of transactions in the block with the given hash.
func (api *EthAPI) GetBlockTransactionCountByHash(ctx context.Context, blockHash common.Hash) *hexutil.Uint {
	return api.kaiaTransactionAPI.GetBlockTransactionCountByHash(ctx, blockHash)
}

// EthRPCTransaction represents a transaction that will serialize to the RPC representation of a transaction
// RPCTransaction in go-ethereum has been renamed to EthRPCTransaction.
// RPCTransaction is defined in go-ethereum's internal package, so RPCTransaction is redefined here as EthRPCTransaction.
type EthRPCTransaction struct {
	BlockHash         *common.Hash                 `json:"blockHash"`
	BlockNumber       *hexutil.Big                 `json:"blockNumber"`
	From              common.Address               `json:"from"`
	Gas               hexutil.Uint64               `json:"gas"`
	GasPrice          *hexutil.Big                 `json:"gasPrice"`
	GasFeeCap         *hexutil.Big                 `json:"maxFeePerGas,omitempty"`
	GasTipCap         *hexutil.Big                 `json:"maxPriorityFeePerGas,omitempty"`
	Hash              common.Hash                  `json:"hash"`
	Input             hexutil.Bytes                `json:"input"`
	Nonce             hexutil.Uint64               `json:"nonce"`
	To                *common.Address              `json:"to"`
	TransactionIndex  *hexutil.Uint64              `json:"transactionIndex"`
	Value             *hexutil.Big                 `json:"value"`
	Type              hexutil.Uint64               `json:"type"`
	Accesses          *types.AccessList            `json:"accessList,omitempty"`
	ChainID           *hexutil.Big                 `json:"chainId,omitempty"`
	AuthorizationList []types.SetCodeAuthorization `json:"authorizationList,omitempty"`
	V                 *hexutil.Big                 `json:"v"`
	R                 *hexutil.Big                 `json:"r"`
	S                 *hexutil.Big                 `json:"s"`
}

// ethTxJSON is the JSON representation of Ethereum transaction.
// ethTxJSON is used by eth namespace APIs which returns Transaction object as it is.
// Because every transaction in Kaia, implements json.Marshaler interface (MarshalJSON), but
// it is marshaled for Kaia format only.
// e.g. Ethereum transaction have V, R, and S field for signature but,
// Kaia transaction have types.TxSignaturesJSON which includes array of signatures which is not
// applicable for Ethereum transaction.
type ethTxJSON struct {
	Type hexutil.Uint64 `json:"type"`

	// Common transaction fields:
	Nonce                *hexutil.Uint64 `json:"nonce"`
	GasPrice             *hexutil.Big    `json:"gasPrice"`
	MaxPriorityFeePerGas *hexutil.Big    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big    `json:"maxFeePerGas"`
	Gas                  *hexutil.Uint64 `json:"gas"`
	Value                *hexutil.Big    `json:"value"`
	Data                 *hexutil.Bytes  `json:"input"`
	V                    *hexutil.Big    `json:"v"`
	R                    *hexutil.Big    `json:"r"`
	S                    *hexutil.Big    `json:"s"`
	To                   *common.Address `json:"to"`

	// Access list transaction fields:
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
	AccessList *types.AccessList `json:"accessList,omitempty"`

	// Set code transaction fields:
	AuthorizationList []types.SetCodeAuthorization `json:"authorizationList,omitempty"`

	// Only used for encoding:
	Hash common.Hash `json:"hash"`
}

// newEthRPCTransactionFromBlockIndex creates an EthRPCTransaction from block and index parameters.
func newEthRPCTransactionFromBlockIndex(b *types.Block, index uint64, config *params.ChainConfig) *EthRPCTransaction {
	txs := b.Transactions()
	if index >= uint64(len(txs)) {
		logger.Error("invalid transaction index", "given index", index, "length of txs", len(txs))
		return nil
	}
	return newEthRPCTransaction(b, txs[index], b.Hash(), b.NumberU64(), index, config)
}

// newEthRPCTransactionFromBlockHash returns a transaction that will serialize to the RPC representation.
func newEthRPCTransactionFromBlockHash(b *types.Block, hash common.Hash, config *params.ChainConfig) *EthRPCTransaction {
	for idx, tx := range b.Transactions() {
		if tx.Hash() == hash {
			return newEthRPCTransactionFromBlockIndex(b, uint64(idx), config)
		}
	}
	return nil
}

// resolveToField returns value which fits to `to` field based on transaction types.
// This function is used when converting Kaia transactions to Ethereum transaction types.
func resolveToField(tx *types.Transaction) *common.Address {
	switch tx.Type() {
	case types.TxTypeAccountUpdate, types.TxTypeFeeDelegatedAccountUpdate, types.TxTypeFeeDelegatedAccountUpdateWithRatio,
		types.TxTypeCancel, types.TxTypeFeeDelegatedCancel, types.TxTypeFeeDelegatedCancelWithRatio,
		types.TxTypeChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoring, types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		// These type of transactions actually do not have `to` address, but Ethereum always have `to` field,
		// so we Kaia developers decided to fill the `to` field with `from` address value in these case.
		from := getFrom(tx)
		return &from
	}
	return tx.To()
}

// newEthRPCTransaction creates an EthRPCTransaction from Kaia transaction.
func newEthRPCTransaction(block *types.Block, tx *types.Transaction, blockHash common.Hash, blockNumber, index uint64, config *params.ChainConfig) *EthRPCTransaction {
	// When an unknown transaction is requested through rpc call,
	// nil is returned by Kaia API, and it is handled.
	if tx == nil {
		return nil
	}

	typeInt := tx.Type()
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		typeInt = types.TxTypeLegacyTransaction
	}

	signature := tx.GetTxInternalData().RawSignatureValues()[0]

	result := &EthRPCTransaction{
		Type:     hexutil.Uint64(byte(typeInt)),
		From:     getFrom(tx),
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: (*hexutil.Big)(tx.GasPrice()),
		Hash:     tx.Hash(),
		Input:    tx.Data(),
		Nonce:    hexutil.Uint64(tx.Nonce()),
		To:       resolveToField(tx),
		Value:    (*hexutil.Big)(tx.Value()),
		V:        (*hexutil.Big)(signature.V),
		R:        (*hexutil.Big)(signature.R),
		S:        (*hexutil.Big)(signature.S),
	}

	if blockHash != (common.Hash{}) {
		result.BlockHash = &blockHash
		result.BlockNumber = (*hexutil.Big)(new(big.Int).SetUint64(blockNumber))
		result.TransactionIndex = (*hexutil.Uint64)(&index)
	}

	switch typeInt {
	case types.TxTypeEthereumAccessList:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
	case types.TxTypeEthereumDynamicFee:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		if block != nil {
			result.GasPrice = (*hexutil.Big)(tx.EffectiveGasPrice(block.Header(), config))
		} else {
			// transaction is not processed yet
			result.GasPrice = (*hexutil.Big)(tx.EffectiveGasPrice(nil, nil))
		}
	case types.TxTypeEthereumSetCode:
		al := tx.AccessList()
		result.Accesses = &al
		result.ChainID = (*hexutil.Big)(tx.ChainId())
		result.GasFeeCap = (*hexutil.Big)(tx.GasFeeCap())
		result.GasTipCap = (*hexutil.Big)(tx.GasTipCap())
		if block != nil {
			result.GasPrice = (*hexutil.Big)(tx.EffectiveGasPrice(block.Header(), config))
		} else {
			// transaction is not processed yet
			result.GasPrice = (*hexutil.Big)(tx.EffectiveGasPrice(nil, nil))
		}
		result.AuthorizationList = tx.AuthList()
	}
	return result
}

// newEthRPCPendingTransaction creates an EthRPCTransaction for pending tx.
func newEthRPCPendingTransaction(tx *types.Transaction, config *params.ChainConfig) *EthRPCTransaction {
	return newEthRPCTransaction(nil, tx, common.Hash{}, 0, 0, config)
}

// formatTxToEthTxJSON formats types.Transaction to ethTxJSON.
// Use this function for only Ethereum typed transaction.
func formatTxToEthTxJSON(tx *types.Transaction) *ethTxJSON {
	var enc ethTxJSON

	enc.Type = hexutil.Uint64(byte(tx.Type()))
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		enc.Type = hexutil.Uint64(types.TxTypeLegacyTransaction)
	}
	enc.Hash = tx.Hash()
	signature := tx.GetTxInternalData().RawSignatureValues()[0]
	// Initialize signature values when it is nil.
	if signature.V == nil {
		signature.V = new(big.Int)
	}
	if signature.R == nil {
		signature.R = new(big.Int)
	}
	if signature.S == nil {
		signature.S = new(big.Int)
	}
	nonce := tx.Nonce()
	gas := tx.Gas()
	enc.Nonce = (*hexutil.Uint64)(&nonce)
	enc.Gas = (*hexutil.Uint64)(&gas)
	enc.Value = (*hexutil.Big)(tx.Value())
	data := tx.Data()
	enc.Data = (*hexutil.Bytes)(&data)
	enc.To = tx.To()
	enc.V = (*hexutil.Big)(signature.V)
	enc.R = (*hexutil.Big)(signature.R)
	enc.S = (*hexutil.Big)(signature.S)

	switch tx.Type() {
	case types.TxTypeEthereumAccessList:
		al := tx.AccessList()
		enc.AccessList = &al
		enc.ChainID = (*hexutil.Big)(tx.ChainId())
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice())
	case types.TxTypeEthereumDynamicFee:
		al := tx.AccessList()
		enc.AccessList = &al
		enc.ChainID = (*hexutil.Big)(tx.ChainId())
		enc.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(tx.GasTipCap())
	case types.TxTypeEthereumSetCode:
		al := tx.AccessList()
		enc.AccessList = &al
		enc.ChainID = (*hexutil.Big)(tx.ChainId())
		enc.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap())
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(tx.GasTipCap())
		enc.AuthorizationList = tx.AuthList()
	default:
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice())
	}
	return &enc
}

// GetTransactionByBlockNumberAndIndex returns the transaction for the given block number and index.
func (api *EthAPI) GetTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) *EthRPCTransaction {
	block, err := api.kaiaTransactionAPI.b.BlockByNumber(ctx, blockNr)
	if err != nil {
		return nil
	}

	return newEthRPCTransactionFromBlockIndex(block, uint64(index), api.kaiaBlockChainAPI.b.ChainConfig())
}

// GetTransactionByBlockHashAndIndex returns the transaction for the given block hash and index.
func (api *EthAPI) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) *EthRPCTransaction {
	block, err := api.kaiaTransactionAPI.b.BlockByHash(ctx, blockHash)
	if err != nil || block == nil {
		return nil
	}
	return newEthRPCTransactionFromBlockIndex(block, uint64(index), api.kaiaBlockChainAPI.b.ChainConfig())
}

// GetRawTransactionByBlockNumberAndIndex returns the bytes of the transaction for the given block number and index.
func (api *EthAPI) GetRawTransactionByBlockNumberAndIndex(ctx context.Context, blockNr rpc.BlockNumber, index hexutil.Uint) hexutil.Bytes {
	rawTx := api.kaiaTransactionAPI.GetRawTransactionByBlockNumberAndIndex(ctx, blockNr, index)
	if rawTx == nil {
		return nil
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		rawTx = rawTx[1:]
	}
	return rawTx
}

// GetRawTransactionByBlockHashAndIndex returns the bytes of the transaction for the given block hash and index.
func (api *EthAPI) GetRawTransactionByBlockHashAndIndex(ctx context.Context, blockHash common.Hash, index hexutil.Uint) hexutil.Bytes {
	rawTx := api.kaiaTransactionAPI.GetRawTransactionByBlockHashAndIndex(ctx, blockHash, index)
	if rawTx == nil {
		return nil
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		rawTx = rawTx[1:]
	}
	return rawTx
}

// GetTransactionCount returns the number of transactions the given address has sent for the given block number.
func (api *EthAPI) GetTransactionCount(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Uint64, error) {
	return api.kaiaTransactionAPI.GetTransactionCount(ctx, address, blockNrOrHash)
}

// GetTransactionByHash returns the transaction for the given hash.
func (api *EthAPI) GetTransactionByHash(ctx context.Context, hash common.Hash) (*EthRPCTransaction, error) {
	txpoolAPI := api.kaiaTransactionAPI.b

	// Try to return an already finalized transaction
	if tx, blockHash, blockNumber, index := txpoolAPI.ChainDB().ReadTxAndLookupInfo(hash); tx != nil {
		block, err := txpoolAPI.BlockByHash(ctx, blockHash)
		if err != nil {
			return nil, err
		}
		if block == nil {
			return nil, errNotFoundBlock
		}
		return newEthRPCTransaction(block, tx, blockHash, blockNumber, index, api.kaiaBlockChainAPI.b.ChainConfig()), nil
	}
	// No finalized transaction, try to retrieve it from the pool
	if tx := txpoolAPI.GetPoolTransaction(hash); tx != nil {
		return newEthRPCPendingTransaction(tx, api.kaiaBlockChainAPI.b.ChainConfig()), nil
	}
	// Transaction unknown, return as such
	return nil, nil
}

// GetRawTransactionByHash returns the bytes of the transaction for the given hash.
func (api *EthAPI) GetRawTransactionByHash(ctx context.Context, hash common.Hash) (hexutil.Bytes, error) {
	rawTx, err := api.kaiaTransactionAPI.GetRawTransactionByHash(ctx, hash)
	if rawTx == nil || err != nil {
		return nil, err
	}
	if rawTx[0] == byte(types.EthereumTxTypeEnvelope) {
		return rawTx[1:], nil
	}
	return rawTx, nil
}

// GetTransactionReceipt returns the transaction receipt for the given transaction hash.
func (api *EthAPI) GetTransactionReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	txpoolAPI := api.kaiaTransactionAPI.b

	// Formats return Kaia transaction Receipt to the Ethereum Transaction Receipt.
	tx, blockHash, blockNumber, index, receipt := txpoolAPI.GetTxLookupInfoAndReceipt(ctx, hash)

	if tx == nil {
		return nil, nil
	}
	receipts := txpoolAPI.GetBlockReceipts(ctx, blockHash)
	cumulativeGasUsed := uint64(0)
	for i := uint64(0); i <= index; i++ {
		cumulativeGasUsed += receipts[i].GasUsed
	}

	// No error handling is required here.
	// Header is checked in the following newEthTransactionReceipt function
	header, _ := txpoolAPI.HeaderByHash(ctx, blockHash)

	ethTx, err := newEthTransactionReceipt(header, tx, txpoolAPI, blockHash, blockNumber, index, cumulativeGasUsed, receipt)
	if err != nil {
		return nil, err
	}
	return ethTx, nil
}

// GetBlockReceipts returns the receipts of all transactions in the block identified by number or hash.
func (api *EthAPI) GetBlockReceipts(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) ([]map[string]interface{}, error) {
	b := api.kaiaBlockChainAPI.b
	block, err := b.BlockByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, err
	}

	var (
		header            = block.Header()
		blockHash         = block.Hash()
		blockNumber       = block.NumberU64()
		txs               = block.Transactions()
		receipts          = b.GetBlockReceipts(ctx, block.Hash())
		cumulativeGasUsed = uint64(0)
		outputList        = make([]map[string]interface{}, 0, len(receipts))
	)
	if receipts.Len() != txs.Len() {
		return nil, fmt.Errorf("the size of transactions and receipts is different in the block (%s)", blockHash.String())
	}
	for index, receipt := range receipts {
		tx := txs[index]
		cumulativeGasUsed += receipt.GasUsed
		output, err := newEthTransactionReceipt(header, tx, b, blockHash, blockNumber, uint64(index), cumulativeGasUsed, receipt)
		if err != nil {
			return nil, err
		}
		outputList = append(outputList, output)
	}
	return outputList, nil
}

// newEthTransactionReceipt creates a transaction receipt in Ethereum format.
func newEthTransactionReceipt(header *types.Header, tx *types.Transaction, b Backend, blockHash common.Hash, blockNumber, index, cumulativeGasUsed uint64, receipt *types.Receipt) (map[string]interface{}, error) {
	// When an unknown transaction receipt is requested through rpc call,
	// nil is returned by Kaia API, and it is handled.
	if tx == nil || receipt == nil {
		return nil, nil
	}

	typeInt := tx.Type()
	// If tx is not Ethereum transaction, the type is converted to TxTypeLegacyTransaction.
	if !tx.IsEthereumTransaction() {
		typeInt = types.TxTypeLegacyTransaction
	}

	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   tx.Hash(),
		"transactionIndex":  hexutil.Uint64(index),
		"from":              getFrom(tx),
		"to":                resolveToField(tx),
		"gasUsed":           hexutil.Uint64(receipt.GasUsed),
		"cumulativeGasUsed": hexutil.Uint64(cumulativeGasUsed),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(byte(typeInt)),
	}

	// After Magma hard fork : return header.baseFee
	// After EthTxType hard fork : use zero baseFee to calculate effective gas price for EthereumDynamicFeeTx :
	//  return gas price of tx.
	// Before EthTxType hard fork : return gas price of tx. (typed ethereum txs are not available.)
	fields["effectiveGasPrice"] = hexutil.Uint64(tx.EffectiveGasPrice(header, b.ChainConfig()).Uint64())

	// Always use the "status" field and Ignore the "root" field.
	if receipt.Status != types.ReceiptStatusSuccessful {
		// In Ethereum, status field can have 0(=Failure) or 1(=Success) only.
		fields["status"] = hexutil.Uint(types.ReceiptStatusFailed)
	} else {
		fields["status"] = hexutil.Uint(receipt.Status)
	}

	if receipt.Logs == nil {
		fields["logs"] = [][]*types.Log{}
	}
	// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
	if receipt.ContractAddress != (common.Address{}) {
		fields["contractAddress"] = receipt.ContractAddress
	}

	return fields, nil
}

// SendTransaction creates a transaction for the given argument, sign it and submit it to the
// transaction pool.
func (api *EthAPI) SendTransaction(ctx context.Context, args EthTransactionArgs) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresses mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		api.kaiaTransactionAPI.nonceLock.LockAddr(args.from())
		defer api.kaiaTransactionAPI.nonceLock.UnlockAddr(args.from())
	}
	if err := args.setDefaults(ctx, api.kaiaTransactionAPI.b); err != nil {
		return common.Hash{}, err
	}
	tx, _ := args.toTransaction()
	signedTx, err := api.kaiaTransactionAPI.sign(args.from(), tx)
	if err != nil {
		return common.Hash{}, err
	}
	// Check if signedTx is RLP-Encodable format.
	_, err = rlp.EncodeToBytes(signedTx)
	if err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, api.kaiaTransactionAPI.b, signedTx)
}

// EthSignTransactionResult represents a RLP encoded signed transaction.
// SignTransactionResult in go-ethereum has been renamed to EthSignTransactionResult.
// SignTransactionResult is defined in go-ethereum's internal package, so SignTransactionResult is redefined here as EthSignTransactionResult.
type EthSignTransactionResult struct {
	Raw hexutil.Bytes `json:"raw"`
	Tx  *ethTxJSON    `json:"tx"`
}

// FillTransaction fills the defaults (nonce, gas, gasPrice or 1559 fields)
// on a given unsigned transaction, and returns it to the caller for further
// processing (signing + broadcast).
func (api *EthAPI) FillTransaction(ctx context.Context, args EthTransactionArgs) (*EthSignTransactionResult, error) { // Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, api.kaiaTransactionAPI.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and obtain rlp
	tx, _ := args.toTransaction()
	data, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	if tx.IsEthTypedTransaction() {
		// Return rawTx except types.TxTypeEthEnvelope: 0x78(= 1 byte)
		return &EthSignTransactionResult{data[1:], formatTxToEthTxJSON(tx)}, nil
	}
	return &EthSignTransactionResult{data, formatTxToEthTxJSON(tx)}, nil
}

// SendRawTransaction will add the signed transaction to the transaction pool.
// The sender is responsible for signing the transaction and using the correct nonce.
func (api *EthAPI) SendRawTransaction(ctx context.Context, input hexutil.Bytes) (common.Hash, error) {
	if len(input) == 0 {
		return common.Hash{}, fmt.Errorf("Empty input")
	}
	if 0 < input[0] && input[0] < 0x7f {
		inputBytes := []byte{byte(types.EthereumTxTypeEnvelope)}
		inputBytes = append(inputBytes, input...)
		return api.kaiaTransactionAPI.SendRawTransaction(ctx, inputBytes)
	}
	// legacy transaction
	return api.kaiaTransactionAPI.SendRawTransaction(ctx, input)
}

// Sign calculates an ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message).
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The account associated with addr must be unlocked.
//
// https://github.com/ethereum/wiki/wiki/JSON-RPC#eth_sign
func (api *EthAPI) Sign(addr common.Address, data hexutil.Bytes) (hexutil.Bytes, error) {
	return api.kaiaTransactionAPI.Sign(addr, data)
}

// SignTransaction will sign the given transaction with the from account.
// The node needs to have the private key of the account corresponding with
// the given from address and it needs to be unlocked.
func (api *EthAPI) SignTransaction(ctx context.Context, args EthTransactionArgs) (*EthSignTransactionResult, error) {
	b := api.kaiaTransactionAPI.b

	if args.Gas == nil {
		return nil, fmt.Errorf("gas not specified")
	}
	if args.GasPrice == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
		return nil, fmt.Errorf("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
	}
	if args.Nonce == nil {
		return nil, fmt.Errorf("nonce not specified")
	}
	if err := args.setDefaults(ctx, b); err != nil {
		return nil, err
	}
	// Before actually sign the transaction, ensure the transaction fee is reasonable.
	tx, _ := args.toTransaction()
	if err := checkTxFee(tx.GasPrice(), tx.Gas(), b.RPCTxFeeCap()); err != nil {
		return nil, err
	}
	signed, err := api.kaiaTransactionAPI.sign(args.from(), tx)
	if err != nil {
		return nil, err
	}
	data, err := signed.MarshalBinary()
	if err != nil {
		return nil, err
	}
	if tx.IsEthTypedTransaction() {
		// Return rawTx except types.TxTypeEthEnvelope: 0x78(= 1 byte)
		return &EthSignTransactionResult{data[1:], formatTxToEthTxJSON(tx)}, nil
	}
	return &EthSignTransactionResult{data, formatTxToEthTxJSON(tx)}, nil
}

// PendingTransactions returns the transactions that are in the transaction pool
// and have a from address that is one of the accounts this node manages.
func (api *EthAPI) PendingTransactions() ([]*EthRPCTransaction, error) {
	b := api.kaiaTransactionAPI.b
	pending, err := b.GetPoolTransactions()
	if err != nil {
		return nil, err
	}
	accounts := getAccountsFromWallets(b.AccountManager().Wallets())
	transactions := make([]*EthRPCTransaction, 0, len(pending))
	for _, tx := range pending {
		from := getFrom(tx)
		if _, exists := accounts[from]; exists {
			ethTx := newEthRPCPendingTransaction(tx, api.kaiaBlockChainAPI.b.ChainConfig())
			if ethTx == nil {
				return nil, nil
			}
			transactions = append(transactions, ethTx)
		}
	}
	return transactions, nil
}

// Resend accepts an existing transaction and a new gas price and limit. It will remove
// the given transaction from the pool and reinsert it with the new gas price and limit.
func (api *EthAPI) Resend(ctx context.Context, sendArgs EthTransactionArgs, gasPrice *hexutil.Big, gasLimit *hexutil.Uint64) (common.Hash, error) {
	return resend(api.kaiaTransactionAPI, ctx, &sendArgs, gasPrice, gasLimit)
}

// Accounts returns the collection of accounts this node manages.
func (api *EthAPI) Accounts() []common.Address {
	return api.kaiaAccountAPI.Accounts()
}

func RpcMarshalEthHeader(head *types.Header, engine consensus.Engine, chainConfig *params.ChainConfig, inclMiner bool) (map[string]interface{}, error) {
	var proposer common.Address
	var err error

	if head.Number.Sign() != 0 && inclMiner {
		proposer, err = engine.Author(head)
		if err != nil {
			// miner is the field Kaia should provide the correct value. It's not the field dummy value is allowed.
			logger.Error("Failed to fetch author during marshaling header", "err", err.Error())
			return nil, err
		}
	}
	result := map[string]interface{}{
		"number":     (*hexutil.Big)(head.Number),
		"hash":       head.Hash(),
		"parentHash": head.ParentHash,
		"nonce":      BlockNonce{},  // There is no block nonce concept in Kaia, so it must be empty.
		"mixHash":    common.Hash{}, // Kaia does not use mixHash, so it must be empty.
		"sha3Uncles": common.HexToHash(EmptySha3Uncles),
		"logsBloom":  head.Bloom,
		"stateRoot":  head.Root,
		"miner":      proposer,
		"difficulty": (*hexutil.Big)(head.BlockScore),
		// extraData always return empty Bytes because actual value of extraData in Kaia header cannot be used as meaningful way because
		// we cannot provide original header of Kaia and this field is used as consensus info which is encoded value of validators addresses, validators signatures, and proposer signature in Kaia.
		"extraData": hexutil.Bytes{},
		"size":      hexutil.Uint64(head.Size()),
		// No block gas limit in Kaia, instead there is computation cost limit per tx.
		"gasLimit":         hexutil.Uint64(params.UpperGasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Big(*head.Time),
		"transactionsRoot": head.TxHash,
		"receiptsRoot":     head.ReceiptHash,
	}

	if chainConfig.IsEthTxTypeForkEnabled(head.Number) {
		if head.BaseFee == nil {
			result["baseFeePerGas"] = (*hexutil.Big)(new(big.Int).SetUint64(params.ZeroBaseFee))
		} else {
			result["baseFeePerGas"] = (*hexutil.Big)(head.BaseFee)
		}
	}
	if chainConfig.IsRandaoForkEnabled(head.Number) {
		result["randomReveal"] = hexutil.Bytes(head.RandomReveal)
		result["mixHash"] = hexutil.Bytes(head.MixHash)
	}
	return result, nil
}

// rpcMarshalHeader marshal block header as Ethereum compatible format.
// It returns error when fetching Author which is block proposer is failed.
func (api *EthAPI) rpcMarshalHeader(head *types.Header, inclMiner bool) (map[string]interface{}, error) {
	return RpcMarshalEthHeader(head, api.kaiaAPI.b.Engine(), api.kaiaAPI.b.ChainConfig(), inclMiner)
}

// rpcMarshalBlock marshal block as Ethereum compatible format
func (api *EthAPI) rpcMarshalBlock(block *types.Block, inclMiner, inclTx, fullTx bool) (map[string]interface{}, error) {
	fields, err := api.rpcMarshalHeader(block.Header(), inclMiner)
	if err != nil {
		return nil, err
	}
	fields["size"] = hexutil.Uint64(block.Size())

	if inclTx {
		formatTx := func(tx *types.Transaction) (interface{}, error) {
			return tx.Hash(), nil
		}
		if fullTx {
			formatTx = func(tx *types.Transaction) (interface{}, error) {
				return newEthRPCTransactionFromBlockHash(block, tx.Hash(), api.kaiaBlockChainAPI.b.ChainConfig()), nil
			}
		}
		txs := block.Transactions()
		transactions := make([]interface{}, len(txs))
		var err error
		for i, tx := range txs {
			if transactions[i], err = formatTx(tx); err != nil {
				return nil, err
			}
		}
		fields["transactions"] = transactions
	}
	// There is no uncles in Kaia
	fields["uncles"] = []common.Hash{}

	return fields, nil
}

func EthDoCall(ctx context.Context, b Backend, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride, timeout time.Duration, globalGasCap uint64) (*blockchain.ExecutionResult, error) {
	defer func(start time.Time) { logger.Debug("Executing EVM call finished", "runtime", time.Since(start)) }(time.Now())

	state, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if state == nil || err != nil {
		return nil, err
	}
	if err := overrides.Apply(state); err != nil {
		return nil, err
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

	// header.BaseFee != nil means magma hardforked
	var baseFee *big.Int
	if header.BaseFee != nil {
		baseFee = header.BaseFee
	} else {
		baseFee = new(big.Int).SetUint64(params.ZeroBaseFee)
	}
	intrinsicGas, err := types.IntrinsicGas(args.data(), args.GetAccessList(), args.GetAuthorizationList(), args.To == nil, b.ChainConfig().Rules(header.Number))
	if err != nil {
		return nil, err
	}
	msg, err := args.ToMessage(globalGasCap, baseFee, intrinsicGas)
	if err != nil {
		return nil, err
	}

	// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
	state.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), msg.EffectiveGasPrice(header, b.ChainConfig())))

	// The intrinsicGas is checked again later in the blockchain.ApplyMessage function,
	// but we check in advance here in order to keep StateTransition.TransactionDb method as unchanged as possible
	// and to clarify error reason correctly to serve eth namespace APIs.
	// This case is handled by EthDoEstimateGas function.
	if msg.Gas() < intrinsicGas {
		return nil, fmt.Errorf("%w: msg.gas %d, want %d", blockchain.ErrIntrinsicGas, msg.Gas(), intrinsicGas)
	}
	evm, vmError, err := b.GetEVM(ctx, msg, state, header, vm.Config{ComputationCostLimit: params.OpcodeComputationCostLimitInfinite, UseConsoleLog: b.IsConsoleLogEnabled()})
	if err != nil {
		return nil, err
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
		return nil, err
	}
	// If the timer caused an abort, return an appropriate error message
	if evm.Cancelled() {
		return nil, fmt.Errorf("execution aborted (timeout = %v)", timeout)
	}
	if err != nil {
		return result, fmt.Errorf("err: %w (supplied gas %d)", err, msg.Gas())
	}
	return result, nil
}

func EthDoEstimateGas(ctx context.Context, b Backend, args EthTransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *EthStateOverride, gasCap uint64) (hexutil.Uint64, error) {
	// Use zero address if sender unspecified.
	if args.From == nil {
		args.From = new(common.Address)
	}

	var gasLimit uint64 = 0
	if args.Gas != nil {
		gasLimit = uint64(*args.Gas)
	}

	// Normalize the max fee per gas the call is willing to spend.
	var feeCap *big.Int = common.Big0
	if args.GasPrice != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return 0, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.GasPrice != nil {
		feeCap = args.GasPrice.ToInt()
	} else if args.MaxFeePerGas != nil {
		feeCap = args.MaxFeePerGas.ToInt()
	}

	state, _, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return 0, err
	}
	if err := overrides.Apply(state); err != nil {
		return 0, err
	}
	balance := state.GetBalance(*args.From) // from can't be nil

	executable := func(gas uint64) (bool, *blockchain.ExecutionResult, error) {
		args.Gas = (*hexutil.Uint64)(&gas)
		result, err := EthDoCall(ctx, b, args, blockNrOrHash, overrides, b.RPCEVMTimeout(), gasCap)
		if err != nil {
			if errors.Is(err, blockchain.ErrIntrinsicGas) || errors.Is(err, blockchain.ErrFloorDataGas) {
				return true, nil, nil // Special case, raise gas limit
			}
			// Returns error when it is not VM error (less balance or wrong nonce, etc...).
			return true, nil, err // Bail out
		}
		// If err is vmError, return vmError with returned data
		return result.Failed(), result, nil
	}

	return blockchain.DoEstimateGas(ctx, gasLimit, gasCap, args.Value.ToInt(), feeCap, balance, executable)
}

// checkTxFee is an internal function used to check whether the fee of
// the given transaction is _reasonable_(under the cap).
func checkTxFee(gasPrice *big.Int, gas uint64, cap float64) error {
	// Short circuit if there is no cap for transaction fee at all.
	if cap == 0 {
		return nil
	}
	feeEth := new(big.Float).Quo(new(big.Float).SetInt(new(big.Int).Mul(gasPrice, new(big.Int).SetUint64(gas))), new(big.Float).SetInt(big.NewInt(params.KAIA)))
	feeFloat, _ := feeEth.Float64()
	if feeFloat > cap {
		return fmt.Errorf("tx fee (%.2f KAIA) exceeds the configured cap (%.2f KAIA)", feeFloat, cap)
	}
	return nil
}

// accessListResult returns an optional accesslist
// It's the result of the `debug_createAccessList` RPC call.
// It contains an error if the transaction itself failed.
type accessListResult struct {
	Accesslist *types.AccessList `json:"accessList"`
	Error      string            `json:"error,omitempty"`
	GasUsed    hexutil.Uint64    `json:"gasUsed"`
}

func doCreateAccessList(ctx context.Context, b Backend, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (interface{}, error) {
	bNrOrHash := rpc.NewBlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
	if blockNrOrHash != nil {
		bNrOrHash = *blockNrOrHash
	}
	acl, gasUsed, vmerr, err := AccessList(ctx, b, bNrOrHash, args)
	if err != nil {
		return nil, err
	}
	result := &accessListResult{Accesslist: &acl, GasUsed: hexutil.Uint64(gasUsed)}
	if vmerr != nil {
		result.Error = vmerr.Error()
	}
	return result, nil
}

// CreateAccessList creates an EIP-2930 type AccessList for the given transaction.
// Reexec and BlockNrOrHash can be specified to create the accessList on top of a certain state.
func (api *EthAPI) CreateAccessList(ctx context.Context, args EthTransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (interface{}, error) {
	return doCreateAccessList(ctx, api.kaiaAPI.b, args, blockNrOrHash)
}

// AccessList creates an access list for the given transaction.
// If the accesslist creation fails an error is returned.
// If the transaction itself fails, an vmErr is returned.
func AccessList(ctx context.Context, b Backend, blockNrOrHash rpc.BlockNumberOrHash, args EthTransactionArgs) (acl types.AccessList, gasUsed uint64, vmErr error, err error) {
	// Retrieve the execution context
	db, header, err := b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
	if db == nil || err != nil {
		return nil, 0, nil, err
	}
	gasCap := uint64(0)
	if rpcGasCap := b.RPCGasCap(); rpcGasCap != nil {
		gasCap = rpcGasCap.Uint64()
	}

	// Retrieve the precompiles since they don't need to be added to the access list
	rules := b.ChainConfig().Rules(header.Number)
	precompiles := vm.ActivePrecompiles(rules)

	toMsg := func() (*types.Transaction, error) {
		intrinsicGas, err := types.IntrinsicGas(args.data(), nil, nil, args.To == nil, rules)
		if err != nil {
			return nil, err
		}
		return args.ToMessage(gasCap, header.BaseFee, intrinsicGas)
	}

	if args.Gas == nil {
		// Set gaslimit to maximum if the gas is not specified
		upperGasLimit := hexutil.Uint64(params.UpperGasLimit)
		args.Gas = &upperGasLimit
	}
	if msg, err := toMsg(); err == nil {
		baseFee := new(big.Int).SetUint64(params.ZeroBaseFee)
		if header.BaseFee != nil {
			baseFee = header.BaseFee
		}
		// Add gas fee to sender for estimating gasLimit/computing cost or calling a function by insufficient balance sender.
		db.AddBalance(msg.ValidatedSender(), new(big.Int).Mul(new(big.Int).SetUint64(msg.Gas()), baseFee))
	}

	// Ensure any missing fields are filled, extract the recipient and input data
	if err := args.setDefaults(ctx, b); err != nil {
		return nil, 0, nil, err
	}
	var to common.Address
	if args.To != nil {
		to = *args.To
	} else {
		to = crypto.CreateAddress(args.from(), uint64(*args.Nonce))
	}

	// Create an initial tracer
	prevTracer := vm.NewAccessListTracer(nil, args.from(), to, precompiles)
	if args.AccessList != nil {
		prevTracer = vm.NewAccessListTracer(*args.AccessList, args.from(), to, precompiles)
	}
	for {
		if err := ctx.Err(); err != nil {
			return nil, 0, nil, err
		}
		// Retrieve the current access list to expand
		accessList := prevTracer.AccessList()
		logger.Trace("Creating access list", "input", accessList)

		// Copy the original db so we don't modify it
		statedb := db.Copy()
		// Set the accesslist to the last al
		args.AccessList = &accessList
		msg, err := toMsg()
		if err != nil {
			return nil, 0, nil, err
		}

		// Apply the transaction with the access list tracer
		tracer := vm.NewAccessListTracer(accessList, args.from(), to, precompiles)
		config := vm.Config{Tracer: tracer, Debug: true}
		vmenv, _, err := b.GetEVM(ctx, msg, statedb, header, config)
		res, err := blockchain.ApplyMessage(vmenv, msg)
		if err != nil {
			tx, _ := args.toTransaction()
			return nil, 0, nil, fmt.Errorf("failed to apply transaction: %v err: %v", tx.Hash().Hex(), err)
		}
		if tracer.Equal(prevTracer) {
			return accessList, res.UsedGas, res.Unwrap(), nil
		}
		prevTracer = tracer
	}
}
