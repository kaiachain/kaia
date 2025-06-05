// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/ethclient_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package client

import kaia "github.com/kaiachain/kaia/v2"

// Verify that Client implements the Kaia interfaces.
var (
	// _ = kaia.Subscription(&Client{})
	_ = kaia.ChainReader(&Client{})
	_ = kaia.TransactionReader(&Client{})
	_ = kaia.ChainStateReader(&Client{})
	_ = kaia.ChainSyncReader(&Client{})
	_ = kaia.ContractCaller(&Client{})
	_ = kaia.LogFilterer(&Client{})
	_ = kaia.TransactionSender(&Client{})
	_ = kaia.GasPricer(&Client{})
	_ = kaia.PendingStateReader(&Client{})
	_ = kaia.PendingContractCaller(&Client{})
	_ = kaia.GasEstimator(&Client{})
	// _ = kaia.PendingStateEventer(&Client{})
)
