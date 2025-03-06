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
	"context"
	"errors"
	"fmt"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/rlp"
)

func (b *BuilderModule) APIs() []rpc.API {
	return []rpc.API{
		{
			Namespace: "kaia",
			Version:   "1.0",
			Service:   NewBuilderAPI(b),
			Public:    true,
		},
	}
}

type BuilderAPI struct {
	b *BuilderModule
}

func NewBuilderAPI(b *BuilderModule) *BuilderAPI {
	return &BuilderAPI{b}
}

func (s *BuilderAPI) SendRawTransactions(ctx context.Context, inputs []hexutil.Bytes) ([]common.Hash, error) {
	hash := []common.Hash{}
	errs := []error{}

	if len(inputs) == 0 {
		hash = append(hash, common.Hash{})
		return hash, errors.New("Empty input")
	}

	for i, input := range inputs {
		if len(input) == 0 {
			hash = append(hash, common.Hash{})
			errs = append(errs, fmt.Errorf("Index %d: empty input", i))
			break
		}
		if 0 < input[0] && input[0] < 0x7f {
			input = append([]byte{byte(types.EthereumTxTypeEnvelope)}, input...)
		}
		tx := new(types.Transaction)
		if err := rlp.DecodeBytes(input, tx); err != nil {
			hash = append(hash, common.Hash{})
			errs = append(errs, fmt.Errorf("Index %d: %w", i, err))
			break
		}
		if err := s.b.Backend.SendTx(ctx, tx); err != nil {
			hash = append(hash, common.Hash{})
			errs = append(errs, fmt.Errorf("Index %d: %w", i, err))
			break
		}
		hash = append(hash, tx.Hash())
	}

	return hash, errors.Join(errs...)
}
