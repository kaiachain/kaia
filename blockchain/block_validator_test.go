// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/block_validator_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	mock_gov "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/valset"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type verifySealsTestSealer struct {
	consensus.Sealer
	author        common.Address
	committers    []common.Address
	quorum        int
	f             int
	qualifiedLen  int
	committeeSize int
}

func (s *verifySealsTestSealer) Author(*types.Header) (common.Address, error) {
	return s.author, nil
}

func (s *verifySealsTestSealer) Committers(*types.Header) ([]common.Address, error) {
	return s.committers, nil
}

func (s *verifySealsTestSealer) Round(*types.Header) (byte, error) {
	return 0, nil
}

func (s *verifySealsTestSealer) Quorum(_ uint64, qualifiedLen, committeeSize int) int {
	s.qualifiedLen = qualifiedLen
	s.committeeSize = committeeSize
	return s.quorum
}

func (s *verifySealsTestSealer) F(_ uint64, _, _ int) int {
	return s.f
}

func TestCountValidCommittedSeals(t *testing.T) {
	var (
		a = common.HexToAddress("0x0001")
		b = common.HexToAddress("0x0002")
		c = common.HexToAddress("0x0003")
		d = common.HexToAddress("0x0004")
	)

	testcases := []struct {
		name          string
		signers       []common.Address
		committers    []common.Address
		expectedCount int
		expectedErr   error
	}{
		{
			name:          "counts unique committed seals",
			signers:       []common.Address{a, b, c},
			committers:    []common.Address{a, b},
			expectedCount: 2,
		},
		{
			name:        "rejects duplicate committer",
			signers:     []common.Address{a, b, c},
			committers:  []common.Address{a, a, b},
			expectedErr: istanbul.ErrInvalidCommittedSeals,
		},
		{
			name:        "rejects non-signer committer",
			signers:     []common.Address{a, b, c},
			committers:  []common.Address{a, d},
			expectedErr: istanbul.ErrInvalidCommittedSeals,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := countValidCommittedSeals(tc.committers, valset.NewAddressSet(tc.signers))

			if tc.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tc.expectedErr)
			}
			assert.Equal(t, tc.expectedCount, count)
		})
	}
}

