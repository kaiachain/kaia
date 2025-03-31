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

package system

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	contracts "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/auction"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestReadAuctioneer(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)
	var (
		senderKey, _ = crypto.GenerateKey()
		sender       = bind.NewKeyedTransactor(senderKey)
		addr         = common.HexToAddress("0xaaaa")

		alloc = blockchain.GenesisAlloc{
			sender.From: {
				Balance: big.NewInt(params.KAIA),
			},
			RegistryAddr: {
				Code:    RegistryMockCode,
				Balance: common.Big0,
			},
			AuctionEntryPointAddrMock: {
				Code:    AuctionEntryPointMockCode,
				Balance: common.Big0,
				Storage: map[common.Hash]common.Hash{
					common.BytesToHash([]byte{0x00}): addr.Hash(),
				},
			},
		}

		backend = backends.NewSimulatedBackend(alloc)
	)

	contract, _ := contracts.NewRegistryMockTransactor(RegistryAddr, backend)
	contract.Register(sender, AuctionEntryPointName, AuctionEntryPointAddrMock, common.Big1)
	backend.Commit()

	auctioneer, err := ReadAuctioneer(backend, AuctionEntryPointAddrMock, common.Big1)
	assert.Nil(t, err)
	assert.Equal(t, addr, auctioneer)
}

var data = auction.BidData{
	TargetTxHash:  common.HexToHash("0xcf5879724726228474db71caf83955a61341f2e5bacaa2b9d37f6e5f48241cbc"),
	BlockNumber:   11,
	Sender:        common.HexToAddress("0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
	To:            common.HexToAddress("0x5FC8d32690cc91D4c39d9d3abcBD16989F875707"),
	Nonce:         0,
	Bid:           new(big.Int).SetBytes(common.Hex2Bytes("8ac7230489e80000")),
	CallGasLimit:  10000000,
	Data:          common.Hex2Bytes("d09de08a"),
	SearcherSig:   common.Hex2Bytes("45eb1182b8ccffc953ebbb192dd5376cf3a9f53c060a3de6e42526ebe0cb91682e571067ac5aa38d5742f6534d0adc95793443056bd2041d551c7f3aff29967e1c"),
	AuctioneerSig: common.Hex2Bytes("27bba77731bbba25dcde852c54d05eb8d095b2d71b05ff9611ef5a123d1b2af01a7b42d0c8f8ab5ba07456932f4b1d29ff5f6e7b7fb8f59d41d22c4294d23bd91c"),
}

var testBid = &auction.Bid{
	BidData: data,
}

func TestEncodeAuctionCallData(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	encoded, err := EncodeAuctionCallData(testBid)
	assert.Nil(t, err)

	assert.Equal(t, common.Hex2Bytes("ca1575540000000000000000000000000000000000000000000000000000000000000020cf5879724726228474db71caf83955a61341f2e5bacaa2b9d37f6e5f48241cbc000000000000000000000000000000000000000000000000000000000000000b00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c80000000000000000000000005fc8d32690cc91d4c39d9d3abcbd16989f87570700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008ac7230489e8000000000000000000000000000000000000000000000000000000000000009896800000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000018000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000004d09de08a00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004145eb1182b8ccffc953ebbb192dd5376cf3a9f53c060a3de6e42526ebe0cb91682e571067ac5aa38d5742f6534d0adc95793443056bd2041d551c7f3aff29967e1c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004127bba77731bbba25dcde852c54d05eb8d095b2d71b05ff9611ef5a123d1b2af01a7b42d0c8f8ab5ba07456932f4b1d29ff5f6e7b7fb8f59d41d22c4294d23bd91c00000000000000000000000000000000000000000000000000000000000000"), encoded)
}
