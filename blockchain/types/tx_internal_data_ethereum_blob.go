// Copyright 2025 The Kaia Authors
// This file is part of the kaia library.
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
// Modified and improved for the Kaia development.

package types

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"reflect"
	"slices"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
)

const (
	// BlobSidecarVersion0 includes a single proof for verifying the entire blob
	// against its commitment. Used when the full blob is available and needs to
	// be checked as a whole.
	BlobSidecarVersion0 = byte(0)

	// BlobSidecarVersion1 includes multiple cell proofs for verifying specific
	// blob elements (cells). Used in scenarios like data availability sampling,
	// where only portions of the blob are verified individually.
	BlobSidecarVersion1 = byte(1)
)

type TxInternalDataEthereumBlob struct {
	ChainID      *uint256.Int
	AccountNonce uint64
	GasTipCap    *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap    *uint256.Int // a.k.a. maxFeePerGas
	GasLimit     uint64
	Recipient    common.Address
	Amount       *uint256.Int
	Payload      []byte
	AccessList   AccessList
	BlobFeeCap   *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes   []common.Hash

	// A blob transaction can optionally contain blobs. This field must be set when BlobTx
	// is used to create a transaction for signing.
	Sidecar *BlobTxSidecar `rlp:"-"`

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataEthereumBlobJSON struct {
	Type                 TxType           `json:"typeInt"`
	TypeStr              string           `json:"type"`
	ChainID              *hexutil.U256    `json:"chainId"`
	AccountNonce         hexutil.Uint64   `json:"nonce"`
	MaxPriorityFeePerGas *hexutil.U256    `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.U256    `json:"maxFeePerGas"`
	GasLimit             hexutil.Uint64   `json:"gas"`
	Recipient            common.Address   `json:"to"`
	Amount               *hexutil.U256    `json:"value"`
	Payload              hexutil.Bytes    `json:"input"`
	AccessList           AccessList       `json:"accessList"`
	BlobFeeCap           *hexutil.U256    `json:"blobFeeCap"`
	BlobHashes           []common.Hash    `json:"blobHashes"`
	Sidecar              *BlobTxSidecar   `json:"sidecar"`
	TxSignatures         TxSignaturesJSON `json:"signatures"`
	Hash                 *common.Hash     `json:"hash"`
}

type TxInternalDataEthereumBlobSerializable struct {
	ChainID      *uint256.Int
	AccountNonce uint64
	GasTipCap    *uint256.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap    *uint256.Int // a.k.a. maxFeePerGas
	GasLimit     uint64
	Recipient    common.Address
	Amount       *uint256.Int
	Payload      []byte
	AccessList   AccessList
	BlobFeeCap   *uint256.Int // a.k.a. maxFeePerBlobGas
	BlobHashes   []common.Hash

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int
}

// BlobTxSidecar contains the blobs of a blob transaction.
type BlobTxSidecar struct {
	Version     byte                 // Version
	Blobs       []kzg4844.Blob       // Blobs needed by the blob pool
	Commitments []kzg4844.Commitment // Commitments needed by the blob pool
	Proofs      []kzg4844.Proof      // Proofs needed by the blob pool
}

// NewBlobTxSidecar initialises the BlobTxSidecar object with the provided parameters.
func NewBlobTxSidecar(version byte, blobs []kzg4844.Blob, commitments []kzg4844.Commitment, proofs []kzg4844.Proof) *BlobTxSidecar {
	return &BlobTxSidecar{
		Version:     version,
		Blobs:       blobs,
		Commitments: commitments,
		Proofs:      proofs,
	}
}

// BlobHashes computes the blob hashes of the given blobs.
func (sc *BlobTxSidecar) BlobHashes() []common.Hash {
	hasher := sha256.New()
	h := make([]common.Hash, len(sc.Commitments))
	for i := range sc.Blobs {
		h[i] = kzg4844.CalcBlobHashV1(hasher, &sc.Commitments[i])
	}
	return h
}

// CellProofsAt returns the cell proofs for blob with index idx.
// This method is only valid for sidecars with version 1.
func (sc *BlobTxSidecar) CellProofsAt(idx int) ([]kzg4844.Proof, error) {
	if sc.Version != BlobSidecarVersion1 {
		return nil, fmt.Errorf("cell proof unsupported, version: %d", sc.Version)
	}
	if idx < 0 || idx >= len(sc.Blobs) {
		return nil, fmt.Errorf("cell proof out of bounds, index: %d, blobs: %d", idx, len(sc.Blobs))
	}
	index := idx * kzg4844.CellProofsPerBlob
	if len(sc.Proofs) < index+kzg4844.CellProofsPerBlob {
		return nil, fmt.Errorf("cell proof is corrupted, index: %d, proofs: %d", idx, len(sc.Proofs))
	}
	return sc.Proofs[index : index+kzg4844.CellProofsPerBlob], nil
}

// ToV1 converts the BlobSidecar to version 1, attaching the cell proofs.
func (sc *BlobTxSidecar) ToV1() error {
	if sc.Version == BlobSidecarVersion1 {
		return nil
	}
	if sc.Version == BlobSidecarVersion0 {
		proofs := make([]kzg4844.Proof, 0, len(sc.Blobs)*kzg4844.CellProofsPerBlob)
		for _, blob := range sc.Blobs {
			cellProofs, err := kzg4844.ComputeCellProofs(&blob)
			if err != nil {
				return err
			}
			proofs = append(proofs, cellProofs...)
		}
		sc.Version = BlobSidecarVersion1
		sc.Proofs = proofs
	}
	return nil
}

// encodedSize computes the RLP size of the sidecar elements. This does NOT return the
// encoded size of the BlobTxSidecar, it's just a helper for tx.Size().
func (sc *BlobTxSidecar) encodedSize() uint64 {
	var blobs, commitments, proofs uint64
	for i := range sc.Blobs {
		blobs += rlp.BytesSize(sc.Blobs[i][:])
	}
	for i := range sc.Commitments {
		commitments += rlp.BytesSize(sc.Commitments[i][:])
	}
	for i := range sc.Proofs {
		proofs += rlp.BytesSize(sc.Proofs[i][:])
	}
	return rlp.ListSize(blobs) + rlp.ListSize(commitments) + rlp.ListSize(proofs)
}

// ValidateBlobCommitmentHashes checks whether the given hashes correspond to the
// commitments in the sidecar
func (sc *BlobTxSidecar) ValidateBlobCommitmentHashes(hashes []common.Hash) error {
	if len(sc.Commitments) != len(hashes) {
		return fmt.Errorf("invalid number of %d blob commitments compared to %d blob hashes", len(sc.Commitments), len(hashes))
	}
	hasher := sha256.New()
	for i, vhash := range hashes {
		computed := kzg4844.CalcBlobHashV1(hasher, &sc.Commitments[i])
		if vhash != computed {
			return fmt.Errorf("blob %d: computed hash %#x mismatches transaction one %#x", i, computed, vhash)
		}
	}
	return nil
}

// Copy returns a deep-copied BlobTxSidecar object.
func (sc *BlobTxSidecar) Copy() *BlobTxSidecar {
	return &BlobTxSidecar{
		Version: sc.Version,

		// The element of these slice is fix-size byte array,
		// therefore slices.Clone will actually deep copy by value.
		Blobs:       slices.Clone(sc.Blobs),
		Commitments: slices.Clone(sc.Commitments),
		Proofs:      slices.Clone(sc.Proofs),
	}
}

// blobTxWithBlobs represents blob tx with its corresponding sidecar.
// This is an interface because sidecars are versioned.
type blobTxWithBlobs interface {
	tx() *TxInternalDataEthereumBlob
	assign(*BlobTxSidecar) error
}

type blobTxWithBlobsV0 struct {
	BlobTx      *TxInternalDataEthereumBlobSerializable
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

type blobTxWithBlobsV1 struct {
	BlobTx      *TxInternalDataEthereumBlobSerializable
	Version     byte
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func (btx *blobTxWithBlobsV0) tx() *TxInternalDataEthereumBlob {
	return btx.BlobTx.toBlobTx()
}

func (btx *blobTxWithBlobsV0) assign(sc *BlobTxSidecar) error {
	sc.Version = BlobSidecarVersion0
	sc.Blobs = btx.Blobs
	sc.Commitments = btx.Commitments
	sc.Proofs = btx.Proofs
	return nil
}

func (btx *blobTxWithBlobsV1) tx() *TxInternalDataEthereumBlob {
	return btx.BlobTx.toBlobTx()
}

func (btx *blobTxWithBlobsV1) assign(sc *BlobTxSidecar) error {
	if btx.Version != BlobSidecarVersion1 {
		return fmt.Errorf("unsupported blob tx version %d", btx.Version)
	}
	sc.Version = BlobSidecarVersion1
	sc.Blobs = btx.Blobs
	sc.Commitments = btx.Commitments
	sc.Proofs = btx.Proofs
	return nil
}

func newTxInternalDataEthereumBlob() *TxInternalDataEthereumBlob {
	return &TxInternalDataEthereumBlob{
		ChainID:      new(uint256.Int),
		AccountNonce: 0,
		GasTipCap:    new(uint256.Int),
		GasFeeCap:    new(uint256.Int),
		GasLimit:     0,
		Recipient:    common.Address{},
		Amount:       new(uint256.Int),
		Payload:      []byte{},
		AccessList:   AccessList{},
		BlobFeeCap:   new(uint256.Int),
		BlobHashes:   []common.Hash{},
		Sidecar:      nil,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
}

func newTxInternalDataEthereumBlobWithValues(nonce uint64, to common.Address, amount *big.Int, gasLimit uint64, gasTipCap *big.Int, gasFeeCap *big.Int, data []byte, accessList AccessList, chainID *big.Int, blobFeeCap *big.Int, blobHashes []common.Hash, sidecar *BlobTxSidecar) *TxInternalDataEthereumBlob {
	d := newTxInternalDataEthereumBlob()

	d.AccountNonce = nonce
	d.Recipient = to
	d.GasLimit = gasLimit

	if chainID != nil {
		d.ChainID.Set(uint256.MustFromBig(chainID))
	}

	if gasTipCap != nil {
		d.GasTipCap.Set(uint256.MustFromBig(gasTipCap))
	}

	if gasFeeCap != nil {
		d.GasFeeCap.Set(uint256.MustFromBig(gasFeeCap))
	}

	if amount != nil {
		d.Amount.Set(uint256.MustFromBig(amount))
	}

	if len(data) > 0 {
		d.Payload = common.CopyBytes(data)
	}

	if accessList != nil {
		d.AccessList = append(d.AccessList, accessList...)
	}

	if blobFeeCap != nil {
		d.BlobFeeCap.Set(uint256.MustFromBig(blobFeeCap))
	}

	if len(blobHashes) > 0 {
		d.BlobHashes = append(d.BlobHashes, blobHashes...)
	}

	if sidecar != nil {
		d.Sidecar = sidecar.Copy()
	}

	return d
}

func newTxInternalDataEthereumBlobWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataEthereumBlob, error) {
	d := newTxInternalDataEthereumBlob()

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
		return nil, errValueKeyToMustAddressPointer
	}

	if v, ok := values[TxValueKeyAmount].(*big.Int); ok {
		d.Amount.Set(uint256.MustFromBig(v))
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
		d.GasFeeCap.Set(uint256.MustFromBig(v))
		delete(values, TxValueKeyGasFeeCap)
	} else {
		return nil, errValueKeyGasFeeCapMustBigInt
	}
	if v, ok := values[TxValueKeyGasTipCap].(*big.Int); ok {
		d.GasTipCap.Set(uint256.MustFromBig(v))
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
	if v, ok := values[TxValueKeyBlobFeeCap].(*big.Int); ok {
		d.BlobFeeCap.Set(uint256.MustFromBig(v))
		delete(values, TxValueKeyBlobFeeCap)
	} else {
		return nil, errValueKeyBlobFeeCapMustBigInt
	}
	if v, ok := values[TxValueKeyBlobHashes].([]common.Hash); ok {
		d.BlobHashes = make([]common.Hash, len(v))
		copy(d.BlobHashes, v)
		delete(values, TxValueKeyBlobHashes)
	} else {
		return nil, errValueKeyDataMustByteSlice
	}
	if v, ok := values[TxValueKeySidecar].(*BlobTxSidecar); ok {
		d.Sidecar = v.Copy()
		delete(values, TxValueKeySidecar)
	} else {
		return nil, errValueKeySidecarInvalid
	}
	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataEthereumBlobSerializable) toBlobTx() *TxInternalDataEthereumBlob {
	return &TxInternalDataEthereumBlob{
		ChainID:      t.ChainID,
		AccountNonce: t.AccountNonce,
		GasTipCap:    t.GasTipCap,
		GasFeeCap:    t.GasFeeCap,
		GasLimit:     t.GasLimit,
		Recipient:    t.Recipient,
		Amount:       t.Amount,
		Payload:      t.Payload,
		AccessList:   t.AccessList,
		BlobFeeCap:   t.BlobFeeCap,
		BlobHashes:   t.BlobHashes,
		V:            t.V,
		R:            t.R,
		S:            t.S,
	}
}

func (t *TxInternalDataEthereumBlob) toBlobTxSerializable() *TxInternalDataEthereumBlobSerializable {
	return &TxInternalDataEthereumBlobSerializable{
		ChainID:      t.ChainID,
		AccountNonce: t.AccountNonce,
		GasTipCap:    t.GasTipCap,
		GasFeeCap:    t.GasFeeCap,
		GasLimit:     t.GasLimit,
		Recipient:    t.Recipient,
		Amount:       t.Amount,
		Payload:      t.Payload,
		AccessList:   t.AccessList,
		BlobFeeCap:   t.BlobFeeCap,
		BlobHashes:   t.BlobHashes,
		V:            t.V,
		R:            t.R,
		S:            t.S,
	}
}

func (t *TxInternalDataEthereumBlob) Type() TxType {
	return TxTypeEthereumBlob
}

func (t *TxInternalDataEthereumBlob) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataEthereumBlob) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataEthereumBlob) GetPrice() *big.Int {
	return new(big.Int).Set(t.GasFeeCap.ToBig())
}

func (t *TxInternalDataEthereumBlob) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataEthereumBlob) GetRecipient() *common.Address {
	return &t.Recipient
}

func (t *TxInternalDataEthereumBlob) GetAmount() *big.Int {
	return new(big.Int).Set(t.Amount.ToBig())
}

func (t *TxInternalDataEthereumBlob) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataEthereumBlob) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataEthereumBlob) GetAccessList() AccessList {
	return t.AccessList
}

func (t *TxInternalDataEthereumBlob) GetGasTipCap() *big.Int {
	return new(big.Int).Set(t.GasTipCap.ToBig())
}

func (t *TxInternalDataEthereumBlob) GetGasFeeCap() *big.Int {
	return new(big.Int).Set(t.GasFeeCap.ToBig())
}

func (t *TxInternalDataEthereumBlob) GetBlobGas() *big.Int {
	return new(big.Int).SetUint64(params.BlobTxBlobGasPerBlob * uint64(len(t.BlobHashes)))
}

func (t *TxInternalDataEthereumBlob) SetHash(hash *common.Hash) {
	t.Hash = hash
}

func (t *TxInternalDataEthereumBlob) SetSignature(signatures TxSignatures) {
	if len(signatures) != 1 {
		logger.Crit("TxInternalDataEthereumBlob can receive only single signature!")
	}

	t.V = signatures[0].V
	t.R = signatures[0].R
	t.S = signatures[0].S
}

func (t *TxInternalDataEthereumBlob) RawSignatureValues() TxSignatures {
	return TxSignatures{&TxSignature{t.V, t.R, t.S}}
}

func (t *TxInternalDataEthereumBlob) ValidateSignature() bool {
	v := byte(t.V.Uint64())
	return crypto.ValidateSignatureValues(v, t.R, t.S, false)
}

func (t *TxInternalDataEthereumBlob) RecoverAddress(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) (common.Address, error) {
	V := vfunc(t.V)
	return recoverPlain(txhash, t.R, t.S, V, homestead)
}

func (t *TxInternalDataEthereumBlob) RecoverPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	V := vfunc(t.V)

	pk, err := recoverPlainPubkey(txhash, t.R, t.S, V, homestead)
	if err != nil {
		return nil, err
	}

	return []*ecdsa.PublicKey{pk}, nil
}

func (t *TxInternalDataEthereumBlob) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	return IntrinsicGas(t.Payload, t.AccessList, nil, false, *fork.Rules(big.NewInt(int64(currentBlockNumber))))
}

func (t *TxInternalDataEthereumBlob) ChainId() *big.Int {
	return t.ChainID.ToBig()
}

func (t *TxInternalDataEthereumBlob) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataEthereumBlob)
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
		t.BlobFeeCap.Cmp(ta.BlobFeeCap) == 0 &&
		reflect.DeepEqual(t.BlobHashes, ta.BlobHashes) &&
		reflect.DeepEqual(t.Sidecar, ta.Sidecar) &&
		t.V.Cmp(ta.V) == 0 &&
		t.R.Cmp(ta.R) == 0 &&
		t.S.Cmp(ta.S) == 0
}

func (t *TxInternalDataEthereumBlob) String() string {
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
		BlobFeeCap: %#x
		BlobHashes: %x
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
		t.BlobFeeCap,
		t.BlobHashes,
		v,
		r,
		s,
		enc,
	)
}

func (t *TxInternalDataEthereumBlob) SerializeForSign() []interface{} {
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
		t.BlobFeeCap,
		t.BlobHashes,
	}
}

func (t *TxInternalDataEthereumBlob) TxHash() common.Hash {
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
		t.BlobFeeCap,
		t.BlobHashes,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumBlob) SenderTxHash() common.Hash {
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
		t.BlobFeeCap,
		t.BlobHashes,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumBlob) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	if common.IsPrecompiledContractAddress(t.Recipient, *fork.Rules(big.NewInt(int64(currentBlockNumber)))) {
		return kerrors.ErrPrecompiledContractAddress
	}
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataEthereumBlob) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataEthereumBlob) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataEthereumBlob) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//start := time.Now()
	//defer func() {
	//	elapsed := time.Since(start)
	//	logger.Debug("[TxInternalDataLegacy] EVM execution done", "elapsed", elapsed)
	//}()
	///////////////////////////////////////////////////////
	stateDB.IncNonce(sender.Address())
	return vm.Call(sender, t.Recipient, t.Payload, gas, value)
}

func (t *TxInternalDataEthereumBlob) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"typeInt":              t.Type(),
		"type":                 t.Type().String(),
		"chainId":              (*hexutil.Big)(t.ChainId()),
		"nonce":                hexutil.Uint64(t.AccountNonce),
		"maxPriorityFeePerGas": (*hexutil.Big)(t.GasTipCap.ToBig()),
		"maxFeePerGas":         (*hexutil.Big)(t.GasFeeCap.ToBig()),
		"gas":                  hexutil.Uint64(t.GasLimit),
		"to":                   t.Recipient,
		"input":                hexutil.Bytes(t.Payload),
		"value":                (*hexutil.Big)(t.Amount.ToBig()),
		"accessList":           t.AccessList,
		"blobFeeCap":           (*hexutil.Big)(t.BlobFeeCap.ToBig()),
		"blobHashes":           t.BlobHashes,
		"sidecar":              t.Sidecar,
		"signatures":           TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
	}
}

func (t *TxInternalDataEthereumBlob) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataEthereumBlobJSON{
		t.Type(),
		t.Type().String(),
		(*hexutil.U256)(t.ChainID),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.U256)(t.GasTipCap),
		(*hexutil.U256)(t.GasFeeCap),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.U256)(t.Amount),
		t.Payload,
		t.AccessList,
		(*hexutil.U256)(t.BlobFeeCap),
		t.BlobHashes,
		t.Sidecar,
		TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
		t.Hash,
	})
}

func (t *TxInternalDataEthereumBlob) UnmarshalJSON(bytes []byte) error {
	js := &TxInternalDataEthereumBlobJSON{}
	if err := json.Unmarshal(bytes, js); err != nil {
		return err
	}

	t.ChainID = (*uint256.Int)(js.ChainID)
	t.AccountNonce = uint64(js.AccountNonce)
	t.GasTipCap = (*uint256.Int)(js.MaxPriorityFeePerGas)
	t.GasFeeCap = (*uint256.Int)(js.MaxFeePerGas)
	t.GasLimit = uint64(js.GasLimit)
	t.Recipient = js.Recipient
	t.Amount = (*uint256.Int)(js.Amount)
	t.Payload = js.Payload
	t.AccessList = js.AccessList
	t.BlobFeeCap = (*uint256.Int)(js.BlobFeeCap)
	t.BlobHashes = js.BlobHashes
	t.Sidecar = js.Sidecar
	t.V = (*big.Int)(js.TxSignatures[0].V)
	t.R = (*big.Int)(js.TxSignatures[0].R)
	t.S = (*big.Int)(js.TxSignatures[0].S)
	t.Hash = js.Hash

	return nil
}

func (t *TxInternalDataEthereumBlob) setSignatureValues(chainID, v, r, s *big.Int) {
	t.ChainID, t.V, t.R, t.S = uint256.MustFromBig(chainID), v, r, s
}

func (tx *TxInternalDataEthereumBlob) withoutSidecar() *TxInternalDataEthereumBlob {
	cpy := *tx
	cpy.Sidecar = nil
	return &cpy
}

func (tx *TxInternalDataEthereumBlob) withSidecar(sideCar *BlobTxSidecar) *TxInternalDataEthereumBlob {
	cpy := *tx
	cpy.Sidecar = sideCar
	return &cpy
}

func (tx *TxInternalDataEthereumBlob) EncodeRLP(w io.Writer) error {
	switch {
	case tx.Sidecar == nil:
		return rlp.Encode(w, tx)

	case tx.Sidecar.Version == BlobSidecarVersion0:
		return rlp.Encode(w, &blobTxWithBlobsV0{
			BlobTx:      tx.toBlobTxSerializable(),
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})

	case tx.Sidecar.Version == BlobSidecarVersion1:
		return rlp.Encode(w, &blobTxWithBlobsV1{
			BlobTx:      tx.toBlobTxSerializable(),
			Version:     tx.Sidecar.Version,
			Blobs:       tx.Sidecar.Blobs,
			Commitments: tx.Sidecar.Commitments,
			Proofs:      tx.Sidecar.Proofs,
		})

	default:
		return errors.New("unsupported sidecar version")
	}
}

// Decode various BlobTx encoding.
// no sidecar: [chainID, nonce, ...]
// sidecar v0: [[chainID, nonce, ...], blobs, commitments, proofs]
// sidecar v1: [[chainID, nonce, ...], version, blobs, commitments, proofs]
func (tx *TxInternalDataEthereumBlob) DecodeRLP(s *rlp.Stream) error {
	// Here we need to support two outer formats: the network protocol encoding of the tx
	// (with blobs) or the canonical encoding without blobs.
	//
	// The canonical encoding is just a list of fields:
	//
	//     [chainID, nonce, ...]
	//
	// The network encoding is a list where the first element is the tx in the canonical encoding,
	// and the remaining elements are the 'sidecar':
	//
	//     [[chainID, nonce, ...], ...]
	//
	// The two outer encodings can be distinguished by checking whether the first element
	// of the input list is itself a list. If it's the canonical encoding, the first
	// element is the chainID, which is a number.

	input, err := s.Raw()
	if err != nil {
		return err
	}
	firstElem, _, err := rlp.SplitList(input)
	if err != nil {
		return err
	}
	firstElemKind, _, secondElem, err := rlp.Split(firstElem)
	if err != nil {
		return err
	}
	if firstElemKind != rlp.List {
		// Blob tx without blobs.
		return rlp.DecodeBytes(input, tx)
	}

	// Now we know it's the network encoding with the blob sidecar. Here we again need to
	// support multiple encodings: legacy sidecars (v0) with a blob proof, and versioned
	// sidecars.
	//
	// The legacy encoding is:
	//
	//     [tx, blobs, commitments, proofs]
	//
	// The versioned encoding is:
	//
	//     [tx, version, blobs, ...]
	//
	// We can tell the two apart by checking whether the second element is the version byte.
	// For legacy sidecar the second element is a list of blobs.

	secondElemKind, _, _, err := rlp.Split(secondElem)
	if err != nil {
		return err
	}
	var payload blobTxWithBlobs
	if secondElemKind == rlp.List {
		// No version byte: blob sidecar v0.
		payload = new(blobTxWithBlobsV0)
	} else {
		// It has a version byte. Decode as v1, version is checked by assign()
		payload = new(blobTxWithBlobsV1)
	}
	if err := rlp.DecodeBytes(input, payload); err != nil {
		return err
	}
	sc := new(BlobTxSidecar)
	if err := payload.assign(sc); err != nil {
		return err
	}
	*tx = *payload.tx()
	tx.Sidecar = sc
	return nil
}
