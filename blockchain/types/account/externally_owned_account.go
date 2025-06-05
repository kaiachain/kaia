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

package account

import (
	"bytes"
	"encoding/json"
	"io"
	"math/big"

	"github.com/kaiachain/kaia/v2/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/common/hexutil"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/rlp"
)

var emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

// ExternallyOwnedAccount represents a Kaia account used by a user.
type ExternallyOwnedAccount struct {
	*AccountCommon
	storageRoot common.ExtHash // merkle root plus optional sequence of the storage trie
	codeHash    []byte
	codeInfo    params.CodeInfo // consists of two information, vmVersion and codeFormat
}

type externallyOwnedAccountSerializable struct {
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.Hash
	CodeHash           []byte
	CodeInfo           params.CodeInfo
}

// smartContractAccountSerializableExt is an internal data structure for RLP serialization.
// This structure inherits accountCommonSerializable.
// nolint: maligned  // Because it is a temporary struct, memory footprint is not important.
type externallyOwnedAccountSerializableExt struct {
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.ExtHash
	CodeHash           []byte
	CodeInfo           params.CodeInfo
}

type externallyOwnedAccountSerializableJSON struct {
	Nonce         uint64                           `json:"nonce"`
	Balance       *hexutil.Big                     `json:"balance"`
	HumanReadable bool                             `json:"humanReadable"`
	Key           *accountkey.AccountKeySerializer `json:"key"`
	StorageRoot   common.Hash                      `json:"storageRoot,omitempty"`
	CodeHash      []byte                           `json:"codeHash,omitempty"`
	CodeFormat    params.CodeFormat                `json:"codeFormat"`
	VmVersion     params.VmVersion                 `json:"vmVersion"`
}

// newExternallyOwnedAccount creates an ExternallyOwnedAccount object with default values.
func newExternallyOwnedAccount() *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		newAccountCommon(),
		emptyRoot.ExtendZero(),
		emptyCodeHash,
		params.CodeInfo(0),
	}
}

// newExternallyOwnedAccountWithMap creates an ExternallyOwnedAccount object initialized with the given values.
func newExternallyOwnedAccountWithMap(values map[AccountValueKeyType]interface{}) *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		newAccountCommonWithMap(values),
		emptyRoot.ExtendZero(),
		emptyCodeHash,
		params.CodeInfo(0),
	}
}

func (e *ExternallyOwnedAccount) Type() AccountType {
	return ExternallyOwnedAccountType
}

func (e *ExternallyOwnedAccount) DeepCopy() Account {
	return &ExternallyOwnedAccount{
		AccountCommon: e.AccountCommon.DeepCopy(),
		storageRoot:   e.storageRoot,
		codeHash:      common.CopyBytes(e.codeHash),
		codeInfo:      e.codeInfo,
	}
}

func (e *ExternallyOwnedAccount) GetStorageRoot() common.ExtHash {
	return e.storageRoot
}

func (e *ExternallyOwnedAccount) GetCodeHash() []byte {
	return e.codeHash
}

func (e *ExternallyOwnedAccount) GetCodeFormat() params.CodeFormat {
	return e.codeInfo.GetCodeFormat()
}

func (e *ExternallyOwnedAccount) GetVmVersion() params.VmVersion {
	return e.codeInfo.GetVmVersion()
}

func (e *ExternallyOwnedAccount) SetStorageRoot(h common.ExtHash) {
	e.storageRoot = h
}

func (e *ExternallyOwnedAccount) SetCodeHash(h []byte) {
	e.codeHash = h
}

func (e *ExternallyOwnedAccount) SetCodeInfo(ci params.CodeInfo) {
	e.codeInfo = ci
}

func (e *ExternallyOwnedAccount) toSerializable() *externallyOwnedAccountSerializable {
	return &externallyOwnedAccountSerializable{
		CommonSerializable: e.AccountCommon.toSerializable(),
		StorageRoot:        e.storageRoot.Unextend(),
		CodeHash:           e.codeHash,
		CodeInfo:           e.codeInfo,
	}
}

func (e *ExternallyOwnedAccount) toSerializableExt() *externallyOwnedAccountSerializableExt {
	return &externallyOwnedAccountSerializableExt{
		CommonSerializable: e.AccountCommon.toSerializable(),
		StorageRoot:        e.storageRoot,
		CodeHash:           e.codeHash,
		CodeInfo:           e.codeInfo,
	}
}

func (e *ExternallyOwnedAccount) fromSerializable(o *externallyOwnedAccountSerializable) {
	e.AccountCommon.fromSerializable(o.CommonSerializable)
	e.storageRoot = o.StorageRoot.ExtendZero()
	e.codeHash = o.CodeHash
	e.codeInfo = o.CodeInfo
}

func (e *ExternallyOwnedAccount) fromSerializableExt(o *externallyOwnedAccountSerializableExt) {
	e.AccountCommon.fromSerializable(o.CommonSerializable)
	e.storageRoot = o.StorageRoot
	e.codeHash = o.CodeHash
	e.codeInfo = o.CodeInfo
}

