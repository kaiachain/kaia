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
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/holiman/uint256"
	mock_bc "github.com/kaiachain/kaia/blockchain/mocks"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/types/accountkey"
	"github.com/kaiachain/kaia/blockchain/vm"
	mock_vm "github.com/kaiachain/kaia/blockchain/vm/mocks"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/fork"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestStateTransition_validateAuthorization(t *testing.T) {
	var (
		authorityKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key, _          = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		authority       = crypto.PubkeyToAddress(authorityKey.PublicKey)
		addr            = crypto.PubkeyToAddress(key.PublicKey)
		aa              = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		zeroAddress     = common.HexToAddress("0x0000000000000000000000000000000000000000")
		toAuthorityTx   = types.NewTransaction(uint64(0), authority, nil, 0, nil, nil)
		toAddrTx        = types.NewTransaction(uint64(0), addr, nil, 0, nil, nil)
	)

	tests := []struct {
		name              string
		makeAuthorization func() types.SetCodeAuthorization
		msg               Message
		expectedAddress   common.Address
		expectedError     error
		expectedMockCalls func(m *mock_vm.MockStateDB)
	}{
		// Cases: Valid
		{
			name: "valid Authorization",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				return auth
			},
			msg:             toAddrTx,
			expectedAddress: authority,
			expectedError:   nil,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
			},
		},
		// Cases: Invalids
		{
			name: "wrong ChainID",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.NewInt(10),
					Address: aa,
					Nonce:   1,
				})
				assert.NoError(t, err)
				return auth
			},
			msg:             toAddrTx,
			expectedAddress: zeroAddress,
			expectedError:   ErrAuthorizationWrongChainID,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "nonce overflow by uint64 max value",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(18446744073709551615),
				})
				assert.NoError(t, err)
				return auth
			},
			msg:             toAddrTx,
			expectedAddress: zeroAddress,
			expectedError:   ErrAuthorizationNonceOverflow,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "invalid Signature in Authority",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				auth.V = uint8(10)
				return auth
			},
			msg:             toAuthorityTx,
			expectedAddress: zeroAddress,
			expectedError:   fmt.Errorf("%w: %v", ErrAuthorizationInvalidSignature, types.ErrInvalidSig),
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		// Cases: Invalids after getting Authority
		{
			name: "destination has code",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				return auth
			},
			msg:             toAddrTx,
			expectedAddress: authority,
			expectedError:   ErrAuthorizationDestinationHasCode,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return([]byte{42, 42})
			},
		},
		{
			name: "nonce mismatch",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(10),
				})
				assert.NoError(t, err)
				return auth
			},
			msg:             toAddrTx,
			expectedAddress: authority,
			expectedError:   ErrAuthorizationNonceMismatch,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStateDB := mock_vm.NewMockStateDB(mockCtrl)
			tt.expectedMockCalls(mockStateDB)

			header := &types.Header{Number: big.NewInt(0), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
			blockContext := NewEVMBlockContext(header, nil, &common.Address{}) // stub author (COINBASE) to 0x0
			txContext := NewEVMTxContext(tt.msg, header, params.TestChainConfig)
			evm := vm.NewEVM(blockContext, txContext, mockStateDB, params.TestChainConfig, &vm.Config{Debug: true})

			// Verify that the expected mockStateDB's calls are being made.
			auth := tt.makeAuthorization()
			actual, err := NewStateTransition(evm, tt.msg).validateAuthorization(&auth)
			require.Equal(t, tt.expectedAddress, actual)
			require.Equal(t, tt.expectedError, err)
		})
	}
}

