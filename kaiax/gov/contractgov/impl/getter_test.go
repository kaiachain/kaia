package impl

import (
	"math/big"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	govcontract "github.com/kaiachain/kaia/contracts/bindings/gov"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSimulateBackend(t *testing.T) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	// Create accounts and simulated blockchain
	accounts := []*bind.TransactOpts{}
	alloc := blockchain.GenesisAlloc{}
	for range 1 {
		key, _ := crypto.GenerateKey()
		account := bind.NewKeyedTransactor(key)
		account.GasLimit = 10000000
		accounts = append(accounts, account)
		alloc[account.From] = blockchain.GenesisAccount{Balance: big.NewInt(params.KAIA)}
	}
	config := &params.ChainConfig{}
	config.SetDefaults()
	config.UnitPrice = 25e9
	config.IstanbulCompatibleBlock = common.Big0
	config.LondonCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0

	sim := backends.NewSimulatedBackendWithDatabase(database.NewMemoryDBManager(), alloc, config)

	// Deploy contract
	owner := accounts[0]
	address, tx, contract, err := govcontract.DeployGovParam(owner, sim)
	require.Nil(t, err)
	sim.Commit()

	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

	return accounts, sim, address, contract
}

func prepareContractGovModule(t *testing.T, bc *blockchain.BlockChain, addr common.Address) *contractGovModule {
	return prepareContractGovModuleWithConfig(t, bc, addr, &params.ChainConfig{KoreCompatibleBlock: big.NewInt(100)})
}

func prepareContractGovModuleWithConfig(t *testing.T, bc *blockchain.BlockChain, addr common.Address, config *params.ChainConfig) *contractGovModule {
	mockHGM := headergov_mock.NewMockHeaderGovModule(gomock.NewController(t))
	cgm := NewContractGovModule()
	err := cgm.Init(&InitOpts{
		Chain:       bc,
		ChainConfig: config,
		Hgm:         mockHGM,
	})
	require.Nil(t, err)
	mockHGM.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{GovParamContract: addr}).AnyTimes()
	return cgm
}

func setParam(t *testing.T, sim *backends.SimulatedBackend, gp *govcontract.GovParam, owner *bind.TransactOpts, name string, val []byte, activation int64) {
	tx, err := gp.SetParam(owner, name, true, val, big.NewInt(activation))
	require.Nil(t, err)
	sim.Commit()
	receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
	require.NotNil(t, receipt)
	require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)
}

func TestGetContractParamsInvalidBlockNumber(t *testing.T) {
	_, sim, addr, _ := createSimulateBackend(t)
	api := NewContractGovAPI(prepareContractGovModule(t, sim.BlockChain(), addr))

	_, err := api.GetContractParams(rpc.BlockNumber(-3), nil)
	require.ErrorIs(t, err, gov.ErrUnknownBlock)

	future := rpc.BlockNumber(sim.BlockChain().CurrentBlock().NumberU64() + 1)
	_, err = api.GetContractParams(future, nil)
	require.ErrorIs(t, err, gov.ErrUnknownBlock)
}

func TestGetParamSet(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	name := string(gov.GovernanceUnitPrice)
	accounts, sim, addr, gp := createSimulateBackend(t)
	cgm := prepareContractGovModule(t, sim.BlockChain(), addr)

	setParam(t, sim, gp, accounts[0], name, []byte{0, 0, 0, 0, 0, 0, 0, 25}, 1000)
	assert.Equal(t, uint64(25), cgm.GetParamSet(1000).UnitPrice)

	setParam(t, sim, gp, accounts[0], name, []byte{0, 0, 0, 0, 0, 0, 0, 125}, 2000)
	assert.Equal(t, uint64(125), cgm.GetParamSet(2000).UnitPrice)
}

// reward.deferredtxfee (gov.AlwaysDeprecated): always dropped from contract governance.
func TestGetParamSetIgnoresDeprecatedParam(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	accounts, sim, addr, gp := createSimulateBackend(t)
	cgm := prepareContractGovModuleWithConfig(t, sim.BlockChain(), addr,
		&params.ChainConfig{KoreCompatibleBlock: big.NewInt(100), PermissionlessCompatibleBlock: big.NewInt(2000)})

	setParam(t, sim, gp, accounts[0], string(gov.RewardDeferredTxFee), []byte{0x01}, 200)

	deferredTxFee := gov.GetDefaultGovernanceParamSet().DeferredTxFee
	assert.Equal(t, deferredTxFee, cgm.GetParamSet(1500).DeferredTxFee) // ignored
	assert.Equal(t, deferredTxFee, cgm.GetParamSet(2500).DeferredTxFee) // ignored
}

// istanbul.committeesize is in gov.PermissionlessDeprecated: allowed before the Permissionless fork, ignored after.
func TestGetParamSetIgnoresPermissionlessDeprecatedParam(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	accounts, sim, addr, gp := createSimulateBackend(t)
	cgm := prepareContractGovModuleWithConfig(t, sim.BlockChain(), addr,
		&params.ChainConfig{KoreCompatibleBlock: big.NewInt(100), PermissionlessCompatibleBlock: big.NewInt(2000)})

	committeeSize := uint64(50)
	setParam(t, sim, gp, accounts[0], string(gov.IstanbulCommitteeSize), new(big.Int).SetUint64(committeeSize).Bytes(), 200)

	assert.Equal(t, committeeSize, cgm.GetParamSet(1500).CommitteeSize)                                    // before fork: applied
	assert.Equal(t, gov.GetDefaultGovernanceParamSet().CommitteeSize, cgm.GetParamSet(2500).CommitteeSize) // after fork: ignored
}
