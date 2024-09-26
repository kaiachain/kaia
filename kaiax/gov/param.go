package gov

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"

	"github.com/kaiachain/kaia/common"
)

type canonicalizerT func(v any) (any, error)

type Param struct {
	ParamSetFieldName string
	Canonicalizer     canonicalizerT
	FormatChecker     func(cv any) bool // validation on canonical value.

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

	addressListCanonicalizer canonicalizerT = func(v any) (any, error) {
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
		case []byte:
			return stringToAddressList(string(v))
		case string:
			return stringToAddressList(v)
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

var Params = map[ParamName]*Param{
	GovernanceDeriveShaImpl: {
		ParamSetFieldName: "DeriveShaImpl",
		Canonicalizer:     uint64Canonicalizer,
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
		ParamSetFieldName: "GovernanceMode",
		Canonicalizer:     stringCanonicalizer,
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
		ParamSetFieldName: "GoverningNode",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker: func(cv any) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	GovernanceGovParamContract: {
		ParamSetFieldName: "GovParamContract",
		Canonicalizer:     addressCanonicalizer,
		FormatChecker: func(cv any) bool {
			_, ok := cv.(common.Address)
			return ok
		},
		DefaultValue:  common.HexToAddress("0x0000000000000000000000000000000000000000"),
		VoteForbidden: false,
	},
	GovernanceUnitPrice: {
		ParamSetFieldName: "UnitPrice",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(250e9),
		VoteForbidden:     false,
	},
	IstanbulCommitteeSize: {
		ParamSetFieldName: "CommitteeSize",
		Canonicalizer:     uint64Canonicalizer,
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
		ParamSetFieldName: "Epoch",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(604800),
		VoteForbidden:     true,
	},
	IstanbulPolicy: {
		ParamSetFieldName: "ProposerPolicy",
		Canonicalizer:     uint64Canonicalizer,
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
		ParamSetFieldName: "BaseFeeDenominator",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker: func(cv any) bool {
			v, ok := cv.(uint64)
			return ok && v != 0
		},
		DefaultValue:  uint64(20),
		VoteForbidden: false,
	},
	Kip71GasTarget: {
		ParamSetFieldName: "GasTarget",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(30000000),
		VoteForbidden:     false,
	},
	Kip71LowerBoundBaseFee: {
		ParamSetFieldName: "LowerBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(25000000000),
		VoteForbidden:     false,
	},
	Kip71MaxBlockGasUsedForBaseFee: {
		ParamSetFieldName: "MaxBlockGasUsedForBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(60000000),
		VoteForbidden:     false,
	},
	Kip71UpperBoundBaseFee: {
		ParamSetFieldName: "UpperBoundBaseFee",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(750000000000),
		VoteForbidden:     false,
	},
	RewardDeferredTxFee: {
		ParamSetFieldName: "DeferredTxFee",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      false,
		VoteForbidden:     true,
	},
	RewardKip82Ratio: {
		ParamSetFieldName: "Kip82Ratio",
		Canonicalizer:     stringCanonicalizer,
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
		ParamSetFieldName: "MintingAmount",
		Canonicalizer:     bigIntCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      big.NewInt(0),
		VoteForbidden:     false,
	},
	RewardMinimumStake: {
		ParamSetFieldName: "MinimumStake",
		Canonicalizer:     bigIntCanonicalizer,
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
		ParamSetFieldName: "ProposerUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(3600),
		VoteForbidden:     true,
	},
	RewardRatio: {
		ParamSetFieldName: "Ratio",
		Canonicalizer:     stringCanonicalizer,
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
		ParamSetFieldName: "StakingUpdateInterval",
		Canonicalizer:     uint64Canonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      uint64(86400),
		VoteForbidden:     true,
	},
	RewardUseGiniCoeff: {
		ParamSetFieldName: "UseGiniCoeff",
		Canonicalizer:     boolCanonicalizer,
		FormatChecker:     noopFormatChecker,
		DefaultValue:      false,
		VoteForbidden:     true,
	},
}

const (
	// Proposer policy
	RoundRobin = iota
	Sticky
	WeightedRandom
	ProposerPolicy_End
)
