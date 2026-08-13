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

package impl

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/faker"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	valset_mock "github.com/kaiachain/kaia/kaiax/valset/mock"
	"github.com/kaiachain/kaia/params"
	chain_mock "github.com/kaiachain/kaia/work/mocks"
)

func init() {
	blockchain.InitDeriveSha(&params.ChainConfig{DeriveShaImpl: 0}) // DeriveSha is irrelevant for this test. Set any value.
}

// rewardModuleTestConfig configures makeTestRewardModuleWithConfig. Zero values are sensible defaults.
type rewardModuleTestConfig struct {
	Simple   bool // false => WeightedRandom, true => RoundRobin
	Deferred bool
	IsPrague bool

	// Proposer address returned by Chain.Sealer().Author(header).
	// Zero value falls back to params.AuthorAddressForTesting.
	Proposer common.Address

	// Block fixture. When Header is nil, block- and number-keyed mocks are
	// replaced with gomock.Any() so the helper can be used by tests that
	// don't fix a block (e.g. VerifyHeader).
	Header   *types.Header
	Txs      []*types.Transaction
	Receipts []*types.Receipt
}

// Make a RewardModule that works with the given block.
func makeTestRewardModule(t *testing.T,
	simple, deferred, isPrague bool,
	header *types.Header, txs []*types.Transaction, receipts []*types.Receipt,
) *RewardModule {
	return makeTestRewardModuleWithConfig(t, rewardModuleTestConfig{
		Simple:   simple,
		Deferred: deferred,
		IsPrague: isPrague,
		Header:   header,
		Txs:      txs,
		Receipts: receipts,
	})
}

func makeTestRewardModuleWithConfig(t *testing.T, cfg rewardModuleTestConfig) *RewardModule {
	proposerPolicy := uint64(istanbul.WeightedRandom)
	var pragueForkBlock *big.Int
	if cfg.Simple {
		proposerPolicy = uint64(istanbul.RoundRobin)
	}
	if cfg.IsPrague {
		pragueForkBlock = big.NewInt(60)
	}
	proposer := cfg.Proposer
	if proposer == (common.Address{}) {
		proposer = params.AuthorAddressForTesting
	}

	var (
		mockCtrl = gomock.NewController(t)
		chain    = chain_mock.NewMockBlockChain(mockCtrl)
		mStaking = staking_mock.NewMockStakingModule(mockCtrl)
		mGov     = gov_mock.NewMockGovModule(mockCtrl)
		mValset  = valset_mock.NewMockValsetModule(mockCtrl)

		chainConfig = &params.ChainConfig{
			ChainID:                  big.NewInt(31337),
			IstanbulCompatibleBlock:  big.NewInt(10),
			LondonCompatibleBlock:    big.NewInt(10),
			EthTxTypeCompatibleBlock: big.NewInt(10),
			MagmaCompatibleBlock:     big.NewInt(20),
			KoreCompatibleBlock:      big.NewInt(30),
			ShanghaiCompatibleBlock:  big.NewInt(40),
			CancunCompatibleBlock:    big.NewInt(50),
			KaiaCompatibleBlock:      big.NewInt(60),
			PragueCompatibleBlock:    pragueForkBlock,
		}

		paramset = gov.ParamSet{
			ProposerPolicy:         proposerPolicy,
			UnitPrice:              25e9,
			MintingAmount:          big.NewInt(6.4e18),
			MinimumStake:           big.NewInt(5_000_000),
			DeferredTxFee:          cfg.Deferred,
			Ratio:                  "50/20/30",
			Kip82Ratio:             "20/80",
			StakingRewardThreshold: big.NewInt(5_000_000),
			UseFlexReward:          false,
		}

		stakingInfo = makeTestStakingInfo([]uint64{5_000_001, 5_000_002, 5_000_003, 5_000_000}, cfg.IsPrague)
	)

	r := NewRewardModule()
	r.Init(&InitOpts{
		ChainConfig:   chainConfig,
		Chain:         chain,
		GovModule:     mGov,
		StakingModule: mStaking,
		ValsetModule:  mValset,
	})

	chain.EXPECT().Sealer().Return(faker.NewFakerWithFixedSealer(proposer)).AnyTimes()
	mStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(stakingInfo, nil).AnyTimes()

	if cfg.Header != nil {
		block := types.NewBlock(cfg.Header, cfg.Txs, cfg.Receipts)
		chain.EXPECT().GetBlockByNumber(cfg.Header.Number.Uint64()).Return(block).AnyTimes()
		chain.EXPECT().GetReceiptsByBlockHash(block.Hash()).Return(cfg.Receipts).AnyTimes()
		mGov.EXPECT().GetParamSet(cfg.Header.Number.Uint64()).Return(paramset).AnyTimes()
	} else {
		mGov.EXPECT().GetParamSet(gomock.Any()).Return(paramset).AnyTimes()
	}

	return r
}

