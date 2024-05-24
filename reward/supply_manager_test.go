// Copyright 2024 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package reward

import (
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/system"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/contracts/contracts/testing/system_contracts"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/storage/statedb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// A test suite with the blockchain having a reward-related history similar to Cypress.
// | Block     	| Fork  	| Minting 	| Ratio    	| KIP82 	| Event                      	|
// |-----------	|-------	|---------	|----------	|-------	|----------------------------	|
// | Genesis   	| None  	| 9.6     	| 34/54/12 	| n/a   	| Launch                     	|
// | Block 100 	| Magma 	| 6.4     	| 50/40/10 	| n/a   	| KGP-3                      	|
// | Block 200 	| Kore  	| 6.4     	| 50/20/30 	| 20/80 	| KGP-6 + KIP-103            	|
// | Block 300 	| Kaia  	| 8.0     	| 50/25/25 	| 10/90 	| KGP-25 + KIP-160 + KIP-162 	|
// | Block 400 is the latest block
type SupplyTestSuite struct {
	suite.Suite

	// Setup per-Suite
	config            *params.ChainConfig
	gov               governanceHelper
	engine            consensus.Engine
	oldStakingManager *StakingManager

	// Setup per-Test
	db      database.DBManager
	genesis *types.Block
	chain   *blockchain.BlockChain
	sm      *supplyManager
}

// ----------------------------------------------------------------------------
// Test cases

// Tests totalSupplyFromState as well as the setupHistory() itself.
// This accounts for (AccMinted - AccBurntFee - kip103Burn - kip160Burn) but not (- zeroBurn - deadBurn).
func (s *SupplyTestSuite) TestFromState() {
	t := s.T()
	s.setupHistory()
	s.sm.Start() // start catchup

	testcases := s.testcases()
	for _, tc := range testcases {
		fromState, err := s.sm.totalSupplyFromState(tc.number)
		require.NoError(t, err)
		bigEqual(t, tc.expectFromState, fromState, tc.number)

		if s.T().Failed() {
			s.dumpState(tc.number)
			break
		}
	}
}

// Test reading canonical burn amounts from the state.
func (s *SupplyTestSuite) TestCanonicalBurn() {
	t := s.T()
	s.setupHistory()

	// Delete state at 199
	root := s.db.ReadBlockByNumber(199).Root()
	s.db.DeleteTrieNode(root.ExtendZero())

	// State unavailable at 199
	zero, dead, err := s.sm.GetCanonicalBurn(199)
	assert.Error(t, err)
	assert.Nil(t, zero)
	assert.Nil(t, dead)

	// State available at 200
	zero, dead, err = s.sm.GetCanonicalBurn(200)
	assert.NoError(t, err)
	assert.Equal(t, "1000000000000000000000000000", zero.String())
	assert.Equal(t, "2000000000000000000000000000", dead.String())
}

// Test reading rebalance memo from the contract.
func (s *SupplyTestSuite) TestRebalanceMemo() {
	t := s.T()
	s.setupHistory()

	// rebalance not configured
	amount, err := s.sm.GetRebalanceBurn(199, nil, common.Address{})
	require.NoError(t, err)
	assert.Equal(t, "0", amount.String())

	// num < forkNum
	amount, err = s.sm.GetRebalanceBurn(199, big.NewInt(200), addrKip103)
	require.NoError(t, err)
	assert.Equal(t, "0", amount.String())

	// num >= forkNum
	amount, err = s.sm.GetRebalanceBurn(200, big.NewInt(200), addrKip103)
	require.NoError(t, err)
	assert.Equal(t, "650511428500000000000", amount.String())

	amount, err = s.sm.GetRebalanceBurn(300, big.NewInt(300), addrKip160)
	require.NoError(t, err)
	assert.Equal(t, "120800000000000000000", amount.String())

	// call failed: bad contract address
	amount, err = s.sm.GetRebalanceBurn(200, big.NewInt(200), addrFund1)
	require.Error(t, err)
	require.Nil(t, amount)
}

// Tests sm.Stop() does not block.
func (s *SupplyTestSuite) TestStop() {
	s.setupHistory() // head block is already 400
	s.sm.Start()     // start catchup, enters big-step mode

	endCh := make(chan struct{})
	go func() {
		s.sm.Stop() // immediately stop. the catchup() returns during big step and has no chance to <-sm.quitCh.
		close(endCh)
	}()

	timer := time.NewTimer(1 * time.Second)
	select {
	case <-timer.C:
		s.T().Fatal("timeout")
	case <-endCh:
		return
	}
}

// Tests the insertBlock -> go catchup, where the catchup enters the big-step mode.
func (s *SupplyTestSuite) TestCatchupBigStep() {
	t := s.T()
	s.setupHistory() // head block is already 400
	s.sm.Start()     // start catchup
	defer s.sm.Stop()
	s.waitAccReward() // Wait for catchup to finish

	testcases := s.testcases()
	for _, tc := range testcases {
		accReward, err := s.sm.GetAccReward(tc.number)
		require.NoError(t, err)
		bigEqual(t, tc.expectTotalSupply.TotalMinted, accReward.Minted, tc.number)
		bigEqual(t, tc.expectTotalSupply.BurntFee, accReward.BurntFee, tc.number)
	}
}

// Tests go catchup -> insertBlock case, where the catchup follows the chain head event.
func (s *SupplyTestSuite) TestCatchupEventSubscription() {
	t := s.T()
	s.sm.Start() // start catchup
	defer s.sm.Stop()
	time.Sleep(10 * time.Millisecond) // yield to the catchup goroutine to start
	s.setupHistory()                  // block is inserted after the catchup started
	s.waitAccReward()

	testcases := s.testcases()
	for _, tc := range testcases {
		accReward, err := s.sm.GetAccReward(tc.number)
		require.NoError(t, err)
		bigEqual(t, tc.expectTotalSupply.TotalMinted, accReward.Minted, tc.number)
		bigEqual(t, tc.expectTotalSupply.BurntFee, accReward.BurntFee, tc.number)
	}
}

// Tests all supply components.
func (s *SupplyTestSuite) TestTotalSupply() {
	t := s.T()
	s.setupHistory()
	s.sm.Start()
	defer s.sm.Stop()
	s.waitAccReward()

	testcases := s.testcases()
	for _, tc := range testcases {
		ts, err := s.sm.GetTotalSupply(tc.number)
		require.NoError(t, err)

		expected := tc.expectTotalSupply
		actual := ts
		bigEqual(t, expected.TotalSupply, actual.TotalSupply, tc.number)
		bigEqual(t, expected.TotalMinted, actual.TotalMinted, tc.number)
		bigEqual(t, expected.TotalBurnt, actual.TotalBurnt, tc.number)
		bigEqual(t, expected.BurntFee, actual.BurntFee, tc.number)
		bigEqual(t, expected.ZeroBurn, actual.ZeroBurn, tc.number)
		bigEqual(t, expected.DeadBurn, actual.DeadBurn, tc.number)
		bigEqual(t, expected.Kip103Burn, actual.Kip103Burn, tc.number)
		bigEqual(t, expected.Kip160Burn, actual.Kip160Burn, tc.number)
	}
}

// Test that when some data are missing, GetTotalSupply leaves some fields nil and returns an error.
func (s *SupplyTestSuite) TestTotalSupplyPartialInfo() {
	t := s.T()
	s.setupHistory()
	s.sm.Start()
	s.waitAccReward()
	s.sm.Stop()

	var num uint64 = 200
	var expected *TotalSupply
	for _, tc := range s.testcases() {
		if tc.number == num {
			expected = tc.expectTotalSupply
			break
		}
	}

	// Missing state trie; no canonical burn amounts
	root := s.db.ReadBlockByNumber(num).Root()
	s.db.DeleteTrieNode(root.ExtendZero())

	ts, err := s.sm.GetTotalSupply(num)
	assert.ErrorContains(t, err, "missing trie node")
	assert.Nil(t, ts.TotalSupply)
	assert.Equal(t, expected.TotalMinted, ts.TotalMinted)
	assert.Nil(t, ts.TotalBurnt)
	assert.Equal(t, expected.BurntFee, ts.BurntFee)
	assert.Nil(t, ts.ZeroBurn)
	assert.Nil(t, ts.DeadBurn)
	assert.Equal(t, expected.Kip103Burn, ts.Kip103Burn)
	assert.Equal(t, expected.Kip160Burn, ts.Kip160Burn)

	// Misconfigured KIP-103
	s.chain.Config().Kip103ContractAddress = addrFund1

	ts, err = s.sm.GetTotalSupply(num)
	assert.ErrorContains(t, err, "missing trie node") // Errors are concatenated
	assert.ErrorContains(t, err, "no contract code")
	assert.Nil(t, ts.TotalSupply)
	assert.Equal(t, expected.TotalMinted, ts.TotalMinted)
	assert.Nil(t, ts.TotalBurnt)
	assert.Equal(t, expected.BurntFee, ts.BurntFee)
	assert.Nil(t, ts.ZeroBurn)
	assert.Nil(t, ts.DeadBurn)
	assert.Nil(t, ts.Kip103Burn)
	assert.Equal(t, expected.Kip160Burn, ts.Kip160Burn)

	// No AccReward
	s.db.WriteLastAccRewardBlockNumber(num - 1)
	s.sm.accRewardCache.Purge()
	ts, err = s.sm.GetTotalSupply(num)
	assert.ErrorIs(t, err, errNoAccReward)
	assert.Nil(t, ts)
}

func (s *SupplyTestSuite) waitAccReward() {
	for i := 0; i < 1000; i++ { // wait 10 seconds until catchup complete
		if s.db.ReadLastAccRewardBlockNumber() >= 400 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if s.db.ReadLastAccRewardBlockNumber() < 400 {
		s.T().Fatal("Catchup not finished in time")
	}
}

func TestSupplyManager(t *testing.T) {
	suite.Run(t, new(SupplyTestSuite))
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
)

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

func (s *SupplyTestSuite) SetupSuite() {
	log.EnableLogForTest(log.LvlCrit, log.LvlError)

	s.config = &params.ChainConfig{
		ChainID:       big.NewInt(31337),
		DeriveShaImpl: 2,
		UnitPrice:     25 * params.Gwei,
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

	s.gov = newSupplyTestGovernance(s.config)
	s.engine = newSupplyTestEngine(s.T(), s.config, s.gov)

	s.oldStakingManager = GetStakingManager()
	SetTestStakingManagerWithStakingInfoCache(&StakingInfo{
		KIFAddr: addrFund1,
		KEFAddr: addrFund2,
	})
}

func (s *SupplyTestSuite) SetupTest() {
	t := s.T()

	s.db = database.NewMemoryDBManager()
	genesis := &blockchain.Genesis{
		Config:     s.config,
		Timestamp:  uint64(time.Now().Unix()),
		BlockScore: common.Big1,
		Alloc: blockchain.GenesisAlloc{
			addrGenesis1: {Balance: bigMult(amount1B, big.NewInt(1))},
			addrGenesis2: {Balance: bigMult(amount1B, big.NewInt(2))},
			addrGenesis3: {Balance: bigMult(amount1B, big.NewInt(3))},
			addrGenesis4: {Balance: bigMult(amount1B, big.NewInt(4))},
			addrKip103:   s.kip103Alloc(),
			addrKip160:   s.kip160Alloc(),
		},
	}
	s.genesis = genesis.MustCommit(s.db)

	cacheConfig := &blockchain.CacheConfig{
		ArchiveMode:         false,
		CacheSize:           512,
		BlockInterval:       128,
		TriesInMemory:       128,
		TrieNodeCacheConfig: statedb.GetEmptyTrieNodeCacheConfig(),
	}
	chain, err := blockchain.NewBlockChain(s.db, cacheConfig, s.config, s.engine, vm.Config{})
	require.NoError(t, err)
	s.chain = chain

	s.sm = NewSupplyManager(s.chain, s.gov, s.db, 1)
}

func (s *SupplyTestSuite) setupHistory() {
	t := s.T()

	var ( // Generate blocks with 1 tx per block. Send 1 wei from Genesis4 to Genesis3.
		signer   = types.LatestSignerForChainID(s.config.ChainID)
		key      = keyGenesis4
		from     = addrGenesis4
		to       = addrGenesis3
		amount   = big.NewInt(1)
		gasLimit = uint64(21000)
		genFunc  = func(i int, b *blockchain.BlockGen) {
			num := b.Number().Uint64()
			var gasPrice *big.Int
			if num < 100 { // Must be equal to UnitPrice before Magma
				gasPrice = big.NewInt(25 * params.Gwei)
			} else if num < 300 { // Double the basefee before kaia
				gasPrice = big.NewInt(50 * params.Gwei)
			} else { // Basefee plus recommended tip since Kaia
				gasPrice = big.NewInt(27 * params.Gwei)
			}
			unsignedTx := types.NewTransaction(b.TxNonce(from), to, amount, gasLimit, gasPrice, nil)
			if unsignedTx != nil {
				signedTx, err := types.SignTx(unsignedTx, signer, key)
				require.NoError(t, err)
				b.AddTx(signedTx)
			}
		}
	)
	blocks, _ := blockchain.GenerateChain(s.config, s.genesis, s.engine, s.db, 400, genFunc)
	require.NotEmpty(t, blocks)

	// Insert s.chain
	_, err := s.chain.InsertChain(blocks)
	require.NoError(t, err)
	expected := blocks[len(blocks)-1]
	actual := s.chain.CurrentBlock()
	assert.Equal(t, expected.Hash(), actual.Hash())
}

type supplyTestTC struct {
	number uint64

	expectTotalSupply *TotalSupply
	expectFromState   *big.Int
}

func (s *SupplyTestSuite) testcases() []supplyTestTC {
	var (
		halfFee = big.NewInt(262500000000000) // Magma burn amount (gasPrice=50, baseFee=25 -> effectivePrice=25)
		fullFee = big.NewInt(525000000000000) // Kore burn amount (gasPrice=50, baseFee=25 -> effectivePrice=25)
		withTip = big.NewInt(567000000000000) // Kaia burn amount (gasPrice=27, baseFee=25 -> effectivePrice=27)

		// Fund1 has 769254566500000000000, Fund2 has 178056862000000000000 at block 199 but burnt
		// Fund1 get   1280000000000000000, Fund2 get   1920000000000000000 at block 200 minted from reward but burnt
		// Fund2 get 100000000000000000000, Fund2 get 200000000000000000000 at block 200 minted from kip103
		// kip103Burn = (769.25 + 178.06 + 1.28 + 1.92) - (100 + 200) = 650.51
		kip103Burn, _ = new(big.Int).SetString("650511428500000000000", 10)

		// Fund1 has 226720000000000000000, Fund2 has 390080000000000000000 at block 299 but burnt
		// Fund1 get   2000000000000000000, Fund2 get   2000000000000000000 at block 300 minted from reward but burnt
		// Fund2 get 200000000000000000000, Fund2 get 300000000000000000000 at block 300 minted from kip160
		// kip160Burn = (226.72 + 390.08 + 2.00 + 2.00) - (200 + 300) = 120.80
		kip160Burn, _ = new(big.Int).SetString("120800000000000000000", 10)

		// Allocated at genesis
		zeroBurn = bigMult(amount1B, big.NewInt(1))
		deadBurn = bigMult(amount1B, big.NewInt(2))
	)
	// accumulated rewards: segment sums
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
	testcases := []supplyTestTC{}
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
			fromState   = bigSub(minted[num], bigAdd(burntFee[num], kip103BurnAtNum, kip160BurnAtNum)) // AccMinted - (AccBurntFee + Kip103Burn + Kip160Burn)
			totalBurnt  = bigAdd(burntFee[num], zeroBurn, deadBurn, kip103BurnAtNum, kip160BurnAtNum)
			totalSupply = bigSub(minted[num], totalBurnt)
		)
		testcases = append(testcases, supplyTestTC{
			number: num,

			expectTotalSupply: &TotalSupply{
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

func (s *SupplyTestSuite) kip103Alloc() blockchain.GenesisAccount {
	// Memo taken from the INFO log of blockchain/system/rebalance.go "successfully executed treasury rebalancing"
	return s.rebalanceAlloc(200, addrKip103, system.Kip103MockCode,
		[]*big.Int{bigMult(amount1, big.NewInt(100)), bigMult(amount1, big.NewInt(200))},
		"{\"retired\":{\"0x000000000000000000000000000000000000a000\":770534566500000000000,\"0x000000000000000000000000000000000000b000\":179976862000000000000},\"newbie\":{\"0x000000000000000000000000000000000000a000\":100000000000000000000,\"0x000000000000000000000000000000000000b000\":200000000000000000000},\"burnt\":650511428500000000000,\"success\":true}",
	)
}

func (s *SupplyTestSuite) kip160Alloc() blockchain.GenesisAccount {
	// Memo taken from the INFO log of blockchain/system/rebalance.go "successfully executed treasury rebalancing"
	return s.rebalanceAlloc(300, addrKip160, system.Kip160MockCode,
		[]*big.Int{bigMult(amount1, big.NewInt(200)), bigMult(amount1, big.NewInt(300))},
		"{\"zeroed\":{\"0x000000000000000000000000000000000000a000\":228720000000000000000,\"0x000000000000000000000000000000000000b000\":392080000000000000000},\"allocated\":{\"0x000000000000000000000000000000000000a000\":200000000000000000000,\"0x000000000000000000000000000000000000b000\":300000000000000000000},\"burnt\":120800000000000000000,\"success\":true}",
	)
}

func (s *SupplyTestSuite) rebalanceAlloc(blockNum uint64, addr common.Address, code []byte, amounts []*big.Int, memo string) blockchain.GenesisAccount {
	// Create a simulated state with the mock contract populated.
	var (
		alloc = blockchain.GenesisAlloc{
			addrGenesis4: {Balance: big.NewInt(params.KAIA)},
			addr:         {Balance: common.Big0, Code: code},
		}
		db          = database.NewMemoryDBManager()
		backend     = backends.NewSimulatedBackendWithDatabase(db, alloc, &params.ChainConfig{})
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
	require.NoError(s.T(), err)
	_, err = contract.TestFinalize( // Set memo before rebalance block number for convenience
		bind.NewKeyedTransactor(keyGenesis4),
		memo,
	)
	require.NoError(s.T(), err)
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

// ----------------------------------------------------------------------------
// Cleanup test

func (s *SupplyTestSuite) TearDownTest() {
	s.chain.Stop()
}

func (s *SupplyTestSuite) TearDownSuite() {
	SetTestStakingManager(s.oldStakingManager)
}

// ----------------------------------------------------------------------------
// Mocks

// Mock governance and consensus engine for testing.
var (
	_ = (governanceHelper)(&supplyTestGovernance{})
	_ = (consensus.Engine)(&supplyTestEngine{})
)

type supplyTestGovernance struct {
	p0 *params.GovParamSet
	p1 *params.GovParamSet
	p2 *params.GovParamSet
	p3 *params.GovParamSet
}

func newSupplyTestGovernance(chainConfig *params.ChainConfig) *supplyTestGovernance {
	gov := new(supplyTestGovernance)
	gov.p0, _ = params.NewGovParamSetChainConfig(chainConfig)

	map1 := gov.p0.IntMap()
	map1[params.MintingAmount] = amount64
	map1[params.Ratio] = "50/40/10"
	gov.p1, _ = params.NewGovParamSetIntMap(map1)

	map2 := gov.p1.IntMap()
	map1[params.MintingAmount] = amount64
	map2[params.Ratio] = "50/20/30"
	map2[params.Kip82Ratio] = "20/80"
	gov.p2, _ = params.NewGovParamSetIntMap(map2)

	map3 := gov.p2.IntMap()
	map3[params.MintingAmount] = amount80
	map3[params.Ratio] = "50/25/25"
	map3[params.Kip82Ratio] = "10/90"
	gov.p3, _ = params.NewGovParamSetIntMap(map3)

	return gov
}

func (g *supplyTestGovernance) CurrentParams() *params.GovParamSet {
	return g.p3
}

func (g *supplyTestGovernance) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	if num < 100 {
		return g.p0, nil
	} else if num < 200 {
		return g.p1, nil
	} else if num < 300 {
		return g.p2, nil
	} else {
		return g.p3, nil
	}
}

// Minimally implements the consensus.Engine interface for testing reward distribution.
type supplyTestEngine struct {
	t      *testing.T
	config *params.ChainConfig
	gov    governanceHelper
}

func newSupplyTestEngine(t *testing.T, config *params.ChainConfig, gov governanceHelper) *supplyTestEngine {
	return &supplyTestEngine{
		t:      t,
		config: config,
		gov:    gov,
	}
}

func (s *supplyTestEngine) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

func (s *supplyTestEngine) CanVerifyHeadersConcurrently() bool {
	return false
}

// Pretend it performed parallel processing and return no error.
func (s *supplyTestEngine) PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	for range headers {
		results <- nil
	}
	return abort, results
}

func (s *supplyTestEngine) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return nil
}

func (s *supplyTestEngine) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return nil, nil
}

func (s *supplyTestEngine) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

func (s *supplyTestEngine) Prepare(chain consensus.ChainReader, header *types.Header) error {
	return nil
}

// Simplfied version of istanbul Finalize for testing native token distribution.
func (s *supplyTestEngine) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	header.BlockScore = common.Big1
	header.Rewardbase = addrProposer

	rules := s.config.Rules(header.Number)
	pset, _ := s.gov.EffectiveParams(header.Number.Uint64())
	rewardSpec, err := CalcDeferredReward(header, txs, receipts, rules, pset)
	if err != nil {
		return nil, err
	}
	DistributeBlockReward(state, rewardSpec.Rewards)

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
}

func (s *supplyTestEngine) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	return nil, nil
}

func (s *supplyTestEngine) CalcBlockScore(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return common.Big0
}

func (s *supplyTestEngine) APIs(chain consensus.ChainReader) []rpc.API { return nil }
func (s *supplyTestEngine) Protocol() consensus.Protocol               { return consensus.Protocol{} }

func (s *supplyTestEngine) CreateSnapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) error {
	return nil
}

func (s *supplyTestEngine) GetConsensusInfo(block *types.Block) (consensus.ConsensusInfo, error) {
	return consensus.ConsensusInfo{}, nil
}

func (s *supplyTestEngine) InitSnapshot() {
	// do nothing
}
