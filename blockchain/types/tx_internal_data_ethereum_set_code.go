// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/rlp"
)

// DelegationPrefix is used by code to denote the account is delegating to
// another account.
var DelegationPrefix = []byte{0xef, 0x01, 0x00}

// ParseDelegation tries to parse the address from a delegation slice.
func ParseDelegation(b []byte) (common.Address, bool) {
	if len(b) != 23 || !bytes.HasPrefix(b, DelegationPrefix) {
		return common.Address{}, false
	}
	var addr common.Address
	copy(addr[:], b[len(DelegationPrefix):])
	return addr, true
}

// AddressToDelegation adds the delegation prefix to the specified address.
func AddressToDelegation(addr common.Address) []byte {
	return append(DelegationPrefix, addr.Bytes()...)
}

// TxInternalDataEthereumSetCode represents a set code transaction.
type TxInternalDataEthereumSetCode struct {
	ChainID           *uint256.Int
	AccountNonce      uint64
	GasTipCap         *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap         *big.Int // a.k.a. maxFeePerGas
	GasLimit          uint64
	Recipient         common.Address
	Amount            *big.Int
	Payload           []byte
	AccessList        AccessList
	AuthorizationList []SetCodeAuthorization

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataEthereumSetCodeJSON struct {
	Type                 TxType                 `json:"typeInt"`
	TypeStr              string                 `json:"type"`
	ChainID              *hexutil.U256          `json:"chainId"`
	AccountNonce         hexutil.Uint64         `json:"nonce"`
	MaxPriorityFeePerGas *hexutil.Big           `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big           `json:"maxFeePerGas"`
	GasLimit             hexutil.Uint64         `json:"gas"`
	Recipient            common.Address         `json:"to"`
	Amount               *hexutil.Big           `json:"value"`
	Payload              hexutil.Bytes          `json:"input"`
	AccessList           AccessList             `json:"accessList"`
	AuthorizationList    []SetCodeAuthorization `json:"authorizationList"`
	TxSignatures         TxSignaturesJSON       `json:"signatures"`
	Hash                 *common.Hash           `json:"hash"`
}

func newTxInternalDataEthereumSetCode() *TxInternalDataEthereumSetCode {
	return &TxInternalDataEthereumSetCode{
		ChainID:           new(uint256.Int),
		AccountNonce:      0,
		GasTipCap:         new(big.Int),
		GasFeeCap:         new(big.Int),
		GasLimit:          0,
		Recipient:         common.Address{},
		Amount:            new(big.Int),
		Payload:           []byte{},
		AccessList:        AccessList{},
		AuthorizationList: []SetCodeAuthorization{},
		V:                 new(big.Int),
		R:                 new(big.Int),
		S:                 new(big.Int),
	}
}

func newTxInternalDataEthereumSetCodeWithValues(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasFeeCap, gasTipCap *big.Int, data []byte, accessList AccessList, authList []SetCodeAuthorization, chainID *big.Int) *TxInternalDataEthereumSetCode {
	d := newTxInternalDataEthereumSetCode()

	d.AccountNonce = nonce
	d.Recipient = to
	d.GasLimit = gasLimit

	if chainID != nil {
		d.ChainID.Set(uint256.MustFromBig(chainID))
	}

	if gasTipCap != nil {
		d.GasTipCap.Set(gasTipCap)
	}

	if gasFeeCap != nil {
		d.GasFeeCap.Set(gasFeeCap)
	}

	if amount != nil {
		d.Amount.Set(amount)
	}

	if len(data) > 0 {
		d.Payload = common.CopyBytes(data)
	}

	if accessList != nil {
		d.AccessList = append(d.AccessList, accessList...)
	}

	if authList != nil {
		d.AuthorizationList = authList
	}

	return d
}

func newTxInternalDataEthereumSetCodeWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataEthereumSetCode, error) {
	d := newTxInternalDataEthereumSetCode()

	if v, ok := values[TxValueKeyChainID].(*big.Int); ok {
		d.ChainID.Set(uint256.MustFromBig(v))
		delete(values, TxValueKeyChainID)
	} else {
		return nil, errValueKeyChainIDInvalid
	}

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		d.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyTo].(common.Address); ok {
		d.Recipient = v
		delete(values, TxValueKeyTo)
	} else {
		return nil, errValueKeyToMustAddress
	}

	if v, ok := values[TxValueKeyAmount].(*big.Int); ok {
		d.Amount.Set(v)
		delete(values, TxValueKeyAmount)
	} else {
		return nil, errValueKeyAmountMustBigInt
	}

	if v, ok := values[TxValueKeyData].([]byte); ok {
		d.Payload = common.CopyBytes(v)
		delete(values, TxValueKeyData)
	} else {
		return nil, errValueKeyDataMustByteSlice
	}

	if v, ok := values[TxValueKeyGasLimit].(uint64); ok {
		d.GasLimit = v
		delete(values, TxValueKeyGasLimit)
	} else {
		return nil, errValueKeyGasLimitMustUint64
	}

	if v, ok := values[TxValueKeyGasFeeCap].(*big.Int); ok {
		d.GasFeeCap.Set(v)
		delete(values, TxValueKeyGasFeeCap)
	} else {
		return nil, errValueKeyGasFeeCapMustBigInt
	}
	if v, ok := values[TxValueKeyGasTipCap].(*big.Int); ok {
		d.GasTipCap.Set(v)
		delete(values, TxValueKeyGasTipCap)
	} else {
		return nil, errValueKeyGasTipCapMustBigInt
	}
	if v, ok := values[TxValueKeyAccessList].(AccessList); ok {
		d.AccessList = make(AccessList, len(v))
		copy(d.AccessList, v)
		delete(values, TxValueKeyAccessList)
	} else {
		return nil, errValueKeyAccessListInvalid
	}
	if v, ok := values[TxValueKeyAuthorizationList].([]SetCodeAuthorization); ok {
		d.AuthorizationList = make([]SetCodeAuthorization, len(v))
		copy(d.AuthorizationList, v)
		delete(values, TxValueKeyAuthorizationList)
	} else {
		return nil, errValueKeyAuthorizationListInvalid
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataEthereumSetCode) Type() TxType {
	return TxTypeEthereumSetCode
}

func (t *TxInternalDataEthereumSetCode) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataEthereumSetCode) GetPrice() *big.Int {
	return t.GasFeeCap
}

func (t *TxInternalDataEthereumSetCode) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataEthereumSetCode) GetRecipient() *common.Address {
	if t.Recipient == (common.Address{}) {
		return nil
	}

	to := common.Address(t.Recipient)
	return &to
}

func (t *TxInternalDataEthereumSetCode) GetAmount() *big.Int {
	return new(big.Int).Set(t.Amount)
}

func (t *TxInternalDataEthereumSetCode) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataEthereumSetCode) GetAccessList() AccessList {
	return t.AccessList
}

func (tx *TxInternalDataEthereumSetCode) GetAuthorizationList() []SetCodeAuthorization {
	return tx.AuthorizationList
}

func (t *TxInternalDataEthereumSetCode) GetGasTipCap() *big.Int {
	return t.GasTipCap
}

func (t *TxInternalDataEthereumSetCode) GetGasFeeCap() *big.Int {
	return t.GasFeeCap
}

func (t *TxInternalDataEthereumSetCode) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataEthereumSetCode) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataEthereumSetCode) SetSignature(signatures TxSignatures) {
	if len(signatures) != 1 {
		logger.Crit("TxTypeEthereum can receive only single signature!")
	}

	t.V = signatures[0].V
	t.R = signatures[0].R
	t.S = signatures[0].S
}

func (t *TxInternalDataEthereumSetCode) RawSignatureValues() TxSignatures {
	return TxSignatures{&TxSignature{t.V, t.R, t.S}}
}

func (t *TxInternalDataEthereumSetCode) ValidateSignature() bool {
	v := byte(t.V.Uint64())
	return crypto.ValidateSignatureValues(v, t.R, t.S, false)
}

func (t *TxInternalDataEthereumSetCode) RecoverAddress(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) (common.Address, error) {
	V := vfunc(t.V)
	return recoverPlain(txhash, t.R, t.S, V, homestead)
}

func (t *TxInternalDataEthereumSetCode) RecoverPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	V := vfunc(t.V)

	pk, err := recoverPlainPubkey(txhash, t.R, t.S, V, homestead)
	if err != nil {
		return nil, err
	}

	return []*ecdsa.PublicKey{pk}, nil
}

func (t *TxInternalDataEthereumSetCode) ChainId() *big.Int {
	return t.ChainID.ToBig()
}

func (t *TxInternalDataEthereumSetCode) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataEthereumSetCode)
	if !ok {
		return false
	}

	return t.ChainID.Cmp(ta.ChainID) == 0 &&
		t.AccountNonce == ta.AccountNonce &&
		t.GasFeeCap.Cmp(ta.GasFeeCap) == 0 &&
		t.GasTipCap.Cmp(ta.GasTipCap) == 0 &&
		t.GasLimit == ta.GasLimit &&
		t.Recipient == ta.Recipient &&
		t.Amount.Cmp(ta.Amount) == 0 &&
		reflect.DeepEqual(t.AccessList, ta.AccessList) &&
		reflect.DeepEqual(t.AuthorizationList, ta.AuthorizationList) &&
		t.V.Cmp(ta.V) == 0 &&
		t.R.Cmp(ta.R) == 0 &&
		t.S.Cmp(ta.S) == 0
}

func (t *TxInternalDataEthereumSetCode) IntrinsicGas(currentBlockNumber uint64) (uint64, uint64, error) {
	return IntrinsicGas(t.Payload, t.AccessList, t.AuthorizationList, false, *fork.Rules(big.NewInt(int64(currentBlockNumber))))
}

func (t *TxInternalDataEthereumSetCode) SerializeForSign() []interface{} {
	// If the chainId has nil or empty value, It will be set signer's chainId.
	return []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.GasTipCap,
		t.GasFeeCap,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
		t.AuthorizationList,
	}
}

func (t *TxInternalDataEthereumSetCode) TxHash() common.Hash {
	return prefixedRlpHash(byte(t.Type()), []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.GasTipCap,
		t.GasFeeCap,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
		t.AuthorizationList,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumSetCode) SenderTxHash() common.Hash {
	return prefixedRlpHash(byte(t.Type()), []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.GasTipCap,
		t.GasFeeCap,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
		t.AuthorizationList,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumSetCode) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	if t.Recipient == (common.Address{}) {
		return kerrors.ErrEmptyRecipient
	} else {
		if common.IsPrecompiledContractAddress(t.Recipient) {
			return kerrors.ErrPrecompiledContractAddress
		}
	}
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataEthereumSetCode) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataEthereumSetCode) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataEthereumSetCode) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataEthereumSetCode) String() string {
	var from, to string
	tx := &Transaction{data: t}

	v, r, s := t.V, t.R, t.S
	if v != nil {
		signer := LatestSignerForChainID(t.ChainId())
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = hex.EncodeToString(f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if t.GetRecipient() == nil {
		to = "[contract creation]"
	} else {
		to = hex.EncodeToString(t.GetRecipient().Bytes())
	}
	enc, _ := rlp.EncodeToBytes(tx)
	return fmt.Sprintf(`
		TX(%x)
		Contract: %v
		Chaind:   %#x
		From:     %s
		To:       %s
		Nonce:    %v
		GasTipCap: %#x
		GasFeeCap: %#x
		GasLimit  %#x
		Value:    %#x
		Data:     0x%x
		AccessList: %x
		AuthorizationList: %x
		V:        %#x
		R:        %#x
		S:        %#x
		Hex:      %x
	`,
		tx.Hash(),
		t.GetRecipient() == nil,
		t.ChainId(),
		from,
		to,
		t.GetAccountNonce(),
		t.GetGasTipCap(),
		t.GetGasFeeCap(),
		t.GetGasLimit(),
		t.GetAmount(),
		t.GetPayload(),
		t.AccessList,
		t.AuthorizationList,
		v,
		r,
		s,
		enc,
	)
}

func (t *TxInternalDataEthereumSetCode) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//start := time.Now()
	//defer func() {
	//	elapsed := time.Since(start)
	//	logger.Debug("[TxInternalDataLegacy] EVM execution done", "elapsed", elapsed)
	//}()
	///////////////////////////////////////////////////////

	// SetCode sender nonce increment is done in TransitionDb to comply with the spec.
	return vm.Call(sender, t.Recipient, t.Payload, gas, value)
}

func (t *TxInternalDataEthereumSetCode) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"typeInt":              t.Type(),
		"type":                 t.Type().String(),
		"chainId":              (*hexutil.Big)(t.ChainId()),
		"nonce":                hexutil.Uint64(t.AccountNonce),
		"maxPriorityFeePerGas": (*hexutil.Big)(t.GasTipCap),
		"maxFeePerGas":         (*hexutil.Big)(t.GasFeeCap),
		"gas":                  hexutil.Uint64(t.GasLimit),
		"to":                   t.Recipient,
		"input":                hexutil.Bytes(t.Payload),
		"value":                (*hexutil.Big)(t.Amount),
		"accessList":           t.AccessList,
		"authorizationList":    t.AuthorizationList,
		"signatures":           TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
	}
}

func (t *TxInternalDataEthereumSetCode) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataEthereumSetCodeJSON{
		t.Type(),
		t.Type().String(),
		(*hexutil.U256)(t.ChainID),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.GasTipCap),
		(*hexutil.Big)(t.GasFeeCap),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.Big)(t.Amount),
		t.Payload,
		t.AccessList,
		t.AuthorizationList,
		TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
		t.Hash,
	})
}

func (t *TxInternalDataEthereumSetCode) UnmarshalJSON(bytes []byte) error {
	js := &TxInternalDataEthereumSetCodeJSON{}
	if err := json.Unmarshal(bytes, js); err != nil {
		return err
	}

	t.ChainID = (*uint256.Int)(js.ChainID)
	t.AccountNonce = uint64(js.AccountNonce)
	t.GasTipCap = (*big.Int)(js.MaxPriorityFeePerGas)
	t.GasFeeCap = (*big.Int)(js.MaxFeePerGas)
	t.GasLimit = uint64(js.GasLimit)
	t.Recipient = js.Recipient
	t.Amount = (*big.Int)(js.Amount)
	t.Payload = js.Payload
	t.AccessList = js.AccessList
	t.AuthorizationList = js.AuthorizationList
	t.V = (*big.Int)(js.TxSignatures[0].V)
	t.R = (*big.Int)(js.TxSignatures[0].R)
	t.S = (*big.Int)(js.TxSignatures[0].S)
	t.Hash = js.Hash

	return nil
}

func (t *TxInternalDataEthereumSetCode) setSignatureValues(chainID, v, r, s *big.Int) {
	t.ChainID, t.V, t.R, t.S = uint256.MustFromBig(chainID), v, r, s
}

type SetCodeAuthorization struct {
	ChainID uint256.Int    `json:"chainId"`
	Address common.Address `json:"address"`
	Nonce   uint64         `json:"nonce"`
	V       uint8          `json:"v"`
	R       *big.Int       `json:"r"`
	S       *big.Int       `json:"s"`
}

// SignSetCode creates a signed the SetCode authorization.
func SignSetCode(prv *ecdsa.PrivateKey, auth SetCodeAuthorization) (SetCodeAuthorization, error) {
	sighash := auth.sigHash()
	sig, err := crypto.Sign(sighash[:], prv)
	if err != nil {
		return SetCodeAuthorization{}, err
	}
	r, s, _ := decodeSignature(sig)
	return SetCodeAuthorization{
		ChainID: auth.ChainID,
		Address: auth.Address,
		Nonce:   auth.Nonce,
		V:       sig[crypto.RecoveryIDOffset],
		R:       r,
		S:       s,
	}, nil
}

func (a *SetCodeAuthorization) sigHash() common.Hash {
	return prefixedRlpHash(0x05, []any{
		a.ChainID,
		a.Address,
		a.Nonce,
	})
}

// Authority recovers the the authorizing account of an authorization.
func (a *SetCodeAuthorization) Authority() (common.Address, error) {
	sighash := a.sigHash()
	if !crypto.ValidateSignatureValues(a.V, a.R, a.S, true) {
		return common.Address{}, ErrInvalidSig
	}
	// encode the signature in uncompressed format
	r, s := a.R.Bytes(), a.S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[64] = a.V
	// recover the public key from the signature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return common.Address{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return common.Address{}, errors.New("invalid public key")
	}
	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}
