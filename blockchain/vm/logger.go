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
// This file is derived from core/vm/logger.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package vm

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/big"
	"sync/atomic"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
)

// Storage represents a contract's storage.
type Storage map[common.Hash]common.Hash

// Copy duplicates the current storage.
func (s Storage) Copy() Storage {
	cpy := make(Storage)
	maps.Copy(cpy, s)

	return cpy
}

// LogConfig are the configuration options for structured logger the EVM
type LogConfig struct {
	DisableMemory  bool // disable memory capture
	DisableStack   bool // disable stack capture
	DisableStorage bool // disable storage capture
	Debug          bool // print output during capture end
	Limit          int  // maximum length of output, but zero means unlimited
}

//go:generate gencodec -type StructLog -field-override structLogMarshaling -out gen_structlog.go

// StructLog is emitted to the EVM each cycle and lists information about the current internal state
// prior to the execution of the statement.
type StructLog struct {
	Pc              uint64                      `json:"pc"`
	Op              OpCode                      `json:"op"`
	Gas             uint64                      `json:"gas"`
	GasCost         uint64                      `json:"gasCost"`
	Memory          []byte                      `json:"memory"`
	MemorySize      int                         `json:"memSize"`
	Stack           []*big.Int                  `json:"stack"`
	Storage         map[common.Hash]common.Hash `json:"-"`
	Depth           int                         `json:"depth"`
	RefundCounter   uint64                      `json:"refund"`
	Computation     uint64                      `json:"computation"`
	ComputationCost uint64                      `json:"computationCost"`
	Err             error                       `json:"-"`
}

// structLogResult is the JSON representation returned by the tracing RPCs.
// Keeping it in the VM package lets the logger encode each entry immediately
// instead of retaining a second set of memory, stack and storage snapshots.
type structLogResult struct {
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

type structLogExecutionResult struct {
	Gas         uint64            `json:"gas"`
	Failed      bool              `json:"failed"`
	ReturnValue string            `json:"returnValue"`
	StructLogs  []json.RawMessage `json:"structLogs"`
}

// overrides for gencodec
type structLogMarshaling struct {
	Stack       []*math.HexOrDecimal256
	Gas         math.HexOrDecimal64
	GasCost     math.HexOrDecimal64
	Memory      hexutil.Bytes
	OpName      string `json:"opName"` // adds call to OpName() in MarshalJSON
	ErrorString string `json:"error"`  // adds call to ErrorString() in MarshalJSON
}

// OpName formats the operand name in a human-readable format.
func (s *StructLog) OpName() string {
	return s.Op.String()
}

// ErrorString formats the log's error as a string.
func (s *StructLog) ErrorString() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}

// Tracer is used to collect execution traces from an EVM transaction
// execution. CaptureState is called for each step of the VM with the
// current VM state.
// Note that reference types are actual VM data structures; make copies
// if you need to retain them beyond the current call.
type Tracer interface {
	// Transaction level
	CaptureTxStart(gasLimit uint64)
	CaptureTxEnd(restGas uint64)
	// Top call frame
	CaptureStart(env *EVM, from common.Address, to common.Address, call bool, input []byte, gas uint64, value *big.Int)
	CaptureEnd(output []byte, gasUsed uint64, err error)
	// Rest of call frames
	CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int)
	CaptureExit(output []byte, gasUsed uint64, err error)
	// Opcode level
	CaptureState(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error)
	CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error)
}

// TxPrestateTracer is implemented by tracers that need to snapshot top-level
// transaction state before StateTransition preCheck and subsequent execution
// mutate balances, nonces, or access lists. It is optional so existing tracers
// do not need to implement it.
type TxPrestateTracer interface {
	CaptureTxStartPreCheck(env *EVM, from common.Address, feePayer common.Address, to common.Address, create bool, input []byte, value *big.Int, authList []types.SetCodeAuthorization)
}

