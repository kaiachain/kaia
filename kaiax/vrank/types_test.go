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

package vrank

import (
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/stretchr/testify/assert"
)

func TestEncodeDecodeCfReport(t *testing.T) {
	cases := []struct {
		name   string
		report CfReport
	}{
		{
			name:   "addresses",
			report: CfReport{common.HexToAddress("0x15d34AAf54267DB7D7cC839724318F2730aC377B"), common.HexToAddress("0x9965507D1a55bcC2695C58ba16FB37d819D0A4DC")},
		},
		{
			name:   "empty",
			report: CfReport{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			enc, err := EncodeCfReport(tc.report)
			assert.NoError(t, err)
			if len(tc.report) == 0 {
				assert.Nil(t, enc)
				return
			}
			assert.NotEmpty(t, enc)
			dec, err := DecodeCfReport(enc)
			assert.NoError(t, err)
			assert.Equal(t, tc.report, dec)
		})
	}
}
