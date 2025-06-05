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
// This file is derived from cmd/geth/consolecmd_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package nodecmd

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/kaiachain/kaia/v2/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 eth:1.0 governance:1.0 istanbul:1.0 kaia:1.0 net:1.0 personal:1.0 rpc:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "eth:1.0 kaia:1.0 net:1.0 rpc:1.0 web3:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	// Start a Kaia console, make sure it's cleaned up and terminate the console
	kaia := runKaia(t,
		"kaia-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none",
		"console")

	// Gather all the infos the welcome message needs to contain
	kaia.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	kaia.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	kaia.SetTemplateFunc("gover", runtime.Version)
	// TODO: Fix as in testAttachWelcome()
	kaia.SetTemplateFunc("klayver", func() string { return params.VersionWithCommit("") })
	kaia.SetTemplateFunc("apis", func() string { return ipcAPIs })
	kaia.SetTemplateFunc("datadir", func() string { return kaia.Datadir })

	// Verify the actual welcome message to the required template
	kaia.Expect(`
Welcome to the Kaia JavaScript console!

 instance: Klaytn/{{klayver}}/{{goos}}-{{goarch}}/{{gover}}
  datadir: {{datadir}}
  modules: {{apis}}

> {{.InputLine "exit"}}
`)
	kaia.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\klay` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "klay.ipc")
	}
	// Note: we need --shh because testAttachWelcome checks for default
	// list of ipc modules and shh is included there.
	kaia := runKaia(t,
		"kaia-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--ipcpath", ipc)

	waitForEndpoint(t, ipc, 10*time.Second)
	testAttachWelcome(t, kaia, "ipc:"+ipc, ipcAPIs)

	kaia.Interrupt()
	kaia.Kill()
}

func TestHTTPAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	kaia := runKaia(t,
		"kaia-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--rpc", "--rpcport", port)

	endpoint := "http://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 10*time.Second)
	testAttachWelcome(t, kaia, endpoint, httpAPIs)

	kaia.Interrupt()
	kaia.Kill()
}

func TestWSAttachWelcome(t *testing.T) {
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	kaia := runKaia(t,
		"kaia-test", "--port", "0", "--maxconnections", "0", "--nodiscover", "--nat", "none", "--ws", "--wsport", port)

	endpoint := "ws://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 10*time.Second)
	testAttachWelcome(t, kaia, endpoint, httpAPIs)

	kaia.Interrupt()
	kaia.Kill()
}

func testAttachWelcome(t *testing.T, klay *testKaia, endpoint, apis string) {
	// Attach to a running Kaia node and terminate immediately
	attach := runKaia(t, "kaia-test", "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	// The node version uses cmd/utils.gitCommit which is always empty.
	// TODO: Fix the cmd/utils.DefaultNodeConfig() to use cmd/utils/nodecmd.gitCommit
	// and then restore "klayver" to use gitCommit.
	attach.SetTemplateFunc("klayver", func() string { return params.VersionWithCommit("") })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return klay.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the Kaia JavaScript console!

 instance: Klaytn/{{klayver}}/{{goos}}-{{goarch}}/{{gover}}{{if ipc}}
  datadir: {{datadir}}{{end}}
  modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
