// Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from eth/tracers/tracers_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package tracers

import (
	"encoding/json"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/common/math"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/klaytn/klaytn/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type tracerTestdata struct {
	Genesis *struct {
		Config *params.ChainConfig     `json:"config"`
		Alloc  blockchain.GenesisAlloc `json:"alloc"`
	} `json:"genesis"`

	Context *struct {
		BaseFee    *math.HexOrDecimal256 `json:"baseFeePerGas"`
		MixHash    hexutil.Bytes         `json:"mixHash"`
		Number     math.HexOrDecimal64   `json:"number"`
		Timestamp  math.HexOrDecimal64   `json:"timestamp"`
		BlockScore *math.HexOrDecimal256 `json:"blockScore"`
	} `json:"context"`

	Input  string          `json:"input"`
	Result json.RawMessage `json:"result"`
}

func TestPrestateTracer(t *testing.T) {
	forEachJson(t, "testdata/prestate_tracer", func(t *testing.T, tc *tracerTestdata) {
		tracer, err := New("prestateTracer", new(Context), false)
		require.NoError(t, err)
		runTracer(t, tc, tracer)
	})
}

func forEachJson(t *testing.T, dir string, f func(t *testing.T, tc *tracerTestdata)) {
	files, err := os.ReadDir(dir)
	require.NoError(t, err)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		blob, err := os.ReadFile(filepath.Join(dir, file.Name()))
		require.NoError(t, err, file.Name())

		tc := new(tracerTestdata)
		require.NoError(t, json.Unmarshal(blob, tc), file.Name())

		t.Run(strings.TrimSuffix(file.Name(), ".json"), func(t *testing.T) {
			f(t, tc)
		})
	}
}

func runTracer(t *testing.T, tc *tracerTestdata, tracer vm.Tracer) (*types.Transaction, *blockchain.ExecutionResult, json.RawMessage) {
	// Parse the raw transaction
	var tx *types.Transaction
	require.NoError(t, rlp.DecodeBytes(common.FromHex(tc.Input), &tx))

	// Create the EVM environment at the point of tx execution
	var (
		config = tc.Genesis.Config
		alloc  = tc.Genesis.Alloc

		header = &types.Header{ // Must have all fields used in NewEVMBlockContext and NewEVMTxContext
			BaseFee:    (*big.Int)(tc.Context.BaseFee),
			MixHash:    tc.Context.MixHash,
			Number:     new(big.Int).SetUint64(uint64(tc.Context.Number)),
			Time:       new(big.Int).SetUint64(uint64(tc.Context.Timestamp)),
			BlockScore: (*big.Int)(tc.Context.BlockScore),
		}

		signer       = types.MakeSigner(config, header.Number)
		blockContext = blockchain.NewEVMBlockContext(header, nil, &common.Address{}) // stub author (COINBASE) to 0x0
		txContext    = blockchain.NewEVMTxContext(tx, header, config)
		statedb      = tests.MakePreState(database.NewMemoryDBManager(), alloc)
		evm          = vm.NewEVM(blockContext, txContext, statedb, config, &vm.Config{Debug: true, Tracer: tracer})
	)

	// Run the transaction with tracer enabled
	fork.SetHardForkBlockNumberConfig(config) // needed by IntrinsicGas()
	msg, err := tx.AsMessageWithAccountKeyPicker(signer, statedb, header.Number.Uint64())
	require.NoError(t, err)

	st := blockchain.NewStateTransition(evm, msg)
	execResult, err := st.TransitionDb()
	require.NoError(t, err)

	var tracerResult json.RawMessage
	switch tracer := tracer.(type) {
	case *Tracer:
		tracerResult, err = tracer.GetResult()
	}
	require.NoError(t, err)
	assert.JSONEq(t, string(tc.Result), string(tracerResult))

	return msg, execResult, tracerResult
}
