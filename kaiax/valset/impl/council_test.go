package impl

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/stretchr/testify/assert"
)

func TestNewCouncil(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vModule, tm, err := newTestVModule(ctrl)
	assert.NoError(t, err)

	blockNum := testIstanbulCompatibleNumber.Uint64()
	tm.prepareMockExpectGovParam(blockNum, uint64(WeightedRandom), testSubGroupSize, tgn)
	tm.prepareMockExpectStakingInfo(blockNum, []uint64{0, 1, 2, 3, 4, 5}, []uint64{aM, aM, 0, 0, aL, aL})

	c, err := newCouncil(vModule, blockNum)
	assert.NoError(t, err)

	assert.Equal(t, valset.AddressList{n[0], n[1]}, c.qualifiedValidators)
	assert.Equal(t, valset.AddressList{n[3], n[2]}, c.demotedValidators)
}
