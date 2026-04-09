package nodecmd

import (
	"flag"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/kaiachain/kaia/cmd/utils"
	"github.com/urfave/cli/v2"
)

func newNodecmdContext(t *testing.T, args ...string) *cli.Context {
	app := cli.NewApp()
	app.Flags = utils.AllNodeFlags()

	set := flag.NewFlagSet("nodecmd-test", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	for _, f := range app.Flags {
		if err := f.Apply(set); err != nil {
			t.Fatalf("failed to apply flag %v: %v", f.Names(), err)
		}
	}
	if err := set.Parse(args); err != nil {
		t.Fatalf("failed to parse flags: %v", err)
	}
	return cli.NewContext(app, set, nil)
}

func stopNodeOnCleanup(t *testing.T, stack interface{ Stop() error }) {
	t.Helper()
	t.Cleanup(func() {
		_ = stack.Stop()
	})
}

func registeredConstructorCount(t *testing.T, stack interface{}, fieldName string) int {
	t.Helper()

	value := reflect.ValueOf(stack)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		t.Fatalf("stack must be a non-nil pointer, got %T", stack)
	}

	field := value.Elem().FieldByName(fieldName)
	if !field.IsValid() {
		t.Fatalf("node field %q not found", fieldName)
	}
	if field.Kind() != reflect.Slice {
		t.Fatalf("node field %q must be a slice, got %s", fieldName, field.Kind())
	}
	return field.Len()
}

func TestDefaultCmd_MakeFullNode(t *testing.T) {
	datadir := tmpdir(t)
	t.Cleanup(func() { _ = os.RemoveAll(datadir) })

	ctx := newNodecmdContext(t, "--datadir", datadir, "--verbosity", "0")
	stack := MakeFullNode(ctx)
	if stack == nil {
		t.Fatal("MakeFullNode returned nil")
	}
	stopNodeOnCleanup(t, stack)
	if got := stack.DataDir(); got != datadir {
		t.Fatalf("unexpected datadir: got %q want %q", got, datadir)
	}
	if got := registeredConstructorCount(t, stack, "coreServiceFuncs"); got != 1 {
		t.Fatalf("unexpected number of core services: got %d want 1", got)
	}
	if got := registeredConstructorCount(t, stack, "serviceFuncs"); got != 0 {
		t.Fatalf("unexpected number of sub services: got %d want 0", got)
	}
}

