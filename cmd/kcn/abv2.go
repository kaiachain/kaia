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
	"path/filepath"
	"slices"
	"strings"
	"time"

	kaia "github.com/kaiachain/kaia"
	"github.com/kaiachain/kaia/accounts/abi"
	"github.com/kaiachain/kaia/accounts/abi/bind"
	"github.com/kaiachain/kaia/blockchain/system"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/client"
	"github.com/kaiachain/kaia/cmd/utils"
	"github.com/kaiachain/kaia/common"
	"github.com/kaiachain/kaia/common/hexutil"
	addressbookv2 "github.com/kaiachain/kaia/contracts/bindings/addressbookv2"
	"github.com/kaiachain/kaia/crypto"
	"github.com/kaiachain/kaia/kaiax/valset"
	"github.com/kaiachain/kaia/networks/rpc"
	"github.com/kaiachain/kaia/node"
	"github.com/kaiachain/kaia/params"
	"github.com/urfave/cli/v2"
)

const defaultTimeout = 2 * time.Minute

// minEpochTimeLeft is the time before the next epoch block below which a ready
// transition prints a warning. The onboarding portal uses the same value.
const minEpochTimeLeft = 5 * time.Minute

// defaultDatadir is the data directory of an rpm install (DATA_DIR in /etc/kcnd/conf/kcnd.conf).
const defaultDatadir = "/var/kcnd/data"

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
		Usage: "RPC or IPC endpoint (default: <datadir>/klay.ipc)",
	}
	abv2PrivateKeyFlag = &cli.StringFlag{
		Name:  "private-key",
		Usage: "Hex-encoded private key, used instead of the node key file. Unsafe: a key on the command line stays in shell history and process listings",
	}
	abv2DatadirFlag = &cli.StringFlag{
		Name:  "datadir",
		Usage: "Node data directory (rpm installs use the default; for a tar install, DATA_DIR in conf/kcnd.conf). The key is loaded from <datadir>/nodekey if it exists, else <datadir>/klay/nodekey, as the node does",
		Value: defaultDatadir,
	}
	abv2NodeIdFlag = &cli.StringFlag{
		Name:     "node-id",
		Usage:    "Address of the target validator node",
		Required: true,
	}
	abv2TimeoutFlag = &cli.DurationFlag{
		Name:  "timeout",
		Usage: "How long to wait for the transaction receipt",
		Value: defaultTimeout,
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
		GenerateKeysCommand,
		SignCreateNodeCommand,
	},
}

func abv2Flags() []cli.Flag {
	return []cli.Flag{abv2EndpointFlag, abv2PrivateKeyFlag, abv2DatadirFlag, abv2TimeoutFlag}
}

func abv2SuspenderFlags() []cli.Flag {
	return []cli.Flag{abv2EndpointFlag, abv2PrivateKeyFlag, abv2DatadirFlag, abv2TimeoutFlag, abv2NodeIdFlag}
}

// nodeFilePath returns the path the node uses for a key file under datadir.
// It resolves the path as the node does (node.Config.ResolvePath):
// <datadir>/<name> if that file exists, else <datadir>/klay/<name>.
func nodeFilePath(datadir, name string) string {
	cfg := node.Config{DataDir: datadir, Name: utils.ClientIdentifier}
	return cfg.ResolvePath(name)
}

// endpointOf returns --endpoint if set, else the IPC socket of the node under --datadir.
func endpointOf(ctx *cli.Context) string {
	if endpoint := ctx.String("endpoint"); endpoint != "" {
		return endpoint
	}
	return filepath.Join(ctx.String("datadir"), "klay.ipc")
}

// datadirUsed reports whether the command reads the key or the endpoint from --datadir.
func datadirUsed(ctx *cli.Context) bool {
	return ctx.String("private-key") == "" || ctx.String("endpoint") == ""
}

// datadirSource describes where the --datadir value came from.
func datadirSource(ctx *cli.Context) string {
	if ctx.IsSet("datadir") {
		return "--datadir"
	}
	return "default"
}

// parseAddress parses a non-zero 20-byte hex address given to --flag.
func parseAddress(flag, value string) (common.Address, error) {
	if !common.IsHexAddress(value) {
		return common.Address{}, fmt.Errorf("invalid --%s %q: expected a 20-byte hex address (40 hex digits, with or without 0x)", flag, value)
	}
	addr := common.HexToAddress(value)
	if addr == (common.Address{}) {
		return common.Address{}, fmt.Errorf("invalid --%s: zero address", flag)
	}
	return addr, nil
}

// formatKAIA formats an amount in kei as KAIA.
func formatKAIA(kei *big.Int) string {
	return new(big.Float).Quo(new(big.Float).SetInt(kei), big.NewFloat(params.KAIA)).Text('f', 6) + " KAIA"
}

