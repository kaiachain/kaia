package impl

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/kaiax/gov/headergov"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
)

func TestEffectiveParams(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlDebug)
	epoch := uint64(1000)
	gov := map[uint64]headergov.GovData{
		1000: headergov.NewGovData(gov.PartialParamSet{
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
		{"Pre-Kore, Block 1000", 999999999, 1000, 250e9},
		{"Pre-Kore, Block 1001", 999999999, 1001, 250e9},
		{"Pre-Kore, Block 2000", 999999999, 2000, 250e9},
		{"Pre-Kore, Block 2001", 999999999, 2001, 750},

		{"Post-Kore, Block 0", 0, 0, 250e9},
		{"Post-Kore, Block 1000", 0, 1000, 250e9},
		{"Post-Kore, Block 1001", 0, 1001, 250e9},
		{"Post-Kore, Block 2000", 0, 2000, 750},
		{"Post-Kore, Block 2001", 0, 2001, 750},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			config := &params.ChainConfig{
				KoreCompatibleBlock: big.NewInt(int64(tc.koreBlock)),
				Istanbul: &params.IstanbulConfig{
					Epoch: epoch,
				},
			}
			h := newHeaderGovModule(t, config)
			for num, g := range gov {
				h.HandleGov(num, g)
			}

			gp := h.EffectiveParamSet(tc.blockNum)
			assert.Equal(t, tc.expectedPrice, gp.UnitPrice)
		})
	}
}

func TestEffectiveParamsPartial(t *testing.T) {
	var (
		epoch = uint64(1000)
		h     = newHeaderGovModule(t, &params.ChainConfig{
			KoreCompatibleBlock: big.NewInt(0),
			Istanbul:            &params.IstanbulConfig{Epoch: epoch},
		})
	)

	testCases := []struct {
		blockNum      uint64
		expectedPrice uint64
	}{
		{0, 0},
		{1000, 0},
		{2000, 1},
		{3000, 2},
		{4000, 3},
	}

	for i, govPrice := range []uint64{0, 1, 2, 3} {
		h.AddGov(uint64(i*1000), headergov.NewGovData(gov.PartialParamSet{gov.GovernanceUnitPrice: govPrice}))
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Block %d", tc.blockNum), func(t *testing.T) {
			gp := h.EffectiveParamsPartial(tc.blockNum)
			assert.Equal(t, tc.expectedPrice, gp[gov.GovernanceUnitPrice])
		})
	}
}

func TestPrevEpochStart(t *testing.T) {
	epoch := uint64(1000)
	type TestCase struct {
		blockNum    uint64
		expectedGov uint64
	}

	preKoreTcs := []TestCase{
		{0, 0},
		{999, 0},
		{1000, 0},
		{1001, 0},
		{1999, 0},
		{2000, 0},
		{2001, 1000},
		{2999, 1000},
		{3000, 1000},
		{3001, 2000},
	}

	for _, tc := range preKoreTcs {
		t.Run(fmt.Sprintf("Pre-Kore Block %d", tc.blockNum), func(t *testing.T) {
			result := PrevEpochStart(tc.blockNum, epoch, false)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}

	postKoreTcs := []TestCase{
		{0, 0},
		{999, 0},
		{1000, 0},
		{1001, 0},
		{1999, 0},
		{2000, 1000},
		{2001, 1000},
		{2999, 1000},
		{3000, 2000},
		{3001, 2000},
	}

	for _, tc := range postKoreTcs {
		t.Run(fmt.Sprintf("Post-Kore Block %d", tc.blockNum), func(t *testing.T) {
			result := PrevEpochStart(tc.blockNum, epoch, true)
			assert.Equal(t, tc.expectedGov, result, "Incorrect governance data block for block %d", tc.blockNum)
		})
	}
}
