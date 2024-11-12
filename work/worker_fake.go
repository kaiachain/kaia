// Modifications Copyright 2024 The Kaia Authors
// Copyright 2020 The klaytn Authors
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

package work

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/kaiax"
)

type FakeWorker struct{}

// NewFakeWorker disables mining and block processing
// worker and istanbul engine will not be started.
func NewFakeWorker() *FakeWorker {
	logger.Warn("worker is disabled; no processing according to consensus logic")
	return &FakeWorker{}
}

func (*FakeWorker) Start()                                                   {}
func (*FakeWorker) Stop()                                                    {}
func (*FakeWorker) Register(Agent)                                           {}
func (*FakeWorker) Mining() bool                                             { return false }
func (*FakeWorker) HashRate() (tot int64)                                    { return 0 }
func (*FakeWorker) SetExtra([]byte) error                                    { return nil }
func (*FakeWorker) Pending() (*types.Block, *state.StateDB)                  { return nil, nil }
func (*FakeWorker) PendingBlock() *types.Block                               { return nil }
func (*FakeWorker) RegisterExecutionModule(modules ...kaiax.ExecutionModule) {}
