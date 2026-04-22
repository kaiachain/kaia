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
	"bytes"
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"strings"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	addressbookv2 "github.com/kaiachain/kaia/contracts_permissionless/contracts/AddressBookV2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/urfave/cli/v2"
)

const (
	defaultNodeKeyPath = "/var/kcnd/data/klay/nodekey"
	defaultEndpoint    = "/var/kcnd/data/klay.ipc"
)

type (
	abv2RunFn func(*addressbookv2.AddressBookV2Transactor, *bind.TransactOpts, *ecdsa.PrivateKey) (*types.Transaction, error)
	abv2TxFn  func(*addressbookv2.AddressBookV2Transactor, *bind.TransactOpts, common.Address) (*types.Transaction, error)
)

var (
	abv2EndpointFlag = &cli.StringFlag{
		Name:  "endpoint",
		Usage: "RPC or IPC endpoint",
		Value: defaultEndpoint,
	}
	abv2PrivateKeyFlag = &cli.StringFlag{
		Name:  "private-key",
		Usage: "Hex-encoded private key (default: load from " + defaultNodeKeyPath + ")",
	}
)

var ABv2Command = &cli.Command{
	Name:     "abv2",
	Usage:    "AddressBookV2 contract commands",
	Category: "PERMISSIONLESS COMMANDS",
	Subcommands: []*cli.Command{
		{
			Name:      "suspend-validator",
			Usage:     "Suspend a validator (requires suspender role)",
			ArgsUsage: "<node-id>",
			Flags:     abv2Flags(),
			Action:    abv2SuspenderAction((*addressbookv2.AddressBookV2Transactor).SuspendValidator),
		},
		{
			Name:      "unsuspend-validator",
			Usage:     "Unsuspend a validator (requires suspender role)",
			ArgsUsage: "<node-id>",
			Flags:     abv2Flags(),
			Action:    abv2SuspenderAction((*addressbookv2.AddressBookV2Transactor).UnsuspendValidator),
		},
		{
			Name:   "ready-candidate",
			Usage:  "Transition node to CandReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).ReadyCandidate),
		},
		{
			Name:   "unready-candidate",
			Usage:  "Transition node out of CandReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).UnreadyCandidate),
		},
		{
			Name:   "ready-validator",
			Usage:  "Transition node to ValReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).ReadyValidator),
		},
		{
			Name:   "unready-validator",
			Usage:  "Transition node out of ValReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).UnreadyValidator),
		},
		{
			Name:   "pause",
			Usage:  "Pause the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).Pause),
		},
		{
			Name:   "resume",
			Usage:  "Resume the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).Resume),
		},
		{
			Name:   "exit",
			Usage:  "Exit the node from the validator set",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).Exit),
		},
		{
			Name:   "offboard",
			Usage:  "Offboard the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction((*addressbookv2.AddressBookV2Transactor).Offboard),
		},
	},
}

func abv2Flags() []cli.Flag {
	return []cli.Flag{abv2EndpointFlag, abv2PrivateKeyFlag}
}

func loadKey(ctx *cli.Context) (*ecdsa.PrivateKey, error) {
	if keyHex := ctx.String("private-key"); keyHex != "" {
		return crypto.HexToECDSA(strings.TrimPrefix(keyHex, "0x"))
	}
	return crypto.LoadECDSA(defaultNodeKeyPath)
}

func abv2Run(ctx *cli.Context, fn abv2RunFn) error {
	key, err := loadKey(ctx)
	if err != nil {
		return fmt.Errorf("load key: %w (use --private-key to specify explicitly)", err)
	}
	ec, transactor, err := dialABv2(ctx.String("endpoint"))
	if err != nil {
		return fmt.Errorf("%w (use --endpoint to specify explicitly)", err)
	}
	defer ec.Close()

	opts := bind.NewKeyedTransactor(key)
	opts.GasLimit = params.UpperGasLimit
	tx, err := fn(transactor, opts, key)
	if err != nil {
		return err
	}
	return printAndWait(ec, tx, opts.From)
}

// abv2SuspenderAction wraps a suspender-role command; parses <node-id> from args.
func abv2SuspenderAction(fn abv2TxFn) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if ctx.Args().Len() != 1 {
			return fmt.Errorf("usage: %s <node-id>", ctx.Command.Name)
		}
		nodeId := common.HexToAddress(ctx.Args().Get(0))
		return abv2Run(ctx, func(t *addressbookv2.AddressBookV2Transactor, opts *bind.TransactOpts, _ *ecdsa.PrivateKey) (*types.Transaction, error) {
			return fn(t, opts, nodeId)
		})
	}
}

// abv2NodeOperatorAction wraps a node-operator command; node-id is derived from the key.
func abv2NodeOperatorAction(fn abv2TxFn) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return abv2Run(ctx, func(t *addressbookv2.AddressBookV2Transactor, opts *bind.TransactOpts, key *ecdsa.PrivateKey) (*types.Transaction, error) {
			return fn(t, opts, crypto.PubkeyToAddress(key.PublicKey))
		})
	}
}

func dialABv2(endpoint string) (*client.KaiaClient, *addressbookv2.AddressBookV2Transactor, error) {
	ec, err := client.Dial(endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("connect to %s: %w", endpoint, err)
	}
	t, err := addressbookv2.NewAddressBookV2Transactor(system.AddressBookAddr, ec)
	if err != nil {
		ec.Close()
		return nil, nil, err
	}
	return ec, t, nil
}

func printAndWait(ec *client.KaiaClient, tx *types.Transaction, from common.Address) error {
	fmt.Printf("tx: %s\n", tx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), ec, tx)
	if err != nil {
		return fmt.Errorf("wait for receipt: %w", err)
	}
	if receipt.Status == types.ReceiptStatusFailed {
		if reason := revertReason(ec, tx, from); reason != "" {
			return fmt.Errorf("transaction failed: %s", reason)
		}
		return errors.New("transaction failed")
	}
	fmt.Println("status: success")
	return nil
}

func revertReason(ec *client.KaiaClient, tx *types.Transaction, from common.Address) string {
	msg := kaia.CallMsg{
		From:     from,
		To:       tx.To(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
		Value:    tx.Value(),
		Data:     tx.Data(),
	}
	_, err := ec.CallContract(context.Background(), msg, nil)
	if err == nil {
		return ""
	}
	return revertReasonFromErr(err)
}

func revertReasonFromErr(err error) string {
	var de rpc.DataError
	if errors.As(err, &de) {
		if hexStr, ok := de.ErrorData().(string); ok {
			if b, decErr := hexutil.Decode(hexStr); decErr == nil {
				return decodeCustomError(b)
			}
		}
	}
	return ""
}

func decodeCustomError(data []byte) string {
	if len(data) < 4 {
		return ""
	}
	if reason, err := abi.UnpackRevert(data); err == nil {
		return reason
	}
	abv2ABI, err := addressbookv2.AddressBookV2MetaData.GetAbi()
	if err != nil {
		return ""
	}
	for _, abiErr := range abv2ABI.Errors {
		if bytes.Equal(data[:4], abiErr.ID[:4]) {
			if len(data) == 4 {
				return abiErr.Name
			}
			unpacked, err := abiErr.Unpack(data)
			if err != nil {
				return abiErr.Name
			}
			return fmt.Sprintf("%s(%v)", abiErr.Name, unpacked)
		}
	}
	return ""
}
