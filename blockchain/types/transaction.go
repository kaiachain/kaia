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
// This file is derived from core/types/transaction.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package types

import (
	"bytes"
	"container/heap"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

var (
	ErrInvalidSig                     = errors.New("invalid transaction v, r, s values")
	ErrInvalidSigSender               = errors.New("invalid transaction v, r, s values of the sender")
	ErrInvalidSigFeePayer             = errors.New("invalid transaction v, r, s values of the fee payer")
	ErrInvalidTxTypeForAnchoredData   = errors.New("invalid transaction type for anchored data")
	ErrNotLegacyAccount               = errors.New("not a legacy account")
	ErrInvalidAccountKey              = errors.New("invalid account key")
	errLegacyTransaction              = errors.New("should not be called by a legacy transaction")
	errNotImplementTxInternalDataFrom = errors.New("not implement TxInternalDataFrom")
	errNotFeeDelegationTransaction    = errors.New("not a fee delegation type transaction")
	errInvalidValueMap                = errors.New("tx fields should be filled with valid values")
	errNotImplementTxInternalEthTyped = errors.New("not implement TxInternalDataEthTyped")
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(V *big.Int) Signer {
	return LatestSignerForChainID(deriveChainId(V))
}

func ErrSender(err error) error {
	return fmt.Errorf("invalid sender: %s", err)
}

func ErrFeePayer(err error) error {
	return fmt.Errorf("invalid fee payer: %s", err)
}

// ValidatedGas holds the intrinsic gas, sig validation gas, and data tokens.
//   - Intrinsic gas is the gas for the tx type + signature validation + data.
//     After Prague, floor gas would be used if intrinsic gas < floor gas.
//     Note that SigValidationGas is already included in IntrinsicGas.
//   - Sig validation gas is the gas for validating sender and feePayer.
//     It is related to Kaia-specific tx types, so it is not part of the floor gas comparison.
type ValidatedGas struct {
	IntrinsicGas   uint64
	SigValidateGas uint64
}

type Transaction struct {
	data TxInternalData
	time time.Time
	// caches
	hash         atomic.Value
	size         atomic.Value
	from         atomic.Value
	feePayer     atomic.Value
	senderTxHash atomic.Value

	// validatedSender represents the sender of the transaction to be used for ApplyTransaction().
	// This value is set in AsMessageWithAccountKeyPicker().
	validatedSender common.Address
	// validatedFeePayer represents the fee payer of the transaction to be used for ApplyTransaction().
	// This value is set in AsMessageWithAccountKeyPicker().
	validatedFeePayer common.Address
	// validatedGas holds intrinsic gas, sig validation gas, and number of tokens for the transaction to be used for ApplyTransaction().
	// This value is set in AsMessageWithAccountKeyPicker().
	validatedGas *ValidatedGas
	// The account's nonce is checked only if `checkNonce` is true.
	checkNonce bool
	// This value is set when the tx is invalidated in block tx validation, and is used to remove pending tx in txPool.
	markedUnexecutable int32

	// lock for protecting fields in Transaction struct
	mu sync.RWMutex
}

// NewTransactionWithMap generates a tx from tx field values.
// One of the return value, retErr, is lastly updated when panic is occurred.
func NewTransactionWithMap(t TxType, values map[TxValueKeyType]interface{}) (tx *Transaction, retErr error) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warn("Got panic and recovered", "panicErr", err)
			retErr = errInvalidValueMap
		}
	}()
	txData, err := NewTxInternalDataWithMap(t, values)
	if err != nil {
		return nil, err
	}
	tx = NewTx(txData)
	return tx, retErr
}

// NewTx creates a new transaction.
func NewTx(data TxInternalData) *Transaction {
	tx := new(Transaction)
	tx.setDecoded(data, 0)

	return tx
}

func NewTransaction(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, &to, amount, gasLimit, gasPrice, data)
}

func NewContractCreation(nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(nonce, nil, amount, gasLimit, gasPrice, data)
}

