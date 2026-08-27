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

package accountkey

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPublicKey(t *testing.T) *AccountKeyPublic {
	t.Helper()
	k, err := crypto.GenerateKey()
	require.NoError(t, err)
	return NewAccountKeyPublicWithValue(&k.PublicKey)
}

// TestAccountKeyRoleBased_UpdateLength verifies that, from the Permissionless fork, a
// role-based key update stores a key set whose length follows the new key, so that any role
// not set by the new key maps to RoleTransaction by default. Explicit AccountKeyNil keeps a
// role unchanged, and growing to a longer key is unaffected. For blocks before the fork, the
// gate retains the previous result unchanged.
func TestAccountKeyRoleBased_UpdateLength(t *testing.T) {
	blockBeforeHF := uint64(4)
	blockAfterHF := uint64(6)
	fork.SetHardForkBlockNumberConfig(&params.ChainConfig{PermissionlessCompatibleBlock: big.NewInt(5)})
	defer fork.ClearHardForkBlockNumberConfig()

	keyA, keyB, keyC := newTestPublicKey(t), newTestPublicKey(t), newTestPublicKey(t)
	keyD := newTestPublicKey(t)

	newThree := func() *AccountKeyRoleBased {
		return NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyA, keyB, keyC})
	}

	t.Run("shorter new key retains previous result before Permissionless", func(t *testing.T) {
		a := newThree()
		require.NoError(t, a.Update(NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyD}), blockBeforeHF))
		require.Equal(t, 3, len(*a))
		assert.True(t, (*a)[RoleTransaction].Equal(keyD))
		assert.True(t, (*a)[RoleAccountUpdate].Equal(keyB))
		assert.True(t, (*a)[RoleFeePayer].Equal(keyC))
	})

	t.Run("shorter new key from Permissionless tracks new key length", func(t *testing.T) {
		a := newThree()
		require.NoError(t, a.Update(NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyD}), blockAfterHF))
		require.Equal(t, 1, len(*a))
		assert.True(t, (*a)[RoleTransaction].Equal(keyD))
		// Omitted roles map to RoleTransaction.
		assert.True(t, a.getDefaultKey().Equal(keyD))
	})

	t.Run("explicit AccountKeyNil keeps a role from Permissionless", func(t *testing.T) {
		a := newThree()
		upd := NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyD, NewAccountKeyNil(), NewAccountKeyNil()})
		require.NoError(t, a.Update(upd, blockAfterHF))
		require.Equal(t, 3, len(*a))
		assert.True(t, (*a)[RoleTransaction].Equal(keyD))
		assert.True(t, (*a)[RoleAccountUpdate].Equal(keyB))
		assert.True(t, (*a)[RoleFeePayer].Equal(keyC))
	})

	t.Run("longer new key appends roles regardless of fork", func(t *testing.T) {
		for _, bn := range []uint64{blockBeforeHF, blockAfterHF} {
			a := NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyA})
			require.NoError(t, a.Update(NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{keyD, keyB, keyC}), bn))
			require.Equal(t, 3, len(*a))
			assert.True(t, (*a)[RoleTransaction].Equal(keyD))
			assert.True(t, (*a)[RoleAccountUpdate].Equal(keyB))
			assert.True(t, (*a)[RoleFeePayer].Equal(keyC))
		}
	})
}
