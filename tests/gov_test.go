package tests

import (
	"os"
	"testing"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_impl "github.com/kaiachain/kaia/kaiax/gov/impl"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainnetGenesisGovernance is a regression test for genesis parameters (kaia.getParams(0)).
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

	// Get params from GovModule
	dbm := database.NewMemoryDBManager()
	govModule := gov_impl.NewGovModule()
	err := govModule.Init(&gov_impl.InitOpts{
		Chain:       chain,
		ChainConfig: config,
		ChainKv:     dbm.GetMemDB(),
	})
	require.NoError(t, err)

	genesisParamsMap := map[string]interface{}{
		"governance.deriveshaimpl":        uint64(2),
		"governance.governancemode":       "single",
		"governance.governingnode":        common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"),
		"governance.govparamcontract":     common.Address{},
		"governance.unitprice":            uint64(25000000000),
		"istanbul.committeesize":          uint64(22),
		"istanbul.epoch":                  uint64(604800),
		"istanbul.policy":                 uint64(2),
		"kip71.basefeedenominator":        uint64(20),
		"kip71.gastarget":                 uint64(30000000),
		"kip71.lowerboundbasefee":         uint64(25000000000),
		"kip71.maxblockgasusedforbasefee": uint64(60000000),
		"kip71.upperboundbasefee":         uint64(750000000000),
		"reward.deferredtxfee":            true,
		"reward.kip82ratio":               "20/80",
		"reward.minimumstake":             "5000000",
		"reward.mintingamount":            "9600000000000000000",
		"reward.proposerupdateinterval":   uint64(3600),
		"reward.ratio":                    "34/54/12",
		"reward.stakingupdateinterval":    uint64(86400),
		"reward.useginicoeff":             true,
	}

	pset := govModule.GetParamSet(0)
	govParamsMap := pset.ToMap()
	assert.Equal(t, len(genesisParamsMap), len(govParamsMap))

	for name, val := range genesisParamsMap {
		govVal, exists := govParamsMap[gov.ParamName(name)]
		assert.True(t, exists, "Key %s missing from GovModule params", name)
		assert.Equal(t, val, govVal, "Key %s mismatch", name)
	}
}
