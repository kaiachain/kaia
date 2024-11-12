// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2019 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/geth/main.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package nodecmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/kaiachain/kaia/accounts"
	"github.com/kaiachain/kaia/accounts/keystore"
	"github.com/kaiachain/kaia/api/debug"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/cmd/utils"
	"github.com/kaiachain/kaia/log"
	metricutils "github.com/kaiachain/kaia/metrics/utils"
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/node/cn"
	"github.com/kaiachain/kaia/params"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

// runKaiaNode is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func RunKaiaNode(ctx *cli.Context) error {
	fullNode := MakeFullNode(ctx)
	startNode(ctx, fullNode)
	fullNode.Wait()
	return nil
}

func MakeFullNode(ctx *cli.Context) *node.Node {
	stack, cfg := utils.MakeConfigNode(ctx)

	if utils.NetworkTypeFlag.Value == utils.SCNNetworkType && cfg.ServiceChain.EnabledSubBridge {
		if !cfg.CN.NoAccountCreation {
			logger.Warn("generated accounts can't be synced with the parent chain since account creation is enabled")
		}
		switch cfg.ServiceChain.ServiceChainConsensus {
		case "istanbul":
			utils.RegisterCNService(stack, &cfg.CN)
		case "clique":
			logger.Crit("using clique consensus type is not allowed anymore!")
		default:
			logger.Crit("unknown consensus type for the service chain", "consensus", cfg.ServiceChain.ServiceChainConsensus)
		}
	} else {
		utils.RegisterCNService(stack, &cfg.CN)
	}
	utils.RegisterService(stack, &cfg.ServiceChain)
	utils.RegisterDBSyncerService(stack, &cfg.DB)
	utils.RegisterChainDataFetcherService(stack, &cfg.ChainDataFetcher)
	return stack
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, stack *node.Node) {
	// Ntp time check
	if err := node.NtpCheckWithLocal(stack); err != nil {
		log.Fatalf("System time should be synchronized: %v", err)
	}

	// Start up the node itself
	utils.StartNode(stack)

	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {
		// Create a chain state reader for self-derivation
		rpcClient, err := stack.Attach()
		if err != nil {
			log.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := client.NewClient(rpcClient)

		// Open any wallets already attached
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				logger.Error("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					logger.Error("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				logger.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				logger.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	if utils.NetworkTypeFlag.Value == utils.SCNNetworkType && utils.ServiceChainConsensusFlag.Value == "clique" {
		logger.Crit("using clique consensus type is not allowed anymore!")
	} else {
		startKaiaAuxiliaryService(ctx, stack)
	}

	// Unlock any account specifically requested
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.String(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			UnlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
}

func startKaiaAuxiliaryService(ctx *cli.Context, stack *node.Node) {
	var cn *cn.CN
	if err := stack.Service(&cn); err != nil {
		log.Fatalf("Kaia service not running: %v", err)
	}

	// TODO-Kaia-NodeCmd disable accept tx before finishing sync.
	if err := cn.StartMining(false); err != nil {
		log.Fatalf("Failed to start mining: %v", err)
	}
}

func CommandNotExist(ctx *cli.Context, s string) {
	cli.ShowAppHelp(ctx)
	fmt.Printf("Error: Unknown command \"%v\"\n", s)
	os.Exit(1)
}

func OnUsageError(context *cli.Context, err error, isSubcommand bool) error {
	cli.ShowAppHelp(context)
	return err
}

func CheckCommands(ctx *cli.Context) error {
	valid := false
	for _, cmd := range ctx.App.Commands {
		if cmd.Name == ctx.Args().First() {
			valid = true
		}
	}

	if !valid && ctx.Args().Present() {
		cli.ShowAppHelp(ctx)
		return fmt.Errorf("Unknown command \"%v\"\n", ctx.Args().First())
	}

	return nil
}

func FlagsFromYaml(ctx *cli.Context) error {
	confFile := "conf" // flag option for yaml file name
	if ctx.String(confFile) != "" {
		if err := altsrc.InitInputSourceWithContext(utils.AllNodeFlags(), altsrc.NewYamlSourceFromFlagFunc(confFile))(ctx); err != nil {
			return err
		}
	}
	return nil
}

func BeforeRunNode(ctx *cli.Context) error {
	if err := FlagsFromYaml(ctx); err != nil {
		return err
	}
	MigrateGlobalFlags(ctx)
	if err := CheckCommands(ctx); err != nil {
		return err
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	logDir := (&node.Config{DataDir: utils.MakeDataDir(ctx)}).ResolvePath("logs")
	debug.CreateLogDir(logDir)
	if err := debug.Setup(ctx); err != nil {
		return err
	}
	metricutils.StartMetricCollectionAndExport(ctx)
	setupNetwork(ctx)
	return nil
}

// setupNetwork configures the system for either the main net or some test network.
func setupNetwork(ctx *cli.Context) {
	params.TargetGasLimit = ctx.Uint64(utils.TargetGasLimitFlag.Name)
}

func BeforeRunBootnode(ctx *cli.Context) error {
	if err := FlagsFromYaml(ctx); err != nil {
		return err
	}
	MigrateGlobalFlags(ctx)
	if err := debug.Setup(ctx); err != nil {
		return err
	}
	metricutils.StartMetricCollectionAndExport(ctx)
	return nil
}

var migrationApplied = map[*cli.Command]struct{}{}

// migrateGlobalFlags makes all global flag values available in the
// context. This should be called as early as possible in app.Before.
//
// Example:
//
//	ken account new --keystore /tmp/mykeystore --lightkdf
//
// is equivalent after calling this method with:
//
//	ken --keystore /tmp/mykeystore --lightkdf account new
//
// i.e. in the subcommand Action function of 'account new', ctx.Bool("lightkdf)
// will return true even if --lightkdf is set as a global option.
//
// This function may become unnecessary when https://github.com/urfave/cli/pull/1245 is merged.
func MigrateGlobalFlags(ctx *cli.Context) {
	var iterate func(cs []*cli.Command, fn func(*cli.Command))
	iterate = func(cs []*cli.Command, fn func(*cli.Command)) {
		for _, cmd := range cs {
			if _, ok := migrationApplied[cmd]; ok {
				continue
			}
			migrationApplied[cmd] = struct{}{}
			fn(cmd)
			iterate(cmd.Subcommands, fn)
		}
	}

	// This iterates over all commands and wraps their action function.
	iterate(ctx.App.Commands, func(cmd *cli.Command) {
		if cmd.Action == nil {
			return
		}

		action := cmd.Action
		cmd.Action = func(ctx *cli.Context) error {
			doMigrateFlags(ctx)
			return action(ctx)
		}
	})
}

func doMigrateFlags(ctx *cli.Context) {
	for _, name := range ctx.FlagNames() {
		for _, parent := range ctx.Lineage()[1:] {
			if parent.IsSet(name) {
				ctx.Set(name, parent.String(name))
				break
			}
		}
	}
}
