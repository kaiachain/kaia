// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
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
// This file is derived from core/block_validator_test.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"math/big"
	"runtime"
	"testing"
	"time"

	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/blockchain/vm"
	"github.com/kaiachain/kaia/v2/common"
	"github.com/kaiachain/kaia/v2/consensus/gxhash"
	"github.com/kaiachain/kaia/v2/crypto"
	"github.com/kaiachain/kaia/v2/params"
	"github.com/kaiachain/kaia/v2/storage/database"
	"github.com/stretchr/testify/assert"
)

// Tests that simple header verification works, for both good and bad blocks.
func TestHeaderVerification(t *testing.T) {
	// Create a simple chain to verify
	var (
		testdb    = database.NewMemoryDBManager()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), testdb, 8, nil)
	)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	// Run the header checker for blocks one-by-one, checking for both valid and invalid nonces
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, gxhash.NewFaker(), vm.Config{})
	defer chain.Stop()

	for i := 0; i < len(blocks); i++ {
		for j, valid := range []bool{true, false} {
			var results <-chan error

			if valid {
				engine := gxhash.NewFaker()
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true})
			} else {
				engine := gxhash.NewFakeFailer(headers[i].Number.Uint64())
				_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true})
			}
			// Wait for the verification result
			select {
			case result := <-results:
				if (result == nil) != valid {
					t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, result, valid)
				}
			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
			// Make sure no more data is returned
			select {
			case result := <-results:
				t.Fatalf("test %d.%d: unexpected result returned: %v", i, j, result)
			case <-time.After(25 * time.Millisecond):
			}
		}
		chain.InsertChain(blocks[i : i+1])
	}
}

func TestVerifyBlockBody(t *testing.T) {
	testcases := []struct {
		baseFee *big.Int
		txData  types.TxInternalData
		err     bool
	}{
		{
			big.NewInt(100),
			&types.TxInternalDataLegacy{
				GasLimit: 1,
				Price:    big.NewInt(200),
				Payload:  []byte("abcdef"),
			},
			false,
		},

		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(10),
				GasTipCap: big.NewInt(10),
				Payload:   []byte("abcdef"),
			},
			true,
		},
		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(100),
				GasTipCap: big.NewInt(10),
				Payload:   []byte("abcdef"),
			},
			false,
		},
		{
			big.NewInt(100),
			&types.TxInternalDataEthereumDynamicFee{
				ChainID:   big.NewInt(1),
				GasLimit:  123457,
				GasFeeCap: big.NewInt(200),
				GasTipCap: big.NewInt(10),
			},
			false,
		},
	}

	// Create a simple chain to verify
	var (
		testdb  = database.NewMemoryDBManager()
		gspec   = &Genesis{Config: params.TestChainConfig}
		genesis = gspec.MustCommit(testdb)
	)

	// We don't need istanbul instance here because validateBody is in BlockValidator instance
	GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), testdb, 8, nil)
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, gxhash.NewFaker(), vm.Config{})
	defer chain.Stop()

	// Generate a batch of accounts to start with
	privKey, _ := crypto.GenerateKey()
	signer := types.LatestSignerForChainID(big.NewInt(1))
	var block *types.Block
	for _, testcase := range testcases {
		// Generate a block header
		header := &types.Header{
			ParentHash: chain.hc.currentHeaderHash,
			Number:     common.Big1,
			GasUsed:    0,
			Extra:      []byte{},
			BaseFee:    testcase.baseFee,
		}

		// Generate a block with tx
		tx := types.NewTx(testcase.txData)
		tx.Sign(signer, privKey)
		block = types.NewBlock(header, append(types.Transactions{}, tx), nil)

		err := chain.validator.ValidateBody(block)
		if errExist := err != nil; errExist != testcase.err {
			assert.Error(t, err)
		}
	}
}

