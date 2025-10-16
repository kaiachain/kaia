package impl

import (
	"math/big"
	"reflect"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/kaiax/gov"
	contractgov_mock "github.com/kaiachain/kaia/kaiax/gov/contractgov/mock"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
)

func newKaiaAPI(t *testing.T, chainConfig *params.ChainConfig) (*KaiaAPI, *mocks.MockBlockChain, *headergov_mock.MockHeaderGovModule, *contractgov_mock.MockContractGovModule) {
	h, mockChain, mockHgm, mockCgm := newGovModule(t, chainConfig)
	return NewKaiaAPI(h), mockChain, mockHgm, mockCgm
}

func TestAPI_kaia_getChainConfig(t *testing.T) {
	var (
		latestGenesisConfig              = params.MainnetChainConfig.Copy()
		api, mockChain, mockHgm, mockCgm = newKaiaAPI(t, latestGenesisConfig)
		num                              = rpc.BlockNumber(0)
	)

	mockHgm.EXPECT().GetPartialParamSet(uint64(0)).Return(gov.PartialParamSet{}).Times(1)
	mockCgm.EXPECT().GetPartialParamSet(uint64(0)).Return(gov.PartialParamSet{}).Times(1)
	mockChain.EXPECT().Config().Return(latestGenesisConfig).Times(2) // isKore, getChainConfig

	// Set all *CompatibleBlock fields to zero.
	{
		v := reflect.ValueOf(latestGenesisConfig).Elem()
		ty := v.Type()

		for i := 0; i < ty.NumField(); i++ {
			field := ty.Field(i)
			if strings.HasSuffix(field.Name, "CompatibleBlock") {
				fieldValue := v.Field(i)
				if fieldValue.Type() == reflect.TypeFor[*big.Int]() {
					fieldValue.Set(reflect.ValueOf(big.NewInt(0)))
				}
			}
		}
	}

	// Check if there is any missing Hardfork block in the API result
	{
		cc := api.GetChainConfig(&num)
		v := reflect.ValueOf(cc).Elem()
		ty := v.Type()

		for i := 0; i < ty.NumField(); i++ {
			field := ty.Field(i)
			if strings.HasSuffix(field.Name, "CompatibleBlock") {
				fieldValue := v.Field(i)
				if fieldValue.Type() == reflect.TypeFor[*big.Int]() {
					assert.Equal(t, fieldValue.Interface().(*big.Int).String(), big.NewInt(0).String())
				}
			}
		}
	}
}