func TestValidateHeader(t *testing.T) {
	testcases := []struct {
		name           string
		maxHardfork    string
		mutateHeader   func(header *types.Header, parent *types.Header)
		expectedErr    error
		expectedErrStr string
	}{
		{
			name:        "invalid block score",
			maxHardfork: "kore",
			mutateHeader: func(header *types.Header, _ *types.Header) {
				header.BlockScore = big.NewInt(2)
			},
			expectedErr: ErrInvalidBlockScore,
		},
		{
			name:        "future block",
			maxHardfork: "kore",
			mutateHeader: func(header *types.Header, _ *types.Header) {
				header.Time = new(big.Int).Add(big.NewInt(time.Now().Unix()), new(big.Int).SetUint64(10))
			},
			expectedErr: consensus.ErrFutureBlock,
		},
		{
			name:        "invalid timestamp",
			maxHardfork: "kore",
			mutateHeader: func(header *types.Header, parent *types.Header) {
				header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(uint64(params.BlockGenerationInterval)-1))
			},
			expectedErr: ErrInvalidTimestamp,
		},
		{
			name:        "eip4844 header before osaka - excess blob gas",
			maxHardfork: "kore",
			mutateHeader: func(header *types.Header, _ *types.Header) {
				excessBlobGas := uint64(0)
				header.ExcessBlobGas = &excessBlobGas
			},
			expectedErr: ErrUnexpectedExcessBlobGasBeforeOsaka,
		},
		{
			name:        "eip4844 header before osaka - blob gas used",
			maxHardfork: "kore",
			mutateHeader: func(header *types.Header, _ *types.Header) {
				blobGasUsed := uint64(0)
				header.BlobGasUsed = &blobGasUsed
			},
			expectedErr: ErrUnexpectedBlobGasUsedBeforeOsaka,
		},
		{
			name:        "invalid eip4844 header",
			maxHardfork: "osaka",
			mutateHeader: func(header *types.Header, _ *types.Header) {
				header.ExcessBlobGas = nil
			},
			expectedErrStr: "header is missing excessBlobGas",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			config := params.TestKaiaConfig(tc.maxHardfork).Copy()
			var (
				testdb  = database.NewMemoryDBManager()
				gspec   = &Genesis{Config: config}
				genesis = gspec.MustCommit(testdb)
				sealer  = faker.NewFaker()
			)
			blocks, _ := GenerateChain(config, genesis, sealer, testdb, 1, nil)

			chain, _ := NewBlockChain(testdb, nil, config, sealer, vm.Config{})
			t.Cleanup(chain.Stop)

			header := types.CopyHeader(blocks[0].Header())
			parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			mGov := mock_gov.NewMockGovModule(ctrl)
			pset := gov.ParamSet{}
			if config != nil && config.Governance != nil && config.Governance.KIP71 != nil {
				pset = gov.ParamSet{
					LowerBoundBaseFee:         config.Governance.KIP71.LowerBoundBaseFee,
					UpperBoundBaseFee:         config.Governance.KIP71.UpperBoundBaseFee,
					GasTarget:                 config.Governance.KIP71.GasTarget,
					MaxBlockGasUsedForBaseFee: config.Governance.KIP71.MaxBlockGasUsedForBaseFee,
					BaseFeeDenominator:        config.Governance.KIP71.BaseFeeDenominator,
				}
			}
			mGov.EXPECT().GetParamSet(gomock.Any()).Return(pset).AnyTimes()

			validator := &BlockValidator{
				config: chain.Config(),
				hc:     chain,
				mGov:   mGov,
			}

			tc.mutateHeader(header, parent)

			err := validator.ValidateHeader(header)
			if tc.expectedErrStr != "" {
				assert.EqualError(t, err, tc.expectedErrStr)
				return
			}
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

// TestVerifySealsPermissionlessRejectsNonCommitteeAuthor checks that an author outside the
// round's committee is rejected even when every committed seal is valid.
func TestVerifySealsPermissionlessRejectsNonCommitteeAuthor(t *testing.T) {
	var (
		author    = common.HexToAddress("0x0001")
		b         = common.HexToAddress("0x0002")
		c         = common.HexToAddress("0x0003")
		blockNum  = uint64(7)
		header    = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		committee = []common.Address{b, c}
		sealer    = &verifySealsTestSealer{
			Sealer:     faker.NewFaker(),
			author:     author,
			committers: []common.Address{b, c},
		}
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return(committee, nil)

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.ErrorIs(t, validator.verifySeals(header, true), consensus.ErrUnauthorized)
}

func TestVerifySealsBeforePermissionlessUsesQualifiedSet(t *testing.T) {
	var (
		author   = common.HexToAddress("0x0001")
		blockNum = uint64(7)
		header   = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		sealer   = faker.NewFakerWithFixedSealer(author)
		config   = params.TestKaiaConfig("permissionless").Copy()
	)
	config.PermissionlessCompatibleBlock = big.NewInt(10)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetQualifiedValidators(blockNum).Return([]common.Address{author}, nil)
	mValset.EXPECT().GetCouncil(blockNum).Return([]common.Address{author}, nil)
	mGov.EXPECT().GetParamSet(blockNum).Return(gov.ParamSet{CommitteeSize: 1})

	validator := &BlockValidator{
		config:  config,
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.NoError(t, validator.verifySeals(header, true))
}

func TestVerifySealsPermissionlessAcceptsCommitteeAuthor(t *testing.T) {
	var (
		author   = common.HexToAddress("0x0001")
		blockNum = uint64(7)
		header   = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		sealer   = faker.NewFakerWithFixedSealer(author)
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return([]common.Address{author}, nil)

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.NoError(t, validator.verifySeals(header, true))
}

// TestVerifySealsPermissionlessUsesSealerQuorum checks that post-permissionless, the
// quorum is taken from sealer.Quorum with the committee size passed for both
// arguments (committee == qualified), and validSeal must meet it.
func TestVerifySealsPermissionlessUsesSealerQuorum(t *testing.T) {
	var (
		author    = common.HexToAddress("0x0001")
		b         = common.HexToAddress("0x0002")
		c         = common.HexToAddress("0x0003")
		d         = common.HexToAddress("0x0004")
		e         = common.HexToAddress("0x0005")
		blockNum  = uint64(7)
		header    = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		committee = []common.Address{author, b, c, d, e}
		sealer    = &verifySealsTestSealer{
			Sealer:     faker.NewFaker(),
			author:     author,
			committers: []common.Address{author, b, c, d},
			quorum:     4,
		}
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return(committee, nil)

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.NoError(t, validator.verifySeals(header, true))
	// Quorum is computed over the committee size, passed as both arguments.
	assert.Equal(t, len(committee), sealer.qualifiedLen)
	assert.Equal(t, len(committee), sealer.committeeSize)
}

// TestVerifySealsPermissionlessRejectsBelowQuorum checks that fewer committed seals
// than sealer.Quorum reports are rejected.
func TestVerifySealsPermissionlessRejectsBelowQuorum(t *testing.T) {
	var (
		author    = common.HexToAddress("0x0001")
		b         = common.HexToAddress("0x0002")
		c         = common.HexToAddress("0x0003")
		d         = common.HexToAddress("0x0004")
		blockNum  = uint64(7)
		header    = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		committee = []common.Address{author, b, c, d}
		sealer    = &verifySealsTestSealer{
			Sealer:     faker.NewFaker(),
			author:     author,
			committers: []common.Address{author, b, c},
			quorum:     4,
		}
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return(committee, nil)

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.ErrorIs(t, validator.verifySeals(header, true), istanbul.ErrInvalidCommittedSeals)
}

// TestVerifySealsPermissionlessRejectsNonCommitteeCommitter checks that a committed
// seal from a validator outside the round's committee is rejected, even if it
// belongs to the broader qualified/council set.
func TestVerifySealsPermissionlessRejectsNonCommitteeCommitter(t *testing.T) {
	var (
		author       = common.HexToAddress("0x0001")
		b            = common.HexToAddress("0x0002")
		c            = common.HexToAddress("0x0003")
		d            = common.HexToAddress("0x0004")
		nonCommittee = common.HexToAddress("0x0005")
		blockNum     = uint64(7)
		header       = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		committee    = []common.Address{author, b, c, d}
		sealer       = &verifySealsTestSealer{
			Sealer:     faker.NewFaker(),
			author:     author,
			committers: []common.Address{author, b, nonCommittee},
		}
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return(committee, nil)

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.ErrorIs(t, validator.verifySeals(header, true), istanbul.ErrInvalidCommittedSeals)
}

// roundBoundSealer stands in for the post-permissionless dynamicSealer: it always recovers
// with the round-bound preimage. The per-fork selection itself is covered by consensus/engine,
// which this package cannot import (engine -> istanbul/backend -> blockchain).
type roundBoundSealer struct {
	*istanbul.IstanbulSealer
}

func (s *roundBoundSealer) Committers(header *types.Header) ([]common.Address, error) {
	return s.IstanbulSealer.CommittersWithRound(header)
}

// TestVerifySealsPermissionlessRejectsMutatedRound runs real committer recovery through
// verifySeals: the round byte lives in the vanity that HeaderHash zeroes, so mutating it
// leaves the block hash and the author seal intact and only the committed-seal preimage
// changes. Post-permissionless that must cost the block its committers.
func TestVerifySealsPermissionlessRejectsMutatedRound(t *testing.T) {
	const round = int64(3)
	var (
		blockNum = uint64(7)
		keys     = make([]*ecdsa.PrivateKey, 4) // Quorum(4, 4) == 3, so 4 valid seals pass
		addrs    = make([]common.Address, 4)
	)
	for i := range keys {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i], addrs[i] = key, crypto.PubkeyToAddress(key.PublicKey)
	}

	config := params.TestKaiaConfig("permissionless")
	sealer := &roundBoundSealer{istanbul.NewSealerImpl(keys[0])} // keys[0] proposes

	header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}
	require.NoError(t, sealer.WriteValidators(header, addrs))
	sealer.WriteRound(header, round)
	authorSeal, err := sealer.MakeAuthorSeal(header)
	require.NoError(t, err)
	require.NoError(t, sealer.WriteAuthorSeal(header, authorSeal))

	seals := make([][]byte, len(keys))
	for i, key := range keys {
		preimage := istanbul.PrepareCommittedSealWithRound(sealer.HeaderHash(header), byte(round))
		seals[i], err = crypto.Sign(crypto.Keccak256(preimage), key)
		require.NoError(t, err)
	}
	require.NoError(t, sealer.WriteCommittedSeals(header, seals))

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, gomock.Any()).Return(addrs, nil).AnyTimes()

	validator := &BlockValidator{config: config, sealer: sealer, mGov: mGov, mValset: mValset}

	require.NoError(t, validator.verifySeals(header, true))

	sealer.WriteRound(header, round+1)
	assert.ErrorIs(t, validator.verifySeals(header, true), istanbul.ErrInvalidCommittedSeals)
}

// TestVerifySealsPermissionlessAcceptsHashLockedAuthor covers the block a round change
// produces: the locked proposal is re-proposed in a later round, so the header carries the
// final round with the original author's seal. Such a block is valid, and deciding it must
// not depend on who was scheduled to propose the final round.
func TestVerifySealsPermissionlessAcceptsHashLockedAuthor(t *testing.T) {
	const round = int64(3)
	var (
		blockNum = uint64(7)
		keys     = make([]*ecdsa.PrivateKey, 4) // Quorum(4, 4) == 3, so 4 valid seals pass
		addrs    = make([]common.Address, 4)
	)
	for i := range keys {
		key, err := crypto.GenerateKey()
		require.NoError(t, err)
		keys[i], addrs[i] = key, crypto.PubkeyToAddress(key.PublicKey)
	}

	sealer := &roundBoundSealer{istanbul.NewSealerImpl(keys[0])} // keys[0] sealed an earlier round

	header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}
	require.NoError(t, sealer.WriteValidators(header, addrs))
	sealer.WriteRound(header, round)
	authorSeal, err := sealer.MakeAuthorSeal(header)
	require.NoError(t, err)
	require.NoError(t, sealer.WriteAuthorSeal(header, authorSeal))

	seals := make([][]byte, len(keys))
	for i, key := range keys {
		preimage := istanbul.PrepareCommittedSealWithRound(sealer.HeaderHash(header), byte(round))
		seals[i], err = crypto.Sign(crypto.Keccak256(preimage), key)
		require.NoError(t, err)
	}
	require.NoError(t, sealer.WriteCommittedSeals(header, seals))

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(round)).Return(addrs, nil)
	// The final round is another validator's turn, and the block must be valid regardless.
	mValset.EXPECT().GetProposer(blockNum, uint64(round)).Return(addrs[1], nil).AnyTimes()

	validator := &BlockValidator{
		config:  params.TestKaiaConfig("permissionless"),
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	require.NoError(t, validator.verifySeals(header, true))
}

func TestVerifyBlockBody(t *testing.T) {
	testcases := []struct {
		baseFee *big.Int
		txData  types.TxInternalData
		err     bool
	}{
		{
			big.NewInt(100),
			&types.TxInternalDataLegacy{
				GasLimit: 1,
				Price:    big.NewInt(200),
				Payload:  []byte("abcdef"),
			},
			false,
		},

		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(10),
				GasTipCap: big.NewInt(10),
				Payload:   []byte("abcdef"),
			},
			true,
		},
		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(10),
				Payload:   []byte("abcdef"),
			},
			false,
		},
		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(200),
				GasTipCap: big.NewInt(10),
			},
			false,
		},
	}

	// Create a simple chain to verify
	var (
		testdb  = database.NewMemoryDBManager()
		gspec   = &Genesis{Config: params.TestChainConfig}
		genesis = gspec.MustCommit(testdb)
	)

	// We don't need istanbul instance here because validateBody is in BlockValidator instance
	GenerateChain(params.TestChainConfig, genesis, faker.NewFaker(), testdb, 8, nil)
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, faker.NewFaker(), vm.Config{})
	defer chain.Stop()

	// Generate a batch of accounts to start with
	privKey, _ := crypto.GenerateKey()
	signer := types.LatestSignerForChainID(big.NewInt(1))
	var block *types.Block
	for _, testcase := range testcases {
		// Generate a block header
		header := &types.Header{
			ParentHash: chain.hc.currentHeaderHash,
			Number:     common.Big1,
			GasUsed:    0,
			Extra:      []byte{},
			BaseFee:    testcase.baseFee,
		}

		// Generate a block with tx
		tx := types.NewTx(testcase.txData)
		tx.Sign(signer, privKey)
		block = types.NewBlock(header, append(types.Transactions{}, tx), nil)

		err := chain.validator.ValidateBody(block)
		if errExist := err != nil; errExist != testcase.err {
			assert.Error(t, err)
		}
	}
}

