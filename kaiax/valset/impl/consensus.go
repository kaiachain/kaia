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
	"slices"

	"github.com/kaiachain/kaia/blockchain/types"
)

// VerifyHeader checks that the validator set written into the header (via the
// IstanbulSealer) matches the qualified validator set this node would compute.
// This is enforced only after the permissionless hardfork.
//
// This guards against a malicious proposer writing an arbitrary validator set
// into header.Extra and getting away with it because committed seals only
// verify against whatever set the header itself declares.
func (v *ValsetModule) VerifyHeader(header *types.Header, _ *types.Header) error {
	if !v.Chain.Config().IsPermissionlessForkEnabled(header.Number) {
		return nil
	}
	num := header.Number.Uint64()
	headerVals, err := v.Chain.Sealer().Validators(header)
	if err != nil {
		return err
	}
	qualified, err := v.getQualifiedValidators(num)
	if err != nil {
		return err
	}
	if !slices.Equal(headerVals, qualified.List()) {
		return errMismatchedValidators
	}
	return nil
}

// PrepareHeader is a no-op. The validator set is written by
// blockchain.PrepareHeader via Sealer().WriteValidators, before this hook runs.
func (v *ValsetModule) PrepareHeader(*types.Header) error { return nil }