// StructLogger is an EVM state logger and implements Tracer.
//
// StructLogger can capture state based on the given Log configuration and also keeps
// a track record of modified storage which is used in reporting snapshots of the
// contract their storage.
type StructLogger struct {
	cfg LogConfig

	logs          []StructLog
	jsonLogs      []json.RawMessage
	changedValues map[common.Address]Storage
	output        []byte
	err           error
	encodeJSON    bool
	maxLogs       int
	maxBytes      uint64
	resultBytes   uint64
	limitReached  bool
	encodeErr     error
	interrupt     atomic.Bool
}

// NewStructLogger returns a new logger
func NewStructLogger(cfg *LogConfig) *StructLogger {
	logger := &StructLogger{changedValues: make(map[common.Address]Storage)}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

// NewStructLoggerWithLimits creates a logger with server-owned capture limits.
func NewStructLoggerWithLimits(cfg *LogConfig, maxLogs int, maxBytes uint64) *StructLogger {
	logger := &StructLogger{
		changedValues: make(map[common.Address]Storage),
		jsonLogs:      make([]json.RawMessage, 0),
		encodeJSON:    true,
		maxLogs:       maxLogs,
		maxBytes:      maxBytes,
	}
	if cfg != nil {
		logger.cfg = *cfg
	}
	return logger
}

func (l *StructLogger) stopAtLimit(env *EVM) {
	l.limitReached = true
	l.interrupt.Store(true)
	if env != nil {
		env.Cancel(CancelByCtxDone)
	}
}

func (l *StructLogger) logCount() int {
	if l.encodeJSON {
		return len(l.jsonLogs)
	}
	return len(l.logs)
}

// encodedEntryMayExceedLimit rejects an entry before allocating its formatted
// strings. The exact encoded length is checked again after marshaling.
func (l *StructLogger) encodedEntryMayExceedLimit(memoryBytes, stackItems, storageItems int) bool {
	if l.maxBytes == 0 {
		return false
	}
	if l.resultBytes >= l.maxBytes {
		return true
	}
	remaining := l.maxBytes - l.resultBytes
	const baseUpperBound uint64 = 512
	if remaining < baseUpperBound {
		return true
	}
	remaining -= baseUpperBound
	for _, field := range []struct {
		count int
		size  uint64
	}{
		{memoryBytes / 32, 67}, // 64 hex digits plus JSON delimiters
		{stackItems, 67},
		{storageItems, 134}, // key and value plus JSON delimiters
	} {
		if field.count > 0 && uint64(field.count) > remaining/field.size {
			return true
		}
		remaining -= uint64(field.count) * field.size
	}
	return false
}

func (l *StructLogger) captureJSON(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, memory *Memory, stack *Stack, contract *Contract, depth int, traceErr error) {
	memoryBytes, stackItems, storageItems := 0, 0, 0
	if !l.cfg.DisableMemory {
		memoryBytes = len(memory.Data())
	}
	if !l.cfg.DisableStack {
		stackItems = len(stack.Data())
	}
	if !l.cfg.DisableStorage {
		storageItems = len(l.changedValues[contract.Address()])
	}
	if l.encodedEntryMayExceedLimit(memoryBytes, stackItems, storageItems) {
		l.stopAtLimit(env)
		return
	}

	result := structLogResult{
		Pc:              pc,
		Op:              op.String(),
		Gas:             gas,
		GasCost:         cost,
		Depth:           depth,
		Error:           traceErr,
		Computation:     ccLeft,
		ComputationCost: ccOpcode,
	}
	if !l.cfg.DisableStack {
		formatted := make([]string, len(stack.Data()))
		for i, item := range stack.Data() {
			if i%256 == 0 && l.interrupt.Load() {
				return
			}
			word := item.Bytes32()
			formatted[i] = hex.EncodeToString(word[:])
		}
		result.Stack = &formatted
	}
	if !l.cfg.DisableMemory {
		data := memory.Data()
		formatted := make([]string, 0, len(data)/32)
		for i := 0; i+32 <= len(data); i += 32 {
			if i%(256*32) == 0 && l.interrupt.Load() {
				return
			}
			formatted = append(formatted, hex.EncodeToString(data[i:i+32]))
		}
		result.Memory = &formatted
	}
	if !l.cfg.DisableStorage {
		formatted := make(map[string]string, len(l.changedValues[contract.Address()]))
		i := 0
		for key, value := range l.changedValues[contract.Address()] {
			if i%256 == 0 && l.interrupt.Load() {
				return
			}
			formatted[hex.EncodeToString(key[:])] = hex.EncodeToString(value[:])
			i++
		}
		result.Storage = &formatted
	}
	entry, err := json.Marshal(&result)
	if err != nil {
		l.encodeErr = err
		l.interrupt.Store(true)
		env.Cancel(CancelByCtxDone)
		return
	}
	entryBytes := uint64(len(entry))
	if len(l.jsonLogs) > 0 {
		entryBytes++ // JSON array separator
	}
	if l.maxBytes > 0 && entryBytes > l.maxBytes-l.resultBytes {
		l.stopAtLimit(env)
		return
	}
	l.jsonLogs = append(l.jsonLogs, entry)
	l.resultBytes += entryBytes
}

// CaptureStart implements the Tracer interface to initialize the tracing operation.
func (l *StructLogger) CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
}

