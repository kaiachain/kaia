// Copyright 2026 The Kaia Authors
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
package params

import (
	"errors"
	"fmt"
	"math/big"
)

var (
	// DefaultOsakaBlobConfig is the default blob configuration for the Osaka fork.
	DefaultOsakaBlobConfig = &BlobConfig{
		Target:         1,
		Max:            1,
		UpdateFraction: 5007716,
	}
	// DefaultBlobSchedule is the latest configured blob schedule for Ethereum mainnet.
	DefaultBlobSchedule = &BlobScheduleConfig{
		Osaka: DefaultOsakaBlobConfig,
	}
)

// BlobConfig specifies the target and max blobs per block for the associated fork.
type BlobConfig struct {
	Target         int    `json:"target"`
	Max            int    `json:"max"`
	UpdateFraction uint64 `json:"baseFeeUpdateFraction"`
}

// String implement fmt.Stringer, returning string format blob config.
func (bcfg *BlobConfig) String() string {
	if bcfg == nil {
		return "nil"
	}
	return fmt.Sprintf("target: %d, max: %d, fraction: %d", bcfg.Target, bcfg.Max, bcfg.UpdateFraction)
}

func (bcfg *BlobConfig) MaxBlobGas() uint64 {
	return uint64(bcfg.Max) * BlobTxBlobGasPerBlob
}

// blobPriceEIP4844 returns the price for EIP-4844 of one blob in Wei.
func (bcfg *BlobConfig) BlobPriceEIP4844(excessBlobGas uint64) *big.Int {
	f := bcfg.BlobBaseFeeEIP4844(excessBlobGas)
	return new(big.Int).Mul(f, big.NewInt(BlobTxBlobGasPerBlob))
}

// blobBaseFeeEIP4844 computes the blob fee for EIP-4844.
func (bcfg *BlobConfig) BlobBaseFeeEIP4844(excessBlobGas uint64) *big.Int {
	return fakeExponential(big.NewInt(BlobTxMinBlobGasprice), new(big.Int).SetUint64(excessBlobGas), new(big.Int).SetUint64(bcfg.UpdateFraction))
}

// VerifyEIP4844Header verifies the presence of the excessBlobGas field and that
// if the current block contains no transactions, the excessBlobGas is updated
// accordingly.
func (bcfg *BlobConfig) VerifyEIP4844Header(parentExcessBlobGas, parentBlobGasUsed *uint64, headerExcessBlobGas, headerBlobGasUsed *uint64) error {
	if headerExcessBlobGas == nil {
		return errors.New("header is missing excessBlobGas")
	}
	if headerBlobGasUsed == nil {
		return errors.New("header is missing blobGasUsed")
	}

	// Verify that the blob gas used remains within reasonable limits.
	if *headerBlobGasUsed > bcfg.MaxBlobGas() {
		return fmt.Errorf("blob gas used %d exceeds maximum allowance %d", *headerBlobGasUsed, bcfg.MaxBlobGas())
	}
	if *headerBlobGasUsed%BlobTxBlobGasPerBlob != 0 {
		return fmt.Errorf("blob gas used %d not a multiple of blob gas per blob %d", *headerBlobGasUsed, BlobTxBlobGasPerBlob)
	}

	// Verify the excessBlobGas is correct based on the parent header
	expectedExcessBlobGas := bcfg.CalcExcessBlobGas(parentExcessBlobGas, parentBlobGasUsed)
	if *headerExcessBlobGas != expectedExcessBlobGas {
		return fmt.Errorf("invalid excessBlobGas: have %d, want %d", *headerExcessBlobGas, expectedExcessBlobGas)
	}
	return nil
}

// CalcExcessBlobGas calculates the excess blob gas for KIP-279.
func (bcfg *BlobConfig) CalcExcessBlobGas(parentExcessBlobGasPtr, parentBlobGasUsedPtr *uint64) uint64 {
	var parentExcessBlobGas, parentBlobGasUsed uint64
	if parentExcessBlobGasPtr != nil {
		parentExcessBlobGas = *parentExcessBlobGasPtr
		parentBlobGasUsed = *parentBlobGasUsedPtr
	}

	// return max(0, parentExcessBlobGas+parentBlobGasUsed-(uint64(bcfg.Target)*params.BlobTxBlobGasPerBlob))
	// Since max(0, negative) cannot be expressed in uint64, it is calculated as follows.
	if parentExcessBlobGas+parentBlobGasUsed < uint64(bcfg.Target)*BlobTxBlobGasPerBlob {
		return 0
	}
	return parentExcessBlobGas + parentBlobGasUsed - (uint64(bcfg.Target) * BlobTxBlobGasPerBlob)
}

