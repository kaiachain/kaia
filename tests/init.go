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
// This file is derived from tests/init.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package tests

import (
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/params"
)

// Forks table defines supported forks and their chain config.
// TODO-Kaia-RemoveLater Remove fork configs that are not meaningful for Kaia
var Forks = map[string]*params.ChainConfig{
	"Frontier": {
		ChainID: big.NewInt(1),
	},
	"Homestead": {
		ChainID: big.NewInt(1),
	},
	"Byzantium": {
		ChainID: big.NewInt(1),
	},
	"Constantinople": {
		ChainID: big.NewInt(1),
	},
	"ConstantinopleFix": {
		ChainID: big.NewInt(1),
	},
	"Istanbul": {
		ChainID:                 big.NewInt(1),
		IstanbulCompatibleBlock: new(big.Int),
	},
	"Berlin": {
		ChainID: big.NewInt(1),
	},
	"London": {
		ChainID:                 big.NewInt(1),
		IstanbulCompatibleBlock: new(big.Int),
		LondonCompatibleBlock:   new(big.Int),
	},
	"EthTxType": {
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  new(big.Int),
		LondonCompatibleBlock:    new(big.Int),
		EthTxTypeCompatibleBlock: new(big.Int),
	},
	"Magma": {
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  new(big.Int),
		LondonCompatibleBlock:    new(big.Int),
		EthTxTypeCompatibleBlock: new(big.Int),
		MagmaCompatibleBlock:     new(big.Int),
	},
	"Merge": {
		ChainID: big.NewInt(1),
	},
	"Shanghai": {
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  new(big.Int),
		LondonCompatibleBlock:    new(big.Int),
		EthTxTypeCompatibleBlock: new(big.Int),
		MagmaCompatibleBlock:     new(big.Int),
		KoreCompatibleBlock:      new(big.Int),
		ShanghaiCompatibleBlock:  new(big.Int),
	},
	"Cancun": {
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  new(big.Int),
		LondonCompatibleBlock:    new(big.Int),
		EthTxTypeCompatibleBlock: new(big.Int),
		MagmaCompatibleBlock:     new(big.Int),
		KoreCompatibleBlock:      new(big.Int),
		Kip103CompatibleBlock:    new(big.Int),
		ShanghaiCompatibleBlock:  new(big.Int),
		CancunCompatibleBlock:    new(big.Int),
		KaiaCompatibleBlock:      new(big.Int),
		Kip160CompatibleBlock:    new(big.Int),
	},
	"Prague": {
		ChainID:                  big.NewInt(1),
		IstanbulCompatibleBlock:  new(big.Int),
		LondonCompatibleBlock:    new(big.Int),
		EthTxTypeCompatibleBlock: new(big.Int),
		MagmaCompatibleBlock:     new(big.Int),
		KoreCompatibleBlock:      new(big.Int),
		Kip103CompatibleBlock:    new(big.Int),
		ShanghaiCompatibleBlock:  new(big.Int),
		CancunCompatibleBlock:    new(big.Int),
		KaiaCompatibleBlock:      new(big.Int),
		Kip160CompatibleBlock:    new(big.Int),
		PragueCompatibleBlock:    new(big.Int),
	},
}

// UnsupportedForkError is returned when a test requests a fork that isn't implemented.
type UnsupportedForkError struct {
	Name string
}

func (e UnsupportedForkError) Error() string {
	return fmt.Sprintf("unsupported fork %q", e.Name)
}