func newTransaction(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *Transaction {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	d := TxInternalDataLegacy{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     gasLimit,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return NewTx(&d)
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return tx.data.ChainId()
}

// SenderTxHash returns (SenderTxHash, true) if the tx is a fee-delegated transaction.
// Otherwise, it returns (nil hash, false).
func (tx *Transaction) SenderTxHash() (common.Hash, bool) {
	if tx.Type().IsFeeDelegatedTransaction() == false {
		// Do not compute SenderTxHash for non-fee-delegated txs
		return common.Hash{}, false
	}
	if senderTxHash := tx.senderTxHash.Load(); senderTxHash != nil {
		return senderTxHash.(common.Hash), tx.Type().IsFeeDelegatedTransaction()
	}
	v := tx.data.SenderTxHash()
	tx.senderTxHash.Store(v)
	return v, tx.Type().IsFeeDelegatedTransaction()
}

// SenderTxHashAll returns SenderTxHash for all tx types.
// If it is not a fee-delegated tx, SenderTxHash and TxHash are the same.
func (tx *Transaction) SenderTxHashAll() common.Hash {
	if senderTxHash := tx.senderTxHash.Load(); senderTxHash != nil {
		return senderTxHash.(common.Hash)
	}
	v := tx.data.SenderTxHash()
	tx.senderTxHash.Store(v)
	return v
}

func validateSignature(v, r, s *big.Int) bool {
	// TODO-Kaia: Need to consider the case v.BitLen() > 64.
	// Since ValidateSignatureValues receives v as type of byte, leave it as a future work.
	if v != nil && !isProtectedV(v) {
		return crypto.ValidateSignatureValues(byte(v.Uint64()-27), r, s, true)
	}

	chainID := deriveChainId(v).Uint64()
	V := byte(v.Uint64() - 35 - 2*chainID)

	return crypto.ValidateSignatureValues(V, r, s, false)
}

func (tx *Transaction) Equal(tb *Transaction) bool {
	return tx.data.Equal(tb.data)
}

// EncodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	serializer := newTxInternalDataSerializerWithValues(tx.data)
	return rlp.Encode(w, serializer)
}

// MarshalBinary returns the canonical encoding of the transaction.
// For legacy transactions, it returns the RLP encoding. For typed
// transactions, it returns the type and payload.
func (tx *Transaction) MarshalBinary() ([]byte, error) {
	var buf bytes.Buffer
	err := tx.EncodeRLP(&buf)
	return buf.Bytes(), err
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	serializer := newTxInternalDataSerializer()
	if err := s.Decode(serializer); err != nil {
		return err
	}

	if !serializer.tx.ValidateSignature() {
		return ErrInvalidSig
	}

	size := calculateTxSize(serializer.tx)
	tx.setDecoded(serializer.tx, int(size))

	return nil
}

// UnmarshalBinary decodes the canonical encoding of transactions.
// It supports legacy RLP transactions and EIP2718 typed transactions.
func (tx *Transaction) UnmarshalBinary(b []byte) error {
	newTx := &Transaction{}
	if err := rlp.DecodeBytes(b, newTx); err != nil {
		return err
	}

	tx.setDecoded(newTx.data, len(b))
	return nil
}

// MarshalJSON encodes the web3 RPC transaction format.
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.data
	data.SetHash(&hash)
	serializer := newTxInternalDataSerializerWithValues(tx.data)
	return json.Marshal(serializer)
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	serializer := newTxInternalDataSerializer()
	if err := json.Unmarshal(input, serializer); err != nil {
		return err
	}
	if !serializer.tx.ValidateSignature() {
		return ErrInvalidSig
	}

	tx.setDecoded(serializer.tx, 0)
	return nil
}

func (tx *Transaction) setDecoded(inner TxInternalData, size int) {
	tx.data = inner
	tx.time = time.Now()

	if size > 0 {
		tx.size.Store(common.StorageSize(size))
	}
}

func (tx *Transaction) Gas() uint64        { return tx.data.GetGasLimit() }
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.data.GetPrice()) }
func (tx *Transaction) GasTipCap() *big.Int {
	if tx.Type() == TxTypeEthereumDynamicFee || tx.Type() == TxTypeEthereumSetCode {
		te := tx.GetTxInternalData().(TxInternalDataBaseFee)
		return te.GetGasTipCap()
	}

	return tx.data.GetPrice()
}

func (tx *Transaction) GasFeeCap() *big.Int {
	if tx.Type() == TxTypeEthereumDynamicFee || tx.Type() == TxTypeEthereumSetCode {
		te := tx.GetTxInternalData().(TxInternalDataBaseFee)
		return te.GetGasFeeCap()
	}

	return tx.data.GetPrice()
}

