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
	"fmt"
	"math/big"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/staking"
	"github.com/kaiachain/kaia/params"
)

const ( // Numeric Type IDs used by AddressBook.getAllAddress(), and in turn by MultiCall contract.
	CN_NODE_ID_TYPE         = 0
	CN_STAKING_ADDRESS_TYPE = 1
	CN_REWARD_ADDRESS_TYPE  = 2
	POC_CONTRACT_TYPE       = 3 // The AddressBook's pocContractAddress field repurposed as PoC, KGF, KFF, KIF
	KIR_CONTRACT_TYPE       = 4 // The AddressBook's kirContractAddress field repurposed as KIR, KCF, KEF
)

type clRegistryResult struct {
	NodeIds        []common.Address
	ClPools        []common.Address
	ClStakings     []common.Address
	StakingAmounts []*big.Int
}

func (s *StakingModule) GetStakingInfo(num uint64) (*staking.StakingInfo, error) {
	isKaia := s.ChainConfig.IsKaiaForkEnabled(new(big.Int).SetUint64(num))
	sourceNum := sourceBlockNum(num, isKaia, s.stakingInterval)

	// Try cache first
	if si, ok := s.stakingInfoCache.Get(sourceNum); ok {
		return si.(*staking.StakingInfo), nil
	}

	// Only before Kaia, try the database
	if !isKaia {
		if si := ReadStakingInfo(s.ChainKv, sourceNum); si != nil {
			s.stakingInfoCache.Add(sourceNum, si)
			return si, nil
		}
	}

	// Read from the state
	si, err := s.getFromStateByNumber(sourceNum)
	if err != nil {
		return nil, err
	}

	// Only before Kaia, write to database
	if !isKaia {
		WriteStakingInfo(s.ChainKv, sourceNum, si)
	}

	// Cache it
	s.stakingInfoCache.Add(sourceNum, si)
	return si, nil
}

// Read the staking status from the blockchain state.
func (s *StakingModule) getFromStateByNumber(num uint64) (*staking.StakingInfo, error) {
	header := s.Chain.GetHeaderByNumber(num)
	if header == nil {
		return nil, fmt.Errorf("failed to get header for block number %d", num)
	}

	// If found in side state, no bother getting from the state.
	if si := s.preloadBuffer.GetInfo(header.Root); si != nil { // Try side state
		return si, nil
	}

	// Otherwise bring up the state from the database.
	statedb, err := s.Chain.StateAt(header.Root)
	if err != nil {
		return nil, fmt.Errorf("failed to get state for block number %d: %v", num, err)
	}
	return s.getFromState(header, statedb)
}

// Efficiently read addresses and balances from the AddressBook in one EVM call.
// Works by temporarily injecting the MultiCallContract to a copied state.
func (s *StakingModule) getFromState(header *types.Header, statedb *state.StateDB) (*staking.StakingInfo, error) {
	isPrague := s.ChainConfig.IsPragueForkEnabled(header.Number)
	num := header.Number.Uint64()

	// Bail out if AddressBook is not installed.
	// This is a common case for private nets.
	if statedb.GetCode(system.AddressBookAddr) == nil {
		logger.Trace("AddressBook not installed", "sourceNum", num)
		return emptyStakingInfo(num), nil
	}

	// Now we're safe to call the MultiCall contract.
	contract, err := system.NewMultiCallContractCaller(statedb, s.Chain, header)
	if err != nil {
		return nil, staking.ErrMultiCallCall(err)
	}

	// Get staking info from AddressBook
	callOpts := &bind.CallOpts{BlockNumber: header.Number}
	abRes, err := contract.MultiCallStakingInfo(callOpts)
	if err != nil {
		return nil, staking.ErrAddressBookCall(err)
	}

	// Get CL registry info if Prague fork is enabled
	var clRes clRegistryResult
	if isPrague {
		// If Registry is not installed, do not handle CL staking info.
		if statedb.GetCode(system.RegistryAddr) == nil {
			logger.Trace("Registry not installed", "sourceNum", num)
		} else {
			// Note that if CLRegistry is not registered in Registry,
			// it will return empty result and no error.
			clRes, err = contract.MultiCallDPStakingInfo(callOpts)
			if err != nil {
				return nil, staking.ErrCLRegistryCall(err)
			}
		}
	}

	return parseCallResult(num, abRes.TypeList, abRes.AddressList, abRes.StakingAmounts, clRes)
}

