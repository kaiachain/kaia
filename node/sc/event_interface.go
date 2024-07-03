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
package sc

import (
	"math/big"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	bridgecontract "github.com/kaiachain/kaia/contracts/contracts/service_chain/bridge"
)

type IRequestValueTransferEvent interface {
	Nonce() uint64
	GetTokenType() uint8
	GetFrom() common.Address
	GetTo() common.Address
	GetTokenAddress() common.Address
	GetValueOrTokenId() *big.Int
	GetRequestNonce() uint64
	GetFee() *big.Int
	GetExtraData() []byte

	GetRaw() types.Log
}

// ////////////////// type RequestValueTransferEvent struct ////////////////////
// RequestValueTransferEvent from Bridge contract
type RequestValueTransferEvent struct {
	*bridgecontract.BridgeRequestValueTransfer
}

func (rEv RequestValueTransferEvent) Nonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEvent) GetTokenType() uint8 {
	return rEv.TokenType
}

func (rEv RequestValueTransferEvent) GetFrom() common.Address {
	return rEv.From
}

func (rEv RequestValueTransferEvent) GetTo() common.Address {
	return rEv.To
}

func (rEv RequestValueTransferEvent) GetTokenAddress() common.Address {
	return rEv.TokenAddress
}

func (rEv RequestValueTransferEvent) GetValueOrTokenId() *big.Int {
	return rEv.ValueOrTokenId
}

func (rEv RequestValueTransferEvent) GetRequestNonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEvent) GetFee() *big.Int {
	return rEv.Fee
}

func (rEv RequestValueTransferEvent) GetExtraData() []byte {
	return rEv.ExtraData
}

func (rEv RequestValueTransferEvent) GetRaw() types.Log {
	return rEv.Raw
}

// ////////////////// type RequestValueTransferEncodedEvent struct ////////////////////
type RequestValueTransferEncodedEvent struct {
	*bridgecontract.BridgeRequestValueTransferEncoded
}

func (rEv RequestValueTransferEncodedEvent) Nonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEncodedEvent) GetTokenType() uint8 {
	return rEv.TokenType
}

func (rEv RequestValueTransferEncodedEvent) GetFrom() common.Address {
	return rEv.From
}

func (rEv RequestValueTransferEncodedEvent) GetTo() common.Address {
	return rEv.To
}

func (rEv RequestValueTransferEncodedEvent) GetTokenAddress() common.Address {
	return rEv.TokenAddress
}

func (rEv RequestValueTransferEncodedEvent) GetValueOrTokenId() *big.Int {
	return rEv.ValueOrTokenId
}

func (rEv RequestValueTransferEncodedEvent) GetRequestNonce() uint64 {
	return rEv.RequestNonce
}

func (rEv RequestValueTransferEncodedEvent) GetFee() *big.Int {
	return rEv.Fee
}

func (rEv RequestValueTransferEncodedEvent) GetExtraData() []byte {
	return rEv.ExtraData
}

func (rEv RequestValueTransferEncodedEvent) GetRaw() types.Log {
	return rEv.Raw
}
