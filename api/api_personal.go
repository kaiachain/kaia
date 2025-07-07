// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package api

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"time"

	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/common/math"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/rlp"
)

// PersonalAPI provides an API to access accounts managed by this node.
// It offers methods to create, (un)lock en list accounts. Some methods accept
// passwords and are therefore considered private by default.
type PersonalAPI struct {
	am        accounts.AccountManager
	nonceLock *AddrLocker
	b         Backend
}

// NewPersonalAPI create a new PersonalAPI.
func NewPersonalAPI(b Backend, nonceLock *AddrLocker) *PersonalAPI {
	return &PersonalAPI{
		am:        b.AccountManager(),
		nonceLock: nonceLock,
		b:         b,
	}
}

// ListAccounts will return a list of addresses for accounts this node manages.
func (s *PersonalAPI) ListAccounts() []common.Address {
	addresses := make([]common.Address, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		for _, account := range wallet.Accounts() {
			addresses = append(addresses, account.Address)
		}
	}
	return addresses
}

// rawWallet is a JSON representation of an accounts.Wallet interface, with its
// data contents extracted into plain fields.
type rawWallet struct {
	URL      string             `json:"url"`
	Status   string             `json:"status"`
	Failure  string             `json:"failure,omitempty"`
	Accounts []accounts.Account `json:"accounts,omitempty"`
}

// ListWallets will return a list of wallets this node manages.
func (s *PersonalAPI) ListWallets() []rawWallet {
	wallets := make([]rawWallet, 0) // return [] instead of nil if empty
	for _, wallet := range s.am.Wallets() {
		status, failure := wallet.Status()

		raw := rawWallet{
			URL:      wallet.URL().String(),
			Status:   status,
			Accounts: wallet.Accounts(),
		}
		if failure != nil {
			raw.Failure = failure.Error()
		}
		wallets = append(wallets, raw)
	}
	return wallets
}

// OpenWallet initiates a hardware wallet opening procedure, establishing a USB
// connection and attempting to authenticate via the provided passphrase. Note,
// the method may return an extra challenge requiring a second open (e.g. the
// Trezor PIN matrix challenge).
func (s *PersonalAPI) OpenWallet(url string, passphrase *string) error {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return err
	}
	pass := ""
	if passphrase != nil {
		pass = *passphrase
	}
	return wallet.Open(pass)
}

// DeriveAccount requests a HD wallet to derive a new account, optionally pinning
// it for later reuse.
func (s *PersonalAPI) DeriveAccount(url string, path string, pin *bool) (accounts.Account, error) {
	wallet, err := s.am.Wallet(url)
	if err != nil {
		return accounts.Account{}, err
	}
	derivPath, err := accounts.ParseDerivationPath(path)
	if err != nil {
		return accounts.Account{}, err
	}
	if pin == nil {
		pin = new(bool)
	}
	return wallet.Derive(derivPath, *pin)
}

// NewAccount will create a new account and returns the address for the new account.
func (s *PersonalAPI) NewAccount(password string) (common.Address, error) {
	acc, err := fetchKeystore(s.am).NewAccount(password)
	if err == nil {
		return acc.Address, nil
	}
	return common.Address{}, err
}

// fetchKeystore retrives the encrypted keystore from the account manager.
func fetchKeystore(am accounts.AccountManager) *keystore.KeyStore {
	return am.Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)
}

func parseKaiaWalletKey(k string) (string, string, *common.Address, error) {
	// if key length is not 110, just return.
	if len(k) != 110 {
		return k, "", nil, nil
	}

	walletKeyType := k[66:68]
	if walletKeyType != "00" {
		return "", "", nil, errors.New("Kaia wallet key type must be 00.")
	}
	a := common.HexToAddress(k[70:110])

	return k[0:64], walletKeyType, &a, nil
}

// ReplaceRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PersonalAPI) ReplaceRawKey(privkey string, passphrase string, newPassphrase string) (common.Address, error) {
	privkey, _, address, err := parseKaiaWalletKey(privkey)
	if err != nil {
		return common.Address{}, err
	}
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := fetchKeystore(s.am).ReplaceECDSAWithAddress(key, passphrase, newPassphrase, address)
	return acc.Address, err
}

// ImportRawKey stores the given hex encoded ECDSA key into the key directory,
// encrypting it with the passphrase.
func (s *PersonalAPI) ImportRawKey(privkey string, password string) (common.Address, error) {
	privkey, _, address, err := parseKaiaWalletKey(privkey)
	if err != nil {
		return common.Address{}, err
	}
	key, err := crypto.HexToECDSA(privkey)
	if err != nil {
		return common.Address{}, err
	}
	acc, err := fetchKeystore(s.am).ImportECDSAWithAddress(key, password, address)
	return acc.Address, err
}