// TestVerifyBlockBodyForOsakaFork test rlp size limit for Osaka fork
// See: tests/osaka/eip7934_block_rlp_limit/test_max_block_rlp_size.py::test_block_at_rlp_size_limit_boundary
func TestVerifyBlockBodyForOsakaFork(t *testing.T) {
	testcases := []struct {
		txData types.TxInternalData
		err    bool
	}{
		{
			&types.TxInternalDataLegacy{
				GasLimit: 1,
				Price:    big.NewInt(200),
				Payload:  make([]byte, params.MaxBlockSize-1), // fork_Osaka-blockchain_test-max_rlp_size_minus_1_byte
			},
			false,
		},
		{
			&types.TxInternalDataLegacy{
				GasLimit: 1,
				Price:    big.NewInt(200),
				Payload:  make([]byte, params.MaxBlockSize), // fork_Osaka-blockchain_test-max_rlp_size
			},
			false,
		},
		{
			&types.TxInternalDataLegacy{
				GasLimit: 1,
				Price:    big.NewInt(200),
				Payload:  make([]byte, params.MaxBlockSize+1), // fork_Osaka-blockchain_test-max_rlp_size_plus_1_byte
			},
			true,
		},
	}

	// Create a simple chain to verify
	config := params.TestChainConfig.Copy()
	config.SetDefaults()
	config.IstanbulCompatibleBlock = common.Big1
	config.LondonCompatibleBlock = common.Big1
	config.EthTxTypeCompatibleBlock = common.Big1
	config.MagmaCompatibleBlock = common.Big1
	config.KoreCompatibleBlock = common.Big1
	config.ShanghaiCompatibleBlock = common.Big1
	config.CancunCompatibleBlock = common.Big1
	config.RandaoCompatibleBlock = common.Big1
	config.KaiaCompatibleBlock = common.Big1
	config.PragueCompatibleBlock = common.Big1
	config.OsakaCompatibleBlock = common.Big1
	var (
		testdb  = database.NewMemoryDBManager()
		gspec   = &Genesis{Config: config}
		genesis = gspec.MustCommit(testdb)
	)
	GenerateChain(config, genesis, faker.NewFaker(), testdb, 8, nil)
	chain, _ := NewBlockChain(testdb, nil, config, faker.NewFaker(), vm.Config{})
	defer chain.Stop()

	// Generate a batch of accounts to start with
	privKey, _ := crypto.GenerateKey()
	signer := types.LatestSignerForChainID(big.NewInt(1))
	var block *types.Block
	for _, testcase := range testcases {
		// Generate a block header
		header := &types.Header{
			ParentHash: chain.hc.currentHeaderHash,
			Number:     common.Big1,
			GasUsed:    0,
			Extra:      []byte{},
			BaseFee:    big.NewInt(100),
		}

		// Generate a block with tx
		tx := types.NewTx(testcase.txData)
		tx.Sign(signer, privKey)
		block = types.NewBlock(header, append(types.Transactions{}, tx), nil)

		err := chain.validator.ValidateBody(block)
		if errExist := err != nil; errExist != testcase.err {
			assert.Error(t, err)
		}
	}
}

