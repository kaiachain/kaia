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
package governance

import (
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

type testBlockChain struct {
	num    uint64
	config *params.ChainConfig
}

func newTestBlockchain(config *params.ChainConfig) *testBlockChain {
	return &testBlockChain{
		config: config,
	}
}

func newTestGovernanceApi() *GovernanceAPI {
	config := params.MainnetChainConfig
	config.Governance.KIP71 = params.GetDefaultKIP71Config()
	govApi := NewGovernanceAPI(NewMixedEngine(config, database.NewMemoryDBManager()), nil)
	govApi.governance.SetNodeAddress(common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"))
	bc := newTestBlockchain(config)
	govApi.governance.SetBlockchain(bc)
	return govApi
}

func TestUpperBoundBaseFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curLowerBoundBaseFee := govApi.governance.CurrentParams().LowerBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidUpperBoundBaseFee := curLowerBoundBaseFee - 100
	_, err := govApi.Vote("kip71.upperboundbasefee", invalidUpperBoundBaseFee)
	assert.Equal(t, err, errInvalidUpperBound)
}

func TestLowerBoundFeeSet(t *testing.T) {
	govApi := newTestGovernanceApi()

	curUpperBoundBaseFee := govApi.governance.CurrentParams().UpperBoundBaseFee()
	// unexpected case : upperboundbasefee < lowerboundbasefee
	invalidLowerBoundBaseFee := curUpperBoundBaseFee + 100
	_, err := govApi.Vote("kip71.lowerboundbasefee", invalidLowerBoundBaseFee)
	assert.Equal(t, err, errInvalidLowerBound)
}

func (bc *testBlockChain) Engine() consensus.Engine                    { return nil }
func (bc *testBlockChain) GetHeader(common.Hash, uint64) *types.Header { return nil }
func (bc *testBlockChain) GetHeaderByNumber(val uint64) *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(val),
	}
}

func (bc *testBlockChain) GetReceiptsByBlockHash(hash common.Hash) types.Receipts {
	return types.Receipts{
		&types.Receipt{GasUsed: 10},
		&types.Receipt{GasUsed: 10},
	}
}

func (bc *testBlockChain) GetBlockByNumber(num uint64) *types.Block {
	return types.NewBlockWithHeader(bc.GetHeaderByNumber(num))
}
func (bc *testBlockChain) StateAt(root common.Hash) (*state.StateDB, error) { return nil, nil }
func (bc *testBlockChain) State() (*state.StateDB, error)                   { return nil, nil }
func (bc *testBlockChain) Config() *params.ChainConfig {
	return bc.config
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(bc.CurrentHeader())
}

func (bc *testBlockChain) CurrentHeader() *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(bc.num),
	}
}

func (bc *testBlockChain) SetBlockNum(num uint64) {
	bc.num = num
}

func (bc *testBlockChain) GetBlock(hash common.Hash, num uint64) *types.Block {
	return bc.GetBlockByNumber(num)
}
