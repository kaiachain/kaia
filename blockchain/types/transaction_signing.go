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
// This file is derived from core/types/transaction_signing.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package types

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
)

var (
	ErrInvalidChainId        = errors.New("invalid chain id for signer")
	errNotTxInternalDataFrom = errors.New("not an TxInternalDataFrom")
)

// sigCache is used to cache the derived sender and contains
// the signer used to derive it.
type sigCache struct {
	signer Signer
	from   common.Address
}

// sigCachePubkey is used to cache the derived public key and contains
// the signer used to derive it.
type sigCachePubkey struct {
	signer Signer
	pubkey []*ecdsa.PublicKey
}

// MakeSigner returns a Signer based on the given chain config and block number.
func MakeSigner(config *params.ChainConfig, blockNumber *big.Int) Signer {
	var signer Signer

	if config.IsOsakaForkEnabled(blockNumber) {
		signer = NewOsakaSigner(config.ChainID)
	} else if config.IsPragueForkEnabled(blockNumber) {
		signer = NewPragueSigner(config.ChainID)
	} else if config.IsEthTxTypeForkEnabled(blockNumber) {
		signer = NewLondonSigner(config.ChainID)
	} else {
		signer = NewEIP155Signer(config.ChainID)
	}

	return signer
}

// LatestSigner returns the 'most permissive' Signer available for the given chain
// configuration. Specifically, this enables support of all types of transactions
// when their respective forks are scheduled to occur at any block number
// in the chain config.
//
// Use this in transaction-handling code where the current block number is unknown. If you
// have the current block number available, use MakeSigner instead.
func LatestSigner(config *params.ChainConfig) Signer {
	// Be aware that it checks whether EthTxTypeCompatibleBlock or PragueCompatible is set,
	// but doesn't check whether it is enabled on a specific block number.
	if config.OsakaCompatibleBlock != nil {
		return NewOsakaSigner(config.ChainID)
	} else if config.PragueCompatibleBlock != nil {
		return NewPragueSigner(config.ChainID)
	} else if config.EthTxTypeCompatibleBlock != nil {
		return NewLondonSigner(config.ChainID)
	}

	return NewEIP155Signer(config.ChainID)
}

// LatestSignerForChainID returns the 'most permissive' Signer available. Specifically,
// this enables support for EIP-155 replay protection and all implemented EIP-2718
// transaction types if chainID is non-nil.
//
// Use this in transaction-handling code where the current block number and fork
// configuration are unknown. If you have a ChainConfig, use LatestSigner instead.
// If you have a ChainConfig and know the current block number, use MakeSigner instead.
func LatestSignerForChainID(chainID *big.Int) Signer {
	return NewOsakaSigner(chainID)
}

// SignTx signs the transaction using the given signer and private key
func SignTx(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h := s.Hash(tx)
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}

	return tx.WithSignature(s, sig)
}

// SignTxAsFeePayer signs the transaction as a fee payer using the given signer and private key
func SignTxAsFeePayer(tx *Transaction, s Signer, prv *ecdsa.PrivateKey) (*Transaction, error) {
	h, err := s.HashFeePayer(tx)
	if err != nil {
		return nil, err
	}
	sig, err := crypto.Sign(h[:], prv)
	if err != nil {
		return nil, err
	}
	return tx.WithFeePayerSignature(s, sig)
}

// AccountKeyPicker has a function GetKey() to retrieve an account key from statedb.
type AccountKeyPicker interface {
	GetKey(address common.Address) accountkey.AccountKey
	Exist(addr common.Address) bool
}

// Sender returns the address of the transaction.
// If an ethereum transaction, it calls SenderFrom().
// Otherwise, it just returns tx.From() because the other transaction types have the field `from`.
// NOTE: this function should not be called if tx signature validation is required.
// In that situtation, you should call ValidateSender().
func Sender(signer Signer, tx *Transaction) (common.Address, error) {
	if tx.IsEthereumTransaction() {
		return SenderFrom(signer, tx)
	}

	return tx.From()
}

// SenderFeePayer returns the fee payer address of the transaction.
// If the transaction is not a fee-delegated transaction, the fee payer is set to
// the address of the `from` of the transaction.
func SenderFeePayer(signer Signer, tx *Transaction) (common.Address, error) {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return Sender(signer, tx)
	}
	return tf.GetFeePayer(), nil
}

