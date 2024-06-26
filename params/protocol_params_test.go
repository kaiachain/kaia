// Copyright 2024 The Kaia Authors
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
	"testing"

	"github.com/stretchr/testify/assert"
)

var testCaseWithBinaryCodeFormat = []struct {
	cf                     CodeInfo
	getCodeFormatExpectVal CodeFormat
	validateExpectVal      bool
	getVmVersionExpectVal  VmVersion
	stringExpectVal        string
}{
	{0b00000000, CodeFormatEVM, true, VmVersion0, "CodeFormatEVM"},
	{0b00010000, CodeFormatEVM, true, VmVersion1, "CodeFormatEVM"},
	{0b00000001, CodeFormatLast, false, VmVersion0, "UndefinedCodeFormat"},
	{0b00010001, CodeFormatLast, false, VmVersion1, "UndefinedCodeFormat"},
}

func TestGetCodeFormat(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.getCodeFormatExpectVal, tc.cf.GetCodeFormat())
	}
}

func TestString(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.stringExpectVal, tc.cf.GetCodeFormat().String())
	}
}

func TestValidate(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.validateExpectVal, tc.cf.GetCodeFormat().Validate())
	}
}

func TestGetVmVersion(t *testing.T) {
	for _, tc := range testCaseWithBinaryCodeFormat {
		assert.Equal(t, tc.getVmVersionExpectVal, tc.cf.GetVmVersion())
	}
}

func TestSetVmVersion(t *testing.T) {
	testCaseSetIstanbulHFField := []struct {
		cf                        CodeFormat
		isDeployedAfterIstanbulHF bool
		ifDeployedAfterLondonHF   bool
		expectCi                  CodeInfo
	}{
		{CodeFormatEVM, false, false, 0b00000000},
		{CodeFormatEVM, true, false, 0b00010000},
		{CodeFormatEVM, true, true, 0b00010000},
	}
	for _, tc := range testCaseSetIstanbulHFField {
		assert.Equal(t, tc.expectCi, NewCodeInfoWithRules(tc.cf, Rules{IsIstanbul: tc.isDeployedAfterIstanbulHF}))
	}
}
