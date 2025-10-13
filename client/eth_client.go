// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package client

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

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

//go:generate go run github.com/fjl/gencodec -type EthHeader -field-override headerMarshaling -out gen_header_json.go
//go:generate go run ../rlp/rlpgen -type EthHeader -out gen_header_rlp.go

// EthHeader represents a block header in the Ethereum blockchain.
type EthHeader struct {
	ParentHash  common.Hash    `gencodec:"required" json:"parentHash"`
	UncleHash   common.Hash    `gencodec:"required" json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `gencodec:"required" json:"stateRoot"`
	TxHash      common.Hash    `gencodec:"required" json:"transactionsRoot"`
	ReceiptHash common.Hash    `gencodec:"required" json:"receiptsRoot"`
	Bloom       types.Bloom    `gencodec:"required" json:"logsBloom"`
	Difficulty  *big.Int       `gencodec:"required" json:"difficulty"`
	Number      *big.Int       `gencodec:"required" json:"number"`
	GasLimit    uint64         `gencodec:"required" json:"gasLimit"`
	GasUsed     uint64         `gencodec:"required" json:"gasUsed"`
	Time        uint64         `gencodec:"required" json:"timestamp"`
	Extra       []byte         `gencodec:"required" json:"extraData"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`

	// BaseFee was added by EIP-1559 and is ignored in legacy headers.
	BaseFee *big.Int `json:"baseFeePerGas" rlp:"optional"`

	// WithdrawalsHash was added by EIP-4895 and is ignored in legacy headers.
	WithdrawalsHash *common.Hash `json:"withdrawalsRoot" rlp:"optional"`

	// BlobGasUsed was added by EIP-4844 and is ignored in legacy headers.
	BlobGasUsed *uint64 `json:"blobGasUsed" rlp:"optional"`

	// ExcessBlobGas was added by EIP-4844 and is ignored in legacy headers.
	ExcessBlobGas *uint64 `json:"excessBlobGas" rlp:"optional"`

	// ParentBeaconRoot was added by EIP-4788 and is ignored in legacy headers.
	ParentBeaconRoot *common.Hash `json:"parentBeaconBlockRoot" rlp:"optional"`

	// RequestsHash was added by EIP-7685 and is ignored in legacy headers.
	RequestsHash *common.Hash `json:"requestsHash" rlp:"optional"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	Difficulty    *hexutil.Big
	Number        *hexutil.Big
	GasLimit      hexutil.Uint64
	GasUsed       hexutil.Uint64
	Time          hexutil.Uint64
	Extra         hexutil.Bytes
	BaseFee       *hexutil.Big
	Hash          common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
	BlobGasUsed   *hexutil.Uint64
	ExcessBlobGas *hexutil.Uint64
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *EthHeader) Hash() common.Hash {
	return rlpHash(h)
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

type EthClient struct {
	c       *rpc.Client
	chainID *big.Int
}

// Dial connects a client to the given URL.
func DialEth(rawurl string) (*EthClient, error) {
	return DialContextEth(context.Background(), rawurl)
}

func DialContextEth(ctx context.Context, rawurl string) (*EthClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewEthClient(c), nil
}

// NewClient creates a client that uses the given RPC client.
func NewEthClient(c *rpc.Client) *EthClient {
	return &EthClient{c, nil}
}

func (ec *EthClient) Close() {
	ec.c.Close()
}

func (ec *EthClient) SetHeader(key, value string) {
	ec.c.SetHeader(key, value)
}

// HeaderByHash returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *EthClient) HeaderByHash(ctx context.Context, hash common.Hash) (*EthHeader, error) {
	var head *EthHeader
	err := ec.c.CallContext(ctx, &head, "eth_getBlockByHash", hash, false)
	if err == nil && head == nil {
		err = kaia.NotFound
	}
	return head, err
}

// HeaderByNumber returns a block header from the current canonical chain. If number is
// nil, the latest known header is returned.
func (ec *EthClient) HeaderByNumber(ctx context.Context, number *big.Int) (*EthHeader, error) {
	var head *EthHeader
	err := ec.c.CallContext(ctx, &head, "eth_getBlockByNumber", toBlockNumArg(number), false)
	if err == nil && head == nil {
		err = kaia.NotFound
	}
	return head, err
}

