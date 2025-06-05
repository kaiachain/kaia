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

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/state"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/v2/blockchain/vm"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/common/hexutil"
	"github.com/kaiachain/kaia/v2/common/math"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/crypto/sha3"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/rlp"
	"github.com/kaiachain/kaia/v2/storage/database"
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
	Root            common.UnprefixedHash `json:"hash"`
	Logs            common.UnprefixedHash `json:"logs"`
	TxBytes         hexutil.Bytes         `json:"txbytes"`
	ExpectException string                `json:"expectException"`
	Indexes         struct {
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
	GasPrice             *big.Int            `json:"gasPrice"`
	MaxFeePerGas         *big.Int            `json:"maxFeePerGas"`
	MaxPriorityFeePerGas *big.Int            `json:"maxPriorityFeePerGas"`
	Nonce                uint64              `json:"nonce"`
	To                   string              `json:"to"`
	Data                 []string            `json:"data"`
	AccessLists          []*types.AccessList `json:"accessLists,omitempty"`
	GasLimit             []uint64            `json:"gasLimit"`
	Value                []string            `json:"value"`
	PrivateKey           []byte              `json:"secretKey"`
	AuthorizationList    []*stAuthorization  `json:"authorizationList"`
}

type stTransactionMarshaling struct {
	GasPrice             *math.HexOrDecimal256
	MaxFeePerGas         *math.HexOrDecimal256
	MaxPriorityFeePerGas *math.HexOrDecimal256
	Nonce                math.HexOrDecimal64
	GasLimit             []math.HexOrDecimal64
	PrivateKey           hexutil.Bytes
}

//go:generate gencodec -type stAuthorization -field-override stAuthorizationMarshaling -out gen_stauthorization.go

// Authorization is an authorization from an account to deploy code at it's
// address.
type stAuthorization struct {
	ChainID *big.Int       `gencodec:"required" json:"chainId"`
	Address common.Address `gencodec:"required" json:"address"`
	Nonce   uint64         `gencodec:"required" json:"nonce"`
	V       uint8          `gencodec:"required" json:"v"`
	R       *big.Int       `gencodec:"required" json:"r"`
	S       *big.Int       `gencodec:"required" json:"s"`
}