// Coverage is detected by calling functions that modify state.
func TestStateTransition_applyAuthorization(t *testing.T) {
	var (
		authorityKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key, _          = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		authority       = crypto.PubkeyToAddress(authorityKey.PublicKey)
		addr            = crypto.PubkeyToAddress(key.PublicKey)
		aa              = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		bb              = common.HexToAddress("0x000000000000000000000000000000000000bbbb")
		zeroAddress     = common.HexToAddress("0x0000000000000000000000000000000000000000")
		toAddrTx        = types.NewTransaction(uint64(0), addr, nil, 0, nil, nil)
		rules           = params.TestChainConfig.Rules(big.NewInt(0))
	)

	tests := []struct {
		name              string
		makeAuthorization func() types.SetCodeAuthorization
		msg               Message
		expectedError     error
		expectedMockCalls func(m *mock_vm.MockStateDB)
	}{
		// Cases: success to set code
		{
			name: "success (minimum)",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				return auth
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(false)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, types.AddressToDelegation(aa), rules)
			},
		},
		{
			name: "success (case of refund)",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   1,
				})
				assert.NoError(t, err)
				return auth
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
				m.EXPECT().SetCodeToEOA(authority, types.AddressToDelegation(aa), rules)
			},
		},
		{
			name: "success (empty address 0x0000000000000000000000000000000000000000)",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: zeroAddress,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				return auth
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				m.EXPECT().AddAddressToAccessList(authority)
				m.EXPECT().GetCode(authority).Return(nil)
				m.EXPECT().GetNonce(authority).Return(uint64(1))
				m.EXPECT().Exist(authority).Return(false)
				m.EXPECT().IncNonce(authority)
				m.EXPECT().SetCodeToEOA(authority, []byte{}, rules)
			},
		},
		{
			name: "success (don't ecrecover authority)",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)

				// The msg is tampered with so a different pubkey is ecrecovered.
				auth.Address = bb
				return auth
			},
			msg: toAddrTx,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				notAuthority := gomock.Not(authority)
				m.EXPECT().AddAddressToAccessList(notAuthority)
				m.EXPECT().GetCode(notAuthority).Return(nil)
				m.EXPECT().GetNonce(notAuthority).Return(uint64(1))
				m.EXPECT().Exist(notAuthority).Return(false)
				m.EXPECT().IncNonce(notAuthority)
				m.EXPECT().SetCodeToEOA(notAuthority, types.AddressToDelegation(bb), rules)
			},
		},
		// Cases: fail to set code
		{
			name: "invalid validateAuthorization",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.NewInt(10),
					Address: aa,
					Nonce:   1,
				})
				assert.NoError(t, err)
				return auth
			},
			msg:           toAddrTx,
			expectedError: ErrAuthorizationWrongChainID,
			expectedMockCalls: func(m *mock_vm.MockStateDB) {
				// nothing
			},
		},
		{
			name: "don't allow account key type: signer's key was updated",
			makeAuthorization: func() types.SetCodeAuthorization {
				auth, err := types.SignSetCode(authorityKey, types.SetCodeAuthorization{
					ChainID: *uint256.MustFromBig(params.TestChainConfig.ChainID),
					Address: aa,
					Nonce:   uint64(1),
				})
				assert.NoError(t, err)
				return auth
			},
			msg:           toAddrTx,
			expectedError: fmt.Errorf("%w: %v", ErrAuthorizationNotAllowAccountKeyType, accountkey.AccountKeyTypePublic),
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

			header := &types.Header{Number: big.NewInt(0), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
			blockContext := NewEVMBlockContext(header, nil, &common.Address{}) // stub author (COINBASE) to 0x0
			txContext := NewEVMTxContext(tt.msg, header, params.TestChainConfig)
			evm := vm.NewEVM(blockContext, txContext, mockStateDB, params.TestChainConfig, &vm.Config{Debug: true})

			// Verify that the expected mockStateDB's calls are being made.
			auth := tt.makeAuthorization()
			actual := NewStateTransition(evm, tt.msg).applyAuthorization(&auth, rules)
			require.Equal(t, tt.expectedError, actual)
		})
	}
}

