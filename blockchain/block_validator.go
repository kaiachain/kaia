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
// This file is derived from core/block_validator.go (2018/06/04).
// Modified and improved for the klaytn development.
// Modified and improved for the Kaia development.

package blockchain

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
	"github.com/kaiachain/kaia/consensus"
	"github.com/kaiachain/kaia/kaiax"
	"github.com/kaiachain/kaia/kaiax/gov"
	"github.com/kaiachain/kaia/params"
)

var (
	// ErrInvalidTimestamp is returned if the timestamp of a block is lower than
	// the previous block's timestamp plus the minimum block interval.
	ErrInvalidTimestamp = errors.New("invalid timestamp")

	// ErrInvalidBlockScore is returned if the BlockScore of a block is not 1.
	ErrInvalidBlockScore = errors.New("invalid blockscore")

	// ErrInvalidBaseFee is returned if a block before Magma has a non-nil base fee.
	ErrInvalidBaseFee = errors.New("invalid baseFee before fork")

	// ErrUnexpectedExcessBlobGasBeforeOsaka is returned if excessBlobGas exists before Osaka.
	ErrUnexpectedExcessBlobGasBeforeOsaka = errors.New("unexpected excessBlobGas before osaka")

	// ErrUnexpectedBlobGasUsedBeforeOsaka is returned if blobGasUsed exists before Osaka.
	ErrUnexpectedBlobGasUsedBeforeOsaka = errors.New("unexpected blobGasUsed before osaka")
)

// BlockValidator is responsible for validating block headers and
// processed state.
//
// BlockValidator implements Validator.
type BlockValidator struct {
	config        *params.ChainConfig   // Chain configuration options
	bc            *BlockChain           // Canonical block chain
	hc            consensus.ChainReader // Header chain
	engine        consensus.Engine      // Consensus engine
	headerModules []kaiax.HeaderModule  // Header modules
	mGov          gov.GovModule         // Governance module
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(config *params.ChainConfig, blockchain *BlockChain) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		bc:     blockchain,
		hc:     blockchain,
		engine: blockchain.Engine(),
	}
	return validator
}

func NewBlockValidatorWithHeaderChain(config *params.ChainConfig, hc consensus.ChainReader) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		hc:     hc,
		engine: hc.Engine(),
	}
	return validator
}

func (v *BlockValidator) RegisterHeaderModules(modules ...kaiax.HeaderModule) {
	v.headerModules = append(v.headerModules, modules...)
}

func (v *BlockValidator) SetupKaiaxModules(mGov gov.GovModule) {
	v.mGov = mGov
}

func (v *BlockValidator) kip71Config(blockNum uint64) (*params.KIP71Config, error) {
	if v.mGov != nil {
		pset := v.mGov.GetParamSet(blockNum)
		return pset.ToKip71Config(), nil
	}
	if v.config == nil || v.config.Governance == nil || v.config.Governance.KIP71 == nil {
		return nil, errors.New("missing KIP-71 config for magma header verification")
	}
	return v.config.Governance.KIP71, nil
}

func (v *BlockValidator) ValidateHeader(header *types.Header) error {
	if header.Number == nil {
		return consensus.ErrUnknownBlock
	}

	return v.validateHeader(header, nil)
}

