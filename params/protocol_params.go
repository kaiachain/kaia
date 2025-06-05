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
// This file is derived from params/protocol_params.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package params

import (
	"fmt"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/v2/common"
)

var TargetGasLimit = GenesisGasLimit // The artificial target

const (
	// Fee schedule parameters

	CallValueTransferGas  uint64 = 9000  // Paid for CALL when the value transfer is non-zero.                  // G_callvalue
	CallNewAccountGas     uint64 = 25000 // Paid for CALL when the destination address didn't exist prior.      // G_newaccount
	TxGas                 uint64 = 21000 // Per transaction not creating a contract. NOTE: Not payable on data of calls between transactions. // G_transaction
	TxGasContractCreation uint64 = 53000 // Per transaction that creates a contract. NOTE: Not payable on data of calls between transactions. // G_transaction + G_create
	TxDataZeroGas         uint64 = 4     // Per byte of data attached to a transaction that equals zero. NOTE: Not payable on data of calls between transactions. // G_txdatazero
	TxTokenPerNonZeroByte uint64 = 4     // Token cost per non-zero byte as specified by EIP-7623.
	TxCostFloorPerToken   uint64 = 10    // Cost floor per byte of data as specified by EIP-7623.
	QuadCoeffDiv          uint64 = 512   // Divisor for the quadratic particle of the memory cost equation.
	SstoreSetGas          uint64 = 20000 // Once per SLOAD operation.                                           // G_sset
	LogDataGas            uint64 = 8     // Per byte in a LOG* operation's data.                                // G_logdata
	CallStipend           uint64 = 2300  // Free gas given at beginning of call.                                // G_callstipend
	Sha3Gas               uint64 = 30    // Once per SHA3 operation.                                                 // G_sha3
	Sha3WordGas           uint64 = 6     // Once per word of the SHA3 operation's data.                              // G_sha3word
	InitCodeWordGas       uint64 = 2     // Once per word of the init code when creating a contract.				 // G_InitCodeWord
	SstoreResetGas        uint64 = 5000  // Once per SSTORE operation if the zeroness changes from zero.             // G_sreset
	SstoreClearGas        uint64 = 5000  // Once per SSTORE operation if the zeroness doesn't change.                // G_sreset
	SstoreRefundGas       uint64 = 15000 // Once per SSTORE operation if the zeroness changes to zero.               // R_sclear

	// gasSStoreEIP2200
	SstoreSentryGasEIP2200            uint64 = 2300  // Minimum gas required to be present for an SSTORE call, not consumed
	SstoreSetGasEIP2200               uint64 = 20000 // Once per SSTORE operation from clean zero to non-zero
	SstoreResetGasEIP2200             uint64 = 5000  // Once per SSTORE operation from clean non-zero to something else
	SstoreClearsScheduleRefundEIP2200 uint64 = 15000 // Once per SSTORE operation for clearing an originally existing storage slot

	ColdAccountAccessCostEIP2929 = uint64(2600) // COLD_ACCOUNT_ACCESS_COST
	ColdSloadCostEIP2929         = uint64(2100) // COLD_SLOAD_COST
	WarmStorageReadCostEIP2929   = uint64(100)  // WARM_STORAGE_READ_COST

	// In EIP-2200: SstoreResetGas was 5000.
	// In EIP-2929: SstoreResetGas was changed to '5000 - COLD_SLOAD_COST'.
	// In EIP-3529: SSTORE_CLEARS_SCHEDULE is defined as SSTORE_RESET_GAS + ACCESS_LIST_STORAGE_KEY_COST
	// Which becomes: 5000 - 2100 + 1900 = 4800
	SstoreClearsScheduleRefundEIP3529 uint64 = SstoreResetGasEIP2200 - ColdSloadCostEIP2929 + TxAccessListStorageKeyGas

	JumpdestGas           uint64 = 1     // Once per JUMPDEST operation.
	CreateDataGas         uint64 = 200   // Paid per byte for a CREATE operation to succeed in placing code into state. // G_codedeposit
	ExpGas                uint64 = 10    // Once per EXP instruction
	LogGas                uint64 = 375   // Per LOG* operation.                                                          // G_log
	CopyGas               uint64 = 3     // Partial payment for COPY operations, multiplied by words copied, rounded up. // G_copy
	CreateGas             uint64 = 32000 // Once per CREATE operation & contract-creation transaction.               // G_create
	Create2Gas            uint64 = 32000 // Once per CREATE2 operation
	SelfdestructRefundGas uint64 = 24000 // Refunded following a selfdestruct operation.                                  // R_selfdestruct
	MemoryGas             uint64 = 3     // Times the address of the (highest referenced byte in memory + 1). NOTE: referencing happens on read, write and in instructions such as RETURN and CALL. // G_memory
	LogTopicGas           uint64 = 375   // Multiplied by the * of the LOG*, per LOG transaction. e.g. LOG0 incurs 0 * c_txLogTopicGas, LOG4 incurs 4 * c_txLogTopicGas.   // G_logtopic

	TxDataNonZeroGasFrontier uint64 = 68 // Per byte of data attached to a transaction that is not equal to zero. NOTE: Not payable on data of calls between transactions.
	TxDataNonZeroGasEIP2028  uint64 = 16 // Per byte of non zero data attached to a transaction after EIP 2028 (part in Istanbul in Geth, part in Prague in Kaia)

	CallGas         uint64 = 700  // Static portion of gas for CALL-derivates after EIP 150 (Tangerine)
	ExtcodeSizeGas  uint64 = 700  // Cost of EXTCODESIZE after EIP 150 (Tangerine)
	SelfdestructGas uint64 = 5000 // Cost of SELFDESTRUCT post EIP 150 (Tangerine)

	// Istanbul version of BalanceGas, SloadGas, ExtcodeHash is added.
	BalanceGasEIP150             uint64 = 400 // Cost of BALANCE     before EIP 1884
	BalanceGasEIP1884            uint64 = 700 // Cost of BALANCE     after  EIP 1884 (part of Istanbul)
	SloadGasEIP150               uint64 = 200 // Cost of SLOAD       before EIP 1884
	SloadGasEIP1884              uint64 = 800 // Cost of SLOAD       after  EIP 1884 (part of Istanbul)
	SloadGasEIP2200              uint64 = 800 // Cost of SLOAD       after  EIP 2200 (part of Istanbul)
	ExtcodeHashGasConstantinople uint64 = 400 // Cost of EXTCODEHASH before EIP 1884
	ExtcodeHashGasEIP1884        uint64 = 700 // Cost of EXTCODEHASH after  EIP 1884 (part in Istanbul)

	// EXP has a dynamic portion depending on the size of the exponent
	// was set to 10 in Frontier, was raised to 50 during Eip158 (Spurious Dragon)
	ExpByte uint64 = 50

	// Extcodecopy has a dynamic AND a static cost. This represents only the
	// static portion of the gas. It was changed during EIP 150 (Tangerine)
	ExtcodeCopyBase uint64 = 700

	// CreateBySelfdestructGas is used when the refunded account is one that does
	// not exist. This logic is similar to call.
	// Introduced in Tangerine Whistle (Eip 150)
	CreateBySelfdestructGas uint64 = 25000

	// Fee for Service Chain
	// TODO-Kaia-ServiceChain The following parameters should be fixed.
	// TODO-Kaia-Governance The following parameters should be able to be modified by governance.
	TxChainDataAnchoringGas uint64 = 21000 // Per transaction anchoring chain data. NOTE: Not payable on data of calls between transactions. // G_transactionchaindataanchoring
	ChainDataAnchoringGas   uint64 = 100   // Per byte of anchoring chain data NOTE: Not payable on data of calls between transactions. // G_chaindataanchoring

	// Precompiled contract gas prices

	EcrecoverGas        uint64 = 3000 // Elliptic curve sender recovery gas price
	Sha256BaseGas       uint64 = 60   // Base price for a SHA256 operation
	Sha256PerWordGas    uint64 = 12   // Per-word price for a SHA256 operation
	Ripemd160BaseGas    uint64 = 600  // Base price for a RIPEMD160 operation
	Ripemd160PerWordGas uint64 = 120  // Per-word price for a RIPEMD160 operation
	IdentityBaseGas     uint64 = 15   // Base price for a data copy operation
	IdentityPerWordGas  uint64 = 3    // Per-work price for a data copy operation

	Bn256AddGasByzantium               uint64 = 500    // Gas needed for an elliptic curve addition
	Bn256AddGasIstanbul                uint64 = 150    // Istanbul version of gas needed for an elliptic curve addition
	Bn256ScalarMulGasByzantium         uint64 = 40000  // Gas needed for an elliptic curve scalar multiplication
	Bn256ScalarMulGasIstanbul          uint64 = 6000   // Istanbul version of gas needed for an elliptic curve scalar multiplication
	Bn256PairingBaseGasByzantium       uint64 = 100000 // Base price for an elliptic curve pairing check
	Bn256PairingBaseGasIstanbul        uint64 = 45000  // Istanbul version of base price for an elliptic curve pairing check
	Bn256PairingPerPointGasByzantium   uint64 = 80000  // Per-point price for an elliptic curve pairing check
	Bn256PairingPerPointGasIstanbul    uint64 = 34000  // Istanbul version of per-point price for an elliptic curve pairing check
	BlobTxPointEvaluationPrecompileGas uint64 = 50000  // Gas price for the point evaluation precompile.
	VMLogBaseGas                       uint64 = 100    // Base price for a VMLOG operation
	VMLogPerByteGas                    uint64 = 20     // Per-byte price for a VMLOG operation
	FeePayerGas                        uint64 = 300    // Gas needed for calculating the fee payer of the transaction in a smart contract.
	ValidateSenderGas                  uint64 = 5000   // Gas needed for validating the signature of a message.

	// The Refund Quotient is the cap on how much of the used gas can be refunded. Before EIP-3529,
	// up to half the consumed gas could be refunded. Redefined as 1/5th in EIP-3529
	RefundQuotient        uint64 = 2
	RefundQuotientEIP3529 uint64 = 5

	GasLimitBoundDivisor uint64 = 1024    // The bound divisor of the gas limit, used in update calculations.
	MinGasLimit          uint64 = 5000    // Minimum the gas limit may ever be.
	GenesisGasLimit      uint64 = 4712388 // Gas limit of the Genesis block.

	MaximumExtraDataSize uint64 = 32 // Maximum size extra data may be after Genesis.

	EpochDuration   uint64 = 30000 // Duration between proof-of-work epochs.
	CallCreateDepth uint64 = 1024  // Maximum depth of call/create stack.
	StackLimit      uint64 = 1024  // Maximum size of VM stack allowed.

	// eip-3860: limit and meter initcode (Shanghai)
	MaxCodeSize     = 24576           // Maximum bytecode to permit for a contract
	MaxInitCodeSize = 2 * MaxCodeSize // Maximum initcode to permit in a creation transaction and create instructions

	// istanbul BFT
	BFTMaximumExtraDataSize uint64 = 65 // Maximum size extra data may be after Genesis.

	// AccountKey
	// TODO-Kaia: Need to fix below values.
	TxAccountCreationGasDefault uint64 = 0
	TxValidationGasDefault      uint64 = 0
	TxAccountCreationGasPerKey  uint64 = 20000 // WARNING: With integer overflow in mind before changing this value.
	TxValidationGasPerKey       uint64 = 15000 // WARNING: With integer overflow in mind before changing this value.

	// Fee for new tx types
	// TODO-Kaia: Need to fix values
	TxGasAccountCreation       uint64 = 21000
	TxGasAccountUpdate         uint64 = 21000
	TxGasFeeDelegated          uint64 = 10000
	TxGasFeeDelegatedWithRatio uint64 = 15000
	TxGasCancel                uint64 = 21000

	// Network Id
	UnusedNetworkId              uint64 = 0
	AspenNetworkId               uint64 = 1000
	KairosNetworkId              uint64 = 1001
	MainnetNetworkId             uint64 = 8217
	ServiceChainDefaultNetworkId uint64 = 3000

	TxGasValueTransfer     uint64 = 21000
	TxGasContractExecution uint64 = 21000

	TxDataGas uint64 = 100

	TxAccessListAddressGas    uint64 = 2400 // Per address specified in EIP 2930 access list
	TxAccessListStorageKeyGas uint64 = 1900 // Per storage key specified in EIP 2930 access list

	TxAuthTupleGas uint64 = 12500 // Per auth tuple code specified in EIP-7702

	Bls12381G1AddGas          uint64 = 375   // Price for BLS12-381 elliptic curve G1 point addition
	Bls12381G1MulGas          uint64 = 12000 // Price for BLS12-381 elliptic curve G1 point scalar multiplication
	Bls12381G2AddGas          uint64 = 600   // Price for BLS12-381 elliptic curve G2 point addition
	Bls12381G2MulGas          uint64 = 22500 // Price for BLS12-381 elliptic curve G2 point scalar multiplication
	Bls12381PairingBaseGas    uint64 = 37700 // Base gas price for BLS12-381 elliptic curve pairing check
	Bls12381PairingPerPairGas uint64 = 32600 // Per-point pair gas price for BLS12-381 elliptic curve pairing check
	Bls12381MapG1Gas          uint64 = 5500  // Gas price for BLS12-381 mapping field element to G1 operation
	Bls12381MapG2Gas          uint64 = 23800 // Gas price for BLS12-381 mapping field element to G2 operation

	// ZeroBaseFee exists for supporting Ethereum compatible data structure.
	ZeroBaseFee uint64 = 0

	HistoryServeWindow = 8192 // Number of blocks to serve historical block hashes for, EIP-2935.
)

