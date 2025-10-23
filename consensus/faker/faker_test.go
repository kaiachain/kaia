// Copyright 2025 The Kaia Authors
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

package faker

import (
	"math/big"
	"testing"
	"time"

	"github.com/kaiachain/kaia/blockchain"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/blockchain/vm"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/params"
	"github.com/kaiachain/kaia/storage/database"
	"github.com/stretchr/testify/assert"
)

// TestNewFaker tests the creation of different faker instances
func TestNewFaker(t *testing.T) {
	// Test NewFaker
	f := NewFaker()
	assert.NotNil(t, f)
	assert.Equal(t, uint64(0), f.failBlock)
	assert.Equal(t, time.Duration(0), f.failDelay)
	assert.False(t, f.fullFake)

	// Test NewFakeFailer
	f2 := NewFakeFailer(10)
	assert.NotNil(t, f2)
	assert.Equal(t, uint64(10), f2.failBlock)

	// Test NewFakeDelayer
	delay := 5 * time.Second
	f3 := NewFakeDelayer(delay)
	assert.NotNil(t, f3)
	assert.Equal(t, delay, f3.failDelay)

	// Test NewFullFaker
	f4 := NewFullFaker()
	assert.NotNil(t, f4)
	assert.True(t, f4.fullFake)
}

// TestAuthor tests the Author method
func TestAuthor(t *testing.T) {
	f := NewFaker()
	header := &types.Header{Number: big.NewInt(1)}

	author, err := f.Author(header)
	assert.NoError(t, err)
	assert.Equal(t, params.AuthorAddressForTesting, author)
}

// TestVerifyHeader tests header verification
func TestVerifyHeader(t *testing.T) {
	f := NewFaker()

	// Test with fullFake mode - should accept everything
	f2 := NewFullFaker()
	err := f2.VerifyHeader(nil, &types.Header{Number: big.NewInt(1)}, true)
	assert.NoError(t, err)

	// Test with failBlock
	f3 := NewFakeFailer(5)
	header := &types.Header{Number: big.NewInt(5)}
	err = f3.VerifyHeader(nil, header, true)
	assert.Equal(t, consensus.ErrUnknownAncestor, err)

	// Test normal case - should pass

	// Create a test blockchain for proper chain context
	db := database.NewMemoryDBManager()
	genesis := blockchain.DefaultKairosGenesisBlock()
	genesis.MustCommit(db)
	chain, _ := blockchain.NewBlockChain(db, nil, params.TestChainConfig, f, vm.Config{})
	defer chain.Stop()

	header = &types.Header{
		Number:     big.NewInt(1),
		ParentHash: chain.CurrentBlock().Hash(),
	}
	err = f.VerifyHeader(chain, header, true)
	assert.NoError(t, err)
}

// TestPrepare tests the Prepare method
func TestPrepare(t *testing.T) {
	f := NewFaker()

	// Create a test blockchain for proper chain context
	db := database.NewMemoryDBManager()
	genesis := blockchain.DefaultKairosGenesisBlock()
	genesis.MustCommit(db)
	chain, _ := blockchain.NewBlockChain(db, nil, params.TestChainConfig, f, vm.Config{})
	defer chain.Stop()

	// Get the current block (genesis block 0) to set proper parent hash
	currentBlock := chain.CurrentBlock()
	header := &types.Header{
		Number:     big.NewInt(1), // Create block 1, not 5
		ParentHash: currentBlock.Hash(),
		Time:       big.NewInt(time.Now().Unix()),
	}

	err := f.Prepare(chain, header)
	assert.NoError(t, err)
	// CalcBlockScore returns 1 for faker consensus (simplified difficulty)
	assert.Equal(t, big.NewInt(1), header.BlockScore)
}

