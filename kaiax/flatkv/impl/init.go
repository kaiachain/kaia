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
	"github.com/erigontech/erigon-lib/kv"
	flatkv "github.com/kaiachain/kaia/kaiax/flatkv"
)

var _ flatkv.FlatKVModule = &FlatKVModule{}

type InitOpts struct {
	chaindb kv.RwDB
}

type FlatKVModule struct {
	InitOpts
}

func NewFlatKVModule() *FlatKVModule {
	return &FlatKVModule{}
}

func (k *FlatKVModule) Init(opts *InitOpts) error {
	if opts == nil || opts.chaindb == nil {
		return ErrInitUnexpectedNil
	}
	k.InitOpts = *opts
	return nil
}

func (k *FlatKVModule) Start() error {
	return nil
}

func (k *FlatKVModule) Stop() {
}
