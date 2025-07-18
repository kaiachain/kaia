// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of go-ethereum.
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
// This file is derived from eth/api.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package cn

import (
	"github.com/kaiachain/kaia/common"
)

// KaiaCNAPI provides an API to access Kaia CN-related
// information.
type KaiaCNAPI struct {
	cn *CN
}

// NewKaiaCNAPI creates a new Kaia protocol API for full nodes.
func NewKaiaCNAPI(e *CN) *KaiaCNAPI {
	return &KaiaCNAPI{e}
}

// Rewardbase is the address that consensus rewards will be send to
func (api *KaiaCNAPI) Rewardbase() (common.Address, error) {
	return api.cn.Rewardbase()
}
