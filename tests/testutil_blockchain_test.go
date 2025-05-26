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
	"crypto/ecdsa"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/consensus/istanbul"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/crypto/bls"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/node/cn"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip32"
	"golang.org/x/crypto/pbkdf2"
)

// Full blockchain test context.
// TODO: replace newBlockchain()
type blockchainTestContext struct {
	numNodes     int
	accountKeys  []*ecdsa.PrivateKey
	accountAddrs []common.Address
	accounts     []*bind.TransactOpts // accounts[0:numNodes] are node keys
	config       *params.ChainConfig
	genesis      *blockchain.Genesis

	workspace string
	nodes     []*blockchainTestNode
}

type blockchainTestNode struct {
	datadir string
	node    *node.Node
	cn      *cn.CN
}

type blockchainTestOverrides struct {
	numNodes    int                     // default: 1
	numAccounts int                     // default: numNodes
	config      *params.ChainConfig     // default: blockchainTestChainConfig
	alloc       blockchain.GenesisAlloc // default: 10_000_000 KAIA for each account
}

var blockchainTestChainConfig = &params.ChainConfig{
	ChainID:       big.NewInt(31337),
	DeriveShaImpl: 2,
	UnitPrice:     25 * params.Gkei,
	Governance: &params.GovernanceConfig{
		GovernanceMode: "none",
		Reward: &params.RewardConfig{
			MintingAmount:          big.NewInt(params.KAIA * 6.4),
			Ratio:                  "100/0/0",
			Kip82Ratio:             "20/80",
			UseGiniCoeff:           false,
			DeferredTxFee:          true,
			StakingUpdateInterval:  60,
			ProposerUpdateInterval: 30,
			MinimumStake:           big.NewInt(5_000_000),
		},
	},
	Istanbul: &params.IstanbulConfig{
		Epoch:          120,
		ProposerPolicy: uint64(istanbul.RoundRobin),
		SubGroupSize:   100,
	},
}

func newBlockchainTestContext(overrides *blockchainTestOverrides) (*blockchainTestContext, error) {
	if overrides == nil {
		overrides = &blockchainTestOverrides{}
	}
	if overrides.numNodes == 0 {
		overrides.numNodes = 1
	}
	if overrides.numAccounts == 0 {
		overrides.numAccounts = overrides.numNodes
	}
	if overrides.numAccounts < overrides.numNodes {
		return nil, errors.New("numAccounts less than numNodes")
	}
	if overrides.config == nil {
		overrides.config = blockchainTestChainConfig
	}
	if overrides.alloc == nil {
		overrides.alloc = make(blockchain.GenesisAlloc)
	}

	ctx := &blockchainTestContext{
		numNodes: overrides.numNodes,
	}
	ctx.setAccounts(overrides.numAccounts)
	ctx.setConfig(overrides.config)
	ctx.setGenesis(overrides.alloc)
	ctx.setWorkspace()
	err := ctx.setNodes(ctx.numNodes)
	return ctx, err
}

func (ctx *blockchainTestContext) setAccounts(count int) {
	ctx.accountKeys = make([]*ecdsa.PrivateKey, count)
	ctx.accountAddrs = make([]common.Address, count)
	ctx.accounts = make([]*bind.TransactOpts, count)
	for i := 0; i < count; i++ {
		privateKey := deriveTestAccount(i)
		ctx.accountKeys[i] = privateKey
		ctx.accountAddrs[i] = crypto.PubkeyToAddress(privateKey.PublicKey)
		ctx.accounts[i] = bind.NewKeyedTransactor(privateKey)
	}
}

func (ctx *blockchainTestContext) setConfig(config *params.ChainConfig) {
	ctx.config = config.Copy()
	ctx.config.Istanbul.SubGroupSize = uint64(ctx.numNodes)
}

func (ctx *blockchainTestContext) setGenesis(alloc blockchain.GenesisAlloc) {
	// Genesis ExtraData from nodeAddrs
	extra, _ := rlp.EncodeToBytes(&types.IstanbulExtra{
		Validators:    ctx.accountAddrs[:ctx.numNodes],
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	})
	vanity := make([]byte, types.IstanbulExtraVanity)

	// Genesis Alloc from overrides.alloc + rich accountAddrs
	richBalance := new(big.Int).Mul(big.NewInt(params.KAIA), big.NewInt(10_000_000))
	for _, addr := range ctx.accountAddrs {
		account := alloc[addr]
		account.Balance = richBalance
		alloc[addr] = account
	}

	ctx.genesis = &blockchain.Genesis{
		Config:     ctx.config,
		Timestamp:  uint64(time.Now().Unix()),
		ExtraData:  append(vanity, extra...),
		BlockScore: common.Big1,
		Alloc:      alloc,
	}
}

func (ctx *blockchainTestContext) setWorkspace() {
	workspace, _ := os.MkdirTemp("", "kaia-test-state-")
	ctx.workspace = workspace
}

