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

package sc

import (
	"math"
	"math/big"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/stretchr/testify/assert"
)

// TestBridgeAccountLockUnlock checks the lock/unlock functionality.
func TestBridgeAccountLockUnlock(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, pPwdStr, cPwdStr := newBridgeAccountsForTest(t, nil, config.DataDir)
	assert.Equal(t, true, bAcc.cAccount.IsUnlockedAccount())
	assert.Equal(t, true, bAcc.pAccount.IsUnlockedAccount())

	lockAccountsWithCheck := func(t *testing.T, bAcc *BridgeAccounts) {
		{
			err := bAcc.cAccount.LockAccount()
			assert.NoError(t, err)
			assert.Equal(t, false, bAcc.cAccount.IsUnlockedAccount())
		}
		{
			err := bAcc.pAccount.LockAccount()
			assert.NoError(t, err)
			assert.Equal(t, false, bAcc.pAccount.IsUnlockedAccount())
		}
	}

	testCases := []struct {
		name        string
		preLockCnt  int
		duration    int64
		parentPass  string
		childPass   string
		expectedErr error
	}{
		{"invalid timeout", 3, int64(uint64(time.Duration(math.MaxInt64)/time.Second) + 1), pPwdStr, cPwdStr, errUnlockDurationTooLarge},
		{"wrong password duration 0", 1, 0, pPwdStr[:3], cPwdStr[:3], keystore.ErrDecrypt},
		{"success duration 0", 1, 0, pPwdStr, cPwdStr, nil},
		{"success duration 1", 1, 1, pPwdStr, cPwdStr, nil},
		{"wrong password duration 1", 1, 1, pPwdStr[:3], cPwdStr[:3], keystore.ErrDecrypt},
		{"success nil duration", 1, -1, pPwdStr, cPwdStr, nil},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var durationPtr *uint64
			if tc.duration >= 0 {
				duration := uint64(tc.duration)
				durationPtr = &duration
			}
			for i := 0; i < tc.preLockCnt; i++ {
				lockAccountsWithCheck(t, bAcc)
			}

			expectedIsUnlock := tc.expectedErr == nil

			err := bAcc.pAccount.UnLockAccount(tc.parentPass, durationPtr)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, expectedIsUnlock, bAcc.pAccount.IsUnlockedAccount())

			err = bAcc.cAccount.UnLockAccount(tc.childPass, durationPtr)
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, expectedIsUnlock, bAcc.cAccount.IsUnlockedAccount())

			if tc.expectedErr != nil || durationPtr == nil || *durationPtr == 0 {
				return
			}

			deadline := time.Now().Add(time.Duration(*durationPtr)*time.Second + 500*time.Millisecond)
			for bAcc.pAccount.IsUnlockedAccount() || bAcc.cAccount.IsUnlockedAccount() {
				if time.Now().After(deadline) {
					t.Fatal("timeout waiting for unlock expiration")
				}
				runtime.Gosched()
			}
		})
	}
}

// TestBridgeAccountInformation checks if the information result is right or not.
func TestBridgeAccountInformation(t *testing.T) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, _, _ := newBridgeAccountsForTest(t, nil, config.DataDir)
	assert.Equal(t, true, bAcc.cAccount.IsUnlockedAccount())
	assert.Equal(t, true, bAcc.pAccount.IsUnlockedAccount())

	bAcc.pAccount.gasPrice = big.NewInt(100)
	bAcc.pAccount.nonce = 10
	bAcc.pAccount.chainID = big.NewInt(200)
	bAcc.pAccount.isNonceSynced = true
	err = bAcc.pAccount.LockAccount()
	assert.NoError(t, err)

	res := bAcc.GetBridgeOperators()
	assert.Equal(t, 2, len(res))

	pRes := res["parentOperator"].(map[string]interface{})
	cRes := res["childOperator"].(map[string]interface{})
	assert.Equal(t, 6, len(pRes))
	assert.Equal(t, 6, len(cRes))

	assert.Equal(t, pRes["address"], bAcc.pAccount.address)
	assert.Equal(t, pRes["nonce"], bAcc.pAccount.nonce)
	assert.Equal(t, pRes["chainID"].(*big.Int).String(), bAcc.pAccount.chainID.String())
	assert.Equal(t, pRes["gasPrice"].(*big.Int).String(), bAcc.pAccount.gasPrice.String())
	assert.Equal(t, pRes["isNonceSynced"], bAcc.pAccount.isNonceSynced)
	assert.Equal(t, pRes["isUnlocked"], bAcc.pAccount.IsUnlockedAccount())

	assert.Equal(t, cRes["address"], bAcc.cAccount.address)
	assert.Equal(t, cRes["nonce"], bAcc.cAccount.nonce)
	assert.Equal(t, cRes["chainID"].(*big.Int).String(), bAcc.cAccount.chainID.String())
	assert.Equal(t, cRes["gasPrice"].(*big.Int).String(), bAcc.cAccount.gasPrice.String())
	assert.Equal(t, cRes["isNonceSynced"], bAcc.cAccount.isNonceSynced)
	assert.Equal(t, cRes["isUnlocked"], bAcc.cAccount.IsUnlockedAccount())
}
