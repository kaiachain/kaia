// Modifications Copyright 2026 The Kaia Authors
// This file is part of the Kaia library.

package accountkey

import (
	"encoding/json"
	"testing"

	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/require"
)

func TestAccountKeyRoleBasedDecodeRLP(t *testing.T) {
	legacy := NewAccountKeySerializerWithAccountKey(NewAccountKeyLegacy())
	valid := NewAccountKeySerializerWithAccountKey(
		NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{NewAccountKeyLegacy()}),
	)
	nested := NewAccountKeySerializerWithAccountKey(
		NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{
			NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{NewAccountKeyLegacy()}),
		}),
	)

	for _, tc := range []struct {
		name string
		key  *AccountKeySerializer
		err  error
	}{
		{name: "legacy", key: legacy},
		{name: "role based", key: valid},
		{name: "nested role based", key: nested, err: kerrors.ErrNestedCompositeType},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := rlp.EncodeToBytes(tc.key)
			require.NoError(t, err)

			decoded := NewAccountKeySerializer()
			err = rlp.DecodeBytes(encoded, decoded)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}

func TestAccountKeyRoleBasedUnmarshalJSON(t *testing.T) {
	valid := NewAccountKeySerializerWithAccountKey(
		NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{NewAccountKeyLegacy()}),
	)
	nested := NewAccountKeySerializerWithAccountKey(
		NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{
			NewAccountKeyRoleBasedWithValues(AccountKeyRoleBased{NewAccountKeyLegacy()}),
		}),
	)

	for _, tc := range []struct {
		name string
		key  *AccountKeySerializer
		err  error
	}{
		{name: "role based", key: valid},
		{name: "nested role based", key: nested, err: kerrors.ErrNestedCompositeType},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.key)
			require.NoError(t, err)

			decoded := NewAccountKeySerializer()
			err = json.Unmarshal(encoded, decoded)
			if tc.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.err)
			}
		})
	}
}