// fakeExponential approximates factor * e ** (numerator / denominator) using
// Taylor expansion.
func fakeExponential(factor, numerator, denominator *big.Int) *big.Int {
	var (
		output = new(big.Int)
		accum  = new(big.Int).Mul(factor, denominator)
	)
	for i := 1; accum.Sign() > 0; i++ {
		output.Add(output, accum)

		accum.Mul(accum, numerator)
		accum.Div(accum, denominator)
		accum.Div(accum, big.NewInt(int64(i)))
	}
	return output.Div(output, denominator)
}

// VerifyEIP4844HeaderForEEST verifies the presence of the excessBlobGas field and that
// if the current block contains no transactions, the excessBlobGas is updated
// accordingly.
func (bcfg *BlobConfig) VerifyEIP4844HeaderForEEST(parentExcessBlobGas, parentBlobGasUsed *uint64, parentBaseFee *big.Int, headerExcessBlobGas, headerBlobGasUsed *uint64) error {
	if headerExcessBlobGas == nil {
		return errors.New("header is missing excessBlobGas")
	}
	if headerBlobGasUsed == nil {
		return errors.New("header is missing blobGasUsed")
	}

	// Verify that the blob gas used remains within reasonable limits.
	if *headerBlobGasUsed > bcfg.MaxBlobGas() {
		return fmt.Errorf("blob gas used %d exceeds maximum allowance %d", *headerBlobGasUsed, bcfg.MaxBlobGas())
	}
	if *headerBlobGasUsed%BlobTxBlobGasPerBlob != 0 {
		return fmt.Errorf("blob gas used %d not a multiple of blob gas per blob %d", *headerBlobGasUsed, BlobTxBlobGasPerBlob)
	}

	// Verify the excessBlobGas is correct based on the parent header.
	// EEST expects the next block's excess to be derived from the parent values.
	expectedExcessBlobGas := bcfg.CalcExcessBlobGasForEEST(parentExcessBlobGas, parentBlobGasUsed, parentBaseFee)
	if *headerExcessBlobGas != expectedExcessBlobGas {
		return fmt.Errorf("invalid excessBlobGas: have %d, want %d", *headerExcessBlobGas, expectedExcessBlobGas)
	}
	return nil
}

// CalcExcessBlobGasForEEST calculates the excess blob gas after applying the set of
// blobs on top of the excess blob gas for EIP-7918.
// NOTE: This is not used because Kaia calculates blobBaseFee by multiplying it with baseFee.
// For compatibility, `blobBaseFeeEIP4844` and `blobPriceEIP4844` is also being kept for CalcExcessBlobGasForEEST.
// ref: https://github.com/kaiachain/kips/blob/fac337aad17db40ed2a6c988e03ad6b2ec9771f2/KIPs/kip-279.md#blob-base-fee-as-a-multiple-of-base-fee
func (bcfg *BlobConfig) CalcExcessBlobGasForEEST(parentExcessBlobGasPtr, parentBlobGasUsedPtr *uint64, parentBaseFee *big.Int) uint64 {
	var parentExcessBlobGas, parentBlobGasUsed uint64
	if parentExcessBlobGasPtr != nil {
		parentExcessBlobGas = *parentExcessBlobGasPtr
		parentBlobGasUsed = *parentBlobGasUsedPtr
	}

	var (
		excessBlobGas = parentExcessBlobGas + parentBlobGasUsed
		targetGas     = uint64(bcfg.Target) * BlobTxBlobGasPerBlob
	)
	if excessBlobGas < targetGas {
		return 0
	}

	// EIP-7918 (post-Osaka) introduces a different formula for computing excess,
	// in cases where the price is lower than a 'reserve price'.
	var (
		baseCost     = big.NewInt(BlobBaseCost)
		reservePrice = baseCost.Mul(baseCost, parentBaseFee)
		blobPrice    = bcfg.BlobPriceEIP4844(parentExcessBlobGas)
	)
	if reservePrice.Cmp(blobPrice) > 0 {
		scaledExcess := parentBlobGasUsed * uint64(bcfg.Max-bcfg.Target) / uint64(bcfg.Max)
		return parentExcessBlobGas + scaledExcess
	}

	// Original EIP-4844 formula.
	return excessBlobGas - targetGas
}

// BlobScheduleConfig determines target and max number of blobs allow per fork.
type BlobScheduleConfig struct {
	Osaka *BlobConfig `json:"osaka,omitempty"`
}

// CalcBlobFee calculates the blobfee for KIP-279 from the header's base fee field.
func CalcBlobFee(baseFee *big.Int) *big.Int {
	// If baseFee is nil, CalcBlobFee will return nil as it cannot be calculated.
	// Returning nil makes the difference from 0 explicit.
	if baseFee == nil {
		return nil
	}
	return new(big.Int).Mul(baseFee, new(big.Int).SetUint64(BlobBaseFeeMultiplier))
}
