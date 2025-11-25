package impl

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	contractgov_mock "github.com/kaiachain/kaia/kaiax/gov/contractgov/mock"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/work/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newKaiaAPI(t *testing.T, chainConfig *params.ChainConfig) (*KaiaAPI, *mocks.MockBlockChain, *headergov_mock.MockHeaderGovModule, *contractgov_mock.MockContractGovModule) {
	h, mockChain, mockHgm, mockCgm := newGovModule(t, chainConfig)
	return NewKaiaAPI(h), mockChain, mockHgm, mockCgm
}

func generateRandomValue(name gov.ParamName) any {
	param := gov.Params[name]
	var randomValue any

	switch param.DefaultValue.(type) {
	case int:
		randomValue = rand.Int()
	case int64:
		randomValue = rand.Int63()
	case uint64:
		if name == gov.IstanbulPolicy || name == gov.GovernanceDeriveShaImpl {
			randomValue = uint64(rand.Intn(3))
		} else {
			randomValue = rand.Uint64()
		}
	case bool:
		// Alternate true/false by hash of name
		randomValue = rand.Intn(2) == 0
	case string:
		if name == gov.GovernanceGovernanceMode {
			randomValue = []string{"none", "single"}[rand.Intn(2)]
		} else if name == gov.RewardRatio {
			a := rand.Intn(101)
			b := rand.Intn(101 - a)
			c := 100 - a - b
			randomValue = fmt.Sprintf("%d/%d/%d", a, b, c)
		} else if name == gov.RewardKip82Ratio {
			a := rand.Intn(101)
			b := 100 - a
			randomValue = fmt.Sprintf("%d/%d", a, b)
		} else {
			panic(fmt.Sprintf("unknown string type %s", name))
		}
	case *big.Int:
		randomValue = big.NewInt(rand.Int63())
	case common.Address:
		randomValue = common.BytesToAddress(common.MakeRandomBytes(common.AddressLength))
	default:
		panic(fmt.Errorf("unknown type %s: %T", name, param.DefaultValue))
	}
	return randomValue
}

// TestAPI_kaia_getChainConfig_CompatibleBlocksCompleteness ensures that
// no hardfork block is omitted from the API output.
// In other words, the test fails if API has not been synced
// after a new hardfork block has been introduced.
func TestAPI_kaia_getChainConfig_CompatibleBlocksCompleteness(t *testing.T) {
	var (
		config                           = params.MainnetChainConfig.Copy()
		api, mockChain, mockHgm, mockCgm = newKaiaAPI(t, config)
		apiArg                           = rpc.BlockNumber(0)
		hardforkBlock                    = big.NewInt(123456789)
		configType                       = reflect.TypeFor[params.ChainConfig]()
	)

	mockHgm.EXPECT().GetPartialParamSet(uint64(0)).Return(gov.PartialParamSet{}).Times(1)
	mockCgm.EXPECT().GetPartialParamSet(uint64(0)).Return(gov.PartialParamSet{}).Times(1)
	mockChain.EXPECT().Config().Return(config).Times(2) // isKore, getChainConfig

	// Set all *CompatibleBlock fields to hardforkBlock
	v := reflect.ValueOf(config).Elem()
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if strings.HasSuffix(field.Name, "CompatibleBlock") {
			fieldValue := v.Field(i)
			if fieldValue.Type() == reflect.TypeFor[*big.Int]() {
				fieldValue.Set(reflect.ValueOf(hardforkBlock))
			}
		}
	}

	// Call the API
	apiChainConfig := api.GetChainConfig(&apiArg)
	require.NotNil(t, apiChainConfig, "kaia_getChainConfig should not return nil")

	// Check the all hardfork block numbers are the hardforkBlock
	apiStruct := reflect.ValueOf(apiChainConfig).Elem()
	for i := 0; i < configType.NumField(); i++ {
		fieldName := configType.Field(i).Name
		if strings.HasSuffix(fieldName, "CompatibleBlock") {
			apiValue := apiStruct.FieldByName(fieldName).Interface().(*big.Int).String()
			assert.Equal(t, hardforkBlock.String(), apiValue, "kaia_getChainConfig did not return %s", fieldName)
		}
	}
}

// TestAPI_kaia_getChainConfig_ParameterCompleteness ensures that
// no governance parameter is omitted from the API output.
// In other words, the test fails if API has not been synced
// after a new governance parameter (`gov.Params`) has been introduced.
func TestAPI_kaia_getChainConfig_ParameterCompleteness(t *testing.T) {
	var (
		config                           = params.MainnetChainConfig.Copy()
		api, mockChain, mockHgm, mockCgm = newKaiaAPI(t, config)
		apiArg                           = rpc.BlockNumber(0)
	)

	// run multiple times to increase reliability (with random ChainConfig values)
	for range 100 {
		latestParams := gov.PartialParamSet{}
		for name := range gov.Params {
			err := latestParams.Add(string(name), generateRandomValue(name))
			require.NoError(t, err, "failed to add %s to PartialParamSet", name)
		}

		mockHgm.EXPECT().GetPartialParamSet(uint64(0)).Return(latestParams).Times(1)
		mockCgm.EXPECT().GetPartialParamSet(uint64(0)).Return(gov.PartialParamSet{}).Times(1)
		mockChain.EXPECT().Config().Return(config).Times(2) // isKore, getChainConfig

		// Call the API
		apiChainConfig := api.GetChainConfig(&apiArg)
		require.NotNil(t, apiChainConfig, "kaia_getChainConfig should not return nil")

		// Check the values are the latest
		for name, param := range gov.Params {
			expectedValue := latestParams[name]
			apiValue, err := param.ChainConfigValue(apiChainConfig)
			require.NoError(t, err, "failed to get chain config value for %s", name)
			require.Equal(t, expectedValue, apiValue, "kaia_getChainConfig did not return %s", name)
		}
	}
}
