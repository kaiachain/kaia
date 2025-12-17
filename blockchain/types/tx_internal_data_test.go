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
	"crypto/sha256"
	"errors"
	"strings"
	"testing"

	"github.com/holiman/uint256"
	"github.com/kaiachain/kaia/blockchain/types/account"
	mock_types "github.com/kaiachain/kaia/blockchain/types/mocks"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/kzg4844"
	"github.com/kaiachain/kaia/kerrors"
	"github.com/kaiachain/kaia/params"
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

func TestSidecar_ValidateWithBlobHashes(t *testing.T) {
	// Helper function to create a valid sidecar with empty blob
	createValidSidecar := func(numBlobs int) (*BlobTxSidecar, []common.Hash) {
		blobs := make([]kzg4844.Blob, numBlobs)
		commitments := make([]kzg4844.Commitment, numBlobs)
		proofs := make([]kzg4844.Proof, 0)
		blobHashes := make([]common.Hash, numBlobs)
		hasher := sha256.New()

		for i := 0; i < numBlobs; i++ {
			blob := kzg4844.Blob{}
			commitment, _ := kzg4844.BlobToCommitment(&blob)
			cellProofs, _ := kzg4844.ComputeCellProofs(&blob)

			blobs[i] = blob
			commitments[i] = commitment
			proofs = append(proofs, cellProofs...)
			blobHashes[i] = kzg4844.CalcBlobHashV1(hasher, &commitment)
		}

		return &BlobTxSidecar{
			Version:     1,
			Blobs:       blobs,
			Commitments: commitments,
			Proofs:      proofs,
		}, blobHashes
	}

	var (
		singleSidecar, singleSidecarHashes     = createValidSidecar(1)
		multipleSidecar, multipleSidecarHashes = createValidSidecar(2)
		singleSidecarCommitEmpty               = singleSidecar.Copy()
		singleSidecarVersion0                  = singleSidecar.Copy()
		singleSidecarProofInvalid              = singleSidecar.Copy()
		wrongHash                              = []common.Hash{common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")}
	)
	singleSidecarCommitEmpty.Commitments = []kzg4844.Commitment{}
	singleSidecarVersion0.Version = 0
	singleSidecarProofInvalid.Proofs = singleSidecarProofInvalid.Proofs[:len(singleSidecarProofInvalid.Proofs)-1]

	tests := []struct {
		name       string
		sidecar    *BlobTxSidecar
		hashes     []common.Hash
		expected   error
		cacheCheck func(t *testing.T, sidecar *BlobTxSidecar, err error)
	}{
		{
			name:    "empty blob",
			sidecar: singleSidecar,
			hashes:  singleSidecarHashes,
			cacheCheck: func(t *testing.T, sidecar *BlobTxSidecar, err error) {
				require.NoError(t, err)
				blobHashesSummary := make([]byte, 0, 32)
				blobHashesSummary = append(blobHashesSummary, singleSidecarHashes[0][:]...)
				expectedSummaryHash := crypto.Keccak256Hash(blobHashesSummary)
				assert.Equal(t, sidecar.validatedSummaryHash, expectedSummaryHash)
			},
		},
		{
			name:    "multiple blobs",
			sidecar: multipleSidecar,
			hashes:  multipleSidecarHashes,
			cacheCheck: func(t *testing.T, sidecar *BlobTxSidecar, err error) {
				require.NoError(t, err)
				blobHashesSummary := make([]byte, 0, len(multipleSidecarHashes)*32)
				for _, hash := range multipleSidecarHashes {
					blobHashesSummary = append(blobHashesSummary, hash[:]...)
				}
				expectedSummaryHash := crypto.Keccak256Hash(blobHashesSummary)
				assert.Equal(t, sidecar.validatedSummaryHash, expectedSummaryHash)
			},
		},
		{
			name:     "mismatched blob count",
			sidecar:  singleSidecar,
			hashes:   multipleSidecarHashes,
			expected: errors.New("invalid number of"),
		},
		{
			name:     "mismatched commitment count",
			sidecar:  singleSidecarCommitEmpty,
			hashes:   singleSidecarHashes,
			expected: errors.New("invalid number of"),
		},
		{
			name:     "mismatched hash",
			sidecar:  singleSidecar,
			hashes:   wrongHash,
			expected: errors.New("mismatches transaction one"),
		},
		{
			name:     "unsupported version 0",
			sidecar:  singleSidecarVersion0,
			hashes:   singleSidecarHashes,
			expected: errors.New("version 0 not supported"),
		},
		{
			name:     "invalid proof count",
			sidecar:  singleSidecarProofInvalid,
			hashes:   singleSidecarHashes,
			expected: errors.New("invalid number of"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.sidecar.ValidateWithBlobHashes(tt.hashes)

			if tt.expected != nil {
				require.Error(t, err)
				require.True(t, strings.Contains(err.Error(), tt.expected.Error()))
			} else {
				require.NoError(t, err)
			}

			if tt.cacheCheck != nil {
				tt.cacheCheck(t, tt.sidecar, err)
			}
		})
	}
}

func BenchmarkSidecar_ValidateWithBlobHashes(b *testing.B) {
	emptyBlob := kzg4844.Blob{}
	emptyBlobCommit, _ := kzg4844.BlobToCommitment(&emptyBlob)
	cellProofs, _ := kzg4844.ComputeCellProofs(&emptyBlob)
	blobHash := kzg4844.CalcBlobHashV1(sha256.New(), &emptyBlobCommit)
	sidecar := &BlobTxSidecar{
		Version:     1,
		Blobs:       []kzg4844.Blob{emptyBlob},
		Commitments: []kzg4844.Commitment{emptyBlobCommit},
		Proofs:      cellProofs,
	}

	for i := 0; i < b.N; i++ {
		sidecar.ValidateWithBlobHashes([]common.Hash{blobHash})
	}
}
