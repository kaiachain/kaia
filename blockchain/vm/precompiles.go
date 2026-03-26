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

package vm

// Precompile routing in Kaia
//
// There are three distinct "views" of the precompiled contract set, each serving
// a different purpose:
//
//  1. common.IsPrecompiledContractAddress (common/types.go)
//     A fast range-check (0x0001–0x03FF) used for tx validation and the
//     evm.Call() guard. It does NOT reflect the exact active set for a given
//     fork or vmversion.
//
//  2. ActivePrecompiles / ActivePrecompiledContracts (this file)
//     Returns the address list / map for the current fork rules, used when
//     constructing the EIP-2929 access list. After Istanbul, the vmversion0
//     compat addresses (0x0a, 0x0b) are appended so that old contracts are
//     not charged cold-access gas for precompiles they can legitimately call.
//
//  3. EVM.GetPrecompiledContractMap / getPrecompiledContractForVersion (this file)
//     Returns the exact map used during EVM execution. The map is keyed on
//     the caller's address, not the precompile's: if the caller was deployed
//     before Istanbul (VmVersion0), it always receives the Byzantium map so
//     that its references to 0x09–0x0b (vmLog, feePayer, validateSender) remain
//     valid. Post-Istanbul callers get the fork-appropriate map where 0x09 is
//     blake2F and those three contracts live at 0x3fd–0x3ff.

import (
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/params"
)

// Fork-specific precompile maps.
// See contracts.go for the PrecompiledContract interface and all
// Run / GetRequiredGasAndComputationCost implementations.

// PrecompiledContractsByzantium contains the default set of pre-compiled Kaia
// contracts based on Ethereum Byzantium.
var PrecompiledContractsByzantium = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):  &ecrecover{},
	common.BytesToAddress([]byte{2}):  &sha256hash{},
	common.BytesToAddress([]byte{3}):  &ripemd160hash{},
	common.BytesToAddress([]byte{4}):  &dataCopy{},
	common.BytesToAddress([]byte{5}):  &bigModExp{eip2565: false, eip7823: false, eip7883: false},
	common.BytesToAddress([]byte{6}):  &bn256AddByzantium{},
	common.BytesToAddress([]byte{7}):  &bn256ScalarMulByzantium{},
	common.BytesToAddress([]byte{8}):  &bn256PairingByzantium{},
	common.BytesToAddress([]byte{9}):  &vmLog{},
	common.BytesToAddress([]byte{10}): &feePayer{},
	common.BytesToAddress([]byte{11}): &validateSender{},
}

// DO NOT USE 0x3FD, 0x3FE, 0x3FF ADDRESSES BEFORE ISTANBUL CHANGE ACTIVATED.

// PrecompiledContractsIstanbul contains the default set of pre-compiled Kaia
// contracts based on Ethereum Istanbul.
var PrecompiledContractsIstanbul = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):      &ecrecover{},
	common.BytesToAddress([]byte{2}):      &sha256hash{},
	common.BytesToAddress([]byte{3}):      &ripemd160hash{},
	common.BytesToAddress([]byte{4}):      &dataCopy{},
	common.BytesToAddress([]byte{5}):      &bigModExp{eip2565: false, eip7823: false, eip7883: false},
	common.BytesToAddress([]byte{6}):      &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):      &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):      &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):      &blake2F{},
	common.BytesToAddress([]byte{3, 253}): &vmLog{},
	common.BytesToAddress([]byte{3, 254}): &feePayer{},
	common.BytesToAddress([]byte{3, 255}): &validateSender{},
}

// PrecompiledContractsKore contains the default set of pre-compiled Kaia
// contracts based on Ethereum Berlin.
var PrecompiledContractsKore = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):      &ecrecover{},
	common.BytesToAddress([]byte{2}):      &sha256hash{},
	common.BytesToAddress([]byte{3}):      &ripemd160hash{},
	common.BytesToAddress([]byte{4}):      &dataCopy{},
	common.BytesToAddress([]byte{5}):      &bigModExp{eip2565: true, eip7823: false, eip7883: false},
	common.BytesToAddress([]byte{6}):      &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):      &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):      &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):      &blake2F{},
	common.BytesToAddress([]byte{3, 253}): &vmLog{},
	common.BytesToAddress([]byte{3, 254}): &feePayer{},
	common.BytesToAddress([]byte{3, 255}): &validateSender{},
}

