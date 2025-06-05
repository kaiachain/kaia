// Modifications Copyright 2024 The Kaia Authors
// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from cmd/utils/cmd.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package utils

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kaiachain/kaia/v2/blockchain"
	"github.com/kaiachain/kaia/v2/blockchain/types"
	"github.com/kaiachain/kaia/v2/log"
	"github.com/kaiachain/kaia/v2/node"
	"github.com/kaiachain/kaia/v2/rlp"
)

const (
	importBatchSize = 2500
)

var logger = log.NewModuleLogger(log.CMDUtils)

func StartNode(stack *node.Node) {
	if err := stack.Start(); err != nil {
		log.Fatalf("Error starting protocol stack: %v", err)
	}
	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(sigc)
		<-sigc
		logger.Info("Got interrupt, shutting down...")
		go stack.Stop()
		for i := 10; i > 0; i-- {
			<-sigc
			if i > 1 {
				logger.Info("Already shutting down, interrupt more to panic.", "times", i-1)
			}
		}
	}()
}

func ImportChain(chain *blockchain.BlockChain, fn string) error {
	// Watch for Ctrl-C while the import is running.
	// If a signal is received, the import will stop at the next batch.
	interrupt := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)
	defer close(interrupt)
	go func() {
		if _, ok := <-interrupt; ok {
			logger.Info("Interrupted during import, stopping at next batch")
		}
		close(stop)
	}()
	checkInterrupt := func() bool {
		select {
		case <-stop:
			return true
		default:
			return false
		}
	}

	logger.Info("Importing blockchain", "file", fn)

	// Open the file handle and potentially unwrap the gzip stream
	fh, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer fh.Close()

	var reader io.Reader = fh
	if strings.HasSuffix(fn, ".gz") {
		if reader, err = gzip.NewReader(reader); err != nil {
			return err
		}
	}
	stream := rlp.NewStream(reader, 0)

	// Run actual the import.
	blocks := make(types.Blocks, importBatchSize)
	n := 0
	for batch := 0; ; batch++ {
		// Load a batch of RLP blocks.
		if checkInterrupt() {
			return fmt.Errorf("interrupted")
		}
		i := 0
		for ; i < importBatchSize; i++ {
			var b types.Block
			if err := stream.Decode(&b); err == io.EOF {
				break
			} else if err != nil {
				return fmt.Errorf("at block %d: %v", n, err)
			}
			// don't import first block
			if b.NumberU64() == 0 {
				i--
				continue
			}
			blocks[i] = &b
			n++
		}
		if i == 0 {
			break
		}
		// Import the batch.
		if checkInterrupt() {
			return fmt.Errorf("interrupted")
		}
		missing := missingBlocks(chain, blocks[:i])
		if len(missing) == 0 {
			logger.Info("Skipping batch as all blocks present", "batch", batch, "first", blocks[0].Hash(), "last", blocks[i-1].Hash())
			continue
		}
		if _, err := chain.InsertChain(missing); err != nil {
			return fmt.Errorf("invalid block %d: %v", n, err)
		}
	}
	return nil
}

func missingBlocks(chain *blockchain.BlockChain, blocks []*types.Block) []*types.Block {
	head := chain.CurrentBlock()
	for i, block := range blocks {
		// If we're behind the chain head, only check block, state is available at head
		if head.NumberU64() > block.NumberU64() {
			if !chain.HasBlock(block.Hash(), block.NumberU64()) {
				return blocks[i:]
			}
			continue
		}
		// If we're above the chain head, state availability is a must
		if !chain.HasBlockAndState(block.Hash(), block.NumberU64()) {
			return blocks[i:]
		}
	}
	return nil
}

// ExportChain exports a blockchain into the specified file, truncating any data
// already present in the file.
func ExportChain(blockchain *blockchain.BlockChain, fn string) error {
	logger.Info("Exporting blockchain", "file", fn)

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}
	// Iterate over the blocks and export them
	if err := blockchain.Export(writer); err != nil {
		return err
	}
	logger.Info("Exported blockchain", "file", fn)

	return nil
}

// ExportAppendChain exports a blockchain into the specified file, appending to
// the file if data already exists in it.
func ExportAppendChain(blockchain *blockchain.BlockChain, fn string, first uint64, last uint64) error {
	logger.Info("Exporting blockchain", "file", fn)

	// Open the file handle and potentially wrap with a gzip stream
	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fh.Close()

	var writer io.Writer = fh
	if strings.HasSuffix(fn, ".gz") {
		writer = gzip.NewWriter(writer)
		defer writer.(*gzip.Writer).Close()
	}
	// Iterate over the blocks and export them
	if err := blockchain.ExportN(writer, first, last); err != nil {
		return err
	}
	logger.Info("Exported blockchain to", "file", fn)
	return nil
}

// TODO-Kaia Commented out due to mismatched interface.
//// ImportPreimages imports a batch of exported hash preimages into the database.
//func ImportPreimages(db *database.LevelDB, fn string) error {
//	logger.Info("Importing preimages", "file", fn)
//
//	// Open the file handle and potentially unwrap the gzip stream
//	fh, err := os.Open(fn)
//	if err != nil {
//		return err
//	}
//	defer fh.Close()
//
//	var reader io.Reader = fh
//	if strings.HasSuffix(fn, ".gz") {
//		if reader, err = gzip.NewReader(reader); err != nil {
//			return err
//		}
//	}
//	stream := rlp.NewStream(reader, 0)
//
//	// Import the preimages in batches to prevent disk trashing
//	preimages := make(map[common.Hash][]byte)
//
//	for {
//		// Read the next entry and ensure it's not junk
//		var blob []byte
//
//		if err := stream.Decode(&blob); err != nil {
//			if err == io.EOF {
//				break
//			}
//			return err
//		}
//		// Accumulate the preimages and flush when enough ws gathered
//		preimages[crypto.Keccak256Hash(blob)] = common.CopyBytes(blob)
//		if len(preimages) > 1024 {
//			rawdb.WritePreimages(db, 0, preimages)
//			preimages = make(map[common.Hash][]byte)
//		}
//	}
//	// Flush the last batch preimage data
//	if len(preimages) > 0 {
//		rawdb.WritePreimages(db, 0, preimages)
//	}
//	return nil
//}
//
//// ExportPreimages exports all known hash preimages into the specified file,
//// truncating any data already present in the file.
//func ExportPreimages(db *database.LevelDB, fn string) error {
//	logger.Info("Exporting preimages", "file", fn)
//
//	// Open the file handle and potentially wrap with a gzip stream
//	fh, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
//	if err != nil {
//		return err
//	}
//	defer fh.Close()
//
//	var writer io.Writer = fh
//	if strings.HasSuffix(fn, ".gz") {
//		writer = gzip.NewWriter(writer)
//		defer writer.(*gzip.Writer).Close()
//	}
//	// Iterate over the preimages and export them
//	it := db.NewIteratorWithPrefix([]byte("secure-key-"))
//	for it.Next() {
//		if err := rlp.Encode(writer, it.Value()); err != nil {
//			return err
//		}
//	}
//	logger.Info("Exported preimages", "file", fn)
//	return nil
//}