func (tx *Transaction) EffectiveGasTip(baseFee *big.Int) *big.Int {
	// effectiveGasPrice - baseFee = min(baseFee + tipCap, feeCap) - baseFee = min(tipCap, feeCap - baseFee)
	if baseFee != nil {
		// For EthereumDynamicFee TxType: min(GasTipCap, Sub(GasFeeCap,baseFee))
		// For Non-EthereumDynamicFee TxType: min(GasPrice, Sub(gasPrice, baseFee)
		tip := math.BigMax(big.NewInt(0), new(big.Int).Sub(tx.GasFeeCap(), baseFee))
		return math.BigMin(tx.GasTipCap(), tip)
	}

	return new(big.Int).Set(tx.GasTipCap())
}

func (tx *Transaction) EffectiveGasPrice(header *Header, config *params.ChainConfig) *big.Int {
	if header == nil || header.BaseFee == nil {
		return new(big.Int).Set(tx.GasPrice())
	}
	if config.Rules(header.Number).IsKaia {
		tip := tx.EffectiveGasTip(header.BaseFee)
		return new(big.Int).Add(tip, header.BaseFee)
	}
	return new(big.Int).Set(header.BaseFee)
}

func (tx *Transaction) AccessList() AccessList {
	if tx.IsEthTypedTransaction() {
		te := tx.GetTxInternalData().(TxInternalDataEthTyped)
		return te.GetAccessList()
	}
	return nil
}

func (tx *Transaction) AuthList() []SetCodeAuthorization {
	if tx.Type() == TxTypeEthereumSetCode {
		te := tx.GetTxInternalData().(*TxInternalDataEthereumSetCode)
		return te.GetAuthorizationList()
	}
	return nil
}

// SetCodeAuthorities returns a list of unique authorities from the
// authorization list.
func (tx *Transaction) SetCodeAuthorities() []common.Address {
	if tx.Type() != TxTypeEthereumSetCode {
		return nil
	}
	var (
		marks = make(map[common.Address]bool)
		auths = make([]common.Address, 0, len(tx.AuthList()))
	)
	for _, auth := range tx.AuthList() {
		if addr, err := auth.Authority(); err == nil {
			if marks[addr] {
				continue
			}
			marks[addr] = true
			auths = append(auths, addr)
		}
	}
	return auths
}

func (tx *Transaction) Value() *big.Int { return new(big.Int).Set(tx.data.GetAmount()) }
func (tx *Transaction) Nonce() uint64   { return tx.data.GetAccountNonce() }
func (tx *Transaction) CheckNonce() bool {
	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return tx.checkNonce
}
func (tx *Transaction) Type() TxType                { return tx.data.Type() }
func (tx *Transaction) IsLegacyTransaction() bool   { return tx.Type().IsLegacyTransaction() }
func (tx *Transaction) IsEthTypedTransaction() bool { return tx.Type().IsEthTypedTransaction() }
func (tx *Transaction) IsEthereumTransaction() bool {
	return tx.Type().IsEthereumTransaction()
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28 && v != 1 && v != 0
	}
	// anything not 27 or 28 is considered protected
	return true
}

// Protected says whether the transaction is replay-protected.
func (tx *Transaction) Protected() bool {
	if tx.IsLegacyTransaction() {
		v := tx.RawSignatureValues()[0].V
		return v != nil && isProtectedV(v)
	}
	return true
}

func (tx *Transaction) ValidatedSender() common.Address {
	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return tx.validatedSender
}

func (tx *Transaction) ValidatedFeePayer() common.Address {
	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return tx.validatedFeePayer
}

func (tx *Transaction) ValidatedGas() *ValidatedGas {
	tx.mu.RLock()
	defer tx.mu.RUnlock()
	return tx.validatedGas
}
func (tx *Transaction) MakeRPCOutput() map[string]interface{} { return tx.data.MakeRPCOutput() }
func (tx *Transaction) GetTxInternalData() TxInternalData     { return tx.data }

func (tx *Transaction) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	return tx.data.IntrinsicGas(currentBlockNumber)
}

func (tx *Transaction) Validate(db StateDB, blockNumber uint64) error {
	return tx.data.Validate(db, blockNumber)
}

