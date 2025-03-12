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
	govcontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/gov"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	headergov_mock "github.com/kaiachain/kaia/kaiax/gov/headergov/mock"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createSimulateBackend(t *testing.T) ([]*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *govcontract.GovParam) {
	// Create accounts and simulated blockchain
	accounts := []*bind.TransactOpts{}
	alloc := blockchain.GenesisAlloc{}
	for i := 0; i < 1; i++ {
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
	mockHGM := headergov_mock.NewMockHeaderGovModule(gomock.NewController(t))
	cgm := NewContractGovModule()
	err := cgm.Init(&InitOpts{
		Chain:       bc,
		ChainConfig: &params.ChainConfig{KoreCompatibleBlock: big.NewInt(100)},
		Hgm:         mockHGM,
	})
	require.Nil(t, err)
	mockHGM.EXPECT().GetParamSet(gomock.Any()).Return(gov.ParamSet{GovParamContract: addr}).AnyTimes()
	return cgm
}

func TestGetParamSet(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)
	paramName := string(gov.GovernanceUnitPrice)
	accounts, sim, addr, gp := createSimulateBackend(t)
	cgm := prepareContractGovModule(t, sim.BlockChain(), addr)

	{
		activation := big.NewInt(1000)
		val := []byte{0, 0, 0, 0, 0, 0, 0, 25}
		tx, err := gp.SetParam(accounts[0], paramName, true, val, activation)
		require.Nil(t, err)
		sim.Commit()

		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		ps := cgm.GetParamSet(activation.Uint64())
		assert.NotNil(t, ps)
		assert.Equal(t, uint64(25), ps.UnitPrice)
	}

	{
		activation := big.NewInt(2000)
		val := []byte{0, 0, 0, 0, 0, 0, 0, 125}
		tx, err := gp.SetParam(accounts[0], paramName, true, val, activation)
		require.Nil(t, err)
		sim.Commit()

		receipt, _ := sim.TransactionReceipt(nil, tx.Hash())
		require.NotNil(t, receipt)
		require.Equal(t, types.ReceiptStatusSuccessful, receipt.Status)

		ps := cgm.GetParamSet(activation.Uint64())
		assert.NotNil(t, ps)
		assert.Equal(t, uint64(125), ps.UnitPrice)
	}
}
