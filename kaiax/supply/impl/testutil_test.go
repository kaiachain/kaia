// Copyright 2024 The Kaia Authors
// This file is part of the Kaia library.
//
// The Kaia library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Kaia library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Kaia library. If not, see <http://www.gnu.org/licenses/>.

package supply

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/consensus/istanbul"
	consensus_mock "github.com/kaiachain/kaia/consensus/mocks"
	"github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/gov"
	gov_mock "github.com/kaiachain/kaia/kaiax/gov/mock"
	"github.com/kaiachain/kaia/kaiax/reward"
	reward_impl "github.com/kaiachain/kaia/kaiax/reward/impl"
	"github.com/kaiachain/kaia/kaiax/staking"
	staking_mock "github.com/kaiachain/kaia/kaiax/staking/mock"
	"github.com/kaiachain/kaia/kaiax/supply"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// ----------------------------------------------------------------------------
// SupplyTestSuite

// A test suite with the blockchain having a reward-related history similar to Mainnet.
// | Block     	| Fork  	| Minting 	| Ratio    	| KIP82 	| Event                      	|
// |-----------	|-------	|---------	|----------	|-------	|----------------------------	|
// | Genesis   	| None  	| 9.6     	| 34/54/12 	| n/a   	| Launch                     	|
// | Block 100 	| Magma 	| 6.4     	| 50/40/10 	| n/a   	| KGP-3                      	|
// | Block 200 	| Kore  	| 6.4     	| 50/20/30 	| 20/80 	| KGP-6 + KIP-103            	|
// | Block 300 	| Kaia  	| 8.0     	| 50/25/25 	| 10/90 	| KGP-25 + KIP-160 + KIP-162 	|
// | Block 400 is the latest block
type SupplyTestSuite struct {
	suite.Suite

	mockCtrl *gomock.Controller
	dbm      database.DBManager
	chain    *blockchain.BlockChain
	s        *SupplyModule
	blocks   []*types.Block
}

