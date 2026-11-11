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
	"fmt"
	"math/big"
	"testing"
)

var GethCancunBlobConfig = &BlobConfig{
	Target:         3,
	Max:            6,
	UpdateFraction: 3338477,
}

func TestCalcExcessBlobGas(t *testing.T) {
	config := TestChainConfig.Copy()
	var (
		targetBlobs   = config.BlobScheduleConfig.Osaka.Target
		targetBlobGas = uint64(targetBlobs) * BlobTxBlobGasPerBlob
	)

	tests := []struct {
		excess uint64
		blobs  int
		want   uint64
	}{
		// In Osaka, all patterns where (excessBlobGas, blobGasUsed)
		// is 0 or targetBlobGas are tested.
		{0, 0, 0},
		{targetBlobGas, 0, 0},
		{0, targetBlobs, 0},
		{targetBlobGas, targetBlobs, targetBlobGas},
	}
	for i, tt := range tests {
		blobGasUsed := uint64(tt.blobs) * BlobTxBlobGasPerBlob
		result := config.LatestBlobConfig(config.OsakaCompatibleBlock).CalcExcessBlobGas(&tt.excess, &blobGasUsed)
		if result != tt.want {
			t.Errorf("test %d: excess blob gas mismatch: have %v, want %v", i, result, tt.want)
		}
	}
}

func TestCalcExcessBlobGasEIP4844(t *testing.T) {
	config := MainnetChainConfig.Copy()
	config.OsakaCompatibleBlock = big.NewInt(0)
	config.BlobScheduleConfig.Osaka = GethCancunBlobConfig
	var (
		targetBlobs   = config.BlobScheduleConfig.Osaka.Target
		targetBlobGas = uint64(targetBlobs) * BlobTxBlobGasPerBlob
	)

	tests := []struct {
		excess uint64
		blobs  int
		want   uint64
	}{
		// The excess blob gas should not increase from zero if the used blob
		// slots are below - or equal - to the target.
		{0, 0, 0},
		{0, 1, 0},
		{0, targetBlobs, 0},

		// If the target blob gas is exceeded, the excessBlobGas should increase
		// by however much it was overshot
		{0, targetBlobs + 1, BlobTxBlobGasPerBlob},
		{1, targetBlobs + 1, BlobTxBlobGasPerBlob + 1},
		{1, targetBlobs + 2, 2*BlobTxBlobGasPerBlob + 1},

		// The excess blob gas should decrease by however much the target was
		// under-shot, capped at zero.
		{targetBlobGas, targetBlobs, targetBlobGas},
		{targetBlobGas, targetBlobs - 1, targetBlobGas - BlobTxBlobGasPerBlob},
		{targetBlobGas, targetBlobs - 2, targetBlobGas - (2 * BlobTxBlobGasPerBlob)},
		{BlobTxBlobGasPerBlob - 1, targetBlobs - 1, 0},
	}
	for i, tt := range tests {
		blobGasUsed := uint64(tt.blobs) * BlobTxBlobGasPerBlob
		result := config.LatestBlobConfig(config.OsakaCompatibleBlock).CalcExcessBlobGasForEEST(&tt.excess, &blobGasUsed, big.NewInt(1))
		if result != tt.want {
			t.Errorf("test %d: excess blob gas mismatch: have %v, want %v", i, result, tt.want)
		}
	}
}

func TestCalcBlobFee(t *testing.T) {
	tests := []struct {
		baseFee int64
		blobfee int64
	}{
		{0, 0},
		{1, 8},
		{25000000000, 200000000000}, // 25gkei
	}
	for i, tt := range tests {
		config := &ChainConfig{OsakaCompatibleBlock: big.NewInt(0), BlobScheduleConfig: DefaultBlobSchedule}
		config.BlobScheduleConfig.Osaka = GethCancunBlobConfig
		have := CalcBlobFee(big.NewInt(tt.baseFee))
		if have.Int64() != tt.blobfee {
			t.Errorf("test %d: blobfee mismatch: have %v want %v", i, have, tt.blobfee)
		}
	}
}

func TestCalcBlobFeePostOsaka(t *testing.T) {
	zero := big.NewInt(0)

	tests := []struct {
		excessBlobGas uint64
		blobGasUsed   uint64
		blobfee       uint64
		basefee       uint64
		parenttime    big.Int
		headertime    big.Int
	}{
		{5149252, 1310720, 6328900, 30, *big.NewInt(1754904516), *big.NewInt(1754904528)},
		{19251039, 2490368, 21610335, 50, *big.NewInt(1755033204), *big.NewInt(1755033216)},
	}
	for i, tt := range tests {
		config := TestChainConfig.Copy()
		config.OsakaCompatibleBlock = zero
		config.BlobScheduleConfig = &BlobScheduleConfig{
			Osaka: DefaultOsakaBlobConfig,
		}

		have := config.LatestBlobConfig(&tt.headertime).CalcExcessBlobGasForEEST(&tt.excessBlobGas, &tt.blobGasUsed, big.NewInt(int64(tt.basefee)))
		if have != tt.blobfee {
			t.Errorf("test %d: blobfee mismatch: have %v want %v", i, have, tt.blobfee)
		}
	}
}

