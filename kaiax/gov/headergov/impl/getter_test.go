package impl

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/v2/kaiax/gov"
	"github.com/kaiachain/kaia/v2/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/stretchr/testify/assert"
)

func TestEffectiveParams(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	gov := map[uint64]headergov.GovData{
		100: headergov.NewGovData(gov.PartialParamSet{
			gov.GovernanceUnitPrice: uint64(750),
		}),
	}

	testCases := []struct {
		desc          string
		koreBlock     uint64
		blockNum      uint64
		expectedPrice uint64
	}{
		{"Pre-Kore, Block 0", 999999999, 0, 250e9},
		{"Pre-Kore, Block 100", 999999999, 100, 250e9},
		{"Pre-Kore, Block 101", 999999999, 101, 250e9},
		{"Pre-Kore, Block 200", 999999999, 200, 250e9},
		{"Pre-Kore, Block 201", 999999999, 201, 750},

		{"Post-Kore, Block 0", 0, 0, 250e9},
		{"Post-Kore, Block 100", 0, 100, 250e9},
		{"Post-Kore, Block 101", 0, 101, 250e9},
		{"Post-Kore, Block 200", 0, 200, 750},
		{"Post-Kore, Block 201", 0, 201, 750},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := getTestChainConfig()
			config.KoreCompatibleBlock = big.NewInt(int64(tc.koreBlock))
			config.UnitPrice = 250e9
			h := newHeaderGovModule(t, config)
			for num, g := range gov {
				h.HandleGov(num, g)
			}

			gp := h.GetParamSet(tc.blockNum)
			assert.Equal(t, tc.expectedPrice, gp.UnitPrice)
		})
	}
}

func TestEffectiveParamsPartial(t *testing.T) {
	config := getTestChainConfigKore()
	h := newHeaderGovModule(t, config)

	testCases := []struct {
		blockNum      uint64
		expectedPrice uint64
	}{
		{0, 0},
		{100, 0},
		{200, 1},
		{300, 2},
		{400, 3},
	}

	for i, govPrice := range []uint64{0, 1, 2, 3} {
		h.AddGov(uint64(i*100), headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: govPrice}))
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Block %d", tc.blockNum), func(t *testing.T) {
			gp := h.GetPartialParamSet(tc.blockNum)
			assert.Equal(t, tc.expectedPrice, gp[gov.GovernanceUnitPrice])
		})
	}
}

func TestPrevEpochStart(t *testing.T) {
	type TestCase struct {
		blockNum    uint64
		expectedGov uint64
	}

	preKoreTcs := []TestCase{
		{0, 0},
		{99, 0},
		{100, 0},
		{101, 0},
		{199, 0},
		{200, 0},
		{201, 100},
		{299, 100},
		{300, 100},
		{301, 200},
	}

	for _, tc := range preKoreTcs {
		t.Run(fmt.Sprintf("Pre-Kore Block %d", tc.blockNum), func(t *testing.T) {
			result := PrevEpochStart(tc.blockNum, 100, false)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}

	postKoreTcs := []TestCase{
		{0, 0},
		{99, 0},
		{100, 0},
		{101, 0},
		{199, 0},
		{200, 100},
		{201, 100},
		{299, 100},
		{300, 200},
		{301, 200},
	}

	for _, tc := range postKoreTcs {
		t.Run(fmt.Sprintf("Post-Kore Block %d", tc.blockNum), func(t *testing.T) {
			result := PrevEpochStart(tc.blockNum, 100, true)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}
}
