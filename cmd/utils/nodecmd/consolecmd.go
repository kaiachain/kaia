// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
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
// This file is derived from cmd/geth/consolecmd.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package nodecmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kaiachain/kaia/v2/cmd/utils"
	"github.com/kaiachain/kaia/v2/console"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/networks/rpc"
	"github.com/kaiachain/kaia/v2/node"
	"github.com/urfave/cli/v2"
)

var AttachCommand = &cli.Command{
	Action:    remoteConsole,
	Name:      "attach",
	Usage:     "Start an interactive JavaScript environment (connect to node)",
	ArgsUsage: "[endpoint]",
	Flags:     append(utils.ConsoleFlags, utils.DataDirFlag),
	Category:  "CONSOLE COMMANDS",
	Description: `
The Klaytn console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console.
This command allows to open a console on a running Klaytn node.`,
}

// GetConsoleCommand returns cli.Command `console` whose flags are initialized with nodeFlags, rpcFlags, and ConsoleFlags.
func GetConsoleCommand(nodeFlags, rpcFlags []cli.Flag) *cli.Command {
	return &cli.Command{
		Action:   localConsole,
		Name:     "console",
		Usage:    "Start an interactive JavaScript environment",
		Flags:    append(append(nodeFlags, rpcFlags...), utils.ConsoleFlags...),
		Category: "CONSOLE COMMANDS",
		Description: `
The Klaytn console is an interactive shell for the JavaScript runtime environment
which exposes a node admin interface as well as the Ðapp JavaScript API.
See https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console.`,
	}
}

// localConsole starts a new Kaia node, attaching a JavaScript console to it at the
// same time.
func localConsole(ctx *cli.Context) error {
	// Create and start the node based on the CLI flags
	node := MakeFullNode(ctx)
	startNode(ctx, node)
	defer node.Stop()

	// Attach to the newly started node and start the JavaScript console
	client, err := node.Attach()
	if err != nil {
		log.Fatalf("Failed to attach to the inproc node: %v", err)
	}
	config := console.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.String(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		log.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	// If only a short execution was requested, evaluate and return
	if script := ctx.String(utils.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}
	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

// remoteConsole will connect to a remote node instance, attaching a JavaScript
// console to it.
func remoteConsole(ctx *cli.Context) error {
	// Attach to a remotely running node instance and start the JavaScript console
	endpoint := rpcEndpoint(ctx)
	client, err := dialRPC(endpoint)
	if err != nil {
		log.Fatalf("Unable to attach to remote node: %v", err)
	}
	config := console.Config{
		DataDir: utils.MakeDataDir(ctx),
		DocRoot: ctx.String(utils.JSpathFlag.Name),
		Client:  client,
		Preload: utils.MakeConsolePreloads(ctx),
	}

	console, err := console.New(config)
	if err != nil {
		log.Fatalf("Failed to start the JavaScript console: %v", err)
	}
	defer console.Stop(false)

	if script := ctx.String(utils.ExecFlag.Name); script != "" {
		console.Evaluate(script)
		return nil
	}

	// Otherwise print the welcome screen and enter interactive mode
	console.Welcome()
	console.Interactive()

	return nil
}

func rpcEndpoint(ctx *cli.Context) string {
	endpoint := ctx.Args().First()
	if endpoint == "" {
		path := node.DefaultDataDir()
		if ctx.IsSet(utils.DataDirFlag.Name) {
			path = ctx.String(utils.DataDirFlag.Name)
		}
		if path != "" {
			if ctx.Bool(utils.KairosFlag.Name) {
				path = filepath.Join(path, "baobab") // TODO: rename to Kairos
			}
		}
		endpoint = fmt.Sprintf("%s/klay.ipc", path)
	}
	return endpoint
}

// dialRPC returns a RPC client which connects to the given endpoint.
// The check for empty endpoint implements the defaulting logic
// for "ken attach" and "ken monitor" with no argument.
func dialRPC(endpoint string) (*rpc.Client, error) {
	if endpoint == "" {
		endpoint = node.DefaultIPCEndpoint(utils.ClientIdentifier)
	} else if strings.HasPrefix(endpoint, "rpc:") || strings.HasPrefix(endpoint, "ipc:") {
		// TODO-Kaia-RemoveLater: The below backward compatibility is not related to Kaia.
		// Backwards compatibility with Klaytn < 1.5 which required
		// these prefixes.
		endpoint = endpoint[4:]
	}
	return rpc.Dial(endpoint)
}
