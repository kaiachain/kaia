package impl

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/kaiax/gov"
	contractgov_mock "github.com/kaiachain/kaia/kaiax/gov/contractgov/mock"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/stretchr/testify/assert"
)

func newHeaderGovModuleMock(t *testing.T) *headergov_mock.MockHeaderGovModule {
	mock := headergov_mock.NewMockHeaderGovModule(gomock.NewController(t))
	return mock
}

func newContractGovModuleMock(t *testing.T) *contractgov_mock.MockContractGovModule {
	mock := contractgov_mock.NewMockContractGovModule(gomock.NewController(t))
	return mock
}

func TestEffectiveParamSet(t *testing.T) {
	hgm := newHeaderGovModuleMock(t)
	cgm := newContractGovModuleMock(t)
	m := &GovModule{
		hgm: hgm,
		cgm: cgm,
	}

	// default value returned
	{
		defaultVal := uint64(250e9)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, defaultVal, ps.UnitPrice)
	}

	// headergov value returned
	{
		headerGovVal := uint64(123)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[gov.ParamEnum]interface{}{gov.GovernanceUnitPrice: headerGovVal}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[gov.ParamEnum]interface{}{}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, headerGovVal, ps.UnitPrice)
	}

	// contractgov value returned
	{
		contractGovVal := uint64(456)
		hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[gov.ParamEnum]interface{}{gov.GovernanceUnitPrice: 0}, nil)
		cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(map[gov.ParamEnum]interface{}{gov.GovernanceUnitPrice: contractGovVal}, nil)
		ps, _ := m.EffectiveParamSet(1)
		assert.Equal(t, contractGovVal, ps.UnitPrice)
	}
}
