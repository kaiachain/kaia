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
	"math/big"
	"strings"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	addressbookv2 "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/params"
	"github.com/urfave/cli/v2"
)

const (
	defaultNodeKeyPath = "/var/kcnd/data/klay/nodekey"
	defaultEndpoint    = "/var/kcnd/data/klay.ipc"
)

// AddressBookV2 method names, shared by the command wiring and the pre-flight
// dry-run (which packs calldata by name).
const (
	methodSuspendValidator   = "suspendValidator"
	methodUnsuspendValidator = "unsuspendValidator"
	methodReadyCandidate     = "readyCandidate"
	methodUnreadyCandidate   = "unreadyCandidate"
	methodReadyValidator     = "readyValidator"
	methodUnreadyValidator   = "unreadyValidator"
	methodPause              = "pause"
	methodResume             = "resume"
	methodExit               = "exit"
	methodOffboard           = "offboard"
)

type abv2TxFn func(*addressbookv2.AddressBookV2Transactor, *bind.TransactOpts, common.Address) (*types.Transaction, error)

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
	abv2NodeIdFlag = &cli.StringFlag{
		Name:     "node-id",
		Usage:    "Address of the target validator node",
		Required: true,
	}
)

var ValOpsCommand = &cli.Command{
	Name:     "valops",
	Usage:    "Validator state transition and suspension commands",
	Category: "PERMISSIONLESS COMMANDS",
	Subcommands: []*cli.Command{
		{
			Name:   "suspend-validator",
			Usage:  "Suspend a validator (requires suspender role)",
			Flags:  abv2SuspenderFlags(),
			Action: abv2SuspenderAction(methodSuspendValidator, (*addressbookv2.AddressBookV2Transactor).SuspendValidator),
		},
		{
			Name:   "unsuspend-validator",
			Usage:  "Unsuspend a validator (requires suspender role)",
			Flags:  abv2SuspenderFlags(),
			Action: abv2SuspenderAction(methodUnsuspendValidator, (*addressbookv2.AddressBookV2Transactor).UnsuspendValidator),
		},
		{
			Name:   "ready-candidate",
			Usage:  "Transition node to CandReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodReadyCandidate, (*addressbookv2.AddressBookV2Transactor).ReadyCandidate),
		},
		{
			Name:   "unready-candidate",
			Usage:  "Transition node out of CandReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodUnreadyCandidate, (*addressbookv2.AddressBookV2Transactor).UnreadyCandidate),
		},
		{
			Name:   "ready-validator",
			Usage:  "Transition node to ValReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodReadyValidator, (*addressbookv2.AddressBookV2Transactor).ReadyValidator),
		},
		{
			Name:   "unready-validator",
			Usage:  "Transition node out of ValReady state",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodUnreadyValidator, (*addressbookv2.AddressBookV2Transactor).UnreadyValidator),
		},
		{
			Name:   "pause",
			Usage:  "Pause the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodPause, (*addressbookv2.AddressBookV2Transactor).Pause),
		},
		{
			Name:   "resume",
			Usage:  "Resume the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodResume, (*addressbookv2.AddressBookV2Transactor).Resume),
		},
		{
			Name:   "exit",
			Usage:  "Exit the node from the validator set",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodExit, (*addressbookv2.AddressBookV2Transactor).Exit),
		},
		{
			Name:   "offboard",
			Usage:  "Offboard the node",
			Flags:  abv2Flags(),
			Action: abv2NodeOperatorAction(methodOffboard, (*addressbookv2.AddressBookV2Transactor).Offboard),
		},
	},
}

func abv2Flags() []cli.Flag {
	return []cli.Flag{abv2EndpointFlag, abv2PrivateKeyFlag}
}

func abv2SuspenderFlags() []cli.Flag {
	return []cli.Flag{abv2EndpointFlag, abv2PrivateKeyFlag, abv2NodeIdFlag}
}

func loadKey(ctx *cli.Context) (*ecdsa.PrivateKey, error) {
	if keyHex := ctx.String("private-key"); keyHex != "" {
		return crypto.HexToECDSA(strings.TrimPrefix(keyHex, "0x"))
	}
	return crypto.LoadECDSA(defaultNodeKeyPath)
}

func abv2Run(ctx *cli.Context, method string, nodeIdFn func(*bind.TransactOpts) common.Address, fn abv2TxFn) error {
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
	nodeId := nodeIdFn(opts)

	// Pre-flight: fail fast on an inevitable revert (e.g. SlotsFull) without gas.
	if reason := abv2DryRun(ec, method, opts.From, nodeId, opts.GasLimit); reason != "" {
		return fmt.Errorf("pre-flight check reverted: %s (no transaction sent)", reason)
	}

	tx, err := fn(transactor, opts, nodeId)
	if err != nil {
		return err
	}
	return printAndWait(ec, tx, opts.From)
}