// SenderFeePayerPubkey returns the public key derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// SenderFeePayerPubkey may cache the public key, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func SenderFeePayerPubkey(signer Signer, tx *Transaction) ([]*ecdsa.PublicKey, error) {
	if sc := tx.feePayer.Load(); sc != nil {
		sigCache := sc.(sigCachePubkey)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.pubkey, nil
		}
	}

	pubkey, err := signer.SenderFeePayer(tx)
	if err != nil {
		return nil, err
	}

	tx.feePayer.Store(sigCachePubkey{signer: signer, pubkey: pubkey})
	return pubkey, nil
}

// SenderFrom returns the address derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// SenderFrom may cache the address, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func SenderFrom(signer Signer, tx *Transaction) (common.Address, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCache)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.from, nil
		}
	}

	addr, err := signer.Sender(tx)
	if err != nil {
		return common.Address{}, err
	}
	tx.from.Store(sigCache{signer: signer, from: addr})
	return addr, nil
}

// SenderPubkey returns the public key derived from the signature (V, R, S) using secp256k1
// elliptic curve and an error if it failed deriving or upon an incorrect
// signature.
//
// SenderPubkey may cache the public key, allowing it to be used regardless of
// signing method. The cache is invalidated if the cached signer does
// not match the signer used in the current call.
func SenderPubkey(signer Signer, tx *Transaction) ([]*ecdsa.PublicKey, error) {
	if sc := tx.from.Load(); sc != nil {
		sigCache := sc.(sigCachePubkey)
		// If the signer used to derive from in a previous
		// call is not the same as used current, invalidate
		// the cache.
		if sigCache.signer.Equal(signer) {
			return sigCache.pubkey, nil
		}
	}

	pubkey, err := signer.SenderPubkey(tx)
	if err != nil {
		return nil, err
	}
	tx.from.Store(sigCachePubkey{signer: signer, pubkey: pubkey})
	return pubkey, nil
}

// Signer encapsulates transaction signature handling. Note that this interface is not a
// stable API and may change at any time to accommodate new protocol rules.
type Signer interface {
	// Sender returns the sender address of the transaction.
	Sender(tx *Transaction) (common.Address, error)
	// SenderPubkey returns the public key derived from tx signature and txhash.
	SenderPubkey(tx *Transaction) ([]*ecdsa.PublicKey, error)
	// SenderFeePayer returns the public key derived from tx signature and txhash.
	SenderFeePayer(tx *Transaction) ([]*ecdsa.PublicKey, error)
	// SignatureValues returns the raw R, S, V values corresponding to the
	// given signature.
	SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error)
	// ChainID returns the chain id.
	ChainID() *big.Int
	// Hash returns 'signature hash', i.e. the transaction hash that is signed by the
	// private key. This hash does not uniquely identify the transaction.
	Hash(tx *Transaction) common.Hash
	// HashFeePayer returns the hash with a fee payer's address to be signed by a fee payer.
	HashFeePayer(tx *Transaction) (common.Hash, error)
	// Equal returns true if the given signer is the same as the receiver.
	Equal(Signer) bool
}

// modernSigner is a unified signer that can handle all eth typed transactions.
// It uses a map to track supported eth typed transaction types.
// Legacy transactions are all enabled by default and handled by the legacy signer.
type modernSigner struct {
	chainId   *big.Int
	legacy    EIP155Signer
	supported map[TxType]bool
}

// NewModernSigner creates a modern signer that supports the specified eth typed transaction types.
// Supported eth typed transaction types are passed as a map.
func NewModernSigner(chainId *big.Int, supported map[TxType]bool) Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return modernSigner{
		chainId:   chainId,
		legacy:    NewEIP155Signer(chainId),
		supported: supported,
	}
}

// ChainID returns the chain id.
func (s modernSigner) ChainID() *big.Int {
	return s.chainId
}

// Equal returns true if the given signer is the same as the receiver.
// All modern signers with the same chain ID are considered equal, which allows
// cache sharing between different fork configurations.
func (s modernSigner) Equal(s2 Signer) bool {
	ms, ok := s2.(modernSigner)
	return ok && ms.chainId.Cmp(s.chainId) == 0 && reflect.DeepEqual(ms.supported, s.supported) && s.legacy.Equal(ms.legacy)
}

// IsSupported returns true if the signer supports the given transaction type.
func (s modernSigner) IsSupported(tx *Transaction) bool {
	if !tx.Type().IsEthTypedTransaction() {
		return true // Legacy transactions are always supported
	}
	return s.supported[tx.Type()]
}

