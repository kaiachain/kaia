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

package core

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus/bft"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/require"
)

func TestAdoptPreparedCertificateOnlyMovesLockForward(t *testing.T) {
	proposal := func(marker int64) *types.Block {
		return types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1), Time: big.NewInt(marker)})
	}
	certificate := func(round int64, block *types.Block) *bft.PreparedCertificate {
		return &bft.PreparedCertificate{
			View:     &bft.View{Sequence: big.NewInt(1), Round: big.NewInt(round)},
			Proposal: block,
		}
	}

	x, y, z := proposal(1), proposal(2), proposal(3)
	state := newRoundState(
		&bft.View{Sequence: big.NewInt(1), Round: big.NewInt(2)},
		valset.NewAddressSet(nil), x.Hash(), nil, certificate(0, x), nil, nil,
	)

	state.AdoptPreparedCertificate(certificate(1, y))
	require.Equal(t, y.Hash(), state.GetLockedHash(), "a higher prepared round must override the older lock")
	require.Equal(t, int64(1), state.LockedRound().Int64())

	state.AdoptPreparedCertificate(certificate(0, x))
	require.Equal(t, y.Hash(), state.GetLockedHash(), "an older certificate must not downgrade the lock")
	require.Equal(t, int64(1), state.LockedRound().Int64())

	state.AdoptPreparedCertificate(certificate(1, z))
	require.Equal(t, y.Hash(), state.GetLockedHash(), "an equal-round conflicting certificate must not replace the lock")
	require.Equal(t, int64(1), state.LockedRound().Int64())

	state.AdoptPreparedCertificate(&bft.PreparedCertificate{Proposal: z})
	require.Equal(t, y.Hash(), state.GetLockedHash(), "an incomplete certificate must not replace the lock")
}