// ValidateMutableValue conducts validation of the sender's account key and additional validation for each transaction type.
func (tx *Transaction) ValidateMutableValue(db StateDB, signer Signer, currentBlockNumber uint64) error {
	// validate the sender's account key
	accKey := db.GetKey(tx.ValidatedSender())
	if tx.IsEthereumTransaction() {
		if !accKey.Type().IsLegacyAccountKey() {
			return ErrNotLegacyAccount
		}
	} else {
		if pubkey, err := SenderPubkey(signer, tx); err != nil {
			return ErrInvalidSigSender
		} else if accountkey.ValidateAccountKey(currentBlockNumber, tx.ValidatedSender(), accKey, pubkey, tx.GetRoleTypeForValidation()) != nil {
			return ErrInvalidAccountKey
		}
	}

	// validate the fee payer's account key
	if tx.IsFeeDelegatedTransaction() {
		feePayerAccKey := db.GetKey(tx.ValidatedFeePayer())
		if feePayerPubkey, err := SenderFeePayerPubkey(signer, tx); err != nil {
			return ErrInvalidSigFeePayer
		} else if accountkey.ValidateAccountKey(currentBlockNumber, tx.ValidatedFeePayer(), feePayerAccKey, feePayerPubkey, accountkey.RoleFeePayer) != nil {
			return ErrInvalidAccountKey
		}
	}

	return tx.data.ValidateMutableValue(db, currentBlockNumber)
}

func (tx *Transaction) GetRoleTypeForValidation() accountkey.RoleType {
	return tx.data.GetRoleTypeForValidation()
}

func (tx *Transaction) Data() []byte {
	tp, ok := tx.data.(TxInternalDataPayload)
	if !ok {
		return []byte{}
	}

	return common.CopyBytes(tp.GetPayload())
}

// IsFeeDelegatedTransaction returns true if the transaction is a fee-delegated transaction.
// A fee-delegated transaction has an address of the fee payer which can be different from `from` of the tx.
func (tx *Transaction) IsFeeDelegatedTransaction() bool {
	_, ok := tx.data.(TxInternalDataFeePayer)

	return ok
}

// AnchoredData returns the anchored data of the chain data anchoring transaction.
// if the tx is not chain data anchoring transaction, it will return error.
func (tx *Transaction) AnchoredData() ([]byte, error) {
	switch tx.Type() {
	case TxTypeChainDataAnchoring:
		txData, ok := tx.data.(*TxInternalDataChainDataAnchoring)
		if ok {
			return txData.Payload, nil
		}
	case TxTypeFeeDelegatedChainDataAnchoring:
		txData, ok := tx.data.(*TxInternalDataFeeDelegatedChainDataAnchoring)
		if ok {
			return txData.Payload, nil
		}
	case TxTypeFeeDelegatedChainDataAnchoringWithRatio:
		txData, ok := tx.data.(*TxInternalDataFeeDelegatedChainDataAnchoringWithRatio)
		if ok {
			return txData.Payload, nil
		}
	}
	return []byte{}, ErrInvalidTxTypeForAnchoredData
}

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *common.Address {
	if tx.data.GetRecipient() == nil {
		return nil
	}
	to := *tx.data.GetRecipient()
	return &to
}

// From returns the from address of the transaction.
// Since a legacy transaction (TxInternalDataLegacy) does not have the field `from`,
// calling From() is failed for `TxInternalDataLegacy`.
func (tx *Transaction) From() (common.Address, error) {
	if tx.IsEthereumTransaction() {
		return common.Address{}, errLegacyTransaction
	}

	tf, ok := tx.data.(TxInternalDataFrom)
	if !ok {
		return common.Address{}, errNotImplementTxInternalDataFrom
	}

	return tf.GetFrom(), nil
}

// FeePayer returns the fee payer address.
// If the tx is a fee-delegated transaction, it returns the specified fee payer.
// Otherwise, it returns `from` of the tx.
func (tx *Transaction) FeePayer() (common.Address, error) {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		// if the tx is not a fee-delegated transaction, the fee payer is `from` of the tx.
		return tx.From()
	}

	return tf.GetFeePayer(), nil
}

// FeeRatio returns the fee ratio of a transaction and a boolean value indicating TxInternalDataFeeRatio implementation.
// If the transaction does not implement TxInternalDataFeeRatio,
// it returns MaxFeeRatio which means the fee payer will be paid all tx fee by default.
func (tx *Transaction) FeeRatio() (FeeRatio, bool) {
	tf, ok := tx.data.(TxInternalDataFeeRatio)
	if !ok {
		// default fee ratio is MaxFeeRatio.
		return MaxFeeRatio, ok
	}

	return tf.GetFeeRatio(), ok
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}

	var v common.Hash
	if tx.IsEthTypedTransaction() {
		te := tx.data.(TxInternalDataEthTyped)
		v = te.TxHash()
	} else {
		v = rlpHash(tx)
	}

	tx.hash.Store(v)
	return v
}