// Tests that header preprocess verification works for both valid and invalid chains.
func TestPreprocessHeaderVerification(t *testing.T) {
	for _, threads := range []int{2, 8, 32} {
		t.Run(fmt.Sprintf("threads_%d", threads), func(t *testing.T) {
			testPreprocessHeaderVerification(t, threads)
		})
	}
}

func testPreprocessHeaderVerification(t *testing.T, threads int) {
	// Create a simple chain to preprocess.
	var (
		testConfig = params.TestChainConfig.Copy()
		testdb     = database.NewMemoryDBManager()
		gspec      = &Genesis{Config: testConfig}
		genesis    = gspec.MustCommit(testdb)
	)
	// Keep this test in pre-Osaka behavior for legacy faker headers.
	testConfig.OsakaCompatibleBlock = nil
	blocks, _ := GenerateChain(testConfig, genesis, faker.NewFaker(), testdb, 8, nil)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}

	// Match the old concurrent test shape, even though preprocess is now validator-owned.
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

	failBlock := headers[len(headers)-2].Number.Uint64()
	for i, valid := range []bool{true, false} {
		var abort chan<- struct{}
		var results <-chan error
		if valid {
			validator := &BlockValidator{sealer: faker.NewFaker()}
			abort, results = validator.Preprocess(headers)
		} else {
			validator := &BlockValidator{sealer: faker.NewFakeFailer(failBlock)}
			abort, results = validator.Preprocess(headers)
		}
		for j := 0; j < len(headers); j++ {
			select {
			case result := <-results:
				if valid || j < len(headers)-2 {
					if result != nil {
						t.Errorf("test %d.%d: expected nil result, got %v", i, j, result)
					}
				} else {
					assert.ErrorIs(t, result, consensus.ErrUnknownAncestor, "test %d.%d: expected ErrUnknownAncestor", i, j)
				}
			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: preprocessing timeout", i, j)
			}
		}
		select {
		case result := <-results:
			t.Fatalf("test %d: unexpected result returned: %v", i, result)
		case <-time.After(25 * time.Millisecond):
		}

		close(abort)
	}
}

