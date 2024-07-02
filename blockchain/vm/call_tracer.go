// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2020 The klaytn Authors
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

package vm

import (
	"errors"
	"math/big"
	"sync/atomic"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

var _ Tracer = (*CallTracer)(nil)

//go:generate go run github.com/fjl/gencodec -type CallFrame -field-override callFrameMarshaling -out gen_callframe_json.go
type CallFrame struct {
	Type         OpCode          `json:"-"` // e.g. CALL, DELEGATECALL, CREATE
	From         common.Address  `json:"from"`
	Gas          uint64          `json:"gas"`          // gasLeft. for top-level call, tx.gasLimit
	GasUsed      uint64          `json:"gasUsed"`      // gasUsed so far. for top-level call, tx.gasLimit - gasLeft = receipt.gasUsed
	To           *common.Address `json:"to,omitempty"` // recipient address, created contract address, or nil for failed contract creation,
	Input        []byte          `json:"input"`
	Output       []byte          `json:"output,omitempty"` // result of an internal call or revert message or runtime bytecode
	Error        string          `json:"error,omitempty"`
	RevertReason string          `json:"revertReason,omitempty"` // decoded revert message in geth style.
	Reverted     *RevertedInfo   `json:"reverted,omitempty"`     // decoded revert message and reverted contract address in klaytn style.
	Calls        []CallFrame     `json:"calls,omitempty"`        // child calls
	Value        *big.Int        `json:"value,omitempty"`
}

func (f CallFrame) TypeString() string { // to satisfy gencodec
	return f.Type.String()
}

func (f CallFrame) ToInternalTxTrace() *InternalTxTrace {
	t := &InternalTxTrace{}

	t.Type = f.Type.String()
	t.From = &f.From
	t.To = f.To
	if f.Value != nil {
		t.Value = "0x" + f.Value.Text(16)
	}

	t.Gas = f.Gas
	t.GasUsed = f.GasUsed

	if len(f.Input) > 0 {
		t.Input = hexutil.Encode(f.Input)
	}
	if len(f.Output) > 0 {
		t.Output = hexutil.Encode(f.Output)
	}
	t.Error = errors.New(f.Error)

	t.Calls = make([]*InternalTxTrace, len(f.Calls))
	for i, call := range f.Calls {
		t.Calls[i] = call.ToInternalTxTrace()
	}

	t.RevertReason = f.RevertReason
	t.Reverted = f.Reverted

	return t
}

// FieldType overrides for callFrame that's used for JSON encoding
// Must rerun gencodec after modifying this struct
type callFrameMarshaling struct {
	TypeString string `json:"type"`
	Gas        hexutil.Uint64
	GasUsed    hexutil.Uint64
	Value      *hexutil.Big
	Input      hexutil.Bytes
	Output     hexutil.Bytes
}

// Populate output, error, and revert-related fields
// 1. no error: {output}
// 2. non-revert error: {to: nil if CREATE, error}
// 3. revert error without message: {to: nil if CREATE, output, error, reverted{contract}}
// 4. revert error with message: {to: nil if CREATE, output, error, reverted{contract, message}, revertReason}
func (c *CallFrame) processOutput(output []byte, err error) {
	// 1: return output
	if err == nil {
		c.Output = common.CopyBytes(output)
		return
	}

	// 2,3,4: to = nil if CREATE failed
	if c.Type == CREATE || c.Type == CREATE2 {
		c.To = nil
	}

	// 2: do not return output
	if !errors.Is(err, ErrExecutionReverted) { // non-revert error
		c.Error = err.Error()
		return
	}

	// 3,4: return output and revert info
	c.Output = common.CopyBytes(output)
	c.Error = "execution reverted"
	c.Reverted = &RevertedInfo{Contract: c.To} // 'To' was recorded when entering this call frame

	// 4: attach revert reason
	if reason, unpackErr := abi.UnpackRevert(output); unpackErr == nil {
		c.RevertReason = reason
		c.Reverted.Message = reason
	}
}

// Implements vm.Tracer interface
type CallTracer struct {
	callstack       []CallFrame
	gasLimit        uint64 // saved tx.gasLimit
	interrupt       atomic.Bool
	interruptReason error
}

func NewCallTracer() *CallTracer {
	return &CallTracer{
		callstack: make([]CallFrame, 1), // empty top-level frame
	}
}

// Transaction start
func (t *CallTracer) CaptureTxStart(gasLimit uint64) {
	t.gasLimit = gasLimit
}

// Transaction end
func (t *CallTracer) CaptureTxEnd(gasLeft uint64) {
	t.callstack[0].GasUsed = t.callstack[0].Gas - gasLeft
}

// Enter top-level call frame
func (t *CallTracer) CaptureStart(env *EVM, from common.Address, to common.Address, create bool, input []byte, gas uint64, value *big.Int) {
	toCopy := to
	t.callstack[0] = CallFrame{
		Type:  CALL,
		From:  from,
		To:    &toCopy,
		Input: common.CopyBytes(input),
		Gas:   t.gasLimit, // ignore 'gas' supplied from EVM. Use tx.gasLimit that includes intrinsic gas.
		Value: value,
	}
	if create {
		t.callstack[0].Type = CREATE
	}
}

// Exit top-level call frame
func (t *CallTracer) CaptureEnd(output []byte, gasUsed uint64, err error) {
	// gasUsed will be filled by CaptureTxEnd; just process the output
	t.callstack[0].processOutput(output, err)
}

// Enter nested call frame
func (t *CallTracer) CaptureEnter(typ OpCode, from common.Address, to common.Address, input []byte, gas uint64, value *big.Int) {
	if t.interrupt.Load() {
		return
	}

	toCopy := to
	call := CallFrame{
		Type:  typ,
		From:  from,
		Gas:   gas,
		To:    &toCopy,
		Value: value,
		Input: common.CopyBytes(input),
	}
	t.callstack = append(t.callstack, call)
}

// Exit nested call frame
func (t *CallTracer) CaptureExit(output []byte, gasUsed uint64, err error) {
	size := len(t.callstack)
	if size <= 1 { // just in case; should never happen though because CaptureExit is only called when depth > 0
		return
	}

	// process output into the currently exiting call
	call := t.callstack[size-1]
	call.GasUsed = gasUsed
	call.processOutput(output, err)

	// pop current frame
	t.callstack = t.callstack[:size-1]

	// append it to the parent frame's Calls
	t.callstack[size-2].Calls = append(t.callstack[size-2].Calls, call)
}

// Each opcode
func (t *CallTracer) CaptureState(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
}

// Fault during opcode execution
func (t *CallTracer) CaptureFault(env *EVM, pc uint64, op OpCode, gas, cost, ccLeft, ccOpcode uint64, scope *ScopeContext, depth int, err error) {
}

func (t *CallTracer) GetResult() (CallFrame, error) {
	if len(t.callstack) != 1 {
		return CallFrame{}, errors.New("incorrect number of top-level calls")
	}

	// Return with interrupt reason if any
	return t.callstack[0], t.interruptReason
}

// Stop terminates execution of the tracer at the first opportune moment.
// For CallTracer, it stops at CaptureEnter, which is the most repetitive operation.
func (t *CallTracer) Stop(err error) {
	t.interrupt.Store(true)
	t.interruptReason = err
}