func (v *BlockValidator) Preprocess(headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		errored := false
		for _, header := range headers {
			var err error
			if errored {
				err = consensus.ErrUnknownAncestor
			} else {
				err = v.engine.VerifySeals(header, true)
			}

			if err != nil {
				errored = true
			}

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

func (v *BlockValidator) validateHeader(header *types.Header, parent *types.Header) error {
	if header.Number == nil {
		return consensus.ErrUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(time.Now().Add(time.Duration(params.DefaultBlockGenerationInterval)*time.Second).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}

	// Ensure that the block's blockscore is meaningful (may not be correct at this point)
	if header.BlockScore == nil || header.BlockScore.Cmp(params.DefaultBlockScore) != 0 {
		return ErrInvalidBlockScore
	}

	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}

	// Ensure that the block's timestamp isn't too close to it's parent
	if parent == nil {
		parent = v.hc.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time.Uint64()+uint64(params.BlockGenerationInterval) > header.Time.Uint64() {
		return ErrInvalidTimestamp
	}

	// Verify the existence / non-existence of osaka-specific header fields.
	osaka := v.config.IsOsakaForkEnabled(header.Number)
	if !osaka {
		switch {
		case header.ExcessBlobGas != nil:
			return ErrUnexpectedExcessBlobGasBeforeOsaka
		case header.BlobGasUsed != nil:
			return ErrUnexpectedBlobGasUsedBeforeOsaka
		}
	} else {
		if header.Number.Uint64() != parent.Number.Uint64()+1 {
			panic("bad header pair")
		}
		bcfg := v.config.LatestBlobConfig(header.Number)
		if bcfg == nil {
			panic("called before EIP-4844 is active")
		}
		if err := bcfg.VerifyEIP4844Header(parent.ExcessBlobGas, parent.BlobGasUsed, header.ExcessBlobGas, header.BlobGasUsed); err != nil {
			return err
		}
	}

	// Verify Magma basefee rule from governance paramset.
	if v.config.IsMagmaForkEnabled(header.Number) {
		kip71Config, err := v.kip71Config(header.Number.Uint64())
		if err != nil {
			return err
		}
		if err := kip71Config.VerifyMagmaHeader(header.BaseFee, parent.Number, parent.BaseFee, parent.GasUsed); err != nil {
			return err
		}
	} else if header.BaseFee != nil {
		return ErrInvalidBaseFee
	}

	// Verify the header's seal and consensus-specific fields.
	if err := v.engine.VerifySeals(header, false); err != nil {
		return err
	}

	// Verify the header with registered header modules.
	for _, module := range v.headerModules {
		if err := module.VerifyHeader(header, parent); err != nil {
			return err
		}
	}
	return nil
}

// ValidateBody verifies the block
// header's transaction. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(block *types.Block) error {
	if v.bc == nil {
		return errors.New("block validator with header chain is not supported")
	}
	// check EIP 7934 RLP-encoded block size cap
	if v.config.IsOsakaForkEnabled(block.Number()) && block.Size() > params.MaxBlockSize {
		return ErrBlockOversized
	}
	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(block.Hash(), block.NumberU64()) {
		return ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(block.ParentHash(), block.NumberU64()-1) {
		if !v.bc.HasBlock(block.ParentHash(), block.NumberU64()-1) {
			logger.Error("unknown ancestor (ValidateBody)", "num", block.NumberU64(),
				"hash", block.Hash(), "parentHash", block.ParentHash())
			return consensus.ErrUnknownAncestor
		}
		return consensus.ErrPrunedAncestor
	}
	// Header validity is known at this point, check the transactions
	header := block.Header()
	if hash := types.DeriveTransactionsRoot(block.Transactions(), block.Number()); hash != header.TxHash {
		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}

	baseFee := block.Header().BaseFee
	// Blob transactions may be present after the Osaka fork.
	var blobs int
	for _, tx := range block.Transactions() {
		// NOTE: Kaia validates tx gasPrice
		if baseFee != nil {
			if baseFee.Cmp(tx.GasPrice()) > 0 {
				return fmt.Errorf("invalid GasPrice: txHash %x, GasPrice %d, BaseFee %d", tx.Hash(), tx.GasPrice(), baseFee)
			}
		}
		// Count the number of blobs to validate against the header's blobGasUsed
		blobs += len(tx.BlobHashes())

		// The individual checks for blob validity (version-check + not empty)
		// happens in state transition.
	}

	// Check blob gas usage.
	if header.BlobGasUsed != nil {
		if want := *header.BlobGasUsed / params.BlobTxBlobGasPerBlob; uint64(blobs) != want { // div because the header is surely good vs the body might be bloated
			return fmt.Errorf("blob gas used mismatch (header %v, calculated %v)", *header.BlobGasUsed, blobs*params.BlobTxBlobGasPerBlob)
		}
	} else {
		if blobs > 0 {
			return errors.New("data blobs present in block body")
		}
	}

	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as amount of used gas, the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(block, parent *types.Block, statedb *state.StateDB, receipts types.Receipts, usedGas uint64) error {
	header := block.Header()
	if block.GasUsed() != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", block.GasUsed(), usedGas)
	}
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveReceiptsRoot(receipts, block.Number())
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(true); header.Root != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	return nil
}
