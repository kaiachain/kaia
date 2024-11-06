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
	"fmt"
	"io"
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
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
		common.ExtHash{},
		emptyCodeHash,
		params.CodeInfo(0),
	}
}

// newExternallyOwnedAccountWithMap creates an ExternallyOwnedAccount object initialized with the given values.
func newExternallyOwnedAccountWithMap(values map[AccountValueKeyType]interface{}) *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		newAccountCommonWithMap(values),
		common.ExtHash{},
		emptyCodeHash,
		params.CodeInfo(0),
	}
}

func (e *ExternallyOwnedAccount) Type() AccountType {
	return ExternallyOwnedAccountType
}

func (e *ExternallyOwnedAccount) Dump() {
	fmt.Println(e.String())
}

func (e *ExternallyOwnedAccount) String() string {
	return fmt.Sprintf("EOA: %s", e.AccountCommon.String())
}

func (e *ExternallyOwnedAccount) DeepCopy() Account {
	return &ExternallyOwnedAccount{
		AccountCommon: e.AccountCommon.DeepCopy(),
	}
}

func (e *ExternallyOwnedAccount) Equal(a Account) bool {
	e2, ok := a.(*ExternallyOwnedAccount)
	if !ok {
		return false
	}

	return e.AccountCommon.Equal(e2.AccountCommon)
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
	if isEOAWithCode(e) {
		return rlp.Encode(w, e.toSerializable())
	}
	return rlp.Encode(w, e.AccountCommon.toSerializable())
}

func (e *ExternallyOwnedAccount) EncodeRLPExt(w io.Writer) error {
	if isEOAWithCode(e) {
		if e.storageRoot.IsZeroExtended() {
			return rlp.Encode(w, e.toSerializable()) // EOA without code. [n, b, hR, k, 0x0, eCH, 0x0]
		}
		return rlp.Encode(w, e.toSerializableExt()) // EOA with code, has ExtHash [n, b, hR, k, sRExt, cH, cI]
	}
	return rlp.Encode(w, e.AccountCommon.toSerializable()) // EOA without code. [n, b, hR, k]
}

func isEOAWithCode(e *ExternallyOwnedAccount) bool {
	zeroHash := common.ExtHash{}
	if e.codeHash == nil { // EOA without code. [n, b, hR, k] (Backwards Compatibility)
		return false
	} else if e.storageRoot == zeroHash && e.storageRoot.IsZeroExtended() && e.codeInfo == params.CodeInfo(0) { // EOA without code. [n, b, hR, k, 0x0, eCH, 0x0]
		return false
	} else if e.storageRoot == zeroHash || bytes.Equal(e.storageRoot.Bytes(), emptyRoot.Bytes()) {
		return true
	} else if bytes.Equal(e.codeHash, emptyCodeHash) || e.codeHash != nil {
		return true
	} else if e.codeInfo == params.CodeInfo(0) {
		return true
	}
	return false
}

func (e *ExternallyOwnedAccount) DecodeRLP(s *rlp.Stream) error {
	// Save original stream data
	savedStream, err := s.Raw()
	if err != nil {
		return err
	}

	// Reset and check the first list structure
	s.Reset(bytes.NewReader(savedStream), 0)
	if _, err := s.List(); err != nil {
		return err
	}

	// Check the type of first element inside the list
	kind, _, err := s.Kind()
	if err != nil {
		return err
	}

	// Reset stream for actual decoding
	s.Reset(bytes.NewReader(savedStream), 0)

	if kind == rlp.List {
		// Case 2: First element is a list - use extended format decoding
		serializedExt := &externallyOwnedAccountSerializableExt{
			CommonSerializable: newAccountCommonSerializable(),
		}
		if err := s.Decode(serializedExt); err == nil {
			e.fromSerializableExt(serializedExt)
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
		return err
	} else {
		// Case 1: First element is not a list - use AccountCommon.DecodeRLP
		return e.AccountCommon.DecodeRLP(s)
	}
}

func (e *ExternallyOwnedAccount) MarshalJSON() ([]byte, error) {
	if e.codeHash != nil {
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
	return e.AccountCommon.MarshalJSON()
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
	e.storageRoot = serialized.StorageRoot.ExtendZero() // API inputs should contain merkle hash
	e.codeHash = serialized.CodeHash
	// e.codeInfo = params.NewCodeInfo(serialized.CodeFormat, serialized.VmVersion)

	return nil
}