// loadKey returns the --private-key key if set, else the raw nodekey under --datadir.
func loadKey(ctx *cli.Context) (*ecdsa.PrivateKey, error) {
	if keyHex := ctx.String("private-key"); keyHex != "" {
		return crypto.HexToECDSA(strings.TrimPrefix(keyHex, "0x"))
	}
	key, err := crypto.LoadECDSA(nodeFilePath(ctx.String("datadir"), node.DatadirPrivateKey))
	if err != nil {
		return nil, fmt.Errorf("load key: %w (check --datadir; the default is "+defaultDatadir+")", err)
	}
	return key, nil
}

func abv2Run(ctx *cli.Context, method string, nodeIdFn func(*bind.TransactOpts) common.Address, fn abv2TxFn) error {
	key, err := loadKey(ctx)
	if err != nil {
		return err
	}
	endpoint := endpointOf(ctx)
	ec, transactor, err := dialABv2(endpoint)
	if err != nil {
		return fmt.Errorf("%w (check --endpoint or --datadir)", err)
	}
	defer ec.Close()
	chainID, err := ec.ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("get chain ID from %s: %w", endpoint, err)
	}

	opts := bind.NewKeyedTransactor(key)
	nodeId := nodeIdFn(opts)

	fmt.Printf("method:   %s\n", method)
	if datadirUsed(ctx) {
		fmt.Printf("datadir:  %s (%s)\n", ctx.String("datadir"), datadirSource(ctx))
	}
	fmt.Printf("endpoint: %s\n", endpoint)
	fmt.Printf("chain ID: %s\n", chainID)
	fmt.Printf("sender:   %s\n", opts.From.Hex())
	fmt.Printf("node ID:  %s\n", nodeId.Hex())

	caller, err := addressbookv2.NewAddressBookV2Caller(system.AddressBookAddr, ec)
	if err != nil {
		return err
	}
	if method == methodReadyCandidate || method == methodReadyValidator {
		if err := abv2PrintEpochWarnings(ec, caller, method, nodeId); err != nil {
			return err
		}
	}

	// Pre-flight: fail fast on an inevitable revert (e.g. SlotsFull) without gas.
	if reason := abv2DryRun(ec, method, opts.From, nodeId); reason != "" {
		return fmt.Errorf("pre-flight check reverted: %s (no transaction sent)", reason)
	}
	if opts.GasLimit, opts.GasPrice, err = abv2GasParams(ec, method, opts.From, nodeId); err != nil {
		return fmt.Errorf("%w (no transaction sent)", err)
	}

	tx, err := fn(transactor, opts, nodeId)
	if err != nil {
		return err
	}
	if err := printAndWait(ec, tx, opts.From, ctx.Duration("timeout")); err != nil {
		return err
	}
	return abv2PrintNodeState(caller, nodeId)
}

// abv2PrintEpochWarnings prints the time left before the next epoch block and a
// warning for each way the ready transition may not take effect at that block.
// It never refuses: the contract has no epoch-time condition.
func abv2PrintEpochWarnings(ec *client.KaiaClient, caller *addressbookv2.AddressBookV2Caller, method string, nodeId common.Address) error {
	head, err := ec.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("get head block: %w", err)
	}
	interval, err := caller.EpochBlockInterval(nil)
	if err != nil {
		return fmt.Errorf("get epoch block interval: %w (AddressBookV2 is active only after the Permissionless fork)", err)
	}
	if interval.Sign() <= 0 {
		return fmt.Errorf("invalid epoch block interval %s", interval)
	}
	nextEpoch, timeLeft := epochTimeLeft(head.Number, interval)
	fmt.Printf("epoch:    next epoch block %s, about %s left\n", nextEpoch, timeLeft)

	// A zero timeoutAt never expires.
	var timeoutIn *time.Duration
	if method == methodReadyValidator {
		info, err := caller.GetNodeInfo(nil, nodeId)
		if err != nil {
			return fmt.Errorf("get node info: %w", err)
		}
		if info.TimeoutAt != nil && info.TimeoutAt.Sign() > 0 {
			d := time.Duration(info.TimeoutAt.Int64()-int64(head.Time.Uint64())) * time.Second
			timeoutIn = &d
			fmt.Printf("timeout:  %s (about %s from now)\n", time.Unix(info.TimeoutAt.Int64(), 0).UTC().Format(time.RFC3339), d)
		}
	}
	epochLength := time.Duration(interval.Int64()) * time.Second
	for _, w := range readyWarnings(method, nextEpoch, timeLeft, epochLength, timeoutIn) {
		fmt.Println("WARNING: " + w)
	}
	return nil
}

