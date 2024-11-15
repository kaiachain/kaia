// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from tests/state_test_util.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package tests

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/kaiachain/kaia/storage/database"
)

// StateTest checks transaction processing without block context.
// See https://github.com/ethereum/EIPs/issues/176 for the test format specification.
type StateTest struct {
	json stJSON
}

// StateSubtest selects a specific configuration of a General State Test.
type StateSubtest struct {
	Fork  string
	Index int
}

func (t *StateTest) UnmarshalJSON(in []byte) error {
	return json.Unmarshal(in, &t.json)
}

type stJSON struct {
	Env  stEnv                    `json:"env"`
	Pre  blockchain.GenesisAlloc  `json:"pre"`
	Tx   stTransaction            `json:"transaction"`
	Out  hexutil.Bytes            `json:"out"`
	Post map[string][]stPostState `json:"post"`
}

type stPostState struct {
	Root    common.UnprefixedHash `json:"hash"`
	Logs    common.UnprefixedHash `json:"logs"`
	Indexes struct {
		Data  int `json:"data"`
		Gas   int `json:"gas"`
		Value int `json:"value"`
	}
}

//go:generate gencodec -type stEnv -field-override stEnvMarshaling -out gen_stenv.go

type stEnv struct {
	Coinbase   common.Address `gencodec:"required" json:"currentCoinbase"`
	BlockScore *big.Int       `gencodec:"required" json:"currentDifficulty"`
	GasLimit   uint64         `gencodec:"required" json:"currentGasLimit"`
	Number     uint64         `gencodec:"required" json:"currentNumber"`
	Timestamp  uint64         `gencodec:"required" json:"currentTimestamp"`
	BaseFee    *big.Int       `gencodec:"optional" json:"currentBaseFee"`
}

type stEnvMarshaling struct {
	Coinbase   common.UnprefixedAddress
	BlockScore *math.HexOrDecimal256
	GasLimit   math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	BaseFee    *math.HexOrDecimal256
}

//go:generate gencodec -type stTransaction -field-override stTransactionMarshaling -out gen_sttransaction.go

