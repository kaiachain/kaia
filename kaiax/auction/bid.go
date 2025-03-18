// Copyright 2025 The Kaia Authors
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

package auction

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"sync/atomic"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
)

const (
	auctionType      = "AuctionTx(bytes32 targetTxHash,uint256 blockNumber,address sender,address to,uint256 nonce,uint256 bid,uint256 callGasLimit,bytes data)"
	EIP712DomainType = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
	auctionName      = "KAIA_AUCTION"
	auctionVersion   = "0.0.1"
)

var (
	auctionTypeHash    = crypto.Keccak256Hash([]byte(auctionType))
	eip712TypeHash     = crypto.Keccak256Hash([]byte(EIP712DomainType))
	auctionNameHash    = crypto.Keccak256Hash([]byte(auctionName))
	auctionVersionHash = crypto.Keccak256Hash([]byte(auctionVersion))
)

// EIP712Encoder interface defines the methods required for EIP-712 encoding
type EIP712Encoder interface {
	EncodeType() []byte
	EncodeData() []byte
}

type EIP712Domain struct {
	EIP712DomainTypeHash common.Hash
	NameHash             common.Hash
	VersionHash          common.Hash
	ChainId              *big.Int
	VerifyingContract    common.Address
}

func (d EIP712Domain) EncodeType() []byte {
	return d.EIP712DomainTypeHash.Bytes()
}

func (d EIP712Domain) EncodeData() []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, d.NameHash.Bytes()...)
	encoded = append(encoded, d.VersionHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(d.ChainId.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(d.VerifyingContract.Bytes(), 32)...)
	return encoded
}

type bidData struct {
	TargetTxHash  common.Hash    `json:"targetTxHash"`
	BlockNumber   uint64         `json:"blockNumber"`
	Sender        common.Address `json:"sender"`
	To            common.Address `json:"to"`
	Nonce         uint64         `json:"nonce"`
	Bid           *big.Int       `json:"bid"`
	CallGasLimit  uint64         `json:"callGasLimit"`
	Data          []byte         `json:"data"`
	SearcherSig   []byte         `json:"searcherSig"`
	AuctioneerSig []byte         `json:"auctioneerSig"`
}

type Bid struct {
	bidData
	hash atomic.Value
}

func (at *Bid) Hash() common.Hash {
	if at.hash.Load() == nil {
		at.hash.Store(crypto.Keccak256Hash(at.EncodeData()))
	}
	return at.hash.Load().(common.Hash)
}

func (at *Bid) EncodeType() []byte {
	return auctionTypeHash.Bytes()
}

func (at *Bid) EncodeData() []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, at.TargetTxHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(at.BlockNumber), 32)...)
	encoded = append(encoded, common.LeftPadBytes(at.Sender.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(at.To.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(at.Nonce), 32)...)
	encoded = append(encoded, common.LeftPadBytes(at.Bid.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(at.CallGasLimit), 32)...)
	encoded = append(encoded, crypto.Keccak256Hash(at.Data).Bytes()...)
	return encoded
}

// encodeEIP712 encodes any EIP712Encoder according to EIP-712 specification
func encodeEIP712(encoder EIP712Encoder) []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, encoder.EncodeType()...)
	encoded = append(encoded, encoder.EncodeData()...)
	return crypto.Keccak256Hash(encoded).Bytes()
}

func (at *Bid) GetHashTypedData(chainId *big.Int, verifyingContract common.Address) []byte {
	if chainId == nil {
		return nil
	}

	domain := EIP712Domain{
		EIP712DomainTypeHash: eip712TypeHash,
		NameHash:             auctionNameHash,
		VersionHash:          auctionVersionHash,
		ChainId:              chainId,
		VerifyingContract:    verifyingContract,
	}

	domainSeparator := encodeEIP712(domain)
	structHash := encodeEIP712(at)

	return crypto.Keccak256([]byte{0x19, 0x01}, domainSeparator, structHash)
}

func (at *Bid) GetEthSignedMessageHash() []byte {
	data := at.SearcherSig
	return crypto.Keccak256(fmt.Appendf(nil, "\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

func (at *Bid) ValidateSearcherSig(chainId *big.Int, verifyingContract common.Address) error {
	if chainId == nil {
		return errors.New("chainId is nil")
	}

	if verifyingContract == (common.Address{}) {
		return errors.New("verifyingContract is empty")
	}

	digest := at.GetHashTypedData(chainId, verifyingContract)

	// Manually convert V from 27/28 to 0/1
	sig := slices.Clone(at.SearcherSig)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27
	}

	pub, err := crypto.Ecrecover(digest, sig)
	if err != nil {
		return fmt.Errorf("failed to recover searcher sig: %v", err)
	}

	recoveredSender := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	if recoveredSender != at.Sender {
		return errors.New("invalid searcher sig")
	}

	return nil
}

func (at *Bid) ValidateAuctioneerSig(auctioneer common.Address) error {
	digest := at.GetEthSignedMessageHash()

	// Manually convert V from 27/28 to 0/1
	sig := slices.Clone(at.AuctioneerSig)
	if sig[crypto.RecoveryIDOffset] == 27 || sig[crypto.RecoveryIDOffset] == 28 {
		sig[crypto.RecoveryIDOffset] -= 27
	}

	pub, err := crypto.Ecrecover(digest, sig)
	if err != nil {
		return fmt.Errorf("failed to recover auctioneer sig: %v", err)
	}

	recoveredAuctioneer := common.BytesToAddress(crypto.Keccak256(pub[1:])[12:])
	if recoveredAuctioneer != auctioneer {
		return fmt.Errorf("invalid auctioneer sig: expected %v, calculated %v", auctioneer, recoveredAuctioneer)
	}

	return nil
}