// abv2DryRun returns the simulated transition's revert reason, or "" if it would
// succeed or can't be decoded (fail-open).
func abv2DryRun(ec *client.KaiaClient, method string, from, nodeId common.Address, gas uint64) string {
	abv2ABI, err := addressbookv2.AddressBookV2MetaData.GetAbi()
	if err != nil {
		return ""
	}
	data, err := abv2ABI.Pack(method, nodeId)
	if err != nil {
		return ""
	}
	to := system.AddressBookAddr
	msg := kaia.CallMsg{From: from, To: &to, Gas: gas, Data: data}
	if _, err := ec.CallContract(context.Background(), msg, nil); err != nil {
		reason := revertReasonFromErr(err)
		if reason == "SlotsFull" {
			reason += abv2SlotInfo(ec, method, nodeId)
		}
		return reason
	}
	return ""
}

// abv2SlotInfo explains a SlotsFull revert in plain terms by re-checking, against
// current chain state, the same guards the contract enforces and naming the ones
// that are violated. Mirrors NodeActions.{readyCandidate,pause,exit}; returns ""
// on any error or if no guard reproduces (fail-open).
func abv2SlotInfo(ec *client.KaiaClient, method string, nodeId common.Address) string {
	c, err := addressbookv2.NewAddressBookV2Caller(system.AddressBookAddr, ec)
	if err != nil {
		return ""
	}
	count := func(s valset.NodeState) (*big.Int, error) {
		return c.GetStateCount(nil, s.ToUint8())
	}
	var reasons []string
	switch method {
	case methodReadyCandidate:
		cand, e1 := count(valset.CandReady)
		maxNode, maxCand, e2 := c.GetMaxCounts(nil)
		nodes, e3 := c.GetAllNodesLength(nil)
		if e1 != nil || e2 != nil || e3 != nil {
			return ""
		}
		if cand.Cmp(maxCand) >= 0 {
			reasons = append(reasons, fmt.Sprintf("candidate-ready slots full (CandReady %s/%s)", cand, maxCand))
		}
		if nodes.Cmp(maxNode) >= 0 {
			reasons = append(reasons, fmt.Sprintf("max node count reached (%s/%s)", nodes, maxNode))
		}
	case methodPause:
		paused, e1 := count(valset.ValPaused)
		active, e2 := count(valset.ValActive)
		limits, e3 := c.GetSlotLimits(nil)
		if e1 != nil || e2 != nil || e3 != nil {
			return ""
		}
		if paused.Cmp(limits.MaxSlotAvailable) >= 0 {
			reasons = append(reasons, fmt.Sprintf("no pause slots available (ValPaused %s/%s)", paused, limits.MaxSlotAvailable))
		}
		if active.Cmp(limits.MinActiveCount) <= 0 {
			reasons = append(reasons, fmt.Sprintf("too few active validators to pause one (ValActive %s, min %s)", active, limits.MinActiveCount))
		}
	case methodExit:
		exiting, e1 := count(valset.ValExiting)
		active, e2 := count(valset.ValActive)
		limits, e3 := c.GetSlotLimits(nil)
		state, e4 := c.GetNodeState(nil, nodeId)
		if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
			return ""
		}
		if exiting.Cmp(limits.MaxSlotAvailable) >= 0 {
			reasons = append(reasons, fmt.Sprintf("no exit slots available (ValExiting %s/%s)", exiting, limits.MaxSlotAvailable))
		}
		// The active-floor guard only applies when exiting from ValActive.
		if state == valset.ValActive.ToUint8() && active.Cmp(limits.MinActiveCount) <= 0 {
			reasons = append(reasons, fmt.Sprintf("too few active validators to exit one (ValActive %s, min %s)", active, limits.MinActiveCount))
		}
	}
	if len(reasons) == 0 {
		return ""
	}
	return ": " + strings.Join(reasons, "; ")
}

// abv2SuspenderAction wraps a suspender-role command; reads node-id from --node-id flag.
func abv2SuspenderAction(method string, fn abv2TxFn) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		nodeId := common.HexToAddress(ctx.String("node-id"))
		return abv2Run(ctx, method, func(*bind.TransactOpts) common.Address { return nodeId }, fn)
	}
}

// abv2NodeOperatorAction wraps a node-operator command; node-id is derived from the key.
func abv2NodeOperatorAction(method string, fn abv2TxFn) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		return abv2Run(ctx, method, func(opts *bind.TransactOpts) common.Address { return opts.From }, fn)
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