func TestValidateHeaderRejectsOutOfUint64Number(t *testing.T) {
	v := &BlockValidator{}
	h := &types.Header{
		Number:     new(big.Int).Lsh(big.NewInt(1), 64), // 2^64, exceeds uint64
		Time:       big.NewInt(1),
		BlockScore: params.DefaultBlockScore,
	}
	assert.ErrorIs(t, v.ValidateHeader(h), ErrInvalidBlockNumber)
}

type roundShiftSealer struct {
	*verifySealsTestSealer
	round byte
}

func (s *roundShiftSealer) Round(*types.Header) (byte, error) { return s.round, nil }

func permissionlessConfig(t *testing.T, on bool, blockNum uint64) *params.ChainConfig {
	t.Helper()
	config := params.TestKaiaConfig("permissionless").Copy()
	if !on {
		config.PermissionlessCompatibleBlock = big.NewInt(int64(blockNum + 100))
	}
	return config
}

// A proposal has no committed seals, but its author must still be authorized.
func TestProposalHeaderRejectsUnauthorizedAuthor(t *testing.T) {
	var (
		outsider = common.HexToAddress("0x00ff")
		a        = common.HexToAddress("0x0001")
		b        = common.HexToAddress("0x0002")
		blockNum = uint64(7)
	)

	for _, permissionless := range []bool{false, true} {
		name := "pre-permissionless"
		if permissionless {
			name = "post-permissionless"
		}
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			sealer := &verifySealsTestSealer{author: outsider, committers: nil}
			header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}

			mGov := mock_gov.NewMockGovModule(ctrl)
			mValset := mock_valset.NewMockValsetModule(ctrl)
			if permissionless {
				mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return([]common.Address{a, b}, nil)
			} else {
				mValset.EXPECT().GetQualifiedValidators(blockNum).Return([]common.Address{a, b}, nil)
			}

			v := &BlockValidator{
				config: permissionlessConfig(t, permissionless, blockNum),
				sealer: sealer, mGov: mGov, mValset: mValset,
			}
			assert.ErrorIs(t, v.verifySeals(header, false), consensus.ErrUnauthorized)
		})
	}
}