func (s modernSigner) Sender(tx *Transaction) (common.Address, error) {
	// Check if this transaction type is supported
	if !s.IsSupported(tx) {
		return common.Address{}, ErrTxTypeNotSupported
}

	// Legacy transactions are handled by the legacy signer
	if !tx.Type().IsEthTypedTransaction() {
		return s.legacy.Sender(tx)
	}

	// Validate chain ID
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}

	// All modern transaction types use the same recovery mechanism
	return RecoverTxSender(s.Hash(tx), tx.data.RawSignatureValues(), true, func(v *big.Int) *big.Int {
		// Modern txs use 0 and 1 as recovery id, add 27 to become equivalent to unprotected Homestead signatures.
		V := new(big.Int).Add(v, big.NewInt(27))
		return V
	})
}

func (s modernSigner) SenderPubkey(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	// Check if this transaction type is supported
	if !s.IsSupported(tx) {
		return nil, ErrTxTypeNotSupported
	}

	// Legacy transactions are handled by the legacy signer
	if !tx.Type().IsEthTypedTransaction() {
		return s.legacy.SenderPubkey(tx)
	}

	// Validate chain ID
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, ErrInvalidChainId
	}

	// All modern transaction types use the same recovery mechanism
	return RecoverTxPubkeys(s.Hash(tx), tx.data.RawSignatureValues(), true, func(v *big.Int) *big.Int {
		// Modern txs use 0 and 1 as recovery id, add 27 to become equivalent to unprotected Homestead signatures.
		V := new(big.Int).Add(v, big.NewInt(27))
		return V
	})
}

func (s modernSigner) SenderFeePayer(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	// Check if this transaction type is supported
	if !s.IsSupported(tx) {
		return nil, ErrTxTypeNotSupported
}

	return s.legacy.SenderFeePayer(tx)
}

func (s modernSigner) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	// Check if this transaction type is supported
	if !s.IsSupported(tx) {
		return nil, nil, nil, ErrTxTypeNotSupported
}

	// Legacy transactions are handled by the legacy signer
	if !tx.Type().IsEthTypedTransaction() {
		return s.legacy.SignatureValues(tx, sig)
	}

	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}

	// Check that chain ID of tx matches the signer. We also accept ID zero or nil here,
	// because it indicates that the chain ID was not specified in the tx.
	// NOTE: Kaia allow chain ID to be nil in this fork
	if tx.data.ChainId() != nil && tx.data.ChainId().Sign() != 0 && tx.data.ChainId().Cmp(s.ChainID()) != 0 {
		return nil, nil, nil, ErrInvalidChainId
	}

	R = new(big.Int).SetBytes(sig[:32])
	S = new(big.Int).SetBytes(sig[32:64])
	V = big.NewInt(int64(sig[crypto.RecoveryIDOffset]))

	return R, S, V, nil
}

func (s modernSigner) Hash(tx *Transaction) common.Hash {
	if !s.IsSupported(tx) {
		return common.Hash{}
	}

	return tx.data.SigHash(s.ChainID())
}

func (s modernSigner) HashFeePayer(tx *Transaction) (common.Hash, error) {
	if !s.IsSupported(tx) {
		return common.Hash{}, ErrTxTypeNotSupported
	}

	return s.legacy.HashFeePayer(tx)
}

// NewOsakaSigner returns a signer that accepts
// - EIP-4844 blob transactions,
// - EIP-7702 set code transactions,
// - EIP-1559 dynamic fee transactions,
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions and
// - legacy transactions
func NewOsakaSigner(chainId *big.Int) Signer {
	supported := map[TxType]bool{
		TxTypeEthereumBlob:       true,
		TxTypeEthereumSetCode:    true,
		TxTypeEthereumDynamicFee: true,
		TxTypeEthereumAccessList: true,
	}
	return NewModernSigner(chainId, supported)
}

// NewPragueSigner returns a signer that accepts
// - EIP-7702 set code transactions,
// - EIP-1559 dynamic fee transactions,
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions and
// - legacy transactions
func NewPragueSigner(chainId *big.Int) Signer {
	supported := map[TxType]bool{
		TxTypeEthereumSetCode:    true,
		TxTypeEthereumDynamicFee: true,
		TxTypeEthereumAccessList: true,
	}
	return NewModernSigner(chainId, supported)
}

// NewLondonSigner returns a signer that accepts
// - EIP-1559 dynamic fee transactions,
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions and
// - legacy transactions
func NewLondonSigner(chainId *big.Int) Signer {
	supported := map[TxType]bool{
		TxTypeEthereumDynamicFee: true,
		TxTypeEthereumAccessList: true,
	}
	return NewModernSigner(chainId, supported)
	}

