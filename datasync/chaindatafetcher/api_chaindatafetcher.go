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

package chaindatafetcher

import (
	"errors"
	"sync/atomic"

	"github.com/kaiachain/kaia/datasync/chaindatafetcher/types"
)

type ChainDataFetcherAPI struct {
	f *ChainDataFetcher
}

func NewChainDataFetcherAPI(f *ChainDataFetcher) *ChainDataFetcherAPI {
	return &ChainDataFetcherAPI{f: f}
}

func (api *ChainDataFetcherAPI) StartFetching() error {
	return api.f.startFetching()
}

func (api *ChainDataFetcherAPI) StopFetching() error {
	return api.f.stopFetching()
}

func (api *ChainDataFetcherAPI) StartRangeFetching(start, end uint64, reqType interface{}) error {
	var t types.RequestType
	switch reqType {
	case "all":
		t = types.RequestTypeGroupAll
	case "block":
		t = types.RequestTypeBlockGroup
	case "trace":
		t = types.RequestTypeTraceGroup
	default:
		ut, ok := reqType.(float64)
		if !ok {
			return errors.New("the request type should be 'all', 'block', 'trace', or uint type")
		}
		t = types.RequestType(ut)
	}

	if !t.IsValid() {
		return errors.New("the request type is not valid")
	}

	return api.f.startRangeFetching(start, end, t)
}

func (api *ChainDataFetcherAPI) StopRangeFetching() error {
	return api.f.stopRangeFetching()
}

func (api *ChainDataFetcherAPI) Status() string {
	return api.f.status()
}

func (api *ChainDataFetcherAPI) ReadCheckpoint() (int64, error) {
	return api.f.checkpointDB.ReadCheckpoint()
}

func (api *ChainDataFetcherAPI) WriteCheckpoint(checkpoint int64) error {
	isRunning := atomic.LoadUint32(&api.f.fetchingStarted)
	if isRunning == running {
		return errors.New("call stopFetching before writing checkpoint manually")
	}

	api.f.checkpoint = checkpoint
	return api.f.checkpointDB.WriteCheckpoint(checkpoint)
}

// GetConfig returns the configuration setting of the launched chaindata fetcher.
func (api *ChainDataFetcherAPI) GetConfig() *ChainDataFetcherConfig {
	return api.f.config
}
