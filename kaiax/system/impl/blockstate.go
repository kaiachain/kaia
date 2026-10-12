// Copyright 2026 The Kaia Authors
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
	bcsystem "github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/params"
)

func (m *SystemModule) InitializeState(header *types.Header, state *state.StateDB) {}

func (m *SystemModule) FinalizeState(header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) error {
	chainConfig := m.Chain.Config()

	// Install Randao registry at fork block.
	if chainConfig.IsRandaoForkBlock(header.Number) {
		if err := bcsystem.InstallRegistry(state, chainConfig.RandaoRegistry); err != nil {
			return err
		}
	}

	// RebalanceTreasury can modify the global state, so use the current stateDB in-place.
	if chainConfig.IsKIP160ForkBlock(header.Number) || chainConfig.IsKIP103ForkBlock(header.Number) {
		rebalanceResult, err := bcsystem.RebalanceTreasury(state, m.Chain, header)
		if err != nil {
			logger.Error("failed to execute treasury rebalancing. State not changed", "err", err)
		} else {
			// Memo format differs between KIP-103 and KIP-160.
			isKIP103 := chainConfig.IsKIP103ForkBlock(header.Number)
			logger.Info("successfully executed treasury rebalancing", "memo", string(rebalanceResult.Memo(isKIP103)))
		}
	}

	// Replace the Mainnet credit contract.
	if chainConfig.IsKaiaForkBlockParent(header.Number) {
		if chainConfig.ChainID.Uint64() == params.MainnetNetworkId && state.GetCode(bcsystem.NonExistentAddress) != nil {
			if err := state.SetCode(bcsystem.NonExistentAddress, bcsystem.MainnetCreditV2Code); err != nil {
				return err
			}
			logger.Info("Replaced CypressCredit with CypressCreditV2", "blockNum", header.Number.Uint64())
		}
	}

	// Restore Mainnet credit contract address (0x0) back to pure EOA at Osaka hardfork.
	if chainConfig.IsOsakaForkBlockParent(header.Number) {
		if chainConfig.ChainID.Uint64() == params.MainnetNetworkId && state.GetCode(bcsystem.NonExistentAddress) != nil {
			prevNonce := state.GetNonce(bcsystem.NonExistentAddress)
			state.CreateEOA(bcsystem.NonExistentAddress, false, accountkey.NewAccountKeyLegacy())
			state.SetNonce(bcsystem.NonExistentAddress, prevNonce) // Preserve account counters across account type migration.
			logger.Info("Restored Mainnet credit address to EOA", "blockNum", header.Number.Uint64())
		}
	}

	return nil
}