// TestSeal tests the Seal method
func TestSeal(t *testing.T) {
	// Test normal seal
	f := NewFaker()
	block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
	stop := make(chan struct{})

	sealed, err := f.Seal(nil, block, stop)
	assert.NoError(t, err)
	assert.Equal(t, block, sealed)

	// Test seal with failure
	f2 := NewFakeFailer(5)
	block2 := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(5)})
	sealed, err = f2.Seal(nil, block2, stop)
	assert.Error(t, err)
	assert.Nil(t, sealed)

	// Test seal with delay
	f3 := NewFakeDelayer(100 * time.Millisecond)
	start := time.Now()
	sealed, err = f3.Seal(nil, block, stop)
	elapsed := time.Since(start)
	assert.NoError(t, err)
	assert.NotNil(t, sealed)
	assert.True(t, elapsed >= 100*time.Millisecond)
}

// TestVerifySeal tests seal verification
func TestVerifySeal(t *testing.T) {
	f := NewFaker()
	header := &types.Header{Number: big.NewInt(3)}

	// Normal case - should pass
	err := f.VerifySeal(nil, header)
	assert.NoError(t, err)

	// Test with failBlock
	f2 := NewFakeFailer(5)
	header2 := &types.Header{Number: big.NewInt(5)}
	err = f2.VerifySeal(nil, header2)
	assert.Error(t, err)
}

// TestVerifyHeaders tests batch header verification
func TestVerifyHeaders(t *testing.T) {
	f := NewFaker()

	// Test with empty headers
	headers := []*types.Header{}
	seals := []bool{}
	abort, results := f.VerifyHeaders(nil, headers, seals)
	assert.NotNil(t, abort)
	assert.NotNil(t, results)

	// Test with fullFake mode
	f2 := NewFullFaker()
	headers = []*types.Header{
		{Number: big.NewInt(1)},
		{Number: big.NewInt(2)},
	}
	seals = []bool{true, true}
	abort, results = f2.VerifyHeaders(nil, headers, seals)
	assert.NotNil(t, abort)
	assert.NotNil(t, results)

	// Read results
	for i := 0; i < len(headers); i++ {
		err := <-results
		assert.NoError(t, err)
	}
}

// TestNewShared tests the NewShared constructor
func TestNewShared(t *testing.T) {
	f := NewShared()
	assert.NotNil(t, f)
	assert.Equal(t, uint64(0), f.failBlock)
	assert.Equal(t, time.Duration(0), f.failDelay)
	assert.False(t, f.fullFake)
}

// TestHeaderValidation tests various header validation scenarios
func TestHeaderValidation(t *testing.T) {
	db := database.NewMemoryDBManager()
	genesis := blockchain.DefaultKairosGenesisBlock()
	genesis.MustCommit(db)
	chain, _ := blockchain.NewBlockChain(db, nil, params.TestChainConfig, NewFaker(), vm.Config{})
	defer chain.Stop()

	tests := []struct {
		name      string
		faker     *Faker
		header    *types.Header
		chain     *blockchain.BlockChain
		expectErr bool
	}{
		{
			name:  "normal header passes",
			faker: NewFaker(),
			header: &types.Header{
				Number:     big.NewInt(1),
				ParentHash: chain.CurrentBlock().Hash(),
				Time:       big.NewInt(time.Now().Unix()),
			},
			chain:     chain,
			expectErr: false,
		},
		{
			name:      "fullFake accepts anything",
			faker:     NewFullFaker(),
			header:    &types.Header{Number: big.NewInt(999)},
			chain:     nil,
			expectErr: false,
		},
		{
			name:      "failBlock triggers error",
			faker:     NewFakeFailer(2),
			header:    &types.Header{Number: big.NewInt(2)},
			chain:     chain,
			expectErr: true,
		},
		{
			name:  "before failBlock passes",
			faker: NewFakeFailer(2),
			header: &types.Header{
				Number:     big.NewInt(1),
				ParentHash: chain.CurrentBlock().Hash(),
				Time:       big.NewInt(time.Now().Unix()),
			},
			chain:     chain,
			expectErr: false,
		},
		{
			name:  "after failBlock passes",
			faker: NewFakeFailer(2),
			header: &types.Header{
				Number:     big.NewInt(3),
				ParentHash: chain.CurrentBlock().Hash(),
				Time:       big.NewInt(time.Now().Unix()),
			},
			chain:     chain,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.faker.VerifyHeader(tt.chain, tt.header, true)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
