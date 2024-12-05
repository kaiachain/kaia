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

package compress

import (
	"github.com/kaiachain/kaia/kaiax/compress"
	"github.com/kaiachain/kaia/storage/database"
)

type HeaderCompressModule struct {
	compress.InitOpts
}

func NewHeaderCompression() *HeaderCompressModule {
	return &HeaderCompressModule{}
}

func (hc *HeaderCompressModule) GetChain() compress.BlockChain {
	return hc.Chain
}

func (hc *HeaderCompressModule) GetDbm() database.DBManager {
	return hc.Dbm
}

func (hc *HeaderCompressModule) Init(opts *compress.InitOpts) error {
	if opts == nil || opts.Chain == nil || opts.Dbm == nil {
		return errHCInitNil
	}
	hc.InitOpts = *opts
	return nil
}

func (hc *HeaderCompressModule) Start() error {
	compress.Logger.Info("[Header Compression] Compression started")
	go hc.Compress()
	return nil
}

func (hc *HeaderCompressModule) Stop() {
	compress.Logger.Info("[Header Compression] Compression Stopped")
}