// UnlockAccount will unlock the account associated with the given address with
// the given password for duration seconds. If duration is nil it will use a
// default of 300 seconds. It returns an indication if the account was unlocked.
func (s *PersonalAPI) UnlockAccount(addr common.Address, password string, duration *uint64) (bool, error) {
	const max = uint64(time.Duration(math.MaxInt64) / time.Second)
	var d time.Duration
	if duration == nil {
		d = 300 * time.Second
	} else if *duration > max {
		return false, errors.New("unlock duration too large")
	} else {
		d = time.Duration(*duration) * time.Second
	}
	err := fetchKeystore(s.am).TimedUnlock(accounts.Account{Address: addr}, password, d)
	return err == nil, err
}

// LockAccount will lock the account associated with the given address when it's unlocked.
func (s *PersonalAPI) LockAccount(addr common.Address) bool {
	return fetchKeystore(s.am).Lock(addr) == nil
}

// signTransaction sets defaults and signs the given transaction.
// NOTE: the caller needs to ensure that the nonceLock is held, if applicable,
// and release it after the transaction has been submitted to the tx pool.
func (s *PersonalAPI) signTransaction(ctx context.Context, args SendTxArgs, passwd string) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: args.From}
	wallet, err := s.am.Find(account)
	if err != nil {
		return nil, err
	}
	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	// Assemble the transaction and sign with the wallet
	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}

	return wallet.SignTxWithPassphrase(account, passwd, tx, s.b.ChainConfig().ChainID)
}

// SendTransaction will create a transaction from the given arguments and try to
// sign it with the key associated with args.From. If the given password isn't
// able to decrypt the key it fails.
func (s *PersonalAPI) SendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {
	if args.AccountNonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}
	signedTx, err := s.SignTransaction(ctx, args, passwd)
	if err != nil {
		return common.Hash{}, err
	}
	return submitTransaction(ctx, s.b, signedTx.Tx)
}

// SendTransactionAsFeePayer will create a transaction from the given arguments and
// try to sign it as a fee payer with the key associated with args.From. If the
// given password isn't able to decrypt the key it fails.
func (s *PersonalAPI) SendTransactionAsFeePayer(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {
	// Don't allow dynamic assign of values from the setDefaults function since the sender already signed on specific values.
	if args.TypeInt == nil {
		return common.Hash{}, errTxArgNilTxType
	}
	if args.AccountNonce == nil {
		return common.Hash{}, errTxArgNilNonce
	}
	if args.GasLimit == nil {
		return common.Hash{}, errTxArgNilGas
	}
	if args.Price == nil {
		return common.Hash{}, errTxArgNilGasPrice
	}

	if args.TxSignatures == nil {
		return common.Hash{}, errTxArgNilSenderSig
	}

	feePayerSignedTx, err := s.SignTransactionAsFeePayer(ctx, args, passwd)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.b, feePayerSignedTx.Tx)
}

func (s *PersonalAPI) signNewTransaction(ctx context.Context, args NewTxArgs, passwd string) (*types.Transaction, error) {
	account := accounts.Account{Address: args.from()}
	wallet, err := s.am.Find(account)
	if err != nil {
		return nil, err
	}

	// Set some sanity defaults and terminate on failure
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}

	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}

	signed, err := wallet.SignTxWithPassphrase(account, passwd, tx, s.b.ChainConfig().ChainID)
	if err != nil {
		return nil, err
	}

	return signed, nil
}

// SendAccountUpdate will create a TxTypeAccountUpdate transaction from the given arguments and
// try to sign it with the key associated with args.From. If the given password isn't able to
// decrypt the key it fails.
func (s *PersonalAPI) SendAccountUpdate(ctx context.Context, args AccountUpdateTxArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	signed, err := s.signNewTransaction(ctx, &args, passwd)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.b, signed)
}

// SendValueTransfer will create a TxTypeValueTransfer transaction from the given arguments and
// try to sign it with the key associated with args.From. If the given password isn't able to
// decrypt the key it fails.
func (s *PersonalAPI) SendValueTransfer(ctx context.Context, args ValueTransferTxArgs, passwd string) (common.Hash, error) {
	if args.Nonce == nil {
		// Hold the addresse's mutex around signing to prevent concurrent assignment of
		// the same nonce to multiple accounts.
		s.nonceLock.LockAddr(args.From)
		defer s.nonceLock.UnlockAddr(args.From)
	}

	signed, err := s.signNewTransaction(ctx, &args, passwd)
	if err != nil {
		return common.Hash{}, err
	}

	return submitTransaction(ctx, s.b, signed)
}

// SignTransaction will create a transaction from the given arguments and
// try to sign it with the key associated with args.From. If the given password isn't able to
// decrypt the key, it fails. The transaction is returned in RLP-form, not broadcast to other nodes
func (s *PersonalAPI) SignTransaction(ctx context.Context, args SendTxArgs, passwd string) (*SignTransactionResult, error) {
	if args.TypeInt != nil && args.TypeInt.IsEthTypedTransaction() {
		if args.Price == nil && (args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil) {
			return nil, errors.New("missing gasPrice or maxFeePerGas/maxPriorityFeePerGas")
		}
	}

	// No need to obtain the noncelock mutex, since we won't be sending this
	// tx into the transaction pool, but right back to the user
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}
	signedTx, err := s.sign(args.From, passwd, tx)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, signedTx}, nil
}

