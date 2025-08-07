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
// This file is derived from rpc/types_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package rpc

import (
	"encoding/json"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/math"
)

func TestBlockNumberJSONUnmarshal(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected BlockNumber
	}{
		0:  {`"0x"`, true, BlockNumber(0)},
		1:  {`"0x0"`, false, BlockNumber(0)},
		2:  {`"0X1"`, false, BlockNumber(1)},
		3:  {`"0x00"`, true, BlockNumber(0)},
		4:  {`"0x01"`, true, BlockNumber(0)},
		5:  {`"0x1"`, false, BlockNumber(1)},
		6:  {`"0x12"`, false, BlockNumber(18)},
		7:  {`"0x7fffffffffffffff"`, false, BlockNumber(math.MaxInt64)},
		8:  {`"0x8000000000000000"`, true, BlockNumber(0)},
		9:  {"0", false, BlockNumber(0)}, // NOTE: This case is rpc.EarliestBlockNumber in Kaia
		10: {`"ff"`, true, BlockNumber(0)},
		11: {`"pending"`, false, PendingBlockNumber},
		12: {`"latest"`, false, LatestBlockNumber},
		13: {`"earliest"`, false, EarliestBlockNumber},
		14: {`"safe"`, false, LatestBlockNumber},      // NOTE: Kaia treats "safe" as "latest"
		15: {`"finalized"`, false, LatestBlockNumber}, // NOTE: Kaia treats "finalized" as "latest"
		16: {`someString`, true, BlockNumber(0)},
		17: {`""`, true, BlockNumber(0)},
		18: {``, true, BlockNumber(0)},

		// NOTE: Kaia originil tests
		19: {"1", false, BlockNumber(1)},
		20: {"10", false, BlockNumber(10)},
		21: {"80000000", false, BlockNumber(80000000)},
		22: {"-1", true, BlockNumber(0)},
		23: {"-2", true, BlockNumber(0)},
		24: {"-3", true, BlockNumber(0)},
		25: {"-4", true, BlockNumber(0)},
		26: {"-5", true, BlockNumber(0)},
	}

	for i, test := range tests {
		var num BlockNumber
		err := json.Unmarshal([]byte(test.input), &num)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}
		if num != test.expected {
			t.Errorf("Test %d got unexpected value, want %d, got %d", i, test.expected, num)
		}
	}
}

func TestBlockNumberOrHash_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected BlockNumberOrHash
	}{
		0:  {`"0x"`, true, BlockNumberOrHash{}},
		1:  {`"0x0"`, false, NewBlockNumberOrHashWithNumber(0)},
		2:  {`"0X1"`, false, NewBlockNumberOrHashWithNumber(1)},
		3:  {`"0x00"`, true, BlockNumberOrHash{}},
		4:  {`"0x01"`, true, BlockNumberOrHash{}},
		5:  {`"0x1"`, false, NewBlockNumberOrHashWithNumber(1)},
		6:  {`"0x12"`, false, NewBlockNumberOrHashWithNumber(18)},
		7:  {`"0x7fffffffffffffff"`, false, NewBlockNumberOrHashWithNumber(math.MaxInt64)},
		8:  {`"0x8000000000000000"`, true, BlockNumberOrHash{}},
		9:  {"0", false, NewBlockNumberOrHashWithNumber(0)}, // NOTE: This case is rpc.EarliestBlockNumber in Kaia
		10: {`"ff"`, true, BlockNumberOrHash{}},
		11: {`"pending"`, false, NewBlockNumberOrHashWithNumber(PendingBlockNumber)},
		12: {`"latest"`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)},
		13: {`"earliest"`, false, NewBlockNumberOrHashWithNumber(EarliestBlockNumber)},
		14: {`"safe"`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)},      // NOTE: Kaia treats "safe" as "latest"
		15: {`"finalized"`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)}, // NOTE: Kaia treats "finalized" as "latest"
		16: {`someString`, true, BlockNumberOrHash{}},
		17: {`""`, true, BlockNumberOrHash{}},
		18: {``, true, BlockNumberOrHash{}},
		19: {`"0x0000000000000000000000000000000000000000000000000000000000000000"`, false, NewBlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		20: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, false, NewBlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		21: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":false}`, false, NewBlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		22: {`{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":true}`, false, NewBlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), true)},
		23: {`{"blockNumber":"0x1"}`, false, NewBlockNumberOrHashWithNumber(1)},
		24: {`{"blockNumber":"pending"}`, false, NewBlockNumberOrHashWithNumber(PendingBlockNumber)},
		25: {`{"blockNumber":"latest"}`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)},
		26: {`{"blockNumber":"earliest"}`, false, NewBlockNumberOrHashWithNumber(EarliestBlockNumber)},
		27: {`{"blockNumber":"safe"}`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)},      // NOTE: Kaia treats "safe" as "latest"
		28: {`{"blockNumber":"finalized"}`, false, NewBlockNumberOrHashWithNumber(LatestBlockNumber)}, // NOTE: Kaia treats "finalized" as "latest"
		29: {`{"blockNumber":"0x1", "blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, true, BlockNumberOrHash{}},

		// NOTE: Kaia originil tests
		30: {`"-1"`, true, BlockNumberOrHash{}},
		31: {`"-2"`, true, BlockNumberOrHash{}},
		32: {`"-3"`, true, BlockNumberOrHash{}},
		33: {`{"blockNumber":"-1"}`, true, BlockNumberOrHash{}},
		34: {`{"blockNumber":"-2"}`, true, BlockNumberOrHash{}},
		35: {`{"blockNumber":"-3"}`, true, BlockNumberOrHash{}},
	}

	for i, test := range tests {
		var bnh BlockNumberOrHash
		err := json.Unmarshal([]byte(test.input), &bnh)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}
		hash, hashOk := bnh.Hash()
		expectedHash, expectedHashOk := test.expected.Hash()
		num, numOk := bnh.Number()
		expectedNum, expectedNumOk := test.expected.Number()
		if bnh.RequireCanonical != test.expected.RequireCanonical ||
			hash != expectedHash || hashOk != expectedHashOk ||
			num != expectedNum || numOk != expectedNumOk {
			t.Errorf("Test %d got unexpected value, want %v, got %v", i, test.expected, bnh)
		}
	}
}
