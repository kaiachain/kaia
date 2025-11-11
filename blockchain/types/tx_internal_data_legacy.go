// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package types

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto/sha3"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

type TxInternalDataLegacy struct {
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`
	GasLimit     uint64          `json:"gas"      gencodec:"required"`
	Recipient    *common.Address `json:"to"       rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"    gencodec:"required"`
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataLegacyJSON struct {
	AccountNonce hexutil.Uint64   `json:"nonce"`
	Price        *hexutil.Big     `json:"gasPrice"`
	GasLimit     hexutil.Uint64   `json:"gas"`
	Recipient    *common.Address  `json:"to"`
	Amount       *hexutil.Big     `json:"value"`
	Payload      hexutil.Bytes    `json:"input"`
	TxSignatures TxSignaturesJSON `json:"signatures"`
	Hash         *common.Hash     `json:"hash"`
}

func newEmptyTxInternalDataLegacy() *TxInternalDataLegacy {
	return &TxInternalDataLegacy{}
}

func newTxInternalDataLegacy() *TxInternalDataLegacy {
	return &TxInternalDataLegacy{
		AccountNonce: 0,
		Recipient:    nil,
		Payload:      []byte{},
		Amount:       new(big.Int),
		GasLimit:     0,
		Price:        new(big.Int),
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
}

func newTxInternalDataLegacyWithValues(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) *TxInternalDataLegacy {
	d := newTxInternalDataLegacy()

	d.AccountNonce = nonce
	d.Recipient = to
	d.GasLimit = gasLimit

	if len(data) > 0 {
		d.Payload = common.CopyBytes(data)
	}
	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return d
}

func newTxInternalDataLegacyWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataLegacy, error) {
	d := newTxInternalDataLegacy()

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		d.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyTo].(common.Address); ok {
		d.Recipient = &v
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

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		d.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		return nil, errValueKeyGasPriceMustBigInt
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataLegacy) Type() TxType {
	return TxTypeLegacyTransaction
}

func (t *TxInternalDataLegacy) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataLegacy) ChainId() *big.Int {
	return deriveChainId(t.V)
}

func (t *TxInternalDataLegacy) GetNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataLegacy) GetGasPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataLegacy) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataLegacy) GetTo() *common.Address {
	return t.Recipient
}

func (t *TxInternalDataLegacy) GetValue() *big.Int {
	return new(big.Int).Set(t.Amount)
}

func (t *TxInternalDataLegacy) GetData() []byte {
	return t.Payload
}

func (t *TxInternalDataLegacy) setHashForMarshaling(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataLegacy) SetSignature(s TxSignatures) {
	if len(s) != 1 {
		logger.Crit("LegacyTransaction receives a single signature only!")
	}

	t.V = s[0].V
	t.R = s[0].R
	t.S = s[0].S
}

func (t *TxInternalDataLegacy) RawSignatureValues() TxSignatures {
	return TxSignatures{&TxSignature{t.V, t.R, t.S}}
}

func (t *TxInternalDataLegacy) RecoverAddress(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) (common.Address, error) {
	V := vfunc(t.V)
	return recoverPlain(txhash, t.R, t.S, V, homestead)
}

func (t *TxInternalDataLegacy) RecoverPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	V := vfunc(t.V)

	pk, err := recoverPlainPubkey(txhash, t.R, t.S, V, homestead)
	if err != nil {
		return nil, err
	}

	return []*ecdsa.PublicKey{pk}, nil
}

func (t *TxInternalDataLegacy) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	return IntrinsicGas(t.Payload, nil, nil, t.Recipient == nil, *fork.Rules(big.NewInt(int64(currentBlockNumber))))
}

func (t *TxInternalDataLegacy) SerializeForSign() []interface{} {
	return []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
	}
}

func (t *TxInternalDataLegacy) SenderTxHash() common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.V,
		t.R,
		t.S,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataLegacy) String() string {
	var from, to string
	tx := &Transaction{data: t}

	v, r, s := t.V, t.R, t.S
	if v != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(v)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if t.GetTo() == nil {
		to = "[contract creation]"
	} else {
		to = hex.EncodeToString(t.GetTo().Bytes())
	}
	enc, _ := rlp.EncodeToBytes(t)
	return fmt.Sprintf(`
	TX(%x)
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit  %#x
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
`,
		tx.Hash(),
		t.GetTo() == nil,
		from,
		to,
		t.GetNonce(),
		t.GetGasPrice(),
		t.GetGasLimit(),
		t.GetValue(),
		t.GetData(),
		v,
		r,
		s,
		enc,
	)
}

func (t *TxInternalDataLegacy) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	if t.Recipient != nil {
		if common.IsPrecompiledContractAddress(*t.Recipient, *fork.Rules(big.NewInt(int64(currentBlockNumber)))) {
			return kerrors.ErrPrecompiledContractAddress
		}
	}
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataLegacy) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataLegacy) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//start := time.Now()
	//defer func() {
	//	elapsed := time.Since(start)
	//	logger.Debug("[TxInternalDataLegacy] EVM execution done", "elapsed", elapsed)
	//}()
	///////////////////////////////////////////////////////
	if t.Recipient == nil {
		// Sender's nonce will be increased in '`vm.Create()`
		ret, _, usedGas, err = vm.Create(sender, t.Payload, gas, value, params.CodeFormatEVM)
	} else {
		stateDB.IncNonce(sender.Address())
		ret, usedGas, err = vm.Call(sender, *t.Recipient, t.Payload, gas, value)
	}
	return ret, usedGas, err
}

func (t *TxInternalDataLegacy) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"typeInt":    t.Type(),
		"type":       t.Type().String(),
		"gas":        hexutil.Uint64(t.GasLimit),
		"gasPrice":   (*hexutil.Big)(t.Price),
		"input":      hexutil.Bytes(t.Payload),
		"nonce":      hexutil.Uint64(t.AccountNonce),
		"to":         t.Recipient,
		"value":      (*hexutil.Big)(t.Amount),
		"signatures": TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
	}
}

func (t *TxInternalDataLegacy) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataLegacyJSON{
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.Big)(t.Amount),
		t.Payload,
		TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
		t.Hash,
	})
}

func (t *TxInternalDataLegacy) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataLegacyJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.Recipient = js.Recipient
	t.Amount = (*big.Int)(js.Amount)
	t.Payload = js.Payload
	t.V = (*big.Int)(js.TxSignatures[0].V)
	t.R = (*big.Int)(js.TxSignatures[0].R)
	t.S = (*big.Int)(js.TxSignatures[0].S)
	t.Hash = js.Hash

	return nil
}