// TransactionCount returns the total number of transactions in the given block.
func (ec *EthClient) TransactionCount(ctx context.Context, blockHash common.Hash) (uint, error) {
	var num hexutil.Uint
	err := ec.c.CallContext(ctx, &num, "eth_getBlockTransactionCountByHash", blockHash)
	return uint(num), err
}

// TransactionReceiptRpcOutput returns the receipt of a transaction by transaction hash as a rpc output.
func (ec *EthClient) TransactionReceiptRpcOutput(ctx context.Context, txHash common.Hash) (r map[string]interface{}, err error) {
	err = ec.c.CallContext(ctx, &r, "eth_getTransactionReceipt", txHash)
	if err == nil && r == nil {
		return nil, kaia.NotFound
	}
	return
}

// State Access

// NetworkID returns the network ID (also known as the chain ID) for this chain.
func (ec *EthClient) NetworkID(ctx context.Context) (*big.Int, error) {
	version := new(big.Int)
	var ver string
	if err := ec.c.CallContext(ctx, &ver, "net_version"); err != nil {
		return nil, err
	}
	if _, ok := version.SetString(ver, 10); !ok {
		return nil, fmt.Errorf("invalid net_version result %q", ver)
	}
	return version, nil
}

// BalanceAt returns the kei balance of the given account.
// The block number can be nil, in which case the balance is taken from the latest known block.
func (ec *EthClient) BalanceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_getBalance", account, toBlockNumArg(blockNumber))
	return (*big.Int)(&result), err
}

// StorageAt returns the value of key in the contract storage of the given account.
// The block number can be nil, in which case the value is taken from the latest known block.
func (ec *EthClient) StorageAt(ctx context.Context, account common.Address, key common.Hash, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "eth_getStorageAt", account, key, toBlockNumArg(blockNumber))
	return result, err
}

// CodeAt returns the contract code of the given account.
// The block number can be nil, in which case the code is taken from the latest known block.
func (ec *EthClient) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "eth_getCode", account, toBlockNumArg(blockNumber))
	return result, err
}

// NonceAt returns the account nonce of the given account.
// The block number can be nil, in which case the nonce is taken from the latest known block.
func (ec *EthClient) NonceAt(ctx context.Context, account common.Address, blockNumber *big.Int) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "eth_getTransactionCount", account, toBlockNumArg(blockNumber))
	return uint64(result), err
}

// Filters

// FilterLogs executes a filter query.
func (ec *EthClient) FilterLogs(ctx context.Context, q kaia.FilterQuery) ([]types.Log, error) {
	var result []types.Log
	err := ec.c.CallContext(ctx, &result, "eth_getLogs", toFilterArg(q))
	return result, err
}

// SubscribeFilterLogs subscribes to the results of a streaming filter query.
func (ec *EthClient) SubscribeFilterLogs(ctx context.Context, q kaia.FilterQuery, ch chan<- types.Log) (kaia.Subscription, error) {
	return ec.c.EthSubscribe(ctx, ch, "logs", toFilterArg(q))
}

// Pending State
// PendingBalanceAt returns the kei balance of the given account in the pending state.
func (ec *EthClient) PendingBalanceAt(ctx context.Context, account common.Address) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_getBalance", account, "pending")
	return (*big.Int)(&result), err
}

// PendingStorageAt returns the value of key in the contract storage of the given account in the pending state.
func (ec *EthClient) PendingStorageAt(ctx context.Context, account common.Address, key common.Hash) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "eth_getStorageAt", account, key, "pending")
	return result, err
}

// PendingCodeAt returns the contract code of the given account in the pending state.
func (ec *EthClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	var result hexutil.Bytes
	err := ec.c.CallContext(ctx, &result, "eth_getCode", account, "pending")
	return result, err
}

// PendingNonceAt returns the account nonce of the given account in the pending state.
// This is the nonce that should be used for the next transaction.
func (ec *EthClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var result hexutil.Uint64
	err := ec.c.CallContext(ctx, &result, "eth_getTransactionCount", account, "pending")
	return uint64(result), err
}