func TestCalcExcessBlobGasEIP7918(t *testing.T) {
	cfg := TestChainConfig.Copy()
	cfg.OsakaCompatibleBlock = big.NewInt(0)
	cfg.BlobScheduleConfig = &BlobScheduleConfig{
		Osaka: DefaultOsakaBlobConfig,
	}
	var (
		targetBlobs      = cfg.BlobScheduleConfig.Osaka.Target
		blobGasTarget    = uint64(targetBlobs) * BlobTxBlobGasPerBlob
		latestBlobConfig = cfg.LatestBlobConfig(cfg.OsakaCompatibleBlock)
	)

	tests := []struct {
		name          string
		excessBlobGas uint64
		blobGasUsed   int
		baseFee       *big.Int
		wantExcessGas uint64
	}{
		{
			name:          "BelowReservePrice",
			excessBlobGas: 0,
			blobGasUsed:   targetBlobs,
			baseFee:       big.NewInt(1_000_000_000),
			wantExcessGas: blobGasTarget * (uint64(latestBlobConfig.Max) - uint64(latestBlobConfig.Target)) / uint64(latestBlobConfig.Max),
		},
		{
			name:          "AboveReservePrice",
			excessBlobGas: 0,
			blobGasUsed:   targetBlobs,
			baseFee:       big.NewInt(1),
			wantExcessGas: 0,
		},
	}
	for _, tc := range tests {
		blobGasUsed := uint64(tc.blobGasUsed) * BlobTxBlobGasPerBlob
		got := latestBlobConfig.CalcExcessBlobGasForEEST(&tc.excessBlobGas, &blobGasUsed, tc.baseFee)
		if got != tc.wantExcessGas {
			t.Fatalf("%s: excess-blob-gas mismatch – have %d, want %d",
				tc.name, got, tc.wantExcessGas)
		}
	}
}

func TestVerifyEIP4844HeaderForEESTUsesParentValues(t *testing.T) {
	config := MainnetChainConfig.Copy()
	config.OsakaCompatibleBlock = big.NewInt(0)
	config.BlobScheduleConfig.Osaka = GethCancunBlobConfig
	bcfg := config.LatestBlobConfig(config.OsakaCompatibleBlock)

	parentExcessBlobGas := uint64(0)
	parentBlobGasUsed := uint64(bcfg.Target) * BlobTxBlobGasPerBlob
	parentBaseFee := big.NewInt(1)

	// Header's excessBlobGas must be derived from parent values.
	headerExcessBlobGas := bcfg.CalcExcessBlobGasForEEST(&parentExcessBlobGas, &parentBlobGasUsed, parentBaseFee)
	headerBlobGasUsed := uint64(0)

	err := bcfg.VerifyEIP4844HeaderForEEST(
		&parentExcessBlobGas, &parentBlobGasUsed, parentBaseFee,
		&headerExcessBlobGas, &headerBlobGasUsed,
	)
	if err != nil {
		t.Fatalf("expected header verification to pass, got error: %v", err)
	}
}

func TestFakeExponential(t *testing.T) {
	tests := []struct {
		factor      int64
		numerator   int64
		denominator int64
		want        int64
	}{
		// When numerator == 0 the return value should always equal the value of factor
		{1, 0, 1, 1},
		{38493, 0, 1000, 38493},
		{0, 1234, 2345, 0}, // should be 0
		{1, 2, 1, 6},       // approximate 7.389
		{1, 4, 2, 6},
		{1, 3, 1, 16}, // approximate 20.09
		{1, 6, 2, 18},
		{1, 4, 1, 49}, // approximate 54.60
		{1, 8, 2, 50},
		{10, 8, 2, 542}, // approximate 540.598
		{11, 8, 2, 596}, // approximate 600.58
		{1, 5, 1, 136},  // approximate 148.4
		{1, 5, 2, 11},   // approximate 12.18
		{2, 5, 2, 23},   // approximate 24.36
		{1, 50000000, 2225652, 5709098764},
	}
	for i, tt := range tests {
		f, n, d := big.NewInt(tt.factor), big.NewInt(tt.numerator), big.NewInt(tt.denominator)
		original := fmt.Sprintf("%d %d %d", f, n, d)
		have := fakeExponential(f, n, d)
		if have.Int64() != tt.want {
			t.Errorf("test %d: fake exponential mismatch: have %v want %v", i, have, tt.want)
		}
		later := fmt.Sprintf("%d %d %d", f, n, d)
		if original != later {
			t.Errorf("test %d: fake exponential modified arguments: have\n%v\nwant\n%v", i, later, original)
		}
	}
}
