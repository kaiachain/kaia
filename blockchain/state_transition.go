// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/state_transition.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/params"
)

var (
	errInsufficientBalanceForGas         = errors.New("insufficient balance of the sender to pay for gas")
	errInsufficientBalanceForGasFeePayer = errors.New("insufficient balance of the fee payer to pay for gas")
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==

	4a) Attempt to run transaction data
	4b) If valid, use result as code for the new state object

== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	gasTipCap  *big.Int
	gasFeeCap  *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	// ValidatedSender returns the sender of the transaction.
	// It should be set by calling AsMessageAccountKeyPicker().
	ValidatedSender() common.Address

	// ValidatedFeePayer returns the fee payer of the transaction.
	// It should be set by calling AsMessageAccountKeyPicker().
	ValidatedFeePayer() common.Address

	// ValidatedGas holds the intrinsic gas, sig validation gas, and number of data tokens for the transaction.
	// It should be set by calling AsMessageAccountKeyPicker().
	ValidatedGas() *types.ValidatedGas

	// FeeRatio returns a ratio of tx fee paid by the fee payer in percentage.
	// For example, if it is 30, 30% of tx fee will be paid by the fee payer.
	// 70% will be paid by the sender.
	FeeRatio() (types.FeeRatio, bool)

	To() *common.Address

	Hash() common.Hash

	GasPrice() *big.Int

	// For TxTypeEthereumDynamicFee
	GasTipCap() *big.Int
	GasFeeCap() *big.Int
	EffectiveGasTip(baseFee *big.Int) *big.Int
	EffectiveGasPrice(header *types.Header, config *params.ChainConfig) *big.Int

	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte

	// IntrinsicGas returns `intrinsic gas` based on the tx type.
	// This value is used to differentiate tx fee based on the tx type.
	IntrinsicGas(currentBlockNumber uint64) (uint64, error)

	// Type returns the transaction type of the message.
	Type() types.TxType

	// Validate performs additional validation for each transaction type
	Validate(stateDB types.StateDB, currentBlockNumber uint64) error

	// Execute performs execution of the transaction according to the transaction type.
	Execute(vm types.VM, stateDB types.StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) ([]byte, uint64, error)

	AccessList() types.AccessList
	AuthList() []types.SetCodeAuthorization
}

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	// Total used gas but include the refunded gas
	UsedGas uint64

	// Indicate status of transaction after execution. If the execution succeed, the status is 1.
	// If it fails, its status value indicates any error encountered during the execution (listed in blockchain/vm/errors.go)
	// This value will be stored in Receipt if Receipt is available.
	// Please see getReceiptStatusFromErrTxFailed() how the status code is derived.
	VmExecutionStatus uint

	// Returned data from evm(function result or data supplied with revert opcode)
	ReturnData []byte
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	if !result.Failed() {
		return nil
	}
	errTxFailed, ok := receiptstatus2errTxFailed[result.VmExecutionStatus]
	if !ok {
		return ErrInvalidReceiptStatus
	}
	return errTxFailed
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool {
	return result.VmExecutionStatus != types.ReceiptStatusSuccessful
}

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Failed() {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.VmExecutionStatus != types.ReceiptStatusErrExecutionReverted {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// getReceiptStatusFromErrTxFailed returns corresponding ReceiptStatus for VM error.
func getReceiptStatusFromErrTxFailed(errTxFailed error) (status uint) {
	status, ok := errTxFailed2receiptstatus[errTxFailed]
	if !ok {
		// No corresponding receiptStatus available for errTxFailed
		status = types.ReceiptStatusErrDefault
	}
	return
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message) *StateTransition {
	// before magma hardfork, effectiveGasPrice is GasPrice of tx
	// after magma hardfork, effectiveGasPrice is BaseFee
	// after kaia hardfork, effectiveGasPrice is BaseFee + effectiveGasTip
	effectiveGasPrice := evm.GasPrice

	return &StateTransition{
		evm:       evm,
		msg:       msg,
		gasPrice:  effectiveGasPrice,
		gasFeeCap: msg.GasFeeCap(),
		gasTipCap: msg.GasTipCap(),
		value:     msg.Value(),
		data:      msg.Data(),
		state:     evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message) (*ExecutionResult, error) {
	return NewStateTransition(evm, msg).TransitionDb()
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return common.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) buyGas() error {
	// st.gasPrice : gasPrice user set before magma hardfork
	// st.gasPrice : BaseFee after magma hardfork
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)

	validatedFeePayer := st.msg.ValidatedFeePayer()
	validatedSender := st.msg.ValidatedSender()
	feeRatio, isRatioTx := st.msg.FeeRatio()
	if isRatioTx {
		feePayerFee, senderFee := types.CalcFeeWithRatio(feeRatio, mgval)

		if st.state.GetBalance(validatedFeePayer).Cmp(feePayerFee) < 0 {
			logger.Debug(errInsufficientBalanceForGasFeePayer.Error(), "feePayer", validatedFeePayer.String(),
				"feePayerBalance", st.state.GetBalance(validatedFeePayer).Uint64(), "feePayerFee", feePayerFee.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGasFeePayer
		}

		if st.state.GetBalance(validatedSender).Cmp(senderFee) < 0 {
			logger.Debug(errInsufficientBalanceForGas.Error(), "sender", validatedSender.String(),
				"senderBalance", st.state.GetBalance(validatedSender).Uint64(), "senderFee", senderFee.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGas
		}

		st.state.SubBalance(validatedFeePayer, feePayerFee)
		st.state.SubBalance(validatedSender, senderFee)
	} else {
		// to make a short circuit, process the special case feeRatio == MaxFeeRatio
		if st.state.GetBalance(validatedFeePayer).Cmp(mgval) < 0 {
			logger.Debug(errInsufficientBalanceForGasFeePayer.Error(), "feePayer", validatedFeePayer.String(),
				"feePayerBalance", st.state.GetBalance(validatedFeePayer).Uint64(), "feePayerFee", mgval.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGasFeePayer
		}

		st.state.SubBalance(validatedFeePayer, mgval)
	}

	st.gas += st.msg.Gas()

	st.initialGas = st.msg.Gas()
	return nil
}

func (st *StateTransition) preCheck() error {
	// when prefetching, skip the nonce and balance check logic.
	// however, st.gas still needs to be set whether it's prefetching or not.
	if st.evm.IsPrefetching() {
		st.gas = st.msg.Gas()
		return nil
	}

	// Make sure this transaction's nonce is correct.
	if st.msg.CheckNonce() {
		stNonce := st.state.GetNonce(st.msg.ValidatedSender())
		txNonce := st.msg.Nonce()
		if stNonce < txNonce {
			logger.Debug(ErrNonceTooHigh.Error(), "account", st.msg.ValidatedSender().String(),
				"accountNonce", stNonce, "txNonce", txNonce, "txHash", st.msg.Hash().String())
			return ErrNonceTooHigh
		} else if stNonce > txNonce {
			logger.Debug(ErrNonceTooLow.Error(), "account", st.msg.ValidatedSender().String(),
				"accountNonce", stNonce, "txNonce", txNonce, "txHash", st.msg.Hash().String())
			return ErrNonceTooLow
		} else if stNonce+1 < stNonce {
			return fmt.Errorf("%w: address %v, nonce: %d", ErrNonceMax,
				st.msg.ValidatedSender().Hex(), stNonce)
		}
	}

	// Make sure that transaction gasFeeCap is greater than the baseFee (post london)
	// NOTE: Kaia adopt the validation post Prague since we're already released London
	if st.evm.ChainConfig().Rules(st.evm.Context.BlockNumber).IsPrague {
		// Skip the checks if gas fields are zero and baseFee was explicitly disabled (eth_call)
		skipCheck := st.msg.GasFeeCap().BitLen() == 0 && st.msg.GasTipCap().BitLen() == 0
		if !skipCheck {
			if l := st.msg.GasFeeCap().BitLen(); l > 256 {
				return fmt.Errorf("%w: address %v, maxFeePerGas bit length: %d", ErrFeeCapVeryHigh,
					st.msg.ValidatedSender().Hex(), l)
			}
			if l := st.msg.GasTipCap().BitLen(); l > 256 {
				return fmt.Errorf("%w: address %v, maxPriorityFeePerGas bit length: %d", ErrTipVeryHigh,
					st.msg.ValidatedSender().Hex(), l)
			}
			if st.msg.GasFeeCap().Cmp(st.msg.GasTipCap()) < 0 {
				return fmt.Errorf("%w: address %v, maxPriorityFeePerGas: %s, maxFeePerGas: %s", ErrTipAboveFeeCap,
					st.msg.ValidatedSender().Hex(), st.msg.GasTipCap(), st.msg.GasFeeCap())
			}
			// This will panic if baseFee is nil, but basefee presence is verified
			// as part of header validation.
			if st.msg.GasFeeCap().Cmp(st.evm.Context.BaseFee) < 0 {
				return fmt.Errorf("%w: address %v, maxFeePerGas: %s, baseFee: %s", ErrFeeCapTooLow,
					st.msg.ValidatedSender().Hex(), st.msg.GasFeeCap(), st.evm.Context.BaseFee)
			}
		}
	}

	// Check that EIP-7702 authorization list signatures are well formed.
	if st.msg.AuthList() != nil {
		if st.msg.To() == nil {
			return fmt.Errorf("%w (sender %v)", ErrSetCodeTxCreate, st.msg.ValidatedSender())
		}
		if len(st.msg.AuthList()) == 0 {
			return fmt.Errorf("%w (sender %v)", ErrEmptyAuthList, st.msg.ValidatedSender())
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the evm execution result with following fields.
//
//   - used gas:
//     total gas used (including gas being refunded)
//   - returndata:
//     the returned data from evm
//   - vm execution status:
//     indicates the execution result of a transaction. if the execution succeed, the status is 1.
//     if it fails, the status indicates various **EVM** errors which abort the execution.
//     e.g. ReceiptStatusErrOutOfGas, ReceiptStatusErrExecutionReverted
//
// However if any consensus issue encountered, return the error directly with
// nil evm execution result.
func (st *StateTransition) TransitionDb() (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas
	// 6. caller has enough balance to cover asset transfer for **topmost** call

	// Check clauses 1-3, buy gas if everything is correct
	if err := st.preCheck(); err != nil {
		return nil, err
	}

	var (
		msg              = st.msg
		msgTo            = msg.To()
		contractCreation = msgTo == nil
		rules            = st.evm.ChainConfig().Rules(st.evm.Context.BlockNumber)
		floorDataGas     uint64
		err              error
	)

	if st.evm.Config.Debug {
		st.evm.Config.Tracer.CaptureTxStart(st.initialGas)
		defer func() {
			st.evm.Config.Tracer.CaptureTxEnd(st.gas)
		}()
	}

	// Check clauses 4-5, subtract intrinsic gas if everything is correct
	validatedGas := msg.ValidatedGas()
	if st.gas < validatedGas.IntrinsicGas {
		return nil, ErrIntrinsicGas
	}
	if rules.IsPrague {
		floorDataGas, err = FloorDataGas(msg.Type(), msg.Data(), validatedGas.SigValidateGas)
		if err != nil {
			return nil, err
		}
		if msg.Gas() < floorDataGas {
			return nil, fmt.Errorf("%w: have %d, want %d", ErrFloorDataGas, st.gas, floorDataGas)
		}
	}
	// SigValidationGas is already inclduded in IntrinsicGas
	st.gas -= validatedGas.IntrinsicGas

	// Check clause 6
	if msg.Value().Sign() > 0 && !st.evm.Context.CanTransfer(st.state, msg.ValidatedSender(), msg.Value()) {
		return nil, vm.ErrInsufficientBalance
	}

	// Execute the preparatory steps for state transition which includes:
	// - prepare accessList(post-berlin)
	// - reset transient storage(eip 1153)
	st.state.Prepare(rules, msg.ValidatedSender(), msg.ValidatedFeePayer(), st.evm.Context.Coinbase, msgTo, vm.ActivePrecompiles(rules), msg.AccessList())

	// skip when creating a new contract
	if !contractCreation {
		// Unlike other transaction types, where the sender nonce is incremented in msg.Execute(),
		// SetCodeTx's sender nonce should be incremented before processing AuthList.
		if msg.Type() == types.TxTypeEthereumSetCode {
			// Increment the nonce for the next transaction.
			// Note: EIP-7702 authorizations can also modify the nonce. We perform
			// this update first to ensure correct validation of authorization nonces.
			st.state.IncNonce(msg.ValidatedSender())
		}

		// Apply EIP-7702 authorizations.
		if msg.AuthList() != nil {
			for _, auth := range msg.AuthList() {
				// Note errors are ignored, we simply skip invalid authorizations here.
				st.applyAuthorization(&auth, rules)
			}
		}

		if addr, ok := types.ParseDelegation(st.state.GetCode(*msgTo)); ok {
			// Perform convenience warming of sender's delegation target. Although the
			// sender is already warmed in Prepare(..), it's possible a delegation to
			// the account was deployed during this transaction. To handle correctly,
			// wait until the final state of delegations is determined before
			// performing the resolution and warming.
			st.state.AddAddressToAccessList(addr)
		}
	}

	// Check whether the init code size has been exceeded.
	if rules.IsShanghai && contractCreation && len(st.data) > params.MaxInitCodeSize {
		return nil, fmt.Errorf("%w: code size %v limit %v", ErrMaxInitCodeSizeExceeded, len(st.data), params.MaxInitCodeSize)
	}

	var (
		ret   []byte
		vmerr error
	)
	ret, st.gas, vmerr = msg.Execute(st.evm, st.state, st.evm.Context.BlockNumber.Uint64(), st.gas, st.value)

	// These tx types does not enter the EVM in the msg.Execute() method. For the purpose of debug traces,
	// those tx types are considered as a harmless zero-value transfer to sender itself. This aligns with
	// how the eth_getTransaction fills the 'to' field for those tx types. See also api_ethereum.go:resolveToField
	if st.evm.Config.Debug {
		txType := msg.Type()
		if txType.IsAccountUpdate() || txType.IsCancelTransaction() || txType.IsChainDataAnchoring() {
			st.evm.Config.Tracer.CaptureStart(st.evm, msg.ValidatedSender(), msg.ValidatedSender(), false, msg.Data(), st.initialGas, msg.Value())
			defer func() {
				st.evm.Config.Tracer.CaptureEnd(ret, st.gasUsed(), vmerr)
			}()
		}
	}

	// time-limit error is not a vm error. This error is returned when the EVM is still running while the
	// block proposer's total execution time of txs for a candidate block reached the predefined limit.
	if vmerr == vm.ErrTotalTimeLimitReached {
		return nil, vm.ErrTotalTimeLimitReached
	}

	// Compute refund
	gasRefund := st.calcRefund()
	st.gas += gasRefund
	if rules.IsPrague {
		// After EIP-7623: Data-heavy transactions pay the floor gas.
		if st.gasUsed() < floorDataGas {
			st.gas = st.initialGas - floorDataGas
		}
	}
	st.returnGas()

	// Defer transferring Tx fee when DeferredTxFee is true
	// DeferredTxFee has never been voted, so it's ok to use the genesis value instead of the latest value from governance.
	if st.evm.ChainConfig().Governance == nil || !st.evm.ChainConfig().Governance.DeferredTxFee() {
		if rules.IsMagma {
			st.state.AddBalance(st.evm.Context.Rewardbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
		} else {
			st.state.AddBalance(st.evm.Context.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice))
		}
	}

	return &ExecutionResult{
		UsedGas:           st.gasUsed(),
		VmExecutionStatus: getReceiptStatusFromErrTxFailed(vmerr), // only vm error reach here.
		ReturnData:        ret,
	}, nil
}

var errTxFailed2receiptstatus = map[error]uint{
	nil:                                             types.ReceiptStatusSuccessful,
	vm.ErrDepth:                                     types.ReceiptStatusErrDepth,
	vm.ErrContractAddressCollision:                  types.ReceiptStatusErrContractAddressCollision,
	vm.ErrCodeStoreOutOfGas:                         types.ReceiptStatusErrCodeStoreOutOfGas,
	vm.ErrMaxCodeSizeExceeded:                       types.ReceiptStatuserrMaxCodeSizeExceed,
	kerrors.ErrOutOfGas:                             types.ReceiptStatusErrOutOfGas,
	vm.ErrWriteProtection:                           types.ReceiptStatusErrWriteProtection,
	vm.ErrExecutionReverted:                         types.ReceiptStatusErrExecutionReverted,
	vm.ErrOpcodeComputationCostLimitReached:         types.ReceiptStatusErrOpcodeComputationCostLimitReached,
	kerrors.ErrAccountAlreadyExists:                 types.ReceiptStatusErrAddressAlreadyExists,
	kerrors.ErrNotProgramAccount:                    types.ReceiptStatusErrNotAProgramAccount,
	kerrors.ErrNotHumanReadableAddress:              types.ReceiptStatusErrNotHumanReadableAddress,
	kerrors.ErrFeeRatioOutOfRange:                   types.ReceiptStatusErrFeeRatioOutOfRange,
	kerrors.ErrAccountKeyFailNotUpdatable:           types.ReceiptStatusErrAccountKeyFailNotUpdatable,
	kerrors.ErrDifferentAccountKeyType:              types.ReceiptStatusErrDifferentAccountKeyType,
	kerrors.ErrAccountKeyNilUninitializable:         types.ReceiptStatusErrAccountKeyNilUninitializable,
	kerrors.ErrNotOnCurve:                           types.ReceiptStatusErrNotOnCurve,
	kerrors.ErrZeroKeyWeight:                        types.ReceiptStatusErrZeroKeyWeight,
	kerrors.ErrUnserializableKey:                    types.ReceiptStatusErrUnserializableKey,
	kerrors.ErrDuplicatedKey:                        types.ReceiptStatusErrDuplicatedKey,
	kerrors.ErrWeightedSumOverflow:                  types.ReceiptStatusErrWeightedSumOverflow,
	kerrors.ErrUnsatisfiableThreshold:               types.ReceiptStatusErrUnsatisfiableThreshold,
	kerrors.ErrZeroLength:                           types.ReceiptStatusErrZeroLength,
	kerrors.ErrLengthTooLong:                        types.ReceiptStatusErrLengthTooLong,
	kerrors.ErrNestedCompositeType:                  types.ReceiptStatusErrNestedRoleBasedKey,
	kerrors.ErrLegacyTransactionMustBeWithLegacyKey: types.ReceiptStatusErrLegacyTransactionMustBeWithLegacyKey,
	kerrors.ErrDeprecated:                           types.ReceiptStatusErrDeprecated,
	kerrors.ErrNotSupported:                         types.ReceiptStatusErrNotSupported,
	kerrors.ErrInvalidCodeFormat:                    types.ReceiptStatusErrInvalidCodeFormat,
}

var receiptstatus2errTxFailed = map[uint]error{
	types.ReceiptStatusSuccessful:                              nil,
	types.ReceiptStatusErrDefault:                              ErrVMDefault,
	types.ReceiptStatusErrDepth:                                vm.ErrDepth,
	types.ReceiptStatusErrContractAddressCollision:             vm.ErrContractAddressCollision,
	types.ReceiptStatusErrCodeStoreOutOfGas:                    vm.ErrCodeStoreOutOfGas,
	types.ReceiptStatuserrMaxCodeSizeExceed:                    vm.ErrMaxCodeSizeExceeded,
	types.ReceiptStatusErrOutOfGas:                             kerrors.ErrOutOfGas,
	types.ReceiptStatusErrWriteProtection:                      vm.ErrWriteProtection,
	types.ReceiptStatusErrExecutionReverted:                    vm.ErrExecutionReverted,
	types.ReceiptStatusErrOpcodeComputationCostLimitReached:    vm.ErrOpcodeComputationCostLimitReached,
	types.ReceiptStatusErrAddressAlreadyExists:                 kerrors.ErrAccountAlreadyExists,
	types.ReceiptStatusErrNotAProgramAccount:                   kerrors.ErrNotProgramAccount,
	types.ReceiptStatusErrNotHumanReadableAddress:              kerrors.ErrNotHumanReadableAddress,
	types.ReceiptStatusErrFeeRatioOutOfRange:                   kerrors.ErrFeeRatioOutOfRange,
	types.ReceiptStatusErrAccountKeyFailNotUpdatable:           kerrors.ErrAccountKeyFailNotUpdatable,
	types.ReceiptStatusErrDifferentAccountKeyType:              kerrors.ErrDifferentAccountKeyType,
	types.ReceiptStatusErrAccountKeyNilUninitializable:         kerrors.ErrAccountKeyNilUninitializable,
	types.ReceiptStatusErrNotOnCurve:                           kerrors.ErrNotOnCurve,
	types.ReceiptStatusErrZeroKeyWeight:                        kerrors.ErrZeroKeyWeight,
	types.ReceiptStatusErrUnserializableKey:                    kerrors.ErrUnserializableKey,
	types.ReceiptStatusErrDuplicatedKey:                        kerrors.ErrDuplicatedKey,
	types.ReceiptStatusErrWeightedSumOverflow:                  kerrors.ErrWeightedSumOverflow,
	types.ReceiptStatusErrUnsatisfiableThreshold:               kerrors.ErrUnsatisfiableThreshold,
	types.ReceiptStatusErrZeroLength:                           kerrors.ErrZeroLength,
	types.ReceiptStatusErrLengthTooLong:                        kerrors.ErrLengthTooLong,
	types.ReceiptStatusErrNestedRoleBasedKey:                   kerrors.ErrNestedCompositeType,
	types.ReceiptStatusErrLegacyTransactionMustBeWithLegacyKey: kerrors.ErrLegacyTransactionMustBeWithLegacyKey,
	types.ReceiptStatusErrDeprecated:                           kerrors.ErrDeprecated,
	types.ReceiptStatusErrNotSupported:                         kerrors.ErrNotSupported,
	types.ReceiptStatusErrInvalidCodeFormat:                    kerrors.ErrInvalidCodeFormat,
}

func (st *StateTransition) validateAuthorization(auth *types.SetCodeAuthorization) (authority common.Address, err error) {
	// Verify chain ID is 0 or equal to current chain ID.
	if !auth.ChainID.IsZero() && auth.ChainID.CmpBig(st.evm.ChainConfig().ChainID) != 0 {
		return authority, ErrAuthorizationWrongChainID
	}
	// Limit nonce to 2^64-1 per EIP-2681.
	if auth.Nonce+1 < auth.Nonce {
		return authority, ErrAuthorizationNonceOverflow
	}
	// Validate signature values and recover authority.
	authority, err = auth.Authority()
	if err != nil {
		return authority, fmt.Errorf("%w: %v", ErrAuthorizationInvalidSignature, err)
	}
	// Check the authority account
	//  1) doesn't have code or has exisiting delegation
	//  2) matches the auth's nonce
	//
	// Note it is added to the access list even if the authorization is invalid.
	st.state.AddAddressToAccessList(authority)
	code := st.state.GetCode(authority)
	if _, ok := types.ParseDelegation(code); len(code) != 0 && !ok {
		return authority, ErrAuthorizationDestinationHasCode
	}
	if have := st.state.GetNonce(authority); have != auth.Nonce {
		return authority, ErrAuthorizationNonceMismatch
	}
	return authority, nil
}

func (st *StateTransition) applyAuthorization(auth *types.SetCodeAuthorization, rules params.Rules) (err error) {
	authority, err := st.validateAuthorization(auth)
	if err != nil {
		return err
	}

	// If the account already exists in state, refund the new account cost
	// charged in the initrinsic calculation.
	if st.state.Exist(authority) {
		// If the account is not AccountKeyTypeLegacy, setcode is not allowed.
		accountKeyType := st.state.GetKey(authority).Type()
		if !accountKeyType.IsLegacyAccountKey() {
			return fmt.Errorf("%w: %v", ErrAuthorizationNotAllowAccountKeyType, accountKeyType)
		}
		st.state.AddRefund(params.CallNewAccountGas - params.TxAuthTupleGas)
	}

	// Update nonce and account code.
	st.state.IncNonce(authority)
	delegation := types.AddressToDelegation(auth.Address)
	if common.EmptyAddress(auth.Address) {
		// Delegation to zero address means clear.
		st.state.SetCodeToEOA(authority, []byte{}, rules)
		return nil
	}

	// Otherwise install delegation to auth.Address.
	// We treat EOA and SCA as separate objects and therefore need to use
	// distinct methods.
	st.state.SetCodeToEOA(authority, delegation, rules)

	return nil
}

// calcRefund computes refund counter, capped to a refund quotient.
func (st *StateTransition) calcRefund() uint64 {
	var gasRefund uint64
	if !st.evm.ChainConfig().Rules(st.evm.Context.BlockNumber).IsKore {
		// Before EIP-3529: refunds were capped to gasUsed / 2
		gasRefund = st.gasUsed() / params.RefundQuotient
	} else {
		// After EIP-3529: refunds are capped to gasUsed / 5
		gasRefund = st.gasUsed() / params.RefundQuotientEIP3529
	}
	if gasRefund > st.state.GetRefund() {
		gasRefund = st.state.GetRefund()
	}
	return gasRefund
}

// returnGas returns KAIA for remaining gas, exchanged at the original rate.
func (st *StateTransition) returnGas() {
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)

	validatedFeePayer := st.msg.ValidatedFeePayer()
	validatedSender := st.msg.ValidatedSender()
	feeRatio, isRatioTx := st.msg.FeeRatio()
	if isRatioTx {
		feePayer, feeSender := types.CalcFeeWithRatio(feeRatio, remaining)

		st.state.AddBalance(validatedFeePayer, feePayer)
		st.state.AddBalance(validatedSender, feeSender)
	} else {
		// To make a short circuit, the below routine processes when feeRatio == 100.
		st.state.AddBalance(validatedFeePayer, remaining)
	}
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}

// FloorDataGas calculates the minimum gas required for a transaction
// based on its data tokens (EIP-7623).
func FloorDataGas(txType types.TxType, data []byte, sigValidateGas uint64) (uint64, error) {
	var (
		z      = uint64(bytes.Count(data, []byte{0}))
		nz     = uint64(len(data)) - z
		tokens = nz*params.TxTokenPerNonZeroByte + z
	)
	// Check for overflow
	// Instead of using parmas.TxGas, we should consider the tx type
	// because Kaia tx type has different tx gas (e.g., fee delegated tx).
	txGas, err := types.GetTxGasForTxType(txType)
	if err != nil {
		return 0, err
	}
	if (math.MaxUint64-txGas-sigValidateGas)/params.TxCostFloorPerToken < tokens {
		return 0, types.ErrGasUintOverflow
	}
	// We add up sig validate gas too, as it's the final floor gas
	return txGas + tokens*params.TxCostFloorPerToken + sigValidateGas, nil
}
