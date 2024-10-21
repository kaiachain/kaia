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

package reward

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRatio(t *testing.T) {
	ratio, err := NewRewardRatio("34/54/12")
	assert.NoError(t, err)
	assert.Equal(t, &RewardRatio{34, 54, 12}, ratio)

	ratio, err = NewRewardRatio("100/0/0")
	assert.NoError(t, err)
	assert.Equal(t, &RewardRatio{100, 0, 0}, ratio)

	_, err = NewRewardRatio("")
	assert.Error(t, err)

	_, err = NewRewardRatio("34/54/12/1")
	assert.Error(t, err)

	_, err = NewRewardRatio("99/88/77")
	assert.Error(t, err)

	_, err = NewRewardRatio("-1/50/51")
	assert.Error(t, err)

	ratio = &RewardRatio{50, 25, 25}
	mintingAmount, _ := new(big.Int).SetString("9600000000000000000", 10)
	g, x, y := ratio.Split(mintingAmount)
	assert.Equal(t, "4800000000000000000", g.String())
	assert.Equal(t, "2400000000000000000", x.String())
	assert.Equal(t, "2400000000000000000", y.String())
}

func TestKip82Ratio(t *testing.T) {
	ratio, err := NewRewardKip82Ratio("20/80")
	assert.NoError(t, err)
	assert.Equal(t, &RewardKip82Ratio{20, 80}, ratio)

	ratio, err = NewRewardKip82Ratio("100/0")
	assert.NoError(t, err)
	assert.Equal(t, &RewardKip82Ratio{100, 0}, ratio)

	_, err = NewRewardKip82Ratio("")
	assert.Error(t, err)

	_, err = NewRewardKip82Ratio("20/80/0")
	assert.Error(t, err)

	_, err = NewRewardKip82Ratio("101/0")
	assert.Error(t, err)

	_, err = NewRewardKip82Ratio("133/-33")
	assert.Error(t, err)

	ratio = &RewardKip82Ratio{20, 80}
	gcAmount, _ := new(big.Int).SetString("4800000000000000000", 10)
	p, s := ratio.Split(gcAmount)
	assert.Equal(t, "960000000000000000", p.String())
	assert.Equal(t, "3840000000000000000", s.String())
}