func parseCallResult(num uint64, types []uint8, addrs []common.Address, amounts []*big.Int, clRes clRegistryResult) (*staking.StakingInfo, error) {
	// Sanity check.
	if len(types) == 0 && len(addrs) == 0 {
		// This is an expected behavior when the AddressBook contract is not activated yet.
		logger.Trace("returning empty staking info because AddressBook is not activated", "sourceNum", num)
		return emptyStakingInfo(num), nil
	}
	if len(types) != len(addrs) {
		logger.Error("length of type list and address list differ", "sourceNum", num, "typeLen", len(types), "addrLen", len(addrs))
		return nil, staking.ErrAddressBookResult
	}
	if len(clRes.NodeIds) != len(clRes.ClPools) || len(clRes.NodeIds) != len(clRes.ClStakings) || len(clRes.NodeIds) != len(clRes.StakingAmounts) {
		logger.Error("length of CL registry result fields differ", "sourceNum", num, "nodeLen", len(clRes.NodeIds), "poolLen", len(clRes.ClPools), "stakingLen", len(clRes.ClStakings), "amountLen", len(clRes.StakingAmounts))
		return nil, staking.ErrCLRegistryResult
	}

	// Collect the AddressBook results to StakingInfo fields.
	var (
		nodeIds          []common.Address
		stakingContracts []common.Address
		rewardAddrs      []common.Address
		kefAddr          common.Address
		kifAddr          common.Address

		stakingAmounts = make([]uint64, len(amounts))
	)
	for i, ty := range types {
		switch ty {
		case CN_NODE_ID_TYPE:
			nodeIds = append(nodeIds, addrs[i])
		case CN_STAKING_ADDRESS_TYPE:
			stakingContracts = append(stakingContracts, addrs[i])
		case CN_REWARD_ADDRESS_TYPE:
			rewardAddrs = append(rewardAddrs, addrs[i])
		// Caution: not to confuse (POC, KIR) order
		case POC_CONTRACT_TYPE:
			kifAddr = addrs[i]
		case KIR_CONTRACT_TYPE:
			kefAddr = addrs[i]
		default:
			logger.Error("unknown entry type", "sourceNum", num, "type", ty)
			return nil, staking.ErrAddressBookResult
		}
	}
	for i, a := range amounts {
		stakingAmounts[i] = big.NewInt(0).Div(a, big.NewInt(params.KAIA)).Uint64()
	}

	// Collect the CL registry results to StakingInfo fields.
	// If there is no CL registry result, it will be nil.
	var clStakingInfos staking.CLStakingInfos
	if len(clRes.NodeIds) > 0 {
		clStakingInfos = make(staking.CLStakingInfos, len(clRes.NodeIds))
		for i := range clRes.NodeIds {
			clStakingInfos[i] = &staking.CLStakingInfo{
				CLNodeId:        clRes.NodeIds[i],
				CLPoolAddr:      clRes.ClPools[i],
				CLRewardAddr:    clRes.ClStakings[i],
				CLStakingAmount: big.NewInt(0).Div(clRes.StakingAmounts[i], big.NewInt(params.KAIA)).Uint64(),
			}
		}
	}

	// Sanity check
	if len(nodeIds) != len(stakingContracts) || len(nodeIds) != len(rewardAddrs) || len(nodeIds) != len(amounts) ||
		common.EmptyAddress(kefAddr) || common.EmptyAddress(kifAddr) {
		// This is an expected behavior when the AddressBook contract is not activated yet.
		logger.Trace("returning empty staking info because AddressBook is not activated", "sourceNum", num)
		return emptyStakingInfo(num), nil
	}

	return &staking.StakingInfo{
		SourceBlockNum:   num,
		NodeIds:          nodeIds,
		StakingContracts: stakingContracts,
		RewardAddrs:      rewardAddrs,
		KEFAddr:          kefAddr,
		KIFAddr:          kifAddr,
		StakingAmounts:   stakingAmounts,
		CLStakingInfos:   clStakingInfos,
	}, nil
}

func emptyStakingInfo(num uint64) *staking.StakingInfo {
	return &staking.StakingInfo{SourceBlockNum: num}
}

func sourceBlockNum(num uint64, isKaia bool, interval uint64) uint64 {
	if isKaia {
		if num == 0 {
			return 0
		} else {
			return num - 1
		}
	} else {
		if num <= 2*interval {
			return 0
		} else {
			// Simplified from the previous implementation:
			// if (num % interval) == 0, return num - 2*interval
			// else return num - interval - (num % interval)
			return roundDown(num-1, interval) - interval
		}
	}
}

func roundDown(n, p uint64) uint64 {
	return n - (n % p)
}
