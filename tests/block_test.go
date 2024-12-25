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

// func execBlockTest(t *testing.T, st *testMatcher, test *StateTest, name string, skipForks []string, isTestExecutionSpecState bool) {
// 	for _, subtest := range test.Subtests() {
// 		subtest := subtest
// 		key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
// 		name := name + "/" + key
// 		t.Run(key, func(t *testing.T) {
// 			for _, skip := range skipForks {
// 				if skip == subtest.Fork {
// 					t.Skipf("%s not supported yet", subtest.Fork)
// 				}
// 			}
// 			withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
// 				err := test.Run(subtest, vmconfig, isTestExecutionSpecState)
// 				return st.checkFailure(t, name, err)
// 			})
// 		})
// 	}
// }

// TestExecutionSpecState runs the state_test fixtures from execution-spec-tests.
type ExecutionSpecBlockTestSuite struct {
	suite.Suite
	originalIsPrecompiledContractAddress func(common.Address, interface{}) bool
}

func (suite *ExecutionSpecBlockTestSuite) SetupSuite() {
	suite.originalIsPrecompiledContractAddress = common.IsPrecompiledContractAddress
	common.IsPrecompiledContractAddress = isPrecompiledContractAddressForEthTest
}

func (suite *ExecutionSpecBlockTestSuite) TearDownSuite() {
	common.IsPrecompiledContractAddress = suite.originalIsPrecompiledContractAddress
}

func (suite *ExecutionSpecBlockTestSuite) TestExecutionSpecBlock() {
	t := suite.T()

	if !common.FileExist(executionSpecBlockTestDir) {
		t.Skipf("directory %s does not exist", executionSpecBlockTestDir)
	}
	bt := new(testMatcher)

	// TODO-Kaia: should remove these skip
	// json format error
	// bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/invalid_tx_invalid_auth_signature.json`)
	// bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/tx_validity_chain_id.json`)
	// bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/tx_validity_nonce.json`)

	// only target cuncun
	bt.skipLoad(`^berlin\/`)
	bt.skipLoad(`^byzantium\/`)
	// bt.skipLoad(`^cancun\/`)
	bt.skipLoad(`^cancun\/eip1153_tstore\/tload.*\/`)
	bt.skipLoad(`^cancun\/eip1153_tstore\/tstorage.*\/`)
	bt.skipLoad(`^cancun\/eip1153_tstore\/tstore_reentrancy\/`)
	bt.skipLoad(`^cancun\/eip5656_mcopy\/`)
	bt.skipLoad(`^cancun\/eip6780_selfdestruct\/`)
	bt.skipLoad(`^cancun\/eip7516_blobgasfee\/`)
	bt.skipLoad(`^constantinople\/`)
	bt.skipLoad(`^frontier\/`)
	bt.skipLoad(`^homestead\/`)
	bt.skipLoad(`^istanbul\/`)
	bt.skipLoad(`^paris\/`)
	bt.skipLoad(`^prague\/`)
	bt.skipLoad(`^shanghai\/`)

	// tests to skip
	// unsupported EIPs
	bt.skipLoad(`^cancun\/eip4788_beacon_root\/`)
	bt.skipLoad(`^cancun\/eip4844_blobs\/`)

	bt.walk(t, executionSpecBlockTestDir, func(t *testing.T, name string, test *BlockTest) {
		if err := bt.checkFailure(t, name, test.Run()); err != nil {
			t.Error(err)
		}
	})
}

func TestExecutionSpecBlockTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSpecBlockTestSuite))
}