func TestDefaultCmd_MakeFullNodeServiceRegistrationFlags(t *testing.T) {
	runner := newMakeFullNodeRunner(t)
	tests := []struct {
		name                 string
		args                 []string
		wantMainBridge       bool
		wantSubBridge        bool
		wantDBSyncer         bool
		wantChainDataFetcher bool
	}{
		{
			name:           "mainbridge",
			args:           []string{"--mainbridge"},
			wantMainBridge: true,
		},
		{
			name:          "subbridge",
			args:          []string{"--subbridge"},
			wantSubBridge: true,
		},
		{
			name:         "dbsyncer",
			args:         []string{"--dbsyncer", "--dbsyncer.db.host", "127.0.0.1", "--dbsyncer.db.user", "tester", "--dbsyncer.db.password", "secret", "--dbsyncer.db.name", "kaia"},
			wantDBSyncer: true,
		},
		{
			name:                 "chaindatafetcher",
			args:                 []string{"--chaindatafetcher", "--chaindatafetcher.mode", "kafka", "--chaindatafetcher.kafka.brokers", "127.0.0.1:9092"},
			wantChainDataFetcher: true,
		},
		{
			name:                 "combined",
			args:                 []string{"--mainbridge", "--subbridge", "--dbsyncer", "--dbsyncer.db.host", "127.0.0.1", "--dbsyncer.db.user", "tester", "--dbsyncer.db.password", "secret", "--dbsyncer.db.name", "kaia", "--chaindatafetcher", "--chaindatafetcher.mode", "kafka", "--chaindatafetcher.kafka.brokers", "127.0.0.1:9092"},
			wantMainBridge:       true,
			wantSubBridge:        true,
			wantDBSyncer:         true,
			wantChainDataFetcher: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			datadir := tmpdir(t)
			t.Cleanup(func() { _ = os.RemoveAll(datadir) })

			args := append([]string{"--datadir", datadir, "--verbosity", "0"}, tc.args...)
			ctx, err := runner.context(args...)
			if err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}
			_, cfg := utils.MakeConfigNode(ctx)
			if cfg.ServiceChain.EnabledMainBridge != tc.wantMainBridge {
				t.Fatalf("unexpected EnabledMainBridge: got %v want %v", cfg.ServiceChain.EnabledMainBridge, tc.wantMainBridge)
			}
			if cfg.ServiceChain.EnabledSubBridge != tc.wantSubBridge {
				t.Fatalf("unexpected EnabledSubBridge: got %v want %v", cfg.ServiceChain.EnabledSubBridge, tc.wantSubBridge)
			}
			if cfg.DB.EnabledDBSyncer != tc.wantDBSyncer {
				t.Fatalf("unexpected EnabledDBSyncer: got %v want %v", cfg.DB.EnabledDBSyncer, tc.wantDBSyncer)
			}
			if cfg.ChainDataFetcher.EnabledChainDataFetcher != tc.wantChainDataFetcher {
				t.Fatalf("unexpected EnabledChainDataFetcher: got %v want %v", cfg.ChainDataFetcher.EnabledChainDataFetcher, tc.wantChainDataFetcher)
			}

			stack := MakeFullNode(ctx)
			if stack == nil {
				t.Fatal("MakeFullNode returned nil")
			}
			stopNodeOnCleanup(t, stack)

			wantSubServices := 0
			if tc.wantMainBridge {
				wantSubServices++
			}
			if tc.wantSubBridge {
				wantSubServices++
			}
			if tc.wantDBSyncer {
				wantSubServices++
			}
			if tc.wantChainDataFetcher {
				wantSubServices++
			}
			if got := registeredConstructorCount(t, stack, "coreServiceFuncs"); got != 1 {
				t.Fatalf("unexpected number of core services: got %d want 1", got)
			}
			if got := registeredConstructorCount(t, stack, "serviceFuncs"); got != wantSubServices {
				t.Fatalf("unexpected number of sub services: got %d want %d", got, wantSubServices)
			}
		})
	}
}