// Size returns the true RLP encoded storage size of the transaction, either by
// encoding and returning it, or returning a previsouly cached value.
func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}

	size := calculateTxSize(tx.data)
	tx.size.Store(size)

	return size
}

func (tx *Transaction) SetTime(t time.Time) {
	tx.time = t
}

// Time returns the time that transaction was created.
func (tx *Transaction) Time() time.Time {
	return tx.time
}

// FillContractAddress fills contract address to receipt. This only works for types deploying a smart contract.
func (tx *Transaction) FillContractAddress(from common.Address, r *Receipt) {
	if filler, ok := tx.data.(TxInternalDataContractAddressFiller); ok {
		filler.FillContractAddress(from, r)
	}
}

// Execute performs execution of the transaction. This function will be called from StateTransition.TransitionDb().
// Since each transaction type performs different execution, this function calls TxInternalData.TransitionDb().
func (tx *Transaction) Execute(vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) ([]byte, uint64, error) {
	sender := NewAccountRefWithFeePayer(tx.ValidatedSender(), tx.ValidatedFeePayer())
	return tx.data.Execute(sender, vm, stateDB, currentBlockNumber, gas, value)
}

// AsMessageWithAccountKeyPicker returns the transaction as a blockchain.Message.
//
// AsMessageWithAccountKeyPicker requires a signer to derive the sender and AccountKeyPicker.
//
// XXX Rename message to something less arbitrary?
// TODO-Kaia: Message is removed and this function will return *Transaction.
func (tx *Transaction) AsMessageWithAccountKeyPicker(s Signer, picker AccountKeyPicker, currentBlockNumber uint64) (*Transaction, error) {
	intrinsicGas, err := tx.IntrinsicGas(currentBlockNumber)
	if err != nil {
		return nil, err
	}

	gasFrom, err := tx.ValidateSender(s, picker, currentBlockNumber)
	if err != nil {
		return nil, ErrSender(err)
	}

	tx.mu.Lock()
	tx.checkNonce = true
	tx.mu.Unlock()

	gasFeePayer := uint64(0)
	if tx.IsFeeDelegatedTransaction() {
		gasFeePayer, err = tx.ValidateFeePayer(s, picker, currentBlockNumber)
		if err != nil {
			return nil, ErrFeePayer(err)
		}
	}

	sigValidationGas := gasFrom + gasFeePayer
	intrinsicGas = intrinsicGas + sigValidationGas

	tx.mu.Lock()
	tx.validatedGas = &ValidatedGas{IntrinsicGas: intrinsicGas, SigValidateGas: sigValidationGas}
	tx.mu.Unlock()

	return tx, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}

	cpy := &Transaction{data: tx.data, time: tx.time}
	if tx.Type().IsEthTypedTransaction() {
		te, ok := cpy.data.(TxInternalDataEthTyped)
		if ok {
			te.setSignatureValues(signer.ChainID(), v, r, s)
		} else {
			return nil, errNotImplementTxInternalEthTyped
		}
	}

	cpy.data.SetSignature(TxSignatures{&TxSignature{v, r, s}})
	return cpy, nil
}

// WithFeePayerSignature returns a new transaction with the given fee payer signature.
func (tx *Transaction) WithFeePayerSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}

	cpy := &Transaction{data: tx.data, time: tx.time}

	feePayerSig := TxSignatures{&TxSignature{v, r, s}}
	if err := cpy.SetFeePayerSignatures(feePayerSig); err != nil {
		return nil, err
	}
	return cpy, nil
}

// Cost returns amount + gasprice * gaslimit.
func (tx *Transaction) Cost() *big.Int {
	total := tx.Fee()
	total.Add(total, tx.data.GetAmount())
	return total
}

func (tx *Transaction) Fee() *big.Int {
	return new(big.Int).Mul(tx.data.GetPrice(), new(big.Int).SetUint64(tx.data.GetGasLimit()))
}

// Sign signs the tx with the given signer and private key.
func (tx *Transaction) Sign(s Signer, prv *ecdsa.PrivateKey) error {
	h := s.Hash(tx)
	sig, err := NewTxSignatureWithValues(s, tx, h, prv)
	if err != nil {
		return err
	}

	tx.SetSignature(TxSignatures{sig})
	return nil
}