// Tests that concurrent header verification works, for both good and bad blocks.
func TestHeaderConcurrentVerification2(t *testing.T)  { testHeaderConcurrentVerification(t, 2) }
func TestHeaderConcurrentVerification8(t *testing.T)  { testHeaderConcurrentVerification(t, 8) }
func TestHeaderConcurrentVerification32(t *testing.T) { testHeaderConcurrentVerification(t, 32) }

func testHeaderConcurrentVerification(t *testing.T, threads int) {
	// Create a simple chain to verify
	var (
		testdb    = database.NewMemoryDBManager()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), testdb, 8, nil)
	)
	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	// Set the number of threads to verify on
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

	// Run the header checker for the entire block chain at once both for a valid and
	// also an invalid chain (enough if one arbitrary block is invalid).
	for i, valid := range []bool{true, false} {
		var results <-chan error

		if valid {
			chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, gxhash.NewFaker(), vm.Config{})
			_, results = chain.engine.VerifyHeaders(chain, headers, seals)
			chain.Stop()
		} else {
			chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, gxhash.NewFakeFailer(uint64(len(headers)-1)), vm.Config{})
			_, results = chain.engine.VerifyHeaders(chain, headers, seals)
			chain.Stop()
		}
		// Wait for all the verification results
		checks := make(map[int]error)
		for j := 0; j < len(blocks); j++ {
			select {
			case result := <-results:
				checks[j] = result

			case <-time.After(time.Second):
				t.Fatalf("test %d.%d: verification timeout", i, j)
			}
		}
		// Check nonce check validity
		for j := 0; j < len(blocks); j++ {
			want := valid || (j < len(blocks)-2) // We chose the last-but-one nonce in the chain to fail
			if (checks[j] == nil) != want {
				t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, checks[j], want)
			}
			if !want {
				// A few blocks after the first error may pass verification due to concurrent
				// workers. We don't care about those in this test, just that the correct block
				// errors out.
				break
			}
		}
		// Make sure no more data is returned
		select {
		case result := <-results:
			t.Fatalf("test %d: unexpected result returned: %v", i, result)
		case <-time.After(25 * time.Millisecond):
		}
	}
}

// Tests that aborting a header validation indeed prevents further checks from being
// run, as well as checks that no left-over goroutines are leaked.
func TestHeaderConcurrentAbortion2(t *testing.T)  { testHeaderConcurrentAbortion(t, 2) }
func TestHeaderConcurrentAbortion8(t *testing.T)  { testHeaderConcurrentAbortion(t, 8) }
func TestHeaderConcurrentAbortion32(t *testing.T) { testHeaderConcurrentAbortion(t, 32) }

func testHeaderConcurrentAbortion(t *testing.T, threads int) {
	// Create a simple chain to verify
	var (
		testdb    = database.NewMemoryDBManager()
		gspec     = &Genesis{Config: params.TestChainConfig}
		genesis   = gspec.MustCommit(testdb)
		blocks, _ = GenerateChain(params.TestChainConfig, genesis, gxhash.NewFaker(), testdb, 1024, nil)
	)
	headers := make([]*types.Header, len(blocks))
	seals := make([]bool, len(blocks))

	for i, block := range blocks {
		headers[i] = block.Header()
		seals[i] = true
	}
	// Set the number of threads to verify on
	old := runtime.GOMAXPROCS(threads)
	defer runtime.GOMAXPROCS(old)

	// Start the verifications and immediately abort
	chain, _ := NewBlockChain(testdb, nil, params.TestChainConfig, gxhash.NewFakeDelayer(time.Millisecond), vm.Config{})
	defer chain.Stop()

	abort, results := chain.engine.VerifyHeaders(chain, headers, seals)
	close(abort)

	// Deplete the results channel
	verified := 0
	for depleted := false; !depleted; {
		select {
		case result := <-results:
			if result != nil {
				t.Errorf("header %d: validation failed: %v", verified, result)
			}
			verified++
		case <-time.After(50 * time.Millisecond):
			depleted = true
		}
	}
	// Check that abortion was honored by not processing too many POWs
	if verified > 2*threads {
		t.Errorf("verification count too large: have %d, want below %d", verified, 2*threads)
	}
}
