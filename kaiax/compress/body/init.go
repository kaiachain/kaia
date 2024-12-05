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

type BodyCompressModule struct {
	compress.InitOpts
}

func NewBodyCompression() *BodyCompressModule {
	return &BodyCompressModule{}
}

func (bc *BodyCompressModule) GetChain() compress.BlockChain {
	return bc.Chain
}

func (bc *BodyCompressModule) GetDbm() database.DBManager {
	return bc.Dbm
}

func (bc *BodyCompressModule) Init(opts *compress.InitOpts) error {
	if opts == nil || opts.Chain == nil || opts.Dbm == nil {
		return errBCInitNil
	}
	bc.InitOpts = *opts
	return nil
}

func (bc *BodyCompressModule) Start() error {
	compress.Logger.Info("[Body Compression] Compression started")
	go bc.Compress()
	return nil
}

func (bc *BodyCompressModule) Stop() {
	compress.Logger.Info("[Body Compression] Compression Stopped")
}
