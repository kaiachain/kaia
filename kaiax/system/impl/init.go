// Copyright 2026 The Kaia Authors
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
	"errors"

	"github.com/kaiachain/kaia/consensus"
	kaiaxsystem "github.com/kaiachain/kaia/kaiax/system"
	"github.com/kaiachain/kaia/log"
)

var (
	_ kaiaxsystem.SystemModule = (*SystemModule)(nil)

	logger = log.NewModuleLogger(log.Blockchain)
)

type InitOpts struct {
	Chain consensus.ChainReaderWithSealer
}

type SystemModule struct {
	InitOpts
}

func NewSystemModule() *SystemModule {
	return &SystemModule{}
}

func (m *SystemModule) Init(opts *InitOpts) error {
	if opts == nil || opts.Chain == nil {
		return errors.New("unexpected nil init option")
	}
	m.InitOpts = *opts
	return nil
}

func (m *SystemModule) Start() error {
	return nil
}

func (m *SystemModule) Stop() {
}
