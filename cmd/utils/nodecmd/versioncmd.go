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
// This file is derived from cmd/geth/misccmd.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package nodecmd

import (
	"fmt"

	"github.com/kaiachain/kaia/v2/params"
	"github.com/urfave/cli/v2"
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""

	// Git tag (set via linker flags if exists)
	gitTag = ""
)

var VersionCommand = &cli.Command{
	Action:    version,
	Name:      "version",
	Usage:     "Show version number",
	ArgsUsage: " ",
	Category:  "MISCELLANEOUS COMMANDS",
}

func version(ctx *cli.Context) error {
	fmt.Print("Kaia ")
	if gitTag != "" {
		// stable version
		fmt.Println(params.Version)
	} else {
		// unstable version
		fmt.Println(params.VersionWithCommit(gitCommit))
	}
	return nil
}

// GetGitCommit returns gitCommit set by linker flags.
func GetGitCommit() string {
	return gitCommit
}
