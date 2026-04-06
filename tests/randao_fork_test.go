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
package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/accounts/abi/bind/backends"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/consensus/istanbul"
	kip149contract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/kip149"
	proxycontract "github.com/kaiachain/kaia/contracts/contracts/system_contracts/proxy"
	testcontract "github.com/kaiachain/kaia/contracts/contracts/testing/system_contracts"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRandao(t *testing.T) {
	log.EnableLogForTest(log.LvlCrit, log.LvlWarn)

	cases := []struct {
		name         string
		forkNum      *big.Int
		fromGenesis  bool
		kip113AtFork common.Address
	}{
		{
			name:        "Deploy",
			forkNum:     big.NewInt(3),
			fromGenesis: false,
		},
		{
			name:         "Genesis",
			forkNum:      big.NewInt(0),
			fromGenesis:  true,
			kip113AtFork: common.HexToAddress("0x0000000000000000000000000000000000000403"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			runRandaoScenario(t, tc.forkNum, tc.fromGenesis, tc.kip113AtFork)
		})
	}
}

func runRandaoScenario(t *testing.T, forkNum *big.Int, fromGenesis bool, genesisKip113Addr common.Address) {
	var (
		numNodes    = 1
		owner       = bind.NewKeyedTransactor(deriveTestAccount(5))
		randomAddr  = common.HexToAddress("0x0000000000000000000000000000000000000404")
		kip113Addr  common.Address
		blockPeriod uint64 = 1
	)
	if fromGenesis {
		kip113Addr = genesisKip113Addr
	} else {
		kip113Addr = crypto.CreateAddress(owner.From, uint64(1)) // predicted deployed address.
	}

	config := testRandao_config(forkNum, owner.From, kip113Addr)
	alloc := testRandao_allocRandom(randomAddr)
	if fromGenesis {
		alloc = system.MergeGenesisAlloc(
			alloc,
			testRandao_allocRegistry(owner.From, kip113Addr),
			testRandao_allocKip113(numNodes, owner.From, kip113Addr),
		)
	}

	ctx, err := newBlockchainTestContext(&blockchainTestOverrides{
		numNodes:    numNodes,
		numAccounts: 8,
		config:      config,
		alloc:       alloc,
		blockPeriod: &blockPeriod,
	})
	require.Nil(t, err)
	ctx.Start()
	defer ctx.Cleanup()

	if !fromGenesis {
		// Wait for consensus start and deploy KIP113 before hardfork.
		ctx.WaitBlock(t, 1)
		_, actualKip113Addr := testRandao_deployKip113(t, ctx, owner)
		assert.Equal(t, kip113Addr, actualKip113Addr)
	}

	// Pass the hardfork block, give each CN a chance to propose.
	ctx.WaitBlock(t, forkNum.Uint64()+uint64(numNodes))

	// Inspect the chain.
	testRandao_checkRegistry(t, ctx, owner.From, kip113Addr)
	testRandao_checkKip113(t, ctx)
	testRandao_checkKip114(t, ctx, randomAddr)
}

// Make ChainConfig that hardforks at `forkNum` and the Registry owner be `owner`.
func testRandao_config(forkNum *big.Int, owner, kip113Addr common.Address) *params.ChainConfig {
	config := blockchainTestChainConfig.Copy()
	config.LondonCompatibleBlock = common.Big0
	config.IstanbulCompatibleBlock = common.Big0
	config.EthTxTypeCompatibleBlock = common.Big0
	config.MagmaCompatibleBlock = common.Big0
	config.KoreCompatibleBlock = common.Big0
	config.ShanghaiCompatibleBlock = common.Big0
	config.CancunCompatibleBlock = forkNum
	config.RandaoCompatibleBlock = forkNum

	// Use WeightedRandom to test KIP-146 random proposer selection
	config.Istanbul.ProposerPolicy = uint64(istanbul.WeightedRandom)

	if forkNum.Sign() != 0 {
		// RandaoRegistry is only effective if forkNum > 0
		config.RandaoRegistry = &params.RegistryConfig{
			Records: map[string]common.Address{
				system.Kip113Name: kip113Addr,
			},
			Owner: owner,
		}
	}
	return config
}

// Deploy a small contract to test RANDOM opcode
func testRandao_allocRandom(randomAddr common.Address) blockchain.GenesisAlloc {
	return blockchain.GenesisAlloc{
		randomAddr: {
			// contract Random { function random() external view returns (uint256) { return block.prevrandao; }}  // 0x44 opcode is block.prevrandao in solc 0.8.18+
			Code:    hexutil.MustDecode("0x6080604052348015600f57600080fd5b506004361060285760003560e01c80635ec01e4d14602d575b600080fd5b60336047565b604051603e91906066565b60405180910390f35b600044905090565b6000819050919050565b606081604f565b82525050565b6000602082019050607960008301846059565b9291505056fea2646970667358221220291164179a7b6e34ccb0821e55e26f9202870c95464cde432863dde9ca55426c64736f6c63430008120033"),
			Balance: common.Big0,
		},
	}
}

