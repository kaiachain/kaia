package cn

import (
	"errors"
	"os"
	"reflect"
	"strconv"
	"testing"

	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/event"
	"github.com/kaiachain/kaia/governance"
	"github.com/kaiachain/kaia/kaiax/compress"
	compress_bc "github.com/kaiachain/kaia/kaiax/compress/body"
	compress_hc "github.com/kaiachain/kaia/kaiax/compress/header"
	compress_interface "github.com/kaiachain/kaia/kaiax/compress/interface"
	compress_rc "github.com/kaiachain/kaia/kaiax/compress/receipts"
	"github.com/stretchr/testify/assert"

	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/kaiachain/kaia/storage/statedb"
)

var (
	ERR_ENV_NOSET_DIR     = errors.New("directory name(TEST_COMPRESS_DIR) is empty")
	ERR_ENV_NOSET_NETWORK = errors.New("network ID(TEST_COMPRESS_NETWORKID) is empty")
)

func setup(t *testing.T) (*blockchain.BlockChain, database.DBManager, error) {
	dir := os.Getenv("TEST_COMPRESS_DIR")
	if dir == "" {
		return nil, nil, ERR_ENV_NOSET_DIR
	}
	networkIDStr := os.Getenv("TEST_COMPRESS_NETWORKID")
	if networkIDStr == "" {
		return nil, nil, ERR_ENV_NOSET_NETWORK
	}
	networkID, err := strconv.ParseUint(networkIDStr, 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	ctx := node.NewServiceContext(&node.Config{Name: "klay", DataDir: dir}, map[reflect.Type]node.Service{}, &event.TypeMux{}, &accounts.Manager{})
	dbc := &database.DBConfig{
		Dir:                 "chaindata",
		DBType:              database.LevelDB,
		SingleDB:            false,
		NumStateTrieShards:  4,
		ParallelDBWrite:     false,
		OpenFilesLimit:      4,
		EnableDBPerfMetrics: false,
		LevelDBCacheSize:    768,
		LevelDBCompression:  0,
		LevelDBBufferPool:   true,
		PebbleDBCacheSize:   768,
	}
	chainDB := ctx.OpenDatabase(dbc)
	chainConfig, _, genesisErr := blockchain.SetupGenesisBlock(chainDB, nil, networkID, false, false)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		t.Fatal(genesisErr)
	}
	cacheConfig := &blockchain.CacheConfig{
		ArchiveMode:          false,
		CacheSize:            512,
		BlockInterval:        blockchain.DefaultBlockInterval,
		TriesInMemory:        blockchain.DefaultTriesInMemory,
		LivePruningRetention: blockchain.DefaultLivePruningRetention,
		TrieNodeCacheConfig:  statedb.GetEmptyTrieNodeCacheConfig(),
		SnapshotCacheSize:    512,
		SnapshotAsyncGen:     true,
	}
	config := GetDefaultConfig()
	governance := governance.NewMixedEngine(chainConfig, chainDB)
	engine := CreateConsensusEngine(ctx, config, chainConfig, chainDB, governance, ctx.NodeType())
	bc, err := blockchain.NewBlockChain(chainDB, cacheConfig, chainConfig, engine, vm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return bc, chainDB, nil
}

func TestCompress(t *testing.T) {
	var (
		copyTempDirHC         = "copy_header"
		copyTempDirBC         = "copy_body"
		copyTempDirRC         = "copy_receipts"
		bc, chainDB, setupErr = setup(t)
		initOpts              = &compress.InitOpts{
			Chain: bc,
			Dbm:   chainDB,
		}
		mCompressHC = compress_hc.NewHeaderCompression()
		mCompressBC = compress_bc.NewBodyCompression()
		mCompressRC = compress_rc.NewReceiptCompression()
		hcInitErr   = mCompressHC.Init(initOpts)
		bcInitErr   = mCompressBC.Init(initOpts)
		rcInitErr   = mCompressRC.Init(initOpts)
	)

	if err := errors.Join(setupErr, hcInitErr, bcInitErr, rcInitErr); err != nil {
		// If no environment varaible set, do not execute compression test
		// TODO-hyunsooda: Change this test to functional test and remove temp storage directory
		if errors.Is(err, ERR_ENV_NOSET_DIR) || errors.Is(err, ERR_ENV_NOSET_NETWORK) {
			return
		}
		t.Fatal(err)
	}
	defer chainDB.Close()

	from, to := uint64(0), uint64(3300)
	// receipt compression test
	assert.Nil(t, compress_interface.TestCompress(mCompressHC, database.HeaderCompressType, chainDB.CompressHeader, from, to, &from, copyTempDirHC))
	assert.Nil(t, compress_interface.TestCompress(mCompressBC, database.BodyCompressType, chainDB.CompressBody, from, to, &from, copyTempDirBC))
	assert.Nil(t, compress_interface.TestCompress(mCompressRC, database.ReceiptCompressType, chainDB.CompressReceipts, from, to, &from, copyTempDirRC))
}
