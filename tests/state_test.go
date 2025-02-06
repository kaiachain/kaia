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
// This file is derived from tests/state_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package tests

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/suite"
)

// TestKaiaSpecState runs the StateTests fixtures from kaia-core-tests
func TestKaiaSpecState(t *testing.T) {
	t.Parallel()

	st := new(testMatcher)
	// Long tests:
	st.skipShortMode(`^stQuadraticComplexityTest/`)
	// Broken tests:
	st.skipLoad(`^stTransactionTest/OverflowGasRequire\.json`) // gasLimit > 256 bits
	st.skipLoad(`^stTransactionTest/zeroSigTransa[^/]*\.json`) // EIP-86 is not supported yet
	// Expected failures:
	st.fails(`^stRevertTest/RevertPrecompiledTouch\.json/Byzantium`, "bug in test")
	st.skipLoad(`^stZeroKnowledge2/ecmul_0-3_5616_28000_96\.json`)

	// Skip since the tests transfer values to precompiled contracts
	st.skipLoad(`^stPreCompiledContracts2/CallSha256_1_nonzeroValue.json`)
	st.skipLoad(`^stPreCompiledContracts2/CallIdentity_1_nonzeroValue.json`)
	st.skipLoad(`^stPreCompiledContracts2/CallEcrecover0_NoGas.json`)
	st.skipLoad(`^stRandom2/randomStatetest644.json`)
	st.skipLoad(`^stRandom2/randomStatetest645.json`)
	st.skipLoad(`^stStaticCall/static_CallIdentity_1_nonzeroValue.json`)
	st.skipLoad(`^stStaticCall/static_CallSha256_1_nonzeroValue.json`)
	st.skipLoad(`^stArgsZeroOneBalance/callNonConst.json`)
	st.skipLoad(`^stPreCompiledContracts2/modexpRandomInput.json`)
	st.skipLoad(`^stRandom2/randomStatetest642.json`)

	st.walk(t, stateTestDir, func(t *testing.T, name string, test *StateTest) {
		execStateTest(t, st, test, name, []string{"Constantinople"}, false)
	})
}

// TestExecutionSpecState runs the state_test fixtures from execution-spec-tests.

type ExecutionSpecStateTestSuite struct {
	suite.Suite
}

func (suite *ExecutionSpecStateTestSuite) SetupSuite() {
	vm.RelaxPrecompileRangeForTest(true)
}

func (suite *ExecutionSpecStateTestSuite) TearDownSuite() {
	vm.RelaxPrecompileRangeForTest(false)
}

func (suite *ExecutionSpecStateTestSuite) TestExecutionSpecState() {
	t := suite.T()

	if !common.FileExist(executionSpecStateTestDir) {
		t.Skipf("directory %s does not exist", executionSpecStateTestDir)
	}
	st := new(testMatcher)

	// TODO-Kaia: should remove these skip
	// executing precompiled contracts with value transferring is not permitted
	st.skipLoad(`^frontier\/opcodes\/all_opcodes\/all_opcodes.json`)

	// tests to skip
	// unsupported EIPs
	st.skipLoad(`^cancun\/eip4788_beacon_root\/`)
	st.skipLoad(`^cancun\/eip4844_blobs\/`)
	// different amount of gas is consumed because 0x0b contract is added to access list by ActivePrecompiles although Cancun doesn't have it as a precompiled contract
	st.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-state_test-31_bytes\]`)
	st.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-state_test-32_bytes\]`)
	st.skipLoad(`^frontier\/precompiles\/precompile_absence\/precompile_absence.json\/tests\/frontier\/precompiles\/test_precompile_absence.py::test_precompile_absence\[fork_Cancun-state_test-empty_calldata\]`)
	st.skipLoad(`^prague\/eip2537_bls_12_381_precompiles\/bls12_precompiles_before_fork\/precompile_before_fork.json\/tests\/prague\/eip2537_bls_12_381_precompiles\/test_bls12_precompiles_before_fork.py::test_precompile_before_fork\[fork_CancunToPragueAtTime15k-state_test--G1ADD\]`)
	// type 3 tx (EIP-4844) is not supported
	st.skipLoad(`^prague\/eip7623_increase_calldata_cost\/.*type_3.*`)

	st.walk(t, executionSpecStateTestDir, func(t *testing.T, name string, test *StateTest) {
		execStateTest(t, st, test, name, []string{
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
			// "Cancun",
			// "Prague",
		}, true)
	})
}

func TestExecutionSpecStateTestSuite(t *testing.T) {
	suite.Run(t, new(ExecutionSpecStateTestSuite))
}

func execStateTest(t *testing.T, st *testMatcher, test *StateTest, name string, skipForks []string, isTestExecutionSpecState bool) {
	for _, subtest := range test.Subtests() {
		subtest := subtest
		key := fmt.Sprintf("%s/%d", subtest.Fork, subtest.Index)
		name := name + "/" + key
		t.Run(key, func(t *testing.T) {
			for _, skip := range skipForks {
				if skip == subtest.Fork {
					t.Skipf("%s not supported yet", subtest.Fork)
				}
			}
			withTrace(t, test.gasLimit(subtest), func(vmconfig vm.Config) error {
				err := test.Run(subtest, vmconfig, isTestExecutionSpecState)
				return st.checkFailure(t, name, err)
			})
		})
	}
}

// Transactions with gasLimit above this value will not get a VM trace on failure.
const traceErrorLimit = 400000

func withTrace(t *testing.T, gasLimit uint64, test func(vm.Config) error) {
	// Set ComputationCostLimit as infinite
	err := test(vm.Config{ComputationCostLimit: params.OpcodeComputationCostLimitInfinite})
	if err == nil {
		return
	}
	t.Error(err)
	if gasLimit > traceErrorLimit {
		t.Log("gas limit too high for EVM trace")
		return
	}
	tracer := vm.NewStructLogger(nil)
	err2 := test(vm.Config{Debug: true, Tracer: tracer, ComputationCostLimit: params.OpcodeComputationCostLimitInfinite})
	if !reflect.DeepEqual(err, err2) {
		t.Errorf("different error for second run: %v", err2)
	}
	buf := new(bytes.Buffer)
	vm.WriteTrace(buf, tracer.StructLogs())
	if buf.Len() == 0 {
		t.Log("no EVM operation logs generated")
	} else {
		t.Log("EVM operation log:\n" + buf.String())
	}
	t.Logf("EVM output: 0x%x", tracer.Output())
	t.Logf("EVM error: %v", tracer.Error())
}
