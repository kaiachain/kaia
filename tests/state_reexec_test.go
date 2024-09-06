package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	testcontract "github.com/kaiachain/kaia/contracts/contracts/testing/reward"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/reward"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test State Regeneration (reexecution) after pruning state trie nodes.
// This test ensures that the state regeneration yields the exact same state as the block's stateRoot.
// Post-Kaia engine.Finalize() relies on the state trie to calculate rewards, so the state regeneration
// can be interfered. This test ensures that the state regeneration is robust against such interference.
func TestStateReexec(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	// Test parameters
	var (
		numNodes = 1
		forkNum  = big.NewInt(4)
		nodeId   = bind.NewKeyedTransactor(deriveTestAccount(0)).From
		owner    = bind.NewKeyedTransactor(deriveTestAccount(5))

		config = testStateReexec_config(forkNum)
		alloc  = testStateReexec_alloc(t, owner, nodeId)
	)

	// Start the chain
	ctx, err := newBlockchainTestContext(&blockchainTestOverrides{
		numNodes:    numNodes,
		numAccounts: 8,
		config:      config,
		alloc:       alloc,
	})
	require.Nil(t, err)
	ctx.Subscribe(t, func(ev *blockchain.ChainEvent) {
		b := ev.Block
		t.Logf("block[%3d] stateRoot=%x", b.NumberU64(), b.Header().Root)
	})
	ctx.Start()
	defer ctx.Cleanup()

	ctx.WaitBlock(t, 6)

	// Clear staking cache to force GetStakingInfo post-Kaia to utilize the state trie.
	reward.PurgeStakingInfoCache()
	// Delete state roots to force historical state regeneration
	testStateReexec_prune(t, ctx.nodes[0], []uint64{2, 3, 4, 5})
	// Test state regeneration
	testStateReexec_run(t, ctx.nodes[0], 3)

	// Repeat for post-Kaia block
	reward.PurgeStakingInfoCache()
	testStateReexec_prune(t, ctx.nodes[0], []uint64{2, 3, 4, 5})
	testStateReexec_run(t, ctx.nodes[0], 5) // post-kaia

	// Ensure preloaded staking info are released after use.
	assert.Equal(t, 0, reward.TestGetStakingPreloadSize())
}

func testStateReexec_config(forkNum *big.Int) *params.ChainConfig {
	config := blockchainTestChainConfig.Copy()
	config.LondonCompatibleBlock = common.Big0
	config.IstanbulCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0
	config.ShanghaiCompatibleBlock = common.Big0
	config.CancunCompatibleBlock = common.Big0
	config.KaiaCompatibleBlock = forkNum

	// Use WeightedRandom to test reward distribution based on StakingInfo
	config.Istanbul.ProposerPolicy = uint64(istanbul.WeightedRandom)
	// Set the reward ratio so that reward distribution is different from the 'all to proposer' fallback.
	// If the GetStakingInfo() fails during state regen, the regenerated state would just give all
	// rewards to the proposer, deviating from the actual historical state.
	config.Governance.Reward.Ratio = "34/54/12"
	return config
}

// Create a genesis state with an AddressBookMock
func testStateReexec_alloc(t *testing.T, owner *bind.TransactOpts, nodeId common.Address) blockchain.GenesisAlloc {
	// Create a simulated state with the mock contract populated.
	var (
		abookAddr   = system.AddressBookAddr
		abookCode   = common.FromHex(testcontract.AddressBookMockBinRuntime)
		stakingAddr = common.HexToAddress("0x1000")
		rewardAddr  = common.HexToAddress("0x2000")
		fund1Addr   = common.HexToAddress("0xa000")
		fund2Addr   = common.HexToAddress("0xb000")

		alloc = blockchain.GenesisAlloc{
			owner.From:             {Balance: big.NewInt(params.KAIA)},
			system.AddressBookAddr: {Balance: common.Big0, Code: abookCode},
		}
		db          = database.NewMemoryDBManager()
		backend     = backends.NewSimulatedBackendWithDatabase(db, alloc, &params.ChainConfig{})
		contract, _ = testcontract.NewAddressBookMockTransactor(abookAddr, backend)
	)
	_, err := contract.ConstructContract(owner, []common.Address{owner.From}, common.Big1)
	backend.Commit()
	require.Nil(t, err)

	_, err = contract.RegisterCnStakingContract(owner, nodeId, stakingAddr, rewardAddr)
	require.Nil(t, err)
	_, err = contract.UpdatePocContract(owner, fund1Addr, common.Big1)
	require.Nil(t, err)
	_, err = contract.UpdateKirContract(owner, fund2Addr, common.Big1)
	backend.Commit()
	require.Nil(t, err)

	_, err = contract.ActivateAddressBook(owner)
	backend.Commit()
	require.Nil(t, err)

	// Copy contract storage from the simulated state to the genesis account.
	abookStorage := make(map[common.Hash]common.Hash)
	stateDB, _ := backend.BlockChain().State()
	stateDB.ForEachStorage(abookAddr, func(key common.Hash, value common.Hash) bool {
		abookStorage[key] = value
		return true
	})
	return blockchain.GenesisAlloc{
		abookAddr: {
			Balance: common.Big0,
			Code:    abookCode,
			Storage: abookStorage,
		},
		stakingAddr: {
			Balance: new(big.Int).Mul(big.NewInt(params.KAIA), big.NewInt(5_000_000)),
		},
	}
}

func testStateReexec_prune(t *testing.T, node *blockchainTestNode, nums []uint64) {
	db := node.cn.ChainDB()

	for _, num := range nums {
		block := node.cn.BlockChain().GetBlockByNumber(num)
		root := block.Header().Root
		db.DeleteTrieNode(root.ExtendZero())
	}
}

func testStateReexec_run(t *testing.T, node *blockchainTestNode, num uint64) {
	block := node.cn.BlockChain().GetBlockByNumber(num)

	t.Logf("Regenerating state at block %d", num)
	state, release, err := node.cn.APIBackend.StateAtBlock(context.Background(), block, 10, nil, false, false)
	if release != nil {
		release()
	}
	require.Nil(t, err)

	// Regenerated state must match the stored block's stateRoot
	assert.Equal(t, block.Header().Root, state.IntermediateRoot(false))
}
