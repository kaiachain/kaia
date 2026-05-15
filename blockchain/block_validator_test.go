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
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	mock_gov "github.com/kaiachain/kaia/kaiax/gov/mock"
	mock_valset "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

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

func TestVerifySealsChecksAuthorMatchesProposer(t *testing.T) {
	var (
		author        = common.HexToAddress("0x0001")
		otherProposer = common.HexToAddress("0x0002")
		blockNum      = uint64(7)
		header        = &types.Header{Number: new(big.Int).SetUint64(blockNum)}
		sealer        = faker.NewFakerWithFixedSealer(author)
	)

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mGov := mock_gov.NewMockGovModule(ctrl)
	mValset := mock_valset.NewMockValsetModule(ctrl)
	mValset.EXPECT().GetProposer(blockNum, uint64(0)).Return(otherProposer, nil)

	validator := &BlockValidator{
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.ErrorIs(t, validator.verifySeals(header), consensus.ErrUnauthorized)
}

func TestVerifySealsAcceptsExpectedProposer(t *testing.T) {
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
	mValset.EXPECT().GetProposer(blockNum, uint64(0)).Return(author, nil)
	mValset.EXPECT().GetQualifiedValidators(blockNum).Return([]common.Address{author}, nil)
	mValset.EXPECT().GetCouncil(blockNum).Return([]common.Address{author}, nil)
	mGov.EXPECT().GetParamSet(blockNum).Return(gov.ParamSet{CommitteeSize: 1})

	validator := &BlockValidator{
		sealer:  sealer,
		mGov:    mGov,
		mValset: mValset,
	}

	assert.NoError(t, validator.verifySeals(header))
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
