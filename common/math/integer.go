// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from common/math/integer.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package math

import (
	"fmt"
	"strconv"
)

const (
	// Integer limit values.
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

// HexOrDecimal64 marshals uint64 as hex or decimal.
type HexOrDecimal64 uint64

// UnmarshalText implements encoding.TextUnmarshaler.
func (i *HexOrDecimal64) UnmarshalText(input []byte) error {
	int, ok := ParseUint64(string(input))
	if !ok {
		return fmt.Errorf("invalid hex or decimal integer %q", input)
	}
	*i = HexOrDecimal64(int)
	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (i HexOrDecimal64) MarshalText() ([]byte, error) {
	return fmt.Appendf(nil, "%#x", uint64(i)), nil
}

// ParseUint64 parses s as an integer in decimal or hexadecimal syntax.
// Leading zeros are accepted. The empty string parses as zero.
func ParseUint64(s string) (uint64, bool) {
	if s == "" {
		return 0, true
	}
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		v, err := strconv.ParseUint(s[2:], 16, 64)
		return v, err == nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	return v, err == nil
}

// MustParseUint64 parses s as an integer and panics if the string is invalid.
func MustParseUint64(s string) uint64 {
	v, ok := ParseUint64(s)
	if !ok {
		panic("invalid unsigned 64 bit integer: " + s)
	}
	return v
}

// NOTE: The following methods need to be optimised using either bit checking or asm

// SafeSub returns subtraction result and whether overflow occurred.
func SafeSub(x, y uint64) (uint64, bool) {
	return x - y, x < y
}

// SafeAdd returns the result and whether overflow occurred.
func SafeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > MaxUint64-x
}

// SafeMul returns multiplication result and whether overflow occurred.
func SafeMul(x, y uint64) (uint64, bool) {
	if x == 0 || y == 0 {
		return 0, false
	}
	return x * y, y > MaxUint64/x
}