const (
	DefaultBlockGenerationInterval  = int64(1) // unit: seconds
	DefaultBlockGenerationTimeLimit = 250 * time.Millisecond
)

var Bls12381G1MultiExpDiscountTable = [128]uint64{1000, 949, 848, 797, 764, 750, 738, 728, 719, 712, 705, 698, 692, 687, 682, 677, 673, 669, 665, 661, 658, 654, 651, 648, 645, 642, 640, 637, 635, 632, 630, 627, 625, 623, 621, 619, 617, 615, 613, 611, 609, 608, 606, 604, 603, 601, 599, 598, 596, 595, 593, 592, 591, 589, 588, 586, 585, 584, 582, 581, 580, 579, 577, 576, 575, 574, 573, 572, 570, 569, 568, 567, 566, 565, 564, 563, 562, 561, 560, 559, 558, 557, 556, 555, 554, 553, 552, 551, 550, 549, 548, 547, 547, 546, 545, 544, 543, 542, 541, 540, 540, 539, 538, 537, 536, 536, 535, 534, 533, 532, 532, 531, 530, 529, 528, 528, 527, 526, 525, 525, 524, 523, 522, 522, 521, 520, 520, 519}

var Bls12381G2MultiExpDiscountTable = [128]uint64{1000, 1000, 923, 884, 855, 832, 812, 796, 782, 770, 759, 749, 740, 732, 724, 717, 711, 704, 699, 693, 688, 683, 679, 674, 670, 666, 663, 659, 655, 652, 649, 646, 643, 640, 637, 634, 632, 629, 627, 624, 622, 620, 618, 615, 613, 611, 609, 607, 606, 604, 602, 600, 598, 597, 595, 593, 592, 590, 589, 587, 586, 584, 583, 582, 580, 579, 578, 576, 575, 574, 573, 571, 570, 569, 568, 567, 566, 565, 563, 562, 561, 560, 559, 558, 557, 556, 555, 554, 553, 552, 552, 551, 550, 549, 548, 547, 546, 545, 545, 544, 543, 542, 541, 541, 540, 539, 538, 537, 537, 536, 535, 535, 534, 533, 532, 532, 531, 530, 530, 529, 528, 528, 527, 526, 526, 525, 524, 524}

