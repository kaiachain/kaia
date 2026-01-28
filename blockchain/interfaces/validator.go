package interfaces

import (
	"github.com/kaiachain/kaia/blockchain/state"
	"github.com/kaiachain/kaia/blockchain/types"
)

// Validator is an interface which defines the standard for block validation. It
// is only responsible for validating block contents, as the header validation is
// done by the specific consensus engines.
type Validator interface {
	// ValidateHeader validates the given header.
	ValidateHeader(header *types.Header) error

	// ValidateHeaders validates the given headers.
	ValidateHeaders(headers []*types.Header) (chan<- struct{}, <-chan error)

	// ValidateBody validates the given block's content.
	ValidateBody(block *types.Block) error

	// ValidateState validates the given statedb and optionally the receipts and
	// gas used.
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas uint64) error
}