// PrecompiledContractsCancun contains the default set of pre-compiled Kaia
// contracts based on Ethereum Cancun.
var PrecompiledContractsCancun = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):      &ecrecover{},
	common.BytesToAddress([]byte{2}):      &sha256hash{},
	common.BytesToAddress([]byte{3}):      &ripemd160hash{},
	common.BytesToAddress([]byte{4}):      &dataCopy{},
	common.BytesToAddress([]byte{5}):      &bigModExp{eip2565: true, eip7823: false, eip7883: false},
	common.BytesToAddress([]byte{6}):      &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):      &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):      &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):      &blake2F{},
	common.BytesToAddress([]byte{0x0a}):   &kzgPointEvaluation{},
	common.BytesToAddress([]byte{3, 253}): &vmLog{},
	common.BytesToAddress([]byte{3, 254}): &feePayer{},
	common.BytesToAddress([]byte{3, 255}): &validateSender{},
}

// PrecompiledContractsPrague contains the set of pre-compiled Ethereum
// contracts used in the Prague release.
var PrecompiledContractsPrague = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):      &ecrecover{},
	common.BytesToAddress([]byte{2}):      &sha256hash{},
	common.BytesToAddress([]byte{3}):      &ripemd160hash{},
	common.BytesToAddress([]byte{4}):      &dataCopy{},
	common.BytesToAddress([]byte{5}):      &bigModExp{eip2565: true, eip7823: false, eip7883: false},
	common.BytesToAddress([]byte{6}):      &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):      &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):      &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):      &blake2F{},
	common.BytesToAddress([]byte{0x0a}):   &kzgPointEvaluation{},
	common.BytesToAddress([]byte{0x0b}):   &bls12381G1Add{},
	common.BytesToAddress([]byte{0x0c}):   &bls12381G1MultiExp{},
	common.BytesToAddress([]byte{0x0d}):   &bls12381G2Add{},
	common.BytesToAddress([]byte{0x0e}):   &bls12381G2MultiExp{},
	common.BytesToAddress([]byte{0x0f}):   &bls12381Pairing{},
	common.BytesToAddress([]byte{0x10}):   &bls12381MapG1{},
	common.BytesToAddress([]byte{0x11}):   &bls12381MapG2{},
	common.BytesToAddress([]byte{3, 253}): &vmLog{},
	common.BytesToAddress([]byte{3, 254}): &feePayer{},
	common.BytesToAddress([]byte{3, 255}): &validateSender{},
}

// PrecompiledContractsOsaka contains the set of pre-compiled Ethereum
// contracts used in the Osaka release.
var PrecompiledContractsOsaka = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}):         &ecrecover{},
	common.BytesToAddress([]byte{2}):         &sha256hash{},
	common.BytesToAddress([]byte{3}):         &ripemd160hash{},
	common.BytesToAddress([]byte{4}):         &dataCopy{},
	common.BytesToAddress([]byte{5}):         &bigModExp{eip2565: true, eip7823: true, eip7883: true},
	common.BytesToAddress([]byte{6}):         &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}):         &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}):         &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}):         &blake2F{},
	common.BytesToAddress([]byte{0x0a}):      &kzgPointEvaluation{},
	common.BytesToAddress([]byte{0x0b}):      &bls12381G1Add{},
	common.BytesToAddress([]byte{0x0c}):      &bls12381G1MultiExp{},
	common.BytesToAddress([]byte{0x0d}):      &bls12381G2Add{},
	common.BytesToAddress([]byte{0x0e}):      &bls12381G2MultiExp{},
	common.BytesToAddress([]byte{0x0f}):      &bls12381Pairing{},
	common.BytesToAddress([]byte{0x10}):      &bls12381MapG1{},
	common.BytesToAddress([]byte{0x11}):      &bls12381MapG2{},
	common.BytesToAddress([]byte{0x1, 0x00}): &p256Verify{},
	common.BytesToAddress([]byte{3, 253}):    &vmLog{},
	common.BytesToAddress([]byte{3, 254}):    &feePayer{},
	common.BytesToAddress([]byte{3, 255}):    &validateSender{},
}

// PrecompiledContractsP256Verify contains the precompiled Ethereum
// contract specified in EIP-7212. This is exported for testing purposes.
var PrecompiledContractsP256Verify = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{0x1, 0x00}): &p256Verify{},
}