var (
	// Dummy Randao fields to be used in a Randao-enabled Genesis block.
	ZeroRandomReveal = make([]byte, 96)
	ZeroMixHash      = make([]byte, 32)

	TxGasHumanReadable uint64 = 4000000000 // NOTE: HumanReadable related functions are inactivated now

	// TODO-Kaia Change the variables used in GXhash to more appropriate values for Kaia network
	BlockScoreBoundDivisor = big.NewInt(2048)   // The bound divisor of the blockscore, used in the update calculations.
	GenesisBlockScore      = big.NewInt(131072) // BlockScore of the Genesis block.
	MinimumBlockScore      = big.NewInt(131072) // The minimum that the blockscore may ever be.
	DurationLimit          = big.NewInt(13)     // The decision boundary on the blocktime duration used to determine whether blockscore should go up or not.

	// SystemAddress is where the system-transaction is sent from as per EIP-4788
	SystemAddress = common.HexToAddress("0xfffffffffffffffffffffffffffffffffffffffe")

	// EIP-2935 - Serve historical block hashes from state
	HistoryStorageAddress = common.HexToAddress("0x0000F90827F1C53a10cb7A02335B175320002935")
	HistoryStorageCode    = common.FromHex("3373fffffffffffffffffffffffffffffffffffffffe14604657602036036042575f35600143038111604257611fff81430311604257611fff9006545f5260205ff35b5f5ffd5b5f35611fff60014303065500")
)