// On the proposal path, absent committed seals are expected and must not raise an error.
func TestProposalHeaderAcceptsAuthorizedAuthorWithoutSeals(t *testing.T) {
	var (
		author   = common.HexToAddress("0x0001")
		b        = common.HexToAddress("0x0002")
		blockNum = uint64(7)
	)

	for _, permissionless := range []bool{false, true} {
		name := "pre-permissionless"
		if permissionless {
			name = "post-permissionless"
		}
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			sealer := &verifySealsTestSealer{author: author, committers: nil}
			header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}

			mGov := mock_gov.NewMockGovModule(ctrl)
			mValset := mock_valset.NewMockValsetModule(ctrl)
			if permissionless {
				mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return([]common.Address{author, b}, nil)
			} else {
				mValset.EXPECT().GetQualifiedValidators(blockNum).Return([]common.Address{author, b}, nil)
				mValset.EXPECT().GetCouncil(blockNum).Return([]common.Address{author, b}, nil)
			}

			v := &BlockValidator{
				config: permissionlessConfig(t, permissionless, blockNum),
				sealer: sealer, mGov: mGov, mValset: mValset,
			}
			assert.NoError(t, v.verifySeals(header, false))
		})
	}
}

// The import path still requires committed seals.
func TestImportPathStillRequiresCommittedSeals(t *testing.T) {
	var (
		author   = common.HexToAddress("0x0001")
		b        = common.HexToAddress("0x0002")
		blockNum = uint64(7)
	)
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	sealer := &verifySealsTestSealer{author: author, committers: nil}
	header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetCommittee(blockNum, uint64(0)).Return([]common.Address{author, b}, nil)

	v := &BlockValidator{
		config: permissionlessConfig(t, true, blockNum),
		sealer: sealer, mGov: mGov, mValset: mValset,
	}
	assert.ErrorIs(t, v.verifySeals(header, true), istanbul.ErrEmptyCommittedSeals)
}