// SignWithKeys signs the tx with the given signer and a slice of private keys.
func (tx *Transaction) SignWithKeys(s Signer, prv []*ecdsa.PrivateKey) error {
	h := s.Hash(tx)
	sig, err := NewTxSignaturesWithValues(s, tx, h, prv)
	if err != nil {
		return err
	}

	tx.SetSignature(sig)
	return nil
}

// SignFeePayer signs the tx with the given signer and private key as a fee payer.
func (tx *Transaction) SignFeePayer(s Signer, prv *ecdsa.PrivateKey) error {
	h, err := s.HashFeePayer(tx)
	if err != nil {
		return err
	}
	sig, err := NewTxSignatureWithValues(s, tx, h, prv)
	if err != nil {
		return err
	}

	if err := tx.SetFeePayerSignatures(TxSignatures{sig}); err != nil {
		return err
	}

	return nil
}

// SignFeePayerWithKeys signs the tx with the given signer and a slice of private keys as a fee payer.
func (tx *Transaction) SignFeePayerWithKeys(s Signer, prv []*ecdsa.PrivateKey) error {
	h, err := s.HashFeePayer(tx)
	if err != nil {
		return err
	}
	sig, err := NewTxSignaturesWithValues(s, tx, h, prv)
	if err != nil {
		return err
	}

	if err := tx.SetFeePayerSignatures(sig); err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) SetFeePayerSignatures(s TxSignatures) error {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return errNotFeeDelegationTransaction
	}

	tf.SetFeePayerSignatures(s)

	return nil
}

// GetFeePayerSignatures returns fee payer signatures of the transaction.
func (tx *Transaction) GetFeePayerSignatures() (TxSignatures, error) {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return nil, errNotFeeDelegationTransaction
	}
	return tf.GetFeePayerRawSignatureValues(), nil
}

func (tx *Transaction) SetSignature(signature TxSignatures) {
	tx.data.SetSignature(signature)
}

func (tx *Transaction) MarkUnexecutable(b bool) {
	v := int32(0)
	if b {
		v = 1
	}
	atomic.StoreInt32(&tx.markedUnexecutable, v)
}

func (tx *Transaction) IsMarkedUnexecutable() bool {
	return atomic.LoadInt32(&tx.markedUnexecutable) == 1
}

func (tx *Transaction) RawSignatureValues() TxSignatures {
	return tx.data.RawSignatureValues()
}

func (tx *Transaction) String() string {
	return tx.data.String()
}

// ValidateSender finds a sender from both legacy and new types of transactions.
// It returns the senders address and gas used for the tx validation.
func (tx *Transaction) ValidateSender(signer Signer, p AccountKeyPicker, currentBlockNumber uint64) (uint64, error) {
	if tx.IsEthereumTransaction() {
		addr, err := Sender(signer, tx)
		// Legacy transaction cannot be executed unless the account has a legacy key.
		if p.GetKey(addr).Type().IsLegacyAccountKey() == false {
			return 0, kerrors.ErrLegacyTransactionMustBeWithLegacyKey
		}
		tx.mu.Lock()
		if tx.validatedSender == (common.Address{}) {
			tx.validatedSender = addr
			tx.validatedFeePayer = addr
		}
		tx.mu.Unlock()
		return 0, err
	}

	pubkey, err := SenderPubkey(signer, tx)
	if err != nil {
		return 0, err
	}
	txfrom, ok := tx.data.(TxInternalDataFrom)
	if !ok {
		return 0, errNotTxInternalDataFrom
	}
	from := txfrom.GetFrom()
	accKey := p.GetKey(from)

	gasKey, err := accKey.SigValidationGas(currentBlockNumber, tx.GetRoleTypeForValidation(), len(pubkey))
	if err != nil {
		return 0, err
	}

	if err := accountkey.ValidateAccountKey(currentBlockNumber, from, accKey, pubkey, tx.GetRoleTypeForValidation()); err != nil {
		return 0, ErrInvalidAccountKey
	}

	tx.mu.Lock()
	if tx.validatedSender == (common.Address{}) {
		tx.validatedSender = from
		tx.validatedFeePayer = from
	}
	tx.mu.Unlock()

	return gasKey, nil
}

