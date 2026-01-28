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
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/suite"
)

// TestExecutionSpecState runs the state_test fixtures from execution-spec-tests.
type ExecutionSpecBlockTestSuite struct {
	suite.Suite
	originalIsPrecompiledContractAddress func(common.Address, interface{}) bool
	originalBlobTxMaxBlobs               int
}

func (suite *ExecutionSpecBlockTestSuite) SetupSuite() {
	suite.originalIsPrecompiledContractAddress = common.IsPrecompiledContractAddress
	common.IsPrecompiledContractAddress = isPrecompiledContractAddressForEthTest
	suite.originalBlobTxMaxBlobs = params.BlobTxMaxBlobs
	params.BlobTxMaxBlobs = EthBlobTxMaxBlobs
}

func (suite *ExecutionSpecBlockTestSuite) TearDownSuite() {
	// Reset global variables for test
	common.IsPrecompiledContractAddress = suite.originalIsPrecompiledContractAddress
	params.BlobTxMaxBlobs = suite.originalBlobTxMaxBlobs
}

func (suite *ExecutionSpecBlockTestSuite) TestExecutionSpecBlock() {
	t := suite.T()

	if !common.FileExist(executionSpecBlockTestDir) {
		t.Skipf("directory %s does not exist", executionSpecBlockTestDir)
	}
	bt := new(testMatcher)

	// should be skipped
	// note: Unsupported EIPs
	bt.skipLoad(`^shanghai/eip4895_withdrawals/`)
	bt.skipLoad(`^cancun/eip4788_beacon_root/`)
	bt.skipLoad(`^prague/eip7251_consolidations`)
	bt.skipLoad(`^prague/eip7685_general_purpose_el_requests`)
	bt.skipLoad(`^prague/eip7002_el_triggerable_withdrawals`)
	bt.skipLoad(`^prague/eip6110_deposits`)
	bt.skipLoad(`^osaka/eip7825_transaction_gas_limit_cap`)

	// should be skipped
	// note: Type 3 tx (EIP-4844) is not supported before Osaka. See tests/block_test_util.go:decode().
	bt.skipLoad(`^static/state_tests/Cancun/stEIP4844_blobtransactions/.*\[fork_(Cancun|Prague)`)
	bt.skipLoad(`^static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhashOutOfRange.json/tests/static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhashOutOfRangeFiller.yml::opcodeBlobhashOutOfRange\[fork_Prague-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhashOutOfRange.json/tests/static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhashOutOfRangeFiller.yml::opcodeBlobhashOutOfRange\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhBounds.json/tests/static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhBoundsFiller.yml::opcodeBlobhBounds\[fork_Prague-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhBounds.json/tests/static/state_tests/Cancun/stEIP4844_blobtransactions/opcodeBlobhBoundsFiller.yml::opcodeBlobhBounds\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^frontier/scenarios/test_scenarios.json/tests/frontier/scenarios/test_scenarios.py::test_scenarios\[fork_Prague-blockchain_test-test_program_program_BLOBBASEFEE-debug\]`)
	bt.skipLoad(`^frontier/scenarios/test_scenarios.json/tests/frontier/scenarios/test_scenarios.py::test_scenarios\[fork_Cancun-blockchain_test-test_program_program_BLOBBASEFEE-debug\]`)
	bt.skipLoad(`^istanbul/eip1344_chainid/test_chainid.json/tests/istanbul/eip1344_chainid/test_chainid.py::test_chainid\[fork_Prague-typed_transaction_3-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^istanbul/eip1344_chainid/test_chainid.json/tests/istanbul/eip1344_chainid/test_chainid.py::test_chainid\[fork_Cancun-typed_transaction_3-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^istanbul/eip1344_chainid/.*\[fork_(Cancun|Prague).*typed_transaction_3`)
	bt.skipLoad(`^cancun/eip4844_blobs/.*\[fork_(Cancun|Prague)`)
	bt.skipLoad(`^prague/eip7702_set_code_tx/test_eoa_tx_after_set_code.json/tests/prague/eip7702_set_code_tx/test_set_code_txs.py::test_eoa_tx_after_set_code\[fork_Prague-tx_type_3-evm_code_type_LEGACY-blockchain_test-same_block\]`)
	bt.skipLoad(`^prague/eip7702_set_code_tx/test_eoa_tx_after_set_code.json/tests/prague/eip7702_set_code_tx/test_set_code_txs.py::test_eoa_tx_after_set_code\[fork_Prague-tx_type_3-evm_code_type_LEGACY-blockchain_test-different_block\]`)
	bt.skipLoad(`^prague/eip7623_increase_calldata_cost/test_transaction_validity_type_3.json/tests/prague/eip7623_increase_calldata_cost/test_transaction_validity.py::test_transaction_validity_type_3\[fork_(Cancun|Prague)`)
	bt.skipLoad(`^prague/eip7623_increase_calldata_cost/.*type_3.*\[fork_Prague`)
	bt.skipLoad(`^prague/eip7623_increase_calldata_cost/.*\[fork_(Cancun|Prague).*type_3`)

	// should be skipped
	// note: Kaia cannot calculate the same block hash as Ethereum.
	bt.skipLoad(`^frontier/scenarios/test_scenarios.json/tests/frontier/scenarios/test_scenarios.py::test_scenarios\[fork_Cancun-blockchain_test-test_program_program_BLOCKHASH-debug\]`)
	bt.skipLoad(`^frontier/scenarios/test_scenarios.json/tests/frontier/scenarios/test_scenarios.py::test_scenarios\[fork_Prague-blockchain_test-test_program_program_BLOCKHASH-debug\]`)
	bt.skipLoad(`^frontier/scenarios/test_scenarios.json/tests/frontier/scenarios/test_scenarios.py::test_scenarios\[fork_Osaka-blockchain_test-test_program_program_BLOCKHASH-debug\]`)

	// should be skipped
	// note: Kaia's MaxBlockSize (10 MiB) higher than Ethereum's (8 MiB), so max_plus_1 is accepted in Kaia.
	bt.skipLoad(`^osaka/eip7934_block_rlp_limit/test_block_at_rlp_size_limit_boundary.json/tests/osaka/eip7934_block_rlp_limit/test_max_block_rlp_size.py::test_block_at_rlp_size_limit_boundary\[fork_Osaka-blockchain_test-max_rlp_size_plus_1_byte\]`)

	// should be skipped
	// note: Withdrawal (eip4895) is not supported.
	bt.skipLoad(`^osaka/eip7934_block_rlp_limit/test_block_at_rlp_limit_with_withdrawals.json`)

	// should be skipped
	// note: There are two causes of failure:
	//   - Different amount of gas is consumed because 0x0b contract is added to access list by ActivePrecompiles although Cancun doesn't have it as a precompiled contract.
	//     - See https://github.com/kaiachain/kaia/blob/a7cb0f8aef2ce813e7ceda38e38353d63437489d/blockchain/vm/contracts.go#L242-L250.
	//   - Executing precompiled contracts with value transferring is not permitted
	//     - See https://github.com/kaiachain/kaia/blob/d44ae2f4269a84bd379b4e992d8e3be46b7e5ad3/blockchain/vm/evm.go#L267-L270
	// But the failure happens in only Cancun fork, so we can skip these.
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Cancun-blockchain_test_from_state_test-all_then_yes_from_prague\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Cancun-blockchain_test_from_state_test-all_then_yes_from_prague-15\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Cancun-blockchain_test_from_state_test-all_then_yes_from_prague-8\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Cancun-blockchain_test_from_state_test-yes.*\]`)

	// should be skipped
	// note: Different amount of gas is consumed because 0x0b contract is added to access list by ActivePrecompiles although Ethereum's Cancun fork doesn't have it as a precompiled contract.
	// See https://github.com/kaiachain/kaia/blob/a7cb0f8aef2ce813e7ceda38e38353d63437489d/blockchain/vm/contracts.go#L242-L250.
	// but the failure happens in only Cancun fork, so we can skip these.
	bt.skipLoad(`^static/state_tests/stSpecialTest/failed_tx_xcf416c53_Paris.json/tests/static/state_tests/stSpecialTest/failed_tx_xcf416c53_ParisFiller.json::failed_tx_xcf416c53_Paris\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^frontier/precompiles/test_precompile_absence.json/tests/frontier/precompiles/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-31_bytes\]`)
	bt.skipLoad(`^frontier/precompiles/test_precompile_absence.json/tests/frontier/precompiles/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-32_bytes\]`)
	bt.skipLoad(`^frontier/precompiles/test_precompile_absence.json/tests/frontier/precompiles/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-blockchain_test_from_state_test-empty_calldata\]`)
	bt.skipLoad(`^frontier/precompiles/test_precompiles.json/tests/frontier/precompiles/test_precompiles.py::test_precompiles\[fork_Cancun-address_0x000000000000000000000000000000000000000b-precompile_exists_False-blockchain_test_from_state_test\]`)

	// should be skipped
	// note: Kaia's Cancun fork allow creating SCA with 0xEF code, but Ethereum's Cancun fork does not.
	// See https://github.com/kaiachain/kaia/blob/e89d0fece3d9ca7428e8ba541afa6394074e8f85/blockchain/vm/evm.go#L588-L591.
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateTransactionRefundEF.json/tests/static/state_tests/stCreateTest/CreateTransactionRefundEFFiller.yml::CreateTransactionRefundEF\[fork_Cancun-blockchain_test_from_state_test-refund_EF\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create-0xef-v0\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create-0xef-v1\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create2-0xef-v0\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create2-0xef-v1\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CREATE2_RefundEF.json/tests/static/state_tests/stCreateTest/CREATE2_RefundEFFiller.yml::CREATE2_RefundEF\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stCreate2/CREATE2_FirstByte_loop.json/tests/static/state_tests/stCreate2/CREATE2_FirstByte_loopFiller.yml::CREATE2_FirstByte_loop\[fork_Cancun-blockchain_test_from_state_test-invalidByte\]`)
	bt.skipLoad(`^frontier/create/test_create_one_byte.json/tests/frontier/create/test_create_one_byte.py::test_create_one_byte\[fork_Cancun-create_opcode_CREATE-evm_code_type_LEGACY-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^frontier/create/test_create_one_byte.json/tests/frontier/create/test_create_one_byte.py::test_create_one_byte\[fork_Cancun-create_opcode_CREATE2-evm_code_type_LEGACY-blockchain_test_from_state_test\]`)

	// should be skipped
	// note: Kaia's Cancun fork recognizes 0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b as a contract (AccountKeyTypeFail).
	// Kaia's after-Prague forks recognizes 0xa94f5374fce5edbc8e2a8697c15331677e6ebf0b as EOA with code (AccountKeyTypePublic).
	// Sender cannot be a contract in Kaia. But the failure happens in only Cancun fork, so we can skip these.
	bt.skipLoad(`^static/state_tests/stTransactionTest/Opcodes_TransactionInit.json/tests/static/state_tests/stTransactionTest/Opcodes_TransactionInitFiller.json::Opcodes_TransactionInit\[fork_Cancun-blockchain_test_from_state_test.*\]`)

	// should be skipped
	// note: `panic: can't encode object at ...: rlp: cannot encode negative big.Int` happens.
	// But the failure happens in only Cancun fork, so we can skip these.
	bt.skipLoad(`^static/state_tests/stEIP1559/lowGasPriceOldTypes.json/tests/static/state_tests/stEIP1559/lowGasPriceOldTypesFiller.yml::lowGasPriceOldTypes\[fork_Cancun-blockchain_test_from_state_test-`)
	bt.skipLoad(`^static/state_tests/stEIP1559/lowFeeCap.json/tests/static/state_tests/stEIP1559/lowFeeCapFiller.yml::lowFeeCap\[fork_Cancun-blockchain_test_from_state_test-declaredKeyWrite\]`)
	bt.skipLoad(`^frontier/validation/test_tx_gas_limit.json/tests/frontier/validation/test_transaction.py::test_tx_gas_limit\[fork_Cancun-`)

	// should be resolved
	// note: Executing precompiled contracts with value transferring is not permitted
	// See https://github.com/kaiachain/kaia/blob/d44ae2f4269a84bd379b4e992d8e3be46b7e5ad3/blockchain/vm/evm.go#L267-L270
	bt.skipLoad(`^frontier/opcodes/all_opcodes/all_opcodes.json`)

	// should be resolved
	// note: "to" is not allow to be a precompiled contract address.
	// See https://github.com/kaiachain/kaia/blob/d44ae2f4269a84bd379b4e992d8e3be46b7e5ad3/blockchain/types/tx_internal_data_legacy.go#L365
	bt.skipLoad(`^static/state_tests/stRandom2/randomStatetest642.json/tests/static/state_tests/stRandom2/randomStatetest642Filler.json::randomStatetest642\[.*\]`)
	bt.skipLoad(`^static/state_tests/stRandom2/randomStatetest644.json/tests/static/state_tests/stRandom2/randomStatetest644Filler.json::randomStatetest644\[.*\]`)
	bt.skipLoad(`^static/state_tests/stRandom2/randomStatetest645.json/tests/static/state_tests/stRandom2/randomStatetest645Filler.json::randomStatetest645\[.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/modexpRandomInput.json/tests/static/state_tests/stPreCompiledContracts2/modexpRandomInputFiller.json::modexpRandomInput\[.*\]`)
	bt.skipLoad(`^cancun/eip4844_blobs/test_tx_entry_point.json`)
	bt.skipLoad(`^prague/eip7702_set_code_tx/test_gas_diff_pointer_vs_direct_call.json/tests/prague/eip7702_set_code_tx/test_set_code_txs_2.py::test_gas_diff_pointer_vs_direct_call\[.*-pointer_definition_PointerDefinition`)
	bt.skipLoad(`^osaka/eip7883_modexp_gas_increase/modexp_thresholds/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-exact_gas\]`)
	bt.skipLoad(`^osaka/eip7883_modexp_gas_increase/modexp_thresholds/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-extra_gas\]`)
	bt.skipLoad(`^osaka/eip7883_modexp_gas_increase/modexp_thresholds/modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[fork_Osaka-blockchain_test_from_state_test-insufficient_gas\]`)
	bt.skipLoad(`^osaka/eip7883_modexp_gas_increase/test_modexp_used_in_transaction_entry_points.json/tests/osaka/eip7883_modexp_gas_increase/test_modexp_thresholds.py::test_modexp_used_in_transaction_entry_points\[.*\]`)
	bt.skipLoad(`^osaka/eip7951_p256verify_precompiles/p256verify/precompile_as_tx_entry_point.json`)
	bt.skipLoad(`^osaka/eip7951_p256verify_precompiles/test_precompile_will_return_success_with_tx_value.json`)
	bt.skipLoad(`^osaka/eip7951_p256verify_precompiles/test_precompile_as_tx_entry_point.json`)

	// should be resolved
	// note: `block insertion should have failed` happens. Current tester cannot recognize invalid block header.
	bt.skipLoad(`^frontier/validation/test_gas_limit_below_minimum.json/tests/frontier/validation/test_header.py::test_gas_limit_below_minimum\[fork_Cancun-`)
	bt.skipLoad(`^frontier/validation/test_gas_limit_below_minimum.json/tests/frontier/validation/test_header.py::test_gas_limit_below_minimum\[fork_Prague-`)
	bt.skipLoad(`^frontier/validation/test_gas_limit_below_minimum.json/tests/frontier/validation/test_header.py::test_gas_limit_below_minimum\[fork_Osaka-`)
	bt.skipLoad(`^london/validation/test_invalid_header.json/tests/london/validation/test_header.py::test_invalid_header\[fork_Cancun-blockchain_test-field_base_fee_per_gas-invalid_value_1-exception_BlockException.INVALID_BASEFEE_PER_GAS\]`)
	bt.skipLoad(`^london/validation/test_invalid_header.json/tests/london/validation/test_header.py::test_invalid_header\[fork_Prague-blockchain_test-field_base_fee_per_gas-invalid_value_1-exception_BlockException.INVALID_BASEFEE_PER_GAS\]`)
	bt.skipLoad(`^london/validation/test_invalid_header.json/tests/london/validation/test_header.py::test_invalid_header\[fork_Osaka-blockchain_test-field_base_fee_per_gas-invalid_value_1-exception_BlockException.INVALID_BASEFEE_PER_GAS\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/transactionIntinsicBug_Paris.json/tests/static/state_tests/stEIP1559/transactionIntinsicBug_ParisFiller.yml::transactionIntinsicBug_Paris\[fork_Cancun-blockchain_test_from_state_test-declaredKeyWrite\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/transactionIntinsicBug_Paris.json/tests/static/state_tests/stEIP1559/transactionIntinsicBug_ParisFiller.yml::transactionIntinsicBug_Paris\[fork_Osaka-blockchain_test_from_state_test-declaredKeyWrite\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/transactionIntinsicBug_Paris.json/tests/static/state_tests/stEIP1559/transactionIntinsicBug_ParisFiller.yml::transactionIntinsicBug_Paris\[fork_Prague-blockchain_test_from_state_test-declaredKeyWrite\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/lowGasLimit.json/tests/static/state_tests/stEIP1559/lowGasLimitFiller.yml::lowGasLimit\[fork_Cancun-blockchain_test_from_state_test--g0\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/lowGasLimit.json/tests/static/state_tests/stEIP1559/lowGasLimitFiller.yml::lowGasLimit\[fork_Prague-blockchain_test_from_state_test--g0\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/lowGasLimit.json/tests/static/state_tests/stEIP1559/lowGasLimitFiller.yml::lowGasLimit\[fork_Osaka-blockchain_test_from_state_test--g0\]`)
	bt.skipLoad(`^static/state_tests/stEIP1559/tipTooHigh.json/tests/static/state_tests/stEIP1559/tipTooHighFiller.yml::tipTooHigh\[fork_Cancun-blockchain_test_from_state_test-declaredKeyWrite\]`)

	// should be resolved
	// note: `signal: killed` happens. It may be due to OOM. A lot of opcodes are executed in one transaction.
	bt.skipLoad(`^static/state_tests/stStaticCall/`)

	// should be resolved
	// note: Executing precompiled contracts with value transferring is not permitted.
	// See https://github.com/kaiachain/kaia/blob/a7cb0f8aef2ce813e7ceda38e38353d63437489d/blockchain/vm/evm.go#L281-L284.
	bt.skipLoad(`^static/state_tests/stArgsZeroOneBalance/callNonConst.json/tests/static/state_tests/stArgsZeroOneBalance/callNonConstFiller.yml::callNonConst\[fork_Cancun-blockchain_test_from_state_test--v1\]`)
	bt.skipLoad(`^static/state_tests/stArgsZeroOneBalance/callNonConst.json/tests/static/state_tests/stArgsZeroOneBalance/callNonConstFiller.yml::callNonConst\[fork_Osaka-blockchain_test_from_state_test--v1\]`)
	bt.skipLoad(`^static/state_tests/stArgsZeroOneBalance/callNonConst.json/tests/static/state_tests/stArgsZeroOneBalance/callNonConstFiller.yml::callNonConst\[fork_Prague-blockchain_test_from_state_test--v1\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create-high-nonce-v1\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Cancun-blockchain_test_from_state_test-new.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Osaka-blockchain_test_from_state_test-all_then_yes_from_prague.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Osaka-blockchain_test_from_state_test-new.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Osaka-blockchain_test_from_state_test-yes.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Prague-blockchain_test_from_state_test-all_then_yes_from_prague.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Prague-blockchain_test_from_state_test-new.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts/precompsEIP2929Cancun.json/tests/static/state_tests/stPreCompiledContracts/precompsEIP2929CancunFiller.yml::precompsEIP2929Cancun\[fork_Prague-blockchain_test_from_state_test-yes.*\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGas.json/tests/static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGasFiller.json::CallEcrecover0_NoGas\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGas.json/tests/static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGasFiller.json::CallEcrecover0_NoGas\[fork_Osaka-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGas.json/tests/static/state_tests/stPreCompiledContracts2/CallEcrecover0_NoGasFiller.json::CallEcrecover0_NoGas\[fork_Prague-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValue.json/tests/static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValueFiller.json::CallSha256_1_nonzeroValue\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValue.json/tests/static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValueFiller.json::CallSha256_1_nonzeroValue\[fork_Osaka-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValue.json/tests/static/state_tests/stPreCompiledContracts2/CallSha256_1_nonzeroValueFiller.json::CallSha256_1_nonzeroValue\[fork_Prague-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stRandom2/randomStatetest650.json/tests/static/state_tests/stRandom2/randomStatetest650Filler.json::randomStatetest650\[fork_Cancun-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stRandom2/randomStatetest650.json/tests/static/state_tests/stRandom2/randomStatetest650Filler.json::randomStatetest650\[fork_Prague-blockchain_test_from_state_test-\]`)
	bt.skipLoad(`^static/state_tests/stCreateTest/CreateAddressWarmAfterFail.json/tests/static/state_tests/stCreateTest/CreateAddressWarmAfterFailFiller.yml::CreateAddressWarmAfterFail\[fork_Cancun-blockchain_test_from_state_test-create-high-nonce-v0\]`)
	bt.skipLoad(`^frontier/identity_precompile/test_call_identity_precompile.json/tests/frontier/identity_precompile/test_identity.py::test_call_identity_precompile\[fork_Cancun-blockchain_test_from_state_test-identity_1_nonzerovalue-call_type_CALL\]`)
	bt.skipLoad(`^frontier/identity_precompile/test_call_identity_precompile.json/tests/frontier/identity_precompile/test_identity.py::test_call_identity_precompile\[fork_Osaka-blockchain_test_from_state_test-identity_1_nonzerovalue-call_type_CALL\]`)
	bt.skipLoad(`^frontier/identity_precompile/test_call_identity_precompile.json/tests/frontier/identity_precompile/test_identity.py::test_call_identity_precompile\[fork_Prague-blockchain_test_from_state_test-identity_1_nonzerovalue-call_type_CALL\]`)
	bt.skipLoad(`^frontier/opcodes/test_all_opcodes.json/tests/frontier/opcodes/test_all_opcodes.py::test_all_opcodes\[fork_Cancun-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^frontier/opcodes/test_all_opcodes.json/tests/frontier/opcodes/test_all_opcodes.py::test_all_opcodes\[fork_Osaka-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^frontier/opcodes/test_all_opcodes.json/tests/frontier/opcodes/test_all_opcodes.py::test_all_opcodes\[fork_Prague-blockchain_test_from_state_test\]`)
	bt.skipLoad(`^frontier/precompiles/precompiles/precompiles.json/tests/frontier/precompiles/test_precompiles.py::test_precompiles\[fork_Cancun-address_0xb-precompile_exists_False-blockchain_test_from_state_test\]`)

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
			"OsakaToBPO1AtTime15k",
			"BPO1ToBPO2AtTime15k",
			"BPO2ToBPO3AtTime15k",
			"BPO3ToBPO4AtTime15k",
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
