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
	compress_impl "github.com/kaiachain/kaia/kaiax/compress/impl"

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
		LivePruningRetention: blockchain.DefaultPruningRetention,
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

func TestCompressFunction(t *testing.T) {
	compress_impl.TestCompressFunction(t, setup)
}

func TestCompressPerformance(t *testing.T) {
	compress_impl.TestCompressPerformance(t, setup)
}
