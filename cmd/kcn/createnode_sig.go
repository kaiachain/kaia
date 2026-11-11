// Copyright 2026 The Kaia Authors
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

package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/urfave/cli/v2"
)

var (
	signChainIDFlag = &cli.Int64Flag{
		Name:     "chain-id",
		Usage:    "Chain ID the node runs on (8217 mainnet, 1001 kairos)",
		Required: true,
	}
	signStakingFlag = &cli.StringFlag{
		Name:     "staking-contract",
		Usage:    "CnStaking contract address bound into the proof",
		Required: true,
	}
	signManagerFlag = &cli.StringFlag{
		Name:     "manager",
		Usage:    "Address that will send createNode (becomes the node's manager)",
		Required: true,
	}
	signNodeIdFlag = &cli.StringFlag{
		Name:  "node-id",
		Usage: "Expected node ID address; fails if it does not match the signing key",
	}
)

// SignCreateNodeCommand produces the nodeId-ownership signature that
// AddressBookV2.createNode requires. Offline: signs with the nodekey and needs
// no live endpoint (pass --chain-id explicitly).
var SignCreateNodeCommand = &cli.Command{
	Name:   "sign-create-node",
	Usage:  "Sign the nodeId-ownership proof required by AddressBookV2.createNode (offline)",
	Flags:  []cli.Flag{abv2PrivateKeyFlag, abv2DatadirFlag, signChainIDFlag, signStakingFlag, signManagerFlag, signNodeIdFlag},
	Action: signCreateNodeAction,
}

func signCreateNodeAction(ctx *cli.Context) error {
	key, err := loadKey(ctx)
	if err != nil {
		return err
	}
	nodeId := crypto.PubkeyToAddress(key.PublicKey)
	if v := ctx.String("node-id"); v != "" {
		expected, err := parseAddress("node-id", v)
		if err != nil {
			return err
		}
		if expected != nodeId {
			return fmt.Errorf("--node-id %s does not match the signing key's address %s", v, nodeId.Hex())
		}
	}
	chainID := ctx.Int64("chain-id")
	if chainID <= 0 {
		return fmt.Errorf("invalid --chain-id %d: must be positive (8217 mainnet, 1001 kairos)", chainID)
	}
	manager, err := parseAddress("manager", ctx.String("manager"))
	if err != nil {
		return err
	}
	staking, err := parseAddress("staking-contract", ctx.String("staking-contract"))
	if err != nil {
		return err
	}
	sig, err := system.SignCreateNodeProof(key, big.NewInt(chainID), manager, nodeId, staking)
	if err != nil {
		return err
	}
	// The inputs go to stderr so that stdout holds only the signature.
	fmt.Fprintf(os.Stderr, "chain ID:         %d\nnode ID:          %s\nmanager:          %s\nstaking contract: %s\n",
		chainID, nodeId.Hex(), manager.Hex(), staking.Hex())
	fmt.Println(hexutil.Encode(sig))
	return nil
}