// field type overrides for gencodec
type stAuthorizationMarshaling struct {
	ChainID *math.HexOrDecimal256
	Nonce   math.HexOrDecimal64
	V       math.HexOrDecimal64
	R       *math.HexOrDecimal256
	S       *math.HexOrDecimal256
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

// checkError checks if the error returned by the state transition matches any expected error.
// A failing expectation returns a wrapped version of the original error, if any,
// or a new error detailing the failing expectation.
// This function does not return or modify the original error, it only evaluates and returns expectations for the error.
func (t *StateTest) checkError(subtest StateSubtest, err error) error {
	expectedError := t.json.Post[subtest.Fork][subtest.Index].ExpectException
	if err == nil && expectedError == "" {
		return nil
	}
	if err == nil && expectedError != "" {
		return fmt.Errorf("expected error %q, got no error", expectedError)
	}
	if err != nil && expectedError == "" {
		return fmt.Errorf("unexpected error: %w", err)
	}
	if err != nil && expectedError != "" {
		// Ignore expected errors (TODO MariusVanDerWijden check error string)
		return nil
	}
	return nil
}

func (t *StateTest) Run(subtest StateSubtest, vmconfig vm.Config, isTestExecutionSpecState bool) (result error) {
	st, root, err := t.RunNoVerify(subtest, vmconfig, isTestExecutionSpecState)

	// kaia-core-tests doesn't support ExpectException
	if isTestExecutionSpecState {
		checkedErr := t.checkError(subtest, err)
		if checkedErr != nil {
			return checkedErr
		}

		// The error has been checked; if it was unexpected, it's already returned.
		if err != nil {
			// Here, an error exists but it was expected.
			// We do not check the post state or logs.
			return nil
		}
	}

	post := t.json.Post[subtest.Fork][subtest.Index]
	// N.B: We need to do this in a two-step process, because the first Commit takes care
	// of self-destructs, and we need to touch the coinbase _after_ it has potentially self-destructed.
	if root != common.Hash(post.Root) {
		return fmt.Errorf("post state root mismatch: got %x, want %x", root, post.Root)
	}
	if logs := rlpHash(st.Logs()); logs != common.Hash(post.Logs) {
		return fmt.Errorf("post state logs hash mismatch: got %x, want %x", logs, post.Logs)
	}
	st, _ = state.New(root, st.Database(), nil, nil)
	return nil
}

// RunNoVerify runs a specific subtest and returns the statedb and post-state root.
// Remember to call state.Close after verifying the test result!
func (t *StateTest) RunNoVerify(subtest StateSubtest, vmconfig vm.Config, isTestExecutionSpecState bool) (st *state.StateDB, root common.Hash, err error) {
	config, eips, err := getVMConfig(subtest.Fork)
	if err != nil {
		return st, common.Hash{}, UnsupportedForkError{subtest.Fork}
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
	rules := config.Rules(block.Number())
	st = MakePreState(memDBManager, t.json.Pre, isTestExecutionSpecState, rules)

	post := t.json.Post[subtest.Fork][subtest.Index]
	msg, err := t.json.Tx.toMessage(post, rules, isTestExecutionSpecState)
	if err != nil {
		return st, common.Hash{}, err
	}

	// Try to recover tx with current signer
	if len(post.TxBytes) != 0 {
		var ttx types.Transaction
		txBytes := post.TxBytes
		if post.TxBytes[0] <= 4 {
			txBytes = append([]byte{byte(types.EthereumTxTypeEnvelope)}, txBytes...)
		}
		err = ttx.UnmarshalBinary(txBytes)
		if err != nil {
			return st, common.Hash{}, err
		}
		if _, err := types.Sender(types.LatestSigner(config), &ttx); err != nil {
			return st, common.Hash{}, err
		}
	}

	txContext := blockchain.NewEVMTxContext(msg, block.Header(), config)
	if isTestExecutionSpecState {
		txContext.GasPrice, err = useEthGasPrice(rules, &t.json)
		if err != nil {
			return st, common.Hash{}, err
		}
	}
	blockContext := blockchain.NewEVMBlockContext(block.Header(), nil, &t.json.Env.Coinbase)
	blockContext.GetHash = vmTestBlockHash
	if isTestExecutionSpecState {
		blockContext.GasLimit = t.json.Env.GasLimit
	}
	evm := vm.NewEVM(blockContext, txContext, st, config, &vmconfig)

	if isTestExecutionSpecState {
		useEthOpCodeGas(rules, evm)
	}

	snapshot := st.Snapshot()
	result, err := blockchain.ApplyMessage(evm, msg)
	if err != nil {
		st.RevertToSnapshot(snapshot)
	}

	if err == nil && isTestExecutionSpecState {
		useEthMiningReward(st, evm.Context.Coinbase, &t.json.Tx, t.json.Env.BaseFee, result.UsedGas, txContext.GasPrice, rules)
	}

	root, _ = st.Commit(true)
	// Add 0-value mining reward. This only makes a difference in the cases
	// where
	// - the coinbase self-destructed, or
	// - there are only 'bad' transactions, which aren't executed. In those cases,
	//   the coinbase gets no txfee, so isn't created, and thus needs to be touched
	st.AddBalance(block.Rewardbase(), new(big.Int))
	// And _now_ get the state root
	root = st.IntermediateRoot(true)

	if err == nil && isTestExecutionSpecState {
		root, err = useEthState(st)
		if err != nil {
			return st, common.Hash{}, err
		}
	}
	return st, root, err
}

func (t *StateTest) gasLimit(subtest StateSubtest) uint64 {
	return t.json.Tx.GasLimit[t.json.Post[subtest.Fork][subtest.Index].Indexes.Gas]
}

func MakePreState(db database.DBManager, accounts blockchain.GenesisAlloc, isTestExecutionSpecState bool, rules params.Rules) *state.StateDB {
	sdb := state.NewDatabase(db)
	statedb, _ := state.New(common.Hash{}, sdb, nil, nil)
	for addr, a := range accounts {
		if isTestExecutionSpecState {
			if _, ok := types.ParseDelegation(a.Code); ok && rules.IsPrague {
				statedb.SetCodeToEOA(addr, a.Code, rules)
			} else if len(a.Code) == 0 && len(a.Storage) != 0 && rules.IsPrague {
				statedb.CreateEOA(addr, false, accountkey.NewAccountKeyLegacy())
			} else if len(a.Code) != 0 {
				statedb.CreateSmartContractAccount(addr, params.CodeFormatEVM, rules)
				statedb.SetCode(addr, a.Code)
			}
		} else {
			if len(a.Code) != 0 {
				statedb.SetCode(addr, a.Code)
			}
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
	if tx.To == "" {
		if tx.AuthorizationList != nil {
			// NOTE: Kaia's `newTxInternalDataEthereumSetCodeWithValues` in `MewMessage` ​​cannot be called with a "to" of "nil",
			// so specify an emptyAddress to generate a test message and test it.
			to = &common.Address{}
		}
	} else {
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

	var accessList types.AccessList
	if tx.AccessLists != nil && tx.AccessLists[ps.Indexes.Data] != nil {
		accessList = *tx.AccessLists[ps.Indexes.Data]
	}

	var authorizationList []types.SetCodeAuthorization
	if tx.AuthorizationList != nil {
		authorizationList = make([]types.SetCodeAuthorization, 0)
		for _, auth := range tx.AuthorizationList {
			authorizationList = append(authorizationList, types.SetCodeAuthorization{
				ChainID: *uint256.MustFromBig(auth.ChainID),
				Address: auth.Address,
				Nonce:   auth.Nonce,
				V:       auth.V,
				R:       *uint256.MustFromBig(auth.R),
				S:       *uint256.MustFromBig(auth.S),
			})
		}
	}

	var intrinsicGas uint64
	if isTestExecutionSpecState {
		intrinsicGas, err = useEthIntrinsicGas(data, accessList, authorizationList, to == nil, r)
	} else {
		intrinsicGas, err = types.IntrinsicGas(data, nil, nil, to == nil, r)
	}

	if err != nil {
		return nil, err
	}

	msg := types.NewMessage(from, to, tx.Nonce, value, gasLimit, tx.GasPrice, tx.MaxFeePerGas, tx.MaxPriorityFeePerGas, data, true, intrinsicGas, accessList, nil, authorizationList)
	return msg, nil
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func useEthGasPrice(r params.Rules, json *stJSON) (*big.Int, error) {
	if json.Tx.MaxFeePerGas == nil {
		json.Tx.MaxFeePerGas = json.Tx.GasPrice
	}
	if json.Tx.MaxFeePerGas == nil {
		json.Tx.MaxFeePerGas = new(big.Int)
	}
	if json.Tx.MaxPriorityFeePerGas == nil {
		json.Tx.MaxPriorityFeePerGas = json.Tx.MaxFeePerGas
	}
	return calculateEthGasPrice(r, json.Tx.GasPrice, json.Env.BaseFee, json.Tx.MaxFeePerGas, json.Tx.MaxPriorityFeePerGas)
}

func calculateEthGasPrice(r params.Rules, envGasPrice, envBaseFee, envMaxFeePerGas, envMaxPriorityFeePerGas *big.Int) (*big.Int, error) {
	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L241-L249
	var baseFee *big.Int
	if r.IsLondon {
		baseFee = envBaseFee
		if baseFee == nil {
			// Retesteth uses `0x10` for genesis baseFee. Therefore, it defaults to
			// parent - 2 : 0xa as the basefee for 'this' context.
			baseFee = big.NewInt(0x0a)
		}
	}

	// https://github.com/ethereum/go-ethereum/blob/v1.14.11/tests/state_test_util.go#L402-L416
	gasPrice := envGasPrice
	if baseFee != nil {
		gasPrice = math.BigMin(new(big.Int).Add(envMaxPriorityFeePerGas, baseFee), envMaxFeePerGas)
	}

	if gasPrice == nil {
		return nil, errors.New("no gas price provided")
	}

	return gasPrice, nil
}

func useEthOpCodeGas(r params.Rules, evm *vm.EVM) {
	if r.IsCancun {
		// EIP-1052 must be activated for backward compatibility on Kaia. But EIP-2929 is activated instead of it on Ethereum
		vm.ChangeGasCostForTest(&evm.Config.JumpTable, vm.EXTCODEHASH, params.WarmStorageReadCostEIP2929)
	}
}

func useEthIntrinsicGas(data []byte, accessList types.AccessList, authorizationList []types.SetCodeAuthorization, contractCreation bool, r params.Rules) (uint64, error) {
	if r.IsIstanbul {
		r.IsPrague = true
	}
	return types.IntrinsicGas(data, accessList, authorizationList, contractCreation, r)
}

func useEthMiningReward(statedb *state.StateDB, coinbase common.Address, tx *stTransaction, envBaseFee *big.Int, usedGas uint64, gasPrice *big.Int, rules params.Rules) {
	fee := calculateEthMiningReward(gasPrice, tx.MaxFeePerGas, tx.MaxPriorityFeePerGas, envBaseFee, usedGas, rules)
	statedb.AddBalance(coinbase, fee)
}

func calculateEthMiningReward(gasPrice, maxFeePerGas, maxPriorityFeePerGas, envBaseFee *big.Int, usedGas uint64, rules params.Rules) *big.Int {
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
		effectiveTip = math.BigMin(maxPriorityFeePerGas, new(big.Int).Sub(maxFeePerGas, baseFee))
	}

	fee := new(big.Int).SetUint64(usedGas)
	return fee.Mul(fee, effectiveTip)
}

func useEthGenesisState(statedb *state.StateDB) (common.Hash, error) {
	return useEthStateRootWithOption(statedb, false)
}

func useEthState(statedb *state.StateDB) (common.Hash, error) {
	return useEthStateRootWithOption(statedb, true)
}

func useEthStateRootWithOption(statedb *state.StateDB, deleteEmptyObjects bool) (common.Hash, error) {
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

	return newState.IntermediateRoot(deleteEmptyObjects), nil
}