func (ctx *blockchainTestContext) setNodes(numNodes int) error {
	ctx.nodes = make([]*blockchainTestNode, numNodes)
	for i := 0; i < numNodes; i++ {
		if err := ctx.setNode(i); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *blockchainTestContext) setNode(nodeIndex int) (err error) {
	tn := &blockchainTestNode{}
	tn.datadir = filepath.Join(ctx.workspace, fmt.Sprintf("node%d", nodeIndex))

	// P2P ports: 32000, 32001, 32002...
	// RPC ports: 38000, 38001, 38002...
	peers := make([]*discover.Node, ctx.numNodes)
	for i := 0; i < ctx.numNodes; i++ {
		id := crypto.FromECDSAPub(&ctx.accountKeys[i].PublicKey)[1:] // strip 0x04 prefix byte
		kni := fmt.Sprintf("kni://%x@127.0.0.1:%d?discport=0&type=cn", id, 32000+i)
		peers[i], err = discover.ParseNode(kni)
		if err != nil {
			return
		}
	}
	peers = append(peers[:nodeIndex], peers[nodeIndex+1:]...) // remove self

	nodeKey := ctx.accountKeys[nodeIndex]
	blsKey, _ := bls.DeriveFromECDSA(nodeKey)
	nodeConf := &node.Config{
		DataDir:           tn.datadir,
		UseLightweightKDF: true,
		P2P: p2p.Config{
			PrivateKey:             nodeKey,
			MaxPhysicalConnections: 100, // big enough
			ConnectionType:         common.CONSENSUSNODE,
			NoDiscovery:            true,
			StaticNodes:            peers,
			ListenAddr:             fmt.Sprintf("0.0.0.0:%d", 32000+nodeIndex),
		},
		BlsKey:           blsKey,
		IPCPath:          "klay.ipc",
		HTTPHost:         "0.0.0.0",
		HTTPPort:         38000 + nodeIndex,
		HTTPVirtualHosts: []string{"*"},
		HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
		NtpRemoteServer:  "",
	}
	if tn.node, err = node.New(nodeConf); err != nil {
		return
	}

	cnConf := cn.GetDefaultConfig()
	cnConf.NetworkId = ctx.config.ChainID.Uint64()
	cnConf.Genesis = ctx.genesis
	cnConf.Rewardbase = ctx.accountAddrs[nodeIndex]
	cnConf.SingleDB = false       // identical to regular CN
	cnConf.NumStateTrieShards = 4 // identical to regular CN
	cnConf.NoPruning = true       // archive mode
	err = tn.node.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return cn.New(ctx, cnConf)
	})
	if err != nil {
		return
	}
	if err = tn.node.Start(); err != nil {
		return
	}
	if err = tn.node.Service(&tn.cn); err != nil {
		return
	}
	ctx.nodes[nodeIndex] = tn
	return
}

func (ctx *blockchainTestContext) forEachNode(f func(*blockchainTestNode) error) error {
	for _, tn := range ctx.nodes {
		if err := f(tn); err != nil {
			return err
		}
	}
	return nil
}

func (ctx *blockchainTestContext) Start() error {
	return ctx.forEachNode(func(tn *blockchainTestNode) error {
		return tn.cn.StartMining(false)
	})
}

func (ctx *blockchainTestContext) Stop() error {
	err := ctx.forEachNode(func(tn *blockchainTestNode) error {
		return tn.node.Stop()
	})
	if err != nil {
		return err
	}

	blockchain.ClearMigrationPrerequisites()
	return nil
}

func (ctx *blockchainTestContext) Restart() error {
	if err := ctx.Stop(); err != nil {
		return err
	}
	// Recreate nodes
	if err := ctx.setNodes(ctx.numNodes); err != nil {
		return err
	}
	return ctx.Start()
}

func (ctx *blockchainTestContext) Cleanup() error {
	if err := ctx.Stop(); err != nil {
		return err
	}
	return os.RemoveAll(ctx.workspace)
}

func (ctx *blockchainTestContext) WaitBlock(t *testing.T, num uint64) {
	block := waitBlock(ctx.nodes[0].cn.BlockChain(), num)
	assert.NotNil(t, block)
}

func (ctx *blockchainTestContext) WaitTx(t *testing.T, txhash common.Hash) {
	rc := waitReceipt(ctx.nodes[0].cn.BlockChain().(*blockchain.BlockChain), txhash)
	assert.NotNil(t, rc)
	if rc != nil {
		assert.Equal(t, types.ReceiptStatusSuccessful, rc.Status)
	}
}

func (ctx *blockchainTestContext) Dump(t *testing.T) {
	for i, node := range ctx.nodes {
		t.Logf("node[%d] http://%s  %s", i, node.node.HTTPEndpoint(), node.node.IPCEndpoint())
	}
	for i, addr := range ctx.accountAddrs {
		t.Logf("account[%d] %s", i, addr.Hex())
	}
}

func (ctx *blockchainTestContext) Subscribe(t *testing.T, logFunc func(ev *blockchain.ChainEvent)) {
	if logFunc == nil {
		logFunc = func(ev *blockchain.ChainEvent) {
			t.Logf("block[%d] txs=%d", ev.Block.NumberU64(), ev.Block.Transactions().Len())
		}
	}

	go func() {
		chain := ctx.nodes[0].cn.BlockChain()
		chainEventCh := make(chan blockchain.ChainEvent)
		subscription := chain.SubscribeChainEvent(chainEventCh)
		defer subscription.Unsubscribe()
		for {
			ev := <-chainEventCh
			logFunc(&ev)
		}
	}()
}

func deriveTestAccount(index int) *ecdsa.PrivateKey {
	// "m/44'/60'/0'/0/0"
	mnemonic := "test test test test test test test test test test test junk"
	seed := pbkdf2.Key([]byte(mnemonic), []byte("mnemonic"), 2048, 64, sha512.New)
	key, _ := bip32.NewMasterKey(seed)
	key, _ = key.NewChildKey(0x8000002c)
	key, _ = key.NewChildKey(0x8000003c)
	key, _ = key.NewChildKey(0x80000000)
	key, _ = key.NewChildKey(0x00000000)

	child, _ := key.NewChildKey(uint32(index))
	privateKey, _ := crypto.ToECDSA(child.Key)
	return privateKey
}