func TestStateTransition_preCheck(t *testing.T) {
	var (
		addr              = common.HexToAddress("0x1234567890123456789012345678901234567890")
		to                = common.HexToAddress("0x000000000000000000000000000000000000aaaa")
		baseFee           = big.NewInt(10 * params.Gkei)
		blobBaseFee       = big.NewInt(10 * params.Gkei)
		nonce             = uint64(5)
		gasLimit          = uint64(100000)
		kaiaChainConfig   = params.TestKaiaConfig("kaia")
		pragueChainConfig = params.TestKaiaConfig("prague")
		osakaChainConfig  = params.TestKaiaConfig("osaka")
	)

	// Generate valid versioned hashes for testing
	validBlobHash := common.Hash{}
	validBlobHash[0] = 0x01 // Version 1 blob hash

	invalidBlobHash := common.Hash{}
	invalidBlobHash[0] = 0xFF // Invalid version

	// Test cases for preCheck
	tests := []struct {
		name               string
		config             *params.ChainConfig
		expectedError      error
		prefetching        bool
		setupMockMsg       func(ctrl *gomock.Controller) *mock_bc.MockMessage
		setupStateMockCall func(m *mock_vm.MockStateDB)
	}{
		// === Successful preCheck ===
		{
			name:          "successful preCheck with legacy tx",
			config:        kaiaChainConfig,
			expectedError: nil,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().ValidatedFeePayer().Return(addr).AnyTimes()
				mockMsg.EXPECT().FeeRatio().Return(types.FeeRatio(0), false).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().AuthList().Return(nil).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
				m.EXPECT().GetBalance(addr).Return(big.NewInt(params.KAIA))
				m.EXPECT().SubBalance(addr, gomock.Any())
			},
		},
		// === Prefetching mode ===
		{
			name:          "prefetching mode skips nonce and balance checks",
			config:        kaiaChainConfig,
			expectedError: nil,
			prefetching:   true,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				// No calls expected in prefetching mode
			},
		},
		// === Nonce checks ===
		{
			name:          "nonce too high",
			config:        kaiaChainConfig,
			expectedError: ErrNonceTooHigh,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce - 1) // state nonce < tx nonce
			},
		},
		{
			name:          "nonce too low",
			config:        kaiaChainConfig,
			expectedError: ErrNonceTooLow,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce + 1) // state nonce > tx nonce
			},
		},
		{
			name:          "nonce max overflow",
			config:        kaiaChainConfig,
			expectedError: ErrNonceMax,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(^uint64(0)).AnyTimes() // max uint64
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(^uint64(0))
			},
		},
		{
			name:          "check nonce disabled - nonce mismatch allowed",
			config:        kaiaChainConfig,
			expectedError: nil,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(false).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().ValidatedFeePayer().Return(addr).AnyTimes()
				mockMsg.EXPECT().FeeRatio().Return(types.FeeRatio(0), false).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().AuthList().Return(nil).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetBalance(addr).Return(big.NewInt(params.KAIA))
				m.EXPECT().SubBalance(addr, gomock.Any())
			},
		},
		// === Prague gas fee checks ===
		{
			name:          "Prague: gasFeeCap bit length too high",
			config:        pragueChainConfig,
			expectedError: ErrFeeCapVeryHigh,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				hugeGasFeeCap := new(big.Int).Lsh(big.NewInt(1), 257) // 2^257
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(hugeGasFeeCap).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Prague: gasTipCap bit length too high",
			config:        pragueChainConfig,
			expectedError: ErrTipVeryHigh,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				hugeGasTipCap := new(big.Int).Lsh(big.NewInt(1), 257) // 2^257
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(params.KAIA)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(hugeGasTipCap).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Prague: tip above fee cap",
			config:        pragueChainConfig,
			expectedError: ErrTipAboveFeeCap,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(10 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(20 * params.Gkei)).AnyTimes() // tip > feeCap
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Prague: fee cap too low (below baseFee)",
			config:        pragueChainConfig,
			expectedError: ErrFeeCapTooLow,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes() // below baseFee
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Prague: skip check when gasFeeCap and gasTipCap are both zero",
			config:        pragueChainConfig,
			expectedError: nil,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().ValidatedFeePayer().Return(addr).AnyTimes()
				mockMsg.EXPECT().FeeRatio().Return(types.FeeRatio(0), false).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(0)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().AuthList().Return(nil).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(0)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
				m.EXPECT().GetBalance(addr).Return(big.NewInt(params.KAIA))
				m.EXPECT().SubBalance(addr, gomock.Any())
			},
		},
		// === Blob validation ===
		{
			name:          "blob tx with nil to (create contract) returns error",
			config:        osakaChainConfig,
			expectedError: ErrBlobTxCreate,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return([]common.Hash{validBlobHash}).AnyTimes()
				mockMsg.EXPECT().To().Return(nil).AnyTimes() // nil to
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "blob tx with empty blob hashes",
			config:        osakaChainConfig,
			expectedError: ErrMissingBlobHashes,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return([]common.Hash{}).AnyTimes() // empty but not nil
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Osaka: too many blobs",
			config:        osakaChainConfig,
			expectedError: ErrTooManyBlobs,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				blobHashes := make([]common.Hash, params.BlobTxMaxBlobs+1)
				for i := range blobHashes {
					blobHashes[i] = validBlobHash
				}
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(blobHashes).AnyTimes()
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "invalid blob hash version",
			config:        osakaChainConfig,
			expectedError: errors.New("blob 0 has invalid hash version"),
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return([]common.Hash{invalidBlobHash}).AnyTimes()
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Osaka: blob fee cap too low",
			config:        osakaChainConfig,
			expectedError: ErrBlobFeeCapTooLow,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return([]common.Hash{validBlobHash}).AnyTimes()
				mockMsg.EXPECT().BlobGasFeeCap().Return(big.NewInt(1)).AnyTimes() // too low
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "Osaka: skip blob fee cap check when blobGasFeeCap is zero",
			config:        osakaChainConfig,
			expectedError: nil,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().ValidatedFeePayer().Return(addr).AnyTimes()
				mockMsg.EXPECT().FeeRatio().Return(types.FeeRatio(0), false).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return([]common.Hash{validBlobHash}).AnyTimes()
				mockMsg.EXPECT().BlobGasFeeCap().Return(big.NewInt(0)).AnyTimes() // zero - skip check
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().AuthList().Return(nil).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
				m.EXPECT().GetBalance(addr).Return(big.NewInt(params.KAIA))
				m.EXPECT().SubBalance(addr, gomock.Any())
			},
		},
		// === EIP-7702 Authorization List checks ===
		{
			name:          "EIP-7702: setcode tx cannot create contract (to is nil)",
			config:        pragueChainConfig,
			expectedError: ErrSetCodeTxCreate,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().AuthList().Return([]types.SetCodeAuthorization{{}}).AnyTimes() // non-empty
				mockMsg.EXPECT().To().Return(nil).AnyTimes()                                    // nil to
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
		{
			name:          "EIP-7702: empty auth list",
			config:        pragueChainConfig,
			expectedError: ErrEmptyAuthList,
			prefetching:   false,
			setupMockMsg: func(ctrl *gomock.Controller) *mock_bc.MockMessage {
				mockMsg := mock_bc.NewMockMessage(ctrl)
				mockMsg.EXPECT().CheckNonce().Return(true).AnyTimes()
				mockMsg.EXPECT().ValidatedSender().Return(addr).AnyTimes()
				mockMsg.EXPECT().Nonce().Return(nonce).AnyTimes()
				mockMsg.EXPECT().Hash().Return(common.Hash{}).AnyTimes()
				mockMsg.EXPECT().Gas().Return(gasLimit).AnyTimes()
				mockMsg.EXPECT().GasPrice().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasFeeCap().Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().GasTipCap().Return(big.NewInt(1 * params.Gkei)).AnyTimes()
				mockMsg.EXPECT().BlobHashes().Return(nil).AnyTimes()
				mockMsg.EXPECT().AuthList().Return([]types.SetCodeAuthorization{}).AnyTimes() // empty
				mockMsg.EXPECT().To().Return(&to).AnyTimes()
				mockMsg.EXPECT().Value().Return(big.NewInt(100)).AnyTimes()
				mockMsg.EXPECT().Data().Return(nil).AnyTimes()
				mockMsg.EXPECT().EffectiveGasPrice(gomock.Any(), gomock.Any()).Return(big.NewInt(25 * params.Gkei)).AnyTimes()
				return mockMsg
			},
			setupStateMockCall: func(m *mock_vm.MockStateDB) {
				m.EXPECT().GetNonce(addr).Return(nonce)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockStateDB := mock_vm.NewMockStateDB(mockCtrl)
			tt.setupStateMockCall(mockStateDB)

			mockMsg := tt.setupMockMsg(mockCtrl)

			fork.SetHardForkBlockNumberConfig(tt.config)

			header := &types.Header{
				Number:     big.NewInt(0),
				Time:       big.NewInt(0),
				BlockScore: big.NewInt(0),
				BaseFee:    baseFee,
			}
			// Create block context manually to avoid nil chain issue
			blockContext := vm.BlockContext{
				CanTransfer: CanTransfer,
				Transfer:    Transfer,
				GetHash:     func(n uint64) common.Hash { return common.Hash{} },
				Coinbase:    common.Address{},
				BlockNumber: new(big.Int).Set(header.Number),
				Time:        new(big.Int).Set(header.Time),
				BlockScore:  new(big.Int).Set(header.BlockScore),
				BaseFee:     baseFee,
				BlobBaseFee: blobBaseFee,
			}

			txContext := NewEVMTxContext(mockMsg, header, tt.config)
			evm := vm.NewEVM(blockContext, txContext, mockStateDB, tt.config, &vm.Config{Prefetching: tt.prefetching})

			st := NewStateTransition(evm, mockMsg)
			err := st.preCheck()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStateTransition_EIP7623(t *testing.T) {
	// Prague fork block at 10
	config := params.TestChainConfig.Copy()
	config.PragueCompatibleBlock = big.NewInt(10)
	// Apply chain config to fork
	fork.SetHardForkBlockNumberConfig(config)

	var (
		key, _    = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr      = crypto.PubkeyToAddress(key.PublicKey)
		amount    = big.NewInt(1000)
		data      = []byte{1, 2, 3, 4, 0, 0, 0, 0} // 4 non-zero bytes, 4 zero bytes
		signer    = types.LatestSigner(config)
		gaslimit1 = uint64(21800) // 21000 + 100 * 8 (100 per byte)
		gaslimit2 = uint64(21200) // 21000 + 10*4*4 + 10*4 (10 per token, 4 tokens per non-zero byte, 1 token per zero byte)
	)

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockStateDB := mock_vm.NewMockStateDB(mockCtrl)
	mockStateDB.EXPECT().GetBalance(gomock.Any()).Return(big.NewInt(params.KAIA)).AnyTimes()
	mockStateDB.EXPECT().SubBalance(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockStateDB.EXPECT().GetKey(gomock.Any()).Return(accountkey.NewAccountKeyLegacy()).AnyTimes()
	mockStateDB.EXPECT().GetNonce(gomock.Any()).Return(uint64(0)).AnyTimes()
	mockStateDB.EXPECT().Prepare(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockStateDB.EXPECT().IncNonce(gomock.Any()).Return().AnyTimes()
	mockStateDB.EXPECT().GetCode(gomock.Any()).Return(nil).AnyTimes()
	mockStateDB.EXPECT().Snapshot().Return(1).AnyTimes()
	mockStateDB.EXPECT().Exist(gomock.Any()).Return(false).AnyTimes()
	mockStateDB.EXPECT().GetRefund().Return(uint64(0)).AnyTimes()
	mockStateDB.EXPECT().AddBalance(gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockStateDB.EXPECT().CreateEOA(gomock.Any(), gomock.Any(), gomock.Any()).Return().AnyTimes()
	mockStateDB.EXPECT().GetVmVersion(gomock.Any()).Return(params.VmVersion0, false).AnyTimes()
	mockStateDB.EXPECT().IsProgramAccount(gomock.Any()).Return(false).AnyTimes()

	var (
		header       *types.Header
		blockContext vm.BlockContext
		txContext    vm.TxContext
		evm          *vm.EVM
		res          *ExecutionResult
		err          error
		tx           *types.Transaction
	)

	// Generate tx before Prague
	tx = types.NewTransaction(0, addr, amount, gaslimit1, big.NewInt(1), data)
	err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key})
	assert.NoError(t, err)
	tx, err = tx.AsMessageWithAccountKeyPicker(signer, mockStateDB, 0)
	assert.NoError(t, err)

	header = &types.Header{Number: big.NewInt(0), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
	blockContext = NewEVMBlockContext(header, nil, &common.Address{})
	txContext = NewEVMTxContext(tx, header, config)
	evm = vm.NewEVM(blockContext, txContext, mockStateDB, config, &vm.Config{})

	res, err = NewStateTransition(evm, tx).TransitionDb()
	assert.NoError(t, err)
	assert.Equal(t, gaslimit1, res.UsedGas)

	// Generate tx after Prague
	tx = types.NewTransaction(0, addr, amount, gaslimit2, big.NewInt(1), data)
	err = tx.SignWithKeys(signer, []*ecdsa.PrivateKey{key})
	assert.NoError(t, err)
	tx, err = tx.AsMessageWithAccountKeyPicker(signer, mockStateDB, 20)
	assert.NoError(t, err)

	header = &types.Header{Number: big.NewInt(20), Time: big.NewInt(0), BlockScore: big.NewInt(0)}
	blockContext = NewEVMBlockContext(header, nil, &common.Address{})
	txContext = NewEVMTxContext(tx, header, config)
	evm = vm.NewEVM(blockContext, txContext, mockStateDB, config, &vm.Config{})

	res, err = NewStateTransition(evm, tx).TransitionDb()
	assert.NoError(t, err)
	assert.Equal(t, gaslimit2, res.UsedGas)
}
