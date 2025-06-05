// Copyright 2024 The Kaia Authors
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

package types

import (
	"testing"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/v2/blockchain/types/account"
	mock_types "github.com/kaiachain/kaia/v2/blockchain/types/mocks"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/kerrors"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

func TestPrefixedRlpHash7702(t *testing.T) {
	var (
		config  = *params.TestChainConfig
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
	)
	// auth1
	auth1 := SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(config.ChainID),
		Address: addr1,
		Nonce:   uint64(1),
	}
	rlpHash1 := prefixedRlpHash(0x05, []any{
		auth1.ChainID,
		auth1.Address,
		auth1.Nonce,
	})
	// auth2
	auth2 := SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(config.ChainID),
		Address: addr2,
		Nonce:   uint64(1),
	}
	rlpHash2 := prefixedRlpHash(0x05, []any{
		auth2.ChainID,
		auth2.Address,
		auth2.Nonce,
	})
	require.NotEqual(t, rlpHash1, rlpHash2)
}

func TestSignSetCodeAndAuthority7702(t *testing.T) {
	var (
		config  = *params.TestChainConfig
		key1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		key2, _ = crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
		addr1   = crypto.PubkeyToAddress(key1.PublicKey)
		addr2   = crypto.PubkeyToAddress(key2.PublicKey)
	)

	// auth1
	auth1 := SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(config.ChainID),
		Address: addr1,
		Nonce:   uint64(1),
	}
	// auth1: SignSetCode and Authority
	actualAuth1, err := SignSetCode(key1, auth1)
	require.NoError(t, err)
	require.Equal(t, auth1.ChainID, actualAuth1.ChainID)
	require.Equal(t, auth1.Address, actualAuth1.Address)
	require.Equal(t, auth1.Nonce, actualAuth1.Nonce)
	actualAuthority1, err := actualAuth1.Authority()
	require.NoError(t, err)
	require.Equal(t, addr1, actualAuthority1)

	// auth2
	auth2 := SetCodeAuthorization{
		ChainID: *uint256.MustFromBig(config.ChainID),
		Address: addr2,
		Nonce:   uint64(1),
	}
	// auth2: SignSetCode and Authority
	actualAuth2, err := SignSetCode(key2, auth2)
	require.NoError(t, err)
	require.Equal(t, auth2.ChainID, actualAuth2.ChainID)
	require.Equal(t, auth2.Address, actualAuth2.Address)
	require.Equal(t, auth2.Nonce, actualAuth2.Nonce)
	actualAuthority2, err := actualAuth2.Authority()
	require.NoError(t, err)
	require.Equal(t, addr2, actualAuthority2)

	// set addr2 to auth1
	{
		actualAuth1.Address = addr2
		actualAuthority, err := actualAuth1.Authority()
		require.NoError(t, err)
		require.NotEqual(t, addr1, actualAuthority)
		require.NotEqual(t, addr2, actualAuthority)
	}
}