func makeTestStakingInfo(stakingAmounts []uint64, isPrague bool) *staking.StakingInfo {
	var clStakingInfos staking.CLStakingInfos
	if isPrague {
		clStakingInfos = staking.CLStakingInfos{
			&staking.CLStakingInfo{
				CLNodeId:        common.HexToAddress("0xa01"), // For Node1
				CLPoolAddr:      common.HexToAddress("0xe01"),
				CLStakingAmount: 999_999, // Total staking amount = 5_000_001 + 999_999 = 6_000_000
			},
			&staking.CLStakingInfo{
				CLNodeId:        common.HexToAddress("0xa02"), // For Node2
				CLPoolAddr:      common.HexToAddress("0xe02"),
				CLStakingAmount: 1_999_998, // Total staking amount = 5_000_002 + 1_999_998 = 7_000_000
			},
			&staking.CLStakingInfo{
				CLNodeId:        common.HexToAddress("0xa03"), // For Node3
				CLPoolAddr:      common.HexToAddress("0xe03"),
				CLStakingAmount: 2_999_997, // Total staking amount = 5_000_003 + 2_999_997 = 8_000_000
			},
		}
	}

	return &staking.StakingInfo{
		NodeIds:          []common.Address{common.HexToAddress("0xa01"), common.HexToAddress("0xa02"), common.HexToAddress("0xa03"), common.HexToAddress("0xa04")},
		StakingContracts: []common.Address{common.HexToAddress("0xb01"), common.HexToAddress("0xb02"), common.HexToAddress("0xb03"), common.HexToAddress("0xb04")},
		RewardAddrs:      []common.Address{common.HexToAddress("0xc01"), common.HexToAddress("0xc02"), common.HexToAddress("0xc03"), common.HexToAddress("0xc04")},
		StakingAmounts:   stakingAmounts,
		KIFAddr:          common.HexToAddress("0xd01"),
		KEFAddr:          common.HexToAddress("0xd02"),
		KPFAddr:          common.HexToAddress("0x0"),
		CLStakingInfos:   clStakingInfos,
	}
}

// Note: totalFee = 0.025e18
func makeTestPreMagmaBlock(num int64) (*types.Header, []*types.Transaction, []*types.Receipt) {
	return &types.Header{
			Number:     big.NewInt(num),
			BaseFee:    nil,
			GasUsed:    1_000_000,
			Rewardbase: common.HexToAddress("0xfff"),
		},
		[]*types.Transaction{
			makeTestTx_type0(25e9, 200_000),
			makeTestTx_type0(25e9, 400_000),
			makeTestTx_type0(25e9, 600_000),
			makeTestTx_type0(25e9, 800_000),
		},
		[]*types.Receipt{
			{GasUsed: 100_000},
			{GasUsed: 200_000},
			{GasUsed: 300_000},
			{GasUsed: 400_000},
		}
}

// Note: totalFee = 0.027e18
func makeTestMagmaBlock(num int64) (*types.Header, []*types.Transaction, []*types.Receipt) {
	return &types.Header{
			Number:     big.NewInt(num),
			BaseFee:    big.NewInt(27e9),
			GasUsed:    1_000_000,
			Rewardbase: common.HexToAddress("0xfff"),
		},
		[]*types.Transaction{
			makeTestTx_type0(27e9, 200_000),
			makeTestTx_type0(28e9, 400_000),
			makeTestTx_type0(29e9, 600_000),
			makeTestTx_type0(30e9, 800_000),
		},
		[]*types.Receipt{
			{GasUsed: 100_000},
			{GasUsed: 200_000},
			{GasUsed: 300_000},
			{GasUsed: 400_000},
		}
}

// Note: totalFee = 0.0376e18
func makeTestKaiaBlock(num int64) (*types.Header, []*types.Transaction, []*types.Receipt) {
	return &types.Header{
			Number:     big.NewInt(num),
			BaseFee:    big.NewInt(27e9),
			GasUsed:    1_000_000,
			Rewardbase: common.HexToAddress("0xfff"),
		},
		[]*types.Transaction{
			makeTestTx_type2(50e9, 1e9, 200_000), // effectivePrice = 28e9, effectiveTip = 1e9
			makeTestTx_type2(29e9, 7e9, 400_000), // effectivePrice = 29e9, effectiveTip = 7e9
			makeTestTx_type0(30e9, 600_000),      // effectivePrice = 30e9, effectiveTip = 3e9
			makeTestTx_type0(50e9, 800_000),      // effectivePrice = 50e9, effectiveTip = 33e9
		},
		[]*types.Receipt{
			{GasUsed: 100_000},
			{GasUsed: 200_000},
			{GasUsed: 300_000},
			{GasUsed: 400_000},
		}
}

// makeTestOsakaBlock creates a Kaia block with one blob transaction.
// execFee = 0.0376e18 (same as makeTestKaiaBlock)
// blobFee = BlobTxBlobGasPerBlob * baseFee * BlobBaseFeeMultiplier = 131072 * 27e9 * 8 = 28311552000000000
func makeTestOsakaBlock(num int64) (*types.Header, []*types.Transaction, []*types.Receipt) {
	blobGasUsed := uint64(params.BlobTxBlobGasPerBlob)
	header, txs, receipts := makeTestKaiaBlock(num)
	header.BlobGasUsed = &blobGasUsed
	return header, txs, receipts
}

func makeTestTx_type0(gasPrice, gasLimit int64) *types.Transaction {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataLegacy{
		Price:     big.NewInt(gasPrice),
		GasLimit:  uint64(gasLimit),
		Recipient: &addr,
	}), types.NewLondonSigner(big.NewInt(31337)), key)
	return tx
}

func makeTestTx_type2(feeCap, tipCap, gasLimit int64) *types.Transaction {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	tx, _ := types.SignTx(types.NewTx(&types.TxInternalDataEthereumDynamicFee{
		GasFeeCap: big.NewInt(feeCap),
		GasTipCap: big.NewInt(tipCap),
		GasLimit:  uint64(gasLimit),
		Recipient: &addr,
	}), types.NewLondonSigner(big.NewInt(31337)), key)
	return tx
}
