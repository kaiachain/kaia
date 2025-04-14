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
	lru "github.com/hashicorp/golang-lru"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/kaiax/randao"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"golang.org/x/sync/singleflight"
)

var (
	_ randao.RandaoModule = &RandaoModule{}

	logger = log.NewModuleLogger(log.KaiaxRandao)
)

type ProtocolManagerDownloader interface {
	Synchronising() bool
}

type InitOpts struct {
	ChainConfig *params.ChainConfig
	Chain       backends.BlockChainForCaller
	Downloader  ProtocolManagerDownloader
}

type RandaoModule struct {
	InitOpts

	blsPubkeyCache   *lru.ARCCache
	storageRootCache *lru.ARCCache
	sfGroup          singleflight.Group
}

func NewRandaoModule() *RandaoModule {
	blsPubkeyCache, _ := lru.NewARC(128)
	storageRootCache, _ := lru.NewARC(1000)
	return &RandaoModule{
		blsPubkeyCache:   blsPubkeyCache,
		storageRootCache: storageRootCache,
		sfGroup:          singleflight.Group{},
	}
}

func (r *RandaoModule) Init(opts *InitOpts) error {
	if opts == nil || opts.ChainConfig == nil || opts.Chain == nil || opts.Downloader == nil {
		return randao.ErrInitUnexpectedNil
	}
	r.InitOpts = *opts

	return nil
}

func (r *RandaoModule) Start() error {
	return nil
}

func (r *RandaoModule) Stop() {
}
