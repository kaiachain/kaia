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

	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	"github.com/kaiachain/kaia/crypto"
	"github.com/urfave/cli/v2"
)

// createNodeTag domain-separates the createNode proof; must match NodeVerifier.CREATE_NODE_TAG.
var createNodeTag = crypto.Keccak256Hash([]byte("KAIA_ADDRESS_BOOK_V2_CREATE_NODE_V1"))

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
		Usage: "Node ID address (default: derived from the signing key)",
	}
)

// SignCreateNodeCommand produces the nodeId-ownership signature that
// AddressBookV2.createNode requires. Offline: signs with the nodekey and needs
// no live endpoint (pass --chain-id explicitly).
var SignCreateNodeCommand = &cli.Command{
	Name:   "sign-create-node",
	Usage:  "Sign the nodeId-ownership proof required by AddressBookV2.createNode (offline)",
	Flags:  []cli.Flag{abv2PrivateKeyFlag, signChainIDFlag, signStakingFlag, signManagerFlag, signNodeIdFlag},
	Action: signCreateNodeAction,
}

func signCreateNodeAction(ctx *cli.Context) error {
	key, err := loadKey(ctx)
	if err != nil {
		return fmt.Errorf("load key: %w (use --private-key to specify explicitly)", err)
	}
	nodeId := crypto.PubkeyToAddress(key.PublicKey)
	if v := ctx.String("node-id"); v != "" {
		nodeId = common.HexToAddress(v)
	}
	digest := createNodeDigest(
		big.NewInt(ctx.Int64("chain-id")),
		common.HexToAddress(ctx.String("manager")),
		nodeId,
		common.HexToAddress(ctx.String("staking-contract")),
	)
	sig, err := crypto.Sign(digest, key)
	if err != nil {
		return err
	}
	sig[64] += 27 // AddressBookV2 verifies via OZ ECDSA.tryRecover, which expects v in {27,28}
	fmt.Println(hexutil.Encode(sig))
	return nil
}

// createNodeDigest reproduces NodeVerifier._verifyNodeIdProof:
// keccak256(abi.encode(TAG, chainId, addressBook, manager, nodeId, stakingContract)).
// Every field is static, so abi.encode is the 32-byte-word concatenation below.
func createNodeDigest(chainID *big.Int, manager, nodeId, staking common.Address) []byte {
	word := func(b []byte) []byte { return common.LeftPadBytes(b, 32) }
	buf := make([]byte, 0, 32*6)
	buf = append(buf, createNodeTag.Bytes()...)
	buf = append(buf, word(chainID.Bytes())...)
	buf = append(buf, word(system.AddressBookAddr.Bytes())...)
	buf = append(buf, word(manager.Bytes())...)
	buf = append(buf, word(nodeId.Bytes())...)
	buf = append(buf, word(staking.Bytes())...)
	return crypto.Keccak256(buf)
}
