// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from core/vm/logger_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package vm

import (
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

type dummyContractRef struct {
	calledForEach bool
}

func (dummyContractRef) ReturnGas(*big.Int)          {}
func (dummyContractRef) Address() common.Address     { return common.Address{} }
func (dummyContractRef) FeePayer() common.Address    { return common.Address{} }
func (dummyContractRef) Value() *big.Int             { return new(big.Int) }
func (dummyContractRef) SetCode(common.Hash, []byte) {}
func (d *dummyContractRef) ForEachStorage(callback func(key, value common.Hash) bool) {
	d.calledForEach = true
}
func (d *dummyContractRef) SubBalance(amount *big.Int) {}
func (d *dummyContractRef) AddBalance(amount *big.Int) {}
func (d *dummyContractRef) SetBalance(*big.Int)        {}
func (d *dummyContractRef) SetNonce(uint64)            {}
func (d *dummyContractRef) Balance() *big.Int          { return new(big.Int) }

type dummyStatedb struct {
	state.StateDB
}

func (*dummyStatedb) GetRefund() uint64 { return 1337 }

func TestStoreCapture(t *testing.T) {
	var (
		env      = NewEVM(BlockContext{}, TxContext{}, &dummyStatedb{}, params.TestChainConfig, &Config{})
		logger   = NewStructLogger(nil)
		mem      = NewMemory()
		stack    = newstack()
		contract = NewContract(&dummyContractRef{}, &dummyContractRef{}, new(big.Int), 0, nil)
	)

	stack.push(uint256.NewInt(1))
	stack.push(uint256.NewInt(0))

	var index common.Hash

	logger.CaptureState(env, 0, SSTORE, 0, 0, 0, 0, &ScopeContext{Memory: mem, Stack: stack, Contract: contract}, 0, nil)
	if len(logger.changedValues[contract.Address()]) == 0 {
		t.Fatalf("expected exactly 1 changed value on address %x, got %d", contract.Address(), len(logger.changedValues[contract.Address()]))
	}
	exp := common.BigToHash(big.NewInt(1))
	if logger.changedValues[contract.Address()][index] != exp {
		t.Errorf("expected %x, got %x", exp, logger.changedValues[contract.Address()][index])
	}
}

func TestStructLoggerStopsAtServerLimit(t *testing.T) {
	var (
		env      = NewEVM(BlockContext{}, TxContext{}, &dummyStatedb{}, params.TestChainConfig, &Config{})
		logger   = NewStructLoggerWithLimits(nil, 10, 255)
		mem      = NewMemory()
		stack    = newstack()
		contract = NewContract(&dummyContractRef{}, &dummyContractRef{}, new(big.Int), 0, nil)
	)

	logger.CaptureState(env, 0, STOP, 0, 0, 0, 0, &ScopeContext{Memory: mem, Stack: stack, Contract: contract}, 0, nil)
	if !logger.LimitReached() {
		t.Fatal("expected capture limit to be reached")
	}
	if !env.Cancelled() {
		t.Fatal("expected capture limit to cancel the EVM")
	}
	if len(logger.StructLogs()) != 0 {
		t.Fatal("over-limit log was retained")
	}
}

func TestStructLoggerEncodesRPCLogsImmediately(t *testing.T) {
	var (
		env      = NewEVM(BlockContext{}, TxContext{}, &dummyStatedb{}, params.TestChainConfig, &Config{})
		logger   = NewStructLoggerWithLimits(&LogConfig{DisableMemory: true}, 10, 4096)
		mem      = NewMemory()
		stack    = newstack()
		contract = NewContract(&dummyContractRef{}, &dummyContractRef{}, new(big.Int), 0, nil)
	)

	logger.CaptureState(env, 1, STOP, 2, 3, 4, 5, &ScopeContext{Memory: mem, Stack: stack, Contract: contract}, 6, nil)
	if logger.LimitReached() {
		t.Fatal("unexpected capture limit")
	}
	if len(logger.StructLogs()) != 0 {
		t.Fatal("RPC logger retained an unencoded snapshot")
	}
	if len(logger.JSONLogs()) != 1 {
		t.Fatalf("expected one encoded log, got %d", len(logger.JSONLogs()))
	}
	want := `{"pc":1,"op":"STOP","gas":2,"gasCost":3,"depth":6,"stack":[],"storage":{},"computation":4,"computationCost":5}`
	if got := string(logger.JSONLogs()[0]); got != want {
		t.Fatalf("encoded log mismatch\nwant: %s\n got: %s", want, got)
	}
}

func TestStructLoggerKeepsLegacySnapshots(t *testing.T) {
	var (
		env      = NewEVM(BlockContext{}, TxContext{}, &dummyStatedb{}, params.TestChainConfig, &Config{})
		logger   = NewStructLogger(&LogConfig{DisableMemory: true, DisableStack: true, DisableStorage: true})
		mem      = NewMemory()
		stack    = newstack()
		contract = NewContract(&dummyContractRef{}, &dummyContractRef{}, new(big.Int), 0, nil)
	)

	logger.CaptureState(env, 0, STOP, 0, 0, 0, 0, &ScopeContext{Memory: mem, Stack: stack, Contract: contract}, 0, nil)
	if len(logger.StructLogs()) != 1 {
		t.Fatalf("expected one legacy snapshot, got %d", len(logger.StructLogs()))
	}
	if len(logger.JSONLogs()) != 0 {
		t.Fatal("legacy logger unexpectedly encoded an RPC log")
	}
}

func TestStructLoggerBoundsCompleteResult(t *testing.T) {
	logger := NewStructLoggerWithLimits(nil, 10, 1024)
	if _, err := logger.GetResult(0, false, make([]byte, 512)); err == nil {
		t.Fatal("expected return data to count toward the result limit")
	}
}
