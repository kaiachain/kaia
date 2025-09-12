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
	// Reset global variables for test
	common.IsPrecompiledContractAddress = suite.originalIsPrecompiledContractAddress
}

func (suite *ExecutionSpecBlockTestSuite) TestExecutionSpecBlock() {
	t := suite.T()

	if !common.FileExist(executionSpecBlockTestDir) {
		t.Skipf("directory %s does not exist", executionSpecBlockTestDir)
	}
	bt := new(testMatcher)

	// TODO-Kaia: should remove these skip
	// executing precompiled contracts with value transferring is not permitted
	bt.skipLoad(`^frontier\/opcodes\/all_opcodes\/all_opcodes.json`)
	// executing precompiled contracts with set code tx is not permitted
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs_2\/gas_diff_pointer_vs_direct_call.json`)

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
	// different amount of gas is consumed because 0x0b contract is added to access list by ActivePrecompiles although Cancun doesn't have it as a precompiled contract
	bt.skipLoad(`^frontier\/precompiles\/precompiles\/precompiles.json\/tests\/frontier\/precompiles\/test_precompiles.py::test_precompiles\[fork_Cancun-address_0xb-precompile_exists_False-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-31_bytes\]`)
	bt.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-32_bytes\]`)
	bt.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-empty_calldata\]`)
	bt.skipLoad(`^prague\/eip2537_bls_12_381_precompiles\/bls12_precompiles_before_fork\/precompile_before_fork.json\/tests\/prague\/eip2537_bls_12_381_precompiles\/test_bls12_precompiles_before_fork.py::test_precompile_before_fork\[fork_Cancun-state_test--G1ADD\]`)
	// type 3 tx (EIP-4844) is not supported
	bt.skipLoad(`^frontier\/scenarios\/scenarios\/scenarios.json\/tests\/frontier\/scenarios\/test_scenarios.py::test_scenarios\[fork_Osaka-blockchain_test-test_program_program_BLOBBASEFEE-debug\]`)
	bt.skipLoad(`^prague\/eip7623_increase_calldata_cost\/.*type_3.*`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/eoa_tx_after_set_code.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs.py::test_eoa_tx_after_set_code\[fork_Osaka-tx_type_3-evm_code_type_LEGACY-blockchain_test-different_block\]`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/eoa_tx_after_set_code.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs.py::test_eoa_tx_after_set_code\[fork_Osaka-tx_type_3-evm_code_type_LEGACY-blockchain_test-same_block\]`)
	// Kaia cannot calculate the same block hash as Ethereum
	bt.skipLoad(`^frontier\/scenarios\/scenarios\/scenarios.json\/tests\/frontier\/scenarios\/test_scenarios.py::test_scenarios\[fork_Osaka-blockchain_test-test_program_program_BLOCKHASH-debug\]`)

	// TODO: Skip EIP tests that are not yet supported; expect to remove them
	bt.skipLoad(`osaka/eip7594_peerdas`)
	bt.skipLoad(`osaka/eip7825_transaction_gas_limit_cap`)
	bt.skipLoad(`osaka/eip7918_blob_reserve_price`)
	bt.skipLoad(`osaka/eip7934_block_rlp_limit`)
	// TODO: When EIP-7951 is imeplemted, this skip should be removed: address_0x0000000000000000000000000000000000000100
	bt.skipLoad(`osaka/eip7951_p256verify_precompiles`)
	bt.skipLoad(`^frontier\/precompiles\/precompiles\/precompiles.json\/tests\/frontier\/precompiles\/test_precompiles.py::test_precompiles\[fork_Osaka-address_0x0000000000000000000000000000000000000100-precompile_exists_True-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs\/set_code_to_precompile.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs.py::test_set_code_to_precompile\[fork_Osaka-precompile_0x0000000000000000000000000000000000000100-`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs_2\/pointer_to_precompile.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs_2.py::test_pointer_to_precompile\[fork_Osaka-precompile_0x0000000000000000000000000000000000000100-`)
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs_2\/call_to_precompile_in_pointer_context.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs_2.py::test_call_to_precompile_in_pointer_context\[fork_Osaka-precompile_0x0000000000000000000000000000000000000100-`)
	// TODO: Why; Cannot run with "to" is address_0x0000000000000000000000000000000000000005 because precompiled contract address validation in TxInternalData#Validate
	bt.skipLoad(`^prague\/eip7702_set_code_tx\/set_code_txs_2\/call_to_precompile_in_pointer_context.json\/tests\/prague\/eip7702_set_code_tx\/test_set_code_txs_2.py::test_call_to_precompile_in_pointer_context\[fork_Osaka-precompile_0x0000000000000000000000000000000000000005-`)
	bt.skipLoad(`^osaka\/eip7883_modexp_gas_increase\/modexp_thresholds\/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-exact_gas\]`)
	bt.skipLoad(`^osaka\/eip7883_modexp_gas_increase\/modexp_thresholds\/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-extra_gas\]`)
	bt.skipLoad(`^osaka\/eip7883_modexp_gas_increase\/modexp_thresholds\/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-insufficient_gas\]`)
	// TODO: Investigate after all Osaka EIPs are applied
	bt.skipLoad(`^frontier\/identity_precompile\/identity\/call_identity_precompile.json\/tests\/frontier\/identity_precompile\/test_identity.py::test_call_identity_precompile\[fork_Osaka-blockchain_test_from_state_test-identity_1_nonzerovalue-call_type_CALL\]`)

	bt.walk(t, executionSpecBlockTestDir, func(t *testing.T, name string, test *BlockTest) {
		skipForks := []string{
			// Even if we skip fork levels, old EIPs are still retrospectively tested against Cancun or later forks.
			// The EEST framework was added when Kaia was at Cancun hardfork.
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
			"ParisToShanghaiAtTime15k",
			"Shanghai",
			"ShanghaiToCancunAtTime15k",
			// "Cancun",
			"CancunToPragueAtTime15k",
			// "Prague",
			"PragueToOsakaAtTime15k",
			// "Osaka",
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
