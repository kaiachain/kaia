// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from cmd/bootnode/main.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiachain/kaia/api/debug"
	"github.com/kaiachain/kaia/cmd/utils"
	"github.com/kaiachain/kaia/cmd/utils/nodecmd"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/log"
	"github.com/kaiachain/kaia/networks/p2p"
	"github.com/kaiachain/kaia/networks/p2p/discover"
	"github.com/kaiachain/kaia/networks/p2p/nat"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/urfave/cli/v2"
)

var logger = log.NewModuleLogger(log.CMDKBN)

const (
	generateNodeKeySpecified = iota
	noPrivateKeyPathSpecified
	nodeKeyDuplicated
	writeOutAddress
	goodToGo
)

func bootnode(ctx *cli.Context) error {
	var (
		// Local variables
		err  error
		bcfg = bootnodeConfig{
			// Config variables
			networkID:    ctx.Uint64(utils.NetworkIdFlag.Name),
			addr:         ctx.String(utils.BNAddrFlag.Name),
			genKeyPath:   ctx.String(utils.GenKeyFlag.Name),
			nodeKeyFile:  ctx.String(utils.NodeKeyFileFlag.Name),
			nodeKeyHex:   ctx.String(utils.NodeKeyHexFlag.Name),
			natFlag:      ctx.String(utils.NATFlag.Name),
			netrestrict:  ctx.String(utils.NetrestrictFlag.Name),
			writeAddress: ctx.Bool(utils.WriteAddressFlag.Name),

			IPCPath:          "klay.ipc",
			DataDir:          ctx.String(utils.DataDirFlag.Name),
			HTTPPort:         DefaultHTTPPort,
			HTTPModules:      []string{"net"},
			HTTPVirtualHosts: []string{"localhost"},
			HTTPTimeouts:     rpc.DefaultHTTPTimeouts,
			WSPort:           DefaultWSPort,
			WSModules:        []string{"net"},
			GRPCPort:         DefaultGRPCPort,

			Logger: log.NewModuleLogger(log.CMDKBN),
		}
	)

	if err = nodecmd.CheckCommands(ctx); err != nil {
		return err
	}

	setIPC(ctx, &bcfg)
	// httptype is http or fasthttp
	if ctx.IsSet(utils.SrvTypeFlag.Name) {
		bcfg.HTTPServerType = ctx.String(utils.SrvTypeFlag.Name)
	}
	setHTTP(ctx, &bcfg)
	setWS(ctx, &bcfg)
	setgRPC(ctx, &bcfg)
	setAuthorizedNodes(ctx, &bcfg)

	// Check exit condition
	switch bcfg.checkCMDState() {
	case generateNodeKeySpecified:
		bcfg.generateNodeKey()
	case noPrivateKeyPathSpecified:
		return errors.New("Use --nodekey or --nodekeyhex to specify a private key")
	case nodeKeyDuplicated:
		return errors.New("Options --nodekey and --nodekeyhex are mutually exclusive")
	case writeOutAddress:
		bcfg.doWriteOutAddress()
	default:
		err = bcfg.readNodeKey()
		if err != nil {
			return err
		}
	}

	err = bcfg.validateNetworkParameter()
	if err != nil {
		return err
	}

	addr, err := net.ResolveUDPAddr("udp", bcfg.listenAddr)
	if err != nil {
		log.Fatalf("Failed to ResolveUDPAddr: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to ListenUDP: %v", err)
	}

	realaddr := conn.LocalAddr().(*net.UDPAddr)
	if bcfg.natm != nil {
		if !realaddr.IP.IsLoopback() {
			go nat.Map(bcfg.natm, nil, "udp", realaddr.Port, realaddr.Port, "Kaia node discovery")
		}
		// TODO: react to external IP changes over time.
		if ext, err := bcfg.natm.ExternalIP(); err == nil {
			realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
		}
	}

	cfg := discover.Config{
		NetworkID:       bcfg.networkID,
		PrivateKey:      bcfg.nodeKey,
		AnnounceAddr:    realaddr,
		NetRestrict:     bcfg.restrictList,
		Conn:            conn,
		Addr:            realaddr,
		Id:              discover.PubkeyID(&bcfg.nodeKey.PublicKey),
		NodeType:        p2p.ConvertNodeType(common.BOOTNODE),
		AuthorizedNodes: bcfg.AuthorizedNodes,
		DiscoverTypes:   discover.DiscoverTypesConfig{CN: true, PN: true, EN: true},
	}

	tab, err := discover.ListenUDP(&cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}

	node, err := New(&bcfg)
	if err != nil {
		return err
	}
	node.appendAPIs(NewBN(tab).APIs())
	if err := startNode(node); err != nil {
		return err
	}
	node.Wait()
	return nil
}

func startNode(node *Node) error {
	if err := node.Start(); err != nil {
		return err
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logger.Info("Got interrupt, shutting down...")
		go node.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				logger.Info("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
	return nil
}

func main() {
	// TODO-Kaia: remove `help` command
	app := utils.NewApp("", "the Kaia's bootnode command line interface")
	app.Name = "kbn"
	app.Copyright = "Copyright 2018-2024 The Kaia Authors"
	app.UsageText = app.Name + " [global options] [commands]"
	app.Flags = append(app.Flags, utils.BNAppFlags()...)
	app.Commands = []*cli.Command{
		nodecmd.VersionCommand,
		nodecmd.AttachCommand,
	}

	app.Action = bootnode

	app.CommandNotFound = nodecmd.CommandNotExist
	app.OnUsageError = nodecmd.OnUsageError

	app.Before = nodecmd.BeforeRunBootnode

	app.After = func(c *cli.Context) error {
		debug.Exit()
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