// SignTransactionAsFeePayer will create a transaction from the given arguments and
// try to sign it as a fee payer with the key associated with args.From. If the given
// password isn't able to decrypt the key, it fails. The transaction is returned in RLP-form,
// not broadcast to other nodes
func (s *PersonalAPI) SignTransactionAsFeePayer(ctx context.Context, args SendTxArgs, passwd string) (*SignTransactionResult, error) {
	// Allows setting a default nonce value of the sender just for the case the fee payer tries to sign a tx earlier than the sender.
	if err := args.setDefaults(ctx, s.b); err != nil {
		return nil, err
	}
	tx, err := args.toTransaction()
	if err != nil {
		return nil, err
	}
	// Don't return errors for nil signature allowing the fee payer to sign a tx earlier than the sender.
	if args.TxSignatures != nil {
		tx.SetSignature(args.TxSignatures.ToTxSignatures())
	}
	feePayer, err := tx.FeePayer()
	if err != nil {
		return nil, errTxArgInvalidFeePayer
	}
	feePayerSignedTx, err := s.signAsFeePayer(feePayer, passwd, tx)
	if err != nil {
		return nil, err
	}
	data, err := rlp.EncodeToBytes(feePayerSignedTx)
	if err != nil {
		return nil, err
	}
	return &SignTransactionResult{data, feePayerSignedTx}, nil
}

// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//
//	keccak256("\x19Klaytn Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Klaytn Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func ethSignHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// sign is a helper function that signs a transaction with the private key of the given address.
// If the given password isn't able to decrypt the key, it fails.
func (s *PersonalAPI) sign(addr common.Address, passwd string, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTxWithPassphrase(account, passwd, tx, s.b.ChainConfig().ChainID)
}

// signAsFeePayer is a helper function that signs a transaction with the private key of the given address.
// If the given password isn't able to decrypt the key, it fails.
func (s *PersonalAPI) signAsFeePayer(addr common.Address, passwd string, tx *types.Transaction) (*types.Transaction, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Request the wallet to sign the transaction
	return wallet.SignTxAsFeePayerWithPassphrase(account, passwd, tx, s.b.ChainConfig().ChainID)
}

// Sign calculates a Kaia ECDSA signature for:
// keccack256("\x19Ethereum Signed Message:\n" + len(message) + message))
//
// Note, the produced signature conforms to the secp256k1 curve R, S and V values,
// where the V value will be 27 or 28 for legacy reasons.
//
// The key used to calculate the signature is decrypted with the given password.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_sign
func (s *PersonalAPI) Sign(ctx context.Context, data hexutil.Bytes, addr common.Address, passwd string) (hexutil.Bytes, error) {
	// Look up the wallet containing the requested signer
	account := accounts.Account{Address: addr}

	wallet, err := s.b.AccountManager().Find(account)
	if err != nil {
		return nil, err
	}
	// Assemble sign the data with the wallet
	signature, err := wallet.SignHashWithPassphrase(account, passwd, ethSignHash(data))
	if err != nil {
		return nil, err
	}
	signature[crypto.RecoveryIDOffset] += 27 // Transform V from 0/1 to 27/28 according to the yellow paper
	return signature, nil
}

// EcRecover returns the address for the account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be 27 or 28 for legacy reasons.
func (s *PersonalAPI) EcRecover(ctx context.Context, data, sig hexutil.Bytes) (common.Address, error) {
	if len(sig) != crypto.SignatureLength {
		return common.Address{}, errors.New("signature must be 65 bytes long")
	}
	if sig[crypto.RecoveryIDOffset] != 27 && sig[crypto.RecoveryIDOffset] != 28 {
		return common.Address{}, errors.New("invalid signature (V is not 27 or 28)")
	}

	// Transform yellow paper V from 27/28 to 0/1
	sig[crypto.RecoveryIDOffset] -= 27

	pubkey, err := ethEcRecover(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pubkey), nil
}

func klayEcRecover(data, sig hexutil.Bytes) (*ecdsa.PublicKey, error) {
	return crypto.SigToPub(signHash(data), sig)
}

func ethEcRecover(data, sig hexutil.Bytes) (*ecdsa.PublicKey, error) {
	return crypto.SigToPub(ethSignHash(data), sig)
}

// SignAndSendTransaction was renamed to SendTransaction. This method is deprecated
// and will be removed in the future. It primary goal is to give clients time to update.
func (s *PersonalAPI) SignAndSendTransaction(ctx context.Context, args SendTxArgs, passwd string) (common.Hash, error) {
	return s.SendTransaction(ctx, args, passwd)
}
