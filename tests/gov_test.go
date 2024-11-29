package tests

import (
	"math/big"
	"os"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainnetGenesisGovernance is a regression test for genesis parameters.
// TODO-kaiax: replace mixedEngine with fixed genesis parameters values.
func TestMainnetGenesisGovernance(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	config := params.MainnetChainConfig.Copy()
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Config = config
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)

	fullNode, node, _, _, workspace := newBlockchain(t, config, genesis)
	chain := node.BlockChain().(*blockchain.BlockChain)
	defer os.RemoveAll(workspace)
	defer fullNode.Stop()

	// Create MixedEngine
	dbm := database.NewMemoryDBManager()
	mixed := governance.NewMixedEngine(config, dbm)

	// Get params from MixedEngine
	blockNum := uint64(0)
	mixedParams, err := mixed.EffectiveParams(blockNum)
	assert.NoError(t, err)

	// Get params from GovModule
	govModule := gov_impl.NewGovModule()
	err = govModule.Init(&gov_impl.InitOpts{
		Chain:       chain,
		ChainConfig: config,
		ChainKv:     dbm.GetMemDB(),
	})
	require.NoError(t, err)
	govParams := govModule.EffectiveParamSet(blockNum)

	// Compare all parameter values
	mixedMap := mixedParams.StrMap()
	govMap := govParams.ToMap()
	assert.Equal(t, len(mixedMap), len(govMap))

	for key, mixedVal := range mixedMap {
		govVal, exists := govMap[gov.ParamName(key)]
		assert.True(t, exists, "Key %s missing from GovModule params at block %d", key, blockNum)
		switch key {
		case string(gov.RewardMintingAmount), string(gov.RewardMinimumStake):
			assert.Equal(t, mixedVal, govVal.(*big.Int).String())
		default:
			assert.Equal(t, mixedVal, govVal, "Key %s mismatch", key)
		}
	}
}
