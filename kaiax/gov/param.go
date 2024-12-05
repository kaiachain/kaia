package gov

import (
	"bytes"
	"errors"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
)

type canonicalizerT func(v any) (any, error)

type Param struct {
	Canonicalizer canonicalizerT
	FormatChecker func(cv any) bool // validation on canonical value.

	DefaultValue  any
	VoteForbidden bool
}

var (
	addressCanonicalizer canonicalizerT = func(v any) (any, error) {
		switch v := v.(type) {
		case []byte:
			if len(v) != common.AddressLength {
				return nil, ErrCanonicalizeByteToAddress
			}
			return common.BytesToAddress(v), nil
		case string:
			if !common.IsHexAddress(v) {
				return nil, ErrCanonicalizeStringToAddress
			}
			return common.HexToAddress(v), nil
		case common.Address:
			return v, nil
		}
		return nil, ErrCanonicalizeToAddress
	}

	validatorAddressListCanonicalizer canonicalizerT = func(v any) (any, error) {
		stringToAddressList := func(v string) ([]common.Address, error) {
			ret := []common.Address{}
			for _, address := range strings.Split(v, ",") {
				if !common.IsHexAddress(address) {
					return nil, ErrCanonicalizeStringToAddress
				}
				ret = append(ret, common.HexToAddress(address))
			}
			return ret, nil
		}

		switch v := v.(type) {
		case []byte: // input from header.Vote
			// There are three types of header.Vote encoding for validator address(es).
			// Type1. Single address, [20]byte. See Mainnet block 5505383.
			// Type2. Multiple addresses, [20*n]byte.
			// Type3. Single address, [42]byte (hex-encoded bytes). See Mainnet block 90915008.

			// Type1
			if len(v) == common.AddressLength {
				return []common.Address{common.BytesToAddress(v)}, nil
			}

			// Type2
			if len(v)%common.AddressLength == 0 {
				addresses := make([]common.Address, len(v)/common.AddressLength)
				for i := 0; i < len(v)/common.AddressLength; i++ {
					addresses[i] = common.BytesToAddress(v[i*common.AddressLength : (i+1)*common.AddressLength])
				}
				return addresses, nil
			}

			// Type3
			v, err := hexutil.Decode(string(v))
			if err != nil {
				return nil, errors.Join(ErrCanonicalizeByteToAddress, err)
			}

			if len(v) == common.AddressLength {
				return []common.Address{common.BytesToAddress(v)}, nil
			}

			return nil, ErrCanonicalizeToAddressList
		case string: // input from API
			return stringToAddressList(v)
		case common.Address:
			return []common.Address{v}, nil
		case []common.Address:
			return v, nil
		}

		return nil, ErrCanonicalizeToAddressList
	}

	bigIntCanonicalizer canonicalizerT = func(v any) (any, error) {
		switch v := v.(type) {
		case []byte:
			cv, ok := new(big.Int).SetString(string(v), 10)
			if !ok {
				return nil, ErrCanonicalizeByteToBigInt
			}
			return cv, nil
		case string:
			cv, ok := new(big.Int).SetString(v, 10)
			if !ok {
				return nil, ErrCanonicalizeStringToBigInt
			}
			return cv, nil
		case *big.Int:
			return v, nil
		}
		return nil, ErrCanonicalizeBigInt
	}

	boolCanonicalizer canonicalizerT = func(v any) (any, error) {
		switch v := v.(type) {
		case []byte:
			if bytes.Equal(v, []byte{0x01}) {
				return true, nil
			} else if bytes.Equal(v, []byte{0x00}) {
				return false, nil
			} else {
				return nil, ErrCanonicalizeByteToBool
			}
		case bool:
			return v, nil
		}
		return nil, ErrCanonicalizeBool
	}

	stringCanonicalizer canonicalizerT = func(v any) (any, error) {
		switch v := v.(type) {
		case []byte:
			return string(v), nil
		case string:
			return v, nil
		}
		return nil, ErrCanonicalizeString
	}

	uint64Canonicalizer canonicalizerT = func(v any) (any, error) {
		switch v := v.(type) {
		case []byte:
			if len(v) > 8 {
				return nil, ErrCanonicalizeByteToUint64
			}
			return new(big.Int).SetBytes(v).Uint64(), nil
		case float64:
			if float64(uint64(v)) != v {
				return nil, ErrCanonicalizeFloatToUint64
			}

			return uint64(v), nil
		case uint64:
			return v, nil
		}
		return nil, ErrCanonicalizeUint64
	}
)

func noopFormatChecker(cv any) bool {
	return true
}

func valSetVoteFormatChecker(cv any) bool {
	_, ok := cv.(common.Address)
	if ok {
		return true
	}

	v, ok := cv.([]common.Address)
	if !ok || len(v) == 0 {
		return false
	}

	// do not allow duplicated addresses
	var duplicateCheckMap map[common.Address]bool
	for _, address := range v {
		if ok, _ := duplicateCheckMap[address]; ok {
			return false
		}
	}
	return true
}

type ParamName string

