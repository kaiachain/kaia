// Modifications Copyright 2026 The Kaia Authors
// Copyright 2022 The go-ethereum Authors
// This file is part of the go-ethereum library.

package vm

import (
	"bytes"
	"encoding/json"
	"math/big"
	"sync/atomic"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
)

var (
	_ Tracer           = (*PrestateTracer)(nil)
	_ TxPrestateTracer = (*PrestateTracer)(nil)
)

// PrestateAccount captures account state for the prestate tracer.
type PrestateAccount struct {
	Balance *big.Int                    `json:"balance,omitempty"`
	Nonce   uint64                      `json:"nonce,omitempty"`
	Code    []byte                      `json:"code,omitempty"`
	Storage map[common.Hash]common.Hash `json:"storage,omitempty"`
}

func (a *PrestateAccount) MarshalJSON() ([]byte, error) {
	type accountJSON struct {
		Balance *hexutil.Big                `json:"balance,omitempty"`
		Nonce   uint64                      `json:"nonce,omitempty"`
		Code    *hexutil.Bytes              `json:"code,omitempty"`
		Storage map[common.Hash]common.Hash `json:"storage,omitempty"`
	}
	enc := accountJSON{Nonce: a.Nonce}
	if a.Balance != nil {
		enc.Balance = (*hexutil.Big)(a.Balance)
	}
	if a.Code != nil {
		code := hexutil.Bytes(a.Code)
		enc.Code = &code
	}
	if a.Storage != nil {
		enc.Storage = a.Storage
	}
	return json.Marshal(enc)
}

func (a *PrestateAccount) empty() bool {
	return (a.Balance == nil || a.Balance.Sign() == 0) && a.Nonce == 0 && len(a.Code) == 0 && len(a.Storage) == 0
}

type prestateState = map[common.Address]*PrestateAccount

// PrestateTracerConfig is parsed from the JSON tracer config.
type PrestateTracerConfig struct {
	DiffMode       bool `json:"diffMode"`
	DisableCode    bool `json:"disableCode"`
	DisableStorage bool `json:"disableStorage"`
}

// PrestateTracer records the prestate, and optionally the post-tx diff, of all
// accounts and storage slots touched during a transaction.
type PrestateTracer struct {
	env *EVM

	pre    prestateState
	post   prestateState
	config PrestateTracerConfig

	created map[common.Address]bool

	synthAddr   common.Address
	synthAmount *big.Int
	synthPrice  *big.Int
	synthRefund *big.Int
	synthFee    *big.Int

	interrupt atomic.Bool
	reason    error
}

// NewPrestateTracer constructs a PrestateTracer. cfg may be nil.
func NewPrestateTracer(cfg json.RawMessage) (*PrestateTracer, error) {
	var config PrestateTracerConfig
	if len(cfg) > 0 {
		if err := json.Unmarshal(cfg, &config); err != nil {
			return nil, err
		}
	}
	return &PrestateTracer{
		pre:     prestateState{},
		post:    prestateState{},
		config:  config,
		created: make(map[common.Address]bool),
	}, nil
}

// SetSyntheticBalance records balance that debug_traceCall added only to make
// simulation executable. The tracer subtracts it from prestate snapshots so the
// output reflects the state selected by the caller plus stateOverrides.
func (t *PrestateTracer) SetSyntheticBalance(addr common.Address, amount *big.Int, gasPrice *big.Int) {
	t.synthAddr = addr
	if amount != nil {
		t.synthAmount = new(big.Int).Set(amount)
	}
	if gasPrice != nil {
		t.synthPrice = new(big.Int).Set(gasPrice)
	}
}

func (t *PrestateTracer) CaptureTxStartPreCheck(env *EVM, from common.Address, feePayer common.Address, to common.Address, create bool, input []byte, value *big.Int) {
	t.env = env
	t.lookupAccount(from)
	t.lookupAccount(feePayer)
	t.lookupAccount(to)
	if create {
		t.created[to] = true
	}
}

func (t *PrestateTracer) CaptureTxStart(gasLimit uint64) {}

