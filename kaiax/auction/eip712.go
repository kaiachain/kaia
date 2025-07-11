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
	"math/big"

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
	return eip712TypeHash.Bytes()
}

func (d EIP712Domain) EncodeData() []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, d.NameHash.Bytes()...)
	encoded = append(encoded, d.VersionHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(d.ChainId.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(d.VerifyingContract.Bytes(), 32)...)
	return encoded
}

func (b *Bid) EncodeType() []byte {
	return auctionTypeHash.Bytes()
}

func (b *Bid) EncodeData() []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, b.TargetTxHash.Bytes()...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(b.BlockNumber), 32)...)
	encoded = append(encoded, common.LeftPadBytes(b.Sender.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(b.To.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(b.Nonce), 32)...)
	encoded = append(encoded, common.LeftPadBytes(b.Bid.Bytes(), 32)...)
	encoded = append(encoded, common.LeftPadBytes(common.Int64ToByteBigEndian(b.CallGasLimit), 32)...)
	encoded = append(encoded, crypto.Keccak256Hash(b.Data).Bytes()...)
	return encoded
}

// EncodeEIP712 encodes any EIP712Encoder according to EIP-712 specification
func EncodeEIP712(encoder EIP712Encoder) []byte {
	encoded := make([]byte, 0)
	encoded = append(encoded, encoder.EncodeType()...)
	encoded = append(encoded, encoder.EncodeData()...)
	return crypto.Keccak256Hash(encoded).Bytes()
}

func (b *Bid) GetHashTypedData(chainId *big.Int, verifyingContract common.Address) []byte {
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

	domainSeparator := EncodeEIP712(domain)
	structHash := EncodeEIP712(b)

	return crypto.Keccak256([]byte{0x19, 0x01}, domainSeparator, structHash)
}