// alphabetically sorted. These are only used in-memory, so the order does not matter.
const (
	GovernanceDeriveShaImpl        ParamName = "governance.deriveshaimpl"
	GovernanceGovernanceMode       ParamName = "governance.governancemode"
	GovernanceGoverningNode        ParamName = "governance.governingnode"
	GovernanceGovParamContract     ParamName = "governance.govparamcontract"
	GovernanceUnitPrice            ParamName = "governance.unitprice"
	IstanbulCommitteeSize          ParamName = "istanbul.committeesize"
	IstanbulEpoch                  ParamName = "istanbul.epoch"
	IstanbulPolicy                 ParamName = "istanbul.policy"
	Kip71BaseFeeDenominator        ParamName = "kip71.basefeedenominator"
	Kip71GasTarget                 ParamName = "kip71.gastarget"
	Kip71LowerBoundBaseFee         ParamName = "kip71.lowerboundbasefee"
	Kip71MaxBlockGasUsedForBaseFee ParamName = "kip71.maxblockgasusedforbasefee"
	Kip71UpperBoundBaseFee         ParamName = "kip71.upperboundbasefee"
	RewardDeferredTxFee            ParamName = "reward.deferredtxfee"
	RewardKip82Ratio               ParamName = "reward.kip82ratio"
	RewardMintingAmount            ParamName = "reward.mintingamount"
	RewardMinimumStake             ParamName = "reward.minimumstake"
	RewardProposerUpdateInterval   ParamName = "reward.proposerupdateinterval"
	RewardRatio                    ParamName = "reward.ratio"
	RewardStakingUpdateInterval    ParamName = "reward.stakingupdateinterval"
	RewardUseGiniCoeff             ParamName = "reward.useginicoeff"
)

const (
	AddValidator    ParamName = "governance.addvalidator"
	RemoveValidator ParamName = "governance.removevalidator"
)

var Params = map[ParamName]*Param{
	GovernanceDeriveShaImpl: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v <= 2 // deriveShaImpl has only three options.
		},
		DefaultValue:  uint64(0),
		VoteForbidden: false,
	},
	GovernanceGovernanceMode: {
		Canonicalizer: stringCanonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			if v == "none" || v == "single" {
				return true
			}
			return false
		},
		DefaultValue:  "none",
		VoteForbidden: true,
	},
	GovernanceGoverningNode: {
		Canonicalizer: addressCanonicalizer,
		FormatChecker: func(cv any) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	GovernanceGovParamContract: {
		Canonicalizer: addressCanonicalizer,
		FormatChecker: func(cv any) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	GovernanceUnitPrice: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(250e9),
		VoteForbidden: false,
	},
	IstanbulCommitteeSize: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v > 0
		},
		DefaultValue:  uint64(21),
		VoteForbidden: false,
	},
	IstanbulEpoch: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(604800),
		VoteForbidden: true,
	},
	IstanbulPolicy: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(uint64)
			if !ok {
				return false
			}
			return v <= 2 // policy has only three options.
		},
		DefaultValue:  uint64(RoundRobin),
		VoteForbidden: true,
	},
	Kip71BaseFeeDenominator: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(uint64)
			return ok && v != 0
		},
		DefaultValue:  uint64(20),
		VoteForbidden: false,
	},
	Kip71GasTarget: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(30000000),
		VoteForbidden: false,
	},
	Kip71LowerBoundBaseFee: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(25000000000),
		VoteForbidden: false,
	},
	Kip71MaxBlockGasUsedForBaseFee: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(60000000),
		VoteForbidden: false,
	},
	Kip71UpperBoundBaseFee: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(750000000000),
		VoteForbidden: false,
	},
	RewardDeferredTxFee: {
		Canonicalizer: boolCanonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  false,
		VoteForbidden: true,
	},
	RewardKip82Ratio: {
		Canonicalizer: stringCanonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			parts := strings.Split(v, "/")
			if len(parts) != 2 {
				return false
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false
				}
				if num < 0 {
					return false
				}
				sum += num
			}

			return sum == 100
		},
		DefaultValue:  "20/80",
		VoteForbidden: false,
	},
	RewardMintingAmount: {
		Canonicalizer: bigIntCanonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  big.NewInt(0),
		VoteForbidden: false,
	},
	RewardMinimumStake: {
		Canonicalizer: bigIntCanonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(*big.Int)
			if !ok {
				return false
			}
			return v.Sign() >= 0
		},
		DefaultValue:  big.NewInt(2000000),
		VoteForbidden: true,
	},
	RewardProposerUpdateInterval: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(3600),
		VoteForbidden: true,
	},
	RewardRatio: {
		Canonicalizer: stringCanonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(string)
			if !ok {
				return false
			}
			parts := strings.Split(v, "/")
			if len(parts) != 3 {
				return false
			}
			sum := 0
			for _, part := range parts {
				num, err := strconv.Atoi(part)
				if err != nil {
					return false
				}
				if num < 0 {
					return false
				}
				sum += num
			}

			return sum == 100
		},
		DefaultValue:  "100/0/0",
		VoteForbidden: false,
	},
	RewardStakingUpdateInterval: {
		Canonicalizer: uint64Canonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  uint64(86400),
		VoteForbidden: true,
	},
	RewardUseGiniCoeff: {
		Canonicalizer: boolCanonicalizer,
		FormatChecker: noopFormatChecker,
		DefaultValue:  false,
		VoteForbidden: true,
	},
}

var ValidatorParams = map[ParamName]*Param{
	AddValidator: {
		Canonicalizer: validatorAddressListCanonicalizer,
		FormatChecker: valSetVoteFormatChecker,
		DefaultValue:  []common.Address{},
		VoteForbidden: false,
	},
	RemoveValidator: {
		Canonicalizer: validatorAddressListCanonicalizer,
		FormatChecker: valSetVoteFormatChecker,
		DefaultValue:  []common.Address{},
		VoteForbidden: false,
	},
}

const (
	// Proposer policy
	RoundRobin = iota
	Sticky
	WeightedRandom
	ProposerPolicy_End
)
