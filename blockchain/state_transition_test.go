// Modifications Copyright 2024 The Kaia Authors
// Copyright 2018 The klaytn Authors
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

package blockchain

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/blockchain/vm"
	mock_vm "github.com/kaiachain/kaia/blockchain/vm/mocks"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestGetVMerrFromReceiptStatus(t *testing.T) {
	testData := []struct {
		status           uint
		expectMatchError error
	}{
		{types.ReceiptStatusFailed, ErrInvalidReceiptStatus},
		{types.ReceiptStatusLast, ErrInvalidReceiptStatus},
		{types.ReceiptStatusSuccessful, nil},
		{types.ReceiptStatusErrDefault, ErrVMDefault},
	}

	for _, tc := range testData {
		result := ExecutionResult{VmExecutionStatus: tc.status}
		assert.Equal(t, tc.expectMatchError, result.Unwrap())
	}
}

func TestGetReceiptStatusFromVMerr(t *testing.T) {
	status := getReceiptStatusFromErrTxFailed(nil)
	expectedStatus := types.ReceiptStatusSuccessful
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}

	status = getReceiptStatusFromErrTxFailed(vm.ErrMaxCodeSizeExceeded)
	expectedStatus = types.ReceiptStatuserrMaxCodeSizeExceed
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}

	// Unknown VM error
	status = getReceiptStatusFromErrTxFailed(errors.New("Unknown VM error"))
	expectedStatus = types.ReceiptStatusErrDefault
	if status != expectedStatus {
		t.Fatalf("Invalid receipt status, want %d, got %d", expectedStatus, status)
	}
}

// TestPrintErrorCodeTable prints the error code table in a format of a markdown table.
func TestPrintErrorCodeTable(t *testing.T) {
	if testing.Verbose() {
		fmt.Println("| ErrorCode | Description |")
		fmt.Println("|---|---|")
		for i := uint(types.ReceiptStatusErrDefault); i < types.ReceiptStatusLast; i++ {
			fmt.Printf("|0x%02x|%s|\n", i, receiptstatus2errTxFailed[i])
		}
	}
}

// Coverage is detected by calling functions that modify state.
func TestStateTransition_checkAuthorizationList(t *testing.T) {
	var (
		authorityKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key, _          = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		authority       = crypto.PubkeyToAddress(authorityKey.PublicKey)
		addr            = crypto.PubkeyToAddress(key.PublicKey)
		aa              = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		bb              = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		zeroAddress     = common.HexToAddress("0x0000000000000000000000000000000000000000")
		toAuthorityTx   = types.NewTransaction(uint64(0), authority, nil, 0, nil, nil)
		toAddrTx        = types.NewTransaction(uint64(0), addr, nil, 0, nil, nil)
	)

	tests := []struct {
		name              string
		makeAuthList      func() types.AuthorizationList
		msg               Message
		expectedMockCalls func(m *mock_vm.MockStateDB)
	}{
		// Cases: success to set code
		{
			name: "success (minimum)",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(false)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, types.AddressToDelegation(aa))
			},
		},
		{
			name: "success (case of refund)",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   1,
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(true)
				m.EXPECT().GetKey(authority).Return(accountkey.NewAccountKeyLegacy())
				m.EXPECT().AddRefund(params.CallNewAccountGas - params.TxAuthTupleGas)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, types.AddressToDelegation(aa))
			},
		},
		{
			name: "success (address 0x0000000000000000000000000000000000000000)",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: zeroAddress,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(false)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, []byte{})
			},
		},
		{
			name: "success (to == authority)",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAuthorityTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(false)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, types.AddressToDelegation(aa))
				m.EXPECT().AddAddressToAccessList(aa)
			},
		},
		{
			name: "success (don't ecrecover authority)",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)

				// The msg is tampered with so a different pubkey is ecrecovered.
				auth.Address = bb
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				notAuthority := gomock.Not(authority)
				m.EXPECT().AddAddressToAccessList(notAuthority)
				m.EXPECT().GetCode(notAuthority).Return(nil)
				m.EXPECT().GetNonce(notAuthority).Return(uint64(1))
				m.EXPECT().Exist(notAuthority).Return(false)
				m.EXPECT().IncNonce(notAuthority)
				m.EXPECT().SetCodeToEOA(notAuthority, types.AddressToDelegation(bb))
			},
		},
		// Cases: fail to set code
		{
			name: "invalid chainId",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: uint64(10),
					Address: aa,
					Nonce:   1,
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "nonce is uint64 max value",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(18446744073709551615),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "error in Authority",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				auth.V = uint8(10)
				return types.AuthorizationList{*auth}
			},
			msg: toAuthorityTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "exist some code",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return([]byte{42, 42})
			},
		},
		{
			name: "invalid nonce",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(10),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
			},
		},
		{
			name: "signer's key was updated",
			makeAuthList: func() types.AuthorizationList {
				auth, err := types.SignAuth(&types.Authorization{
					ChainID: params.TestChainConfig.ChainID.Uint64(),
					Address: aa,
					Nonce:   uint64(1),
				}, authorityKey)
				assert.NoError(t, err)
				return types.AuthorizationList{*auth}
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(true)
				m.EXPECT().GetKey(authority).Return(accountkey.NewAccountKeyPublic())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStateDB := mock_vm.NewMockStateDB(mockCtrl)
			tt.expectedMockCalls(mockStateDB)

			authList := tt.makeAuthList()

			header := &types.Header{Number: big.NewInt(0), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
			blockContext := NewEVMBlockContext(header, nil, &common.Address{}) // stub author (COINBASE) to 0x0
			txContext := NewEVMTxContext(tt.msg, header, params.TestChainConfig)
			evm := vm.NewEVM(blockContext, txContext, mockStateDB, params.TestChainConfig, &vm.Config{Debug: true})

			// Verify that the expected mockStateDB's calls are being made.
			NewStateTransition(evm, tt.msg).checkAuthorizationList(authList, *tt.msg.To())
		})
	}
}