// ValidateFeePayer finds a fee payer from a transaction.
// If the transaction is not a fee-delegated transaction, it returns an error.
func (tx *Transaction) ValidateFeePayer(signer Signer, p AccountKeyPicker, currentBlockNumber uint64) (uint64, error) {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return 0, errUndefinedTxType
	}

	pubkey, err := SenderFeePayerPubkey(signer, tx)
	if err != nil {
		return 0, err
	}

	feePayer := tf.GetFeePayer()
	accKey := p.GetKey(feePayer)

	gasKey, err := accKey.SigValidationGas(currentBlockNumber, accountkey.RoleFeePayer, len(pubkey))
	if err != nil {
		return 0, err
	}

	if err := accountkey.ValidateAccountKey(currentBlockNumber, feePayer, accKey, pubkey, accountkey.RoleFeePayer); err != nil {
		return 0, ErrInvalidAccountKey
	}

	tx.mu.Lock()
	if tx.validatedFeePayer == tx.validatedSender {
		tx.validatedFeePayer = feePayer
	}
	tx.mu.Unlock()

	return gasKey, nil
}

// Transactions is a Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s.
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s.
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp.
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// TxDifference returns a new set t which is the difference between a to b.
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// FilterTransactionWithBaseFee returns a list of transactions for each account that filters transactions
// that are greater than or equal to baseFee.
func FilterTransactionWithBaseFee(pending map[common.Address]Transactions, baseFee *big.Int) map[common.Address]Transactions {
	txMap := make(map[common.Address]Transactions)
	for addr, list := range pending {
		txs := list
		for i, tx := range list {
			if tx.GasPrice().Cmp(baseFee) < 0 {
				txs = list[:i]
				break
			}
		}

		if len(txs) > 0 {
			txMap[addr] = txs
		}
	}
	return txMap
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int { return len(s) }
func (s TxByNonce) Less(i, j int) bool {
	return s[i].data.GetAccountNonce() < s[j].data.GetAccountNonce()
}
func (s TxByNonce) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// txWithMinerFee wraps a transaction with its gas price or effective miner gasTipCap
type txWithMinerFee struct {
	tx   *Transaction
	from common.Address
	fees *big.Int
}

// newTxWithMinerFee creates a wrapped transaction, calculating the effective
// miner gasTipCap if a base fee is provided.
// Returns error in case of a negative effective miner gasTipCap.
func newTxWithMinerFee(tx *Transaction, from common.Address, baseFee *big.Int) (*txWithMinerFee, error) {
	tip := new(big.Int).Set(tx.GasTipCap())
	if baseFee != nil {
		if tx.GasFeeCap().Cmp(baseFee) < 0 {
			return nil, errors.New("invalid gas fee cap. It must be set to value greater than or equal to baseFee")
		}
		tip = new(big.Int).Sub(tx.GasFeeCap(), baseFee)
		if tip.Cmp(tx.GasTipCap()) == 1 {
			tip = tx.GasTipCap()
		}
	}
	return &txWithMinerFee{
		tx:   tx,
		from: from,
		fees: tip,
	}, nil
}

// txByPriceAndTime implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type txByPriceAndTime []*txWithMinerFee

func (s txByPriceAndTime) Len() int { return len(s) }
func (s txByPriceAndTime) Less(i, j int) bool {
	// If the prices are equal, use the time the transaction was first seen for
	// deterministic sorting
	cmp := s[i].fees.Cmp(s[j].fees)
	if cmp == 0 {
		return s[i].tx.Time().Before(s[j].tx.Time())
	}
	return cmp > 0
}
func (s txByPriceAndTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *txByPriceAndTime) Push(x interface{}) {
	*s = append(*s, x.(*txWithMinerFee))
}

func (s *txByPriceAndTime) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	old[n-1] = nil
	*s = old[0 : n-1]
	return x
}

