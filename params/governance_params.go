// Modifications Copyright 2024 The Kaia Authors
// Copyright 2019 The klaytn Authors
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

package params

import (
	"math/big"
)

const (
	// Block reward will be separated by three pieces and distributed
	RewardSliceCount      = 3
	RewardKip82SliceCount = 2
	// GovernanceConfig is stored in a cache which has below capacity
	GovernanceCacheLimit    = 512
	GovernanceIdxCacheLimit = 1000
	// The prefix for governance cache
	GovernanceCachePrefix = "governance"

	CheckpointInterval       = 1024 // For Istanbul snapshot
	SupplyCheckpointInterval = 128  // for SupplyManager tracking native token supply
)

const (
	// Governance Key
	GovernanceMode int = iota
	GoverningNode
	Epoch
	Policy
	CommitteeSize
	UnitPrice
	MintingAmount
	Ratio
	UseGiniCoeff
	DeferredTxFee
	MinimumStake
	AddValidator
	RemoveValidator
	StakeUpdateInterval
	ProposerRefreshInterval
	ConstTxGasHumanReadable
	CliqueEpoch
	Timeout
	LowerBoundBaseFee
	UpperBoundBaseFee
	GasTarget
	MaxBlockGasUsedForBaseFee
	BaseFeeDenominator
	GovParamContract
	Kip82Ratio
	DeriveShaImpl
)

const (
	GovernanceMode_None = iota
	GovernanceMode_Single
	GovernanceMode_Ballot
)

const (
	// Proposer policy
	// At the moment this is duplicated in istanbul/config.go, not to make a cross reference
	// TODO-Klatn-Governance: Find a way to manage below constants at single location
	RoundRobin = iota
	Sticky
	WeightedRandom
)

var (
	// Default Values: Constants used for getting default values for configuration
	DefaultGovernanceMode            = "none"
	DefaultGoverningNode             = "0x0000000000000000000000000000000000000000"
	DefaultGovParamContract          = "0x0000000000000000000000000000000000000000"
	DefaultEpoch                     = uint64(604800)
	DefaultProposerPolicy            = uint64(RoundRobin)
	DefaultSubGroupSize              = uint64(21)
	DefaultUnitPrice                 = uint64(250000000000)
	DefaultLowerBoundBaseFee         = uint64(25000000000)
	DefaultUpperBoundBaseFee         = uint64(750000000000)
	DefaultGasTarget                 = uint64(30000000)
	DefaultMaxBlockGasUsedForBaseFee = uint64(60000000)
	DefaultBaseFeeDenominator        = uint64(20)
	DefaultMintingAmount             = big.NewInt(0)
	DefaultRatio                     = "100/0/0"
	DefaultKip82Ratio                = "20/80"
	DefaultUseGiniCoeff              = false
	DefaultDeferredTxFee             = false
	DefaultMinimumStake              = big.NewInt(2000000)
	DefaultStakeUpdateInterval       = uint64(86400) // 1 day
	DefaultProposerRefreshInterval   = uint64(3600)  // 1 hour
	DefaultDeriveShaImpl             = uint64(0)     // Orig
)

func IsCheckpointInterval(blockNum uint64) bool {
	return blockNum != 0 && blockNum%CheckpointInterval == 0
}
