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
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
)

func (v *ValsetModule) VerifyHeader(header *types.Header) error {
	logger.Info("NoopModule VerifyHeader", "blockNum", header.Number.Uint64())
	return nil
}

func (v *ValsetModule) PrepareHeader(header *types.Header) error {
	logger.Info("NoopModule PrepareHeader", "blockNum", header.Number.Uint64())
	return nil
}

func (v *ValsetModule) FinalizeHeader(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	logger.Info("NoopModule FinalizeHeader", "blockNum", header.Number.Uint64())
	return nil
}