// RandaoRegistry must be allocated at Genesis if forkNum == 0
func testRandao_allocRegistry(ownerAddr, kip113Addr common.Address) blockchain.GenesisAlloc {
	return blockchain.GenesisAlloc{
		system.RegistryAddr: {
			Code:    system.RegistryCode,
			Balance: common.Big0,
			Storage: system.AllocRegistry(&params.RegistryConfig{
				Records: map[string]common.Address{
					system.Kip113Name: kip113Addr,
				},
				Owner: ownerAddr,
			}),
		},
	}
}

// Allocate the KIP-113 with all node BLS public keys
func testRandao_allocKip113(numNodes int, ownerAddr, kip113Addr common.Address) blockchain.GenesisAlloc {
	infos := make(system.BlsPublicKeyInfos)
	for i := range numNodes {
		var (
			key   = deriveTestAccount(i)
			addr  = crypto.PubkeyToAddress(key.PublicKey)
			sk, _ = bls.DeriveFromECDSA(key)
			pk    = sk.PublicKey().Marshal()
			pop   = bls.PopProve(sk).Marshal()
		)
		infos[addr] = system.BlsPublicKeyInfo{PublicKey: pk, Pop: pop}
	}

	var (
		logicAddr = common.HexToAddress("0x0000000000000000000000000000000000000402")
		owner     = crypto.PubkeyToAddress(deriveTestAccount(5).PublicKey)

		proxyStorage       = system.AllocProxy(logicAddr)
		kip113ProxyStorage = system.AllocKip113Proxy(system.AllocKip113Init{
			Infos: infos,
			Owner: owner,
		})
		kip113LogicStorage = system.AllocKip113Logic()
		storage            = system.MergeStorage(proxyStorage, kip113ProxyStorage)
	)

	return blockchain.GenesisAlloc{
		logicAddr: {
			Code:    system.Kip113MockCode,
			Storage: kip113LogicStorage,
			Balance: common.Big0,
		},
		kip113Addr: {
			Code:    system.ERC1967ProxyCode,
			Storage: storage,
			Balance: common.Big0,
		},
	}
}

// Deploy KIP-113 contract
func testRandao_deployKip113(t *testing.T, ctx *blockchainTestContext, owner *bind.TransactOpts) (*testcontract.KIP113Mock, common.Address) {
	var (
		abi, _      = testcontract.KIP113MockMetaData.GetAbi()
		initData, _ = abi.Pack("initialize")

		chain   = ctx.nodes[0].cn.BlockChain()
		txpool  = ctx.nodes[0].cn.TxPool().(*blockchain.TxPool)
		backend = backends.NewBlockchainContractBackend(chain, txpool, nil)
	)

	startNonce, err := backend.PendingNonceAt(context.Background(), owner.From)
	require.NoError(t, err)

	// Send deploy/register txs with fixed nonces so they can be mined in a single block.
	implAddr, implTx, _, err := testcontract.DeployKIP113Mock(randaoTxOpts(owner, startNonce, 8_000_000), backend)
	require.NoError(t, err)

	proxyAddr, proxyTx, _, err := proxycontract.DeployERC1967Proxy(randaoTxOpts(owner, startNonce+1, 2_000_000), backend, implAddr, initData)
	require.NoError(t, err)

	t.Logf("Kip113 impl=%s proxy=%s", implAddr.Hex(), proxyAddr.Hex())
	kip113, _ := testcontract.NewKIP113Mock(proxyAddr, backend)

	// Register node BLS public keys
	txs := []*types.Transaction{implTx, proxyTx}
	for i := 0; i < ctx.numNodes; i++ {
		var (
			addr  = ctx.accountAddrs[i]
			sk, _ = bls.DeriveFromECDSA(ctx.accountKeys[i])
			pk    = sk.PublicKey().Marshal()
			pop   = bls.PopProve(sk).Marshal()
		)
		t.Logf("node[%2d] addr=%x blsPub=%x", i, addr, pk)

		tx, err := kip113.Register(randaoTxOpts(owner, startNonce+2+uint64(i), 500_000), addr, pk, pop)
		require.NoError(t, err)
		txs = append(txs, tx)
	}

	// Wait for all txs after sending all of them.
	for _, tx := range txs {
		ctx.WaitTx(t, tx.Hash())
	}

	infos, _ := system.ReadKip113All(backend, proxyAddr, nil)
	t.Logf("Kip113 getAllBlsInfo().length=%d", len(infos))

	return kip113, proxyAddr
}

func randaoTxOpts(base *bind.TransactOpts, nonce uint64, gasLimit uint64) *bind.TransactOpts {
	opts := *base
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.GasLimit = gasLimit
	return &opts
}

