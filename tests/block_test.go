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
// This file is derived from tests/block_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package tests

import (
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/suite"
)

func TestBlockchain(t *testing.T) {
	t.Parallel()

	bt := new(testMatcher)
	// General state tests are 'exported' as blockchain tests, but we can run them natively.
	bt.skipLoad(`^GeneralStateTests/`)
	// Skip random failures due to selfish mining test.
	// bt.skipLoad(`^bcForgedTest/bcForkUncle\.json`)
	bt.skipLoad(`^bcMultiChainTest/(ChainAtoChainB_blockorder|CallContractFromNotBestBlock)`)
	bt.skipLoad(`^bcTotalDifficultyTest/(lotsOfLeafs|lotsOfBranches|sideChainWithMoreTransactions)`)
	// This test is broken
	bt.fails(`blockhashNonConstArg_Constantinople`, "Broken test")

	// Still failing tests
	// bt.skipLoad(`^bcWalletTest.*_Byzantium$`)

	// TODO-Kaia Update BlockchainTests first to enable this test, since block header has been changed in Kaia.
	//bt.walk(t, blockTestDir, func(t *testing.T, name string, test *BlockTest) {
	//	if err := bt.checkFailure(t, name, test.Run()); err != nil {
	//		t.Error(err)
	//	}
	//})
}

// TestExecutionSpecState runs the state_test fixtures from execution-spec-tests.
type ExecutionSpecBlockTestSuite struct {
	suite.Suite
	originalIsPrecompiledContractAddress func(common.Address, interface{}) bool
}

func (suite *ExecutionSpecBlockTestSuite) SetupSuite() {
	suite.originalIsPrecompiledContractAddress = common.IsPrecompiledContractAddress
	common.IsPrecompiledContractAddress = isPrecompiledContractAddressForEthTest
	blockchain.UseKaiaCancunExtCodeHashFee = true
	blockchain.CreateContractWithCodeFormatInExecutionSpecTest = true
}

func (suite *ExecutionSpecBlockTestSuite) TearDownSuite() {
	// Reset global variables for test
	common.IsPrecompiledContractAddress = suite.originalIsPrecompiledContractAddress
	blockchain.UseKaiaCancunExtCodeHashFee = false
	blockchain.GasLimitInExecutionSpecTest = 0
	blockchain.CreateContractWithCodeFormatInExecutionSpecTest = false
	types.IsPragueInExecutionSpecTest = false
}

func (suite *ExecutionSpecBlockTestSuite) TestExecutionSpecBlock() {
	t := suite.T()

	if !common.FileExist(executionSpecBlockTestDir) {
		t.Skipf("directory %s does not exist", executionSpecBlockTestDir)
	}
	bt := new(testMatcher)

	// TODO-Kaia: should remove these skip
	// json format error
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/invalid_tx_invalid_auth_signature.json`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/tx_validity_chain_id.json`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/tx_validity_nonce.json`)

	// only target after shanghai
	bt.skipLoad(`^frontier\/`)
	bt.skipLoad(`^homestead\/`)
	bt.skipLoad(`^byzantium\/`)
	bt.skipLoad(`^constantinople\/`)
	bt.skipLoad(`^istanbul\/`)
	bt.skipLoad(`^berlin\/`)
	bt.skipLoad(`^paris\/`)

	bt.skipLoad(`^prague\/eip2537_bls_12_381_precompiles`) // gas error
	bt.skipLoad(`^prague\/eip7702_set_code_tx`)            // state, gas (after update we should do it)

	// tests to skip
	// unsupported EIPs
	bt.skipLoad(`^shanghai\/eip4895_withdrawals\/`)
	bt.skipLoad(`^cancun\/eip4788_beacon_root\/`)
	bt.skipLoad(`^cancun\/eip4844_blobs\/`)
	bt.skipLoad(`^cancun\/eip7516_blobgasfee\/`)
	bt.skipLoad(`^prague\/eip7251_consolidations`)
	bt.skipLoad(`^prague\/eip7685_general_purpose_el_requests`)
	bt.skipLoad(`^prague\/eip7002_el_triggerable_withdrawals`)
	bt.skipLoad(`^prague\/eip6110_deposits`)
	// type 3 tx (EIP-4844) is not supported
	bt.skipLoad(`^prague\/eip7623_increase_calldata_cost\/.*type_3.*`)

	bt.walk(t, executionSpecBlockTestDir, func(t *testing.T, name string, test *BlockTest) {
		skipForks := []string{
			"Frontier",
			"Homestead",
			"Byzantium",
			"Constantinople",
			"ConstantinopleFix",
			"Istanbul",
			"Berlin",
			"London",
			"Merge",
			"Paris",
			"Shanghai",
			"ShanghaiToCancunAtTime15k",
			"CancunToPragueAtTime15k",
			// "Cancun",
			// "Prague",
		}
		for _, fork := range skipForks {
			if test.json.Network == fork {
				t.Skip()
			}
		}

		if err := bt.checkFailure(t, name, test.Run()); err != nil {
			t.Error(err)
		}
	})
}

func TestExecutionSpecBlockTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSpecBlockTestSuite))
}
