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
	tm.prepareMockExpectHeader(blockNum-1, nil, nil, n[1])
	tm.prepareMockExpectGovParam(blockNum-1, uint64(WeightedRandom), testSubGroupSize, tgn)
	tm.prepareMockExpectStakingInfo(blockNum-1, []uint64{0, 1, 2, 3}, []uint64{aM, aM, aL, aL})

	valCtx, err := newValSetContext(vModule, blockNum)
	assert.NoError(t, err)

	cList, err := vModule.GetCouncilAddressList(blockNum - 1)
	assert.NoError(t, err)

	c, err := newCouncil(valCtx, cList)
	assert.NoError(t, err)

	assert.Equal(t, valset.AddressList{n[0], n[1]}, c.qualifiedValidators)
	assert.Equal(t, valset.AddressList{n[5], n[3], n[2], n[4]}, c.demotedValidators)
}
