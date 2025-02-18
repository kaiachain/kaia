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
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/require"
)

func TestIsApproveTx(t *testing.T) {
	log.EnableLogForTest(log.LvlTrace, log.LvlTrace)
	testcases := []struct {
		tx *types.Transaction
		ok bool
	}{
		{ // Legacy TestToken.approve(SwapRouter, 1000000)
			types.NewTransaction(0, common.HexToAddress("0xabcd"), big.NewInt(0), 1000000, big.NewInt(0),
				hexutil.MustDecode("0x095ea7b3000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000f4240")),
			true,
		},
	}

	g := NewGaslessModule()
	key, _ := crypto.GenerateKey()
	g.Init(&InitOpts{
		ChainConfig: &params.ChainConfig{ChainID: big.NewInt(1)},
		NodeKey:     key,
	})
	for _, tc := range testcases {
		ok := g.IsApproveTx(tc.tx)
		require.Equal(t, tc.ok, ok)
	}
}
