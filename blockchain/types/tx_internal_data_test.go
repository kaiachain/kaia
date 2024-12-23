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

	"github.com/kaiachain/kaia/blockchain/types/account"
	mock_types "github.com/kaiachain/kaia/blockchain/types/mocks"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kerrors"
	"go.uber.org/mock/gomock"
	"gotest.tools/assert"
)

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