var (
	PrecompiledAddressOsaka       []common.Address
	PrecompiledAddressPrague      []common.Address
	PrecompiledAddressCancun      []common.Address
	PrecompiledAddressIstanbul    []common.Address
	PrecompiledAddressesByzantium []common.Address
)

func init() {
	for k := range PrecompiledContractsByzantium {
		PrecompiledAddressesByzantium = append(PrecompiledAddressesByzantium, k)
	}
	for k := range PrecompiledContractsIstanbul {
		PrecompiledAddressIstanbul = append(PrecompiledAddressIstanbul, k)
	}
	for k := range PrecompiledContractsCancun {
		PrecompiledAddressCancun = append(PrecompiledAddressCancun, k)
	}
	for k := range PrecompiledContractsPrague {
		PrecompiledAddressPrague = append(PrecompiledAddressPrague, k)
	}
	for k := range PrecompiledContractsOsaka {
		PrecompiledAddressOsaka = append(PrecompiledAddressOsaka, k)
	}
}

// ActivePrecompiles returns the precompiles enabled with the current configuration.
// This is a variable so that it can be overridden in tests (e.g. EEST) to exclude
// the vmversion0 compatibility addresses from the access list.
var ActivePrecompiles = func(rules params.Rules) []common.Address {
	var precompiledContractAddrs []common.Address
	for addr := range ActivePrecompiledContracts(rules) {
		precompiledContractAddrs = append(precompiledContractAddrs, addr)
	}
	// After istanbulCompatible hf, need to support for vmversion0 contracts, too.
	// VmVersion0 contracts are deployed before istanbulCompatible and they use byzantiumCompatible precompiled contracts.
	// VmVersion0 contracts are the contracts deployed before istanbulCompatible hf.
	if rules.IsIstanbul {
		return append(precompiledContractAddrs,
			[]common.Address{common.BytesToAddress([]byte{10}), common.BytesToAddress([]byte{11})}...)
	} else {
		return precompiledContractAddrs
	}
}

// ActivePrecompiledContracts returns the precompiled contracts enabled with the current configuration.
// This function doesn't support for vmversion0 contracts, it only supports for istanbulCompatible hf.
func ActivePrecompiledContracts(rules params.Rules) map[common.Address]PrecompiledContract {
	switch {
	case rules.IsOsaka:
		return PrecompiledContractsOsaka
	case rules.IsPrague:
		return PrecompiledContractsPrague
	case rules.IsCancun:
		return PrecompiledContractsCancun
	case rules.IsIstanbul:
		return PrecompiledContractsIstanbul
	default:
		return PrecompiledContractsByzantium
	}
}

// GetPrecompiledContractMap returns the precompiled contract map for the given
// caller address, accounting for vmversion and the active fork rules.
// If console.log is enabled in the EVM config, the consoleLog precompile is
// included as well.
func (evm *EVM) GetPrecompiledContractMap(addr common.Address) map[common.Address]PrecompiledContract {
	precompiles := evm.getPrecompiledContractForVersion(addr)
	// If console.log is enabled, add console.log precompile too.
	if evm.Config.UseConsoleLog {
		precompiles[consoleLogContractAddress] = &consoleLog{}
	}
	return precompiles
}

func (evm *EVM) getPrecompiledContractForVersion(addr common.Address) map[common.Address]PrecompiledContract {
	// VmVersion means that the contract uses the precompiled contract map at the deployment time.
	// Also, it follows old map's gas price & computation cost.

	// Get vmVersion from addr only if the addr is a contract address.
	// If new "VmVersion" is added, add new if clause below
	if vmVersion, ok := evm.StateDB.GetVmVersion(addr); ok && vmVersion == params.VmVersion0 {
		// Without VmVersion0, precompiled contract address 0x09-0x0b won't work properly
		// with the contracts deployed before istanbulHF
		return PrecompiledContractsByzantium
	}

	switch {
	case evm.chainRules.IsOsaka:
		return PrecompiledContractsOsaka
	case evm.chainRules.IsPrague:
		return PrecompiledContractsPrague
	case evm.chainRules.IsCancun:
		return PrecompiledContractsCancun
	case evm.chainRules.IsKore:
		return PrecompiledContractsKore
	case evm.chainRules.IsIstanbul:
		return PrecompiledContractsIstanbul
	default:
		return PrecompiledContractsByzantium
	}
}