func (e *ExternallyOwnedAccount) EncodeRLP(w io.Writer) error {
	if !e.shouldUseCompactEncoding() {
		return rlp.Encode(w, e.toSerializable())
	}
	return rlp.Encode(w, e.AccountCommon.toSerializable())
}

func (e *ExternallyOwnedAccount) EncodeRLPExt(w io.Writer) error {
	var serializable interface{}

	switch {
	case !e.shouldUseCompactEncoding() && e.storageRoot.IsZeroExtended():
		serializable = e.toSerializable() // Full encoding. [n, b, hR, k, sR, cH, cI]
	case !e.shouldUseCompactEncoding():
		serializable = e.toSerializableExt() // Full encoding, has ExtHash [n, b, hR, k, sRExt, cH, cI]
	default:
		serializable = e.AccountCommon.toSerializable() // Compact encoding. [n, b, hR, k]
	}

	return rlp.Encode(w, serializable)
}

func (e *ExternallyOwnedAccount) shouldUseCompactEncoding() bool {
	zeroHash := common.ExtHash{}

	isEmptyStorageRoot := e.storageRoot == emptyRoot.ExtendZero() || e.storageRoot == zeroHash
	hasEmptyCodeHash := e.codeHash == nil || bytes.Equal(e.codeHash, emptyCodeHash)
	hasDefaultCodeInfo := e.codeInfo == params.CodeInfo(0)

	if isEmptyStorageRoot && hasEmptyCodeHash && hasDefaultCodeInfo {
		return true
	}

	return false
}

func (e *ExternallyOwnedAccount) setDefaultValues() {
	e.SetStorageRoot(emptyRoot.ExtendZero())
	e.SetCodeHash(emptyCodeHash)
	e.SetCodeInfo(params.CodeInfo(0))
}

func (e *ExternallyOwnedAccount) DecodeRLP(s *rlp.Stream) error {
	savedStream, err := s.Raw()
	if err != nil {
		return err
	}

	s.Reset(bytes.NewReader(savedStream), 0)
	if _, err := s.List(); err != nil {
		return err
	}

	kind, _, err := s.Kind()
	if err != nil {
		return err
	}

	s.Reset(bytes.NewReader(savedStream), 0)

	// full EOA encoding (the first element is a list)
	if kind == rlp.List {
		serializedExt := &externallyOwnedAccountSerializableExt{
			CommonSerializable: newAccountCommonSerializable(),
		}
		if err := s.Decode(serializedExt); err == nil {
			e.fromSerializableExt(serializedExt)
			if e.shouldUseCompactEncoding() {
				e.setDefaultValues()
			}
			return nil
		}

		s.Reset(bytes.NewReader(savedStream), 0)
		serialized := &externallyOwnedAccountSerializable{
			CommonSerializable: newAccountCommonSerializable(),
		}
		if err := s.Decode(serialized); err == nil {
			e.fromSerializable(serialized)
			return nil
		}

		s.Reset(bytes.NewReader(savedStream), 0)
		if err := e.AccountCommon.DecodeRLP(s); err != nil {
			return err
		}
		e.setDefaultValues()
		return nil
	}

	// compact EOA encoding
	if err := e.AccountCommon.DecodeRLP(s); err != nil {
		return err
	}

	e.setDefaultValues()
	return nil
}

func (e *ExternallyOwnedAccount) MarshalJSON() ([]byte, error) {
	if e.shouldUseCompactEncoding() {
		e.setDefaultValues()
	}

	return json.Marshal(&externallyOwnedAccountSerializableJSON{
		Nonce:         e.nonce,
		Balance:       (*hexutil.Big)(e.balance),
		HumanReadable: e.humanReadable,
		Key:           accountkey.NewAccountKeySerializerWithAccountKey(e.key),
		StorageRoot:   e.storageRoot.Unextend(), // Unextend for API compatibility
		CodeHash:      e.codeHash,
		CodeFormat:    e.codeInfo.GetCodeFormat(),
		VmVersion:     e.codeInfo.GetVmVersion(),
	})
}

func (e *ExternallyOwnedAccount) UnmarshalJSON(b []byte) error {
	serialized := &externallyOwnedAccountSerializableJSON{}

	if err := json.Unmarshal(b, serialized); err != nil {
		return err
	}

	e.nonce = serialized.Nonce
	e.balance = (*big.Int)(serialized.Balance)
	e.humanReadable = serialized.HumanReadable
	e.key = serialized.Key.GetKey()

	e.storageRoot = serialized.StorageRoot.ExtendZero()
	e.codeHash = serialized.CodeHash
	e.codeInfo = params.NewCodeInfo(serialized.CodeFormat, serialized.VmVersion)

	if e.shouldUseCompactEncoding() {
		e.setDefaultValues()
	}

	return nil
}
