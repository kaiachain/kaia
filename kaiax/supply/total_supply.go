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

package supply

import (
	"math/big"

	"github.com/kaiachain/kaia/common/hexutil"
)

// AccReward is a subset of TotalSupply that comprises the minted and burnt amounts by the block reward mechanism.
type AccReward struct {
	TotalMinted *big.Int // Genesis + Minted[1..n]
	BurntFee    *big.Int // BurntFee[1..n]
}

func (ar *AccReward) Copy() *AccReward {
	return &AccReward{
		TotalMinted: new(big.Int).Set(ar.TotalMinted),
		BurntFee:    new(big.Int).Set(ar.BurntFee),
	}
}

func (ar *AccReward) ToTotalSupply(zeroBurn, deadBurn, kip103Burn, kip160Burn *big.Int) *TotalSupply {
	ts := &TotalSupply{
		TotalSupply: nil, // will be filled below

		TotalMinted: new(big.Int).Set(ar.TotalMinted),

		TotalBurnt: nil, // will be filled below
		BurntFee:   new(big.Int).Set(ar.BurntFee),
		ZeroBurn:   zeroBurn,
		DeadBurn:   deadBurn,
		Kip103Burn: kip103Burn,
		Kip160Burn: kip160Burn,
	}

	if ar.TotalMinted != nil && ar.BurntFee != nil && zeroBurn != nil && deadBurn != nil && kip103Burn != nil && kip160Burn != nil {
		totalBurnt := new(big.Int).Set(ar.BurntFee)
		totalBurnt.Add(totalBurnt, zeroBurn)
		totalBurnt.Add(totalBurnt, deadBurn)
		totalBurnt.Add(totalBurnt, kip103Burn)
		totalBurnt.Add(totalBurnt, kip160Burn)
		totalSupply := new(big.Int).Sub(ar.TotalMinted, totalBurnt)

		ts.TotalSupply = totalSupply
		ts.TotalBurnt = totalBurnt
	}
	return ts
}

type TotalSupply struct {
	TotalSupply *big.Int // TotalMinted - TotalBurnt

	TotalMinted *big.Int // AccReward.TotalMinted. It covers all minting amounts.

	// Tokens are burnt by various mechanisms.
	TotalBurnt *big.Int // Sum of all burnt amounts: BurntFee[1..n] + CanonicalBurn[n] + RebalanceBurn[n]
	BurntFee   *big.Int // AccReward.BurntFee
	ZeroBurn   *big.Int // CanonicalBurn[n] at 0x0
	DeadBurn   *big.Int // CanonicalBurn[n] at 0xdead
	Kip103Burn *big.Int // RebalanceBurn[n] by KIP-103
	Kip160Burn *big.Int // RebalanceBurn[n] by KIP-160
}

type TotalSupplyResponse struct {
	// Block number in which the total supply was calculated.
	Number *hexutil.Big `json:"number"`
	// Errors that occurred while fetching the components, thus failed to deliver the total supply. Only applies when showPartial is true.
	Error *string `json:"error,omitempty"`

	// The total supply of the native token. i.e. TotalMinted - TotalBurnt.
	TotalSupply *hexutil.Big `json:"totalSupply"`

	// Accumulated block rewards minted up to number plus the genesis total supply.
	TotalMinted *hexutil.Big `json:"totalMinted"`

	// Accumulated burn amount. Sum of all burns. It is null if some components cannot be determined.
	TotalBurnt *hexutil.Big `json:"totalBurnt"`
	// Accumulated transaction fees burnt
	BurntFee *hexutil.Big `json:"burntFee"`
	// The balance of the address 0x0, which is a canonical burn address. It is null if the state at the requested block is not available.
	ZeroBurn *hexutil.Big `json:"zeroBurn"`
	// The balance of the address 0xdead, which is another canonical burn address. It is null if the state at the requested block is not available.
	DeadBurn *hexutil.Big `json:"deadBurn"`
	// The net amounts burnt by KIP-103 TreasuryRebalance, if executed. It is null if KIP-103 is configured and the hardfork block passed but memo field is not set in the contract. It is 0 if KIP-103 is not configured or the hardfork block is larger than the requested block.
	Kip103Burn *hexutil.Big `json:"kip103Burn"`
	// The net amounts burnt by KIP-160 TreasuryRebalanceV2, if executed. It is null if KIP-160 is configured and the hardfork block passed but memo field is not set in the contract. It is 0 if KIP-160 is not configured or the hardfork block is larger than the requested block.
	Kip160Burn *hexutil.Big `json:"kip160Burn"`
}

func (ts *TotalSupply) ToResponse(num uint64, err error) *TotalSupplyResponse {
	var pErrStr *string
	if err != nil {
		errStr := err.Error()
		pErrStr = &errStr
	}
	return &TotalSupplyResponse{
		Number:      (*hexutil.Big)(new(big.Int).SetUint64(num)),
		Error:       pErrStr,
		TotalSupply: (*hexutil.Big)(ts.TotalSupply),
		TotalMinted: (*hexutil.Big)(ts.TotalMinted),
		TotalBurnt:  (*hexutil.Big)(ts.TotalBurnt),
		BurntFee:    (*hexutil.Big)(ts.BurntFee),
		ZeroBurn:    (*hexutil.Big)(ts.ZeroBurn),
		DeadBurn:    (*hexutil.Big)(ts.DeadBurn),
		Kip103Burn:  (*hexutil.Big)(ts.Kip103Burn),
		Kip160Burn:  (*hexutil.Big)(ts.Kip160Burn),
	}
}
