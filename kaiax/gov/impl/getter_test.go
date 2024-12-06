package impl

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/kaiax/gov"
	contractgov_mock "github.com/kaiachain/kaia/kaiax/gov/contractgov/mock"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	blockchain_mock "github.com/kaiachain/kaia/kaiax/gov/impl/mock"
	"github.com/kaiachain/kaia/params"
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

func newGovModuleMock(t *testing.T, config *params.ChainConfig) (*headergov_mock.MockHeaderGovModule, *contractgov_mock.MockContractGovModule, *GovModule) {
	hgm := newHeaderGovModuleMock(t)
	cgm := newContractGovModuleMock(t)
	chain := blockchain_mock.NewMockBlockChain(gomock.NewController(t))
	chain.EXPECT().Config().Return(config).AnyTimes()
	m := NewGovModule()

	// don't use Init, we need to use mock Hgm and Cgm.
	m.Chain = chain
	m.ChainConfig = config
	m.Hgm = hgm
	m.Cgm = cgm
	return hgm, cgm, m
}

func TestEffectiveParamSet(t *testing.T) {
	var (
		defaultVal     = uint64(250e9)
		headerGovVal   = uint64(123)
		contractGovVal = uint64(456)
	)

	t.Run("pre-kore", func(t *testing.T) {
		hgm, cgm, m := newGovModuleMock(t, &params.ChainConfig{KoreCompatibleBlock: nil})
		t.Run("default", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil)
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil)
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, defaultVal, ps.UnitPrice)
		})

		t.Run("headergov", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: headerGovVal})
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{})
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, headerGovVal, ps.UnitPrice)
		})

		t.Run("contractgov ignored", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: headerGovVal})
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: contractGovVal})
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, headerGovVal, ps.UnitPrice)
		})
	})

	t.Run("post-kore", func(t *testing.T) {
		hgm, cgm, m := newGovModuleMock(t, &params.ChainConfig{KoreCompatibleBlock: big.NewInt(0)})

		t.Run("default", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil)
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(nil)
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, defaultVal, ps.UnitPrice)
		})

		t.Run("headergov", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: headerGovVal})
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{})
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, headerGovVal, ps.UnitPrice)
		})

		t.Run("contractgov", func(t *testing.T) {
			hgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: headerGovVal})
			cgm.EXPECT().EffectiveParamsPartial(gomock.Any()).Return(gov.PartialParamSet{gov.GovernanceUnitPrice: contractGovVal})
			ps := m.EffectiveParamSet(1)
			assert.Equal(t, contractGovVal, ps.UnitPrice)
		})
	})
}
