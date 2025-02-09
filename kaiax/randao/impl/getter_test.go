// Copyright 2025 The Kaia Authors
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

	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/datasync/downloader"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

func testRandaoForkChainConfig(forkNum *big.Int) *params.ChainConfig {
	var config *params.ChainConfig

	config = &params.ChainConfig{
		ChainID: common.Big1,
		Governance: &params.GovernanceConfig{
			Reward: &params.RewardConfig{
				UseGiniCoeff:          true,
				StakingUpdateInterval: 86400,
			},
			KIP71: params.GetDefaultKIP71Config(),
		},
		RandaoRegistry: &params.RegistryConfig{
			Records: map[string]common.Address{
				system.Kip113Name: system.Kip113LogicAddrMock,
			},
			Owner: common.HexToAddress("0xffff"),
		},
	}
	config.LondonCompatibleBlock = big.NewInt(0)
	config.IstanbulCompatibleBlock = big.NewInt(0)
	config.EthTxTypeCompatibleBlock = big.NewInt(0)
	config.MagmaCompatibleBlock = big.NewInt(0)
	config.KoreCompatibleBlock = big.NewInt(0)
	config.ShanghaiCompatibleBlock = big.NewInt(0)
	config.CancunCompatibleBlock = big.NewInt(0)
	config.KaiaCompatibleBlock = big.NewInt(0)
	config.RandaoCompatibleBlock = forkNum

	return config
}

func TestGetBlsPubKey(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		db           = database.NewMemoryDBManager()
		allocStorage = system.AllocRegistry(&params.RegistryConfig{
			Records: map[string]common.Address{
				"KIP113": system.Kip113LogicAddrMock,
			},
			Owner: common.HexToAddress("0xffff"),
		})
		alloc = blockchain.GenesisAlloc{
			system.RegistryAddr: {
				Code:    system.RegistryMockCode,
				Balance: big.NewInt(0),
				Storage: allocStorage,
			},
			system.Kip113LogicAddrMock: {
				Code:    system.Kip113MockThreeCNCode,
				Balance: big.NewInt(0),
			},
		}
		config = testRandaoForkChainConfig(big.NewInt(0))

		pubKey1, _ = bls.PublicKeyFromBytes(hexutil.MustDecode("0x" + "b716443d8d1b3c1230d1d186b1db0db80f79f72805646ba8135b98242df276bdbfb5dea0201c0258d6b60f30724f28e3"))
		pubKey2, _ = bls.PublicKeyFromBytes(hexutil.MustDecode("0x" + "a5b6d96a1bb2bd8ec5480d112dc6bbad46ec08937b9320187928c0ed27339791186f581397c5a9679e49f6ac459d5a48"))
		pubKey3, _ = bls.PublicKeyFromBytes(hexutil.MustDecode("0x" + "a2093da481a55e7e374de2fa19a8d9acbf055a52048d697d87de864fab9d334bbd4d838c68d53022f355c06fb4cd6722"))

		expected = map[common.Address]bls.PublicKey{
			common.HexToAddress("0x0000000000000000000000000000000000000001"): pubKey1,
			common.HexToAddress("0x0000000000000000000000000000000000000002"): pubKey2,
			common.HexToAddress("0x0000000000000000000000000000000000000003"): pubKey3,
		}
	)

	backend := backends.NewSimulatedBackendWithDatabase(db, alloc, config)

	mRandao := NewRandaoModule()
	fakeDownloader := &downloader.FakeDownloader{}
	mRandao.Init(&InitOpts{
		ChainConfig: config,
		Chain:       backend.BlockChain(),
		Downloader:  fakeDownloader,
	})

	for addr, pk := range expected {
		expected, err := mRandao.GetBlsPubkey(addr, big.NewInt(1))
		assert.NoError(t, err)
		assert.Equal(t, expected, pk)
	}
}
