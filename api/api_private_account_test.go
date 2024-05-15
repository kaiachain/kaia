// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package api

import (
	"errors"
	"os"
	"testing"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/require"
)

// TestPrivateAccountAPI_ImportRawKey tests ImportRawKey() and ReplaceRawKey().
func TestPrivateAccountAPI_ImportRawKey(t *testing.T) {
	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	// To get JSON files use below.
	// keydir := filepath.Join(".", "keystore")
	keydir, err := os.MkdirTemp("", "kaia-test-api-")
	require.NoError(t, err)
	defer os.RemoveAll(keydir)

	// Assemble the account manager and supported backends
	backends := []accounts.Backend{
		keystore.NewKeyStore(keydir, scryptN, scryptP),
	}

	api := PrivateAccountAPI{
		am:        accounts.NewManager(backends...),
		nonceLock: new(AddrLocker),
		b:         nil,
	}

	// 1. Import private key only.
	{
		addr, err := api.ImportRawKey("aebb680a5e596c1d1a01bac78a3985b62c685c5e995d780c176138cb2679ba3e", "1234")
		require.NoError(t, err)

		require.Equal(t, common.HexToAddress("0x819104a190255e0cedbdd9d5f59a557633d79db1"), addr)
	}

	// 2. Import Kaia Wallet Key. Since the same address is already registered, it should fail.
	{
		_, err := api.ImportRawKey("f8cc7c3813ad23817466b1802ee805ee417001fcce9376ab8728c92dd8ea0a6b0x000x819104a190255e0cedbdd9d5f59a557633d79db1", "1234")
		require.Equal(t, errors.New("account already exists"), err)
	}

	// 3. Replace Kaia Wallet key. It should work.
	{
		addr, err := api.ReplaceRawKey("f8cc7c3813ad23817466b1802ee805ee417001fcce9376ab8728c92dd8ea0a6b0x000x819104a190255e0cedbdd9d5f59a557633d79db1", "1234", "1234")
		require.NoError(t, err)

		require.Equal(t, common.HexToAddress("0x819104a190255e0cedbdd9d5f59a557633d79db1"), addr)
	}

	// 4. Allowable Wallet key type is 0x00 only.
	{
		_, err := api.ImportRawKey("f8cc7c3813ad23817466b1802ee805ee417001fcce9376ab8728c92dd8ea0a6b0x010x819104a190255e0cedbdd9d5f59a557633d79db1", "1234")
		require.Equal(t, errors.New("Kaia wallet key type must be 00."), err)
	}

	// 5. Should return an error if wrong length.
	{
		_, err := api.ImportRawKey("1ea7b7bc7f525cc936ec65e0e93f146bd6fad4b3158067ad64560defd9bba0b0x010x3b3d49ebac925797b2471c7b01108ba16bb36950", "1234")
		require.Equal(t, errors.New("invalid hex string"), err)
	}

	// 6. Should return an error if wrong length.
	{
		_, err := api.ImportRawKey("1ea7b7bc7f525cc936ec65e0e93f146bd6fad4b3158067ad64560defd9bba0b", "1234")
		require.Equal(t, errors.New("invalid hex string"), err)
	}

	// 7. Import Kaia Wallet Key.
	{
		addr, err := api.ImportRawKey("f8cc7c3813ad23817466b1802ee805ee417001fcce9376ab8728c92dd8ea0a6b0x000x819104a190255e0cedbdd9d5f59a557633d79db2", "1234")
		require.NoError(t, err)

		require.Equal(t, common.HexToAddress("0x819104a190255e0cedbdd9d5f59a557633d79db2"), addr)
	}
}