type stTransaction struct {
	GasPrice             *big.Int `json:"gasPrice"`
	MaxFeePerGas         *big.Int `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int `json:"maxPriorityFeePerGas"`
	Nonce                uint64   `json:"nonce"`
	To                   string   `json:"to"`
	Data                 []string `json:"data"`
	GasLimit             []uint64 `json:"gasLimit"`
	Value                []string `json:"value"`
	PrivateKey           []byte   `json:"secretKey"`
}

type stTransactionMarshaling struct {
	GasPrice            *math.HexOrDecimal256
	MaxFeePerGas        *math.HexOrDecimal256
	MaxPriorityFeePerGa *math.HexOrDecimal256
	Nonce               math.HexOrDecimal64
	GasLimit            []math.HexOrDecimal64
	PrivateKey          hexutil.Bytes
}

// getVMConfig takes a fork definition and returns a chain config.
// The fork definition can be
// - a plain forkname, e.g. `Byzantium`,
// - a fork basename, and a list of EIPs to enable; e.g. `Byzantium+1884+1283`.
func getVMConfig(forkString string) (baseConfig *params.ChainConfig, eips []int, err error) {
	var (
		splitForks            = strings.Split(forkString, "+")
		ok                    bool
		baseName, eipsStrings = splitForks[0], splitForks[1:]
	)
	if baseConfig, ok = Forks[baseName]; !ok {
		return nil, nil, UnsupportedForkError{baseName}
	}
	for _, eip := range eipsStrings {
		if eipNum, err := strconv.Atoi(eip); err != nil {
			return nil, nil, fmt.Errorf("syntax error, invalid eip number %v", eipNum)
		} else {
			eips = append(eips, eipNum)
		}
	}
	return baseConfig, eips, nil
}

// Subtests returns all valid subtests of the test.
func (t *StateTest) Subtests() []StateSubtest {
	var sub []StateSubtest
	for fork, pss := range t.json.Post {
		for i := range pss {
			sub = append(sub, StateSubtest{fork, i})
		}
	}
	return sub
}

// Run executes a specific subtest.
func (t *StateTest) Run(subtest StateSubtest, vmconfig vm.Config, isTestExecutionSpecState bool) (*state.StateDB, error) {
	config, eips, err := getVMConfig(subtest.Fork)
	if err != nil {
		return nil, UnsupportedForkError{subtest.Fork}
	}

	if isTestExecutionSpecState {
		config.Governance = &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				DeferredTxFee: true,
			},
		}
	}

	vmconfig.ExtraEips = eips
	blockchain.InitDeriveSha(config)
	block := t.genesis(config).ToBlock(common.Hash{}, nil)
	memDBManager := database.NewMemoryDBManager()
	statedb := MakePreState(memDBManager, t.json.Pre)

	post := t.json.Post[subtest.Fork][subtest.Index]
	msg, err := t.json.Tx.toMessage(post, config.Rules(block.Number()), isTestExecutionSpecState)
	if err != nil {
		return nil, err
	}
	txContext := blockchain.NewEVMTxContext(msg, block.Header(), config)
	if isTestExecutionSpecState {
		txContext.GasPrice, err = useEthGasPrice(config, &t.json)
		if err != nil {
			return nil, err
		}
	}
	blockContext := blockchain.NewEVMBlockContext(block.Header(), nil, &t.json.Env.Coinbase)
	blockContext.GetHash = vmTestBlockHash
	evm := vm.NewEVM(blockContext, txContext, statedb, config, &vmconfig)

	snapshot := statedb.Snapshot()
	result, err := blockchain.ApplyMessage(evm, msg)
	if err != nil {
		statedb.RevertToSnapshot(snapshot)
	}

	if isTestExecutionSpecState && err == nil {
		useEthMiningReward(statedb, evm, &t.json.Tx, t.json.Env.BaseFee, result.UsedGas, txContext.GasPrice)
	}

	if logs := rlpHash(statedb.Logs()); logs != common.Hash(post.Logs) {
		return statedb, fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, post.Logs)
	}

	root, _ := statedb.Commit(true)
	// Add 0-value mining reward. This only makes a difference in the cases
	// where
	// - the coinbase self-destructed, or
	// - there are only 'bad' transactions, which aren't executed. In those cases,
	//   the coinbase gets no txfee, so isn't created, and thus needs to be touched
	statedb.AddBalance(block.Rewardbase(), new(big.Int))
	// And _now_ get the state root
	root = statedb.IntermediateRoot(true)

	if isTestExecutionSpecState {
		root, err = useEthStateRoot(statedb)
		if err != nil {
			return nil, err
		}
	}
	if root != common.Hash(post.Root) {
		return statedb, fmt.Errorf("post state root mismatch: got %x, want %x", root, post.Root)
	}
	return statedb, nil
}

func (t *StateTest) gasLimit(subtest StateSubtest) uint64 {
	return t.json.Tx.GasLimit[t.json.Post[subtest.Fork][subtest.Index].Indexes.Gas]
}

func MakePreState(db database.DBManager, accounts blockchain.GenesisAlloc) *state.StateDB {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb, nil, nil)
	for addr, a := range accounts {
		if len(a.Code) != 0 {
			statedb.SetCode(addr, a.Code)
		}
		for k, v := range a.Storage {
			statedb.SetState(addr, k, v)
		}
		statedb.SetNonce(addr, a.Nonce)
		statedb.SetBalance(addr, a.Balance)
	}
	// Commit and re-open to start with a clean state.
	root, _ := statedb.Commit(false)
	statedb, _ = state.New(root, sdb, nil, nil)
	return statedb
}

func (t *StateTest) genesis(config *params.ChainConfig) *blockchain.Genesis {
	return &blockchain.Genesis{
		Config:     config,
		BlockScore: t.json.Env.BlockScore,
		Number:     t.json.Env.Number,
		Timestamp:  t.json.Env.Timestamp,
		Alloc:      t.json.Pre,
	}
}

func (tx *stTransaction) toMessage(ps stPostState, r params.Rules, isTestExecutionSpecState bool) (blockchain.Message, error) {
	// Derive sender from private key if present.
	var from common.Address
	if len(tx.PrivateKey) > 0 {
		key, err := crypto.ToECDSA(tx.PrivateKey)
		if err != nil {
			return nil, fmt.Errorf("invalid private key: %v", err)
		}
		from = crypto.PubkeyToAddress(key.PublicKey)
	}
	// Parse recipient if present.
	var to *common.Address
	if tx.To != "" {
		to = new(common.Address)
		if err := to.UnmarshalText([]byte(tx.To)); err != nil {
			return nil, fmt.Errorf("invalid to address: %v", err)
		}
	}

	// Get values specific to this post state.
	if ps.Indexes.Data > len(tx.Data) {
		return nil, fmt.Errorf("tx data index %d out of bounds", ps.Indexes.Data)
	}
	if ps.Indexes.Value > len(tx.Value) {
		return nil, fmt.Errorf("tx value index %d out of bounds", ps.Indexes.Value)
	}
	if ps.Indexes.Gas > len(tx.GasLimit) {
		return nil, fmt.Errorf("tx gas limit index %d out of bounds", ps.Indexes.Gas)
	}
	dataHex := tx.Data[ps.Indexes.Data]
	valueHex := tx.Value[ps.Indexes.Value]
	gasLimit := tx.GasLimit[ps.Indexes.Gas]
	// Value, Data hex encoding is messy: https://github.com/ethereum/tests/issues/203
	value := new(big.Int)
	if valueHex != "0x" {
		v, ok := math.ParseBig256(valueHex)
		if !ok {
			return nil, fmt.Errorf("invalid tx value %q", valueHex)
		}
		value = v
	}
	data, err := hex.DecodeString(strings.TrimPrefix(dataHex, "0x"))
	if err != nil {
		return nil, fmt.Errorf("invalid tx data %q", dataHex)
	}

	var intrinsicGas uint64
	if isTestExecutionSpecState {
		intrinsicGas, err = useEthIntrinsicGas(data, to == nil, r)
		if err != nil {
			return nil, err
		}
	} else {
		intrinsicGas, err = types.IntrinsicGas(data, nil, to == nil, r)
		if err != nil {
			return nil, err
		}
	}

	msg := types.NewMessage(from, to, tx.Nonce, value, gasLimit, tx.GasPrice, data, true, intrinsicGas, nil)
	return msg, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// simulate gas price for Ethereum
func useEthGasPrice(config *params.ChainConfig, json *stJSON) (*big.Int, error) {
	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L241-L249
	var baseFee *big.Int
	if config.IsLondonForkEnabled(new(big.Int)) {
		baseFee = json.Env.BaseFee
		if baseFee == nil {
			// Retesteth uses `0x10` for genesis baseFee. Therefore, it defaults to
			// parent - 2 : 0xa as the basefee for 'this' context.
			baseFee = big.NewInt(0x0a)
		}
	}

	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L402-L416
	gasPrice := json.Tx.GasPrice
	if baseFee != nil {
		if json.Tx.MaxFeePerGas == nil {
			json.Tx.MaxFeePerGas = gasPrice
		}
		if json.Tx.MaxFeePerGas == nil {
			json.Tx.MaxFeePerGas = new(big.Int)
		}
		if json.Tx.MaxPriorityFeePerGas == nil {
			json.Tx.MaxPriorityFeePerGas = json.Tx.MaxFeePerGas
		}
		gasPrice = math.BigMin(new(big.Int).Add(json.Tx.MaxPriorityFeePerGas, baseFee), json.Tx.MaxFeePerGas)
	}

	if gasPrice == nil {
		return nil, errors.New("no gas price provided")
	}

	return gasPrice, nil
}

// simulate intrinsic gas amount for Ethereum
func useEthIntrinsicGas(data []byte, contractCreation bool, r params.Rules) (uint64, error) {
	if r.IsIstanbul {
		r.IsPrague = true
	}
	return types.IntrinsicGas(data, nil, contractCreation, r)
}

// simulate mining reward for Ethereum
func useEthMiningReward(statedb *state.StateDB, evm *vm.EVM, tx *stTransaction, envBaseFee *big.Int, usedGas uint64, gasPrice *big.Int) {
	rules := evm.ChainConfig().Rules(evm.Context.BlockNumber)
	effectiveTip := new(big.Int).Set(gasPrice)

	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L241-L249
	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/core/state_transition.go#L462-L465
	if rules.IsLondon {
		baseFee := new(big.Int).Set(envBaseFee)
		if baseFee == nil {
			// Retesteth uses `0x10` for genesis baseFee. Therefore, it defaults to
			// parent - 2 : 0xa as the basefee for 'this' context.
			baseFee = big.NewInt(0x0a)
		}
		effectiveTip = math.BigMin(tx.MaxPriorityFeePerGas, new(big.Int).Sub(tx.MaxFeePerGas, baseFee))
	}

	fee := new(big.Int).SetUint64(usedGas)
	fee.Mul(fee, effectiveTip)
	statedb.AddBalance(evm.Context.Coinbase, fee)
}

// simulate state root for Ethereum
func useEthStateRoot(statedb *state.StateDB) (common.Hash, error) {
	memDb := database.NewMemoryDBManager()
	db := state.NewDatabase(memDb)
	newState, _ := state.New(common.Hash{}, db, nil, nil)

	for addr, acc := range statedb.RawDump().Accounts {
		b, ok := new(big.Int).SetString(acc.Balance, 10)
		if !ok {
			return common.Hash{}, errors.New("balance is not decimal")
		}
		newState.SetLegacyAccountForTest(
			common.HexToAddress(addr),
			acc.Nonce,
			b,
			common.HexToHash(acc.Root),
			common.HexToHash(acc.CodeHash).Bytes(),
		)
	}

	return newState.IntermediateRoot(true), nil
}
