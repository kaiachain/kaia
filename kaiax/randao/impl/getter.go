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
	"math/big"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/kaiax/randao"
)

func (r *RandaoModule) GetBlsPubkey(proposer common.Address, num *big.Int) (bls.PublicKey, error) {
	infos, err := r.getAllCached(num)
	if err != nil {
		return nil, err
	}

	info, ok := infos[proposer]
	if !ok {
		return nil, randao.ErrNoBlsPub
	}
	if info.VerifyErr != nil {
		return nil, info.VerifyErr
	}
	return bls.PublicKeyFromBytes(info.PublicKey)
}

// getParentBlockStateDB retrieves both the StateDB and block number of the parent block
// for the given current block number. This function handles the calculation of the parent
// block number and retrieval of its state in a single operation.
func (r *RandaoModule) getParentBlockStateDB(blockNum *big.Int) (*state.StateDB, *big.Int, error) {
	if common.Big0.Cmp(blockNum) == 0 {
		return nil, nil, randao.ErrZeroBlockNumber
	}

	parentNum := new(big.Int).Sub(blockNum, common.Big1)
	pHeader := r.Chain.GetHeaderByNumber(parentNum.Uint64())

	// Validate parent header existence
	if pHeader == nil {
		return nil, nil, consensus.ErrUnknownAncestor
	}

	// If no state exist at block number `parentNum`,
	// return the error `consensus.ErrPrunedAncestor`
	statedb, err := r.Chain.StateAt(pHeader.Root)
	if err != nil {
		return nil, nil, consensus.ErrPrunedAncestor
	}

	return statedb, parentNum, nil
}

func (r *RandaoModule) getAllCached(num *big.Int) (system.BlsPublicKeyInfos, error) {
	// First check the block number based cache
	if item, ok := r.blsPubkeyCache.Get(num.Uint64()); ok {
		logger.Trace("BlsPublicKeyInfos cache hit", "number", num.Uint64())
		return item.(system.BlsPublicKeyInfos), nil
	}

	res, err, _ := r.sfGroup.Do(num.String(), func() (interface{}, error) {
		start := time.Now()

		backend := backends.NewBlockchainContractBackend(r.Chain, nil, nil)

		var kip113Addr common.Address
		var statedb *state.StateDB
		var parentNum *big.Int

		// Get the parent block's statedb and block number
		statedb, parentNum, err := r.getParentBlockStateDB(num)
		if err != nil {
			return nil, err
		}

		// Because the system contract Registry is installed at Finalize() of RandaoForkBlock,
		// it is not possible to read KIP113 address from the Registry at RandaoForkBlock.
		// Hence the ChainConfig fallback.
		if r.ChainConfig.IsRandaoForkBlock(num) {
			var ok bool
			kip113Addr, ok = r.ChainConfig.RandaoRegistry.Records[system.Kip113Name]
			if !ok {
				return nil, randao.ErrMissingKIP113
			}
		} else if r.ChainConfig.IsRandaoForkEnabled(num) {
			kip113Addr, err = system.ReadActiveAddressFromRegistry(backend, system.Kip113Name, parentNum)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, randao.ErrBeforeRandaoFork
		}

		// Check storage root cache
		systemRegistryAddr := system.RegistryAddr
		kip113Root := statedb.GetStorageRoot(kip113Addr)
		systemRegistryRoot := statedb.GetStorageRoot(systemRegistryAddr)
		storageKey := kip113Root.Hex() + ":" + systemRegistryRoot.Hex()

		if item, ok := r.storageRootCache.Get(storageKey); ok {
			logger.Trace("BlsPublicKeyInfos storage root cache hit",
				"number", num.Uint64(),
				"kip113Root", kip113Root.Hex(),
				"systemRegistryRoot", systemRegistryRoot.Hex())

			infos := item.(system.BlsPublicKeyInfos)
			r.blsPubkeyCache.Add(num.Uint64(), infos)
			return infos, nil
		}

		// Cache miss - read data from contracts
		infos, err := system.ReadKip113All(backend, kip113Addr, parentNum)
		if err != nil {
			return nil, err
		}

		logger.Trace("BlsPublicKeyInfos cache miss",
			"number", num.Uint64(),
			"kip113Root", kip113Root.Hex(),
			"systemRegistryRoot", systemRegistryRoot.Hex(),
			"elapsed", time.Since(start))

		// Update both caches
		r.blsPubkeyCache.Add(num.Uint64(), infos)
		r.storageRootCache.Add(storageKey, infos)

		return infos, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(system.BlsPublicKeyInfos), nil
}