// New2930Signer returns a signer that accepts
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions and
// - legacy transactions
func NewEIP2930Signer(chainId *big.Int) Signer {
	supported := map[TxType]bool{
		TxTypeEthereumAccessList: true,
	}
	return NewModernSigner(chainId, supported)
}

type FrontierSigner struct{}

func (s FrontierSigner) ChainID() *big.Int {
	return nil
}

func (s FrontierSigner) Equal(s2 Signer) bool {
	_, ok := s2.(FrontierSigner)
	return ok
}

func (fs FrontierSigner) Sender(tx *Transaction) (common.Address, error) {
	if !tx.IsLegacyTransaction() {
		return common.Address{}, ErrTxTypeNotSupported
	}
	sigs := tx.RawSignatureValues()
	if len(sigs) != 1 {
		return common.Address{}, ErrShouldBeSingleSignature
	}
	v, r, s := sigs[0].V, sigs[0].R, sigs[0].S
	return recoverPlain(fs.Hash(tx), r, s, v, false)
}

func (fs FrontierSigner) SenderPubkey(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	return nil, ErrSenderPubkeyNotSupported
}

func (fs FrontierSigner) SenderFeePayer(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	return nil, ErrSenderFeePayerNotSupported
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (fs FrontierSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	if tx.IsLegacyTransaction() {
		return nil, nil, nil, ErrTxTypeNotSupported
	}
	r, s, v = decodeSignature(sig)
	return r, s, v, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (fs FrontierSigner) Hash(tx *Transaction) common.Hash {
	return rlpHash([]interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
	})
}

func (fs FrontierSigner) HashFeePayer(tx *Transaction) (common.Hash, error) {
	return common.Hash{}, ErrHashFeePayerNotSupported
}

// HomesteadTransaction implements TransactionInterface using the
// homestead rules.
type HomesteadSigner struct{ FrontierSigner }

func (s HomesteadSigner) ChainID() *big.Int {
	return nil
}

func (s HomesteadSigner) Equal(s2 Signer) bool {
	_, ok := s2.(HomesteadSigner)
	return ok
}

// SignatureValues returns signature values. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (hs HomesteadSigner) SignatureValues(tx *Transaction, sig []byte) (r, s, v *big.Int, err error) {
	return hs.FrontierSigner.SignatureValues(tx, sig)
}

func (hs HomesteadSigner) Sender(tx *Transaction) (common.Address, error) {
	if !tx.IsLegacyTransaction() {
		return common.Address{}, ErrTxTypeNotSupported
	}
	sigs := tx.RawSignatureValues()
	if len(sigs) != 1 {
		return common.Address{}, ErrShouldBeSingleSignature
	}
	v, r, s := sigs[0].V, sigs[0].R, sigs[0].S

	return recoverPlain(hs.Hash(tx), r, s, v, true)
}

// EIP155Transaction implements Signer using the EIP155 rules.
type EIP155Signer struct {
	chainId, chainIdMul *big.Int
}

func NewEIP155Signer(chainId *big.Int) EIP155Signer {
	if chainId == nil {
		chainId = new(big.Int)
	}
	return EIP155Signer{
		chainId:    chainId,
		chainIdMul: new(big.Int).Mul(chainId, common.Big2),
	}
}

// ChainID returns the chain id.
func (s EIP155Signer) ChainID() *big.Int {
	return s.chainId
}

func (s EIP155Signer) Equal(s2 Signer) bool {
	eip155, ok := s2.(EIP155Signer)
	return ok && eip155.chainId.Cmp(s.chainId) == 0
}

var big8 = big.NewInt(8)

func (s EIP155Signer) Sender(tx *Transaction) (common.Address, error) {
	if tx.IsEthTypedTransaction() {
		return common.Address{}, ErrTxTypeNotSupported
	}

	if !tx.Protected() {
		return HomesteadSigner{}.Sender(tx)
	}

	if !tx.IsLegacyTransaction() {
		logger.Warn("No need to execute Sender!", "tx", tx.String())
	}

	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	return RecoverTxSender(s.Hash(tx), tx.data.RawSignatureValues(), true, func(v *big.Int) *big.Int {
		V := new(big.Int).Sub(v, s.chainIdMul)
		return V.Sub(V, big8)
	})
}

