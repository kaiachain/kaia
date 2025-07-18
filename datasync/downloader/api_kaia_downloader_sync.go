// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/downloader/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package downloader

// KaiaDownloaderSyncAPI provides an API which gives syncing staking information.
type KaiaDownloaderSyncAPI struct {
	d downloader
}

func NewKaiaDownloaderSyncAPI(d downloader) *KaiaDownloaderSyncAPI {
	api := &KaiaDownloaderSyncAPI{
		d: d,
	}
	return api
}

func (api *KaiaDownloaderSyncAPI) SyncStakingInfo(id string, from, to uint64) error {
	return api.d.SyncStakingInfo(id, from, to)
}

func (api *KaiaDownloaderSyncAPI) SyncStakingInfoStatus() *SyncingStatus {
	return api.d.SyncStakingInfoStatus()
}