func TestFlagSetReuse_LeaksIsSet(t *testing.T) {
	app := cli.NewApp()
	app.Flags = utils.AllNodeFlags()

	set := flag.NewFlagSet("nodecmd-reuse-leak", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	for _, f := range app.Flags {
		if err := f.Apply(set); err != nil {
			t.Fatalf("failed to apply flag %v: %v", f.Names(), err)
		}
	}
	reset := func() {
		set.VisitAll(func(f *flag.Flag) {
			_ = f.Value.Set(f.DefValue)
		})
	}

	reset()
	if err := set.Parse([]string{"--mainbridge"}); err != nil {
		t.Fatalf("failed to parse first args: %v", err)
	}
	ctx1 := cli.NewContext(app, set, nil)
	if !ctx1.IsSet(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q to be set in first parse", utils.MainBridgeFlag.Name)
	}
	if !ctx1.Bool(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q bool value to be true in first parse", utils.MainBridgeFlag.Name)
	}

	reset()
	if err := set.Parse([]string{"--subbridge"}); err != nil {
		t.Fatalf("failed to parse second args: %v", err)
	}
	ctx2 := cli.NewContext(app, set, nil)
	if ctx2.Bool(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q bool value to be false after reset", utils.MainBridgeFlag.Name)
	}
	if !ctx2.IsSet(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q IsSet to leak when reusing FlagSet", utils.MainBridgeFlag.Name)
	}
}

func TestMakeFullNodeRunner_Context_NoIsSetLeak(t *testing.T) {
	runner := newMakeFullNodeRunner(t)

	datadir1 := tmpdir(t)
	t.Cleanup(func() { _ = os.RemoveAll(datadir1) })
	ctx1, err := runner.context("--datadir", datadir1, "--verbosity", "0", "--mainbridge")
	if err != nil {
		t.Fatalf("failed to parse first context: %v", err)
	}
	if !ctx1.IsSet(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q to be set in first parse", utils.MainBridgeFlag.Name)
	}

	datadir2 := tmpdir(t)
	t.Cleanup(func() { _ = os.RemoveAll(datadir2) })
	ctx2, err := runner.context("--datadir", datadir2, "--verbosity", "0", "--subbridge")
	if err != nil {
		t.Fatalf("failed to parse second context: %v", err)
	}
	if ctx2.IsSet(utils.MainBridgeFlag.Name) {
		t.Fatalf("expected %q IsSet to be false in second parse", utils.MainBridgeFlag.Name)
	}
	if !ctx2.IsSet(utils.SubBridgeFlag.Name) {
		t.Fatalf("expected %q to be set in second parse", utils.SubBridgeFlag.Name)
	}
}

const (
	FlagTypeBoolean = iota
	FlagTypeArgument
)

const (
	ErrorIncorrectUsage = iota
	ErrorInvalidValue
	ErrorFatal
	// TODO-Kaia-Node fix the configuration to filter wrong input flags before the Kaia server is launched
	NonError // This error case expects an error, but currently it does not filter the wrong value.
)

var (
	commonThreeErrors = []string{"abcdefg", "1234567", "!@#$%^&"}
	commonTwoErrors   = []string{"abcdefg", "!@#$%^&"}
)

var flagsWithValues = []struct {
	flag        string
	flagType    uint16
	values      []string
	wrongValues []string
	errors      []int
}{
	// TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:     "--networktype",
	//	flagType: FlagTypeArgument,
	//	// values: []string{"mn", "scn"},
	//	values: []string{},
	//	wrongValues: []string{"kairos", "abcdefg", "1234567", "!@#$%^&"},
	//	errors: []string{},
	//},
	{
		flag:        "--dbtype",
		flagType:    FlagTypeArgument,
		values:      []string{"LevelDB", "BadgerDB", "MemoryDB", "DynamoDBS3"},
		wrongValues: append(commonThreeErrors, "oracle"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--srvtype",
		flagType:    FlagTypeArgument,
		values:      []string{"http", "fasthttp"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--keystore",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--networkid",
		flagType:    FlagTypeArgument,
		values:      []string{"1", "1000", "1001", "12312"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--identity",
		flagType:    FlagTypeArgument,
		values:      []string{"abc", "abde-", "oai121"},
		wrongValues: []string{},
		errors:      []int{},
	},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--docroot",
	//	flagType:    FlagTypeBoolean,
	//	values:      []string{},
	//},
	{
		flag:        "--syncmode",
		flagType:    FlagTypeArgument,
		values:      []string{"full"}, //[]string{"fast", "full"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--gcmode",
		flagType:    FlagTypeArgument,
		values:      []string{"full", "archive"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:     "--lightkdf",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--txpool.nolocals",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--txpool.journal",
		flagType:    FlagTypeArgument,
		values:      []string{"transactions.rlp"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--txpool.journal-interval",
		flagType:    FlagTypeArgument,
		values:      []string{"1h0m0s", "0h0m0s", "0h0m1s", "0h1m0s", "0.5h0.5m0.5s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.pricelimit",
		flagType:    FlagTypeArgument,
		values:      []string{"1"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.pricebump",
		flagType:    FlagTypeArgument,
		values:      []string{"10"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.exec-slots.account",
		flagType:    FlagTypeArgument,
		values:      []string{"16"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.exec-slots.all",
		flagType:    FlagTypeArgument,
		values:      []string{"4096"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.nonexec-slots.account",
		flagType:    FlagTypeArgument,
		values:      []string{"64"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txpool.nonexec-slots.all",
		flagType:    FlagTypeArgument,
		values:      []string{"1024"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--txpool.keeplocals",
	//	flagType:    FlagTypeBoolean,
	//	values:      []string{},
	//},
	{
		flag:        "--txpool.lifetime",
		flagType:    FlagTypeArgument,
		values:      []string{"0h20m0s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--db.single",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--db.num-statetrie-shards",
		flagType: FlagTypeArgument,
		// values:    []string{"1", "2"},
		values:      []string{"1"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--db.leveldb.cache-size",
		flagType:    FlagTypeArgument,
		values:      []string{"768"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--db.leveldb.compression",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2", "3"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--db.leveldb.no-buffer-pool",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--db.pebbledb.cache-size",
		flagType:    FlagTypeArgument,
		values:      []string{"768"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--db.no-parallel-write",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--state.cache-size",
		flagType:    FlagTypeArgument,
		values:      []string{"64", "128", "256", "512"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--state.block-interval",
		flagType:    FlagTypeArgument,
		values:      []string{"64", "128", "256"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.type",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.scale",
		flagType:    FlagTypeArgument,
		values:      []string{"100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.level",
		flagType:    FlagTypeBoolean,
		values:      []string{"saving", "normal", "extreme"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--cache.memory",
		flagType:    FlagTypeArgument,
		values:      []string{"16"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--state.trie-cache-limit",
		flagType:    FlagTypeArgument,
		values:      []string{"512", "1024", "2048", "4096"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--sendertxhashindexing",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--childchainindexing",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--targetgaslimit",
		flagType:    FlagTypeArgument,
		values:      []string{"4712388"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--scsigner",
		flagType:    FlagTypeArgument,
		values:      []string{"0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--rewardbase",
		flagType:    FlagTypeArgument,
		values:      []string{"0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--extradata",
		flagType:    FlagTypeArgument,
		values:      []string{"0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--txresend.interval",
		flagType:    FlagTypeArgument,
		values:      []string{"3", "5", "7"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--txresend.max-count",
		flagType:    FlagTypeArgument,
		values:      []string{"1000", "2000"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--txresend.use-legacy",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--unlock",
		flagType:    FlagTypeArgument,
		values:      []string{"", "0x0", "0x777fd033b5e3bcaad6006bc9f481ffed6b83cf5a"},
		wrongValues: []string{"abcdefg", "!@#$%^&", "0x921jfinowaae333"},
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--password",
		flagType:    FlagTypeArgument,
		values:      []string{"abcd", "aoije091"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--vmdebug",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--vmlog",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "1", "2", "3"},
		wrongValues: append(commonThreeErrors, "4"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:     "--metrics",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--prometheus",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--prometheusport",
		flagType:    FlagTypeArgument,
		values:      []string{"61001"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--rpc",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--rpcaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123"},
		wrongValues: append(commonThreeErrors, "123.123.123.256"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--rpcport",
		flagType:    FlagTypeArgument,
		values:      []string{"8551"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--rpccorsdomain",
		flagType:    FlagTypeArgument,
		values:      []string{"", "localhost", "123.123.123.123"},
		wrongValues: append(commonThreeErrors, "123.123.123.256"),
		errors:      []int{NonError, NonError, NonError, NonError},
	},
	{
		flag:        "--rpcvhosts",
		flagType:    FlagTypeArgument,
		values:      []string{"*"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--rpcapi",
		flagType:    FlagTypeArgument,
		values:      []string{"", "kaia", "kaia,personal,istanbul,debug,miner"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--ipcdisable",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--ipcpath",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--ws",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--wsaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--wsport",
		flagType:    FlagTypeArgument,
		values:      []string{"8552"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--grpc",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--grpcaddr",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--grpcport",
		flagType:    FlagTypeArgument,
		values:      []string{"8553"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--wsapi",
		flagType:    FlagTypeArgument,
		values:      []string{"", "kaia", "kaia,personal,istanbul,debug,miner"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--wsorigins",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--exec",
		flagType:    FlagTypeArgument,
		values:      []string{"klay.blockNumber", "klat.getBlock(0)", "governance.getParams()[\"reward.proposerupdateinterval\"]"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--preload",
		flagType:    FlagTypeArgument,
		values:      []string{"abc.js", "tmp.js", "tmp"},
		wrongValues: []string{},
		errors:      []int{},
	},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--nodetype",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{"cn", "pn", "en"},
	//  wrongValues: []string{},
	//  errors:      []int{},
	//},
	{
		flag:        "--maxconnections",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "30", "25000"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--maxpendpeers",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "30", "50"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--port",
		flagType:    FlagTypeArgument,
		values:      []string{"32323", "30303"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--subport",
		flagType:    FlagTypeArgument,
		values:      []string{"32324", "32325", "32327"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--multichannel",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--bootnodes",
		flagType:    FlagTypeArgument,
		values:      []string{"0xf4316f69d9522667c0674afcd8638288489f0333", "", "0xf4316f69d9522667c0674afcd8638288489f0333, d473284239f704adccd24647c7ca132992a28973"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--nodekey",
		flagType:    FlagTypeArgument,
		values:      []string{""},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--nodekeyhex",
		flagType:    FlagTypeArgument,
		values:      []string{"8da4ef21b864d2cc526dbdb2a120bd2874c36c9d0a1fb7f8c63d7f7a8b41de8f"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorFatal, ErrorFatal, ErrorFatal},
	},
	{
		flag:        "--nat",
		flagType:    FlagTypeArgument,
		values:      []string{"any", "none", "upnp", "pmp", "extip:127.0.0.1"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--nodiscover",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--netrestrict",
		flagType:    FlagTypeArgument,
		values:      []string{},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:        "--chaintxperiod",
		flagType:    FlagTypeArgument,
		values:      []string{"1", "5", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--chaintxlimit",
		flagType:    FlagTypeArgument,
		values:      []string{"100", "200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--jspath",
		flagType:    FlagTypeArgument,
		values:      []string{".", "root/abc/efg"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--kairos",
		flagType: FlagTypeBoolean,
	},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--bnaddr",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{},
	//	wrongValues: []string{},
	//	errors:      []int{},
	//},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--genkey",
	//	flagType:    FlagTypeArgument,
	//	values:      []string{},
	//	wrongValues: []string{},
	//	errors:      []int{},
	//},
	//TODO-Kaia-Node the flag is not defined on any Kaia binaries
	//{
	//	flag:        "--writeaddress",
	//	flagType:    FlagTypeBoolean,
	//},
	{
		flag:     "--mainbridge",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--subbridge",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--mainbridgeport",
		flagType:    FlagTypeArgument,
		values:      []string{"50505", "23232"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:        "--subbridgeport",
		flagType:    FlagTypeArgument,
		values:      []string{"50505", "23232"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, NonError, ErrorInvalidValue},
	},
	{
		flag:     "--vtrecovery",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--vtrecoveryinterval",
		flagType:    FlagTypeArgument,
		values:      []string{"60", "100", "200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:     "--scnewaccount",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--dbsyncer",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.host",
		flagType:    FlagTypeArgument,
		values:      []string{"localhost", "123.123.123.123", "127.0.0.1"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.db.port",
		flagType:    FlagTypeArgument,
		values:      []string{"3306", "32323"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--dbsyncer.db.name",
		flagType: FlagTypeBoolean,
	},
	{
		flag:     "--dbsyncer.db.user",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.password",
		flagType:    FlagTypeArgument,
		values:      []string{"aboaise", "jaooao122!@", "18231@#!@412!"},
		wrongValues: []string{},
		errors:      []int{},
	},
	{
		flag:     "--dbsyncer.logmode",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--dbsyncer.db.max.idle",
		flagType:    FlagTypeArgument,
		values:      []string{"50", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.db.max.open",
		flagType:    FlagTypeArgument,
		values:      []string{"30", "50", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.db.max.lifetime",
		flagType:    FlagTypeArgument,
		values:      []string{"1h0m0s", "0h0m10s"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.block.channel.size",
		flagType:    FlagTypeArgument,
		values:      []string{"5"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.mode",
		flagType:    FlagTypeArgument,
		values:      []string{"single", "multi", "context"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.genquery.th",
		flagType:    FlagTypeArgument,
		values:      []string{"50", "123"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.insert.th",
		flagType:    FlagTypeArgument,
		values:      []string{"30", "123"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.bulk.size",
		flagType:    FlagTypeArgument,
		values:      []string{"200"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--dbsyncer.event.mode",
		flagType:    FlagTypeArgument,
		values:      []string{"block", "head"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:        "--dbsyncer.max.block.diff",
		flagType:    FlagTypeArgument,
		values:      []string{"0", "5", "100"},
		wrongValues: commonTwoErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue},
	},
	{
		flag:        "--autorestart.daemon.path",
		flagType:    FlagTypeArgument,
		values:      []string{"~/klaytn/bin/kcnd", "~/klaytn/bin/kpnd", "~/klaytn/bin/kend"},
		wrongValues: commonThreeErrors,
		errors:      []int{NonError, NonError, NonError},
	},
	{
		flag:     "--autorestart.enable",
		flagType: FlagTypeBoolean,
	},
	{
		flag:        "--autorestart.timeout",
		flagType:    FlagTypeArgument,
		values:      []string{"10m", "60s", "1h"},
		wrongValues: commonThreeErrors,
		errors:      []int{ErrorInvalidValue, ErrorInvalidValue, ErrorInvalidValue},
	},
}

type makeFullNodeRunner struct {
	app *cli.App
}

func newMakeFullNodeRunner(t *testing.T) *makeFullNodeRunner {
	t.Helper()
	app := cli.NewApp()
	app.Flags = utils.AllNodeFlags()
	return &makeFullNodeRunner{app: app}
}

func (r *makeFullNodeRunner) context(args ...string) (*cli.Context, error) {
	set := flag.NewFlagSet("nodecmd-makefullnode-runner", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	for _, f := range r.app.Flags {
		if err := f.Apply(set); err != nil {
			return nil, err
		}
	}
	if err := set.Parse(args); err != nil {
		return nil, err
	}
	return cli.NewContext(r.app, set, nil), nil
}

type parseOnlyRunner struct {
	set *flag.FlagSet
}

func newParseOnlyRunner(t *testing.T) *parseOnlyRunner {
	t.Helper()
	app := cli.NewApp()
	app.Flags = utils.AllNodeFlags()
	set := flag.NewFlagSet("nodecmd-parseonly-inproc", flag.ContinueOnError)
	set.SetOutput(io.Discard)
	for _, f := range app.Flags {
		if err := f.Apply(set); err != nil {
			t.Fatalf("failed to apply flag %v: %v", f.Names(), err)
		}
	}
	return &parseOnlyRunner{set: set}
}

func (r *parseOnlyRunner) parse(args ...string) error {
	// Reuse the same parser for speed; reset values before each parse.
	r.set.VisitAll(func(f *flag.Flag) {
		_ = f.Value.Set(f.DefValue)
	})
	baseArgs := []string{"--verbosity", "0"}
	return r.set.Parse(append(baseArgs, args...))
}

func sanitizeTestName(s string) string {
	replacer := strings.NewReplacer("/", "_", " ", "_", "\t", "_")
	return replacer.Replace(s)
}

func parseExpectationFromLegacyResult(resultCode int) (strict bool, wantErr bool) {
	switch resultCode {
	case ErrorIncorrectUsage, ErrorInvalidValue:
		return true, true
	case ErrorFatal:
		// Fatal cases are expected to fail in config/setup stage, not parsing stage.
		return true, false
	case NonError:
		// Legacy matrix marks these as currently unfiltered.
		return false, false
	default:
		return false, false
	}
}

func TestDefaultCmd_ParseOnlyAllFlagCases(t *testing.T) {
	runner := newParseOnlyRunner(t)
	for _, fwv := range flagsWithValues {
		switch fwv.flagType {
		case FlagTypeBoolean:
			t.Run("parseonly"+sanitizeTestName(fwv.flag), func(t *testing.T) {
				if err := runner.parse(fwv.flag, ""); err != nil {
					t.Fatalf("expected parse success for %s, got error: %v", fwv.flag, err)
				}
			})
		case FlagTypeArgument:
			for _, item := range fwv.values {
				t.Run("parseonly"+sanitizeTestName(fwv.flag)+"-"+sanitizeTestName(item), func(t *testing.T) {
					if err := runner.parse(fwv.flag, item); err != nil {
						t.Fatalf("expected parse success for %s %q, got error: %v", fwv.flag, item, err)
					}
				})
			}
			for idx2, wrongItem := range fwv.wrongValues {
				t.Run("parseonly"+sanitizeTestName(fwv.flag)+"-"+sanitizeTestName(wrongItem), func(t *testing.T) {
					strict, wantErr := parseExpectationFromLegacyResult(fwv.errors[idx2])
					err := runner.parse(fwv.flag, wrongItem)
					if !strict {
						return
					}
					if wantErr && err == nil {
						t.Fatalf("expected parse failure for %s %q, but got success", fwv.flag, wrongItem)
					}
					if !wantErr && err != nil {
						t.Fatalf("expected parse success for %s %q, got error: %v", fwv.flag, wrongItem, err)
					}
				})
			}
		}
	}
}