// Inspect the given chain for Registry contract
func testRandao_checkRegistry(t *testing.T, ctx *blockchainTestContext, ownerAddr, kip113Addr common.Address) {
	var (
		forkNum     = int64(ctx.config.RandaoCompatibleBlock.Uint64())
		bgctx       = context.Background()
		chain       = ctx.nodes[0].cn.BlockChain()
		backend     = backends.NewBlockchainContractBackend(chain, nil, nil)
		registry, _ = kip149contract.NewRegistryCaller(system.RegistryAddr, backend)

		before *big.Int // Largest num without Registry
		after  *big.Int // Smallest num with Registry
	)

	if forkNum == 0 {
		after = common.Big0
	} else {
		before = big.NewInt(forkNum - 1)
		after = big.NewInt(forkNum)
	}

	// Registry code is installed exactly at forkParentNum
	if before != nil {
		code, err := backend.CodeAt(bgctx, system.RegistryAddr, before)
		assert.Nil(t, err)
		assert.Empty(t, code)

		addr, err := system.ReadActiveAddressFromRegistry(backend, system.Kip113Name, before)
		assert.ErrorIs(t, err, system.ErrRegistryNotInstalled)
		assert.Empty(t, addr)
	}

	// Inspect code
	code, err := backend.CodeAt(bgctx, system.RegistryAddr, after)
	assert.Nil(t, err)
	assert.NotNil(t, code)

	// Inspect contract contents
	names, err := registry.GetAllNames(&bind.CallOpts{BlockNumber: after})
	t.Logf("Registry.getAllNames()=%v", names)
	assert.Nil(t, err)
	assert.Equal(t, []string{system.Kip113Name}, names)

	addr, err := registry.GetActiveAddr(&bind.CallOpts{BlockNumber: after}, system.Kip113Name)
	t.Logf("Registry.getActiveAddr('KIP113')=%s", addr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, kip113Addr, addr)

	addr, err = registry.Owner(&bind.CallOpts{BlockNumber: after})
	t.Logf("Registry.owner()=%s", ownerAddr.Hex())
	assert.Nil(t, err)
	assert.Equal(t, ownerAddr, addr)

	// Inspect via system contract accessors
	addr, err = system.ReadActiveAddressFromRegistry(backend, system.Kip113Name, after)
	assert.Nil(t, err)
	assert.Equal(t, kip113Addr, addr)
}

// Inspect the given chain for KIP-113 contract
func testRandao_checkKip113(t *testing.T, ctx *blockchainTestContext) {
	var (
		forkNum = ctx.config.RandaoCompatibleBlock
		chain   = ctx.nodes[0].cn.BlockChain()
		backend = backends.NewBlockchainContractBackend(chain, nil, nil)
	)

	kip113Addr, err := system.ReadActiveAddressFromRegistry(backend, system.Kip113Name, forkNum)
	assert.Nil(t, err)

	// Inspect via system contract accessors
	// BLS public keys of every nodes are registered
	infos, err := system.ReadKip113All(backend, kip113Addr, forkNum)
	t.Logf("Kip113.getAllBlsInfo()=%v", infos.String())
	assert.Nil(t, err)
	assert.Len(t, infos, ctx.numNodes)
	for i := 0; i < ctx.numNodes; i++ {
		addr := ctx.accountAddrs[i]
		assert.Contains(t, infos, addr)
	}
}

// Inspect the given chain for KIP-114 header fields and RANDOM opcode
func testRandao_checkKip114(t *testing.T, ctx *blockchainTestContext, randomAddr common.Address) {
	var (
		chain   = ctx.nodes[0].cn.BlockChain()
		backend = backends.NewBlockchainContractBackend(chain, nil, nil)

		forkNum = ctx.config.RandaoCompatibleBlock.Uint64()
		headNum = chain.CurrentBlock().NumberU64()
	)

	// Call the contract to check RANDOM opcode result
	callRandom := func(num uint64) []byte {
		tx := kaia.CallMsg{
			To:   &randomAddr,
			Data: hexutil.MustDecode("0x5ec01e4d"), // random()
		}
		out, err := backend.CallContract(context.Background(), tx, new(big.Int).SetUint64(num))
		assert.Nil(t, err)
		return out
	}

	for num := uint64(1); num <= headNum; num++ {
		header := chain.GetHeaderByNumber(num)
		require.NotNil(t, header)

		random := callRandom(num)
		t.Logf("block[%3d] opRandom=%x", num, random)

		if num < forkNum {
			assert.Nil(t, header.RandomReveal, num)
			assert.Nil(t, header.MixHash, num)
			assert.Equal(t, header.ParentHash.Bytes(), random, num)
		} else {
			assert.NotNil(t, header.RandomReveal, num)
			assert.NotNil(t, header.MixHash, num)
			assert.Equal(t, header.MixHash, random, num)
		}
	}
}