func (t *PrestateTracer) CaptureTxEnd(restGas uint64) {
	if t.env == nil {
		return
	}
	t.setSyntheticGasUsage(restGas)
	if !t.config.DiffMode {
		for addr := range t.created {
			if acc, ok := t.pre[addr]; ok && acc.empty() {
				delete(t.pre, addr)
			}
		}
		return
	}

	t.processDiffState()

	for addr := range t.created {
		if acc, ok := t.pre[addr]; ok && acc.empty() {
			delete(t.pre, addr)
		}
	}
}

func (t *PrestateTracer) CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	if t.env == nil {
		t.env = env
	}
}

func (t *PrestateTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {}

func (t *PrestateTracer) CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (t *PrestateTracer) CaptureExit(output []byte, gasUsed uint64, err error) {}

func (t *PrestateTracer) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
	if err != nil || t.interrupt.Load() {
		return
	}
	if t.env == nil {
		t.env = env
	}
	stack := scope.Stack
	stackData := stack.Data()
	stackLen := len(stackData)
	caller := scope.Contract.Address()
	switch {
	case stackLen >= 1 && (op == SLOAD || op == SSTORE):
		t.lookupAccount(caller)
		if !t.config.DisableStorage {
			slot := common.Hash(stackData[stackLen-1].Bytes32())
			t.lookupStorage(caller, slot)
		}
	case stackLen >= 1 && (op == EXTCODECOPY || op == EXTCODEHASH || op == EXTCODESIZE || op == BALANCE || op == SELFDESTRUCT):
		addr := common.Address(stackData[stackLen-1].Bytes20())
		t.lookupAccount(addr)
		if op == SELFDESTRUCT {
			t.lookupAccount(caller)
		}
	case stackLen >= 5 && (op == DELEGATECALL || op == CALL || op == STATICCALL || op == CALLCODE):
		addr := common.Address(stackData[stackLen-2].Bytes20())
		t.lookupAccount(addr)
	case op == CREATE:
		t.lookupAccount(caller)
		nonce := t.env.StateDB.GetNonce(caller)
		addr := crypto.CreateAddress(caller, nonce)
		t.lookupAccount(addr)
		t.created[addr] = true
	case stackLen >= 4 && op == CREATE2:
		t.lookupAccount(caller)
		offset := stackData[stackLen-2]
		size := stackData[stackLen-3]
		init := scope.Memory.GetCopy(int64(offset.Uint64()), int64(size.Uint64()))
		inithash := crypto.Keccak256(init)
		salt := stackData[stackLen-4]
		addr := crypto.CreateAddress2(caller, salt.Bytes32(), inithash)
		t.lookupAccount(addr)
		t.created[addr] = true
	}
}

func (t *PrestateTracer) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
}

func (t *PrestateTracer) GetResult() (json.RawMessage, error) {
	var (
		res []byte
		err error
	)
	if t.config.DiffMode {
		res, err = json.Marshal(struct {
			Pre  prestateState `json:"pre"`
			Post prestateState `json:"post"`
		}{t.pre, t.post})
	} else {
		res, err = json.Marshal(t.pre)
	}
	if err != nil {
		return nil, err
	}
	return res, t.reason
}

func (t *PrestateTracer) Stop(err error) {
	t.reason = err
	t.interrupt.Store(true)
}

