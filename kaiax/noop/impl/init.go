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

package impl

import (
	"github.com/kaiachain/kaia/kaiax/noop"
	"github.com/kaiachain/kaia/log"
)

var (
	_ (noop.NoopModule) = (*NoopModule)(nil)

	logger = log.NewModuleLogger(log.KaiaxMeta)
)

type NoopModule struct{}

func NewNoopModule() *NoopModule {
	return &NoopModule{}
}

func (m *NoopModule) Init() error {
	logger.Info("NoopModule Init")
	return nil
}

func (m *NoopModule) Start() error {
	logger.Info("NoopModule Start")
	return nil
}

func (m *NoopModule) Stop() {
	logger.Info("NoopModule Stop")
}