// CaptureState logs a new structured log message and pushes it out to the environment
//
// CaptureState also tracks SSTORE ops to track dirty values.
func (l *StructLogger) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
	memory := scope.Memory
	stack := scope.Stack
	contract := scope.Contract
	if l.limitReached || l.interrupt.Load() {
		return
	}
	if l.maxLogs > 0 && l.logCount() >= l.maxLogs {
		l.stopAtLimit(env)
		return
	}
	// check if already accumulated the specified number of logs
	if l.cfg.Limit != 0 && l.cfg.Limit <= l.logCount() {
		return
	}

	// initialise new changed values storage container for this contract
	// if not present.
	if l.changedValues[contract.Address()] == nil {
		l.changedValues[contract.Address()] = make(Storage)
	}

	// capture SSTORE opcodes and determine the changed value and store
	// it in the local storage container.
	if op == SSTORE && stack.len() >= 2 {
		var (
			value   = common.Hash(stack.data[stack.len()-2].Bytes32())
			address = common.Hash(stack.data[stack.len()-1].Bytes32())
		)
		l.changedValues[contract.Address()][address] = value
	}
	if l.encodeJSON {
		l.captureJSON(env, pc, op, gas, cost, ccLeft, ccOpcode, memory, stack, contract, depth, err)
		return
	}
	// Copy a snapshot of the current memory state to a new buffer
	var mem []byte
	if !l.cfg.DisableMemory {
		mem = make([]byte, len(memory.Data()))
		copy(mem, memory.Data())
	}
	// Copy a snapshot of the current stack state to a new buffer
	var stck []*big.Int
	if !l.cfg.DisableStack {
		stck = make([]*big.Int, len(stack.Data()))
		for i, item := range stack.Data() {
			stck[i] = new(big.Int).Set(item.ToBig())
		}
	}
	// Copy a snapshot of the current storage to a new container
	var storage Storage
	if !l.cfg.DisableStorage {
		storage = l.changedValues[contract.Address()].Copy()
	}
	// create a new snapshot of the EVM.
	log := StructLog{pc, op, gas, cost, mem, memory.Len(), stck, storage, depth, env.StateDB.GetRefund(), ccLeft, ccOpcode, err}

	l.logs = append(l.logs, log)
}

// CaptureFault implements the Tracer interface to trace an execution fault
// while running an opcode.
func (l *StructLogger) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
}

// CaptureEnd is called after the call finishes to finalize the tracing.
func (l *StructLogger) CaptureEnd(output []byte, gasUsed uint64, err error) {
	l.output = output
	l.err = err
	if l.cfg.Debug {
		fmt.Printf("0x%x\n", output)
		if err != nil {
			fmt.Printf(" error: %v\n", err)
		}
	}
}

