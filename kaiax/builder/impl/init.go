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

package impl

import (
	"time"

	"github.com/kaiachain/kaia/api"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/builder"
	"github.com/kaiachain/kaia/log"
)

var (
	_      builder.BuilderModule = (*BuilderModule)(nil)
	logger                       = log.NewModuleLogger(log.KaiaxBuilder)
)

type InitOpts struct {
	Backend api.Backend
	Modules []builder.TxBundlingModule
}

type txAndTime struct {
	time time.Time
	tx   *types.Transaction
}

type BuilderModule struct {
	InitOpts

	txAndTimes map[common.Hash]txAndTime
}

func NewBuilderModule() *BuilderModule {
	return &BuilderModule{}
}

func (b *BuilderModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Backend == nil || opts.Modules == nil {
		return ErrInitUnexpectedNil
	}
	b.InitOpts = *opts
	b.txAndTimes = make(map[common.Hash]txAndTime)
	return nil
}

func (b *BuilderModule) Start() error {
	return nil
}

func (b *BuilderModule) Stop() {
}