// readyWarnings returns the warnings for a ready transition sent timeLeft before
// the epoch block nextEpoch. timeoutIn is the time until the node's idle timeout,
// or nil if it has none. A ValReady or ValInactive node becomes Registered at the
// first block whose time reaches its idle timeout. At an epoch block the epoch
// transition runs first, so a timeout on the epoch block itself does not stop a
// promotion to ValActive; it still applies if the node is not promoted.
func readyWarnings(method string, nextEpoch *big.Int, timeLeft, epochLength time.Duration, timeoutIn *time.Duration) []string {
	var ws []string
	late := timeLeft < minEpochTimeLeft
	if late {
		ws = append(ws, fmt.Sprintf("less than %s is left before epoch block %s. If the transaction is included after that block, the transition takes effect one epoch later (about %s)",
			minEpochTimeLeft, nextEpoch, epochLength))
	}
	if method != methodReadyValidator || timeoutIn == nil {
		return ws
	}
	switch {
	case *timeoutIn < timeLeft:
		ws = append(ws, fmt.Sprintf("the node's idle timeout comes before epoch block %s. The node will become Registered, not ValActive, and must go through ready-candidate again", nextEpoch))
	case late && *timeoutIn < timeLeft+epochLength:
		ws = append(ws, fmt.Sprintf("if the transaction is included after epoch block %s, the node's idle timeout comes before the next epoch block. The node will then become Registered, not ValActive", nextEpoch))
	}
	return ws
}

// epochTimeLeft returns the next epoch block after head and the time until it,
// counting one second per block.
func epochTimeLeft(head, interval *big.Int) (*big.Int, time.Duration) {
	blocksLeft := new(big.Int).Sub(interval, new(big.Int).Mod(head, interval))
	return new(big.Int).Add(head, blocksLeft), time.Duration(blocksLeft.Int64()) * time.Second
}

// abv2PrintNodeState prints the node's state and suspension after a transaction.
func abv2PrintNodeState(caller *addressbookv2.AddressBookV2Caller, nodeId common.Address) error {
	state, err := caller.GetNodeState(nil, nodeId)
	if err != nil {
		return fmt.Errorf("get node state: %w", err)
	}
	suspended, err := caller.GetSuspendedValidators(nil)
	if err != nil {
		return fmt.Errorf("get suspended validators: %w", err)
	}
	fmt.Printf("state:    %s\n", valset.NodeState(state))
	fmt.Printf("suspended: %t\n", slices.Contains(suspended, nodeId))
	return nil
}

// abv2CallMsg builds the call of method(nodeId) on AddressBookV2 sent by from.
func abv2CallMsg(method string, from, nodeId common.Address) (kaia.CallMsg, error) {
	abv2ABI, err := addressbookv2.AddressBookV2MetaData.GetAbi()
	if err != nil {
		return kaia.CallMsg{}, err
	}
	data, err := abv2ABI.Pack(method, nodeId)
	if err != nil {
		return kaia.CallMsg{}, err
	}
	to := system.AddressBookAddr
	return kaia.CallMsg{From: from, To: &to, Data: data}, nil
}

// abv2GasParams returns the gas limit (estimate plus 50%) and gas price for the
// transition, and fails if the sender cannot pay gasLimit * gasPrice up front.
func abv2GasParams(ec *client.KaiaClient, method string, from, nodeId common.Address) (uint64, *big.Int, error) {
	msg, err := abv2CallMsg(method, from, nodeId)
	if err != nil {
		return 0, nil, err
	}
	estimate, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		return 0, nil, fmt.Errorf("estimate gas: %w", err)
	}
	gasLimit := estimate + estimate/2
	gasPrice, err := ec.SuggestGasPrice(context.Background())
	if err != nil {
		return 0, nil, fmt.Errorf("get gas price: %w", err)
	}
	balance, err := ec.BalanceAt(context.Background(), from, nil)
	if err != nil {
		return 0, nil, fmt.Errorf("get balance of %s: %w", from.Hex(), err)
	}
	required := new(big.Int).Mul(new(big.Int).SetUint64(gasLimit), gasPrice)
	if balance.Cmp(required) < 0 {
		return 0, nil, fmt.Errorf("sender %s needs %s for gas but holds %s; send KAIA to this address",
			from.Hex(), formatKAIA(required), formatKAIA(balance))
	}
	return gasLimit, gasPrice, nil
}

// abv2DryRun returns the simulated transition's revert reason, or "" if it would
// succeed or can't be decoded (fail-open).
func abv2DryRun(ec *client.KaiaClient, method string, from, nodeId common.Address) string {
	msg, err := abv2CallMsg(method, from, nodeId)
	if err != nil {
		return ""
	}
	msg.Gas = params.UpperGasLimit
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
		nodeId, err := parseAddress("node-id", ctx.String("node-id"))
		if err != nil {
			return err
		}
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

func printAndWait(ec *client.KaiaClient, tx *types.Transaction, from common.Address, timeout time.Duration) error {
	fmt.Printf("tx: %s\n", tx.Hash().Hex())
	fmt.Println("waiting for receipt...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	receipt, err := bind.WaitMined(ctx, ec, tx)
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("no receipt for %s after %s; the outcome is unknown. Check the transaction hash before sending again", tx.Hash().Hex(), timeout)
	}
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
