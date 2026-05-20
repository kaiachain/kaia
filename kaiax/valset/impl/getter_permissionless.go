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
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/valset"
)

// getCouncil: Council = {ValActive, ValPaused}.
func (v *ValsetModule) getCouncil(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	return nodes.Council().Addresses(), nil
}

// committeeWithFallback returns the committee at the given block.
// Committee = {VA} - {Suspended}.
//
// Safety fallback: when len(committee)=0 due to {ValActive} ⊆ SuspendedSet, fallback to all VA.
func committeeWithFallback(nodes valset.NodeMap) (committee valset.NodeMap, fellBack bool) {
	committee = nodes.Committee()
	if len(committee) > 0 {
		return committee, false
	}

	// ignore SuspendedSet
	return nodes.FilterByState(valset.ValActive), true
}

// getQualifiedValidators: qualified = committee = {VA} − {Suspended} (with safety fallback).
func (v *ValsetModule) getQualifiedValidators(num uint64) (*valset.AddressSet, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	committee, fellBack := committeeWithFallback(nodes)
	if fellBack && v.lastSuspendFallbackLog != num {
		logger.Warn("all ValActive are suspended, ignoring suspended set for committee", "num", num)
		v.lastSuspendFallbackLog = num
	}
	return valset.NewAddressSet(committee.Addresses()), nil
}

// getCommittee: committee: qualified = {VA} − suspended
func (v *ValsetModule) getCommittee(num uint64) ([]common.Address, error) {
	qualified, err := v.getQualifiedValidators(num)
	if err != nil {
		return nil, err
	}
	return qualified.List(), nil
}

// getDemoted: demoted = council − qualified = {VP} ∪ {Suspended}.
func (v *ValsetModule) getDemoted(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	council := valset.NewAddressSet(nodes.Council().Addresses())
	committee, _ := committeeWithFallback(nodes)
	return council.Subtract(valset.NewAddressSet(committee.Addresses())).List(), nil
}

// getCandTesting: CandTesting = {CT}
func (v *ValsetModule) getCandTesting(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	return nodes.FilterByState(valset.CandTesting).Addresses(), nil
}

// getCNPeers: CNPeers = {VA, VR, VP, CR, CT}.
func (v *ValsetModule) getCNPeers(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	return nodes.CNPeers().Addresses(), nil
}

// getHeaderGovVoters: HeaderGovVoters = {VA} - {Suspended}.
func (v *ValsetModule) getHeaderGovVoters(num uint64) ([]common.Address, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	return nodes.HeaderGovVoters().Addresses(), nil
}

// getNodeByState filters nodes by state. Empty states means "all".
// Returns a defensive copy so callers can mutate the result freely.
func (v *ValsetModule) getNodeByState(num uint64, states []valset.NodeState) (valset.NodeMap, error) {
	nodes, err := v.getNodes(num)
	if err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nodes.Copy(), nil
	}
	return nodes.FilterByState(states...).Copy(), nil
}