// PendingTransactionCount returns the total number of transactions in the pending state.
func (ec *EthClient) PendingTransactionCount(ctx context.Context) (uint, error) {
	var num hexutil.Uint
	err := ec.c.CallContext(ctx, &num, "eth_getBlockTransactionCountByNumber", "pending")
	return uint(num), err
}

// Contract Calling

// CallContract executes a message call transaction, which is directly executed in the VM
// of the node, but never mined into the blockchain.
//
// blockNumber selects the block height at which the call runs. It can be nil, in which
// case the code is taken from the latest known block. Note that state from very old
// blocks might not be available.
func (ec *EthClient) CallContract(ctx context.Context, msg kaia.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.c.CallContext(ctx, &hex, "eth_call", toCallArg(msg), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// PendingCallContract executes a message call transaction using the EVM.
// The state seen by the contract call is the pending state.
func (ec *EthClient) PendingCallContract(ctx context.Context, msg kaia.CallMsg) ([]byte, error) {
	var hex hexutil.Bytes
	err := ec.c.CallContext(ctx, &hex, "eth_call", toCallArg(msg), "pending")
	if err != nil {
		return nil, err
	}
	return hex, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (ec *EthClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var hex hexutil.Big
	if err := ec.c.CallContext(ctx, &hex, "eth_gasPrice"); err != nil {
		return nil, err
	}
	return (*big.Int)(&hex), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific transaction based on
// the latest state of the backend blockchain. There is no guarantee that this is
// the true gas limit requirement as other transactions may be added or removed by miners,
// but it should provide a basis for setting a reasonable default.
func (ec *EthClient) EstimateGas(ctx context.Context, msg kaia.CallMsg) (uint64, error) {
	var hex hexutil.Uint64
	err := ec.c.CallContext(ctx, &hex, "eth_estimateGas", toCallArg(msg))
	if err != nil {
		return 0, err
	}
	return uint64(hex), nil
}

// SendTransaction injects a signed transaction into the pending pool for execution.
//
// If the transaction was a contract creation use the TransactionReceipt method to get the
// contract address after the transaction has been mined.
func (ec *EthClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	_, err := ec.SendRawTransaction(ctx, tx)
	return err
}

// SendRawTransaction injects a signed transaction into the pending pool for execution.
//
// This function can return the transaction hash and error.
func (ec *EthClient) SendRawTransaction(ctx context.Context, tx *types.Transaction) (common.Hash, error) {
	var hex hexutil.Bytes
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return common.Hash{}, err
	}
	if data[0] == byte(types.EthereumTxTypeEnvelope) {
		data = data[1:]
	}
	if err := ec.c.CallContext(ctx, &hex, "eth_sendRawTransaction", hexutil.Encode(data)); err != nil {
		return common.Hash{}, err
	}
	hash := common.BytesToHash(hex)
	return hash, nil
}

// BlockNumber can get the latest block number.
func (ec *EthClient) BlockNumber(ctx context.Context) (*big.Int, error) {
	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_blockNumber")
	return (*big.Int)(&result), err
}

// ChainID can return the chain ID of the chain.
func (ec *EthClient) ChainID(ctx context.Context) (*big.Int, error) {
	if ec.chainID != nil {
		return ec.chainID, nil
	}

	var result hexutil.Big
	err := ec.c.CallContext(ctx, &result, "eth_chainId")
	if err == nil {
		ec.chainID = (*big.Int)(&result)
	}
	return ec.chainID, err
}

// CreateAccessList tries to create an access list for a specific transaction based on the
// current pending state of the blockchain.
func (ec *EthClient) CreateAccessList(ctx context.Context, msg kaia.CallMsg) (*types.AccessList, uint64, string, error) {
	type AccessListResult struct {
		Accesslist *types.AccessList `json:"accessList"`
		Error      string            `json:"error,omitempty"`
		GasUsed    hexutil.Uint64    `json:"gasUsed"`
	}
	var result AccessListResult
	if err := ec.c.CallContext(ctx, &result, "eth_createAccessList", toCallArg(msg)); err != nil {
		return nil, 0, "", err
	}
	return result.Accesslist, uint64(result.GasUsed), result.Error, nil
}