// ----------------------------------------------------------------------------
// Setup test
var (
	amount1, _  = new(big.Int).SetString("1000000000000000000", 10)
	amount64, _ = new(big.Int).SetString("6400000000000000000", 10)
	amount96, _ = new(big.Int).SetString("9600000000000000000", 10)
	amount80, _ = new(big.Int).SetString("8000000000000000000", 10)
	amount1B, _ = new(big.Int).SetString("1000000000000000000000000000", 10)

	keyGenesis4, _ = crypto.HexToECDSA("ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	addrGenesis1   = common.HexToAddress("0x0000") // zero address
	addrGenesis2   = common.HexToAddress("0xdead") // burn address
	addrGenesis3   = common.HexToAddress("0x3000")
	addrGenesis4   = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	addrProposer   = common.HexToAddress("0x9000")
	addrFund1      = common.HexToAddress("0xa000")
	addrFund2      = common.HexToAddress("0xb000")
	addrKip103     = common.HexToAddress("0xff1030")
	addrKip160     = common.HexToAddress("0xff1600")

	testConfig = &params.ChainConfig{
		ChainID:       big.NewInt(31337),
		DeriveShaImpl: 2,
		UnitPrice:     25 * params.Gkei,
		Governance: &params.GovernanceConfig{
			GovernanceMode: "none",
			Reward: &params.RewardConfig{
				MintingAmount:          amount96,
				Ratio:                  "34/54/12",
				Kip82Ratio:             "20/80",
				UseGiniCoeff:           true,
				DeferredTxFee:          true,
				StakingUpdateInterval:  60,
				ProposerUpdateInterval: 30,
				MinimumStake:           big.NewInt(5_000_000),
			},
			KIP71: params.GetDefaultKIP71Config(),
		},
		Istanbul: &params.IstanbulConfig{
			Epoch:          10,
			ProposerPolicy: uint64(istanbul.WeightedRandom),
			SubGroupSize:   30,
		},

		IstanbulCompatibleBlock:  big.NewInt(0),
		LondonCompatibleBlock:    big.NewInt(0),
		EthTxTypeCompatibleBlock: big.NewInt(0),
		MagmaCompatibleBlock:     big.NewInt(100),
		KoreCompatibleBlock:      big.NewInt(200),
		ShanghaiCompatibleBlock:  big.NewInt(300),
		CancunCompatibleBlock:    big.NewInt(300),
		KaiaCompatibleBlock:      big.NewInt(300),

		Kip103CompatibleBlock: big.NewInt(200),
		Kip103ContractAddress: addrKip103,
		Kip160CompatibleBlock: big.NewInt(300),
		Kip160ContractAddress: addrKip160,
	}
)

func (s *SupplyTestSuite) SetupTest() {
	t := s.T()
	log.EnableLogForTest(log.LvlCrit, log.LvlError)

	var (
		mockCtrl = gomock.NewController(t)
		config   = testConfig.Copy()
		genesis  = makeGenesis(t, config)

		engine   = consensus_mock.NewMockEngine(mockCtrl)
		mGov     = gov_mock.NewMockGovModule(mockCtrl)
		mStaking = staking_mock.NewMockStakingModule(mockCtrl)
		mReward  = reward_impl.NewRewardModule()
		mSupply  = NewSupplyModule()

		dbm         = database.NewMemoryDBManager()
		cacheConfig = &blockchain.CacheConfig{
			ArchiveMode:         false,
			CacheSize:           512,
			BlockInterval:       128,
			TriesInMemory:       128,
			TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
		}
	)
	genesisBlock := genesis.MustCommit(dbm)
	chain, err := blockchain.NewBlockChain(dbm, cacheConfig, config, engine, vm.Config{})
	require.NoError(t, err)

	setupMockEngine(engine, chain, mReward)
	setupMockGov(mGov, config)
	setupMockStaking(mStaking)

	mReward.Init(&reward_impl.InitOpts{
		ChainConfig:   config,
		Chain:         chain,
		GovModule:     mGov,
		StakingModule: mStaking,
	})
	mSupply.Init(&InitOpts{
		ChainKv:      dbm.GetMiscDB(),
		ChainConfig:  config,
		Chain:        chain,
		RewardModule: mReward,
	})
	chain.RegisterExecutionModule(mSupply)

	// With every modules ready, generate the blocks
	blocks := setupBlocks(config, genesisBlock, engine, dbm)

	s.mockCtrl = mockCtrl
	s.dbm = dbm
	s.chain = chain
	s.s = mSupply
	s.blocks = blocks
}

func (s *SupplyTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *SupplyTestSuite) insertBlocks() {
	t := s.T()

	// Insert blocks to chain, triggers ExecutionModule.PostInsertBlock.
	for _, block := range s.blocks {
		_, err := s.chain.InsertChain([]*types.Block{block})
		require.NoError(t, err)
	}
	expected := s.blocks[len(s.blocks)-1]
	actual := s.chain.CurrentBlock()
	assert.Equal(t, expected.Hash(), actual.Hash())
}

func makeGenesis(t *testing.T, config *params.ChainConfig) *blockchain.Genesis {
	// Memo taken from the INFO log of blockchain/system/rebalance.go "successfully executed treasury rebalancing"
	kip103Alloc := rebalanceAlloc(t, 200, addrKip103, system.Kip103MockCode,
		[]*big.Int{bigMult(amount1, big.NewInt(100)), bigMult(amount1, big.NewInt(200))},
		"{\"retirees\":[{\"retired\":\"0x000000000000000000000000000000000000a000\",\"balance\":770534566500000000000},{\"retired\":\"0x000000000000000000000000000000000000b000\",\"balance\":179976862000000000000}],\"newbies\":[{\"newbie\":\"0x000000000000000000000000000000000000a000\",\"fundAllocated\":100000000000000000000},{\"newbie\":\"0x000000000000000000000000000000000000b000\",\"fundAllocated\":200000000000000000000}],\"burnt\":650511428500000000000,\"success\":true}",
	)
	kip160Alloc := rebalanceAlloc(t, 300, addrKip160, system.Kip160MockCode,
		[]*big.Int{bigMult(amount1, big.NewInt(300)), bigMult(amount1, big.NewInt(400))},
		`{"before": {"allocated": {"0x000000000000000000000000000000000000a000": 226720000000000000000, "0x000000000000000000000000000000000000b000": 390080000000000000000}, "zeroed": {"0x000000000000000000000000000000000000a000": 226720000000000000000, "0x000000000000000000000000000000000000b000": 390080000000000000000}},"after": {"allocated": {"0x000000000000000000000000000000000000a000": 300000000000000000000, "0x000000000000000000000000000000000000b000": 400000000000000000000}, "zeroed": {"0x000000000000000000000000000000000000a000": 0, "0x000000000000000000000000000000000000b000": 0}}, "burnt": -79200000000000000000, "success": true}`,
	)

	return &blockchain.Genesis{
		Config:     config,
		Timestamp:  uint64(time.Now().Unix()),
		BlockScore: common.Big1,
		Alloc: blockchain.GenesisAlloc{
			addrGenesis1: {Balance: bigMult(amount1B, big.NewInt(1))},
			addrGenesis2: {Balance: bigMult(amount1B, big.NewInt(2))},
			addrGenesis3: {Balance: bigMult(amount1B, big.NewInt(3))},
			addrGenesis4: {Balance: bigMult(amount1B, big.NewInt(4))},
			addrKip103:   kip103Alloc,
			addrKip160:   kip160Alloc,
		},
	}
}

func rebalanceAlloc(t *testing.T, blockNum uint64, addr common.Address, code []byte, amounts []*big.Int, memo string) blockchain.GenesisAccount {
	// Create a simulated state with the mock contract populated.
	var (
		alloc = blockchain.GenesisAlloc{
			addrGenesis4: {Balance: big.NewInt(params.KAIA)},
			addr:         {Balance: common.Big0, Code: code},
		}
		dbm         = database.NewMemoryDBManager()
		backend     = backends.NewSimulatedBackendWithDatabase(dbm, alloc, &params.ChainConfig{})
		contract, _ = system_contracts.NewTreasuryRebalanceMockV2Transactor(addr, backend)
	)
	_, err := contract.TestSetAll(
		bind.NewKeyedTransactor(keyGenesis4),
		[]common.Address{addrFund1, addrFund2}, // zeroed
		[]common.Address{addrFund1, addrFund2}, // allocated to the same address for simplicity
		amounts,
		new(big.Int).SetUint64(blockNum),
		system.EnumRebalanceStatus_Approved,
	)
	require.NoError(t, err)
	_, err = contract.TestFinalize( // Set memo before rebalance block number for convenience
		bind.NewKeyedTransactor(keyGenesis4),
		memo,
	)
	require.NoError(t, err)
	backend.Commit()

	// Copy contract storage from the simulated state to the genesis account.
	storage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(addr, func(key common.Hash, value common.Hash) bool {
		storage[key] = value
		return true
	})
	return blockchain.GenesisAccount{
		Balance: common.Big0,
		Code:    code,
		Storage: storage,
	}
}

func setupMockEngine(engine *consensus_mock.MockEngine, chain *blockchain.BlockChain, mReward reward.RewardModule) {
	engine.EXPECT().Author(gomock.Any()).Return(addrProposer, nil).AnyTimes()
	engine.EXPECT().CalcBlockScore(gomock.Any(), gomock.Any(), gomock.Any()).Return(common.Big0).AnyTimes()
	engine.EXPECT().CanVerifyHeadersConcurrently().Return(false).AnyTimes()

	engine.EXPECT().PreprocessHeaderVerification(gomock.Any()).DoAndReturn(
		func(headers []*types.Header) (chan<- struct{}, <-chan error) {
			abort := make(chan struct{})
			results := make(chan error, len(headers))
			for range headers {
				results <- nil
			}
			return abort, results
		},
	).AnyTimes()

	engine.EXPECT().VerifyHeader(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()
	engine.EXPECT().CreateSnapshot(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	engine.EXPECT().Finalize(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
			header.BlockScore = common.Big1
			header.Rewardbase = addrProposer
			if err := mReward.FinalizeHeader(header, state, txs, receipts); err != nil {
				return nil, err
			}

			if chain.Config().IsKIP103ForkBlock(header.Number) || chain.Config().IsKIP160ForkBlock(header.Number) {
				result, err := system.RebalanceTreasury(state, chain, header)
				_ = result
				// resultJson, _ := json.MarshalIndent(result, "", "  ")
				// s.t.Log("RebalanceTreasury", string(resultJson), err)
				if err != nil {
					return nil, err
				}
			}

			header.Root = state.IntermediateRoot(true)
			return types.NewBlock(header, txs, receipts), nil
		},
	).AnyTimes()
}

func setupMockGov(mGov *gov_mock.MockGovModule, config *params.ChainConfig) {
	mGov.EXPECT().EffectiveParamSet(gomock.Any()).DoAndReturn(func(num uint64) (gov.ParamSet, error) {
		p0 := gov.ParamSet{
			ProposerPolicy: uint64(config.Istanbul.ProposerPolicy),
			UnitPrice:      config.UnitPrice,
			MintingAmount:  config.Governance.Reward.MintingAmount,
			MinimumStake:   config.Governance.Reward.MinimumStake,
			DeferredTxFee:  config.Governance.Reward.DeferredTxFee,
			Ratio:          config.Governance.Reward.Ratio,
			Kip82Ratio:     config.Governance.Reward.Kip82Ratio,
		}
		p100 := gov.ParamSet{
			ProposerPolicy: uint64(config.Istanbul.ProposerPolicy),
			UnitPrice:      config.UnitPrice,
			MintingAmount:  amount64,
			MinimumStake:   config.Governance.Reward.MinimumStake,
			DeferredTxFee:  config.Governance.Reward.DeferredTxFee,
			Ratio:          "50/40/10",
			Kip82Ratio:     config.Governance.Reward.Kip82Ratio,
		}
		p200 := gov.ParamSet{
			ProposerPolicy: uint64(config.Istanbul.ProposerPolicy),
			UnitPrice:      config.UnitPrice,
			MintingAmount:  amount64,
			MinimumStake:   config.Governance.Reward.MinimumStake,
			DeferredTxFee:  config.Governance.Reward.DeferredTxFee,
			Ratio:          "50/20/30",
			Kip82Ratio:     "20/80",
		}
		p300 := gov.ParamSet{
			ProposerPolicy: uint64(config.Istanbul.ProposerPolicy),
			UnitPrice:      config.UnitPrice,
			MintingAmount:  amount80,
			MinimumStake:   config.Governance.Reward.MinimumStake,
			DeferredTxFee:  config.Governance.Reward.DeferredTxFee,
			Ratio:          "50/25/25",
			Kip82Ratio:     "10/90",
		}

		if num < 100 {
			return p0, nil
		} else if num < 200 {
			return p100, nil
		} else if num < 300 {
			return p200, nil
		}
		return p300, nil
	}).AnyTimes()
}

func setupMockStaking(mStaking *staking_mock.MockStakingModule) {
	mStaking.EXPECT().GetStakingInfo(gomock.Any()).Return(&staking.StakingInfo{
		KIFAddr: addrFund1,
		KEFAddr: addrFund2,
	}, nil).AnyTimes()
}

func setupBlocks(config *params.ChainConfig, genesisBlock *types.Block, engine *consensus_mock.MockEngine, dbm database.DBManager) []*types.Block {
	var ( // Generate blocks with 1 tx per block. Send 1 kei from Genesis4 to Genesis3.
		signer   = types.LatestSignerForChainID(config.ChainID)
		key      = keyGenesis4
		from     = addrGenesis4
		to       = addrGenesis3
		amount   = big.NewInt(1)
		gasLimit = uint64(21000)
		genFunc  = func(i int, b *blockchain.BlockGen) {
			num := b.Number().Uint64()
			var gasPrice *big.Int
			if num < 100 { // Must be equal to UnitPrice before Magma
				gasPrice = big.NewInt(25 * params.Gkei)
			} else if num < 300 { // Double the basefee before kaia
				gasPrice = big.NewInt(50 * params.Gkei)
			} else { // Basefee plus recommended tip since Kaia
				gasPrice = big.NewInt(27 * params.Gkei)
			}
			unsignedTx := types.NewTransaction(b.TxNonce(from), to, amount, gasLimit, gasPrice, nil)
			if unsignedTx != nil {
				signedTx, err := types.SignTx(unsignedTx, signer, key)
				if err != nil {
					panic(err)
				}
				b.AddTx(signedTx)
			}
		}
	)
	blocks, _ := blockchain.GenerateChain(config, genesisBlock, engine, dbm, 400, genFunc)
	return blocks
}

func bigEqual(t *testing.T, expected, actual *big.Int, msg ...interface{}) {
	assert.Equal(t, expected.String(), actual.String(), msg...)
}

func bigAdd(nums ...*big.Int) *big.Int {
	result := new(big.Int)
	for _, num := range nums {
		result.Add(result, num)
	}
	return result
}

func bigSub(a, b *big.Int) *big.Int {
	return new(big.Int).Sub(a, b)
}

func bigMult(a, b *big.Int) *big.Int {
	return new(big.Int).Mul(a, b)
}

type supplyTC struct {
	number            uint64
	expectTotalSupply *supply.TotalSupply
	expectFromState   *big.Int
}

func (s *SupplyTestSuite) testcases() []supplyTC {
	var (
		halfFee = big.NewInt(262500000000000) // Magma burn amount (gasPrice=50, baseFee=25 -> effectivePrice=25)
		fullFee = big.NewInt(525000000000000) // Kore burn amount (gasPrice=50, baseFee=25 -> effectivePrice=25)
		withTip = big.NewInt(567000000000000) // Kaia burn amount (gasPrice=27, baseFee=25 -> effectivePrice=27)

		// Fund1 has 769254566500000000000, Fund2 has 178056862000000000000 at block 199 but burnt
		// Fund1 get   1280000000000000000, Fund2 get   1920000000000000000 at block 200 minted from reward but burnt
		// Fund1 get 100000000000000000000, Fund2 get 200000000000000000000 at block 200 minted from kip103
		// kip103Burn = (769.25 + 178.06 + 1.28 + 1.92) - (100 + 200) = 650.51
		kip103Burn, _ = new(big.Int).SetString("650511428500000000000", 10)

		// Fund1 has 226720000000000000000, Fund2 has 390080000000000000000 at block 299 but burnt
		// Fund1 get   2000000000000000000, Fund2 get   2000000000000000000 at block 300 minted from reward but burnt
		// Fund1 get 300000000000000000000, Fund2 get 400000000000000000000 at block 300 minted from kip160
		// kip160Burn = (226.72 + 390.08 + 2.00 + 2.00) - (300 + 400) = -79.2
		kip160Burn, _ = new(big.Int).SetString("-79200000000000000000", 10)

		// Allocated at genesis
		zeroBurn = bigMult(amount1B, big.NewInt(1))
		deadBurn = bigMult(amount1B, big.NewInt(2))
	)
	// supply checkpoints: segment sums
	minted := make(map[uint64]*big.Int)
	burntFee := make(map[uint64]*big.Int)

	minted[0] = bigMult(amount1B, big.NewInt(10)) // Genesis: 1B + 2B + 3B + 4B = 10B
	burntFee[0] = common.Big0

	for i := uint64(1); i <= 99; i++ {
		minted[i] = bigAdd(minted[i-1], amount96)
		burntFee[i] = bigAdd(burntFee[i-1], common.Big0)
	}
	for i := uint64(100); i <= 199; i++ {
		minted[i] = bigAdd(minted[i-1], amount64)
		burntFee[i] = bigAdd(burntFee[i-1], halfFee)
	}
	for i := uint64(200); i <= 299; i++ {
		minted[i] = bigAdd(minted[i-1], amount64)
		burntFee[i] = bigAdd(burntFee[i-1], fullFee)
	}
	for i := uint64(300); i <= 400; i++ {
		minted[i] = bigAdd(minted[i-1], amount80)
		burntFee[i] = bigAdd(burntFee[i-1], withTip)
	}

	nums := []uint64{0, 1, 99, 100, 199, 200, 299, 300, 399, 400}
	testcases := []supplyTC{}
	for _, num := range nums {
		kip103BurnAtNum := common.Big0
		kip160BurnAtNum := common.Big0
		if num >= 200 {
			kip103BurnAtNum = kip103Burn
		}
		if num >= 300 {
			kip160BurnAtNum = kip160Burn
		}

		var (
			// fromState = AccMinted - (AccBurntFee + Kip103Burn + Kip160Burn)
			fromState   = bigSub(minted[num], bigAdd(burntFee[num], kip103BurnAtNum, kip160BurnAtNum))
			totalBurnt  = bigAdd(burntFee[num], zeroBurn, deadBurn, kip103BurnAtNum, kip160BurnAtNum)
			totalSupply = bigSub(minted[num], totalBurnt)
		)
		testcases = append(testcases, supplyTC{
			number: num,
			expectTotalSupply: &supply.TotalSupply{
				TotalSupply: totalSupply,
				TotalMinted: minted[num],
				TotalBurnt:  totalBurnt,
				BurntFee:    burntFee[num],
				ZeroBurn:    zeroBurn,
				DeadBurn:    deadBurn,
				Kip103Burn:  kip103BurnAtNum,
				Kip160Burn:  kip160BurnAtNum,
			},
			expectFromState: fromState,
		})
	}
	return testcases
}

func (s *SupplyTestSuite) dumpState(num uint64) {
	header := s.chain.GetHeaderByNumber(num)
	stateDB, _ := s.chain.StateAt(header.Root)
	dump := stateDB.RawDump()

	out := fmt.Sprintf("num=%d root=%s\n", num, header.Root.Hex())
	keys := []string{}
	for k := range dump.Accounts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := dump.Accounts[k]
		if common.HexToAddress(k) == addrGenesis4 {
			initial := bigMult(amount1B, big.NewInt(4))
			left, _ := new(big.Int).SetString(v.Balance, 10)
			out += fmt.Sprintf("  %s: %30s -%s\n", k, v.Balance, initial.Sub(initial, left))
		} else {
			out += fmt.Sprintf("  %s: %30s\n", k, v.Balance)
		}
	}
	s.T().Log(out)
}