// Parameters for execution time limit
// These parameters will be re-assigned by init options
var (
	// Execution time limit for all txs in a block
	BlockGenerationTimeLimit = DefaultBlockGenerationTimeLimit

	// TODO-Kaia-Governance Change the following variables to governance items which requires consensus of CCN
	// Block generation interval in seconds. It should be equal or larger than 1
	BlockGenerationInterval = DefaultBlockGenerationInterval
)

// istanbul BFT
func GetMaximumExtraDataSize() uint64 {
	return BFTMaximumExtraDataSize
}

// CodeFormat is the version of the interpreter that smart contract uses
type CodeFormat uint8

// Supporting CodeFormat
// CodeFormatLast should be equal or less than 16 because only the last 4 bits of CodeFormat are used for CodeInfo.
const (
	CodeFormatEVM CodeFormat = iota
	CodeFormatLast
)

func (t CodeFormat) Validate() bool {
	if t < CodeFormatLast {
		return true
	}
	return false
}

func (t CodeFormat) String() string {
	switch t {
	case CodeFormatEVM:
		return "CodeFormatEVM"
	}

	return "UndefinedCodeFormat"
}

// VmVersion contains the information of the contract deployment time (ex. 0x0(constantinople), 0x1(istanbul,...))
type VmVersion uint8

