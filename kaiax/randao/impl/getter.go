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

func (r *RandaoModule) getAllCached(num *big.Int) (system.BlsPublicKeyInfos, error) {
	if item, ok := r.blsPubkeyCache.Get(num.Uint64()); ok {
		logger.Trace("BlsPublicKeyInfos cache hit", "number", num.Uint64())
		return item.(system.BlsPublicKeyInfos), nil
	}

	res, err, _ := r.sfGroup.Do(num.String(), func() (interface{}, error) {
		start := time.Now()

		backend := backends.NewBlockchainContractBackend(r.Chain, nil, nil)
		if common.Big0.Cmp(num) == 0 {
			return nil, randao.ErrZeroBlockNumber
		}
		parentNum := new(big.Int).Sub(num, common.Big1)

		var kip113Addr common.Address
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
			// If no state exist at block number `parentNum`,
			// return the error `consensus.ErrPrunedAncestor`
			pHeader := r.Chain.GetHeaderByNumber(parentNum.Uint64())
			if pHeader == nil {
				return nil, consensus.ErrUnknownAncestor
			}
			_, err := r.Chain.StateAt(pHeader.Root)
			if err != nil {
				return nil, consensus.ErrPrunedAncestor
			}
			kip113Addr, err = system.ReadActiveAddressFromRegistry(backend, system.Kip113Name, parentNum)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, randao.ErrBeforeRandaoFork
		}

		infos, err := system.ReadKip113All(backend, kip113Addr, parentNum)
		if err != nil {
			return nil, err
		}
		logger.Trace("BlsPublicKeyInfos cache miss", "number", num.Uint64(), "elapsed", time.Since(start))
		r.blsPubkeyCache.Add(num.Uint64(), infos)

		return infos, nil
	})

	if err != nil {
		return nil, err
	}

	return res.(system.BlsPublicKeyInfos), nil
}
