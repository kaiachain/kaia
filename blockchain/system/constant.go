// Modifications Copyright 2024 The Kaia Authors
// Copyright 2023 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.
// Modified and improved for the Kaia development.

package system

import (
	"errors"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	kip113contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip113"
	kip149contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip149"
	misccontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/misc"
	"github.com/kaiachain/kaia/contracts/contracts/system_contracts/multicall"
	proxycontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/proxy"
	"github.com/kaiachain/kaia/contracts/contracts/testing/reward"
	testcontract "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/log"
)

var (
	logger = log.NewModuleLogger(log.Blockchain)

	// Canonical system contract names registered in Registry.
	AddressBookName = "AddressBook"
	GovParamName    = "GovParam"
	Kip103Name      = "KIP103"
	Kip113Name      = "KIP113"
	Kip160Name      = "KIP160"
	CLRegistryName  = "CLRegistry"

	AllContractNames = []string{
		AddressBookName,
		GovParamName,
		Kip103Name,
		Kip113Name,
		Kip160Name,
	}

	// This is the keccak-256 hash of "eip1967.proxy.implementation" subtracted by 1 used in the
	// EIP-1967 proxy contract. See https://eips.ethereum.org/EIPS/eip-1967#implementation-slot
	ImplementationSlot = common.Hex2Bytes("360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")

	// Some system contracts are allocated at special addresses.
	MainnetCreditAddr = common.HexToAddress("0x0000000000000000000000000000000000000000")
	AddressBookAddr   = common.HexToAddress("0x0000000000000000000000000000000000000400")
	RegistryAddr      = common.HexToAddress("0x0000000000000000000000000000000000000401")
	MultiCallAddr     = common.HexToAddress("0x0000000000000000000000000000000000000402")
	// The following addresses are only used for testing.
	Kip113ProxyAddrMock = common.HexToAddress("0x0000000000000000000000000000000000000402")
	Kip113LogicAddrMock = common.HexToAddress("0x0000000000000000000000000000000000000403")

	// Registry will return zero address for non-existent system contract.
	NonExistentAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

	// System contract binaries to be injected at hardfork or used in testing.
	MainnetCreditCode   = hexutil.MustDecode("0x" + misccontract.CypressCreditBinRuntime)
	MainnetCreditV2Code = hexutil.MustDecode("0x" + misccontract.CypressCreditV2BinRuntime)
	RegistryCode        = hexutil.MustDecode("0x" + kip149contract.RegistryBinRuntime)
	RegistryMockCode    = hexutil.MustDecode("0x" + testcontract.RegistryMockBinRuntime)
	Kip160MockCode      = hexutil.MustDecode("0x" + testcontract.TreasuryRebalanceMockV2BinRuntime)
	Kip103MockCode      = hexutil.MustDecode("0x" + testcontract.TreasuryRebalanceMockBinRuntime)
	Kip113Code          = hexutil.MustDecode("0x" + kip113contract.SimpleBlsRegistryBinRuntime)
	Kip113MockCode      = hexutil.MustDecode("0x" + testcontract.KIP113MockBinRuntime)

	ERC1967ProxyCode = hexutil.MustDecode("0x" + proxycontract.ERC1967ProxyBinRuntime)

	AddressBookMockTwoCNCode = hexutil.MustDecode("0x" + reward.AddressBookMockTwoCNBinRuntime)
	Kip113MockThreeCNCode    = hexutil.MustDecode("0x" + testcontract.KIP113MockThreeCNBinRuntime)

	MultiCallCode     = hexutil.MustDecode("0x" + multicall.MultiCallContractBinRuntime)
	MultiCallMockCode = hexutil.MustDecode("0x" + testcontract.MultiCallContractMockBinRuntime)

	// Mock for CLRegistry testing
	CLRegistryMockThreeCLAddr = common.HexToAddress("0x0000000000000000000000000000000000000Ff0")
	WrappedKaiaMockAddr       = common.HexToAddress("0x0000000000000000000000000000000000000Ff1")
	RegistryMockForCLCode     = hexutil.MustDecode("0x" + reward.RegistryMockForCLBinRuntime)
	RegistryMockZero          = hexutil.MustDecode("0x" + reward.RegistryMockZeroBinRuntime)
	CLRegistryMockThreeCLCode = hexutil.MustDecode("0x" + reward.CLRegistryMockThreeCLBinRuntime)
	WrappedKaiaMockCode       = hexutil.MustDecode("0x" + reward.WrappedKaiaMockBinRuntime)

	// Errors
	ErrRegistryNotInstalled      = errors.New("Registry contract not installed")
	ErrRebalanceIncorrectBlock   = errors.New("cannot find a proper target block number")
	ErrRebalanceNotEnoughBalance = errors.New("the sum of zeroed balances are less than the sum of allocated balances")
	ErrRebalanceBadStatus        = errors.New("rebalance contract is not in proper status")
	ErrKip113BadResult           = errors.New("KIP113 call returned bad data")
	ErrKip113BadPop              = errors.New("KIP113 PoP verification failed")
)

// Solidity enums are not generated by abigen. Add them manually.
const (
	EnumRebalanceStatus_Initialized uint8 = iota
	EnumRebalanceStatus_Registered
	EnumRebalanceStatus_Approved
	EnumRebalanceStatus_Finalized
)
