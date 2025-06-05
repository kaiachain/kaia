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
	"encoding/hex"
	"math/big"

	"github.com/kaiachain/kaia/v2/kaiax/randao"
	"github.com/kaiachain/kaia/v2/networks/rpc"
)

func (r *RandaoModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   newRandaoAPI(r),
			Public:    true,
		},
	}
}

type RandaoAPI struct {
	r *RandaoModule
}

func newRandaoAPI(r *RandaoModule) *RandaoAPI {
	return &RandaoAPI{r: r}
}

func (api *RandaoAPI) GetBlsInfos(number rpc.BlockNumber) (map[string]interface{}, error) {
	bn := big.NewInt(number.Int64())
	if number == rpc.LatestBlockNumber || number == rpc.PendingBlockNumber {
		bn = big.NewInt(api.r.Chain.CurrentBlock().Number().Int64())
	}
	if !api.r.ChainConfig.IsRandaoForkEnabled(bn) {
		return nil, randao.ErrBeforeRandaoFork
	}

	infos, err := api.r.getAllCached(bn)
	if err != nil {
		return nil, err
	}

	blsInfos := make(map[string]interface{})
	for addr, info := range infos {
		// hexlify publicKey and pop
		blsInfos[addr.Hex()] = map[string]interface{}{
			"publicKey": hex.EncodeToString(info.PublicKey),
			"pop":       hex.EncodeToString(info.Pop),
			"verifyErr": info.VerifyErr,
		}
	}
	return blsInfos, nil
}