func (t *PrestateTracer) processDiffState() {
	for addr, preAcc := range t.pre {
		if t.env.StateDB.HasSelfDestructed(addr) {
			continue
		}

		modified := false
		postAcc := &PrestateAccount{}
		newBalance := t.adjustPostBalance(addr, t.env.StateDB.GetBalance(addr))
		newNonce := t.env.StateDB.GetNonce(addr)
		newCode := t.env.StateDB.GetCode(addr)

		if newBalance.Cmp(preAcc.Balance) != 0 {
			modified = true
			postAcc.Balance = new(big.Int).Set(newBalance)
		}
		if newNonce != preAcc.Nonce {
			modified = true
			postAcc.Nonce = newNonce
		}
		if !t.config.DisableCode && !bytes.Equal(newCode, preAcc.Code) {
			modified = true
			postAcc.Code = common.CopyBytes(newCode)
			if postAcc.Code == nil {
				postAcc.Code = []byte{}
			}
		}

		if t.config.DisableStorage {
			if modified {
				t.post[addr] = postAcc
			} else {
				delete(t.pre, addr)
			}
			continue
		}

		for key, val := range preAcc.Storage {
			if val == (common.Hash{}) {
				delete(preAcc.Storage, key)
			}
			newVal := t.env.StateDB.GetState(addr, key)
			if val != newVal {
				modified = true
				if newVal != (common.Hash{}) {
					if postAcc.Storage == nil {
						postAcc.Storage = make(map[common.Hash]common.Hash)
					}
					postAcc.Storage[key] = newVal
				}
			} else {
				delete(preAcc.Storage, key)
			}
		}
		if len(preAcc.Storage) == 0 {
			preAcc.Storage = nil
		}

		if modified {
			t.post[addr] = postAcc
		} else {
			delete(t.pre, addr)
		}
	}
}

func (t *PrestateTracer) lookupAccount(addr common.Address) {
	if _, ok := t.pre[addr]; ok {
		return
	}
	balance := t.adjustPreBalance(addr, t.env.StateDB.GetBalance(addr))
	acc := &PrestateAccount{
		Balance: balance,
		Nonce:   t.env.StateDB.GetNonce(addr),
	}
	if !t.config.DisableStorage {
		acc.Storage = make(map[common.Hash]common.Hash)
	}
	if !t.config.DisableCode {
		acc.Code = common.CopyBytes(t.env.StateDB.GetCode(addr))
		if acc.Code == nil {
			acc.Code = []byte{}
		}
	}
	t.pre[addr] = acc
}

func (t *PrestateTracer) setSyntheticGasUsage(restGas uint64) {
	if t.synthAmount == nil || t.synthPrice == nil {
		return
	}
	t.synthRefund = new(big.Int).Mul(new(big.Int).SetUint64(restGas), t.synthPrice)
	t.synthFee = new(big.Int).Sub(new(big.Int).Set(t.synthAmount), t.synthRefund)
	if t.synthFee.Sign() < 0 {
		t.synthFee.SetInt64(0)
	}
}

func (t *PrestateTracer) adjustPreBalance(addr common.Address, stateBalance *big.Int) *big.Int {
	balance := new(big.Int).Set(stateBalance)
	if t.synthAmount != nil && addr == t.synthAddr {
		balance.Sub(balance, t.synthAmount)
	}
	return balance
}

func (t *PrestateTracer) adjustPostBalance(addr common.Address, stateBalance *big.Int) *big.Int {
	balance := new(big.Int).Set(stateBalance)
	if t.synthRefund != nil && addr == t.synthAddr {
		balance.Sub(balance, t.synthRefund)
	}
	if t.synthFee != nil {
		if recipient, ok := t.syntheticFeeRecipient(); ok && addr == recipient {
			balance.Sub(balance, t.synthFee)
		}
	}
	return balance
}

func (t *PrestateTracer) syntheticFeeRecipient() (common.Address, bool) {
	if t.env.ChainConfig().Governance != nil && t.env.ChainConfig().Governance.DeferredTxFee() {
		return common.Address{}, false
	}
	if t.env.ChainConfig().Rules(t.env.Context.BlockNumber).IsMagma {
		return t.env.Context.Rewardbase, true
	}
	return t.env.Context.Coinbase, true
}

func (t *PrestateTracer) lookupStorage(addr common.Address, key common.Hash) {
	acc, ok := t.pre[addr]
	if !ok {
		t.lookupAccount(addr)
		acc = t.pre[addr]
	}
	if acc.Storage == nil {
		acc.Storage = make(map[common.Hash]common.Hash)
	}
	if _, ok := acc.Storage[key]; ok {
		return
	}
	acc.Storage[key] = t.env.StateDB.GetState(addr, key)
}