// Coverage is detected by calling functions that modify state.
func TestValidate7702(t *testing.T) {
	// prepare target txType
	typesToMustBeEoaWithoutCode := []TxType{
		TxTypeValueTransfer, TxTypeFeeDelegatedValueTransfer, TxTypeFeeDelegatedValueTransferWithRatio,
		TxTypeValueTransferMemo, TxTypeFeeDelegatedValueTransferMemo, TxTypeFeeDelegatedValueTransferMemoWithRatio,
	}
	typesFromMustBeEoaWithoutCode := []TxType{
		TxTypeAccountUpdate, TxTypeFeeDelegatedAccountUpdate, TxTypeFeeDelegatedAccountUpdateWithRatio,
	}
	typesToMustBeEOAWithCodeOrSCA := []TxType{
		TxTypeSmartContractExecution, TxTypeFeeDelegatedSmartContractExecution, TxTypeFeeDelegatedSmartContractExecutionWithRatio,
	}

	// prepare account
	nonEmptyCodeHash := common.HexToHash("aaaaaaaabbbbbbbbccccccccddddddddaaaaaaaabbbbbbbbccccccccdddddddd").Bytes()
	eoaWithoutCode, _ := account.NewAccountWithType(account.ExternallyOwnedAccountType)
	eoa, _ := account.NewAccountWithType(account.ExternallyOwnedAccountType)
	eoaWithCode := account.GetProgramAccount(eoa)
	if eoaWithCode != nil {
		eoaWithCode.SetCodeHash(nonEmptyCodeHash)
	}
	sca, _ := account.NewAccountWithType(account.SmartContractAccountType)

	// Because we are testing with mocks, the consistency of the key and address is not important.
	// Therefore, we assign simple addresses.
	eoaWithoutCodeAddress := common.HexToAddress("0x000000000000000000000000000000000000aaaa")
	eoaWithCodeAddress := common.HexToAddress("0x000000000000000000000000000000000000bbbb")
	scaAddress := common.HexToAddress("0x000000000000000000000000000000000000cccc")

	type Args struct {
		targetTypes []TxType
		from        common.Address
		to          common.Address
	}

	tests := []struct {
		name              string
		args              Args
		expectedErr       error
		expectedMockCalls func(m *mock_types.MockStateDB)
	}{
		// Group1: type of value transfer
		{
			name: "valid value transfer (to account not exist)",
			args: Args{
				targetTypes: typesToMustBeEoaWithoutCode,
				from:        common.Address{},
				to:          eoaWithoutCodeAddress,
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(nil)
			},
		},
		{
			name: "valid value transfer (to account is EOA without a code)",
			args: Args{
				targetTypes: typesToMustBeEoaWithoutCode,
				from:        common.Address{},
				to:          eoaWithoutCodeAddress,
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(eoaWithoutCode)
			},
		},
		{
			name: "invalid value transfer (to account is EOA with a code)",
			args: Args{
				targetTypes: typesToMustBeEoaWithoutCode,
				from:        common.Address{},
				to:          eoaWithCodeAddress,
			},
			expectedErr: kerrors.ErrToMustBeEOAWithoutCode,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithCodeAddress).Return(eoaWithCode)
			},
		},
		{
			name: "invalid value transfer (to account is SCA)",
			args: Args{
				targetTypes: typesToMustBeEoaWithoutCode,
				from:        common.Address{},
				to:          scaAddress,
			},
			expectedErr: kerrors.ErrToMustBeEOAWithoutCode,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(scaAddress).Return(sca)
			},
		},

		// Group2: type of account update
		{
			name: "valid account update (from account not exist)",
			args: Args{
				targetTypes: typesFromMustBeEoaWithoutCode,
				from:        eoaWithoutCodeAddress,
				to:          eoaWithoutCodeAddress,
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(nil)
			},
		},
		{
			name: "valid account update (from account is EOA without a code)",
			args: Args{
				targetTypes: typesFromMustBeEoaWithoutCode,
				from:        eoaWithoutCodeAddress,
				to:          common.Address{},
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(eoaWithoutCode)
			},
		},
		{
			name: "invalid account update (from account is EOA with a code)",
			args: Args{
				targetTypes: typesFromMustBeEoaWithoutCode,
				from:        eoaWithCodeAddress,
				to:          common.Address{},
			},
			expectedErr: kerrors.ErrFromMustBeEOAWithoutCode,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithCodeAddress).Return(eoaWithCode)
			},
		},
		{
			name: "invalid account update (from account is SCA)",
			args: Args{
				targetTypes: typesFromMustBeEoaWithoutCode,
				from:        scaAddress,
				to:          common.Address{},
			},
			expectedErr: kerrors.ErrFromMustBeEOAWithoutCode,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(scaAddress).Return(sca)
			},
		},

		// Group3: type of smart contract execution
		{
			name: "valid smart contract execution (to account is SCA)",
			args: Args{
				targetTypes: typesToMustBeEOAWithCodeOrSCA,
				from:        common.Address{},
				to:          scaAddress,
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(scaAddress).Return(sca)
			},
		},
		{
			name: "valid smart contract execution (to account is EOA with a code)",
			args: Args{
				targetTypes: typesToMustBeEOAWithCodeOrSCA,
				from:        common.Address{},
				to:          eoaWithCodeAddress,
			},
			expectedErr: nil,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithCodeAddress).Return(eoaWithCode)
			},
		},
		{
			name: "invalid smart contract execution (to account not exist)",
			args: Args{
				targetTypes: typesToMustBeEOAWithCodeOrSCA,
				from:        common.Address{},
				to:          eoaWithoutCodeAddress,
			},
			expectedErr: kerrors.ErrToMustBeEOAWithCodeOrSCA,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(nil)
			},
		},
		{
			name: "invalid smart contract execution (to account is EOA without a code)",
			args: Args{
				targetTypes: typesToMustBeEOAWithCodeOrSCA,
				from:        common.Address{},
				to:          eoaWithoutCodeAddress,
			},
			expectedErr: kerrors.ErrToMustBeEOAWithCodeOrSCA,
			expectedMockCalls: func(m *mock_types.MockStateDB) {
				m.EXPECT().GetAccount(eoaWithoutCodeAddress).Return(eoaWithoutCode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, txType := range tt.args.targetTypes {
				mockCtrl := gomock.NewController(t)
				defer mockCtrl.Finish()

				mockStateDB := mock_types.NewMockStateDB(mockCtrl)
				tt.expectedMockCalls(mockStateDB)

				err := validate7702(mockStateDB, txType, tt.args.from, tt.args.to)
				assert.Equal(t, tt.expectedErr, err)
			}
		})
	}
}