// SortTxsByPriceAndTime is used to sort the txs by expected effectiveGasTip and then arrival time.
// It is called on the process of txs broadcasting. There's three points when this function called.
// (1) BroadcastTxs: before broadcasting txs to the peers
// (2) RebroadcastTxs: before rebroadcasting the remaining pending txs to the peers
// (3) syncTransactions: before sending the all pending txs to the newly connected peer
func SortTxsByPriceAndTime(txs Transactions, baseFee *big.Int) Transactions {
	sortedTxsWithMinerFee := make(txByPriceAndTime, len(txs))
	for i, tx := range txs {
		// fee cannot be negative
		sortedTxsWithMinerFee[i] = &txWithMinerFee{tx, common.Address{}, math.BigMax(tx.EffectiveGasTip(baseFee), big.NewInt(0))}
	}

	// If already sorted, just return original txs.
	if sort.IsSorted(sortedTxsWithMinerFee) {
		return txs
	}

	// Sort the batch of txs and derive sortedTxs to return it.
	sort.Sort(sortedTxsWithMinerFee)
	sortedTxs := make(Transactions, len(txs))
	for i, tx := range sortedTxsWithMinerFee {
		sortedTxs[i] = tx.tx
	}
	return sortedTxs
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximizing sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce struct {
	txs     map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads   txByPriceAndTime                // Next transaction for each unique account (price heap)
	signer  Signer                          // Signer for the set of transactions
	baseFee *big.Int                        // Current base fee
}

// NewTransactionsByPriceAndNonce creates a transaction set that can retrieve
// price sorted transactions in a nonce-honouring way.
//
// Note, the input map is reowned so the caller should not interact any more with
// if after providing it to the constructor.
func NewTransactionsByPriceAndNonce(signer Signer, txs map[common.Address]Transactions, baseFee *big.Int) *TransactionsByPriceAndNonce {
	// Initialize a price and received time based heap with the head transactions
	heads := make(txByPriceAndTime, 0, len(txs))
	for from, accTxs := range txs {
		wrapped, err := newTxWithMinerFee(accTxs[0], from, baseFee)
		if err != nil {
			delete(txs, from)
			continue
		}
		heads = append(heads, wrapped)
		txs[from] = accTxs[1:]
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce{
		txs:     txs,
		heads:   heads,
		signer:  signer,
		baseFee: baseFee,
	}
}

// Peek returns the next transaction by price and nonce.
func (t *TransactionsByPriceAndNonce) Peek() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0].tx
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce) Shift() {
	if len(t.heads) == 0 {
		return
	}
	acc := t.heads[0].from
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		if wrapped, err := newTxWithMinerFee(txs[0], acc, t.baseFee); err == nil {
			t.heads[0], t.txs[acc] = wrapped, txs[1:]
			heap.Fix(&t.heads, 0)
			return
		}
	}
	heap.Pop(&t.heads)
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriceAndNonce) Pop() {
	if len(t.heads) == 0 {
		return
	}
	heap.Pop(&t.heads)
}

// Empty returns if the price heap is empty. It can be used to check it simpler
// than calling peek and checking for nil return.
func (t *TransactionsByPriceAndNonce) Empty() bool {
	return len(t.heads) == 0
}

// Clear removes the entire content of the heap.
func (t *TransactionsByPriceAndNonce) Clear() {
	t.heads, t.txs = nil, nil
}

// Copy the current object.
func (t *TransactionsByPriceAndNonce) Copy() *TransactionsByPriceAndNonce {
	txsCopy := make(map[common.Address]Transactions)
	for addr, txList := range t.txs {
		txsCopy[addr] = txList
	}

	headsCopy := make(txByPriceAndTime, len(t.heads))
	copy(headsCopy, t.heads)

	return &TransactionsByPriceAndNonce{
		txs:     txsCopy, // shift changes it.
		heads:   t.heads, // pop, shift changes it.
		signer:  t.signer,
		baseFee: t.baseFee, // read-only
	}
}

// NewMessage returns a `*Transaction` object with the given arguments.
// Care must be taken when creating SetCodeTx because if you assign nil to `to`,
// a panic will occur because `newTxInternalDataEthereumSetCodeWithValues` reference the pointer of `to`.
func NewMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice, gasFeeCap, gasTipCap *big.Int, data []byte, checkNonce bool, intrinsicGas uint64, list AccessList, chainId *big.Int, auth []SetCodeAuthorization) *Transaction {
	transaction := &Transaction{
		validatedGas:      &ValidatedGas{IntrinsicGas: intrinsicGas, SigValidateGas: 0},
		validatedFeePayer: from,
		validatedSender:   from,
		checkNonce:        checkNonce,
	}

	// Call supports EthereumAccessList, EthereumSetCode and Legacy txTypes only.
	if auth != nil {
		internalData := newTxInternalDataEthereumSetCodeWithValues(nonce, *to, amount, gasLimit, gasFeeCap, gasTipCap, data, list, chainId, auth)
		transaction.setDecoded(internalData, 0)
	} else if list != nil {
		internalData := newTxInternalDataEthereumAccessListWithValues(nonce, to, amount, gasLimit, gasPrice, data, list, chainId)
		transaction.setDecoded(internalData, 0)
	} else {
		internalData := newTxInternalDataLegacyWithValues(nonce, to, amount, gasLimit, gasPrice, data)
		transaction.setDecoded(internalData, 0)
	}

	return transaction
}
