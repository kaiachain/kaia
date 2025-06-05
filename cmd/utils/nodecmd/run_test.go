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
// This file is derived from cmd/geth/run_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package nodecmd

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sort"
	"testing"
	"time"

	"github.com/docker/docker/pkg/reexec"
	"github.com/kaiachain/kaia/v2/api/debug"
	"github.com/kaiachain/kaia/v2/cmd/utils"
	"github.com/kaiachain/kaia/v2/console"
	metricutils "github.com/kaiachain/kaia/v2/metrics/utils"
	"github.com/kaiachain/kaia/v2/networks/rpc"
	"github.com/kaiachain/kaia/v2/node"
	"github.com/urfave/cli/v2"
)

func tmpdir(t *testing.T) string {
	dir, err := os.MkdirTemp("", "kaia-test-")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

type testKaia struct {
	*utils.TestCmd

	// template variables for expect
	Datadir    string
	Rewardbase string
}

var (
	// The app that holds all commands and flags.
	app = utils.NewApp(GetGitCommit(), "the Kaia command line interface")

	// flags that configure the node
	nodeFlags = utils.CommonNodeFlags

	rpcFlags = utils.CommonRPCFlags
)

func init() {
	// Initialize the CLI app and start Kaia
	app.Action = RunKaiaNode
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2018-2024 The Kaia Authors"
	app.Commands = []*cli.Command{
		// See chaincmd.go:
		InitCommand,

		// See accountcmd.go
		AccountCommand,

		// See consolecmd.go:
		GetConsoleCommand(nodeFlags, rpcFlags),
		AttachCommand,

		// See versioncmd.go:
		VersionCommand,

		// See dumpconfigcmd.go:
		GetDumpConfigCommand(nodeFlags, rpcFlags),
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = utils.AllNodeFlags()

	app.Before = func(ctx *cli.Context) error {
		MigrateGlobalFlags(ctx)
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

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}

	// Run the app if we've been exec'd as "kaia-test" in runKaia.
	reexec.Register("kaia-test", func() {
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
	reexec.Register("kaia-test-flag", func() {
		app.Action = RunTestKaiaNode
		if err := app.Run(os.Args); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	})
}

func TestMain(m *testing.M) {
	// check if we have been reexec'd
	if reexec.Init() {
		return
	}
	os.Exit(m.Run())
}

// spawns Kaia with the given command line args. If the args don't set --datadir, the
// child g gets a temporary data directory.
func runKaia(t *testing.T, name string, args ...string) *testKaia {
	tt := &testKaia{}
	tt.TestCmd = utils.NewTestCmd(t, tt)
	for i, arg := range args {
		switch {
		case arg == "-datadir" || arg == "--datadir":
			if i < len(args)-1 {
				tt.Datadir = args[i+1]
			}
		case arg == "-rewardbase" || arg == "--rewardbase":
			if i < len(args)-1 {
				tt.Rewardbase = args[i+1]
			}
		}
	}
	if tt.Datadir == "" {
		tt.Datadir = tmpdir(t)
		tt.Cleanup = func() { os.RemoveAll(tt.Datadir) }
		args = append([]string{"--datadir", tt.Datadir}, args...)
		// Remove the temporary datadir if something fails below.
		defer func() {
			if t.Failed() {
				tt.Cleanup()
			}
		}()
	}

	// Boot "Kaia". This actually runs the test binary but the TestMain
	// function will prevent any tests from running.
	tt.Run(name, args...)

	return tt
}

// waitForEndpoint attempts to connect to an RPC endpoint until it succeeds.
func waitForEndpoint(t *testing.T, endpoint string, timeout time.Duration) {
	probe := func() bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		c, err := rpc.DialContext(ctx, endpoint)
		if c != nil {
			_, err = c.SupportedModules()
			c.Close()
		}
		return err == nil
	}

	start := time.Now()
	for {
		if probe() {
			return
		}
		if time.Since(start) > timeout {
			t.Fatal("endpoint", endpoint, "did not open within", timeout)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func RunTestKaiaNode(ctx *cli.Context) error {
	fullNode := MakeFullNode(ctx)
	fullNode.Wait()
	return nil
}