// A hash-locked proposal keeps its original author's seal while the header's round
// advances. It must stay valid on both paths.
func TestRoundShiftedLockedProposalStaysValid(t *testing.T) {
	var (
		round0Proposer = common.HexToAddress("0x0001")
		round1Proposer = common.HexToAddress("0x0002")
		third          = common.HexToAddress("0x0003")
		blockNum       = uint64(7)
		committee      = []common.Address{round0Proposer, round1Proposer, third}
	)

	for _, withSeals := range []bool{false, true} {
		name := "proposal path"
		if withSeals {
			name = "import path"
		}
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			t.Cleanup(ctrl.Finish)

			base := &verifySealsTestSealer{author: round0Proposer, committers: committee, quorum: 3}
			sealer := &roundShiftSealer{verifySealsTestSealer: base, round: 1}
			header := &types.Header{Number: new(big.Int).SetUint64(blockNum)}

			mGov := mock_gov.NewMockGovModule(ctrl)
			mValset := mock_valset.NewMockValsetModule(ctrl)
			mValset.EXPECT().GetCommittee(blockNum, uint64(1)).Return(committee, nil)

			v := &BlockValidator{
				config: params.TestKaiaConfig("permissionless"),
				sealer: sealer, mGov: mGov, mValset: mValset,
			}
			assert.NoError(t, v.verifySeals(header, withSeals))
		})
	}
}