// Supporting VmVersion
const (
	VmVersion0 VmVersion = iota // Deployed at Constantinople
	VmVersion1                  // Deployed at Istanbul, ...(later HFs would be added)
)

func (t VmVersion) String() string {
	return "VmVersion" + string(t)
}

// CodeInfo consists of 8 bits, and has information of the contract code.
// Originally, codeInfo only contains codeFormat information(interpreter version), but now it is divided into two parts.
// First four bit contains the deployment time (ex. 0x00(constantinople), 0x10(istanbul,...)), so it is called vmVersion.
// Last four bit contains the interpreter version (ex. 0x00(EVM), 0x01(EWASM)), so it is called codeFormat.
type CodeInfo uint8

const (
	// codeFormatBitMask filters only the codeFormat. It means the interpreter version used by the contract.
	// Mask result 1. [x x x x 0 0 0 1]. The contract uses EVM interpreter.
	codeFormatBitMask = 0b00001111

	// vmVersionBitMask filters only the vmVersion. It means deployment time of the contract.
	// Mask result 1. [0 0 0 0 x x x x]. The contract is deployed at constantinople
	// Mask result 2. [0 0 0 1 x x x x]. The contract is deployed after istanbulHF
	vmVersionBitMask = 0b11110000
)

func NewCodeInfo(codeFormat CodeFormat, vmVersion VmVersion) CodeInfo {
	return CodeInfo(codeFormat&codeFormatBitMask) | CodeInfo(vmVersion)<<4
}

func NewCodeInfoWithRules(codeFormat CodeFormat, r Rules) CodeInfo {
	var vmVersion VmVersion
	switch {
	// If new HF is added, please add new case below
	// case r.IsNextHF:          // If this HF is backward compatible with vmVersion1.
	case r.IsLondon:
		fallthrough
	case r.IsIstanbul:
		vmVersion = VmVersion1
	default:
		vmVersion = VmVersion0
	}
	return NewCodeInfo(codeFormat, vmVersion)
}

func (t CodeInfo) GetCodeFormat() CodeFormat {
	return CodeFormat(t & codeFormatBitMask)
}

func (t CodeInfo) GetVmVersion() VmVersion {
	return VmVersion(t & vmVersionBitMask >> 4)
}

func (t CodeInfo) String() string {
	return fmt.Sprintf("[%s, %s]", t.GetCodeFormat().String(), t.GetVmVersion().String())
}