func (s EIP155Signer) SenderPubkey(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	if tx.IsEthTypedTransaction() {
		return nil, ErrTxTypeNotSupported
	}

	if tx.IsLegacyTransaction() {
		logger.Warn("No need to execute SenderPubkey!", "tx", tx.String())
	}

	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, ErrInvalidChainId
	}
	return RecoverTxPubkeys(s.Hash(tx), tx.data.RawSignatureValues(), true, func(v *big.Int) *big.Int {
		V := new(big.Int).Sub(v, s.chainIdMul)
		return V.Sub(V, big8)
	})
}

func (s EIP155Signer) SenderFeePayer(tx *Transaction) ([]*ecdsa.PublicKey, error) {
	if tx.IsEthTypedTransaction() {
		return nil, ErrTxTypeNotSupported
	}

	if tx.IsLegacyTransaction() {
		logger.Warn("No need to execute SenderFeePayer!", "tx", tx.String())
	}

	if tx.ChainId().Cmp(s.chainId) != 0 {
		return nil, ErrInvalidChainId
	}

	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return nil, errNotFeeDelegationTransaction
	}

	hash, err := s.HashFeePayer(tx)
	if err != nil {
		return nil, err
	}

	return RecoverTxPubkeys(hash, tf.GetFeePayerRawSignatureValues(), true, func(v *big.Int) *big.Int {
		V := new(big.Int).Sub(v, s.chainIdMul)
		return V.Sub(V, big8)
	})
}

// SignatureValues returns a new transaction with the given signature. This signature
// needs to be in the [R || S || V] format where V is 0 or 1.
func (s EIP155Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	if tx.Type().IsEthTypedTransaction() {
		return nil, nil, nil, ErrTxTypeNotSupported
	}

	R, S, _ = decodeSignature(sig)
	V = big.NewInt(int64(sig[crypto.RecoveryIDOffset] + 35))
	V.Add(V, s.chainIdMul)

	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s EIP155Signer) Hash(tx *Transaction) common.Hash {
	return tx.data.SigHash(s.ChainID())
}

// HashFeePayer returns the hash with a fee payer's address to be signed by a fee payer.
// It does not uniquely identify the transaction.
func (s EIP155Signer) HashFeePayer(tx *Transaction) (common.Hash, error) {
	tf, ok := tx.data.(TxInternalDataFeePayer)
	if !ok {
		return common.Hash{}, errNotFeeDelegationTransaction
	}

	return tf.FeePayerSigHash(s.ChainID()), nil
}

func recoverPlainCommon(sighash common.Hash, R, S, Vb *big.Int, homestead bool) ([]byte, error) {
	if Vb.BitLen() > 8 {
		return []byte{}, ErrInvalidSig
	}
	V := byte(Vb.Uint64() - 27)
	if !crypto.ValidateSignatureValues(V, R, S, homestead) {
		return []byte{}, ErrInvalidSig
	}
	// encode the snature in uncompressed format
	r, s := R.Bytes(), S.Bytes()
	sig := make([]byte, crypto.SignatureLength)
	copy(sig[32-len(r):32], r)
	copy(sig[64-len(s):64], s)
	sig[crypto.RecoveryIDOffset] = V
	// recover the public key from the snature
	pub, err := crypto.Ecrecover(sighash[:], sig)
	if err != nil {
		return []byte{}, err
	}
	if len(pub) == 0 || pub[0] != 4 {
		return []byte{}, errors.New("invalid public key")
	}
	return pub, nil
}

func recoverPlain(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (common.Address, error) {
	pub, err := recoverPlainCommon(sighash, R, S, Vb, homestead)
	if err != nil {
		return common.Address{}, err
	}

	var addr common.Address
	copy(addr[:], crypto.Keccak256(pub[1:])[12:])
	return addr, nil
}

func recoverPlainPubkey(sighash common.Hash, R, S, Vb *big.Int, homestead bool) (*ecdsa.PublicKey, error) {
	pub, err := recoverPlainCommon(sighash, R, S, Vb, homestead)
	if err != nil {
		return nil, err
	}

	pubkey, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return nil, err
	}

	return pubkey, nil
}

// deriveChainId derives the chain id from the given v parameter
func deriveChainId(v *big.Int) *big.Int {
	if v.BitLen() <= 64 {
		v := v.Uint64()
		if v == 27 || v == 28 {
			return new(big.Int)
		}
		return new(big.Int).SetUint64((v - 35) / 2)
	}
	v = new(big.Int).Sub(v, big.NewInt(35))
	return v.Div(v, common.Big2)
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}