func (l *StructLogger) CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
}

func (l *StructLogger) CaptureExit(output []byte, gasUsed uint64, err error) {}

func (l *StructLogger) CaptureTxStart(gasLimit uint64) {}

func (l *StructLogger) CaptureTxEnd(restGas uint64) {}

// StructLogs returns the captured log entries.
func (l *StructLogger) StructLogs() []StructLog { return l.logs }

// JSONLogs returns entries encoded in their final RPC representation.
func (l *StructLogger) JSONLogs() []json.RawMessage { return l.jsonLogs }

// GetResult encodes the complete structured trace result.
func (l *StructLogger) GetResult(gas uint64, failed bool, returnValue []byte) (json.RawMessage, error) {
	if l.maxBytes > 0 {
		const envelopeUpperBound uint64 = 512
		returnBytes := uint64(len(returnValue)) * 2
		if l.resultBytes > l.maxBytes || returnBytes > l.maxBytes-l.resultBytes || envelopeUpperBound > l.maxBytes-l.resultBytes-returnBytes {
			return nil, ErrTraceResultLimitReached
		}
	}
	result, err := json.Marshal(&structLogExecutionResult{
		Gas:         gas,
		Failed:      failed,
		ReturnValue: hex.EncodeToString(returnValue),
		StructLogs:  l.jsonLogs,
	})
	if err != nil {
		return nil, err
	}
	if l.maxBytes > 0 && uint64(len(result)) > l.maxBytes {
		return nil, ErrTraceResultLimitReached
	}
	return result, nil
}

// LimitReached reports whether a server-owned capture limit stopped the trace.
func (l *StructLogger) LimitReached() bool { return l.limitReached }

// EncodingError reports an error encountered while encoding a trace entry.
func (l *StructLogger) EncodingError() error { return l.encodeErr }

// Stop interrupts trace collection.
func (l *StructLogger) Stop() { l.interrupt.Store(true) }

// Error returns the VM error captured by the trace.
func (l *StructLogger) Error() error { return l.err }

// Output returns the VM return value captured by the trace.
func (l *StructLogger) Output() []byte { return l.output }

// WriteTrace writes a formatted trace to the given writer
func WriteTrace(writer io.Writer, logs []StructLog) {
	for _, log := range logs {
		fmt.Fprintf(writer, "%-16spc=%08d gas=%v cost=%v", log.Op, log.Pc, log.Gas, log.GasCost)
		if log.Err != nil {
			fmt.Fprintf(writer, " ERROR: %v", log.Err)
		}
		fmt.Fprintln(writer)

		if len(log.Stack) > 0 {
			fmt.Fprintln(writer, "Stack:")
			for i := len(log.Stack) - 1; i >= 0; i-- {
				fmt.Fprintf(writer, "%08d  %x\n", len(log.Stack)-i-1, math.PaddedBigBytes(log.Stack[i], 32))
			}
		}
		if len(log.Memory) > 0 {
			fmt.Fprintln(writer, "Memory:")
			fmt.Fprint(writer, hex.Dump(log.Memory))
		}
		if len(log.Storage) > 0 {
			fmt.Fprintln(writer, "Storage:")
			for h, item := range log.Storage {
				fmt.Fprintf(writer, "%x: %x\n", h, item)
			}
		}
		fmt.Fprintln(writer)
	}
}

// WriteLogs writes vm logs in a readable format to the given writer
func WriteLogs(writer io.Writer, logs []*types.Log) {
	for _, log := range logs {
		fmt.Fprintf(writer, "LOG%d: %x bn=%d txi=%x\n", len(log.Topics), log.Address, log.BlockNumber, log.TxIndex)

		for i, topic := range log.Topics {
			fmt.Fprintf(writer, "%08d  %x\n", i, topic)
		}

		fmt.Fprint(writer, hex.Dump(log.Data))
		fmt.Fprintln(writer)
	}
}
